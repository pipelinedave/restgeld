package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	sessionCookieName = "restgeld_session"
	sessionDuration   = 30 * 24 * time.Hour
	magicLinkDuration = 15 * time.Minute
)

type User struct {
	ID                   string    `json:"id"`
	Email                string    `json:"email"`
	CreatedAt            time.Time `json:"createdAt"`
	LastLoginAt          time.Time `json:"lastLoginAt"`
	DefaultMonthlyBudget float64   `json:"defaultMonthlyBudget"`
	DefaultPeriodDays    int       `json:"defaultPeriodDays"`
	Theme                string    `json:"theme"`
	IsActive             bool      `json:"isActive"`
}

type MagicLinkRequest struct {
	Email string `json:"email"`
}

type VerifyMagicLinkRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	User      *User  `json:"user"`
	Token     string `json:"token,omitempty"`
	DebugLink string `json:"debugLink,omitempty"`
	IsNewUser bool   `json:"isNewUser"`
}

type UpdateUserSettingsRequest struct {
	DefaultMonthlyBudget float64 `json:"defaultMonthlyBudget,omitempty"`
	DefaultPeriodDays    int     `json:"defaultPeriodDays,omitempty"`
	Theme                string  `json:"theme,omitempty"`
}

type authService struct {
	db  *sql.DB
	now func() time.Time
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func jsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *authService) getUserIDFromRequest(r *http.Request) string {
	var rawToken string
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		rawToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if rawToken == "" {
		if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
			rawToken = cookie.Value
		}
	}

	if rawToken == "" {
		return ""
	}

	tokenHash := hashToken(rawToken)
	var userID string
	err := s.db.QueryRow(
		"SELECT user_id FROM auth_sessions WHERE token_hash = $1 AND expires_at > $2",
		tokenHash, s.now(),
	).Scan(&userID)

	if err != nil {
		return ""
	}
	return userID
}

func (s *authService) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error", "db": "disconnected"})
		return
	}
	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "auth-service", "db": "connected"})
}

func (s *authService) handleMagicLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	var req MagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültiges json")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if _, err := mail.ParseAddress(email); err != nil || email == "" {
		writeError(w, http.StatusBadRequest, "ungültige e-mail-adresse")
		return
	}

	token, err := generateSecureToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "serverfehler beim erstellen des login-links")
		return
	}

	tokenHash := hashToken(token)
	expiresAt := s.now().Add(magicLinkDuration)

	if s.db == nil {
		writeError(w, http.StatusInternalServerError, "keine datenbank-verbindung")
		return
	}

	_, err = s.db.Exec(
		"INSERT INTO magic_links (email, token_hash, expires_at) VALUES ($1, $2, $3)",
		email, tokenHash, expiresAt,
	)
	if err != nil {
		log.Printf("fehler beim speichern des magic links: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim erstellen des login-links")
		return
	}

	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		origin := r.Header.Get("Origin")
		if origin != "" {
			appBaseURL = origin
		} else {
			appBaseURL = "http://localhost:5173"
		}
	}
	appBaseURL = strings.TrimRight(appBaseURL, "/")
	magicLinkURL := fmt.Sprintf("%s/?auth_token=%s", appBaseURL, token)

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost != "" {
		go sendMagicLinkEmail(email, magicLinkURL)
	} else {
		log.Printf("[DEV AUTH-SVC] Magic Link für %s: %s", email, magicLinkURL)
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"message":   "Login-Link wurde gesendet. Bitte prüfe deinen Posteingang.",
		"debugLink": magicLinkURL,
	})
}

func (s *authService) handleVerifyMagicLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	var rawToken string
	if r.Method == http.MethodGet {
		rawToken = r.URL.Query().Get("token")
	} else {
		var req VerifyMagicLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			rawToken = req.Token
		}
	}

	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		writeError(w, http.StatusBadRequest, "token erforderlich")
		return
	}

	tokenHash := hashToken(rawToken)
	var email string
	now := s.now()
	err := s.db.QueryRow(
		`UPDATE magic_links 
		 SET used_at = $1 
		 WHERE token_hash = $2 AND expires_at > $1 AND used_at IS NULL 
		 RETURNING email`,
		now, tokenHash,
	).Scan(&email)

	if err != nil {
		writeError(w, http.StatusUnauthorized, "ungültiger oder abgelaufener login-link")
		return
	}

	var u User
	isNew := false
	err = s.db.QueryRow(
		`SELECT id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)

	if err == sql.ErrNoRows {
		isNew = true
		err = s.db.QueryRow(
			`INSERT INTO users (email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active) 
			 VALUES ($1, $2, $2, 450.00, 30, 'emerald', TRUE) 
			 RETURNING id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active`,
			email, now,
		).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)
	} else if err == nil {
		_, _ = s.db.Exec("UPDATE users SET last_login_at = $1 WHERE id = $2", now, u.ID)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim laden des benutzerprofils")
		return
	}

	sessionToken, err := generateSecureToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim erstellen der sitzung")
		return
	}

	sessionHash := hashToken(sessionToken)
	sessionExpires := now.Add(sessionDuration)

	_, err = s.db.Exec(
		"INSERT INTO auth_sessions (user_id, token_hash, user_agent, ip_address, expires_at) VALUES ($1, $2, $3, $4, $5)",
		u.ID, sessionHash, r.UserAgent(), r.RemoteAddr, sessionExpires,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim speichern der sitzung")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  sessionExpires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	jsonHeader(w)
	writeJSON(w, http.StatusOK, AuthResponse{
		User:      &u,
		Token:     sessionToken,
		IsNewUser: isNew,
	})
}

func (s *authService) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	switch r.Method {
	case http.MethodGet:
		var u User
		err := s.db.QueryRow(
			`SELECT id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active 
			 FROM users WHERE id = $1`,
			userID,
		).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)
		if err != nil {
			writeError(w, http.StatusNotFound, "benutzer nicht gefunden")
			return
		}
		jsonHeader(w)
		writeJSON(w, http.StatusOK, u)

	case http.MethodDelete:
		_, err := s.db.Exec("DELETE FROM users WHERE id = $1", userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "fehler beim löschen des accounts")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		jsonHeader(w)
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "account gelöscht"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
	}
}

func (s *authService) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	var req UpdateUserSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültiges json")
		return
	}

	_, err := s.db.Exec(
		`UPDATE users SET 
		 default_monthly_budget = CASE WHEN $2 > 0 THEN $2 ELSE default_monthly_budget END,
		 default_period_days = CASE WHEN $3 > 0 THEN $3 ELSE default_period_days END,
		 theme = CASE WHEN $4 <> '' THEN $4 ELSE theme END
		 WHERE id = $1`,
		userID, req.DefaultMonthlyBudget, req.DefaultPeriodDays, req.Theme,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim speichern der einstellungen")
		return
	}

	var u User
	_ = s.db.QueryRow(
		`SELECT id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active 
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)

	jsonHeader(w)
	writeJSON(w, http.StatusOK, u)
}

func (s *authService) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	var rawToken string
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		rawToken = strings.TrimPrefix(authHeader, "Bearer ")
	}
	if rawToken == "" {
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			rawToken = cookie.Value
		}
	}

	if rawToken != "" {
		_, _ = s.db.Exec("DELETE FROM auth_sessions WHERE token_hash = $1", hashToken(rawToken))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "ausgeloggt"})
}

func (s *authService) handleMigrateGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	var req struct {
		Expenses []map[string]any `json:"expenses"`
		Periods  []map[string]any `json:"periods"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültiges json")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "ok",
		"migratedCount": len(req.Expenses),
	})
}

func sendMagicLinkEmail(toEmail, magicLinkURL string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "1025"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "Restgeld App <noreply@restgeld.stillon.top>"
	}

	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	subject := "Subject: Dein Login-Link für Restgeld\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: system-ui, sans-serif; background: #0a0a0c; color: #f4f4f6; padding: 24px;">
  <div style="max-width: 480px; margin: 0 auto; background: #121216; border: 1px solid #222; border-radius: 16px; padding: 24px; text-align: center;">
    <h1 style="color: #22c55e; margin: 0 0 12px;">restgeld.</h1>
    <p style="color: #8e8e9c; font-size: 15px; margin-bottom: 24px;">Klicke auf den Button unten, um dich einzuloggen:</p>
    <a href="%s" style="display: inline-block; background: #22c55e; color: #000; font-weight: 700; text-decoration: none; padding: 12px 28px; border-radius: 9999px; font-size: 15px;">In Restgeld einloggen &rarr;</a>
  </div>
</body>
</html>`, magicLinkURL)

	msg := []byte("From: " + from + "\r\n" + "To: " + toEmail + "\r\n" + subject + mime + body)
	addr := fmt.Sprintf("%s:%s", host, port)
	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg); err != nil {
		log.Printf("fehler beim mail-versand: %v", err)
	}
}

func (s *authService) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/health", s.handleHealth)
	mux.HandleFunc("/api/auth/magic-link", s.handleMagicLink)
	mux.HandleFunc("/api/auth/verify", s.handleVerifyMagicLink)
	mux.HandleFunc("/api/auth/me", s.handleMe)
	mux.HandleFunc("/api/auth/settings", s.handleSettings)
	mux.HandleFunc("/api/auth/logout", s.handleLogout)
	mux.HandleFunc("/api/auth/migrate-guest", s.handleMigrateGuest)
	return corsMiddleware(mux)
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := envOrDefault("DB_HOST", "localhost")
		port := envOrDefault("DB_PORT", "5432")
		user := envOrDefault("DB_USER", "restgeld")
		password := envOrDefault("DB_PASSWORD", "restgeld")
		dbname := envOrDefault("DB_NAME", "restgeld_auth")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("fehler beim db-verbindungsaufbau: %v", err)
	}
	defer db.Close()

	if err := runAuthMigrations(db); err != nil {
		log.Printf("warnung: fehler bei auth-migrationen: %v", err)
	}

	port := envOrDefault("PORT", "8081")
	svc := &authService{db: db, now: time.Now}

	log.Printf("auth-service startet auf :%s", port)
	if err := http.ListenAndServe(":"+port, svc.router()); err != nil {
		log.Fatalf("server-fehler: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
