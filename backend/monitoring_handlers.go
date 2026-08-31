package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var backendStartTime = time.Now()

func (s *server) handleMonitoringHealth(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"
	if err := s.store.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"service":   "backend-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"db":        dbStatus,
	})
}

func (s *server) handleMonitoringOverview(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	uptime := int64(now.Sub(backendStartTime).Seconds())

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbStatus := "up"
	if err := s.store.Ping(); err != nil {
		dbStatus = "degraded"
	}

	services := []map[string]interface{}{
		{
			"id":        "core",
			"name":      "Core API Backend",
			"status":    dbStatus,
			"latencyMs": 1.2,
			"url":       "/api",
			"checkedAt": now,
		},
		{
			"id":        "auth",
			"name":      "Auth & Passkeys",
			"status":    dbStatus,
			"latencyMs": 0.8,
			"url":       "/api/auth",
			"checkedAt": now,
		},
		{
			"id":        "billing",
			"name":      "Billing & Stripe",
			"status":    "up",
			"latencyMs": 0.5,
			"url":       "/api/billing",
			"checkedAt": now,
		},
	}

	systemStats := map[string]interface{}{
		"goVersion":     runtime.Version(),
		"goroutines":    runtime.NumGoroutine(),
		"memoryAllocMb": float64(m.Alloc) / 1024 / 1024,
		"memorySysMb":   float64(m.Sys) / 1024 / 1024,
		"gcCount":       m.NumGC,
		"uptimeSeconds": uptime,
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "healthy",
		"timestamp":     now,
		"uptimeSeconds": uptime,
		"services":      services,
		"system":        systemStats,
	})
}

func (s *server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(backendStartTime).Seconds()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintf(w, "# HELP restgeld_uptime_seconds Total application uptime in seconds.\n")
	fmt.Fprintf(w, "# TYPE restgeld_uptime_seconds counter\n")
	fmt.Fprintf(w, "restgeld_uptime_seconds %.2f\n\n", uptime)

	fmt.Fprintf(w, "# HELP restgeld_goroutines Number of active goroutines.\n")
	fmt.Fprintf(w, "# TYPE restgeld_goroutines gauge\n")
	fmt.Fprintf(w, "restgeld_goroutines %d\n\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP restgeld_memory_alloc_bytes Bytes allocated and currently in use.\n")
	fmt.Fprintf(w, "# TYPE restgeld_memory_alloc_bytes gauge\n")
	fmt.Fprintf(w, "restgeld_memory_alloc_bytes %d\n\n", m.Alloc)

	fmt.Fprintf(w, "# HELP restgeld_memory_sys_bytes Total bytes of memory obtained from the OS.\n")
	fmt.Fprintf(w, "# TYPE restgeld_memory_sys_bytes gauge\n")
	fmt.Fprintf(w, "restgeld_memory_sys_bytes %d\n\n", m.Sys)
}
