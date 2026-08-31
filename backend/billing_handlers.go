package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (s *server) handleBillingHealth(w http.ResponseWriter, r *http.Request) {
	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "billing-service"})
}

func (s *server) handleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
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

func (s *server) handleCustomerPortal(w http.ResponseWriter, r *http.Request) {
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

func (s *server) handleBillingWebhook(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("[BILLING] Webhook Event empfangen: %s", eventType)

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"received": "true", "type": eventType})
}
