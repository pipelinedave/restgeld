# Restgeld – SaaS Daily Allowance Tracker

**Restgeld** ist ein minimalistischer, hochperformanter **Multi-Tenant SaaS Daily Allowance Tracker** als modern mobile-first Progressive Web App (PWA).

[![Go Core](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/backend.yml)
[![Auth Service](https://github.com/pipelinedave/restgeld/actions/workflows/auth-service.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/auth-service.yml)
[![Vue Frontend](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/frontend.yml)
[![Docker](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/docker.yml)
[![GitHub Pages](https://github.com/pipelinedave/restgeld/actions/workflows/pages.yml/badge.svg)](https://github.com/pipelinedave/restgeld/actions/workflows/pages.yml)
[![codecov](https://codecov.io/gh/pipelinedave/restgeld/branch/main/graph/badge.svg)](https://codecov.io/gh/pipelinedave/restgeld)

🌐 **Landing Page & Live Demo:** [https://pipelinedave.github.io/restgeld](https://pipelinedave.github.io/restgeld)

---

## 🚀 Key SaaS Features

- 🌍 **Multi-Language Engine (i18n)**: Nahtloser Sprachwechsel in Echtzeit für **Deutsch 🇩🇪, Englisch 🇬🇧, Spanisch 🇪🇸 und Französisch 🇫🇷** (inkl. automatischer Browser-Spracherkennung und lokalisierter Magic-Link-E-Mails).
- 🔐 **Passwortlose Authentifizierung**: Standalone Auth-Microservice mit **Passkeys / WebAuthn** und **Magic-Links** für blitzschnellen, sicheren Login ohne Passwörter.
- ⚡ **Multi-Tenant Cloud Sync**: Mandantentrennende Datenarchitektur über Postgres-Schlüssel für synchronisierte Nutzung über Smartphone, Tablet und Desktop hinweg.
- 📊 **Tägliche Restgeld-Kalkulation**: Statt unübersichtlicher 40-Kategorien-Tabellen berechnet Restgeld jeden Tag präzise eine Kennzahl: *Dein verbleibendes Tagesbudget*.
- 🔄 **Automatischer Sparpuffer-Rollover**: Ersparnisse fließen direkt in den Folgetages-Puffer; Überziehungen werden sanft über die Restlaufzeit ausgedehnt.
- 🔥 **Gamification & Spar-Streaks**: Interaktive Streaks, Spar-Puffer-Rechner und Monatsende-Projektionen (🔮) in der Header-Leiste.
- 📲 **Zero-Scroll Mobile PWA**: Optimiert auf `100dvh` Viewport ohne vertikales Scrollen auf der Hauptseite, mit Haptik-Feedback, Numpad-Buchen und Offline-Outbox.
- 🔓 **Full Data Ownership**: Lokale Offline-Nutzung möglich mit vollständigem CSV & JSON Export / Import.

---

## 🛠️ Microservice Tech Stack

| Layer | Technologie & Architektur |
|---|---|
| **Frontend PWA** | Vue 3, Vite, TypeScript, `useI18n` Composable, Mobile-First Zero-Scroll CSS |
| **Core API Backend** | Go 1.22 REST Microservice (`net/http`, high-throughput, native SQL `lib/pq`) |
| **Auth Microservice** | Go Standalone Service (Magic Links, Passkeys / WebAuthn, OAuth/JWT Sessions) |
| **Billing Service** | Go Standalone Subscriptions & Quota Microservice |
| **Datenbank Layer** | PostgreSQL 16 (Multi-Tenant, getrennte Instanzen `restgeld_core` & `restgeld_auth`) |
| **Dev SMTP Mailer** | Mailpit (In-Memory Mailserver für lokale Auth-Tests) |
| **Cloud & Hosting** | Vercel (Preview & Edge Deployments), Docker, Kubernetes (k3s) |

---

## 🌲 Umgebungs- & Deployment-Modell

- **`main` → Production**: Produktiv-Release.
- **`develop` / Feature-Branches → Preview Environments**: Automatisches **Vercel Preview Deployment** bei jedem Git Push mit getrennter Preview-Datenbank.

---

## ⚡ Lokale Entwicklung

### 1-Befehl-Start: Live-Entwicklung mit Compose

Startet die Datenbanken, Backend-Microservices, Auth-Service und den Vite Dev-Server mit Hot-Module-Replacement (HMR):

```bash
docker compose -f docker-compose.dev.yml up --build
```

- **Frontend (Live-HMR):** [http://localhost:5173](http://localhost:5173)
- **Core API Backend:** [http://localhost:8080](http://localhost:8080)
- **Auth Microservice:** [http://localhost:8081](http://localhost:8081)
- **Billing Microservice:** [http://localhost:8082](http://localhost:8082)
- **PostgreSQL DB:** `localhost:5432`

---

## 📄 Lizenz

MIT License. 100% Open Source.
