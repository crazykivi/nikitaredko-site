const STORAGE_KEY = 'mc-articles-read'
const LEGACY_PREFIX = 'article-xp-'
const MAX_ENTRIES = 5

function readIds(): string[] {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        if (!raw) return []
        const parsed: unknown = JSON.parse(raw)
        return Array.isArray(parsed)
            ? parsed.filter((x): x is string => typeof x === 'string')
            : []
    } catch {
        return []
    }
}

function saveIds(ids: string[]): void {
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(ids.slice(-MAX_ENTRIES)))
    } catch {
        // пропуск
    }
}

let migrated = false

function migrateLegacyKeys(): void {
    if (migrated) return
    migrated = true
    try {
        const legacy: string[] = []
        for (let i = localStorage.length - 1; i >= 0; i--) {
            const key = localStorage.key(i)
            if (key && key.startsWith(LEGACY_PREFIX)) legacy.push(key)
        }
        if (legacy.length === 0) return

        const ids = new Set(readIds())
        for (const key of legacy) {
            ids.add(key.slice(LEGACY_PREFIX.length))
            localStorage.removeItem(key)
        }
        saveIds([...ids])
    } catch {
        // пропуск
    }
}

export function isArticleRead(id: string): boolean {
    migrateLegacyKeys()
    return readIds().includes(id)
}

export function markArticleRead(id: string): boolean {
    migrateLegacyKeys()
    const ids = readIds()
    if (ids.includes(id)) return false
    saveIds([...ids, id])
    return true
}