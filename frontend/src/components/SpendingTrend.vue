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
        <span class="detail-day">Tag {{ selectedStat.day }} ({{ formatDate(selectedStat.date) }}):</span>
        <span
          class="detail-spent"
          :class="selectedStat.spent > baseBudget ? 'spent-over' : 'spent-ok'"
        >
          {{ selectedStat.spent > 0 ? formatAmount(selectedStat.spent) + ' €' : '0,00 € (Spar-Tag! 🎯)' }}
        </span>
      </template>
      <template v-else>
        <span class="detail-hint">Tippe auf einen Tag für Details</span>
      </template>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { DailyStat } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'

const props = defineProps<{
  stats?: DailyStat[]
  baseBudget: number
  currentDay: number
}>()

const haptics = useHaptics()
const selectedDay = ref<number | null>(null)

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
  } else {
    selectedDay.value = stat.day
  }
}

function formatAmount(amount: number) {
  return amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
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
  padding: 16px;
  background: var(--bg-card, #112240);
  border: 1px solid #233554;
  border-radius: 14px;
  margin: 0 16px 16px 16px;
}

.trend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.trend-title {
  font-size: 0.85rem;
  color: var(--text-dim, #8892b0);
  margin: 0;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.average-badge {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--accent, #64ffda);
  background: rgba(100, 255, 218, 0.1);
  padding: 3px 8px;
  border-radius: 12px;
  border: 1px solid rgba(100, 255, 218, 0.2);
}

.chart-container {
  width: 100%;
  overflow-x: auto;
  padding-bottom: 6px;
  -webkit-overflow-scrolling: touch;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  height: 90px;
  min-width: 100%;
  padding-top: 10px;
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
  background: var(--accent, #64ffda);
  box-shadow: 0 0 6px rgba(100, 255, 218, 0.3);
}

.bar-over {
  background: #ff6b6b;
  box-shadow: 0 0 6px rgba(255, 107, 107, 0.3);
}

.bar-zero {
  background: #495670;
  opacity: 0.5;
}

.base-budget-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: rgba(204, 214, 246, 0.25);
  pointer-events: none;
  z-index: 2;
}

.bar-label {
  font-size: 0.65rem;
  color: var(--text-dim, #8892b0);
  margin-top: 4px;
  line-height: 1;
}

.is-today .bar-label {
  color: var(--accent, #64ffda);
  font-weight: 700;
}

.is-selected .bar-track {
  outline: 1.5px solid var(--accent, #64ffda);
}

.detail-preview {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #1d2d50;
  font-size: 0.8rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 24px;
}

.detail-hint {
  color: #495670;
  font-style: italic;
  font-size: 0.75rem;
}

.detail-day {
  color: var(--text-dim, #8892b0);
}

.detail-spent {
  font-weight: 600;
}

.spent-ok {
  color: var(--accent, #64ffda);
}

.spent-over {
  color: #ff6b6b;
}
</style>
