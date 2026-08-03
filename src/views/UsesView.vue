<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useHead } from '@unhead/vue'
import { getUses, type UsesCategory } from '../services/api'

const categories = ref<UsesCategory[]>([])
const lastUpdated = ref('')
const loading = ref(true)
const error = ref<string | null>(null)

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

onMounted(async () => {
  try {
    const data = await getUses()
    categories.value = data.categories
    lastUpdated.value = data.lastUpdated
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load'
    console.error(e)
  } finally {
    loading.value = false
  }
})

useHead({
  title: 'Uses | Nikita Redko',
  meta: [
    { name: 'description', content: 'Моё рабочее место: железо, софт и инструменты.' },
  ],
})
</script>

<template>
  <div class="min-h-[calc(100vh-7.3rem)] pt-10 pb-20 px-4">
    <div class="max-w-5xl mx-auto">
      <div v-if="loading" class="space-y-8 animate-pulse">
        <div class="h-12 bg-muted/30 rounded w-1/2 mb-4"></div>
        <div class="h-6 bg-muted/30 rounded w-3/4 mb-8"></div>
        <div v-for="i in 3" :key="i" class="space-y-3">
          <div class="h-8 bg-muted/30 rounded w-1/3"></div>
          <div class="grid gap-3 sm:grid-cols-2">
            <div v-for="j in 4" :key="j" class="h-24 bg-muted/20 rounded-xl"></div>
          </div>
        </div>
      </div>
      <div v-else-if="error" class="text-center py-20">
        <p class="text-red-500 mb-2">Не удалось загрузить страницу</p>
        <p class="text-muted text-sm">{{ error }}</p>
        <p class="text-muted/60 text-xs mt-4">
          Проверьте, что в Outline создан документ с заголовком «Uses»
          или прописан <code class="font-mono bg-muted/30 px-1 py-0.5 rounded">USES_DOCUMENT_ID</code> в .env
        </p>
      </div>
      <div v-else-if="categories.length === 0" class="text-center py-20">
        <p class="text-muted">Страница пока пустая</p>
        <p class="text-muted/60 text-sm mt-2">Добавь контент в Outline в формате Markdown</p>
      </div>
      <div v-else>
        <header class="mb-16 animate-fade-in">
          <div class="flex items-center gap-2 text-sm text-muted mb-4">
            <router-link to="/" class="hover:text-foreground transition-colors">Главная</router-link>
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
            <span class="text-foreground">Uses</span>
          </div>
          <h1 class="text-5xl md:text-6xl font-bold tracking-tight mb-4">
            Что я использую
          </h1>
          <p class="text-lg text-muted max-w-2xl leading-relaxed">
            Список железа, софта и инструментов, которые помогают мне кодить каждый день.
            Вдохновлено <a href="https://uses.tech/" target="_blank" rel="noopener noreferrer" class="underline hover:text-foreground transition-colors">uses.tech</a>.
          </p>
          <p v-if="lastUpdated" class="text-sm text-muted/60 mt-4">
            Последнее обновление: {{ formatDate(lastUpdated) }}
          </p>
        </header>
        <div class="space-y-12">
          <section
            v-for="(category, idx) in categories"
            :key="category.id"
            class="animate-fade-in-delayed"
            :style="{ animationDelay: `${idx * 100}ms` }"
          >
            <div class="flex items-start gap-4 mb-6">
              <div class="shrink-0 w-12 h-12 rounded-xl bg-muted/10 border border-border flex items-center justify-center">
                <svg class="w-6 h-6 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-semibold mb-1">{{ category.title }}</h2>
                <p v-if="category.description" class="text-sm text-muted">{{ category.description }}</p>
              </div>
            </div>
            <div class="grid gap-3 sm:grid-cols-2">
              <a
                v-for="item in category.items"
                :key="item.name"
                :href="item.url || '#'"
                :target="item.url ? '_blank' : undefined"
                :rel="item.url ? 'noopener noreferrer' : undefined"
                :class="[
                  'group block p-4 rounded-xl border border-border bg-background transition-all duration-200',
                  item.url
                    ? 'hover:border-foreground/30 hover:bg-muted/5 cursor-pointer'
                    : 'cursor-default'
                ]"
              >
                <div class="flex items-start justify-between gap-2 mb-1">
                  <h3 class="font-medium text-foreground flex items-center gap-1.5">
                    {{ item.name }}
                    <svg
                      v-if="item.url"
                      class="w-3.5 h-3.5 text-muted group-hover:text-foreground group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </h3>
                </div>
                <p class="text-sm text-muted leading-relaxed">{{ item.description }}</p>
              </a>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>