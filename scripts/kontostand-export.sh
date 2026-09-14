#!/usr/bin/env bash
# Kontostand-Export (JSON) — RESTGELD, NUR READONLY
# Ruft GET /api/export (JSON) auf und zeigt den aktuellen Kontostand-Export als JSON.
#
# Nutzung:
#   ./scripts/kontostand-export.sh                     # gegen http://localhost:8080
#   BASE_URL=... ./scripts/kontostand-export.sh
#   RESTGELD_TOKEN=... ./scripts/kontostand-export.sh  # per Bearer-Token
#   RESTGELD_COOKIE=... ./scripts/kontostand-export.sh # per Session-Cookie (restgeld_session)
#
# Ausgabe (readonly, keine Mutation):
#   - Volles JSON-Backup (Kontostand / Periode / Ausgaben) als JSON-Dokument.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BASE_URL="${BASE_URL:-http://localhost:8080}"
TOKEN="${RESTGELD_TOKEN:-}"
COOKIE="${RESTGELD_COOKIE:-}"

url="${BASE_URL}/api/export"

args=(-sS --fail)
if [ -n "$TOKEN" ]; then
  args+=( -H "Authorization: Bearer ${TOKEN}" )
elif [ -n "$COOKIE" ]; then
  args+=( -H "Cookie: restgeld_session=${COOKIE}" )
fi

echo "Kontostand-Export (JSON) von ${BASE_URL}" >&2
echo "----------------------------------------" >&2

# JSON-Export anfordern; Ausgabe direkt als JSON-Dokument
curl "${args[@]}" "$url"
echo
