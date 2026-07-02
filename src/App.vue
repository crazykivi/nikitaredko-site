<script setup lang="ts">
import AppHeader from './components/AppHeader.vue'
import AppFooter from './components/AppFooter.vue'
import LoadingBar from './components/LoadingBar.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background text-foreground">
    <LoadingBar />
    <AppHeader /> 
    <main class="flex-1 flex flex-col pt-16">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <AppFooter />
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