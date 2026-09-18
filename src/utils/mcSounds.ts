import { unlockAudio } from './digSound'

export type McSoundName = 'dig' | 'break' | 'click' | 'xp' | 'cave' | 'piston'

const SOURCES: Record<McSoundName, string[]> = {
    dig: ['/sounds/gravel1.mp3', '/sounds/gravel2.mp3', '/sounds/gravel3.mp3', '/sounds/gravel4.mp3'],
    break: ['/sounds/break.mp3'],
    click: ['/sounds/click.mp3'],
    xp: ['/sounds/xp.mp3'],
    cave: ['/sounds/cave1.mp3','/sounds/cave2.mp3',
        '/sounds/cave3.mp3',
        '/sounds/cave4.mp3',
        '/sounds/cave5.mp3',
        '/sounds/cave6.mp3',
    ],
    piston: ['/sounds/piston.mp3'],
}

const players = new Map<string, HTMLAudioElement>()

function getPlayer(name: McSoundName): HTMLAudioElement | null {
    const sources = SOURCES[name]
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
    opts: { volume?: number; rate?: number } = {},
): void {
    unlockAudio()
    const el = getPlayer(name)
    if (!el) return

    el.currentTime = 0
    el.volume = opts.volume ?? 1
    el.playbackRate = opts.rate ?? 1
    void el.play().catch(() => { })
}