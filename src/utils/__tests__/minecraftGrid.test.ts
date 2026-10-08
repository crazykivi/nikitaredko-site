import { describe, it, expect } from 'vitest'
import { hash2, materialAtCell, cellAtDoc, oreAtCell, CELL } from '../minecraftGrid'

describe('minecraftGrid', () => {
    it('hash2 returns deterministic value between 0 and 1', () => {
        const val1 = hash2(10, 20)
        const val2 = hash2(10, 20)
        expect(val1).toBe(val2)
        expect(val1).toBeGreaterThanOrEqual(0)
        expect(val1).toBeLessThan(1)
    })

    it('materialAtCell returns dirt or stone', () => {
        const mat = materialAtCell(5, 5)
        expect(['dirt', 'stone']).toContain(mat)
    })

    it('cellAtDoc calculates correct grid coordinates', () => {
        const cell = cellAtDoc(CELL * 2 + 10, CELL * 3 + 5)
        expect(cell.col).toBe(2)
        expect(cell.row).toBe(3)
    })

    it('oreAtCell returns null or valid ore', () => {
        const ore = oreAtCell(0, 0, 2000)
        if (ore !== null) {
            expect(['coal', 'iron', 'gold', 'lapis', 'diamond', 'emerald']).toContain(ore)
        }
    })
})