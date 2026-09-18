import { watch, onUnmounted, type Ref } from 'vue'
import type { ThemeMode } from './useTheme'
import { playMcSound } from '../utils/mcSounds'

export function useCaveSounds(mode: Ref<ThemeMode>) {
    let timeoutId: ReturnType<typeof setTimeout> | null = null
    let isMounted = true

    const playRandomCave = () => {
        if (!isMounted || mode.value !== 'charcoal') return
        playMcSound('cave', { volume: 0.3 })
        scheduleNext()
    }

    const scheduleNext = () => {
        if (!isMounted || mode.value !== 'charcoal') {
            if (timeoutId) clearTimeout(timeoutId)
            timeoutId = null
            return
        }
        const delay = 45000 + Math.random() * 75000
        timeoutId = setTimeout(playRandomCave, delay)
    }

    watch(
        mode,
        (newMode) => {
            if (newMode === 'charcoal') {
                if (timeoutId) clearTimeout(timeoutId)
                const initialDelay = 10000 + Math.random() * 20000
                timeoutId = setTimeout(playRandomCave, initialDelay)
            } else {
                if (timeoutId) {
                    clearTimeout(timeoutId)
                    timeoutId = null
                }
            }
        },
        { immediate: true },
    )

    onUnmounted(() => {
        isMounted = false
        if (timeoutId) clearTimeout(timeoutId)
    })
}