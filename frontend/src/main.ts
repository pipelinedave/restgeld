import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

// Automatic reload recovery on stale chunk 404s after new deployments
window.addEventListener('vite:preloadError', (event) => {
  console.warn('Stale chunk detected after deployment, reloading...', event)
  window.location.reload()
})

if ('serviceWorker' in navigator && !(window as any).__DISABLE_SW__) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .then((reg) => {
        // Check for updates on load
        reg.update().catch(() => {})
      })
      .catch(() => {})
  })
}

createApp(App).mount('#app')
