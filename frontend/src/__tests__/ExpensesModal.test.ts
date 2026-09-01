import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ExpensesModal from '../components/ExpensesModal.vue'
import { useApi, type PaginatedExpenses } from '../composables/useApi'

vi.mock('../composables/useApi')
const mockApi = vi.mocked(useApi)

function makePaginated(overrides: Partial<PaginatedExpenses> = {}): PaginatedExpenses {
  return {
    items: [
      {
        id: 'exp-1',
        periodId: '2026-08',
        amount: 12.5,
        note: 'Mittagessen',
        createdAt: '2026-08-18T12:30:00Z',
      },
      {
        id: 'exp-2',
        periodId: '2026-08',
        amount: 3.2,
        note: '',
        createdAt: '2026-08-18T08:15:00Z',
      },
    ],
    total: 8,
    page: 1,
    limit: 6,
    totalPages: 2,
    ...overrides,
  }
}

describe('ExpensesModal (Bottom Sheet Overlay)', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockApi.mockReturnValue({
      getExpenses: vi.fn().mockResolvedValue(makePaginated()),
      deleteExpense: vi.fn().mockResolvedValue({ status: 'ok' }),
    } as any)
  })

  it('rendert nichts wenn nicht sichtbar', () => {
    const wrapper = mount(ExpensesModal, { props: { visible: false } })
    expect(wrapper.find('.sheet-overlay').exists()).toBe(false)
  })

  it('laedt und rendert ausgaben wenn sichtbar', async () => {
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(mockApi().getExpenses).toHaveBeenCalledWith(1, 12)
    expect(wrapper.find('.bottom-sheet h2').text()).toBe('Alle Ausgaben')
    expect(wrapper.find('.badge').text()).toBe('8')
    expect(wrapper.find('.drag-handle-pill').exists()).toBe(true)
    expect(wrapper.text()).toContain('Mittagessen')
    expect(wrapper.text()).toContain('Ausgabe') // Fallback für leere Notiz
    expect(wrapper.text()).toContain('-12,50')
    expect(wrapper.text()).toContain('-3,20')
  })

  it('zeigt leere nachricht wenn keine ausgaben vorhanden', async () => {
    mockApi.mockReturnValue({
      getExpenses: vi.fn().mockResolvedValue({ items: [], total: 0, page: 1, limit: 12, totalPages: 1 }),
      deleteExpense: vi.fn(),
    } as any)

    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.empty-state').text()).toContain('Keine Ausgaben in dieser Periode vorhanden')
  })

  it('laedt mehr daten bei infinite scroll / loadNextPage', async () => {
    const getExpensesMock = vi.fn()
      .mockResolvedValueOnce(makePaginated({ page: 1, totalPages: 2 }))
      .mockResolvedValueOnce(makePaginated({ page: 2, totalPages: 2, items: [{ id: 'exp-3', periodId: '2026-08', amount: 5.0, note: 'Kaffee', createdAt: '2026-08-19T10:00:00Z' }] }))

    mockApi.mockReturnValue({
      getExpenses: getExpensesMock,
      deleteExpense: vi.fn(),
    } as any)

    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.load-more-hint').exists()).toBe(true)

    // Trigger more
    await wrapper.find('.load-more-hint').trigger('click')
    await flushPromises()

    expect(getExpensesMock).toHaveBeenCalledWith(2, 12)
    expect(wrapper.text()).toContain('Kaffee')
  })

  it('loescht ausgabe und emittet expense-deleted', async () => {
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    const deleteBtn = wrapper.find('.delete-btn')
    await deleteBtn.trigger('click')
    await flushPromises()

    expect(mockApi().deleteExpense).toHaveBeenCalledWith('exp-1')
    expect(wrapper.emitted('expense-deleted')).toBeTruthy()
    expect(wrapper.emitted('expense-deleted')![0]).toEqual(['exp-1'])
  })

  it('schliesst overlay bei klick auf schliessen-button', async () => {
    vi.useFakeTimers()
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    await wrapper.find('.close-btn').trigger('click')
    vi.advanceTimersByTime(300)
    expect(wrapper.emitted('close')).toBeTruthy()
    vi.useRealTimers()
  })

  it('schliesst overlay bei klick auf overlay-hintergrund', async () => {
    vi.useFakeTimers()
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    await wrapper.find('.sheet-overlay').trigger('click')
    vi.advanceTimersByTime(300)
    expect(wrapper.emitted('close')).toBeTruthy()
    vi.useRealTimers()
  })
})
