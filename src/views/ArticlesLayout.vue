<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticlesStructured, searchArticles, cleanCollectionName } from '../services/api'
import type { CollectionWithArticles, Article } from '../services/api'

const route = useRoute()
const router = useRouter()

const collections = ref<CollectionWithArticles[]>([])
const searchResults = ref<Article[]>([])
const isSearching = ref(false)
const expandedCollections = ref<Set<string>>(new Set())
const isSidebarCollapsed = ref(false)
const searchQuery = ref('')
const searchAbortController = ref<AbortController | null>(null)

const isSearchActive = computed(() => searchQuery.value.trim().length > 0)

const filteredCollections = computed(() => {
    if (isSearchActive.value) {
        const grouped = new Map<string, CollectionWithArticles>()
        
        for (const article of searchResults.value) {
            if (!grouped.has(article.collectionId)) {
                const coll = collections.value.find(c => c.id === article.collectionId)
                if (coll) {
                    grouped.set(article.collectionId, {
                        ...coll,
                        articles: [],
                        articleCount: 0
                    })
                }
            }
            const group = grouped.get(article.collectionId)
            if (group) {
                group.articles.push(article)
                group.articleCount++
            }
        }
        
        return Array.from(grouped.values())
    }
    return collections.value
})

const SIDEBAR_INITIAL_LIMIT = 5
const SIDEBAR_LOAD_MORE_STEP = 20
const expandedSidebarLimits = ref<Map<string, number>>(new Map())

const getSidebarArticles = (coll: CollectionWithArticles) => {
    if (isSearchActive.value) {
        return coll.articles
    }
    const currentLimit = expandedSidebarLimits.value.get(coll.id) || SIDEBAR_INITIAL_LIMIT
    return coll.articles.slice(0, currentLimit)
}

const hasMoreArticles = (coll: CollectionWithArticles) => {
    if (isSearchActive.value) return false
    const currentLimit = expandedSidebarLimits.value.get(coll.id) || SIDEBAR_INITIAL_LIMIT
    return coll.articles.length > currentLimit
}

const loadMoreArticles = (id: string) => {
    const currentLimit = expandedSidebarLimits.value.get(id) || SIDEBAR_INITIAL_LIMIT
    expandedSidebarLimits.value.set(id, currentLimit + SIDEBAR_LOAD_MORE_STEP)
}

const resetSidebarLimits = () => {
    expandedSidebarLimits.value.clear()
}

const toggleCollection = (id: string) => {
    if (expandedCollections.value.has(id)) {
        expandedCollections.value.delete(id)
    } else {
        expandedCollections.value.add(id)
    }
}

const selectCollection = (id: string | null) => {
    if (id) {
        router.push({ path: '/articles', query: { collection: id } })
    } else {
        router.push('/articles')
    }
}

const toggleSidebar = () => {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
}

const onCollapsedCategoryClick = (id: string) => {
    isSidebarCollapsed.value = false
    expandedCollections.value.add(id)
    router.push({ path: '/articles', query: { collection: id } })
}

const isCollectionExpanded = (id: string) => expandedCollections.value.has(id)

const getCategoryLabel = (coll: CollectionWithArticles) => {
    if (coll.icon) return coll.icon
    return cleanCollectionName(coll.name).charAt(0).toUpperCase()
}

const getSelectedCollectionId = () => {
    return route.query.collection as string | undefined
}

let searchTimeout: ReturnType<typeof setTimeout> | null = null

const performSearch = async (query: string) => {
    if (searchTimeout) {
        clearTimeout(searchTimeout)
    }

    if (!query.trim()) {
        searchResults.value = []
        isSearching.value = false
        return
    }

    searchTimeout = setTimeout(async () => {
        if (searchAbortController.value) {
            searchAbortController.value.abort()
        }
        
        searchAbortController.value = new AbortController()
        isSearching.value = true
        
        try {
            const results = await searchArticles(query, searchAbortController.value.signal)
            searchResults.value = Array.isArray(results) ? results : []
            
            for (const coll of collections.value) {
                const hasResults = searchResults.value.some(r => r.collectionId === coll.id)
                if (hasResults) {
                    expandedCollections.value.add(coll.id)
                }
            }
        } catch (e) {
            if (e instanceof Error && e.name !== 'AbortError') {
                console.error('Search failed:', e)
                searchResults.value = []
            }
        } finally {
            isSearching.value = false
        }
    }, 300)
}

watch(searchQuery, (newQuery) => {
    performSearch(newQuery)
    resetSidebarLimits()
})

onMounted(async () => {
    const collectionParam = route.query.collection as string | undefined
    try {
        collections.value = await getArticlesStructured()
        for (const coll of collections.value) {
            expandedCollections.value.add(coll.id)
        }
        if (collectionParam) {
            expandedCollections.value.add(collectionParam)
        }
    } catch (e) {
        console.error('Failed to load collections:', e)
    }
})

watch(
    () => route.query.collection,
    (newCollection) => {
        if (newCollection) {
            expandedCollections.value.add(newCollection as string)
        }
    }
)

onUnmounted(() => {
    if (searchTimeout) {
        clearTimeout(searchTimeout)
    }
    if (searchAbortController.value) {
        searchAbortController.value.abort()
    }
})
</script>

<template>
  <div class="min-h-[calc(100vh-5rem)] pt-10 pb-10 px-4">
    <div class="max-w-7xl mx-auto flex gap-6">
      <aside
        :class="[
          'shrink-0 hidden lg:block transition-all duration-300',
          isSidebarCollapsed ? 'w-16' : 'w-72'
        ]"
      >
        <div class="sticky top-24">
          <div class="flex items-center justify-between mb-4 px-2">
            <h2
              v-if="!isSidebarCollapsed"
              class="text-sm font-semibold uppercase tracking-wider text-muted"
            >
              Категории
            </h2>
            <button
              @click="toggleSidebar"
              class="p-1.5 rounded-lg hover:bg-muted/50 transition-colors text-muted hover:text-foreground"
              :aria-label="isSidebarCollapsed ? 'Развернуть панель' : 'Свернуть панель'"
              :title="isSidebarCollapsed ? 'Развернуть' : 'Свернуть'"
            >
              <svg
                class="w-4 h-4 transition-transform duration-300"
                :class="{ 'rotate-180': isSidebarCollapsed }"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
              </svg>
            </button>
          </div>
          <div v-if="!isSidebarCollapsed" class="mb-4 px-2">
              <div class="relative">
                  <svg
                      class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                  >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  <input
                      v-model="searchQuery"
                      type="text"
                      placeholder="Поиск по статьям..."
                      class="w-full pl-9 pr-9 py-2 text-sm rounded-lg border border-border bg-background focus:outline-none focus:ring-2 focus:ring-foreground/20 transition-all"
                  />
                  <button
                      v-if="searchQuery || isSearching"
                      @click="searchQuery = ''"
                      class="absolute right-2 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-muted/50 transition-colors text-muted hover:text-foreground"
                  >
                      <svg v-if="!isSearching" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      <svg v-else class="w-3 h-3 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                      </svg>
                  </button>
              </div>
          </div>
          <div v-if="isSidebarCollapsed" class="space-y-1">
            <button
              @click="selectCollection(null)"
              :class="[
                'w-full flex items-center justify-center p-2 rounded-lg transition-colors text-sm relative',
                !getSelectedCollectionId()
                  ? 'bg-muted/80 text-foreground font-medium'
                  : 'text-muted hover:text-foreground hover:bg-muted/30'
              ]"
              title="Все статьи"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <span
                v-if="collections.reduce((sum, c) => sum + c.articleCount, 0) > 0"
                class="absolute -top-1 -right-1 w-4 h-4 text-[10px] bg-foreground text-background rounded-full flex items-center justify-center"
              >
                {{ collections.reduce((sum, c) => sum + c.articleCount, 0) }}
              </span>
            </button>
            <button
              v-for="coll in collections"
              :key="coll.id"
              @click="onCollapsedCategoryClick(coll.id)"
              :class="[
                'w-full flex items-center justify-center p-2 rounded-lg transition-colors text-sm relative',
                getSelectedCollectionId() === coll.id
                  ? 'bg-muted/80 text-foreground font-medium'
                  : 'text-muted hover:text-foreground hover:bg-muted/30'
              ]"
              :title="cleanCollectionName(coll.name)"
            >
              <span class="text-base">{{ getCategoryLabel(coll) }}</span>
              <span
                v-if="coll.articleCount > 0"
                class="absolute -top-1 -right-1 w-4 h-4 text-[10px] bg-foreground text-background rounded-full flex items-center justify-center"
              >
                {{ coll.articleCount }}
              </span>
            </button>
          </div>
          <div v-else>
            <button
              @click="selectCollection(null)"
              :class="[
                'w-full text-left px-3 py-2 rounded-lg mb-1 transition-colors text-sm',
                !getSelectedCollectionId()
                  ? 'bg-muted/80 text-foreground font-medium'
                  : 'text-muted hover:text-foreground hover:bg-muted/30'
              ]"
            >
              Все статьи
              <span class="float-right text-xs opacity-60">
                {{ collections.reduce((sum, c) => sum + c.articleCount, 0) }}
              </span>
            </button>
            <div class="space-y-1 mt-2">
              <div v-for="coll in filteredCollections" :key="coll.id">
                <div class="flex items-center gap-1">
                  <button
                    @click="toggleCollection(coll.id)"
                    class="p-1 rounded hover:bg-muted/50 transition-colors text-muted"
                    :aria-label="isCollectionExpanded(coll.id) ? 'Свернуть' : 'Развернуть'"
                  >
                    <svg
                      class="w-3 h-3 transition-transform"
                      :class="{ 'rotate-90': isCollectionExpanded(coll.id) }"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                  <button
                    @click="selectCollection(coll.id)"
                    :class="[
                      'flex-1 text-left px-2 py-1.5 rounded-lg transition-colors text-sm flex items-center gap-2',
                      getSelectedCollectionId() === coll.id
                        ? 'bg-muted/80 text-foreground font-medium'
                        : 'text-muted hover:text-foreground hover:bg-muted/30'
                    ]"
                  >
                    <svg
                      v-if="!coll.icon"
                      class="w-4 h-4 shrink-0 opacity-60"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                    </svg>
                    <span v-else class="text-base shrink-0">{{ coll.icon }}</span>
                    
                    <span class="flex-1 truncate">{{ cleanCollectionName(coll.name) }}</span>
                    <span
                      v-if="coll.articleCount > 0"
                      class="text-[10px] px-1.5 py-0.5 rounded-full bg-muted/50 text-muted font-mono shrink-0"
                    >
                      {{ coll.articleCount }}
                    </span>
                  </button>
                </div>
                <div
                    v-if="isCollectionExpanded(coll.id)"
                    class="ml-6 mt-1 space-y-0.5 border-l border-border pl-2"
                >
                    <router-link
                        v-for="article in getSidebarArticles(coll)"
                        :key="article.id"
                        :to="`/articles/${article.id}`"
                        class="block px-2 py-1 text-xs transition-colors truncate rounded"
                        :class="isSearchActive ? 'text-foreground hover:bg-muted/50' : 'text-muted hover:text-foreground'"
                        :title="article.title"
                    >
                        {{ article.title }}
                        <span v-if="isSearchActive && article.tags.includes('title-match')" class="text-[9px] ml-1 text-muted">
                            (заголовок)
                        </span>
                    </router-link>
                    <p v-if="coll.articles.length === 0" class="text-xs text-muted/50 px-2 py-1 italic">
                        {{ isSearchActive ? 'Ничего не найдено' : 'Пока пусто' }}
                    </p>
                    <button
                        v-if="hasMoreArticles(coll)"
                        @click="loadMoreArticles(coll.id)"
                        class="w-full text-left px-2 py-1 text-xs text-muted hover:text-foreground transition-colors rounded flex items-center gap-1 mt-1 opacity-60 hover:opacity-100"
                    >
                        <span>Показать ещё {{ coll.articles.length - (expandedSidebarLimits.get(coll.id) || SIDEBAR_INITIAL_LIMIT) }}</span>
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                        </svg>
                    </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <main class="flex-1 min-w-0">
        <router-view v-slot="{ Component, route }">
          <Suspense>
            <template #default>
              <keep-alive :include="['ArticlesView']">
                <component :is="Component" :key="route.name === 'articles' ? 'articles-list' : route.fullPath" />
              </keep-alive>
            </template>
            
            <template #fallback>
              <div class="flex items-center justify-center py-20">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-foreground"></div>
              </div>
            </template>
          </Suspense>
        </router-view>
      </main>
    </div>
  </div>
</template>