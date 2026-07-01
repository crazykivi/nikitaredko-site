<script setup lang="ts">
import { ref, onMounted } from 'vue'

const phrases = [
  'о коде и багах',
  'о жизни и кофе',
  'о проектах',
  'о мыслях вслух',
  'о всём подряд'
]

const currentPhrase = ref(0)
const displayedText = ref('')
const isDeleting = ref(false)
const typingSpeed = ref(100)

const typeEffect = () => {
  const fullText = phrases[currentPhrase.value]
  
  if (!isDeleting.value) {
    displayedText.value = fullText.substring(0, displayedText.value.length + 1)
    
    if (displayedText.value === fullText) {
      isDeleting.value = true
      typingSpeed.value = 2000
    } else {
      typingSpeed.value = 100
    }
  } else {
    displayedText.value = fullText.substring(0, displayedText.value.length - 1)
    
    if (displayedText.value === '') {
      isDeleting.value = false
      currentPhrase.value = (currentPhrase.value + 1) % phrases.length
      typingSpeed.value = 500
    } else {
      typingSpeed.value = 50
    }
  }
  
  setTimeout(typeEffect, typingSpeed.value)
}

onMounted(() => {
  setTimeout(typeEffect, 1000)
})
</script>

<template>
  <div class="min-h-[calc(100vh-5rem)] flex items-center justify-center px-4">
    <div class="text-center space-y-8 max-w-4xl">
      <div class="space-y-4 min-h-[8.3rem] md:min-h-[9rem]">
        <p class="text-sm font-mono text-muted uppercase tracking-widest animate-fade-in">
          Привет, я Никита
        </p>
        <h1 class="text-5xl md:text-7xl min-w-0 md:min-w-[800px] font-bold tracking-tight">
          Пишу
          <span class="text-foreground/60">
            {{ displayedText }}
          </span>
          <span class="animate-pulse">|</span>
        </h1>
      </div>
      
      <p class="text-lg md:text-xl text-muted leading-relaxed max-w-2xl mx-auto animate-fade-in-delayed">
        Здесь я делюсь опытом разработки, рассказываю о своих проектах,
        размышляю о жизни и просто записываю мысли, чтобы не забыть.
        Иногда получается полезно, иногда — просто интересно.
      </p>
      
      <div class="flex gap-4 justify-center pt-4 animate-fade-in-delayed-3">
        <router-link
          to="/articles"
          class="group px-8 py-4 bg-foreground text-background rounded-lg font-medium hover:opacity-90 transition-all hover:scale-105 flex items-center gap-2"
        >
          Читать статьи
          <svg
            class="w-4 h-4 transition-transform group-hover:translate-x-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </router-link>
      </div>
      
      <div class="pt-8 text-xs text-muted/40 max-w-md mx-auto animate-fade-in-delayed-3">
        <p class="leading-relaxed">
          Всё написанное здесь — мои личные мысли, опыт и правки от ИИ.
          Не принимайте всё за истину — всегда проверяйте информацию самостоятельно.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.6s ease-out;
}

.animate-fade-in-delayed {
  animation: fadeIn 0.6s ease-out 0.2s both;
}

.animate-fade-in-delayed-2 {
  animation: fadeIn 0.6s ease-out 0.4s both;
}

.animate-fade-in-delayed-3 {
  animation: fadeIn 0.6s ease-out 0.6s both;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>