package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBillingServiceHealth(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodGet, "/api/billing/health", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 OK für health, bekommen %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("fehler beim dekodieren: %v", err)
	}

	if resp["service"] != "billing-service" {
		t.Fatalf("erwartet service=billing-service, bekommen %s", resp["service"])
	}
}

func TestBillingCreateCheckoutUnauthorized(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/billing/create-checkout-session", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("erwartet 401 Unauthorized ohne Auth, bekommen %d", rec.Code)
	}
}

func TestBillingCustomerPortalUnauthorized(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/billing/customer-portal", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("erwartet 401 Unauthorized ohne Auth, bekommen %d", rec.Code)
	}
}

func TestBillingMethodNotAllowed(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodGet, "/api/billing/create-checkout-session", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("erwartet 405 Method Not Allowed für GET, bekommen %d", rec.Code)
	}
}

func TestBillingWebhookValid(t *testing.T) {
	svc := &billingService{now: time.Now}
	body := `{"type": "checkout.session.completed", "data": {"id": "cs_test_123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/billing/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 OK für Webhook, bekommen %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode fehler: %v", err)
	}

	if resp["received"] != "true" || resp["type"] != "checkout.session.completed" {
		t.Fatalf("unerwartete webhook antwort: %+v", resp)
	}
}

func TestBillingWebhookMalformedJSON(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/billing/webhook", strings.NewReader(`{malformed`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("erwartet 400 Bad Request für ungültiges JSON, bekommen %d", rec.Code)
	}
}

func TestBillingCORSOptions(t *testing.T) {
	svc := &billingService{now: time.Now}
	req := httptest.NewRequest(http.MethodOptions, "/api/billing/create-checkout-session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 für OPTIONS, bekommen %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("falscher CORS origin header: %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}
