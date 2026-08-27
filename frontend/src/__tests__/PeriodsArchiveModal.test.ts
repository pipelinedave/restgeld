import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PeriodsArchiveModal from '../components/PeriodsArchiveModal.vue'

const mockGetPeriods = vi.fn()
const mockGetExpenses = vi.fn()
vi.mock('../composables/useApi', () => ({
  useApi: () => ({
    getPeriods: mockGetPeriods,
    getExpenses: mockGetExpenses,
  }),
}))

describe('PeriodsArchiveModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('rendert nichts wenn visible=false', () => {
    const wrapper = mount(PeriodsArchiveModal, { props: { visible: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('laedt und rendert Perioden-Historie mit Datumsbereich', async () => {
    mockGetPeriods.mockResolvedValueOnce([
      {
        id: '2026-08',
        startDate: '2026-08-01T00:00:00Z',
        monthDays: 31,
        baseBudget: 15.0,
        monthlyTotal: 465.0,
        totalSpent: 400.0,
        savings: 65.0,
        expenseCount: 14,
      },
    ])

    const wrapper = mount(PeriodsArchiveModal, { props: { visible: true } })
    expect(mockGetPeriods).toHaveBeenCalledTimes(1)

    await flushPromises()

    expect(wrapper.find('.period-card').exists()).toBe(true)
    expect(wrapper.find('.period-daterange').text()).toContain('31 Tage')
    expect(wrapper.find('.period-badge').classes()).toContain('saving')
    expect(wrapper.find('.period-badge').text()).toContain('+65,00 €')
  })

  it('klappt Abschlussbericht und Buchungen bei Klick auf', async () => {
    mockGetPeriods.mockResolvedValueOnce([
      {
        id: '2026-08',
        startDate: '2026-08-01T00:00:00Z',
        monthDays: 31,
        baseBudget: 15.0,
        monthlyTotal: 465.0,
        totalSpent: 12.5,
        savings: 452.5,
        expenseCount: 1,
      },
    ])
    mockGetExpenses.mockResolvedValueOnce({
      items: [{ id: 'exp-1', periodId: '2026-08', amount: 12.5, note: 'Kaffee', createdAt: '2026-08-05T10:00:00Z' }],
      total: 1,
      page: 1,
      limit: 100,
      totalPages: 1,
    })

    const wrapper = mount(PeriodsArchiveModal, { props: { visible: true } })
    await flushPromises()

    await wrapper.find('.period-card').trigger('click')
    await flushPromises()

    expect(mockGetExpenses).toHaveBeenCalledWith(1, 100, '2026-08')
    expect(wrapper.find('.report-section').exists()).toBe(true)
    expect(wrapper.find('.report-expense-item').text()).toContain('Kaffee')
    expect(wrapper.find('.report-expense-item').text()).toContain('-12,50 €')
  })

  it('emittet close bei klick auf Schliessen', async () => {
    mockGetPeriods.mockResolvedValueOnce([])
    const wrapper = mount(PeriodsArchiveModal, { props: { visible: true } })
    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
