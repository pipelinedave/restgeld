# Restgeld – Daily Allowance Tracker (PWA)

## Projektbeschreibung
Minimalistische Web-App (Mobile-first PWA) zum Tracken eines täglichen "Spending Money" Budgets.
Betrieben auf lokalem single-node k3s-Cluster.

## Architektur (3-Tier)
1. **Frontend:** Vue 3 + Vite + TypeScript + PWA (Service Worker via vite-plugin-pwa), Nginx
2. **Backend:** Go REST API (net/http + lib/pq), port 8080
3. **Datenbank:** PostgreSQL 16-alpine mit PVC

## Domain & Infrastruktur
- **Domain:** restgeld.stillon.top
- **Ingress:** nginx-ingress, TLS via letsencrypt-prod (cert-manager)
- **Namespace:** restgeld (manuell via kubectl create namespace restgeld)
- **Flux Kustomization:** apps/restgeld.yaml → kustomize/restgeld/
- **Secrets:** SealedSecret (bitnami.com/v1alpha1)

## API-Endpunkte

| Methode | Pfad | Funktion |
|---|---|---|
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
│   │       ├── MonthProgress.vue
│   │       ├── BudgetDisplay.vue
│   │       ├── Numpad.vue
│   │       └── RecentExpenses.vue
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
- AGENTS.md erstellt (Workflow-Definition)
- CONTEXT.md erstellt (Plan-Dokument)

### 🔄 Nächste Schritte
1. **Backend-Implementierung** (models.go, db.go, handlers.go, main.go, go.mod, Dockerfile)
2. **Backend-Tests** (httptest-Handler-Tests)
3. **Frontend-Implementierung** (Vue 3 Komponenten + PWA)
4. **Docker-Build-Test** (beide Images)
5. **Lokaler Testlauf** (Postgres + Backend + Frontend + DevTools)
6. **Kubernetes-Manifeste** (k3s-config-ready)
7. **Deployment** (k3s + Flux + Domain verifizieren)
