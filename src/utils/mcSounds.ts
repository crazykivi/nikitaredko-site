import { unlockAudio } from './digSound'
import { isSoundMuted } from '../composables/useSoundSettings'

export type McSoundName = 'dig' | 'break' | 'click' | 'xp' | 'cave' | 'piston' | 'ladder' | 'hiss' | 'explode'

const SOURCES: Record<McSoundName, string[]> = {
    dig: ['/sounds/gravel1.mp3', '/sounds/gravel2.mp3', '/sounds/gravel3.mp3', '/sounds/gravel4.mp3'],
    break: ['/sounds/break.mp3'],
    click: ['/sounds/click.mp3'],
    xp: ['/sounds/xp.mp3'],
    cave: ['/sounds/cave1.mp3', '/sounds/cave2.mp3', '/sounds/cave3.mp3', '/sounds/cave4.mp3', '/sounds/cave5.mp3', '/sounds/cave6.mp3'],
    piston: ['/sounds/piston.mp3'],
    ladder: ['/sounds/ladder1.mp3', '/sounds/ladder2.mp3', '/sounds/ladder3.mp3'],
    hiss: ['/sounds/hiss.mp3'],
    explode: ['/sounds/explode.mp3'],
}

const players = new Map<string, HTMLAudioElement>()

function getPlayer(name: McSoundName): HTMLAudioElement | null {
    const sources = SOURCES[name]
    if (!sources || sources.length === 0) return null
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