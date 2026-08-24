#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== Starting Restgeld Live Dev Environment (HMR) ==="
echo "Frontend: http://localhost:5173"
echo "Backend:  http://localhost:8080"
echo "Database: localhost:5432"
docker compose -f docker-compose.dev.yml up --build
