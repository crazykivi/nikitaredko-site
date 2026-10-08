import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

vi.mock('../../utils/mcTransition', () => ({
    mcBlockTransition: vi.fn((apply: () => void) => apply()),
}))

describe('useTheme', () => {
    let mockMq: { matches: boolean; addEventListener: ReturnType<typeof vi.fn>; removeEventListener: ReturnType<typeof vi.fn> }

    beforeEach(() => {
        vi.resetModules()
        localStorage.clear()
        document.documentElement.className = ''
        document.documentElement.classList.remove('dark', 'charcoal')

        mockMq = {
            matches: false,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        }
        vi.stubGlobal('matchMedia', vi.fn(() => mockMq))
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it('defaults to auto mode', async () => {
        const { useTheme } = await import('../useTheme')
        const { mode } = useTheme()
        expect(mode.value).toBe('auto')
    })

    it('label describes auto mode', async () => {
        const { useTheme } = await import('../useTheme')
        const { label } = useTheme()
        expect(label.value).toContain('Автоматическая')
    })

    it('initTheme reads stored theme from localStorage', async () => {
        localStorage.setItem('theme', 'dark')
        const { useTheme } = await import('../useTheme')
        const { mode, initTheme } = useTheme()
        initTheme()
        expect(mode.value).toBe('dark')
        expect(document.documentElement.classList.contains('dark')).toBe(true)
    })

    it('initTheme falls back to auto for unknown value', async () => {
        localStorage.setItem('theme', 'neon')
        const { useTheme } = await import('../useTheme')
        const { mode, initTheme } = useTheme()
        initTheme()
        expect(mode.value).toBe('auto')
    })

    it('initTheme registers matchMedia listener', async () => {
        const { useTheme } = await import('../useTheme')
        const { initTheme } = useTheme()
        initTheme()
        expect(mockMq.addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    })

    it('destroyTheme removes listener', async () => {
        const { useTheme } = await import('../useTheme')
        const { initTheme, destroyTheme } = useTheme()
        initTheme()
        destroyTheme()
        expect(mockMq.removeEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    })

    it('cycles auto -> light -> dark -> auto', async () => {
        const { useTheme } = await import('../useTheme')
        const { mode, cycleTheme, initTheme } = useTheme()
        initTheme()

        expect(mode.value).toBe('auto')
        cycleTheme()
        expect(mode.value).toBe('light')
        cycleTheme()
        expect(mode.value).toBe('dark')
        cycleTheme()
        expect(mode.value).toBe('auto')
    })

    it('persists to localStorage on cycle', async () => {
        const { useTheme } = await import('../useTheme')
        const { cycleTheme, initTheme } = useTheme()
        initTheme()

        cycleTheme() // auto -> light
        expect(localStorage.getItem('theme')).toBe('light')

        cycleTheme() // light -> dark
        expect(localStorage.getItem('theme')).toBe('dark')

        cycleTheme() // dark -> auto
        expect(localStorage.getItem('theme')).toBeNull() // auto удаляет ключ
    })

    it('applies dark class to html element', async () => {
        const { useTheme } = await import('../useTheme')
        const { cycleTheme, initTheme } = useTheme()
        initTheme()

        cycleTheme() // -> light
        expect(document.documentElement.classList.contains('dark')).toBe(false)

        cycleTheme() // -> dark
        expect(document.documentElement.classList.contains('dark')).toBe(true)
    })

    it('enterCharcoal sets charcoal mode and class', async () => {
        const { useTheme } = await import('../useTheme')
        const { mode, enterCharcoal } = useTheme()
        enterCharcoal()
        expect(mode.value).toBe('charcoal')
        expect(document.documentElement.classList.contains('charcoal')).toBe(true)
        expect(localStorage.getItem('theme')).toBe('charcoal')
    })

    it('exitCharcoal restores previous mode', async () => {
        localStorage.setItem('theme', 'dark')
        const { useTheme } = await import('../useTheme')
        const { mode, initTheme, enterCharcoal, exitCharcoal } = useTheme()
        initTheme()
        expect(mode.value).toBe('dark')

        enterCharcoal()
        expect(mode.value).toBe('charcoal')

        exitCharcoal()
        expect(mode.value).toBe('dark')
        expect(document.documentElement.classList.contains('charcoal')).toBe(false)
    })

    it('exitCharcoal falls back to auto when no previous mode', async () => {
        const { useTheme } = await import('../useTheme')
        const { mode, enterCharcoal, exitCharcoal } = useTheme()
        enterCharcoal()
        exitCharcoal()
        expect(mode.value).toBe('auto')
    })

    it('effectiveDark is false for light', async () => {
        const { useTheme } = await import('../useTheme')
        const { effectiveDark, cycleTheme, initTheme } = useTheme()
        initTheme()
        cycleTheme() // -> light
        expect(effectiveDark.value).toBe(false)
    })

    it('effectiveDark is true for dark', async () => {
        const { useTheme } = await import('../useTheme')
        const { effectiveDark, cycleTheme, initTheme } = useTheme()
        initTheme()
        cycleTheme() // -> light
        cycleTheme() // -> dark
        expect(effectiveDark.value).toBe(true)
    })

    it('effectiveDark follows systemDark in auto mode', async () => {
        mockMq.matches = true
        const { useTheme } = await import('../useTheme')
        const { effectiveDark, initTheme } = useTheme()
        initTheme()
        expect(effectiveDark.value).toBe(true)
    })

    it('returns correct labels for each mode', async () => {
        const { useTheme } = await import('../useTheme')
        const { label, cycleTheme, initTheme, enterCharcoal } = useTheme()
        initTheme()

        expect(label.value).toContain('Автоматическая')
        cycleTheme()
        expect(label.value).toContain('Светлая')
        cycleTheme()
        expect(label.value).toContain('Тёмная')
        enterCharcoal()
        expect(label.value).toContain('Minecraft')
    })
})