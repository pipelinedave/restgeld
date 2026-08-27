<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Einstellungen</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <!-- Konto & Cloud-Sync Sektion -->
        <section class="setting-section account-section">
          <div class="account-header">
            <div class="account-info">
              <span class="section-title">{{ auth.isLoggedIn.value ? 'Konto & Cloud-Sync' : 'Cloud-Sync & Backup' }}</span>
              <p class="description">
                {{ auth.isLoggedIn.value ? `Angemeldet als ${auth.user.value?.email}` : 'Synchronisiere deine Daten sicher & passwortlos über alle deine Geräte.' }}
              </p>
            </div>
            <button class="account-btn" @click="$emit('open-auth')">
              {{ auth.isLoggedIn.value ? 'Konto verwalten ⚙️' : 'Anmelden / Registrieren ☁️' }}
            </button>
          </div>
        </section>

        <!-- Interaktiver Budget- & Perioden-Konfigurator -->
        <section class="setting-section config-section">
          <!-- Monatsbudget -->
          <div class="slider-group">
            <div class="slider-header">
              <label for="monthly-budget-input" class="section-title">Monatsbudget</label>
              <div class="input-badge">
                <input
                  id="monthly-budget-input"
                  v-model.number="budgetInput"
                  type="number"
                  min="10"
                  max="5000"
                  step="1"
                  placeholder="600"
                  class="num-input"
                  @keyup.enter="handleSaveSettings"
                />
                <span class="currency-symbol">&euro;</span>
              </div>
            </div>

            <input
              v-model.number="budgetInput"
              type="range"
              min="50"
              max="2500"
              step="1"
              class="range-slider"
            />

            <!-- Presets -->
            <div class="preset-row">
              <button
                v-for="p in [300, 450, 600, 1000]"
                :key="p"
                type="button"
                class="preset-chip"
                :class="{ active: budgetInput === p }"
                @click="setBudget(p)"
              >
                {{ p }} &euro;
              </button>
            </div>
          </div>

          <!-- Periodendauer -->
          <div class="slider-group" style="margin-top: 14px;">
            <div class="slider-header">
              <label for="period-days-input" class="section-title">Periodendauer</label>
              <div class="input-badge">
                <input
                  id="period-days-input"
                  v-model.number="daysInput"
                  type="number"
                  min="1"
                  max="365"
                  step="1"
                  placeholder="30"
                  class="num-input"
                  @keyup.enter="handleSaveSettings"
                />
                <span class="currency-symbol">Tage</span>
              </div>
            </div>

            <input
              v-model.number="daysInput"
              type="range"
              min="7"
              max="45"
              step="1"
              class="range-slider"
            />

            <!-- Days Presets -->
            <div class="preset-row">
              <button
                v-for="d in [14, 28, 30, 31]"
                :key="d"
                type="button"
                class="preset-chip"
                :class="{ active: daysInput === d }"
                @click="setDays(d)"
              >
                {{ d }} Tage
              </button>
            </div>
          </div>

          <!-- Live Kalkulator Card: Tages-Restgeld bidirektional gekoppelt -->
          <div class="calc-card">
            <div class="calc-header">
              <span class="calc-label">Tages-Restgeld</span>
              <div class="input-badge">
                <input
                  id="day-budget-input"
                  v-model.number="dayEditor"
                  type="number"
                  min="0.1"
                  step="0.1"
                  placeholder="0,00"
                  class="num-input"
                  @change="setDayBudget(dayEditor)"
                />
                <span class="currency-symbol">&euro;</span>
              </div>
            </div>
            <input
              v-model.number="dayEditor"
              type="range"
              min="0.5"
              :max="Math.max(50, Math.round(budgetInput / Math.max(daysInput, 1)) * 2)"
              step="0.5"
              class="range-slider"
              @change="setDayBudget(dayEditor)"
            />
            <div class="calc-footer">
              <span class="calc-value">&empty; {{ calculatedDailyBudget }} &euro; / Tag</span>
              <button
                class="save-btn"
                :disabled="!isValidSettings"
                @click="handleSaveSettings"
              >
                Speichern
              </button>
            </div>
          </div>

          <p v-if="budgetSavedMsg" class="success-msg">{{ budgetSavedMsg }}</p>
        </section>

        <!-- Design & Theming -->
        <section class="setting-section theme-section">
          <span class="section-title">Design & Akzentfarbe</span>
          <p class="description">
            Wähle dein Lieblings-Farbschema oder stelle eine eigene Farbe ein.
          </p>

          <div class="theme-palette-row">
            <button
              v-for="preset in theme.presets"
              :key="preset.id"
              type="button"
              class="theme-color-btn"
              :class="{ active: theme.currentAccent.value.toLowerCase() === preset.accent.toLowerCase() }"
              :style="{ backgroundColor: preset.accent }"
              :title="preset.name"
              @click="theme.applyTheme(preset.accent)"
            >
              <span v-if="theme.currentAccent.value.toLowerCase() === preset.accent.toLowerCase()" class="check-icon">✓</span>
            </button>

            <!-- Custom Color Picker -->
            <label class="custom-color-picker" title="Eigene Farbe wählen">
              <input
                type="color"
                :value="theme.currentAccent.value"
                class="color-picker-input"
                @input="handleCustomColorChange"
              />
              <span class="custom-color-indicator" :style="{ backgroundColor: theme.currentAccent.value }">
                🎨
              </span>
            </label>
          </div>
        </section>

        <!-- Backup & Archiv -->
        <section class="setting-section backup-zone">
          <span class="section-title">Daten & Archiv</span>
          <p class="description">
            Sichere deine Ausgaben oder wirf einen Blick in frühere Perioden.
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

          <div style="margin-top: 8px;">
            <button type="button" class="archive-trigger-btn" @click="handleOpenArchive">
              📜 Frühere Monate / Archiv ansehen
            </button>
          </div>
        </section>

        <!-- Über Restgeld / Info Zone -->
        <section class="setting-section about-zone">
          <span class="section-title">App-Info & Philosophie</span>
          <p class="description">
            Erfahre mehr über die Prinzipien von Restgeld, Shortcuts und Open-Source-Quellcode.
          </p>
          <button type="button" class="about-trigger-btn" @click="handleOpenAbout">
            ℹ️ Über Restgeld öffnen
          </button>
        </section>

        <!-- Danger Zone -->
        <section class="setting-section danger-zone">
          <span class="section-title danger-title">Neue Periode ab heute starten</span>
          <p class="description">
            Startet deinen Abrechnungszyklus ab heute bei Tag 1 mit dem konfigurierten Budget.
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
import { ref, computed, watch, nextTick } from 'vue'
import { useHaptics } from '../composables/useHaptics'
import { useApi } from '../composables/useApi'
import { useTheme } from '../composables/useTheme'
import { useAuth } from '../composables/useAuth'

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
  (e: 'open-about'): void
  (e: 'open-auth'): void
}>()

const api = useApi()
const haptics = useHaptics()
const theme = useTheme()
const auth = useAuth()
const budgetInput = ref<number>(props.currentMonthlyBudget ?? 450)
const daysInput = ref<number>(props.currentMonthDays ?? 30)
const confirmReset = ref(false)
const budgetSavedMsg = ref('')

const isExporting = ref(false)
const isImporting = ref(false)
const backupMsg = ref('')
const backupMsgType = ref<'success' | 'error'>('success')

// Bidirektionale Tages-Restgeld-Kopplung (Epic 3 #1):
// Tages-Restgeld ist editierbar (Slider/Nummer) und treibt beim Ändern das
// Monatsbudget (Tag * Tage), während das Monatsbudget umgekehrt das Tagesbudget anzeigt.
const dayEditor = ref(0)
const isAdjustingDay = ref(false)

function syncDayEditor() {
  if (isAdjustingDay.value) return
  const days = daysInput.value
  if (days > 0) {
    dayEditor.value = Math.round((budgetInput.value / days) * 100) / 100
  }
}

function setDayBudget(val: number) {
  if (typeof val !== 'number' || !isFinite(val) || val <= 0) return
  isAdjustingDay.value = true
  dayEditor.value = val
  const days = daysInput.value
  if (days > 0) {
    budgetInput.value = Math.round(val * days * 100) / 100
  }
  nextTick(() => {
    isAdjustingDay.value = false
  })
}

watch([budgetInput, daysInput], () => {
  syncDayEditor()
})

watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      isAdjustingDay.value = false
      nextTick(syncDayEditor)
    }
  }
)

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
      if (props.currentMonthlyBudget) {
        budgetInput.value = props.currentMonthlyBudget
      }
      if (props.currentMonthDays) {
        daysInput.value = props.currentMonthDays
      }
      confirmReset.value = false
      budgetSavedMsg.value = ''
      backupMsg.value = ''
    }
  }
)

const isValidSettings = computed(() => {
  const b = budgetInput.value
  const d = daysInput.value
  return typeof b === 'number' && b > 0 && typeof d === 'number' && d > 0 && d <= 365
})

const calculatedDailyBudget = computed(() => {
  if (!isValidSettings.value) return '0,00'
  const daily = budgetInput.value / daysInput.value
  return daily.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
})

function setBudget(val: number) {
  haptics.tap()
  budgetInput.value = val
}

function setDays(val: number) {
  haptics.tap()
  daysInput.value = val
}

function handleSaveSettings() {
  if (!isValidSettings.value) return
  haptics.success()
  emit('update-budget', budgetInput.value, daysInput.value)
  budgetSavedMsg.value = '✓ Einstellungen gespeichert'
  setTimeout(() => {
    budgetSavedMsg.value = ''
  }, 2500)
}

async function handleExport(format: 'json' | 'csv') {
  isExporting.value = true
  haptics.tap()

  try {
    const blob = await api.exportData(format)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const dateStr = new Date().toISOString().slice(0, 10)
    a.download = `restgeld-backup-${dateStr}.${format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    haptics.success()
    backupMsgType.value = 'success'
    backupMsg.value = `Export als ${format.toUpperCase()} erfolgreich!`
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

function handleCustomColorChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target && target.value) {
    theme.applyTheme(target.value)
  }
}

function handleOpenArchive() {
  haptics.tap()
  emit('open-archive')
}

function handleOpenAbout() {
  haptics.tap()
  emit('open-about')
}

function handleResetPeriod() {
  haptics.warning()
  confirmReset.value = false
  emit('new-period', budgetInput.value, daysInput.value)
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
  animation: modal-fade 0.2s ease-out;
}

@keyframes modal-fade {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: #121216;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 20px;
  width: 100%;
  max-width: 440px;
  max-height: 85dvh;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.9);
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
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
}

.close-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  font-size: 1.25rem;
  line-height: 1;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.close-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
}

.modal-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  overflow-y: auto;
}

.setting-section {
  background: #18181e;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 14px;
}

.account-section {
  background: rgba(34, 197, 94, 0.04);
  border: 1px solid rgba(34, 197, 94, 0.15);
}

.account-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.account-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-btn {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.15));
  border: 1px solid var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
  border-radius: 10px;
  padding: 8px 14px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
}

.account-btn:hover {
  background: var(--accent-green, #22c55e);
  color: #000;
}

.slider-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.slider-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 0.8rem;
  color: var(--text-muted, #8e8e9c);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.input-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 8px;
  padding: 2px 8px;
}

.num-input {
  width: 60px;
  background: transparent;
  border: none;
  color: var(--text-main, #f4f4f6);
  font-family: var(--font-mono, monospace);
  font-size: 0.95rem;
  font-weight: 700;
  text-align: right;
  outline: none;
}

.currency-symbol {
  font-size: 0.8rem;
  color: var(--text-dim, #5c5c6e);
  font-weight: 600;
}

.range-slider {
  width: 100%;
  accent-color: var(--accent-green, #22c55e);
  cursor: pointer;
}

.preset-row {
  display: flex;
  gap: 6px;
}

.preset-chip {
  flex: 1;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  padding: 5px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.preset-chip.active,
.preset-chip:hover {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.12));
  color: var(--accent-green, #22c55e);
  border-color: rgba(34, 197, 94, 0.3);
}

.calc-card {
  margin-top: 12px;
  padding: 10px 12px;
  background: rgba(34, 197, 94, 0.06);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.calc-header,
.calc-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.calc-left {
  display: flex;
  flex-direction: column;
}

.calc-label {
  font-size: 0.68rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.calc-value {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--accent-green, #22c55e);
  font-family: var(--font-mono, monospace);
}

.save-btn {
  background: var(--accent-green, #22c55e);
  color: #05200e;
  border: none;
  padding: 7px 16px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.success-msg {
  color: var(--accent-green, #22c55e);
  font-size: 0.8rem;
  margin: 6px 0 0;
  font-weight: 600;
}

.description {
  color: var(--text-dim, #5c5c6e);
  font-size: 0.75rem;
  margin: 4px 0 10px;
  line-height: 1.35;
}

.backup-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.backup-btn {
  flex: 1;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  padding: 8px;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.backup-btn:hover:not(:disabled) {
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}

.import-wrap {
  width: 100%;
}

.import-btn {
  display: block;
  width: 100%;
  background: rgba(255, 255, 255, 0.02);
  border: 1px dashed var(--border-color, rgba(255, 255, 255, 0.12));
  color: var(--text-muted, #8e8e9c);
  padding: 8px;
  border-radius: 8px;
  font-size: 0.78rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.15s;
}

.import-btn:hover {
  border-color: var(--text-main, #f4f4f6);
  color: var(--text-main, #f4f4f6);
}

.file-input-hidden {
  display: none;
}

.backup-msg {
  font-size: 0.75rem;
  margin-top: 6px;
}

.backup-msg.success { color: var(--accent-green, #22c55e); }
.backup-msg.error { color: var(--accent-red, #ef4444); }

.archive-trigger-btn,
.about-trigger-btn {
  width: 100%;
  display: block;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-main, #f4f4f6);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
  text-align: center;
}

.archive-trigger-btn:hover,
.about-trigger-btn:hover {
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}

.theme-palette-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.theme-color-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #05200e;
  font-weight: 800;
  font-size: 0.9rem;
  transition: transform 0.15s, border-color 0.15s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
}

.theme-color-btn:hover {
  transform: scale(1.1);
}

.theme-color-btn.active {
  border-color: #ffffff;
  transform: scale(1.12);
  box-shadow: 0 0 12px rgba(255, 255, 255, 0.4);
}

.check-icon {
  font-weight: 900;
}

.custom-color-picker {
  position: relative;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed var(--border-color, rgba(255, 255, 255, 0.2));
  background: rgba(255, 255, 255, 0.04);
}

.color-picker-input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.custom-color-indicator {
  font-size: 1rem;
}

.danger-title {
  color: var(--accent-red, #ef4444);
}

.danger-btn {
  width: 100%;
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  border: 1px solid rgba(239, 68, 68, 0.25);
  color: var(--accent-red, #ef4444);
  padding: 9px;
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

.danger-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

.confirm-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: rgba(239, 68, 68, 0.08);
  padding: 10px;
  border-radius: 8px;
}

.confirm-text {
  font-size: 0.78rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  text-align: center;
}

.confirm-actions {
  display: flex;
  gap: 6px;
}

.cancel-btn {
  flex: 1;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
  padding: 6px;
  border-radius: 6px;
  font-size: 0.78rem;
  cursor: pointer;
}
</style>
