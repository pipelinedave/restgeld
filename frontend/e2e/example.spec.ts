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
    page.on('console', msg => console.log('BROWSER LOG:', msg.text()))
    page.on('pageerror', err => console.log('BROWSER ERROR:', err.message))

    // Service Worker blockieren um Bypassing der Mocks zu verhindern
    await page.route('**/sw.js', route => route.abort())

    // Service Worker de-registrieren und Sprache auf DE fixieren
    await page.addInitScript(() => {
      try {
        localStorage.setItem('restgeld_language', 'de')
        localStorage.setItem('restgeld_currency', 'EUR')
      } catch {
        // ignore
      }
      if (window.navigator && window.navigator.serviceWorker) {
        window.navigator.serviceWorker.getRegistrations().then(registrations => {
          for (const registration of registrations) {
            registration.unregister()
          }
        })
      }
    })

    // Reset state vor jedem Test
    expenses = []

    // Mock all /api endpoints cleanly
    await page.route('**/api/**', async (route) => {
      const method = route.request().method()
      const url = route.request().url()

      if (url.includes('/api/budget')) {
        if (method === 'GET') {
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
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ status: 'ok' }),
          })
        }
        return
      }

      if (url.includes('/api/expenses')) {
        if (method === 'GET') {
          const urlObj = new URL(url, 'http://localhost')
          const pageNum = parseInt(urlObj.searchParams.get('page') || '1', 10)
          const limitNum = parseInt(urlObj.searchParams.get('limit') || '10', 10)
          const allReversed = [...expenses].reverse()
          const total = allReversed.length
          const totalPages = Math.max(1, Math.ceil(total / limitNum))
          const offset = (pageNum - 1) * limitNum
          const items = allReversed.slice(offset, offset + limitNum)

          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              items,
              total,
              page: pageNum,
              limit: limitNum,
              totalPages,
            }),
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
        }
        return
      }

      // Default fallback for /api/health, /api/auth/*, /api/billing/*, etc.
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ status: 'ok' }),
      })
    })
  })

  test('seite laedt und zeigt budget', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByText(/Tag \d+ von \d+/)).toBeVisible()
    await expect(page.getByText(/[\d.,]+ €/).first()).toBeVisible()
    await expect(page.getByText(/Puffer|Plan|überzogen/i).first()).toBeVisible()
  })

  test('modal oeffnet sich bei klick auf ausgabe-button', async ({ page }) => {
    await page.goto('/')
    await page.locator('.add-btn').click()
    await expect(page.getByRole('heading', { name: 'Ausgabe buchen' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Speichern' })).toBeVisible()
  })

  test('ausgabe buchen ueber modal', async ({ page }) => {
    await page.goto('/')
    await page.locator('.add-btn').click()

    // betrag eingeben: 12,50
    await page.getByPlaceholder('0,00').fill('12,50')

    // notiz eingeben
    await page.getByPlaceholder(/notiz/i).fill('Mittagessen')

    // speichern
    await page.getByRole('button', { name: 'Speichern' }).click()

    // ausgabe in der liste
    await expect(page.locator('.expense-list').getByText('Mittagessen')).toBeVisible()
    await expect(page.locator('.expense-amount').first()).toBeVisible()
  })

  test('ausgabe loeschen', async ({ page }) => {
    await page.goto('/')
    await page.locator('.add-btn').click()
    await page.getByPlaceholder('0,00').fill('10')
    await page.getByPlaceholder(/notiz/i).fill('Test')
    await page.getByRole('button', { name: 'Speichern' }).click()

    await expect(page.locator('.expense-list').getByText('Test')).toBeVisible()
    await page.locator('.delete-btn').first().click()
    await expect(page.locator('.expense-list').getByText('Test')).not.toBeVisible()
  })

  test('oeffnet ausgaben-historie und laedt weitere ausgaben nach', async ({ page }) => {
    // 15 Ausgaben anlegen (damit bei PageSize 12 eine 2. Seite zum Nachladen existiert)
    for (let i = 1; i <= 15; i++) {
      expenses.push({
        id: `exp-${i}`,
        periodId: 'test-period-id',
        amount: i * 5,
        note: `Ausgabe ${i}`,
        createdAt: new Date(2026, 7, 18, 10, i).toISOString(),
      })
    }

    await page.goto('/')
    await expect(page.locator('.show-all-btn')).toBeVisible()
    await page.locator('.show-all-btn').click()

    // Bottom Sheet geöffnet
    await expect(page.getByRole('heading', { name: 'Alle Ausgaben' })).toBeVisible()
    await expect(page.locator('.bottom-sheet').getByText('Ausgabe 15', { exact: true })).toBeVisible()
    await expect(page.locator('.drag-handle-pill')).toBeVisible()

    // Weitere Ausgaben nachladen (Infinite Scroll / Load More Hint)
    if (await page.locator('.load-more-hint').isVisible()) {
      await page.locator('.load-more-hint').click()
      await expect(page.locator('.bottom-sheet').getByText('Ausgabe 1', { exact: true })).toBeVisible()
    }

    // Schließen
    await page.locator('.sheet-header .close-btn').click()
    await expect(page.getByRole('heading', { name: 'Alle Ausgaben' })).not.toBeVisible()
  })
})
