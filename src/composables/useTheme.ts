import { computed, ref } from 'vue'
import { mcBlockTransition } from '../utils/mcTransition'

export type ThemeMode = 'auto' | 'light' | 'dark' | 'charcoal'
type BaseMode = Exclude<ThemeMode, 'charcoal'>

const BASE_CYCLE: Record<BaseMode, BaseMode> = {
  auto: 'light',
  light: 'dark',
  dark: 'auto',
}

const mode = ref<ThemeMode>('auto')
const systemDark = ref(false)
const preCharcoal = ref<BaseMode | null>(null)
let mediaQuery: MediaQueryList | null = null

const effectiveDark = computed(() => {
  if (mode.value === 'light') return false
  if (mode.value === 'dark' || mode.value === 'charcoal') return true
  return systemDark.value
})

const applyTheme = () => {
  const root = document.documentElement
  root.classList.toggle('dark', effectiveDark.value)
  root.classList.toggle('charcoal', mode.value === 'charcoal')
}

const persist = (next: ThemeMode) => {
  mode.value = next
  if (next === 'auto') localStorage.removeItem('theme')
  else localStorage.setItem('theme', next)
  applyTheme()
}

const onSystemThemeChange = (e: MediaQueryListEvent | MediaQueryList) => {
  systemDark.value = e.matches
  if (mode.value === 'auto') applyTheme()
}

export function useTheme() {
  const initTheme = () => {
    if (!mediaQuery) {
      mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      mediaQuery.addEventListener('change', onSystemThemeChange)
    }
    systemDark.value = mediaQuery.matches
    const stored = localStorage.getItem('theme')
    mode.value = stored === 'light' || stored === 'dark' || stored === 'charcoal' ? stored : 'auto'
    applyTheme()
  }

  const destroyTheme = () => {
    mediaQuery?.removeEventListener('change', onSystemThemeChange)
    mediaQuery = null
  }

  const enterCharcoal = () => {
    if (mode.value !== 'charcoal') preCharcoal.value = mode.value as BaseMode
    persist('charcoal')
  }

  const exitCharcoal = () => {
    const back = preCharcoal.value ?? 'auto'
    preCharcoal.value = null
    persist(back)
  }

  const cycleTheme = (e?: MouseEvent) => {
    if (mode.value === 'charcoal') {
      mcBlockTransition(() => exitCharcoal())
      return
    }

    const apply = () => persist(BASE_CYCLE[mode.value as BaseMode])

    const supported =
      typeof document.startViewTransition === 'function' &&
      !window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const x = e?.clientX
    const y = e?.clientY
    if (!supported || typeof x !== 'number' || typeof y !== 'number') {
      apply()
      return
    }

    const transition = document.startViewTransition(apply)
    transition.ready.then(() => {
      const endRadius = Math.hypot(Math.max(x, window.innerWidth - x), Math.max(y, window.innerHeight - y))
      document.documentElement.animate(
        { clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`] },
        { duration: 400, easing: 'ease-in-out', pseudoElement: '::view-transition-new(root)' },
      )
    })
  }

  const label = computed(() => {
    if (mode.value === 'light') return 'Светлая тема'
    if (mode.value === 'dark') return 'Тёмная тема'
    if (mode.value === 'charcoal') return 'Угольная тема'
    return 'Автоматическая тема (системная)'
  })

  return { mode, effectiveDark, label, cycleTheme, enterCharcoal, exitCharcoal, initTheme, destroyTheme }
}