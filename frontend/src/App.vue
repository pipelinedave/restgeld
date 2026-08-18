<template>
  <div class="app-shell">
    <MonthProgress :day="budget?.day ?? 1" :monthDays="budget?.monthDays ?? 30" />

    <div class="hero-area">
      <BudgetDisplay
        v-if="budget"
        :currentBudget="budget.currentBudget"
        :baseBudget="budget.baseBudget"
        :savings="budget.savings"
        :color="budget.color"
      />
      <div v-else class="loading">Lade...</div>

      <button class="add-btn" @click="showNumpad = true">
        &minus; Ausgabe
      </button>
    </div>

    <div class="history-area">
      <RecentExpenses
        :expenses="budget?.expenses ?? []"
        @delete="handleDelete"
      />
    </div>

    <Numpad
      :visible="showNumpad"
      @confirm="handleConfirm"
      @cancel="showNumpad = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi, type BudgetData } from './composables/useApi'
import MonthProgress from './components/MonthProgress.vue'
import BudgetDisplay from './components/BudgetDisplay.vue'
import Numpad from './components/Numpad.vue'
import RecentExpenses from './components/RecentExpenses.vue'

const api = useApi()
const budget = ref<BudgetData | null>(null)
const showNumpad = ref(false)

async function loadBudget() {
  try {
    budget.value = await api.getBudget()
  } catch (e: any) {
    console.error('Fehler beim Laden:', e.message)
  }
}

async function handleConfirm(amount: number, note: string) {
  showNumpad.value = false
  try {
    await api.addExpense(amount, note)
    await loadBudget()
  } catch (e: any) {
    console.error('Fehler beim Speichern:', e.message)
  }
}

async function handleDelete(id: string) {
  try {
    await api.deleteExpense(id)
    await loadBudget()
  } catch (e: any) {
    console.error('Fehler beim Löschen:', e.message)
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
