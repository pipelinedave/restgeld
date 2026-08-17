# Restgeld – Agent Workflow

## Rolle
Du bist ein autonomer Full-Stack Agent, der die Restgeld-App entwickelt.
Du startest das Projekt von Null, schreibst Code, deployst auf k3s und iterierst.

## Workflow (Zyklen)

Jeder Zyklus folgt diesem Ablauf:

### 1. Plan prüfen
- `CONTEXT.md` lesen → aktuellen Projektstatus verstehen
- Nächsten unerledigten Schritt identifizieren

### 2. Git-Clean-Check
- `git status` prüfen
- Wenn uncommitted changes und kein Zyklus aktiv: user fragen ob commit/stash
- Vor Änderungen: `git pull --rebase` (falls remote existiert)

### 3. Implementieren
- Code schreiben gemäss Plan
- Maximal eine logische Änderung pro Zyklus (SRP)
- Nach jeder Änderung: `git add` relevante Dateien

### 4. Testen (mehrstufig)

#### 4a. Backend-Tests
```bash
cd backend && go test ./... -v
```
Wenn keine Tests existieren: `httptest`-basierte Handler-Tests schreiben.

#### 4b. Frontend-Build
```bash
cd frontend && npm run build
```
Keine Fehler erlaubt.

#### 4c. Docker-Build
```bash
docker build -t restgeld-backend:test backend/
docker build -t restgeld-frontend:test frontend/
```

#### 4d. Lokaler Testlauf (required vor Deploy)
1. Postgres starten: `docker run -d --name restgeld-db -e POSTGRES_USER=restgeld -e POSTGRES_PASSWORD=restgeld -e POSTGRES_DB=restgeld -p 5432:5432 postgres:16-alpine`
2. Backend starten: `docker run --rm --network host restgeld-backend:test`
3. Frontend starten: `docker run --rm -p 8080:80 restgeld-frontend:test`
4. Chrome DevTools (wsl-devtools) nutzen:
   - `wsl-devtools_navigate_page` → App URL
   - `wsl-devtools_take_snapshot` (a11y tree, verbose)
   - Screenshots NUR wenn Snapshot nicht ausreicht
   - Network-Requests prüfen (`wsl-devtools_list_network_requests`)
   - Console auf Fehler checken (`wsl-devtools_list_console_messages`)

#### 4e. Backend-API manuell testen
```bash
curl -s http://localhost:8080/api/budget | jq .
curl -s -X POST http://localhost:8080/api/expenses -H 'Content-Type: application/json' -d '{"amount": 5.50, "note": "Test"}'
curl -s http://localhost:8080/api/expenses | jq .
```

### 5. Fehler beheben
- Bei Test-Fail: Stacktrace analysieren → fix → zurück zu Schritt 4
- Bei DevTools-Fehlern: Console/Network-Tab auswerten → fix → zurück zu Schritt 4

### 6. Commit (Entwicklerlog)
- Commit nach **jeder logischen Einheit** (Backend, Frontend, K8s, Config)
- Nicht warten bis alles fertig ist – Git-Log = Entwicklungstagebuch
- Voraussetzung: Die Einheit ist für sich funktionsfähig (kompiliert, Tests grün)
- Commit-Message: `git commit -m "restgeld: <scope> - <kurzbeschreibung>"`
- Scope: `backend`, `frontend`, `k8s`, `docker`, `docs`

### 7. Plan aktualisieren
- `CONTEXT.md` lesen, `## Status`-Sektion updaten
- Erledigte Schritte markieren, neuen Status dokumentieren

## Konventionen

### Code
- Sprache: **DE** (Logs, Kommentare, Commit-Messages)
- Backend: Go, `net/http` + `lib/pq`, keine Frameworks
- Frontend: Vue 3 + Vite + TypeScript, keine UI-Library
- Mobile-first CSS, keine Tailwind/Bootstrap

### Secrets
- Postgres-Passwort: erstes Deployment via kubeseal
- `kubectl create secret generic restgeld-secrets -n restgeld ... | kubeseal --format yaml > kustomize/restgeld/restgeld-secrets-sealed.yaml`

### Debugging
- Bei Frontend-Problemen: `wsl-devtools` Tools nutzen
- Bei Backend-Problemen: `debugmcp` Breakpoints setzen
- Nie native Tastatur im Numpad-Kontext verwenden

## Projekt-Root
`/home/dhallmann/projects/restgeld`
