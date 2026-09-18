import MarkdownIt from 'markdown-it'
import DOMPurify, { type Config as DOMPurifyConfig } from 'dompurify'

const inline = new MarkdownIt({ html: false, linkify: false, breaks: false })

const excerptRules: DOMPurifyConfig = {
  FORBID_TAGS: ['a', 'style', 'form', 'input', 'button', 'textarea', 'select'],
}

export function renderExcerpt(markdown: string): string {
  return DOMPurify.sanitize(inline.renderInline(markdown), excerptRules) as unknown as string
}