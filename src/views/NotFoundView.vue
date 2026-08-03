<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const glitchText = ref('404')
const isGlitching = ref(false)
let glitchInterval: ReturnType<typeof setInterval> | null = null

const characters = '!<>-_\\/[]{}—=+*^?#________'

const startGlitch = () => {
  glitchInterval = setInterval(() => {
    isGlitching.value = true
    let iteration = 0
    const target = '404'

    const animate = setInterval(() => {
      glitchText.value = target
        .split('')
        .map((char, index) => {
          if (index < iteration) return char
          return characters[Math.floor(Math.random() * characters.length)]
        })
        .join('')

      iteration += 1 / 3

      if (iteration >= target.length) {
        clearInterval(animate)
        glitchText.value = target
        isGlitching.value = false
      }
    }, 40)
  }, 4000)
}

onMounted(() => {
  startGlitch()
})

onUnmounted(() => {
  if (glitchInterval) clearInterval(glitchInterval)
})
</script>

<template>
  <div class="min-h-[calc(100vh-7.3rem)] flex items-center justify-center px-4 select-none">
    <div class="text-center space-y-8 max-w-2xl">
      <div class="relative">
        <h1
          class="text-[8rem] md:text-[12rem] font-bold font-mono leading-none tracking-tighter"
          :class="{ 'text-foreground': !isGlitching, 'text-red-500': isGlitching }"
        >
          {{ glitchText }}
        </h1>
        <div
          class="absolute inset-0 flex items-center justify-center pointer-events-none"
          aria-hidden="true"
        >
          <span
            v-if="isGlitching"
            class="text-[8rem] md:text-[12rem] font-bold font-mono leading-none tracking-tighter text-blue-500/30 translate-x-1 translate-y-[-2px]"
          >
            404
          </span>
        </div>
      </div>
      <div class="space-y-3 animate-fade-in-delayed">
        <h2 class="text-2xl md:text-3xl font-semibold">
          Страница не найдена
        </h2>
        <p class="text-muted text-lg leading-relaxed max-w-md mx-auto">
          Похоже, вы забрели не туда.
          Эта страница либо удалена, либо никогда не существовала.
        </p>
      </div>
      <div class="flex flex-col sm:flex-row gap-4 justify-center pt-4 animate-fade-in-delayed-2">
        <router-link
          to="/"
          class="group px-8 py-4 bg-foreground text-background rounded-lg font-medium hover:opacity-90 transition-all hover:scale-105 flex items-center justify-center gap-2"
        >
          <svg class="w-4 h-4 transition-transform group-hover:-translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          На главную
        </router-link>
        <router-link
          to="/articles"
          class="group px-8 py-4 border border-border rounded-lg font-medium hover:bg-muted/50 transition-all hover:scale-105 flex items-center justify-center gap-2"
        >
          К статьям
          <svg class="w-4 h-4 transition-transform group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </router-link>
      </div>
      <div class="pt-8 animate-fade-in-delayed-3">
        <p class="font-mono text-xs text-muted/40">
          // error: page_not_found
        </p>
        <p class="font-mono text-xs text-muted/40 mt-1">
          // status: 404 | route: {{ $route.fullPath }}
        </p>
      </div>
    </div>
  </div>
</template>