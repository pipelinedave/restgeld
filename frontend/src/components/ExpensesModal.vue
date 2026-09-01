<template>
  <div
    v-if="visible"
    class="sheet-overlay"
    :class="{ 'sheet-closing': isClosing }"
    @click.self="closeSheet"
  >
    <div
      ref="sheetRef"
      class="bottom-sheet"
      :class="{ 'sheet-animating': isClosing }"
      :style="sheetStyle"
    >
      <!-- Drag Handle Bar (Material You / Google Style) -->
      <div
        class="drag-handle-zone"
        @touchstart.passive="onTouchStart"
        @touchmove="onTouchMove"
        @touchend="onTouchEnd"
        @mousedown="onMouseDown"
      >
        <div class="drag-handle-pill"></div>
      </div>

      <div class="sheet-header">
        <div class="header-title-wrap">
          <h2>{{ i18n.t('expenses.title') }}</h2>
          <span v-if="totalCount > 0" class="badge">{{ totalCount }}</span>
        </div>
        <button class="close-btn" :aria-label="i18n.t('common.close')" @click="closeSheet">&times;</button>
      </div>

      <div
        ref="listScrollRef"
        class="sheet-body"
        @scroll="handleScroll"
      >
        <div v-if="isLoading && expenses.length === 0" class="loading-state">
          Lade Ausgaben...
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

        <!-- Infinite Scroll Sentinel & Loading Indicator -->
        <div ref="sentinelRef" class="scroll-sentinel">
          <div v-if="isLoadingMore" class="loading-more">
            <div class="spinner-small"></div>
            <span>Lade weitere Ausgaben...</span>
          </div>
          <div v-else-if="hasMore && expenses.length > 0" class="load-more-hint" @click="loadNextPage">
            <span>{{ i18n.t('expenses.load_more') || 'Mehr laden' }}</span>
          </div>
          <div v-else-if="expenses.length > 0 && !hasMore" class="end-of-list">
            <span>✓ Alle {{ totalCount }} Ausgaben geladen</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
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
const isLoadingMore = ref(false)
const deletingId = ref<string | null>(null)
const pageSize = 12

const sheetRef = ref<HTMLElement | null>(null)
const listScrollRef = ref<HTMLElement | null>(null)
const sentinelRef = ref<HTMLElement | null>(null)
const isClosing = ref(false)

// Drag gesture state
const translateY = ref(0)
const isDragging = ref(false)
let startY = 0
let currentY = 0

const hasMore = computed(() => currentPage.value < totalPages.value)

const sheetStyle = computed(() => {
  if (translateY.value > 0) {
    return {
      transform: `translateY(${translateY.value}px)`,
      transition: isDragging.value ? 'none' : 'transform 0.25s cubic-bezier(0.16, 1, 0.3, 1)',
    }
  }
  return {}
})

async function loadExpenses(page: number, append = false) {
  if (append) {
    isLoadingMore.value = true
  } else {
    isLoading.value = true
  }

  try {
    const res = await api.getExpenses(page, pageSize)
    if (Array.isArray(res)) {
      expenses.value = append ? [...expenses.value, ...res] : res
      currentPage.value = 1
      totalPages.value = 1
      totalCount.value = expenses.value.length
    } else if (res && Array.isArray(res.items)) {
      expenses.value = append ? [...expenses.value, ...res.items] : res.items
      currentPage.value = res.page || page
      totalPages.value = res.totalPages || 1
      totalCount.value = res.total ?? expenses.value.length
    } else {
      if (!append) expenses.value = []
      currentPage.value = 1
      totalPages.value = 1
      totalCount.value = 0
    }
  } catch (err: any) {
    console.error('Fehler beim Laden der Ausgaben-Historie:', err.message)
    if (!append) expenses.value = []
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

async function loadNextPage() {
  if (isLoading.value || isLoadingMore.value || !hasMore.value) return
  await loadExpenses(currentPage.value + 1, true)
}

function handleScroll() {
  if (!listScrollRef.value || isLoading.value || isLoadingMore.value || !hasMore.value) return
  const el = listScrollRef.value
  // Wenn weniger als 150px bis zum Ende gescrollt sind, nächste Seite laden
  if (el.scrollHeight - el.scrollTop - el.clientHeight < 150) {
    loadNextPage()
  }
}

function closeSheet() {
  haptics.tap()
  isClosing.value = true
  translateY.value = 0
  setTimeout(() => {
    isClosing.value = false
    emit('close')
  }, 220)
}

// Touch Drag Handlers
function onTouchStart(e: TouchEvent) {
  if (e.touches.length !== 1) return
  isDragging.value = true
  startY = e.touches[0].clientY
  currentY = startY
}

function onTouchMove(e: TouchEvent) {
  if (!isDragging.value || e.touches.length !== 1) return
  currentY = e.touches[0].clientY
  const delta = currentY - startY
  if (delta > 0) {
    translateY.value = delta
    if (e.cancelable) e.preventDefault()
  } else {
    // Leichter elastischer Widerstand nach oben
    translateY.value = delta * 0.15
  }
}

function onTouchEnd() {
  if (!isDragging.value) return
  isDragging.value = false
  if (translateY.value > 90) {
    closeSheet()
  } else {
    translateY.value = 0
  }
}

// Mouse Drag Fallback
function onMouseDown(e: MouseEvent) {
  isDragging.value = true
  startY = e.clientY
  currentY = startY

  const onMouseMove = (moveEvent: MouseEvent) => {
    if (!isDragging.value) return
    currentY = moveEvent.clientY
    const delta = currentY - startY
    if (delta > 0) {
      translateY.value = delta
    }
  }

  const onMouseUp = () => {
    isDragging.value = false
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
    if (translateY.value > 90) {
      closeSheet()
    } else {
      translateY.value = 0
    }
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      isClosing.value = false
      translateY.value = 0
      currentPage.value = 1
      loadExpenses(1, false)
    }
  },
  { immediate: true }
)

async function handleDelete(id: string) {
  if (deletingId.value) return
  deletingId.value = id
  haptics.tap()
  try {
    await api.deleteExpense(id)
    haptics.success()
    emit('expense-deleted', id)
    expenses.value = expenses.value.filter((e) => e.id !== id)
    totalCount.value = Math.max(0, totalCount.value - 1)
    if (expenses.value.length === 0 && hasMore.value) {
      await loadExpenses(1, false)
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
.sheet-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 10, 12, 0.8);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
  animation: overlay-fade-in 0.25s ease-out;
}

.sheet-overlay.sheet-closing {
  animation: overlay-fade-out 0.22s ease-in forwards;
}

@keyframes overlay-fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes overlay-fade-out {
  from { opacity: 1; }
  to { opacity: 0; }
}

.bottom-sheet {
  background: var(--bg-card, #14141a);
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.12));
  border-left: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-right: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 28px 28px 0 0;
  width: 100%;
  max-width: 520px;
  height: 82dvh;
  max-height: 85dvh;
  box-shadow: 0 -10px 40px -10px rgba(0, 0, 0, 0.8), 0 0 1px 1px rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: sheet-slide-up 0.28s cubic-bezier(0.16, 1, 0.3, 1);
  will-change: transform;
}

.sheet-animating {
  animation: sheet-slide-down 0.22s ease-in forwards !important;
}

@keyframes sheet-slide-up {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

@keyframes sheet-slide-down {
  from { transform: translateY(0); }
  to { transform: translateY(100%); }
}

/* Material You Handle */
.drag-handle-zone {
  width: 100%;
  padding: 12px 0 6px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: grab;
  touch-action: none;
  flex-shrink: 0;
}

.drag-handle-zone:active {
  cursor: grabbing;
}

.drag-handle-pill {
  width: 38px;
  height: 4px;
  background: var(--text-dim, rgba(255, 255, 255, 0.25));
  border-radius: 9999px;
  transition: background 0.15s ease, transform 0.15s ease;
}

.drag-handle-zone:hover .drag-handle-pill {
  background: var(--text-muted, rgba(255, 255, 255, 0.45));
  transform: scaleX(1.1);
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 20px 14px;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  flex-shrink: 0;
}

.header-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sheet-header h2 {
  font-size: 1.2rem;
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
  background: rgba(255, 255, 255, 0.04);
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

.sheet-body {
  padding: 16px 20px 24px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
  -webkit-overflow-scrolling: touch;
}

.sheet-body::-webkit-scrollbar {
  width: 4px;
}

.sheet-body::-webkit-scrollbar-track {
  background: transparent;
}

.sheet-body::-webkit-scrollbar-thumb {
  background: var(--border-color, rgba(255, 255, 255, 0.12));
  border-radius: 4px;
}

.loading-state,
.empty-state {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.85rem;
  text-align: center;
  padding: 36px 16px;
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
  padding: 12px 14px;
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  transition: transform 0.15s ease, background 0.15s ease;
}

.expense-item:active {
  transform: scale(0.99);
  background: #23232e;
}

.expense-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  flex: 1;
}

.category-icon {
  font-size: 1.15rem;
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
  font-size: 0.92rem;
  font-weight: 600;
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
  gap: 10px;
  flex-shrink: 0;
}

.expense-amount {
  color: var(--accent-red, #ef4444);
  font-weight: 700;
  font-family: var(--font-mono, monospace);
  font-size: 0.95rem;
  white-space: nowrap;
}

.delete-btn {
  background: none;
  border: none;
  color: var(--accent-red, #ef4444);
  cursor: pointer;
  padding: 6px;
  display: flex;
  align-items: center;
  opacity: 0.5;
  border-radius: 6px;
  transition: opacity 0.15s, background 0.15s;
}

.delete-btn:hover,
.delete-btn:active {
  opacity: 1;
  background: rgba(239, 68, 68, 0.1);
}

/* Scroll Sentinel & Loading More */
.scroll-sentinel {
  padding: 12px 0 6px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 0.78rem;
  color: var(--text-dim, #5c5c6e);
}

.loading-more {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted, #8e8e9c);
}

.spinner-small {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.15);
  border-top-color: var(--accent-green, #22c55e);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.load-more-hint {
  cursor: pointer;
  padding: 6px 14px;
  border-radius: 9999px;
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--accent-green, #22c55e);
  font-weight: 600;
  transition: all 0.15s;
}

.load-more-hint:hover {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
}

.end-of-list {
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
}
</style>

