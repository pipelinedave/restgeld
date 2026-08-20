package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetBudget(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp BudgetResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}

	if resp.Day != 17 {
		t.Errorf("erwartet Tag 17, bekommen %d", resp.Day)
	}
	if resp.MonthDays != 31 {
		t.Errorf("erwartet 31 Tage, bekommen %d", resp.MonthDays)
	}
	if resp.BaseBudget != 14.52 {
		t.Errorf("erwartet base_budget 14.52, bekommen %.2f", resp.BaseBudget)
	}
	if resp.PeriodID != "2026-08" {
		t.Errorf("erwartet period 2026-08, bekommen %s", resp.PeriodID)
	}
}

func TestGetBudgetWithExpenses(t *testing.T) {
	store := newMemoryStore()
	store.AddExpense("2026-08", 10.00, "Pizza")
	store.AddExpense("2026-08", 5.50, "Kaffee")

	now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	var resp BudgetResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Savings <= 0 {
		t.Errorf("ersparnis sollte positiv sein (15.5€ ausgegeben von 246.84€), bekommen %.2f", resp.Savings)
	}

	if len(resp.Expenses) != 2 {
		t.Errorf("erwartet 2 ausgaben, bekommen %d", len(resp.Expenses))
	}
}

func TestCreateExpense(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	body := `{"amount": 8.50, "note": "Bus"}`
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var exp Expense
	json.NewDecoder(rec.Body).Decode(&exp)

	if exp.Amount != 8.50 {
		t.Errorf("erwartet amount 8.50, bekommen %.2f", exp.Amount)
	}
	if exp.Note != "Bus" {
		t.Errorf("erwartet note 'Bus', bekommen '%s'", exp.Note)
	}
	if exp.PeriodID != "2026-08" {
		t.Errorf("erwartet period 2026-08, bekommen %s", exp.PeriodID)
	}
}

func TestCreateExpenseInvalid(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"amount": 0, "note": "nix"}`
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 für amount=0, bekommen %d", rec.Code)
	}
}

func TestDeleteExpense(t *testing.T) {
	store := newMemoryStore()
	exp, _ := store.AddExpense("2026-08", 5.00, "Test")
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodDelete, "/api/expenses/"+exp.ID, nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	total, _ := store.GetTotalExpenses("2026-08")
	if total != 0 {
		t.Errorf("erwartet 0 nach löschen, bekommen %.2f", total)
	}
}

func TestDeleteExpenseNotFound(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodDelete, "/api/expenses/nonexistent", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("erwartet 404, bekommen %d", rec.Code)
	}
}

func TestNewPeriod(t *testing.T) {
	store := newMemoryStore()
	store.AddExpense("2026-08", 50.00, "Alt")
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodPost, "/api/period", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d", rec.Code)
	}

	total, _ := store.GetTotalExpenses("2026-08")
	if total != 0 {
		t.Errorf("erwartet 0 nach reset, bekommen %.2f", total)
	}
}

func TestUpdateBudget(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"monthlyTotal": 600}`
	req := httptest.NewRequest(http.MethodPatch, "/api/budget", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	p, _ := store.GetOrCreatePeriod()
	if p.MonthlyTotal != 600 {
		t.Errorf("erwartet monthly_total 600, bekommen %.2f", p.MonthlyTotal)
	}
}

func TestCORSHeaders(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodOptions, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("erwartet 200 für OPTIONS, bekommen %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header fehlt")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodPut, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("erwartet 405, bekommen %d", rec.Code)
	}
}

func TestGetBudgetOnFirstDay(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 8, 1, 6, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	var resp BudgetResponse
	decodeJSON(t, rec.Body, &resp)

	if resp.Day != 1 {
		t.Errorf("erwartet Tag 1, bekommen %d", resp.Day)
	}
	if resp.Savings != resp.BaseBudget {
		t.Errorf("erwartet savings = baseBudget (%.2f) an tag 1, bekommen %.2f", resp.BaseBudget, resp.Savings)
	}
	if resp.CurrentBudget <= resp.BaseBudget {
		t.Errorf("erwartet current > base (%.2f) wegen rollover, bekommen %.2f", resp.BaseBudget, resp.CurrentBudget)
	}
}

func TestGetBudgetAfterAllSpent(t *testing.T) {
	store := newMemoryStore()
	for i := 0; i < 18; i++ {
		store.AddExpense("2026-08", 14.52, "tag")
	}
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	var resp BudgetResponse
	decodeJSON(t, rec.Body, &resp)

	if resp.Color != "white" {
		t.Errorf("erwartet white bei genau aufgebraucht, bekommen %s", resp.Color)
	}
	if resp.Savings != 0 {
		t.Errorf("erwartet savings 0 bei genau aufgebraucht, bekommen %.2f", resp.Savings)
	}
}

func TestGetBudgetInDebt(t *testing.T) {
	store := newMemoryStore()
	for i := 0; i < 18; i++ {
		store.AddExpense("2026-08", 14.52, "tag")
	}
	store.AddExpense("2026-08", 50, "extra")
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	var resp BudgetResponse
	decodeJSON(t, rec.Body, &resp)

	if resp.Color != "red" {
		t.Errorf("erwartet red bei überschuss, bekommen %s", resp.Color)
	}
	if resp.Savings >= 0 {
		t.Errorf("erwartet negative savings, bekommen %.2f", resp.Savings)
	}
}

func TestUpdateBudgetInvalid(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"monthlyTotal": 0}`
	req := httptest.NewRequest(http.MethodPatch, "/api/budget", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 für monthlyTotal=0, bekommen %d", rec.Code)
	}

	body = `{"monthlyTotal": -100}`
	req = httptest.NewRequest(http.MethodPatch, "/api/budget", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 für monthlyTotal=-100, bekommen %d", rec.Code)
	}
}

func TestDeleteExpenseInvalidUUID(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodDelete, "/api/expenses/keine-uuid", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("erwartet 404 für ungültige UUID, bekommen %d", rec.Code)
	}
}

func decodeJSON(t *testing.T, body interface{ Read([]byte) (int, error) }, v interface{}) {
	t.Helper()
	b := make([]byte, 4096)
	n, _ := body.Read(b)
	if n > 0 {
		json.Unmarshal(b[:n], v)
	}
}
