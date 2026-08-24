package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type postgresStore struct {
	db *sql.DB
}

func newPostgresStoreFromDB(db *sql.DB) Store {
	createTables(db)
	return &postgresStore{db: db}
}

func newPostgresStore() Store {
	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "5432")
	user := envOrDefault("DB_USER", "restgeld")
	password := envOrDefault("DB_PASSWORD", "restgeld")
	dbname := envOrDefault("DB_NAME", "restgeld")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("fehler beim db-verbindungsaufbau: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("fehler beim db-ping: %v", err)
	}

	createTables(db)
	log.Println("db-verbindung hergestellt und tabellen bereit")

	return &postgresStore{db: db}
}

func createTables(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS periods (
		id TEXT PRIMARY KEY,
		start_date TIMESTAMPTZ NOT NULL,
		month_days INTEGER NOT NULL,
		base_budget NUMERIC(10,2) NOT NULL,
		monthly_total NUMERIC(10,2) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS expenses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		period_id TEXT NOT NULL REFERENCES periods(id) ON DELETE CASCADE,
		amount NUMERIC(10,2) NOT NULL,
		note TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_expenses_period ON expenses(period_id);
	CREATE INDEX IF NOT EXISTS idx_expenses_created ON expenses(created_at DESC);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("fehler beim tabellen-erstellen: %v", err)
	}
}

func (s *postgresStore) GetOrCreatePeriod() (*Period, error) {
	now := time.Now()
	periodID := fmt.Sprintf("%d-%02d", now.Year(), now.Month())

	var p Period
	err := s.db.QueryRow(
		"SELECT id, start_date, month_days, base_budget, monthly_total FROM periods WHERE id = $1",
		periodID,
	).Scan(&p.ID, &p.StartDate, &p.MonthDays, &p.BaseBudget, &p.MonthlyTotal)

	if err == sql.ErrNoRows {
		return s.createPeriodInternal(now)
	}

	if err != nil {
		return nil, fmt.Errorf("period lesen: %w", err)
	}

	return &p, nil
}

func (s *postgresStore) createPeriodInternal(now time.Time) (*Period, error) {
	periodID := fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := startOfMonth.AddDate(0, 1, 0)
	monthDays := nextMonth.Sub(startOfMonth).Hours() / 24
	defaultBudget := 450.0
	baseBudget := parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))
	baseDaily := mathRound(baseBudget/monthDays, 2)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction start: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO periods (id, start_date, month_days, base_budget, monthly_total) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET start_date=$2, month_days=$3, base_budget=$4, monthly_total=$5",
		periodID, startOfMonth, int(monthDays), baseDaily, baseBudget,
	)
	if err != nil {
		return nil, fmt.Errorf("period upsert: %w", err)
	}

	_, err = tx.Exec("DELETE FROM expenses WHERE period_id = $1", periodID)
	if err != nil {
		return nil, fmt.Errorf("expenses löschen: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", err)
	}

	return &Period{
		ID:           periodID,
		StartDate:    startOfMonth,
		MonthDays:    int(monthDays),
		BaseBudget:   baseDaily,
		MonthlyTotal: baseBudget,
	}, nil
}

func (s *postgresStore) CreatePeriod() (*Period, error) {
	return s.createPeriodInternal(time.Now())
}

func (s *postgresStore) UpdateBudget(newTotal float64) error {
	now := time.Now()
	periodID := fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := startOfMonth.AddDate(0, 1, 0)
	monthDays := nextMonth.Sub(startOfMonth).Hours() / 24
	baseDaily := mathRound(newTotal/monthDays, 2)

	_, err := s.db.Exec(
		"INSERT INTO periods (id, start_date, month_days, base_budget, monthly_total) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET base_budget=$4, monthly_total=$5",
		periodID, startOfMonth, int(monthDays), baseDaily, newTotal,
	)
	return err
}

func (s *postgresStore) GetTotalExpenses(periodID string) (float64, error) {
	var total sql.NullFloat64
	err := s.db.QueryRow(
		"SELECT SUM(amount) FROM expenses WHERE period_id = $1",
		periodID,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

func (s *postgresStore) GetRecentExpenses(periodID string, limit int) ([]Expense, error) {
	rows, err := s.db.Query(
		"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 ORDER BY created_at DESC LIMIT $2",
		periodID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	if expenses == nil {
		expenses = []Expense{}
	}

	return expenses, nil
}

func (s *postgresStore) AddExpense(periodID string, amount float64, note string) (*Expense, error) {
	var e Expense
	err := s.db.QueryRow(
		"INSERT INTO expenses (period_id, amount, note) VALUES ($1, $2, $3) RETURNING id, period_id, amount, note, created_at",
		periodID, amount, note,
	).Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *postgresStore) DeleteExpense(expenseID string) error {
	result, err := s.db.Exec("DELETE FROM expenses WHERE id = $1", expenseID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("ausgabe nicht gefunden")
	}
	return nil
}

func (s *postgresStore) Ping() error {
	return s.db.Ping()
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
