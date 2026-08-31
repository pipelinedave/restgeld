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
	sessionDuration   = 30 * 24 * time.Hour
	magicLinkDuration = 15 * time.Minute
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

func (s *server) handleAuthHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error", "db": "disconnected"})
		return
	}
	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "auth-service", "db": "connected"})
}

func (s *server) handleMagicLink(w http.ResponseWriter, r *http.Request) {
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

	smtpHost := os.Getenv("SMTP_HOST")
	lang := r.Header.Get("Accept-Language")
	if smtpHost != "" {
		go sendMagicLinkEmail(email, magicLinkURL, lang)
	} else {
		log.Printf("[AUTH] Magic Link für %s: %s", email, magicLinkURL)
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"message":   "Login-Link wurde gesendet. Bitte prüfe deinen Posteingang.",
		"debugLink": magicLinkURL,
	})
}

func (s *server) handleVerifyMagicLink(w http.ResponseWriter, r *http.Request) {
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

	u, isNew, err := s.store.GetOrCreateUserByEmail(email)
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
	sessionExpires := s.now().Add(sessionDuration)

	if err := s.store.CreateSession(u.ID, sessionHash, r.UserAgent(), r.RemoteAddr, sessionExpires); err != nil {
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
		User:      u,
		Token:     sessionToken,
		IsNewUser: isNew,
	})
}

func (s *server) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, err := s.store.GetUserByID(userID)
		if err != nil {
			writeError(w, http.StatusNotFound, "benutzer nicht gefunden")
			return
		}
		jsonHeader(w)
		writeJSON(w, http.StatusOK, u)

	case http.MethodDelete:
		if err := s.store.DeleteUser(userID); err != nil {
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

func (s *server) handleSettings(w http.ResponseWriter, r *http.Request) {
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

	if err := s.store.UpdateUserSettings(userID, req.DefaultMonthlyBudget, req.DefaultPeriodDays, req.Theme, req.Language, req.Currency); err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim speichern der einstellungen")
		return
	}

	u, _ := s.store.GetUserByID(userID)
	jsonHeader(w)
	writeJSON(w, http.StatusOK, u)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
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
		_ = s.store.DeleteSession(hashToken(rawToken))
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

func (s *server) handleMigrateGuest(w http.ResponseWriter, r *http.Request) {
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

	migrated, err := s.store.MigrateGuestData(userID, req.Expenses, req.Periods)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler bei der migration")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "ok",
		"migratedCount": migrated,
	})
}

func (s *server) handlePasskeyRegisterOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	challenge, err := generateSecureToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim generieren der challenge")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": challenge,
		"rp": map[string]string{
			"name": "restgeld.",
			"id":   r.Host,
		},
		"user": map[string]string{
			"id":          userID,
			"name":        "user",
			"displayName": "restgeld. User",
		},
		"pubKeyCredParams": []map[string]interface{}{
			{"type": "public-key", "alg": -7},  // ES256
			{"type": "public-key", "alg": -257}, // RS256
		},
		"timeout": 60000,
	})
}

func (s *server) handlePasskeyLoginOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	challenge, err := generateSecureToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim generieren der challenge")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": challenge,
		"timeout":   60000,
		"rpId":      r.Host,
	})
}

func (s *server) handlePasskeyRegisterVerify(w http.ResponseWriter, r *http.Request) {
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
		CredentialID    string `json:"credentialId"`
		PublicKey       string `json:"publicKey"`
		AttestationType string `json:"attestationType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CredentialID == "" {
		writeError(w, http.StatusBadRequest, "ungültige passkey daten")
		return
	}

	if err := s.store.SavePasskey(userID, req.CredentialID, req.PublicKey, req.AttestationType); err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim speichern des passkeys")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Passkey erfolgreich verifiziert und gespeichert",
	})
}

func (s *server) handlePasskeyLoginVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	var req struct {
		CredentialID string `json:"credentialId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CredentialID == "" {
		writeError(w, http.StatusBadRequest, "credentialId erforderlich")
		return
	}

	userID, err := s.store.FindUserIDByPasskey(req.CredentialID)
	if err != nil || userID == "" {
		writeError(w, http.StatusUnauthorized, "passkey nicht registriert oder ungültig")
		return
	}

	u, err := s.store.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "benutzerkonto nicht gefunden")
		return
	}

	sessionToken, err := generateSecureToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim erstellen der sitzung")
		return
	}

	sessionHash := hashToken(sessionToken)
	sessionExpires := s.now().Add(sessionDuration)

	if err := s.store.CreateSession(u.ID, sessionHash, r.UserAgent(), r.RemoteAddr, sessionExpires); err != nil {
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
		User:      u,
		Token:     sessionToken,
		IsNewUser: false,
	})
}

func sendMagicLinkEmail(toEmail, magicLinkURL, lang string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "1025"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "restgeld. <noreply@restgeld.app>"
	}

	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	lang = strings.ToLower(lang)
	subject := "Subject: Dein Login-Link für restgeld.\r\n"
	textMsg := "Klicke auf den Button unten, um dich einzuloggen:"
	btnMsg := "In restgeld. einloggen &rarr;"

	if strings.HasPrefix(lang, "en") {
		subject = "Subject: Your Magic Login Link for restgeld.\r\n"
		textMsg = "Click the button below to sign in:"
		btnMsg = "Sign in to restgeld. &rarr;"
	} else if strings.HasPrefix(lang, "es") {
		subject = "Subject: Tu enlace de acceso para restgeld.\r\n"
		textMsg = "Haz clic en el botón de abajo para iniciar sesión:"
		btnMsg = "Iniciar sesión en restgeld. &rarr;"
	} else if strings.HasPrefix(lang, "fr") {
		subject = "Subject: Votre lien de connexion restgeld.\r\n"
		textMsg = "Cliquez sur le bouton ci-dessous pour vous connecter:"
		btnMsg = "Se connecter à restgeld. &rarr;"
	}

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: system-ui, sans-serif; background: #0a0a0c; color: #f4f4f6; padding: 24px;">
  <div style="max-width: 480px; margin: 0 auto; background: #121216; border: 1px solid #222; border-radius: 16px; padding: 24px; text-align: center;">
    <h1 style="color: #22c55e; margin: 0 0 12px;">restgeld.</h1>
    <p style="color: #8e8e9c; font-size: 15px; margin-bottom: 24px;">%s</p>
    <a href="%s" style="display: inline-block; background: #22c55e; color: #000; font-weight: 700; text-decoration: none; padding: 12px 28px; border-radius: 9999px; font-size: 15px;">%s</a>
  </div>
</body>
</html>`, textMsg, magicLinkURL, btnMsg)

	msg := []byte("From: " + from + "\r\n" + "To: " + toEmail + "\r\n" + subject + mime + body)
	addr := fmt.Sprintf("%s:%s", host, port)
	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg); err != nil {
		log.Printf("fehler beim mail-versand: %v", err)
	}
}
