<template>
  <section class="month-progress" aria-label="Monatsfortschritt">
    <div class="progress-row">
      <div class="progress-label">Tag {{ day }} von {{ monthDays }}</div>
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="{ width: pct + '%' }"></div>
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
.month-progress {
  padding: 4px 16px 8px;
}

.progress-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.progress-label {
  font-size: 0.8rem;
  color: #8892b0;
}

.progress-bar {
  height: 4px;
  background: #233554;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #64ffda;
  border-radius: 2px;
  transition: width 0.3s ease;
}
</style>
