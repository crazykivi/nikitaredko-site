<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted, nextTick, ref } from "vue";
import { useCommandPalette } from "../composables/useCommandPalette";
import { useCommandPaletteCommands } from "../composables/useCommandPaletteCommands";
import { useCommandPaletteKeyboard } from "../composables/useCommandPaletteKeyboard";
import type { Command } from "../types/commandPalette";
import CommandPaletteList from "./CommandPaletteList.vue";

const {
  isOpen,
  query,
  selectedIndex,
  expandedId,
  open,
  close,
  onCmdPointerDown,
  onCmdPointerEnd,
} = useCommandPalette();

const inputRef = ref<HTMLInputElement | null>(null);
const queryComputed = computed(() => query.value);
const { filtered, grouped, loadArticles, abort } = useCommandPaletteCommands(
  queryComputed
);

const execute = (cmd: Command) => {
  if (!cmd.keepOpen) close();
  try {
    cmd.action();
  } catch (e) {
    console.error("[CommandPalette] action failed:", e);
  }
};

const onOpen = () => {
  open();
  loadArticles();
};

const onClose = () => {
  close();
  abort();
};

const { handleInputKeydown, handleGlobalKeydown } = useCommandPaletteKeyboard({
  isOpen,
  selectedIndex,
  flatCommands: filtered,
  onExecute: execute,
  onOpen: onOpen,
  onClose: onClose,
});

watch(isOpen, (opened) => {
  if (opened) {
    nextTick(() => {
      inputRef.value?.focus({ preventScroll: true });
    });
  }
});

onMounted(() => window.addEventListener("keydown", handleGlobalKeydown));
onUnmounted(() => window.removeEventListener("keydown", handleGlobalKeydown));
</script>

<template>
  <Transition name="palette">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[9998] flex items-start justify-center pt-[15vh] px-4"
      @click.self="onClose"
    >
      <div class="absolute inset-0 modal-overlay bg-background/70 backdrop-blur-sm" />
      <div
        class="modal-panel relative w-full max-w-xl rounded-2xl border border-border bg-background shadow-2xl overflow-hidden animate-slide-down modal-control-panel"
        role="dialog"
        aria-modal="true"
        aria-label="Command palette"
      >
        <!-- INPUT -->
        <div class="flex items-center gap-3 px-5 py-4 border-b border-border">
          <svg
            class="w-5 h-5 text-muted shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            :value="query"
            @input="query = ($event.target as HTMLInputElement).value"
            type="text"
            placeholder="Поиск страниц, статей, действий..."
            class="flex-1 bg-transparent outline-none text-foreground placeholder:text-muted/60 text-base"
            autocomplete="off"
            spellcheck="false"
            @keydown="handleInputKeydown"
            autofocus
          />
          <button
            type="button"
            @click="onClose"
            class="esc-badge hidden sm:inline-flex items-center px-1.5 py-0.5 text-[10px] font-mono rounded border border-border text-muted hover:text-foreground"
            aria-label="Закрыть (Esc)"
            title="Закрыть (Esc)"
          >
            ESC
          </button>
        </div>

        <!-- LIST -->
        <CommandPaletteList
          :groups="grouped"
          :selected-index="selectedIndex"
          :expanded-id="expandedId"
          @execute="execute"
          @select="selectedIndex = $event"
          @pointer-down="onCmdPointerDown"
          @pointer-end="onCmdPointerEnd"
        />

        <!-- FOOTER -->
        <div
          class="hidden sm:flex items-center justify-between px-5 py-2 border-t border-border text-[10px] font-mono text-muted/70 bg-muted/10"
        >
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">↑↓</kbd>
              навигация</span
            >
            <span class="hidden sm:flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">Tab</kbd> цикл</span
            >
            <span class="flex items-center gap-1"
              ><kbd class="px-1 py-0.5 rounded border border-border">↵</kbd> выбрать</span
            >
          </div>
          <span class="flex items-center gap-2">
            <span
              >открыть:
              <kbd class="px-1 py-0.5 rounded border border-border">/</kbd></span
            >
            <span class="opacity-50">или</span>
            <kbd class="px-1 py-0.5 rounded border border-border">⌘⇧K</kbd>
          </span>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@keyframes slide-down {
  from {
    opacity: 0;
    transform: translateY(-12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
.animate-slide-down {
  animation: slide-down 0.18s ease-out;
}
.palette-enter-active,
.palette-leave-active {
  transition: opacity 0.15s ease;
}
.palette-enter-from,
.palette-leave-to {
  opacity: 0;
}
</style>
