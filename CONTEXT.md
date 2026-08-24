# Restgeld – Daily Allowance Tracker (PWA)

## Projektbeschreibung
Minimalistische Web-App (Mobile-first PWA) zum Tracken eines täglichen "Spending Money" Budgets.
Betrieben auf lokalem single-node k3s-Cluster.

## Architektur (3-Tier)
1. **Frontend:** Vue 3 + Vite + TypeScript + PWA (Service Worker via vite-plugin-pwa), Nginx
2. **Backend:** Go REST API (net/http + lib/pq), port 8080
3. **Datenbank:** PostgreSQL 16-alpine mit PVC

## Domain & Infrastruktur
- **Domain:** restgeld.stillon.top
- **Ingress:** nginx-ingress, TLS via letsencrypt-prod (cert-manager)
- **Namespace:** restgeld (manuell via kubectl create namespace restgeld)
- **Flux Kustomization:** apps/restgeld.yaml → kustomize/restgeld/
- **Secrets:** SealedSecret (bitnami.com/v1alpha1)

## API-Endpunkte

| Methode | Pfad | Funktion |
|---|---|---|
| GET | /api/health | Health-Check (DB-Connectivity) |
| GET | /api/budget | Aktuelle Periode + Tagesbudget + Ersparnis |
| POST | /api/period | Neue Periode starten (reset) |
| PATCH | /api/budget | Monatsbudget ändern |
| GET | /api/expenses | Letzte 3 Ausgaben |
| POST | /api/expenses | Ausgabe buchen ({amount, note}) |
| DELETE | /api/expenses/{id} | Ausgabe löschen |

## Kernlogik
- Basis-Budget = Monatsbudget / Tage im Monat
- savings = (baseBudget * day) - totalSpentSoFar
- currentBudget = baseBudget + (savings / remainingDays)
- Farbcodierung: green (savings > 0), white (savings == 0), red (savings < 0)

## Verzeichnisstruktur

```
restgeld/
├── AGENTS.md        # Agent-Workflow
├── CONTEXT.md       # Projektstatus (diese Datei)
├── prompt.md        # Originale Anforderung
├── backend/         # Go REST API
│   ├── main.go
│   ├── handlers.go
│   ├── models.go
│   ├── db.go
│   ├── go.mod
│   └── Dockerfile
├── frontend/        # Vue 3 + Vite + PWA
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── style.css
│   │   ├── composables/useApi.ts
│   │   └── components/
│   │       ├── AppHeader.vue
│   │       ├── AppFooter.vue
│   │       ├── MonthProgress.vue
│   │       ├── BudgetDisplay.vue
│   │       ├── Numpad.vue
│   │       ├── RecentExpenses.vue
│   │       └── SettingsModal.vue
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   ├── nginx.conf
│   └── Dockerfile
└── k8s/
    └── restgeld/     # Kubernetes-Manifeste (für k3s-config)
        ├── kustomization.yaml
        ├── postgres-pvc.yaml
        ├── postgres-deployment.yaml
        ├── postgres-service.yaml
        ├── backend-deployment.yaml
        ├── backend-service.yaml
        ├── frontend-deployment.yaml
        ├── frontend-service.yaml
        ├── restgeld-secrets-sealed.yaml
        ├── ingress.yaml
        └── apps/
            └── restgeld.yaml
```

## Status

### ✅ Erledigt
- AGENTS.md erstellt (Workflow-Definition mit EZE-Dev-Zyklen)
- CONTEXT.md erstellt (Plan-Dokument)
- Backend: Go REST API (models, store-interface, postgres-store, handlers, main, Dockerfile)
- Backend-Unit-Tests: 15 httptest-Tests (in-memory store) ✅
- Backend-Integration-Tests: 8 Tests mit echter PostgreSQL (db_test.go, build-tag: integration) ✅
- Docker-Build backend:test + frontend:test ✅
- Frontend: Vue 3 + Vite + PWA (MonthProgress, BudgetDisplay, Numpad, RecentExpenses)
- Frontend-Unit-Tests: 35 vitest-Tests (vue-test-utils, happy-dom) ✅
- Frontend-Coverage: vitest + @vitest/coverage-v8 (lcov) ✅
- K8s-Manifeste: postgres/backend/frontend/ingress/flux-kustomization in k8s/restgeld/
- Git-Log mit 7 Commits (backend, frontend, k8s, ci, docs, fixes, qa)
- AGENTS.md: Commit-Regel auf "nach jeder logischen Einheit"
- **Lokaler Testlauf vollstaendig** ✅
  - GET/POST/PATCH/DELETE API getestet
  - Neue Periode starten (reset)
  - Numpad-Eingabe via DevTools (keine native Tastatur) + Notiz-Feld
  - Ausgabe buchen + löschen via UI
  - Budget-Konfiguration aendern
  - SPA-Routing funktioniert
  - Console 0 errors/0 warnings
  - PWA Service Worker registriert
- GitHub Actions Workflows:
  - backend.yml: go vet, unit+integration tests, codecov
  - frontend.yml: vue-tsc, vite build, vitest + coverage, codecov
  - docker.yml: docker build backend + frontend
- GitHub Repo: pipelinedave/restgeld (private)
- README.md mit CI/CD-Badges und Codecov
- TESTSTRATEGIE.md dokumentiert
- Numpad: 2-Step Flow (Betrag → Notiz → Speichern)
- Backend-Fixes:
  - createPeriod: ON CONFLICT + expenses löschen
  - DELETE: ungültige UUID → 404 statt 500
- Backend: Health-Endpoint `/api/health` implementiert (DB-Ping) und K8s-Probes umgestellt ✅
- Backend: Datenbank-Migrationen via `schema_migrations` und `embed.FS` mit Start-Migration `001_initial.sql` implementiert und getestet ✅
- Frontend: Playwright E2E-Tests repariert (Dev-Server-Anbindung, Mock-API mit dynamischem State & GitHub Actions CI-Integration) ✅
- 1-Befehl lokales Dev-Environment via Docker Compose (`docker-compose.yml`, `scripts/dev.sh`, `scripts/dev.ps1`) ✅
- UI-Konfiguration: Settings-Modal (`SettingsModal.vue`) zum Einstellen des Monatsbudgets und Zurücksetzen der Periode mit Bestätigungsdialog implementiert und getestet ✅
- UI: Stylischer App-Header mit Titel "restgeld." in eigenständige Komponente `AppHeader.vue` extrahiert und sauber von `MonthProgress.vue` entkoppelt ✅
- UI: Minimalistischer App-Footer (`AppFooter.vue`) mit dezentem Trennstrich, Tagline und Versionsanzeige für verbesserte visuelle Struktur implementiert und getestet ✅
- Flexible Perioden & Payday-Cycles: Start einer neuen Periode ab dem heutigen Tag (beginnt sofort bei Tag 1 von N), nahtloser Monatsübergang und Migration `002_flexible_periods.sql` implementiert und verifiziert ✅
- Korrektur der Budget- & Sparbetragsformel: Heutiges Budget zählt nicht mehr vorab als Ersparnis; an Tag 1 startet der Tracker sauber mit genau $300 / 31 = 9,68 €$ und 0 € Ersparnis ✅
- Ausgaben-Historie mit Paginierung: Backend-Endpoint `/api/expenses?page=X&limit=Y` mit `PaginatedExpenses`, Frontend "Alle anzeigen"-Button in `RecentExpenses.vue` und stylisches, voll funktionales `ExpensesModal.vue` mit Zeitstempel-Formatierung, Paginierung und Einzel-Löschung implementiert und vollständig mit Unit-, Snapshot- und E2E-Tests abgedeckt ✅
- Hero-Card Redesign & intuitive Tages-Restgeld-Logik: Große Hero-Zahl zeigt nun immer exakt das heute noch verfügbare Budget (`HEUTE VERFÜGBAR`), subtiler Status-Badge zeigt den echten Monats-Puffer (Grün) bzw. Überzug (Rot). Ausgabe buchen zieht direkt vom heutigen Restgeld ab ✅
- Mobile Numpad & Keyboard UX: Ausgaben-Erfassung auf native Mobile-Tastatur (`inputmode="decimal"`) umgestellt, 1-stufiger kompakter Dialog mit Betrag + Notiz, Vermeidung von Überdeckungen durch `interactive-widget=resizes-content` und `dvh`-Viewport-Handling ✅
- Konfigurierbare Periodendauer (Tage): Periodendauer kann im Einstellungs-Dialog sowohl mit sofortiger Wirkung für die laufende Periode angepasst als auch für neue Abrechnungszyklen frei definiert werden (z. B. 14 Tage, 20 Tage, etc.) ✅
- UX- & Micro-Feedback Upgrade: Taktiles haptisches Feedback (`useHaptics` via Web Vibration API) bei Tastendruck, erfolgreichem Buchen/Löschen und Warnungen; schwebende Toast-Benachrichtigungen (`ToastNotification.vue`), Lade-Spinner & Button-Blockierung beim Buchen zur Vermeidung von Doppelbuchungen sowie Zahlen-Puls-Animation bei Budgetänderungen implementiert und mit 74 Vitest-Tests abgedeckt ✅
- Trend-Chart & Sparkline-Verlauf: Backend `GetDailyExpenses` & `dailyStats` in `BudgetResponse`, neue interaktive Dashboard-Komponente `SpendingTrend.vue` mit Tagesbalken, Farbcodierung (Grün/Rot/Spar-Tag), Basis-Budget-Referenzlinie, Ø Tagesschnitt und Tap-Details implementiert und mit 77 Vitest- & 23 Go-Tests abgedeckt ✅

### 🔄 Nächste Schritte
1. **Deployment auf k3s**
   - `kubectl create namespace restgeld`
   - SealedSecret für Postgres-Passwort erstellen
   - Flux Kustomization in k3s-config repo übernehmen
2. **Domain verifizieren** (restgeld.stillon.top)
3. **PWA-End-to-End-Test** im Produktions-Setup
