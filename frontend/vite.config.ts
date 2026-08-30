import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'

let gitCommit = process.env.VITE_GIT_COMMIT || process.env.VERCEL_GIT_COMMIT_SHA?.slice(0, 7) || 'dev'
if (gitCommit === 'dev') {
  try {
    gitCommit = execSync('git rev-parse --short HEAD').toString().trim()
  } catch {
    // fallback if git not available
  }
}

export default defineConfig({
  plugins: [vue()],
  define: {
    __GIT_COMMIT__: JSON.stringify(gitCommit),
  },
  server: {
    host: '0.0.0.0',
    proxy: {
      '/api/monitoring': {
        target: process.env.VITE_MONITORING_PROXY || 'http://localhost:8083',
        changeOrigin: true,
      },
      '/api/billing': {
        target: process.env.VITE_BILLING_PROXY || 'http://localhost:8082',
        changeOrigin: true,
      },
      '/api/auth': {
        target: process.env.VITE_AUTH_PROXY || 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api': {
        target: process.env.VITE_API_PROXY || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
    watch: {
      usePolling: true
    }
  }
})

