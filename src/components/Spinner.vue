<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg'
}>()

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-8 h-8',
  lg: 'w-12 h-12'
}
</script>

<template>
  <div 
    :class="[
      'animate-spin rounded-full border-2 border-muted border-t-foreground',
      sizeClasses[size || 'md']
    ]"
  />
</template>