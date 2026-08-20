//go:build integration

package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

func sqlOpen(driver, connStr string) (*sql.DB, error) {
	return sql.Open(driver, connStr)
}

func TestIntegrationCreateAndReadPeriod(t *testing.T) {
	store := newIntegrationStore(t)
	p, err := store.GetOrCreatePeriod()
	if err != nil {
		t.Fatalf("get or create: %v", err)
	}
	if p.ID == "" {
		t.Fatal("period id sollte nicht leer sein")
	}
	if p.MonthDays <= 0 {
		t.Fatalf("ungültige month_days: %d", p.MonthDays)
	}
	if p.BaseBudget <= 0 {
		t.Fatalf("ungültiger base_budget: %.2f", p.BaseBudget)
	}
}

func TestIntegrationCreateDuplicatePeriod(t *testing.T) {
	store := newIntegrationStore(t)
	_, err := store.CreatePeriod()
	if err != nil {
		t.Fatalf("erste periode: %v", err)
	}
	_, err = store.CreatePeriod()
	if err != nil {
		t.Fatalf("zweite periode (on conflict): %v", err)
	}
}

func TestIntegrationAddAndListExpenses(t *testing.T) {
	store := newIntegrationStore(t)
	p, _ := store.GetOrCreatePeriod()

	e1, err := store.AddExpense(p.ID, 8.50, "Frühstück")
	if err != nil {
		t.Fatalf("add expense 1: %v", err)
	}
	if e1.Amount != 8.50 {
		t.Errorf("amount: erwartet 8.50, bekommen %.2f", e1.Amount)
	}
	if e1.Note != "Frühstück" {
		t.Errorf("note: erwartet 'Frühstück', bekommen '%s'", e1.Note)
	}
	if e1.ID == "" {
		t.Error("id sollte generiert werden")
	}

	store.AddExpense(p.ID, 3.50, "Kaffee")

	expenses, err := store.GetRecentExpenses(p.ID, 2)
	if err != nil {
		t.Fatalf("list expenses: %v", err)
	}
	if len(expenses) != 2 {
		t.Fatalf("erwartet 2 expenses, bekommen %d", len(expenses))
	}
	if expenses[0].Note != "Kaffee" {
		t.Errorf("erwartet neueste zuerst (Kaffee), bekommen '%s'", expenses[0].Note)
	}
}

func TestIntegrationDeleteExpense(t *testing.T) {
	store := newIntegrationStore(t)
	p, _ := store.GetOrCreatePeriod()

	e, _ := store.AddExpense(p.ID, 5.00, "Test")
	err := store.DeleteExpense(e.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	total, _ := store.GetTotalExpenses(p.ID)
	if total != 0 {
		t.Errorf("erwartet 0 nach delete, bekommen %.2f", total)
	}
}

func TestIntegrationDeleteExpenseNotFound(t *testing.T) {
	store := newIntegrationStore(t)
	err := store.DeleteExpense("00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("erwartet fehler bei nicht-existenter id")
	}
}

func TestIntegrationGetTotalExpenses(t *testing.T) {
	store := newIntegrationStore(t)
	p, _ := store.GetOrCreatePeriod()

	store.AddExpense(p.ID, 10.00, "")
	store.AddExpense(p.ID, 20.00, "")
	store.AddExpense(p.ID, 30.00, "")

	total, err := store.GetTotalExpenses(p.ID)
	if err != nil {
		t.Fatalf("get total: %v", err)
	}
	if total != 60.00 {
		t.Errorf("erwartet 60.00, bekommen %.2f", total)
	}
}

func TestIntegrationUpdateBudget(t *testing.T) {
	store := newIntegrationStore(t)
	err := store.UpdateBudget(600)
	if err != nil {
		t.Fatalf("update budget: %v", err)
	}

	p, _ := store.GetOrCreatePeriod()
	if p.MonthlyTotal != 600 {
		t.Errorf("erwartet monthly_total 600, bekommen %.2f", p.MonthlyTotal)
	}
}

func TestIntegrationCreatePeriodClearsExpenses(t *testing.T) {
	store := newIntegrationStore(t)
	p, _ := store.GetOrCreatePeriod()
	store.AddExpense(p.ID, 100.00, "Alt")

	_, err := store.CreatePeriod()
	if err != nil {
		t.Fatalf("create period: %v", err)
	}

	expenses, _ := store.GetRecentExpenses(p.ID, 10)
	if len(expenses) != 0 {
		t.Errorf("erwartet 0 expenses nach reset, bekommen %d", len(expenses))
	}

	total, _ := store.GetTotalExpenses(p.ID)
	if total != 0 {
		t.Errorf("erwartet total 0 nach reset, bekommen %.2f", total)
	}
}

func newIntegrationStore(t *testing.T) Store {
	t.Helper()
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "restgeld")
	password := getEnv("DB_PASSWORD", "restgeld")
	dbname := getEnv("DB_NAME", "restgeld")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlOpen("postgres", connStr)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("db ping: %v", err)
	}

	store := newPostgresStoreFromDB(db)

	// Clean tables
	db.Exec("DELETE FROM expenses")
	db.Exec("DELETE FROM periods")
	time.Sleep(50 * time.Millisecond)

	return store
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
