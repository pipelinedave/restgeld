//go:build integration

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func newIntegrationServer(t *testing.T) *server {
	t.Helper()
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "restgeld")
	password := getEnv("DB_PASSWORD", "restgeld")
	dbname := getEnv("DB_NAME", "restgeld")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("db ping: %v", err)
	}

	store := newPostgresStoreFromDB(db)
	os.Setenv("DEFAULT_MONTHLY_BUDGET", "450")

	// Clean
	db.Exec("DELETE FROM expenses")
	db.Exec("DELETE FROM periods")

	return &server{store: store, now: time.Now}
}

func TestAPI_Health(t *testing.T) {
	srv := newIntegrationServer(t)
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

func TestAPI_GetBudget_CreatesPeriod(t *testing.T) {
	srv := newIntegrationServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var resp BudgetResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.PeriodID == "" {
		t.Error("period id erwartet")
	}
	if resp.MonthDays <= 0 {
		t.Errorf("ungültige month_days: %d", resp.MonthDays)
	}
	if resp.BaseBudget <= 0 {
		t.Errorf("ungültiger base_budget: %.2f", resp.BaseBudget)
	}
}

func TestAPI_FullFlow(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	// 1. Anfangsbudget holen
	rec := get(router, "/api/budget")
	var initial BudgetResponse
	json.NewDecoder(rec.Body).Decode(&initial)
	initialBudget := initial.CurrentBudget

	// 2. Ausgabe buchen
	body := `{"amount": 15.00, "note": "Mittagessen"}`
	rec = post(router, "/api/expenses", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var expense Expense
	json.NewDecoder(rec.Body).Decode(&expense)
	if expense.Amount != 15.00 {
		t.Errorf("amount: erwartet 15, bekommen %.2f", expense.Amount)
	}
	if expense.Note != "Mittagessen" {
		t.Errorf("note: erwartet 'Mittagessen', bekommen '%s'", expense.Note)
	}

	// 3. Budget nach Ausgabe
	rec = get(router, "/api/budget")
	var afterExpense BudgetResponse
	json.NewDecoder(rec.Body).Decode(&afterExpense)

	if afterExpense.CurrentBudget >= initialBudget {
		t.Errorf("budget sollte nach ausgabe sinken: vorher %.2f, nachher %.2f",
			initialBudget, afterExpense.CurrentBudget)
	}
	if len(afterExpense.Expenses) != 1 {
		t.Errorf("erwartet 1 expense, bekommen %d", len(afterExpense.Expenses))
	}

	// 4. Ausgabe löschen
	rec = del(router, "/api/expenses/"+expense.ID)
	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 bei delete, bekommen %d", rec.Code)
	}

	// 5. Budget nach Löschen
	rec = get(router, "/api/budget")
	var afterDelete BudgetResponse
	json.NewDecoder(rec.Body).Decode(&afterDelete)

	if afterDelete.CurrentBudget != initialBudget {
		t.Errorf("budget sollte nach delete wieder initial sein: erwartet %.2f, bekommen %.2f",
			initialBudget, afterDelete.CurrentBudget)
	}
}

func TestAPI_ConcurrentExpenses(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	// Initialisiere die Periode vorab, damit nicht mehrere Goroutines gleichzeitig upsert ausführen
	srv.store.GetOrCreatePeriod("")

	// 3 parallele Requests gleichzeitig
	n := 3
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			body := fmt.Sprintf(`{"amount": %.2f, "note": "gleichzeitig-%d"}`, 10.0, idx)
			rec := post(router, "/api/expenses", body)
			if rec.Code != http.StatusCreated {
				errs <- fmt.Errorf("request %d: status %d", idx, rec.Code)
				return
			}
			errs <- nil
		}(i)
	}

	for i := 0; i < n; i++ {
		if err := <-errs; err != nil {
			t.Error(err)
		}
	}

	// Prüfen ob alle 3 Expenses vorhanden sind
	rec := get(router, "/api/budget")
	var resp BudgetResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Expenses) != 3 {
		t.Errorf("erwartet 3 expenses nach parallelen requests, bekommen %d", len(resp.Expenses))
	}
}

func TestAPI_UpdateBudgetAndReset(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	// Auf 600 ändern
	body := `{"monthlyTotal": 600}`
	rec := patch(router, "/api/budget", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", rec.Code)
	}

	rec = get(router, "/api/budget")
	var afterUpdate BudgetResponse
	json.NewDecoder(rec.Body).Decode(&afterUpdate)
	if afterUpdate.BaseBudget <= 14.52 {
		t.Errorf("baseBudget sollte nach update > 14.52 sein, bekommen %.2f", afterUpdate.BaseBudget)
	}

	// Neue Periode starten (reset)
	rec = post(router, "/api/period", "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201 bei new period, bekommen %d", rec.Code)
	}

	rec = get(router, "/api/budget")
	var afterReset BudgetResponse
	json.NewDecoder(rec.Body).Decode(&afterReset)
	if afterReset.BaseBudget != 14.52 {
		t.Errorf("baseBudget sollte nach reset wieder 14.52 sein, bekommen %.2f", afterReset.BaseBudget)
	}
}

func TestAPI_DeleteNonexistent(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	rec := del(router, "/api/expenses/00000000-0000-0000-0000-000000000000")
	if rec.Code != http.StatusNotFound {
		t.Errorf("erwartet 404, bekommen %d", rec.Code)
	}
}

func TestAPI_InvalidAmount(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	body := `{"amount": -5}`
	rec := post(router, "/api/expenses", body)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 für negative amount, bekommen %d", rec.Code)
	}

	body = `{"amount": 0}`
	rec = post(router, "/api/expenses", body)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 für amount=0, bekommen %d", rec.Code)
	}
}

func TestAPI_GetDayExpenses(t *testing.T) {
	srv := newIntegrationServer(t)
	router := srv.router()

	// Ausgabe buchen
	body := `{"amount": 4.50, "note": "Snack"}`
	rec := post(router, "/api/expenses", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("erwartet 201, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var expense Expense
	json.NewDecoder(rec.Body).Decode(&expense)
	day := expense.CreatedAt.Format("2006-01-02")

	// Buchung des Tages holen
	rec = get(router, "/api/expenses/day?date="+day)
	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	var list []Expense
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("erwartet 1 tagesausgabe, bekommen %d", len(list))
	}
	if list[0].Amount != 4.50 {
		t.Errorf("amount: erwartet 4.50, bekommen %.2f", list[0].Amount)
	}

	// Anderer Tag ohne Buchungen -> leere Liste
	oldDay := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	rec = get(router, "/api/expenses/day?date="+oldDay)
	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200 fuer leeren tag, bekommen %d", rec.Code)
	}
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 0 {
		t.Errorf("erwartet 0 tagesausgaben, bekommen %d", len(list))
	}

	// Fehlerfall: ungueltiges Datum
	rec = get(router, "/api/expenses/day?date=kaputt")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 fuer ungueltiges datum, bekommen %d", rec.Code)
	}

	// Fehlerfall: fehlendes Datum
	rec = get(router, "/api/expenses/day")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("erwartet 400 fuer fehlendes datum, bekommen %d", rec.Code)
	}
}
func get(router http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func post(router http.Handler, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(http.MethodPost, path, nil)
	} else {
		req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func del(router http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func patch(router http.Handler, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
