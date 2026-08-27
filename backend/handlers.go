package main

import (
	"encoding/json"
	"fmt"
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

	streak := calcStreakInfo(dailyStats, period.BaseBudget)
	projection := calcProjection(totalSpent, day, period.MonthDays, period.MonthlyTotal)

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
		Streak:        streak,
		Projection:    projection,
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
	periodID := r.URL.Query().Get("period_id")
	if periodID == "" {
		period, err := s.store.GetOrCreatePeriod()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
			return
		}
		periodID = period.ID
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

	expenses, err := s.store.GetExpenses(periodID, page, limit)
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

type ExportBackup struct {
	Version    string    `json:"version"`
	ExportedAt time.Time `json:"exportedAt"`
	Period     Period    `json:"period"`
	Expenses   []Expense `json:"expenses"`
}

type ImportRequest struct {
	Expenses []Expense `json:"expenses"`
}

func (s *server) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	period, err := s.store.GetOrCreatePeriod()
	if err != nil {
		log.Printf("fehler beim laden der periode fuer export: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
		return
	}

	expenses, err := s.store.GetAllExpenses(period.ID)
	if err != nil {
		log.Printf("fehler beim laden der ausgaben fuer export: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der ausgaben")
		return
	}

	format := strings.ToLower(r.URL.Query().Get("format"))
	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"restgeld-export-%s.csv\"", period.ID))
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Datum;Uhrzeit;Betrag;Notiz\n"))
		for _, exp := range expenses {
			dateStr := exp.CreatedAt.Format("2006-01-02")
			timeStr := exp.CreatedAt.Format("15:04:05")
			amountStr := fmt.Sprintf("%.2f", exp.Amount)
			escapedNote := strings.ReplaceAll(exp.Note, "\"", "\"\"")
			line := fmt.Sprintf("%s;%s;%s;\"%s\"\n", dateStr, timeStr, amountStr, escapedNote)
			w.Write([]byte(line))
		}
		return
	}

	// JSON Export (Standard)
	backup := ExportBackup{
		Version:    "1.0",
		ExportedAt: s.now(),
		Period:     *period,
		Expenses:   expenses,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"restgeld-backup-%s.json\"", period.ID))
	writeJSON(w, http.StatusOK, backup)
}

func (s *server) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	period, err := s.store.GetOrCreatePeriod()
	if err != nil {
		log.Printf("fehler beim laden der periode fuer import: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim laden der periode")
		return
	}

	contentType := r.Header.Get("Content-Type")
	var expensesToImport []Expense

	if strings.Contains(contentType, "text/csv") {
		// CSV Import
		bodyBytes, err := ioReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "fehler beim lesen des request-bodys")
			return
		}
		lines := strings.Split(string(bodyBytes), "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || (i == 0 && (strings.HasPrefix(strings.ToLower(line), "datum") || strings.HasPrefix(strings.ToLower(line), "date"))) {
				continue
			}
			// Trennzeichen: Semikolon oder Komma
			sep := ";"
			if !strings.Contains(line, ";") && strings.Contains(line, ",") {
				sep = ","
			}
			parts := strings.Split(line, sep)
			if len(parts) >= 2 {
				dateStr := strings.TrimSpace(parts[0])
				var amountStr string
				var noteStr string
				if len(parts) >= 4 {
					amountStr = strings.TrimSpace(parts[2])
					noteStr = strings.Trim(strings.TrimSpace(parts[3]), "\"")
				} else {
					amountStr = strings.TrimSpace(parts[1])
					if len(parts) > 2 {
						noteStr = strings.Trim(strings.TrimSpace(parts[2]), "\"")
					}
				}

				amountStr = strings.ReplaceAll(amountStr, ",", ".")
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err == nil && amount > 0 {
					createdAt := s.now()
					if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
						createdAt = parsedDate
					}
					expensesToImport = append(expensesToImport, Expense{
						Amount:    amount,
						Note:      noteStr,
						CreatedAt: createdAt,
					})
				}
			}
		}
	} else {
		// JSON Import (akzeptiert ExportBackup oder Array []Expense oder ImportRequest)
		bodyBytes, err := ioReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "fehler beim lesen des request-bodys")
			return
		}

		var backup ExportBackup
		if err := json.Unmarshal(bodyBytes, &backup); err == nil && len(backup.Expenses) > 0 {
			expensesToImport = backup.Expenses
		} else {
			var list []Expense
			if err := json.Unmarshal(bodyBytes, &list); err == nil {
				expensesToImport = list
			} else {
				var req ImportRequest
				if err := json.Unmarshal(bodyBytes, &req); err == nil {
					expensesToImport = req.Expenses
				}
			}
		}
	}

	if len(expensesToImport) == 0 {
		writeError(w, http.StatusBadRequest, "keine gueltigen ausgaben zum importieren gefunden")
		return
	}

	imported, err := s.store.ImportExpenses(period.ID, expensesToImport)
	if err != nil {
		log.Printf("fehler beim importieren: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim importieren der ausgaben")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"imported": imported,
	})
}

type ioReader interface {
	Read(p []byte) (n int, err error)
}

func ioReadAll(r ioReader) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	var buf []byte
	b := make([]byte, 1024)
	for {
		n, err := r.Read(b)
		if n > 0 {
			buf = append(buf, b[:n]...)
		}
		if err != nil {
			break
		}
	}
	return buf, nil
}

func (s *server) handleGetPeriods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	periods, err := s.store.GetAllPeriods()
	if err != nil {
		log.Printf("fehler beim abrufen aller perioden: %v", err)
		writeError(w, http.StatusInternalServerError, "fehler beim abrufen der perioden")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, periods)
}

func (s *server) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/budget", s.handleBudget)
	mux.HandleFunc("/api/expenses", s.handleExpenses)
	mux.HandleFunc("/api/expenses/", s.handleDeleteExpense)
	mux.HandleFunc("/api/period", s.handleNewPeriod)
	mux.HandleFunc("/api/periods", s.handleGetPeriods)
	mux.HandleFunc("/api/export", s.handleExport)
	mux.HandleFunc("/api/import", s.handleImport)
	return corsMiddleware(mux)
}
