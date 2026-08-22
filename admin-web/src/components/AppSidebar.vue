<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { FileStack, LogOut, SquarePen, X } from '@lucide/vue'
import { useAuthStore } from '@/stores/auth'

defineProps<{ mobileOpen: boolean }>()
const emit = defineEmits<{ close: []; logout: [] }>()

const route = useRoute()
const auth = useAuthStore()

const isPostsActive = computed(() =>
  ['posts', 'post-new', 'post-edit'].includes(String(route.name)),
)

const initial = computed(() => (auth.username ? auth.username.slice(0, 1).toUpperCase() : '?'))
</script>

<template>
  <!-- 移动端遮罩 -->
  <Transition name="fade">
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-30 bg-ink/40 lg:hidden"
      @click="emit('close')"
    />
  </Transition>

  <aside
    class="fixed inset-y-0 left-0 z-40 flex w-64 shrink-0 -translate-x-full flex-col border-r border-border bg-surface-alt transition-transform duration-200 lg:static lg:translate-x-0"
    :class="{ 'translate-x-0': mobileOpen }"
  >
    <div class="flex h-16 items-center justify-between border-b border-border px-5">
      <div class="flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-ink text-paper">
          <SquarePen :size="16" :stroke-width="2.25" />
        </div>
        <span class="font-serif text-lg font-semibold tracking-tight text-ink">笔耕</span>
      </div>
      <button
        type="button"
        class="rounded-md p-1.5 text-ink-soft hover:bg-surface lg:hidden"
        aria-label="关闭菜单"
        @click="emit('close')"
      >
        <X :size="18" />
      </button>
    </div>

    <nav class="flex-1 space-y-1 overflow-y-auto p-3">
      <router-link
        :to="{ name: 'posts' }"
        class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
        :class="
          isPostsActive
            ? 'bg-primary-soft text-primary'
            : 'text-ink-soft hover:bg-surface hover:text-ink'
        "
        @click="emit('close')"
      >
        <FileStack :size="17" :stroke-width="2.1" />
        全部文章
      </router-link>

      <router-link
        :to="{ name: 'post-new' }"
        class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
        :class="
          route.name === 'post-new'
            ? 'bg-primary-soft text-primary'
            : 'text-ink-soft hover:bg-surface hover:text-ink'
        "
        @click="emit('close')"
      >
        <SquarePen :size="17" :stroke-width="2.1" />
        写文章
      </router-link>
    </nav>

    <div class="border-t border-border p-3">
      <div class="flex items-center gap-2.5 rounded-lg px-2 py-2">
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-ink font-serif text-sm font-semibold text-paper"
        >
          {{ initial }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-ink">{{ auth.username || '管理员' }}</p>
          <p class="text-xs text-ink-faint">博客管理员</p>
        </div>
        <button
          type="button"
          class="shrink-0 rounded-md p-1.5 text-ink-faint transition-colors hover:bg-surface hover:text-danger"
          aria-label="退出登录"
          title="退出登录"
          @click="emit('logout')"
        >
          <LogOut :size="16" />
        </button>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
