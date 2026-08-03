<script setup lang="ts">
import { scrollToHeading } from '../utils/scroll'

export interface TOCItem {
  id: string
  text: string
  level: number
}

withDefaults(defineProps<{
  items: TOCItem[]
  activeId: string
  hideTitle?: boolean
}>(), {
  hideTitle: false
})

const emit = defineEmits<{
  navigate: []
}>()

const onClick = (id: string) => {
  scrollToHeading(id)
  emit('navigate')
}
</script>

<template>
  <nav aria-label="Содержание статьи" class="select-none">
    <p v-if="!hideTitle" class="text-xs font-semibold uppercase tracking-wider text-muted mb-3">
    </p>
    <ul class="space-y-0.5">
      <li v-for="item in items" :key="item.id">
        <button
          @click="onClick(item.id)"
          :data-active="activeId === item.id"
          :class="[
            'block w-full text-left py-1 px-2 rounded-md text-xs transition-all duration-150 truncate',
            item.level === 3 ? 'ml-3' : '',
            activeId === item.id
              ? 'text-foreground font-bold bg-muted/60'
              : 'text-muted hover:text-foreground hover:bg-muted/30'
          ]"
          :title="item.text"
        >
          {{ item.text }}
        </button>
      </li>
    </ul>
  </nav>
</template>