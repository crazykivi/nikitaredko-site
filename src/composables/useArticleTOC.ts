import { computed, type Ref } from 'vue'
import { slugify } from '../utils/markdownRenderer'

export interface TOCItem {
    id: string
    text: string
    level: number
}

export function useArticleTOC(content: Ref<string | undefined>) {
    const tocItems = computed<TOCItem[]>(() => {
        if (!content.value) return []
        return extractHeadings(content.value)
    })
    return { tocItems }
}

function extractHeadings(markdown: string): TOCItem[] {
    const headings: TOCItem[] = []
    let inCodeBlock = false

    for (const line of markdown.split('\n')) {
        const trimmed = line.trim()
        if (trimmed.startsWith('```')) {
            inCodeBlock = !inCodeBlock
            continue
        }
        if (inCodeBlock) continue

        const match = trimmed.match(/^(#{2,3})\s+(.+)$/)
        if (match) {
            headings.push({
                id: slugify(match[2].trim()),
                text: match[2].trim(),
                level: match[1].length,
            })
        }
    }
    return headings
}