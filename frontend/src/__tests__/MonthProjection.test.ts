import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MonthProjection from '../components/MonthProjection.vue'

describe('MonthProjection', () => {
  it('rendert nichts wenn projection prop fehlt', () => {
    const wrapper = mount(MonthProjection, { props: { projection: undefined } })
    expect(wrapper.find('.projection-strip').exists()).toBe(false)
  })

  it('rendert Spar-Prognose im Plus korrekt', () => {
    const wrapper = mount(MonthProjection, {
      props: {
        projection: {
          projectedSavings: 42.5,
          projectedTotalSpent: 257.5,
          avgDailySpend: 8.5,
          status: 'saving',
        },
      },
    })

    expect(wrapper.find('.projection-strip').exists()).toBe(true)
    expect(wrapper.find('.projection-strip').classes()).toContain('saving')
    expect(wrapper.find('.projection-highlight').text()).toContain('+42,50 €')
    expect(wrapper.find('.projection-tag').text()).toBe('gespart')
    expect(wrapper.find('.projection-rate').text()).toContain('Ø 8,50 € / Tag')
  })

  it('rendert Defizit-Prognose im Minus korrekt', () => {
    const wrapper = mount(MonthProjection, {
      props: {
        projection: {
          projectedSavings: -15.0,
          projectedTotalSpent: 315.0,
          avgDailySpend: 10.5,
          status: 'deficit',
        },
      },
    })

    expect(wrapper.find('.projection-strip').classes()).toContain('deficit')
    expect(wrapper.find('.projection-highlight').text()).toContain('15,00 €')
    expect(wrapper.find('.projection-tag').text()).toBe('Defizit')
    expect(wrapper.find('.projection-rate').text()).toContain('Ø 10,50 € / Tag')
  })
})
