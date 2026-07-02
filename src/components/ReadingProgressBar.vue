<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const progress = ref(0);
const hasScroll = ref(false);
let ticking = false;

const calculateProgress = () => {
  const scrollTop = window.scrollY || document.documentElement.scrollTop;
  const docHeight = document.documentElement.scrollHeight;
  const winHeight = document.documentElement.clientHeight;
  
  const scrollableHeight = docHeight - winHeight;
  
  if (scrollableHeight <= 0) {
    hasScroll.value = false;
    progress.value = 0;
  } else {
    hasScroll.value = true;
    progress.value = (scrollTop / scrollableHeight) * 100;
  }
};

const onScroll = () => {
  if (!ticking) {
    window.requestAnimationFrame(() => {
      calculateProgress();
      ticking = false;
    });
    ticking = true;
  }
};

onMounted(() => {
  calculateProgress();
  window.addEventListener('scroll', onScroll, { passive: true });
  window.addEventListener('resize', calculateProgress);
});

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll);
  window.removeEventListener('resize', calculateProgress);
});
</script>

<template>
  <div v-if="hasScroll" class="fixed top-16  left-0 right-0 z-40">
    <div class="h-1 bg-foreground/20 relative">
      <div 
        class="h-full bg-foreground transition-[width] duration-150 ease-out"
        :style="{ width: `${progress}%` }"
      />
    </div>
  </div>
</template>