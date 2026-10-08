import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

describe('useServerStatus', () => {
    let mockFetch: ReturnType<typeof vi.fn>

    beforeEach(() => {
        vi.resetModules()
        vi.useFakeTimers()
        mockFetch = vi.fn()
        vi.stubGlobal('fetch', mockFetch)
    })

    afterEach(() => {
        vi.useRealTimers()
        vi.unstubAllGlobals()
    })

    const mountHelper = async () => {
        const { useServerStatus } = await import('../useServerStatus')
        const { defineComponent, h, nextTick } = await import('vue')
        const { mount } = await import('@vue/test-utils')

        let status!: ReturnType<typeof useServerStatus>
        const Comp = defineComponent({
            setup() {
                status = useServerStatus()
                return () => h('div')
            },
        })

        const wrapper = mount(Comp)
        return { status, wrapper, nextTick }
    }

    it('starts with isServerReachable=true', async () => {
        mockFetch.mockResolvedValue({ ok: true })
        const { status, wrapper } = await mountHelper()

        expect(status.isServerReachable.value).toBe(true)
        wrapper.unmount()
    })

    it('sets isServerReachable=false when fetch rejects', async () => {
        mockFetch.mockRejectedValue(new Error('Network error'))
        const { status, wrapper } = await mountHelper()

        await vi.advanceTimersByTimeAsync(0)

        expect(status.isServerReachable.value).toBe(false)
        wrapper.unmount()
    })

    it('sets isServerReachable=false on non-ok response', async () => {
        mockFetch.mockResolvedValue({ ok: false, status: 503 })
        const { status, wrapper } = await mountHelper()

        await vi.advanceTimersByTimeAsync(0)

        expect(status.isServerReachable.value).toBe(false)
        wrapper.unmount()
    })

    it('restores isServerReachable=true on successful recheck', async () => {
        mockFetch.mockRejectedValueOnce(new Error('down'))
        const { status, wrapper } = await mountHelper()
        await vi.advanceTimersByTimeAsync(0)
        expect(status.isServerReachable.value).toBe(false)

        mockFetch.mockResolvedValueOnce({ ok: true })
        status.recheck()
        await vi.advanceTimersByTimeAsync(0)

        expect(status.isServerReachable.value).toBe(true)
        wrapper.unmount()
    })

    it('updates lastCheck after ping', async () => {
        mockFetch.mockResolvedValue({ ok: true })
        const { status, wrapper } = await mountHelper()

        await vi.advanceTimersByTimeAsync(0)

        expect(status.lastCheck.value).toBeInstanceOf(Date)
        wrapper.unmount()
    })

    it('stops interval on unmount (no further fetches)', async () => {
        mockFetch.mockResolvedValue({ ok: true })
        const { wrapper } = await mountHelper()
        await vi.advanceTimersByTimeAsync(0)

        const callCount = mockFetch.mock.calls.length
        wrapper.unmount()

        await vi.advanceTimersByTimeAsync(60_000)
        expect(mockFetch.mock.calls.length).toBe(callCount)
    })


    it('does not re-notify when fetch fails while already unreachable', async () => {
        mockFetch.mockRejectedValueOnce(new Error('down'))
        const { status, wrapper } = await mountHelper()
        await vi.advanceTimersByTimeAsync(0)
        expect(status.isServerReachable.value).toBe(false)

        mockFetch.mockRejectedValueOnce(new Error('still down'))
        await status.recheck()
        await vi.advanceTimersByTimeAsync(0)

        expect(status.isServerReachable.value).toBe(false)
        wrapper.unmount()
    })
})