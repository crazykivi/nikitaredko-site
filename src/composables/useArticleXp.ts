import { nextTick, onMounted, onUnmounted, watch, type Ref } from 'vue'
import { playMcSound } from '../utils/mcSounds'
import { markArticleRead } from '../utils/readProgress'

const BOTTOM_OFFSET_PX = 500

export function useArticleXp(articleId: Ref<string | null>) {
    let earnedForCurrent = false

    const isAtBottom = () =>
        window.scrollY + window.innerHeight >=
        document.documentElement.scrollHeight - BOTTOM_OFFSET_PX

    const check = () => {
        if (earnedForCurrent) return
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