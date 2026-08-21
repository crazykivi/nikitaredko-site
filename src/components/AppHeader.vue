<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

type ThemeMode = 'auto' | 'light' | 'dark'

const mode = ref<ThemeMode>('auto')
const systemDark = ref(false)

let mediaQuery: MediaQueryList | null = null
const onSystemThemeChange = (e: MediaQueryListEvent | MediaQueryList) => {
  systemDark.value = e.matches
  if (mode.value === 'auto') updateTheme()
}

const effectiveDark = computed(() => {
  if (mode.value === 'light') return false
  if (mode.value === 'dark') return true
  return systemDark.value
})

const updateTheme = () => {
  document.documentElement.classList.toggle('dark', effectiveDark.value)
}

const cycleTheme = () => {
  const next: Record<ThemeMode, ThemeMode> = {
    auto: 'light',
    light: 'dark',
    dark: 'auto',
  }
  mode.value = next[mode.value]

  if (mode.value === 'auto') {
    localStorage.removeItem('theme')
  } else {
    localStorage.setItem('theme', mode.value)
  }
  updateTheme()
}

const label = computed(() => {
  if (mode.value === 'light') return 'Светлая тема'
  if (mode.value === 'dark') return 'Тёмная тема'
  return 'Автоматическая тема (системная)'
})

onMounted(() => {
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  systemDark.value = mediaQuery.matches

  const stored = localStorage.getItem('theme')
  if (stored === 'light' || stored === 'dark') {
    mode.value = stored
  } else {
    mode.value = 'auto'
  }
  updateTheme()

  mediaQuery.addEventListener('change', onSystemThemeChange)
})

onUnmounted(() => {
  mediaQuery?.removeEventListener('change', onSystemThemeChange)
})
</script>

<template>
  <header class="fixed top-0 left-0 z-50 w-full bg-background/80 backdrop-blur-sm select-none">
    <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
      <router-link to="/" class="text-xl font-bold tracking-tight hover:opacity-70 transition-opacity">
      Nikita Redko
      </router-link>

      <div class="flex items-center gap-6">
        <router-link
          to="/articles"
          class="text-sm font-medium text-muted hover:text-foreground transition-colors"
        >
          Статьи
        </router-link>
        <router-link
          to="/about"
          class="text-sm font-medium text-muted hover:text-foreground transition-colors"
        >
          О себе
        </router-link>
        
        <button
          @click="cycleTheme"
          class="p-2 rounded-lg hover:bg-border/50 transition-colors"
          :aria-label="label"
          :title="label"
        >
          <svg v-if="mode === 'auto'" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 3a9 9 0 0 0 0 18z" fill="currentColor" />
          </svg>
          <svg v-else-if="mode === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </button>
      </div>
    </nav>
  </header>
</template>