<script setup lang="ts">
import { Compass } from '@lucide/vue'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-paper px-4 text-center">
    <div class="flex h-14 w-14 items-center justify-center rounded-full bg-surface-alt text-ink-faint">
      <Compass :size="26" :stroke-width="1.75" />
    </div>
    <div>
      <h1 class="font-serif text-xl font-semibold text-ink">页面不存在</h1>
      <p class="mt-1 text-sm text-ink-soft">请检查地址是否正确</p>
    </div>
    <router-link
      to="/"
      class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-primary-hover"
    >
      返回首页
    </router-link>
  </div>
</template>
