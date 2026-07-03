<script setup lang="ts">
import ReadingProgressBar from '../components/ReadingProgressBar.vue'; 
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getArticle, getArticlesStructured } from "../services/api";
import type { Article, CollectionWithArticles } from "../services/api";
import { useHead } from '@unhead/vue';
import MarkdownIt from "markdown-it";
import markdownItContainer from "markdown-it-container";
import hljs from "highlight.js";
import "highlight.js/styles/github-dark.css";

const route = useRoute();
const router = useRouter();

const article = ref<Article | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);
const allCollections = ref<CollectionWithArticles[]>([]);
const flatArticles = ref<Article[]>([]);

let abortController: AbortController | null = null;
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const flattenArticles = (articles: Article[]): Article[] => {
  const result: Article[] = [];
  for (const a of articles) {
    result.push(a);
    if (a.children && a.children.length > 0) {
      result.push(...flattenArticles(a.children));
    }
  }
  return result;
};
const currentCollectionId = computed(() => {
  const queryCollection = route.query.collection as string | undefined;
  if (queryCollection) return queryCollection;
  return article.value?.collectionId;
});
const sortedArticles = computed(() => {
  const collId = currentCollectionId.value;
  let articles: Article[] = [];

  if (!collId) {
    for (const coll of allCollections.value) {
      articles.push(...flattenArticles(coll.articles));
    }
  } else {
    const coll = allCollections.value.find((c) => c.id === collId);
    if (coll) articles = flattenArticles(coll.articles);
  }

  return articles.sort(
		(a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );
});

const currentIndex = computed(() => {
  if (!article.value) return -1;
  return sortedArticles.value.findIndex((a) => a.id === article.value!.id);
});

const prevArticle = computed(() => {
  const idx = currentIndex.value;
  if (idx <= 0) return null;
  return sortedArticles.value[idx - 1];
});

const nextArticle = computed(() => {
  const idx = currentIndex.value;
  if (idx === -1 || idx >= sortedArticles.value.length - 1) return null;
  return sortedArticles.value[idx + 1];
});

const goToArticle = (id: string) => {
  router.push(`/articles/${id}`);
};

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: (str: string, lang: string) => {
    const language = lang && hljs.getLanguage(lang) ? lang : "plaintext";
    let highlighted: string;

    try {
      highlighted = hljs.highlight(str, { language }).value;
    } catch {
      highlighted = hljs.highlightAuto(str).value;
    }

    const safeLang = language.replace(/"/g, "&quot;");

    return (
      `<div class="code-block-wrapper relative group" data-language="${safeLang}">` +
      `<div class="code-header flex items-center justify-between px-4 py-2 bg-muted/80 border-b border-border rounded-t-lg">` +
      `<span class="text-xs text-muted font-mono uppercase">${safeLang}</span>` +
      `<button ` +
      `class="copy-code-btn text-xs px-3 py-1 rounded-md bg-muted hover:bg-muted/80 text-muted-foreground transition-all opacity-0 group-hover:opacity-100 focus:opacity-100 flex items-center gap-1"` +
      `onclick="window.__copyCode(this)"` +
      `title="Копировать код">` +
      `<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">` +
      `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />` +
      `</svg>` +
      `<span>Копировать</span>` +
      `</button>` +
      `</div>` +
      `<pre class="!mt-0 !mb-0 !rounded-t-none"><code class="hljs language-${safeLang}">${highlighted}</code></pre>` +
      `</div>`
    );
  },
});

const blockTypes = ["warning", "info", "success", "danger", "tip", "note"];

blockTypes.forEach((blockType) => {
  md.use(markdownItContainer, blockType, {
    validate: (params: string) => {
      return params.trim() === blockType;
    },
    render: (tokens: any[], idx: number) => {
      if (tokens[idx].nesting === 1) {
        const title = blockType.charAt(0).toUpperCase() + blockType.slice(1);
        return `<div class="custom-block ${blockType}">
  <p class="custom-block-title">${title}</p>
`;
      } else {
        return "</div>\n";
      }
    },
  });
});

if (typeof window !== "undefined") {
  (window as any).__copyCode = (button: HTMLButtonElement) => {
    const codeBlock = button.closest(".code-block-wrapper");
    if (!codeBlock) return;

    const codeElement = codeBlock.querySelector("code");
    if (!codeElement) return;

    const code = codeElement.textContent || "";

    navigator.clipboard
      .writeText(code)
      .then(() => {
        const originalHTML = button.innerHTML;
        button.innerHTML = `
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span>Скопировано!</span>
      `;
        button.classList.add("bg-green-500/20", "text-green-600");

        setTimeout(() => {
          button.innerHTML = originalHTML;
          button.classList.remove("bg-green-500/20", "text-green-600");
        }, 2000);
      })
      .catch((err) => {
        console.error("Failed to copy:", err);
      });
  };
}

const loadArticle = async (id: string) => {
  if (abortController) abortController.abort();
  abortController = new AbortController();
  const signal = abortController.signal;

  loading.value = true;
  error.value = null;
  article.value = null;

  try {
    const [data, collections] = await Promise.all([
      getArticle(id, signal),
      getArticlesStructured(signal),
    ]);

    if (signal.aborted) return;

    article.value = data;
    allCollections.value = collections;
    const collId = (route.query.collection as string | undefined) || data.collectionId;
    let articles: Article[] = [];
    if (!collId) {
      for (const coll of collections) {
        articles.push(...flattenArticles(coll.articles));
      }
    } else {
      const coll = collections.find((c) => c.id === collId);
      if (coll) articles = flattenArticles(coll.articles);
    }
    flatArticles.value = articles.sort(
      (a, b) => new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime()
    );

    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (e: any) {
    if (e.name === "AbortError") return;
    error.value = e instanceof Error ? e.message : "Failed to load article";
    console.error("Failed to load article:", e);
  } finally {
    if (!signal.aborted) loading.value = false;
  }
};

const goBack = () => {
  const returnUrl = sessionStorage.getItem("return_url_articles");
  if (returnUrl) {
    router.push(returnUrl);
  } else {
    router.push("/articles");
  }
};

const goToCollection = (collectionId: string) => {
  router.push({ path: "/articles", query: { collection: collectionId } });
};

const siteName = 'Nikita Redko'
const defaultOgImage = 'https://www.nikitaredko.ru/favicon.svg' 

const seoData = computed(() => {
  const currentUrl = window.location.href
  
  if (!article.value) {
    return { title: siteName }
  }

  const title = `${article.value.title} | ${siteName}`
  
  const rawDesc = article.value.excerpt || article.value.content || ''
  const description = rawDesc.replace(/<[^>]*>?/gm, '').substring(0, 150).trim() || 'Статья на сайте Никиты Редко'

  return {
    title,
    meta: [
      { name: 'description', content: description },
      
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
      { property: 'og:type', content: 'article' },
      { property: 'og:url', content: currentUrl },
      { property: 'og:image', content: defaultOgImage },
      { property: 'og:site_name', content: siteName },
      
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: title },
      { name: 'twitter:description', content: description },
      { name: 'twitter:image', content: defaultOgImage },
    ],
    link: [
      { rel: 'canonical' as const, href: currentUrl }
    ]
  }
})

useHead(seoData)

watch(
  () => route.params.id,
  (newId, oldId) => {
    if (newId && newId !== oldId) loadArticle(newId as string);
  }
);

onMounted(async () => {
  await loadArticle(route.params.id as string);
});

onUnmounted(() => {
  if (abortController) abortController.abort();
});
</script>

<template>
  <ReadingProgressBar />
  <div class="flex-1 min-w-0">
    <button
      @click="goBack"
      class="mb-8 text-muted hover:text-foreground transition-colors flex items-center gap-2"
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
        <div class="flex items-center gap-4 text-sm text-muted">
          <time>{{ formatDate(article.createdAt) }}</time>
          <span>•</span>
          <span>{{ article.readTime }} мин чтения</span>
        </div>
      </div>

      <div
        class="prose prose-neutral dark:prose-invert max-w-none break-words article-content"
      >
        <div v-html="md.render(article.content || '')"></div>
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
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>