import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useOfflineSync } from '../composables/useOfflineSync'

describe('useOfflineSync', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  it('initialisiert mit leerer Outbox', () => {
    const { isOnline, pendingCount, getOutbox } = useOfflineSync()
    expect(isOnline.value).toBe(true)
    expect(pendingCount.value).toBe(0)
    expect(getOutbox()).toEqual([])
  })

  it('fuegt Ausgaben zur Outbox hinzu', () => {
    const { enqueueExpense, pendingCount, getOutbox } = useOfflineSync()
    const item = enqueueExpense(14.5, 'Test Offline')

    expect(item.amount).toBe(14.5)
    expect(item.note).toBe('Test Offline')
    expect(pendingCount.value).toBe(1)
    expect(getOutbox().length).toBe(1)
    expect(getOutbox()[0].id).toBe(item.id)
  })

  it('synchronisiert ausstehende Ausgaben erfolgreich', async () => {
    const { enqueueExpense, syncPendingExpenses, pendingCount } = useOfflineSync()
    enqueueExpense(5.0, 'Kaffee')
    enqueueExpense(12.0, 'Lunch')

    expect(pendingCount.value).toBe(2)

    const mockApi = {
      addExpense: vi.fn().mockResolvedValue({ id: '1', amount: 5.0 }),
    }

    const synced = await syncPendingExpenses(mockApi)
    expect(synced).toBe(2)
    expect(mockApi.addExpense).toHaveBeenCalledTimes(2)
    expect(pendingCount.value).toBe(0)
  })

  it('behaelt fehlgeschlagene Ausgaben in der Outbox', async () => {
    const { enqueueExpense, syncPendingExpenses, pendingCount, getOutbox } = useOfflineSync()
    enqueueExpense(5.0, 'Kaffee')

    const mockApi = {
      addExpense: vi.fn().mockRejectedValue(new Error('Network error')),
    }

    const synced = await syncPendingExpenses(mockApi)
    expect(synced).toBe(0)
    expect(pendingCount.value).toBe(1)
    expect(getOutbox().length).toBe(1)
  })
})
