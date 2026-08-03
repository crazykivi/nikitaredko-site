<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const progress = ref(0);
const hasScroll = ref(false);
let ticking = false;
let articleEl: HTMLElement | null = null;

const calculateProgress = () => {
  if (!articleEl) {
    articleEl = document.querySelector('.article-content');
    if (!articleEl) {
      hasScroll.value = false;
      progress.value = 0;
      return;
    }
  }

  const scrollTop = window.scrollY;
  const winHeight = window.innerHeight;
  
  const articleTop = articleEl.offsetTop;
  const articleHeight = articleEl.offsetHeight;
  const readableHeight = articleHeight - winHeight;

  if (readableHeight <= 0) {
    hasScroll.value = false;
    progress.value = 100;
    return;
  }

  hasScroll.value = true;

  const scrolled = scrollTop - articleTop;
  const raw = (scrolled / readableHeight) * 100;

  progress.value = Math.min(100, Math.max(0, raw));
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
  setTimeout(() => {
    calculateProgress();
  }, 100);

  window.addEventListener('scroll', onScroll, { passive: true });
  window.addEventListener('resize', calculateProgress);
});

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll);
  window.removeEventListener('resize', calculateProgress);
});
</script>

<template>
  <div v-if="hasScroll" class="fixed top-16 left-0 right-0 z-40">
    <div class="h-1 bg-foreground/20 relative">
      <div
        class="h-full bg-foreground transition-[width] duration-150 ease-out"
        :style="{ width: `${progress}%` }"
      />
    </div>
  </div>
</template>