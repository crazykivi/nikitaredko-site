import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle, getArticlesStructured } from '../services/api'
import type { Article, CollectionWithArticles } from '../services/api'

function flattenArticles(articles: Article[]): Article[] {
    const result: Article[] = []
    for (const a of articles) {
        result.push(a)
        if (a.children?.length) result.push(...flattenArticles(a.children))
    }
    return result
}

export function useArticleLoader() {
    const route = useRoute()

    const article = ref<Article | null>(null)
    const loading = ref(true)
    const error = ref<string | null>(null)
    const allCollections = ref<CollectionWithArticles[]>([])

    let abortController: AbortController | null = null

    const currentCollectionId = computed(() => {
        return (route.query.collection as string) ?? article.value?.collectionId
    })

    const sortedArticles = computed(() => {
        const collId = currentCollectionId.value
        let articles: Article[] = []
        if (!collId) {
            for (const coll of allCollections.value) {
                articles.push(...flattenArticles(coll.articles))
            }
        } else {
            const coll = allCollections.value.find((c) => c.id === collId)
            if (coll) articles = flattenArticles(coll.articles)
        }
        return articles.sort(
            (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
        )
    })

    const currentIndex = computed(() => {
        if (!article.value) return -1
        return sortedArticles.value.findIndex((a) => a.id === article.value!.id)
    })

    const prevArticle = computed(() => {
        const idx = currentIndex.value
        return idx > 0 ? sortedArticles.value[idx - 1] : null
    })

    const nextArticle = computed(() => {
        const idx = currentIndex.value
        if (idx === -1 || idx >= sortedArticles.value.length - 1) return null
        return sortedArticles.value[idx + 1]
    })

    const loadArticle = async (id: string) => {
        abortController?.abort()
        abortController = new AbortController()
        const signal = abortController.signal

        loading.value = true
        error.value = null
        article.value = null

        try {
            const [data, collections] = await Promise.all([
                getArticle(id, signal),
                getArticlesStructured(signal),
            ])
            if (signal.aborted) return
            article.value = data
            allCollections.value = collections
            window.scrollTo({ top: 0, behavior: 'smooth' })
        } catch (e: any) {
            if (e.name === 'AbortError') return
            error.value = e instanceof Error ? e.message : 'Failed to load article'
        } finally {
            if (!signal.aborted) loading.value = false
        }
    }

    watch(
        () => route.params.id,
        (newId, oldId) => {
            if (newId && newId !== oldId) loadArticle(newId as string)
        },
    )

    return {
        article, loading, error,
        prevArticle, nextArticle,
        loadArticle,
    }
}