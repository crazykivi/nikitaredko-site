<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useTheme } from '../composables/useTheme'

const { mode, label, cycleTheme, initTheme, destroyTheme } = useTheme()

const openCommandPalette = () => {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', bubbles: true, cancelable: true }))
}
const handleToggleThemeEvent = () => cycleTheme()

onMounted(() => {
  initTheme()
  window.addEventListener('toggle-app-theme', handleToggleThemeEvent)
})
onUnmounted(() => {
  destroyTheme()
  window.removeEventListener('toggle-app-theme', handleToggleThemeEvent)
})
</script>

<template>
  <header class="fixed top-0 left-0 z-50 w-full bg-background/80 backdrop-blur-sm select-none">
    <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between header">
      <router-link to="/" class="text-xl font-bold tracking-tight hover:opacity-70 transition-opacity">
        Nikita Redko
      </router-link>

      <div class="flex items-center gap-4 sm:gap-6">
        <button
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs text-muted hover:text-foreground hover:bg-muted/30 transition-colors border border-transparent hover:border-border"
          @click="openCommandPalette" title="Command Palette (/)" type="button">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <span>Поиск</span>
          <kbd class="ml-1 px-1 py-0.5 text-[9px] font-mono rounded border border-border text-muted/80">/</kbd>
        </button>

        <router-link to="/articles" class="text-sm font-medium text-muted hover:text-foreground transition-colors">
          Статьи
        </router-link>
        <router-link to="/about" class="text-sm font-medium text-muted hover:text-foreground transition-colors">
          О себе
        </router-link>

        <button @click="cycleTheme" class="p-2 rounded-lg hover:bg-border/50 transition-colors" :aria-label="label"
          :title="label" type="button">
          <svg v-if="mode === 'auto'" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2" aria-hidden="true">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 3a9 9 0 0 0 0 18z" fill="currentColor" />
          </svg>

          <svg v-else-if="mode === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>

          <svg v-else-if="mode === 'dark'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
          <svg v-else class="w-6 h-6" viewBox="0 0 16 16" shape-rendering="crispEdges" aria-hidden="true">
            <path d="M8 1 L14 4 L8 7 L2 4 Z" fill="#5b8731" />
            <rect x="7" y="3" width="1" height="1" fill="#4c7227" />
            <rect x="9" y="5" width="1" height="1" fill="#4c7227" />
            <rect x="5" y="4" width="1" height="1" fill="#6a9f39" />
            <rect x="10" y="4" width="1" height="1" fill="#6a9f39" />

            <path d="M2 4 L8 7 L8 15 L2 12 Z" fill="#735135" />
            <path d="M8 7 L14 4 L14 12 L8 15 Z" fill="#866041" />

            <path d="M2 4 L8 7 L8 9 L2 6 Z" fill="#4c7227" />
            <rect x="4" y="7" width="1" height="1" fill="#4c7227" />
            <rect x="6" y="8" width="1" height="1" fill="#4c7227" />

            <path d="M8 7 L14 4 L14 6 L8 9 Z" fill="#5b8731" />
            <rect x="10" y="8" width="1" height="1" fill="#5b8731" />
            <rect x="12" y="7" width="1" height="1" fill="#5b8731" />

            <rect x="4" y="10" width="1" height="1" fill="#4d3322" />
            <rect x="6" y="12" width="1" height="1" fill="#4d3322" />
            <rect x="10" y="11" width="1" height="1" fill="#5c402b" />
            <rect x="12" y="10" width="1" height="1" fill="#5c402b" />
          </svg>
        </button>
      </div>
    </nav>
  </header>
</template>