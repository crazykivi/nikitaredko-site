import { ref, readonly, onMounted } from 'vue'

const needRefresh = ref(false)
const offlineReady = ref(false)
const updating = ref(false)

let initialized = false

export function applyUpdate() {
  if (updating.value) return
  updating.value = true
  window.location.reload()
}

export function dismissUpdate() {
  needRefresh.value = false
}

export function usePwaUpdate() {
  onMounted(() => {
    if (initialized) return
    initialized = true

    if (import.meta.env.DEV || import.meta.env.VITE_CAPACITOR_BUILD === 'true' || !('serviceWorker' in navigator)) {
      return
    }

    import('virtual:pwa-register').then(({ registerSW }) => {
      registerSW({
        immediate: true,
        onNeedReload() {
          needRefresh.value = true
        },
        onOfflineReady() {
          offlineReady.value = true
          console.info('[PWA] Приложение готово к офлайн-работе')
        },
        onRegisterError(error) {
          console.error('[PWA] Ошибка регистрации service worker:', error)
        },
      })
    })
  })

  return {
    needRefresh: readonly(needRefresh),
    offlineReady: readonly(offlineReady),
    updating: readonly(updating),
    applyUpdate,
    dismissUpdate,
  }
}
