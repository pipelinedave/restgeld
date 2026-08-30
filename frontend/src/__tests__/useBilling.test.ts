import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useBilling } from '../composables/useBilling'
import { useAuth, type User } from '../composables/useAuth'

function makeUser(overrides: Partial<User> = {}): User {
  return {
    id: 'u1',
    email: 'test@example.com',
    createdAt: new Date().toISOString(),
    lastLoginAt: new Date().toISOString(),
    defaultMonthlyBudget: 450,
    defaultPeriodDays: 30,
    theme: 'emerald',
    language: 'de',
    currency: 'EUR',
    plan: 'free',
    isActive: true,
    ...overrides,
  }
}

describe('useBilling composable', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  it('erkennt Pro Plan korrekt', () => {
    const auth = useAuth()
    const billing = useBilling()

    expect(billing.isPro.value).toBe(false)

    auth.user.value = makeUser({ plan: 'pro' })
    expect(billing.isPro.value).toBe(true)
  })

  it('erstellt Checkout-Session erfolgreich', async () => {
    const billing = useBilling()
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ checkoutUrl: 'https://checkout.stripe.com/session123' }),
    })
    global.fetch = mockFetch

    const res = await billing.createCheckoutSession()
    expect(mockFetch).toHaveBeenCalled()
    expect(res.checkoutUrl).toBe('https://checkout.stripe.com/session123')
    expect(billing.error.value).toBeNull()
    expect(billing.isLoading.value).toBe(false)
  })

  it('behandelt API-Fehler bei Checkout-Session', async () => {
    const billing = useBilling()
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'Unauthorized user' }),
    })
    global.fetch = mockFetch

    const res = await billing.createCheckoutSession()
    expect(res.error).toBe('Unauthorized user')
    expect(billing.error.value).toBe('Unauthorized user')
    expect(billing.isLoading.value).toBe(false)
  })

  it('behandelt Netzwerkfehler bei Checkout-Session', async () => {
    const billing = useBilling()
    const mockFetch = vi.fn().mockRejectedValue(new Error('Network offline'))
    global.fetch = mockFetch

    const res = await billing.createCheckoutSession()
    expect(res.error).toBe('Netzwerkfehler beim Checkout')
    expect(billing.error.value).toBe('Netzwerkfehler beim Checkout')
  })

  it('oeffnet Customer Portal erfolgreich', async () => {
    const billing = useBilling()
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ portalUrl: 'https://billing.stripe.com/portal123' }),
    })
    global.fetch = mockFetch

    const res = await billing.openCustomerPortal()
    expect(res.portalUrl).toBe('https://billing.stripe.com/portal123')
    expect(billing.error.value).toBeNull()
  })

  it('behandelt API-Fehler beim Customer Portal', async () => {
    const billing = useBilling()
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'No active stripe customer' }),
    })
    global.fetch = mockFetch

    const res = await billing.openCustomerPortal()
    expect(res.error).toBe('No active stripe customer')
    expect(billing.error.value).toBe('No active stripe customer')
  })
})
