package main

import "time"

type Store interface {
	// Period & Budget (Scoped to user, or "" for guest/default)
	GetOrCreatePeriod(userID string) (*Period, error)
	CreatePeriod(userID string) (*Period, error)
	CreatePeriodWithStart(userID string, start time.Time, monthlyTotal float64, days int) (*Period, error)
	UpdateBudget(userID string, newTotal float64, days int) error
	GetTotalExpenses(userID, periodID string) (float64, error)
	GetTodayExpenses(userID, periodID string, now time.Time) (float64, error)
	GetRecentExpenses(userID, periodID string, limit int) ([]Expense, error)
	GetExpenses(userID, periodID string, page, limit int) (*PaginatedExpenses, error)
	AddExpense(userID, periodID string, amount float64, note string) (*Expense, error)
	AddExpenseWithDate(userID, periodID string, amount float64, note string, createdAt time.Time) (*Expense, error)
	GetDailyExpenses(userID, periodID string, start time.Time, upToDay int) ([]DailyStat, error)
	GetAllExpenses(userID, periodID string) ([]Expense, error)
	GetAllPeriods(userID string) ([]PeriodSummary, error)
	ImportExpenses(userID, periodID string, expenses []Expense) (int, error)
	DeleteExpense(userID, expenseID string) error

	// Auth & User Management
	CreateMagicLink(email, tokenHash string, expiresAt time.Time) error
	ValidateAndConsumeMagicLink(tokenHash string) (string, error) // Returns email
	GetOrCreateUserByEmail(email string) (*User, bool, error)    // Returns user, isNew, error
	GetUserByID(userID string) (*User, error)
	UpdateUserSettings(userID string, defaultBudget float64, defaultDays int, theme string) error
	CreateSession(userID, tokenHash, userAgent, ipAddress string, expiresAt time.Time) error
	ValidateSession(tokenHash string) (string, error) // Returns userID
	DeleteSession(tokenHash string) error
	DeleteUser(userID string) error
	MigrateGuestData(targetUserID string, guestExpenses []Expense, guestPeriods []Period) (int, error)

	Ping() error
}
