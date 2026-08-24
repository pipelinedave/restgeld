package main

import "time"

type Store interface {
	GetOrCreatePeriod() (*Period, error)
	CreatePeriod() (*Period, error)
	CreatePeriodWithStart(start time.Time, monthlyTotal float64, days int) (*Period, error)
	UpdateBudget(newTotal float64, days int) error
	GetTotalExpenses(periodID string) (float64, error)
	GetTodayExpenses(periodID string, now time.Time) (float64, error)
	GetRecentExpenses(periodID string, limit int) ([]Expense, error)
	GetExpenses(periodID string, page, limit int) (*PaginatedExpenses, error)
	AddExpense(periodID string, amount float64, note string) (*Expense, error)
	GetDailyExpenses(periodID string, start time.Time, upToDay int) ([]DailyStat, error)
	GetAllExpenses(periodID string) ([]Expense, error)
	ImportExpenses(periodID string, expenses []Expense) (int, error)
	DeleteExpense(expenseID string) error
	Ping() error
}
