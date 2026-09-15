import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import MonthlyTrend from '../components/MonthlyTrend.vue'
import type { TrendResponse } from '../composables/useApi'

const mockGetTrend = vi.fn()
vi.mock('../composables/useApi', () => ({
  useApi: () => ({
    getTrend: mockGetTrend,
  }),
}))

function makeTrend(): TrendResponse {
  return {
    months: [
      {
        month: '2026-06',
        startDate: '2026-06-01T00:00:00Z',
        monthlyTotal: 450,
        totalSpent: 420,
        savings: 30,
        expenseCount: 11,
        avgDailySpend: 14,
      },
      {
        month: '2026-07',
        startDate: '2026-07-01T00:00:00Z',
        monthlyTotal: 465,
        totalSpent: 480,
        savings: -15,
        expenseCount: 16,
        avgDailySpend: 15.48,
      },
    ],
    summary: {
      monthCount: 2,
      totalBudget: 915,
      totalSpent: 900,
      totalSaved: 15,
      avgSavings: 7.5,
      bestSavingsMonth: '2026-06',
    },
  }
}

describe('MonthlyTrend', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockGetTrend.mockResolvedValue(makeTrend())
  })

  it('laedt keinen Trend solange visible=false', () => {
    mount(MonthlyTrend, { props: { visible: false } })
    expect(mockGetTrend).not.toHaveBeenCalled()
  })

  it('laedt und rendert Monats-Verlauf mit Summary', async () => {
    const wrapper = mount(MonthlyTrend, { props: { visible: true } })
    expect(mockGetTrend).toHaveBeenCalledTimes(1)

    await flushPromises()

    expect(wrapper.find('.trend-month').exists()).toBe(true)
    expect(wrapper.find('.trend-month-title').text()).toContain('Monats-Verlauf')

    const columns = wrapper.findAll('.bar-column')
    expect(columns.length).toBe(2)

    // Spar-Übersicht: 2 Perioden, +15,00 € gesamt
    expect(wrapper.find('.avg-badge').text()).toContain('7,50')
  })

  it('zeigt Detail-KPIs bei Auswahl eines Monats', async () => {
    const wrapper = mount(MonthlyTrend, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.detail-hint').exists()).toBe(true)

    const columns = wrapper.findAll('.bar-column')
    // Monat 1 (2026-06) mit +30 € Ersparnis
    await columns[0].trigger('click')

    expect(wrapper.find('.detail-kpis').exists()).toBe(true)
    const kpis = wrapper.findAll('.detail-kpi')
    expect(kpis.length).toBe(5)
    // Savings-KPI zeigt +30,00 €
    expect(wrapper.find('.detail-kpi .detail-val.val-good').text()).toContain('+30,00')
  })

  it('zeigt positiven und negativen Balken farblich korrekt', async () => {
    const wrapper = mount(MonthlyTrend, { props: { visible: true } })
    await flushPromises()

    const columns = wrapper.findAll('.bar-column')
    // 2026-06 (savings 30 > 0) -> bar-good
    expect(columns[0].find('.bar-fill').classes()).toContain('bar-good')
    // 2026-07 (savings -15 < 0) -> bar-over
    expect(columns[1].find('.bar-fill').classes()).toContain('bar-over')
  })

  it('zeigt leeren Zustand ohne Monate', async () => {
    mockGetTrend.mockResolvedValueOnce({
      months: [],
      summary: {
        monthCount: 0,
        totalBudget: 0,
        totalSpent: 0,
        totalSaved: 0,
        avgSavings: 0,
        bestSavingsMonth: '',
      },
    })

    const wrapper = mount(MonthlyTrend, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.empty-state').exists()).toBe(true)
    expect(wrapper.find('.bar-column').exists()).toBe(false)
  })

  it('zeigt Fehler-Zustand bei API-Fehler mit Retry', async () => {
    mockGetTrend.mockRejectedValueOnce(new Error('Netzwerkfehler'))

    const wrapper = mount(MonthlyTrend, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.error-state').exists()).toBe(true)

    // Retry lädt erneut
    mockGetTrend.mockResolvedValueOnce(makeTrend())
    await wrapper.find('.retry-btn').trigger('click')
    await flushPromises()

    expect(mockGetTrend).toHaveBeenCalledTimes(2)
    expect(wrapper.find('.trend-month-title').exists()).toBe(true)
  })
})
