<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>📜 Frühere Perioden</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div v-if="loading" class="loading-state">
          <span class="spinner"></span>
          <p>Lade Archiv...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <p>{{ error }}</p>
          <button class="retry-btn" @click="loadPeriods">Erneut versuchen</button>
        </div>

        <div v-else-if="periods.length === 0" class="empty-state">
          <p>Noch keine vergangenen Perioden im Archiv.</p>
        </div>

        <div v-else class="periods-list">
          <div
            v-for="period in periods"
            :key="period.id"
            class="period-card"
            :class="{ expanded: selectedPeriodId === period.id }"
            @click="togglePeriodDetails(period)"
          >
            <!-- Card Header -->
            <div class="period-card-header">
              <div class="period-title-group">
                <span class="period-title">{{ formatPeriodTitle(period.startDate) }}</span>
                <span class="period-daterange">{{ formatDateRange(period.startDate, period.monthDays) }}</span>
              </div>
              <div class="badge-group">
                <span class="period-badge" :class="period.savings >= 0 ? 'saving' : 'deficit'">
                  {{ period.savings >= 0 ? '+' : '' }}{{ formatAmount(period.savings) }} €
                </span>
                <span class="expand-arrow" :class="{ 'is-open': selectedPeriodId === period.id }">▾</span>
              </div>
            </div>

            <!-- Summary KPI Grid -->
            <div class="period-kpi-grid">
              <div class="kpi-col">
                <span class="kpi-label">Budget</span>
                <span class="kpi-val">{{ formatAmount(period.monthlyTotal) }} €</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">Ausgaben</span>
                <span class="kpi-val">{{ formatAmount(period.totalSpent) }} €</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">&Oslash; / Tag</span>
                <span class="kpi-val">{{ formatAvgDaily(period.totalSpent, period.monthDays) }} €</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">Buchungen</span>
                <span class="kpi-val">{{ period.expenseCount }}</span>
              </div>
            </div>

            <!-- Abschlussbericht Detail-Bereich -->
            <div v-if="selectedPeriodId === period.id" class="report-section" @click.stop>
              <div class="report-header">
                <span class="report-title">📊 Abschlussbericht & Buchungen</span>
              </div>

              <!-- Lade-Status für Buchungen -->
              <div v-if="loadingExpenses" class="report-loading">
                <span class="spinner-small"></span>
                <span>Lade Buchungen...</span>
              </div>

              <!-- Keine Ausgaben in Periode -->
              <div v-else-if="periodExpenses.length === 0" class="report-empty">
                <span>Keine Buchungen in dieser Periode erfasst (0 € ausgegeben).</span>
              </div>

              <!-- Ausgaben-Tabelle -->
              <div v-else class="report-expenses-list">
                <div
                  v-for="exp in periodExpenses"
                  :key="exp.id"
                  class="report-expense-item"
                >
                  <div class="exp-left">
                    <span class="exp-date">{{ formatExpenseDate(exp.createdAt) }}</span>
                    <span class="exp-note">{{ exp.note || 'Ausgabe' }}</span>
                  </div>
                  <span class="exp-amount">-{{ formatAmount(exp.amount) }} €</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi, type PeriodSummary, type Expense } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'

const props = defineProps<{
  visible: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const api = useApi()
const haptics = useHaptics()
const periods = ref<PeriodSummary[]>([])
const loading = ref(false)
const error = ref('')

const selectedPeriodId = ref<string | null>(null)
const periodExpenses = ref<Expense[]>([])
const loadingExpenses = ref(false)

async function loadPeriods() {
  loading.value = true
  error.value = ''
  selectedPeriodId.value = null
  try {
    periods.value = await api.getPeriods()
  } catch (err: any) {
    error.value = err.message || 'Fehler beim Laden des Archivs'
    haptics.error()
  } finally {
    loading.value = false
  }
}

async function togglePeriodDetails(period: PeriodSummary) {
  haptics.tap()
  if (selectedPeriodId.value === period.id) {
    selectedPeriodId.value = null
    return
  }

  selectedPeriodId.value = period.id
  loadingExpenses.value = true
  periodExpenses.value = []

  try {
    const res = await api.getExpenses(1, 100, period.id)
    periodExpenses.value = res.items
  } catch {
    periodExpenses.value = []
  } finally {
    loadingExpenses.value = false
  }
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      loadPeriods()
    }
  },
  { immediate: true }
)

function formatPeriodTitle(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('de-DE', { month: 'long', year: 'numeric' })
  } catch {
    return dateStr
  }
}

function formatDateRange(startDateStr: string, days: number): string {
  try {
    const start = new Date(startDateStr)
    const end = new Date(start)
    end.setDate(start.getDate() + days - 1)

    const startFormatted = start.toLocaleDateString('de-DE', { day: 'numeric', month: 'short' })
    const endFormatted = end.toLocaleDateString('de-DE', { day: 'numeric', month: 'short', year: 'numeric' })
    return `${startFormatted} – ${endFormatted} (${days} Tage)`
  } catch {
    return `${days} Tage`
  }
}

function formatExpenseDate(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('de-DE', { day: '2-digit', month: '2-digit' })
  } catch {
    return ''
  }
}

function formatAvgDaily(spent: number, days: number): string {
  if (days <= 0) return '0,00'
  return (spent / days).toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

function formatAmount(val: number): string {
  return val.toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
  animation: modal-fade 0.2s ease-out;
}

@keyframes modal-fade {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: #121216;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 20px;
  width: 100%;
  max-width: 440px;
  max-height: 85dvh;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.9);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
}

.close-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  font-size: 1.25rem;
  line-height: 1;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.close-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
}

.modal-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  gap: 12px;
}

.loading-state,
.empty-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  color: var(--text-dim, #5c5c6e);
  gap: 12px;
  text-align: center;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--accent-green, #22c55e);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.spinner-small {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--accent-green, #22c55e);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.retry-btn {
  background: transparent;
  border: 1px solid var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
}

.periods-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.period-card {
  background: #18181e;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.period-card:hover {
  border-color: rgba(255, 255, 255, 0.14);
}

.period-card.expanded {
  border-color: var(--accent-green, #22c55e);
  box-shadow: 0 4px 20px -5px rgba(34, 197, 94, 0.2);
}

.period-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.period-title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.period-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  text-transform: capitalize;
}

.period-daterange {
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
  font-family: var(--font-mono, monospace);
}

.badge-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.period-badge {
  font-size: 0.8rem;
  font-weight: 700;
  font-family: var(--font-mono, monospace);
  padding: 3px 8px;
  border-radius: 6px;
}

.period-badge.saving {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.period-badge.deficit {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  color: var(--accent-red, #ef4444);
  border: 1px solid rgba(239, 68, 68, 0.25);
}

.expand-arrow {
  font-size: 0.9rem;
  color: var(--text-dim, #5c5c6e);
  transition: transform 0.2s;
}

.expand-arrow.is-open {
  transform: rotate(180deg);
  color: var(--accent-green, #22c55e);
}

.period-kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  padding-top: 8px;
}

.kpi-col {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 0.65rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.kpi-val {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-main, #f4f4f6);
  font-family: var(--font-mono, monospace);
}

.report-section {
  margin-top: 6px;
  border-top: 1px dashed rgba(255, 255, 255, 0.08);
  padding-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.report-title {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-muted, #8e8e9c);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.report-loading,
.report-empty {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--text-dim, #5c5c6e);
  padding: 8px 0;
}

.report-expenses-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 180px;
  overflow-y: auto;
  padding-right: 4px;
}

.report-expense-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 8px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 6px;
  font-size: 0.78rem;
}

.exp-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.exp-date {
  font-family: var(--font-mono, monospace);
  font-size: 0.7rem;
  color: var(--text-dim, #5c5c6e);
}

.exp-note {
  color: var(--text-main, #f4f4f6);
}

.exp-amount {
  font-family: var(--font-mono, monospace);
  font-weight: 600;
  color: var(--text-muted, #8e8e9c);
}
</style>
