<template>
  <div class="app-shell">
    <ToastNotification
      :visible="toast.visible"
      :message="toast.message"
      :type="toast.type"
    />

    <AppHeader
      :isOffline="!offlineSync.isOnline.value"
      :pendingSyncCount="offlineSync.pendingCount.value"
      :streak="budget?.streak"
      :projection="budget?.projection"
      :day="budget?.day"
      :monthDays="budget?.monthDays"
      :baseBudget="budget?.baseBudget"
      @open-settings="openSettings"
    />

    <!-- Top Loading Streak -->
    <div v-if="isLoading" class="loading-streak" aria-hidden="true"></div>

    <main class="dashboard-viewport">
      <!-- Obere Leiste: Monatsfortschritt -->
      <section class="meta-section">
        <MonthProgress
          :day="budget?.day ?? 1"
          :monthDays="budget?.monthDays ?? 30"
        />
      </section>

      <!-- Hero Restgeld & Action Button -->
      <section class="hero-section">
        <BudgetDisplay
          v-if="budget"
          :currentBudget="budget.currentBudget"
          :baseBudget="budget.baseBudget"
          :savings="budget.savings"
          :color="budget.color"
          :spentToday="spentToday"
        />
        <div v-else class="loading">Lade Budget...</div>

        <button class="add-btn" @click="openNumpad">
          <span class="btn-icon">&minus;</span> Ausgabe buchen
        </button>
      </section>

      <!-- Widgets: Trend-Sparkline -->
      <section class="widgets-section">
        <SpendingTrend
          :stats="budget?.dailyStats"
          :baseBudget="budget?.baseBudget ?? 15"
          :currentDay="budget?.day ?? 1"
        />
      </section>

      <!-- Ausgaben-Schnellansicht -->
      <section class="history-section">
        <RecentExpenses
          :expenses="budget?.expenses ?? []"
          @delete="handleDelete"
          @open-all="openExpensesModal"
        />
      </section>
    </main>

    <AppFooter />

    <Numpad
      :visible="showNumpad"
      :isSaving="isSavingExpense"
      :currentBudget="budget?.currentBudget"
      :savings="budget?.savings"
      :recentExpenses="budget?.expenses"
      @confirm="handleConfirm"
      @cancel="showNumpad = false"
    />

    <ExpensesModal
      :visible="showExpensesModal"
      @expense-deleted="handleExpenseDeletedFromModal"
      @close="showExpensesModal = false"
    />

    <SettingsModal
      :visible="showSettings"
      :currentMonthlyBudget="budget ? Math.round(budget.baseBudget * budget.monthDays) : undefined"
      :currentMonthDays="budget?.monthDays"
      @update-budget="handleUpdateBudget"
      @new-period="handleNewPeriod"
      @data-imported="handleDataImported"
      @open-archive="openArchiveModal"
      @open-about="openAboutModal"
      @open-auth="openAuthModal"
      @close="showSettings = false"
    />

    <AuthModal
      :visible="showAuthModal"
      :guestExpenses="budget?.expenses"
      @login-success="handleLoginSuccess"
      @logout-success="handleLogoutSuccess"
      @migration-complete="handleMigrationComplete"
      @close="showAuthModal = false"
    />

    <PeriodsArchiveModal
      :visible="showArchiveModal"
      @close="showArchiveModal = false"
    />

    <AboutModal
      :visible="showAboutModal"
      @close="showAboutModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useApi, type BudgetData } from './composables/useApi'
import { useHaptics } from './composables/useHaptics'
import { useOfflineSync } from './composables/useOfflineSync'
import { useTheme } from './composables/useTheme'
import { useAuth } from './composables/useAuth'
import AppHeader from './components/AppHeader.vue'
import MonthProgress from './components/MonthProgress.vue'
import BudgetDisplay from './components/BudgetDisplay.vue'
import SpendingTrend from './components/SpendingTrend.vue'
import Numpad from './components/Numpad.vue'
import RecentExpenses from './components/RecentExpenses.vue'
import ExpensesModal from './components/ExpensesModal.vue'
import SettingsModal from './components/SettingsModal.vue'
import PeriodsArchiveModal from './components/PeriodsArchiveModal.vue'
import AboutModal from './components/AboutModal.vue'
import AuthModal from './components/AuthModal.vue'
import AppFooter from './components/AppFooter.vue'
import ToastNotification from './components/ToastNotification.vue'

const api = useApi()
const haptics = useHaptics()
const offlineSync = useOfflineSync()
const theme = useTheme()
const auth = useAuth()
const budget = ref<BudgetData | null>(null)

const spentToday = computed(() => {
  if (!budget.value || !budget.value.dailyStats) return 0
  const todayStat = budget.value.dailyStats.find((s) => s.day === budget.value!.day)
  return todayStat?.spent ?? 0
})
const showNumpad = ref(false)
const showSettings = ref(false)
const showExpensesModal = ref(false)
const showArchiveModal = ref(false)
const showAboutModal = ref(false)
const showAuthModal = ref(false)
const isSavingExpense = ref(false)
const isLoading = ref(false)

const toast = reactive({
  visible: false,
  message: '',
  type: 'success' as 'success' | 'error' | 'info',
})

let toastTimer: any = null

function showToast(message: string, type: 'success' | 'error' | 'info' = 'success') {
  if (toastTimer) clearTimeout(toastTimer)
  toast.message = message
  toast.type = type
  toast.visible = true
  toastTimer = setTimeout(() => {
    toast.visible = false
  }, 2500)
}

function openNumpad() {
  haptics.tap()
  showNumpad.value = true
}

function openSettings() {
  haptics.tap()
  showSettings.value = true
}

function openExpensesModal() {
  haptics.tap()
  showExpensesModal.value = true
}

function openArchiveModal() {
  haptics.tap()
  showSettings.value = false
  showArchiveModal.value = true
}

function openAboutModal() {
  haptics.tap()
  showSettings.value = false
  showAboutModal.value = true
}

function openAuthModal() {
  haptics.tap()
  showSettings.value = false
  showAuthModal.value = true
}

async function handleLoginSuccess() {
  showAuthModal.value = false
  await loadBudget()
  showToast('✓ Erfolgreich angemeldet', 'success')
}

async function handleLogoutSuccess() {
  showAuthModal.value = false
  await loadBudget()
  showToast('✓ Erfolgreich abgemeldet', 'info')
}

async function handleMigrationComplete(count: number) {
  showAuthModal.value = false
  await loadBudget()
  showToast(`✓ ${count} Ausgaben erfolgreich synchronisiert`, 'success')
}

async function loadBudget() {
  isLoading.value = true
  try {
    budget.value = await api.getBudget()
  } catch (e: any) {
    console.error('Fehler beim Laden:', e.message)
  } finally {
    isLoading.value = false
  }
}

async function handleConfirm(amount: number, note: string) {
  isSavingExpense.value = true
  const formatted = amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  const noteText = note ? ` (${note})` : ''

  // Fallback / Offline Handling
  if (!offlineSync.isOnline.value) {
    offlineSync.enqueueExpense(amount, note)
    // Optimistic Update
    if (budget.value) {
      budget.value.currentBudget -= amount
      budget.value.expenses.unshift({
        id: `offline-${Date.now()}`,
        periodId: budget.value.periodId,
        amount,
        note,
        createdAt: new Date().toISOString(),
      })
    }
    showNumpad.value = false
    isSavingExpense.value = false
    haptics.success()
    showToast(`Offline gespeichert: ${formatted} €${noteText}`, 'info')
    return
  }

  try {
    await api.addExpense(amount, note)
    await loadBudget()
    showNumpad.value = false
    haptics.success()
    showToast(`✓ ${formatted} € gebucht${noteText}`, 'success')
  } catch (e: any) {
    // If request failed (e.g. lost network during request)
    offlineSync.enqueueExpense(amount, note)
    if (budget.value) {
      budget.value.currentBudget -= amount
      budget.value.expenses.unshift({
        id: `offline-${Date.now()}`,
        periodId: budget.value.periodId,
        amount,
        note,
        createdAt: new Date().toISOString(),
      })
    }
    showNumpad.value = false
    haptics.warning()
    showToast(`Offline gespeichert: ${formatted} €${noteText}`, 'info')
  } finally {
    isSavingExpense.value = false
  }
}

async function handleAutoSync() {
  const synced = await offlineSync.syncPendingExpenses(api)
  if (synced > 0) {
    await loadBudget()
    haptics.success()
    showToast(`✓ ${synced} Offline-Ausgabe(n) synchronisiert`, 'success')
  }
}

async function handleDelete(id: string) {
  haptics.tap()
  try {
    await api.deleteExpense(id)
    await loadBudget()
    haptics.success()
    showToast('✓ Ausgabe gelöscht', 'info')
  } catch (e: any) {
    haptics.error()
    showToast('Fehler beim Löschen', 'error')
    console.error('Fehler beim Löschen:', e.message)
  }
}

async function handleExpenseDeletedFromModal() {
  await loadBudget()
  showToast('✓ Ausgabe gelöscht', 'info')
}

async function handleUpdateBudget(monthlyTotal: number, days?: number) {
  try {
    await api.updateBudget(monthlyTotal, days)
    await loadBudget()
    haptics.success()
    showToast(`✓ Einstellungen angepasst`, 'success')
  } catch (e: any) {
    haptics.error()
    showToast('Fehler beim Aktualisieren der Einstellungen', 'error')
    console.error('Fehler beim Aktualisieren:', e.message)
  }
}

async function handleNewPeriod(monthlyTotal?: number, days?: number) {
  showSettings.value = false
  try {
    await api.newPeriod(monthlyTotal, undefined, days)
    await loadBudget()
    haptics.success()
    showToast('✓ Neue Periode ab heute gestartet', 'success')
  } catch (e: any) {
    haptics.error()
    showToast('Fehler beim Starten der neuen Periode', 'error')
    console.error('Fehler beim Starten der neuen Periode:', e.message)
  }
}

async function handleDataImported(count: number) {
  await loadBudget()
  haptics.success()
  showToast(`✓ ${count} Ausgabe(n) erfolgreich importiert`, 'success')
}

let cleanupListeners: (() => void) | undefined

onMounted(async () => {
  theme.initTheme()

  // Magic Link Token in URL abfangen
  const autoLoggedIn = await auth.checkUrlForAuthToken()
  if (autoLoggedIn) {
    showToast('✓ Erfolgreich eingeloggt', 'success')
  } else {
    await auth.fetchMe()
  }

  await loadBudget()
  cleanupListeners = offlineSync.initListeners(handleAutoSync)

  // Try auto sync on load if pending
  if (offlineSync.pendingCount.value > 0 && offlineSync.isOnline.value) {
    handleAutoSync()
  }

  // PWA Shortcut Check: ?action=add-expense
  if (typeof window !== 'undefined') {
    const urlParams = new URLSearchParams(window.location.search)
    if (urlParams.get('action') === 'add-expense' || window.location.hash === '#add-expense') {
      showNumpad.value = true
      window.history.replaceState({}, document.title, window.location.pathname)
    }
  }
})

onUnmounted(() => {
  cleanupListeners?.()
})
</script>

<style scoped>
.app-shell {
  height: 100dvh;
  max-height: 100dvh;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  z-index: 1;
}

/* Top Loading Streak */
.loading-streak {
  height: 2px;
  width: 100%;
  background: linear-gradient(90deg, transparent, var(--accent-green, #22c55e), transparent);
  background-size: 200% 100%;
  animation: loading-pulse 1s linear infinite;
  flex-shrink: 0;
}

@keyframes loading-pulse {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.dashboard-viewport {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 2px 0;
  gap: 4px;
}

.meta-section {
  flex-shrink: 0;
}

.hero-section {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 4px 16px;
}

.loading {
  color: var(--text-dim, #5c5c6e);
  font-size: 1rem;
  font-family: var(--font-mono, monospace);
}

.add-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  max-width: 240px;
  padding: 10px 20px;
  font-size: 0.95rem;
  font-weight: 700;
  border-radius: 9999px;
  background-color: var(--accent-green, #22c55e);
  color: #05200e;
  border: 1px solid transparent;
  cursor: pointer;
  box-shadow: 0 4px 18px rgba(34, 197, 94, 0.25);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.add-btn:hover {
  background-color: #2ed66b;
  transform: translateY(-1px);
  box-shadow: 0 6px 22px rgba(34, 197, 94, 0.4);
}

.add-btn:active {
  transform: scale(0.97);
  background-color: #1eb854;
}

.btn-icon {
  font-size: 1.1rem;
  line-height: 1;
  font-weight: 800;
}

.widgets-section {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-section {
  flex-shrink: 0;
}
</style>
