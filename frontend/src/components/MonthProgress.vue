<template>
  <section class="month-progress-card" aria-label="Monatsfortschritt">
    <div class="progress-info">
      <span class="progress-title">Aktuelle Periode</span>
      <span class="progress-days">Tag {{ day }} von {{ monthDays }}</span>
    </div>
    <div class="progress-bar-bg">
      <div class="progress-bar-fill" :style="{ width: pct + '%' }"></div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  day: number
  monthDays: number
}>()

const pct = computed(() => Math.round((props.day / props.monthDays) * 100))
</script>

<style scoped>
.month-progress-card {
  margin: 2px 16px 4px;
  background: var(--bg-card, #121216);
  padding: 6px 12px;
  border-radius: 14px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.72rem;
  margin-bottom: 4px;
}

.progress-title {
  color: var(--text-muted, #8e8e9c);
  font-weight: 500;
}

.progress-days {
  color: var(--text-main, #f4f4f6);
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  font-size: 0.72rem;
}

.progress-bar-bg {
  width: 100%;
  height: 6px;
  background: #1c1c24;
  border-radius: 9999px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--accent-green, #22c55e), #4ade80);
  border-radius: 9999px;
  transition: width 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
</style>
