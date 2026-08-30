package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAuthServiceHealth(t *testing.T) {
	svc := &authService{now: time.Now}

	// Direct call health handler without pinging db
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/magic-link", strings.NewReader(`{"email": "invalid"}`))
	rec2 := httptest.NewRecorder()
	svc.router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("erwartet 400 für ungültige email, bekommen %d", rec2.Code)
	}

	req3 := httptest.NewRequest(http.MethodPost, "/api/auth/magic-link", strings.NewReader(`{"email": "user@example.com"}`))
	rec3 := httptest.NewRecorder()
	svc.router().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusInternalServerError {
		// Nil db will return 500 when inserting magic link
		t.Fatalf("erwartet 500 ohne db, bekommen %d", rec3.Code)
	}
}

func TestPasskeyLoginOptions(t *testing.T) {
	svc := &authService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/passkey/login-options", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 OK für passkey login options, bekommen %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("fehler beim dekodieren des json: %v", err)
	}

	if resp["challenge"] == "" {
		t.Fatalf("challenge sollte nicht leer sein")
	}
}

func TestPasskeyRegisterUnauthorized(t *testing.T) {
	svc := &authService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/passkey/register-options", nil)
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("erwartet 401 Unauthorized ohne Auth-Header, bekommen %d", rec.Code)
	}
}

func TestPasskeyRegisterVerifyInvalid(t *testing.T) {
	svc := &authService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/passkey/register-verify", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("erwartet 401 Unauthorized ohne Login, bekommen %d", rec.Code)
	}
}

func TestPasskeyLoginVerifyInvalid(t *testing.T) {
	svc := &authService{now: time.Now}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/passkey/login-verify", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	svc.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("erwartet 400 Bad Request ohne credentialId, bekommen %d", rec.Code)
	}
}
