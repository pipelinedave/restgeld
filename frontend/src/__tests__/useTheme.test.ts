import { describe, it, expect, beforeEach } from 'vitest'
import { useTheme, hexToRgba, THEME_PRESETS } from '../composables/useTheme'

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
})
