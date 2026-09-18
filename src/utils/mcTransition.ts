const CELL = 56 // размер «блока» в px
const MAX_DELAY = 420 // макс. случайная задержка появления/исчезновения
const DUR = 180 // длительность анимации одного блока

// палитра «блоков»: уголь / камень / земля
const PALETTE = ['#1c1c1c', '#242424', '#2e2e2e', '#26190f', '#171717']

export function mcBlockTransition(apply: () => void): void {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
        apply()
        return
    }

    const cols = Math.max(1, Math.ceil(window.innerWidth / CELL))
    const rows = Math.max(1, Math.ceil(window.innerHeight / CELL))

    const overlay = document.createElement('div')
    overlay.className = 'mc-transition'
    overlay.style.gridTemplateColumns = `repeat(${cols}, 1fr)`
    overlay.style.gridTemplateRows = `repeat(${rows}, 1fr)`

    const cells: HTMLDivElement[] = []
    for (let i = 0; i < cols * rows; i++) {
        const cell = document.createElement('div')
        cell.className = 'mc-transition-cell'
        cell.style.backgroundColor = PALETTE[(Math.random() * PALETTE.length) | 0]
        cell.style.animationDelay = `${(Math.random() * MAX_DELAY) | 0}ms`
        cell.style.animationDuration = `${DUR}ms`
        overlay.appendChild(cell)
        cells.push(cell)
    }
    document.documentElement.appendChild(overlay)

    const swapAt = MAX_DELAY + DUR
    window.setTimeout(apply, swapAt)

    window.setTimeout(() => {
        for (const cell of cells) {
            cell.style.animationDelay = `${(Math.random() * MAX_DELAY) | 0}ms`
            cell.classList.add('mc-out')
        }
    }, swapAt + 80)

    window.setTimeout(() => overlay.remove(), swapAt + 80 + MAX_DELAY + DUR + 100)
}