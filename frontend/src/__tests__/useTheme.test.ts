import { describe, it, expect, beforeEach } from 'vitest'
import { useTheme, hexToRgba, THEME_PRESETS, THEME_SCHEMES } from '../composables/useTheme'

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('konvertiert Hex-Farben korrekt in RGBA', () => {
    expect(hexToRgba('#22c55e', 0.12)).toBe('rgba(34, 197, 94, 0.12)')
    expect(hexToRgba('#fff', 0.5)).toBe('rgba(255, 255, 255, 0.5)')
  })

  it('initialisiert Standard-Akzentfarbe Emerald Green', () => {
    const theme = useTheme()
    theme.initTheme()
    expect(theme.currentAccent.value).toBe('#22c55e')
    expect(theme.presets.length).toBeGreaterThan(3)
  })

  it('wendet Farbaenderung an und speichert in localStorage', () => {
    const theme = useTheme()
    theme.applyTheme('#06b6d4')
    expect(theme.currentAccent.value).toBe('#06b6d4')
    expect(localStorage.getItem('restgeld_custom_theme')).toBe('#06b6d4')
  })

  it('bietet mehrere vollstaendige Farbschemata an', () => {
    expect(THEME_SCHEMES.length).toBeGreaterThanOrEqual(3)
    for (const s of THEME_SCHEMES) {
      expect(s.id).toBeTruthy()
      expect(s.bg).toBeTruthy()
      expect(s.accent).toBeTruthy()
    }
  })

  it('initialisiert das Standard-Schema OLED Black', () => {
    const theme = useTheme()
    theme.initTheme()
    expect(theme.currentScheme.value.id).toBe('oled')
    expect(theme.currentScheme.value.bg).toBe('#0a0a0c')
  })

  it('wechslet das Schema und setzt CSS-Variablen', () => {
    const theme = useTheme()
    theme.applyScheme('cyberpunk')
    expect(theme.currentScheme.value.id).toBe('cyberpunk')
    expect(localStorage.getItem('restgeld_theme_scheme')).toBe('cyberpunk')
  })

  it('behaelt Custom-Akzent beim Schema-Wechsel', () => {
    const theme = useTheme()
    theme.applyTheme('#ff0000')
    theme.applyScheme('sunset')
    expect(theme.currentAccent.value).toBe('#ff0000')
    expect(theme.currentScheme.value.id).toBe('sunset')
  })

  it('faellt bei unbekanntem Schema auf Standard zurueck', () => {
    const theme = useTheme()
    theme.applyScheme('does-not-exist')
    expect(theme.currentScheme.value.id).toBe('oled')
  })
})
