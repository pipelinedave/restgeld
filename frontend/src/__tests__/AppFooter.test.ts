import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppFooter from '../components/AppFooter.vue'

describe('AppFooter', () => {
  it('rendert den Tagline, Energy-Referenz und den Commit/Dev-Hash', () => {
    const wrapper = mount(AppFooter)
    expect(wrapper.find('.footer-tagline').text()).toContain('Track daily. Stay in budget.')
    expect(wrapper.find('.footer-energy').text()).toContain('lowlifehigh.tech')
    expect(wrapper.find('.commit-badge').exists()).toBe(true)
  })
})
