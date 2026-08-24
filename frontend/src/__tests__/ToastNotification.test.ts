import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ToastNotification from '../components/ToastNotification.vue'

describe('ToastNotification', () => {
  it('rendert nichts wenn visible=false', () => {
    const wrapper = mount(ToastNotification, {
      props: { visible: false, message: 'Test' },
    })
    expect(wrapper.find('.toast-container').exists()).toBe(false)
  })

  it('rendert Meldung und Standard-Erfolgs-Klasse wenn visible=true', () => {
    const wrapper = mount(ToastNotification, {
      props: { visible: true, message: '✓ 12,50 € gebucht', type: 'success' },
    })
    expect(wrapper.find('.toast-container').exists()).toBe(true)
    expect(wrapper.find('.toast-container').classes()).toContain('toast-success')
    expect(wrapper.find('.toast-message').text()).toBe('✓ 12,50 € gebucht')
  })

  it('rendert Fehlermeldung mit Fehler-Klasse', () => {
    const wrapper = mount(ToastNotification, {
      props: { visible: true, message: 'Fehler aufgetreten', type: 'error' },
    })
    expect(wrapper.find('.toast-container').classes()).toContain('toast-error')
    expect(wrapper.find('.toast-message').text()).toBe('Fehler aufgetreten')
  })
})
