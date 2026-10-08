import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('readProgress', () => {
  beforeEach(() => {
    vi.resetModules()
    localStorage.clear()
  })

  it('returns false for unread article', async () => {
    const { isArticleRead } = await import('../readProgress')
    expect(isArticleRead('art-1')).toBe(false)
  })

  it('marks article as read', async () => {
    const { isArticleRead, markArticleRead } = await import('../readProgress')
    expect(markArticleRead('art-1')).toBe(true)
    expect(isArticleRead('art-1')).toBe(true)
  })

  it('returns false when marking already read article', async () => {
    const { markArticleRead } = await import('../readProgress')
    markArticleRead('art-2')
    expect(markArticleRead('art-2')).toBe(false)
  })

  it('returns false for non-array JSON in storage', async () => {
    localStorage.setItem('mc-articles-read', '{"key":"value"}')
    const { isArticleRead } = await import('../readProgress')
    expect(isArticleRead('art-1')).toBe(false)
  })

  it('returns false for invalid JSON in storage (catch branch)', async () => {
    localStorage.setItem('mc-articles-read', '!!!{broken')
    const { isArticleRead } = await import('../readProgress')
    expect(isArticleRead('art-1')).toBe(false)
  })

  it('filters non-string entries from stored array', async () => {
    localStorage.setItem('mc-articles-read', '["art-1", 42, null, true, "art-2"]')
    const { isArticleRead } = await import('../readProgress')
    expect(isArticleRead('art-1')).toBe(true)
    expect(isArticleRead('art-2')).toBe(true)
    expect(isArticleRead('42')).toBe(false)
  })

  it('migrates legacy keys and removes them from storage', async () => {
    localStorage.setItem('article-xp-art-old', 'true')
    localStorage.setItem('article-xp-art-ancient', 'true')
    const { isArticleRead } = await import('../readProgress')

    expect(isArticleRead('art-old')).toBe(true)
    expect(isArticleRead('art-ancient')).toBe(true)
    // легаси-ключи удалены
    expect(localStorage.getItem('article-xp-art-old')).toBeNull()
    expect(localStorage.getItem('article-xp-art-ancient')).toBeNull()
  })

  it('preserves existing entries when migrating legacy keys', async () => {
    localStorage.setItem('mc-articles-read', '["art-keep"]')
    localStorage.setItem('article-xp-art-new', 'true')
    const { isArticleRead } = await import('../readProgress')

    expect(isArticleRead('art-keep')).toBe(true)
    expect(isArticleRead('art-new')).toBe(true)
  })

  it('evicts oldest entries beyond MAX_ENTRIES (5)', async () => {
    const { markArticleRead, isArticleRead } = await import('../readProgress')
    for (let i = 1; i <= 7; i++) {
      markArticleRead(`art-${i}`)
    }
    expect(isArticleRead('art-1')).toBe(false)
    expect(isArticleRead('art-2')).toBe(false)
    for (let i = 3; i <= 7; i++) {
      expect(isArticleRead(`art-${i}`)).toBe(true)
    }
  })

  it('does not re-run migration on subsequent calls', async () => {
    localStorage.setItem('article-xp-art-x', 'true')
    const { isArticleRead, markArticleRead } = await import('../readProgress')

    expect(isArticleRead('art-x')).toBe(true)
    expect(localStorage.getItem('article-xp-art-x')).toBeNull()
    markArticleRead('art-y')
    expect(isArticleRead('art-y')).toBe(true)
    expect(isArticleRead('art-x')).toBe(true)
  })
})