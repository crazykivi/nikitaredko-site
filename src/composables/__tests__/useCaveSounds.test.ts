import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import type { ThemeMode } from '../useTheme'
import type { Ref } from 'vue'

vi.mock('../../utils/mcSounds', () => ({
    playMcSound: vi.fn(),
}))

describe('useCaveSounds', () => {
    beforeEach(() => {
        vi.resetModules()
        vi.clearAllMocks()
        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.useRealTimers()
        vi.restoreAllMocks()
    })

    const mountWithMode = async (initial: 'auto' | 'light' | 'dark' | 'charcoal') => {
        const { ref, defineComponent, h, nextTick } = await import('vue')
        const { mount } = await import('@vue/test-utils')
        const { useCaveSounds } = await import('../useCaveSounds')

        let modeRef!: Ref<ThemeMode>
        const Comp = defineComponent({
            setup() {
                modeRef = ref(initial)
                useCaveSounds(modeRef)
                return () => h('div')
            },
        })
        const wrapper = mount(Comp)
        await nextTick()
        return { modeRef, wrapper }
    }

    it('schedules cave sound when mode is charcoal', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const { wrapper } = await mountWithMode('charcoal')

        await vi.advanceTimersByTimeAsync(31_000)
        expect(playMcSound).toHaveBeenCalledWith('cave', { volume: 0.3 })
        wrapper.unmount()
    })

    it('does not play sound when mode is not charcoal', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const { wrapper } = await mountWithMode('light')

        await vi.advanceTimersByTimeAsync(200_000)
        expect(playMcSound).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('clears timer when switching away from charcoal', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const { modeRef, wrapper } = await mountWithMode('charcoal')
        const { nextTick } = await import('vue')

        modeRef.value = 'light'
        await nextTick()
        await vi.advanceTimersByTimeAsync(200_000)
        expect(playMcSound).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('schedules next sound after each play', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const { wrapper } = await mountWithMode('charcoal')

        await vi.advanceTimersByTimeAsync(31_000)
        expect(playMcSound).toHaveBeenCalledTimes(1)

        await vi.advanceTimersByTimeAsync(121_000)
        expect(playMcSound).toHaveBeenCalledTimes(2)
        wrapper.unmount()
    })
})