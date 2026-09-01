import { test, expect } from '@playwright/test'

test('Live E2E: oeffnet echte Ausgaben-Historie und prueft Console & Modal', async ({ page }) => {
  test.skip(!process.env.LIVE_BACKEND, 'Skip live integration test unless LIVE_BACKEND is set')
  const consoleErrors: string[] = []
  page.on('console', (msg) => {
    console.log('LIVE BROWSER LOG:', msg.text())
    if (msg.type() === 'error') {
      consoleErrors.push(msg.text())
    }
  })
  page.on('pageerror', (err) => {
    console.log('LIVE BROWSER ERROR:', err.message)
    consoleErrors.push(err.message)
  })

  await page.goto('http://localhost:5173')
  await expect(page.getByText('restgeld.')).toBeVisible()

  // Ausgabe hinzufügen
  await page.getByRole('button', { name: /ausgabe/i }).click()
  await page.getByPlaceholder('0,00').fill('5')
  await page.getByPlaceholder(/notiz/i).fill('Live-Test Kaffee')
  await page.getByRole('button', { name: 'Speichern' }).click()

  // "Alle anzeigen" Button anklicken
  const showAllBtn = page.locator('.show-all-btn')
  await expect(showAllBtn).toBeVisible()
  await showAllBtn.click()

  // Bottom Sheet muss sichtbar sein und keine Fehler werfen
  await expect(page.getByRole('heading', { name: 'Alle Ausgaben' })).toBeVisible()
  await expect(page.locator('.bottom-sheet').getByText('Live-Test Kaffee').first()).toBeVisible()

  // Schließen
  await page.locator('.sheet-header .close-btn').click()
  await expect(page.getByRole('heading', { name: 'Alle Ausgaben' })).not.toBeVisible()

  // Console darf 0 Errors enthalten
  expect(consoleErrors).toEqual([])
})
