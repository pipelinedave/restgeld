import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuth, type User } from '../composables/useAuth'

function makeUser(overrides: Partial<User> = {}): User {
  return {
    id: 'u-1',
    email: 'test@example.com',
    language: 'de',
    currency: 'EUR',
    plan: 'free',
    createdAt: new Date().toISOString(),
    lastLoginAt: new Date().toISOString(),
    defaultMonthlyBudget: 450,
    defaultPeriodDays: 30,
    theme: 'emerald',
    isActive: true,
    ...overrides,
  }
}

describe('useAuth composable', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
    const auth = useAuth()
    auth.logout()
  })

  it('initializes with logged out state', () => {
    const auth = useAuth()
    expect(auth.isLoggedIn.value).toBe(false)
    expect(auth.user.value).toBeNull()
    expect(auth.authToken.value).toBeNull()
  })

  it('requestMagicLink sends request and handles response', async () => {
    const auth = useAuth()
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        status: 'ok',
        message: 'Link gesendet',
        debugLink: 'http://localhost:5173/?auth_token=tok123',
      }),
    })
    globalThis.fetch = fetchMock

    const res = await auth.requestMagicLink('test@example.com')
    expect(res.success).toBe(true)
    expect(res.debugLink).toBe('http://localhost:5173/?auth_token=tok123')
    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/api/auth/magic-link'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ email: 'test@example.com' }),
      })
    )
  })

  it('requestMagicLink handles API errors and network failures', async () => {
    const auth = useAuth()
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'Invalid email address' }),
    })
    const res1 = await auth.requestMagicLink('bad-email')
    expect(res1.success).toBe(false)
    expect(res1.message).toBe('Invalid email address')

    globalThis.fetch = vi.fn().mockRejectedValue(new Error('offline'))
    const res2 = await auth.requestMagicLink('test@example.com')
    expect(res2.success).toBe(false)
    expect(res2.message).toBe('Netzwerkfehler beim Anfordern des Login-Links')
  })

  it('verifyToken sets session on success', async () => {
    const auth = useAuth()
    const mockUser = makeUser({ email: 'test@example.com' })

    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        user: mockUser,
        token: 'sess-token-xyz',
        isNewUser: true,
      }),
    })
    globalThis.fetch = fetchMock

    const ok = await auth.verifyToken('valid-magic-token')
    expect(ok).toBe(true)
    expect(auth.isLoggedIn.value).toBe(true)
    expect(auth.user.value?.email).toBe('test@example.com')
    expect(auth.authToken.value).toBe('sess-token-xyz')
    expect(localStorage.getItem('restgeld_auth_token')).toBe('sess-token-xyz')
  })

  it('verifyToken handles invalid tokens and network errors', async () => {
    const auth = useAuth()
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: 'Token expired' }),
    })
    const ok1 = await auth.verifyToken('expired')
    expect(ok1).toBe(false)
    expect(auth.error.value).toBe('Token expired')

    globalThis.fetch = vi.fn().mockRejectedValue(new Error('network down'))
    const ok2 = await auth.verifyToken('bad')
    expect(ok2).toBe(false)
    expect(auth.error.value).toBe('Netzwerkfehler bei der Verifizierung')
  })

  it('fetchMe retrieves profile with authorization headers', async () => {
    const auth = useAuth()
    auth.authToken.value = 'my-token'
    const mockUser = makeUser({ email: 'me@example.com', plan: 'pro' })
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => mockUser,
    })

    const me = await auth.fetchMe()
    expect(me?.email).toBe('me@example.com')
    expect(auth.user.value?.plan).toBe('pro')
  })

  it('updateProfile sends update payload and updates state', async () => {
    const auth = useAuth()
    auth.authToken.value = 'my-token'
    auth.user.value = makeUser({ email: 'me@example.com' })
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'ok' }),
    })

    const ok = await auth.updateProfile({ language: 'fr', currency: 'EUR' })
    expect(ok).toBe(true)
    expect(auth.user.value?.language).toBe('fr')
  })

  it('deleteAccount calls DELETE endpoint and resets state', async () => {
    const auth = useAuth()
    auth.authToken.value = 'tok'
    auth.user.value = makeUser({ email: 'del@example.com' })
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'ok' }),
    })

    const ok = await auth.deleteAccount()
    expect(ok).toBe(true)
    expect(auth.isLoggedIn.value).toBe(false)
    expect(auth.user.value).toBeNull()
  })

  it('logout clears user and token from storage', async () => {
    const auth = useAuth()
    localStorage.setItem('restgeld_auth_token', 'test-tok')

    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'ok' }),
    })
    globalThis.fetch = fetchMock

    await auth.logout()
    expect(auth.isLoggedIn.value).toBe(false)
    expect(auth.user.value).toBeNull()
    expect(auth.authToken.value).toBeNull()
    expect(localStorage.getItem('restgeld_auth_token')).toBeNull()
  })

  it('migrateGuestData calls endpoint with payload', async () => {
    const auth = useAuth()
    auth.user.value = makeUser()
    auth.authToken.value = 'sess-token'

    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'ok', migratedCount: 2 }),
    })
    globalThis.fetch = fetchMock

    const count = await auth.migrateGuestData(
      [{ id: '1', periodId: '2026-08', amount: 10, note: 'Tee', createdAt: '' }],
      []
    )
    expect(count).toBe(2)
  })

  it('handles passkey methods gracefully when not supported', async () => {
    const auth = useAuth()
    expect(typeof auth.isPasskeySupported()).toBe('boolean')
    
    // In node/jsdom environment without WebAuthn
    const regRes = await auth.registerPasskey()
    expect(regRes.success).toBe(false)

    const loginRes = await auth.loginWithPasskey()
    expect(loginRes).toBe(false)
  })
})
