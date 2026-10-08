import { computed, ref, type ComputedRef } from 'vue'
import { useRouter } from 'vue-router'
import { getArticlesFeed, type Article } from '../services/api'
import { useBlockBreak } from './useBlockBreak'
import { playMcSound } from '../utils/mcSounds'
import { unlockAudio } from '../utils/digSound'
import type { Command, CommandGroup } from '../types/commandPalette'

const SECRET_WORDS = ['minecraft', 'майнкрафт']
const ARTICLES_LIMIT = 20

export function useCommandPaletteCommands(query: ComputedRef<string>) {
    const router = useRouter()
    const { unlocked, toggle } = useBlockBreak()

    const recentArticles = ref<Article[]>([])
    const articlesLoaded = ref(false)
    const articlesLoading = ref(false)
    let abortController: AbortController | null = null

    const staticCommands = computed<Command[]>(() => [
        { id: 'nav-home', label: 'Главная', group: 'Навигация', icon: 'home', shortcut: ['⌘', '1'], action: () => router.push('/') },
        { id: 'nav-articles', label: 'Статьи', group: 'Навигация', icon: 'article', shortcut: ['⌘', '2'], action: () => router.push('/articles') },
        { id: 'nav-about', label: 'О себе', group: 'Навигация', icon: 'user', shortcut: ['⌘', '3'], action: () => router.push('/about') },
        { id: 'nav-uses', label: 'Uses', group: 'Навигация', icon: 'tool', shortcut: ['⌘', '4'], action: () => router.push('/uses') },
        { id: 'action-theme', label: 'Сменить тему', description: 'Переключить светлую/тёмную тему/системная', group: 'Действия', icon: 'theme', keepOpen: true, action: () => window.dispatchEvent(new CustomEvent('toggle-app-theme')) },
        { id: 'action-github', label: 'GitHub репозиторий', description: 'crazykivi/nikitaredko-site', group: 'Действия', icon: 'github', action: () => window.open('https://github.com/crazykivi/nikitaredko-site', '_blank', 'noopener,noreferrer') },
        { id: 'action-rss', label: 'RSS-лента', description: 'Подписаться на обновления', group: 'Действия', icon: 'rss', action: () => window.open('/api/rss.xml', '_blank', 'noopener,noreferrer') },
    ])

    const secretCommand = computed<Command>(() => ({
        id: 'secret-block-break',
        label: unlocked.value ? 'Выключить режим добычи блока' : 'Включить режим добычи блока',
        description: unlocked.value ? 'Секретный режим Minecraft активен' : 'Секретная вкладка: добыча фона как в Minecraft',
        group: 'Секреты',
        icon: 'tool',
        keepOpen: true,
        action: () => { if (toggle()) { unlockAudio(); playMcSound('break') } },
    }))

    const allCommands = computed<Command[]>(() => {
        const articleCmds: Command[] = recentArticles.value.map((a: Article) => ({
            id: `article-${a.id}`,
            label: a.title,
            description: a.collectionName,
            group: 'Статьи',
            icon: 'article',
            action: () => router.push(`/articles/${a.id}`),
        }))
        return [...staticCommands.value, ...articleCmds]
    })

    const filtered = computed<Command[]>(() => {
        const q = query.value.trim().toLowerCase()
        const base = q
            ? allCommands.value.filter(c =>
                c.label.toLowerCase().includes(q) ||
                (c.description ?? '').toLowerCase().includes(q))
            : allCommands.value

        if (SECRET_WORDS.some(w => q.includes(w))) {
            return [...base, secretCommand.value]
        }
        return base
    })

    const grouped = computed<CommandGroup[]>(() => {
        const map = new Map<string, Command[]>()
        for (const cmd of filtered.value) {
            if (!map.has(cmd.group)) map.set(cmd.group, [])
            map.get(cmd.group)!.push(cmd)
        }
        return Array.from(map.entries()).map(([name, commands]) => ({ name, commands }))
    })

    const loadArticles = async () => {
        if (articlesLoaded.value || articlesLoading.value) return
        articlesLoading.value = true
        abortController = new AbortController()
        try {
            const feed = await getArticlesFeed(1, ARTICLES_LIMIT, undefined, abortController.signal)
            recentArticles.value = feed.articles
            articlesLoaded.value = true
        } catch (e) {
            if (e instanceof Error && e.name === 'AbortError') return
            console.error('[CommandPalette] Failed to load articles:', e)
        } finally {
            articlesLoading.value = false
            abortController = null
        }
    }

    const abort = () => abortController?.abort()

    return { filtered, grouped, loadArticles, abort }
}