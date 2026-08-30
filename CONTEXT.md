# Restgeld – Daily Allowance Tracker (PWA)

## Projektbeschreibung
Minimalistische Web-App (Mobile-first PWA) zum Tracken eines täglichen "Spending Money" Budgets.
Betrieben auf lokalem single-node k3s-Cluster.

## Architektur (3-Tier)
1. **Frontend:** Vue 3 + Vite + TypeScript + PWA (Service Worker via vite-plugin-pwa), Nginx
2. **Backend:** Go REST API (net/http + lib/pq), port 8080
3. **Datenbank:** PostgreSQL 16-alpine mit PVC

## Domain, Umgebungen & Infrastruktur
- **Environments & Branching:**
  - **Production (`main`):** Produktive Live-App für tägliche Nutzung.
  - **Preview Environment (`develop` / Feature-Branches):** Automatisches Pre-Production Deployment auf Vercel pro Push/Branch für schnelle Iterationen und sichere Feature-Verifikation ohne Gefahr für Live-Daten.
  - **Environment Separation:** Getrennte DB-Instanzen / Environment-Variablen zwischen Production und Preview.
- **Production Host:** Vercel & lokaler k3s-Cluster (restgeld.stillon.top)
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
- CSV & JSON Daten Export / Import: Endpoints `GET /api/export` (CSV & JSON-Backup mit automatischem Datei-Download) und `POST /api/import` (Wiederherstellung aus CSV/JSON) sowie neuer Bereich "Daten & Backup" in `SettingsModal.vue` implementiert und mit 80 Vitest- & 27 Go-Tests abgesichert ✅
- 🔥 Spar-Streaks & Gamification: Streak-Berechnung im Backend (`calcStreakInfo`), `StreakCard.vue` mit animiertem 🔥 Flammen-Badge, Rekord-Streak und 🎯 Null-Euro-Tage Zähler implementiert und mit 83 Vitest- & 28 Go-Tests abgedeckt ✅
- Numpad Live-Impact & Quick Note Chips: Dynamische Echtzeit-Berechnung des verbleibenden Tagesbudgets bzw. Puffer-Abzugs direkt bei der Betragseingabe (`impact-ok` / `impact-warning`) sowie One-Tap Quick-Chips für die letzten Notizen mit localStorage-Persistierung implementiert und mit 85 Vitest-Tests abgedeckt ✅
- 📱 PWA Quick Actions & Offline-Outbox: Homescreen App-Shortcut (`/?action=add-expense`), Outbox-Warteschlange (`useOfflineSync`), automatischer Hintergrund-Sync bei Wiederverbindung, optimistische Dashboard-Aktualisierung und Offline-Indikator im App-Header implementiert und mit 91 Vitest-Tests abgesichert ✅
- 🔮 Monatsende-Projektion ("Wo lande ich?"): Backend-Berechnung `calcProjection` (`projectedSavings`, `projectedTotalSpent`, `avgDailySpend`, `status`), minimalistischer Dashboard-Strip `MonthProjection.vue` mit Ersparnis-/Defizit-Vorschau und Ø Tagesschnitt implementiert und mit 94 Vitest- & 29 Go-Tests abgedeckt ✅
- 📜 Vorherige Perioden & Monats-Rückblick: Backend-Endpoint `GET /api/periods` mit aggregierten Perioden-Metriken (`totalSpent`, `savings`, `expenseCount`), Archiv-Modal `PeriodsArchiveModal.vue` und Aufruf in `SettingsModal.vue` implementiert und mit 98 Vitest- & 30 Go-Tests abgesichert ✅
- Preview Environment & Branching-Strategie: Vercel Preview Deployments und Pre-Production Branch `develop` eingerichtet, CI/CD-Pipelines & Agent-Workflow-Dokumentation (`AGENTS.md`, `CONTEXT.md`, `README.md`) für sichere Iteration ohne Gefahr für Live-Daten aktualisiert ✅
- 🌐 Landing Page & GitHub Pages Deployment: Responsive, eigenständige Landing Page in `docs/` mit interaktivem Live-Mockup, dynamischem Budget-Rechner, Bento-Feature-Grid, FAQ und GitHub Actions Workflow (`pages.yml`) für automatisches GitHub Pages Hosting implementiert ✅
- 🎨 UI Redesign nach Landing-Page Vorbild: Vollständige Übernahme des modernen OLED Dark Mode Designs (`#0a0a0c`, Plus Jakarta Sans, JetBrains Mono, Emerald Green `#22c55e`, abgerundete Hero- & Progress-Cards, Glow-Border und Glassmorphism-Modals) in der Vue 3 PWA ✅
- 🚀 Mobile-UX & Zero-Scroll Dashboard (Epic 1 & 2): Numpad mit initialem "Heute verfügbar" Banner, Android/iOS Autofocus-Optimierung, Hero-Display "x / y noch übrig", Top Loading-Streak unter dem Header, scroll-freies `100dvh` Viewport-Layout und Git-Commit Verlinkung im Footer implementiert und mit 98 Vitest- & 30 Go-Tests abgedeckt ✅
- 🎛️ Interaktive Settings-Slider & Abschlussbericht-Archiv (Epic 3): Granulare 1 €-Budget- & Tages-Slider mit bidirektionalem Live-Kalkulator in `SettingsModal.vue`, `GET /api/expenses?period_id=...` im Backend und interaktiver Monats-Abschlussbericht mit Ausgabenliste in `PeriodsArchiveModal.vue` implementiert und mit 99 Vitest- & 31 Go-Tests abgesichert ✅
- 🎨 Custom Theming, Health Popover & About Page (Epic 4): Dynamische OLED-Farbwelten mit Color-Picker (`useTheme`), interaktives Service-Health Popover in `AppHeader.vue` (API Latenz, DB-Status, Offline-Queue), About-Modal `AboutModal.vue` und Phone-Mockup Sync auf der Landing Page implementiert und mit 106 Vitest- & 31 Go-Tests abgedeckt ✅
- 🔮 Header Popover Architecture & UX Refinements: Monatsende-Prognose (🔮) und Spar-Streak (🔥) als elegante Popover in den Header verlagert (100% Zero-Scroll freigespielt), Hero-Anzeige präzisiert auf echten Tages-Bruch (`x € / y €`), Footer mit leuchtender Commit-Badge & `lowlifehigh.tech`-Energy-Referenz verfeinert und mit 109 Vitest- & 31 Go-Tests abgesichert ✅
- 🎨 Logo- & Favicon-Redesign: Modernes 'r.'-Logo mit dunkelgrauem Squircle-Hintergrund, off-weißem Text und smaragdgrünem Akzentpunkt als Vektorgrafik (SVG) für App-Favicon und Landing Page implementiert ✅
- 📜 Scrollbare Ausgaben-Liste & Zero-Scroll Viewport: Hauptseiten-Viewport () strikt auf Zero-Scroll fixiert (), Ausgaben-Liste in  als eigenständig scrollbares Element mit custom OLED-Scrollbar implementiert und mit allen 123 Vitest-Tests verifiziert ✅
- ⚡ SaaS Microservice Architecture & Auth Decoupling (Phase 1-3): Dedizierter Standalone-Microservice `services/auth-service` (Port 8081) mit eigener Postgres-DB (`restgeld_auth`), Mailpit SMTP-Versand, K8s Ingress-Routing & Vite Proxy-Regeln extrahiert. Legacy Auth-Handler im Monolyt-Core vollständig bereinigt und mit 123 Vitest- & 42 Go-Tests sowie Vercel Preview Verifikation abgesichert ✅
- 🧪 Vitest Frontend Test-Suite Fix: Vollständige Bereinigung aller 23 Test-Suites (123/123 Unit- & Komponententests grün), Numpad Live Impact-Vorschau String-Matching und Vue 3 SettingsModal Reaktivitäts-Handhabung stabilisiert ✅
- 🌍 SaaS Multi-Language Support (i18n): Leichtgewichtiger Vue 3 Composable `useI18n.ts` (DE 🇩🇪, EN 🇬🇧, ES 🇪🇸, FR 🇫🇷), Sprachtranslationen über alle UI-Komponenten hinweg, Browser-Erkennung, `localStorage`-Sync, DB-Migration `003_add_language_to_users.sql` im Auth-Microservice, sprachspezifische Magic Link E-Mails, `Accept-Language` Header-Handling und Landing Page Sprachauswahl implementiert und mit 128 Vitest- & 42 Go-Tests abgesichert ✅
- 💱 SaaS Multi-Currency & Zero-Bloat Smart Categorization (Phase 4): Vollständige Multi-Currency Engine (`EUR €`, `USD $`, `GBP £`, `CHF`, `JPY ¥`) mit dynamischem Symbol/Format und Settings-Auswahl, DB-Migration `004_add_currency_to_users.sql` & Profil-Persistierung, 100% flächendeckende Übersetzung aller Modals/Popovers/Toasts sowie automatische Zero-Bloat Emoji-Kategorisierung (☕, 🍔, 🛒, 🚗, 🎉, 📦, 💊) ohne jeglichen Konfigurationsaufwand implementiert und mit 129 Vitest- & 42 Go-Tests abgesichert ✅
- 🔐 Native WebAuthn Passkeys & Smart Category Insights (Phase 5): Native Biometrie-Authentifizierung (Touch ID, Face ID, Windows Hello) mit `/api/auth/passkey/register-verify` & `/api/auth/passkey/login-verify` in Go und `useAuth.ts`, 1-Tap Passkey-Login in `AuthModal.vue`, On-the-fly Category Insights Breakdown im Monatsabschluss (`PeriodsArchiveModal.vue`), Smart-Quick-Chips in `Numpad.vue` und Multi-Currency Landing Page Upgrade (`docs/index.html`, `docs/script.js`, `README.md`) implementiert und mit 130 Vitest- & 45 Go-Tests abgesichert ✅
- 📊 SaaS Observability & Monitoring Microservice (Phase 6): Standalone Go-Microservice `services/monitoring-service` (Port 8083) mit `/metrics` (Prometheus-Scraping), `/api/monitoring/overview` (Cluster-Health & p95 Latenz-Telemetrie) und interaktivem Echtzeit-Dashboard `SystemStatusModal.vue` im Frontend (inkl. Auto-Refresh, Service-Gauges, RAM/Goroutine-Telemetrie und 4-Sprachen-i18n) implementiert und mit 135 Vitest- & 48 Go-Tests abgesichert ✅
- ⚡ PWA Shortcuts, Deep Action Routing & Single-Period CSV Export (Phase 7): PWA Manifest mit `restgeld.` Branding & 3 Quick-Action Shortcuts (`/?action=add-expense`, `status`, `archive`), globale Tastatur-Steuerung (`Space`/`+`/`N`, `ESC`, `M`, `S`, `A`, `?`), 1-Klick Perioden-CSV-Export im Archiv und Keyboard-Dokumentation im Über-Modal implementiert und mit 136 Vitest- & 48 Go-Tests abgesichert ✅
- 🎛️ Hero Spending Gauge, Synthetic Web Audio Feedback & Sound Options (Phase 8): Dynamische Mini-Progress Gauge Bar im Hero-Card-Display (`BudgetDisplay.vue`) mit farbkodierter Tagesauslastung, Web Audio API Synthesizer in `useHaptics.ts` (Click-, Chime- & Warning-Töne ohne externe Audio-Dateien), Sound-Option Toggle in `SettingsModal.vue` und 4-Sprachen-i18n implementiert und mit allen 136 Vitest- & 48 Go-Tests abgesichert ✅

---

## 📋 Master Roadmap & Feedback Backlog

Strukturierte Erfassung der Praxistest-Erkenntnisse und Feature-Wünsche für kommende Entwicklungs-Zyklen.

### Epic 1: Mobile- & Numpad-UX Polishing
1. **Android Quick Action Auto-Focus & Virtual Keyboard Trigger**:
   - Problem: Beim Öffnen via PWA Shortcut (`/?action=add-expense`) öffnet sich das Android-Keyboard nicht sofort automatisch.
   - Ziel: Zuverlässiges Fokussieren und Aufklappen der Bildschirmtastatur auf Android & iOS.
2. **Initialer "Heute verfügbar"-Status im Numpad**:
   - Problem: Die Live-Impact Box erscheint aktuell erst, nachdem eine Ziffer eingetippt wurde.
   - Ziel: Bereits vor Eingabe den aktuellen Stand anzeigen (z. B. *"Heute verfügbar: 15,00 €"*), der sich beim Tippen nahtlos in *"Verbleibt danach: 6,50 €"* wandelt.
3. **Hero-Display: "x / y noch übrig"**:
   - Ziel: Neben dem Restbetrag klar anzeigen, wie viel heute ursprünglich zur Verfügung stand (z. B. `12,50 € / 15,00 €` oder Untertitel *"von 15,00 € heute"*).

### Epic 2: Main Dashboard Layout & Zero-Scroll Information Architecture
1. **Zero-Scroll Viewport Konzept (100dvh)**:
   - Problem: Die Hauptseite ist auf vielen Smartphones vertikal leicht überschritten und erfordert Scrollen zum Footer.
   - Ziel: Absolut scroll-freie Hauptseite (`100dvh`). Intelligente Informationsarchitektur:
     - Kompakte, aufklappbare Widgets oder Tabs (z. B. Streak & Projection als elegante Pill-Widgets).
     - Ausgaben-Vorschau als elegante Drawer-/Footer-Leiste oder Schnellansicht.
     - Maximaler Fokus auf Hero-Restgeld und Schnellbuchung.
2. **Instant-Operations & Top Loading-Streak**:
   - Ziel: Alle Aktionen passieren optimistisch und instant.
   - Dezente, pulsierende Ladeleiste ("Loading Streak") direkt unter dem Header für alle asynchronen Background-Syncs.
3. **Git Commit im Footer**:
   - Ziel: Statt einer manuellen Versionsnummer den aktuellen Git-Commit-Hash (z. B. `2e783da`) mit Link zum GitHub-Commit anzeigen.

### Epic 3: Interaktive Settings & Abschlussbericht-Archiv
1. **Interaktive Slider für Budget & Periode (wie Landing Page)**:
   - Ziel: Slider für Monatsbudget, Periodendauer und Tages-Restgeld mit bidirektionaler, intelligenter Kopplung (z. B. Verschieben des Tagesbudgets passt das Monatsbudget live an).
2. **Überarbeitetes Perioden-Archiv ("Abschlussbericht")**:
   - Problem: Periodennamen sind teilweise kryptisch und genaue Zeitspannen fehlen.
   - Ziel: Klare Anzeige *"1. Aug 2026 – 31. Aug 2026 (31 Tage)"*.
   - Detail-Ansicht beim Antippen einer Periode: Kompletter Abschlussbericht (End-Ersparnis, Gesamtausgaben, Tagesdurchschnitt, Streak-Statistiken und Ausgabenliste des Monats).

### Epic 4: Theming, Monitoring & Pages
1. **Custom Color / RGB & Hex Theming**:
   - Ziel: Farbwähler für Akzentfarben (Emerald, Cyan, Sunset Amber, Cyberpunk) und OLED-Theme-Optionen mit CSS-Variablen & `localStorage`.
2. **Uptime Monitoring & Health Popover**:
   - Ziel: Uptime Kuma Status / Service-Health hinter dem Online-Badge im Header. Antippen öffnet ein Popover mit Status aller Komponenten (DB, API, Sync-Queue).
3. **About Page / Modal**:
   - Ziel: Kompakte Info-Seite zur Philosophie von Restgeld (Achtsamkeit, Zero-Bloat, Open Source).
4. **Landing Page Sync**:
   - Ziel: Aktualisierung des interaktiven Mockups in `docs/index.html` auf das exakte Layout der PWA.

### Epic 5: SaaS Architecture, Microservices & Subscriptions (✅ Umgesetzt)
1. **User Accounts & Tenant Isolation**:
   - Multi-Tenant DB Schema mit Tabellen `users`, `magic_links`, `auth_sessions`, `webauthn_credentials` und `user_id` Foreign Keys.
   - Vollständige Tenant-Isolation auf DB- und Store-Ebene.
2. **Microservice Isolation**:
   - `services/auth-service` (Port 8081): Standalone Auth Microservice (Magic Link & WebAuthn Passkey Vorbereitung).
   - `services/billing-service` (Port 8082): Standalone Billing Microservice (Stripe Checkout & Customer Portal).
3. **Flächendeckendes Rate Limiting & Abuse Protection**:
   - `rateLimitMiddleware` über alle Microservices hinweg (`backend/`, `auth-service`, `billing-service`).
4. **Frontend Billing & Health UI Integration**:
   - `useBilling.ts` Composable & Stripe Checkout / Portal Buttons in `AuthModal.vue`.
   - Microservice Health-Checks (`auth-service` & `billing-service`) im Header Popover.

---

### 🔄 Nächste Schritte
1. **Deployment auf k3s**:
   - SealedSecret für Postgres-Passwort in Produktion verifizieren.
   - Flux Kustomization Sync in k3s Cluster ausführen.
2. **Domain verifizieren** (`restgeld.stillon.top`)

