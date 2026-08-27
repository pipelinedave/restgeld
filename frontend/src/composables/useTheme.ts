import { ref, watch } from 'vue'

export interface ThemePreset {
  id: string
  name: string
  accent: string
  subtle: string
}

export const THEME_PRESETS: ThemePreset[] = [
  {
    id: 'emerald',
    name: 'Emerald Green',
    accent: '#22c55e',
    subtle: 'rgba(34, 197, 94, 0.12)',
  },
  {
    id: 'cyan',
    name: 'Cyan Mint',
    accent: '#06b6d4',
    subtle: 'rgba(6, 182, 212, 0.12)',
  },
  {
    id: 'amber',
    name: 'Sunset Amber',
    accent: '#f59e0b',
    subtle: 'rgba(245, 158, 11, 0.12)',
  },
  {
    id: 'violet',
    name: 'Neon Synthwave',
    accent: '#a855f7',
    subtle: 'rgba(168, 85, 247, 0.12)',
  },
  {
    id: 'rose',
    name: 'Ruby Rose',
    accent: '#f43f5e',
    subtle: 'rgba(244, 63, 94, 0.12)',
  },
]

const STORAGE_KEY = 'restgeld_custom_theme'

export function hexToRgba(hex: string, alpha: number): string {
  let c = hex.replace('#', '')
  if (c.length === 3) {
    c = c.split('').map((char) => char + char).join('')
  }
  const num = parseInt(c, 16)
  const r = (num >> 16) & 255
  const g = (num >> 8) & 255
  const b = num & 255
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

const currentAccent = ref<string>('#22c55e')

export function useTheme() {
  function applyTheme(color: string) {
    currentAccent.value = color
    if (typeof document !== 'undefined') {
      const root = document.documentElement
      root.style.setProperty('--accent-green', color)
      root.style.setProperty('--accent-green-subtle', hexToRgba(color, 0.12))
      root.style.setProperty('--accent', color)
    }
    try {
      localStorage.setItem(STORAGE_KEY, color)
    } catch {
      // Ignore storage errors
    }
  }

  function initTheme() {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        applyTheme(saved)
        return
      }
    } catch {
      // Ignore storage errors
    }
    applyTheme('#22c55e')
  }

  return {
    currentAccent,
    presets: THEME_PRESETS,
    applyTheme,
    initTheme,
  }
}
