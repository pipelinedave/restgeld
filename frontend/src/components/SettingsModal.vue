<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{ i18n.t('settings.title') }}</h2>
        <button class="close-btn" :aria-label="i18n.t('common.close')" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <!-- Konto & Cloud-Sync Sektion -->
        <section class="setting-section account-section">
          <div class="account-header">
            <div class="account-info">
              <span class="section-title">{{ auth.isLoggedIn.value ? (auth.user.value?.email) : i18n.t('auth.title') }}</span>
              <p class="description">
                {{ auth.isLoggedIn.value ? (auth.user.value?.plan === 'pro' ? 'Restgeld PRO Subscribed 👑' : 'Restgeld Free Tier ☁️') : i18n.t('auth.subtitle') }}
              </p>
            </div>
            <button class="account-btn" @click="$emit('open-auth')">
              {{ auth.isLoggedIn.value ? '⚙️' : '☁️' }}
            </button>
          </div>
        </section>

        <!-- Interaktiver Budget- & Perioden-Konfigurator -->
        <section class="setting-section config-section">
          <!-- Monatsbudget -->
          <div class="slider-group">
            <div class="slider-header">
              <label for="monthly-budget-input" class="section-title">{{ i18n.t('settings.monthly_budget') }}</label>
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
                <span class="currency-symbol">{{ i18n.currencySymbol }}</span>
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
                {{ p }} {{ i18n.currencySymbol }}
              </button>
            </div>
          </div>

          <!-- Periodendauer -->
          <div class="slider-group" style="margin-top: 14px;">
            <div class="slider-header">
              <label for="period-days-input" class="section-title">{{ i18n.t('settings.period_days') }}</label>
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
                <span class="currency-symbol">{{ i18n.t('streak.days_unit') }}</span>
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

          <!-- Tages-Restgeld Slider (Bidirektional) -->
          <div class="slider-group" style="margin-top: 14px;">
            <div class="slider-header">
              <label for="daily-budget-input" class="section-title">{{ i18n.t('settings.desired_daily') }}</label>
              <div class="input-badge">
                <input
                  id="daily-budget-input"
                  v-model.number="dailyInput"
                  type="number"
                  min="1"
                  max="200"
                  step="0.5"
                  placeholder="15"
                  class="num-input"
                  @input="handleDailyInputChange"
                />
                <span class="currency-symbol">{{ i18n.currencySymbol }}/{{ i18n.t('streak.days_unit') }}</span>
              </div>
            </div>

            <input
              v-model.number="dailyInput"
              type="range"
              min="5"
              max="100"
              step="0.5"
              class="range-slider"
              @input="handleDailyInputChange"
            />
          </div>

          <!-- Live Kalkulator Card -->
          <div class="calc-card">
            <div class="calc-left">
              <span class="calc-label">{{ i18n.t('settings.calculated_monthly') }}</span>
              <span class="calc-value">{{ calculatedMonthlyTotal }} {{ i18n.currencySymbol }}</span>
            </div>
            <button
              class="save-btn"
              :disabled="!isValidSettings"
              @click="handleSaveSettings"
            >
              {{ i18n.t('settings.save_btn') }}
            </button>
          </div>

          <p v-if="budgetSavedMsg" class="success-msg">{{ budgetSavedMsg }}</p>
        </section>

        <!-- Design & Theming -->
        <section class="setting-section theme-section">
          <span class="section-title">{{ i18n.t('settings.theme_heading') }}</span>
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

        <!-- Sound & Haptik-Feedback -->
        <section class="setting-section sound-section">
          <span class="section-title">{{ i18n.t('settings.sound_heading') }}</span>
          <p class="description">
            {{ i18n.t('settings.sound_desc') }}
          </p>
          <label class="toggle-control-label">
            <input
              type="checkbox"
              :checked="haptics.soundEnabled.value"
              @change="haptics.setSoundEnabled(($event.target as HTMLInputElement).checked)"
            />
            <span>{{ i18n.t('settings.sound_toggle') }}</span>
          </label>
        </section>

        <!-- Sprache / Language -->
        <section class="setting-section language-zone">
          <span class="section-title">{{ i18n.t('settings.language_heading') }}</span>
          <div class="language-chips">
            <button
              v-for="lang in i18n.languages"
              :key="lang.code"
              type="button"
              class="lang-chip"
              :class="{ active: i18n.currentLocale.value === lang.code }"
              @click="i18n.setLocale(lang.code)"
            >
              <span class="lang-flag">{{ lang.flag }}</span>
              <span class="lang-name">{{ lang.name }}</span>
            </button>
          </div>
        </section>

        <!-- Währung / Currency -->
        <section class="setting-section currency-zone">
          <span class="section-title">{{ i18n.t('settings.currency_heading') }}</span>
          <div class="language-chips">
            <button
              v-for="curr in i18n.currencies"
              :key="curr.code"
              type="button"
              class="lang-chip"
              :class="{ active: i18n.currentCurrency.value === curr.code }"
              @click="i18n.setCurrency(curr.code)"
            >
              <span class="lang-flag">{{ curr.symbol }}</span>
              <span class="lang-name">{{ curr.name }}</span>
            </button>
          </div>
        </section>

        <!-- Backup & Archiv -->
        <section class="setting-section backup-zone">
          <span class="section-title">{{ i18n.t('settings.backup_heading') }}</span>
          <p class="description">
            {{ i18n.t('settings.backup_desc') }}
          </p>

          <div class="backup-actions">
            <button class="backup-btn" :disabled="isExporting" @click="handleExport('csv')">
              {{ i18n.t('settings.export_csv') }}
            </button>
            <button class="backup-btn" :disabled="isExporting" @click="handleExport('json')">
              {{ i18n.t('settings.export_json') }}
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
              <span>{{ isImporting ? i18n.t('settings.importing') : i18n.t('settings.import_backup') }}</span>
            </label>
          </div>
          <p v-if="backupMsg" class="backup-msg" :class="backupMsgType">{{ backupMsg }}</p>

          <div style="margin-top: 8px;">
            <button type="button" class="archive-trigger-btn" @click="handleOpenArchive">
              {{ i18n.t('settings.archive_trigger') }}
            </button>
          </div>
        </section>

        <!-- Über Restgeld / Info Zone -->
        <section class="setting-section about-zone">
          <span class="section-title">{{ i18n.t('settings.about_heading') }}</span>
          <p class="description">
            {{ i18n.t('settings.about_desc') }}
          </p>
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <button type="button" class="about-trigger-btn" @click="handleOpenAbout">
              {{ i18n.t('settings.about_trigger') }}
            </button>
            <button type="button" class="about-trigger-btn" @click="handleOpenStatus">
              📊 {{ i18n.t('monitoring.title') }}
            </button>
          </div>
        </section>

        <!-- Danger Zone -->
        <section class="setting-section danger-zone">
          <span class="section-title danger-title">{{ i18n.t('settings.reset_heading') }}</span>
          <p class="description">
            {{ i18n.t('settings.reset_desc') }}
          </p>

          <div v-if="!confirmReset">
            <button class="danger-btn" @click="confirmReset = true">
              {{ i18n.t('settings.reset_btn') }}
            </button>
          </div>
          <div v-else class="confirm-box">
            <p class="confirm-text">{{ i18n.t('settings.reset_confirm_body') }}</p>
            <div class="confirm-actions">
              <button class="danger-btn confirm" @click="handleResetPeriod">
                {{ i18n.t('settings.reset_confirm_btn') }}
              </button>
              <button class="cancel-btn" @click="confirmReset = false">
                {{ i18n.t('settings.reset_cancel_btn') }}
              </button>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, watchEffect } from 'vue'
import { useHaptics } from '../composables/useHaptics'
import { useApi } from '../composables/useApi'
import { useTheme } from '../composables/useTheme'
import { useAuth } from '../composables/useAuth'
import { useI18n } from '../composables/useI18n'

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
  (e: 'open-status'): void
}>()

const api = useApi()
const haptics = useHaptics()
const theme = useTheme()
const auth = useAuth()
const i18n = useI18n()
const budgetInput = ref<number>(props.currentMonthlyBudget ?? 450)
const daysInput = ref<number>(props.currentMonthDays ?? 31)
const dailyInput = ref<number>(Math.round(((props.currentMonthlyBudget ?? 450) / (props.currentMonthDays ?? 31)) * 100) / 100)
const confirmReset = ref(false)
const budgetSavedMsg = ref('')

const isExporting = ref(false)
const isImporting = ref(false)
const backupMsg = ref('')
const backupMsgType = ref<'success' | 'error'>('success')

watch(
  () => props.currentMonthDays,
  (d) => {
    if (typeof d === 'number' && d > 0) {
      daysInput.value = d
    }
  },
  { immediate: true }
)

watch(
  () => props.currentMonthlyBudget,
  (b) => {
    if (typeof b === 'number' && b > 0) {
      budgetInput.value = b
    }
  },
  { immediate: true }
)

watch(
  () => props.visible,
  (isVis) => {
    if (isVis) {
      if (typeof props.currentMonthlyBudget === 'number' && props.currentMonthlyBudget > 0) {
        budgetInput.value = props.currentMonthlyBudget
      }
      if (typeof props.currentMonthDays === 'number' && props.currentMonthDays > 0) {
        daysInput.value = props.currentMonthDays
      }
      if (budgetInput.value && daysInput.value) {
        dailyInput.value = Math.round((budgetInput.value / daysInput.value) * 100) / 100
      }
      confirmReset.value = false
      budgetSavedMsg.value = ''
      backupMsg.value = ''
    }
  }
)

watch(
  () => budgetInput.value,
  (b) => {
    if (b && daysInput.value && daysInput.value > 0) {
      dailyInput.value = Math.round((b / daysInput.value) * 100) / 100
    }
  }
)

function handleDailyInputChange() {
  if (dailyInput.value && daysInput.value && daysInput.value > 0) {
    budgetInput.value = Math.round(dailyInput.value * daysInput.value)
  }
}

const isValidSettings = computed(() => {
  const b = budgetInput.value
  const d = daysInput.value
  return typeof b === 'number' && b > 0 && typeof d === 'number' && d > 0 && d <= 365
})

const calculatedMonthlyTotal = computed(() => {
  if (!isValidSettings.value) return '0,00'
  return budgetInput.value.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
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

function handleOpenStatus() {
  haptics.tap()
  emit('open-status')
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

.toggle-control-label {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.8rem;
  color: var(--text-main, #f4f4f6);
  cursor: pointer;
  user-select: none;
  font-weight: 500;
}

.toggle-control-label input[type="checkbox"] {
  accent-color: var(--accent-green, #22c55e);
  width: 16px;
  height: 16px;
  cursor: pointer;
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

.language-chips {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-top: 6px;
}

.lang-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 8px;
  padding: 8px 12px;
  color: var(--text-main, #f4f4f6);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

.lang-chip:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.15);
}

.lang-chip.active {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.15));
  border-color: var(--accent-green, #22c55e);
  color: var(--accent-green, #22c55e);
}

.lang-flag {
  font-size: 1.1rem;
}
</style>
