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
| Deployment | Docker, Kubernetes (k3s), Flux |

## Lokale Entwicklung

### 1-Befehl-Start (Gesamter Stack)

Der gesamte Stack (PostgreSQL, Backend-API und Frontend) kann mit einem einzigen Befehl gestartet werden:

```bash
# Mit Docker Compose
docker compose up --build

# Oder via PowerShell (Windows)
.\scripts\dev.ps1

# Oder via Bash (Linux/macOS/WSL)
./scripts/dev.sh
```

- **Frontend:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8080](http://localhost:8080) (z.B. `/api/health`, `/api/budget`)
- **PostgreSQL:** `localhost:5432`

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
