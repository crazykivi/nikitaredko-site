import { describe, it, expect } from 'vitest'
import {
    hash2,
    materialAtCell,
    cellAtDoc,
    cellAtPointer,
    oreAtCell,
    cellClipPath,
    CELL,
    BG_SIZE,
    PX,
    JITTER_PX,
    CLIP_MARGIN,
    ORE_PX,
} from '../minecraftGrid'

const VALID_ORES = ['coal', 'iron', 'gold', 'lapis', 'diamond', 'emerald']

describe('minecraftGrid', () => {
    describe('constants', () => {
        it('CELL equals BG_SIZE', () => {
            expect(CELL).toBe(192)
            expect(BG_SIZE).toBe(CELL)
        })

        it('ORE_PX equals CELL', () => {
            expect(ORE_PX).toBe(CELL)
        })

        it('CLIP_MARGIN derives from PX and JITTER_PX', () => {
            expect(CLIP_MARGIN).toBe(JITTER_PX * PX)
        })

        it('SUB divides CELL evenly by PX', () => {
            expect(CELL % PX).toBe(0)
        })
    })

    describe('hash2', () => {
        it('is deterministic', () => {
            expect(hash2(10, 20)).toBe(hash2(10, 20))
        })

        it('returns value in [0, 1)', () => {
            for (let i = 0; i < 50; i++) {
                const v = hash2(i, i * 7)
                expect(v).toBeGreaterThanOrEqual(0)
                expect(v).toBeLessThan(1)
            }
        })

        it('produces different values for different inputs', () => {
            expect(hash2(1, 2)).not.toBe(hash2(3, 4))
            expect(hash2(0, 0)).not.toBe(hash2(0, 1))
        })
    })

    describe('materialAtCell', () => {
        it('always returns dirt at surface rows (absY < 320)', () => {

            for (let col = 0; col < 10; col++) {
                expect(materialAtCell(col, 0)).toBe('dirt')
                expect(materialAtCell(col, 1)).toBe('dirt')
            }
        })

        it('returns stone for most cells at deep levels', () => {
            let stoneCount = 0
            for (let col = 0; col < 50; col++) {
                if (materialAtCell(col, 50) === 'stone') stoneCount++
            }

            expect(stoneCount).toBeGreaterThan(40)
        })

        it('returns a valid Material type', () => {
            const mat = materialAtCell(5, 5)
            expect(['dirt', 'stone']).toContain(mat)
        })
    })

    describe('cellAtDoc', () => {
        it('calculates correct col, row, left, top', () => {
            const cell = cellAtDoc(CELL * 2 + 10, CELL * 3 + 5)
            expect(cell.col).toBe(2)
            expect(cell.row).toBe(3)
            expect(cell.left).toBe(CELL * 2)
            expect(cell.top).toBe(CELL * 3)
        })

        it('returns dirt at surface', () => {
            const cell = cellAtDoc(10, 10)
            expect(cell.material).toBe('dirt')
        })

        it('handles exact cell boundaries', () => {
            const cell = cellAtDoc(CELL, CELL)
            expect(cell.col).toBe(1)
            expect(cell.row).toBe(1)
        })

        it('handles zero coordinates', () => {
            const cell = cellAtDoc(0, 0)
            expect(cell.col).toBe(0)
            expect(cell.row).toBe(0)
        })
    })

    describe('cellAtPointer', () => {
        it('delegates to cellAtDoc with default scroll (0,0)', () => {
            const cell = cellAtPointer(CELL + 5, CELL * 2 + 10)
            expect(cell.col).toBe(1)
            expect(cell.row).toBe(2)
        })
    })

    describe('cellClipPath', () => {
        it('returns a valid polygon string', () => {
            const path = cellClipPath(0, 0)
            expect(path).toMatch(/^polygon\(/)
            expect(path).toContain('px')
        })

        it('is deterministic for the same cell', () => {
            expect(cellClipPath(3, 3)).toBe(cellClipPath(3, 3))
        })

        it('produces polygon with multiple points', () => {
            const path = cellClipPath(5, 30)
            const points = path.replace(/^polygon\(/, '').replace(/\)$/, '').split(',')
            expect(points.length).toBeGreaterThanOrEqual(4)
        })

        it('handles cells in the dirt→stone transition zone', () => {

            for (let row = 3; row <= 10; row++) {
                const path = cellClipPath(2, row)
                expect(path).toMatch(/^polygon\(/)
            }
        })

        it('handles deep stone cells', () => {
            const path = cellClipPath(10, 60)
            expect(path).toMatch(/^polygon\(/)
        })
    })

    describe('oreAtCell', () => {
        it('returns null for dirt cells', () => {
            expect(oreAtCell(0, 0, 2000)).toBeNull()
            expect(oreAtCell(0, 1, 2000)).toBeNull()
            expect(oreAtCell(0, 0, 100000)).toBeNull()
        })

        it('returns null or a valid ore for deep cells', () => {
            const ore = oreAtCell(0, 50, 20000)
            if (ore !== null) {
                expect(VALID_ORES).toContain(ore)
            }
        })

        it('finds at least one ore among many deep stone cells', () => {
            let foundOre = false

            for (let col = 0; col < 300 && !foundOre; col++) {
                if (materialAtCell(col, 50) !== 'stone') continue
                const ore = oreAtCell(col, 50, 50000)
                if (ore !== null) {
                    expect(VALID_ORES).toContain(ore)
                    foundOre = true
                }
            }
            expect(foundOre).toBe(true)
        })

        it('respects depth-based ore distribution', () => {
            for (const docH of [5000, 20000, 80000]) {
                for (let col = 0; col < 20; col++) {
                    const ore = oreAtCell(col, 40, docH)
                    if (ore !== null) {
                        expect(VALID_ORES).toContain(ore)
                    }
                }
            }
        })
    })
})