import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

beforeEach(() => {
  mockFetch.mockReset()
})

import { useApi, type BudgetData, type Expense } from '../composables/useApi'

const api = useApi()

function mockResponse(data: unknown, status = 200) {
  return Promise.resolve({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(data)
  } as Response)
}

describe('useApi', () => {
  it('getBudget ruft /api/budget auf', async () => {
    const data: BudgetData = {
      day: 17, monthDays: 31, baseBudget: 14.52, currentBudget: 19.24,
      savings: 61.36, color: 'green', periodId: '2026-08', expenses: []
    }
    mockFetch.mockResolvedValueOnce(mockResponse(data))
    const result = await api.getBudget()
    expect(mockFetch).toHaveBeenCalledWith('/api/budget', expect.any(Object))
    expect(result).toEqual(data)
  })

  it('getExpenses ruft /api/expenses mit Parametern auf', async () => {
    const data = {
      items: [{ id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Test', createdAt: '2026-08-18T06:35:43Z' }],
      total: 1,
      page: 2,
      limit: 5,
      totalPages: 1
    }
    mockFetch.mockResolvedValueOnce(mockResponse(data))
    const result = await api.getExpenses(2, 5)
    expect(mockFetch).toHaveBeenCalledWith('/api/expenses?page=2&limit=5', expect.any(Object))
    expect(result).toEqual(data)
  })

  it('addExpense ruft POST /api/expenses auf', async () => {
    const expense: Expense = { id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Test', createdAt: '2026-08-18T06:35:43Z' }
    mockFetch.mockResolvedValueOnce(mockResponse(expense, 201))
    const result = await api.addExpense(8.5, 'Test')
    expect(mockFetch).toHaveBeenCalledWith('/api/expenses', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ amount: 8.5, note: 'Test' })
    }))
    expect(result).toEqual(expense)
  })

  it('deleteExpense ruft DELETE /api/expenses/{id} auf', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ status: 'ok' }))
    const result = await api.deleteExpense('e1')
    expect(mockFetch).toHaveBeenCalledWith('/api/expenses/e1', expect.objectContaining({
      method: 'DELETE'
    }))
    expect(result).toEqual({ status: 'ok' })
  })

  it('newPeriod ruft POST /api/period auf', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ id: '2026-08' }))
    const result = await api.newPeriod()
    expect(mockFetch).toHaveBeenCalledWith('/api/period', expect.objectContaining({
      method: 'POST'
    }))
    expect(result).toEqual({ id: '2026-08' })
  })

  it('updateBudget ruft PATCH /api/budget auf', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ status: 'ok' }))
    const result = await api.updateBudget(600)
    expect(mockFetch).toHaveBeenCalledWith('/api/budget', expect.objectContaining({
      method: 'PATCH',
      body: JSON.stringify({ monthlyTotal: 600 })
    }))
    expect(result).toEqual({ status: 'ok' })
  })

  it('exportData ruft GET /api/export auf', async () => {
    const dummyBlob = new Blob(['test'])
    mockFetch.mockResolvedValueOnce({
      ok: true,
      blob: () => Promise.resolve(dummyBlob),
    })
    const result = await api.exportData('json')
    expect(mockFetch).toHaveBeenCalledWith('/api/export?format=json')
    expect(result).toBe(dummyBlob)
  })

  it('importData ruft POST /api/import auf', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ status: 'ok', imported: 3 }))
    const result = await api.importData('{"expenses":[]}', false)
    expect(mockFetch).toHaveBeenCalledWith('/api/import', expect.objectContaining({
      method: 'POST',
      body: '{"expenses":[]}',
      headers: { 'Content-Type': 'application/json' },
    }))
    expect(result).toEqual({ status: 'ok', imported: 3 })
  })

  it('wirft fehler bei nicht-ok response', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ error: 'kaputt' }, 400))
    await expect(api.getBudget()).rejects.toThrow('kaputt')
  })

  it('wirft fehler bei leerem error-body', async () => {
    mockFetch.mockResolvedValueOnce(Promise.resolve({
      ok: false,
      status: 500,
      json: () => Promise.reject(new Error('kein json'))
    } as Response))
    await expect(api.getBudget()).rejects.toThrow('api-fehler')
  })
})
