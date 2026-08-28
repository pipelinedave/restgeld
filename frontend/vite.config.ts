import { defineConfig, type Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'
import { readFileSync, writeFileSync } from 'fs'
import { join } from 'path'

let gitCommit = process.env.VITE_GIT_COMMIT || process.env.VERCEL_GIT_COMMIT_SHA?.slice(0, 7) || 'dev'
if (gitCommit === 'dev') {
  try {
    gitCommit = execSync('git rev-parse --short HEAD').toString().trim()
  } catch {
    // fallback if git not available
  }
}

// Injiziert eine Build-Version oben in den Service Worker, damit der Browser
// bei jedem Deployment einen "neuen" SW erkennt (Byte-Differenz) und so alte
// Caches zuverlaessig verwirft - verhindert veraltete App nach Deploy.
function injectSwVersion(): Plugin {
  const version = `${gitCommit}-${Date.now()}`
  return {
    name: 'inject-sw-version',
    apply: 'build',
    closeBundle() {
      const out = join(__dirname, 'dist', 'sw.js')
      try {
        const sw = readFileSync(out, 'utf-8')
        if (!sw.startsWith('// SW_VERSION')) {
          writeFileSync(out, `// SW_VERSION: ${version}\n` + sw)
        }
      } catch {
        // sw.js fehlt (z. B. dev) - nichts tun
      }
    },
  }
}

export default defineConfig({
  plugins: [vue(), injectSwVersion()],
  define: {
    __GIT_COMMIT__: JSON.stringify(gitCommit),
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': process.env.VITE_API_PROXY || 'http://localhost:8080'
    },
    watch: {
      usePolling: true
    }
  }
})

