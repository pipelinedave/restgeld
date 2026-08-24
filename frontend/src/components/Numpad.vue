<template>
  <div v-if="visible" class="modal-overlay" @click.self="cancel">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Ausgabe buchen</h2>
        <button class="close-btn" aria-label="Schließen" @click="cancel">&times;</button>
      </div>

      <form class="modal-body" @submit.prevent="confirmAll">
        <div class="form-group">
          <label for="expense-amount-input" class="form-label">Betrag (€)</label>
          <div class="amount-input-wrap">
            <input
              id="expense-amount-input"
              ref="amountInputRef"
              v-model="amountInput"
              type="text"
              inputmode="decimal"
              placeholder="0,00"
              enterkeyhint="next"
              autocomplete="off"
              class="amount-input"
              @keydown="handleAmountKeydown"
            />
            <span class="currency-symbol">&euro;</span>
          </div>
        </div>

        <div class="form-group">
          <label for="expense-note-input" class="form-label">Notiz (optional)</label>
          <input
            id="expense-note-input"
            v-model="noteInput"
            type="text"
            inputmode="text"
            placeholder="z. B. Kaffee, Mittagessen"
            maxlength="50"
            enterkeyhint="done"
            class="note-input"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-cancel" @click="cancel">
            Abbrechen
          </button>
          <button type="submit" class="btn btn-confirm" :disabled="!isValidAmount">
            Speichern
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  confirm: [value: number, note: string]
  cancel: []
}>()

const amountInput = ref('')
const noteInput = ref('')
const amountInputRef = ref<HTMLInputElement | null>(null)

const parsedAmount = computed(() => {
  if (!amountInput.value) return 0
  const normalized = amountInput.value.trim().replace(',', '.')
  const num = parseFloat(normalized)
  return isNaN(num) ? 0 : num
})

const isValidAmount = computed(() => parsedAmount.value > 0)

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      amountInput.value = ''
      noteInput.value = ''
      nextTick(() => {
        amountInputRef.value?.focus()
      })
    }
  },
  { immediate: true }
)

function handleAmountKeydown(e: KeyboardEvent) {
  if ((e.key === ',' || e.key === '.') && (amountInput.value.includes(',') || amountInput.value.includes('.'))) {
    e.preventDefault()
  }
}

function cancel() {
  amountInput.value = ''
  noteInput.value = ''
  emit('cancel')
}

function confirmAll() {
  if (!isValidAmount.value) return
  emit('confirm', parsedAmount.value, noteInput.value.trim())
  amountInput.value = ''
  noteInput.value = ''
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
  max-width: 400px;
  max-height: 90dvh;
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
  font-size: 1.25rem;
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
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text, #ccd6f6);
}

.amount-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.amount-input {
  width: 100%;
  padding: 14px 44px 14px 16px;
  font-size: 1.5rem;
  font-weight: 600;
  border: 2px solid #233554;
  border-radius: 10px;
  background: var(--bg, #0a192f);
  color: var(--text, #ccd6f6);
  outline: none;
  transition: border-color 0.15s;
}

.amount-input:focus {
  border-color: var(--accent, #64ffda);
}

.amount-input::placeholder {
  color: #495670;
  font-weight: 400;
}

.currency-symbol {
  position: absolute;
  right: 16px;
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--accent, #64ffda);
  pointer-events: none;
}

.note-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 1rem;
  border: 2px solid #233554;
  border-radius: 10px;
  background: var(--bg, #0a192f);
  color: var(--text, #ccd6f6);
  outline: none;
  transition: border-color 0.15s;
}

.note-input:focus {
  border-color: var(--accent, #64ffda);
}

.note-input::placeholder {
  color: #495670;
}

.form-actions {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 10px;
  margin-top: 4px;
}

.btn {
  padding: 12px 16px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  border: none;
}

.btn-cancel {
  background: transparent;
  border: 1px solid #233554;
  color: var(--text-dim, #8892b0);
}

.btn-cancel:hover,
.btn-cancel:active {
  color: var(--text, #ccd6f6);
  border-color: var(--text-dim, #8892b0);
}

.btn-confirm {
  background: var(--accent, #64ffda);
  color: #0a192f;
}

.btn-confirm:hover:not(:disabled) {
  opacity: 0.95;
}

.btn-confirm:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
