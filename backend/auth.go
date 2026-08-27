package main

import (
	"crypto/rand"
	"crypto/sha256"
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
)

const (
	sessionCookieName = "restgeld_session"
	sessionDuration   = 30 * 24 * time.Hour // 30 Tage Session-Gültigkeit
	magicLinkDuration = 15 * time.Minute    // 15 Minuten Magic-Link-Gültigkeit
)

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

func (s *server) getUserIDFromRequest(r *http.Request) string {
	var rawToken string

	// 1. Aus Authorization Bearer Header
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		rawToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 2. Fallback: Aus HttpOnly Cookie
	if rawToken == "" {
		if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
			rawToken = cookie.Value
		}
	}

	if rawToken == "" {
		return ""
	}

	tokenHash := hashToken(rawToken)
	userID, err := s.store.ValidateSession(tokenHash)
	if err != nil {
		return ""
	}

	return userID
}

func (s *server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := s.getUserIDFromRequest(r)
		if userID == "" {
			writeError(w, http.StatusUnauthorized, "authentifizierung erforderlich")
			return
		}
		next(w, r)
	}
}

// POST /api/auth/magic-link
func (s *server) requestMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("fehler beim generieren des tokens: %v", err)
		writeError(w, http.StatusInternalServerError, "serverfehler beim erstellen des login-links")
		return
	}

	tokenHash := hashToken(token)
	expiresAt := s.now().Add(magicLinkDuration)

	if err := s.store.CreateMagicLink(email, tokenHash, expiresAt); err != nil {
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

	// E-Mail versenden oder Dev-Logging
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost != "" {
		go sendMagicLinkEmail(email, magicLinkURL)
	} else {
		log.Printf("[DEV/PREVIEW AUTH] Magic Link für %s: %s", email, magicLinkURL)
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"message":   "Login-Link wurde gesendet. Bitte prüfe deinen Posteingang.",
		"debugLink": magicLinkURL, // Erleichtert Dev & Preview Tests
	})
}

// POST /api/auth/verify
func (s *server) verifyMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	email, err := s.store.ValidateAndConsumeMagicLink(tokenHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "ungültiger oder abgelaufener login-link")
		return
	}

	user, isNew, err := s.store.GetOrCreateUserByEmail(email)
	if err != nil {
		log.Printf("fehler beim laden/erstellen des users %s: %v", email, err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden des benutzerprofils")
		return
	}

	sessionToken, err := generateSecureToken()
	if err != nil {
		log.Printf("fehler beim generieren des session-tokens: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim erstellen der sitzung")
		return
	}

	sessionHash := hashToken(sessionToken)
	sessionExpires := s.now().Add(sessionDuration)

	if err := s.store.CreateSession(user.ID, sessionHash, r.UserAgent(), r.RemoteAddr, sessionExpires); err != nil {
		log.Printf("fehler beim speichern der session: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim speichern der sitzung")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  sessionExpires,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production" || strings.HasPrefix(r.Header.Get("Origin"), "https://"),
		SameSite: http.SameSiteLaxMode,
	})

	jsonHeader(w)
	writeJSON(w, http.StatusOK, AuthResponse{
		User:      user,
		Token:     sessionToken,
		IsNewUser: isNew,
	})
}

// GET /api/auth/me
func (s *server) getMeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "benutzer nicht gefunden")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, user)
}

// PATCH /api/auth/settings
func (s *server) updateUserSettingsHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := s.store.UpdateUserSettings(userID, req.DefaultMonthlyBudget, req.DefaultPeriodDays, req.Theme); err != nil {
		log.Printf("fehler beim aktualisieren der user-einstellungen: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim speichern der einstellungen")
		return
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, user)
}

// POST /api/auth/logout
func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
		tokenHash := hashToken(rawToken)
		_ = s.store.DeleteSession(tokenHash)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "erfolgreich ausgeloggt"})
}

// DELETE /api/auth/me (DSGVO-Löschung)
func (s *server) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	if err := s.store.DeleteUser(userID); err != nil {
		log.Printf("fehler beim löschen des accounts %s: %v", userID, err)
		writeError(w, http.StatusInternalServerError, "fehler beim löschen des accounts")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "account und alle daten gelöscht"})
}

// POST /api/auth/migrate-guest
func (s *server) migrateGuestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	var req MigrateGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültiges json")
		return
	}

	count, err := s.store.MigrateGuestData(userID, req.Expenses, req.Periods)
	if err != nil {
		log.Printf("fehler bei gast-migration fuer user %s: %v", userID, err)
		writeError(w, http.StatusInternalServerError, "fehler beim importieren der gast-daten")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "ok",
		"migratedCount": count,
	})
}

func sendMagicLinkEmail(toEmail, magicLinkURL string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "Restgeld App <noreply@restgeld.stillon.top>"
	}

	auth := smtp.PlainAuth("", user, pass, host)
	subject := "Subject: Dein Login-Link für Restgeld\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0a0a0c; color: #f4f4f6; padding: 24px;">
  <div style="max-width: 480px; margin: 0 auto; background: #121216; border: 1px solid #222; border-radius: 16px; padding: 24px; text-align: center;">
    <h1 style="color: #22c55e; margin: 0 0 12px;">restgeld.</h1>
    <p style="color: #8e8e9c; font-size: 15px; margin-bottom: 24px;">Klicke auf den Button unten, um dich direkt in Restgeld einzuloggen:</p>
    <a href="%s" style="display: inline-block; background: #22c55e; color: #000; font-weight: 700; text-decoration: none; padding: 12px 28px; border-radius: 9999px; font-size: 15px;">In Restgeld einloggen &rarr;</a>
    <p style="color: #5c5c6e; font-size: 12px; margin-top: 24px;">Dieser Link ist 15 Minuten lang gültig. Falls du den Link nicht angefordert hast, kannst du diese E-Mail ignorieren.</p>
  </div>
</body>
</html>`, magicLinkURL)

	msg := []byte("From: " + from + "\r\n" + "To: " + toEmail + "\r\n" + subject + mime + body)
	addr := fmt.Sprintf("%s:%s", host, port)
	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg); err != nil {
		log.Printf("fehler beim versenden der login-email an %s: %v", toEmail, err)
	} else {
		log.Printf("magic link email erfolgreich an %s gesendet", toEmail)
	}
}
