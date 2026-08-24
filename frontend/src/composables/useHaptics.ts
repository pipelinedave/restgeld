/**
 * Web Vibration API helper for mobile haptic feedback.
 * Safely falls back to no-op on unsupported browsers / iOS.
 */
export function useHaptics() {
  const isSupported = typeof window !== 'undefined' && 'vibrate' in navigator

  /**
   * Subtle light tap for button presses and interactions
   */
  function tap() {
    if (isSupported) {
      try {
        navigator.vibrate(15)
      } catch {
        // Ignore browser vibration errors
      }
    }
  }

  /**
   * Crisp double-pulse for successful operations (e.g. booked expense, saved budget)
   */
  function success() {
    if (isSupported) {
      try {
        navigator.vibrate([30, 40, 30])
      } catch {
        // Ignore browser vibration errors
      }
    }
  }

  /**
   * Warning vibration pulse for destructive actions / confirmation
   */
  function warning() {
    if (isSupported) {
      try {
        navigator.vibrate([40, 60, 40])
      } catch {
        // Ignore browser vibration errors
      }
    }
  }

  /**
   * Warning / error vibration pulse
   */
  function error() {
    if (isSupported) {
      try {
        navigator.vibrate([50, 70, 50])
      } catch {
        // Ignore browser vibration errors
      }
    }
  }

  return {
    isSupported,
    tap,
    success,
    warning,
    error,
  }
}
