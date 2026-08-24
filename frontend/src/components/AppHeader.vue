<template>
  <header class="app-header">
    <div class="brand">
      <span class="brand-title">restgeld</span>
      <span class="brand-dot">.</span>
      <span v-if="isOffline" class="offline-badge" title="Offline - Daten werden lokal gespeichert">
        <span class="offline-dot"></span> Offline
      </span>
      <span v-else-if="pendingSyncCount > 0" class="syncing-badge" title="Ausstehende Synchronisierung">
        {{ pendingSyncCount }} ungesynct
      </span>
    </div>
    <button class="settings-btn" aria-label="Einstellungen" title="Einstellungen" @click="$emit('open-settings')">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
      </svg>
    </button>
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
  padding: 16px 16px 8px;
}

.brand {
  display: flex;
  align-items: baseline;
  user-select: none;
}

.brand-title {
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text, #ccd6f6);
  text-transform: lowercase;
}

.brand-dot {
  font-size: 1.5rem;
  font-weight: 900;
  color: var(--accent, #64ffda);
  line-height: 0;
  margin-left: 1px;
}

.settings-btn {
  background: transparent;
  border: none;
  color: #8892b0;
  cursor: pointer;
  padding: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: color 0.15s, background-color 0.15s, transform 0.2s ease-in-out;
}

.settings-btn:hover {
  color: #64ffda;
  background: rgba(100, 255, 218, 0.08);
  transform: rotate(45deg);
}

.offline-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: rgba(255, 107, 107, 0.15);
  border: 1px solid rgba(255, 107, 107, 0.3);
  color: #ff6b6b;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  margin-left: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.offline-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #ff6b6b;
  animation: pulse-dot 1.5s infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(0.8); }
}

.syncing-badge {
  background: rgba(100, 255, 218, 0.1);
  border: 1px solid rgba(100, 255, 218, 0.3);
  color: var(--accent, #64ffda);
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  margin-left: 10px;
}
</style>
