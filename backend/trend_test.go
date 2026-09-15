package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildMonthlyTrend(t *testing.T) {
	loc := time.UTC
	periods := []PeriodSummary{
		{
			ID:           "2026-07-01",
			StartDate:    time.Date(2026, 7, 1, 0, 0, 0, 0, loc),
			MonthDays:    31,
			MonthlyTotal: 450,
			TotalSpent:   300,
			ExpenseCount: 20,
		},
		{
			ID:           "2026-08-01",
			StartDate:    time.Date(2026, 8, 1, 0, 0, 0, 0, loc),
			MonthDays:    31,
			MonthlyTotal: 450,
			TotalSpent:   500,
			ExpenseCount: 31,
		},
		{
			ID:           "2026-09-01",
			StartDate:    time.Date(2026, 9, 1, 0, 0, 0, 0, loc),
			MonthDays:    30,
			MonthlyTotal: 500,
			TotalSpent:   250,
			ExpenseCount: 15,
		},
	}

	trend := buildMonthlyTrend(periods, loc)

	if len(trend.Months) != 3 {
		t.Fatalf("erwartet 3 Monate, bekommen %d", len(trend.Months))
	}

	// Aufsteigend sortiert
	if trend.Months[0].Month != "2026-07" || trend.Months[1].Month != "2026-08" || trend.Months[2].Month != "2026-09" {
		t.Errorf("sortierung falsch: %s, %s, %s", trend.Months[0].Month, trend.Months[1].Month, trend.Months[2].Month)
	}

	// Savings pro Monat
	if trend.Months[0].Savings != 150 {
		t.Errorf("juli savings erwartet 150, bekommen %.2f", trend.Months[0].Savings)
	}
	if trend.Months[1].Savings != -50 {
		t.Errorf("august savings erwartet -50, bekommen %.2f", trend.Months[1].Savings)
	}

	// AvgDailySpend
	if trend.Months[2].AvgDailySpend != 8.33 {
		t.Errorf("september avg erwartet 8.33, bekommen %.2f", trend.Months[2].AvgDailySpend)
	}

	// Summary
	s := trend.Summary
	if s.MonthCount != 3 {
		t.Errorf("monthCount erwartet 3, bekommen %d", s.MonthCount)
	}
	if s.TotalSaved != 350 {
		t.Errorf("totalSaved erwartet 350, bekommen %.2f", s.TotalSaved)
	}
	if s.TotalBudget != 1400 {
		t.Errorf("totalBudget erwartet 1400, bekommen %.2f", s.TotalBudget)
	}
	if s.TotalSpent != 1050 {
		t.Errorf("totalSpent erwartet 1050, bekommen %.2f", s.TotalSpent)
	}
	// 350/3 = 116.67
	if s.AvgSavings != 116.67 {
		t.Errorf("avgSavings erwartet 116.67, bekommen %.2f", s.AvgSavings)
	}
	if s.BestSavingsMonth != "2026-09" {
		t.Errorf("bestMonth erwartet 2026-09, bekommen %s", s.BestSavingsMonth)
	}
}

func TestBuildMonthlyTrendEmpty(t *testing.T) {
	trend := buildMonthlyTrend(nil, time.UTC)
	if trend.Months == nil {
		t.Error("months sollte leeres Array sein, nicht nil")
	}
	if len(trend.Months) != 0 {
		t.Errorf("erwartet 0 Monate, bekommen %d", len(trend.Months))
	}
	if trend.Summary.MonthCount != 0 || trend.Summary.BestSavingsMonth != "" {
		t.Errorf("summary leer erwartet, bekommen %+v", trend.Summary)
	}
}

func TestBuildMonthlyTrendAggregatesSameMonth(t *testing.T) {
	loc := time.UTC
	periods := []PeriodSummary{
		{ID: "a", StartDate: time.Date(2026, 8, 10, 0, 0, 0, 0, loc), MonthDays: 31, MonthlyTotal: 450, TotalSpent: 200, ExpenseCount: 5},
		{ID: "b", StartDate: time.Date(2026, 8, 20, 0, 0, 0, 0, loc), MonthDays: 31, MonthlyTotal: 500, TotalSpent: 400, ExpenseCount: 9},
	}

	trend := buildMonthlyTrend(periods, loc)
	if len(trend.Months) != 1 {
		t.Fatalf("erwartet 1 Monat, bekommen %d", len(trend.Months))
	}
	// Letzte Periode im Monat gewinnt
	if trend.Months[0].MonthlyTotal != 500 || trend.Months[0].TotalSpent != 400 {
		t.Errorf("aggregation falsch: %+v", trend.Months[0])
	}
}

func TestTrendEndpoint(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	store.AddExpense("", "2026-08-25", 50.00, "Alt")
	srv := &server{store: store, now: func() time.Time { return now }}

	req := httptest.NewRequest(http.MethodGet, "/api/trend", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", rec.Code)
	}

	var resp TrendResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decode: %v", err)
	}

	if len(resp.Months) != 1 {
		t.Fatalf("erwartet 1 Monat, bekommen %d", len(resp.Months))
	}
	if resp.Months[0].Month != "2026-08" {
		t.Errorf("erwartet monat 2026-08, bekommen %s", resp.Months[0].Month)
	}
	// Default-Periode: monthly 450, totalSpent 50
	if resp.Months[0].TotalSpent != 50 {
		t.Errorf("erwartet totalSpent 50, bekommen %.2f", resp.Months[0].TotalSpent)
	}
	if resp.Summary.BestSavingsMonth != "2026-08" {
		t.Errorf("erwartet bestMonth 2026-08, bekommen %s", resp.Summary.BestSavingsMonth)
	}
}

func TestTrendEndpointMethodNotAllowed(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store, now: time.Now}

	req := httptest.NewRequest(http.MethodPost, "/api/trend", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("erwartet 405, bekommen %d", rec.Code)
	}
}
