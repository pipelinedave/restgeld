package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePushSubscribeAndList(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store}

	// POST /api/push/subscribe
	body := `{"endpoint":"https://push.example.com/sub/1","keys":{"p256dh":"AAA","auth":"BBB"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/push/subscribe", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}

	subs, err := store.ListPushSubscriptions("")
	if err != nil {
		t.Fatalf("abonnements laden: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("erwartet 1 abonnement, bekommen %d", len(subs))
	}
	if subs[0].Endpoint != "https://push.example.com/sub/1" {
		t.Fatalf("falscher endpoint: %s", subs[0].Endpoint)
	}
}

func TestHandlePushSubscribeInvalid(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store}

	// Fehlende p256dh/auth -> 400
	body := `{"endpoint":"https://push.example.com/sub/2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/push/subscribe", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("erwartet 400, bekommen %d", rec.Code)
	}
}

func TestHandlePushDelete(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store}

	store.SavePushSubscription("", "https://push.example.com/sub/3", "AAA", "BBB")

	body := `{"endpoint":"https://push.example.com/sub/3"}`
	req := httptest.NewRequest(http.MethodDelete, "/api/push/subscribe", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", rec.Code)
	}

	subs, _ := store.ListPushSubscriptions("")
	if len(subs) != 0 {
		t.Fatalf("erwartet 0 abonnements nach loeschung, bekommen %d", len(subs))
	}
}

func TestHandlePushVapidKeys(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/push/vapid", nil)
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decodieren: %v", err)
	}
	if resp["publicKey"] == "" {
		t.Fatalf("publicKey fehlt")
	}
}

func TestHandlePushSendWithInvalidEndpoint(t *testing.T) {
	store := newMemoryStore()
	srv := &server{store: store}

	// Ungueltiger Endpoint -> senden schlaegt fehl -> failed=1, sent=0
	store.SavePushSubscription("", "https://invalid.invalid/sub/9", "AAA", "BBB")

	body := `{"title":"Test","body":"Hi"}`
	req := httptest.NewRequest(http.MethodPost, "/api/push/send", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	srv.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]int
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("json decodieren: %v", err)
	}
	if resp["total"] != 1 {
		t.Fatalf("erwartet total=1, bekommen %d", resp["total"])
	}
}
