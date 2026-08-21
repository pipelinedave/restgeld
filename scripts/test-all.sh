#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== Restgeld Test Suite ==="
echo ""

# 1. Backend Unit-Tests
echo "--- Backend Unit-Tests ---"
cd "$ROOT/backend"
go test -short ./... -v -count=1
echo ""

# 2. Frontend Unit-Tests
echo "--- Frontend Unit-Tests ---"
cd "$ROOT/frontend"
npm run test:unit
echo ""

# 3. Frontend Build (vue-tsc + vite)
echo "--- Frontend Build ---"
cd "$ROOT/frontend"
npm run build
echo ""

# 4. Docker Build (optional)
if docker info > /dev/null 2>&1; then
    echo "--- Docker Build ---"
    cd "$ROOT"
    docker build -t restgeld-backend:test backend/
    docker build -t restgeld-frontend:test frontend/
else
    echo "--- Docker Build skipped ---"
fi

echo ""
echo "=== Alle Tests bestanden ==="
