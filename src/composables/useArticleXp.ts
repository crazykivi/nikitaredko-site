import { nextTick, onMounted, onUnmounted, watch, type Ref } from 'vue'
import { playMcSound } from '../utils/mcSounds'
import { markArticleRead } from '../utils/readProgress'
import { useTheme } from './useTheme'

const BOTTOM_OFFSET_PX = 500

export function useArticleXp(articleId: Ref<string | null>) {
    const { mode } = useTheme()
    let earnedForCurrent = false

    const isAtBottom = () => {
        const docHeight = document.documentElement.scrollHeight
        const winHeight = window.innerHeight

        if (docHeight <= winHeight + 200) return false

        return window.scrollY + winHeight >= docHeight - BOTTOM_OFFSET_PX
    }

    const check = () => {
        if (earnedForCurrent) return
        if (mode.value !== 'charcoal') return

        const id = articleId.value
        if (!id || !isAtBottom()) return

        earnedForCurrent = true
        if (markArticleRead(id)) {
            window.setTimeout(() => playMcSound('xp', { volume: 0.1 }), 300)
        }
    }

    watch(articleId, () => {
        earnedForCurrent = false
        nextTick(() => window.setTimeout(check, 800))
    })

    onMounted(() => window.addEventListener('scroll', check, { passive: true }))
    onUnmounted(() => window.removeEventListener('scroll', check))
}