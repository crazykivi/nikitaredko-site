import { ref } from 'vue'

let loadingCount = 0
export const isGlobalLoading = ref(false)

export function startLoading() {
  loadingCount++
  isGlobalLoading.value = true
}

export function stopLoading() {
  loadingCount--
  if (loadingCount <= 0) {
    loadingCount = 0
    isGlobalLoading.value = false
  }
}