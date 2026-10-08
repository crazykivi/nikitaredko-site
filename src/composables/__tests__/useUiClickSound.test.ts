import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

vi.mock('../../utils/mcSounds', () => ({
    playMcSound: vi.fn(),
}))

describe('useUiClickSound', () => {
    beforeEach(() => {
        vi.resetModules()
        vi.clearAllMocks()
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    const mountHelper = async (mode: 'auto' | 'light' | 'dark' | 'charcoal') => {
        vi.doMock('../useTheme', () => ({
            useTheme: () => ({ mode: { value: mode } }),
        }))
        const { useUiClickSound } = await import('../useUiClickSound')
        const { defineComponent, h } = await import('vue')
        const { mount } = await import('@vue/test-utils')

        const Comp = defineComponent({
            setup() {
                useUiClickSound()
                return () => h('div')
            },
        })
        return mount(Comp)
    }

    it('plays click sound on button click in charcoal mode', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('charcoal')

        const btn = document.createElement('button')
        document.body.appendChild(btn)
        btn.click()

        expect(playMcSound).toHaveBeenCalledWith('click', { volume: 0.5 })
        document.body.removeChild(btn)
        wrapper.unmount()
    })

    it('does not play sound when not in charcoal mode', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('light')

        const btn = document.createElement('button')
        document.body.appendChild(btn)
        btn.click()

        expect(playMcSound).not.toHaveBeenCalled()
        document.body.removeChild(btn)
        wrapper.unmount()
    })

    it('ignores clicks inside .prose (article content)', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('charcoal')

        const prose = document.createElement('div')
        prose.className = 'prose'
        const btn = document.createElement('button')
        prose.appendChild(btn)
        document.body.appendChild(prose)

        btn.click()
        expect(playMcSound).not.toHaveBeenCalled()
        document.body.removeChild(prose)
        wrapper.unmount()
    })

    it('ignores clicks inside header and footer', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('charcoal')

        const header = document.createElement('header')
        const btn = document.createElement('button')
        header.appendChild(btn)
        document.body.appendChild(header)

        btn.click()
        expect(playMcSound).not.toHaveBeenCalled()
        document.body.removeChild(header)
        wrapper.unmount()
    })

    it('ignores clicks on non-clickable elements', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('charcoal')

        const span = document.createElement('span')
        document.body.appendChild(span)
        span.click()

        expect(playMcSound).not.toHaveBeenCalled()
        document.body.removeChild(span)
        wrapper.unmount()
    })

    it('plays sound on <a href> clicks', async () => {
        const { playMcSound } = await import('../../utils/mcSounds')
        const wrapper = await mountHelper('charcoal')

        const link = document.createElement('a')
        link.href = '/test'
        document.body.appendChild(link)
        link.click()

        expect(playMcSound).toHaveBeenCalledWith('click', { volume: 0.5 })
        document.body.removeChild(link)
        wrapper.unmount()
    })

    it('removes listener on unmount', async () => {
        const wrapper = await mountHelper('charcoal')
        const removeSpy = vi.spyOn(document, 'removeEventListener')
        wrapper.unmount()
        expect(removeSpy).toHaveBeenCalledWith('click', expect.any(Function), true)
    })
})