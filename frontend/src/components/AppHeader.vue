<template>
  <header class="app-header">
    <div class="brand">
      <span class="brand-title">restgeld</span>
      <span class="brand-dot">.</span>
    </div>

    <div class="header-right">
      <div v-if="isOffline" class="status-badge status-offline" title="Offline - Daten werden lokal gespeichert">
        <span class="status-dot dot-offline"></span>
        <span class="status-text">Offline</span>
      </div>
      <div v-else-if="pendingSyncCount > 0" class="status-badge status-syncing" title="Ausstehende Synchronisierung">
        <span class="status-dot dot-syncing"></span>
        <span class="status-text">{{ pendingSyncCount }} ungesynct</span>
      </div>
      <div v-else class="status-badge status-online" title="Online - Alles synchron">
        <span class="status-dot dot-online"></span>
        <span class="status-text">Online</span>
      </div>

      <button class="settings-btn" aria-label="Einstellungen" title="Einstellungen" @click="$emit('open-settings')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
withDefaults(
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
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 18px 10px;
  position: relative;
  z-index: 10;
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
</style>
