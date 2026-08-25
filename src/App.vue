<script setup lang="ts">
import AppHeader from './components/AppHeader.vue'
import AppFooter from './components/AppFooter.vue'
import LoadingBar from './components/LoadingBar.vue'
import ScrollToTop from './components/ScrollToTop.vue'
import CommandPalette from './components/CommandPalette.vue'
import OfflineBanner from './components/OfflineBanner.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background text-foreground">
    <LoadingBar />
    <AppHeader /> 
    <OfflineBanner />
    <main class="flex-1 flex flex-col pt-16">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <AppFooter />
    <ScrollToTop />
    <CommandPalette />
  </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>