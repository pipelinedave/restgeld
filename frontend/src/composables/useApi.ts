export interface Expense {
  id: string
  periodId: string
  amount: number
  note: string
  createdAt: string
}

export interface PaginatedExpenses {
  items: Expense[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export interface DailyStat {
  day: number
  date: string
  spent: number
}

export interface StreakInfo {
  currentStreak: number
  longestStreak: number
  noSpendDays: number
  underBudgetDays: number
}

export interface ProjectionInfo {
  projectedSavings: number
  projectedTotalSpent: number
  avgDailySpend: number
  status: 'saving' | 'deficit'
}

export interface BudgetData {
  day: number
  monthDays: number
  baseBudget: number
  currentBudget: number
  savings: number
  color: string
  periodId: string
  expenses: Expense[]
  dailyStats?: DailyStat[]
  streak?: StreakInfo
  projection?: ProjectionInfo
}

const BASE = import.meta.env.PROD
  ? window.location.origin
  : ''

async function api<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || 'api-fehler')
  }
  return res.json()
}

export function useApi() {
  return {
    getBudget: () => api<BudgetData>('/api/budget'),
    getExpenses: (page = 1, limit = 10) =>
      api<PaginatedExpenses>(`/api/expenses?page=${page}&limit=${limit}`),
    addExpense: (amount: number, note: string) =>
      api<Expense>('/api/expenses', {
        method: 'POST',
        body: JSON.stringify({ amount, note }),
      }),
    deleteExpense: (id: string) =>
      api<{ status: string }>(`/api/expenses/${id}`, { method: 'DELETE' }),
    newPeriod: (monthlyTotal?: number, startDate?: string, days?: number) =>
      api<{ id: string }>('/api/period', {
        method: 'POST',
        body: JSON.stringify({ monthlyTotal, startDate, days }),
      }),
    updateBudget: (monthlyTotal?: number, days?: number) =>
      api<{ status: string }>('/api/budget', {
        method: 'PATCH',
        body: JSON.stringify({ monthlyTotal, days }),
      }),
    exportData: async (format: 'json' | 'csv') => {
      const res = await fetch(`${BASE}/api/export?format=${format}`)
      if (!res.ok) throw new Error('Export fehlgeschlagen')
      return res.blob()
    },
    importData: async (content: string, isCsv = false) => {
      const res = await fetch(`${BASE}/api/import`, {
        method: 'POST',
        headers: {
          'Content-Type': isCsv ? 'text/csv' : 'application/json',
        },
        body: content,
      })
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }))
        throw new Error(err.error || 'Import fehlgeschlagen')
      }
      return res.json() as Promise<{ status: string; imported: number }>
    },
  }
}
