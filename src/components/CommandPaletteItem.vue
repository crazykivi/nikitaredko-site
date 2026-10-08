<script setup lang="ts">
import type { Command } from "../types/commandPalette";
import { ICON_PATHS, FILL_ICONS } from "../data/commandPaletteIcons";

defineProps<{
  command: Command;
  selected: boolean;
  expanded: boolean;
}>();

const emit = defineEmits<{
  execute: [];
  "pointer-down": [];
  "pointer-end": [];
}>();
</script>

<template>
  <button
    :data-cmd-id="command.id"
    :data-selected="selected"
    :title="command.label"
    class="relative w-full flex items-center gap-3 px-5 py-2.5 text-left transition-colors duration-100 select-none palette-cmd"
    :class="
      selected
        ? 'bg-foreground/[0.08] text-foreground'
        : 'text-foreground/80 hover:bg-foreground/[0.04]'
    "
    @click="emit('execute')"
    @pointerdown="emit('pointer-down')"
    @pointerup="emit('pointer-end')"
    @pointercancel="emit('pointer-end')"
    @pointerleave="emit('pointer-end')"
    @contextmenu.prevent
  >
    <span
      v-if="selected"
      class="absolute left-0 top-1 bottom-1 w-0.5 bg-foreground rounded-r-full"
      aria-hidden="true"
    />
    <div class="w-7 h-7 rounded-md bg-muted/30 flex items-center justify-center shrink-0">
      <svg
        class="w-4 h-4 text-foreground/80"
        :fill="FILL_ICONS.has(command.icon) ? 'currentColor' : 'none'"
        :stroke="FILL_ICONS.has(command.icon) ? 'none' : 'currentColor'"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <path
          v-if="!FILL_ICONS.has(command.icon)"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          :d="ICON_PATHS[command.icon]"
        />
        <path v-else :d="ICON_PATHS[command.icon]" />
      </svg>
    </div>
    <div class="flex-1 min-w-0">
      <div
        class="text-sm font-medium"
        :class="expanded ? 'whitespace-normal break-words' : 'truncate'"
      >
        {{ command.label }}
      </div>
      <div
        v-if="command.description"
        class="text-xs text-muted"
        :class="expanded ? 'whitespace-normal break-words' : 'truncate'"
      >
        {{ command.description }}
      </div>
    </div>
    <div
      v-if="command.shortcut"
      class="hidden sm:flex items-center gap-1 text-[10px] font-mono text-muted/80"
    >
      <kbd
        v-for="(k, i) in command.shortcut"
        :key="i"
        class="px-1.5 py-0.5 rounded border border-border"
        >{{ k }}</kbd
      >
    </div>
    <svg
      v-else
      class="w-3.5 h-3.5 text-muted/40 shrink-0"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M13 7l5 5m0 0l-5 5m5-5H6"
      />
    </svg>
  </button>
</template>
