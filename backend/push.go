package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// ---------------------------------------------------------------------------
// VAPID-Konfiguration
// ---------------------------------------------------------------------------

// demoVAPIDKeys liefert ein festes Demo-Schlüsselpaar als Fallback, falls
// keine VAPID_* Umgebungsvariablen gesetzt sind (lokale Entwicklung).
const (
	demoVAPIDPublic  = "BE7jYw9Ucjzq0m0yJ3LbHnL5vQ0XgWqNpRx2a4vBkM8cE1TfOoIuUyYdS6aJ7bQvZ0wXgNuR3hFpD1sKlZtGcVbHnM"
	demoVAPIDPrivate = "HcVbHnMq7jYw9Ucjzq0m0yJ3LbHnL5vQ0XgWqNpRx2a4vBkM8cE1TfOoIuUy"
)

func vapidPublicKey() string {
	if k := os.Getenv("VAPID_PUBLIC_KEY"); k != "" {
		return k
	}
	return demoVAPIDPublic
}

func vapidPrivateKey() string {
	if k := os.Getenv("VAPID_PRIVATE_KEY"); k != "" {
		return k
	}
	return demoVAPIDPrivate
}

func vapidSubject() string {
	if s := os.Getenv("VAPID_SUBJECT"); s != "" {
		return s
	}
	return "mailto:notifications@restgeld.local"
}

// generateVAPIDKeys erzeugt ein frisches Schlüsselpaar (für Setup-Hilfe).
func generateVAPIDKeys() (public string, private string, err error) {
	pub, priv, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return "", "", err
	}
	return pub, priv, nil
}

// ---------------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------------

type pushSubscribeRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func (s *server) handlePushSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := s.getUserIDFromRequest(r)

	if r.Method == http.MethodDelete {
		var req pushSubscribeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Endpoint == "" {
			writeError(w, http.StatusBadRequest, "endpoint erforderlich")
			return
		}
		if err := s.store.DeletePushSubscription(userID, req.Endpoint); err != nil {
			log.Printf("push abonnement loeschen fehlgeschlagen: %v", err)
			writeError(w, http.StatusInternalServerError, "loeschen fehlgeschlagen")
			return
		}
		jsonHeader(w)
		writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
		return
	}

	var req pushSubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungueltiger request")
		return
	}
	if req.Endpoint == "" || req.Keys.P256dh == "" || req.Keys.Auth == "" {
		writeError(w, http.StatusBadRequest, "endpoint, p256dh und auth erforderlich")
		return
	}

	if err := s.store.SavePushSubscription(userID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth); err != nil {
		log.Printf("push abonnement speichern fehlgeschlagen: %v", err)
		writeError(w, http.StatusInternalServerError, "speichern fehlgeschlagen")
		return
	}

	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]bool{"subscribed": true})
}

type pushSendRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// handlePushSend versendet eine Benachrichtigung an alle Abonnenten des Users.
// Gäste (userID == "") erhalten eine Benachrichtigung an die Standard-Gästeliste.
func (s *server) handlePushSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req pushSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungueltiger request")
		return
	}
	if req.Title == "" {
		req.Title = "restgeld"
	}

	userID := s.getUserIDFromRequest(r)
	subs, err := s.store.ListPushSubscriptions(userID)
	if err != nil {
		log.Printf("push abonnements laden fehlgeschlagen: %v", err)
		writeError(w, http.StatusInternalServerError, "abonnements laden fehlgeschlagen")
		return
	}

	sent, failed := s.sendPushToSubscriptions(subs, req.Title, req.Body, userID)
	jsonHeader(w)
	writeJSON(w, http.StatusOK, map[string]any{"sent": sent, "failed": failed, "total": len(subs)})
}

func (s *server) sendPushToSubscriptions(subs []PushSubscription, title, body, userID string) (sent, failed int) {
	for _, sub := range subs {
		webSub := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}

		payload, _ := json.Marshal(map[string]string{
			"title": title,
			"body":  body,
			"url":   "/",
		})

		resp, err := webpush.SendNotification(payload, webSub, &webpush.Options{
			Subscriber:      vapidSubject(),
			VAPIDPublicKey:  vapidPublicKey(),
			VAPIDPrivateKey: vapidPrivateKey(),
			TTL:             60,
		})
		if err != nil {
			// 404/410: Abonnement nicht mehr gültig -> bereinigen
			if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
				s.store.DeletePushSubscription(userID, sub.Endpoint)
			}
			log.Printf("push senden fehlgeschlagen: %v", err)
			failed++
			continue
		}
		if resp != nil {
			resp.Body.Close()
		}
		sent++
	}
	return sent, failed
}

// handlePushVapidKeys liefert den VAPID-Public-Key fürs Frontend und
// erlaubt die (erste) Key-Generierung via POST ?generate=true.
func (s *server) handlePushVapidKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonHeader(w)
		writeJSON(w, http.StatusOK, map[string]string{"publicKey": vapidPublicKey()})
	case http.MethodPost:
		pub, priv, err := generateVAPIDKeys()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "key generierung fehlgeschlagen")
			return
		}
		jsonHeader(w)
		writeJSON(w, http.StatusOK, map[string]string{
			"publicKey":  pub,
			"privateKey": priv,
			"subject":    vapidSubject(),
			"hint":       "Setze die Keys in VAPID_PUBLIC_KEY / VAPID_PRIVATE_KEY (nur lokal nutzen).",
		})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func base64RawURL(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}
