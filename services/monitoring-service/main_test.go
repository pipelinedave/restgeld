package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMonitoringServiceHealth(t *testing.T) {
	svc := NewMonitoringService()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	svc.HandleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var res map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}

	if res["status"] != "ok" || res["service"] != "monitoring-service" {
		t.Fatalf("unexpected health payload: %v", res)
	}
}

func TestMonitoringServiceOverview(t *testing.T) {
	// Mock core API server
	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": "1.0"})
	}))
	defer mockCore.Close()

	svc := &MonitoringService{
		startTime: time.Now().Add(-10 * time.Minute),
		targets: []ServiceTarget{
			{ID: "core", Name: "Core API", URL: mockCore.URL, Required: true},
			{ID: "auth", Name: "Auth Service", URL: "http://127.0.0.1:59999", Required: true}, // Offline target
		},
		httpClient: &http.Client{Timeout: 500 * time.Millisecond},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/monitoring/overview", nil)
	rr := httptest.NewRecorder()

	svc.HandleOverview(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var overview ClusterOverview
	if err := json.Unmarshal(rr.Body.Bytes(), &overview); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}

	if len(overview.Services) != 3 { // self (monitoring) + 2 targets
		t.Fatalf("expected 3 services in overview, got %d", len(overview.Services))
	}

	// First is self (monitoring)
	if overview.Services[0].ID != "monitoring" || overview.Services[0].Status != "up" {
		t.Errorf("expected monitoring service to be up, got %v", overview.Services[0])
	}

	// Second is mockCore
	if overview.Services[1].ID != "core" || overview.Services[1].Status != "up" {
		t.Errorf("expected core service to be up, got %v", overview.Services[1])
	}

	// Third is offline target
	if overview.Services[2].ID != "auth" || overview.Services[2].Status != "down" {
		t.Errorf("expected auth service to be down, got %v", overview.Services[2])
	}

	if overview.System.Goroutines <= 0 {
		t.Errorf("expected positive goroutine count, got %d", overview.System.Goroutines)
	}
}

func TestMonitoringServiceMetrics(t *testing.T) {
	mockTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockTarget.Close()

	svc := &MonitoringService{
		startTime: time.Now().Add(-5 * time.Minute),
		targets: []ServiceTarget{
			{ID: "core", Name: "Core API", URL: mockTarget.URL, Required: true},
		},
		httpClient: &http.Client{Timeout: 500 * time.Millisecond},
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	svc.HandleMetrics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "restgeld_service_up") {
		t.Errorf("expected restgeld_service_up metric in output")
	}
	if !strings.Contains(body, "restgeld_goroutines") {
		t.Errorf("expected restgeld_goroutines metric in output")
	}
	if !strings.Contains(body, "restgeld_memory_alloc_mb") {
		t.Errorf("expected restgeld_memory_alloc_mb metric in output")
	}
}
