import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PeriodsArchiveModal from '../components/PeriodsArchiveModal.vue'

const mockGetPeriods = vi.fn()
vi.mock('../composables/useApi', () => ({
  useApi: () => ({
    getPeriods: mockGetPeriods,
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

  it('laedt und rendert Perioden-Historie', async () => {
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

    expect(wrapper.find('.period-item').exists()).toBe(true)
    expect(wrapper.find('.period-badge').classes()).toContain('saving')
    expect(wrapper.find('.period-badge').text()).toContain('+65,00 €')
  })

  it('emittet close bei klick auf Schliessen', async () => {
    mockGetPeriods.mockResolvedValueOnce([])
    const wrapper = mount(PeriodsArchiveModal, { props: { visible: true } })
    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
