<template>
  <div v-if="visible" class="modal-backdrop" @click="handleBackdropClick">
    <div class="modal-card auth-modal" @click.stop>
      <div class="modal-header">
        <h2 class="modal-title">
          {{ auth.isLoggedIn.value ? (i18n.currentLocale.value === 'en' ? 'Account & Cloud Sync' : 'Konto & Cloud-Sync') : i18n.t('auth.title') }}
        </h2>
        <button class="close-btn" :aria-label="i18n.t('common.close')" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <!-- 1. Nicht eingeloggt: Magic Link anfordern / Passkey Login -->
        <div v-if="!auth.isLoggedIn.value" class="auth-guest-view">
          <div class="auth-hero">
            <span class="auth-hero-icon">☁️</span>
            <h3 class="auth-hero-title">{{ i18n.currentLocale.value === 'en' ? 'Passwordless & Secure' : 'Passwortlos & Sicher' }}</h3>
            <p class="auth-hero-desc">
              {{ i18n.t('auth.subtitle') }}
            </p>
          </div>

          <div v-if="!magicLinkSent" class="auth-form">
            <!-- 1-Tap Passkey Biometric Login -->
            <button
              v-if="auth.isPasskeySupported()"
              class="auth-passkey-login-btn"
              :disabled="auth.isLoading.value"
              @click="handlePasskeyLogin"
            >
              {{ i18n.t('auth.passkey_btn') }}
            </button>

            <div v-if="auth.isPasskeySupported()" class="auth-or-divider">
              <span>{{ i18n.currentLocale.value === 'en' ? 'or with email' : 'oder per E-Mail' }}</span>
            </div>

            <div class="input-group">
              <label for="email-input" class="input-label">{{ i18n.t('auth.email_label') }}</label>
              <input
                id="email-input"
                v-model="email"
                type="email"
                :placeholder="i18n.t('auth.email_placeholder')"
                class="auth-input"
                :disabled="auth.isLoading.value"
                @keyup.enter="handleRequestMagicLink"
              />
            </div>

            <div v-if="auth.error.value" class="auth-error">
              {{ auth.error.value }}
            </div>

            <button
              class="auth-primary-btn"
              :disabled="!isValidEmail || auth.isLoading.value"
              @click="handleRequestMagicLink"
            >
              {{ auth.isLoading.value ? i18n.t('auth.sending') : i18n.t('auth.send_link') + ' ✉️' }}
            </button>
          </div>

          <!-- Magic Link gesendet Bestätigung -->
          <div v-else class="magic-link-sent-box">
            <div class="sent-icon">✉️</div>
            <h4 class="sent-title">{{ i18n.currentLocale.value === 'en' ? 'Link is on its way!' : 'Link ist unterwegs!' }}</h4>
            <p class="sent-text">
              {{ i18n.currentLocale.value === 'en' ? `We sent a login link to ${email}. Just click the link in your email.` : `Wir haben einen Login-Link an ${email} gesendet. Klicke einfach auf den Link in der E-Mail.` }}
            </p>

            <div v-if="debugLink" class="dev-debug-box">
              <span class="dev-badge">Dev / Preview Schnell-Login</span>
              <button class="dev-login-btn" @click="handleDevDirectLogin">
                ⚡ Direkt einloggen (Dev-Token)
              </button>
            </div>

            <button class="auth-secondary-btn" @click="magicLinkSent = false">
              {{ i18n.currentLocale.value === 'en' ? 'Use another email' : 'Andere E-Mail verwenden' }}
            </button>
          </div>
        </div>

        <!-- 2. Eingeloggt: Account & Sync Übersicht -->
        <div v-else class="auth-user-view">
          <div class="user-profile-box">
            <div class="user-avatar">
              {{ auth.user.value?.email.charAt(0).toUpperCase() }}
            </div>
            <div class="user-info">
              <span class="user-email">{{ auth.user.value?.email }}</span>
              <div class="user-badges">
                <span class="user-status-pill">🟢 Cloud-Sync aktiv</span>
                <span v-if="billing.isPro.value || auth.user.value?.plan === 'pro'" class="user-pro-pill">⭐ {{ i18n.t('auth.pro_badge') }}</span>
                <span v-else class="user-free-pill">{{ i18n.t('auth.free_badge') }}</span>
              </div>
            </div>
          </div>

          <!-- Pro Upgrade Banner -->
          <div v-if="!billing.isPro.value && auth.user.value?.plan !== 'pro'" class="pro-upgrade-card">
            <div class="pro-card-content">
              <span class="pro-card-title">⭐⭐⭐ Restgeld PRO ⭐⭐⭐</span>
              <p class="pro-card-desc">{{ i18n.currentLocale.value === 'en' ? 'Unlock unlimited cloud backup, multi-device sync, and custom exports.' : 'Schalte unbegrenztes Cloud-Backup, Multi-Geräte Sync und erweiterte CSV-Exporte frei.' }}</p>
              <button class="pro-checkout-btn" :disabled="billing.isLoading.value" @click="handleUpgradePro">
                {{ billing.isLoading.value ? '...' : i18n.t('auth.upgrade_pro') }}
              </button>
            </div>
          </div>
          <div v-else class="pro-active-card">
            <span>{{ i18n.currentLocale.value === 'en' ? 'Your Restgeld PRO plan is active!' : 'Dein Restgeld PRO Paket ist aktiv!' }}</span>
            <button class="pro-portal-btn" :disabled="billing.isLoading.value" @click="handleOpenPortal">
              {{ i18n.t('auth.manage_sub') }}
            </button>
          </div>

          <div class="account-details-list">
            <div class="detail-row">
              <span class="detail-label">Konto-ID</span>
              <span class="detail-val mono">{{ auth.user.value?.id.slice(0, 8) }}...</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ i18n.currentLocale.value === 'en' ? 'Member since' : 'Mitglied seit' }}</span>
              <span class="detail-val">{{ formatDate(auth.user.value?.createdAt) }}</span>
            </div>
          </div>

          <!-- Gast-Daten Migration Banner -->
          <div v-if="hasGuestData" class="migrate-banner">
            <div class="migrate-info">
              <span class="migrate-icon">📥</span>
              <div class="migrate-text">
                <strong>{{ i18n.currentLocale.value === 'en' ? 'Local Data Found' : 'Lokale Daten gefunden' }}</strong>
                <span>{{ i18n.currentLocale.value === 'en' ? `Transfer your ${guestExpenses?.length || 0} local expenses into this account?` : `Möchtest du deine ${guestExpenses?.length || 0} lokalen Ausgaben in diesen Account übertragen?` }}</span>
              </div>
            </div>
            <button class="migrate-btn" :disabled="isMigrating" @click="handleMigrateData">
              {{ isMigrating ? '...' : (i18n.currentLocale.value === 'en' ? 'Sync Now' : 'Jetzt synchronisieren') }}
            </button>
          </div>

          <div class="auth-actions">
            <div v-if="auth.isPasskeySupported()" class="passkey-section">
              <span class="passkey-label">{{ i18n.currentLocale.value === 'en' ? 'Biometrics & Fast Login' : 'Biometrie & Schnell-Login' }}</span>
              <button class="passkey-btn" :disabled="auth.isLoading.value" @click="handleRegisterPasskey">
                {{ i18n.t('auth.register_passkey_btn') }}
              </button>
            </div>

            <button class="auth-logout-btn" @click="handleLogout">
              {{ i18n.t('auth.logout') }}
            </button>
            <button class="auth-delete-btn" @click="handleDeleteAccount">
              {{ i18n.t('auth.delete_account') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from '../composables/useAuth'
import { useBilling } from '../composables/useBilling'
import { useHaptics } from '../composables/useHaptics'
import { useI18n } from '../composables/useI18n'
import type { Expense, Period } from '../composables/useApi'

const props = defineProps<{
  visible: boolean
  guestExpenses?: Expense[]
  guestPeriods?: Period[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'login-success'): void
  (e: 'logout-success'): void
  (e: 'migration-complete', count: number): void
}>()

const auth = useAuth()
const billing = useBilling()
const haptics = useHaptics()
const i18n = useI18n()

async function handleUpgradePro() {
  const res = await billing.createCheckoutSession()
  if (res.checkoutUrl) {
    window.location.href = res.checkoutUrl
  }
}

async function handleOpenPortal() {
  const res = await billing.openCustomerPortal()
  if (res.portalUrl) {
    window.location.href = res.portalUrl
  }
}

const email = ref('')
const magicLinkSent = ref(false)
const debugLink = ref<string | null>(null)
const isMigrating = ref(false)

const isValidEmail = computed(() => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())
})

const hasGuestData = computed(() => {
  return (props.guestExpenses && props.guestExpenses.length > 0) || (props.guestPeriods && props.guestPeriods.length > 0)
})

async function handleRequestMagicLink() {
  if (!isValidEmail.value || auth.isLoading.value) return
  haptics.tap()

  const res = await auth.requestMagicLink(email.value)
  if (res.success) {
    magicLinkSent.value = true
    debugLink.value = res.debugLink || null
    haptics.success()
  } else {
    haptics.warning()
  }
}

async function handlePasskeyLogin() {
  haptics.tap()
  const ok = await auth.loginWithPasskey()
  if (ok) {
    haptics.success()
    emit('login-success')
  } else {
    haptics.warning()
  }
}

async function handleDevDirectLogin() {
  if (!debugLink.value) return
  const parts = debugLink.value.split('auth_token=')
  if (parts.length >= 2) {
    haptics.tap()
    const token = parts[1]
    const ok = await auth.verifyToken(token)
    if (ok) {
      haptics.success()
      emit('login-success')
    }
  }
}

async function handleRegisterPasskey() {
  haptics.tap()
  const res = await auth.registerPasskey()
  if (res.success) {
    haptics.success()
    alert(res.message || 'Passkey erfolgreich registriert!')
  } else {
    haptics.warning()
    alert(res.message || 'Passkey-Registrierung fehlgeschlagen')
  }
}

async function handleMigrateData() {
  if (!props.guestExpenses) return
  haptics.tap()
  isMigrating.value = true
  try {
    const count = await auth.migrateGuestData(props.guestExpenses, props.guestPeriods || [])
    haptics.success()
    emit('migration-complete', count)
  } finally {
    isMigrating.value = false
  }
}

async function handleLogout() {
  haptics.tap()
  await auth.logout()
  emit('logout-success')
  emit('close')
}

async function handleDeleteAccount() {
  if (confirm('Bist du sicher, dass du deinen Account und alle zugehörigen Daten unwiderruflich löschen möchtest?')) {
    haptics.warning()
    const ok = await auth.deleteAccount()
    if (ok) {
      emit('logout-success')
      emit('close')
    }
  }
}

function handleBackdropClick() {
  emit('close')
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString(i18n.currentLocale.value === 'en' ? 'en-US' : 'de-DE', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 150;
  padding: 16px;
  animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.auth-modal {
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 24px;
  width: 100%;
  max-width: 400px;
  padding: 24px;
  box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.9), 0 0 1px 1px rgba(255, 255, 255, 0.05);
  animation: slideUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.modal-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--text-dim, #5c5c6e);
  font-size: 1.4rem;
  line-height: 1;
  cursor: pointer;
}

.auth-hero {
  text-align: center;
  margin-bottom: 20px;
}

.auth-hero-icon {
  font-size: 2rem;
  display: block;
  margin-bottom: 8px;
}

.auth-hero-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  margin: 0 0 6px;
}

.auth-hero-desc {
  font-size: 0.8rem;
  color: var(--text-muted, #8e8e9c);
  margin: 0;
  line-height: 1.4;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.auth-passkey-login-btn {
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: var(--accent-green, #22c55e);
  font-size: 0.88rem;
  font-weight: 700;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.15s ease;
}

.auth-passkey-login-btn:hover:not(:disabled) {
  background: rgba(34, 197, 94, 0.2);
  border-color: var(--accent-green, #22c55e);
  transform: translateY(-1px);
}

.auth-or-divider {
  display: flex;
  align-items: center;
  text-align: center;
  color: var(--text-dim, #5c5c6e);
  font-size: 0.72rem;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.5px;
  margin: 2px 0;
}

.auth-or-divider::before,
.auth-or-divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.auth-or-divider span {
  padding: 0 10px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-label {
  font-size: 0.75rem;
  color: var(--text-dim, #5c5c6e);
  font-weight: 600;
  text-transform: uppercase;
}

.auth-input {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 12px;
  padding: 12px 14px;
  color: var(--text-main, #f4f4f6);
  font-size: 0.9rem;
  outline: none;
  transition: all 0.15s ease;
}

.auth-input:focus {
  border-color: var(--accent-green, #22c55e);
  background: rgba(255, 255, 255, 0.06);
}

.auth-error {
  font-size: 0.78rem;
  color: var(--accent-red, #ef4444);
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.1));
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.auth-primary-btn {
  background: var(--accent-green, #22c55e);
  color: #000;
  border: none;
  border-radius: 12px;
  padding: 12px;
  font-size: 0.9rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.auth-primary-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  filter: brightness(1.1);
}

.auth-primary-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.magic-link-sent-box {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px 0;
}

.sent-icon {
  font-size: 2.2rem;
}

.sent-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  margin: 0;
}

.sent-text {
  font-size: 0.82rem;
  color: var(--text-muted, #8e8e9c);
  line-height: 1.4;
  margin: 0;
}

.dev-debug-box {
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: 10px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: center;
}

.dev-badge {
  font-size: 0.65rem;
  color: #f59e0b;
  text-transform: uppercase;
  font-weight: 700;
}

.dev-login-btn {
  background: rgba(245, 158, 11, 0.2);
  border: 1px solid #f59e0b;
  color: #f59e0b;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
}

.auth-secondary-btn {
  background: transparent;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
  border-radius: 10px;
  padding: 8px;
  font-size: 0.8rem;
  cursor: pointer;
}

/* User view */
.auth-user-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.user-profile-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 12px 14px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.15));
  color: var(--accent-green, #22c55e);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  font-weight: 800;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-email {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
}

.user-status-pill {
  font-size: 0.7rem;
  color: var(--accent-green, #22c55e);
}

.user-pro-pill {
  font-size: 0.7rem;
  color: #f59e0b;
  font-weight: 700;
  margin-left: 6px;
}

.user-free-pill {
  font-size: 0.7rem;
  color: var(--text-dim, #5c5c6e);
  margin-left: 6px;
}

.account-details-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  padding: 10px 12px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
}

.detail-label { color: var(--text-dim, #5c5c6e); }
.detail-val { color: var(--text-muted, #8e8e9c); }
.detail-val.mono { font-family: var(--font-mono, monospace); }

.migrate-banner {
  background: rgba(34, 197, 94, 0.06);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 12px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.migrate-info {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.migrate-icon {
  font-size: 1.4rem;
}

.migrate-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 0.75rem;
  color: var(--text-muted, #8e8e9c);
}

.migrate-text strong {
  color: var(--text-main, #f4f4f6);
}

.migrate-btn {
  background: var(--accent-green, #22c55e);
  color: #000;
  font-size: 0.78rem;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  padding: 8px;
  cursor: pointer;
}

.passkey-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 10px;
}

.passkey-label {
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
  font-weight: 600;
}

.passkey-btn {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.25);
  color: var(--accent-green, #22c55e);
  border-radius: 8px;
  padding: 8px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
}

.auth-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 4px;
}

.auth-logout-btn {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  border-radius: 10px;
  padding: 10px;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
}

.auth-delete-btn {
  background: transparent;
  border: none;
  color: var(--accent-red, #ef4444);
  font-size: 0.72rem;
  padding: 6px;
  cursor: pointer;
  opacity: 0.8;
}

.auth-delete-btn:hover {
  opacity: 1;
  text-decoration: underline;
}
</style>
