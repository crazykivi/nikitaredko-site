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

// Пиксельная карта трещин на сетке 16×16 — как текстуры destroy_stage в Minecraft.
// [x, y, stage] — пиксель появляется на указанной стадии разрушения (1..9).
// Стадии растут из центра наружу ветвями, в конце добавляются отдельные сколы.
const CRACK_PIXELS: Array<[number, number, number]> = [
  // стадия 1 — первая отметка в центре
  [8, 8, 1], [7, 7, 1],
  // стадия 2 — маленький уголок
  [9, 9, 2], [9, 7, 2], [7, 9, 2], [8, 7, 2],
  // стадия 3 — трещина начинает расходиться
  [6, 6, 3], [10, 10, 3], [8, 6, 3], [8, 10, 3], [7, 8, 3],
  // стадия 4
  [5, 6, 4], [11, 10, 4], [8, 5, 4], [8, 11, 4], [6, 8, 4], [10, 8, 4], [8, 9, 4],
  // стадия 5
  [4, 5, 5], [12, 11, 5], [7, 4, 5], [9, 12, 5], [5, 8, 5], [11, 8, 5], [10, 7, 5], [6, 9, 5],
  // стадия 6
  [3, 4, 6], [13, 12, 6], [7, 3, 6], [9, 13, 6], [4, 7, 6], [12, 9, 6], [5, 10, 6], [11, 6, 6],
  // стадия 7 — боковые ответвления
  [2, 3, 7], [12, 5, 7], [4, 11, 7], [13, 9, 7],
  [5, 5, 7], [11, 5, 7], [5, 11, 7], [11, 11, 7],
  [6, 4, 7], [10, 12, 7], [4, 6, 7], [12, 10, 7],
  // стадия 8 — трещины доходят до краёв
  [2, 2, 8], [12, 4, 8], [4, 12, 8], [14, 10, 8],
  [6, 2, 8], [10, 14, 8], [2, 6, 8], [13, 3, 8],
  [3, 13, 8], [13, 13, 8], [4, 4, 8], [12, 12, 8],
  // стадия 9 — финальные сколы
  [3, 3, 9], [13, 5, 9], [3, 12, 9], [14, 11, 9],
  [5, 3, 9], [11, 13, 9], [3, 5, 9], [13, 11, 9],
  [6, 13, 9], [10, 2, 9], [2, 9, 9], [14, 7, 9],
]

const STAGES = 10
const STEP_MS = HOLD_MS / STAGES

const { unlocked } = useBlockBreak()
const { mode, effectiveDark, enterCharcoal, exitCharcoal } = useTheme()

const active = ref(false)
const healing = ref(false)
const x = ref(0)
const y = ref(0)
const stage = ref(0)
const particles = ref<Array<{ id: number; dx: number; dy: number; color: string; delay: number }>>([])
const visiblePixels = computed(() => CRACK_PIXELS.filter((p) => p[2] <= stage.value))
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
  <div
    v-if="visible"
    class="bb-overlay"
    :class="{ 'bb-healing': healing }"
    :style="{ left: `${x}px`, top: `${y}px` }"
    aria-hidden="true"
  >
    <template v-if="active || healing">
      <span :key="stage" class="bb-frame" />
      <!-- Пиксельные трещины: сетка 16×16, каждый пиксель — rect -->
      <svg class="bb-cracks" viewBox="0 0 16 16" shape-rendering="crispEdges">
        <rect
          v-for="p in visiblePixels"
          :key="p[0] + '-' + p[1]"
          :x="p[0]"
          :y="p[1]"
          width="1"
          height="1"
          :style="{ transitionDelay: healing ? `${(STAGES - p[2]) * 25}ms` : '0ms' }"
        />
      </svg>
    </template>
    <i
      v-for="p in particles"
      :key="p.id"
      class="bb-particle"
      :style="{
        '--dx': `${p.dx}px`,
        '--dy': `${p.dy}px`,
        '--c': p.color,
        animationDelay: `${p.delay}ms`,
      }"
    />
  </div>
</template>