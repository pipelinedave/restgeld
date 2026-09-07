import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

// Automatic reload recovery on stale chunk 404s after new deployments
window.addEventListener('vite:preloadError', (event) => {
  console.warn('Stale chunk detected after deployment, reloading...', event)
  window.location.reload()
})

if ('serviceWorker' in navigator && !(window as any).__DISABLE_SW__) {
  let refreshing = false
  navigator.serviceWorker.addEventListener('controllerchange', () => {
    if (!refreshing) {
      refreshing = true
      window.location.reload()
    }
  })

  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .then((reg) => {
        reg.update().catch(() => {})

        if (reg.waiting) {
          reg.waiting.postMessage({ type: 'SKIP_WAITING' })
        }

        reg.addEventListener('updatefound', () => {
          const newWorker = reg.installing
          if (newWorker) {
            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                newWorker.postMessage({ type: 'SKIP_WAITING' })
              }
            })
          }
        })
      })
      .catch(() => {})
  })
}

createApp(App).mount('#app')
