import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

vi.mock('../../utils/mcSounds', () => ({
    playMcSound: vi.fn(),
}))
vi.mock('../../utils/readProgress', () => ({
    markArticleRead: vi.fn(() => true),
}))

describe('useArticleXp', () => {
    beforeEach(() => {
        vi.resetModules()
        vi.clearAllMocks()
        vi.useFakeTimers()
        Object.defineProperty(window, 'scrollY', { value: 0, writable: true, configurable: true })
        Object.defineProperty(window, 'innerHeight', { value: 800, writable: true, configurable: true })
        Object.defineProperty(document.documentElement, 'scrollHeight', {
            value: 2000,
            writable: true,
            configurable: true,
        })
    })

    afterEach(() => {
        vi.useRealTimers()
        vi.restoreAllMocks()
    })

    const mountWithMode = async (mode: 'auto' | 'light' | 'dark' | 'charcoal', articleId: string | null) => {
        vi.doMock('../useTheme', () => ({
            useTheme: () => ({ mode: { value: mode } }),
        }))
        const { ref, defineComponent, h, nextTick } = await import('vue')
        const { mount } = await import('@vue/test-utils')
        const { useArticleXp } = await import('../useArticleXp')

        let idRef!: ReturnType<typeof ref<string | null>>
        const Comp = defineComponent({
            setup() {
                idRef = ref<string | null>(articleId)
                useArticleXp(idRef)
                return () => h('div')
            },
        })
        const wrapper = mount(Comp)
        await nextTick()
        return { idRef, wrapper }
    }

    const dispatchScroll = async (y: number) => {
        (window as any).scrollY = y
        window.dispatchEvent(new Event('scroll'))
        await vi.advanceTimersByTimeAsync(0)
    }

    it('marks article as read when scrolled to bottom in charcoal mode', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { playMcSound } = await import('../../utils/mcSounds')
        const { wrapper } = await mountWithMode('charcoal', 'art-1')

        await dispatchScroll(1300)

        expect(markArticleRead).toHaveBeenCalledWith('art-1')
        vi.advanceTimersByTime(300)
        expect(playMcSound).toHaveBeenCalledWith('xp', { volume: 0.1 })
        wrapper.unmount()
    })

    it('does not mark read when not in charcoal mode', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { wrapper } = await mountWithMode('light', 'art-1')

        await dispatchScroll(1300)
        expect(markArticleRead).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('does not mark read when not scrolled to bottom', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { wrapper } = await mountWithMode('charcoal', 'art-1')

        await dispatchScroll(100)
        expect(markArticleRead).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('does not mark read twice for the same article', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { wrapper } = await mountWithMode('charcoal', 'art-1')

        await dispatchScroll(1300)
        await dispatchScroll(1300)
        expect(markArticleRead).toHaveBeenCalledTimes(1)
        wrapper.unmount()
    })

    it('resets state when articleId changes', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { idRef, wrapper } = await mountWithMode('charcoal', 'art-1')

        await dispatchScroll(1300)
        expect(markArticleRead).toHaveBeenCalledTimes(1)

        idRef.value = 'art-2'
        await vi.advanceTimersByTimeAsync(1000)

        await dispatchScroll(1300)
        expect(markArticleRead).toHaveBeenCalledWith('art-2')
        wrapper.unmount()
    })

    it('does not mark read when articleId is null', async () => {
        const { markArticleRead } = await import('../../utils/readProgress')
        const { wrapper } = await mountWithMode('charcoal', null)

        await dispatchScroll(1300)
        expect(markArticleRead).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('does not mark read for short pages (docHeight <= innerHeight + 200)', async () => {
        Object.defineProperty(document.documentElement, 'scrollHeight', {
            value: 900,
            writable: true,
            configurable: true,
        })
        const { markArticleRead } = await import('../../utils/readProgress')
        const { wrapper } = await mountWithMode('charcoal', 'art-1')

        await dispatchScroll(500)
        expect(markArticleRead).not.toHaveBeenCalled()
        wrapper.unmount()
    })

    it('removes scroll listener on unmount', async () => {
        const { wrapper } = await mountWithMode('charcoal', 'art-1')
        const removeSpy = vi.spyOn(window, 'removeEventListener')
        wrapper.unmount()
        expect(removeSpy).toHaveBeenCalledWith('scroll', expect.any(Function))
    })
})