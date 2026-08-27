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
	return m.CreatePeriodWithStart(time.Now(), m.period.MonthlyTotal, m.period.MonthDays)
}

func (m *memoryStore) CreatePeriodWithStart(start time.Time, monthlyTotal float64, days int) (*Period, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expenses = nil
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	monthDays := days
	if monthDays <= 0 {
		monthDays = calcPeriodDays(startDay)
	}
	if monthlyTotal <= 0 {
		monthlyTotal = 450
	}
	m.period = &Period{
		ID:           startDay.Format("2006-01-02"),
		StartDate:    startDay,
		MonthDays:    monthDays,
		BaseBudget:   mathRound(monthlyTotal/float64(monthDays), 2),
		MonthlyTotal: monthlyTotal,
	}
	cp := *m.period
	return &cp, nil
}

func (m *memoryStore) UpdateBudget(newTotal float64, days int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if newTotal > 0 {
		m.period.MonthlyTotal = newTotal
	}
	if days > 0 {
		m.period.MonthDays = days
	}
	m.period.BaseBudget = mathRound(m.period.MonthlyTotal/float64(m.period.MonthDays), 2)
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

func (m *memoryStore) GetTodayExpenses(periodID string, now time.Time) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.AddDate(0, 0, 1)
	var total float64
	for _, e := range m.expenses {
		if !e.CreatedAt.Before(startOfToday) && e.CreatedAt.Before(endOfToday) {
			total += e.Amount
		}
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

func (m *memoryStore) GetExpenses(periodID string, page, limit int) (*PaginatedExpenses, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var periodExpenses []Expense
	for _, e := range m.expenses {
		if e.PeriodID == periodID {
			periodExpenses = append(periodExpenses, e)
		}
	}

	total := len(periodExpenses)
	offset := (page - 1) * limit

	var items []Expense
	// Rückwärts sortiert (neueste zuerst)
	if offset < total {
		end := offset + limit
		if end > total {
			end = total
		}
		for i := offset; i < end; i++ {
			items = append(items, periodExpenses[total-1-i])
		}
	}

	if items == nil {
		items = []Expense{}
	}

	totalPages := 1
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &PaginatedExpenses{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *memoryStore) AddExpense(periodID string, amount float64, note string) (*Expense, error) {
	return m.AddExpenseWithDate(periodID, amount, note, time.Now()), nil
}

func (m *memoryStore) AddExpenseWithDate(periodID string, amount float64, note string, createdAt time.Time) *Expense {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	e := Expense{
		ID:        fmt.Sprintf("exp-%d", m.nextID),
		PeriodID:  periodID,
		Amount:    amount,
		Note:      note,
		CreatedAt: createdAt,
	}
	m.expenses = append(m.expenses, e)
	return &e
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

func (m *memoryStore) GetDailyExpenses(periodID string, start time.Time, upToDay int) ([]DailyStat, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if upToDay < 1 {
		upToDay = 1
	}

	startInLoc := start
	startDay := time.Date(startInLoc.Year(), startInLoc.Month(), startInLoc.Day(), 0, 0, 0, 0, start.Location())

	stats := make([]DailyStat, upToDay)
	for d := 1; d <= upToDay; d++ {
		currentDate := startDay.AddDate(0, 0, d-1)
		dateStr := currentDate.Format("2006-01-02")
		nextDate := currentDate.AddDate(0, 0, 1)

		var spent float64
		for _, e := range m.expenses {
			if e.PeriodID == periodID && !e.CreatedAt.Before(currentDate) && e.CreatedAt.Before(nextDate) {
				spent += e.Amount
			}
		}

		stats[d-1] = DailyStat{
			Day:   d,
			Date:  dateStr,
			Spent: mathRound(spent, 2),
		}
	}

	return stats, nil
}

func (m *memoryStore) GetAllExpenses(periodID string) ([]Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []Expense
	for _, e := range m.expenses {
		if e.PeriodID == periodID {
			result = append(result, e)
		}
	}
	if result == nil {
		result = []Expense{}
	}
	return result, nil
}

func (m *memoryStore) ImportExpenses(periodID string, expenses []Expense) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, exp := range expenses {
		if exp.Amount <= 0 {
			continue
		}
		m.nextID++
		id := exp.ID
		if id == "" {
			id = fmt.Sprintf("exp-%d", m.nextID)
		}
		createdAt := exp.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		e := Expense{
			ID:        id,
			PeriodID:  periodID,
			Amount:    exp.Amount,
			Note:      exp.Note,
			CreatedAt: createdAt,
		}
		m.expenses = append(m.expenses, e)
		count++
	}
	return count, nil
}

func (m *memoryStore) GetAllPeriods() ([]PeriodSummary, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var totalSpent float64
	for _, e := range m.expenses {
		if e.PeriodID == m.period.ID {
			totalSpent += e.Amount
		}
	}

	summary := PeriodSummary{
		ID:           m.period.ID,
		StartDate:    m.period.StartDate,
		MonthDays:    m.period.MonthDays,
		BaseBudget:   m.period.BaseBudget,
		MonthlyTotal: m.period.MonthlyTotal,
		TotalSpent:   mathRound(totalSpent, 2),
		Savings:      mathRound(m.period.MonthlyTotal-totalSpent, 2),
		ExpenseCount: len(m.expenses),
	}

	return []PeriodSummary{summary}, nil
}

func (m *memoryStore) Ping() error {
	return nil
}
