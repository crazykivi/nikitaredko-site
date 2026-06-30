export interface Article {
  id: string
  title: string
  excerpt: string
  content: string
  publishedAt: string
  readTime: number
  tags: string[]
  coverImage?: string
}