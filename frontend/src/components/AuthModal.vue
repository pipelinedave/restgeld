<template>
  <div v-if="visible" class="modal-backdrop" @click="handleBackdropClick">
    <div class="modal-card auth-modal" @click.stop>
      <div class="modal-header">
        <h2 class="modal-title">
          {{ auth.isLoggedIn.value ? 'Konto & Cloud-Sync' : 'Anmelden / Registrieren' }}
        </h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <!-- 1. Nicht eingeloggt: Magic Link anfordern -->
        <div v-if="!auth.isLoggedIn.value" class="auth-guest-view">
          <div class="auth-hero">
            <span class="auth-hero-icon">☁️</span>
            <h3 class="auth-hero-title">Passwortlos & Sicher</h3>
            <p class="auth-hero-desc">
              Gib deine E-Mail-Adresse ein. Wir senden dir einen sicheren Magic Link, mit dem du dich sofort ohne Passwort einloggen kannst.
            </p>
          </div>

          <div v-if="!magicLinkSent" class="auth-form">
            <div class="input-group">
              <label for="email-input" class="input-label">E-Mail-Adresse</label>
              <input
                id="email-input"
                v-model="email"
                type="email"
                placeholder="deine.email@beispiel.de"
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
              {{ auth.isLoading.value ? 'Wird gesendet...' : 'Magic Link anfordern ✉️' }}
            </button>
          </div>

          <!-- Magic Link gesendet Bestätigung -->
          <div v-else class="magic-link-sent-box">
            <div class="sent-icon">✉️</div>
            <h4 class="sent-title">Link ist unterwegs!</h4>
            <p class="sent-text">
              Wir haben einen Login-Link an <strong>{{ email }}</strong> gesendet. Klicke einfach auf den Link in der E-Mail.
            </p>

            <div v-if="debugLink" class="dev-debug-box">
              <span class="dev-badge">Dev / Preview Schnell-Login</span>
              <button class="dev-login-btn" @click="handleDevDirectLogin">
                ⚡ Direkt einloggen (Dev-Token)
              </button>
            </div>

            <button class="auth-secondary-btn" @click="magicLinkSent = false">
              Andere E-Mail verwenden
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
              <span class="user-status-pill">🟢 Cloud-Sync aktiv</span>
            </div>
          </div>

          <div class="account-details-list">
            <div class="detail-row">
              <span class="detail-label">Konto-ID</span>
              <span class="detail-val mono">{{ auth.user.value?.id.slice(0, 8) }}...</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Mitglied seit</span>
              <span class="detail-val">{{ formatDate(auth.user.value?.createdAt) }}</span>
            </div>
          </div>

          <!-- Gast-Daten Migration Banner -->
          <div v-if="hasGuestData" class="migrate-banner">
            <div class="migrate-info">
              <span class="migrate-icon">📥</span>
              <div class="migrate-text">
                <strong>Lokale Daten gefunden</strong>
                <span>Möchtest du deine {{ guestExpenses?.length || 0 }} lokalen Ausgaben in diesen Account übertragen?</span>
              </div>
            </div>
            <button class="migrate-btn" :disabled="isMigrating" @click="handleMigrateData">
              {{ isMigrating ? 'Wird übertragen...' : 'Jetzt synchronisieren' }}
            </button>
          </div>

          <div class="auth-actions">
            <button class="auth-logout-btn" @click="handleLogout">
              Abmelden
            </button>
            <button class="auth-delete-btn" @click="handleDeleteAccount">
              Account & alle Daten löschen (DSGVO)
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
import { useHaptics } from '../composables/useHaptics'
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
const haptics = useHaptics()

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
  return d.toLocaleDateString('de-DE', { day: '2-digit', month: '2-digit', year: 'numeric' })
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
