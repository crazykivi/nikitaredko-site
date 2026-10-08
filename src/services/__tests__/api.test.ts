import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { cleanCollectionName } from '../api'

describe('api service', () => {
  let mockFetch: ReturnType<typeof vi.fn>

  beforeEach(() => {
    vi.resetModules()
    mockFetch = vi.fn()
    vi.stubGlobal('fetch', mockFetch)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe('cleanCollectionName', () => {
    it('removes "collection " prefix (case-insensitive)', () => {
      expect(cleanCollectionName('Collection Blog')).toBe('Blog')
      expect(cleanCollectionName('COLLECTION Notes')).toBe('Notes')
      expect(cleanCollectionName('collection articles')).toBe('articles')
    })

    it('does not modify names without prefix', () => {
      expect(cleanCollectionName('Blog')).toBe('Blog')
      expect(cleanCollectionName('My Stuff')).toBe('My Stuff')
    })

    it('handles empty string', () => {
      expect(cleanCollectionName('')).toBe('')
    })

    it('trims whitespace', () => {
      expect(cleanCollectionName('Collection  Blog ')).toBe('Blog')
    })
  })

  describe('getCollections', () => {
    it('returns collections on success', async () => {
      const data = [{ id: 'c1', name: 'Blog', color: '#fff', icon: null }]
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getCollections } = await import('../api')
      const result = await getCollections()

      expect(result).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith('/api/collections', { signal: undefined })
    })

    it('throws on HTTP error', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const { getCollections } = await import('../api')
      await expect(getCollections()).rejects.toThrow('API error: 500')
    })

    it('passes abort signal', async () => {
      const ac = new AbortController()
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) })

      const { getCollections } = await import('../api')
      await getCollections(ac.signal)

      expect(mockFetch).toHaveBeenCalledWith('/api/collections', { signal: ac.signal })
    })

    it('propagates network errors', async () => {
      mockFetch.mockRejectedValueOnce(new TypeError('Failed to fetch'))

      const { getCollections } = await import('../api')
      await expect(getCollections()).rejects.toThrow('Failed to fetch')
    })
  })

  describe('getArticles', () => {
    it('returns articles list', async () => {
      const data = [{ id: 'a1', title: 'Hello' }]
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getArticles } = await import('../api')
      expect(await getArticles()).toEqual(data)
    })
  })

  describe('getArticlesStructured', () => {
    it('returns structured collections', async () => {
      const data = [{ id: 'c1', name: 'Blog', articles: [], articleCount: 0 }]
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getArticlesStructured } = await import('../api')
      expect(await getArticlesStructured()).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith('/api/articles/structured', { signal: undefined })
    })
  })

  describe('getArticle', () => {
    it('fetches single article by id', async () => {
      const data = { id: 'abc', title: 'Test', excerpt: 'Excerpt' }
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getArticle } = await import('../api')
      expect(await getArticle('abc')).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith('/api/articles/abc', { signal: undefined })
    })

    it('throws 404 for missing article', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 404 })

      const { getArticle } = await import('../api')
      await expect(getArticle('nope')).rejects.toThrow('API error: 404')
    })
  })

  describe('searchArticles', () => {
    it('returns empty array for empty query (no fetch)', async () => {
      const { searchArticles } = await import('../api')
      expect(await searchArticles('')).toEqual([])
      expect(mockFetch).not.toHaveBeenCalled()
    })

    it('returns empty array for whitespace-only query', async () => {
      const { searchArticles } = await import('../api')
      expect(await searchArticles('   ')).toEqual([])
      expect(mockFetch).not.toHaveBeenCalled()
    })

    it('encodes query and fetches results', async () => {
      const data = [{ id: 'a1', title: 'Result' }]
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { searchArticles } = await import('../api')
      const result = await searchArticles('hello world')

      expect(result).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith(
        '/api/articles/search?q=hello%20world',
        { signal: undefined },
      )
    })
  })

  describe('getArticlesFeed', () => {
    it('builds URL with page and limit', async () => {
      const data = { articles: [], total: 0, page: 1, limit: 10 }
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getArticlesFeed } = await import('../api')
      await getArticlesFeed(1, 10)

      expect(mockFetch).toHaveBeenCalledWith('/api/articles/feed?page=1&limit=10', { signal: undefined })
    })

    it('appends collection param when provided', async () => {
      const data = { articles: [], total: 0, page: 2, limit: 5 }
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getArticlesFeed } = await import('../api')
      await getArticlesFeed(2, 5, 'coll-42')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/articles/feed?page=2&limit=5&collection=coll-42',
        { signal: undefined },
      )
    })
  })

  describe('getUses', () => {
    it('returns uses page data', async () => {
      const data = { categories: [], lastUpdated: '2025-01-01' }
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getUses } = await import('../api')
      expect(await getUses()).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith('/api/uses', { signal: undefined })
    })
  })

  describe('getAbout', () => {
    it('returns about page data', async () => {
      const data = { intro: 'Hi', facts: [], career: [], stack: [] }
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(data) })

      const { getAbout } = await import('../api')
      expect(await getAbout()).toEqual(data)
      expect(mockFetch).toHaveBeenCalledWith('/api/about', { signal: undefined })
    })
  })

  describe('loading state', () => {
    it('resets loading after successful request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) })

      const { getCollections } = await import('../api')
      const { isGlobalLoading } = await import('../../utils/loading')

      await getCollections()
      expect(isGlobalLoading.value).toBe(false)
    })

    it('resets loading after failed request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

      const { getCollections } = await import('../api')
      const { isGlobalLoading } = await import('../../utils/loading')

      await expect(getCollections()).rejects.toThrow()
      expect(isGlobalLoading.value).toBe(false)
    })
  })
})