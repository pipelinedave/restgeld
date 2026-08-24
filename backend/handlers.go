package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type server struct {
	store Store
	now   func() time.Time
}

func newServer(store Store) *server {
	return &server{store: store, now: time.Now}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *server) handleBudget(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getBudget(w, r)
	case http.MethodPatch:
		s.updateBudgetHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) getBudget(w http.ResponseWriter, r *http.Request) {
	period, err := s.store.GetOrCreatePeriod()
	if err != nil {
		log.Printf("fehler beim laden der periode: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
		return
	}

	totalSpent, err := s.store.GetTotalExpenses(period.ID)
	if err != nil {
		log.Printf("fehler beim laden der ausgaben: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der ausgaben")
		return
	}

	now := s.now()
	todaySpent, err := s.store.GetTodayExpenses(period.ID, now)
	if err != nil {
		log.Printf("fehler beim laden der heutigen ausgaben: %v", err)
		todaySpent = 0
	}

	expenses, err := s.store.GetRecentExpenses(period.ID, 3)
	if err != nil {
		log.Printf("fehler beim laden der letzten ausgaben: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der ausgaben")
		return
	}

	day := period.dayOfMonth(now)
	currentBudget, savings, color := period.calcBudget(totalSpent, todaySpent, now)

	dailyStats, err := s.store.GetDailyExpenses(period.ID, period.StartDate, day)
	if err != nil {
		log.Printf("fehler beim laden der tages-statistiken: %v", err)
		dailyStats = []DailyStat{}
	}

	resp := BudgetResponse{
		Day:           day,
		MonthDays:     period.MonthDays,
		BaseBudget:    period.BaseBudget,
		CurrentBudget: currentBudget,
		Savings:       savings,
		Color:         color,
		PeriodID:      period.ID,
		Expenses:      expenses,
		DailyStats:    dailyStats,
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, resp)
}

func (s *server) updateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültige anfrage")
		return
	}

	if req.MonthlyTotal <= 0 && req.Days <= 0 {
		writeError(w, http.StatusBadRequest, "mindestens budget oder tage müssen > 0 sein")
		return
	}

	if err := s.store.UpdateBudget(req.MonthlyTotal, req.Days); err != nil {
		log.Printf("fehler beim aktualisieren des budgets: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim aktualisieren")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) handleExpenses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getExpenses(w, r)
	case http.MethodPost:
		s.createExpense(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) getExpenses(w http.ResponseWriter, r *http.Request) {
	period, err := s.store.GetOrCreatePeriod()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
		return
	}

	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	expenses, err := s.store.GetExpenses(period.ID, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim laden der ausgaben")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, expenses)
}

func (s *server) createExpense(w http.ResponseWriter, r *http.Request) {
	var req ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungültige anfrage")
		return
	}

	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "betrag muss > 0 sein")
		return
	}

	period, err := s.store.GetOrCreatePeriod()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
		return
	}

	expense, err := s.store.AddExpense(period.ID, req.Amount, req.Note)
	if err != nil {
		log.Printf("fehler beim erstellen der ausgabe: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim speichern")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusCreated, expense)
}

func (s *server) handleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "fehlende ausgaben-id")
		return
	}
	expenseID := parts[3]

	if err := s.store.DeleteExpense(expenseID); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "nicht gefunden") || strings.Contains(msg, "uuid") {
			writeError(w, http.StatusNotFound, "ausgabe nicht gefunden")
		} else {
			log.Printf("fehler beim löschen: %v", err)
			writeError(w, http.StatusInternalServerError, "fehler beim löschen")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) handleNewPeriod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req NewPeriodRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	startDate := s.now()
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = t
		}
	}

	var period *Period
	var err error
	period, err = s.store.CreatePeriodWithStart(startDate, req.MonthlyTotal, req.Days)

	if err != nil {
		log.Printf("fehler beim erstellen der periode: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim erstellen der periode")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusCreated, period)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.store.Ping(); err != nil {
		log.Printf("health check fehler: %v", err)
		writeError(w, http.StatusServiceUnavailable, "datenbank nicht erreichbar")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/budget", s.handleBudget)
	mux.HandleFunc("/api/expenses", s.handleExpenses)
	mux.HandleFunc("/api/expenses/", s.handleDeleteExpense)
	mux.HandleFunc("/api/period", s.handleNewPeriod)
	return corsMiddleware(mux)
}
