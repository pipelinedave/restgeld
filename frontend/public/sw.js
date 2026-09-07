const CACHE_NAME = 'restgeld-v2'
const PRECACHE_URLS = ['/manifest.webmanifest', '/favicon.svg']

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(PRECACHE_URLS))
  )
  self.skipWaiting()
})

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting()
  }
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys.map((key) => {
          if (key !== CACHE_NAME) {
            return caches.delete(key)
          }
        })
      )
    ).then(() => self.clients.claim())
  )
})

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return
  const url = new URL(event.request.url)

  // API calls are handled separately / network-first by app
  if (url.pathname.startsWith('/api/')) return

  const isNavigation = event.request.mode === 'navigate' || event.request.destination === 'document'

  if (isNavigation) {
    // Network-First for HTML navigation to always serve the latest build and Vite chunk hashes
    event.respondWith(
      fetch(event.request)
        .then((networkResponse) => {
          if (networkResponse && networkResponse.status === 200) {
            const clone = networkResponse.clone()
            caches.open(CACHE_NAME).then((cache) => cache.put(event.request, clone))
          }
          return networkResponse
        })
        .catch(async () => {
          const cached = await caches.match(event.request)
          if (cached) return cached
          const rootCached = await caches.match('/')
          if (rootCached) return rootCached
          return new Response('<h1>Offline</h1><p>restgeld ist momentan offline.</p>', {
            headers: { 'Content-Type': 'text/html; charset=utf-8' },
          })
        })
    )
    return
  }

  // Cache-First for versioned / static assets
  event.respondWith(
    caches.match(event.request).then((cachedResponse) => {
      if (cachedResponse) {
        return cachedResponse
      }

      return fetch(event.request)
        .then((networkResponse) => {
          if (networkResponse && networkResponse.status === 200) {
            const clone = networkResponse.clone()
            caches.open(CACHE_NAME).then((cache) => cache.put(event.request, clone))
          }
          return networkResponse
        })
        .catch(() => {
          // Never return null from respondWith
          return new Response('', { status: 408, statusText: 'Request timed out' })
        })
    })
  )
})
