import { unlockAudio } from './digSound'
import { isSoundMuted } from '../composables/useSoundSettings'

export type McSoundName = 'dig' | 'break' | 'click' | 'xp' | 'cave' | 'piston_in' | 'piston_out' |  'ladder' | 'hiss' | 'explode'

const SOURCES: Record<McSoundName, string[]> = {
    dig: ['/sounds/gravel1.mp3', '/sounds/gravel2.mp3', '/sounds/gravel3.mp3', '/sounds/gravel4.mp3'],
    break: ['/sounds/break.mp3'],
    click: ['/sounds/click.mp3'],
    xp: ['/sounds/xp.mp3'],
    cave: ['/sounds/cave1.mp3', '/sounds/cave2.mp3', '/sounds/cave3.mp3', '/sounds/cave4.mp3', '/sounds/cave5.mp3', '/sounds/cave6.mp3'],
    piston_in: ['/sounds/piston_in.mp3'],
    piston_out: ['/sounds/piston_out.mp3'],
    ladder: ['/sounds/ladder1.mp3', '/sounds/ladder2.mp3', '/sounds/ladder3.mp3'],
    hiss: ['/sounds/hiss.mp3'],
    explode: ['/sounds/explode.mp3'],
}

const players = new Map<string, HTMLAudioElement>()
const POOL_SOUNDS: Set<McSoundName> = new Set(['ladder'])
const POOL_SIZE = 4
const pools = new Map<string, HTMLAudioElement[]>()
const poolCursor = new Map<string, number>()

function getPlayer(name: McSoundName): HTMLAudioElement | null {
    const sources = SOURCES[name]
    if (!sources || sources.length === 0) return null
    if (POOL_SOUNDS.has(name)) {
        const key = `${name}-pool`
        let pool = pools.get(key)
        if (!pool) {
            pool = []
            for (let i = 0; i < POOL_SIZE; i++) {
                const src = sources[i % sources.length] 
                const el = new Audio(src)
                el.preload = 'auto'
                pool.push(el)
            }
            pools.set(key, pool)
            poolCursor.set(key, -1)
        }
        let idx: number
        const prev = poolCursor.get(key)!
        if (pool.length <= 1) {
            idx = 0
        } else {
            do {
                idx = Math.floor(Math.random() * pool.length)
            } while (idx === prev)
        }
        poolCursor.set(key, idx)
        return pool[idx]
    }
    const randomSource = sources[Math.floor(Math.random() * sources.length)]
    const key = `${name}-${randomSource}`
    let el = players.get(key)
    if (!el) {
        el = new Audio(randomSource)
        el.preload = 'auto'
        players.set(key, el)
    }
    return el
}

export function primeMcSounds(): void {
    ; (Object.keys(SOURCES) as McSoundName[]).forEach((name) => getPlayer(name))
    try {
        unlockAudio()
    } catch {
        /* AudioContext не поддерживается — игнорирование */
    }
}

export function playMcSound(
    name: McSoundName,
    opts: { volume?: number; rate?: number; loop?: boolean } = {},
): void {
    if (isSoundMuted.value) return
    unlockAudio()
    const el = getPlayer(name)
    if (!el) return
    el.currentTime = 0
    el.volume = opts.volume ?? 1
    el.playbackRate = opts.rate ?? 1
    el.loop = opts.loop ?? false
    void el.play().catch(() => { })
}

export function stopMcSound(name: McSoundName): void {
    const sources = SOURCES[name]
    if (!sources) return
    sources.forEach((source) => {
        const el = players.get(`${name}-${source}`)
        if (!el) return
        el.pause()
        el.currentTime = 0
        el.loop = false
    })
}