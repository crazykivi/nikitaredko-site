import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('useSoundSettings', () => {
  beforeEach(() => {
    vi.resetModules()
    localStorage.clear()
  })

  it('defaults to unmuted when localStorage is empty', async () => {
    const { useSoundSettings } = await import('../useSoundSettings')
    const { isSoundMuted } = useSoundSettings()
    expect(isSoundMuted.value).toBe(false)
  })

  it('reads initial muted state from localStorage', async () => {
    localStorage.setItem('mc-sounds-muted', 'true')
    const { useSoundSettings } = await import('../useSoundSettings')
    const { isSoundMuted } = useSoundSettings()
    expect(isSoundMuted.value).toBe(true)
  })

  it('toggleMute flips the value', async () => {
    const { useSoundSettings } = await import('../useSoundSettings')
    const { isSoundMuted, toggleMute } = useSoundSettings()

    expect(isSoundMuted.value).toBe(false)
    toggleMute()
    expect(isSoundMuted.value).toBe(true)
    toggleMute()
    expect(isSoundMuted.value).toBe(false)
  })

  it('persists mute state to localStorage via watch', async () => {
    const { useSoundSettings } = await import('../useSoundSettings')
    const { nextTick } = await import('vue')
    const { toggleMute } = useSoundSettings()

    toggleMute()
    await nextTick()
    expect(localStorage.getItem('mc-sounds-muted')).toBe('true')

    toggleMute()
    await nextTick()
    expect(localStorage.getItem('mc-sounds-muted')).toBe('false')
  })

  it('shares state across multiple useSoundSettings calls', async () => {
    const { useSoundSettings } = await import('../useSoundSettings')
    const a = useSoundSettings()
    const b = useSoundSettings()
    a.toggleMute()
    expect(b.isSoundMuted.value).toBe(true)
  })
})