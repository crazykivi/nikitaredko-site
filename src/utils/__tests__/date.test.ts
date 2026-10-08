import { describe, it, expect } from 'vitest'
import { formatDateNumeric, formatDateRu } from '../date'

describe('date utils', () => {
    it('formatDateNumeric returns non-empty string for valid date', () => {
        const res = formatDateNumeric('2025-03-15T12:00:00Z')
        expect(res.length).toBeGreaterThan(0)
        expect(res).toContain('2025')
    })

    it('formatDateNumeric returns empty for invalid date', () => {
        expect(formatDateNumeric('invalid')).toBe('')
    })

    it('formatDateRu returns non-empty string for valid date', () => {
        const res = formatDateRu('2025-03-15T12:00:00Z')
        expect(res.length).toBeGreaterThan(0)
        expect(res).toContain('2025')
    })

    it('formatDateRu returns empty for invalid date', () => {
        expect(formatDateRu('invalid')).toBe('')
    })
})