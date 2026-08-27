import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import AuthModal from '../components/AuthModal.vue'
import { useAuth } from '../composables/useAuth'

describe('AuthModal component', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
    const auth = useAuth()
    auth.logout()
  })

  it('renders guest login view when not logged in', () => {
    const wrapper = mount(AuthModal, {
      props: {
        visible: true,
      },
    })

    expect(wrapper.text()).toContain('Anmelden / Registrieren')
    expect(wrapper.text()).toContain('Passwortlos & Sicher')
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('.auth-primary-btn').exists()).toBe(true)
  })

  it('disables request button when email is invalid', async () => {
    const wrapper = mount(AuthModal, {
      props: {
        visible: true,
      },
    })

    const input = wrapper.find('input[type="email"]')
    const btn = wrapper.find('.auth-primary-btn')

    await input.setValue('invalid-email')
    expect(btn.attributes('disabled')).toBeDefined()

    await input.setValue('valid@example.com')
    expect(btn.attributes('disabled')).toBeUndefined()
  })

  it('submits magic link request on valid email', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        status: 'ok',
        message: 'Link gesendet',
        debugLink: 'http://localhost:5173/?auth_token=tok123',
      }),
    })
    globalThis.fetch = fetchMock

    const wrapper = mount(AuthModal, {
      props: {
        visible: true,
      },
    })

    const input = wrapper.find('input[type="email"]')
    await input.setValue('test@example.com')
    await wrapper.find('.auth-primary-btn').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Link ist unterwegs!')
    expect(wrapper.text()).toContain('test@example.com')
    expect(wrapper.find('.dev-login-btn').exists()).toBe(true)
  })

  it('renders user details and logout button when logged in', () => {
    const auth = useAuth()
    auth.user.value = {
      id: 'usr-12345678',
      email: 'member@test.de',
      createdAt: '2026-08-01T12:00:00Z',
      lastLoginAt: '2026-08-27T12:00:00Z',
      defaultMonthlyBudget: 450,
      defaultPeriodDays: 30,
      theme: 'emerald',
      isActive: true,
    }

    const wrapper = mount(AuthModal, {
      props: {
        visible: true,
      },
    })

    expect(wrapper.text()).toContain('Konto & Cloud-Sync')
    expect(wrapper.text()).toContain('member@test.de')
    expect(wrapper.text()).toContain('🟢 Cloud-Sync aktiv')
    expect(wrapper.find('.auth-logout-btn').exists()).toBe(true)
  })
})
