import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

vi.mock('../../utils/mcSounds', () => ({
    playMcSound: vi.fn(),
}))

describe('useLadderSound', () => {
    beforeEach(() => {
        vi.resetModules()
        vi.clearAllMocks()
        vi.useFakeTimers()
        Object.defineProperty(window, 'scrollY', { value: 0, writable: true, configurable: true })
    })

    afterEach(() => {
        vi.useRealTimers()
        vi.restoreAllMocks()
    })

    const mountWithMode = async (mode: 'auto' | 'light' | 'dark' | 'charcoal') => {
        vi.doMock('../useTheme', () => ({
            useTheme: () => ({ mode: { value: mode } }),
        }))
        const { defineComponent, h } = await import('vue')
        const { mount } = await import('@vue/test-utils')
        const { useLadderSound } = await import('../useLadderSound')

        const Comp = defineComponent({
            setup() {
                useLadderSound()
                return () => h('div')
            },
        })
        return mount(Comp)
    }

    const dispatchScroll = async (newScrollY: number) => {
        (window as any).scrollY = newScrollY
        window.dispatchEvent(new Event('scroll'))
        await vi.advanceTimersByTimeAsync(0)
    }

    it('plays ladder sound after 250px scroll in charcoal mode', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountWithMode('charcoal')

        await dispatchScroll(300)
        expect(playMcSound).toHaveBeenCalledWith('ladder', expect.objectContaining({
            volume: 0.10,
        }))
        wrapper.unmount()
    })

    it('does not play sound when not in charcoal mode', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountWithMode('light')

        await dispatchScroll(500)
        expect(playMcSound).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('accumulates small scrolls before playing', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountWithMode('charcoal')

        await dispatchScroll(100)
        expect(playMcSound).not.toHaveBeenCalled()

        await dispatchScroll(200)
        expect(playMcSound).not.toHaveBeenCalled()

        await dispatchScroll(260)
        expect(playMcSound).toHaveBeenCalledTimes(1)
        wrapper.unmount()
    })

    it('respects 350ms cooldown between plays', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountWithMode('charcoal')

        await dispatchScroll(300)
        expect(playMcSound).toHaveBeenCalledTimes(1)

        await dispatchScroll(600)
        expect(playMcSound).toHaveBeenCalledTimes(1)

        vi.advanceTimersByTime(400)
        await dispatchScroll(900)
        expect(playMcSound).toHaveBeenCalledTimes(2)
        wrapper.unmount()
    })

    it('removes scroll listener on unmount', async () => {
        const wrapper = await mountWithMode('charcoal')
        const removeSpy = vi.spyOn(window, 'removeEventListener')
        wrapper.unmount()
        expect(removeSpy).toHaveBeenCalledWith('scroll', expect.any(Function))
    })
})