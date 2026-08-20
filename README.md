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

# Docker
docker build -t restgeld-backend backend/
docker build -t restgeld-frontend frontend/
```

## Lizenz

MIT
