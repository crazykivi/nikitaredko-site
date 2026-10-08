<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
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
  shouldSuppressClick,
} = useCommandPalette();

const panelRef = ref<HTMLElement | null>(null);
const inputRef = ref<HTMLInputElement | null>(null);

const queryComputed = computed(() => query.value);

const {
  filtered,
  grouped,
  loadArticles,
  abort,
  quickNavigate,
} = useCommandPaletteCommands(queryComputed);

const focusInputImmediate = () => {
  inputRef.value?.focus({ preventScroll: true });
};

let focusFrame = 0;

const focusInput = () => {
  cancelAnimationFrame(focusFrame);
  focusFrame = requestAnimationFrame(() => {
    if (!isOpen.value) return;
    focusInputImmediate();
  });
};

const execute = (cmd: Command) => {
  const shouldClose = !cmd.keepOpen;

  if (shouldClose) {
    close();
  }

  try {
    cmd.action();
  } catch (e) {
    console.error("[CommandPalette] action failed:", e);
  }

  if (!shouldClose) {
    nextTick(focusInputImmediate);
  }
};

const executeFromPointer = (cmd: Command) => {
  if (shouldSuppressClick()) return;
  execute(cmd);
};

const onOpen = () => {
  open();
  loadArticles();
};

const onClose = () => {
  close();
  abort();
};

const { handleGlobalKeydown } = useCommandPaletteKeyboard({
  isOpen,
  selectedIndex,
  flatCommands: filtered,
  onExecute: execute,
  onOpen,
  onClose,
  onQuickNavigate: quickNavigate,
});

watch(filtered, () => {
  selectedIndex.value = 0;
});

const keepFocusInsidePalette = (e: FocusEvent) => {
  if (!isOpen.value || !panelRef.value) return;

  const target = e.target as Node | null;
  if (!target) return;

  if (panelRef.value.contains(target)) return;

  focusInput();
};

const lockBodyScroll = () => {
  document.documentElement.style.overflow = "hidden";
  document.body.style.overflow = "hidden";
};

const unlockBodyScroll = () => {
  document.documentElement.style.overflow = "";
  document.body.style.overflow = "";
};

watch(isOpen, (opened) => {
  if (opened) {
    lockBodyScroll();
    nextTick(focusInput);
  } else {
    unlockBodyScroll();
  }
});

const onQueryInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  query.value = target.value;
};

onMounted(() => {
  window.addEventListener("keydown", handleGlobalKeydown, true);
  document.addEventListener("focusin", keepFocusInsidePalette);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleGlobalKeydown, true);
  document.removeEventListener("focusin", keepFocusInsidePalette);

  cancelAnimationFrame(focusFrame);
  unlockBodyScroll();
  abort();
});
</script>

<template>
  <Transition name="palette">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[9998] flex items-start justify-center pt-[15vh] px-4"
    >
      <div
        class="absolute inset-0 modal-overlay bg-background/70 backdrop-blur-sm touch-none"
        aria-hidden="true"
        @pointerdown.prevent="onClose"
        @wheel.prevent
        @touchmove.prevent
      />
      <div
        ref="panelRef"
        class="modal-panel relative z-10 w-full max-w-xl rounded-2xl border border-border bg-background shadow-2xl overflow-hidden animate-slide-down modal-control-panel"
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
            ref="inputRef"
            :value="query"
            type="text"
            placeholder="Поиск страниц, статей, действий..."
            class="flex-1 bg-transparent outline-none text-foreground placeholder:text-muted/60 text-base"
            autocomplete="off"
            spellcheck="false"
            @input="onQueryInput"
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
          @execute="executeFromPointer"
          @select="selectedIndex = $event"
          @pointer-down="onCmdPointerDown"
          @pointer-end="onCmdPointerEnd"
        />

        <!-- FOOTER -->
        <div
          class="hidden sm:flex items-center justify-between px-5 py-2 border-t border-border text-[10px] font-mono text-muted/70 bg-muted/10"
        >
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1">
              <kbd class="px-1 py-0.5 rounded border border-border">↑↓</kbd>
              навигация
            </span>
            <span class="hidden sm:flex items-center gap-1">
              <kbd class="px-1 py-0.5 rounded border border-border">Tab</kbd>
              цикл
            </span>
            <span class="flex items-center gap-1">
              <kbd class="px-1 py-0.5 rounded border border-border">↵</kbd>
              выбрать
            </span>
          </div>

          <span class="flex items-center gap-2">
            <span>
              открыть:
              <kbd class="px-1 py-0.5 rounded border border-border">/</kbd>
            </span>
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
