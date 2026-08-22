<script setup lang="ts">
import { computed } from 'vue'
import { PenLine, Stamp } from '@lucide/vue'
import type { PostStatus } from '@/types/post'

const props = defineProps<{
  status: PostStatus
  size?: 'sm' | 'md'
}>()

const isPublished = computed(() => props.status === 1)
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full font-medium leading-none"
    :class="[
      size === 'sm' ? 'px-2 py-1 text-xs' : 'px-2.5 py-1.5 text-xs',
      isPublished
        ? 'bg-success-soft text-success border border-transparent'
        : 'bg-transparent text-warning border border-dashed border-warning/60',
    ]"
  >
    <Stamp v-if="isPublished" :size="12" :stroke-width="2.25" />
    <PenLine v-else :size="12" :stroke-width="2.25" />
    {{ isPublished ? '已发布' : '草稿' }}
  </span>
</template>
