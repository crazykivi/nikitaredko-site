import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchFab from '../SearchFab.vue'

describe('SearchFab', () => {
    afterEach(() => vi.restoreAllMocks())

    it('renders a button with search icon', () => {
        const wrapper = mount(SearchFab)
        expect(wrapper.find('button').exists()).toBe(true)
        expect(wrapper.find('svg').exists()).toBe(true)
    })

    it('has sm:hidden class (mobile-only FAB)', () => {
        const wrapper = mount(SearchFab)
        expect(wrapper.find('button').classes()).toContain('sm:hidden')
    })

    it('dispatches "/" keydown on click', async () => {
        const spy = vi.spyOn(window, 'dispatchEvent')
        const wrapper = mount(SearchFab)

        await wrapper.find('button').trigger('click')

        expect(spy).toHaveBeenCalledTimes(1)
        const event = spy.mock.calls[0][0] as KeyboardEvent
        expect(event.key).toBe('/')
        expect(event.bubbles).toBe(true)
    })
})