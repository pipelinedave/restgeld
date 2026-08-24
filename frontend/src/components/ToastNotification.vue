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
  padding: 10px 18px;
  border-radius: 24px;
  font-size: 0.9rem;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(2, 12, 27, 0.45);
  pointer-events: none;
  max-width: 90vw;
  backdrop-filter: blur(8px);
}

.toast-success {
  background: rgba(16, 42, 60, 0.95);
  color: #64ffda;
  border: 1px solid rgba(100, 255, 218, 0.35);
}

.toast-error {
  background: rgba(45, 18, 25, 0.95);
  color: #ff6b6b;
  border: 1px solid rgba(255, 107, 107, 0.35);
}

.toast-info {
  background: rgba(17, 34, 64, 0.95);
  color: #ccd6f6;
  border: 1px solid rgba(204, 214, 246, 0.2);
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
