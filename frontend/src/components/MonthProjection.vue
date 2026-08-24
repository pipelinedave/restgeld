<template>
  <div v-if="projection" class="projection-strip" :class="projection.status">
    <div class="projection-content">
      <span class="projection-icon">🔮</span>
      <div class="projection-text-wrap">
        <span class="projection-title">Prognose Monatsende:</span>
        <span class="projection-highlight" :class="projection.status">
          {{ projection.projectedSavings >= 0 ? '+' : '' }}{{ formatAmount(projection.projectedSavings) }} €
          <span class="projection-tag">{{ projection.status === 'saving' ? 'gespart' : 'Defizit' }}</span>
        </span>
      </div>
    </div>
    <span class="projection-rate" title="Erwartete Gesamtausgaben: ~{{ formatAmount(projection.projectedTotalSpent) }} €">
      Ø {{ formatAmount(projection.avgDailySpend) }} € / Tag
    </span>
  </div>
</template>

<script setup lang="ts">
import type { ProjectionInfo } from '../composables/useApi'

defineProps<{
  projection?: ProjectionInfo
}>()

function formatAmount(val: number): string {
  return Math.abs(val).toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}
</script>

<style scoped>
.projection-strip {
  margin: 0 16px 14px 16px;
  background: rgba(17, 34, 64, 0.6);
  border: 1px solid #233554;
  border-radius: 10px;
  padding: 8px 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  gap: 8px;
}

.projection-strip.saving {
  border-color: rgba(100, 255, 218, 0.2);
}

.projection-strip.deficit {
  border-color: rgba(255, 107, 107, 0.25);
  background: rgba(255, 107, 107, 0.04);
}

.projection-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.projection-icon {
  font-size: 1rem;
  line-height: 1;
}

.projection-text-wrap {
  display: flex;
  align-items: baseline;
  gap: 5px;
}

.projection-title {
  color: var(--text-dim, #8892b0);
  font-size: 0.75rem;
}

.projection-highlight {
  font-weight: 700;
  font-size: 0.85rem;
}

.projection-highlight.saving {
  color: var(--accent, #64ffda);
}

.projection-highlight.deficit {
  color: #ff6b6b;
}

.projection-tag {
  font-size: 0.7rem;
  font-weight: 400;
  opacity: 0.85;
}

.projection-rate {
  font-size: 0.75rem;
  color: var(--text-dim, #8892b0);
  white-space: nowrap;
}
</style>
