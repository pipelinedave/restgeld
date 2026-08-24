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
    expect(wrapper.find('.offline-badge').exists()).toBe(true)
    expect(wrapper.find('.offline-badge').text()).toContain('Offline')
  })

  it('zeigt Sync-Badge an wenn pendingSyncCount > 0', () => {
    const wrapper = mount(AppHeader, { props: { pendingSyncCount: 3 } })
    expect(wrapper.find('.syncing-badge').exists()).toBe(true)
    expect(wrapper.find('.syncing-badge').text()).toContain('3 ungesynct')
  })
})
