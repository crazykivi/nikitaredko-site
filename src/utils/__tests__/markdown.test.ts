import { describe, it, expect, vi } from 'vitest'
import { renderExcerpt } from '../markdown'

vi.mock('dompurify', () => ({
    default: {
        sanitize: vi.fn((html: string) => html),
    },
}))

describe('renderExcerpt', () => {
    it('renders italic as <em> tag', () => {
        const html = renderExcerpt('*italic*')
        expect(html).toContain('<em>italic</em>')
    })

    it('preserves text content for bold and italic', () => {
        const html = renderExcerpt('**bold** and *italic*')
        expect(html).toContain('bold')
        expect(html).toContain('italic')
    })

    it('handles empty string', () => {
        expect(renderExcerpt('')).toBe('')
    })

    it('strips raw HTML input (html option is false)', () => {
        const html = renderExcerpt('text with <b>raw</b> html')
        expect(html).toContain('text with')
        expect(html).not.toContain('<b>raw</b>')
    })
})