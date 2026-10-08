import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ScrollToTop from '../ScrollToTop.vue'

describe('ScrollToTop', () => {
    beforeEach(() => {
        Object.defineProperty(window, 'scrollY', { value: 0, writable: true, configurable: true })
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    const dispatchScroll = async (y: number) => {
        (window as any).scrollY = y
        window.dispatchEvent(new Event('scroll'))
        await new Promise(r => requestAnimationFrame(r))
        await nextTick()
    }

    it('is hidden when scrollY <= 300', async () => {
        const wrapper = mount(ScrollToTop)
        expect(wrapper.find('button').exists()).toBe(false)
    })

    it('shows button when scrollY > 300', async () => {
        const wrapper = mount(ScrollToTop)
        await dispatchScroll(500)
        expect(wrapper.find('button').exists()).toBe(true)
    })

    it('hides button when scrolling back up', async () => {
        const wrapper = mount(ScrollToTop)
        await dispatchScroll(500)
        expect(wrapper.find('button').exists()).toBe(true)

        await dispatchScroll(100)
        expect(wrapper.find('button').exists()).toBe(false)
    })

    it('scrolls to top on click', async () => {
        const wrapper = mount(ScrollToTop)
        await dispatchScroll(500)

        const spy = vi.spyOn(window, 'scrollTo')
        await wrapper.find('button').trigger('click')

        expect(spy).toHaveBeenCalledWith({ top: 0, behavior: 'smooth' })
    })

    it('removes scroll listener on unmount', async () => {
        const wrapper = mount(ScrollToTop)
        const removeSpy = vi.spyOn(window, 'removeEventListener')
        wrapper.unmount()
        expect(removeSpy).toHaveBeenCalledWith('scroll', expect.any(Function))
    })
})