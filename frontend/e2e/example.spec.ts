import { test, expect } from '@playwright/test'

test.describe('Restgeld E2E', () => {
  test('seite laedt und zeigt budget', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByText(/Tag \d+ von \d+/)).toBeVisible()
    await expect(page.getByText(/[\d.,]+ €/)).toBeVisible()
    await expect(page.getByText(/angespart/)).toBeVisible()
  })

  test('numpad oeffnet sich bei klick auf ausgabe-button', async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /ausgabe/i }).click()
    await expect(page.getByText('Bestätigen')).toBeVisible()
  })

  test('ausgabe buchen ueber numpad', async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /ausgabe/i }).click()

    // betrag eingeben: 12,50
    await page.getByRole('button', { name: '1' }).click()
    await page.getByRole('button', { name: '2' }).click()
    await page.getByRole('button', { name: ',' }).click()
    await page.getByRole('button', { name: '5' }).click()
    await page.getByRole('button', { name: 'Bestätigen' }).click()

    // notiz eingeben
    await page.getByPlaceholder('Notiz').fill('Mittagessen')

    // speichern
    await page.getByRole('button', { name: 'Speichern' }).click()

    // ausgabe in der liste
    await expect(page.getByText('Mittagessen')).toBeVisible()
    await expect(page.getByText('-12,50')).toBeVisible()
  })

  test('ausgabe loeschen', async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /ausgabe/i }).click()
    await page.getByRole('button', { name: '1' }).click()
    await page.getByRole('button', { name: '0' }).click()
    await page.getByRole('button', { name: 'Bestätigen' }).click()
    await page.getByPlaceholder('Notiz').fill('Test')
    await page.getByRole('button', { name: 'Speichern' }).click()

    await expect(page.getByText('-10,00')).toBeVisible()
    await page.getByRole('button', { name: 'Löschen' }).click()
    await expect(page.getByText('-10,00')).not.toBeVisible()
  })
})
