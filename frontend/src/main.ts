import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

if ('serviceWorker' in navigator && !(window as any).__DISABLE_SW__) {
  navigator.serviceWorker.register('/sw.js').catch(() => {})
}

createApp(App).mount('#app')
