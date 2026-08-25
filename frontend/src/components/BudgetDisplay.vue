<template>
  <section class="mockup-hero-card">
    <span class="hero-label">HEUTE VERFÜGBAR</span>
    <div class="hero-amount" :class="['color-' + color, { 'budget-pulse': isPulsing }]">
      {{ formatted }}
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
      <span class="hero-base">Basis: {{ baseBudgetFormatted }} &euro;/Tag</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  currentBudget: number
  baseBudget: number
  savings: number
  color: string
}>()

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

const formatted = computed(() =>
  props.currentBudget.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' €'
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
  padding: 22px 16px;
  border-radius: 20px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  text-align: center;
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.7), 0 0 1px 1px rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.hero-label {
  font-size: 0.72rem;
  letter-spacing: 1.5px;
  font-weight: 700;
  color: var(--text-muted, #8e8e9c);
  text-transform: uppercase;
}

.hero-amount {
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: 2.8rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  margin: 6px 0 10px;
  line-height: 1.1;
  letter-spacing: -1px;
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), color 0.3s ease;
}

.budget-pulse {
  animation: budget-bump 0.35s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes budget-bump {
  0% { transform: scale(1); }
  50% { transform: scale(1.08); color: var(--accent-green, #22c55e); }
  100% { transform: scale(1); }
}

.color-green { color: var(--text-main, #f4f4f6); }
.color-white { color: var(--text-main, #f4f4f6); }
.color-red { color: var(--accent-red, #ef4444); }

.hero-meta {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.hero-badge-pill {
  font-size: 0.78rem;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 9999px;
  transition: all 0.3s ease;
}

.badge-green {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.badge-white {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-muted, #8e8e9c);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
}

.badge-red {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  color: var(--accent-red, #ef4444);
  border: 1px solid rgba(239, 68, 68, 0.25);
}

.hero-base {
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
  font-weight: 500;
}
</style>
