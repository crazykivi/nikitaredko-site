<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, reactive } from 'vue'
import PetSprite, { type PetMode, type PetFace } from './PetSprite.vue'
import { usePetBody, type PetBounds } from '../composables/usePetBody'
import {
  petStats,
  startPetSession,
  tickPet,
  feedPet,
  petPet,
  togglePet,
} from '../composables/usePet'

const petButton = ref<HTMLButtonElement | null>(null)
const bubble = ref('')
const isPanelOpen = ref(false)
const currentMode = ref<PetMode>('idle')
const facingRight = ref(true)

const reducedMotion =
  typeof window !== 'undefined' &&
  window.matchMedia('(prefers-reduced-motion: reduce)').matches

const size = 96

const body = usePetBody()

let raf = 0
let tickInterval: ReturnType<typeof setInterval> | null = null
let bubbleTimer: ReturnType<typeof setTimeout> | null = null
let lastTickAt = Date.now()
let lastFrameAt = 0
let modeUntil = 0
let strongImpact = false

const pos = reactive({ x: 0, y: 0, vx: 55 })

let pointerSamples: { x: number; y: number; t: number }[] = []

const face = computed<PetFace>(() => {
  if (currentMode.value === 'sleep') return 'sleep'
  if (currentMode.value === 'dizzy') return 'dizzy'
  if (currentMode.value === 'scruff' || currentMode.value === 'thrown') return 'surprised'
  if (petStats.mood >= 70) return 'happy'
  if (petStats.mood >= 35) return 'neutral'
  return 'sad'
})

function clampValue(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value))
}

function random(min: number, max: number): number {
  return min + Math.random() * (max - min)
}

function showMessage(message: string): void {
  bubble.value = message

  if (bubbleTimer !== null) {
    clearTimeout(bubbleTimer)
  }

  bubbleTimer = setTimeout(() => {
    bubble.value = ''
  }, 2500)
}

function getBounds(): PetBounds {
  return {
    minX: 8,
    maxX: window.innerWidth - size - 8,
    minY: 76,
    ground: groundY(),
  }
}

function applyTransform(): void {
  if (!petButton.value) return

  if (Number.isNaN(pos.x) || Number.isNaN(pos.y)) {
    pos.x = 8
    pos.y = groundY()
  }

  const s = body.squash.value
  const transform =
    `translate3d(${pos.x}px, ${pos.y}px, 0) ` +
    `rotate(${body.spin.value}deg) ` +
    `scale(${1 + s * 0.18}, ${1 - s * 0.22})`

  petButton.value.style.transform = transform
}

function groundY(): number {
  return window.innerHeight - size - 12
}

body.setHandlers(
  (strength) => {
    if (strength > 400) {
      strongImpact = true
      currentMode.value = 'dizzy'
      modeUntil = performance.now() + 2500
      showMessage(strength > 600 ? 'Бум! Ай! Звёздочки...' : 'Ой! Стена...')
    } else {
      showMessage('Ой!')
    }
  },
  () => {
    if (performance.now() >= modeUntil) {
      currentMode.value = strongImpact ? 'dizzy' : 'idle'
      strongImpact = false
      modeUntil = performance.now() + (currentMode.value === 'dizzy' ? 2400 : 1200)
    }

    if (currentMode.value === 'dizzy') {
      showMessage('Мяу... звёздочки...')
    }
  },
)

function chooseMode(now: number): void {
  if (petStats.energy < 14) {
    currentMode.value = 'sleep'
    modeUntil = now + random(8000, 16000)
    return
  }

  const roll = Math.random()

  if (petStats.energy < 40 && roll < 0.1) {
    currentMode.value = 'sleep'
    modeUntil = now + random(6000, 12000)
    return
  }

  if (petStats.energy < 30) {
    currentMode.value = roll < 0.6 ? 'sit' : 'idle'
    modeUntil = now + random(5000, 10000)
    return
  }

  if (roll < 0.34) {
    currentMode.value = 'walk'
    const goRight = Math.random() > 0.5
    pos.vx = goRight ? 55 : -55
    facingRight.value = goRight
    modeUntil = now + random(5000, 9000)
  } else if (roll < 0.5) {
    currentMode.value = 'zoomies'
    const goRight = Math.random() > 0.5
    pos.vx = goRight ? 160 : -160
    facingRight.value = goRight
    modeUntil = now + random(2500, 4000)
  } else if (roll < 0.68) {
    currentMode.value = 'sit'
    modeUntil = now + random(5000, 10000)
  } else if (roll < 0.86) {
    currentMode.value = 'play'
    modeUntil = now + random(4000, 7000)
  } else {
    currentMode.value = 'idle'
    modeUntil = now + random(3000, 6000)
  }
}

function updatePosition(now: number, dt: number): void {
  if (currentMode.value === 'scruff') return

  if (now >= modeUntil) {
    chooseMode(now)
  }

  const { minX, maxX } = getBounds()

  if (currentMode.value === 'walk' || currentMode.value === 'zoomies') {
    pos.x += pos.vx * dt

    if (pos.x <= minX) {
      pos.x = minX
      pos.vx = Math.abs(pos.vx)
      facingRight.value = true
    } else if (pos.x >= maxX) {
      pos.x = maxX
      pos.vx = -Math.abs(pos.vx)
      facingRight.value = false
    }
  }

  pos.y += (groundY() - pos.y) * Math.min(1, dt * 3.5)
  
  if (Number.isNaN(pos.x)) pos.x = minX
  if (Number.isNaN(pos.y)) pos.y = groundY()
  
  pos.x = clampValue(pos.x, minX, maxX)
  pos.y = clampValue(pos.y, 76, window.innerHeight - size)

  applyTransform()
}

function frame(now: number): void {
  if (!lastFrameAt) {
    lastFrameAt = now
  }

  const dt = Math.min((now - lastFrameAt) / 1000, 0.05)
  lastFrameAt = now

  if (!document.hidden && !reducedMotion) {
    if (body.isThrown.value) {
      if (currentMode.value !== 'dizzy') {
        currentMode.value = 'thrown'
      }
      
      body.step(dt, getBounds())
      pos.x = body.pos.x
      pos.y = body.pos.y

      if (Math.abs(body.vel.x) > 60) {
        facingRight.value = body.vel.x > 0
      }

      applyTransform()
    } else {
      updatePosition(now, dt)
    }
  }

  raf = requestAnimationFrame(frame)
}

function onResize(): void {
  if (body.isThrown.value) return

  pos.x = clampValue(pos.x, 8, window.innerWidth - size - 8)
  pos.y = groundY()
  applyTransform()
}

function onVisibilityChange(): void {
  if (!document.hidden) {
    lastTickAt = Date.now()
    lastFrameAt = 0
  }
}

const dragState = {
  active: false,
  moved: false,
  pointerId: -1,
  startX: 0,
  startY: 0,
  originX: 0,
  originY: 0,
}

function onPointerDown(event: PointerEvent): void {
  if (!petButton.value) return

  if (body.isThrown.value) {
    body.cancelThrow()
  }

  dragState.active = true
  dragState.moved = false
  dragState.pointerId = event.pointerId
  dragState.startX = event.clientX
  dragState.startY = event.clientY
  dragState.originX = pos.x
  dragState.originY = pos.y

  pointerSamples = [{ x: event.clientX, y: event.clientY, t: performance.now() }]

  petButton.value.setPointerCapture(event.pointerId)
}

function onPointerMove(event: PointerEvent): void {
  if (!dragState.active || dragState.pointerId !== event.pointerId) return

  const dx = event.clientX - dragState.startX
  const dy = event.clientY - dragState.startY

  if (Math.abs(dx) > 4 || Math.abs(dy) > 4) {
    if (!dragState.moved) {
      currentMode.value = 'scruff'
      showMessage('Мяу!')
    }
    dragState.moved = true
  }

  pos.x = clampValue(dragState.originX + dx, 8, window.innerWidth - size - 8)
  pos.y = clampValue(dragState.originY + dy, 76, window.innerHeight - size)

  pointerSamples.push({ x: event.clientX, y: event.clientY, t: performance.now() })

  if (pointerSamples.length > 8) {
    pointerSamples.shift()
  }

  applyTransform()
}

function throwVelocity(): { vx: number; vy: number } | null {
  if (pointerSamples.length < 2) return null

  const last = pointerSamples[pointerSamples.length - 1]
  let first = last

  for (let i = pointerSamples.length - 1; i >= 0; i--) {
    if (last.t - pointerSamples[i].t <= 150) {
      first = pointerSamples[i]
    } else {
      break
    }
  }

  const dtMs = Math.max(16, last.t - first.t)
  const vx = ((last.x - first.x) / dtMs) * 1000
  const vy = ((last.y - first.y) / dtMs) * 1000

  if (Math.hypot(vx, vy) < 250) return null

  return { vx, vy }
}

function onPointerUp(event: PointerEvent): void {
  if (!dragState.active || dragState.pointerId !== event.pointerId) return

  dragState.active = false

  if (petButton.value && petButton.value.hasPointerCapture(event.pointerId)) {
    petButton.value.releasePointerCapture(event.pointerId)
  }

  if (!dragState.moved) {
    currentMode.value = 'idle'
    modeUntil = performance.now() + 2000
    isPanelOpen.value = !isPanelOpen.value
    return
  }

  const v = throwVelocity()
  const isInAir = pos.y < groundY() - 10

  if (v || isInAir) {
    body.pos.x = pos.x
    body.pos.y = pos.y
    body.throwBody(v ? v.vx : 0, v ? v.vy : 0)
  } else {
    currentMode.value = 'idle'
    modeUntil = performance.now() + 2000
  }
}

function onPointerCancel(): void {
  dragState.active = false
  currentMode.value = 'idle'
  modeUntil = performance.now() + 2000
}

function onKeydown(event: KeyboardEvent): void {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    isPanelOpen.value = !isPanelOpen.value
  }
}

function onFeed(): void {
  const result = feedPet()
  showMessage(result.message)

  if (result.ok) {
    currentMode.value = 'sit'
    modeUntil = performance.now() + 4000
  }
}

function onPet(): void {
  const result = petPet()
  showMessage(result.message)
}

function onHide(): void {
  isPanelOpen.value = false
  togglePet()
}

function barStyle(value: number): { width: string } {
  return {
    width: `${Math.min(100, Math.max(0, Math.round(value)))}%`,
  }
}

onMounted(() => {
  startPetSession()

  pos.x = Math.max(8, Math.min(window.innerWidth / 2, window.innerWidth - size - 8))
  pos.y = groundY()

  applyTransform()

  if (!reducedMotion) {
    modeUntil = performance.now() + 1200
    raf = requestAnimationFrame(frame)
  } else {
    currentMode.value = 'idle'
  }

  tickInterval = setInterval(() => {
    const now = Date.now()
    tickPet(now - lastTickAt)
    lastTickAt = now
  }, 30_000)

  window.addEventListener('resize', onResize)
  document.addEventListener('visibilitychange', onVisibilityChange)

  showMessage('Мяу! Клик — открыть меню')
})

onBeforeUnmount(() => {
  if (raf) cancelAnimationFrame(raf)

  if (tickInterval !== null) {
    clearInterval(tickInterval)
  }

  if (bubbleTimer !== null) {
    clearTimeout(bubbleTimer)
  }

  window.removeEventListener('resize', onResize)
  document.removeEventListener('visibilitychange', onVisibilityChange)
})
</script>

<template>
  <button
    ref="petButton"
    type="button"
    aria-label="Котик-компаньон"
    class="fixed left-0 top-0 z-40 select-none outline-none touch-none"
    :style="{
      width: `${size}px`,
      height: `${size}px`,
      transform: 'translate3d(-300px, -300px, 0)',
    }"
    @pointerdown="onPointerDown"
    @pointermove="onPointerMove"
    @pointerup="onPointerUp"
    @pointercancel="onPointerCancel"
    @keydown="onKeydown"
  >
    <span
      v-if="bubble"
      class="pointer-events-none absolute bottom-full left-1/2 mb-2 w-max max-w-56 rounded-xl border border-border bg-background/95 px-3 py-2 text-xs text-foreground shadow-xl"
      :style="{ transform: `translateX(-50%) rotate(${-body.spin.value}deg)` }"
    >
      {{ bubble }}
    </span>

    <span class="block h-full w-full" :class="{ 'facing-left': !facingRight && currentMode !== 'thrown' }">
      <PetSprite :mode="currentMode" :face="face" />
    </span>
  </button>

  <div
    v-if="isPanelOpen"
    class="fixed bottom-28 right-4 z-40 w-64 rounded-2xl border border-border bg-background/95 p-4 shadow-2xl backdrop-blur-sm"
  >
    <div class="mb-3 flex items-center justify-between gap-2">
      <p class="text-sm font-semibold">{{ petStats.name }}</p>
      <span class="text-xs text-muted">Ур. {{ petStats.level }}</span>
    </div>

    <div class="space-y-3 text-xs">
      <div>
        <div class="mb-1 flex items-center justify-between text-muted">
          <span>Настроение</span>
          <span>{{ Math.round(petStats.mood) }}</span>
        </div>
        <div class="h-1.5 overflow-hidden rounded-full bg-muted/30">
          <div
            class="h-full rounded-full bg-foreground"
            :style="barStyle(petStats.mood)"
          />
        </div>
      </div>

      <div>
        <div class="mb-1 flex items-center justify-between text-muted">
          <span>Сытость</span>
          <span>{{ Math.round(petStats.hunger) }}</span>
        </div>
        <div class="h-1.5 overflow-hidden rounded-full bg-muted/30">
          <div
            class="h-full rounded-full bg-foreground"
            :style="barStyle(petStats.hunger)"
          />
        </div>
      </div>

      <div>
        <div class="mb-1 flex items-center justify-between text-muted">
          <span>Энергия</span>
          <span>{{ Math.round(petStats.energy) }}</span>
        </div>
        <div class="h-1.5 overflow-hidden rounded-full bg-muted/30">
          <div
            class="h-full rounded-full bg-foreground"
            :style="barStyle(petStats.energy)"
          />
        </div>
      </div>
    </div>

    <div class="mt-4 grid grid-cols-2 gap-2">
      <button
        type="button"
        class="rounded-lg border border-border px-3 py-2 text-xs font-medium transition-colors hover:bg-muted/40"
        @click="onPet"
      >
        Погладить
      </button>

      <button
        type="button"
        class="rounded-lg border border-border px-3 py-2 text-xs font-medium transition-colors hover:bg-muted/40"
        @click="onFeed"
      >
        Покормить
      </button>

      <button
        type="button"
        class="rounded-lg border border-border px-3 py-2 text-xs font-medium transition-colors hover:bg-muted/40"
        @click="isPanelOpen = false"
      >
        Закрыть
      </button>

      <button
        type="button"
        class="rounded-lg border border-border px-3 py-2 text-xs font-medium text-muted transition-colors hover:bg-muted/40 hover:text-foreground"
        @click="onHide"
      >
        Спрятать
      </button>
    </div>
  </div>
</template>

<style scoped>
.facing-left {
  transform: scaleX(-1);
}
</style>