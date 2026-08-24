# Restgeld – Teststrategie

## Zielsetzung
Die App wird auf allen Ebenen getestet (ausser Last- und Performance-Tests).
Testlauf erfolgt automatisch via GitHub Actions bei jedem Push/PR auf `main`.

## Testpyramide

```
        ╱╲
       ╱  ╲
      ╱ E2E╲          ← manuell (DevTools) + optional Playwright
     ╱──────╲
    ╱Integration╲      ← Backend: postgresStore + echte DB
   ╱────────────╲      ← Frontend: useApi + mock-fetch
  ╱   Unit-Tests  ╲    ← Handler (httptest + memoryStore)
 ╱─────────────────╲   ← Komponenten (vitest + vue-test-utils)
╱─────────────────────╲
```

## Backend – Gesamtüberblick

| Test-Typ | Datei | Tag | DB? | Count |
|---|---|---|---|---|
| Handler-Unit | `handlers_test.go` | – | memoryStore | 15 |
| Store-Integration | `db_test.go` | `integration` | echte Postgres | 8 |
| API-Integration | `api_test.go` | `integration` | echte Postgres | 6 |
| Konkurrenz (Race) | `api_test.go` | `integration` | echte Postgres | 1 |

## Frontend – Gesamtüberblick

| Test-Typ | Datei(en) | Count |
|---|---|---|
| Komponenten-Unit | `MonthProgress.test.ts`, `BudgetDisplay.test.ts`, `RecentExpenses.test.ts`, `Numpad.test.ts` | 31 |
| Snapshot | `snapshots.test.ts` | 9 |
| App-Integration | `App.test.ts` | 7 |
| API-Mock | `useApi.test.ts` | 7 |
| E2E | `e2e/example.spec.ts` | 4 |

---

## 1. Backend – Handler-Unit-Tests

**Framework:** `testing` + `net/http/httptest`  
**Store:** `memoryStore` (in-memory, kein DB-Zugriff)  
**Dateien:** `handlers_test.go`, `store_test.go`

| Test | Szenario |
|---|---|
| `TestGetBudget` | Basis-Budget an Tag 17, keine Ausgaben |
| `TestGetBudgetWithExpenses` | Mit 2 Ausgaben, Ersparnis prüfen |
| `TestGetBudgetOnFirstDay` | Tag 1, Budget = baseBudget |
| `TestGetBudgetAfterAllSpent` | savings = 0 → weiß, currentBudget = baseBudget |
| `TestGetBudgetInDebt` | savings < 0 → rot |
| `TestCreateExpense` | POST gültige Ausgabe → 201 |
| `TestCreateExpenseInvalid` | amount=0 → 400 |
| `TestDeleteExpense` | DELETE existierende ID → 200 |
| `TestDeleteExpenseNotFound` | DELETE nicht-existente ID → 404 |
| `TestDeleteExpenseInvalidUUID` | DELETE ungültiges UUID-Format → 404 |
| `TestNewPeriod` | POST period → zurücksetzen, Ausgaben weg |
| `TestUpdateBudget` | PATCH monthlyTotal → 200, neuer Wert |
| `TestUpdateBudgetInvalid` | monthlyTotal=0 → 400 |
| `TestCORSHeaders` | OPTIONS → 200 + CORS-Header |
| `TestMethodNotAllowed` | PUT auf /api/budget → 405 |

**Ausführung:** `cd backend && go test ./... -v -count=1`

---

## 2. Backend – Store-Integrationstests (echte DB)

**Framework:** `testing` (Tag: `integration`)  
**Store:** `postgresStore` mit echter PostgreSQL-Verbindung  
**Dateien:** `db_test.go`

Postgres wird via GitHub Actions `services.postgres` gestartet.
Lokal: Docker-Container `restgeld-db-test` (siehe `scripts/test-db.sh`).

| Test | Szenario |
|---|---|
| `TestCreateAndReadPeriod` | Period anlegen und lesen |
| `TestCreateDuplicatePeriod` | ON CONFLICT Verhalten |
| `TestAddAndListExpenses` | Ausgaben hinzufügen, limit, Reihenfolge |
| `TestDeleteExpense` | Löschen mit und ohne Cascade |
| `TestGetTotalExpenses` | SUM über mehrere Ausgaben |
| `TestUpdateBudget` | Upsert monthlyTotal |
| `TestCreatePeriodClearsExpenses` | Neustart löscht alte Daten |

**Ausführung:**
```bash
# Lokal (Postgres-Container muss laufen)
go test -tags=integration ./... -v

# Oder mit Script
./scripts/test-db.sh
```

---

## 3. Frontend – Komponenten-Tests

**Framework:** `vitest` + `@vue/test-utils` + `happy-dom`  
**Dateien:** `src/components/*.test.ts`

### Numpad.test.ts
| Test | Beschreibung |
|---|---|
| Rendert Buttons | Alle 12 Ziffern/Steuer-Buttons sichtbar |
| Zifferneingabe | Klick 1 → Display "1" |
| Komma-Eingabe | Klick "," → Display "0," |
| Doppeltes Komma | Zweites "," wird ignoriert |
| Löschen | "12" + ⌫ → "1" |
| Bestätigen mit gültigem Betrag | confirm-Emit mit number |
| Bestätigen mit leerem Input | Kein Emit |
| Abbrechen | cancel-Emit + Input reset |
| Notiz-Textfeld nach Bestätigen | showNote = true |
| Notiz eingeben und speichern | confirm mit amount + note |
| Zurück vom Notiz-Modus | showNote = false |

### BudgetDisplay.test.ts
| Test | Beschreibung |
|---|---|
| Rendert Betrag | "14,52 €" |
| Farbe grün | savings > 0 → class color-green |
| Farbe weiß | savings == 0 → class color-white |
| Farbe rot | savings < 0 → class color-red |
| Ersparnis anzeigen | "+5,00 € angespart" |
| Keine Ersparnis bei 0 | savings == 0 → kein Text |

### MonthProgress.test.ts
| Test | Beschreibung |
|---|---|
| Label Tag X von Y | "Tag 17 von 31" |
| Fortschrittsbalken | 17/31 = 54% Breite |
| Erster Tag | 1/31 = 3% |
| Letzter Tag | 31/31 = 100% |

### RecentExpenses.test.ts
| Test | Beschreibung |
|---|---|
| Liste rendern | 3 Ausgaben anzeigen |
| Leere Liste | "Noch keine Ausgaben heute" |
| Löschen-Button | Klick → delete-Emit mit ID |
| Notiz anzeigen | note-Text statt "Ausgabe" |
| Betrag formatieren | "-8,50 €" |

### useApi.test.ts
| Test | Beschreibung |
|---|---|
| getBudget | Fetch GET /api/budget → BudgetData |
| addExpense | Fetch POST /api/expenses → Expense |
| deleteExpense | Fetch DELETE /api/expenses/{id} |
| newPeriod | Fetch POST /api/period |
| updateBudget | Fetch PATCH /api/budget |
| API-Fehler | fetch.ok=false → throw Error |

**Ausführung:** `cd frontend && npm run test:unit`

---

## 4. CI/CD – GitHub Actions

Siehe `.github/workflows/`.

### workflow: backend.yml
```yaml
on: [push, pull_request]
paths: ['backend/**']
jobs:
  test:
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: restgeld
          POSTGRES_PASSWORD: restgeld
          POSTGRES_DB: restgeld
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go mod download
      - run: go vet ./...
      - run: go test -tags=integration ./... -v -coverprofile=coverage.out
      - run: go tool cover -func=coverage.out
```

### workflow: frontend.yml
```yaml
on: [push, pull_request]
paths: ['frontend/**']
jobs:
  test:
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22' }
      - run: npm ci
      - run: npm run build        # vue-tsc + vite build
      - run: npm run test:unit     # vitest
```

### workflow: docker.yml
```yaml
on: [push, pull_request]
paths: ['backend/**', 'frontend/**']
jobs:
  build:
    steps:
      - uses: actions/checkout@v4
      - run: docker build backend/
      - run: docker build frontend/
```

---

## 5. Backend – API-Integrationstests

**Framework:** `testing` (Tag: `integration`)  
**Datei:** `api_test.go`

Echte HTTP-Requests durch den ganzen Stack (Handler → store → Postgres).

| Test | Szenario |
|---|---|
| `TestAPI_GetBudget_CreatesPeriod` | GET erzeugt Periode bei Leerstand |
| `TestAPI_FullFlow` | Budget → Ausgabe → Budget ↓ → Löschen → Budget = Initial |
| `TestAPI_ConcurrentExpenses` | 3 parallele POST → alle 3 gespeichert |
| `TestAPI_UpdateBudgetAndReset` | PATCH Budget → POST period → Budget zurückgesetzt |
| `TestAPI_DeleteNonexistent` | DELETE ungültige UUID → 404 |
| `TestAPI_InvalidAmount` | POST mit negativem/0 amount → 400 |

---

## 6. Frontend – Snapshot-Tests

**Framework:** `vitest` + `expect().toMatchSnapshot()`  
**Datei:** `snapshots.test.ts`

Erfasst HTML-Struktur jeder Komponente in verschiedenen Zuständen.  
Bei unbeabsichtigten Änderungen schlagen Tests an – Update via `vitest --update`.

| Komponente | Snapshots |
|---|---|
| MonthProgress | Tag 17/31, Tag 1/28 |
| BudgetDisplay | grün, weiss, rot |
| RecentExpenses | leer, mit Ausgaben |
| Numpad | versteckt, sichtbar |

---

## 7. Frontend – App-Integrationstests

**Framework:** `vitest` + `@vue/test-utils`  
**Datei:** `App.test.ts`  
**Mock:** `useApi` via `vi.mock`

| Test | Szenario |
|---|---|
| Budget laden beim Mount | getBudget aufgerufen, Budget sichtbar |
| Budget und Ersparnis anzeigen | Tag, Betrag, Savings im DOM |
| Numpad öffnen | Klick "Ausgabe" → Overlay sichtbar |
| Ausgabe buchen | Numpad-Eingabe → addExpense aufgerufen |
| Ausgabe löschen | Klick delete-btn → deleteExpense aufgerufen |
| Ladestatus | Initial "Lade..." sichtbar |
| API-Fehler | console.error bei Fehlschlag |

---

## 8. Frontend – E2E-Tests (Playwright)

**Framework:** `@playwright/test`  
**Datei:** `e2e/example.spec.ts`  
**Setup:** `playwright.config.ts` startet `npm run preview` als webServer

Laufen im CI via `npm run test:e2e` (separater Workflow oder manuell).

| Test | Szenario |
|---|---|
| Seite lädt | Budget-Zahl und Ersparnis sichtbar |
| Numpad öffnen | Button-Klick → Bestätigen-Button sichtbar |
| Ausgabe buchen | Zifferneingabe → Notiz → Speichern → Liste aktualisiert |
| Ausgabe löschen | Anlegen → Löschen → verschwunden |

```bash
# Lokal (erfordert laufendes Backend + Postgres)
cd frontend && npx playwright install chromium && npm run test:e2e
```

---

## 9. Coverage

- **Backend:** `go test -coverprofile=coverage.out` → `go tool cover -func=coverage.out`
- **Frontend:** `vitest --coverage` → `coverage/lcov.info`

Coverage-Badges via **Codecov**:
- README.md: `[![codecov](https://codecov.io/gh/pipelinedave/restgeld/branch/main/graph/badge.svg)](https://codecov.io/gh/pipelinedave/restgeld)`

---

## 10. CI/CD – GitHub Actions Workflows

### backend.yml
```yaml
- go vet ./...
- go test -short ./...          # Unit-Tests (15)
- go test -tags=integration ./... # Integration-Tests (14)
- codecov/codecov-action
```

### frontend.yml
```yaml
- npm run build                  # vue-tsc + vite build
- npm run test:unit               # vitest (51 Tests)
- codecov/codecov-action
```

---

## 11. Lokale Ausführung

```bash
# Backend Unit
cd backend && go test -short ./... -v -count=1

# Backend Integration (Postgres erforderlich)
docker run -d --name restgeld-db-test -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -p 5433:5432 postgres:16-alpine
DB_PORT=5433 go test -tags=integration ./... -v

# Frontend Unit + Snapshot
cd frontend && npm run test:unit

# Frontend E2E (Playwright, erfordert Backend)
cd frontend && npx playwright install chromium && npm run test:e2e

# Alles
./scripts/test-all.sh  # TODO
```

---

## 12. Ausgeschlossen (vorerst)
- Last- / Performance-Tests
- Security-Scans (Trivy/Snyk) – optional
