<script setup lang="ts">
import { onMounted } from 'vue'
import { primeMcSounds } from './utils/mcSounds'
import { useTheme } from './composables/useTheme'
import { useCaveSounds } from './composables/useCaveSounds'
import { useUiClickSound } from './composables/useUiClickSound'
import { useLadderSound } from './composables/useLadderSound'
import AppHeader from './components/AppHeader.vue'
import AppFooter from './components/AppFooter.vue'
import SearchFab from './components/SearchFab.vue'
import LoadingBar from './components/LoadingBar.vue'
import ScrollToTop from './components/ScrollToTop.vue'
import SoundToggle from './components/SoundToggle.vue'
import CommandPalette from './components/CommandPalette.vue'
import OfflineBanner from './components/OfflineBanner.vue'
import UpdateToast from './components/UpdateToast.vue'
import BlockBreakOverlay from './components/BlockBreakOverlay.vue'


const { mode } = useTheme()
useCaveSounds(mode)
useUiClickSound()
useLadderSound()

onMounted(() => {
  window.addEventListener('pointerdown', primeMcSounds, { once: true })
})
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
    <div class="fixed bottom-6 right-6 z-[90] flex flex-col gap-3 items-end">
      <SoundToggle />
      <ScrollToTop />
      <SearchFab />
    </div>
    <UpdateToast />
    <BlockBreakOverlay />
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