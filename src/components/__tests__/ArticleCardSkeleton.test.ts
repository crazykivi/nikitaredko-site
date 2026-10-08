import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ArticleCardSkeleton from '../ArticleCardSkeleton.vue'

describe('ArticleCardSkeleton', () => {
    it('renders with animate-pulse', () => {
        const wrapper = mount(ArticleCardSkeleton)
        expect(wrapper.find('.animate-pulse').exists()).toBe(true)
    })

    it('contains placeholder blocks', () => {
        const wrapper = mount(ArticleCardSkeleton)
        const blocks = wrapper.findAll('.bg-muted\\/50')
        expect(blocks.length).toBeGreaterThanOrEqual(4)
    })
})