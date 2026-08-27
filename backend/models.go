package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type Period struct {
	ID           string    `json:"id"`
	StartDate    time.Time `json:"startDate"`
	MonthDays    int       `json:"monthDays"`
	BaseBudget   float64   `json:"baseBudget"`
	MonthlyTotal float64   `json:"monthlyTotal"`
}

type PeriodSummary struct {
	ID           string    `json:"id"`
	StartDate    time.Time `json:"startDate"`
	MonthDays    int       `json:"monthDays"`
	BaseBudget   float64   `json:"baseBudget"`
	MonthlyTotal float64   `json:"monthlyTotal"`
	TotalSpent   float64   `json:"totalSpent"`
	Savings      float64   `json:"savings"`
	ExpenseCount int       `json:"expenseCount"`
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

type DailyStat struct {
	Day   int     `json:"day"`
	Date  string  `json:"date"`
	Spent float64 `json:"spent"`
}

type StreakInfo struct {
	CurrentStreak   int `json:"currentStreak"`
	LongestStreak   int `json:"longestStreak"`
	NoSpendDays     int `json:"noSpendDays"`
	UnderBudgetDays int `json:"underBudgetDays"`
}

type ProjectionInfo struct {
	ProjectedSavings    float64 `json:"projectedSavings"`
	ProjectedTotalSpent float64 `json:"projectedTotalSpent"`
	AvgDailySpend       float64 `json:"avgDailySpend"`
	Status              string  `json:"status"` // "saving" | "deficit"
}

type BudgetResponse struct {
	Day           int            `json:"day"`
	MonthDays     int            `json:"monthDays"`
	BaseBudget    float64        `json:"baseBudget"`
	CurrentBudget float64        `json:"currentBudget"`
	Savings       float64        `json:"savings"`
	Color         string         `json:"color"`
	PeriodID      string         `json:"periodId"`
	Expenses      []Expense      `json:"expenses"`
	DailyStats    []DailyStat    `json:"dailyStats"`
	Streak        StreakInfo     `json:"streak"`
	Projection    ProjectionInfo `json:"projection"`
}

func calcProjection(totalSpent float64, day, monthDays int, monthlyTotal float64) ProjectionInfo {
	if day <= 0 || monthDays <= 0 {
		return ProjectionInfo{Status: "saving"}
	}

	avgDailySpend := totalSpent / float64(day)
	remainingDays := monthDays - day
	if remainingDays < 0 {
		remainingDays = 0
	}

	projectedTotalSpent := totalSpent + (avgDailySpend * float64(remainingDays))
	projectedSavings := monthlyTotal - projectedTotalSpent

	status := "saving"
	if projectedSavings < 0 {
		status = "deficit"
	}

	return ProjectionInfo{
		ProjectedSavings:    math.Round(projectedSavings*100) / 100,
		ProjectedTotalSpent: math.Round(projectedTotalSpent*100) / 100,
		AvgDailySpend:       math.Round(avgDailySpend*100) / 100,
		Status:              status,
	}
}

func calcStreakInfo(stats []DailyStat, baseBudget float64, currentDay int) StreakInfo {
	if len(stats) == 0 {
		return StreakInfo{}
	}

	currentStreak := 0
	longestStreak := 0
	noSpendDays := 0
	underBudgetDays := 0
	runningStreak := 0

	for _, s := range stats {
		// Nur abgeschlossene Vortage zählen als Spar-Tage und Null-Ausgaben-Tage
		if s.Day < currentDay {
			if s.Spent == 0 {
				noSpendDays++
			}
			if s.Spent <= baseBudget {
				underBudgetDays++
				runningStreak++
				if runningStreak > longestStreak {
					longestStreak = runningStreak
				}
			} else {
				runningStreak = 0
			}
		} else if s.Day == currentDay {
			// Aktiver Tag heute: Wenn heute bereits das Basisbudget überzogen wurde, bricht der Streak ab
			if s.Spent > baseBudget {
				runningStreak = 0
			}
		}
	}

	currentStreak = runningStreak

	return StreakInfo{
		CurrentStreak:   currentStreak,
		LongestStreak:   longestStreak,
		NoSpendDays:     noSpendDays,
		UnderBudgetDays: underBudgetDays,
	}
}

type ExpenseRequest struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

type UpdateBudgetRequest struct {
	MonthlyTotal float64 `json:"monthlyTotal,omitempty"`
	Days         int     `json:"days,omitempty"`
}

type NewPeriodRequest struct {
	MonthlyTotal float64 `json:"monthlyTotal,omitempty"`
	StartDate    string  `json:"startDate,omitempty"`
	Days         int     `json:"days,omitempty"`
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

func (p Period) calcBudget(totalSpent, todaySpent float64, now time.Time) (currentBudget float64, savings float64, color string) {
	day := p.dayOfMonth(now)
	remainingDays := p.MonthDays - day + 1 // inklusive heute

	// Ausgaben vor dem heutigen Tag
	spentPriorToToday := totalSpent - todaySpent
	remainingPrior := p.MonthlyTotal - spentPriorToToday

	// Ersparnis aus abgeschlossenen Vortagen (Tag 1 bis day-1)
	var priorSavings float64
	if day > 1 {
		allowedPrior := p.BaseBudget * float64(day-1)
		priorSavings = allowedPrior - spentPriorToToday
	} else {
		priorSavings = 0.0
	}

	// Tagesbudget für den heutigen Tag vor heutigen Ausgaben
	var startOfTodayDaily float64
	if remainingDays <= 1 {
		startOfTodayDaily = mathRound(remainingPrior, 2)
	} else {
		startOfTodayDaily = mathRound(remainingPrior/float64(remainingDays), 2)
	}
	if startOfTodayDaily < 0 {
		startOfTodayDaily = 0
	}

	// Heutiges Restbudget nach Abzug der heutigen Ausgaben
	todayRemaining := mathRound(startOfTodayDaily-todaySpent, 2)
	if todayRemaining < 0 {
		currentBudget = 0
		overdraw := todaySpent - startOfTodayDaily
		savings = mathRound(priorSavings-overdraw, 2)
	} else {
		currentBudget = todayRemaining
		savings = mathRound(priorSavings, 2)
	}

	switch {
	case savings > 0:
		color = "green"
	case savings < 0 || (currentBudget == 0 && todaySpent > 0):
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
