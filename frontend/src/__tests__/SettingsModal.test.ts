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

  it('renders correctly when visible is true and populates budget and days', () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true, currentMonthlyBudget: 600, currentMonthDays: 30 }
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    const input = wrapper.find<HTMLInputElement>('input#monthly-budget-input')
    expect(input.exists()).toBe(true)
    expect(input.element.value).toBe('600')
    const daysInput = wrapper.find<HTMLInputElement>('input#period-days-input')
    expect(daysInput.exists()).toBe(true)
    expect(daysInput.element.value).toBe('30')
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

  it('emits update-budget event with days when clicking save', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true, currentMonthlyBudget: 600, currentMonthDays: 30 }
    })
    const input = wrapper.find<HTMLInputElement>('input#monthly-budget-input')
    await input.setValue(750)
    const daysInput = wrapper.find<HTMLInputElement>('input#period-days-input')
    await daysInput.setValue(14)
    await wrapper.find('.save-btn').trigger('click')

    expect(wrapper.emitted('update-budget')).toBeTruthy()
    expect(wrapper.emitted('update-budget')?.[0]).toEqual([750, 14])
  })

  it('requires confirmation before emitting new-period', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })

    const initialDangerBtn = wrapper.find('.danger-btn')
    expect(initialDangerBtn.text()).toContain('Neue Periode ab heute starten')
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

  it('renders export buttons for CSV and JSON', () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })

    const backupBtns = wrapper.findAll('.backup-btn')
    expect(backupBtns.length).toBe(2)
    expect(backupBtns[0].text()).toContain('CSV (Excel)')
    expect(backupBtns[1].text()).toContain('JSON Backup')
  })

  it('renders account section and emits open-auth on button click', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true }
    })

    const accountBtn = wrapper.find('.account-btn')
    expect(accountBtn.exists()).toBe(true)
    await accountBtn.trigger('click')
    expect(wrapper.emitted('open-auth')).toBeTruthy()
  })

  it('emits open-about and open-status when clicking info triggers', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true },
    })

    const aboutBtns = wrapper.findAll('.about-trigger-btn')
    expect(aboutBtns.length).toBeGreaterThanOrEqual(2)

    await aboutBtns[0].trigger('click')
    expect(wrapper.emitted('open-about')).toBeTruthy()

    await aboutBtns[1].trigger('click')
    expect(wrapper.emitted('open-status')).toBeTruthy()
  })

  it('allows choosing theme presets', async () => {
    const wrapper = mount(SettingsModal, {
      props: { visible: true },
    })

    const themeButtons = wrapper.findAll('.theme-color-btn')
    expect(themeButtons.length).toBeGreaterThanOrEqual(5)
    await themeButtons[1].trigger('click')
    expect(themeButtons[1].exists()).toBe(true)
  })
})
