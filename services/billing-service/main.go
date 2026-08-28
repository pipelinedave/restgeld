package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const sessionCookieName = "restgeld_session"

type billingService struct {
	db  *sql.DB
	now func() time.Time
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
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

func (s *billingService) getUserIDFromRequest(r *http.Request) string {
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

	if rawToken == "" || s.db == nil {
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

func (s *billingService) handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "billing-service"})
}

func (s *billingService) handleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	stripePublishableKey := os.Getenv("STRIPE_PUBLISHABLE_KEY")
	appBaseURL := envOrDefault("APP_BASE_URL", "http://localhost:5173")

	mockCheckoutURL := fmt.Sprintf("%s/?billing_success=true&session_id=mock_session_%s", appBaseURL, userID)

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":         "ok",
		"checkoutUrl":    mockCheckoutURL,
		"publishableKey": stripePublishableKey,
		"isMock":         stripePublishableKey == "",
	})
}

func (s *billingService) handleCustomerPortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	userID := s.getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "nicht eingeloggt")
		return
	}

	appBaseURL := envOrDefault("APP_BASE_URL", "http://localhost:5173")
	mockPortalURL := fmt.Sprintf("%s/?portal_return=true", appBaseURL)

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"portalUrl": mockPortalURL,
	})
}

func (s *billingService) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "methode nicht erlaubt")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültiges json")
		return
	}

	eventType, _ := req["type"].(string)
	log.Printf("[BILLING-SVC] Webhook Event empfangen: %s", eventType)

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"received": "true", "type": eventType})
}

func (s *billingService) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/billing/health", s.handleHealth)
	mux.HandleFunc("/api/billing/create-checkout-session", s.handleCreateCheckoutSession)
	mux.HandleFunc("/api/billing/customer-portal", s.handleCustomerPortal)
	mux.HandleFunc("/api/billing/webhook", s.handleWebhook)
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

	db, _ := sql.Open("postgres", connStr)
	if db != nil {
		defer db.Close()
	}

	port := envOrDefault("PORT", "8082")
	svc := &billingService{db: db, now: time.Now}

	log.Printf("billing-service startet auf :%s", port)
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
