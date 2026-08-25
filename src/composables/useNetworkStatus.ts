import { ref, readonly, onMounted } from 'vue'

const isOnline = ref(navigator.onLine)
const lastTransition = ref<'online' | 'offline' | null>(null)

let initialized = false
let previousStatus = navigator.onLine

function updateStatus(online: boolean) {
  if (online !== previousStatus) {
    isOnline.value = online
    lastTransition.value = online ? 'online' : 'offline'
    previousStatus = online
  }
}

export function useNetworkStatus() {
  onMounted(() => {
    if (!initialized) {
      window.addEventListener('online', () => updateStatus(true))
      window.addEventListener('offline', () => updateStatus(false))
      initialized = true
    }
  })

  return {
    isOnline: readonly(isOnline),
    lastTransition: readonly(lastTransition),
  }
}