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
            placeholder="Notiz (z. B. Kaffee, Mittagessen)"
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
  if (props.currentBudget === undefined) return null

  if (!isValidAmount.value) {
    return {
      type: 'impact-neutral',
      icon: '💶',
      text: `Heute verfügbar: ${props.currentBudget.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} €`,
    }
  }

  const diff = props.currentBudget - parsedAmount.value
  if (diff >= 0) {
    return {
      type: 'impact-ok',
      icon: '✓',
      text: `Verbleibt danach: ${diff.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} €`,
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

// Versucht die virtuelle Tastatur programmatisch zu öffnen.
// Mobile Browser blockieren Autofocus ohne User-Gesture; die VirtualKeyboard-API
// (Chrome/Android) hebt diese Einschränkung in standalone-PWAs oft auf.
function requestSoftKeyboard() {
  const navigatorAny = navigator as Navigator & {
    virtualKeyboard?: { show: () => void }
    keyboard?: { show: () => Promise<void> }
  }
  if (navigatorAny.virtualKeyboard?.show) {
    try {
      navigatorAny.virtualKeyboard.show()
      return
    } catch {
      // Fallback unten
    }
  }
  if (navigatorAny.keyboard?.show) {
    navigatorAny.keyboard.show().catch(() => {})
  }
}

function triggerFocus() {
  nextTick(() => {
    if (amountInputRef.value) {
      amountInputRef.value.focus({ preventScroll: false })
      requestSoftKeyboard()
      // Android / mobile virtual keyboard trigger
      amountInputRef.value.click()
    }
  })
}

// Wiederholte Versuche, damit die Tastatur auch beim Cold-Start
// (PWA-Shortcut ohne User-Gesture + Service-Worker-Rendering) zuverlässig aufgeht.
function scheduleFocusRetries(totalAttempts = 8, delayMs = 150) {
  for (let i = 0; i < totalAttempts; i++) {
    setTimeout(triggerFocus, i * delayMs)
  }
}

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      amountInput.value = ''
      noteInput.value = ''
      loadStoredNotes()
      triggerFocus()
      // Fallback delay for Android WebView PWA shortcut rendering
      setTimeout(triggerFocus, 100)
      // Mehrere spätere Versuche gegen das ohne-Gesture-Öffnungslimit
      scheduleFocusRetries()
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
  max-width: 400px;
  max-height: 90dvh;
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
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  overflow-y: auto;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-muted, #8e8e9c);
}

.amount-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.amount-input {
  width: 100%;
  padding: 12px 44px 12px 16px;
  font-size: 1.4rem;
  font-family: var(--font-mono, monospace);
  font-weight: 700;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 12px;
  background: var(--bg-subtle, #1c1c24);
  color: var(--text-main, #f4f4f6);
  outline: none;
  transition: all 0.2s ease;
}

.amount-input:focus {
  border-color: var(--accent-green, #22c55e);
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.15);
}

.amount-input::placeholder {
  color: var(--text-dim, #5c5c6e);
  font-weight: 400;
}

.currency-symbol {
  position: absolute;
  right: 16px;
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--accent-green, #22c55e);
  pointer-events: none;
}

.note-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 0.95rem;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 12px;
  background: var(--bg-subtle, #1c1c24);
  color: var(--text-main, #f4f4f6);
  outline: none;
  transition: all 0.2s ease;
}

.note-input:focus {
  border-color: var(--accent-green, #22c55e);
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.15);
}

.note-input::placeholder {
  color: var(--text-dim, #5c5c6e);
}

.form-actions {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 10px;
  margin-top: 6px;
}

.btn {
  padding: 12px 16px;
  font-size: 0.95rem;
  font-weight: 600;
  border-radius: 9999px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.btn-cancel {
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
}

.btn-cancel:hover,
.btn-cancel:active {
  color: var(--text-main, #f4f4f6);
  background: #242430;
}

.btn-confirm {
  background: var(--accent-green, #22c55e);
  color: #05200e;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(34, 197, 94, 0.25);
}

.btn-confirm:hover:not(:disabled) {
  background: #2ed66b;
  box-shadow: 0 6px 20px rgba(34, 197, 94, 0.35);
}

.btn-confirm:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  box-shadow: none;
}

.spinner-inline {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(5, 32, 14, 0.3);
  border-top-color: #05200e;
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
  font-size: 0.78rem;
  font-weight: 600;
  padding: 5px 10px;
  border-radius: 8px;
}

.impact-neutral {
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-muted, #9494a8);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
}

.impact-ok {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.impact-warning {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  color: var(--accent-red, #ef4444);
  border: 1px solid rgba(239, 68, 68, 0.25);
}

.quick-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 4px;
}

.quick-chip {
  background: #1f1f28;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  font-size: 0.72rem;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.quick-chip:hover {
  background: #2a2a36;
  border-color: var(--accent-green, #22c55e);
}

.quick-chip.active {
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.15));
}
</style>
