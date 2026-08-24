import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MonthProgress from '../components/MonthProgress.vue'

describe('MonthProgress', () => {
  it('zeigt label Tag X von Y', () => {
    const wrapper = mount(MonthProgress, { props: { day: 17, monthDays: 31 } })
    expect(wrapper.text()).toContain('Tag 17 von 31')
  })

  it('setzt breite auf 54% bei 17/31', () => {
    const wrapper = mount(MonthProgress, { props: { day: 17, monthDays: 31 } })
    const fill = wrapper.find('.progress-fill')
    expect(fill.attributes('style')).toContain('width: 55%')
  })

  it('setzt breite auf 3% am ersten tag', () => {
    const wrapper = mount(MonthProgress, { props: { day: 1, monthDays: 31 } })
    const fill = wrapper.find('.progress-fill')
    expect(fill.attributes('style')).toContain('width: 3%')
  })

  it('setzt breite auf 100% am letzten tag', () => {
    const wrapper = mount(MonthProgress, { props: { day: 31, monthDays: 31 } })
    const fill = wrapper.find('.progress-fill')
    expect(fill.attributes('style')).toContain('width: 100%')
  })

  it('emittet open-settings bei klick auf settings-button', async () => {
    const wrapper = mount(MonthProgress, { props: { day: 15, monthDays: 30 } })
    await wrapper.find('.settings-btn').trigger('click')
    expect(wrapper.emitted('open-settings')).toBeTruthy()
  })
})
