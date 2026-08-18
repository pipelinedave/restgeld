import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import RecentExpenses from '../components/RecentExpenses.vue'
import type { Expense } from '../composables/useApi'

function makeExpense(overrides: Partial<Expense> = {}): Expense {
  return {
    id: 'test-1',
    periodId: '2026-08',
    amount: 8.5,
    note: 'Frühstück',
    createdAt: '2026-08-18T06:35:43Z',
    ...overrides
  }
}

describe('RecentExpenses', () => {
  it('rendert ausgaben-liste', () => {
    const expenses = [makeExpense(), makeExpense({ id: 'test-2', amount: 5.5, note: 'Kaffee' })]
    const wrapper = mount(RecentExpenses, { props: { expenses } })
    expect(wrapper.text()).toContain('Frühstück')
    expect(wrapper.text()).toContain('Kaffee')
    expect(wrapper.text()).toContain('8,50')
    expect(wrapper.text()).toContain('5,50')
  })

  it('zeigt leere nachricht bei keiner ausgaben', () => {
    const wrapper = mount(RecentExpenses, { props: { expenses: [] } })
    expect(wrapper.text()).toContain('Noch keine Ausgaben heute')
  })

  it('emittet delete beim loeschen-button', async () => {
    const expenses = [makeExpense()]
    const wrapper = mount(RecentExpenses, { props: { expenses } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('delete')).toBeTruthy()
    expect(wrapper.emitted('delete')![0]).toEqual(['test-1'])
  })

  it('zeigt notiz statt "Ausgabe" wenn vorhanden', () => {
    const expenses = [makeExpense()]
    const wrapper = mount(RecentExpenses, { props: { expenses } })
    expect(wrapper.find('.expense-note').text()).toBe('Frühstück')
  })

  it('zeigt fallback "Ausgabe" wenn keine notiz', () => {
    const expenses = [makeExpense({ note: '' })]
    const wrapper = mount(RecentExpenses, { props: { expenses } })
    expect(wrapper.find('.expense-note').text()).toBe('Ausgabe')
  })

  it('zeigt betrag mit vorzeichen', () => {
    const expenses = [makeExpense()]
    const wrapper = mount(RecentExpenses, { props: { expenses } })
    expect(wrapper.text()).toContain('-8,50')
  })
})
