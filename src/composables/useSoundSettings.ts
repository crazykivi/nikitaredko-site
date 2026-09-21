import { ref, watch } from 'vue'

export const isSoundMuted = ref(false)

if (typeof window !== 'undefined') {
  isSoundMuted.value = localStorage.getItem('mc-sounds-muted') === 'true'
}

watch(isSoundMuted, (val) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('mc-sounds-muted', String(val))
  }
})

export function useSoundSettings() {
  const toggleMute = () => {
    isSoundMuted.value = !isSoundMuted.value
  }
  return { isSoundMuted, toggleMute }
}