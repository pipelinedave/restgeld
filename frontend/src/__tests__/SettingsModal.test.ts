import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SettingsModal from '../components/SettingsModal.vue'

describe('SettingsModal', () => {
  it('does not render when visible is false', () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: false }
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('renders correctly when visible is true and populates budget', () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true, currentMonthlyBudget: 600 }
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    const input = wrapper.find<HTMLInputElement>('input#monthly-budget-input')
    expect(input.exists()).toBe(true)
    expect(input.element.value).toBe('600')
  })

  it('emits close event when clicking close button or backdrop', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })
    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()

    await wrapper.find('.modal-overlay').trigger('click')
    expect(wrapper.emitted('close')?.length).toBe(2)
  })

  it('emits update-budget event when clicking save with valid value', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true, currentMonthlyBudget: 600 }
    })
    const input = wrapper.find<HTMLInputElement>('input#monthly-budget-input')
    await input.setValue(750)
    await wrapper.find('.action-btn').trigger('click')

    expect(wrapper.emitted('update-budget')).toBeTruthy()
    expect(wrapper.emitted('update-budget')?.[0]).toEqual([750])
  })

  it('requires confirmation before emitting new-period', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })

    const initialDangerBtn = wrapper.find('.danger-btn')
    expect(initialDangerBtn.text()).toContain('Neue Periode starten')
    await initialDangerBtn.trigger('click')

    expect(wrapper.find('.confirm-box').exists()).toBe(true)
    expect(wrapper.emitted('new-period')).toBeFalsy()

    const confirmBtn = wrapper.find('.danger-btn.confirm')
    await confirmBtn.trigger('click')

    expect(wrapper.emitted('new-period')).toBeTruthy()
  })

  it('allows canceling reset confirmation', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })

    await wrapper.find('.danger-btn').trigger('click')
    expect(wrapper.find('.confirm-box').exists()).toBe(true)

    await wrapper.find('.cancel-btn').trigger('click')
    expect(wrapper.find('.confirm-box').exists()).toBe(false)
    expect(wrapper.emitted('new-period')).toBeFalsy()
  })
})
