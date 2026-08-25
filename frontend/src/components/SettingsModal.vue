<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Einstellungen</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <section class="setting-section">
          <label for="monthly-budget-input" class="section-title">Budget (€)</label>
          <div class="input-row">
            <input
              id="monthly-budget-input"
              v-model.number="budgetInput"
              type="number"
              min="1"
              step="10"
              placeholder="z. B. 600"
              @keyup.enter="handleSaveSettings"
            />
          </div>

          <label for="period-days-input" class="section-title" style="margin-top: 14px;">Dauer der Periode (Tage)</label>
          <div class="input-row">
            <input
              id="period-days-input"
              v-model.number="daysInput"
              type="number"
              min="1"
              max="365"
              step="1"
              placeholder="z. B. 30"
              @keyup.enter="handleSaveSettings"
            />
            <button class="action-btn" :disabled="!isValidSettings" @click="handleSaveSettings">
              Speichern
            </button>
          </div>
          <p v-if="budgetSavedMsg" class="success-msg">{{ budgetSavedMsg }}</p>
        </section>

        <section class="setting-section backup-zone">
          <span class="section-title">Daten & Backup</span>
          <p class="description">
            Exportiere deine Daten als CSV für Excel oder erstelle ein JSON-Backup zur Wiederherstellung.
          </p>

          <div class="backup-actions">
            <button class="backup-btn" :disabled="isExporting" @click="handleExport('csv')">
              CSV (Excel)
            </button>
            <button class="backup-btn" :disabled="isExporting" @click="handleExport('json')">
              JSON Backup
            </button>
          </div>

          <div class="import-wrap">
            <label class="import-btn" :class="{ disabled: isImporting }">
              <input
                type="file"
                accept=".json,.csv"
                class="file-input-hidden"
                :disabled="isImporting"
                @change="handleFileInput"
              />
              <span>{{ isImporting ? 'Importiere...' : 'Backup importieren (JSON/CSV)' }}</span>
            </label>
          </div>
          <p v-if="backupMsg" class="backup-msg" :class="backupMsgType">{{ backupMsg }}</p>

          <div style="margin-top: 6px;">
            <button type="button" class="archive-trigger-btn" @click="handleOpenArchive">
              📜 Frühere Monate / Archiv ansehen
            </button>
          </div>
        </section>

        <section class="setting-section danger-zone">
          <span class="section-title danger-title">Neue Periode ab heute starten</span>
          <p class="description">
            Startet deinen Gehalts-/Abrechnungszyklus ab heute bei Tag 1 mit dem oben konfigurierten Budget und der Dauer. Bisherige Ausgaben dieser Periode werden zurückgesetzt.
          </p>

          <div v-if="!confirmReset">
            <button class="danger-btn" @click="confirmReset = true">
              Neue Periode ab heute starten
            </button>
          </div>
          <div v-else class="confirm-box">
            <p class="confirm-text">Wirklich ab heute neu starten?</p>
            <div class="confirm-actions">
              <button class="danger-btn confirm" @click="handleResetPeriod">
                Ja, ab heute starten
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
import { useHaptics } from '../composables/useHaptics'
import { useApi } from '../composables/useApi'

const props = defineProps<{
  visible: boolean
  currentMonthlyBudget?: number
  currentMonthDays?: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'update-budget', monthlyTotal: number, days?: number): void
  (e: 'new-period', monthlyTotal?: number, days?: number): void
  (e: 'data-imported', count: number): void
  (e: 'open-archive'): void
}>()

const api = useApi()
const haptics = useHaptics()
const budgetInput = ref<number | null>(props.currentMonthlyBudget ?? null)
const daysInput = ref<number | null>(props.currentMonthDays ?? null)
const confirmReset = ref(false)
const budgetSavedMsg = ref('')

const isExporting = ref(false)
const isImporting = ref(false)
const backupMsg = ref('')
const backupMsgType = ref<'success' | 'error'>('success')

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
  () => props.currentMonthDays,
  (newVal) => {
    if (newVal) {
      daysInput.value = newVal
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
      backupMsg.value = ''
      if (props.currentMonthlyBudget) {
        budgetInput.value = props.currentMonthlyBudget
      }
      if (props.currentMonthDays) {
        daysInput.value = props.currentMonthDays
      }
    }
  }
)

const isValidSettings = computed(() => {
  const validBudget = typeof budgetInput.value === 'number' && budgetInput.value > 0
  const validDays = daysInput.value === null || (typeof daysInput.value === 'number' && daysInput.value > 0)
  return validBudget && validDays
})

function handleSaveSettings() {
  if (isValidSettings.value && budgetInput.value) {
    haptics.success()
    emit('update-budget', budgetInput.value, daysInput.value || undefined)
    budgetSavedMsg.value = 'Einstellungen erfolgreich gespeichert!'
    setTimeout(() => {
      budgetSavedMsg.value = ''
    }, 2500)
  }
}

async function handleExport(format: 'json' | 'csv') {
  isExporting.value = true
  haptics.tap()
  try {
    const blob = await api.exportData(format)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    const dateStr = new Date().toISOString().split('T')[0]
    a.href = url
    a.download = `restgeld-${format === 'csv' ? 'export' : 'backup'}-${dateStr}.${format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    haptics.success()
    backupMsgType.value = 'success'
    backupMsg.value = `${format.toUpperCase()} erfolgreich heruntergeladen!`
    setTimeout(() => {
      backupMsg.value = ''
    }, 3000)
  } catch (err: any) {
    haptics.error()
    backupMsgType.value = 'error'
    backupMsg.value = 'Export fehlgeschlagen'
  } finally {
    isExporting.value = false
  }
}

async function handleFileInput(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  isImporting.value = true
  haptics.tap()

  try {
    const content = await file.text()
    const isCsv = file.name.endsWith('.csv') || content.includes(';') || content.includes('Datum')
    const res = await api.importData(content, isCsv)

    haptics.success()
    backupMsgType.value = 'success'
    backupMsg.value = `${res.imported} Ausgabe(n) erfolgreich importiert!`
    emit('data-imported', res.imported)
    setTimeout(() => {
      backupMsg.value = ''
    }, 3500)
  } catch (err: any) {
    haptics.error()
    backupMsgType.value = 'error'
    backupMsg.value = err.message || 'Import fehlgeschlagen'
  } finally {
    isImporting.value = false
    target.value = ''
  }
}

function handleOpenArchive() {
  haptics.tap()
  emit('open-archive')
}

function handleResetPeriod() {
  haptics.warning()
  confirmReset.value = false
  emit(
    'new-period',
    isValidSettings.value && budgetInput.value ? budgetInput.value : undefined,
    isValidSettings.value && daysInput.value ? daysInput.value : undefined
  )
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
  gap: 22px;
  overflow-y: auto;
}

.setting-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.3px;
  color: var(--text-main, #f4f4f6);
}

.input-row {
  display: flex;
  gap: 8px;
}

input[type='number'] {
  flex: 1;
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 12px;
  padding: 10px 14px;
  color: var(--text-main, #f4f4f6);
  font-size: 0.95rem;
  font-family: var(--font-mono, monospace);
  outline: none;
  transition: all 0.2s ease;
}

input[type='number']:focus {
  border-color: var(--accent-green, #22c55e);
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.15);
}

.action-btn {
  background: var(--accent-green, #22c55e);
  border: 1px solid transparent;
  color: #05200e;
  padding: 10px 18px;
  border-radius: 9999px;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover:not(:disabled) {
  background: #2ed66b;
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.success-msg {
  font-size: 0.8rem;
  color: var(--accent-green, #22c55e);
  margin-top: 4px;
}

.danger-zone {
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  padding-top: 16px;
}

.danger-title {
  color: var(--accent-red, #ef4444);
}

.description {
  font-size: 0.8rem;
  color: var(--text-dim, #5c5c6e);
  line-height: 1.4;
}

.danger-btn {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: var(--accent-red, #ef4444);
  padding: 10px 16px;
  border-radius: 9999px;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
  width: 100%;
  margin-top: 4px;
  transition: all 0.15s;
}

.danger-btn:hover,
.danger-btn:active {
  background: var(--accent-red, #ef4444);
  color: #fff;
}

.confirm-box {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  padding: 12px;
  margin-top: 6px;
}

.confirm-text {
  font-size: 0.85rem;
  color: var(--accent-red, #ef4444);
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
  background: var(--bg-subtle, #1c1c24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
  border-radius: 9999px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
}

.cancel-btn:hover {
  color: var(--text-main, #f4f4f6);
}

.backup-zone {
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  padding-top: 16px;
}

.backup-actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.backup-btn {
  flex: 1;
  background: #1f1f28;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  padding: 9px 12px;
  border-radius: 10px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.backup-btn:hover:not(:disabled),
.backup-btn:active:not(:disabled) {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.15));
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}

.backup-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.import-wrap {
  margin-top: 8px;
}

.import-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  background: transparent;
  border: 1px dashed var(--border-color, rgba(255, 255, 255, 0.15));
  color: var(--text-muted, #8e8e9c);
  padding: 10px 12px;
  border-radius: 10px;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.import-btn:hover,
.import-btn:active {
  border-color: var(--accent-green, #22c55e);
  color: var(--text-main, #f4f4f6);
}

.import-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

.file-input-hidden {
  display: none;
}

.backup-msg {
  font-size: 0.8rem;
  margin-top: 6px;
}

.backup-msg.success {
  color: var(--accent-green, #22c55e);
}

.backup-msg.error {
  color: var(--accent-red, #ef4444);
}

.archive-trigger-btn {
  width: 100%;
  background: #1f1f28;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  padding: 10px 12px;
  border-radius: 10px;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  text-align: center;
}

.archive-trigger-btn:hover {
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}
</style>
