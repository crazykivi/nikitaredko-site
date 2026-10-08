import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Spinner from '../Spinner.vue'

describe('Spinner', () => {
    it('renders with default size md', () => {
        const wrapper = mount(Spinner)
        expect(wrapper.find('div').exists()).toBe(true)
        expect(wrapper.classes()).toContain('w-8')
        expect(wrapper.classes()).toContain('h-8')
        expect(wrapper.classes()).toContain('animate-spin')
    })

    it('applies sm size', () => {
        const wrapper = mount(Spinner, { props: { size: 'sm' } })
        expect(wrapper.classes()).toContain('w-4')
        expect(wrapper.classes()).toContain('h-4')
    })

    it('applies lg size', () => {
        const wrapper = mount(Spinner, { props: { size: 'lg' } })
        expect(wrapper.classes()).toContain('w-12')
        expect(wrapper.classes()).toContain('h-12')
    })

    it('falls back to md for unknown size', () => {
        const wrapper = mount(Spinner, { props: { size: undefined } })
        expect(wrapper.classes()).toContain('w-8')
    })
})