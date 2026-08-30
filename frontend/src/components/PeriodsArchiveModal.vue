<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>📜 {{ i18n.t('archive.title') }}</h2>
        <button class="close-btn" :aria-label="i18n.t('common.close')" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div v-if="loading" class="loading-state">
          <span class="spinner"></span>
          <p>{{ i18n.t('archive.loading') }}</p>
        </div>

        <div v-else-if="error" class="error-state">
          <p>{{ error }}</p>
          <button class="retry-btn" @click="loadPeriods">{{ i18n.t('archive.retry_btn') }}</button>
        </div>

        <div v-else-if="periods.length === 0" class="empty-state">
          <p>{{ i18n.t('archive.empty') }}</p>
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
                  {{ period.savings >= 0 ? '+' : '' }}{{ i18n.formatMoney(period.savings) }}
                </span>
                <span class="expand-arrow" :class="{ 'is-open': selectedPeriodId === period.id }">▾</span>
              </div>
            </div>

            <!-- Summary KPI Grid -->
            <div class="period-kpi-grid">
              <div class="kpi-col">
                <span class="kpi-label">{{ i18n.t('archive.kpi_budget') }}</span>
                <span class="kpi-val">{{ i18n.formatMoney(period.monthlyTotal) }}</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">{{ i18n.t('archive.kpi_spent') }}</span>
                <span class="kpi-val">{{ i18n.formatMoney(period.totalSpent) }}</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">{{ i18n.t('archive.kpi_avg') }}</span>
                <span class="kpi-val">{{ i18n.formatMoney(calcAvgDaily(period.totalSpent, period.monthDays)) }}</span>
              </div>
              <div class="kpi-col">
                <span class="kpi-label">{{ i18n.t('archive.kpi_count') }}</span>
                <span class="kpi-val">{{ period.expenseCount }}</span>
              </div>
            </div>

            <!-- Abschlussbericht Detail-Bereich -->
            <div v-if="selectedPeriodId === period.id" class="report-section" @click.stop>
              <div class="report-header">
                <span class="report-title">📊 {{ i18n.t('archive.report_details') }}</span>
              </div>

              <!-- Lade-Status für Buchungen -->
              <div v-if="loadingExpenses" class="report-loading">
                <span class="spinner-small"></span>
                <span>...</span>
              </div>

              <!-- Keine Ausgaben in Periode -->
              <div v-else-if="periodExpenses.length === 0" class="report-empty">
                <span>{{ i18n.t('expenses.empty') }}</span>
              </div>

              <!-- Ausgaben-Tabelle -->
              <div v-else class="report-expenses-list">
                <div
                  v-for="exp in periodExpenses"
                  :key="exp.id"
                  class="report-expense-item"
                >
                  <div class="exp-left">
                    <span class="category-icon">{{ detectCategoryIcon(exp.note) }}</span>
                    <span class="exp-date">{{ formatExpenseDate(exp.createdAt) }}</span>
                    <span class="exp-note">{{ exp.note || i18n.t('recent.default_note') }}</span>
                  </div>
                  <span class="exp-amount">-{{ i18n.formatMoney(exp.amount) }}</span>
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
import { useI18n, detectCategoryIcon } from '../composables/useI18n'

const props = defineProps<{
  visible: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const api = useApi()
const haptics = useHaptics()
const i18n = useI18n()
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
    const loc = i18n.currentLocale.value === 'en' ? 'en-US' : (i18n.currentLocale.value === 'es' ? 'es-ES' : (i18n.currentLocale.value === 'fr' ? 'fr-FR' : 'de-DE'))
    return d.toLocaleDateString(loc, { month: 'long', year: 'numeric' })
  } catch {
    return dateStr
  }
}

function formatDateRange(startDateStr: string, days: number): string {
  try {
    const start = new Date(startDateStr)
    const end = new Date(start)
    end.setDate(start.getDate() + days - 1)

    const loc = i18n.currentLocale.value === 'en' ? 'en-US' : (i18n.currentLocale.value === 'es' ? 'es-ES' : (i18n.currentLocale.value === 'fr' ? 'fr-FR' : 'de-DE'))
    const startFormatted = start.toLocaleDateString(loc, { day: 'numeric', month: 'short' })
    const endFormatted = end.toLocaleDateString(loc, { day: 'numeric', month: 'short', year: 'numeric' })
    return `${startFormatted} – ${endFormatted} (${days} ${i18n.t('streak.days_unit')})`
  } catch {
    return `${days} ${i18n.t('streak.days_unit')}`
  }
}

function formatExpenseDate(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    const loc = i18n.currentLocale.value === 'en' ? 'en-US' : (i18n.currentLocale.value === 'es' ? 'es-ES' : (i18n.currentLocale.value === 'fr' ? 'fr-FR' : 'de-DE'))
    return d.toLocaleDateString(loc, { day: '2-digit', month: '2-digit' })
  } catch {
    return ''
  }
}

function calcAvgDaily(totalSpent: number, days: number): number {
  if (!days || days <= 0) return 0
  return totalSpent / days
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 10, 12, 0.85);
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
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  max-height: 85vh;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.8), 0 0 1px 1px rgba(255, 255, 255, 0.05);
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
  flex-shrink: 0;
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.close-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  font-size: 1.3rem;
  line-height: 1;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all 0.15s ease;
}

.close-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
}

.modal-body {
  padding: 16px 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-state,
.empty-state,
.error-state {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.85rem;
  text-align: center;
  padding: 32px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.retry-btn {
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.12));
  color: var(--text-main, #f4f4f6);
  padding: 6px 14px;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  margin-top: 6px;
}

.periods-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.period-card {
  background: var(--bg-subtle, #181820);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.period-card:hover {
  border-color: rgba(255, 255, 255, 0.12);
  background: #1c1c24;
}

.period-card.expanded {
  border-color: var(--accent-green, #22c55e);
  background: #191922;
}

.period-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
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
  font-size: 0.76rem;
  font-weight: 700;
  font-family: var(--font-mono, monospace);
  padding: 3px 8px;
  border-radius: 9999px;
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
  font-size: 0.85rem;
  color: var(--text-dim, #5c5c6e);
  transition: transform 0.2s ease;
  display: inline-block;
}

.expand-arrow.is-open {
  transform: rotate(180deg);
  color: var(--accent-green, #22c55e);
}

.period-kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  background: rgba(0, 0, 0, 0.25);
  border-radius: 10px;
  padding: 10px 8px;
  text-align: center;
}

.kpi-col {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kpi-label {
  font-size: 0.65rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  font-weight: 600;
}

.kpi-val {
  font-size: 0.78rem;
  color: var(--text-main, #f4f4f6);
  font-family: var(--font-mono, monospace);
  font-weight: 600;
}

/* Report Detail Section */
.report-section {
  border-top: 1px dashed var(--border-color, rgba(255, 255, 255, 0.08));
  padding-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  cursor: default;
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.report-title {
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--text-muted, #8e8e9c);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.report-loading,
.report-empty {
  font-size: 0.76rem;
  color: var(--text-dim, #5c5c6e);
  padding: 12px 0;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.report-expenses-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 180px;
  overflow-y: auto;
  padding-right: 4px;
}

.report-expense-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  font-size: 0.78rem;
}

.exp-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.category-icon {
  font-size: 0.95rem;
  line-height: 1;
}

.exp-date {
  font-size: 0.68rem;
  color: var(--text-dim, #5c5c6e);
  font-family: var(--font-mono, monospace);
  flex-shrink: 0;
}

.exp-note {
  color: var(--text-main, #f4f4f6);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.exp-amount {
  color: var(--accent-red, #ef4444);
  font-family: var(--font-mono, monospace);
  font-weight: 600;
  flex-shrink: 0;
}
</style>
