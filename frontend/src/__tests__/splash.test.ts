import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

// Modulzustand (internes `hidden`-Flag) pro Test zuruecksetzen, damit die
// Idempotenz sauber geprueft werden kann.
async function loadSplash() {
  vi.resetModules()
  return (await import('../splash')).default
}

describe('splash.ts Splash-Screen', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <div id="app-splash">
        <div class="splash-brand">restgeld.</div>
        <div class="splash-bar"><div class="splash-bar-inner"></div></div>
      </div>
      <div id="app"></div>
    `
    vi.useFakeTimers({ toFake: ['setTimeout', 'Date'] })
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('setzt die hide-Klasse nach der Mindestanzeigezeit und entfernt das Element', async () => {
    const hideSplash = await loadSplash()
    hideSplash()

    // Vor Ablauf bleibt der Splash sichtbar.
    vi.advanceTimersByTime(200)
    let splash = document.getElementById('app-splash')
    expect(splash).not.toBeNull()
    expect(splash!.classList.contains('hide')).toBe(false)

    // Nach Ablauf der Mindestanzeigezeit wird hide gesetzt.
    vi.advanceTimersByTime(500)
    splash = document.getElementById('app-splash')
    expect(splash!.classList.contains('hide')).toBe(true)

    // transitionend entfernt das Element schliesslich aus dem DOM.
    splash!.dispatchEvent(new Event('transitionend'))
    expect(document.getElementById('app-splash')).toBeNull()
  })

  it('entfernt das Element per Fallback, falls transitionend nicht feuert', async () => {
    const hideSplash = await loadSplash()
    hideSplash()
    vi.advanceTimersByTime(700)
    expect(document.getElementById('app-splash')!.classList.contains('hide')).toBe(true)
    // Fallback-Timeout (500ms) entfernt es trotzdem.
    vi.advanceTimersByTime(600)
    expect(document.getElementById('app-splash')).toBeNull()
  })

  it('ist idempotent - zweiter Aufruf ist ein No-op', async () => {
    const hideSplash = await loadSplash()
    hideSplash()
    const first = document.getElementById('app-splash')
    hideSplash()
    expect(document.getElementById('app-splash')).toBe(first)
  })

  it('funktioniert fehlerfrei, wenn kein Splash im DOM existiert', async () => {
    const hideSplash = await loadSplash()
    document.body.innerHTML = '<div id="app"></div>'
    expect(() => hideSplash()).not.toThrow()
  })
})
