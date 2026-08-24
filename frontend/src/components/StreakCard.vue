<template>
  <section v-if="streak" class="streak-card">
    <div class="streak-main">
      <div class="streak-flame-wrap" :class="{ 'flame-active': streak.currentStreak > 0 }">
        <span class="flame-icon">🔥</span>
        <div class="streak-info">
          <span class="streak-count">{{ streak.currentStreak }} {{ streak.currentStreak === 1 ? 'Tag' : 'Tage' }}</span>
          <span class="streak-label">Aktuelle Spar-Streak</span>
        </div>
      </div>

      <div class="streak-badges">
        <div class="mini-badge" title="Tage mit 0 € Ausgaben">
          <span class="badge-icon">🎯</span>
          <span class="badge-val">{{ streak.noSpendDays }}</span>
          <span class="badge-name">Null-Euro</span>
        </div>

        <div class="mini-badge" title="Längste Serie im Budget">
          <span class="badge-icon">🏆</span>
          <span class="badge-val">{{ streak.longestStreak }}</span>
          <span class="badge-name">Rekord</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { StreakInfo } from '../composables/useApi'

defineProps<{
  streak?: StreakInfo
}>()
</script>

<style scoped>
.streak-card {
  margin: 0 16px 16px 16px;
  background: var(--bg-card, #112240);
  border: 1px solid #233554;
  border-radius: 14px;
  padding: 12px 16px;
}

.streak-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.streak-flame-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.flame-icon {
  font-size: 1.6rem;
  filter: grayscale(0.8) opacity(0.5);
  transition: all 0.3s ease;
  line-height: 1;
}

.flame-active .flame-icon {
  filter: grayscale(0) opacity(1);
  text-shadow: 0 0 12px rgba(255, 165, 0, 0.6);
  animation: pulse-flame 2s infinite ease-in-out;
}

@keyframes pulse-flame {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.streak-info {
  display: flex;
  flex-direction: column;
}

.streak-count {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text, #ccd6f6);
  line-height: 1.2;
}

.streak-label {
  font-size: 0.72rem;
  color: var(--text-dim, #8892b0);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.streak-badges {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mini-badge {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: var(--bg, #0a192f);
  border: 1px solid #233554;
  border-radius: 8px;
  padding: 4px 8px;
  min-width: 54px;
}

.badge-icon {
  font-size: 0.8rem;
  line-height: 1;
  margin-bottom: 2px;
}

.badge-val {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--accent, #64ffda);
  line-height: 1.1;
}

.badge-name {
  font-size: 0.6rem;
  color: var(--text-dim, #8892b0);
  text-transform: uppercase;
}
</style>
