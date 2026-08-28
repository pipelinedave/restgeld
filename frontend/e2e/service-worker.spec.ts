import { test, expect } from '@playwright/test'

/**
 * E2E-Tests fuer den Restgeld Service Worker (PWA-Cache-Strategie).
 *
 * Verifiziert das Kern-Problem aus dem Produktionsfeedback: Nach einem
 * Deployment durfte die App nicht mehr laden, bis man manuell die
 * Browser-Daten geloescht hat. Die neue Strategie ist NETWORK-FIRST fuer
 * HTML/Navigation + CACHE-FIRST fuer gehashte Assets + Cache-Purge beim
 * Aktivieren.
 *
 * Diese Tests laufen gegen den PRODUKTIONS-Build (vite preview), weil nur dort
 * gehashte Assets und der echte Service Worker greifen.
 */

test.describe('Service Worker / PWA', () => {
  test('registriert den Service Worker und cached statische Assets', async ({ page }) => {
    const errors: string[] = []
    page.on('console', (msg) => {
      if (msg.type() === 'error') errors.push(msg.text())
    })
    page.on('pageerror', (err) => errors.push(err.message))

    await page.goto('/')
    await expect(page.getByText(/Tag \d+ von \d+/)).toBeVisible()

    // Auf die (asynchrone) SW-Registrierung warten.
    await page.waitForFunction(async () => {
      if (!('serviceWorker' in navigator)) return false
      const registration = await navigator.serviceWorker.getRegistration('/')
      return !!registration && !!registration.active
    })

    // Service Worker der App muss registriert worden sein.
    const reg = await page.evaluate(async () => {
      const registration = await navigator.serviceWorker.getRegistration('/')
      return registration
        ? { active: !!registration.active, scope: registration.scope }
        : null
    })
    expect(reg).not.toBeNull()
    expect(reg!.active).toBe(true)

    // Der Asset-Cache muss angelegt worden sein.
    const cacheInfo = await page.evaluate(async () => {
      const keys = await caches.keys()
      const entries: string[] = []
      for (const key of keys) {
        const c = await caches.open(key)
        const urls = (await c.keys()).map((r) => new URL(r.url).pathname)
        entries.push(...urls)
      }
      return { keys, entries }
    })
    expect(cacheInfo.keys.length).toBeGreaterThan(0)
    expect(cacheInfo.entries.some((p) => p.startsWith('/assets/'))).toBe(true)

    expect(errors.length).toBe(0)
  })

  test('laedt die App nach einem Netzwerk-Ausfall aus dem Cache (Offline)', async ({ page }) => {
    // Erst online laden, damit der Cache gefuellt wird und der SW kontrolliert.
    await page.goto('/')
    await expect(page.getByText(/Tag \d+ von \d+/)).toBeVisible()

    // Warten bis der SW die Seite kontrolliert (aktiv + installiert).
    await page.waitForFunction(async () => {
      const reg = await navigator.serviceWorker.getRegistration('/')
      return !!(reg && reg.active && navigator.serviceWorker.controller)
    })

    // Kontrolle an den SW abgeben und die Seite neu laden, damit echte
    // Fetch-Events durch den SW laufen.
    await page.reload()
    await expect(page.getByText(/Tag \d+ von \d+/)).toBeVisible()

    // Netz offline schalten (nur die Navigation unterbinden - API darf fuer
    // diesen Test weiterhin durch, da wir nur das HTML-Shell-Offline-Verhalten
    // pruefen). HTML-Dokumente muessen aus dem Cache bedient werden.
    await page.context().setOffline(true)
    await page.reload().catch(() => {})
    await page.waitForLoadState('domcontentloaded').catch(() => {})

    const offlineBody = await page.evaluate(() => document.body.innerText)
    expect(offlineBody).toContain('restgeld')

    await page.context().setOffline(false)
  })
})
