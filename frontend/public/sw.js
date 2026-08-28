// Restgeld Service Worker (PWA)
//
// Ziele:
//  - Nach jedem Deployment sofort die neue Version anzeigen (kein "Browser-Daten
//    loeschen" mehr noetig): Navigations-/HTML-Requests sind NETWORK-FIRST.
//  - Schnelles Laden aller statischen, gehashten Assets (JS/CSS Bundles, Icons):
//    CACHE-FIRST mit Netzwerk-Fallback.
//  - Automatisches Aufraeumen alter Caches beim Aktivieren.
//
// Assets von Vite werden in `dist/assets/` mit Inhalts-Hash im Dateinamen gebaut.
// Dadurch koennen wir sie unbegrenzt cachen, ohne dass eine neue Version veraltet:
// jeder neue Build erzeugt neue Dateinamen und das alte HTML laedt die neuen.

const OFFLINE_ASSET_CACHE = 'restgeld-assets'
const PRECACHE_FALLBACK_URLS = ['/', '/favicon.svg', '/manifest.webmanifest']

// Statische Vite-Assets (gehasht) + Icons sicher cachen. Nicht-cache-bar:
// HTML/Navigation und API-Requests.

function isNavigation(req) {
  return req.mode === 'navigate'
}

function isSameOrigin(req) {
  return new URL(req.url).origin === self.location.origin
}

// API-/Daten-Requests niemals aus dem Cache bedienen.
function isApiRequest(req) {
  return new URL(req.url).pathname.startsWith('/api')
}

self.addEventListener('install', (event) => {
  // Offline-Fallback-Grundgeruest vorhalten (kein addAll auf '/' - das geht
  // network-first und darf hier nicht den Cache "frieren").
  event.waitUntil(
    caches.open(OFFLINE_ASSET_CACHE).then((cache) =>
      cache.addAll(PRECACHE_FALLBACK_URLS.filter((u) => u !== '/'))
    )
  )
  // Sofort den neuen SW aktiv werden lassen, ohne auf das Schliessen aller
  // Tabs zu warten.
  self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      // Alte Caches dieser App entfernen - verhindert veraltete Assets.
      const keys = await caches.keys()
      await Promise.all(
        keys
          .filter((k) => k !== OFFLINE_ASSET_CACHE)
          .map((k) => caches.delete(k))
      )
      await self.clients.claim()
    })()
  )
})

self.addEventListener('fetch', (event) => {
  const req = event.request
  if (req.method !== 'GET') return

  const url = new URL(req.url)

  // Nur gleiche-Herkunft bearbeiten; Cross-Origin (Fonts, API-Extern) direkt.
  if (!isSameOrigin(req)) return

  // API-Requests immer live (nie gecacht).
  if (isApiRequest(req)) return

  // Navigation (HTML): NETWORK-FIRST - garantiert frische Version nach Deploy.
  if (isNavigation(req) || (req.mode === 'same-origin' && url.pathname === '/')) {
    event.respondWith(
      fetch(req)
        .then((res) => {
          if (res.ok) {
            const copy = res.clone()
            caches.open(OFFLINE_ASSET_CACHE).then((cache) => cache.put('/', copy))
          }
          return res
        })
        .catch(() =>
          caches.match('/').then(
            (cached) => cached || caches.match('/index.html')
          )
        )
    )
    return
  }

  // Statische Assets (gehasht): CACHE-FIRST mit Netzwerk-Fallback.
  event.respondWith(
    caches.match(req).then((cached) => {
      if (cached) return cached
      return fetch(req).then((res) => {
        if (res && res.ok) {
          const copy = res.clone()
          caches.open(OFFLINE_ASSET_CACHE).then((cache) => cache.put(req, copy))
        }
        return res
      })
    })
  )
})
