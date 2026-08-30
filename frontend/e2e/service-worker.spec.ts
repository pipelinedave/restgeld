import { test, expect } from '@playwright/test'

/**
 * E2E-Tests fuer den Restgeld Service Worker (PWA-Cache-Strategie).
 */

test.describe('Service Worker / PWA', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      try {
        localStorage.setItem('restgeld_language', 'de')
        localStorage.setItem('restgeld_currency', 'EUR')
      } catch {
        // ignore
      }
    })
  })

  test('registriert den Service Worker und cached statische Assets', async ({ page }) => {
    const errors: string[] = []
    page.on('console', (msg) => {
      const txt = msg.text()
      if (msg.type() === 'error' && !txt.includes('500') && !txt.includes('/api/') && !txt.includes('Failed to load resource')) {
        errors.push(txt)
      }
    })
    page.on('pageerror', (err) => errors.push(err.message))

    await page.goto('/')
    await expect(page.locator('.brand-title')).toBeVisible()

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
    expect(cacheInfo.entries.length).toBeGreaterThan(0)

    expect(errors.length).toBe(0)
  })

  test('laedt die App nach einem Netzwerk-Ausfall aus dem Cache (Offline)', async ({ page }) => {
    // Erst online laden, damit der Cache gefuellt wird und der SW kontrolliert.
    await page.goto('/')
    await expect(page.locator('.brand-title')).toBeVisible()

    // Warten bis der SW die Seite kontrolliert (aktiv + installiert).
    await page.waitForFunction(async () => {
      const reg = await navigator.serviceWorker.getRegistration('/')
      return !!(reg && reg.active && navigator.serviceWorker.controller)
    })

    // Kontrolle an den SW abgeben und die Seite neu laden
    await page.reload()
    await expect(page.locator('.brand-title')).toBeVisible()

    // Netz offline schalten
    await page.context().setOffline(true)
    await page.reload().catch(() => {})
    await page.waitForLoadState('domcontentloaded').catch(() => {})

    const offlineBody = await page.evaluate(() => document.body.innerText)
    expect(offlineBody).toContain('restgeld')

    await page.context().setOffline(false)
  })
})
