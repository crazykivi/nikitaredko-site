import { onMounted, onUnmounted } from 'vue'
import { playMcSound } from '../utils/mcSounds'
import { useTheme } from './useTheme'

const CLICKABLE_SELECTOR = 'button:not([disabled]), a[href], [role="button"]'

export function useUiClickSound(): void {
    const { mode } = useTheme()

    const handleClick = (e: MouseEvent) => {
        if (mode.value !== 'charcoal') return

        const target = e.target as Element | null
        if (!target || typeof target.closest !== 'function') return

        // Исключаем зоны, где клики не должны звучать
        if (target.closest('.prose, .gsc-comment, .gsc-timeline, giscus-widget, header, footer, .modal-panel, .bb-overlay, [role="dialog"], .command-palette')) return

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