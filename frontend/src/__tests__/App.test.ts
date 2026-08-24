import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import App from '../App.vue'
import { useApi, type BudgetData } from '../composables/useApi'

vi.mock('../composables/useApi')

const mockApi = vi.mocked(useApi)

function makeBudget(overrides: Partial<BudgetData> = {}): BudgetData {
  return {
    day: 18,
    monthDays: 31,
    baseBudget: 14.52,
    currentBudget: 19.24,
    savings: 61.36,
    color: 'green',
    periodId: '2026-08',
    expenses: [],
    ...overrides,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  mockApi.mockReturnValue({
    getBudget: vi.fn().mockResolvedValue(makeBudget()),
    addExpense: vi.fn().mockResolvedValue({ id: 'new', periodId: '', amount: 0, note: '', createdAt: '' }),
    deleteExpense: vi.fn().mockResolvedValue({ status: 'ok' }),
    newPeriod: vi.fn(),
    updateBudget: vi.fn(),
    getExpenses: vi.fn(),
  } as any)
})

describe('App', () => {
  it('laedt budget beim mount', async () => {
    const wrapper = mount(App)
    await flushPromises()
    expect(mockApi().getBudget).toHaveBeenCalledOnce()
    expect(wrapper.text()).toContain('19,24')
  })

  it('zeigt budget und ersparnis', async () => {
    const wrapper = mount(App)
    await flushPromises()
    expect(wrapper.text()).toContain('Tag 18 von 31')
    expect(wrapper.text()).toContain('19,24')
    expect(wrapper.text()).toContain('61,36')
  })

  it('oeffnet numpad bei klick auf ausgabe-button', async () => {
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.add-btn').trigger('click')
    expect(wrapper.find('.numpad-overlay').exists()).toBe(true)
  })

  it('bucht ausgabe via numpad confirm', async () => {
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.add-btn').trigger('click')
    const buttons = wrapper.findAll('.numpad-btn')
    await buttons[6].trigger('click') // 1
    await buttons[10].trigger('click') // 0
    await buttons[9].trigger('click') // ,
    await buttons[11].trigger('click') // 0
    await buttons[13].trigger('click') // bestätigen
    await flushPromises()
    const noteInput = wrapper.find('.note-input')
    await noteInput.setValue('Kaffee')
    await wrapper.findAll('.note-actions button')[1].trigger('click')
    await flushPromises()
    expect(mockApi().addExpense).toHaveBeenCalledWith(10, 'Kaffee')
    expect(mockApi().getBudget).toHaveBeenCalledTimes(2)
  })

  it('loescht ausgabe', async () => {
    const api = mockApi()
    ;(api.getBudget as any).mockResolvedValue(makeBudget({
      expenses: [{ id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Test', createdAt: '2026-08-18T06:35:43Z' }]
    }))
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.delete-btn').trigger('click')
    await flushPromises()
    expect(mockApi().deleteExpense).toHaveBeenCalledWith('e1')
  })

  it('zeigt ladestatus initial', () => {
    const wrapper = mount(App)
    expect(wrapper.text()).toContain('Lade')
  })

  it('oeffnet settings bei klick auf settings-button', async () => {
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.settings-btn').trigger('click')
    expect(wrapper.find('.modal-content h2').text()).toBe('Einstellungen')
  })

  it('aktualisiert budget ueber settings modal', async () => {
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.settings-btn').trigger('click')
    const input = wrapper.find<HTMLInputElement>('input#monthly-budget-input')
    await input.setValue(800)
    await wrapper.find('.setting-section .action-btn').trigger('click')
    await flushPromises()
    expect(mockApi().updateBudget).toHaveBeenCalledWith(800)
    expect(mockApi().getBudget).toHaveBeenCalledTimes(2)
  })

  it('setzt periode zurueck ueber settings modal', async () => {
    const wrapper = mount(App)
    await flushPromises()
    await wrapper.find('.settings-btn').trigger('click')
    await wrapper.find('.danger-btn').trigger('click') // Erstes "Neue Periode starten"
    await wrapper.find('.danger-btn.confirm').trigger('click') // Bestätigen
    await flushPromises()
    expect(mockApi().newPeriod).toHaveBeenCalledOnce()
    expect(mockApi().getBudget).toHaveBeenCalledTimes(2)
  })

  it('oeffnet ausgaben-historie modal bei klick auf alle anzeigen', async () => {
    const api = mockApi()
    ;(api.getBudget as any).mockResolvedValue(makeBudget({
      expenses: [{ id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Test', createdAt: '2026-08-18T06:35:43Z' }]
    }))
    ;(api.getExpenses as any).mockResolvedValue({
      items: [{ id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Test', createdAt: '2026-08-18T06:35:43Z' }],
      total: 1,
      page: 1,
      limit: 6,
      totalPages: 1
    })

    const wrapper = mount(App)
    await flushPromises()

    const showAllBtn = wrapper.find('.show-all-btn')
    expect(showAllBtn.exists()).toBe(true)
    await showAllBtn.trigger('click')
    await flushPromises()

    expect(wrapper.find('.modal-content h2').text()).toBe('Alle Ausgaben')
    expect(api.getExpenses).toHaveBeenCalledWith(1, 6)
  })

  it('zeigt fehler bei api-fehler', async () => {
    mockApi.mockReturnValue({
      getBudget: vi.fn().mockRejectedValue(new Error('kaputt')),
      addExpense: vi.fn(),
      deleteExpense: vi.fn(),
      newPeriod: vi.fn(),
      updateBudget: vi.fn(),
      getExpenses: vi.fn(),
    } as any)
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    mount(App)
    await flushPromises()
    expect(consoleSpy).toHaveBeenCalled()
    consoleSpy.mockRestore()
  })
})
