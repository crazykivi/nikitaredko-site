<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useTheme } from '../composables/useTheme'
import { mcBlockTransition } from '../utils/mcTransition'
import { useBlockBreak } from '../composables/useBlockBreak'
import { playMcSound } from '../utils/mcSounds'
import { unlockAudio } from '../utils/digSound'

const HOLD_MS = 2000
const CANCEL_PX = 12
const NON_BACKGROUND = [
    'a', 'button', 'input', 'textarea', 'select', 'summary', 'label',
    'p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'li', 'pre', 'code', 'time',
    'header', 'footer', 'nav', 'aside', 'article',
    '[role="button"]', '[role="link"]', '[role="dialog"]',
    '.prose', '.fixed', '.mc-transition', '[data-no-break]',
].join(', ')

const CRACK_PATHS = [
    'M50 50 L41 42',
    'M50 50 L60 40 L68 44',
    'M50 50 L42 61',
    'M41 42 L30 36 M68 44 L74 32',
    'M42 61 L32 70 L24 66',
    'M50 50 L59 59 L70 62',
    'M30 36 L20 26 M74 32 L80 20',
    'M32 70 L28 82 M70 62 L82 68',
    'M50 50 L47 30 L41 18 M59 59 L64 78',
    'M24 66 L12 60 M80 20 L88 12 M28 82 L20 92 M82 68 L92 74',
]
const STAGES = CRACK_PATHS.length
const STEP_MS = HOLD_MS / STAGES
const { unlocked } = useBlockBreak()

const { mode, effectiveDark, enterCharcoal, exitCharcoal } = useTheme()

const active = ref(false)
const healing = ref(false)
const x = ref(0)
const y = ref(0)
const stage = ref(0)
const particles = ref<Array<{ id: number; dx: number; dy: number; color: string; delay: number }>>([])

const visible = computed(() => active.value || healing.value || particles.value.length > 0)

let pointerId: number | null = null
let startX = 0
let startY = 0
let stepTimer: number | null = null
let healTimer: number | null = null
let particleTimer: number | null = null
let nextParticleId = 0

const stopStepTimer = () => {
    if (stepTimer !== null) {
        clearInterval(stepTimer)
        stepTimer = null
    }
}

const shakeScreen = () => {
    const root = document.documentElement
    root.classList.add('bb-shake')
    window.setTimeout(() => root.classList.remove('bb-shake'), 450)
}

const spawnParticles = () => {
    const palette =
        mode.value === 'charcoal'
            ? ['#6f6f6f', '#8b8b8b', '#55ff55']
            : effectiveDark.value
                ? ['#27272a', '#3f3f46', '#52525b']
                : ['#d4d4d8', '#a1a1aa', '#71717a']
    particles.value = Array.from({ length: 12 }, (_, i) => {
        const angle = Math.random() * Math.PI * 2
        const dist = 50 + Math.random() * 60
        return {
            id: nextParticleId++,
            dx: Math.cos(angle) * dist,
            dy: Math.sin(angle) * dist - 20,
            color: palette[i % palette.length],
            delay: i * 8,
        }
    })
    if (particleTimer !== null) clearTimeout(particleTimer)
    particleTimer = window.setTimeout(() => (particles.value = []), 800)
}

const complete = () => {
    stopStepTimer()
    active.value = false
    pointerId = null
    document.documentElement.classList.remove('is-breaking')
    playMcSound('break')
    spawnParticles()
    shakeScreen()
    mcBlockTransition(() => {
        if (mode.value === 'charcoal') exitCharcoal()
        else enterCharcoal()
    })
}

const cancel = () => {
    if (!active.value) return
    stopStepTimer()
    active.value = false
    healing.value = true
    pointerId = null
    document.documentElement.classList.remove('is-breaking')
    if (healTimer !== null) clearTimeout(healTimer)
    healTimer = window.setTimeout(() => {
        healing.value = false
        stage.value = 0
    }, 350)
}

const step = () => {
    stage.value++
    if (stage.value >= STAGES) {
        complete()
        return
    }
    playMcSound('dig', { volume: 0.45, rate: 1.4 })
}

const onPointerDown = (e: PointerEvent) => {
    const canBreak = unlocked.value || mode.value === 'charcoal'
    if (!canBreak) return
    if (active.value || !e.isPrimary) return
    if (e.pointerType === 'mouse' && e.button !== 0) return
    if (document.querySelector('[role="dialog"]')) return
    const target = e.target as Element | null
    if (target && target.closest(NON_BACKGROUND)) return
    if (healTimer !== null) {
        clearTimeout(healTimer)
        healTimer = null
        healing.value = false
    }
    unlockAudio()
    pointerId = e.pointerId
    startX = e.clientX
    startY = e.clientY
    x.value = e.clientX
    y.value = e.clientY
    stage.value = 0
    particles.value = []
    active.value = true
    document.documentElement.classList.add('is-breaking')
    playMcSound('dig', { volume: 0.45, rate: 1.4 })
    stopStepTimer()
    stepTimer = window.setInterval(step, STEP_MS)
}

const onPointerMove = (e: PointerEvent) => {
    if (e.pointerId !== pointerId) return
    if (Math.hypot(e.clientX - startX, e.clientY - startY) > CANCEL_PX) cancel()
}

const onPointerEnd = (e: PointerEvent) => {
    if (e.pointerId !== pointerId) return
    cancel()
}

const onWheel = () => cancel()
const onKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') cancel()
}
const onContextMenu = (e: MouseEvent) => {
    if (pointerId !== null) e.preventDefault()
}

onMounted(() => {
    window.addEventListener('pointerdown', onPointerDown)
    window.addEventListener('pointermove', onPointerMove)
    window.addEventListener('pointerup', onPointerEnd)
    window.addEventListener('pointercancel', onPointerEnd)
    window.addEventListener('wheel', onWheel, { passive: true })
    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('contextmenu', onContextMenu, true)
})

onUnmounted(() => {
    window.removeEventListener('pointerdown', onPointerDown)
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', onPointerEnd)
    window.removeEventListener('pointercancel', onPointerEnd)
    window.removeEventListener('wheel', onWheel)
    window.removeEventListener('keydown', onKeyDown)
    window.removeEventListener('contextmenu', onContextMenu, true)
    stopStepTimer()
    if (healTimer !== null) clearTimeout(healTimer)
    if (particleTimer !== null) clearTimeout(particleTimer)
})
</script>

<template>
    <div v-if="visible" class="bb-overlay" :class="{ 'bb-healing': healing }" :style="{ left: `${x}px`, top: `${y}px` }"
        aria-hidden="true">
        <template v-if="active || healing">
            <span :key="stage" class="bb-frame" />
            <svg class="bb-cracks" viewBox="0 0 100 100">
                <path v-for="(d, i) in CRACK_PATHS" :key="i" :d="d" :class="{ 'bb-on': i < stage }"
                    :style="{ transitionDelay: healing ? `${(STAGES - i) * 25}ms` : '0ms' }" />
            </svg>
        </template>
        <i v-for="p in particles" :key="p.id" class="bb-particle" :style="{
            '--dx': `${p.dx}px`,
            '--dy': `${p.dy}px`,
            '--c': p.color,
            animationDelay: `${p.delay}ms`,
        }" />
    </div>
</template>