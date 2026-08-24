package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Period struct {
	ID           string    `json:"id"`
	StartDate    time.Time `json:"startDate"`
	MonthDays    int       `json:"monthDays"`
	BaseBudget   float64   `json:"baseBudget"`
	MonthlyTotal float64   `json:"monthlyTotal"`
}

type Expense struct {
	ID        string    `json:"id"`
	PeriodID  string    `json:"periodId"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedExpenses struct {
	Items      []Expense `json:"items"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}

type BudgetResponse struct {
	Day           int       `json:"day"`
	MonthDays     int       `json:"monthDays"`
	BaseBudget    float64   `json:"baseBudget"`
	CurrentBudget float64   `json:"currentBudget"`
	Savings       float64   `json:"savings"`
	Color         string    `json:"color"`
	PeriodID      string    `json:"periodId"`
	Expenses      []Expense `json:"expenses"`
}

type ExpenseRequest struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

type UpdateBudgetRequest struct {
	MonthlyTotal float64 `json:"monthlyTotal"`
}

type NewPeriodRequest struct {
	MonthlyTotal float64 `json:"monthlyTotal,omitempty"`
	StartDate    string  `json:"startDate,omitempty"`
}

func calcPeriodDays(start time.Time) int {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	nextMonth := startDay.AddDate(0, 1, 0)
	days := int(nextMonth.Sub(startDay).Hours() / 24)
	if days < 28 {
		days = 30
	}
	return days
}

func (p Period) dayOfMonth(now time.Time) int {
	startInLoc := p.StartDate.In(now.Location())
	start := time.Date(startInLoc.Year(), startInLoc.Month(), startInLoc.Day(), 0, 0, 0, 0, now.Location())
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	diff := nowDay.Sub(start)
	day := int(diff.Hours()/24) + 1
	if day < 1 {
		day = 1
	}
	if day > p.MonthDays {
		day = p.MonthDays
	}
	return day
}

func (p Period) calcBudget(totalSpent float64, now time.Time) (currentBudget float64, savings float64, color string) {
	day := p.dayOfMonth(now)
	remainingDays := p.MonthDays - day + 1 // inklusive heute

	// Ersparnis bis heute: Was stand mir bis einschl. heute zu minus was ausgegeben wurde
	allowedUpToToday := p.BaseBudget * float64(day)
	savings = mathRound(allowedUpToToday-totalSpent, 2)

	// Restbudget verteilt auf verbleibende Tage (inkl. heute)
	remainingBudget := p.MonthlyTotal - totalSpent
	if remainingDays <= 1 {
		currentBudget = mathRound(remainingBudget, 2)
	} else {
		currentBudget = mathRound(remainingBudget/float64(remainingDays), 2)
	}

	if currentBudget < 0 {
		currentBudget = 0
	}

	switch {
	case savings > 0 && currentBudget > 0:
		color = "green"
	case savings < 0 || currentBudget == 0:
		color = "red"
	default:
		color = "white"
	}

	return
}

func mathRound(val float64, places int) float64 {
	shift := 1.0
	for i := 0; i < places; i++ {
		shift *= 10
	}
	return float64(int(val*shift+0.5)) / shift
}

func (p Period) JSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (e Expense) JSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (p Period) StringID() string {
	return fmt.Sprintf("%d-%02d", p.StartDate.Year(), p.StartDate.Month())
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
