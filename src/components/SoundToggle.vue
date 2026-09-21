<script setup lang="ts">
import { useSoundSettings } from '../composables/useSoundSettings'
import { useTheme } from '../composables/useTheme'

const { isSoundMuted, toggleMute } = useSoundSettings()
const { mode } = useTheme()
</script>

<template>
    <Transition name="slide-up">
        <button v-if="mode === 'charcoal'" @click="toggleMute"
            class="w-11 h-11 flex items-center justify-center rounded-full bg-foreground text-background shadow-lg hover:opacity-85 hover:scale-110 transition-all duration-200"
            :aria-label="isSoundMuted ? 'Включить звуки' : 'Выключить звуки'"
            :title="isSoundMuted ? 'Включить звуки' : 'Выключить звуки'" type="button">
            <svg v-if="!isSoundMuted" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
                aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2" />
            </svg>
        </button>
    </Transition>
</template>