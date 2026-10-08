import { describe, it, expect, beforeEach } from 'vitest'
import { isArticleRead, markArticleRead } from '../readProgress'

describe('readProgress', () => {
    beforeEach(() => {
        localStorage.clear()
    })

    it('returns false for unread article', () => {
        expect(isArticleRead('art-1')).toBe(false)
    })

    it('marks article as read', () => {
        expect(markArticleRead('art-1')).toBe(true)
        expect(isArticleRead('art-1')).toBe(true)
    })

    it('returns false when marking already read article', () => {
        markArticleRead('art-2')
        expect(markArticleRead('art-2')).toBe(false)
    })
})