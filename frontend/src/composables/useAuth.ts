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
  language?: string
  currency?: string
  isActive: boolean
  plan?: 'free' | 'pro'
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

function bufferToBase64url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

function base64urlToBuffer(base64url: string): Uint8Array {
  const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/')
  const bin = atob(base64)
  const bytes = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; i++) {
    bytes[i] = bin.charCodeAt(i)
  }
  return bytes
}

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
    const lang = typeof window !== 'undefined' ? localStorage.getItem('restgeld_language') || 'de' : 'de'
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept-Language': lang,
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

  function isPasskeySupported(): boolean {
    return typeof window !== 'undefined' && !!window.PublicKeyCredential && typeof navigator.credentials !== 'undefined'
  }

  async function registerPasskey(): Promise<{ success: boolean; message?: string }> {
    if (!isPasskeySupported()) {
      return { success: false, message: 'Passkeys werden von diesem Gerät nicht unterstützt.' }
    }

    isLoading.value = true
    try {
      const optRes = await fetch(`${BASE}/api/auth/passkey/register-options`, {
        method: 'POST',
        headers: getAuthHeaders(),
      })
      if (!optRes.ok) {
        return { success: false, message: 'Fehler beim Laden der Passkey-Optionen' }
      }

      const options = await optRes.json()
      const challengeBuf = base64urlToBuffer(options.challenge)
      const userBuf = new TextEncoder().encode(options.user?.id || user.value?.id || 'user')

      const credential = (await navigator.credentials.create({
        publicKey: {
          challenge: challengeBuf.buffer as ArrayBuffer,
          rp: {
            name: options.rp?.name || 'restgeld.',
            id: typeof window !== 'undefined' ? window.location.hostname : 'restgeld.stillon.top',
          },
          user: {
            id: userBuf.buffer as ArrayBuffer,
            name: user.value?.email || 'user',
            displayName: user.value?.email || 'Restgeld User',
          },
          pubKeyCredParams: [
            { type: 'public-key', alg: -7 },
            { type: 'public-key', alg: -257 },
          ],
          timeout: 60000,
          attestation: 'none',
        },
      })) as PublicKeyCredential | null

      if (!credential) {
        return { success: false, message: 'Registrierung abgebrochen' }
      }

      const credentialId = bufferToBase64url(credential.rawId)
      const verifyRes = await fetch(`${BASE}/api/auth/passkey/register-verify`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          credentialId,
          publicKey: credentialId,
          attestationType: 'none',
        }),
      })

      if (!verifyRes.ok) {
        return { success: false, message: 'Fehler bei der Passkey-Bestätigung' }
      }

      return { success: true, message: 'Passkey erfolgreich registriert!' }
    } catch (e: any) {
      return { success: false, message: e.message || 'Passkey-Registrierung fehlgeschlagen' }
    } finally {
      isLoading.value = false
    }
  }

  async function loginWithPasskey(): Promise<boolean> {
    if (!isPasskeySupported()) return false

    isLoading.value = true
    error.value = null
    try {
      const optRes = await fetch(`${BASE}/api/auth/passkey/login-options`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      })
      if (!optRes.ok) return false

      const options = await optRes.json()
      const challengeBuf = base64urlToBuffer(options.challenge)

      const assertion = (await navigator.credentials.get({
        publicKey: {
          challenge: challengeBuf.buffer as ArrayBuffer,
          timeout: 60000,
          rpId: typeof window !== 'undefined' ? window.location.hostname : 'restgeld.stillon.top',
        },
      })) as PublicKeyCredential | null

      if (!assertion) return false

      const credentialId = bufferToBase64url(assertion.rawId)
      const verifyRes = await fetch(`${BASE}/api/auth/passkey/login-verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ credentialId }),
      })

      if (!verifyRes.ok) {
        error.value = 'Passkey-Login fehlgeschlagen'
        return false
      }

      const authData = await verifyRes.json()
      setSession(authData.user, authData.token)
      return true
    } catch (e: any) {
      error.value = e.message || 'Passkey-Login abgebrochen'
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function checkUrlForAuthToken(): Promise<boolean> {
    if (typeof window === 'undefined') return false
    const params = new URLSearchParams(window.location.search)
    const token = params.get('auth_token')
    if (token) {
      const success = await verifyToken(token)
      if (success) {
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
    isPasskeySupported,
    registerPasskey,
    loginWithPasskey,
    checkUrlForAuthToken,
  }
}
