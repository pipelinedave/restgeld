<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <div class="header-title-wrap">
          <h2>{{ i18n.t('expenses.title') }}</h2>
          <span v-if="totalCount > 0" class="badge">{{ totalCount }}</span>
        </div>
        <button class="close-btn" :aria-label="i18n.t('common.close')" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div v-if="isLoading" class="loading-state">
          ...
        </div>

        <div v-else-if="expenses.length === 0" class="empty-state">
          {{ i18n.t('expenses.empty') }}
        </div>

        <ul v-else class="expense-list">
          <li v-for="exp in expenses" :key="exp.id" class="expense-item">
            <div class="expense-left">
              <span class="category-icon" :title="exp.note || i18n.t('recent.default_note')">{{ detectCategoryIcon(exp.note) }}</span>
              <div class="expense-details">
                <span class="expense-note">{{ exp.note || i18n.t('recent.default_note') }}</span>
                <span class="expense-date">{{ formatDate(exp.createdAt) }}</span>
              </div>
            </div>
            <div class="expense-right">
              <span class="expense-amount">-{{ i18n.formatMoney(exp.amount) }}</span>
              <button class="delete-btn" @click="handleDelete(exp.id)" :aria-label="i18n.t('expenses.delete')">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
                </svg>
              </button>
            </div>
          </li>
        </ul>

        <!-- Paginierung -->
        <div v-if="totalPages > 1" class="pagination-bar">
          <button
            class="pagination-btn"
            :disabled="currentPage <= 1 || isLoading"
            aria-label="Vorherige Seite"
            @click="changePage(currentPage - 1)"
          >
            &larr; {{ i18n.t('common.back') }}
          </button>
          <span class="pagination-info">{{ i18n.t('expenses.page_info', { page: currentPage, total: totalPages }) }}</span>
          <button
            class="pagination-btn"
            :disabled="currentPage >= totalPages || isLoading"
            aria-label="Nächste Seite"
            @click="changePage(currentPage + 1)"
          >
            Weiter &rarr;
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi, type Expense } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'
import { useI18n, detectCategoryIcon } from '../composables/useI18n'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'expense-deleted', id: string): void
}>()

const api = useApi()
const haptics = useHaptics()
const i18n = useI18n()
const expenses = ref<Expense[]>([])
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)
const isLoading = ref(false)
const deletingId = ref<string | null>(null)
const pageSize = 6

async function loadExpenses(page: number) {
  isLoading.value = true
  try {
    const res = await api.getExpenses(page, pageSize)
    if (Array.isArray(res)) {
      expenses.value = res
      currentPage.value = 1
      totalPages.value = 1
      totalCount.value = res.length
    } else if (res && Array.isArray(res.items)) {
      expenses.value = res.items
      currentPage.value = res.page || 1
      totalPages.value = res.totalPages || 1
      totalCount.value = res.total ?? res.items.length
    } else {
      expenses.value = []
      currentPage.value = 1
      totalPages.value = 1
      totalCount.value = 0
    }
  } catch (err: any) {
    console.error('Fehler beim Laden der Ausgaben-Historie:', err.message)
    expenses.value = []
  } finally {
    isLoading.value = false
  }
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      currentPage.value = 1
      loadExpenses(1)
    }
  },
  { immediate: true }
)

function changePage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    haptics.tap()
    loadExpenses(page)
  }
}

async function handleDelete(id: string) {
  if (deletingId.value) return
  deletingId.value = id
  haptics.tap()
  try {
    await api.deleteExpense(id)
    haptics.success()
    emit('expense-deleted', id)
    if (expenses.value.length === 1 && currentPage.value > 1) {
      await loadExpenses(currentPage.value - 1)
    } else {
      await loadExpenses(currentPage.value)
    }
  } catch (err: any) {
    haptics.error()
    console.error('Fehler beim Löschen der Ausgabe:', err.message)
  } finally {
    deletingId.value = null
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return date.toLocaleString(i18n.currentLocale.value === 'en' ? 'en-US' : 'de-DE', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 10, 12, 0.85);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal-content {
  background: var(--bg-card, #121216);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 20px;
  width: 100%;
  max-width: 440px;
  max-height: 85vh;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.8), 0 0 1px 1px rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  flex-shrink: 0;
}

.header-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.badge {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  font-size: 0.72rem;
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  padding: 2px 8px;
  border-radius: 9999px;
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.close-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  font-size: 1.3rem;
  line-height: 1;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all 0.15s ease;
}

.close-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
}

.modal-body {
  padding: 16px 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-state,
.empty-state {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.85rem;
  text-align: center;
  padding: 32px 16px;
}

.expense-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.expense-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 12px;
}

.expense-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.category-icon {
  font-size: 1.1rem;
  line-height: 1;
  flex-shrink: 0;
}

.expense-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.expense-note {
  color: var(--text-main, #f4f4f6);
  font-size: 0.9rem;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.expense-date {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.72rem;
  font-family: var(--font-mono, monospace);
}

.expense-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
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

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 10px;
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  margin-top: 4px;
}

.pagination-btn {
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  padding: 6px 12px;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.pagination-btn:hover:not(:disabled) {
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}

.pagination-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 0.75rem;
  color: var(--text-dim, #5c5c6e);
}
</style>
