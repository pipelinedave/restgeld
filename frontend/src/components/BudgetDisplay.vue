<template>
  <section class="mockup-hero-card">
    <div class="hero-top-row">
      <span class="hero-label">HEUTE VERFÜGBAR</span>
      <span class="hero-base-tag">Basis: {{ baseBudgetFormatted }} &euro;</span>
    </div>
    
    <div class="hero-amount-row" :class="['color-' + color, { 'budget-pulse': isPulsing }]">
      <span class="current-amount">{{ currentFormatted }}</span>
      <span class="fraction-slash">/</span>
      <span class="start-amount">{{ startTodayFormatted }}</span>
    </div>

    <div class="hero-meta">
      <div class="hero-badge-pill" :class="'badge-' + color">
        <span v-if="savings > 0 && currentBudget > 0">
          +{{ savingsFormatted }} &euro; Spar-Puffer
        </span>
        <span v-else-if="savings < 0">
          {{ savingsFormatted }} &euro; überzogen
        </span>
        <span v-else-if="currentBudget === 0">
          Kein Tagesbudget mehr übrig
        </span>
        <span v-else>
          Perfekt im Plan
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    currentBudget: number
    baseBudget: number
    savings: number
    color: string
    spentToday?: number
  }>(),
  {
    spentToday: 0,
  }
)

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

const currentFormatted = computed(() =>
  props.currentBudget.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' €'
)

const startTodayAmount = computed(() => {
  return props.currentBudget + (props.spentToday || 0)
})

const startTodayFormatted = computed(() =>
  startTodayAmount.value.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' €'
)

const savingsFormatted = computed(() =>
  props.savings.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
)

const baseBudgetFormatted = computed(() =>
  props.baseBudget.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
)
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
