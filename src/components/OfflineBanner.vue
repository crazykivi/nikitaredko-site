<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useNetworkStatus } from '../composables/useNetworkStatus'
import { useServerStatus } from '../composables/useServerStatus'

const { isOnline, lastTransition } = useNetworkStatus()
const { isServerReachable, recheck } = useServerStatus()

const showToast = ref(false)
const toastType = ref<'online' | 'offline' | 'server-down'>('offline')
let toastTimer: ReturnType<typeof setTimeout> | null = null

const bannerState = computed<'online' | 'offline' | 'server-down'>(() => {
  if (!isOnline.value) return 'offline'
  if (!isServerReachable.value) return 'server-down'
  return 'online'
})

watch(lastTransition, (status) => {
  if (!status) return
  toastType.value = status
  showToast.value = true
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    showToast.value = false
  }, 4000)
})

watch(isServerReachable, (reachable, prev) => {
  if (prev === undefined) return
  toastType.value = reachable ? 'online' : 'server-down'
  showToast.value = true
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    showToast.value = false
  }, 4000)
})

const bannerText = computed(() => {
  if (bannerState.value === 'offline') return 'Вы офлайн — данные загружены из кэша'
  if (bannerState.value === 'server-down') return 'Сервер недоступен — показываем кэш'
  return ''
})

const toastText = computed(() => {
  if (toastType.value === 'offline') return 'Нет соединения'
  if (toastType.value === 'server-down') return 'Сервер не отвечает'
  return 'Соединение восстановлено'
})
</script>

<template>
  <Transition name="slide-down">
    <div
      v-if="bannerState !== 'online'"
      class="fixed top-16 left-0 right-0 z-40 backdrop-blur-sm border-b select-none"
      :class="bannerState === 'offline'
        ? 'bg-amber-500/95 dark:bg-amber-600/95 border-amber-600/30'
        : 'bg-rose-500/95 dark:bg-rose-600/95 border-rose-600/30'"
    >
      <div class="max-w-6xl mx-auto px-4 py-2 flex items-center justify-center gap-2 text-sm font-medium text-white">
        <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M18.364 5.636a9 9 0 11-12.728 0M12 9v4m0 4h.01"
          />
        </svg>
        <span>{{ bannerText }}</span>
        <button
          v-if="bannerState === 'server-down'"
          @click="recheck"
          class="ml-2 px-2 py-0.5 text-xs rounded border border-white/30 hover:bg-white/10 transition-colors"
        >
          Проверить снова
        </button>
      </div>
    </div>
  </Transition>
  <Transition name="toast">
    <div
      v-if="showToast"
      class="fixed top-24 right-4 z-[100] max-w-sm"
      role="alert"
      aria-live="polite"
    >
      <div
        class="flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg border backdrop-blur-sm"
        :class="{
          'bg-amber-500/95 border-amber-600/30 text-amber-950 dark:text-amber-50': toastType === 'offline',
          'bg-rose-500/95 border-rose-600/30 text-white': toastType === 'server-down',
          'bg-emerald-500/95 border-emerald-600/30 text-emerald-950 dark:text-emerald-50': toastType === 'online',
        }"
      >
        <svg
          v-if="toastType !== 'online'"
          class="w-5 h-5 shrink-0"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M18.364 5.636a9 9 0 11-12.728 0M12 9v4m0 4h.01"
          />
        </svg>
        <svg
          v-else
          class="w-5 h-5 shrink-0"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"
          />
        </svg>
        <div class="text-sm font-medium">{{ toastText }}</div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.25s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-100%);
}

.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.21, 1.02, 0.73, 1);
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>