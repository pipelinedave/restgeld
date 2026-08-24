# Projekt: Daily Allowance Tracker (PWA)


Der Name der App ist "Restgeld".

## Ziel
Entwicklung einer minimalistischen Web-App (Mobile-first PWA) zum Tracken eines täglichen "Spending Money" Budgets. Das Projekt wird auf einem lokalen single-node k3s-Cluster betrieben.

## Architektur (3-Tier)
1. **Frontend:** PWA (z. B. Vue.js oder React mit Vite), bereitgestellt über Nginx.
2. **Backend:** REST API (z. B. Node.js/Express oder Go).
3. **Datenbank:** PostgreSQL-Container mit Persistent Volume Claim (PVC) für sichere Datenspeicherung.

## Kern-Logik
- **Basis-Budget:** Ein fixes Monatsbudget wird durch die Anzahl der Tage in der Periode geteilt.
- **Rollover:** Nicht ausgegebenes Tagesbudget wird auf die verbleibenden Tage der Periode addiert (Budget wächst).
- **Gamification-Metrik:** Anzeige der aktuellen Ersparnis im Vergleich zum normalen Basis-Budget (z. B. "+ 14,50 € angespart").

## UI/UX-Anforderungen (Fokus auf maximale Reduktion)
- **Header:** Dezenter Fortschrittsbalken für den Monat (z. B. "Tag 17 von 30").
- **Hero-Bereich:** Sehr große, zentrierte Zahl für das aktuelle Tagesbudget. Farbcodiert (Grün = mehr als Basis-Budget, Weiß = im Budget, Rot = überzogen).
- **Eingabe:** Ein großer Minus-Button. Öffnet ein *in der App gerendertes Numpad* (keine native Handy-Tastatur aufrufen!), um Beträge in Sekundenbruchteilen einzugeben.
- **Historie & Korrektur:** Liste der 3 aktuellsten Ausgaben direkt auf dem Startscreen mit "Löschen"-Funktion (Mülleimer-Icon) für fehlerhafte Einträge.

## DevOps & Output-Anforderungen
- Schreibe den App-Code für Frontend und Backend.
- Erstelle `Dockerfile`s für beide Services.
- Generiere die passenden Kubernetes-Manifeste (`Deployment`, `Service`, `Ingress`, `PVC` für die DB).
- Die Manifeste müssen so strukturiert sein, dass sie direkt in ein GitOps-Repository (`pipelinedave/k3s-config`) übernommen werden können.
