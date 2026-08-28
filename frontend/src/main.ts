import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

if ('serviceWorker' in navigator && !navigator.webdriver) {
  navigator.serviceWorker.register('/sw.js')
}

createApp(App).mount('#app')
