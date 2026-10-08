import MarkdownIt from 'markdown-it'
import markdownItContainer from 'markdown-it-container'
import hljs from 'highlight.js'

const BLOCK_TYPES = ['warning', 'info', 'success', 'danger', 'tip', 'note']
const STRUCTURAL_BLOCKS = new Set(['stats', 'timeline', 'stack'])

export const slugify = (text: string): string =>
    text
        .toLowerCase()
        .replace(/[^\wа-яё\s-]/gi, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')

function renderHighlight(str: string, lang: string): string {
    const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext'
    let highlighted: string
    try {
        highlighted = hljs.highlight(str, { language }).value
    } catch {
        highlighted = hljs.highlightAuto(str).value
    }
    const safeLang = language.replace(/"/g, '&quot;')
    return (
        `<div class="code-block-wrapper relative group" data-language="${safeLang}">` +
        `<div class="code-header justify-between flex items-center gap-3 px-4 py-2 bg-muted/80 border-b border-border rounded-t-lg">` +
        `<span class="text-xs text-muted font-mono uppercase">${safeLang}</span>` +
        `</div>` +
        `<pre class="!mt-0 !mb-0 !rounded-t-none"><code class="hljs language-${safeLang}">${highlighted}</code></pre>` +
        `</div>`
    )
}

function registerContainers(md: MarkdownIt): void {
    for (const blockType of BLOCK_TYPES) {
        md.use(markdownItContainer, blockType, {
            validate: (params: string) => params.trim() === blockType,
            render: (tokens: any[], idx: number) => {
                if (tokens[idx].nesting === 1) {
                    const showTitle = !STRUCTURAL_BLOCKS.has(blockType)
                    const title = blockType.charAt(0).toUpperCase() + blockType.slice(1)
                    const titleHtml = showTitle ? `<p class="custom-block-title">${title}</p>\n` : ''
                    return `<div class="custom-block ${blockType}">\n${titleHtml}`
                }
                return '</div>\n'
            },
        })
    }
}

function overrideHeadings(md: MarkdownIt): void {
    const defaultHeadingOpen =
        md.renderer.rules.heading_open ??
        ((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options))

    md.renderer.rules.heading_open = (tokens, idx, options, env, self) => {
        const nextToken = tokens[idx + 1]
        if (nextToken?.type === 'inline') {
            tokens[idx].attrSet('id', slugify(nextToken.content))
        }
        return defaultHeadingOpen(tokens, idx, options, env, self)
    }
}

export function createMarkdownRenderer(): MarkdownIt {
    const md = new MarkdownIt({
        html: true,
        linkify: true,
        typographer: true,
        breaks: true,
        highlight: renderHighlight,
    })
    registerContainers(md)
    overrideHeadings(md)
    return md
}