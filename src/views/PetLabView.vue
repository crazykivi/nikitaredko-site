<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, computed, reactive } from 'vue'
import { useHead } from '@unhead/vue'
import PetSprite, { type PetMode, type PetFace } from '../components/PetSprite.vue'
import { usePetBody, type PetBounds } from '../composables/usePetBody'

const CAT = 120

const stageRef = ref<HTMLDivElement | null>(null)
const catRef = ref<HTMLDivElement | null>(null)

const mode = ref<PetMode>('idle')
const faceOverride = ref<'auto' | PetFace>('auto')
const facingRight = ref(true)
const movementOn = ref(false)
const speed = ref(1)
const log = ref<string[]>([])

const body = usePetBody()

const pos = reactive({ x: 60, y: 0, vx: 60 })
let raf = 0
let lastFrameAt = 0
let pointerSamples: { x: number; y: number; t: number }[] = []
let modeUntil = 0
let strongImpact = false

const MODES: PetMode[] = ['idle', 'walk', 'zoomies', 'sit', 'sleep', 'play', 'scruff', 'thrown', 'dizzy']
const FACES: PetFace[] = ['happy', 'neutral', 'sad', 'surprised', 'dizzy', 'sleep']

const face = computed<PetFace>(() => {
  if (mode.value === 'sleep') return 'sleep'
  if (mode.value === 'dizzy') return 'dizzy'
  if (mode.value === 'scruff' || mode.value === 'thrown') return 'surprised'
  return 'happy'
})

function pushLog(message: string): void {
  log.value.unshift(`${new Date().toLocaleTimeString('ru-RU')} — ${message}`)

  if (log.value.length > 8) {
    log.value.pop()
  }
}

function stageBounds(): PetBounds {
  const el = stageRef.value

  if (!el) {
    return { minX: 0, maxX: 400, minY: 20, ground: 300 }
  }

  return {
    minX: 0,
    maxX: el.clientWidth - CAT,
    minY: 20,
    ground: el.clientHeight - CAT - 8,
  }
}

function applyTransform(): void {
  if (!catRef.value) return

  if (Number.isNaN(pos.x) || Number.isNaN(pos.y)) {
    pos.x = 0
    pos.y = 0
  }

  const s = body.squash.value
  catRef.value.style.transform =
    `translate3d(${pos.x}px, ${pos.y}px, 0) ` +
    `rotate(${body.spin.value}deg) ` +
    `scale(${1 + s * 0.18}, ${1 - s * 0.22})`
}

body.setHandlers(
  (strength) => {
    if (strength > 400) {
      strongImpact = true
      mode.value = 'dizzy'
      modeUntil = performance.now() + 2500
      pushLog(`Сильный удар (${Math.round(strength)} px/s) -> Оглушён`)
    } else {
      pushLog(`Удар (${Math.round(strength)} px/s)`)
    }
  },
  () => {
    if (performance.now() >= modeUntil) {
      mode.value = strongImpact ? 'dizzy' : 'idle'
      strongImpact = false
      modeUntil = performance.now() + (mode.value === 'dizzy' ? 2400 : 1200)
    }
    pushLog('Приземлился')
  },
)

function frame(now: number): void {
  if (!lastFrameAt) lastFrameAt = now

  const dt = Math.min((now - lastFrameAt) / 1000, 0.05)
  lastFrameAt = now

  const bounds = stageBounds()

  if (body.isThrown.value) {
    if (mode.value !== 'dizzy') {
      mode.value = 'thrown'
    }
    
    body.step(dt, bounds)
    pos.x = body.pos.x
    pos.y = body.pos.y

    if (Math.abs(body.vel.x) > 60) {
      facingRight.value = body.vel.x > 0
    }
  } else if (!dragState.active) {
    if (mode.value === 'dizzy') {
    } else if (movementOn.value && (mode.value === 'walk' || mode.value === 'zoomies')) {
      const v = (mode.value === 'zoomies' ? 160 : 55) * speed.value
      pos.x += (facingRight.value ? v : -v) * dt

      if (pos.x <= bounds.minX) {
        pos.x = bounds.minX
        facingRight.value = true
      } else if (pos.x >= bounds.maxX) {
        pos.x = bounds.maxX
        facingRight.value = false
      }
    }

    pos.y += (bounds.ground - pos.y) * Math.min(1, dt * 3.5)
  }

  applyTransform()
  raf = requestAnimationFrame(frame)
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
  if (!catRef.value) return

  if (body.isThrown.value) body.cancelThrow()

  dragState.active = true
  dragState.moved = false
  dragState.pointerId = event.pointerId
  dragState.startX = event.clientX
  dragState.startY = event.clientY
  dragState.originX = pos.x
  dragState.originY = pos.y

  pointerSamples = [{ x: event.clientX, y: event.clientY, t: performance.now() }]

  catRef.value.setPointerCapture(event.pointerId)
}

function onPointerMove(event: PointerEvent): void {
  if (!dragState.active || dragState.pointerId !== event.pointerId) return

  const dx = event.clientX - dragState.startX
  const dy = event.clientY - dragState.startY

  if (Math.abs(dx) > 4 || Math.abs(dy) > 4) {
    if (!dragState.moved) {
      mode.value = 'scruff'
    }
    dragState.moved = true
  }

  const bounds = stageBounds()
  pos.x = Math.min(bounds.maxX, Math.max(bounds.minX, dragState.originX + dx))
  pos.y = Math.min(bounds.ground + 60, Math.max(0, dragState.originY + dy))

  pointerSamples.push({ x: event.clientX, y: event.clientY, t: performance.now() })

  if (pointerSamples.length > 8) pointerSamples.shift()

  applyTransform()
}

function getThrowVelocity(): { vx: number; vy: number } {
  if (pointerSamples.length < 2) return { vx: 0, vy: 0 }

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

  return { vx, vy }
}

function onPointerUp(event: PointerEvent): void {
  if (!dragState.active || dragState.pointerId !== event.pointerId) return

  dragState.active = false

  if (catRef.value && catRef.value.hasPointerCapture(event.pointerId)) {
    catRef.value.releasePointerCapture(event.pointerId)
  }

  if (!dragState.moved) return

  const bounds = stageBounds()
  const isInAir = pos.y < bounds.ground - 10
  const v = getThrowVelocity()
  const hasVelocity = Math.hypot(v.vx, v.vy) >= 250

  if (hasVelocity || isInAir) {
    body.pos.x = pos.x
    body.pos.y = pos.y
    body.throwBody(hasVelocity ? v.vx : 0, hasVelocity ? v.vy : 0)
    pushLog(isInAir && !hasVelocity ? 'Отпущен в воздухе' : 'Бросок!')
  } else {
    mode.value = 'idle'
  }
}

function setMode(m: PetMode): void {
  body.cancelThrow()
  mode.value = m
  modeUntil = performance.now() + 3000
  pushLog(`Режим: ${m}`)
}

function fling(): void {
  body.pos.x = pos.x
  body.pos.y = pos.y
  body.throwBody(
    (Math.random() > 0.5 ? 1 : -1) * (500 + Math.random() * 400),
    -(400 + Math.random() * 400),
  )
  pushLog('Шпулянут кнопкой')
}

function reset(): void {
  body.cancelThrow()
  const bounds = stageBounds()
  pos.x = bounds.maxX / 2
  pos.y = bounds.ground
  mode.value = 'idle'
  applyTransform()
  pushLog('Сброс позиции')
}

onMounted(() => {
  const bounds = stageBounds()
  pos.x = bounds.maxX / 2
  pos.y = bounds.ground
  applyTransform()
  raf = requestAnimationFrame(frame)
})

onBeforeUnmount(() => {
  if (raf) cancelAnimationFrame(raf)
})

useHead({ title: 'Пет-лаборатория | Nikita Redko' })
</script>

<template>
  <div class="min-h-[calc(100vh-7.3rem)] pt-10 pb-20 px-4">
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center gap-2 text-sm text-muted mb-4">
        <router-link to="/" class="hover:text-foreground transition-colors">Главная</router-link>
        <span>/</span>
        <span class="text-foreground">pet-lab</span>
      </div>

      <h1 class="text-3xl font-bold mb-2">Пет-лаборатория</h1>
      <p class="text-muted text-sm mb-6">
        Дебаг-песочница котика: переключай позы, таскай его мышью, шпуляй в стены, крути скорость анимаций.
      </p>

      <div class="grid gap-6 lg:grid-cols-[1fr_300px]">
        <div
          ref="stageRef"
          class="relative h-[420px] rounded-xl border border-border bg-muted/5 overflow-hidden select-none"
        >
          <div class="absolute bottom-8 left-0 right-0 border-t border-border/60" />

          <div
            ref="catRef"
            class="absolute left-0 top-0 touch-none cursor-grab active:cursor-grabbing"
            :style="{ width: `${CAT}px`, height: `${CAT}px` }"
            @pointerdown="onPointerDown"
            @pointermove="onPointerMove"
            @pointerup="onPointerUp"
            @pointercancel="onPointerUp"
          >
            <span class="block h-full w-full" :class="{ 'facing-left': !facingRight && mode !== 'thrown' }">
              <PetSprite :mode="mode" :face="faceOverride === 'auto' ? face : faceOverride" :speed="speed" />
            </span>
          </div>

          <p class="absolute top-2 left-3 text-[10px] font-mono text-muted/60">
            mode: {{ mode }} | thrown: {{ body.isThrown.value }}
          </p>
        </div>
        <div class="space-y-5">
          <section>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted mb-2">Позы</h2>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="m in MODES"
                :key="m"
                @click="setMode(m)"
                class="px-2.5 py-1.5 rounded-lg border text-xs font-mono transition-colors"
                :class="mode === m ? 'bg-foreground text-background border-foreground' : 'border-border text-muted hover:text-foreground'"
              >
                {{ m }}
              </button>
            </div>
          </section>

          <section>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted mb-2">Морда</h2>
            <div class="flex flex-wrap gap-1.5">
              <button
                @click="faceOverride = 'auto'"
                class="px-2.5 py-1.5 rounded-lg border text-xs font-mono transition-colors"
                :class="faceOverride === 'auto' ? 'bg-foreground text-background border-foreground' : 'border-border text-muted hover:text-foreground'"
              >
                auto
              </button>
              <button
                v-for="f in FACES"
                :key="f"
                @click="faceOverride = f"
                class="px-2.5 py-1.5 rounded-lg border text-xs font-mono transition-colors"
                :class="faceOverride === f ? 'bg-foreground text-background border-foreground' : 'border-border text-muted hover:text-foreground'"
              >
                {{ f }}
              </button>
            </div>
          </section>

          <section class="space-y-2">
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted mb-2">Параметры</h2>

            <label class="flex items-center justify-between text-xs text-muted">
              Скорость анимаций: <span class="font-mono text-foreground">{{ speed.toFixed(2) }}×</span>
            </label>
            <input
              v-model.number="speed"
              type="range"
              min="0.25"
              max="2"
              step="0.05"
              class="w-full"
            />

            <div class="flex flex-wrap gap-1.5 pt-1">
              <button
                @click="facingRight = !facingRight"
                class="px-2.5 py-1.5 rounded-lg border border-border text-xs text-muted hover:text-foreground transition-colors"
              >
                Смотрит: {{ facingRight ? 'вправо' : 'влево' }}
              </button>

              <button
                @click="movementOn = !movementOn"
                class="px-2.5 py-1.5 rounded-lg border text-xs transition-colors"
                :class="movementOn ? 'bg-foreground text-background border-foreground' : 'border-border text-muted hover:text-foreground'"
              >
                Движение: {{ movementOn ? 'вкл' : 'выкл' }}
              </button>
            </div>

            <div class="flex flex-wrap gap-1.5">
              <button
                @click="fling"
                class="px-2.5 py-1.5 rounded-lg bg-foreground text-background text-xs font-medium hover:opacity-85 transition-opacity"
              >
                Шпулянуть
              </button>
              <button
                @click="reset"
                class="px-2.5 py-1.5 rounded-lg border border-border text-xs text-muted hover:text-foreground transition-colors"
              >
                Сброс
              </button>
            </div>
          </section>

          <section>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted mb-2">Лог</h2>
            <div class="rounded-lg border border-border bg-muted/5 p-2 font-mono text-[10px] text-muted space-y-1 min-h-[90px]">
              <p v-for="(line, i) in log" :key="i">{{ line }}</p>
              <p v-if="log.length === 0" class="opacity-50">// пусто</p>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { watch } from 'vue'

export default { name: 'PetLabView' }
</script>

<style scoped>
.facing-left {
  transform: scaleX(-1);
}
</style>