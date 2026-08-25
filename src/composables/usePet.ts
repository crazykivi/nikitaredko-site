import { reactive, watch } from 'vue'

export interface PetStats {
  enabled: boolean
  name: string
  mood: number
  energy: number
  hunger: number
  xp: number
  level: number
  visits: number
  lastSeenAt: number
  lastFedAt: number
  lastPetAt: number
}

const STORAGE_KEY = 'nikitaredko.pet.v1'

function clamp(value: number, min = 0, max = 100): number {
  if (Number.isNaN(value)) return min
  return Math.min(max, Math.max(min, value))
}

function createDefaultStats(): PetStats {
  return {
    enabled: false,
    name: 'Кодик',
    mood: 72,
    energy: 78,
    hunger: 64,
    xp: 0,
    level: 1,
    visits: 0,
    lastSeenAt: Date.now(),
    lastFedAt: 0,
    lastPetAt: 0,
  }
}

function sanitizeStats(value: Record<string, unknown>): PetStats {
  const d = createDefaultStats()

  const toNumber = (input: unknown, fallback: number): number => {
    const num = Number(input)
    return Number.isFinite(num) ? num : fallback
  }

  const xp = Math.max(0, toNumber(value.xp, d.xp))

  return {
    enabled: typeof value.enabled === 'boolean' ? value.enabled : d.enabled,
    name: typeof value.name === 'string' && value.name.trim() ? value.name : d.name,
    mood: clamp(toNumber(value.mood, d.mood)),
    energy: clamp(toNumber(value.energy, d.energy)),
    hunger: clamp(toNumber(value.hunger, d.hunger)),
    xp,
    level: Math.max(1, Math.floor(xp / 60) + 1),
    visits: Math.max(0, Math.floor(toNumber(value.visits, d.visits))),
    lastSeenAt: Math.max(0, toNumber(value.lastSeenAt, Date.now())),
    lastFedAt: Math.max(0, toNumber(value.lastFedAt, 0)),
    lastPetAt: Math.max(0, toNumber(value.lastPetAt, 0)),
  }
}

function applyOfflineDecay(stats: PetStats): void {
  const now = Date.now()
  const elapsedMs = Math.max(0, now - stats.lastSeenAt)

  if (elapsedMs < 5 * 60_000) {
    stats.lastSeenAt = now
    return
  }

  const minutes = Math.min(elapsedMs / 60_000, 60 * 24)

  stats.hunger = clamp(stats.hunger - Math.min(65, minutes * 0.35))

  const moodPenalty = Math.min(
    45,
    minutes * 0.12 + (stats.hunger < 35 ? 12 : 0),
  )

  stats.mood = clamp(stats.mood - moodPenalty)
  stats.energy = clamp(stats.energy + Math.min(40, minutes * 0.2))
  stats.lastSeenAt = now
}

function loadStats(): PetStats {
  if (typeof window === 'undefined') {
    return createDefaultStats()
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)

    if (!raw) {
      return createDefaultStats()
    }

    const parsed: unknown = JSON.parse(raw)

    if (typeof parsed !== 'object' || parsed === null) {
      return createDefaultStats()
    }

    const stats = sanitizeStats(parsed as Record<string, unknown>)

    if (stats.enabled) {
      applyOfflineDecay(stats)
    }

    return stats
  } catch {
    return createDefaultStats()
  }
}

export const petStats = reactive<PetStats>(loadStats())

let saveTimer: ReturnType<typeof setTimeout> | null = null

function scheduleSave(): void {
  if (typeof window === 'undefined') return

  if (saveTimer !== null) {
    clearTimeout(saveTimer)
  }

  saveTimer = setTimeout(() => {
    try {
      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(petStats))
    } catch {
      // localStorage может быть недоступен в приватном режиме или переполнен
    }
  }, 400)
}

watch(petStats, scheduleSave, { deep: true })

if (typeof window !== 'undefined') {
  window.addEventListener('storage', (event) => {
    if (event.key !== STORAGE_KEY || event.newValue === null) return

    try {
      const parsed: unknown = JSON.parse(event.newValue)

      if (typeof parsed === 'object' && parsed !== null) {
        Object.assign(petStats, sanitizeStats(parsed as Record<string, unknown>))
      }
    } catch {
      // игнорирование битых данных из другой вкладки
    }
  })
}

function recalcLevel(): void {
  petStats.level = Math.floor(petStats.xp / 60) + 1
}

function addXp(amount: number): void {
  petStats.xp = Math.max(0, petStats.xp + Math.round(amount))
  recalcLevel()
}

export function startPetSession(): void {
  if (!petStats.enabled) return

  petStats.visits += 1
  petStats.lastSeenAt = Date.now()

  scheduleSave()
}

export function togglePet(): void {
  petStats.enabled = !petStats.enabled

  if (petStats.enabled) {
    petStats.lastSeenAt = Date.now()
    petStats.mood = clamp(petStats.mood + 8)
    petStats.energy = clamp(petStats.energy + 8)
  }

  scheduleSave()
}

export function feedPet(): { ok: boolean; message: string } {
  const now = Date.now()

  if (now - petStats.lastFedAt < 20_000) {
    return {
      ok: false,
      message: 'Я только что ел. Подожди чуть-чуть',
    }
  }

  petStats.hunger = clamp(petStats.hunger + 30)
  petStats.mood = clamp(petStats.mood + 6)
  petStats.energy = clamp(petStats.energy + 3)
  petStats.lastFedAt = now

  addXp(10)

  return {
    ok: true,
    message: 'Ням! Спасибо',
  }
}

export function petPet(): { ok: boolean; message: string } {
  const now = Date.now()

  if (now - petStats.lastPetAt < 5_000) {
    return {
      ok: false,
      message: 'Хватит меня тискать!',
    }
  }

  petStats.mood = clamp(petStats.mood + 8)
  petStats.lastPetAt = now

  addXp(5)

  return {
    ok: true,
    message: 'Мурр... то есть бип!',
  }
}

export function tickPet(deltaMs: number): void {
  if (!petStats.enabled || typeof document !== 'undefined' && document.hidden) {
    return
  }

  const minutes = Math.min(Math.max(deltaMs, 0) / 60_000, 5)

  if (minutes <= 0) return

  petStats.hunger = clamp(petStats.hunger - minutes * 1.1)
  petStats.energy = clamp(petStats.energy - minutes * 0.45)

  if (petStats.hunger < 25) {
    petStats.mood = clamp(petStats.mood - minutes * 1.2)
  } else if (petStats.energy < 20) {
    petStats.mood = clamp(petStats.mood - minutes * 0.5)
  } else {
    petStats.mood = clamp(petStats.mood + minutes * 0.18)
  }

  petStats.lastSeenAt = Date.now()

  if (Math.random() < 0.35) {
    addXp(1)
  }
}