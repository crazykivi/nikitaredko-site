<script setup lang="ts">
import { onMounted, nextTick, ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useHead } from "@unhead/vue";
import DOMPurify from "dompurify";
import type { Config as DOMPurifyConfig } from "dompurify";
import Giscus from "@giscus/vue";

import { useArticleLoader } from "../composables/useArticleLoader";
import { useArticleTOC } from "../composables/useArticleTOC";
import { useCopyButtons } from "../composables/useCopyButtons";
import { useTheme } from "../composables/useTheme";
import { createMarkdownRenderer } from "../utils/markdownRenderer";
import { buildArticleSeo } from "../utils/articleSeo";

import ReadingProgressBar from "../components/ReadingProgressBar.vue";
import FloatingTOC from "../components/FloatingTOC.vue";

const route = useRoute();
const router = useRouter();
const { effectiveDark } = useTheme();

const md = createMarkdownRenderer();

const purifyConfig: DOMPurifyConfig = {
  USE_PROFILES: { html: true, svg: true },
  FORBID_TAGS: ["style", "form", "input", "button", "textarea", "select"],
  FORBID_ATTR: ["onerror", "onload", "onclick", "onmouseover"],
  ADD_ATTR: ["allow", "allowfullscreen", "frameborder", "scrolling"],
};

const {
  article,
  loading,
  error,
  prevArticle,
  nextArticle,
  loadArticle,
} = useArticleLoader();

const articleContent = computed(() => article.value?.content);
const { tocItems } = useArticleTOC(articleContent);

const { copyToast, injectCopyButtons } = useCopyButtons();

const activeHeadingId = ref("");

const giscusTheme = computed(() => (effectiveDark.value ? "dark" : "light"));

const renderMarkdown = (content: string): string =>
  (DOMPurify.sanitize(md.render(content), purifyConfig) as unknown) as string;

useHead(() => buildArticleSeo(article.value));

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const goBack = () => {
  const returnUrl = sessionStorage.getItem("return_url_articles");
  router.push(returnUrl ?? "/articles");
};

const goToCollection = (id: string) =>
  router.push({ path: "/articles", query: { collection: id } });

const goToArticle = (id: string) => router.push(`/articles/${id}`);

onMounted(async () => {
  await loadArticle(route.params.id as string);
  nextTick(() => injectCopyButtons());
});
</script>

<template>
  <ReadingProgressBar />
  <div class="flex-1 min-w-0">
    <button
      @click="goBack"
      class="mb-8 text-muted hover:text-foreground transition-colors flex items-center gap-2 back-btn"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 19l-7-7m0 0l7-7m-7 7h18"
        />
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
      <p class="text-muted mt-2">
        Не удалось загрузить данные. Попробуйте обновить страницу.
      </p>
    </div>

    <article v-else-if="article" class="animate-fade-in max-w-6xl mx-auto">
      <div class="mb-8">
        <button
          v-if="article.collectionName && article.collectionId"
          @click="goToCollection(article.collectionId)"
          class="inline-flex items-center gap-1.5 mb-4 text-sm font-bold text-foreground hover:opacity-70 transition-opacity group back-btn"
          :title="`Все статьи из категории «${article.collectionName}»`"
        >
          <svg
            class="w-3.5 h-3.5 text-muted group-hover:text-foreground transition-colors"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
            />
          </svg>
          {{ article.collectionName }}
          <svg
            class="w-3 h-3 text-muted opacity-0 group-hover:opacity-100 transition-opacity"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 5l7 7-7 7"
            />
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
        <div class="article-meta flex items-center gap-4 text-sm text-muted">
          <time>{{ formatDate(article.createdAt) }}</time>
          <span>•</span>
          <span>{{ article.readTime }} мин чтения</span>
        </div>
      </div>

      <div
        class="prose prose-neutral dark:prose-invert max-w-none break-words article-content"
      >
        <div v-html="renderMarkdown(article.content || '')"></div>
      </div>

      <div class="mt-16 pt-8 border-t border-border">
        <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
            />
          </svg>
          Обсуждение
        </h2>
        <Giscus
          repo="crazykivi/nikitaredko-site"
          repo-id="R_kgDOTJQbZA"
          category="Show and tell"
          category-id="DIC_kwDOTJQbZM4DCjCb"
          mapping="specific"
          :term="article.id"
          reactions-enabled="1"
          emit-metadata="0"
          input-position="top"
          :theme="giscusTheme"
          lang="ru"
          loading="lazy"
        />
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
          <div
            class="shrink-0 w-10 h-10 rounded-full bg-muted/30 group-hover:bg-muted/60 flex items-center justify-center transition-colors"
          >
            <svg
              class="w-5 h-5 text-muted group-hover:text-foreground transition-colors"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-xs text-muted uppercase tracking-wider mb-1">Предыдущая</div>
            <div
              class="font-medium text-foreground truncate group-hover:translate-x-[-2px] transition-transform"
            >
              {{ prevArticle.title }}
            </div>
            <div class="text-xs text-muted mt-1">
              {{ formatDate(prevArticle.createdAt) }}
            </div>
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
            <div
              class="font-medium text-foreground truncate group-hover:translate-x-[2px] transition-transform"
            >
              {{ nextArticle.title }}
            </div>
            <div class="text-xs text-muted mt-1">
              {{ formatDate(nextArticle.createdAt) }}
            </div>
          </div>
          <div
            class="shrink-0 w-10 h-10 rounded-full bg-muted/30 group-hover:bg-muted/60 flex items-center justify-center transition-colors"
          >
            <svg
              class="w-5 h-5 text-muted group-hover:text-foreground transition-colors"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              />
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
  <FloatingTOC
    v-if="article && tocItems.length > 1"
    :items="tocItems"
    :active-id="activeHeadingId"
  />
  <Transition name="toast">
    <div
      v-if="copyToast"
      class="fixed top-12 inset-x-0 z-[9999] flex justify-center pointer-events-none"
      role="status"
    >
      <div
        class="flex items-center gap-2.5 px-4 py-2.5 rounded-xl border border-border bg-background/95 backdrop-blur-sm shadow-2xl text-sm text-foreground"
      >
        <svg
          v-if="copyToast.type === 'success'"
          class="w-4 h-4 text-green-500 shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <svg
          v-else
          class="w-4 h-4 text-red-500 shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        {{
          copyToast.type === "success"
            ? "Копирование в буфер обмена прошло успешно!"
            : "Не удалось скопировать код"
        }}
      </div>
    </div>
  </Transition>
</template>
