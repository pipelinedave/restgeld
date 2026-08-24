<template>
  <section class="budget-display">
    <div class="budget-header-label">HEUTE VERFÜGBAR</div>
    <div class="budget-amount" :class="'color-' + color">
      {{ formatted }}
    </div>
    <div class="savings-badge" :class="'badge-' + color">
      <span v-if="savings > 0 && currentBudget > 0">
        +{{ savingsFormatted }} &euro; Monats-Puffer
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
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentBudget: number
  baseBudget: number
  savings: number
  color: string
}>()

const formatted = computed(() =>
  props.currentBudget.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' €'
)

const savingsFormatted = computed(() =>
  props.savings.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
)
</script>

<style scoped>
.budget-display {
  text-align: center;
  padding: 20px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.budget-header-label {
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 1.5px;
  color: var(--text-dim);
  text-transform: uppercase;
}

.budget-amount {
  font-size: 3.6rem;
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.5px;
  transition: color 0.3s ease;
}

.color-green { color: #64ffda; }
.color-white { color: #ccd6f6; }
.color-red { color: #ff6b6b; }

.savings-badge {
  display: inline-flex;
  align-items: center;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  margin-top: 4px;
  transition: all 0.3s ease;
}

.badge-green {
  background: rgba(100, 255, 218, 0.1);
  color: #64ffda;
  border: 1px solid rgba(100, 255, 218, 0.2);
}

.badge-white {
  background: rgba(204, 214, 246, 0.08);
  color: #8892b0;
  border: 1px solid rgba(204, 214, 246, 0.15);
}

.badge-red {
  background: rgba(255, 107, 107, 0.1);
  color: #ff6b6b;
  border: 1px solid rgba(255, 107, 107, 0.25);
}
</style>
