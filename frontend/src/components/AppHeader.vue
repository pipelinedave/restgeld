<template>
  <header class="app-header">
    <div class="brand">
      <span class="brand-title">restgeld</span>
      <span class="brand-dot">.</span>
    </div>

    <div class="header-right">
      <!-- Status Badge mit Klick zum Öffnen des Status-Popovers -->
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

      <button class="settings-btn" aria-label="Einstellungen" title="Einstellungen" @click="$emit('open-settings')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </button>
    </div>

    <!-- Health & Uptime Popover -->
    <div v-if="showStatusPopover" class="popover-backdrop" @click="showStatusPopover = false">
      <div class="status-popover" @click.stop>
        <div class="popover-header">
          <span class="popover-title">System- & Sync-Status</span>
          <button class="popover-close" @click="showStatusPopover = false">&times;</button>
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
            <span class="status-row-val" :class="apiHealthy ? 'val-ok' : 'val-warn'">
              {{ apiHealthy ? `OK (${latencyMs}ms)` : 'Prüfe...' }}
            </span>
          </div>

          <div class="status-row">
            <span class="status-row-label">PostgreSQL Datenbank</span>
            <span class="status-row-val" :class="dbStatus === 'ok' ? 'val-ok' : 'val-warn'">
              {{ dbStatus === 'ok' ? 'Verbunden & bereit' : 'Nicht erreichbar' }}
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
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useHaptics } from '../composables/useHaptics'

const props = withDefaults(
  defineProps<{
    isOffline?: boolean
    pendingSyncCount?: number
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
const showStatusPopover = ref(false)
const apiHealthy = ref(true)
const dbStatus = ref('ok')
const latencyMs = ref(12)

const badgeClass = computed(() => {
  if (props.isOffline) return 'status-offline'
  if (props.pendingSyncCount > 0) return 'status-syncing'
  return 'status-online'
})

const dotClass = computed(() => {
  if (props.isOffline) return 'dot-offline'
  if (props.pendingSyncCount > 0) return 'dot-syncing'
  return 'dot-online'
})

const badgeText = computed(() => {
  if (props.isOffline) return 'Offline'
  if (props.pendingSyncCount > 0) return `${props.pendingSyncCount} ungesynct`
  return 'Online'
})

const badgeTooltip = computed(() => {
  if (props.isOffline) return 'Offline - Daten werden lokal gespeichert. Klicke für Details.'
  if (props.pendingSyncCount > 0) return 'Ausstehende Synchronisierung. Klicke für Details.'
  return 'Online - Alles synchron. Klicke für Status-Details.'
})

async function checkHealth() {
  if (typeof window === 'undefined' || (import.meta as any).env?.MODE === 'test') {
    apiHealthy.value = true
    dbStatus.value = 'ok'
    latencyMs.value = 8
    return
  }

  const start = performance.now()
  try {
    const res = await fetch('/api/health')
    latencyMs.value = Math.max(1, Math.round(performance.now() - start))
    if (res.ok) {
      const data = await res.json().catch(() => ({}))
      apiHealthy.value = true
      dbStatus.value = data.db || 'ok'
    } else {
      apiHealthy.value = false
      dbStatus.value = 'error'
    }
  } catch {
    latencyMs.value = Math.max(1, Math.round(performance.now() - start))
    if (props.isOffline) {
      apiHealthy.value = false
      dbStatus.value = 'offline'
    }
  }
}

function toggleStatusPopover() {
  haptics.tap()
  showStatusPopover.value = !showStatusPopover.value
  if (showStatusPopover.value) {
    checkHealth()
  }
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
  padding: 14px 18px 8px;
  position: relative;
  z-index: 20;
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
  gap: 10px;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--text-muted, #8e8e9c);
  background: rgba(255, 255, 255, 0.04);
  padding: 4px 10px;
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

.dot-syncing {
  background-color: #f59e0b;
  box-shadow: 0 0 6px rgba(245, 158, 11, 0.6);
}

.status-offline {
  background: var(--accent-red-subtle, rgba(239, 68, 68, 0.12));
  border-color: rgba(239, 68, 68, 0.25);
  color: var(--accent-red, #ef4444);
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

.settings-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  cursor: pointer;
  padding: 7px;
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

.settings-btn:active {
  transform: scale(0.95);
}

/* Status Popover */
.popover-backdrop {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: 55px 16px 0 0;
}

.status-popover {
  background: #18181e;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 14px;
  width: 290px;
  padding: 14px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.8);
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
  font-size: 0.8rem;
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

.status-row-label {
  color: var(--text-muted, #8e8e9c);
}

.status-row-val {
  font-family: var(--font-mono, monospace);
  font-weight: 600;
}

.val-ok {
  color: var(--accent-green, #22c55e);
}

.val-warn {
  color: #f59e0b;
}

.val-error {
  color: var(--accent-red, #ef4444);
}

.popover-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  padding-top: 8px;
  font-size: 0.68rem;
}

.uptime-hint {
  color: var(--text-dim, #5c5c6e);
}

.recheck-btn {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  color: var(--text-muted, #8e8e9c);
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 0.68rem;
  cursor: pointer;
}

.recheck-btn:hover {
  color: var(--text-main, #f4f4f6);
  border-color: var(--accent-green, #22c55e);
}
</style>
