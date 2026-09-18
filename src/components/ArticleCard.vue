<script setup lang="ts">
import { computed } from 'vue'
import type { Article } from '../services/api'
import { renderExcerpt } from '../utils/markdown'
import { formatDateNumeric } from '../utils/date'

const props = withDefaults(
  defineProps<{
    article: Article
    accent?: string | null
    number?: number
  }>(),
  { accent: null, number: 0 },
)

const MAX_TAGS = 3
const tags = computed(() => props.article.tags ?? [])
const excerptHtml = computed(() => renderExcerpt(props.article.excerpt))
const visibleTags = computed(() => tags.value.slice(0, MAX_TAGS))
</script>

<template>
  <article class="group relative grid gap-4 py-7 first:pt-0 md:grid-cols-[88px_1fr] md:gap-8">
    <div class="flex items-baseline gap-3 font-mono md:flex-col md:items-start md:gap-1.5 md:pt-1">
      <time class="text-[11px] text-muted" :datetime="article.createdAt">
        {{ formatDateNumeric(article.createdAt) }}
      </time>
    </div>

    <div class="min-w-0">
      <div class="flex items-start justify-between gap-4">
        <h3
          class="text-xl font-semibold leading-snug tracking-tight text-foreground [text-wrap:balance] transition-transform duration-300 group-hover:translate-x-1 md:text-[22px]"
        >
          {{ article.title }}
        </h3>
        <svg
          class="mt-2 h-4 w-4 shrink-0 text-border transition-all duration-300 group-hover:translate-x-1 group-hover:text-foreground"
          fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
        </svg>
      </div>
      <p
        class="mt-2.5 max-w-2xl text-[15px] leading-relaxed text-muted line-clamp-2 md:text-base"
        v-html="excerptHtml"
      />
      <div class="mt-4 flex flex-wrap items-center gap-x-3 gap-y-1.5 font-mono text-[11px] text-muted">
        <span v-if="article.collectionName" class="inline-flex items-center gap-1.5 uppercase tracking-wider">
          <span class="h-1.5 w-1.5 rounded-full" :style="{ backgroundColor: accent ?? 'currentColor' }" aria-hidden="true" />
          {{ article.collectionName }}
        </span>
        <span aria-hidden="true">·</span>
        <span>{{ article.readTime }} мин чтения</span>
        <template v-if="visibleTags.length">
          <span aria-hidden="true">·</span>
          <span v-for="tag in visibleTags" :key="tag">#{{ tag }}</span>
        </template>
      </div>
    </div>
    <router-link
      :to="`/articles/${article.id}`"
      :aria-label="article.title"
      class="absolute inset-0 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-foreground"
    />
  </article>
</template>