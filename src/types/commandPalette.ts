export type CommandIcon =
    | 'home' | 'article' | 'user' | 'tool'
    | 'theme' | 'github' | 'rss'

export interface Command {
    id: string
    label: string
    description?: string
    group: string
    icon: CommandIcon
    shortcut?: string[]
    keepOpen?: boolean
    action: () => void
}

export interface CommandGroup {
    name: string
    commands: Command[]
}