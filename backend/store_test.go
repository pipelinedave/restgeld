package main

import (
	"fmt"
	"sync"
	"time"
)

type memoryStore struct {
	mu       sync.Mutex
	period   *Period
	expenses []Expense
	nextID   int
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		period: &Period{
			ID:           "2026-08",
			StartDate:    time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			MonthDays:    31,
			BaseBudget:   14.52,
			MonthlyTotal: 450,
		},
		expenses: []Expense{},
	}
}

func (m *memoryStore) GetOrCreatePeriod() (*Period, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *m.period
	return &cp, nil
}

func (m *memoryStore) CreatePeriod() (*Period, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expenses = nil
	cp := *m.period
	return &cp, nil
}

func (m *memoryStore) UpdateBudget(newTotal float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.period.MonthlyTotal = newTotal
	m.period.BaseBudget = mathRound(newTotal/float64(m.period.MonthDays), 2)
	return nil
}

func (m *memoryStore) GetTotalExpenses(periodID string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var total float64
	for _, e := range m.expenses {
		total += e.Amount
	}
	return total, nil
}

func (m *memoryStore) GetRecentExpenses(periodID string, limit int) ([]Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.expenses)
	if limit > n {
		limit = n
	}
	result := make([]Expense, limit)
	for i := 0; i < limit; i++ {
		result[i] = m.expenses[n-1-i]
	}
	return result, nil
}

func (m *memoryStore) AddExpense(periodID string, amount float64, note string) (*Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	e := Expense{
		ID:        fmt.Sprintf("exp-%d", m.nextID),
		PeriodID:  periodID,
		Amount:    amount,
		Note:      note,
		CreatedAt: time.Now(),
	}
	m.expenses = append(m.expenses, e)
	return &e, nil
}

func (m *memoryStore) DeleteExpense(expenseID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, e := range m.expenses {
		if e.ID == expenseID {
			m.expenses = append(m.expenses[:i], m.expenses[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("ausgabe nicht gefunden")
}

func (m *memoryStore) Ping() error {
	return nil
}
