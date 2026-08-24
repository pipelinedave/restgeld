<template>
  <section class="recent-expenses">
    <div class="section-header">
      <h3 class="section-title">Letzte Ausgaben</h3>
      <button v-if="expenses.length > 0" class="show-all-btn" @click="$emit('open-all')">
        <span>Alle anzeigen</span>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"></polyline>
        </svg>
      </button>
    </div>
    <ul v-if="expenses.length > 0" class="expense-list">
      <li v-for="exp in expenses" :key="exp.id" class="expense-item">
        <span class="expense-note">{{ exp.note || 'Ausgabe' }}</span>
        <span class="expense-amount">-{{ formatted(exp.amount) }} &euro;</span>
        <button class="delete-btn" @click="$emit('delete', exp.id)" aria-label="Löschen">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
          </svg>
        </button>
      </li>
    </ul>
    <p v-else class="empty-state">Noch keine Ausgaben heute</p>
  </section>
</template>

<script setup lang="ts">
import type { Expense } from '../composables/useApi'

defineProps<{
  expenses: Expense[]
}>()

defineEmits<{
  (e: 'delete', id: string): void
  (e: 'open-all'): void
}>()

function formatted(amount: number) {
  return amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.recent-expenses {
  padding: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.section-title {
  font-size: 0.85rem;
  color: var(--text-dim);
  margin: 0;
  font-weight: 500;
}

.show-all-btn {
  background: transparent;
  border: none;
  color: var(--accent);
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 4px;
  border-radius: 4px;
  transition: opacity 0.15s;
}

.show-all-btn:hover {
  opacity: 0.8;
}

.expense-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.expense-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #112240;
  border-radius: 8px;
  margin-bottom: 6px;
}

.expense-note {
  flex: 1;
  color: #ccd6f6;
  font-size: 0.95rem;
}

.expense-amount {
  color: #ff6b6b;
  font-weight: 600;
  font-size: 0.95rem;
  white-space: nowrap;
}

.delete-btn {
  background: none;
  border: none;
  color: #ff6b6b;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  opacity: 0.6;
  transition: opacity 0.15s;
}

.delete-btn:hover,
.delete-btn:active {
  opacity: 1;
}

.empty-state {
  color: #495670;
  font-size: 0.85rem;
  text-align: center;
  padding: 20px;
}
</style>
