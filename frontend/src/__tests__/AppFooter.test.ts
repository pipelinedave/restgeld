import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppFooter from '../components/AppFooter.vue'

describe('AppFooter', () => {
  it('rendert den Tagline und die Version', () => {
    const wrapper = mount(AppFooter)
    expect(wrapper.find('.footer-tagline').text()).toContain('Track daily. Stay in budget.')
    expect(wrapper.find('.footer-version').text()).toBe('v0.1.0')
  })
})
