<script setup lang="ts">
import { ref, watch } from 'vue'
import { isGlobalLoading } from '../utils/loading'

const progress = ref(0)
const isVisible = ref(false)
const isFadingOut = ref(false)

let timer: ReturnType<typeof setTimeout> | null = null
let hideTimer: ReturnType<typeof setTimeout> | null = null

const start = () => {
  if (hideTimer) clearTimeout(hideTimer)
  
  isVisible.value = true
  isFadingOut.value = false
  progress.value = 0
  
  const increase = () => {
    if (isGlobalLoading.value) {
      progress.value += Math.random() * 15
      if (progress.value > 90) progress.value = 90
      timer = setTimeout(increase, 200)
    }
  }
  increase()
}

const finish = () => {
  if (timer) clearTimeout(timer)
  progress.value = 100
  
  hideTimer = setTimeout(() => {
    isFadingOut.value = true
    setTimeout(() => {
      isVisible.value = false
      isFadingOut.value = false
      progress.value = 0
    }, 300)
  }, 200)
}

watch(isGlobalLoading, (newVal) => {
  if (newVal) {
    start()
  } else {
    finish()
  }
})
</script>

<template>
  <div
    v-show="isVisible"
    class="fixed top-0 left-0 w-full h-1 z-[100] overflow-hidden transition-opacity duration-300"
    :class="{ 'opacity-0': isFadingOut }"
  >
    <div
      class="h-full bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 transition-[width] duration-200 ease-out"
      :style="{ width: `${progress}%` }"
    />
  </div>
</template>