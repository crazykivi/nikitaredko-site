import { onMounted, onUnmounted } from 'vue'
import { playMcSound } from '../utils/mcSounds'
import { useTheme } from './useTheme'

export function useLadderSound() {
    const { mode } = useTheme()
    let lastY = typeof window !== 'undefined' ? window.scrollY : 0
    let lastPlayed = 0
    let accumulatedScroll = 0

    const handleScroll = () => {
        if (mode.value !== 'charcoal') return

        const currentY = window.scrollY
        const delta = Math.abs(currentY - lastY)
        lastY = currentY
        accumulatedScroll += delta

        const now = Date.now()
        if (accumulatedScroll > 250 && now - lastPlayed > 350) {
            playMcSound('ladder', { volume: 0.10, rate: 0.95 + Math.random() * 0.1 })
            lastPlayed = now
            accumulatedScroll = 0
        }
    }

    onMounted(() => {
        window.addEventListener('scroll', handleScroll, { passive: true })
    })

    onUnmounted(() => {
        window.removeEventListener('scroll', handleScroll)
    })
}