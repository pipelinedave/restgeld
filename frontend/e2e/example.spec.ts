import { test, expect } from '@playwright/test'

interface Expense {
  id: string
  periodId: string
  amount: number
  note: string
  createdAt: string
}

test.describe('Restgeld E2E', () => {
  let expenses: Expense[] = []

  test.beforeEach(async ({ page }) => {
    // Reset state vor jedem Test
    expenses = []

    // Mock GET /api/budget
    await page.route('**/api/budget', async (route) => {
      if (route.request().method() === 'GET') {
        const totalSpent = expenses.reduce((acc, exp) => acc + exp.amount, 0)
        const day = 15
        const monthDays = 30
        const baseBudget = 20.0
        const savings = baseBudget * day - totalSpent
        const remainingDays = monthDays - day + 1
        const currentBudget = baseBudget + savings / remainingDays

        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            day,
            monthDays,
            baseBudget,
            currentBudget,
            savings,
            color: savings > 0 ? 'green' : savings === 0 ? 'white' : 'red',
            periodId: 'test-period-id',
            expenses: [...expenses].reverse().slice(0, 3),
          }),
        })
      } else {
        await route.continue()
      }
    })

    // Mock /api/expenses endpoints
    await page.route('**/api/expenses**', async (route) => {
      const method = route.request().method()
      const url = route.request().url()

      if (method === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([...expenses].reverse().slice(0, 3)),
        })
      } else if (method === 'POST') {
        const postData = route.request().postDataJSON()
        const newExpense: Expense = {
          id: 'exp-' + Date.now() + '-' + Math.random().toString(36).substring(2, 7),
          periodId: 'test-period-id',
          amount: Number(postData.amount),
          note: postData.note || '',
          createdAt: new Date().toISOString(),
        }
        expenses.push(newExpense)
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify(newExpense),
        })
      } else if (method === 'DELETE') {
        const parts = url.split('/api/expenses/')
        const id = parts[1]?.split('?')[0]
        expenses = expenses.filter((e) => e.id !== id)
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ status: 'deleted' }),
        })
      } else {
        await route.continue()
      }
    })
  })

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
    await page.getByPlaceholder(/notiz/i).fill('Mittagessen')

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
    await page.getByPlaceholder(/notiz/i).fill('Test')
    await page.getByRole('button', { name: 'Speichern' }).click()

    await expect(page.getByText('-10,00')).toBeVisible()
    await page.getByRole('button', { name: 'Löschen' }).click()
    await expect(page.getByText('-10,00')).not.toBeVisible()
  })
})
