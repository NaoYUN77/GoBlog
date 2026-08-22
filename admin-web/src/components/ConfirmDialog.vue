<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { TriangleAlert } from '@lucide/vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    description?: string
    confirmLabel?: string
    cancelLabel?: string
    danger?: boolean
    loading?: boolean
  }>(),
  {
    description: '',
    confirmLabel: '确认',
    cancelLabel: '取消',
    danger: false,
    loading: false,
  },
)

const emit = defineEmits<{ confirm: []; cancel: [] }>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open && !props.loading) emit('cancel')
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Transition name="dialog">
    <div
      v-if="open"
      class="fixed inset-0 z-[90] flex items-center justify-center bg-ink/40 px-4 backdrop-blur-[1px]"
      @mousedown.self="!loading && emit('cancel')"
    >
      <div
        class="w-full max-w-sm rounded-xl border border-border bg-surface p-5 shadow-pop"
        role="alertdialog"
        aria-modal="true"
        :aria-label="title"
      >
        <div class="flex items-start gap-3">
          <div
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full"
            :class="danger ? 'bg-danger-soft text-danger' : 'bg-primary-soft text-primary'"
          >
            <TriangleAlert :size="18" />
          </div>
          <div class="min-w-0 flex-1 pt-1">
            <h2 class="font-serif text-base font-semibold text-ink">{{ title }}</h2>
            <p v-if="description" class="mt-1 text-sm leading-5 text-ink-soft">{{ description }}</p>
          </div>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <button
            type="button"
            class="rounded-lg border border-border px-3.5 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-surface-alt disabled:opacity-50"
            :disabled="loading"
            @click="emit('cancel')"
          >
            {{ cancelLabel }}
          </button>
          <button
            type="button"
            class="rounded-lg px-3.5 py-2 text-sm font-medium text-white transition-colors disabled:opacity-60"
            :class="danger ? 'bg-danger hover:bg-danger/90' : 'bg-primary hover:bg-primary-hover'"
            :disabled="loading"
            @click="emit('confirm')"
          >
            {{ loading ? '处理中…' : confirmLabel }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.15s ease;
}
.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}
.dialog-enter-active > div,
.dialog-leave-active > div {
  transition: transform 0.15s ease;
}
.dialog-enter-from > div,
.dialog-leave-to > div {
  transform: scale(0.97);
}
</style>
