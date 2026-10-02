<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useTheme } from '../composables/useTheme'
import {
    CELL, BG_SIZE, CLIP_MARGIN, ORE_PX,
    cellClipPath, materialAtCell, oreAtCell,
    type Material, type Ore,
} from '../utils/minecraftGrid'

const ORE_TEX: Record<Ore, string> = {
    coal: `${import.meta.env.BASE_URL}stone_coal.svg`,
    iron: `${import.meta.env.BASE_URL}stone_iron.svg`,
    gold: `${import.meta.env.BASE_URL}stone_gold.svg`,
    lapis: `${import.meta.env.BASE_URL}stone_lapis.svg`,
    diamond: `${import.meta.env.BASE_URL}stone_diamond.svg`,
    emerald: `${import.meta.env.BASE_URL}stone_emerald.svg`,
}

const BUFFER_ROWS = 2
const { mode } = useTheme()
const stones = ref<Stone[]>([])

interface Stone {
    key: string
    left: number
    top: number
    px: number
    py: number
    clip: string
    material: Material
    ore: string | null
}

let rafId = 0
let resizeTimer: ReturnType<typeof setTimeout> | null = null
let resizeObserver: ResizeObserver | null = null

let lastStartRow = -1
let lastEndRow = -1
let lastCols = -1

const recomputeVisible = () => {
    const doc = document.documentElement
    const width = window.innerWidth
    const height = window.innerHeight
    const scrollY = window.scrollY
    const scrollHeight = doc.scrollHeight

    const nextCols = Math.max(1, Math.ceil(width / CELL))
    const maxRow = Math.floor(scrollHeight / CELL)

    const startRow = Math.max(0, Math.floor((scrollY - BUFFER_ROWS * CELL) / CELL))
    const endRow = Math.min(maxRow, Math.ceil((scrollY + height + BUFFER_ROWS * CELL) / CELL))

    if (startRow === lastStartRow && endRow === lastEndRow && nextCols === lastCols) {
        rafId = 0
        return
    }

    lastStartRow = startRow
    lastEndRow = endRow
    lastCols = nextCols

    const next: Stone[] = []

    for (let absRow = startRow; absRow <= endRow; absRow++) {
        const absY = absRow * CELL
        for (let x = 0; x < nextCols; x++) {
            const mat = materialAtCell(x, absRow)
            if (mat === 'stone') {
                const M = CLIP_MARGIN
                const ox = x * CELL - M
                const oy = absY - M
                const ore = oreAtCell(x, absRow, scrollHeight)
                next.push({
                    key: `${absRow}:${x}`,
                    left: ox,
                    top: oy,
                    px: -(((ox % BG_SIZE) + BG_SIZE) % BG_SIZE),
                    py: -(((oy % BG_SIZE) + BG_SIZE) % BG_SIZE),
                    clip: cellClipPath(x, absRow),
                    material: mat,
                    ore: ore ? ORE_TEX[ore] : null,
                })
            }
        }
    }

    stones.value = next
    rafId = 0
}

const schedule = () => {
    if (!rafId) rafId = requestAnimationFrame(recomputeVisible)
}

const onWindowResize = () => {
    if (resizeTimer) clearTimeout(resizeTimer)
    resizeTimer = setTimeout(() => {
        lastStartRow = -1
        schedule()
    }, 120)
}

onMounted(() => {
    recomputeVisible()

    window.addEventListener('scroll', schedule, { passive: true })
    window.addEventListener('resize', onWindowResize)

    if (typeof ResizeObserver !== 'undefined') {
        resizeObserver = new ResizeObserver(schedule)
        resizeObserver.observe(document.documentElement)
    }
})

onUnmounted(() => {
    window.removeEventListener('scroll', schedule)
    window.removeEventListener('resize', onWindowResize)
    if (resizeTimer) clearTimeout(resizeTimer)
    if (resizeObserver) resizeObserver.disconnect()
    if (rafId) cancelAnimationFrame(rafId)
})
</script>

<template>
    <div v-if="mode === 'charcoal'" class="mc-depth-layer" :style="{ '--bg': `${BG_SIZE}px` }" aria-hidden="true">
        <div v-for="s in stones" :key="s.key" :class="['mc-depth-stone', `mc-depth-${s.material}`]" :style="{
            left: `${s.left}px`,
            top: `${s.top}px`,
            width: `${CELL + CLIP_MARGIN * 2}px`,
            height: `${CELL + CLIP_MARGIN * 2}px`,
            backgroundPosition: `${s.px}px ${s.py}px`,
            clipPath: s.clip,
        }">
            <div v-if="s.ore" class="mc-depth-ore" :style="{
                left: `${CLIP_MARGIN}px`,
                top: `${CLIP_MARGIN}px`,
                width: `${ORE_PX}px`,
                height: `${ORE_PX}px`,
                backgroundImage: `url(${s.ore})`,
            }"></div>
        </div>
    </div>
</template>

<style scoped>
.mc-depth-ore {
    position: absolute;
    background-size: 100% 100%;
    background-repeat: no-repeat;
    image-rendering: pixelated;
    pointer-events: none;
}
</style>