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
      @open-settings="openSettings"
    />

    <MonthProgress
      :day="budget?.day ?? 1"
      :monthDays="budget?.monthDays ?? 30"
    />

    <MonthProjection :projection="budget?.projection" />

    <div class="hero-area">
      <BudgetDisplay
        v-if="budget"
        :currentBudget="budget.currentBudget"
        :baseBudget="budget.baseBudget"
        :savings="budget.savings"
        :color="budget.color"
      />
      <div v-else class="loading">Lade...</div>

      <button class="add-btn" @click="openNumpad">
        &minus; Ausgabe
      </button>
    </div>

    <StreakCard :streak="budget?.streak" />

    <SpendingTrend
      :stats="budget?.dailyStats"
      :baseBudget="budget?.baseBudget ?? 15"
      :currentDay="budget?.day ?? 1"
    />

    <div class="history-area">
      <RecentExpenses
        :expenses="budget?.expenses ?? []"
        @delete="handleDelete"
        @open-all="openExpensesModal"
      />
    </div>

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
      @close="showSettings = false"
    />

    <PeriodsArchiveModal
      :visible="showArchiveModal"
      @close="showArchiveModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useApi, type BudgetData } from './composables/useApi'
import { useHaptics } from './composables/useHaptics'
import { useOfflineSync } from './composables/useOfflineSync'
import AppHeader from './components/AppHeader.vue'
import MonthProgress from './components/MonthProgress.vue'
import MonthProjection from './components/MonthProjection.vue'
import BudgetDisplay from './components/BudgetDisplay.vue'
import StreakCard from './components/StreakCard.vue'
import SpendingTrend from './components/SpendingTrend.vue'
import Numpad from './components/Numpad.vue'
import RecentExpenses from './components/RecentExpenses.vue'
import ExpensesModal from './components/ExpensesModal.vue'
import SettingsModal from './components/SettingsModal.vue'
import PeriodsArchiveModal from './components/PeriodsArchiveModal.vue'
import AppFooter from './components/AppFooter.vue'
import ToastNotification from './components/ToastNotification.vue'

const api = useApi()
const haptics = useHaptics()
const offlineSync = useOfflineSync()
const budget = ref<BudgetData | null>(null)
const showNumpad = ref(false)
const showSettings = ref(false)
const showExpensesModal = ref(false)
const showArchiveModal = ref(false)
const isSavingExpense = ref(false)

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

async function loadBudget() {
  try {
    budget.value = await api.getBudget()
  } catch (e: any) {
    console.error('Fehler beim Laden:', e.message)
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
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
}

.hero-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding: 12px 16px 16px;
}

.loading {
  color: var(--text-dim, #5c5c6e);
  font-size: 1.1rem;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  max-width: 260px;
  padding: 13px 24px;
  font-size: 1rem;
  font-weight: 700;
  border-radius: 9999px;
  background-color: var(--accent-green, #22c55e);
  color: #05200e;
  border: 1px solid transparent;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(34, 197, 94, 0.25);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.add-btn:hover {
  background-color: #2ed66b;
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(34, 197, 94, 0.4);
}

.add-btn:active {
  transform: scale(0.97);
  background-color: #1eb854;
}

.history-area {
  flex-shrink: 0;
}
</style>
