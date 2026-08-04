<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useHead } from '@unhead/vue'
import { getAbout, type AboutResponse, type StackItem } from '../services/api'

const data = ref<AboutResponse | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const active = ref<StackItem | null>(null)

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const CAREER_TYPES = {
  work:       { label: 'Работа',      dot: 'border-sky-400',     chip: 'text-sky-400 border-sky-400/30 bg-sky-400/10' },
  pet:        { label: 'Пет-проект',  dot: 'border-emerald-400', chip: 'text-emerald-400 border-emerald-400/30 bg-emerald-400/10' },
  education:  { label: 'Учёба',       dot: 'border-amber-400',   chip: 'text-amber-400 border-amber-400/30 bg-amber-400/10' },
  opensource: { label: 'Open Source', dot: 'border-violet-400',  chip: 'text-violet-400 border-violet-400/30 bg-violet-400/10' },
  freelance:  { label: 'Фриланс',     dot: 'border-rose-400',    chip: 'text-rose-400 border-rose-400/30 bg-rose-400/10' },
} as const

type CareerType = keyof typeof CAREER_TYPES

const stageType = (stage: AboutResponse['career'][number]): CareerType => {
  const t = (stage as { type?: string }).type
  if (t && t in CAREER_TYPES) return t as CareerType
  const text = `${stage.role} ${stage.company}`.toLowerCase()
  if (/пет|pet|свои проекты/.test(text)) return 'pet'
  if (/учёб|учеб|курс|универ/.test(text)) return 'education'
  return 'work'
}

const career = computed(() =>
  (data.value?.career ?? []).map(stage => {
    const type = stageType(stage)
    return { ...stage, meta: CAREER_TYPES[type] }
  }),
)

const legend = computed(() =>
  [...new Set(career.value.map(s => stageType(s)))].map(t => ({ type: t, ...CAREER_TYPES[t] })),
)

onMounted(async () => {
  try {
    data.value = await getAbout()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load'
    console.error(e)
  } finally {
    loading.value = false
  }
})

useHead({
  title: 'Обо мне | Nikita Redko',
  meta: [
    { name: 'description', content: 'Кто я такой: карьерный путь, технологии и факты.' },
  ],
})
</script>

<template>
  <div class="min-h-[calc(100vh-7.3rem)] pt-10 pb-20 px-4">
    <div class="max-w-5xl mx-auto">
      <div class="flex items-center gap-2 text-sm text-muted mb-4">
        <router-link to="/" class="hover:text-foreground transition-colors"
          >Главная</router-link
        >
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 5l7 7-7 7"
          />
        </svg>
        <span class="text-foreground">Обо мне</span>
      </div>

      <div v-if="loading" class="space-y-8 animate-pulse">
        <div class="h-12 bg-muted/30 rounded w-1/3 mb-4"></div>
        <div class="h-6 bg-muted/30 rounded w-2/3 mb-8"></div>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
          <div v-for="i in 4" :key="i" class="h-24 bg-muted/20 rounded-xl"></div>
        </div>
        <div class="h-8 bg-muted/30 rounded w-1/4 mb-4"></div>
        <div class="space-y-4">
          <div v-for="i in 3" :key="i" class="h-32 bg-muted/10 rounded-xl"></div>
        </div>
      </div>

      <div v-else-if="error" class="text-center py-20">
        <p class="text-red-500 mb-2">Не удалось загрузить страницу</p>
        <p class="text-muted text-sm">{{ error }}</p>
        <p class="text-muted/60 text-xs mt-4">
          Проверьте, что в Outline создан документ с заголовком «About» или «Обо мне» или
          прописан
          <code class="font-mono bg-muted/30 px-1 py-0.5 rounded">ABOUT_DOCUMENT_ID</code>
          в .env
        </p>
      </div>
      <div v-else-if="data">
        <header class="mb-14 animate-fade-in">
          <p class="font-mono text-xs uppercase tracking-widest text-muted mb-4">
            // обо мне
          </p>
          <h1 class="text-4xl md:text-5xl font-bold tracking-tight mb-6">
            Привет, я Никита
          </h1>
          <div
            v-if="data.intro"
            class="text-lg text-muted leading-relaxed max-w-3xl space-y-4"
          >
            <p v-for="(p, i) in data.intro.split('\n').filter((l) => l.trim())" :key="i">
              {{ p }}
            </p>
          </div>
          <p v-if="data.lastUpdated" class="text-sm text-muted/60 mt-8">
            Последнее обновление: {{ formatDate(data.lastUpdated) }}
          </p>
        </header>
        <section
          v-if="data.facts && data.facts.length > 0"
          class="mb-16 animate-fade-in-delayed"
        >
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div
              v-for="fact in data.facts"
              :key="fact.label"
              class="p-5 rounded-xl border border-border bg-background hover:bg-muted/5 transition-colors"
            >
              <div class="text-2xl md:text-3xl font-bold text-foreground mb-1">
                {{ fact.value }}
              </div>
              <div class="text-sm text-muted leading-snug">{{ fact.label }}</div>
            </div>
          </div>
        </section>
        <section v-if="career.length > 0" class="mb-16 animate-fade-in-delayed-2">
          <h2
            class="text-2xl md:text-3xl font-semibold tracking-tight mb-8 flex items-center gap-3"
          >
            <span class="font-mono text-sm text-muted">01</span>
            Карьера и опыт
            <div class="flex-1 border-t border-border"></div>
          </h2>
          <div
            class="flex flex-wrap items-center gap-x-5 gap-y-2 mb-8 text-xs text-muted"
          >
            <div v-for="item in legend" :key="item.type" class="flex items-center gap-2">
              <span
                class="w-3 h-3 rounded-full bg-background border-2"
                :class="item.dot"
              ></span>
              {{ item.label }}
            </div>
          </div>
          <div class="relative">
            <div
              class="absolute left-[7px] top-2 bottom-2 w-0.5 bg-border rounded-full"
            ></div>
            <div class="space-y-10">
              <div
                v-for="stage in career"
                :key="stage.period + stage.company"
                class="relative pl-8 group"
              >
                <div
                  class="absolute left-0 top-1.5 w-4 h-4 rounded-full bg-background border-2 group-hover:scale-110 transition-transform"
                  :class="stage.meta.dot"
                  :title="stage.meta.label"
                ></div>
                <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1 mb-3">
                  <h3 class="text-lg font-semibold text-foreground">{{ stage.role }}</h3>
                  <span class="text-muted text-sm font-medium">{{ stage.company }}</span>
                  <span
                    class="text-[11px] font-mono px-2 py-0.5 rounded-full border"
                    :class="stage.meta.chip"
                  >
                    {{ stage.meta.label }}
                  </span>
                  <span
                    class="ml-auto font-mono text-xs text-muted/80 bg-muted/20 px-2.5 py-1 rounded-full flex items-center gap-1.5"
                  >
                    {{ stage.period }}
                    <span
                      v-if="stage.current"
                      class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"
                    ></span>
                  </span>
                </div>
                <p
                  v-if="stage.description"
                  class="text-muted text-sm mb-4 leading-relaxed max-w-3xl"
                >
                  {{ stage.description }}
                </p>

                <ul
                  v-if="stage.highlights && stage.highlights.length > 0"
                  class="space-y-2 text-sm text-muted"
                >
                  <li
                    v-for="(hl, i) in stage.highlights"
                    :key="i"
                    class="flex items-start gap-2.5"
                  >
                    <svg
                      class="w-4 h-4 text-foreground/40 shrink-0 mt-0.5"
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
                    <span class="leading-relaxed">{{ hl }}</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </section>
        <section
        v-if="data.stack && data.stack.length > 0"
        class="mb-16 animate-fade-in-delayed-3"
        >
        <h2 class="text-2xl md:text-3xl font-semibold tracking-tight mb-4 flex items-center gap-3">
            <span class="font-mono text-sm text-muted">02</span>
            Технологический стек
            <div class="flex-1 border-t border-border"></div>
        </h2>
        <div class="mb-8 min-h-[3rem] md:min-h-[1.5rem] font-mono text-sm flex items-start gap-2">
            <span class="text-foreground/40 select-none shrink-0">{{ active ? '>' : '//' }}</span>
            <transition name="fade" mode="out-in">
            <p v-if="active" :key="active.name" class="text-muted leading-relaxed">
                <span class="text-foreground">{{ active.name }}</span>
                — {{ active.description }}
                <span class="animate-pulse text-foreground/40">▍</span>
            </p>
            <p v-else key="hint" class="text-muted/60 leading-relaxed">
                наведи на тег — расскажу, что с ним делал
            </p>
            </transition>
        </div>
        <div class="space-y-8" @mouseleave="active = null">
            <div
            v-for="group in data.stack"
            :key="group.id"
            class="grid gap-3 md:grid-cols-[170px_1fr] md:gap-8"
            >
            <div class="flex items-baseline gap-2 md:pt-2">
                <span class="font-mono text-sm text-foreground/30 select-none">//</span>
                <h3 class="font-mono text-sm text-muted">{{ group.title }}</h3>
                <span class="font-mono text-[10px] text-muted/50">{{ group.items.length }}</span>
            </div>

            <div class="flex flex-wrap gap-2">
                <template v-for="item in group.items" :key="item.name">
                <a
                    v-if="item.url"
                    :href="item.url"
                    target="_blank"
                    rel="noopener noreferrer"
                    @mouseenter="active = item"
                    @focus="active = item"
                    class="group px-3 py-1.5 rounded-lg border text-sm transition-all duration-150 hover:-translate-y-0.5 flex items-center gap-1.5"
                    :class="
                    active?.name === item.name
                        ? 'border-foreground/40 bg-muted/15 text-foreground'
                        : 'border-border bg-muted/5 text-foreground/80 hover:border-foreground/30 hover:text-foreground'
                    "
                >
                    {{ item.name }}
                    <svg
                    class="w-3 h-3 text-muted/60 group-hover:text-muted transition-colors"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                    />
                    </svg>
                </a>
                <button
                    v-else
                    @mouseenter="active = item"
                    @focus="active = item"
                    @click="active = item"
                    class="px-3 py-1.5 rounded-lg border text-sm transition-all duration-150 hover:-translate-y-0.5 focus:outline-none"
                    :class="
                    active?.name === item.name
                        ? 'border-foreground/40 bg-muted/15 text-foreground'
                        : 'border-border bg-muted/5 text-foreground/80 hover:border-foreground/30 hover:text-foreground'
                    "
                >
                    {{ item.name }}
                </button>
                </template>
            </div>
            </div>
        </div>
        </section>
        <div
          class="mt-16 pt-8 border-t border-border flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"
        >
          <p class="text-sm text-muted leading-relaxed">
            Железо, софт и рабочее окружение — на странице
            <router-link
              to="/uses"
              class="text-foreground underline underline-offset-4 hover:opacity-70 transition-opacity"
            >
              /uses </router-link
            >.
          </p>
          <router-link
            to="/articles"
            class="group inline-flex items-center gap-2 text-sm font-medium hover:opacity-70 transition-opacity shrink-0"
          >
            Читать статьи
            <svg
              class="w-4 h-4 transition-transform group-hover:translate-x-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 7l5 5m0 0l-5 5m5-5H6"
              />
            </svg>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
