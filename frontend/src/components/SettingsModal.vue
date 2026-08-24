<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Einstellungen</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <section class="setting-section">
          <label for="monthly-budget-input" class="section-title">Monatsbudget (€)</label>
          <div class="input-row">
            <input
              id="monthly-budget-input"
              v-model.number="budgetInput"
              type="number"
              min="1"
              step="10"
              placeholder="z. B. 600"
              @keyup.enter="handleSaveBudget"
            />
            <button class="action-btn" :disabled="!isValidBudget" @click="handleSaveBudget">
              Speichern
            </button>
          </div>
          <p v-if="budgetSavedMsg" class="success-msg">{{ budgetSavedMsg }}</p>
        </section>

        <section class="setting-section danger-zone">
          <span class="section-title danger-title">Periode zurücksetzen</span>
          <p class="description">
            Startet eine neue Periode ab Tag 1 und archiviert/löscht die aktuellen Ausgaben.
          </p>

          <div v-if="!confirmReset">
            <button class="danger-btn" @click="confirmReset = true">
              Neue Periode starten
            </button>
          </div>
          <div v-else class="confirm-box">
            <p class="confirm-text">Wirklich zurücksetzen?</p>
            <div class="confirm-actions">
              <button class="danger-btn confirm" @click="handleResetPeriod">
                Ja, zurücksetzen
              </button>
              <button class="cancel-btn" @click="confirmReset = false">
                Abbrechen
              </button>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  visible: boolean
  currentMonthlyBudget?: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'update-budget', monthlyTotal: number): void
  (e: 'new-period'): void
}>()

const budgetInput = ref<number | null>(props.currentMonthlyBudget ?? null)
const confirmReset = ref(false)
const budgetSavedMsg = ref('')

watch(
  () => props.currentMonthlyBudget,
  (newVal) => {
    if (newVal) {
      budgetInput.value = newVal
    }
  },
  { immediate: true }
)

watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      confirmReset.value = false
      budgetSavedMsg.value = ''
      if (props.currentMonthlyBudget) {
        budgetInput.value = props.currentMonthlyBudget
      }
    }
  }
)

const isValidBudget = computed(() => {
  return typeof budgetInput.value === 'number' && budgetInput.value > 0
})

function handleSaveBudget() {
  if (isValidBudget.value && budgetInput.value) {
    emit('update-budget', budgetInput.value)
    budgetSavedMsg.value = 'Budget erfolgreich gespeichert!'
    setTimeout(() => {
      budgetSavedMsg.value = ''
    }, 2500)
  }
}

function handleResetPeriod() {
  confirmReset.value = false
  emit('new-period')
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
  max-width: 400px;
  box-shadow: 0 10px 30px -10px rgba(2, 12, 27, 0.7);
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #233554;
}

.modal-header h2 {
  font-size: 1.25rem;
  color: var(--text);
  margin: 0;
  font-weight: 600;
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
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.setting-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text);
}

.input-row {
  display: flex;
  gap: 8px;
}

input[type='number'] {
  flex: 1;
  background: var(--bg);
  border: 1px solid #233554;
  border-radius: 8px;
  padding: 10px 12px;
  color: var(--text);
  font-size: 1rem;
  outline: none;
}

input[type='number']:focus {
  border-color: var(--accent);
}

.action-btn {
  background: transparent;
  border: 1px solid var(--accent);
  color: var(--accent);
  padding: 10px 16px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.action-btn:not(:disabled):active {
  background: var(--accent);
  color: var(--bg);
}

.success-msg {
  font-size: 0.85rem;
  color: var(--accent);
  margin-top: 4px;
}

.danger-zone {
  border-top: 1px solid #233554;
  padding-top: 16px;
}

.danger-title {
  color: var(--danger);
}

.description {
  font-size: 0.85rem;
  color: var(--text-dim);
  line-height: 1.4;
}

.danger-btn {
  background: transparent;
  border: 1px solid var(--danger);
  color: var(--danger);
  padding: 10px 16px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  width: 100%;
  margin-top: 4px;
  transition: all 0.15s;
}

.danger-btn:active {
  background: var(--danger);
  color: #fff;
}

.confirm-box {
  background: rgba(255, 107, 107, 0.08);
  border: 1px solid rgba(255, 107, 107, 0.3);
  border-radius: 8px;
  padding: 12px;
  margin-top: 6px;
}

.confirm-text {
  font-size: 0.9rem;
  color: var(--danger);
  font-weight: 600;
  margin-bottom: 8px;
  text-align: center;
}

.confirm-actions {
  display: flex;
  gap: 8px;
}

.confirm-actions .danger-btn {
  margin-top: 0;
  flex: 1;
}

.cancel-btn {
  flex: 1;
  background: transparent;
  border: 1px solid #233554;
  color: var(--text-dim);
  border-radius: 8px;
  cursor: pointer;
}

.cancel-btn:hover {
  color: var(--text);
  border-color: var(--text-dim);
}
</style>
