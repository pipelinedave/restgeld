import { ref } from 'vue'

const BASE = import.meta.env.PROD ? window.location.origin : ''

export function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(base64)
  const outputArray = new Uint8Array(rawData.length)
  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i)
  }
  return outputArray
}

export function isPushSupported(): boolean {
  return typeof window !== 'undefined' && 'serviceWorker' in navigator && 'PushManager' in window
}

async function registerWithBackend(subscription: PushSubscription | null, active: boolean): Promise<boolean> {
  try {
    if (!subscription) return true
    const payload = {
      endpoint: subscription.endpoint,
      keys: subscription.toJSON().keys as { p256dh: string; auth: string },
    }
    const res = await fetch(`${BASE}/api/push/subscribe`, {
      method: active ? 'POST' : 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    return res.ok
  } catch {
    return false
  }
}

export function usePushNotifications() {
  const supported = ref(isPushSupported())
  const enabled = ref(false)
  const registering = ref(false)
  const error = ref('')

  async function getPublicKey(): Promise<string | null> {
    try {
      const res = await fetch(`${BASE}/api/push/vapid`)
      if (!res.ok) return null
      const data = await res.json()
      return data.publicKey || null
    } catch {
      return null
    }
  }

  async function enable(): Promise<boolean> {
    if (!supported.value) {
      error.value = 'Notifications nicht unterstützt'
      return false
    }
    registering.value = true
    error.value = ''
    try {
      const permission = await Notification.requestPermission()
      if (permission !== 'granted') {
        error.value = 'Berechtigung verweigert'
        return false
      }

      const reg = await navigator.serviceWorker.ready
      const publicKey = await getPublicKey()
      if (!publicKey) {
        error.value = 'VAPID-Key nicht verfügbar'
        return false
      }

      let subscription = await reg.pushManager.getSubscription()
      if (!subscription) {
        subscription = await reg.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: urlBase64ToUint8Array(publicKey) as BufferSource,
        })
      }

      const ok = await registerWithBackend(subscription, true)
      if (!ok) {
        error.value = 'Abonnement konnte nicht gespeichert werden'
        return false
      }
      enabled.value = true
      return true
    } catch (e: any) {
      error.value = e?.message || 'Aktivierung fehlgeschlagen'
      return false
    } finally {
      registering.value = false
    }
  }

  async function disable(): Promise<boolean> {
    try {
      const reg = await navigator.serviceWorker.ready
      const subscription = await reg.pushManager.getSubscription()
      await registerWithBackend(subscription, false)
      if (subscription) {
        await subscription.unsubscribe()
      }
      enabled.value = false
      return true
    } catch {
      return false
    }
  }

  /** Sendet eine Test-Benachrichtigung an dieses Gerät. */
  async function sendTest(): Promise<boolean> {
    try {
      const res = await fetch(`${BASE}/api/push/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: 'restgeld 🔔',
          body: 'Test-Benachrichtigung erfolgreich!',
        }),
      })
      return res.ok
    } catch {
      return false
    }
  }

  async function restore() {
    if (!supported.value) return
    try {
      const reg = await navigator.serviceWorker.ready
      const subscription = await reg.pushManager.getSubscription()
      enabled.value = !!subscription
    } catch {
      enabled.value = false
    }
  }

  return {
    supported,
    enabled,
    registering,
    error,
    enable,
    disable,
    sendTest,
    restore,
  }
}
