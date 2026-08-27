<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>ℹ️ Über Restgeld</h2>
        <button class="close-btn" aria-label="Schließen" @click="$emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <div class="about-hero">
          <div class="brand-large">
            <span class="brand-title">restgeld</span>
            <span class="brand-dot">.</span>
          </div>
          <p class="about-tagline">Track daily. Stay in budget.</p>
        </div>

        <section class="about-section">
          <h3 class="section-heading">🌱 Die Philosophie</h3>
          <p class="section-text">
            <strong>Restgeld</strong> ist kein überladenes Haushaltsbuch mit hunderten Kategorien. Es fokussiert sich auf das Wesentliche: <em>Wie viel Spaßgeld habe ich heute noch übrig?</em>
          </p>
          <p class="section-text">
            Achtsamkeit statt Buchhaltungs-Stress. Jeder Tag startet mit deinem gerechten Tagesbudget &ndash; gibst du weniger aus, wächst dein Sparpuffer; gibst du mehr aus, ziehst du es vom Puffer ab.
          </p>
        </section>

        <section class="about-section">
          <h3 class="section-heading">⚡ Key Features</h3>
          <ul class="features-list">
            <li><strong>100% Offline-First:</strong> Buche Ausgaben überall, automatischer Sync bei Wiederverbindung.</li>
            <li><strong>Zero-Bloat:</strong> 1-Tap Schnellbuchung, keine Pflichtkategorien.</li>
            <li><strong>Homescreen Shortcut:</strong> Long-Press auf das Icon öffnet sofort das Numpad.</li>
            <li><strong>OLED Dark Design:</strong> Modernes Glassmorphism mit anpassbaren Akzentfarben.</li>
          </ul>
        </section>

        <section class="about-section meta-section">
          <div class="meta-row">
            <span class="meta-label">Version</span>
            <a
              :href="`https://github.com/pipelinedave/restgeld/commit/${gitCommit}`"
              target="_blank"
              rel="noopener noreferrer"
              class="commit-link"
            >
              {{ gitCommit }} ↗
            </a>
          </div>
          <div class="meta-row">
            <span class="meta-label">Lizenz</span>
            <span class="meta-val">MIT License &bull; Open Source</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">GitHub</span>
            <a
              href="https://github.com/pipelinedave/restgeld"
              target="_blank"
              rel="noopener noreferrer"
              class="github-link"
            >
              pipelinedave/restgeld ↗
            </a>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const gitCommit = typeof __GIT_COMMIT__ !== 'undefined' ? __GIT_COMMIT__ : 'dev'
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
  animation: modal-fade 0.2s ease-out;
}

@keyframes modal-fade {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: #121216;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: 20px;
  width: 100%;
  max-width: 440px;
  max-height: 85dvh;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.9);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
}

.modal-header h2 {
  font-size: 1.15rem;
  color: var(--text-main, #f4f4f6);
  margin: 0;
  font-weight: 700;
}

.close-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  color: var(--text-muted, #8e8e9c);
  font-size: 1.25rem;
  line-height: 1;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.close-btn:hover {
  color: var(--text-main, #f4f4f6);
  background: rgba(255, 255, 255, 0.08);
}

.modal-body {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

.about-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 8px 0;
}

.brand-large {
  display: flex;
  align-items: baseline;
}

.brand-title {
  font-size: 2rem;
  font-weight: 900;
  letter-spacing: -1px;
  color: var(--text-main, #f4f4f6);
}

.brand-dot {
  font-size: 2.2rem;
  font-weight: 900;
  color: var(--accent-green, #22c55e);
  line-height: 0;
  margin-left: 2px;
}

.about-tagline {
  font-size: 0.82rem;
  color: var(--text-dim, #5c5c6e);
  margin: 4px 0 0;
  letter-spacing: 0.5px;
}

.about-section {
  background: #18181e;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.06));
  border-radius: 14px;
  padding: 14px;
}

.section-heading {
  font-size: 0.82rem;
  color: var(--accent-green, #22c55e);
  margin: 0 0 8px;
  font-weight: 700;
}

.section-text {
  font-size: 0.78rem;
  color: var(--text-muted, #8e8e9c);
  line-height: 1.45;
  margin: 0 0 6px;
}

.section-text:last-child {
  margin-bottom: 0;
}

.features-list {
  margin: 0;
  padding-left: 18px;
  font-size: 0.76rem;
  color: var(--text-muted, #8e8e9c);
  display: flex;
  flex-direction: column;
  gap: 4px;
  line-height: 1.35;
}

.meta-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
}

.meta-label {
  color: var(--text-dim, #5c5c6e);
}

.meta-val {
  color: var(--text-muted, #8e8e9c);
  font-weight: 600;
}

.github-link,
.commit-link {
  color: var(--accent-green, #22c55e);
  text-decoration: none;
  font-weight: 600;
  font-family: var(--font-mono, monospace);
}

.github-link:hover,
.commit-link:hover {
  text-decoration: underline;
}
</style>
