<script setup lang="ts">
import type { Article } from '../services/api'

defineProps<{
  article: Article
}>()

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<template>
  <router-link 
    :to="`/articles/${article.id}`"
    class="group block p-6 rounded-xl border border-border hover:border-foreground/20 bg-background hover:bg-muted/5 transition-all duration-300"
  >
    <div v-if="article.collectionName" class="mb-3">
      <span class="text-xs font-medium text-muted/70 uppercase tracking-wider">
        {{ article.collectionName }}
      </span>
    </div>

    <div class="flex flex-wrap gap-2 mb-3">
      <span 
        v-for="tag in article.tags" 
        :key="tag"
        class="text-xs px-2 py-1 rounded-md bg-muted/50 text-muted font-mono"
      >
        {{ tag }}
      </span>
    </div>
    
    <h3 class="text-2xl font-semibold mb-2 group-hover:text-foreground transition-colors">
      {{ article.title }}
    </h3>
    
    <p class="text-muted mb-4 line-clamp-2">
      {{ article.excerpt }}
    </p>
    
    <div class="flex items-center gap-4 text-sm text-muted">
      <time>{{ formatDate(article.publishedAt) }}</time>
      <span>•</span>
      <span>{{ article.readTime }} мин чтения</span>
    </div>
  </router-link>
</template>