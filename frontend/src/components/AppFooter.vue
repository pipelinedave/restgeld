<template>
  <footer class="app-footer">
    <div class="footer-divider"></div>
    <div class="footer-content">
      <div class="footer-left">
        <span class="footer-tagline">Track daily. Stay in budget.</span>
        <span class="footer-energy">
          ⚡ made with high energy by
          <a
            href="https://lowlifehigh.tech"
            target="_blank"
            rel="noopener noreferrer"
            class="lowlife-link"
          >
            lowlifehigh.tech
          </a>
        </span>
      </div>

      <a
        v-if="commitHash && commitHash !== 'dev'"
        :href="`https://github.com/pipelinedave/restgeld/commit/${commitHash}`"
        target="_blank"
        rel="noopener noreferrer"
        class="commit-badge"
        title="Commit auf GitHub ansehen"
      >
        <span class="commit-icon">#</span>{{ commitHash }}
      </a>
      <span v-else class="commit-badge">dev</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
const commitHash = typeof __GIT_COMMIT__ !== 'undefined' ? __GIT_COMMIT__ : 'dev'
</script>

<style scoped>
.app-footer {
  flex-shrink: 0;
  padding: 6px 16px 10px;
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
  font-size: 0.68rem;
  color: var(--text-dim, #5c5c6e);
}

.footer-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.footer-tagline {
  letter-spacing: 0.02em;
  color: var(--text-muted, #8e8e9c);
  font-weight: 500;
}

.footer-energy {
  font-size: 0.64rem;
  color: var(--text-dim, #5c5c6e);
}

.lowlife-link {
  color: var(--accent-green, #22c55e);
  text-decoration: none;
  font-weight: 600;
}

.lowlife-link:hover {
  text-decoration: underline;
}

.commit-badge {
  font-family: var(--font-mono, monospace);
  font-size: 0.68rem;
  font-weight: 600;
  color: var(--accent-green, #22c55e);
  background: var(--accent-green-subtle, rgba(34, 197, 94, 0.1));
  border: 1px solid rgba(34, 197, 94, 0.25);
  padding: 3px 8px;
  border-radius: 6px;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 2px;
  transition: all 0.15s ease;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
}

.commit-badge:hover {
  background: rgba(34, 197, 94, 0.2);
  border-color: var(--accent-green, #22c55e);
  transform: translateY(-1px);
}

.commit-icon {
  opacity: 0.6;
}
</style>
