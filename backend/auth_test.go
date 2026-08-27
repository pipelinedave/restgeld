package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRequestMagicLink_Valid(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"email": "david@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/magic-link", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("erwartet status ok, bekommen %v", resp["status"])
	}
	if resp["debugLink"] == nil || resp["debugLink"] == "" {
		t.Errorf("erwartet debugLink")
	}
}

func TestRequestMagicLink_InvalidEmail(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"email": "ungueltige-email"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/magic-link", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("erwartet 400 für ungültige email, bekommen %d", rec.Code)
	}
}

func TestVerifyMagicLink_FullFlow(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	// 1. Magic Link anfordern
	body := `{"email": "test@restgeld.app"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/magic-link", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	var mlResp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&mlResp)
	debugLink := mlResp["debugLink"].(string)

	// Token aus debugLink extrahieren (?auth_token=...)
	tokenParts := strings.Split(debugLink, "auth_token=")
	if len(tokenParts) < 2 {
		t.Fatalf("kein auth_token im debugLink: %s", debugLink)
	}
	token := tokenParts[1]

	// 2. Magic Link verifizieren
	verifyBody := fmt.Sprintf(`{"token": "%s"}`, token)
	vReq := httptest.NewRequest(http.MethodPost, "/api/auth/verify", strings.NewReader(verifyBody))
	vReq.Header.Set("Content-Type", "application/json")
	vRec := httptest.NewRecorder()
	srv.router().ServeHTTP(vRec, vReq)

	if vRec.Code != http.StatusOK {
		t.Fatalf("verify erwartet 200, bekommen %d: %s", vRec.Code, vRec.Body.String())
	}

	var authResp AuthResponse
	if err := json.NewDecoder(vRec.Body).Decode(&authResp); err != nil {
		t.Fatalf("decode auth resp: %v", err)
	}

	if authResp.User == nil || authResp.User.Email != "test@restgeld.app" {
		t.Errorf("unerwarteter user: %+v", authResp.User)
	}
	if authResp.Token == "" {
		t.Errorf("erwartet session token")
	}

	// 3. GET /api/auth/me mit Bearer Token
	meReq := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+authResp.Token)
	meRec := httptest.NewRecorder()
	srv.router().ServeHTTP(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("get me erwartet 200, bekommen %d: %s", meRec.Code, meRec.Body.String())
	}

	var meUser User
	json.NewDecoder(meRec.Body).Decode(&meUser)
	if meUser.Email != "test@restgeld.app" {
		t.Errorf("erwartet email test@restgeld.app, bekommen %s", meUser.Email)
	}

	// 4. Token kann nicht zweimal verwendet werden (Single-Use)
	vReq2 := httptest.NewRequest(http.MethodPost, "/api/auth/verify", strings.NewReader(verifyBody))
	vReq2.Header.Set("Content-Type", "application/json")
	vRec2 := httptest.NewRecorder()
	srv.router().ServeHTTP(vRec2, vReq2)

	if vRec2.Code != http.StatusUnauthorized {
		t.Errorf("erwartet 401 bei wiederverwendetem token, bekommen %d", vRec2.Code)
	}
}

func TestMultiTenant_DataIsolation(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	// 1. User A registrieren & einloggen
	userA, _, _ := store.GetOrCreateUserByEmail("user_a@test.com")
	sessionTokenA := "sess-token-a-12345"
	store.CreateSession(userA.ID, hashToken(sessionTokenA), "agent", "127.0.0.1", time.Now().Add(time.Hour))

	// 2. User B registrieren & einloggen
	userB, _, _ := store.GetOrCreateUserByEmail("user_b@test.com")
	sessionTokenB := "sess-token-b-67890"
	store.CreateSession(userB.ID, hashToken(sessionTokenB), "agent", "127.0.0.1", time.Now().Add(time.Hour))

	// 3. User A bucht eine Ausgabe über 25 €
	expABody := `{"amount": 25.0, "note": "User A Kaffee"}`
	reqA := httptest.NewRequest(http.MethodPost, "/api/expenses", strings.NewReader(expABody))
	reqA.Header.Set("Authorization", "Bearer "+sessionTokenA)
	reqA.Header.Set("Content-Type", "application/json")
	recA := httptest.NewRecorder()
	srv.router().ServeHTTP(recA, reqA)

	if recA.Code != http.StatusCreated {
		t.Fatalf("User A create expense: %d %s", recA.Code, recA.Body.String())
	}

	// 4. User B bucht eine Ausgabe über 50 €
	expBBody := `{"amount": 50.0, "note": "User B Lunch"}`
	reqB := httptest.NewRequest(http.MethodPost, "/api/expenses", strings.NewReader(expBBody))
	reqB.Header.Set("Authorization", "Bearer "+sessionTokenB)
	reqB.Header.Set("Content-Type", "application/json")
	recB := httptest.NewRecorder()
	srv.router().ServeHTTP(recB, reqB)

	if recB.Code != http.StatusCreated {
		t.Fatalf("User B create expense: %d %s", recB.Code, recB.Body.String())
	}

	// 5. User A holt seine Ausgaben
	getAReq := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	getAReq.Header.Set("Authorization", "Bearer "+sessionTokenA)
	getARec := httptest.NewRecorder()
	srv.router().ServeHTTP(getARec, getAReq)

	var pExpA PaginatedExpenses
	json.NewDecoder(getARec.Body).Decode(&pExpA)
	if pExpA.Total != 1 || pExpA.Items[0].Note != "User A Kaffee" || pExpA.Items[0].Amount != 25.0 {
		t.Errorf("User A sieht unerwartete Daten: %+v", pExpA)
	}

	// 6. User B holt seine Ausgaben
	getBReq := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	getBReq.Header.Set("Authorization", "Bearer "+sessionTokenB)
	getBRec := httptest.NewRecorder()
	srv.router().ServeHTTP(getBRec, getBReq)

	var pExpB PaginatedExpenses
	json.NewDecoder(getBRec.Body).Decode(&pExpB)
	if pExpB.Total != 1 || pExpB.Items[0].Note != "User B Lunch" || pExpB.Items[0].Amount != 50.0 {
		t.Errorf("User B sieht unerwartete Daten: %+v", pExpB)
	}

	// 7. Gast (ohne Auth) hat keine dieser Ausgaben
	guestReq := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	guestRec := httptest.NewRecorder()
	srv.router().ServeHTTP(guestRec, guestReq)

	var pExpGuest PaginatedExpenses
	json.NewDecoder(guestRec.Body).Decode(&pExpGuest)
	if pExpGuest.Total != 0 {
		t.Errorf("Gast sollte 0 Ausgaben sehen, hat %d", pExpGuest.Total)
	}
}

func TestUpdateUserSettings(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	user, _, _ := store.GetOrCreateUserByEmail("settings@test.com")
	sessionToken := "sess-settings-123"
	store.CreateSession(user.ID, hashToken(sessionToken), "agent", "127.0.0.1", time.Now().Add(time.Hour))

	body := `{"defaultMonthlyBudget": 600.0, "defaultPeriodDays": 14, "theme": "cyberpunk"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/auth/settings", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("settings update: %d %s", rec.Code, rec.Body.String())
	}

	var updatedUser User
	json.NewDecoder(rec.Body).Decode(&updatedUser)
	if updatedUser.DefaultMonthlyBudget != 600.0 || updatedUser.DefaultPeriodDays != 14 || updatedUser.Theme != "cyberpunk" {
		t.Errorf("unerwartete aktualisierte Settings: %+v", updatedUser)
	}
}

func TestMigrateGuestData(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	user, _, _ := store.GetOrCreateUserByEmail("migrate@test.com")
	sessionToken := "sess-migrate-123"
	store.CreateSession(user.ID, hashToken(sessionToken), "agent", "127.0.0.1", time.Now().Add(time.Hour))

	migrateReq := MigrateGuestRequest{
		Expenses: []Expense{
			{Amount: 12.50, Note: "Gast Kaffee", CreatedAt: time.Now()},
			{Amount: 34.00, Note: "Gast Supermarkt", CreatedAt: time.Now()},
		},
		Periods: []Period{
			{StartDate: time.Now(), MonthDays: 30, BaseBudget: 15.0, MonthlyTotal: 450.0},
		},
	}

	b, _ := json.Marshal(migrateReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/migrate-guest", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("migrate guest: %d %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["migratedCount"].(float64) != 2 {
		t.Errorf("erwartet 2 migrierte Ausgaben, bekommen %v", resp["migratedCount"])
	}

	// Prüfen, ob die Ausgaben nun unter dem User abrufbar sind
	getReq := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	getReq.Header.Set("Authorization", "Bearer "+sessionToken)
	getRec := httptest.NewRecorder()
	srv.router().ServeHTTP(getRec, getReq)

	var pExp PaginatedExpenses
	json.NewDecoder(getRec.Body).Decode(&pExp)
	if pExp.Total != 2 {
		t.Errorf("erwartet 2 Ausgaben beim User nach Migration, bekommen %d", pExp.Total)
	}
}

func TestDeleteAccount_GDPR(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	user, _, _ := store.GetOrCreateUserByEmail("delete_me@test.com")
	sessionToken := "sess-delete-123"
	store.CreateSession(user.ID, hashToken(sessionToken), "agent", "127.0.0.1", time.Now().Add(time.Hour))
	store.AddExpense(user.ID, "2026-08", 10.0, "Test Ausgabe")

	// Account löschen
	delReq := httptest.NewRequest(http.MethodDelete, "/api/auth/me", nil)
	delReq.Header.Set("Authorization", "Bearer "+sessionToken)
	delRec := httptest.NewRecorder()
	srv.router().ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("delete account: %d %s", delRec.Code, delRec.Body.String())
	}

	// User sollte nicht mehr existieren
	_, err := store.GetUserByID(user.ID)
	if err == nil {
		t.Errorf("user sollte gelöscht sein")
	}

	// Session sollte ungültig sein
	_, err = store.ValidateSession(hashToken(sessionToken))
	if err == nil {
		t.Errorf("session sollte gelöscht sein")
	}
}
