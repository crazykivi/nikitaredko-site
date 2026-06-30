<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, getArticlesStructured } from '../services/api'
import type { Article, CollectionWithArticles } from '../services/api'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()

const article = ref<Article | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const allCollections = ref<CollectionWithArticles[]>([])
const flatArticles = ref<Article[]>([])

let abortController: AbortController | null = null
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const flattenArticles = (articles: Article[]): Article[] => {
  const result: Article[] = []
  for (const a of articles) {
    result.push(a)
    if (a.children && a.children.length > 0) {
      result.push(...flattenArticles(a.children))
    }
  }
  return result
}
const currentCollectionId = computed(() => {
  const queryCollection = route.query.collection as string | undefined
  if (queryCollection) return queryCollection
  return article.value?.collectionId
})
const sortedArticles = computed(() => {
  const collId = currentCollectionId.value
  let articles: Article[] = []

  if (!collId) {
    for (const coll of allCollections.value) {
      articles.push(...flattenArticles(coll.articles))
    }
  } else {
    const coll = allCollections.value.find(c => c.id === collId)
    if (coll) articles = flattenArticles(coll.articles)
  }

  return articles.sort((a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime())
})

const currentIndex = computed(() => {
  if (!article.value) return -1
  return sortedArticles.value.findIndex(a => a.id === article.value!.id)
})

const prevArticle = computed(() => {
  const idx = currentIndex.value
  if (idx <= 0) return null
  return sortedArticles.value[idx - 1]
})

const nextArticle = computed(() => {
  const idx = currentIndex.value
  if (idx === -1 || idx >= sortedArticles.value.length - 1) return null
  return sortedArticles.value[idx + 1]
})

const goToArticle = (id: string) => {
  router.push(`/articles/${id}`)
}
const loadArticle = async (id: string) => {
  if (abortController) abortController.abort()
  abortController = new AbortController()
  const signal = abortController.signal

  loading.value = true
  error.value = null
  article.value = null

  try {
    const [data, collections] = await Promise.all([
      getArticle(id, signal),
      getArticlesStructured(signal)
    ])

    if (signal.aborted) return

    article.value = data
    allCollections.value = collections
    const collId = route.query.collection as string | undefined || data.collectionId
    let articles: Article[] = []
    if (!collId) {
      for (const coll of collections) {
        articles.push(...flattenArticles(coll.articles))
      }
    } else {
      const coll = collections.find(c => c.id === collId)
      if (coll) articles = flattenArticles(coll.articles)
    }
    flatArticles.value = articles.sort((a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime())

    window.scrollTo({ top: 0, behavior: 'smooth' })
  } catch (e: any) {
    if (e.name === 'AbortError') return
    error.value = e instanceof Error ? e.message : 'Failed to load article'
    console.error('Failed to load article:', e)
  } finally {
    if (!signal.aborted) loading.value = false
  }
}

const goBack = () => {
  const returnUrl = sessionStorage.getItem('return_url_articles')
  if (returnUrl) {
    router.push(returnUrl)
  } else {
    router.push('/articles')
  }
}
const goToCollection = (collectionId: string) => {
  router.push({ path: '/articles', query: { collection: collectionId } })
}

marked.setOptions({ breaks: true, gfm: true })

watch(() => route.params.id, (newId, oldId) => {
  if (newId && newId !== oldId) loadArticle(newId as string)
})

onMounted(async () => {
  await loadArticle(route.params.id as string)
})

onUnmounted(() => {
  if (abortController) abortController.abort()
})
</script>

<template>
  <div class="flex-1 min-w-0">
    <button
      @click="goBack"
      class="mb-8 text-muted hover:text-foreground transition-colors flex items-center gap-2"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
      Назад к статьям
    </button>

    <div v-if="loading" class="space-y-6 animate-pulse">
      <div class="h-4 bg-muted/50 rounded w-1/4 mb-4"></div>
      <div class="flex gap-2 mb-4">
        <div class="h-6 bg-muted/50 rounded w-20"></div>
        <div class="h-6 bg-muted/50 rounded w-24"></div>
      </div>
      <div class="h-12 bg-muted/50 rounded w-3/4 mb-4"></div>
      <div class="flex gap-4">
        <div class="h-4 bg-muted/50 rounded w-32"></div>
        <div class="h-4 bg-muted/50 rounded w-24"></div>
      </div>
      <div class="space-y-3 mt-8">
        <div v-for="i in 7" :key="i" class="h-4 bg-muted/50 rounded w-full"></div>
      </div>
    </div>

    <div v-else-if="error" class="text-center py-20">
      <p class="text-red-500">{{ error }}</p>
    </div>

    <article v-else-if="article" class="animate-fade-in">
      <div class="mb-8">
        <button
          v-if="article.collectionName && article.collectionId"
          @click="goToCollection(article.collectionId)"
          class="inline-flex items-center gap-1.5 mb-4 text-sm font-bold text-foreground hover:opacity-70 transition-opacity group"
          :title="`Все статьи из категории «${article.collectionName}»`"
        >
          <svg class="w-3.5 h-3.5 text-muted group-hover:text-foreground transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          {{ article.collectionName }}
          <svg class="w-3 h-3 text-muted opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
        <div class="flex flex-wrap gap-2 mb-4">
          <span
            v-for="tag in article.tags"
            :key="tag"
            class="text-xs px-2 py-1 rounded-md bg-muted/50 text-muted font-mono break-all"
          >
            {{ tag }}
          </span>
        </div>
        <h1 class="text-4xl font-bold mb-4 break-words">{{ article.title }}</h1>
        <div class="flex items-center gap-4 text-sm text-muted">
          <time>{{ formatDate(article.publishedAt) }}</time>
          <span>•</span>
          <span>{{ article.readTime }} мин чтения</span>
        </div>
      </div>

      <div class="prose prose-neutral dark:prose-invert max-w-none break-words">
        <div v-html="marked(article.content)"></div>
      </div>

      <nav
        v-if="prevArticle || nextArticle"
        class="mt-16 pt-8 border-t border-border grid gap-4"
        :class="{ 'grid-cols-1 md:grid-cols-2': prevArticle && nextArticle }"
      >
        <button
          v-if="prevArticle"
          @click="goToArticle(prevArticle.id)"
          class="group relative p-5 rounded-xl border border-border hover:border-foreground/30 bg-background hover:bg-muted/5 transition-all text-left flex items-center gap-4"
        >
          <div class="shrink-0 w-10 h-10 rounded-full bg-muted/30 group-hover:bg-muted/60 flex items-center justify-center transition-colors">
            <svg class="w-5 h-5 text-muted group-hover:text-foreground transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-xs text-muted uppercase tracking-wider mb-1">Предыдущая</div>
            <div class="font-medium text-foreground truncate group-hover:translate-x-[-2px] transition-transform">
              {{ prevArticle.title }}
            </div>
            <div class="text-xs text-muted mt-1">{{ formatDate(prevArticle.publishedAt) }}</div>
          </div>
        </button>
        <div v-else-if="!prevArticle && !nextArticle"></div>

        <button
          v-if="nextArticle"
          @click="goToArticle(nextArticle.id)"
          class="group relative p-5 rounded-xl border border-border hover:border-foreground/30 bg-background hover:bg-muted/5 transition-all text-left flex items-center gap-4 md:justify-end md:text-right"
        >
          <div class="flex-1 min-w-0">
            <div class="text-xs text-muted uppercase tracking-wider mb-1">Следующая</div>
            <div class="font-medium text-foreground truncate group-hover:translate-x-[2px] transition-transform">
              {{ nextArticle.title }}
            </div>
            <div class="text-xs text-muted mt-1">{{ formatDate(nextArticle.publishedAt) }}</div>
          </div>
          <div class="shrink-0 w-10 h-10 rounded-full bg-muted/30 group-hover:bg-muted/60 flex items-center justify-center transition-colors">
            <svg class="w-5 h-5 text-muted group-hover:text-foreground transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </button>
      </nav>
    </article>

    <div v-else class="text-center py-20">
      <h2 class="text-2xl font-semibold mb-2">Статья не найдена</h2>
      <p class="text-muted">Возможно, она была удалена или перемещена</p>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.prose {
  max-width: 100%;
  overflow-wrap: break-word;
  word-wrap: break-word;
  word-break: break-word;
  hyphens: auto;
}
.prose :deep(*) { max-width: 100%; overflow-wrap: break-word; word-wrap: break-word; }
.prose :deep(pre) { max-width: 100%; overflow-x: auto; white-space: pre-wrap; word-wrap: break-word; }
.prose :deep(code) { word-break: break-word; }
.prose :deep(table) { display: block; max-width: 100%; overflow-x: auto; }
.prose :deep(img) { max-width: 100%; height: auto; }
.prose :deep(h1) { @apply text-3xl font-bold mt-8 mb-4; }
.prose :deep(h2) { @apply text-2xl font-semibold mt-6 mb-3; }
.prose :deep(h3) { @apply text-xl font-semibold mt-4 mb-2; }
.prose :deep(p) { @apply mb-4 leading-relaxed; }
.prose :deep(ul), .prose :deep(ol) { @apply mb-4 ml-6; }
.prose :deep(li) { @apply mb-2; }
.prose :deep(code) { @apply px-2 py-1 rounded font-mono text-sm; background-color: rgb(113 113 122 / 0.5); }
.prose :deep(pre) { @apply p-4 rounded-lg bg-muted overflow-x-auto mb-4; }
.prose :deep(pre code) { @apply p-0 bg-transparent; }
.prose :deep(strong) { @apply font-semibold; }
.prose :deep(a) { @apply text-foreground underline hover:opacity-70 transition-opacity; }
</style>