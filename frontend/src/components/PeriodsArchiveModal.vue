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
  background: rgba(10, 25, 47, 0.85);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal-content {
  background: var(--bg-card, #112240);
  border: 1px solid #233554;
  border-radius: 16px;
  width: 100%;
  max-width: 440px;
  max-height: 85dvh;
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

.modal-header h2 {
  font-size: 1.2rem;
  color: var(--text, #ccd6f6);
  margin: 0;
  font-weight: 600;
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--text-dim, #8892b0);
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.close-btn:hover {
  color: var(--text, #ccd6f6);
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
  color: var(--text-dim, #8892b0);
  gap: 12px;
  text-align: center;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #233554;
  border-top-color: var(--accent, #64ffda);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.retry-btn {
  background: transparent;
  border: 1px solid var(--accent, #64ffda);
  color: var(--accent, #64ffda);
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
}

.periods-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.period-item {
  background: var(--bg, #0a192f);
  border: 1px solid #233554;
  border-radius: 12px;
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
  font-size: 1rem;
  font-weight: 700;
  color: var(--text, #ccd6f6);
  text-transform: capitalize;
}

.period-badge {
  font-size: 0.85rem;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 6px;
}

.period-badge.saving {
  background: rgba(100, 255, 218, 0.1);
  color: var(--accent, #64ffda);
  border: 1px solid rgba(100, 255, 218, 0.25);
}

.period-badge.deficit {
  background: rgba(255, 107, 107, 0.1);
  color: #ff6b6b;
  border: 1px solid rgba(255, 107, 107, 0.25);
}

.period-details {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  border-top: 1px solid #1b2a47;
  padding-top: 8px;
}

.detail-col {
  display: flex;
  flex-direction: column;
}

.detail-label {
  font-size: 0.65rem;
  color: var(--text-dim, #8892b0);
  text-transform: uppercase;
}

.detail-val {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text, #ccd6f6);
}
</style>
