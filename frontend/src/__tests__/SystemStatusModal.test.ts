import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import SystemStatusModal from '../components/SystemStatusModal.vue'
import { useI18n } from '../composables/useI18n'

describe('SystemStatusModal.vue', () => {
  beforeEach(() => {
    localStorage.clear()
    const { setLocale } = useI18n()
    setLocale('de')
  })

  it('does not render when visible is false', () => {
    const wrapper = mount(SystemStatusModal, {
      props: { visible: false },
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('renders modal header and telemetry grid when visible is true', async () => {
    const wrapper = mount(SystemStatusModal, {
      props: { visible: true },
    })
    await flushPromises()
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toContain('System Observability & Monitoring')
    expect(wrapper.find('.telemetry-grid').exists()).toBe(true)
  })

  it('renders services list with health status tags', async () => {
    const wrapper = mount(SystemStatusModal, {
      props: { visible: true },
    })
    await flushPromises()
    const serviceCards = wrapper.findAll('.service-card')
    expect(serviceCards.length).toBeGreaterThanOrEqual(4)
    expect(wrapper.text()).toContain('Core API Backend')
    expect(wrapper.text()).toContain('Auth Service')
    expect(wrapper.text()).toContain('Observability Service')
  })

  it('emits close event when clicking close button or overlay', async () => {
    const wrapper = mount(SystemStatusModal, {
      props: { visible: true },
    })
    await flushPromises()
    await wrapper.find('.close-btn').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()

    await wrapper.find('.modal-overlay').trigger('click')
    expect(wrapper.emitted('close')?.length).toBe(2)
  })

  it('triggers telemetry refresh on button click', async () => {
    const wrapper = mount(SystemStatusModal, {
      props: { visible: true },
    })
    await flushPromises()
    const refreshBtn = wrapper.find('.refresh-btn')
    expect(refreshBtn.exists()).toBe(true)
    await refreshBtn.trigger('click')
    await flushPromises()
    expect(wrapper.find('.services-list').exists()).toBe(true)
  })
})
