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
	if err := RunMigrations(db); err != nil {
		log.Fatalf("fehler bei db-migration: %v", err)
	}
	return &postgresStore{db: db}
}

func newPostgresStore() Store {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := envOrDefault("DB_HOST", "localhost")
		port := envOrDefault("DB_PORT", "5432")
		user := envOrDefault("DB_USER", "restgeld")
		password := envOrDefault("DB_PASSWORD", "restgeld")
		dbname := envOrDefault("DB_NAME", "restgeld")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

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

	if err := RunMigrations(db); err != nil {
		log.Fatalf("fehler bei db-migration: %v", err)
	}
	log.Println("db-verbindung hergestellt und migrationen ausgefuehrt")

	return &postgresStore{db: db}
}

func (s *postgresStore) GetOrCreatePeriod() (*Period, error) {
	now := time.Now()

	var p Period
	err := s.db.QueryRow(
		`SELECT id, start_date, month_days, base_budget, monthly_total 
		 FROM periods 
		 WHERE start_date <= $1 AND (start_date + (month_days || ' days')::interval) > $1 
		 ORDER BY start_date DESC LIMIT 1`,
		now,
	).Scan(&p.ID, &p.StartDate, &p.MonthDays, &p.BaseBudget, &p.MonthlyTotal)

	if err == sql.ErrNoRows {
		// Falls keine aktive Periode existiert, prüfe ob es eine Vorperiode gibt, um das Budget zu übernehmen
		var lastTotal float64
		err = s.db.QueryRow("SELECT monthly_total FROM periods ORDER BY start_date DESC LIMIT 1").Scan(&lastTotal)
		if err != nil {
			defaultBudget := 450.0
			lastTotal = parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))
		}
		return s.CreatePeriodWithStart(now, lastTotal, 0)
	}

	if err != nil {
		return nil, fmt.Errorf("period lesen: %w", err)
	}

	return &p, nil
}

func (s *postgresStore) CreatePeriodWithStart(start time.Time, monthlyTotal float64, days int) (*Period, error) {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	periodID := startDay.Format("2006-01-02")
	monthDays := days
	if monthDays <= 0 {
		monthDays = calcPeriodDays(startDay)
	}

	if monthlyTotal <= 0 {
		defaultBudget := 450.0
		monthlyTotal = parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))
	}
	baseDaily := mathRound(monthlyTotal/float64(monthDays), 2)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction start: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO periods (id, start_date, month_days, base_budget, monthly_total) 
		 VALUES ($1, $2, $3, $4, $5) 
		 ON CONFLICT (id) DO UPDATE SET start_date=$2, month_days=$3, base_budget=$4, monthly_total=$5`,
		periodID, startDay, monthDays, baseDaily, monthlyTotal,
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
		StartDate:    startDay,
		MonthDays:    monthDays,
		BaseBudget:   baseDaily,
		MonthlyTotal: monthlyTotal,
	}, nil
}

func (s *postgresStore) CreatePeriod() (*Period, error) {
	var lastTotal float64
	var lastDays int
	err := s.db.QueryRow("SELECT monthly_total, month_days FROM periods ORDER BY start_date DESC LIMIT 1").Scan(&lastTotal, &lastDays)
	if err != nil {
		defaultBudget := 450.0
		lastTotal = parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))
		lastDays = 0
	}
	return s.CreatePeriodWithStart(time.Now(), lastTotal, lastDays)
}

func (s *postgresStore) UpdateBudget(newTotal float64, days int) error {
	p, err := s.GetOrCreatePeriod()
	if err != nil {
		return err
	}

	monthDays := p.MonthDays
	if days > 0 {
		monthDays = days
	}

	if newTotal <= 0 {
		newTotal = p.MonthlyTotal
	}

	baseDaily := mathRound(newTotal/float64(monthDays), 2)
	_, err = s.db.Exec(
		"UPDATE periods SET base_budget = $1, monthly_total = $2, month_days = $3 WHERE id = $4",
		baseDaily, newTotal, monthDays, p.ID,
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

func (s *postgresStore) GetTodayExpenses(periodID string, now time.Time) (float64, error) {
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.AddDate(0, 0, 1)

	var total sql.NullFloat64
	err := s.db.QueryRow(
		"SELECT SUM(amount) FROM expenses WHERE period_id = $1 AND created_at >= $2 AND created_at < $3",
		periodID, startOfToday, endOfToday,
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

func (s *postgresStore) GetExpenses(periodID string, page, limit int) (*PaginatedExpenses, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM expenses WHERE period_id = $1", periodID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("anzahl ausgaben ermitteln: %w", err)
	}

	offset := (page - 1) * limit
	rows, err := s.db.Query(
		"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		periodID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("ausgaben laden: %w", err)
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("ausgabe scannen: %w", err)
		}
		expenses = append(expenses, e)
	}

	if expenses == nil {
		expenses = []Expense{}
	}

	totalPages := 1
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &PaginatedExpenses{
		Items:      expenses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
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
