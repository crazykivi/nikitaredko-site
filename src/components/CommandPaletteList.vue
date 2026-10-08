<script setup lang="ts">
import { computed, nextTick, ref, watch } from "vue";
import type { Command, CommandGroup } from "../types/commandPalette";
import CommandPaletteItem from "./CommandPaletteItem.vue";

const props = defineProps<{
  groups: CommandGroup[];
  selectedIndex: number;
  expandedId: string | null;
}>();

const emit = defineEmits<{
  execute: [cmd: Command];
  select: [index: number];
  "pointer-down": [id: string];
  "pointer-end": [];
}>();

const listRef = ref<HTMLElement | null>(null);
let lastHoveredId: string | null = null;

const flatCommands = computed(() => props.groups.flatMap((group) => group.commands));

const indexById = computed(() => {
  const map = new Map<string, number>();

  flatCommands.value.forEach((cmd, index) => {
    map.set(cmd.id, index);
  });

  return map;
});

const onMouseMove = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  const btn = target.closest("[data-cmd-id]") as HTMLElement | null;
  const id = btn?.getAttribute("data-cmd-id") ?? null;

  if (id === lastHoveredId) return;

  lastHoveredId = id;

  if (!id) return;

  const index = indexById.value.get(id);

  if (index !== undefined && index !== props.selectedIndex) {
    emit("select", index);
  }
};

const onMouseLeave = () => {
  lastHoveredId = null;
};

const scrollToSelected = () => {
  nextTick(() => {
    const el = listRef.value?.querySelector(
      '[data-selected="true"]'
    ) as HTMLElement | null;
    el?.scrollIntoView({ block: "nearest" });
  });
};

watch(() => props.selectedIndex, scrollToSelected);
</script>

<template>
  <div
    ref="listRef"
    class="max-h-[50vh] overflow-y-auto py-2"
    @mousemove="onMouseMove"
    @mouseleave="onMouseLeave"
  >
    <div v-if="groups.length === 0" class="px-5 py-8 text-center text-muted text-sm">
      Ничего не найдено
    </div>
    <div v-for="group in groups" :key="group.name" class="mb-2">
      <div
        class="px-5 py-1.5 text-[10px] font-semibold uppercase tracking-wider text-muted/70"
      >
        {{ group.name }}
      </div>
      <CommandPaletteItem
        v-for="cmd in group.commands"
        :key="cmd.id"
        :command="cmd"
        :selected="indexById.get(cmd.id) === selectedIndex"
        :expanded="expandedId === cmd.id"
        @execute="emit('execute', cmd)"
        @pointer-down="emit('pointer-down', cmd.id)"
        @pointer-end="emit('pointer-end')"
      />
    </div>
  </div>
</template>
