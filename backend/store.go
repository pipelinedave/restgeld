package main

import "time"

type Store interface {
	GetOrCreatePeriod() (*Period, error)
	CreatePeriod() (*Period, error)
	CreatePeriodWithStart(start time.Time, monthlyTotal float64) (*Period, error)
	UpdateBudget(newTotal float64) error
	GetTotalExpenses(periodID string) (float64, error)
	GetRecentExpenses(periodID string, limit int) ([]Expense, error)
	AddExpense(periodID string, amount float64, note string) (*Expense, error)
	DeleteExpense(expenseID string) error
	Ping() error
}
