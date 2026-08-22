<script setup lang="ts">
import { CheckCircle2, X, XCircle, Info } from '@lucide/vue'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
</script>

<template>
  <div
    class="pointer-events-none fixed inset-x-0 top-4 z-[100] flex flex-col items-center gap-2 px-4 sm:items-end sm:right-4 sm:left-auto"
  >
    <TransitionGroup name="toast">
      <div
        v-for="item in toast.items"
        :key="item.id"
        class="pointer-events-auto flex w-full max-w-sm items-start gap-2.5 rounded-lg border bg-surface px-4 py-3 shadow-pop sm:w-96"
        :class="{
          'border-success/30': item.kind === 'success',
          'border-danger/30': item.kind === 'error',
          'border-border': item.kind === 'info',
        }"
        role="status"
      >
        <CheckCircle2 v-if="item.kind === 'success'" :size="18" class="mt-0.5 shrink-0 text-success" />
        <XCircle v-else-if="item.kind === 'error'" :size="18" class="mt-0.5 shrink-0 text-danger" />
        <Info v-else :size="18" class="mt-0.5 shrink-0 text-ink-soft" />
        <p class="flex-1 text-sm leading-5 text-ink">{{ item.message }}</p>
        <button
          type="button"
          class="shrink-0 rounded p-0.5 text-ink-faint transition-colors hover:bg-surface-alt hover:text-ink"
          aria-label="关闭提示"
          @click="toast.dismiss(item.id)"
        >
          <X :size="15" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition:
    transform 0.2s ease,
    opacity 0.2s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(8px);
}
.toast-leave-active {
  position: absolute;
}
</style>
