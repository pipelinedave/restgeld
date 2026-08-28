// Splash-Screen-Helfer: blendet den statischen Lade-Bildschirm aus, sobald die
// App startklar ist (Budget geladen). Der Splash selbst lebt in index.html,
// damit er sofort beim Seiteneinstieg erscheint (vor dem Vue-Mount) und kein
// "unprofessioneller" Flicker von Halbzustand-Elementen sichtbar wird.

const SPLASH_SELECTOR = '#app-splash'
const MIN_DISPLAY_MS = 600

const startTs = Date.now()
let hidden = false

function hideSplash(): void {
  if (hidden || typeof document === 'undefined') return
  hidden = true
  const splash = document.querySelector<HTMLElement>(SPLASH_SELECTOR)
  if (!splash) return
  // Sicherstellen, dass der Splash kurz genug sichtbar bleibt, damit auch bei
  // blitzschnellem Laden kein haesslicher Bildwechsel entsteht.
  const elapsed = Date.now() - startTs
  const wait = Math.max(0, MIN_DISPLAY_MS - elapsed)
  window.setTimeout(() => {
    splash.classList.add('hide')
    splash.addEventListener('transitionend', () => splash.remove(), { once: true })
    // Fallback, falls transitionend nicht feuert.
    window.setTimeout(() => splash.remove(), 500)
  }, wait)
}

export default hideSplash
