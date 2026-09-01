<template>
  <section class="recent-expenses">
    <div class="section-header">
      <h3 class="section-title">{{ i18n.t('recent.title') }}</h3>
      <button v-if="expenses.length > 0" class="show-all-btn" @click="$emit('open-all')">
        <span>{{ i18n.t('recent.show_all') }}</span>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"></polyline>
        </svg>
      </button>
    </div>
    <ul v-if="expenses.length > 0" class="expense-list">
      <li v-for="exp in expenses" :key="exp.id" class="expense-item">
        <span class="category-icon" :title="exp.note || i18n.t('recent.default_note')">{{ detectCategoryIcon(exp.note) }}</span>
        <span class="expense-note">{{ exp.note || i18n.t('recent.default_note') }}</span>
        <span class="expense-amount">-{{ i18n.formatMoney(exp.amount) }}</span>
        <button class="delete-btn" :aria-label="i18n.t('recent.delete_title')" @click="$emit('delete', exp.id)">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
          </svg>
        </button>
      </li>
    </ul>
    <p v-else class="empty-state">{{ i18n.t('recent.empty') }}</p>
  </section>
</template>

<script setup lang="ts">
import { useI18n, detectCategoryIcon } from '../composables/useI18n'
import type { Expense } from '../composables/useApi'

defineProps<{
  expenses: Expense[]
}>()

defineEmits<{
  (e: 'delete', id: string): void
  (e: 'open-all'): void
}>()

const i18n = useI18n()
</script>

<style scoped>
.recent-expenses {
  padding: 8px 16px 12px;
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.section-title {
  font-size: 0.78rem;
  color: var(--text-muted, #8e8e9c);
  margin: 0;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.show-all-btn {
  background: transparent;
  border: none;
  color: var(--accent-green, #22c55e);
  font-size: 0.78rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.show-all-btn:hover {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
}

.expense-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
  padding-right: 4px;
}

.expense-list::-webkit-scrollbar {
  width: 4px;
}

.expense-list::-webkit-scrollbar-track {
  background: transparent;
}

.expense-list::-webkit-scrollbar-thumb {
  background: var(--border-color, rgba(255, 255, 255, 0.15));
  border-radius: 4px;
}

.expense-list::-webkit-scrollbar-thumb:hover {
  background: var(--text-dim, #5c5c6e);
}

.expense-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 12px;
}

.category-icon {
  font-size: 1.1rem;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.expense-note {
  flex: 1;
  color: var(--text-main, #f4f4f6);
  font-size: 0.9rem;
  font-weight: 500;
}

.expense-amount {
  color: var(--accent-red, #ef4444);
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  font-size: 0.9rem;
  white-space: nowrap;
}

.delete-btn {
  background: none;
  border: none;
  color: var(--accent-red, #ef4444);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  opacity: 0.5;
  transition: opacity 0.15s;
}

.delete-btn:hover,
.delete-btn:active {
  opacity: 1;
}

.empty-state {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.85rem;
  text-align: center;
  padding: 18px;
  background: var(--bg-card, #121216);
  border: 1px dashed var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 12px;
}
</style>
