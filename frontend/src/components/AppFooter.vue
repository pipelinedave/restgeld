<template>
  <footer class="app-footer">
    <div class="footer-divider"></div>
    <div class="footer-content">
      <span class="footer-tagline">Track daily. Stay in budget.</span>
      <span class="footer-version">v0.1.0</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
</script>

<style scoped>
.app-footer {
  flex-shrink: 0;
  padding: 8px 16px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.footer-divider {
  height: 1px;
  background: var(--border-color, rgba(255, 255, 255, 0.06));
  width: 100%;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.72rem;
  color: var(--text-dim, #5c5c6e);
}

.footer-tagline {
  letter-spacing: 0.02em;
}

.footer-version {
  font-family: var(--font-mono, monospace);
  font-size: 0.68rem;
}
</style>
