import { ref, computed } from 'vue'
import { useAuth } from './useAuth'
import { useHaptics } from './useHaptics'

const BASE = typeof window !== 'undefined' && import.meta.env.PROD ? window.location.origin : ''

const isLoading = ref(false)
const error = ref<string | null>(null)

export function useBilling() {
  const auth = useAuth()
  const haptics = useHaptics()

  const isPro = computed(() => {
    return auth.user.value?.plan === 'pro'
  })

  async function createCheckoutSession(): Promise<{ checkoutUrl?: string; error?: string }> {
    isLoading.value = true
    error.value = null
    haptics.tap()

    try {
      const res = await fetch(`${BASE}/api/billing/create-checkout-session`, {
        method: 'POST',
        headers: auth.getAuthHeaders(),
      })

      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        error.value = data.error || 'Fehler beim Erstellen der Checkout-Session'
        haptics.warning()
        return { error: error.value || undefined }
      }

      haptics.success()
      return { checkoutUrl: data.checkoutUrl }
    } catch {
      error.value = 'Netzwerkfehler beim Checkout'
      haptics.warning()
      return { error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  async function openCustomerPortal(): Promise<{ portalUrl?: string; error?: string }> {
    isLoading.value = true
    error.value = null
    haptics.tap()

    try {
      const res = await fetch(`${BASE}/api/billing/customer-portal`, {
        method: 'POST',
        headers: auth.getAuthHeaders(),
      })

      const data = await res.json().catch(() => ({}))
      if (!res.ok) {
        error.value = data.error || 'Fehler beim Öffnen des Kundenportals'
        haptics.warning()
        return { error: error.value || undefined }
      }

      haptics.success()
      return { portalUrl: data.portalUrl }
    } catch {
      error.value = 'Netzwerkfehler beim Kundenportal'
      haptics.warning()
      return { error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  return {
    isPro,
    isLoading,
    error,
    createCheckoutSession,
    openCustomerPortal,
  }
}
