<template>
  <div v-if="visible" class="modal-overlay" @click.self="cancel">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Ausgabe buchen</h2>
        <button class="close-btn" aria-label="Schließen" :disabled="isSaving" @click="cancel">&times;</button>
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
              :disabled="isSaving"
              class="amount-input"
              @keydown="handleAmountKeydown"
            />
            <span class="currency-symbol">&euro;</span>
          </div>

          <!-- Live Budget Impact Vorschau -->
          <div v-if="liveImpact" class="impact-preview" :class="liveImpact.type">
            <span class="impact-icon">{{ liveImpact.icon }}</span>
            <span class="impact-text">{{ liveImpact.text }}</span>
          </div>
        </div>

        <div class="form-group">
          <label for="expense-note-input" class="form-label">Notiz (optional)</label>
          
          <!-- Quick Note Chips -->
          <div v-if="availableChips.length > 0" class="quick-chips">
            <button
              v-for="chip in availableChips"
              :key="chip"
              type="button"
              class="quick-chip"
              :class="{ active: noteInput === chip }"
              :disabled="isSaving"
              @click="selectChip(chip)"
            >
              {{ chip }}
            </button>
          </div>

          <input
            id="expense-note-input"
            v-model="noteInput"
            type="text"
            inputmode="text"
            placeholder="z. B. Kaffee, Mittagessen"
            maxlength="50"
            enterkeyhint="done"
            :disabled="isSaving"
            class="note-input"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-cancel" :disabled="isSaving" @click="cancel">
            Abbrechen
          </button>
          <button type="submit" class="btn btn-confirm" :disabled="!isValidAmount || isSaving">
            <span v-if="isSaving" class="spinner-inline" aria-label="Wird gespeichert..."></span>
            <span>{{ isSaving ? 'Wird gebucht...' : 'Speichern' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useHaptics } from '../composables/useHaptics'
import type { Expense } from '../composables/useApi'

const props = withDefaults(
  defineProps<{
    visible: boolean
    isSaving?: boolean
    currentBudget?: number
    savings?: number
    recentExpenses?: Expense[]
  }>(),
  {
    isSaving: false,
    currentBudget: undefined,
    savings: undefined,
    recentExpenses: () => [],
  }
)

const emit = defineEmits<{
  confirm: [value: number, note: string]
  cancel: []
}>()

const haptics = useHaptics()
const amountInput = ref('')
const noteInput = ref('')
const amountInputRef = ref<HTMLInputElement | null>(null)
const storedNotes = ref<string[]>([])

const RECENT_NOTES_KEY = 'restgeld_recent_notes'
const DEFAULT_CHIPS = ['Kaffee', 'Mittagessen', 'Einkauf', 'Snack']

function loadStoredNotes() {
  try {
    const raw = localStorage.getItem(RECENT_NOTES_KEY)
    if (raw) {
      storedNotes.value = JSON.parse(raw)
    }
  } catch {
    storedNotes.value = []
  }
}

function saveNoteToStorage(note: string) {
  if (!note || note.trim().length === 0) return
  const clean = note.trim()
  const list = [clean, ...storedNotes.value.filter((n) => n.toLowerCase() !== clean.toLowerCase())].slice(0, 6)
  storedNotes.value = list
  try {
    localStorage.setItem(RECENT_NOTES_KEY, JSON.stringify(list))
  } catch {
    // Ignore storage errors
  }
}

const availableChips = computed(() => {
  const fromExpenses = props.recentExpenses
    ? props.recentExpenses.map((e) => e.note).filter((n): n is string => Boolean(n && n.trim()))
    : []
  const combined = Array.from(new Set([...storedNotes.value, ...fromExpenses, ...DEFAULT_CHIPS]))
  return combined.slice(0, 5)
})

const parsedAmount = computed(() => {
  if (!amountInput.value) return 0
  const normalized = amountInput.value.trim().replace(',', '.')
  const num = parseFloat(normalized)
  return isNaN(num) ? 0 : num
})

const isValidAmount = computed(() => parsedAmount.value > 0)

const liveImpact = computed(() => {
  if (!isValidAmount.value || props.currentBudget === undefined) return null

  const diff = props.currentBudget - parsedAmount.value
  if (diff >= 0) {
    return {
      type: 'impact-ok',
      icon: '✓',
      text: `Heute verbleiben: ${diff.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} €`,
    }
  } else {
    const over = Math.abs(diff)
    return {
      type: 'impact-warning',
      icon: '⚠️',
      text: `Überzieht Tagesbudget um ${over.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} €`,
    }
  }
})

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      amountInput.value = ''
      noteInput.value = ''
      loadStoredNotes()
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

function selectChip(chip: string) {
  haptics.tap()
  if (noteInput.value === chip) {
    noteInput.value = ''
  } else {
    noteInput.value = chip
  }
}

function cancel() {
  if (props.isSaving) return
  haptics.tap()
  amountInput.value = ''
  noteInput.value = ''
  emit('cancel')
}

function confirmAll() {
  if (!isValidAmount.value || props.isSaving) return
  haptics.tap()
  const note = noteInput.value.trim()
  if (note) {
    saveNoteToStorage(note)
  }
  emit('confirm', parsedAmount.value, note)
}

onMounted(loadStoredNotes)
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-confirm:hover:not(:disabled) {
  opacity: 0.95;
}

.btn-confirm:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.spinner-inline {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(10, 25, 47, 0.3);
  border-top-color: #0a192f;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.impact-preview {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.8rem;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 6px;
}

.impact-ok {
  background: rgba(100, 255, 218, 0.08);
  color: var(--accent, #64ffda);
  border: 1px solid rgba(100, 255, 218, 0.2);
}

.impact-warning {
  background: rgba(255, 107, 107, 0.08);
  color: #ff6b6b;
  border: 1px solid rgba(255, 107, 107, 0.25);
}

.quick-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 4px;
}

.quick-chip {
  background: rgba(204, 214, 246, 0.05);
  border: 1px solid #233554;
  color: var(--text-dim, #8892b0);
  font-size: 0.75rem;
  padding: 4px 10px;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.15s;
}

.quick-chip:hover,
.quick-chip.active {
  border-color: var(--accent, #64ffda);
  color: var(--accent, #64ffda);
  background: rgba(100, 255, 218, 0.08);
}
</style>
