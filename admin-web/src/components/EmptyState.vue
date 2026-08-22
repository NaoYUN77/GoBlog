<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  icon: Component
  title: string
  description?: string
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed border-border-strong py-16 text-center">
    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-surface-alt text-ink-faint">
      <component :is="icon" :size="22" :stroke-width="1.75" />
    </div>
    <div class="max-w-xs">
      <p class="text-sm font-medium text-ink">{{ title }}</p>
      <p v-if="description" class="mt-1 text-sm text-ink-soft">{{ description }}</p>
    </div>
    <slot name="action" />
  </div>
</template>
