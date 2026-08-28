package main

import (
	"testing"
	"time"
)

// periodOnDay erzeugt eine Periode, deren StartDate so gesetzt ist, dass `now`
// auf einen bestimmten Tag (1-basiert) innerhalb der Periode faellt.
func periodOnDay(t *testing.T, day, monthDays int, monthlyTotal, baseBudget float64, now time.Time) Period {
	t.Helper()
	start := now.AddDate(0, 0, -(day - 1))
	return Period{
		ID:           "test-period",
		MonthlyTotal: monthlyTotal,
		BaseBudget:   baseBudget,
		MonthDays:    monthDays,
		StartDate:    start,
	}
}

func TestCalcBudget_HeutigesTagesbudgetBleibtKonstantBeiUeberzug(t *testing.T) {
	// Reproduktion des gemeldeten Bugs: Beim Ueberschreiten des Tagesbudgets
	// soll das "Start-Tagesbudget heute" (todayBase) konstant bleiben, auch wenn
	// currentBudget auf 0 geklemmt wird. Vorher sprang das UI von "8,78 / 10,41"
	// auf "0,00 / 16,63" - weil das Frontend todayBase als
	// currentBudget + spentToday rekonstruierte.
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	day := 30
	baseBudget := 10.41
	// Monatsbudget so, dass bei Tag 30 mit exakt auf Kurs liegenden Vortagen
	// todayBase = 10.41 ist.
	monthlyTotal := baseBudget*float64(day-1) + baseBudget // = 312.30

	p := periodOnDay(t, day, 30, monthlyTotal, baseBudget, now)

	// Vortage exakt auf Kurs: totalSpent_vorVortag = BaseBudget*(day-1)
	spentPrior := baseBudget * float64(day-1)
	// Vor der Ausgabe: heute bereits 1,63 ausgegeben -> "8,78 / 10,41"
	beforeSpent := 1.63

	curBefore, savBefore, _, baseBefore := p.calcBudget(spentPrior+beforeSpent, beforeSpent, now)
	wantBase := 10.41
	if baseBefore != wantBase {
		t.Fatalf("todayBase vor Ausgabe = %.2f, erwartet %.2f", baseBefore, wantBase)
	}
	if curBefore != 8.78 {
		t.Errorf("currentBudget vor Ausgabe = %.2f, erwartet 8.78", curBefore)
	}
	if savBefore != 0 {
		t.Errorf("savings vor Ausgabe = %.2f, erwartet 0", savBefore)
	}

	// Nach einer 15-Euro-Ausgabe: heute 16,63 ausgegeben.
	afterSpent := beforeSpent + 15.00
	curAfter, savAfter, _, baseAfter := p.calcBudget(spentPrior+afterSpent, afterSpent, now)

	// KERN: todayBase muss unveraendert bleiben.
	if baseAfter != wantBase {
		t.Fatalf("todayBase nach Ueberzug = %.2f, erwartet unveraendert %.2f (Bug: UI zeigte 16,63)", baseAfter, wantBase)
	}
	// currentBudget wird auf 0 geklemmt (Anzeige "0,00").
	if curAfter != 0 {
		t.Errorf("currentBudget nach Ueberzug = %.2f, erwartet 0", curAfter)
	}
	// Ueberzug 15 - 0 + ... = todaySpent(16,63) - base(10,41) = 6,22 auf savings.
	wantSav := mathRound(wantBase-afterSpent, 2)
	if savAfter != wantSav {
		t.Errorf("savings nach Ueberzug = %.2f, erwartet %.2f", savAfter, wantSav)
	}
}

func TestCalcBudget_TagesbudgetRegular(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	day := 15
	baseBudget := 10.00
	monthlyTotal := 300.00
	p := periodOnDay(t, day, 30, monthlyTotal, baseBudget, now)

	// Vortage exakt auf Kurs, heute noch nichts ausgegeben.
	total := baseBudget * float64(day-1)
	cur, sav, color, base := p.calcBudget(total, 0, now)

	if base != 10.00 {
		t.Errorf("todayBase = %.2f, erwartet 10.00", base)
	}
	if cur != 10.00 {
		t.Errorf("currentBudget = %.2f, erwartet 10.00", cur)
	}
	if sav != 0 {
		t.Errorf("savings = %.2f, erwartet 0", sav)
	}
	if color != "white" {
		t.Errorf("color = %s, erwartet white", color)
	}
}

func TestCalcBudget_TagesbudgetBeiVortagsersparnis(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	day := 3
	baseBudget := 10.00
	monthlyTotal := 300.00
	p := periodOnDay(t, day, 30, monthlyTotal, baseBudget, now)

	// An den ersten beiden Tagen je 2 statt 10 ausgegeben -> 16 Ersparnis.
	spentPrior := 4.0
	total := spentPrior // heuta noch nichts
	cur, sav, color, base := p.calcBudget(total, 0, now)

	// todayBase = (300-4)/(30-3+1) = 296/28 = 10.57
	wantBase := mathRound(296.0/28.0, 2)
	if base != wantBase {
		t.Errorf("todayBase = %.2f, erwartet %.2f", base, wantBase)
	}
	if cur != wantBase {
		t.Errorf("currentBudget = %.2f, erwartet %.2f", cur, wantBase)
	}
	if sav != 16 {
		t.Errorf("savings = %.2f, erwartet 16", sav)
	}
	if color != "green" {
		t.Errorf("color = %s, erwartet green", color)
	}
}
