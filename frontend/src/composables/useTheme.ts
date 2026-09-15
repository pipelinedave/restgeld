import { ref } from 'vue'

export interface ThemePreset {
  id: string
  name: string
  accent: string
  subtle: string
}

export interface ThemeScheme {
  id: string
  name: string
  /** Basis-Hintergrundfarbe (Page) */
  bg: string
  /** Karten-Hintergrund */
  bgCard: string
  /** Subtiler Karten-/Hover-Hintergrund */
  bgSubtle: string
  /** Erhöhter/Unerhöhte Panel-Hintergrund */
  bgElevated: string
  /** Eingabe-/Input-Hintergrund */
  bgInput: string
  /** Schriftfarbe */
  text: string
  /** Sekundäre Schriftfarbe */
  textMuted: string
  /** Rahmenfarbe */
  border: string
  /** Akzentfarbe */
  accent: string
  /** Glow der Akzentfarbe */
  accentGlow: string
  /** Subtiler Akzent-Hintergrund */
  accentSubtle: string
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

/**
 * Vollständige Farbschemata (Epic 3.1): Ändern nicht nur den Akzent,
 * sondern die gesamte Ambient-Palette (Hintergrund, Karten, Rahmen).
 */
export const THEME_SCHEMES: ThemeScheme[] = [
  {
    id: 'oled',
    name: 'OLED Black',
    bg: '#0a0a0c',
    bgCard: '#121216',
    bgSubtle: '#1c1c24',
    bgElevated: '#242430',
    bgInput: '#1f1f28',
    text: '#f4f4f6',
    textMuted: '#8e8e9c',
    border: 'rgba(255, 255, 255, 0.08)',
    accent: '#22c55e',
    accentGlow: 'rgba(34, 197, 94, 0.25)',
    accentSubtle: 'rgba(34, 197, 94, 0.12)',
  },
  {
    id: 'cyberpunk',
    name: 'Cyberpunk',
    bg: '#0a0a16',
    bgCard: '#12121f',
    bgSubtle: '#1b1b2e',
    bgElevated: '#23233a',
    bgInput: '#1e1e30',
    text: '#eaf6ff',
    textMuted: '#9aa7c7',
    border: 'rgba(0, 229, 255, 0.18)',
    accent: '#00e5ff',
    accentGlow: 'rgba(0, 229, 255, 0.28)',
    accentSubtle: 'rgba(0, 229, 255, 0.12)',
  },
  {
    id: 'sunset',
    name: 'Sunset',
    bg: '#140d0a',
    bgCard: '#1c1410',
    bgSubtle: '#271d16',
    bgElevated: '#332618',
    bgInput: '#2b2117',
    text: '#fff4ec',
    textMuted: '#c9a997',
    border: 'rgba(255, 122, 51, 0.18)',
    accent: '#ff7a33',
    accentGlow: 'rgba(255, 122, 51, 0.28)',
    accentSubtle: 'rgba(255, 122, 51, 0.12)',
  },
  {
    id: 'forest',
    name: 'Forest Mint',
    bg: '#050d09',
    bgCard: '#0b1710',
    bgSubtle: '#13241a',
    bgElevated: '#1a3324',
    bgInput: '#162a1e',
    text: '#ecfdf3',
    textMuted: '#87b49a',
    border: 'rgba(52, 211, 153, 0.16)',
    accent: '#34d399',
    accentGlow: 'rgba(52, 211, 153, 0.26)',
    accentSubtle: 'rgba(52, 211, 153, 0.12)',
  },
]

const ACCENT_STORAGE_KEY = 'restgeld_custom_theme'
const SCHEME_STORAGE_KEY = 'restgeld_theme_scheme'

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

function isPrefixedColor(value: string): boolean {
  return value.startsWith('rgba(') || value.startsWith('rgb(')
}

const currentAccent = ref<string>('#22c55e')
const currentScheme = ref<ThemeScheme>(THEME_SCHEMES[0])

function applySchemeVariables(scheme: ThemeScheme, accent?: string) {
  if (typeof document !== 'undefined') {
    const root = document.documentElement
    const set = (k: string, v: string) => root.style.setProperty(k, v)

    const acc = accent ?? scheme.accent
    const accGlow = isPrefixedColor(acc) ? acc : hexToRgba(acc, 0.28)
    const accSubtle = isPrefixedColor(acc) ? acc : hexToRgba(acc, 0.12)

    set('--bg-color', scheme.bg)
    set('--bg-card', scheme.bgCard)
    set('--bg-card-hover', scheme.bgSubtle)
    set('--bg-subtle', scheme.bgSubtle)
    set('--bg-elevated', scheme.bgElevated)
    set('--bg-input', scheme.bgInput)
    set('--border-color', scheme.border)
    set('--text-main', scheme.text)
    set('--text-muted', scheme.textMuted)
    set('--accent-green', acc)
    set('--accent-green-glow', accGlow)
    set('--accent-green-subtle', accSubtle)
    set('--accent', acc)
  }
}

export function useTheme() {
  function applyTheme(color: string) {
    currentAccent.value = color
    applySchemeVariables(currentScheme.value, color)
    try {
      localStorage.setItem(ACCENT_STORAGE_KEY, color)
    } catch {
      // Ignore storage errors
    }
  }

  function applyScheme(id: string, accent?: string) {
    const scheme = THEME_SCHEMES.find((s) => s.id === id) ?? THEME_SCHEMES[0]
    currentScheme.value = scheme
    // Behalte eine evtl. gesetzte Custom-Akzentfarbe, sonst Schema-Akzent
    const acc = accent ?? currentAccent.value
    applySchemeVariables(scheme, acc)
    try {
      localStorage.setItem(SCHEME_STORAGE_KEY, scheme.id)
    } catch {
      // Ignore storage errors
    }
  }

  function initTheme() {
    let scheme = THEME_SCHEMES[0]
    let accent: string | undefined

    try {
      const savedScheme = localStorage.getItem(SCHEME_STORAGE_KEY)
      if (savedScheme) {
        const found = THEME_SCHEMES.find((s) => s.id === savedScheme)
        if (found) scheme = found
      }
      const savedAccent = localStorage.getItem(ACCENT_STORAGE_KEY)
      if (savedAccent) accent = savedAccent
    } catch {
      // Ignore storage errors
    }

    currentScheme.value = scheme
    currentAccent.value = accent ?? scheme.accent
    applySchemeVariables(scheme, accent)
  }

  return {
    currentAccent,
    currentScheme,
    presets: THEME_PRESETS,
    schemes: THEME_SCHEMES,
    applyTheme,
    applyScheme,
    initTheme,
  }
}
