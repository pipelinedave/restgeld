import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppHeader from '../components/AppHeader.vue'

describe('AppHeader', () => {
  it('rendert den App-Titel "restgeld"', () => {
    const wrapper = mount(AppHeader)
    expect(wrapper.find('.brand-title').text()).toBe('restgeld')
  })

  it('emittet open-settings bei klick auf settings-button', async () => {
    const wrapper = mount(AppHeader)
    await wrapper.find('.settings-btn').trigger('click')
    expect(wrapper.emitted('open-settings')).toBeTruthy()
  })

  it('zeigt Offline-Badge an wenn isOffline=true', () => {
    const wrapper = mount(AppHeader, { props: { isOffline: true } })
    expect(wrapper.find('.status-offline').exists()).toBe(true)
    expect(wrapper.find('.status-offline').text()).toContain('Offline')
  })

  it('zeigt Sync-Badge an wenn pendingSyncCount > 0', () => {
    const wrapper = mount(AppHeader, { props: { pendingSyncCount: 3 } })
    expect(wrapper.find('.status-syncing').exists()).toBe(true)
    expect(wrapper.find('.status-syncing').text()).toContain('3 ungesynct')
  })

  it('oeffnet Health-Popover bei Klick auf Status-Badge', async () => {
    const wrapper = mount(AppHeader, { props: { isOffline: false } })
    expect(wrapper.find('.header-popover').exists()).toBe(false)
    await wrapper.find('.status-badge').trigger('click')
    expect(wrapper.find('.header-popover').exists()).toBe(true)
    expect(wrapper.find('.popover-title').text()).toContain('System- & Sync-Status')
  })

  it('oeffnet Streak-Popover bei Klick auf Flammen-Button', async () => {
    const wrapper = mount(AppHeader, {
      props: {
        streak: { currentStreak: 5, longestStreak: 12, noSpendDays: 3, underBudgetDays: 7 },
      },
    })
    expect(wrapper.find('.streak-btn').text()).toContain('5')
    await wrapper.find('.streak-btn').trigger('click')
    expect(wrapper.find('.streak-popover').exists()).toBe(true)
    expect(wrapper.find('.streak-big-val').text()).toContain('5 Tage')
  })

  it('oeffnet Prognose-Popover bei Klick auf Glaskugel-Button', async () => {
    const wrapper = mount(AppHeader, {
      props: {
        projection: { projectedSavings: 42.5, projectedTotalSpent: 400, avgDailySpend: 13.5, status: 'saving' },
        baseBudget: 15,
      },
    })
    await wrapper.find('.projection-btn').trigger('click')
    expect(wrapper.find('.projection-popover').exists()).toBe(true)
    expect(wrapper.find('.proj-hero-val').text()).toContain('+42,50')
  })
})
