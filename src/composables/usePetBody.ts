import { reactive, ref } from 'vue'

export interface PetBounds {
  minX: number
  maxX: number
  minY: number
  ground: number
}

const GRAVITY = 1400
const RESTITUTION = 0.55
const STOP_SPEED = 60
const MAX_THROW = 1500

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value))
}

export function usePetBody() {
  const pos = reactive({ x: 0, y: 0 })
  const vel = reactive({ x: 0, y: 0 })
  const spin = ref(0)
  const squash = ref(0)
  const isThrown = ref(false)

  let squashTimer: ReturnType<typeof setTimeout> | null = null
  let impactHandler: ((strength: number) => void) | null = null
  let landHandler: (() => void) | null = null

  function setHandlers(
    onImpact: (strength: number) => void,
    onLand: () => void,
  ): void {
    impactHandler = onImpact
    landHandler = onLand
  }

  function triggerSquash(): void {
    squash.value = 1

    if (squashTimer !== null) {
      clearTimeout(squashTimer)
    }

    squashTimer = setTimeout(() => {
      squash.value = 0
    }, 140)
  }

  function throwBody(vx: number, vy: number): void {
    vel.x = clamp(vx, -MAX_THROW, MAX_THROW)
    vel.y = clamp(vy, -MAX_THROW, MAX_THROW)
    isThrown.value = true
  }

  function cancelThrow(): void {
    isThrown.value = false
    vel.x = 0
    vel.y = 0
    spin.value = 0
  }

  function step(dt: number, bounds: PetBounds): void {
    if (!isThrown.value) return

    const gravityMult = vel.y < 0 ? 0.6 : 1.0
    vel.y += GRAVITY * gravityMult * dt
    vel.x *= Math.max(0, 1 - 2.0 * dt)
    vel.y *= Math.max(0, 1 - 1.0 * dt)

    pos.x += vel.x * dt
    pos.y += vel.y * dt

    const speed = Math.hypot(vel.x, vel.y)
    if (speed > 40) {
      const angle = Math.atan2(vel.y, vel.x) * (180 / Math.PI)
      let targetSpin = 90 - angle
      if (targetSpin > 180) targetSpin -= 360
      if (targetSpin < -180) targetSpin += 360
      
      spin.value = clamp(targetSpin, -180, 180)
    } else {
      spin.value *= 0.85
    }

    let impacted = 0

    if (pos.x <= bounds.minX) {
      pos.x = bounds.minX
      impacted = Math.abs(vel.x)
      vel.x = Math.abs(vel.x) * RESTITUTION
    } else if (pos.x >= bounds.maxX) {
      pos.x = bounds.maxX
      impacted = Math.abs(vel.x)
      vel.x = -Math.abs(vel.x) * RESTITUTION
    }

    if (pos.y <= bounds.minY) {
      pos.y = bounds.minY
      impacted = Math.max(impacted, Math.abs(vel.y))
      vel.y = Math.abs(vel.y) * RESTITUTION
    }

    if (pos.y >= bounds.ground) {
      pos.y = bounds.ground

      if (Math.abs(vel.y) > 160) {
        impacted = Math.max(impacted, Math.abs(vel.y))
        vel.y = -Math.abs(vel.y) * RESTITUTION
        vel.x *= 0.8
      } else {
        vel.y = 0
        vel.x *= 0.86
      }
    }

    if (impacted > 140) {
      triggerSquash()
      impactHandler?.(impacted)
    }

    if (pos.y >= bounds.ground - 0.5 && speed < STOP_SPEED) {
      isThrown.value = false
      spin.value = 0
      landHandler?.()
    }
  }

  return {
    pos,
    vel,
    spin,
    squash,
    isThrown,
    setHandlers,
    throwBody,
    cancelThrow,
    step,
  }
}

export type PetBody = ReturnType<typeof usePetBody>