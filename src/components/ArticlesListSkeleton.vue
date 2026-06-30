<script setup lang="ts">
import ArticleCardSkeleton from './ArticleCardSkeleton.vue'
</script>

<template>
  <div class="space-y-6">
    <ArticleCardSkeleton v-for="i in 3" :key="i" />
  </div>
</template>