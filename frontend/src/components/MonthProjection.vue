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
  margin: 0 16px 12px 16px;
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 12px;
  padding: 8px 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  gap: 8px;
}

.projection-strip.saving {
  border-color: rgba(34, 197, 94, 0.2);
}

.projection-strip.deficit {
  border-color: rgba(239, 68, 68, 0.25);
  background: rgba(239, 68, 68, 0.04);
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
  gap: 6px;
}

.projection-title {
  color: var(--text-muted, #8e8e9c);
  font-size: 0.75rem;
}

.projection-highlight {
  font-weight: 700;
  font-size: 0.85rem;
  font-family: var(--font-mono, monospace);
}

.projection-highlight.saving {
  color: var(--accent-green, #22c55e);
}

.projection-highlight.deficit {
  color: var(--accent-red, #ef4444);
}

.projection-tag {
  font-size: 0.7rem;
  font-weight: 400;
  opacity: 0.85;
  font-family: var(--font-sans, sans-serif);
}

.projection-rate {
  font-size: 0.75rem;
  color: var(--text-dim, #5c5c6e);
  font-family: var(--font-mono, monospace);
  white-space: nowrap;
}
</style>
