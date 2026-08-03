<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from 'vue'
import type { TOCItem } from './ArticleTOC.vue'
import ArticleTOC from './ArticleTOC.vue'
import { scrollToHeading } from '../utils/scroll'

const props = defineProps<{
  items: TOCItem[]
  activeId: string
}>()

const isPopupOpen = ref(false)
const closePopup = () => { isPopupOpen.value = false }

const isHovered = ref(false)
const isPinned = ref(false)
const isExpanded = computed(() => isHovered.value || isPinned.value)

let hoverTimer: ReturnType<typeof setTimeout> | null = null
const onEnter = () => {
  if (hoverTimer) clearTimeout(hoverTimer)
  isHovered.value = true
}
const onLeave = () => {
  hoverTimer = setTimeout(() => { isHovered.value = false }, 200)
}
const togglePin = () => { isPinned.value = !isPinned.value }

onUnmounted(() => { if (hoverTimer) clearTimeout(hoverTimer) })

const activeIndex = computed(() => {
  const idx = props.items.findIndex(i => i.id === props.activeId)
  return idx === -1 ? 0 : idx
})
const prevItem = computed<TOCItem | null>(() => props.items[activeIndex.value - 1] ?? null)
const currentItem = computed<TOCItem | null>(() => props.items[activeIndex.value] ?? null)
const nextItem = computed<TOCItem | null>(() => props.items[activeIndex.value + 1] ?? null)

const listRef = ref<HTMLElement | null>(null)
watch([activeIndex, isExpanded], () => {
  if (!isExpanded.value) return
  nextTick(() => {
    const container = listRef.value
    const el = container?.querySelector('[data-active="true"]') as HTMLElement | null
    if (!container || !el) return
    const elTop = el.offsetTop
    const elBottom = elTop + el.offsetHeight
    if (elTop < container.scrollTop) container.scrollTop = elTop
    else if (elBottom > container.scrollTop + container.clientHeight) {
      container.scrollTop = elBottom - container.clientHeight
    }
  })
})
</script>

<template>
  <div
    class="hidden 2xl:block fixed right-6 top-24 w-64 z-30"
    @mouseenter="onEnter"
    @mouseleave="onLeave"
  >
    <div class="rounded-xl border border-border bg-background/90 backdrop-blur-sm shadow-lg overflow-hidden">
      <button
        class="w-full flex items-center justify-between px-4 py-2.5 select-none group"
        @click="togglePin"
        :title="isPinned ? 'Свернуть (открепить)' : 'Раскрыть и закрепить'"
      >
        <span class="text-xs font-semibold uppercase tracking-wider text-muted group-hover:text-foreground transition-colors">
          Содержание
        </span>
        <svg
          class="w-3.5 h-3.5 text-muted transition-transform duration-200"
          :class="{ 'rotate-180': isExpanded }"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>
      <div v-if="!isExpanded" class="px-2 pb-2 space-y-0.5">
        <button
          v-if="prevItem"
          @click="scrollToHeading(prevItem.id)"
          class="block w-full text-left px-2 py-1 rounded-md text-xs text-muted/60 hover:text-foreground transition-colors truncate"
          :title="prevItem.text"
        >
          {{ prevItem.text }}
        </button>
        <button
          v-if="currentItem"
          @click="scrollToHeading(currentItem.id)"
          class="block w-full text-left px-2 py-1 rounded-md text-xs text-foreground font-medium bg-muted/60 truncate"
          :title="currentItem.text"
        >
          {{ currentItem.text }}
        </button>
        <button
          v-if="nextItem"
          @click="scrollToHeading(nextItem.id)"
          class="block w-full text-left px-2 py-1 rounded-md text-xs text-muted/60 hover:text-foreground transition-colors truncate"
          :title="nextItem.text"
        >
          {{ nextItem.text }}
        </button>
      </div>
      <div v-else ref="listRef" class="px-2 pb-2 max-h-[60vh] overflow-y-auto">
        <ArticleTOC :items="items" :active-id="activeId" hide-title />
      </div>
    </div>
  </div>
  <div class="2xl:hidden">
    <Transition name="popup">
      <div
        v-if="isPopupOpen"
        class="fixed inset-0 z-[95] flex items-end sm:items-center justify-center p-4"
      >
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="closePopup" />
        <div class="relative w-full max-w-sm max-h-[70vh] overflow-y-auto bg-background border border-border rounded-2xl shadow-2xl p-4 animate-slide-up">
          <div class="flex items-center justify-between mb-3">
            <p class="text-sm font-semibold">Содержание</p>
            <button
              @click="closePopup"
              class="p-1.5 rounded-lg hover:bg-muted/50 transition-colors text-muted hover:text-foreground"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <ArticleTOC :items="items" :active-id="activeId" @navigate="closePopup" />
        </div>
      </div>
    </Transition>
    <button
      @click="isPopupOpen = !isPopupOpen"
      class="fixed bottom-6 left-6 z-[90] w-11 h-11 flex items-center justify-center rounded-full bg-foreground text-background shadow-lg hover:opacity-85 hover:scale-110 transition-all duration-200"
      aria-label="Содержание"
    >
      <svg v-if="!isPopupOpen" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
      </svg>
      <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
</template>

<style scoped>
@keyframes slide-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-slide-up { animation: slide-up 0.2s ease-out; }
.popup-enter-active, .popup-leave-active { transition: opacity 0.2s ease; }
.popup-enter-from, .popup-leave-to { opacity: 0; }
</style>