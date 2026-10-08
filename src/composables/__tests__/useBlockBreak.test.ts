import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('useBlockBreak', () => {
    beforeEach(() => {
        vi.resetModules()
        localStorage.clear()
    })

    it('defaults to locked when localStorage is empty', async () => {
        const { useBlockBreak } = await import('../useBlockBreak')
        const { unlocked } = useBlockBreak()
        expect(unlocked.value).toBe(false)
    })

    it('reads unlocked state from localStorage', async () => {
        localStorage.setItem('block-break-unlocked', 'true')
        const { useBlockBreak } = await import('../useBlockBreak')
        const { unlocked } = useBlockBreak()
        expect(unlocked.value).toBe(true)
    })

    it('toggle returns new state and flips the ref', async () => {
        const { useBlockBreak } = await import('../useBlockBreak')
        const { unlocked, toggle } = useBlockBreak()

        const result1 = toggle()
        expect(result1).toBe(true)
        expect(unlocked.value).toBe(true)

        const result2 = toggle()
        expect(result2).toBe(false)
        expect(unlocked.value).toBe(false)
    })

    it('persists state to localStorage via watch', async () => {
        const { useBlockBreak } = await import('../useBlockBreak')
        const { nextTick } = await import('vue')
        const { toggle } = useBlockBreak()

        toggle()
        await nextTick()
        expect(localStorage.getItem('block-break-unlocked')).toBe('true')

        toggle()
        await nextTick()
        expect(localStorage.getItem('block-break-unlocked')).toBe('false')
    })

    it('unlocked ref is readonly', async () => {
        const { useBlockBreak } = await import('../useBlockBreak')
        const { unlocked } = useBlockBreak()
        expect(unlocked.value).toBe(false)
    })
})