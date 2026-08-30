package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ServiceTarget struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Required bool   `json:"required"`
}

type ServiceTelemetry struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Status    string         `json:"status"` // "up", "degraded", "down"
	LatencyMs float64        `json:"latencyMs"`
	URL       string         `json:"url"`
	CheckedAt time.Time      `json:"checkedAt"`
	Details   map[string]any `json:"details,omitempty"`
	Error     string         `json:"error,omitempty"`
}

type SystemStats struct {
	GoVersion     string  `json:"goVersion"`
	Goroutines    int     `json:"goroutines"`
	MemoryAllocMB float64 `json:"memoryAllocMb"`
	MemorySysMB   float64 `json:"memorySysMb"`
	GCCount       uint32  `json:"gcCount"`
	UptimeSeconds int64   `json:"uptimeSeconds"`
}

type ClusterOverview struct {
	Status        string             `json:"status"` // "healthy", "degraded", "critical"
	Timestamp     time.Time          `json:"timestamp"`
	UptimeSeconds int64              `json:"uptimeSeconds"`
	Services      []ServiceTelemetry `json:"services"`
	System        SystemStats        `json:"system"`
}

type MonitoringService struct {
	startTime  time.Time
	targets    []ServiceTarget
	httpClient *http.Client
}

func NewMonitoringService() *MonitoringService {
	coreURL := getEnvOrDefault("CORE_API_URL", "http://localhost:8080")
	authURL := getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8081")
	billingURL := getEnvOrDefault("BILLING_SERVICE_URL", "http://localhost:8082")

	targets := []ServiceTarget{
		{ID: "core", Name: "Core API Backend", URL: coreURL, Required: true},
		{ID: "auth", Name: "Auth Service", URL: authURL, Required: true},
		{ID: "billing", Name: "Billing Service", URL: billingURL, Required: false},
	}

	return &MonitoringService{
		startTime: time.Now(),
		targets:   targets,
		httpClient: &http.Client{
			Timeout: 2500 * time.Millisecond,
		},
	}
}

func getEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return strings.TrimRight(val, "/")
}

func (s *MonitoringService) probeTarget(ctx context.Context, target ServiceTarget) ServiceTelemetry {
	healthURL := fmt.Sprintf("%s/health", target.URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	if err != nil {
		return ServiceTelemetry{
			ID:        target.ID,
			Name:      target.Name,
			Status:    "down",
			LatencyMs: 0,
			URL:       target.URL,
			CheckedAt: time.Now(),
			Error:     err.Error(),
		}
	}

	start := time.Now()
	resp, err := s.httpClient.Do(req)
	latency := time.Since(start).Seconds() * 1000.0

	if err != nil {
		return ServiceTelemetry{
			ID:        target.ID,
			Name:      target.Name,
			Status:    "down",
			LatencyMs: latency,
			URL:       target.URL,
			CheckedAt: time.Now(),
			Error:     "connection refused / timeout",
		}
	}
	defer resp.Body.Close()

	status := "up"
	if resp.StatusCode != http.StatusOK {
		status = "degraded"
	}

	var details map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&details)

	return ServiceTelemetry{
		ID:        target.ID,
		Name:      target.Name,
		Status:    status,
		LatencyMs: mathRound(latency, 2),
		URL:       target.URL,
		CheckedAt: time.Now(),
		Details:   details,
	}
}

func (s *MonitoringService) CollectTelemetry(ctx context.Context) ClusterOverview {
	var wg sync.WaitGroup
	results := make([]ServiceTelemetry, len(s.targets))

	for i, target := range s.targets {
		wg.Add(1)
		go func(idx int, t ServiceTarget) {
			defer wg.Done()
			results[idx] = s.probeTarget(ctx, t)
		}(i, target)
	}

	wg.Wait()

	// Add self monitoring telemetry
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	selfTelemetry := ServiceTelemetry{
		ID:        "monitoring",
		Name:      "Observability Service",
		Status:    "up",
		LatencyMs: 0.1,
		URL:       "http://localhost:8083",
		CheckedAt: time.Now(),
		Details: map[string]any{
			"status":     "ok",
			"goroutines": runtime.NumGoroutine(),
			"version":    "1.0.0",
		},
	}
	allServices := append([]ServiceTelemetry{selfTelemetry}, results...)

	// Evaluate overall cluster health
	clusterStatus := "healthy"
	for _, svc := range allServices {
		if svc.Status == "down" {
			clusterStatus = "degraded"
		}
	}

	uptime := int64(time.Since(s.startTime).Seconds())
	system := SystemStats{
		GoVersion:     runtime.Version(),
		Goroutines:    runtime.NumGoroutine(),
		MemoryAllocMB: mathRound(float64(memStats.Alloc)/(1024*1024), 2),
		MemorySysMB:   mathRound(float64(memStats.Sys)/(1024*1024), 2),
		GCCount:       memStats.NumGC,
		UptimeSeconds: uptime,
	}

	return ClusterOverview{
		Status:        clusterStatus,
		Timestamp:     time.Now(),
		UptimeSeconds: uptime,
		Services:      allServices,
		System:        system,
	}
}

func mathRound(val float64, precision int) float64 {
	p := 1.0
	for i := 0; i < precision; i++ {
		p *= 10.0
	}
	return float64(int(val*p+0.5)) / p
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *MonitoringService) HandleOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	overview := s.CollectTelemetry(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(overview)
}

func (s *MonitoringService) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	overview := s.CollectTelemetry(r.Context())
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	var sb strings.Builder
	sb.WriteString("# HELP restgeld_service_up Status of Restgeld microservice (1 = up, 0 = down)\n")
	sb.WriteString("# TYPE restgeld_service_up gauge\n")
	for _, svc := range overview.Services {
		upVal := 1
		if svc.Status != "up" {
			upVal = 0
		}
		sb.WriteString(fmt.Sprintf("restgeld_service_up{service=\"%s\",id=\"%s\"} %d\n", svc.Name, svc.ID, upVal))
	}

	sb.WriteString("\n# HELP restgeld_service_latency_ms Measured HTTP latency in milliseconds\n")
	sb.WriteString("# TYPE restgeld_service_latency_ms gauge\n")
	for _, svc := range overview.Services {
		sb.WriteString(fmt.Sprintf("restgeld_service_latency_ms{service=\"%s\",id=\"%s\"} %.2f\n", svc.Name, svc.ID, svc.LatencyMs))
	}

	sb.WriteString("\n# HELP restgeld_uptime_seconds Total service uptime in seconds\n")
	sb.WriteString("# TYPE restgeld_uptime_seconds counter\n")
	sb.WriteString(fmt.Sprintf("restgeld_uptime_seconds %d\n", overview.UptimeSeconds))

	sb.WriteString("\n# HELP restgeld_goroutines Current active Go goroutines\n")
	sb.WriteString("# TYPE restgeld_goroutines gauge\n")
	sb.WriteString(fmt.Sprintf("restgeld_goroutines %d\n", overview.System.Goroutines))

	sb.WriteString("\n# HELP restgeld_memory_alloc_mb Current allocated memory in MB\n")
	sb.WriteString("# TYPE restgeld_memory_alloc_mb gauge\n")
	sb.WriteString(fmt.Sprintf("restgeld_memory_alloc_mb %.2f\n", overview.System.MemoryAllocMB))

	_, _ = w.Write([]byte(sb.String()))
}

func (s *MonitoringService) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"service":   "monitoring-service",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(s.startTime).String(),
	})
}

func main() {
	service := NewMonitoringService()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", service.HandleHealth)
	mux.HandleFunc("/metrics", service.HandleMetrics)
	mux.HandleFunc("/api/monitoring/overview", service.HandleOverview)
	mux.HandleFunc("/api/monitoring/health", service.HandleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	handler := corsMiddleware(mux)
	log.Printf("🚀 Restgeld Monitoring & Observability Service listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Fatal server error: %v", err)
	}
}
