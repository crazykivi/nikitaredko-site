import { describe, it, expect, vi, beforeEach } from 'vitest'

describe('useNetworkStatus', () => {
    beforeEach(() => {
        vi.resetModules()
    })

    const mountHelper = async () => {
        const { useNetworkStatus } = await import('../useNetworkStatus')
        const { defineComponent, h, nextTick } = await import('vue')
        const { mount } = await import('@vue/test-utils')

        let status!: ReturnType<typeof useNetworkStatus>
        const Comp = defineComponent({
            setup() {
                status = useNetworkStatus()
                return () => h('div')
            },
        })

        const wrapper = mount(Comp)
        return { status, wrapper, nextTick }
    }

    it('reflects initial navigator.onLine (true in happy-dom)', async () => {
        const { status, wrapper } = await mountHelper()
        expect(status.isOnline.value).toBe(true)
        expect(status.lastTransition.value).toBeNull()
        wrapper.unmount()
    })

    it('sets isOnline=false and lastTransition=offline on offline event', async () => {
        const { status, wrapper, nextTick } = await mountHelper()

        window.dispatchEvent(new Event('offline'))
        await nextTick()

        expect(status.isOnline.value).toBe(false)
        expect(status.lastTransition.value).toBe('offline')
        wrapper.unmount()
    })

    it('restores isOnline=true on online event', async () => {
        const { status, wrapper, nextTick } = await mountHelper()

        window.dispatchEvent(new Event('offline'))
        await nextTick()
        expect(status.isOnline.value).toBe(false)

        window.dispatchEvent(new Event('online'))
        await nextTick()
        expect(status.isOnline.value).toBe(true)
        expect(status.lastTransition.value).toBe('online')
        wrapper.unmount()
    })

    it('ignores duplicate events (no state change)', async () => {
        const { status, wrapper, nextTick } = await mountHelper()
        window.dispatchEvent(new Event('online'))
        await nextTick()
        expect(status.lastTransition.value).toBeNull()
        expect(status.isOnline.value).toBe(true)
        wrapper.unmount()
    })

    it('does not register duplicate listeners on second mount', async () => {
        const { useNetworkStatus } = await import('../useNetworkStatus')
        const { defineComponent, h, nextTick } = await import('vue')
        const { mount } = await import('@vue/test-utils')

        let status1!: ReturnType<typeof useNetworkStatus>
        let status2!: ReturnType<typeof useNetworkStatus>

        const Comp = defineComponent({
            setup() {
                status1 = useNetworkStatus()
                status2 = useNetworkStatus()
                return () => h('div')
            },
        })

        const wrapper = mount(Comp)

        window.dispatchEvent(new Event('offline'))
        await nextTick()
        expect(status1.isOnline.value).toBe(false)
        expect(status2.isOnline.value).toBe(false)
        wrapper.unmount()
    })
})