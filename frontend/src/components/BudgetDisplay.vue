<template>
  <section class="budget-display">
    <div class="budget-amount" :class="'color-' + color">
      {{ formatted }}
    </div>
    <div v-if="savings !== 0" class="savings" :class="'color-' + color">
      {{ savings > 0 ? '+' : '' }}{{ savingsFormatted }} &euro; angespart
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
  padding: 24px 16px;
}

.budget-amount {
  font-size: 3.5rem;
  font-weight: 700;
  line-height: 1.1;
  transition: color 0.3s ease;
}

.color-green { color: #64ffda; }
.color-white { color: #ccd6f6; }
.color-red { color: #ff6b6b; }

.savings {
  font-size: 1rem;
  margin-top: 8px;
  opacity: 0.85;
}
</style>
