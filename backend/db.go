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

func (s *postgresStore) GetOrCreatePeriod(userID string) (*Period, error) {
	now := time.Now()

	var p Period
	var err error
	if userID != "" {
		err = s.db.QueryRow(
			`SELECT id, user_id, start_date, month_days, base_budget, monthly_total 
			 FROM periods 
			 WHERE user_id = $1 AND start_date <= $2 AND (start_date + (month_days || ' days')::interval) > $2 
			 ORDER BY start_date DESC LIMIT 1`,
			userID, now,
		).Scan(&p.ID, &p.UserID, &p.StartDate, &p.MonthDays, &p.BaseBudget, &p.MonthlyTotal)
	} else {
		var uid sql.NullString
		err = s.db.QueryRow(
			`SELECT id, user_id, start_date, month_days, base_budget, monthly_total 
			 FROM periods 
			 WHERE user_id IS NULL AND start_date <= $1 AND (start_date + (month_days || ' days')::interval) > $1 
			 ORDER BY start_date DESC LIMIT 1`,
			now,
		).Scan(&p.ID, &uid, &p.StartDate, &p.MonthDays, &p.BaseBudget, &p.MonthlyTotal)
		if uid.Valid {
			p.UserID = uid.String
		}
	}

	if err == sql.ErrNoRows {
		// Default budget aus user settings oder env
		var lastTotal float64
		defaultBudget := 450.0
		lastTotal = parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))

		if userID != "" {
			var uBudget float64
			var uDays int
			if uErr := s.db.QueryRow("SELECT default_monthly_budget, default_period_days FROM users WHERE id = $1", userID).Scan(&uBudget, &uDays); uErr == nil && uBudget > 0 {
				lastTotal = uBudget
				return s.CreatePeriodWithStart(userID, now, lastTotal, uDays)
			}
		}

		return s.CreatePeriodWithStart(userID, now, lastTotal, 0)
	}

	if err != nil {
		return nil, fmt.Errorf("period lesen: %w", err)
	}

	return &p, nil
}

func (s *postgresStore) CreatePeriodWithStart(userID string, start time.Time, monthlyTotal float64, days int) (*Period, error) {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	periodID := startDay.Format("2006-01-02")
	if userID != "" {
		periodID = fmt.Sprintf("%s_%s", userID, startDay.Format("2006-01-02"))
	}

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

	if userID != "" {
		_, err = tx.Exec(
			`INSERT INTO periods (id, user_id, start_date, month_days, base_budget, monthly_total) 
			 VALUES ($1, $2, $3, $4, $5, $6) 
			 ON CONFLICT (id) DO UPDATE SET start_date=$3, month_days=$4, base_budget=$5, monthly_total=$6`,
			periodID, userID, startDay, monthDays, baseDaily, monthlyTotal,
		)
	} else {
		_, err = tx.Exec(
			`INSERT INTO periods (id, start_date, month_days, base_budget, monthly_total) 
			 VALUES ($1, $2, $3, $4, $5) 
			 ON CONFLICT (id) DO UPDATE SET start_date=$2, month_days=$3, base_budget=$4, monthly_total=$5`,
			periodID, startDay, monthDays, baseDaily, monthlyTotal,
		)
	}
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
		UserID:       userID,
		StartDate:    startDay,
		MonthDays:    monthDays,
		BaseBudget:   baseDaily,
		MonthlyTotal: monthlyTotal,
	}, nil
}

func (s *postgresStore) CreatePeriod(userID string) (*Period, error) {
	var lastTotal float64
	var lastDays int
	var err error
	if userID != "" {
		err = s.db.QueryRow("SELECT monthly_total, month_days FROM periods WHERE user_id = $1 ORDER BY start_date DESC LIMIT 1", userID).Scan(&lastTotal, &lastDays)
	} else {
		err = s.db.QueryRow("SELECT monthly_total, month_days FROM periods WHERE user_id IS NULL ORDER BY start_date DESC LIMIT 1").Scan(&lastTotal, &lastDays)
	}
	if err != nil {
		defaultBudget := 450.0
		lastTotal = parseFloat(envOrDefault("DEFAULT_MONTHLY_BUDGET", fmt.Sprintf("%.0f", defaultBudget)))
		lastDays = 0
	}
	return s.CreatePeriodWithStart(userID, time.Now(), lastTotal, lastDays)
}

func (s *postgresStore) UpdateBudget(userID string, newTotal float64, days int) error {
	p, err := s.GetOrCreatePeriod(userID)
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

func (s *postgresStore) GetTotalExpenses(userID, periodID string) (float64, error) {
	var total sql.NullFloat64
	var err error
	if userID != "" {
		err = s.db.QueryRow("SELECT SUM(amount) FROM expenses WHERE period_id = $1 AND user_id = $2", periodID, userID).Scan(&total)
	} else {
		err = s.db.QueryRow("SELECT SUM(amount) FROM expenses WHERE period_id = $1 AND user_id IS NULL", periodID).Scan(&total)
	}

	if err != nil {
		return 0, err
	}
	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

func (s *postgresStore) GetTodayExpenses(userID, periodID string, now time.Time) (float64, error) {
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.AddDate(0, 0, 1)

	var total sql.NullFloat64
	var err error
	if userID != "" {
		err = s.db.QueryRow(
			"SELECT SUM(amount) FROM expenses WHERE period_id = $1 AND user_id = $2 AND created_at >= $3 AND created_at < $4",
			periodID, userID, startOfToday, endOfToday,
		).Scan(&total)
	} else {
		err = s.db.QueryRow(
			"SELECT SUM(amount) FROM expenses WHERE period_id = $1 AND user_id IS NULL AND created_at >= $2 AND created_at < $3",
			periodID, startOfToday, endOfToday,
		).Scan(&total)
	}

	if err != nil {
		return 0, err
	}
	if total.Valid {
		return total.Float64, nil
	}
	return 0, nil
}

func (s *postgresStore) GetRecentExpenses(userID, periodID string, limit int) ([]Expense, error) {
	var rows *sql.Rows
	var err error
	if userID != "" {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id = $2 ORDER BY created_at DESC LIMIT $3",
			periodID, userID, limit,
		)
	} else {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id IS NULL ORDER BY created_at DESC LIMIT $2",
			periodID, limit,
		)
	}
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
		e.UserID = userID
		expenses = append(expenses, e)
	}

	if expenses == nil {
		expenses = []Expense{}
	}
	return expenses, nil
}

func (s *postgresStore) GetExpenses(userID, periodID string, page, limit int) (*PaginatedExpenses, error) {
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
	var err error
	if userID != "" {
		err = s.db.QueryRow("SELECT COUNT(*) FROM expenses WHERE period_id = $1 AND user_id = $2", periodID, userID).Scan(&total)
	} else {
		err = s.db.QueryRow("SELECT COUNT(*) FROM expenses WHERE period_id = $1 AND user_id IS NULL", periodID).Scan(&total)
	}
	if err != nil {
		return nil, fmt.Errorf("anzahl ausgaben ermitteln: %w", err)
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}
	offset := (page - 1) * limit

	var rows *sql.Rows
	if userID != "" {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4",
			periodID, userID, limit, offset,
		)
	} else {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3",
			periodID, limit, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("ausgaben abrufen: %w", err)
	}
	defer rows.Close()

	var items []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("ausgabe scannen: %w", err)
		}
		e.UserID = userID
		items = append(items, e)
	}

	if items == nil {
		items = []Expense{}
	}

	return &PaginatedExpenses{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *postgresStore) AddExpense(userID, periodID string, amount float64, note string) (*Expense, error) {
	return s.AddExpenseWithDate(userID, periodID, amount, note, time.Now())
}

func (s *postgresStore) AddExpenseWithDate(userID, periodID string, amount float64, note string, createdAt time.Time) (*Expense, error) {
	var e Expense
	var err error
	if userID != "" {
		err = s.db.QueryRow(
			`INSERT INTO expenses (period_id, user_id, amount, note, created_at) 
			 VALUES ($1, $2, $3, $4, $5) 
			 RETURNING id, period_id, amount, note, created_at`,
			periodID, userID, amount, note, createdAt,
		).Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt)
	} else {
		err = s.db.QueryRow(
			`INSERT INTO expenses (period_id, amount, note, created_at) 
			 VALUES ($1, $2, $3, $4) 
			 RETURNING id, period_id, amount, note, created_at`,
			periodID, amount, note, createdAt,
		).Scan(&e.ID, &e.PeriodID, &e.Amount, &e.Note, &e.CreatedAt)
	}

	if err != nil {
		return nil, fmt.Errorf("ausgabe speichern: %w", err)
	}
	e.UserID = userID
	return &e, nil
}

func (s *postgresStore) GetDailyExpenses(userID, periodID string, start time.Time, upToDay int) ([]DailyStat, error) {
	if upToDay < 1 {
		upToDay = 1
	}

	stats := make([]DailyStat, upToDay)
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	for i := 0; i < upToDay; i++ {
		currentDate := startDay.AddDate(0, 0, i)
		stats[i] = DailyStat{
			Day:   i + 1,
			Date:  currentDate.Format("2006-01-02"),
			Spent: 0,
		}
	}

	var rows *sql.Rows
	var err error
	if userID != "" {
		query := `SELECT DATE(created_at AT TIME ZONE 'UTC') AS exp_date, SUM(amount) 
		          FROM expenses 
		          WHERE period_id = $1 AND user_id = $2
		          GROUP BY exp_date`
		rows, err = s.db.Query(query, periodID, userID)
	} else {
		query := `SELECT DATE(created_at AT TIME ZONE 'UTC') AS exp_date, SUM(amount) 
		          FROM expenses 
		          WHERE period_id = $1 AND user_id IS NULL
		          GROUP BY exp_date`
		rows, err = s.db.Query(query, periodID)
	}
	if err != nil {
		return stats, fmt.Errorf("tagesausgaben abfragen: %w", err)
	}
	defer rows.Close()

	spentMap := make(map[string]float64)
	for rows.Next() {
		var dateStr string
		var total float64
		if err := rows.Scan(&dateStr, &total); err != nil {
			continue
		}
		if len(dateStr) >= 10 {
			dateStr = dateStr[:10]
		}
		spentMap[dateStr] = total
	}

	for i := 0; i < upToDay; i++ {
		if val, exists := spentMap[stats[i].Date]; exists {
			stats[i].Spent = val
		}
	}

	return stats, nil
}

func (s *postgresStore) GetAllExpenses(userID, periodID string) ([]Expense, error) {
	var rows *sql.Rows
	var err error
	if userID != "" {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id = $2 ORDER BY created_at ASC",
			periodID, userID,
		)
	} else {
		rows, err = s.db.Query(
			"SELECT id, period_id, amount, note, created_at FROM expenses WHERE period_id = $1 AND user_id IS NULL ORDER BY created_at ASC",
			periodID,
		)
	}
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
		e.UserID = userID
		expenses = append(expenses, e)
	}

	if expenses == nil {
		expenses = []Expense{}
	}
	return expenses, nil
}

func (s *postgresStore) GetAllPeriods(userID string) ([]PeriodSummary, error) {
	var rows *sql.Rows
	var err error
	query := `SELECT p.id, p.start_date, p.month_days, p.base_budget, p.monthly_total,
	                 COALESCE(SUM(e.amount), 0) as total_spent,
	                 COUNT(e.id) as expense_count
	          FROM periods p
	          LEFT JOIN expenses e ON p.id = e.period_id
	          WHERE (p.user_id = $1 OR ($1 = '' AND p.user_id IS NULL))
	          GROUP BY p.id, p.start_date, p.month_days, p.base_budget, p.monthly_total
	          ORDER BY p.start_date DESC`

	rows, err = s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("perioden abfragen: %w", err)
	}
	defer rows.Close()

	var summaries []PeriodSummary
	for rows.Next() {
		var p PeriodSummary
		if err := rows.Scan(&p.ID, &p.StartDate, &p.MonthDays, &p.BaseBudget, &p.MonthlyTotal, &p.TotalSpent, &p.ExpenseCount); err != nil {
			return nil, fmt.Errorf("periode scannen: %w", err)
		}
		p.UserID = userID
		p.Savings = mathRound(p.MonthlyTotal-p.TotalSpent, 2)
		summaries = append(summaries, p)
	}

	if summaries == nil {
		summaries = []PeriodSummary{}
	}
	return summaries, nil
}

func (s *postgresStore) ImportExpenses(userID, periodID string, expenses []Expense) (int, error) {
	if len(expenses) == 0 {
		return 0, nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("transaction start: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO expenses (period_id, user_id, amount, note, created_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	var count int
	for _, e := range expenses {
		createdAt := e.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		var uid interface{}
		if userID != "" {
			uid = userID
		} else {
			uid = nil
		}
		if _, err := stmt.Exec(periodID, uid, e.Amount, e.Note, createdAt); err != nil {
			return count, fmt.Errorf("ausgabe importieren: %w", err)
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("transaction commit: %w", err)
	}

	return count, nil
}

func (s *postgresStore) DeleteExpense(userID, expenseID string) error {
	var res sql.Result
	var err error
	if userID != "" {
		res, err = s.db.Exec("DELETE FROM expenses WHERE id = $1 AND user_id = $2", expenseID, userID)
	} else {
		res, err = s.db.Exec("DELETE FROM expenses WHERE id = $1 AND user_id IS NULL", expenseID)
	}
	if err != nil {
		return fmt.Errorf("ausgabe löschen: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("ausgabe nicht gefunden")
	}
	return nil
}

// Auth & User Management

func (s *postgresStore) CreateMagicLink(email, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.Exec(
		"INSERT INTO magic_links (email, token_hash, expires_at) VALUES ($1, $2, $3)",
		email, tokenHash, expiresAt,
	)
	return err
}

func (s *postgresStore) ValidateAndConsumeMagicLink(tokenHash string) (string, error) {
	now := time.Now()
	var email string
	err := s.db.QueryRow(
		`UPDATE magic_links 
		 SET used_at = $1 
		 WHERE token_hash = $2 AND expires_at > $1 AND used_at IS NULL 
		 RETURNING email`,
		now, tokenHash,
	).Scan(&email)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("ungültiger oder abgelaufener magic link")
	}
	return email, err
}

func (s *postgresStore) GetOrCreateUserByEmail(email string) (*User, bool, error) {
	var u User
	err := s.db.QueryRow(
		`SELECT id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)

	if err == nil {
		_, _ = s.db.Exec("UPDATE users SET last_login_at = $1 WHERE id = $2", time.Now(), u.ID)
		return &u, false, nil
	}

	if err != sql.ErrNoRows {
		return nil, false, err
	}

	// Neuer User erstellen
	now := time.Now()
	err = s.db.QueryRow(
		`INSERT INTO users (email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active) 
		 VALUES ($1, $2, $2, 450.00, 30, 'emerald', TRUE) 
		 RETURNING id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active`,
		email, now,
	).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)

	if err != nil {
		return nil, false, fmt.Errorf("user erstellen: %w", err)
	}

	return &u, true, nil
}

func (s *postgresStore) GetUserByID(userID string) (*User, error) {
	var u User
	err := s.db.QueryRow(
		`SELECT id, email, created_at, last_login_at, default_monthly_budget, default_period_days, theme, is_active 
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.LastLoginAt, &u.DefaultMonthlyBudget, &u.DefaultPeriodDays, &u.Theme, &u.IsActive)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *postgresStore) UpdateUserSettings(userID string, defaultBudget float64, defaultDays int, theme string) error {
	query := `UPDATE users SET 
	          default_monthly_budget = CASE WHEN $2 > 0 THEN $2 ELSE default_monthly_budget END,
	          default_period_days = CASE WHEN $3 > 0 THEN $3 ELSE default_period_days END,
	          theme = CASE WHEN $4 <> '' THEN $4 ELSE theme END
	          WHERE id = $1`
	_, err := s.db.Exec(query, userID, defaultBudget, defaultDays, theme)
	return err
}

func (s *postgresStore) CreateSession(userID, tokenHash, userAgent, ipAddress string, expiresAt time.Time) error {
	_, err := s.db.Exec(
		"INSERT INTO auth_sessions (user_id, token_hash, user_agent, ip_address, expires_at) VALUES ($1, $2, $3, $4, $5)",
		userID, tokenHash, userAgent, ipAddress, expiresAt,
	)
	return err
}

func (s *postgresStore) ValidateSession(tokenHash string) (string, error) {
	now := time.Now()
	var userID string
	err := s.db.QueryRow(
		"SELECT user_id FROM auth_sessions WHERE token_hash = $1 AND expires_at > $2",
		tokenHash, now,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("ungültige session")
	}
	return userID, err
}

func (s *postgresStore) DeleteSession(tokenHash string) error {
	_, err := s.db.Exec("DELETE FROM auth_sessions WHERE token_hash = $1", tokenHash)
	return err
}

func (s *postgresStore) DeleteUser(userID string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}

func (s *postgresStore) MigrateGuestData(targetUserID string, guestExpenses []Expense, guestPeriods []Period) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 1. Perioden migrieren
	for _, p := range guestPeriods {
		newPeriodID := fmt.Sprintf("%s_%s", targetUserID, p.StartDate.Format("2006-01-02"))
		_, _ = tx.Exec(
			`INSERT INTO periods (id, user_id, start_date, month_days, base_budget, monthly_total) 
			 VALUES ($1, $2, $3, $4, $5, $6) 
			 ON CONFLICT (id) DO NOTHING`,
			newPeriodID, targetUserID, p.StartDate, p.MonthDays, p.BaseBudget, p.MonthlyTotal,
		)
	}

	// 2. Ausgaben migrieren
	migratedCount := 0
	for _, e := range guestExpenses {
		periodID := e.PeriodID
		if !pContainsUser(periodID, targetUserID) {
			periodID = fmt.Sprintf("%s_%s", targetUserID, e.CreatedAt.Format("2006-01-02"))
		}
		_, err := tx.Exec(
			"INSERT INTO expenses (period_id, user_id, amount, note, created_at) VALUES ($1, $2, $3, $4, $5)",
			periodID, targetUserID, e.Amount, e.Note, e.CreatedAt,
		)
		if err == nil {
			migratedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return migratedCount, nil
}

func pContainsUser(periodID, userID string) bool {
	return len(periodID) > len(userID) && periodID[:len(userID)] == userID
}

func (s *postgresStore) Ping() error {
	return s.db.Ping()
}
