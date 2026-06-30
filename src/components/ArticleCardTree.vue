<script setup lang="ts">
import { ref } from 'vue'
import type { Article } from '../services/api'
import ArticleCard from './ArticleCard.vue'

defineProps<{
  article: Article
}>()

const isExpanded = ref(true)

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
}
</script>

<template>
  <div>
    <div :style="{ marginLeft: `${(article.level || 0) * 24}px` }">
      <div class="flex items-start gap-2 mb-4">
        <button
          v-if="article.children && article.children.length > 0"
          @click="toggleExpand"
          class="mt-2 p-1 rounded hover:bg-muted/50 transition-colors text-muted hover:text-foreground shrink-0"
        >
          <svg
            class="w-4 h-4 transition-transform duration-200"
            :class="{ 'rotate-90': isExpanded }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
        <span v-else class="w-6 shrink-0"></span>
        <div class="flex-1">
          <ArticleCard :article="article" />
        </div>
      </div>
    </div>
    <div
      v-if="article.children && article.children.length > 0 && isExpanded"
      class="mt-4 space-y-4"
    >
      <ArticleCardTree
        v-for="child in article.children"
        :key="child.id"
        :article="child"
      />
    </div>
  </div>
</template>