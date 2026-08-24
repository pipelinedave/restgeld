<template>
  <div class="app-shell">
    <ToastNotification
      :visible="toast.visible"
      :message="toast.message"
      :type="toast.type"
    />

    <AppHeader @open-settings="openSettings" />

    <MonthProgress
      :day="budget?.day ?? 1"
      :monthDays="budget?.monthDays ?? 30"
    />

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
      @close="showSettings = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useApi, type BudgetData } from './composables/useApi'
import { useHaptics } from './composables/useHaptics'
import AppHeader from './components/AppHeader.vue'
import MonthProgress from './components/MonthProgress.vue'
import BudgetDisplay from './components/BudgetDisplay.vue'
import Numpad from './components/Numpad.vue'
import RecentExpenses from './components/RecentExpenses.vue'
import ExpensesModal from './components/ExpensesModal.vue'
import SettingsModal from './components/SettingsModal.vue'
import AppFooter from './components/AppFooter.vue'
import ToastNotification from './components/ToastNotification.vue'

const api = useApi()
const haptics = useHaptics()
const budget = ref<BudgetData | null>(null)
const showNumpad = ref(false)
const showSettings = ref(false)
const showExpensesModal = ref(false)
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

async function loadBudget() {
  try {
    budget.value = await api.getBudget()
  } catch (e: any) {
    console.error('Fehler beim Laden:', e.message)
  }
}

async function handleConfirm(amount: number, note: string) {
  isSavingExpense.value = true
  try {
    await api.addExpense(amount, note)
    await loadBudget()
    showNumpad.value = false
    haptics.success()
    const formatted = amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    const noteText = note ? ` (${note})` : ''
    showToast(`✓ ${formatted} € gebucht${noteText}`, 'success')
  } catch (e: any) {
    haptics.error()
    showToast('Fehler beim Speichern der Ausgabe', 'error')
    console.error('Fehler beim Speichern:', e.message)
  } finally {
    isSavingExpense.value = false
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

onMounted(loadBudget)
</script>

<style scoped>
.app-shell {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.hero-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 16px;
}

.loading {
  color: #495670;
  font-size: 1.2rem;
}

.add-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 200px;
  padding: 16px;
  font-size: 1.2rem;
  font-weight: 600;
  border: 2px solid #64ffda;
  border-radius: 12px;
  background: transparent;
  color: #64ffda;
  cursor: pointer;
  transition: all 0.15s;
}

.add-btn:active {
  background: #64ffda;
  color: #0a192f;
}

.history-area {
  flex-shrink: 0;
}
</style>
