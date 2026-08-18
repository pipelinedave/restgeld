import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BudgetDisplay from '../components/BudgetDisplay.vue'

describe('BudgetDisplay', () => {
  it('rendert betrag', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 14.52, baseBudget: 14.52, savings: 0, color: 'white' }
    })
    expect(wrapper.text()).toContain('14,52')
  })

  it('farbe gruen bei ersparnis > 0', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 19.24, baseBudget: 14.52, savings: 61.36, color: 'green' }
    })
    expect(wrapper.find('.budget-amount').classes()).toContain('color-green')
  })

  it('farbe weiss bei ersparnis == 0', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 14.52, baseBudget: 14.52, savings: 0, color: 'white' }
    })
    expect(wrapper.find('.budget-amount').classes()).toContain('color-white')
  })

  it('farbe rot bei ersparnis < 0', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 11.55, baseBudget: 14.52, savings: -38.64, color: 'red' }
    })
    expect(wrapper.find('.budget-amount').classes()).toContain('color-red')
  })

  it('zeigt ersparnis an', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 19.24, baseBudget: 14.52, savings: 61.36, color: 'green' }
    })
    expect(wrapper.text()).toContain('61,36')
    expect(wrapper.text()).toContain('angespart')
  })

  it('zeigt kein ersparnis-text bei 0', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 14.52, baseBudget: 14.52, savings: 0, color: 'white' }
    })
    expect(wrapper.find('.savings').exists()).toBe(false)
  })

  it('zeigt minus bei negativer ersparnis', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 11.55, baseBudget: 14.52, savings: -38.64, color: 'red' }
    })
    expect(wrapper.text()).toContain('-38,64')
  })
})
