import { ref, onMounted, onUnmounted } from 'vue'

const COPY_BTN_ICON = `<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>`

export function useCopyButtons() {
    const copyToast = ref<{ type: 'success' | 'error' } | null>(null)
    let copyToastTimer: ReturnType<typeof setTimeout> | null = null

    const showCopyToast = (type: 'success' | 'error') => {
        copyToast.value = { type }
        if (copyToastTimer) clearTimeout(copyToastTimer)
        copyToastTimer = setTimeout(() => { copyToast.value = null }, 2500)
    }

    const handleCopyClick = (e: MouseEvent) => {
        const target = e.target as Element | null
        if (!target) return

        const button = target.closest('.copy-code-btn') as HTMLButtonElement | null
        if (!button) return

        const codeBlock = button.closest('.code-block-wrapper')
        const codeElement = codeBlock?.querySelector('code')
        if (!codeElement) return

        const code = codeElement.textContent || ''
        navigator.clipboard
            .writeText(code)
            .then(() => showCopyToast('success'))
            .catch(() => showCopyToast('error'))
    }

    const injectCopyButtons = () => {
        document.querySelectorAll('.code-header').forEach((header) => {
            if (header.querySelector('.copy-code-btn')) return
            const btn = document.createElement('button')
            btn.type = 'button'
            btn.title = 'Копировать код'
            btn.className =
                'copy-code-btn text-xs px-3 py-1 rounded-md text-muted hover:text-foreground transition-all flex items-center gap-1'
            btn.innerHTML = `${COPY_BTN_ICON}<span>Копировать</span>`
            header.appendChild(btn)
        })
    }

    onMounted(() => document.addEventListener('click', handleCopyClick))
    onUnmounted(() => {
        document.removeEventListener('click', handleCopyClick)
        if (copyToastTimer) clearTimeout(copyToastTimer)
    })

    return { copyToast, injectCopyButtons }
}