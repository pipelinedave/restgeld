<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <div class="header-title-wrap">
          <h2>Alle Ausgaben</h2>
          <span v-if="totalCount > 0" class="badge">{{ totalCount }}</span>
        </div>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div v-if="isLoading" class="loading-state">
          Lade Ausgaben...
        </div>

        <div v-else-if="expenses.length === 0" class="empty-state">
          Keine Ausgaben in dieser Periode vorhanden.
        </div>

        <ul v-else class="expense-list">
          <li v-for="exp in expenses" :key="exp.id" class="expense-item">
            <div class="expense-details">
              <span class="expense-note">{{ exp.note || 'Ausgabe' }}</span>
              <span class="expense-date">{{ formatDate(exp.createdAt) }}</span>
            </div>
            <div class="expense-right">
              <span class="expense-amount">-{{ formatAmount(exp.amount) }} &euro;</span>
              <button class="delete-btn" @click="handleDelete(exp.id)" aria-label="Löschen">
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
            &larr; Zurück
          </button>
          <span class="pagination-info">Seite {{ currentPage }} von {{ totalPages }}</span>
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

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'expense-deleted', id: string): void
}>()

const api = useApi()
const expenses = ref<Expense[]>([])
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)
const isLoading = ref(false)
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
    loadExpenses(page)
  }
}

async function handleDelete(id: string) {
  try {
    await api.deleteExpense(id)
    emit('expense-deleted', id)
    // Wenn das letzte Element der Seite gelöscht wurde und wir auf einer höheren Seite sind
    if (expenses.value.length === 1 && currentPage.value > 1) {
      loadExpenses(currentPage.value - 1)
    } else {
      loadExpenses(currentPage.value)
    }
  } catch (err: any) {
    console.error('Fehler beim Löschen der Ausgabe:', err.message)
  }
}

function formatAmount(amount: number) {
  return amount.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return date.toLocaleString('de-DE', {
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
  background: rgba(10, 25, 47, 0.85);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal-content {
  background: var(--bg-card);
  border: 1px solid #233554;
  border-radius: 16px;
  width: 100%;
  max-width: 440px;
  max-height: 85vh;
  box-shadow: 0 10px 30px -10px rgba(2, 12, 27, 0.7);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #233554;
  flex-shrink: 0;
}

.header-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-header h2 {
  font-size: 1.25rem;
  color: var(--text);
  margin: 0;
  font-weight: 600;
}

.badge {
  background: rgba(100, 255, 218, 0.15);
  color: var(--accent);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid rgba(100, 255, 218, 0.3);
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--text-dim);
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.close-btn:hover {
  color: var(--text);
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
  color: var(--text-dim);
  font-size: 0.9rem;
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
  padding: 10px 12px;
  background: var(--bg);
  border: 1px solid #233554;
  border-radius: 8px;
}

.expense-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.expense-note {
  color: var(--text);
  font-size: 0.95rem;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.expense-date {
  color: var(--text-dim);
  font-size: 0.75rem;
}

.expense-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.expense-amount {
  color: var(--danger);
  font-weight: 600;
  font-size: 0.95rem;
  white-space: nowrap;
}

.delete-btn {
  background: none;
  border: none;
  color: var(--danger);
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

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 8px;
  border-top: 1px solid #233554;
  margin-top: 4px;
}

.pagination-btn {
  background: transparent;
  border: 1px solid #233554;
  color: var(--text);
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.15s;
}

.pagination-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
}

.pagination-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 0.8rem;
  color: var(--text-dim);
}
</style>
