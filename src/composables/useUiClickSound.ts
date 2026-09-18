import { onMounted, onUnmounted } from 'vue'
import { playMcSound } from '../utils/mcSounds'

const CLICKABLE_SELECTOR = 'button:not([disabled]), a[href], [role="button"]'

export function useUiClickSound(): void {
    const handleClick = (e: MouseEvent) => {
        const target = e.target as Element | null
        if (!target || typeof target.closest !== 'function') return
        if (!target.closest(CLICKABLE_SELECTOR)) return
        playMcSound('click', { volume: 0.5 })
    }

    onMounted(() => {
        document.addEventListener('click', handleClick, true)
    })
    onUnmounted(() => {
        document.removeEventListener('click', handleClick, true)
    })
}