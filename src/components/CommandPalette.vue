<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { getArticlesFeed, type Article } from '../services/api'

const router = useRouter()

const isOpen = ref(false)
const query = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)
const listRef = ref<HTMLElement | null>(null)
const paletteRef = ref<HTMLElement | null>(null)

let previouslyFocused: HTMLElement | null = null

const restoreFocus = () => {
  const el = previouslyFocused
  previouslyFocused = null
  if (el && el.isConnected) {
    el.focus({ preventScroll: true })
    return
  }
  if (document.activeElement instanceof HTMLElement) {
    document.activeElement.blur()
  }
}

const onGlobalFocusIn = (e: FocusEvent) => {
  if (!isOpen.value) return
  if (e.target instanceof HTMLElement && paletteRef.value?.contains(e.target)) return
  inputRef.value?.focus({ preventScroll: true })
}

const recentArticles = ref<Article[]>([])
const articlesLoaded = ref(false)
const articlesLoading = ref(false)
let articlesAbort: AbortController | null = null

let lastHoveredId: string | null = null
const FILL_ICONS = new Set<Command['icon']>(['github', 'rss'])

type Command = {
  id: string
  label: string
  description?: string
  group: string
  icon: 'home' | 'article' | 'user' | 'tool' | 'theme' | 'github' | 'rss'
  shortcut?: string[]
  action: () => void
}

const ICONS: Record<Command['icon'], string> = {
  home: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
  article: 'M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2',
  user: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
  tool: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
  theme: 'M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z',
  github: 'M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z',
rss: 'M19.199 24C19.199 13.467 10.533 4.8 0 4.8V0c13.165 0 24 10.835 24 24h-4.801zM3.291 17.415c1.814 0 3.293 1.479 3.293 3.295 0 1.813-1.485 3.29-3.301 3.29C1.47 24 0 22.526 0 20.71s1.475-3.294 3.291-3.295zM15.909 24h-4.665c0-6.169-5.075-11.245-11.244-11.245V8.09c8.727 0 15.909 7.184 15.909 15.91z',}

const buildStaticCommands = (): Command[] => [
  {
    id: 'nav-home',
    label: 'Главная',
    group: 'Навигация',
    icon: 'home',
    shortcut: ['⌘', '1'],
    action: () => router.push('/'),
  },
  {
    id: 'nav-articles',
    label: 'Статьи',
    group: 'Навигация',
    icon: 'article',
    shortcut: ['⌘', '2'],
    action: () => router.push('/articles'),
  },
  {
    id: 'nav-about',
    label: 'О себе',
    group: 'Навигация',
    icon: 'user',
    shortcut: ['⌘', '3'],
    action: () => router.push('/about'),
  },
  {
    id: 'nav-uses',
    label: 'Uses',
    group: 'Навигация',
    icon: 'tool',
    shortcut: ['⌘', '4'],
    action: () => router.push('/uses'),
  },
  {
    id: 'action-theme',
    label: 'Сменить тему',
    description: 'Переключить светлую/тёмную тему/системная',
    group: 'Действия',
    icon: 'theme',
    action: () => window.dispatchEvent(new CustomEvent('toggle-app-theme')),
  },
  {
    id: 'action-github',
    label: 'GitHub репозиторий',
    description: 'crazykivi/nikitaredko-site',
    group: 'Действия',
    icon: 'github',
    action: () => window.open('https://github.com/crazykivi/nikitaredko-site', '_blank', 'noopener,noreferrer'),
  },
  {
    id: 'action-rss',
    label: 'RSS-лента',
    description: 'Подписаться на обновления',
    group: 'Действия',
    icon: 'rss',
    action: () => window.open('/api/rss.xml', '_blank', 'noopener,noreferrer'),
  },
]

const commands = computed<Command[]>(() => {
  const staticCmds = buildStaticCommands()
  const articleCmds: Command[] = recentArticles.value.map((a) => ({
    id: `article-${a.id}`,
    label: a.title,
    description: a.collectionName,
    group: 'Статьи',
    icon: 'article',
    action: () => router.push(`/articles/${a.id}`),
  }))
  return [...staticCmds, ...articleCmds]
})

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return commands.value
  return commands.value.filter(
    (c) =>
      c.label.toLowerCase().includes(q) ||
      (c.description || '').toLowerCase().includes(q),
  )
})

const grouped = computed(() => {
  const map = new Map<string, Command[]>()
  for (const cmd of filtered.value) {
    if (!map.has(cmd.group)) map.set(cmd.group, [])
    map.get(cmd.group)!.push(cmd)
  }
  return Array.from(map.entries())
})

const flatFiltered = computed(() => filtered.value)

watch(flatFiltered, () => {
  selectedIndex.value = 0
})

const open = () => {
  previouslyFocused =
    document.activeElement instanceof HTMLElement && document.activeElement !== document.body
      ? document.activeElement
      : null
  isOpen.value = true
  query.value = ''
  selectedIndex.value = 0
  lastHoveredId = null
  nextTick(() => inputRef.value?.focus())
  if (!articlesLoaded.value && !articlesLoading.value) loadArticles()
}

const close = () => {
  isOpen.value = false
  lastHoveredId = null
  if (articlesAbort) {
    articlesAbort.abort()
    articlesAbort = null
  }
  restoreFocus()
}

const execute = (cmd: Command) => {
  close()
  try {
    cmd.action()
  } catch (e) {
    console.error('[CommandPalette] action failed:', e)
  }
}

const scrollToSelected = () => {
  nextTick(() => {
    const el = listRef.value?.querySelector('[data-selected="true"]') as HTMLElement | null
    el?.scrollIntoView({ block: 'nearest' })
  })
}

const onInputKeydown = (e: KeyboardEvent) => {
  const len = flatFiltered.value.length

  if (e.key === 'Escape') {
    e.preventDefault()
    close()
    return
  }

  if (len === 0) return

  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value + 1) % len
      scrollToSelected()
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value - 1 + len) % len
      scrollToSelected()
      break
    case 'Tab':
      e.preventDefault()
      if (e.shiftKey) {
        selectedIndex.value = (selectedIndex.value - 1 + len) % len
      } else {
        selectedIndex.value = (selectedIndex.value + 1) % len
      }
      scrollToSelected()
      break
    case 'Home':
      e.preventDefault()
      selectedIndex.value = 0
      scrollToSelected()
      break
    case 'End':
      e.preventDefault()
      selectedIndex.value = len - 1
      scrollToSelected()
      break
    case 'Enter':
      e.preventDefault()
      {
        const cmd = flatFiltered.value[selectedIndex.value]
        if (cmd) execute(cmd)
      }
      break
  }
}

const onListMouseMove = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  const button = target.closest('[data-cmd-id]') as HTMLElement | null
  const id = button?.getAttribute('data-cmd-id') ?? null
  if (id === lastHoveredId) return
  lastHoveredId = id
  if (!id) return
  const idx = flatFiltered.value.findIndex((c) => c.id === id)
  if (idx >= 0 && idx !== selectedIndex.value) {
    selectedIndex.value = idx
  }
}

const onListMouseLeave = () => {
  lastHoveredId = null
}

const isTypingField = (target: EventTarget | null): boolean => {
  if (!(target instanceof HTMLElement)) return false
  const tag = target.tagName.toLowerCase()
  if (tag === 'input' || tag === 'textarea' || tag === 'select') return true
  if (target.isContentEditable) return true
  return false
}

const handleGlobalKeydown = (e: KeyboardEvent) => {
  const isMac = navigator.platform.toUpperCase().includes('MAC')
  const mod = isMac ? e.metaKey : e.ctrlKey

  // Открытие панели: / (вне печати)
  if (e.key === '/' && !mod && !e.altKey && !e.shiftKey) {
    if (isTypingField(e.target)) return
    e.preventDefault()
    isOpen.value ? close() : open()
    return
  }

  // Toggle: command+shift+K / Ctrl+Shift+K
  if (mod && e.shiftKey && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    isOpen.value ? close() : open()
    return
  }

  // Быстрые переходы command+1..4
  if (!isOpen.value && mod && !e.shiftKey && /^[1-4]$/.test(e.key)) {
    if (isTypingField(e.target)) return
    e.preventDefault()
    const idx = parseInt(e.key, 10) - 1
    const staticCmds = buildStaticCommands()
    if (staticCmds[idx]) {
      try { staticCmds[idx].action() } catch (err) { console.error(err) }
    }
    return
  }

  // Выход из панели (ESC)
  if (isOpen.value && e.key === 'Escape') {
    e.preventDefault()
    close()
  }
}

const loadArticles = async () => {
  articlesLoading.value = true
  articlesAbort = new AbortController()
  try {
    const feed = await getArticlesFeed(1, 20, undefined, articlesAbort.signal)
    recentArticles.value = feed.articles
    articlesLoaded.value = true
  } catch (e) {
    if (e instanceof Error && e.name === 'AbortError') return
    console.error('[CommandPalette] Failed to load articles:', e)
  } finally {
    articlesLoading.value = false
    articlesAbort = null
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
  document.addEventListener('focusin', onGlobalFocusIn)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
  document.removeEventListener('focusin', onGlobalFocusIn)
  if (articlesAbort) articlesAbort.abort()
})
</script>

<template>
  <Transition name="palette">
    <div
      v-if="isOpen"
      ref="paletteRef"
      class="fixed inset-0 z-[9998] flex items-start justify-center pt-[15vh] px-4"
      @click.self="close"
    >
      <div class="absolute inset-0 bg-background/70 backdrop-blur-sm" />
      <div
        class="relative w-full max-w-xl rounded-2xl border border-border bg-background shadow-2xl overflow-hidden animate-slide-down"
        role="dialog"
        aria-modal="true"
        aria-label="Command palette"
      >
        <div class="flex items-center gap-3 px-5 py-4 border-b border-border">
          <svg
            class="w-5 h-5 text-muted shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Поиск страниц, статей, действий..."
            class="flex-1 bg-transparent outline-none text-foreground placeholder:text-muted/60 text-base"
            autocomplete="off"
            spellcheck="false"
            @keydown="onInputKeydown"
          />
          <kbd
            class="hidden sm:inline-flex items-center px-1.5 py-0.5 text-[10px] font-mono rounded border border-border text-muted"
          >
            ESC
          </kbd>
        </div>
        <div
          ref="listRef"
          class="max-h-[50vh] overflow-y-auto py-2"
          @mousemove="onListMouseMove"
          @mouseleave="onListMouseLeave"
        >
          <div
            v-if="grouped.length === 0"
            class="px-5 py-8 text-center text-muted text-sm"
          >
            Ничего не найдено
          </div>
          <div v-for="[groupName, groupCmds] in grouped" :key="groupName" class="mb-2">
            <div
              class="px-5 py-1.5 text-[10px] font-semibold uppercase tracking-wider text-muted/70"
            >
              {{ groupName }}
            </div>
            <button
              v-for="cmd in groupCmds"
              :key="cmd.id"
              :data-cmd-id="cmd.id"
              :data-selected="flatFiltered[selectedIndex]?.id === cmd.id"
              @click="execute(cmd)"
              class="relative w-full flex items-center gap-3 px-5 py-2.5 text-left transition-colors duration-100"
              :class="
                flatFiltered[selectedIndex]?.id === cmd.id
                  ? 'bg-foreground/[0.08] text-foreground'
                  : 'text-foreground/80 hover:bg-foreground/[0.04]'
              "
            >
              <span
                v-if="flatFiltered[selectedIndex]?.id === cmd.id"
                class="absolute left-0 top-1 bottom-1 w-0.5 bg-foreground rounded-r-full transition-all"
                aria-hidden="true"
              />
              <div
                class="w-7 h-7 rounded-md bg-muted/30 flex items-center justify-center shrink-0"
              >
                <svg
                  v-if="FILL_ICONS.has(cmd.icon)"
                  class="w-4 h-4 text-foreground/80"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path :d="ICONS[cmd.icon]" />
                </svg>
                <svg
                  v-else
                  class="w-3.5 h-3.5 text-foreground/80"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    :d="ICONS[cmd.icon]"
                  />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium truncate">{{ cmd.label }}</div>
                <div v-if="cmd.description" class="text-xs text-muted truncate">
                  {{ cmd.description }}
                </div>
              </div>
              <div
                v-if="cmd.shortcut"
                class="hidden sm:flex items-center gap-1 text-[10px] font-mono text-muted/80"
              >
                <kbd
                  v-for="(k, i) in cmd.shortcut"
                  :key="i"
                  class="px-1.5 py-0.5 rounded border border-border"
                >
                  {{ k }}
                </kbd>
              </div>
              <svg
                v-else
                class="w-3.5 h-3.5 text-muted/40 shrink-0"
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
            </button>
          </div>
        </div>
        <div
          class="flex items-center justify-between px-5 py-2 border-t border-border text-[10px] font-mono text-muted/70 bg-muted/10"
        >
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">↑↓</kbd>
              навигация</span
            >
            <span class="hidden sm:flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">Tab</kbd> цикл</span
            >
            <span class="flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">↵</kbd> выбрать</span
            >
          </div>
          <span class="flex items-center gap-2">
            <span
              >открыть:
              <kbd class="px-1 py-0.5 rounded border border-border">/</kbd></span
            >
            <span class="opacity-50">или</span>
            <kbd class="px-1 py-0.5 rounded border border-border">⌘⇧K</kbd>
          </span>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@keyframes slide-down {
  from {
    opacity: 0;
    transform: translateY(-12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
.animate-slide-down {
  animation: slide-down 0.18s ease-out;
}
.palette-enter-active,
.palette-leave-active {
  transition: opacity 0.15s ease;
}
.palette-enter-from,
.palette-leave-to {
  opacity: 0;
}
</style>
