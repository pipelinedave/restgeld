import { describe, it, expect, vi, beforeEach } from 'vitest'
import { usePushNotifications, urlBase64ToUint8Array } from '../composables/usePushNotifications'

describe('urlBase64ToUint8Array', () => {
  it('konvertiert eine Base64-URL-Key in ein Uint8Array', () => {
    // 'AB' im Raw-URL-Base64 = 0x00 0x01
    const arr = urlBase64ToUint8Array('AAE')
    expect(Array.from(arr)).toEqual([0, 1])
  })
})

describe('usePushNotifications', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('meldet nicht unterstuetzt ohne PushManager', () => {
    // jsdom mockt standardmaessig kein PushManager -> supported=false
    const push = usePushNotifications()
    expect(push.supported.value).toBe(false)
    expect(push.enabled.value).toBe(false)
  })

  it('enable() schlaegt fehl wenn nicht unterstuetzt', async () => {
    const push = usePushNotifications()
    const ok = await push.enable()
    expect(ok).toBe(false)
    expect(push.error.value).toBeTruthy()
  })
})
