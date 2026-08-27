<template>
  <section v-if="stats && stats.length > 0" class="spending-trend">
    <div class="trend-header">
      <h3 class="trend-title">Tages-Verlauf</h3>
      <span class="average-badge">
        &Oslash; {{ averageSpentFormatted }} &euro; / Tag
      </span>
    </div>

    <!-- Interaktiver Balken-/Sparkline-Graph -->
    <div class="chart-container">
      <div class="chart-bars">
        <div
          v-for="stat in stats"
          :key="stat.day"
          class="bar-column"
          :class="{
            'is-selected': selectedDay === stat.day,
            'is-today': stat.day === currentDay,
          }"
          @click="selectDay(stat)"
        >
          <!-- Balken-Wrapper für Höhe -->
          <div class="bar-track">
            <div
              class="bar-fill"
              :class="getBarClass(stat.spent)"
              :style="{ height: getBarHeightPercent(stat.spent) + '%' }"
            ></div>
            <!-- Referenz-Marker für Basis-Tagesbudget -->
            <div
              class="base-budget-line"
              :style="{ bottom: getBaseLinePercent() + '%' }"
              title="Basis-Tagesbudget"
            ></div>
          </div>
          <span class="bar-label">{{ stat.day }}</span>
        </div>
      </div>
    </div>

    <!-- Detail-Anzeige bei Auswahl / Tooltip -->
    <div class="detail-preview" :class="{ 'has-selection': !!selectedStat }">
      <template v-if="selectedStat">
        <div class="detail-top">
          <span class="detail-day">Tag {{ selectedStat.day }} ({{ formatDate(selectedStat.date) }}):</span>
          <span
            class="detail-spent"
            :class="selectedStat.spent > baseBudget ? 'spent-over' : 'spent-ok'"
          >
            {{ getDayDetailText(selectedStat) }}
          </span>
        </div>
        <ul v-if="dayExpenses.length > 0" class="day-expenses">
          <li v-for="exp in dayExpenses" :key="exp.id" class="day-expense">
            <span class="exp-note">{{ exp.note || 'Ohne Notiz' }}</span>
            <span class="exp-amount">−{{ formatAmount(exp.amount) }} €</span>
          </li>
        </ul>
        <span v-else-if="loadingDay" class="detail-hint">Lade Buchungen…</span>
      </template>
      <template v-else>
        <span class="detail-hint">Tippe auf einen Tag für Details</span>
      </template>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { DailyStat, Expense } from '../composables/useApi'
import { useApi } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'

const props = defineProps<{
  stats?: DailyStat[]
  baseBudget: number
  currentDay: number
}>()

const haptics = useHaptics()
const api = useApi()
const selectedDay = ref<number | null>(null)
const dayExpenses = ref<Expense[]>([])
const loadingDay = ref(false)

const selectedStat = computed(() => {
  if (!props.stats || selectedDay.value === null) return null
  return props.stats.find((s) => s.day === selectedDay.value) || null
})

// Maximaler Betrag für Skalierung
const maxSpent = computed(() => {
  if (!props.stats || props.stats.length === 0) return props.baseBudget * 1.5 || 10
  const maxInStats = Math.max(...props.stats.map((s) => s.spent), 0)
  const ceiling = Math.max(maxInStats, props.baseBudget * 1.3)
  return ceiling > 0 ? ceiling : 10
})

const averageSpent = computed(() => {
  if (!props.stats || props.stats.length === 0) return 0
  const total = props.stats.reduce((acc, curr) => acc + curr.spent, 0)
  return total / props.stats.length
})

const averageSpentFormatted = computed(() =>
  averageSpent.value.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
)

function getBarHeightPercent(spent: number): number {
  if (spent <= 0) return 4 // Minimaler sichtbarer Marker für 0 €
  const pct = (spent / maxSpent.value) * 100
  return Math.min(Math.max(pct, 8), 100)
}

function getBaseLinePercent(): number {
  const pct = (props.baseBudget / maxSpent.value) * 100
  return Math.min(Math.max(pct, 5), 95)
}

function getBarClass(spent: number): string {
  if (spent === 0) return 'bar-zero'
  if (spent <= props.baseBudget) return 'bar-good'
  return 'bar-over'
}

function selectDay(stat: DailyStat) {
  haptics.tap()
  if (selectedDay.value === stat.day) {
    selectedDay.value = null
    dayExpenses.value = []
    return
  }
  selectedDay.value = stat.day
  loadDayExpenses(stat.date)
}

async function loadDayExpenses(date: string) {
  loadingDay.value = true
  dayExpenses.value = []
  try {
    dayExpenses.value = await api.getDayExpenses(date)
  } catch {
    dayExpenses.value = []
  } finally {
    loadingDay.value = false
  }
}

function formatAmount(amount: number) {
  return amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function getDayDetailText(stat: DailyStat): string {
  if (stat.day === props.currentDay) {
    return stat.spent > 0 ? `${formatAmount(stat.spent)} € (heute aktiv)` : '0,00 € (heute aktiv)'
  }
  if (stat.spent === 0) {
    return '0,00 € (Null-Ausgaben-Tag 🎯)'
  }
  if (stat.spent <= props.baseBudget) {
    return `${formatAmount(stat.spent)} € (im Budget 👍)`
  }
  return `${formatAmount(stat.spent)} € (über Budget ⚠️)`
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    return `${parts[2]}.${parts[1]}.`
  }
  return dateStr
}
</script>

<style scoped>
.spending-trend {
  padding: 10px 16px 8px;
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 16px;
  margin: 0 16px 12px 16px;
}

.trend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.trend-title {
  font-size: 0.78rem;
  color: var(--text-muted, #8e8e9c);
  margin: 0;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.average-badge {
  font-size: 0.72rem;
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  color: var(--accent-green, #22c55e);
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  padding: 3px 8px;
  border-radius: 9999px;
  border: 1px solid rgba(34, 197, 94, 0.2);
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
  height: 48px;
  min-width: 100%;
  padding-top: 4px;
}

.bar-column {
  flex: 1;
  min-width: 14px;
  max-width: 28px;
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
  background: #2a2a36;
  opacity: 0.6;
}

.base-budget-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: rgba(255, 255, 255, 0.15);
  pointer-events: none;
  z-index: 2;
}

.bar-label {
  font-size: 0.6rem;
  color: var(--text-dim, #5c5c6e);
  font-family: var(--font-mono, monospace);
  margin-top: 2px;
  line-height: 1;
}

.is-today .bar-label {
  color: var(--accent-green, #22c55e);
  font-weight: 700;
}

.is-selected .bar-track {
  outline: 1.5px solid var(--accent-green, #22c55e);
}

.detail-preview {
  margin-top: 4px;
  padding-top: 4px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 0.75rem;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 3px;
  min-height: 18px;
}

.detail-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.day-expenses {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.day-expense {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 0;
}

.exp-note {
  color: var(--text, #e5e7eb);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 70%;
}

.exp-amount {
  color: var(--accent-red, #ef4444);
  font-family: var(--font-mono, monospace);
  font-weight: 600;
  flex-shrink: 0;
}

.detail-hint {
  color: var(--text-dim, #5c5c6e);
  font-style: italic;
  font-size: 0.72rem;
}

.detail-day {
  color: var(--text-muted, #8e8e9c);
}

.detail-spent {
  font-weight: 600;
  font-family: var(--font-mono, monospace);
}

.spent-ok {
  color: var(--accent-green, #22c55e);
}

.spent-over {
  color: var(--accent-red, #ef4444);
}
</style>
