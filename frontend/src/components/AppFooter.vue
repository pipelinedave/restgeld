<template>
  <footer class="app-footer">
    <div class="footer-divider"></div>
    <div class="footer-content">
      <span class="footer-tagline">Track daily. Stay in budget.</span>
      <a
        v-if="commitHash && commitHash !== 'dev'"
        :href="`https://github.com/pipelinedave/restgeld/commit/${commitHash}`"
        target="_blank"
        rel="noopener noreferrer"
        class="footer-commit"
        title="Commit auf GitHub anzeigen"
      >
        <span class="commit-icon">#</span>{{ commitHash }}
      </a>
      <span v-else class="footer-commit">dev</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
const commitHash = typeof __GIT_COMMIT__ !== 'undefined' ? __GIT_COMMIT__ : 'dev'
</script>

<style scoped>
.app-footer {
  flex-shrink: 0;
  padding: 6px 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
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
  font-size: 0.7rem;
  color: var(--text-dim, #5c5c6e);
}

.footer-tagline {
  letter-spacing: 0.02em;
}

.footer-commit {
  font-family: var(--font-mono, monospace);
  font-size: 0.68rem;
  color: var(--text-dim, #5c5c6e);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 2px;
  transition: color 0.15s;
}

.footer-commit:hover {
  color: var(--accent, #22c55e);
}

.commit-icon {
  opacity: 0.6;
}
</style>
