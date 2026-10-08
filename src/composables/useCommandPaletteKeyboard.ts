import { type Ref, type ComputedRef } from 'vue'
import type { Command } from '../types/commandPalette'

interface KeyboardOptions {
    isOpen: Ref<boolean>
    selectedIndex: Ref<number>
    flatCommands: ComputedRef<Command[]>
    onExecute: (cmd: Command) => void
    onOpen: () => void
    onClose: () => void
}

export function useCommandPaletteKeyboard(opts: KeyboardOptions) {
    const { isOpen, selectedIndex, flatCommands, onExecute, onOpen, onClose } = opts

    const handleInputKeydown = (e: KeyboardEvent) => {
        const len = flatCommands.value.length
        if (e.key === 'Escape') { e.preventDefault(); onClose(); return }
        if (len === 0) return

        switch (e.key) {
            case 'ArrowDown': e.preventDefault(); selectedIndex.value = (selectedIndex.value + 1) % len; break
            case 'ArrowUp': e.preventDefault(); selectedIndex.value = (selectedIndex.value - 1 + len) % len; break
            case 'Tab':
                e.preventDefault()
                selectedIndex.value = e.shiftKey
                    ? (selectedIndex.value - 1 + len) % len
                    : (selectedIndex.value + 1) % len
                break
            case 'Home': e.preventDefault(); selectedIndex.value = 0; break
            case 'End': e.preventDefault(); selectedIndex.value = len - 1; break
            case 'Enter': {
                e.preventDefault()
                const cmd = flatCommands.value[selectedIndex.value]
                if (cmd) onExecute(cmd)
                break
            }
        }
    }

    const handleGlobalKeydown = (e: KeyboardEvent) => {
        const isMac = navigator.platform.toUpperCase().includes('MAC')
        const mod = isMac ? e.metaKey : e.ctrlKey

        if (e.key === '/' && !mod && !e.altKey && !e.shiftKey) {
            if (isTypingField(e.target)) return
            e.preventDefault()
            isOpen.value ? onClose() : onOpen()
            return
        }

        if (mod && e.shiftKey && e.key.toLowerCase() === 'k') {
            e.preventDefault()
            isOpen.value ? onClose() : onOpen()
            return
        }

        if (!isOpen.value && mod && !e.shiftKey && /^[1-4]$/.test(e.key)) {
            if (isTypingField(e.target)) return
            e.preventDefault()
            // Быстрые переходы обрабатываются через staticCommands
            return
        }

        if (isOpen.value && e.key === 'Escape') {
            e.preventDefault()
            onClose()
        }
    }

    return { handleInputKeydown, handleGlobalKeydown }
}

function isTypingField(target: EventTarget | null): boolean {
    if (!(target instanceof HTMLElement)) return false
    const tag = target.tagName.toLowerCase()
    return tag === 'input' || tag === 'textarea' || tag === 'select' || target.isContentEditable
}