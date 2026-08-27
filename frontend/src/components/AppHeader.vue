<template>
  <header class="app-header">
    <div class="brand">
      <span class="brand-title">restgeld</span>
      <span class="brand-dot">.</span>
    </div>

    <div class="header-right">
      <!-- 1. Status Badge -->
      <button
        type="button"
        class="status-badge"
        :class="badgeClass"
        :title="badgeTooltip"
        @click="toggleStatusPopover"
      >
        <span class="status-dot" :class="dotClass"></span>
        <span class="status-text">{{ badgeText }}</span>
      </button>

      <!-- 2. Streak Badge & Button (🔥) -->
      <button
        type="button"
        class="header-icon-btn streak-btn"
        :class="{ active: (streak?.currentStreak ?? 0) > 0 }"
        title="Streak & Spartage ansehen"
        @click="toggleStreakPopover"
      >
        <span class="icon-emoji">🔥</span>
        <span class="btn-badge">{{ streak?.currentStreak ?? 0 }}</span>
      </button>

      <!-- 3. Monatsende-Prognose Glaskugel (🔮) -->
      <button
        type="button"
        class="header-icon-btn projection-btn"
        :class="projection?.status === 'saving' ? 'proj-saving' : 'proj-deficit'"
        title="Monatsende-Prognose anzeigen"
        @click="toggleProjectionPopover"
      >
        <span class="icon-emoji">🔮</span>
      </button>

      <!-- 4. Settings Button (⚙️) -->
      <button class="settings-btn" aria-label="Einstellungen" title="Einstellungen" @click="$emit('open-settings')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </button>
    </div>

    <!-- 1. Health & Uptime Popover -->
    <div v-if="activePopover === 'status'" class="popover-backdrop" @click="activePopover = null">
      <div class="header-popover" @click.stop>
        <div class="popover-header">
          <span class="popover-title">System- & Sync-Status</span>
          <button class="popover-close" @click="activePopover = null">&times;</button>
        </div>

        <div class="status-list">
          <div class="status-row">
            <span class="status-row-label">Netzwerk-Verbindung</span>
            <span class="status-row-val" :class="isOffline ? 'val-error' : 'val-ok'">
              {{ isOffline ? 'Offline' : 'Verbunden (Online)' }}
            </span>
          </div>

          <div class="status-row">
            <span class="status-row-label">API Service Health</span>
            <span class="status-row-val" :class="apiHealthy ? 'val-ok' : 'val-error'">
              {{ apiHealthy ? `OK (${latencyMs}ms)` : 'Nicht erreichbar' }}
            </span>
          </div>

          <div class="status-row">
            <span class="status-row-label">PostgreSQL Datenbank</span>
            <span class="status-row-val" :class="isDbConnected ? 'val-ok' : 'val-warn'">
              {{ isDbConnected ? 'Verbunden & bereit' : (dbStatus === 'offline' ? 'Offline' : 'Nicht erreichbar') }}
            </span>
          </div>

          <div class="status-row">
            <span class="status-row-label">Offline-Outbox Queue</span>
            <span class="status-row-val" :class="pendingSyncCount > 0 ? 'val-warn' : 'val-ok'">
              {{ pendingSyncCount > 0 ? `${pendingSyncCount} ausstehend` : 'Keine ausstehenden Syncs' }}
            </span>
          </div>
        </div>

        <div class="popover-footer">
          <span class="uptime-hint">Uptime & Monitoring aktiv</span>
          <button class="recheck-btn" @click="checkHealth">Aktualisieren</button>
        </div>
      </div>
    </div>

    <!-- 2. Streak Popover -->
    <div v-if="activePopover === 'streak'" class="popover-backdrop" @click="activePopover = null">
      <div class="header-popover streak-popover" @click.stop>
        <div class="popover-header">
          <span class="popover-title">🔥 Spar-Streak & Disziplin</span>
          <button class="popover-close" @click="activePopover = null">&times;</button>
        </div>

        <div class="streak-popover-body">
          <div class="streak-hero-stat">
            <span class="streak-flame-icon">🔥</span>
            <div class="streak-hero-text">
              <span class="streak-big-val">{{ streak?.currentStreak ?? 0 }} Tage</span>
              <span class="streak-sub-val">im Tagesbudget geblieben</span>
            </div>
          </div>

          <div class="streak-stats-grid">
            <div class="streak-stat-box">
              <span class="stat-box-label">Längster Streak</span>
              <span class="stat-box-num">{{ streak?.longestStreak ?? 0 }} Tage</span>
            </div>
            <div class="streak-stat-box">
              <span class="stat-box-label">Null-Ausgaben-Tage</span>
              <span class="stat-box-num">🎯 {{ streak?.noSpendDays ?? 0 }}</span>
            </div>
          </div>

          <p class="streak-motivation-text">
            {{ getStreakMotivation(streak?.currentStreak ?? 0) }}
          </p>
        </div>
      </div>
    </div>

    <!-- 3. Monatsende-Prognose Popover (🔮) -->
    <div v-if="activePopover === 'projection'" class="popover-backdrop" @click="activePopover = null">
      <div class="header-popover projection-popover" @click.stop>
        <div class="popover-header">
          <span class="popover-title">🔮 Monatsende-Prognose</span>
          <button class="popover-close" @click="activePopover = null">&times;</button>
        </div>

        <div class="projection-popover-body">
          <div class="proj-hero-card" :class="isProjectedSaving ? 'card-saving' : 'card-deficit'">
            <span class="proj-hero-label">PROJIZIERTE ERSPARNIS</span>
            <span class="proj-hero-val">{{ isProjectedSaving ? '+' : '' }}{{ formatAmount(projection?.projectedSavings ?? 0) }} &euro;</span>
            <span class="proj-hero-sub">
              {{ isProjectedSaving ? 'Voraussichtlicher Sparpuffer zum Periodenende' : 'Voraussichtliches Monats-Defizit' }}
            </span>
          </div>

          <div class="proj-stats-grid">
            <div class="proj-stat-box">
              <span class="stat-box-label">Vorauss. Ausgaben</span>
              <span class="stat-box-num">{{ formatAmount(projection?.projectedTotalSpent ?? 0) }} &euro;</span>
            </div>
            <div class="proj-stat-box">
              <span class="stat-box-label">&Oslash; Ausgaben / Tag</span>
              <span class="stat-box-num">{{ formatAmount(projection?.avgDailySpend ?? 0) }} &euro;</span>
            </div>
          </div>

          <div class="proj-tip-box">
            <span class="tip-icon">💡</span>
            <span class="tip-text">{{ getProjectionTip() }}</span>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useHaptics } from '../composables/useHaptics'
import type { StreakInfo, ProjectionInfo } from '../composables/useApi'

const props = withDefaults(
  defineProps<{
    isOffline?: boolean
    pendingSyncCount?: number
    streak?: StreakInfo
    projection?: ProjectionInfo
    day?: number
    monthDays?: number
    baseBudget?: number
  }>(),
  {
    isOffline: false,
    pendingSyncCount: 0,
  }
)

defineEmits<{
  (e: 'open-settings'): void
}>()

const haptics = useHaptics()
const activePopover = ref<'status' | 'streak' | 'projection' | null>(null)
const apiHealthy = ref(true)
const dbStatus = ref('connected')
const latencyMs = ref(12)

const isDbConnected = computed(() => {
  return dbStatus.value === 'connected' || dbStatus.value === 'ok'
})

const badgeClass = computed(() => {
  if (props.isOffline) return 'status-offline'
  if (!apiHealthy.value) return 'status-offline'
  if (!isDbConnected.value) return 'status-degraded'
  if (props.pendingSyncCount > 0) return 'status-syncing'
  return 'status-online'
})

const dotClass = computed(() => {
  if (props.isOffline) return 'dot-offline'
  if (!apiHealthy.value) return 'dot-offline'
  if (!isDbConnected.value) return 'dot-degraded'
  if (props.pendingSyncCount > 0) return 'dot-syncing'
  return 'dot-online'
})

const badgeText = computed(() => {
  if (props.isOffline) return 'Offline'
  if (!apiHealthy.value) return 'Server getrennt'
  if (!isDbConnected.value) return 'DB getrennt'
  if (props.pendingSyncCount > 0) return `${props.pendingSyncCount} ungesynct`
  return 'Online'
})

const badgeTooltip = computed(() => {
  if (props.isOffline) return 'Offline - Daten werden lokal gespeichert. Klicke für Details.'
  if (!apiHealthy.value) return 'API-Server nicht erreichbar. Klicke für Details.'
  if (!isDbConnected.value) return 'Datenbank-Verbindung unterbrochen. Klicke für Details.'
  if (props.pendingSyncCount > 0) return 'Ausstehende Synchronisierung. Klicke für Details.'
  return 'Online - Alle Systeme einsatzbereit. Klicke für Status-Details.'
})

const isProjectedSaving = computed(() => {
  return (props.projection?.projectedSavings ?? 0) >= 0
})

function formatAmount(val: number): string {
  return val.toLocaleString('de-DE', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

function getStreakMotivation(current: number): string {
  if (current >= 7) return '🔥 Unglaublich stark! Über eine Woche diszipliniert im Sparplan.'
  if (current >= 3) return '💪 Starker Lauf! Halte die Serie am Laufen.'
  if (current >= 1) return '🌱 Guter Start in den Tag. Bleib dran!'
  return 'Starte heute deinen neuen Spar-Streak!'
}

function getProjectionTip(): string {
  const savings = props.projection?.projectedSavings ?? 0
  const avg = props.projection?.avgDailySpend ?? 0
  const base = props.baseBudget ?? 15

  if (savings > 50) {
    return `Exzellent! Du gibst im Schnitt nur ${formatAmount(avg)} €/Tag aus (${formatAmount(base)} € Basis). Am Periodenende winkt ein toller Puffer.`
  }
  if (savings >= 0) {
    return `Gut im Kurs! Dein Schnitt von ${formatAmount(avg)} €/Tag passt perfekt zu deinem Basisbudget von ${formatAmount(base)} €/Tag.`
  }
  return `Vorsicht: Bei aktuellem Schnitt (${formatAmount(avg)} €/Tag) droht ein Defizit. Ein paar Null-Euro-Tage gleichen das schnell aus!`
}

const BASE = typeof window !== 'undefined' && import.meta.env.PROD ? window.location.origin : ''

async function checkHealth() {
  if (typeof window === 'undefined' || (import.meta as any).env?.MODE === 'test') {
    apiHealthy.value = true
    dbStatus.value = 'connected'
    latencyMs.value = 8
    return
  }

  const start = performance.now()
  try {
    const res = await fetch(`${BASE}/api/health`)
    latencyMs.value = Math.max(1, Math.round(performance.now() - start))
    if (res.ok) {
      const data = await res.json().catch(() => ({}))
      apiHealthy.value = true
      dbStatus.value = data.db || 'connected'
    } else {
      apiHealthy.value = false
      dbStatus.value = 'disconnected'
    }
  } catch {
    latencyMs.value = Math.max(1, Math.round(performance.now() - start))
    if (props.isOffline) {
      apiHealthy.value = false
      dbStatus.value = 'offline'
    } else {
      apiHealthy.value = false
      dbStatus.value = 'disconnected'
    }
  }
}

function toggleStatusPopover() {
  haptics.tap()
  activePopover.value = activePopover.value === 'status' ? null : 'status'
  if (activePopover.value === 'status') {
    checkHealth()
  }
}

function toggleStreakPopover() {
  haptics.tap()
  activePopover.value = activePopover.value === 'streak' ? null : 'streak'
}

function toggleProjectionPopover() {
  haptics.tap()
  activePopover.value = activePopover.value === 'projection' ? null : 'projection'
}

onMounted(() => {
  if (!props.isOffline) {
    checkHealth()
  }
})
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px 6px;
  position: relative;
  z-index: 30;
}

.brand {
  display: flex;
  align-items: baseline;
  user-select: none;
}

.brand-title {
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.5px;
  color: var(--text-main, #f4f4f6);
  text-transform: lowercase;
}

.brand-dot {
  font-size: 1.45rem;
  font-weight: 900;
  color: var(--accent-green, #22c55e);
  line-height: 0;
  margin-left: 1px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--text-muted, #8e8e9c);
  background: rgba(255, 255, 255, 0.04);
  padding: 4px 8px;
  border-radius: 9999px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  cursor: pointer;
  transition: all 0.15s ease;
}

.status-badge:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.dot-online {
  background-color: var(--accent-green, #22c55e);
  box-shadow: 0 0 6px rgba(34, 197, 94, 0.6);
}

.dot-offline {
  background-color: var(--accent-red, #ef4444);
  box-shadow: 0 0 6px rgba(239, 68, 68, 0.6);
  animation: pulse-dot 1.5s infinite;
}

.dot-degraded {
  background-color: #f59e0b;
  box-shadow: 0 0 6px rgba(245, 158, 11, 0.6);
  animation: pulse-dot 2s infinite;
}

.dot-syncing {
  background-color: #f59e0b;
  box-shadow: 0 0 6px rgba(245, 158, 11, 0.6);
}

.status-online { color: var(--text-muted, #8e8e9c); }
.status-offline {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  border-color: rgba(239, 68, 68, 0.25);
  color: var(--accent-red, #ef4444);
}
.status-degraded {
  background: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.25);
  color: #f59e0b;
}
.status-syncing {
  background: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.25);
  color: #f59e0b;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(0.8); }
}

.header-icon-btn {
  display: flex;
  align-items: center;
  gap: 3px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  padding: 5px 8px;
  border-radius: 9999px;
  cursor: pointer;
  font-size: 0.76rem;
  transition: all 0.15s ease;
}

.header-icon-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
}

.icon-emoji {
  font-size: 0.85rem;
  line-height: 1;
}

.btn-badge {
  font-family: var(--font-mono, monospace);
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
}

.projection-btn {
  padding: 5px 7px;
}

.settings-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  cursor: pointer;
  padding: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.settings-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--border-focus, rgba(255, 255, 255, 0.2));
  transform: rotate(45deg);
}

/* Popover Modals */
.popover-backdrop {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: 52px 14px 0;
}

.header-popover {
  background: #18181e;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 16px;
  width: 310px;
  padding: 14px;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.85);
  display: flex;
  flex-direction: column;
  gap: 10px;
  animation: popover-slide 0.15s ease-out;
}

@keyframes popover-slide {
  from { opacity: 0; transform: translateY(-8px); }
  to { opacity: 1; transform: translateY(0); }
}

.popover-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  padding-bottom: 8px;
}

.popover-title {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
}

.popover-close {
  background: transparent;
  border: none;
  color: var(--text-dim, #5c5c6e);
  font-size: 1.1rem;
  cursor: pointer;
  line-height: 1;
}

/* Health status list */
.status-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.72rem;
}

.status-row-label { color: var(--text-muted, #8e8e9c); }
.status-row-val {
  font-family: var(--font-mono, monospace);
  font-weight: 600;
}

.val-ok { color: var(--accent-green, #22c55e); }
.val-warn { color: #f59e0b; }
.val-error { color: var(--accent-red, #ef4444); }

.popover-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  padding-top: 8px;
  font-size: 0.68rem;
}

.uptime-hint { color: var(--text-dim, #5c5c6e); }
.recheck-btn {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 0.68rem;
  cursor: pointer;
}

/* Streak Popover details */
.streak-popover-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.streak-hero-stat {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: 10px;
  padding: 10px;
}

.streak-flame-icon {
  font-size: 1.8rem;
}

.streak-hero-text {
  display: flex;
  flex-direction: column;
}

.streak-big-val {
  font-size: 1.1rem;
  font-weight: 800;
  color: #f59e0b;
  font-family: var(--font-mono, monospace);
}

.streak-sub-val {
  font-size: 0.7rem;
  color: var(--text-dim, #5c5c6e);
}

.streak-stats-grid,
.proj-stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.streak-stat-box,
.proj-stat-box {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-box-label {
  font-size: 0.65rem;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
}

.stat-box-num {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-main, #f4f4f6);
  font-family: var(--font-mono, monospace);
}

.streak-motivation-text {
  font-size: 0.74rem;
  color: var(--text-muted, #8e8e9c);
  margin: 0;
  line-height: 1.35;
}

/* Projection Popover details */
.projection-popover-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.proj-hero-card {
  padding: 10px;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 2px;
}

.card-saving {
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.1));
  border: 1px solid rgba(34, 197, 94, 0.25);
}

.card-deficit {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.1));
  border: 1px solid rgba(239, 68, 68, 0.25);
}

.proj-hero-label {
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--text-dim, #5c5c6e);
  text-transform: uppercase;
}

.proj-hero-val {
  font-size: 1.35rem;
  font-weight: 800;
  font-family: var(--font-mono, monospace);
}

.card-saving .proj-hero-val { color: var(--accent-green, #22c55e); }
.card-deficit .proj-hero-val { color: var(--accent-red, #ef4444); }

.proj-hero-sub {
  font-size: 0.7rem;
  color: var(--text-dim, #5c5c6e);
}

.proj-tip-box {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 8px 10px;
  display: flex;
  align-items: flex-start;
  gap: 6px;
}

.tip-icon {
  font-size: 0.85rem;
  line-height: 1.2;
}

.tip-text {
  font-size: 0.72rem;
  color: var(--text-muted, #8e8e9c);
  line-height: 1.35;
}
</style>
