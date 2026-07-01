<script setup lang="ts">
import AppHeader from './components/AppHeader.vue'
import LoadingBar from './components/LoadingBar.vue'
</script>

<template>
  <div class="min-h-screen bg-background text-foreground">
    <LoadingBar />
    <AppHeader />
    <main class="pt-16">
      <router-view v-slot="{ Component }">
        <transition name="fade">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
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