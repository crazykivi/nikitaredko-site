<script setup lang="ts">
import { usePwaUpdate } from '../composables/usePwaUpdate'

const { needRefresh, updating, applyUpdate } = usePwaUpdate()
</script>

<template>
  <Transition name="update-toast">
    <div
      v-if="needRefresh"
      class="fixed bottom-24 right-4 z-[100] w-80 max-w-[calc(100vw-2rem)]"
      role="alert"
      aria-live="polite"
    >
      <div class="rounded-xl border border-border bg-background/95 backdrop-blur-sm shadow-lg p-4">
        <div class="flex items-start gap-3">
          <button
            class="update-icon-btn flex w-10 h-14 shrink-0 items-center justify-center rounded-full bg-foreground/10 hover:bg-foreground/15 active:scale-95 cursor-pointer transition-colors"
            :class="{ updating }"
            :disabled="updating"
            aria-label="Обновить страницу"
            title="Обновить страницу"
            @click="applyUpdate"
          >
            <svg class="update-icon w-[18px] h-[18px] text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
          </button>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-foreground">Сайт обновлён</p>
            <p class="text-xs text-muted mt-1">
              Версия устарела — перезагрузите страницу, чтобы увидеть изменения
            </p>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
