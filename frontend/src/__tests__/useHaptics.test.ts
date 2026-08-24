import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useHaptics } from '../composables/useHaptics'

describe('useHaptics composable', () => {
  const originalNavigator = global.navigator

  beforeEach(() => {
    vi.restoreAllMocks()
  })

  afterEach(() => {
    Object.defineProperty(global, 'navigator', {
      value: originalNavigator,
      configurable: true,
      writable: true,
    })
  })

  it('erkennt Unterstuetzung der Vibration API', () => {
    const vibrateMock = vi.fn()
    Object.defineProperty(global, 'navigator', {
      value: { vibrate: vibrateMock },
      configurable: true,
      writable: true,
    })

    const { isSupported, tap, success, warning, error } = useHaptics()
    expect(isSupported).toBe(true)

    tap()
    expect(vibrateMock).toHaveBeenCalledWith(15)

    success()
    expect(vibrateMock).toHaveBeenCalledWith([30, 40, 30])

    warning()
    expect(vibrateMock).toHaveBeenCalledWith([40, 60, 40])

    error()
    expect(vibrateMock).toHaveBeenCalledWith([50, 70, 50])
  })

  it('funktioniert fehlerfrei auch wenn Vibration API nicht existiert', () => {
    Object.defineProperty(global, 'navigator', {
      value: {},
      configurable: true,
      writable: true,
    })

    const { isSupported, tap, success, warning, error } = useHaptics()
    expect(isSupported).toBe(false)

    expect(() => tap()).not.toThrow()
    expect(() => success()).not.toThrow()
    expect(() => warning()).not.toThrow()
    expect(() => error()).not.toThrow()
  })
})
