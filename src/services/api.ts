import { startLoading, stopLoading } from '../utils/loading'
const API_BASE = '/api'

export function cleanCollectionName(name: string): string {
  return name.replace(/^collection\s+/i, '').trim()
}

export interface Article {
  id: string
  title: string
  excerpt: string
  content: string
  publishedAt: string
  readTime: number
  tags: string[]
  collectionId: string
  collectionName: string
  isDraft: boolean
  children?: Article[]
  level?: number
}

export interface Collection {
  id: string
  name: string
  color: string | null
  icon: string | null
}

export interface CollectionWithArticles {
  id: string
  name: string
  color: string | null
  icon: string | null
  articles: Article[]
  articleCount: number
}

async function request<T>(endpoint: string, signal?: AbortSignal): Promise<T> {
  startLoading()
  try {
    const response = await fetch(`${API_BASE}${endpoint}`, { signal })
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`)
    }
    return await response.json()
  } finally {
    stopLoading()
  }
}
export async function getCollections(signal?: AbortSignal): Promise<Collection[]> {
  return request<Collection[]>('/collections', signal)
}

export async function getArticles(signal?: AbortSignal): Promise<Article[]> {
  return request<Article[]>('/articles', signal)
}

export async function getArticlesStructured(signal?: AbortSignal): Promise<CollectionWithArticles[]> {
  return request<CollectionWithArticles[]>('/articles/structured', signal)
}

export async function getArticle(id: string, signal?: AbortSignal): Promise<Article> {
  return request<Article>(`/articles/${id}`, signal)
}

export async function searchArticles(query: string, signal?: AbortSignal): Promise<Article[]> {
    if (!query.trim()) return []
    return request<Article[]>(`/articles/search?q=${encodeURIComponent(query)}`, signal)
}