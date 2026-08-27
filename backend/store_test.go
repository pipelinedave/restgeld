package main

import (
	"fmt"
	"sync"
	"time"
)

type memoryStore struct {
	mu           sync.Mutex
	periods      map[string]*Period      // key: userID
	expenses     map[string][]Expense    // key: userID
	users        map[string]*User        // key: userID
	usersByEmail map[string]*User        // key: email
	magicLinks   map[string]*MagicLink   // key: tokenHash
	sessions     map[string]*AuthSession // key: tokenHash
	nextID       int
}

func newMemoryStore() *memoryStore {
	ms := &memoryStore{
		periods:      make(map[string]*Period),
		expenses:     make(map[string][]Expense),
		users:        make(map[string]*User),
		usersByEmail: make(map[string]*User),
		magicLinks:   make(map[string]*MagicLink),
		sessions:     make(map[string]*AuthSession),
	}

	// Default guest period
	ms.periods[""] = &Period{
		ID:           "2026-08",
		StartDate:    time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		MonthDays:    31,
		BaseBudget:   14.52,
		MonthlyTotal: 450,
	}
	ms.expenses[""] = []Expense{}

	return ms
}

func (m *memoryStore) GetOrCreatePeriod(userID string) (*Period, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.periods[userID]
	if !exists || p == nil {
		monthlyTotal := 450.0
		monthDays := 31
		if u, ok := m.users[userID]; ok && u != nil {
			if u.DefaultMonthlyBudget > 0 {
				monthlyTotal = u.DefaultMonthlyBudget
			}
			if u.DefaultPeriodDays > 0 {
				monthDays = u.DefaultPeriodDays
			}
		}

		startDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		periodID := startDay.Format("2006-01-02")
		if userID != "" {
			periodID = fmt.Sprintf("%s_%s", userID, periodID)
		}

		p = &Period{
			ID:           periodID,
			UserID:       userID,
			StartDate:    startDay,
			MonthDays:    monthDays,
			BaseBudget:   mathRound(monthlyTotal/float64(monthDays), 2),
			MonthlyTotal: monthlyTotal,
		}
		m.periods[userID] = p
	}

	cp := *p
	return &cp, nil
}

func (m *memoryStore) CreatePeriod(userID string) (*Period, error) {
	p, _ := m.GetOrCreatePeriod(userID)
	return m.CreatePeriodWithStart(userID, time.Now(), p.MonthlyTotal, p.MonthDays)
}

func (m *memoryStore) CreatePeriodWithStart(userID string, start time.Time, monthlyTotal float64, days int) (*Period, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.expenses[userID] = nil
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	monthDays := days
	if monthDays <= 0 {
		monthDays = calcPeriodDays(startDay)
	}
	if monthlyTotal <= 0 {
		monthlyTotal = 450
	}
	periodID := startDay.Format("2006-01-02")
	if userID != "" {
		periodID = fmt.Sprintf("%s_%s", userID, periodID)
	}

	p := &Period{
		ID:           periodID,
		UserID:       userID,
		StartDate:    startDay,
		MonthDays:    monthDays,
		BaseBudget:   mathRound(monthlyTotal/float64(monthDays), 2),
		MonthlyTotal: monthlyTotal,
	}
	m.periods[userID] = p

	cp := *p
	return &cp, nil
}

func (m *memoryStore) UpdateBudget(userID string, newTotal float64, days int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.periods[userID]
	if !exists || p == nil {
		p = &Period{
			ID:           "2026-08",
			UserID:       userID,
			StartDate:    time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			MonthDays:    31,
			BaseBudget:   14.52,
			MonthlyTotal: 450,
		}
		m.periods[userID] = p
	}

	if newTotal > 0 {
		p.MonthlyTotal = newTotal
	}
	if days > 0 {
		p.MonthDays = days
	}
	p.BaseBudget = mathRound(p.MonthlyTotal/float64(p.MonthDays), 2)
	return nil
}

func (m *memoryStore) GetTotalExpenses(userID, periodID string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var total float64
	for _, e := range m.expenses[userID] {
		total += e.Amount
	}
	return total, nil
}

func (m *memoryStore) GetTodayExpenses(userID, periodID string, now time.Time) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := startOfToday.AddDate(0, 0, 1)
	var total float64
	for _, e := range m.expenses[userID] {
		if !e.CreatedAt.Before(startOfToday) && e.CreatedAt.Before(endOfToday) {
			total += e.Amount
		}
	}
	return total, nil
}

func (m *memoryStore) GetRecentExpenses(userID, periodID string, limit int) ([]Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	exps := m.expenses[userID]
	n := len(exps)
	if n == 0 {
		return []Expense{}, nil
	}
	start := n - limit
	if start < 0 {
		start = 0
	}
	result := make([]Expense, n-start)
	for i, j := n-1, 0; i >= start; i, j = i-1, j+1 {
		result[j] = exps[i]
	}
	return result, nil
}

func (m *memoryStore) GetExpenses(userID, periodID string, page, limit int) (*PaginatedExpenses, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	allExps := m.expenses[userID]
	var exps []Expense
	for _, e := range allExps {
		if periodID == "" || e.PeriodID == periodID {
			exps = append(exps, e)
		}
	}

	total := len(exps)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	reversed := make([]Expense, total)
	for i, e := range exps {
		reversed[total-1-i] = e
	}

	offset := (page - 1) * limit
	end := offset + limit
	if offset > total {
		offset = total
	}
	if end > total {
		end = total
	}

	return &PaginatedExpenses{
		Items:      reversed[offset:end],
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *memoryStore) GetDayExpenses(userID, periodID string, day time.Time) ([]Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.AddDate(0, 0, 1)

	var items []Expense
	for _, e := range m.expenses[userID] {
		if periodID != "" && e.PeriodID != periodID {
			continue
		}
		if e.CreatedAt.After(start) && e.CreatedAt.Before(end) {
			items = append(items, e)
		}
	}
	return items, nil
}

func (m *memoryStore) AddExpense(userID, periodID string, amount float64, note string) (*Expense, error) {
	return m.AddExpenseWithDate(userID, periodID, amount, note, time.Now())
}

func (m *memoryStore) AddExpenseWithDate(userID, periodID string, amount float64, note string, createdAt time.Time) (*Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextID++
	e := Expense{
		ID:        fmt.Sprintf("e-%d", m.nextID),
		UserID:    userID,
		PeriodID:  periodID,
		Amount:    amount,
		Note:      note,
		CreatedAt: createdAt,
	}
	m.expenses[userID] = append(m.expenses[userID], e)
	return &e, nil
}

func (m *memoryStore) GetDailyExpenses(userID, periodID string, start time.Time, upToDay int) ([]DailyStat, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

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

	for _, e := range m.expenses[userID] {
		expDay := time.Date(e.CreatedAt.Year(), e.CreatedAt.Month(), e.CreatedAt.Day(), 0, 0, 0, 0, e.CreatedAt.Location())
		diffDays := int(expDay.Sub(startDay).Hours() / 24)
		if diffDays >= 0 && diffDays < upToDay {
			stats[diffDays].Spent += e.Amount
		}
	}

	return stats, nil
}

func (m *memoryStore) GetAllExpenses(userID, periodID string) ([]Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	exps := m.expenses[userID]
	res := make([]Expense, len(exps))
	copy(res, exps)
	return res, nil
}

func (m *memoryStore) GetAllPeriods(userID string, now time.Time) ([]PeriodSummary, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.periods[userID]
	if !exists || p == nil {
		return []PeriodSummary{}, nil
	}

	var totalSpent float64
	for _, e := range m.expenses[userID] {
		totalSpent += e.Amount
	}

	summary := PeriodSummary{
		ID:           p.ID,
		UserID:       userID,
		StartDate:    p.StartDate,
		EndDate:      now,
		MonthDays:    p.MonthDays,
		ActualDays:   calcActualDays(p.StartDate, now),
		BaseBudget:   p.BaseBudget,
		MonthlyTotal: p.MonthlyTotal,
		TotalSpent:   totalSpent,
		Savings:      mathRound(p.MonthlyTotal-totalSpent, 2),
		ExpenseCount: len(m.expenses[userID]),
	}
	return []PeriodSummary{summary}, nil
}

func (m *memoryStore) ImportExpenses(userID, periodID string, expenses []Expense) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, e := range expenses {
		m.nextID++
		e.ID = fmt.Sprintf("e-%d", m.nextID)
		e.UserID = userID
		e.PeriodID = periodID
		m.expenses[userID] = append(m.expenses[userID], e)
	}
	return len(expenses), nil
}

func (m *memoryStore) DeleteExpense(userID, expenseID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	exps := m.expenses[userID]
	for i, e := range exps {
		if e.ID == expenseID {
			m.expenses[userID] = append(exps[:i], exps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("ausgabe nicht gefunden")
}

// Auth & User Methods for memoryStore

func (m *memoryStore) CreateMagicLink(email, tokenHash string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.magicLinks[tokenHash] = &MagicLink{
		ID:        fmt.Sprintf("ml-%d", len(m.magicLinks)+1),
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return nil
}

func (m *memoryStore) ValidateAndConsumeMagicLink(tokenHash string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ml, exists := m.magicLinks[tokenHash]
	if !exists || ml == nil || ml.UsedAt != nil || ml.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("ungültiger oder abgelaufener link")
	}

	now := time.Now()
	ml.UsedAt = &now
	return ml.Email, nil
}

func (m *memoryStore) GetOrCreateUserByEmail(email string) (*User, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, exists := m.usersByEmail[email]; exists && u != nil {
		u.LastLoginAt = time.Now()
		return u, false, nil
	}

	now := time.Now()
	u := &User{
		ID:                   fmt.Sprintf("usr-%d", len(m.users)+1),
		Email:                email,
		CreatedAt:            now,
		LastLoginAt:          now,
		DefaultMonthlyBudget: 450.0,
		DefaultPeriodDays:    30,
		Theme:                "emerald",
		IsActive:             true,
	}

	m.users[u.ID] = u
	m.usersByEmail[email] = u
	return u, true, nil
}

func (m *memoryStore) GetUserByID(userID string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, exists := m.users[userID]
	if !exists || u == nil {
		return nil, fmt.Errorf("user nicht gefunden")
	}
	return u, nil
}

func (m *memoryStore) UpdateUserSettings(userID string, defaultBudget float64, defaultDays int, theme string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, exists := m.users[userID]
	if !exists || u == nil {
		return fmt.Errorf("user nicht gefunden")
	}
	if defaultBudget > 0 {
		u.DefaultMonthlyBudget = defaultBudget
	}
	if defaultDays > 0 {
		u.DefaultPeriodDays = defaultDays
	}
	if theme != "" {
		u.Theme = theme
	}
	return nil
}

func (m *memoryStore) CreateSession(userID, tokenHash, userAgent, ipAddress string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[tokenHash] = &AuthSession{
		ID:        fmt.Sprintf("sess-%d", len(m.sessions)+1),
		UserID:    userID,
		TokenHash: tokenHash,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return nil
}

func (m *memoryStore) ValidateSession(tokenHash string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, exists := m.sessions[tokenHash]
	if !exists || s == nil || s.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("ungültige session")
	}
	return s.UserID, nil
}

func (m *memoryStore) DeleteSession(tokenHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, tokenHash)
	return nil
}

func (m *memoryStore) DeleteUser(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, exists := m.users[userID]; exists && u != nil {
		delete(m.usersByEmail, u.Email)
		delete(m.users, userID)
	}
	delete(m.periods, userID)
	delete(m.expenses, userID)
	for hash, s := range m.sessions {
		if s.UserID == userID {
			delete(m.sessions, hash)
		}
	}
	return nil
}

func (m *memoryStore) MigrateGuestData(targetUserID string, guestExpenses []Expense, guestPeriods []Period) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range guestPeriods {
		p.UserID = targetUserID
		m.periods[targetUserID] = &p
	}

	count := 0
	for _, e := range guestExpenses {
		m.nextID++
		e.ID = fmt.Sprintf("e-%d", m.nextID)
		e.UserID = targetUserID
		m.expenses[targetUserID] = append(m.expenses[targetUserID], e)
		count++
	}
	return count, nil
}

func (m *memoryStore) Ping() error {
	return nil
}
