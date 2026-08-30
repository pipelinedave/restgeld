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

export interface Period {
  id?: string
  startDate: string
  monthDays: number
  baseBudget: number
  monthlyTotal: number
}

export interface PeriodSummary {
  id: string
  startDate: string
  monthDays: number
  baseBudget: number
  monthlyTotal: number
  totalSpent: number
  savings: number
  expenseCount: number
}

const BASE = import.meta.env.PROD
  ? window.location.origin
  : ''

async function api<T>(path: string, options?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('restgeld_auth_token') : null
  const lang = typeof window !== 'undefined' ? localStorage.getItem('restgeld_language') || 'de' : 'de'
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Accept-Language': lang,
    ...(options?.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || 'api-fehler')
  }
  return res.json()
}

export interface ServiceTelemetry {
  id: string
  name: string
  status: 'up' | 'degraded' | 'down'
  latencyMs: number
  url: string
  checkedAt: string
  details?: Record<string, any>
  error?: string
}

export function useApi() {
  return {
    getBudget: () => api<BudgetData>('/api/budget'),
    getPeriods: () => api<PeriodSummary[]>('/api/periods'),
    getMonitoringOverview: () => api<{
      status: 'healthy' | 'degraded' | 'critical'
      timestamp: string
      uptimeSeconds: number
      services: ServiceTelemetry[]
      system: {
        goVersion: string
        goroutines: number
        memoryAllocMb: number
        memorySysMb: number
        gcCount: number
        uptimeSeconds: number
      }
    }>('/api/monitoring/overview'),
    getExpenses: (page = 1, limit = 10, periodId?: string) => {
      const query = new URLSearchParams({ page: String(page), limit: String(limit) })
      if (periodId) query.set('period_id', periodId)
      return api<PaginatedExpenses>(`/api/expenses?${query.toString()}`)
    },
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
