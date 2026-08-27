import { ref, computed } from 'vue'
import type { Expense, Period } from './useApi'

export interface User {
  id: string
  email: string
  createdAt: string
  lastLoginAt: string
  defaultMonthlyBudget: number
  defaultPeriodDays: number
  theme: string
  isActive: boolean
}

const STORAGE_KEY_TOKEN = 'restgeld_auth_token'
const STORAGE_KEY_USER = 'restgeld_auth_user'

const user = ref<User | null>(null)
const authToken = ref<string | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

// Initial aus localStorage laden
if (typeof window !== 'undefined') {
  const savedToken = localStorage.getItem(STORAGE_KEY_TOKEN)
  if (savedToken) {
    authToken.value = savedToken
  }
  const savedUser = localStorage.getItem(STORAGE_KEY_USER)
  if (savedUser) {
    try {
      user.value = JSON.parse(savedUser)
    } catch {
      // ignore
    }
  }
}

const BASE = typeof window !== 'undefined' && import.meta.env.PROD ? window.location.origin : ''

export function useAuth() {
  const isLoggedIn = computed(() => !!user.value)

  function setSession(newUser: User | null, token: string | null) {
    user.value = newUser
    authToken.value = token

    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem(STORAGE_KEY_TOKEN, token)
      } else {
        localStorage.removeItem(STORAGE_KEY_TOKEN)
      }

      if (newUser) {
        localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(newUser))
      } else {
        localStorage.removeItem(STORAGE_KEY_USER)
      }
    }
  }

  function getAuthHeaders(): HeadersInit {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (authToken.value) {
      headers['Authorization'] = `Bearer ${authToken.value}`
    }
    return headers
  }

  async function requestMagicLink(email: string): Promise<{ success: boolean; message?: string; debugLink?: string }> {
    isLoading.value = true
    error.value = null

    try {
      const res = await fetch(`${BASE}/api/auth/magic-link`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      })

      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        error.value = data.error || 'Fehler beim Anfordern des Login-Links'
        return { success: false, message: error.value || undefined }
      }

      return {
        success: true,
        message: data.message || 'Login-Link gesendet.',
        debugLink: data.debugLink,
      }
    } catch (e: any) {
      error.value = 'Netzwerkfehler beim Anfordern des Login-Links'
      return { success: false, message: error.value }
    } finally {
      isLoading.value = false
    }
  }

  async function verifyToken(token: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      const res = await fetch(`${BASE}/api/auth/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token }),
      })

      if (!res.ok) {
        const data = await res.json().catch(() => ({}))
        error.value = data.error || 'Ungültiger oder abgelaufener Login-Link'
        return false
      }

      const authData = await res.json()
      setSession(authData.user, authData.token)
      return true
    } catch (e: any) {
      error.value = 'Netzwerkfehler bei der Verifizierung'
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMe(): Promise<User | null> {
    if (!authToken.value && typeof document !== 'undefined' && !document.cookie.includes('restgeld_session')) {
      return null
    }

    try {
      const res = await fetch(`${BASE}/api/auth/me`, {
        headers: getAuthHeaders(),
      })
      if (res.ok) {
        const me = await res.json()
        setSession(me, authToken.value)
        return me
      } else if (res.status === 401) {
        setSession(null, null)
      }
    } catch {
      // Offline fallback: use local user
    }
    return user.value
  }

  async function logout(): Promise<void> {
    try {
      await fetch(`${BASE}/api/auth/logout`, {
        method: 'POST',
        headers: getAuthHeaders(),
      })
    } catch {
      // ignore
    }
    setSession(null, null)
  }

  async function deleteAccount(): Promise<boolean> {
    try {
      const res = await fetch(`${BASE}/api/auth/me`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
      })
      if (res.ok) {
        setSession(null, null)
        return true
      }
    } catch {
      // ignore
    }
    return false
  }

  async function migrateGuestData(expenses: Expense[], periods: Period[]): Promise<number> {
    if (!isLoggedIn.value) return 0
    try {
      const res = await fetch(`${BASE}/api/auth/migrate-guest`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ expenses, periods }),
      })
      if (res.ok) {
        const data = await res.json()
        return data.migratedCount || 0
      }
    } catch {
      // ignore
    }
    return 0
  }

  async function checkUrlForAuthToken(): Promise<boolean> {
    if (typeof window === 'undefined') return false
    const params = new URLSearchParams(window.location.search)
    const token = params.get('auth_token')
    if (token) {
      const success = await verifyToken(token)
      if (success) {
        // Token aus URL entfernen
        params.delete('auth_token')
        const newSearch = params.toString() ? `?${params.toString()}` : ''
        window.history.replaceState({}, document.title, `${window.location.pathname}${newSearch}`)
      }
      return success
    }
    return false
  }

  return {
    user,
    authToken,
    isLoggedIn,
    isLoading,
    error,
    getAuthHeaders,
    requestMagicLink,
    verifyToken,
    fetchMe,
    logout,
    deleteAccount,
    migrateGuestData,
    checkUrlForAuthToken,
  }
}
