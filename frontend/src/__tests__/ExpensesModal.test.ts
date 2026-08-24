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

describe('ExpensesModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockApi.mockReturnValue({
      getExpenses: vi.fn().mockResolvedValue(makePaginated()),
      deleteExpense: vi.fn().mockResolvedValue({ status: 'ok' }),
    } as any)
  })

  it('rendert nichts wenn nicht sichtbar', () => {
    const wrapper = mount(ExpensesModal, { props: { visible: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('laedt und rendert ausgaben wenn sichtbar', async () => {
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(mockApi().getExpenses).toHaveBeenCalledWith(1, 6)
    expect(wrapper.find('.modal-content h2').text()).toBe('Alle Ausgaben')
    expect(wrapper.find('.badge').text()).toBe('8')
    expect(wrapper.text()).toContain('Mittagessen')
    expect(wrapper.text()).toContain('Ausgabe') // Fallback für leere Notiz
    expect(wrapper.text()).toContain('-12,50')
    expect(wrapper.text()).toContain('-3,20')
  })

  it('zeigt leere nachricht wenn keine ausgaben vorhanden', async () => {
    mockApi.mockReturnValue({
      getExpenses: vi.fn().mockResolvedValue({ items: [], total: 0, page: 1, limit: 6, totalPages: 1 }),
      deleteExpense: vi.fn(),
    } as any)

    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.empty-state').text()).toContain('Keine Ausgaben in dieser Periode vorhanden')
    expect(wrapper.find('.pagination-bar').exists()).toBe(false)
  })

  it('blaettert zur naechsten und vorherigen seite', async () => {
    const getExpensesMock = vi.fn()
      .mockResolvedValueOnce(makePaginated({ page: 1, totalPages: 2 }))
      .mockResolvedValueOnce(makePaginated({ page: 2, totalPages: 2 }))

    mockApi.mockReturnValue({
      getExpenses: getExpensesMock,
      deleteExpense: vi.fn(),
    } as any)

    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    expect(wrapper.find('.pagination-info').text()).toBe('Seite 1 von 2')
    const nextBtn = wrapper.find('button[aria-label="Nächste Seite"]')
    expect(nextBtn.attributes('disabled')).toBeUndefined()

    // Klick auf Weiter
    await nextBtn.trigger('click')
    await flushPromises()

    expect(getExpensesMock).toHaveBeenCalledWith(2, 6)
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

  it('schliesst modal bei klick auf schliessen-button', async () => {
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('schliesst modal bei klick auf overlay', async () => {
    const wrapper = mount(ExpensesModal, { props: { visible: true } })
    await flushPromises()

    await wrapper.find('.modal-overlay').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
