# Restgeld

Daily Allowance Tracker – Mobile-first PWA

[![Go Core](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml)
[![Auth Service](https://github.com/pipelinedave/restgeld/actions/workflows/auth-service.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/auth-service.yml)
[![Vue Frontend](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml)
[![Docker](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml)
[![GitHub Pages](https://github.com/pipelinedave/restgeld/actions/workflows/pages.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/pages.yml)
[![codecov](https://codecov.io/gh/pipelinedave/restgeld/branch/main/graph/badge.svg)](https://codecov.io/gh/pipelinedave/restgeld)

🌐 **Landing Page & Live Demo:** [https://pipelinedave.github.io/restgeld](https://pipelinedave.github.io/restgeld)

## Features

- Tägliches Budget auf Basis eines flexiblen Abrechnungszeitraums (z.B. Monatsbudget)
- Nicht ausgegebenes Budget rolliert automatisch auf verbleibende Tage
- Interaktiver Tages-Sparpuffer & Monatsende-Projektion ("Wo lande ich?")
- Gamification: 🔥 Spar-Streaks, Rekord-Levels & 🎯 Null-Euro-Tage Zähler
- Mobile-first PWA mit Schnell-Erfassung (Numpad, Quick Note Chips & Haptik)
- Offline-Outbox mit automatischer Hintergrund-Synchronisation
- Datenhoheit: CSV & JSON Export / Import
- Passkeys & Magic-Link Multi-Tenant Authentifizierung
- Ausgaben-Historie mit Filter- & Löschfunktion

## Tech Stack

| Layer | Technologie |
|---|---|
| Frontend | Vue 3 + Vite + TypeScript + PWA |
| Core Backend | Go 1.22 (REST API, net/http + lib/pq) |
| Auth Service | Standalone Go Microservice (Magic Links, Passkeys / WebAuthn, Sessions) |
| Database | PostgreSQL 16 (Multi-Tenant, Multi-DB `restgeld_core` & `restgeld_auth`) |
| Dev SMTP | Mailpit (In-Memory Mailer) |
| Deployment | Vercel (Preview & Production), Docker, Kubernetes (k3s), Flux |

## Umgebungen & Branching-Modell

- **Production (`main`):** Produktive Live-App für die tägliche Budgetverwaltung.
- **Preview Environment (`develop` / Feature-Branches):** Automatisches [Vercel Preview Deployment](https://vercel.com/docs/deployments/environments#preview-environment-pre-production) bei jedem Push/PR. Ermöglicht schnelles Testen und Iterieren neuer Features mit getrennter Preview-Datenbank, ohne die Produktivdaten zu gefährden.

## Lokale Entwicklung

### 1-Befehl-Start: Live-Entwicklung mit Hot-Reload (HMR)

Startet die Datenbank, das Backend und den Vite Dev-Server mit Live-Reloading & HMR bei Code-Änderungen:

```bash
# Mit Docker Compose
docker compose -f docker-compose.dev.yml up --build

# Oder via PowerShell (Windows)
.\scripts\dev.ps1

# Oder via Bash (Linux/macOS/WSL)
./scripts/dev.sh
```

- **Frontend (Live-HMR):** [http://localhost:5173](http://localhost:5173)
- **Backend API:** [http://localhost:8080](http://localhost:8080) (z.B. `/api/health`, `/api/budget`)
- **PostgreSQL:** `localhost:5432`

### Produktionsnaher Stack (Nginx + Static Build)

```bash
docker compose up --build
```
- **Frontend (Nginx):** [http://localhost:3000](http://localhost:3000)

### Manuelle Entwicklung (Einzelkomponenten)

```bash
# Backend
cd backend
go test ./... -v -short

# Backend mit Integrationstests (Postgres erforderlich)
go test -tags=integration ./... -v

# Frontend
cd frontend
npm install
npm run test:unit
npm run dev

# Docker einzeln
docker build -t restgeld-backend backend/
docker build -t restgeld-frontend frontend/
```

## Lizenz

MIT
