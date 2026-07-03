<script setup lang="ts">
defineOptions({ name: 'ArticlesView' })
import { ref, computed, onMounted, onUnmounted, onActivated, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { getArticlesFeed } from '../services/api'
import type { Article } from '../services/api'
import ArticleCard from '../components/ArticleCard.vue'
import Spinner from '../components/Spinner.vue'

const route = useRoute()
const articles = ref<Article[]>([])
const total = ref(0)
const currentPage = ref(1)
const limit = 10
const loading = ref(false)
const error = ref<string | null>(null)

const isMobile = ref(window.innerWidth < 768)
const sentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const updateIsMobile = () => { isMobile.value = window.innerWidth < 768 }
onMounted(() => window.addEventListener('resize', updateIsMobile))
onUnmounted(() => {
  window.removeEventListener('resize', updateIsMobile)
  if (observer) observer.disconnect()
})

const totalPages = computed(() => Math.ceil(total.value / limit))
const hasMore = computed(() => currentPage.value * limit < total.value)

const fetchArticles = async (page: number, append: boolean = false) => {
  if (loading.value) return
  loading.value = true
  error.value = null
  
  const collectionParam = route.query.collection as string | undefined
  
  try {
    const res = await getArticlesFeed(page, limit, collectionParam)
    total.value = res.total
    
    if (append) {
      articles.value = [...articles.value, ...res.articles]
    } else {
      articles.value = res.articles
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load articles'
    console.error('Failed to load articles:', e)
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (hasMore.value && !loading.value) {
    currentPage.value++
    fetchArticles(currentPage.value, true)
  }
}

const goToPage = (page: number) => {
  currentPage.value = page
  window.scrollTo({ top: 0, behavior: 'smooth' })
  fetchArticles(page, false)
}

const paginationRange = computed(() => {
  const totalP = totalPages.value
  const current = currentPage.value
  const delta = 2
  const range: (number | string)[] = []
  for (let i = 1; i <= totalP; i++) {
    if (i === 1 || i === totalP || (i >= current - delta && i <= current + delta)) {
      range.push(i)
    } else if (range[range.length - 1] !== '...') {
      range.push('...')
    }
  }
  return range
})

const getWord = (count: number) => {
  const n = count % 100
  if (n >= 11 && n <= 14) return 'статей'
  const last = n % 10
  if (last === 1) return 'статья'
  if (last >= 2 && last <= 4) return 'статьи'
  return 'статей'
}

const collectionName = computed(() => {
  return route.query.collection ? 'Статьи в категории' : 'Все статьи'
})
watch(() => route.query.collection, () => {
  currentPage.value = 1
  articles.value = []
  window.scrollTo({ top: 0, behavior: 'auto' })
  fetchArticles(1, false)
})
watch([sentinel, isMobile], ([newSentinel, mobile]) => {
  if (observer) observer.disconnect()
  if (mobile && newSentinel) {
    observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) loadMore()
    }, { rootMargin: '200px' })
    observer.observe(newSentinel)
  }
}, { immediate: true })

onMounted(async () => {
  await fetchArticles(1, false)
})
onActivated(() => {
  const saved = sessionStorage.getItem('scroll_articles')
  if (saved) {
    setTimeout(() => window.scrollTo({ top: Number(saved), behavior: 'auto' }), 150)
  }
})
</script>
<template>
  <div>
    <div class="mb-10">
      <h1 class="text-5xl font-bold mb-4">{{ collectionName }}</h1>
      <p v-if="!loading" class="text-muted text-lg">
        {{ `${total} ${getWord(total)}` }}
      </p>
    </div>
    <div v-if="loading && articles.length === 0" class="space-y-6">
      <div v-for="i in 3" :key="i" class="p-6 rounded-xl border border-border bg-background animate-pulse">
        <div class="h-4 bg-muted/50 rounded w-1/4 mb-4"></div>
        <div class="flex gap-2 mb-3">
          <div class="h-5 bg-muted/50 rounded w-16"></div>
          <div class="h-5 bg-muted/50 rounded w-20"></div>
        </div>
        <div class="h-7 bg-muted/50 rounded w-3/4 mb-3"></div>
        <div class="h-4 bg-muted/50 rounded w-full mb-2"></div>
      </div>
    </div>

    <div v-else-if="error" class="text-center py-20">
      <p class="text-red-500">{{ error }}</p>
      <p class="text-muted mt-2">Проверьте, что бекенд запущен</p>
    </div>

    <div v-else-if="total === 0" class="text-center py-20">
      <p class="text-muted">В этой категории пока нет статей</p>
    </div>

    <div v-else>
      <div class="space-y-6">
        <ArticleCard
          v-for="article in articles"
          :key="article.id"
          :article="article"
        />
      </div>

      <div v-if="isMobile && hasMore" ref="sentinel" class="mt-8 h-16 flex items-center justify-center">
        <Spinner size="sm" />
        <span class="ml-2 text-sm text-muted">Загрузка...</span>
      </div>
      <div v-else-if="isMobile && !hasMore && total > limit" class="mt-8 text-center text-sm text-muted">
        Вы достигли конца списка
      </div>

      <div v-if="!isMobile && totalPages > 1" class="mt-12 flex items-center justify-center gap-2">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="px-4 py-2 rounded-lg border border-border text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-muted/50"
        >
          Назад
        </button>
        <div class="flex items-center gap-1">
          <template v-for="page in paginationRange" :key="page">
            <button
              v-if="page !== '...'"
              @click="goToPage(page as number)"
              :class="[
                'w-10 h-10 rounded-lg text-sm font-medium transition-colors',
                currentPage === page
                  ? 'bg-foreground text-background'
                  : 'hover:bg-muted/50 text-muted hover:text-foreground'
              ]"
            >
              {{ page }}
            </button>
            <span v-else class="w-10 h-10 flex items-center justify-center text-muted">...</span>
          </template>
        </div>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="px-4 py-2 rounded-lg border border-border text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-muted/50"
        >
          Вперед
        </button>
      </div>
    </div>
  </div>
</template>