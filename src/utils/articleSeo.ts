import type { ReactiveHead } from '@unhead/vue'
import type { Article } from '../services/api'

const SITE_NAME = 'Nikita Redko'
const FORBIDDEN_NODES = 'script,style,noscript,template,iframe,object,embed,svg,math'

function toPlainText(input: string): string {
    const doc = new DOMParser().parseFromString(input, 'text/html')
    doc.querySelectorAll(FORBIDDEN_NODES).forEach((el) => el.remove())
    return (doc.body?.textContent ?? '').split(/\s+/).join(' ').trim()
}

export function buildArticleSeo(article: Article | null): ReactiveHead {
    const origin = typeof window !== 'undefined' ? window.location.origin : ''
    const currentUrl = typeof window !== 'undefined' ? window.location.href : ''
    const ogImage = `${origin}/api/og/${article?.id ?? ''}`

    if (!article) {
        return { title: SITE_NAME, script: [] }
    }

    const title = `${article.title} | ${SITE_NAME}`
    const rawDesc = article.excerpt || article.content || ''
    const description =
        toPlainText(rawDesc).substring(0, 150).trim() || 'Статья на сайте Никиты Редко'

    const jsonLd = {
        '@context': 'https://schema.org',
        '@type': 'Article',
        headline: article.title,
        description,
        datePublished: article.publishedAt || article.createdAt,
        dateModified: article.createdAt,
        author: { '@type': 'Person', name: SITE_NAME, url: origin },
        publisher: {
            '@type': 'Organization',
            name: SITE_NAME,
            logo: { '@type': 'ImageObject', url: ogImage },
        },
        mainEntityOfPage: { '@type': 'WebPage', '@id': currentUrl },
        image: ogImage,
        articleSection: article.collectionName || '',
        keywords: article.tags?.join(', ') || '',
        inLanguage: 'ru',
    }

    return {
        title,
        meta: [
            { name: 'description', content: description },
            { property: 'og:title', content: title },
            { property: 'og:description', content: description },
            { property: 'og:type', content: 'article' },
            { property: 'og:url', content: currentUrl },
            { property: 'og:image', content: ogImage },
            { property: 'og:site_name', content: SITE_NAME },
            { name: 'twitter:card', content: 'summary_large_image' },
            { name: 'twitter:title', content: title },
            { name: 'twitter:description', content: description },
            { name: 'twitter:image', content: ogImage },
        ],
        link: [{ rel: 'canonical' as const, href: currentUrl }],
        script: [
            {
                type: 'application/ld+json',
                innerHTML: JSON.stringify(jsonLd).replace(/</g, '\\u003c'),
            },
        ],
    }
}