export const CELL = 192
export const BG_SIZE = 192
const SURFACE_PX = 800
const TRANSITION_PX = 800
const STONE_MAX = 0.985

const DEEPSLATE_START_PX = 6000
const DEEPSLATE_TRANSITION_PX = 800
const DEEPSLATE_MAX = 0.985

const NOISE = 0.14
const clamp = (v: number, min = 0, max = 1) => Math.min(max, Math.max(min, v))
const smoothstep = (t: number) => t * t * (3 - 2 * t)

export const hash2 = (x: number, y: number) => {
    let h = Math.imul(x ^ 0x9e3779b9, 0x85ebca6b)
    h ^= Math.imul(y ^ 0xc2b2ae35, 0x27d4eb2f)
    h ^= h >>> 15
    h = Math.imul(h, 0x2c1b3c6d)
    h ^= h >>> 12
    h = Math.imul(h, 0x297a2d39)
    h ^= h >>> 15
    return (h >>> 0) / 4294967296
}

const stoneChance = (x: number, absY: number) => {
    if (absY < SURFACE_PX * 0.4) return 0
    const t = clamp((absY - SURFACE_PX) / TRANSITION_PX)
    const base = smoothstep(t) * STONE_MAX
    const jitter = (hash2(x * 7 + 3, absY * 11 + 5) - 0.5) * NOISE
    return clamp(base + jitter)
}

const deepslateChance = (x: number, absY: number) => {
    if (absY < DEEPSLATE_START_PX) return 0
    const t = clamp((absY - DEEPSLATE_START_PX) / DEEPSLATE_TRANSITION_PX)
    const base = smoothstep(t) * DEEPSLATE_MAX
    const jitter = (hash2(x * 13 + 7, absY * 17 + 11) - 0.5) * NOISE
    return clamp(base + jitter)
}

export type Material = 'dirt' | 'stone' | 'deepslate'

export const materialAtCell = (col: number, row: number): Material => {
    const absY = row * CELL
    const r = hash2(col, row)
    
    const dChance = deepslateChance(col, absY)
    if (r < dChance) return 'deepslate'

    const sChance = stoneChance(col, absY)
    if (r < sChance) return 'stone'

    return 'dirt'
}

export interface GridCell {
    col: number
    row: number
    left: number
    top: number
    material: Material
}

export const cellAtDoc = (docX: number, docY: number): GridCell => {
    const col = Math.floor(docX / CELL)
    const row = Math.floor(docY / CELL)
    return { col, row, left: col * CELL, top: row * CELL, material: materialAtCell(col, row) }
}

export const cellAtPointer = (clientX: number, clientY: number): GridCell =>
    cellAtDoc(clientX + window.scrollX, clientY + window.scrollY)

export const PX = 16
export const JITTER_PX = 2
export const CLIP_MARGIN = JITTER_PX * PX
const SUB = CELL / PX

const profileCache = new Map<string, number[]>()

const NOISE_STEP = 3

function borderProfile(seedA: number, seedB: number): number[] {
    const key = `${seedA}:${seedB}`
    const hit = profileCache.get(key)
    if (hit) return hit

    const node = (n: number) => hash2(seedA * 13 + n * 5, seedB * 7 + n * 11) * 2 - 1

    const out: number[] = new Array(SUB)
    for (let i = 0; i < SUB; i++) {
        const x = i / NOISE_STEP
        const i0 = Math.floor(x)
        const t = smoothstep(x - i0)
        const n = node(i0) * (1 - t) + node(i0 + 1) * t
        out[i] = Math.max(-JITTER_PX, Math.min(JITTER_PX, Math.round(n * (JITTER_PX + 0.45))))
    }

    profileCache.set(key, out)
    return out
}

const hBorder = (col: number, rb: number) => borderProfile(col * 2 + 1, rb * 4 + 3)
const vBorder = (cb: number, row: number) => borderProfile(cb * 4 + 2, row * 2 + 5)

export function cellClipPath(col: number, row: number): string {
    const M = CLIP_MARGIN
    const own = materialAtCell(col, row)
    const differs = (c: number, r: number) => materialAtCell(c, r) !== own

    const t = differs(col, row - 1) ? hBorder(col, row) : null
    const b = differs(col, row + 1) ? hBorder(col, row + 1) : null
    const l = differs(col - 1, row) ? vBorder(col, row) : null
    const r = differs(col + 1, row) ? vBorder(col + 1, row) : null

    const yT: number[] = [], yB: number[] = [], xL: number[] = [], xR: number[] = []
    for (let i = 0; i < SUB; i++) {
        yT[i] = M - (t ? t[i] : 0) * PX
        yB[i] = M + CELL - (b ? b[i] : 0) * PX
        xL[i] = M - (l ? l[i] : 0) * PX
        xR[i] = M + CELL - (r ? r[i] : 0) * PX
    }

    const pts: string[] = []
    const p = (x: number, y: number) => pts.push(`${x}px ${y}px`)

    p(xL[0], yT[0])
    for (let i = 0; i < SUB; i++) {                 // верх →
        p(M + (i + 1) * PX, yT[i])
        if (i < SUB - 1) p(M + (i + 1) * PX, yT[i + 1])
    }
    p(xR[0], yT[SUB - 1])
    for (let i = 0; i < SUB; i++) {                 // право ↓
        p(xR[i], M + (i + 1) * PX)
        if (i < SUB - 1) p(xR[i + 1], M + (i + 1) * PX)
    }
    p(xR[SUB - 1], yB[SUB - 1])
    for (let i = SUB - 1; i >= 0; i--) {            // низ ←
        p(M + i * PX, yB[i])
        if (i > 0) p(M + i * PX, yB[i - 1])
    }
    p(xL[SUB - 1], yB[0])
    for (let i = SUB - 1; i >= 0; i--) {            // лево ↑
        p(xL[i], M + i * PX)
        if (i > 0) p(xL[i - 1], M + i * PX)
    }

    return `polygon(${pts.join(',')})`
}