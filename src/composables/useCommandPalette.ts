import { ref, nextTick } from 'vue'

export function useCommandPalette() {
    const isOpen = ref(false)
    const query = ref('')
    const selectedIndex = ref(0)
    const expandedId = ref<string | null>(null)

    let previouslyFocused: HTMLElement | null = null
    let holdTimer: ReturnType<typeof setTimeout> | null = null
    let suppressClick = false

    const open = () => {
        previouslyFocused = document.activeElement instanceof HTMLElement && document.activeElement !== document.body
            ? document.activeElement
            : null
        isOpen.value = true
        query.value = ''
        selectedIndex.value = 0
        nextTick()
    }

    const close = () => {
        isOpen.value = false
        expandedId.value = null
        clearHoldTimer()
        suppressClick = false
        restoreFocus()
    }

    const restoreFocus = () => {
        const el = previouslyFocused
        previouslyFocused = null
        if (el?.isConnected) { el.focus({ preventScroll: true }); return }
        if (document.activeElement instanceof HTMLElement) document.activeElement.blur()
    }

    const clearHoldTimer = () => {
        if (holdTimer !== null) { clearTimeout(holdTimer); holdTimer = null }
    }

    const onCmdPointerDown = (id: string) => {
        suppressClick = false
        clearHoldTimer()
        holdTimer = setTimeout(() => {
            suppressClick = true
            expandedId.value = expandedId.value === id ? null : id
        }, 400)
    }

    const onCmdPointerEnd = () => clearHoldTimer()

    const shouldSuppressClick = (): boolean => {
        if (suppressClick) { suppressClick = false; return true }
        return false
    }

    return {
        isOpen, query, selectedIndex, expandedId,
        open, close,
        onCmdPointerDown, onCmdPointerEnd, shouldSuppressClick,
    }
}