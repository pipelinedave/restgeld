import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  retries: 0,
  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      // Standard-E2E gegen den Dev-Server (Mock-API, ohne Service Worker).
      name: 'chromium',
      // Der SW-Spec braucht den Produktions-Build (gehashte Assets + echten SW),
      // daher laeuft er in einem separaten Projekt gegen `vite preview`.
      testIgnore: /service-worker\.spec\.ts/,
      use: { ...devices['Desktop Chrome'] },
    },
    {
      // Service-Worker-/PWA-Tests gegen den PRODUKTIONS-Build.
      // Der Dev-Server hat keine gehashten Assets und kein echtes Offline-Verhalten.
      name: 'sw-preview',
      testMatch: /service-worker\.spec\.ts/,
      use: {
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:4173',
        // Der echte Service Worker registriert sich nur, wenn die App den
        // Automation-Mode nicht erkennt (navigator.webdriver === false).
        launchOptions: {
          args: ['--disable-blink-features=AutomationControlled'],
        },
      },
    },
  ],
  webServer: [
    {
      command: 'npm run dev -- --port 5173',
      port: 5173,
      reuseExistingServer: !process.env.CI,
      timeout: 15000,
    },
    {
      // Produktions-Build servieren (fuer SW-Tests). Wird vorher gebaut.
      command: 'npm run build && npm run preview -- --port 4173',
      port: 4173,
      reuseExistingServer: !process.env.CI,
      timeout: 15000,
    },
  ],
})
