<template>
  <section v-if="error" class="trend-month">
    <div class="error-state">
      <p>{{ error }}</p>
      <button class="retry-btn" @click="loadTrend">{{ i18n.t('archive.retry_btn') }}</button>
    </div>
  </section>

  <section v-else-if="loading" class="trend-month">
    <div class="loading-state">
      <span class="spinner"></span>
      <p>{{ i18n.t('archive.loading') }}</p>
    </div>
  </section>

  <section v-else-if="trend && trend.months.length > 0" class="trend-month">
    <div class="trend-month-header">
      <div>
        <h3 class="trend-month-title">{{ i18n.t('trendMonth.title') }}</h3>
        <p class="trend-month-subtitle">{{ i18n.t('trendMonth.subtitle') }}</p>
      </div>
      <span class="avg-badge">
        {{ i18n.t('trendMonth.avg_savings', { amount: i18n.formatMoney(trend.summary.avgSavings) }) }}
      </span>
    </div>

    <!-- Summary-Badges -->
    <div class="summary-row">
      <div class="summary-tile">
        <span class="summary-label">{{ i18n.t('trendMonth.periods', { count: trend.summary.monthCount }) }}</span>
        <span class="summary-val">🗓</span>
      </div>
      <div class="summary-tile">
        <span class="summary-label">{{ i18n.t('trendMonth.total_saved', { amount: i18n.formatMoney(trend.summary.totalSaved) }) }}</span>
        <span class="summary-val" :class="trend.summary.totalSaved >= 0 ? 'val-good' : 'val-bad'">
          {{ trend.summary.totalSaved >= 0 ? '+' : '' }}{{ i18n.formatMoney(trend.summary.totalSaved) }}
        </span>
      </div>
      <div v-if="trend.summary.bestSavingsMonth" class="summary-tile">
        <span class="summary-label">{{ i18n.t('trendMonth.best_month') }}</span>
        <span class="summary-val val-good">{{ formatMonthLabel(trend.summary.bestSavingsMonth) }}</span>
      </div>
    </div>

    <!-- Interaktives Monats-Balkendiagramm -->
    <div class="chart-container">
      <div class="chart-bars">
        <div
          v-for="stat in trend.months"
          :key="stat.month"
          class="bar-column"
          :class="{ 'is-selected': selectedMonth === stat.month, 'is-current': stat.month === currentMonth }"
          @click="selectMonth(stat)"
        >
          <div class="bar-track">
            <div
              class="bar-fill"
              :class="getBarClass(stat.savings)"
              :style="{ height: getBarHeight(stat.savings) + '%' }"
            ></div>
            <div
              class="budget-line"
              :style="{ bottom: getBudgetLinePercent(stat) + '%' }"
              :title="i18n.t('trendMonth.kpi_budget')"
            ></div>
          </div>
          <span class="bar-label">{{ formatMonthLabel(stat.month) }}</span>
        </div>
      </div>
    </div>

    <!-- Detail ausgewählter Monat -->
    <div class="detail-preview" :class="{ 'has-selection': !!selectedStat }">
      <template v-if="selectedStat">
        <div class="detail-kpis">
          <div class="detail-kpi">
            <span class="detail-label">{{ i18n.t('trendMonth.kpi_budget') }}</span>
            <span class="detail-val">{{ i18n.formatMoney(selectedStat.monthlyTotal) }}</span>
          </div>
          <div class="detail-kpi">
            <span class="detail-label">{{ i18n.t('trendMonth.kpi_spent') }}</span>
            <span class="detail-val val-bad">{{ i18n.formatMoney(selectedStat.totalSpent) }}</span>
          </div>
          <div class="detail-kpi">
            <span class="detail-label">{{ i18n.t('trendMonth.kpi_savings') }}</span>
            <span class="detail-val" :class="selectedStat.savings >= 0 ? 'val-good' : 'val-bad'">
              {{ selectedStat.savings >= 0 ? '+' : '' }}{{ i18n.formatMoney(selectedStat.savings) }}
            </span>
          </div>
          <div class="detail-kpi">
            <span class="detail-label">{{ i18n.t('trendMonth.kpi_avg') }}</span>
            <span class="detail-val">{{ i18n.formatMoney(selectedStat.avgDailySpend) }}</span>
          </div>
          <div class="detail-kpi">
            <span class="detail-label">{{ i18n.t('archive.kpi_count') }}</span>
            <span class="detail-val">{{ selectedStat.expenseCount }}</span>
          </div>
        </div>
      </template>
      <template v-else>
        <span class="detail-hint">...</span>
      </template>
    </div>
  </section>

  <section v-else class="trend-month">
    <div class="empty-state">
      <p>{{ i18n.t('trendMonth.empty') }}</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useApi, type TrendResponse, type MonthlyTrendStat } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'
import { useI18n } from '../composables/useI18n'

const props = defineProps<{
  visible: boolean
}>()

const api = useApi()
const haptics = useHaptics()
const i18n = useI18n()

const trend = ref<TrendResponse | null>(null)
const loading = ref(false)
const error = ref('')
const selectedMonth = ref<string | null>(null)

const selectedStat = computed(() => {
  if (!trend.value || !selectedMonth.value) return null
  return trend.value.months.find((m) => m.month === selectedMonth.value) || null
})

const currentMonth = computed(() => {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
})

const maxAbsSavings = computed(() => {
  if (!trend.value || trend.value.months.length === 0) return 100
  let max = 0
  for (const m of trend.value.months) {
    const abs = Math.abs(m.savings)
    if (abs > max) max = abs
  }
  return max > 0 ? max : 100
})

function getBarHeight(savings: number): number {
  const pct = (Math.abs(savings) / maxAbsSavings.value) * 100
  return Math.min(Math.max(pct, 10), 100)
}

function getBudgetLinePercent(stat: MonthlyTrendStat): number {
  if (stat.monthlyTotal <= 0) return 100
  const pct = (stat.monthlyTotal / maxAbsSavings.value) * 100
  return Math.min(Math.max(pct, 6), 96)
}

function getBarClass(savings: number): string {
  if (savings === 0) return 'bar-zero'
  return savings > 0 ? 'bar-good' : 'bar-over'
}

function selectMonth(stat: MonthlyTrendStat) {
  haptics.tap()
  if (selectedMonth.value === stat.month) {
    selectedMonth.value = null
  } else {
    selectedMonth.value = stat.month
  }
}

function formatMonthLabel(monthKey: string): string {
  try {
    const [year, month] = monthKey.split('-').map(Number)
    const d = new Date(year, month - 1, 1)
    const loc = i18n.currentLocale.value === 'en' ? 'en-US' : (i18n.currentLocale.value === 'es' ? 'es-ES' : (i18n.currentLocale.value === 'fr' ? 'fr-FR' : 'de-DE'))
    return d.toLocaleDateString(loc, { month: 'short', year: '2-digit' })
  } catch {
    return monthKey
  }
}

async function loadTrend() {
  loading.value = true
  error.value = ''
  selectedMonth.value = null
  try {
    trend.value = await api.getTrend()
  } catch (err: any) {
    error.value = err.message || 'Fehler beim Laden des Monats-Verlaufs'
    haptics.error()
  } finally {
    loading.value = false
  }
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      loadTrend()
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.trend-month {
  padding: 16px;
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.trend-month-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.trend-month-title {
  font-size: 0.82rem;
  color: var(--text-muted, #8e8e9c);
  margin: 0;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.trend-month-subtitle {
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
  margin: 2px 0 0 0;
}

.avg-badge {
  font-size: 0.7rem;
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  color: var(--accent-green, #22c55e);
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  padding: 3px 8px;
  border-radius: 9999px;
  border: 1px solid rgba(34, 197, 94, 0.2);
  white-space: nowrap;
}

.summary-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.summary-tile {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-subtle, rgba(255, 255, 255, 0.03));
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 6px 10px;
  font-size: 0.72rem;
  flex: 1;
  min-width: 0;
}

.summary-label {
  color: var(--text-dim, #5c5c6e);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.summary-val {
  font-family: var(--font-mono, monospace);
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  white-space: nowrap;
}

.val-good {
  color: var(--accent-green, #22c55e);
}

.val-bad {
  color: var(--accent-red, #ef4444);
}

.chart-container {
  width: 100%;
  overflow-x: auto;
  padding-bottom: 4px;
  -webkit-overflow-scrolling: touch;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  height: 90px;
  min-width: 100%;
  padding-top: 8px;
}

.bar-column {
  flex: 1;
  min-width: 22px;
  max-width: 40px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  position: relative;
  transition: transform 0.15s ease;
}

.bar-column:active,
.bar-column.is-selected {
  transform: translateY(-2px);
}

.bar-track {
  flex: 1;
  width: 100%;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 4px;
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar-fill {
  width: 100%;
  border-radius: 3px;
  transition: height 0.3s ease, background-color 0.2s;
}

.bar-good {
  background: var(--accent-green, #22c55e);
  box-shadow: 0 0 6px rgba(34, 197, 94, 0.3);
}

.bar-over {
  background: var(--accent-red, #ef4444);
  box-shadow: 0 0 6px rgba(239, 68, 68, 0.3);
}

.bar-zero {
  background: var(--bg-elevated, #2a2a36);
  opacity: 0.6;
}

.budget-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: rgba(255, 255, 255, 0.18);
  pointer-events: none;
  z-index: 2;
}

.bar-label {
  font-size: 0.62rem;
  color: var(--text-dim, #5c5c6e);
  font-family: var(--font-mono, monospace);
  margin-top: 4px;
  line-height: 1;
  white-space: nowrap;
}

.is-current .bar-label {
  color: var(--accent-green, #22c55e);
  font-weight: 700;
}

.is-selected .bar-track {
  outline: 1.5px solid var(--accent-green, #22c55e);
}

.detail-preview {
  margin-top: 2px;
  padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.detail-hint {
  color: var(--text-dim, #5c5c6e);
  font-style: italic;
  font-size: 0.72rem;
}

.detail-kpis {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px;
}

.detail-kpi {
  display: flex;
  flex-direction: column;
  gap: 2px;
  text-align: center;
}

.detail-label {
  font-size: 0.6rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  font-weight: 600;
}

.detail-val {
  font-size: 0.74rem;
  color: var(--text-main, #f4f4f6);
  font-family: var(--font-mono, monospace);
  font-weight: 600;
}

.loading-state,
.empty-state,
.error-state {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.82rem;
  text-align: center;
  padding: 20px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.15);
  border-top-color: var(--accent-green, #22c55e);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
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
</style>
