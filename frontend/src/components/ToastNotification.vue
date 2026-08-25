<template>
  <Transition name="toast">
    <div
      v-if="visible"
      class="toast-container"
      :class="[`toast-${type}`]"
      role="status"
      aria-live="polite"
    >
      <span class="toast-icon">
        <svg v-if="type === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
        <svg v-else-if="type === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="16" x2="12" y2="12"></line>
          <line x1="12" y1="8" x2="12.01" y2="8"></line>
        </svg>
      </span>
      <span class="toast-message">{{ message }}</span>
    </div>
  </Transition>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
  message: string
  type?: 'success' | 'error' | 'info'
}>()
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 200;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 9999px;
  font-size: 0.85rem;
  font-weight: 600;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.8);
  pointer-events: none;
  max-width: 90vw;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.toast-success {
  background: rgba(18, 18, 22, 0.95);
  color: var(--accent-green, #22c55e);
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.toast-error {
  background: rgba(18, 18, 22, 0.95);
  color: var(--accent-red, #ef4444);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.toast-info {
  background: rgba(18, 18, 22, 0.95);
  color: var(--text-main, #f4f4f6);
  border: 1px solid var(--border-focus, rgba(255, 255, 255, 0.2));
}

.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.toast-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translate(-50%, -12px) scale(0.96);
}

.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -8px) scale(0.96);
}
</style>
