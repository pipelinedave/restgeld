# Restgeld

Daily Allowance Tracker – Mobile-first PWA

[![Go](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml)
[![Vue](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml)
[![Docker](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml)
[![codecov](https://codecov.io/gh/pipelinedave/restgeld/branch/main/graph/badge.svg)](https://codecov.io/gh/pipelinedave/restgeld)

## Features

- Tägliches Budget auf Basis eines Monatsbudgets
- Nicht ausgegebenes Budget rolliert auf verbleibende Tage
- Gamification: Ersparnis-Anzeige
- Mobile-first PWA mit In-App-Numpad (keine native Tastatur)
- Ausgaben-Historie mit Löschfunktion

## Tech Stack

| Layer | Technologie |
|---|---|
| Frontend | Vue 3 + Vite + TypeScript + PWA |
| Backend | Go 1.22 (net/http + lib/pq) |
| Database | PostgreSQL 16 |
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
