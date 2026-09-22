<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useTheme } from '../composables/useTheme'
import { playMcSound, stopMcSound } from '../utils/mcSounds'

const { mode, label, cycleTheme, initTheme, destroyTheme } = useTheme()
const isFusing = ref(false)
const isExploding = ref(false)

const openCommandPalette = () => {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', bubbles: true, cancelable: true }))
}
const handleToggleThemeEvent = () => cycleTheme()

const onPointerDown = () => {
  if (mode.value !== 'charcoal') return
  isFusing.value = true
  playMcSound('hiss', { volume: 0.5, loop: true })
}

const onPointerUp = () => {
  if (mode.value === 'charcoal') {
    if (!isFusing.value) return
    isFusing.value = false
    stopMcSound('hiss')

    isExploding.value = true
    playMcSound('explode', { volume: 0.7 })
    document.documentElement.classList.add('bb-shake')
    setTimeout(() => document.documentElement.classList.remove('bb-shake'), 450)

    setTimeout(() => {
      cycleTheme()
      isExploding.value = false
    }, 350)
  } else {
    cycleTheme()
  }
}

const onPointerCancel = () => {
  if (isFusing.value) {
    isFusing.value = false
    stopMcSound('hiss')
  }
}

onMounted(() => {
  initTheme()
  window.addEventListener('toggle-app-theme', handleToggleThemeEvent)
})
onUnmounted(() => {
  destroyTheme()
  stopMcSound('hiss')
  window.removeEventListener('toggle-app-theme', handleToggleThemeEvent)
})
</script>

<template>
  <header class="fixed top-0 left-0 z-50 w-full bg-background/80 backdrop-blur-sm select-none mc-header">
    <div class="mc-clouds-container">
      <svg class="mc-cloud mc-cloud-1" viewBox="0 0 100 30" fill="currentColor" shape-rendering="crispEdges">
        <rect x="10" y="10" width="80" height="10" />
        <rect x="20" y="0" width="40" height="10" />
        <rect x="0" y="20" width="100" height="10" />
      </svg>
      <svg class="mc-cloud mc-cloud-2" viewBox="0 0 120 40" fill="currentColor" shape-rendering="crispEdges">
        <rect x="20" y="10" width="80" height="20" />
        <rect x="40" y="0" width="40" height="10" />
      </svg>
      <svg class="mc-cloud mc-cloud-3" viewBox="0 0 80 20" fill="currentColor" shape-rendering="crispEdges">
        <rect x="10" y="5" width="60" height="10" />
        <rect x="0" y="10" width="80" height="10" />
      </svg>
    </div>

    <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between header relative z-10">
      <router-link to="/" aria-label="Nikita Redko"
        class="text-xl font-bold tracking-tight hover:opacity-70 transition-opacity mc-logo">
        <span>Nikita Redko</span>
      </router-link>

      <div class="header-actions flex items-center gap-4 sm:gap-6 mc-hotbar">
        <button type="button" @click="openCommandPalette" title="Command Palette (/)" aria-label="Поиск"
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs text-muted hover:text-foreground hover:bg-muted/30 transition-colors border border-transparent hover:border-border mc-slot mc-search-slot">
          <svg class="w-3.5 h-3.5 mc-hide-in-charcoal" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <span class="mc-hide-in-charcoal">Поиск</span>
          <kbd
            class="ml-1 px-1 py-0.5 text-[9px] font-mono rounded border border-border text-muted/80 mc-hide-in-charcoal">/</kbd>
          <svg class="w-8 h-8 mc-show-in-charcoal" viewBox="0 0 16 16" shape-rendering="crispEdges" aria-hidden="true">
            <rect x="4" y="1" width="6" height="2" fill="#e8e8e8" />
            <rect x="2" y="3" width="2" height="6" fill="#e8e8e8" />
            <rect x="10" y="3" width="2" height="6" fill="#e8e8e8" />
            <rect x="4" y="9" width="6" height="2" fill="#e8e8e8" />
            <rect x="4" y="3" width="6" height="6" fill="#6fb7d8" />
            <rect x="5" y="4" width="2" height="2" fill="#a8d8ea" />
            <rect x="10" y="10" width="2" height="2" fill="#e8e8e8" />
            <rect x="12" y="12" width="3" height="2" fill="#e8e8e8" />
          </svg>
          <span class="mc-tip mc-show-in-charcoal">Поиск [/]</span>
        </button>

        <router-link to="/articles" aria-label="Статьи"
          class="text-sm font-medium text-muted hover:text-foreground transition-colors mc-slot">
          <span class="mc-hide-in-charcoal">Статьи</span>
          <svg class="w-8 h-8 mc-show-in-charcoal" viewBox="0 0 16 16" shape-rendering="crispEdges" aria-hidden="true">
            <rect x="3" y="2" width="10" height="12" fill="#8b5a2b" />
            <rect x="3" y="2" width="2" height="12" fill="#5f3c1e" />
            <rect x="6" y="4" width="6" height="8" fill="#e8e8e8" />
            <rect x="7" y="6" width="4" height="1" fill="#8f8f8f" />
            <rect x="7" y="8" width="4" height="1" fill="#8f8f8f" />
            <rect x="7" y="10" width="3" height="1" fill="#8f8f8f" />
          </svg>
          <span class="mc-tip mc-show-in-charcoal">Статьи</span>
        </router-link>

        <router-link to="/about" aria-label="О себе"
          class="text-sm font-medium text-muted hover:text-foreground transition-colors mc-slot">
          <span class="mc-hide-in-charcoal">О себе</span>
          <svg class="w-8 h-8 mc-show-in-charcoal" viewBox="0 0 8 8" shape-rendering="crispEdges" aria-hidden="true">
            <rect width="8" height="8" fill="#c69c6d"></rect>

            <rect x="0" y="2" width="1" height="1" fill="#4a2f1b"></rect>
            <rect x="7" y="2" width="1" height="1" fill="#4a2f1b"></rect>

            <rect width="8" height="2" fill="#4a2f1b"></rect>
            <rect x="1" width="1" height="1" fill="#ffffff" y="4"></rect>
            <rect x="2" width="1" height="1" fill="#4444cc" y="4"></rect>
            <rect x="5" width="1" height="1" fill="#4444cc" y="4"></rect>
            <rect x="6" width="1" height="1" fill="#ffffff" y="4"></rect>
            <rect x="3" width="2" fill="#a06a50" height="1" y="5"></rect>
            <rect x="2" y="6" width="1" height="1" fill="#4a2f1b"></rect>
            <rect x="3" y="6" width="2" height="1" fill="#a55f52"></rect>
            <rect x="5" y="6" width="1" height="1" fill="#4a2f1b"></rect>
            <rect height="1" fill="#4a2f1b" x="2" y="7" width="4"></rect>
          </svg>
          <span class="mc-tip mc-show-in-charcoal">О себе</span>
        </router-link>

        <button type="button" @pointerdown="onPointerDown" @pointerup="onPointerUp" @pointerleave="onPointerCancel"
          @pointercancel="onPointerCancel" :class="{ 'is-fusing': isFusing, 'is-exploding': isExploding }"
          :aria-label="label" :title="label"
          class="p-2 rounded-lg hover:bg-border/50 transition-colors minecraft-icon mc-slot">
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
          <svg v-else class="w-8 h-8 mc-creeper" viewBox="0 0 16 16" shape-rendering="crispEdges" aria-hidden="true">
            <rect width="16" height="16" fill="#55b55a" />
            <rect x="1" y="1" width="2" height="1" fill="#6fd66f" />
            <rect x="12" y="13" width="2" height="1" fill="#3f9445" />
            <rect x="2" y="4" width="4" height="4" fill="#0b0b0b" />
            <rect x="10" y="4" width="4" height="4" fill="#0b0b0b" />
            <rect x="6" y="8" width="4" fill="#0b0b0b" height="5"></rect>
            <rect x="4" y="10" width="2" height="4" fill="#0b0b0b" />
            <rect x="10" y="10" width="2" height="4" fill="#0b0b0b" />
          </svg>
          <span class="mc-tip mc-show-in-charcoal">{{ label }}</span>
        </button>
      </div>
    </nav>
  </header>
</template>