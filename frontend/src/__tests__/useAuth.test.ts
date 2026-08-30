import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuth } from '../composables/useAuth'

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

  it('verifyToken sets session on success', async () => {
    const auth = useAuth()
    const mockUser = {
      id: 'u-1',
      email: 'test@example.com',
      createdAt: new Date().toISOString(),
      lastLoginAt: new Date().toISOString(),
      defaultMonthlyBudget: 450,
      defaultPeriodDays: 30,
      theme: 'emerald',
      isActive: true,
    }

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
    auth.verifyToken = vi.fn().mockImplementation(async () => {
      auth.user.value = {
        id: 'u-1',
        email: 'test@example.com',
        createdAt: '',
        lastLoginAt: '',
        defaultMonthlyBudget: 450,
        defaultPeriodDays: 30,
        theme: '',
        isActive: true,
      }
      auth.authToken.value = 'sess-token'
      return true
    })
    await auth.verifyToken('token')

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
