<script setup lang="ts">
import { inject, useSlots } from 'vue'
import { Menu } from '@lucide/vue'
import SessionCountdown from './SessionCountdown.vue'
import { MOBILE_MENU_KEY } from '@/utils/injectionKeys'

defineProps<{ title: string; subtitle?: string }>()
const slots = useSlots()
const openMobileMenu = inject(MOBILE_MENU_KEY, () => {})
</script>

<template>
  <header
    class="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-border bg-paper/90 px-4 backdrop-blur sm:px-6"
  >
    <button
      type="button"
      class="rounded-md p-1.5 text-ink-soft hover:bg-surface-alt lg:hidden"
      aria-label="打开菜单"
      @click="openMobileMenu"
    >
      <Menu :size="20" />
    </button>

    <div class="min-w-0 flex-1">
      <h1 class="truncate font-serif text-lg font-semibold text-ink sm:text-xl">{{ title }}</h1>
      <p v-if="subtitle" class="truncate text-xs text-ink-soft">{{ subtitle }}</p>
    </div>

    <SessionCountdown />

    <div class="hidden items-center gap-2 sm:flex">
      <slot name="actions" />
    </div>
  </header>

  <div v-if="slots.actions" class="border-b border-border bg-paper px-4 pb-3 pt-1 sm:hidden">
    <slot name="actions" />
  </div>
</template>
