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
})
