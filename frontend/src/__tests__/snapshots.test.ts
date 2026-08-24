import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppHeader from '../components/AppHeader.vue'
import MonthProgress from '../components/MonthProgress.vue'
import BudgetDisplay from '../components/BudgetDisplay.vue'
import RecentExpenses from '../components/RecentExpenses.vue'
import Numpad from '../components/Numpad.vue'

describe('Snapshot: AppHeader', () => {
  it('standard', () => {
    const wrapper = mount(AppHeader)
    expect(wrapper.html()).toMatchSnapshot()
  })
})

describe('Snapshot: MonthProgress', () => {
  it('tag 17/31', () => {
    const wrapper = mount(MonthProgress, { props: { day: 17, monthDays: 31 } })
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('tag 1/28 (februar)', () => {
    const wrapper = mount(MonthProgress, { props: { day: 1, monthDays: 28 } })
    expect(wrapper.html()).toMatchSnapshot()
  })
})

describe('Snapshot: BudgetDisplay', () => {
  it('gruen mit ersparnis', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 19.24, baseBudget: 14.52, savings: 61.36, color: 'green' }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('weiss bei 0 ersparnis', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 14.52, baseBudget: 14.52, savings: 0, color: 'white' }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('rot bei minus', () => {
    const wrapper = mount(BudgetDisplay, {
      props: { currentBudget: 11.55, baseBudget: 14.52, savings: -38.64, color: 'red' }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
})

describe('Snapshot: RecentExpenses', () => {
  it('leer', () => {
    const wrapper = mount(RecentExpenses, { props: { expenses: [] } })
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('mit ausgaben', () => {
    const wrapper = mount(RecentExpenses, {
      props: {
        expenses: [
          { id: 'e1', periodId: '2026-08', amount: 8.5, note: 'Frühstück', createdAt: '2026-08-18T06:35:43Z' },
          { id: 'e2', periodId: '2026-08', amount: 5.5, note: '', createdAt: '2026-08-18T07:00:00Z' },
        ]
      }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
})

describe('Snapshot: Numpad', () => {
  it('versteckt', () => {
    const wrapper = mount(Numpad, { props: { visible: false } })
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('sichtbar', () => {
    const wrapper = mount(Numpad, { props: { visible: true } })
    expect(wrapper.html()).toMatchSnapshot()
  })
})
