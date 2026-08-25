<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>📜 Frühere Perioden</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div v-if="loading" class="loading-state">
          <span class="spinner"></span>
          <p>Lade Archiv...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <p>{{ error }}</p>
          <button class="retry-btn" @click="loadPeriods">Erneut versuchen</button>
        </div>

        <div v-else-if="periods.length === 0" class="empty-state">
          <p>Noch keine vergangenen Perioden im Archiv.</p>
        </div>

        <div v-else class="periods-list">
          <div v-for="period in periods" :key="period.id" class="period-item">
            <div class="period-item-header">
              <span class="period-title">{{ formatPeriodTitle(period.startDate) }}</span>
              <span class="period-badge" :class="period.savings >= 0 ? 'saving' : 'deficit'">
                {{ period.savings >= 0 ? '+' : '' }}{{ formatAmount(period.savings) }} €
              </span>
            </div>

            <div class="period-details">
              <div class="detail-col">
                <span class="detail-label">Budget</span>
                <span class="detail-val">{{ formatAmount(period.monthlyTotal) }} €</span>
              </div>
              <div class="detail-col">
                <span class="detail-label">Ausgaben</span>
                <span class="detail-val">{{ formatAmount(period.totalSpent) }} €</span>
              </div>
              <div class="detail-col">
                <span class="detail-label">Buchungen</span>
                <span class="detail-val">{{ period.expenseCount }}</span>
              </div>
              <div class="detail-col">
                <span class="detail-label">Dauer</span>
                <span class="detail-val">{{ period.monthDays }} Tage</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useApi, type PeriodSummary } from '../composables/useApi'
import { useHaptics } from '../composables/useHaptics'

const props = defineProps<{
  visible: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const api = useApi()
const haptics = useHaptics()
const periods = ref<PeriodSummary[]>([])
const loading = ref(false)
const error = ref('')

async function loadPeriods() {
  loading.value = true
  error.value = ''
  try {
    periods.value = await api.getPeriods()
  } catch (err: any) {
    error.value = err.message || 'Fehler beim Laden des Archivs'
    haptics.error()
  } finally {
    loading.value = false
  }
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      loadPeriods()
    }
  },
  { immediate: true }
)

function formatPeriodTitle(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('de-DE', { month: 'long', year: 'numeric' })
  } catch {
    return dateStr
  }
}

function formatAmount(val: number): string {
  return val.toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
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
  max-height: 85dvh;
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

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
  letter-spacing: -0.3px;
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
  padding: 16px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  gap: 12px;
}

.loading-state,
.empty-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  color: var(--text-dim, #5c5c6e);
  gap: 12px;
  text-align: center;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--bg-subtle, #1c1c24);
  border-top-color: var(--accent-green, #22c55e);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.retry-btn {
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  color: var(--accent-green, #22c55e);
  padding: 6px 14px;
  border-radius: 9999px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
}

.periods-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.period-item {
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.period-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.period-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  text-transform: capitalize;
}

.period-badge {
  font-size: 0.8rem;
  font-weight: 700;
  font-family: var(--font-mono, monospace);
  padding: 3px 8px;
  border-radius: 9999px;
}

.period-badge.saving {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.period-badge.deficit {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  color: var(--accent-red, #ef4444);
  border: 1px solid rgba(239, 68, 68, 0.25);
}

.period-details {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  padding-top: 10px;
}

.detail-col {
  display: flex;
  flex-direction: column;
}

.detail-label {
  font-size: 0.65rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  font-weight: 600;
}

.detail-val {
  font-size: 0.82rem;
  font-weight: 600;
  font-family: var(--font-mono, monospace);
  color: var(--text-main, #f4f4f6);
}
</style>
