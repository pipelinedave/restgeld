import { ref } from 'vue'

export interface QueuedExpense {
  id: string
  amount: number
  note: string
  createdAt: string
}

const OUTBOX_KEY = 'restgeld_offline_outbox'

export function useOfflineSync() {
  const isOnline = ref<boolean>(typeof navigator !== 'undefined' ? navigator.onLine : true)
  const pendingCount = ref<number>(getOutboxCount())

  function getOutboxCount(): number {
    try {
      const raw = localStorage.getItem(OUTBOX_KEY)
      return raw ? JSON.parse(raw).length : 0
    } catch {
      return 0
    }
  }

  function getOutbox(): QueuedExpense[] {
    try {
      const raw = localStorage.getItem(OUTBOX_KEY)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  function saveOutbox(items: QueuedExpense[]) {
    try {
      localStorage.setItem(OUTBOX_KEY, JSON.stringify(items))
      pendingCount.value = items.length
    } catch {
      // Ignore localStorage errors
    }
  }

  function enqueueExpense(amount: number, note: string): QueuedExpense {
    const item: QueuedExpense = {
      id: `offline-${Date.now()}-${Math.random().toString(36).substring(2, 7)}`,
      amount,
      note,
      createdAt: new Date().toISOString(),
    }
    const current = getOutbox()
    current.push(item)
    saveOutbox(current)
    return item
  }

  async function syncPendingExpenses(api: { addExpense: (amount: number, note: string) => Promise<any> }): Promise<number> {
    const items = getOutbox()
    if (items.length === 0) return 0

    let synced = 0
    const remaining: QueuedExpense[] = []

    for (const item of items) {
      try {
        await api.addExpense(item.amount, item.note)
        synced++
      } catch (err) {
        remaining.push(item)
      }
    }

    saveOutbox(remaining)
    return synced
  }

  function initListeners(onOnlineCallback?: () => void) {
    if (typeof window === 'undefined') return

    const handleOnline = () => {
      isOnline.value = true
      onOnlineCallback?.()
    }

    const handleOffline = () => {
      isOnline.value = false
    }

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)

    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', handleOffline)
    }
  }

  return {
    isOnline,
    pendingCount,
    getOutbox,
    enqueueExpense,
    syncPendingExpenses,
    initListeners,
  }
}
