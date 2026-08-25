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
  margin: 0 16px 12px 16px;
  background: var(--bg-card, #121216);
  border: 1px solid rgba(249, 115, 22, 0.18);
  border-radius: 16px;
  padding: 12px 16px;
  box-shadow: 0 4px 20px -5px rgba(0, 0, 0, 0.5);
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
  gap: 12px;
}

.flame-icon {
  font-size: 1.6rem;
  filter: grayscale(0.8) opacity(0.5);
  transition: all 0.3s ease;
  line-height: 1;
}

.flame-active .flame-icon {
  filter: grayscale(0) opacity(1);
  animation: flamePulse 2s infinite ease-in-out;
}

@keyframes flamePulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15) rotate(3deg); }
}

.streak-info {
  display: flex;
  flex-direction: column;
}

.streak-count {
  font-size: 0.95rem;
  font-weight: 700;
  color: #fed7aa;
  line-height: 1.2;
}

.streak-label {
  font-size: 0.7rem;
  color: var(--text-muted, #8e8e9c);
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
  background: #1a1a22;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 10px;
  padding: 4px 8px;
  min-width: 52px;
}

.badge-icon {
  font-size: 0.8rem;
  line-height: 1;
  margin-bottom: 2px;
}

.badge-val {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--accent-green, #22c55e);
  font-family: var(--font-mono, monospace);
  line-height: 1.1;
}

.badge-name {
  font-size: 0.6rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  font-weight: 600;
}
</style>
