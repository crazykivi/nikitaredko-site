<template>
  <div class="divide-y divide-border">
    <div v-for="i in 4" :key="i" class="grid animate-pulse gap-4 py-7 first:pt-0 md:grid-cols-[88px_1fr] md:gap-8">
      <div class="space-y-2">
        <div class="h-3.5 w-6 rounded bg-border" />
        <div class="h-3 w-14 rounded bg-border" />
      </div>
      <div class="space-y-3">
        <div class="h-6 w-3/4 rounded bg-border" />
        <div class="h-4 w-full max-w-2xl rounded bg-border" />
        <div class="h-4 w-2/3 max-w-2xl rounded bg-border" />
        <div class="h-3 w-1/3 rounded bg-border" />
      </div>
    </div>
  </div>
</template>