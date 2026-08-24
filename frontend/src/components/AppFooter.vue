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
  background: rgba(136, 146, 176, 0.15);
  width: 100%;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
  color: var(--text-dim, #8892b0);
  opacity: 0.7;
}

.footer-tagline {
  letter-spacing: 0.02em;
}

.footer-version {
  font-family: monospace;
  font-size: 0.7rem;
}
</style>
