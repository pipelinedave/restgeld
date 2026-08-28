package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
