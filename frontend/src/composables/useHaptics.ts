import { ref } from 'vue'

const SOUND_KEY = 'restgeld_sound_enabled'
const soundEnabled = ref<boolean>(
  typeof window !== 'undefined' ? localStorage.getItem(SOUND_KEY) === 'true' : false
)

let audioCtx: AudioContext | null = null

function getAudioContext(): AudioContext | null {
  if (typeof window === 'undefined') return null
  if (!audioCtx) {
    const AudioContextClass = window.AudioContext || (window as any).webkitAudioContext
    if (AudioContextClass) {
      audioCtx = new AudioContextClass()
    }
  }
  if (audioCtx && audioCtx.state === 'suspended') {
    audioCtx.resume().catch(() => {})
  }
  return audioCtx
}

function playTone(freq: number, durationSec: number, type: OscillatorType = 'sine', gainVal: number = 0.05) {
  if (!soundEnabled.value) return
  try {
    const ctx = getAudioContext()
    if (!ctx) return
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()

    osc.type = type
    osc.frequency.setValueAtTime(freq, ctx.currentTime)

    gain.gain.setValueAtTime(gainVal, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.0001, ctx.currentTime + durationSec)

    osc.connect(gain)
    gain.connect(ctx.destination)

    osc.start()
    osc.stop(ctx.currentTime + durationSec)
  } catch {
    // Ignore audio errors
  }
}

/**
 * Web Vibration & Micro-Audio API helper for haptic feedback.
 */
export function useHaptics() {
  const isSupported = typeof window !== 'undefined' && 'vibrate' in navigator

  function setSoundEnabled(val: boolean) {
    soundEnabled.value = val
    try {
      localStorage.setItem(SOUND_KEY, val ? 'true' : 'false')
    } catch {}
  }

  function tap() {
    if (isSupported) {
      try {
        navigator.vibrate(15)
      } catch {}
    }
    playTone(600, 0.02, 'sine', 0.03)
  }

  function success() {
    if (isSupported) {
      try {
        navigator.vibrate([30, 40, 30])
      } catch {}
    }
    if (soundEnabled.value) {
      playTone(523.25, 0.08, 'sine', 0.05)
      setTimeout(() => playTone(659.25, 0.08, 'sine', 0.05), 60)
      setTimeout(() => playTone(783.99, 0.12, 'sine', 0.06), 120)
    }
  }

  function warning() {
    if (isSupported) {
      try {
        navigator.vibrate([40, 60, 40])
      } catch {}
    }
    if (soundEnabled.value) {
      playTone(440, 0.1, 'triangle', 0.05)
      setTimeout(() => playTone(349.23, 0.12, 'triangle', 0.05), 80)
    }
  }

  function error() {
    if (isSupported) {
      try {
        navigator.vibrate([50, 70, 50])
      } catch {}
    }
    if (soundEnabled.value) {
      playTone(220, 0.15, 'sawtooth', 0.04)
    }
  }

  return {
    isSupported,
    soundEnabled,
    setSoundEnabled,
    tap,
    success,
    warning,
    error,
  }
}
