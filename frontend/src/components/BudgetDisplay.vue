<template>
  <section class="mockup-hero-card">
    <div class="hero-top-row">
      <span class="hero-label">{{ i18n.t('budget.available_today') }}</span>
      <span class="hero-base-tag">{{ i18n.t('budget.base_label', { amount: i18n.formatMoney(baseBudget) }) }}</span>
    </div>
    
    <div class="hero-amount-row" :class="['color-' + color, { 'budget-pulse': isPulsing }]">
      <span class="current-amount">{{ i18n.formatMoney(currentBudget) }}</span>
      <span class="fraction-slash">/</span>
      <span class="start-amount">{{ i18n.formatMoney(startTodayAmount) }}</span>
    </div>

    <!-- Today Spending Gauge Bar -->
    <div
      class="today-gauge-track"
      :title="`${i18n.formatMoney(currentBudget)} / ${i18n.formatMoney(startTodayAmount)}`"
    >
      <div
        class="today-gauge-fill"
        :class="'fill-' + color"
        :style="{ width: `${todayProgressPct}%` }"
      ></div>
    </div>

    <div class="hero-meta">
      <div class="hero-badge-pill" :class="'badge-' + color">
        <span v-if="savings > 0 && currentBudget > 0">
          {{ i18n.t('budget.puffer_plus', { amount: i18n.formatMoney(savings) }) }}
        </span>
        <span v-else-if="savings < 0">
          {{ i18n.t('budget.overdrawn', { amount: i18n.formatMoney(Math.abs(savings)) }) }}
        </span>
        <span v-else-if="currentBudget === 0">
          {{ i18n.t('budget.empty_today') }}
        </span>
        <span v-else>
          {{ i18n.t('budget.on_track') }}
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from '../composables/useI18n'

const props = withDefaults(
  defineProps<{
    currentBudget: number
    baseBudget: number
    savings: number
    color: string
    todayBase?: number
    spentToday?: number
  }>(),
  {
    spentToday: 0,
  }
)

const i18n = useI18n()
const isPulsing = ref(false)

watch(
  () => props.currentBudget,
  (newVal, oldVal) => {
    if (oldVal !== undefined && newVal !== oldVal) {
      isPulsing.value = true
      setTimeout(() => {
        isPulsing.value = false
      }, 400)
    }
  }
)

const startTodayAmount = computed(() => {
  if (props.todayBase !== undefined && props.todayBase > 0) {
    return props.todayBase
  }
  return props.currentBudget + (props.spentToday || 0)
})

const todayProgressPct = computed(() => {
  if (startTodayAmount.value <= 0) return 0
  const pct = (props.currentBudget / startTodayAmount.value) * 100
  return Math.min(100, Math.max(0, Math.round(pct)))
})
</script>

<style scoped>
.mockup-hero-card {
  width: 100%;
  max-width: 380px;
  background: radial-gradient(circle at top, #1c1c24 0%, #121216 100%);
  padding: 18px 16px;
  border-radius: 20px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  text-align: center;
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.7), 0 0 1px 1px rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.hero-top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 0 4px;
}

.hero-label {
  font-size: 0.72rem;
  letter-spacing: 1.5px;
  font-weight: 700;
  color: var(--text-muted, #8e8e9c);
  text-transform: uppercase;
}

.hero-base-tag {
  font-size: 0.72rem;
  font-family: var(--font-mono, monospace);
  color: var(--text-dim, #5c5c6e);
  font-weight: 600;
}

.hero-amount-row {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 6px;
  font-family: var(--font-mono, monospace);
  transition: transform 0.15s ease, filter 0.15s ease;
  line-height: 1;
}

.current-amount {
  font-size: 2.35rem;
  font-weight: 800;
  letter-spacing: -1px;
}

.fraction-slash {
  font-size: 1.5rem;
  color: var(--text-dim, #5c5c6e);
  font-weight: 400;
}

.start-amount {
  font-size: 1.4rem;
  color: var(--text-muted, #8e8e9c);
  font-weight: 600;
}

.budget-pulse {
  animation: pulse-glow 0.4s ease-out;
}

@keyframes pulse-glow {
  0% { transform: scale(1); filter: brightness(1); }
  50% { transform: scale(1.04); filter: brightness(1.3); }
  100% { transform: scale(1); filter: brightness(1); }
}

.hero-meta {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.today-gauge-track {
  width: 100%;
  max-width: 200px;
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 9999px;
  overflow: hidden;
  position: relative;
  margin: -2px 0 2px;
}

.today-gauge-fill {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.3s cubic-bezier(0.16, 1, 0.3, 1), background-color 0.3s ease;
}

.fill-green {
  background: var(--accent-green, #22c55e);
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.5);
}

.fill-amber {
  background: #f59e0b;
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.5);
}

.fill-red {
  background: var(--accent-red, #ef4444);
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.hero-badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 9999px;
  font-size: 0.76rem;
  font-weight: 700;
  font-family: var(--font-mono, monospace);
  border: 1px solid transparent;
}

.color-green .current-amount { color: var(--accent-green, #22c55e); text-shadow: 0 0 20px rgba(34, 197, 94, 0.3); }
.color-amber .current-amount { color: #f59e0b; text-shadow: 0 0 20px rgba(245, 158, 11, 0.3); }
.color-red .current-amount { color: var(--accent-red, #ef4444); text-shadow: 0 0 20px rgba(239, 68, 68, 0.3); }

.badge-green {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border-color: rgba(34, 197, 94, 0.25);
}

.badge-amber {
  background: rgba(245, 158, 11, 0.12);
  color: #f59e0b;
  border-color: rgba(245, 158, 11, 0.25);
}

.badge-red {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  color: var(--accent-red, #ef4444);
  border-color: rgba(239, 68, 68, 0.25);
}
</style>
