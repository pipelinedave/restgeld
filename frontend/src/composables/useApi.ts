export interface Expense {
  id: string
  periodId: string
  amount: number
  note: string
  createdAt: string
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
    getExpenses: () => api<Expense[]>('/api/expenses'),
    addExpense: (amount: number, note: string) =>
      api<Expense>('/api/expenses', {
        method: 'POST',
        body: JSON.stringify({ amount, note }),
      }),
    deleteExpense: (id: string) =>
      api<{ status: string }>(`/api/expenses/${id}`, { method: 'DELETE' }),
    newPeriod: (monthlyTotal?: number, startDate?: string) =>
      api<{ id: string }>('/api/period', {
        method: 'POST',
        body: JSON.stringify({ monthlyTotal, startDate }),
      }),
    updateBudget: (monthlyTotal: number) =>
      api<{ status: string }>('/api/budget', {
        method: 'PATCH',
        body: JSON.stringify({ monthlyTotal }),
      }),
  }
}
