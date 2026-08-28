import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

// Service-Worker-Registrierung mit Auto-Update: Nach einem Deployment wird die
// neue App-Version automatisch in den Vordergrund geholt (kein manuelles
// "Browser-Daten loeschen" mehr noetig). updateViaCache 'none' stellt sicher,
// dass immer die frische sw.js vom Server geholt wird.
if ('serviceWorker' in navigator && !navigator.webdriver) {
  navigator.serviceWorker.register('/sw.js', { updateViaCache: 'none' }).then((reg) => {
    // Wenn ein neuer (aktiver) SW die Kontrolle uebernimmt, lade die Seite neu,
    // damit neue Assets/Bundles sofort greifen.
    reg.addEventListener('updatefound', () => {
      const newWorker = reg.installing
      if (!newWorker) return
      newWorker.addEventListener('statechange', () => {
        if (newWorker.state === 'activated' && navigator.serviceWorker.controller) {
          window.location.reload()
        }
      })
    })
  })
}

createApp(App).mount('#app')
