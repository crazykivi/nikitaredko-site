<script setup lang="ts">
import type { Article } from '../services/api'

const props = defineProps<{
  article: Article
  expandedArticles: Set<string>
}>()

const emit = defineEmits<{
  toggle: [id: string]
}>()

const isExpanded = (id: string) => props.expandedArticles.has(id)

const countChildren = (article: Article): number => {
  let count = 0
  if (article.children) {
    count += article.children.length
    for (const child of article.children) {
      count += countChildren(child)
    }
  }
  return count
}

const totalChildren = countChildren(props.article)
</script>

<template>
  <div>
    <div class="flex items-center gap-1 group">
      <button
        v-if="article.children && article.children.length > 0"
        @click="emit('toggle', article.id)"
        class="p-0.5 rounded hover:bg-muted/50 transition-colors text-muted hover:text-foreground"
      >
        <svg
          class="w-3 h-3 transition-transform duration-200"
          :class="{ 'rotate-90': isExpanded(article.id) }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
      <span v-else class="w-4"></span>
      <router-link
        :to="`/articles/${article.id}`"
        class="flex-1 text-left px-2 py-1 text-xs transition-colors truncate rounded hover:bg-muted/30 flex items-center gap-1.5"
        :class="isExpanded(article.id) ? 'text-foreground' : 'text-muted hover:text-foreground'"
        :title="article.title"
        :style="{ paddingLeft: `${4 + (article.level || 0) * 12}px` }"
      >
        <svg
          v-if="!article.children || article.children.length === 0"
          class="w-3 h-3 shrink-0 opacity-40"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <svg
          v-else
          class="w-3 h-3 shrink-0 opacity-60"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>

        <span class="flex-1 truncate">{{ article.title }}</span>
        <span
          v-if="article.children && article.children.length > 0 && !isExpanded(article.id) && totalChildren > 0"
          class="text-[9px] px-1 py-0.5 rounded-full bg-muted/30 text-muted/70 font-mono shrink-0"
        >
          +{{ totalChildren }}
        </span>
      </router-link>
    </div>
    <div
      v-if="article.children && article.children.length > 0 && isExpanded(article.id)"
      class="mt-0.5"
    >
      <ArticleTreeItem
        v-for="child in article.children"
        :key="child.id"
        :article="child"
        :expanded-articles="expandedArticles"
        @toggle="emit('toggle', $event)"
      />
    </div>
  </div>
</template>