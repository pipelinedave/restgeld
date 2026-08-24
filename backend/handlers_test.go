package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("erwartet status 'ok', bekommen '%s'", resp["status"])
	}
}

type failingStore struct {
	*memoryStore
	pingErr error
}

func (f *failingStore) Ping() error {
	return f.pingErr
}

func TestHealthEndpointError(t *testing.T) {
	store := &failingStore{
		memoryStore: newMemoryStore(),
		pingErr:     fmt.Errorf("db verbindung unterbrochen"),
	}
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("erwartet 503, bekommen %d: %s", rec.Code, rec.Body.String())
	}
}

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

	if len(resp.DailyStats) != 17 {
		t.Errorf("erwartet 17 Tage in DailyStats, bekommen %d", len(resp.DailyStats))
	}
}

func TestGetDailyExpensesStats(t *testing.T) {
	store := newMemoryStore()
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	// Ausgabe an Tag 1
	store.AddExpenseWithDate("2026-08", 12.0, "Einkauf", time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC))
	// Ausgabe an Tag 3
	store.AddExpenseWithDate("2026-08", 8.5, "Kino", time.Date(2026, 8, 3, 18, 0, 0, 0, time.UTC))

	stats, err := store.GetDailyExpenses("2026-08", start, 3)
	if err != nil {
		t.Fatalf("fehler: %v", err)
	}

	if len(stats) != 3 {
		t.Fatalf("erwartet 3 Tage, bekommen %d", len(stats))
	}
	if stats[0].Spent != 12.0 || stats[0].Date != "2026-08-01" {
		t.Errorf("Tag 1 unerwartet: %+v", stats[0])
	}
	if stats[1].Spent != 0.0 || stats[1].Date != "2026-08-02" {
		t.Errorf("Tag 2 unerwartet: %+v", stats[1])
	}
	if stats[2].Spent != 8.5 || stats[2].Date != "2026-08-03" {
		t.Errorf("Tag 3 unerwartet: %+v", stats[2])
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
	store.AddExpense("2026-08-01", 50.00, "Alt")
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodPost, "/api/period", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d", rec.Code)
	}

	var p Period
	json.NewDecoder(rec.Body).Decode(&p)
	if p.ID != "2026-08-25" {
		t.Errorf("erwartet period ID '2026-08-25', bekommen '%s'", p.ID)
	}

	total, _ := store.GetTotalExpenses("2026-08-25")
	if total != 0 {
		t.Errorf("erwartet 0 nach reset, bekommen %.2f", total)
	}
}

func TestNewPeriodWithCustomBudget(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	srv := &server{store: store, now: func() time.Time { return now }}

	body := `{"monthlyTotal": 400, "startDate": "2026-08-25"}`
	req := httptest.NewRequest(http.MethodPost, "/api/period", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d", rec.Code)
	}

	var p Period
	json.NewDecoder(rec.Body).Decode(&p)
	if p.MonthlyTotal != 400 {
		t.Errorf("erwartet monthlyTotal 400, bekommen %.2f", p.MonthlyTotal)
	}
	if p.ID != "2026-08-25" {
		t.Errorf("erwartet period ID '2026-08-25', bekommen '%s'", p.ID)
	}
}

func TestNewPeriodWithCustomBudgetAndDays(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"monthlyTotal": 280, "startDate": "2026-08-25", "days": 14}`
	req := httptest.NewRequest(http.MethodPost, "/api/period", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	p, _ := store.GetOrCreatePeriod()
	if p.MonthlyTotal != 280 {
		t.Errorf("erwartet monthlyTotal 280, bekommen %.2f", p.MonthlyTotal)
	}
	if p.MonthDays != 14 {
		t.Errorf("erwartet monthDays 14, bekommen %d", p.MonthDays)
	}
	if p.BaseBudget != 20.0 {
		t.Errorf("erwartet baseBudget 20.00, bekommen %.2f", p.BaseBudget)
	}
}

func TestUpdateBudgetAndDays(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"monthlyTotal": 300, "days": 15}`
	req := httptest.NewRequest(http.MethodPatch, "/api/budget", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	p, _ := store.GetOrCreatePeriod()
	if p.MonthlyTotal != 300 {
		t.Errorf("erwartet monthly_total 300, bekommen %.2f", p.MonthlyTotal)
	}
	if p.MonthDays != 15 {
		t.Errorf("erwartet monthDays 15, bekommen %d", p.MonthDays)
	}
	if p.BaseBudget != 20.0 {
		t.Errorf("erwartet baseBudget 20.00, bekommen %.2f", p.BaseBudget)
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
		t.Errorf("erwartet savings = baseBudget (%.2f) an tag 1 (0 ausgegeben), bekommen %.2f", resp.BaseBudget, resp.Savings)
	}
	if resp.CurrentBudget != resp.BaseBudget {
		t.Errorf("erwartet current == base (%.2f) an tag 1, bekommen %.2f", resp.BaseBudget, resp.CurrentBudget)
	}
	if resp.Color != "green" {
		t.Errorf("erwartet green an Tag 1 mit vollem Budget, bekommen %s", resp.Color)
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

func TestGetExpensesDefaultPagination(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp PaginatedExpenses
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}

	if resp.Total != 0 || len(resp.Items) != 0 || resp.Page != 1 || resp.Limit != 10 || resp.TotalPages != 1 {
		t.Errorf("unerwartete antwort fuer leere liste: %+v", resp)
	}
}

func TestGetExpensesCustomPagination(t *testing.T) {
	store := newMemoryStore()
	for i := 1; i <= 25; i++ {
		store.AddExpense("2026-08", float64(i), fmt.Sprintf("Ausgabe %d", i))
	}

	srv := &server{store: store, now: time.Now}

	// Seite 2 mit Limit 10 -> sollte 10 Elemente enthalten (Ausgabe 15 bis 6)
	req := httptest.NewRequest(http.MethodGet, "/api/expenses?page=2&limit=10", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp PaginatedExpenses
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}

	if resp.Total != 25 {
		t.Errorf("erwartet Total=25, bekommen %d", resp.Total)
	}
	if resp.Page != 2 {
		t.Errorf("erwartet Page=2, bekommen %d", resp.Page)
	}
	if resp.Limit != 10 {
		t.Errorf("erwartet Limit=10, bekommen %d", resp.Limit)
	}
	if resp.TotalPages != 3 {
		t.Errorf("erwartet TotalPages=3, bekommen %d", resp.TotalPages)
	}
	if len(resp.Items) != 10 {
		t.Errorf("erwartet 10 Items auf Seite 2, bekommen %d", len(resp.Items))
	}
	// Letztes Element auf Seite 2 (bei 25 Einträgen, Seite 2: 11. bis 20. neuste Ausgabe -> Ausgabe 15 bis 6)
	if resp.Items[0].Note != "Ausgabe 15" {
		t.Errorf("erwartet erstes Element auf Seite 2 'Ausgabe 15', bekommen '%s'", resp.Items[0].Note)
	}

	// Seite 3 mit Limit 10 -> sollte die restlichen 5 Elemente enthalten
	req3 := httptest.NewRequest(http.MethodGet, "/api/expenses?page=3&limit=10", nil)
	rec3 := httptest.NewRecorder()
	srv.router().ServeHTTP(rec3, req3)

	var resp3 PaginatedExpenses
	json.NewDecoder(rec3.Body).Decode(&resp3)
	if len(resp3.Items) != 5 {
		t.Errorf("erwartet 5 Items auf Seite 3, bekommen %d", len(resp3.Items))
	}
}

func TestGetExpensesInvalidParams(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/expenses?page=-5&limit=abc", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", rec.Code)
	}

	var resp PaginatedExpenses
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Page != 1 || resp.Limit != 10 {
		t.Errorf("erwartet fallback auf Page=1, Limit=10, bekommen Page=%d, Limit=%d", resp.Page, resp.Limit)
	}
}

func TestExportJSON(t *testing.T) {
	store := newMemoryStore()
	store.AddExpense("2026-08", 14.50, "Kaffee")
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/export?format=json", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Errorf("erwartet application/json content-type, bekommen %s", rec.Header().Get("Content-Type"))
	}

	var backup ExportBackup
	if err := json.NewDecoder(rec.Body).Decode(&backup); err != nil {
		t.Fatalf("json decode fehler: %v", err)
	}
	if len(backup.Expenses) != 1 || backup.Expenses[0].Amount != 14.50 {
		t.Errorf("unerwartete exportierte Ausgaben: %+v", backup.Expenses)
	}
}

func TestExportCSV(t *testing.T) {
	store := newMemoryStore()
	store.AddExpense("2026-08", 12.00, "Mittagessen")
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodGet, "/api/export?format=csv", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/csv") {
		t.Errorf("erwartet text/csv content-type, bekommen %s", rec.Header().Get("Content-Type"))
	}

	csvContent := rec.Body.String()
	if !strings.Contains(csvContent, "Datum;Uhrzeit;Betrag;Notiz") {
		t.Errorf("erwartet CSV Header, bekommen: %s", csvContent)
	}
	if !strings.Contains(csvContent, "12.00;\"Mittagessen\"") {
		t.Errorf("erwartet CSV Eintrag, bekommen: %s", csvContent)
	}
}

func TestImportJSON(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	body := `{"expenses":[{"amount": 9.99, "note": "Buch"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/import", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["imported"].(float64) != 1 {
		t.Errorf("erwartet 1 importiert, bekommen %v", resp["imported"])
	}

	exp, _ := store.GetAllExpenses("2026-08")
	if len(exp) != 1 || exp[0].Amount != 9.99 {
		t.Errorf("unerwartete Ausgaben im Store: %+v", exp)
	}
}

func TestImportCSV(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	csv := "Datum;Uhrzeit;Betrag;Notiz\n2026-08-10;12:00:00;15.50;\"Supermarkt\"\n2026-08-11;14:00:00;5.20;\"Eis\""
	req := httptest.NewRequest(http.MethodPost, "/api/import", strings.NewReader(csv))
	req.Header.Set("Content-Type", "text/csv")
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["imported"].(float64) != 2 {
		t.Errorf("erwartet 2 importiert, bekommen %v", resp["imported"])
	}

	exp, _ := store.GetAllExpenses("2026-08")
	if len(exp) != 2 {
		t.Errorf("erwartet 2 Ausgaben im Store, bekommen %d", len(exp))
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
