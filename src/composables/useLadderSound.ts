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
        // Звук каждые 300px скролла, но не чаще раза в 1.2 секунды
        if (accumulatedScroll > 300 && now - lastPlayed > 1200) {
            playMcSound('ladder', { volume: 0.15, rate: 0.9 + Math.random() * 0.2 })
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