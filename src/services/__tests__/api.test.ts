import { describe, it, expect } from 'vitest'
import { cleanCollectionName } from '../api'

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