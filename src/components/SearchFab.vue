<script setup lang="ts">
const openCommandPalette = () => {
    window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', bubbles: true, cancelable: true }))
}
</script>

<template>
    <button @click="openCommandPalette"
        class="sm:hidden w-11 h-11 flex items-center justify-center rounded-full border border-border bg-background/90 backdrop-blur-sm text-muted shadow-lg hover:text-foreground transition-colors"
        aria-label="Поиск" title="Поиск" type="button">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
    </button>
</template>