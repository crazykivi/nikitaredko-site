const CELL = 56 // размер «блока» в px
const MAX_DELAY = 420 // макс. случайная задержка появления/исчезновения
const DUR = 180 // длительность анимации одного блока

const BASE = import.meta.env.BASE_URL
const BLOCKS: Array<{ tex: string; color: string }> = [
    { tex: `${BASE}dirt.svg`, color: '#26190f' },
    { tex: `${BASE}stone.svg`, color: '#4f4f4f' }
]

export function mcBlockTransition(apply: () => void): void {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
        apply()
        return
    }

    const cols = Math.max(1, Math.ceil(window.innerWidth / CELL))
    const rows = Math.max(1, Math.ceil(window.innerHeight / CELL))

    const overlay = document.createElement('div')
    overlay.className = 'mc-transition'
    overlay.style.gridTemplateColumns = `repeat(${cols}, ${CELL}px)`
    overlay.style.gridTemplateRows = `repeat(${rows}, ${CELL}px)`

    const cells: HTMLDivElement[] = []
    for (let i = 0; i < cols * rows; i++) {
        const cell = document.createElement('div')
        cell.className = 'mc-transition-cell'
        const block = BLOCKS[(Math.random() * BLOCKS.length) | 0]
        cell.style.backgroundColor = block.color
        cell.style.backgroundImage = `url("${block.tex}")`
        const inDelay = (Math.random() * MAX_DELAY) | 0
        cell.style.animationDelay = `${inDelay}ms`
        cell.dataset.inDelay = String(inDelay)
        overlay.appendChild(cell)
        cells.push(cell)
    }
    document.documentElement.appendChild(overlay)

    const swapAt = MAX_DELAY + DUR
    window.setTimeout(apply, swapAt)

    window.setTimeout(() => {
        for (const cell of cells) {
            const inDelay = Number(cell.dataset.inDelay ?? 0)
            cell.style.animationDelay = `${MAX_DELAY - inDelay}ms`
            cell.classList.add('mc-out')
        }
    }, swapAt + 80)

    window.setTimeout(() => overlay.remove(), swapAt + 80 + MAX_DELAY + DUR + 100)
}