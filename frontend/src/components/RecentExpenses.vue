<template>
  <section class="recent-expenses">
    <h3 class="section-title">Letzte Ausgaben</h3>
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
  delete: [id: string]
}>()

function formatted(amount: number) {
  return amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.recent-expenses {
  padding: 16px;
}

.section-title {
  font-size: 0.85rem;
  color: #8892b0;
  margin-bottom: 8px;
  font-weight: 500;
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
