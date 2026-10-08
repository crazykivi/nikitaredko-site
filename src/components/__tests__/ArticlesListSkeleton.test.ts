import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ArticlesListSkeleton from '../ArticlesListSkeleton.vue'

describe('ArticlesListSkeleton', () => {
    it('renders exactly 4 skeleton rows', () => {
        const wrapper = mount(ArticlesListSkeleton)
        const rows = wrapper.findAll('.animate-pulse')
        expect(rows).toHaveLength(4)
    })

    it('uses grid layout for md breakpoint', () => {
        const wrapper = mount(ArticlesListSkeleton)
        const row = wrapper.findAll('.animate-pulse')[0]
        expect(row.classes()).toContain('md:grid-cols-[88px_1fr]')
    })
})