import { readonly, ref, watch } from 'vue'

const stored = typeof localStorage !== 'undefined' && localStorage.getItem('block-break-unlocked') === 'true'
const unlocked = ref(stored)

watch(unlocked, (val) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('block-break-unlocked', String(val))
  }
})

export function useBlockBreak() {
  const toggle = (): boolean => {
    unlocked.value = !unlocked.value
    return unlocked.value
  }

  return { unlocked: readonly(unlocked), toggle }
}