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

      <div class="header-actions flex items-center gap-4 sm:gap-6">
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

        <button @click="cycleTheme" class="p-2 rounded-lg hover:bg-border/50 transition-colors minecraft-icon"
          :aria-label="label" :title="label" type="button">
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
          <svg v-else class="w-8 h-8" viewBox="0 0 20 20">
            <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
              <g transform="translate(-420.000000, -6319.000000)" fill="#000000">
                <g id="icons" transform="translate(56.000000, 160.000000)">
                  <path
                    d="M376,6169 L380,6169 L380,6165 L376,6165 L376,6169 Z M368,6169 L372,6169 L372,6165 L368,6165 L368,6169 Z M382,6177 L378,6177 L378,6171 L376,6171 L376,6169 L372,6169 L372,6171 L370,6171 L370,6177 L366,6177 L366,6161 L382,6161 L382,6177 Z M372,6177 L376,6177 L376,6175 L372,6175 L372,6177 Z M364,6179 L384,6179 L384,6159 L364,6159 L364,6179 Z"
                    id="emoji_minecraft_square-[#409]">
                  </path>
                </g>
              </g>
            </g>
          </svg>
        </button>
      </div>
    </nav>
  </header>
</template>