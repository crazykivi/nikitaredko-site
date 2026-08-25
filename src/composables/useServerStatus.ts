import { ref, onMounted, onUnmounted } from 'vue'

const isServerReachable = ref(true)
const lastCheck = ref<Date | null>(null)
const listeners = new Set<() => void>()

let pingTimer: ReturnType<typeof setInterval> | null = null
let inFlight: AbortController | null = null

const HEALTH_URL = '/api/cache/health'
const PING_INTERVAL = 30_000 // 30 секунд
const PING_TIMEOUT = 5_000   // 5 секунд на ответ

async function ping() {
  if (inFlight) inFlight.abort()
  const controller = new AbortController()
  inFlight = controller

  const timeout = setTimeout(() => controller.abort(), PING_TIMEOUT)

  try {
    const res = await fetch(HEALTH_URL, {
      method: 'GET',
      cache: 'no-store',
      signal: controller.signal,
    })
    const reachable = res.ok
    if (reachable !== isServerReachable.value) {
      isServerReachable.value = reachable
      listeners.forEach((cb) => cb())
    }
    lastCheck.value = new Date()
  } catch {
    if (isServerReachable.value !== false) {
      isServerReachable.value = false
      listeners.forEach((cb) => cb())
    }
    lastCheck.value = new Date()
  } finally {
    clearTimeout(timeout)
    if (inFlight === controller) inFlight = null
  }
}

function start() {
  if (pingTimer) return
  ping()
  pingTimer = setInterval(ping, PING_INTERVAL)
}

function stop() {
  if (pingTimer) {
    clearInterval(pingTimer)
    pingTimer = null
  }
  if (inFlight) {
    inFlight.abort()
    inFlight = null
  }
}

let refCount = 0

export function useServerStatus() {
  onMounted(() => {
    if (refCount === 0) start()
    refCount++
  })

  onUnmounted(() => {
    refCount--
    if (refCount === 0) stop()
  })

  return {
    isServerReachable,
    lastCheck,
    recheck: ping,
  }
}