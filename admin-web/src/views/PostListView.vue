<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Eye,
  FileQuestion,
  Loader2,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  Trash2,
} from '@lucide/vue'
import AppTopbar from '@/components/AppTopbar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import EmptyState from '@/components/EmptyState.vue'
import { usePostsStore } from '@/stores/posts'
import { useToastStore } from '@/stores/toast'
import { deletePost } from '@/api/posts'
import { formatRelative, formatDateTime } from '@/utils/format'
import type { Post } from '@/types/post'

const router = useRouter()
const store = usePostsStore()
const toast = useToastStore()

const searchTerm = ref('')
const statusFilter = ref<'all' | 'published' | 'draft'>('all')
const sortBy = ref<'created' | 'updated' | 'views'>('created')

const deleteTarget = ref<Post | null>(null)
const deleting = ref(false)

onMounted(() => {
  store.loadFirstPage()
})

const filteredPosts = computed(() => {
  let list = store.posts

  if (statusFilter.value !== 'all') {
    const wanted = statusFilter.value === 'published' ? 1 : 0
    list = list.filter((p) => p.status === wanted)
  }

  const term = searchTerm.value.trim().toLowerCase()
  if (term) {
    list = list.filter(
      (p) => p.title.toLowerCase().includes(term) || p.slug.toLowerCase().includes(term),
    )
  }

  const sorted = [...list]
  if (sortBy.value === 'views') {
    sorted.sort((a, b) => b.viewCount - a.viewCount)
  } else if (sortBy.value === 'updated') {
    sorted.sort((a, b) => (b.updatedAt || '').localeCompare(a.updatedAt || ''))
  } else {
    sorted.sort((a, b) => (b.createdAt || '').localeCompare(a.createdAt || ''))
  }
  return sorted
})

const publishedCount = computed(() => store.posts.filter((p) => p.status === 1).length)
const draftCount = computed(() => store.posts.filter((p) => p.status === 0).length)

function goEdit(post: Post) {
  router.push({ name: 'post-edit', params: { id: post.id } })
}

function requestDelete(post: Post) {
  deleteTarget.value = post
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await deletePost(deleteTarget.value.id)
    store.removeLocal(deleteTarget.value.id)
    toast.success(`已删除《${deleteTarget.value.title}》`)
    deleteTarget.value = null
  } catch (err) {
    toast.error(err instanceof Error ? err.message : '删除失败')
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <AppTopbar title="全部文章" :subtitle="`共加载 ${store.posts.length} 篇${store.hasMore ? '（还有更多）' : ''}`">
    <template #actions>
      <button
        type="button"
        class="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-surface-alt disabled:opacity-60"
        :disabled="store.loading"
        @click="store.loadFirstPage(true)"
      >
        <RefreshCw :size="15" :class="{ 'animate-spin': store.loading }" />
        刷新
      </button>
      <router-link
        :to="{ name: 'post-new' }"
        class="inline-flex items-center gap-1.5 rounded-lg bg-primary px-3.5 py-2 text-sm font-semibold text-white transition-colors hover:bg-primary-hover"
      >
        <Plus :size="16" />
        新建文章
      </router-link>
    </template>
  </AppTopbar>

  <main class="flex-1 px-4 py-5 sm:px-6 sm:py-6">
    <!-- 统计条 -->
    <div class="mb-5 grid grid-cols-3 gap-3">
      <div class="rounded-xl border border-border bg-surface px-4 py-3">
        <p class="text-xs text-ink-soft">已加载</p>
        <p class="mt-1 font-mono text-xl font-semibold tabular text-ink">{{ store.posts.length }}</p>
      </div>
      <div class="rounded-xl border border-border bg-surface px-4 py-3">
        <p class="text-xs text-ink-soft">已发布</p>
        <p class="mt-1 font-mono text-xl font-semibold tabular text-success">{{ publishedCount }}</p>
      </div>
      <div class="rounded-xl border border-border bg-surface px-4 py-3">
        <p class="text-xs text-ink-soft">草稿</p>
        <p class="mt-1 font-mono text-xl font-semibold tabular text-warning">{{ draftCount }}</p>
      </div>
    </div>

    <!-- 工具栏 -->
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="relative w-full sm:max-w-xs">
        <Search :size="15" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-faint" />
        <input
          v-model="searchTerm"
          type="search"
          placeholder="搜索标题或 slug"
          aria-label="搜索文章"
          class="w-full rounded-lg border border-border bg-surface py-2 pl-9 pr-3 text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
        />
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <div class="inline-flex rounded-lg border border-border bg-surface p-0.5 text-sm">
          <button
            v-for="opt in [
              { key: 'all', label: '全部' },
              { key: 'published', label: '已发布' },
              { key: 'draft', label: '草稿' },
            ]"
            :key="opt.key"
            type="button"
            class="rounded-md px-3 py-1.5 font-medium transition-colors"
            :class="statusFilter === opt.key ? 'bg-ink text-paper' : 'text-ink-soft hover:text-ink'"
            @click="statusFilter = opt.key as typeof statusFilter"
          >
            {{ opt.label }}
          </button>
        </div>

        <select
          v-model="sortBy"
          aria-label="排序方式"
          class="rounded-lg border border-border bg-surface px-2.5 py-2 text-sm text-ink-soft focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
        >
          <option value="created">最新创建</option>
          <option value="updated">最近更新</option>
          <option value="views">阅读量最高</option>
        </select>
      </div>
    </div>

    <!-- 加载态骨架 -->
    <div v-if="store.loading" class="space-y-2">
      <div v-for="i in 6" :key="i" class="h-16 animate-pulse rounded-xl bg-surface-alt"></div>
    </div>

    <!-- 错误态 -->
    <div v-else-if="store.error" class="rounded-xl border border-danger/30 bg-danger-soft px-4 py-8 text-center text-sm text-danger">
      {{ store.error }}
      <div class="mt-3">
        <button
          type="button"
          class="rounded-lg border border-danger/40 px-3 py-1.5 text-sm font-medium text-danger hover:bg-danger/10"
          @click="store.loadFirstPage(true)"
        >
          重试
        </button>
      </div>
    </div>

    <!-- 空态 -->
    <EmptyState
      v-else-if="filteredPosts.length === 0"
      :icon="FileQuestion"
      :title="store.posts.length === 0 ? '还没有文章' : '没有匹配的文章'"
      :description="store.posts.length === 0 ? '点击右上角「新建文章」开始写作' : '试试调整搜索词或筛选条件'"
    />

    <template v-else>
      <!-- 桌面表格 -->
      <div class="hidden overflow-hidden rounded-xl border border-border bg-surface md:block">
        <table class="w-full text-left text-sm">
          <thead class="border-b border-border bg-surface-alt text-xs text-ink-soft">
            <tr>
              <th class="px-4 py-3 font-medium">文章</th>
              <th class="w-28 px-4 py-3 font-medium">状态</th>
              <th class="w-24 px-4 py-3 font-medium">阅读量</th>
              <th class="w-40 px-4 py-3 font-medium">更新时间</th>
              <th class="w-28 px-4 py-3 text-right font-medium">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="post in filteredPosts" :key="post.id" class="group transition-colors hover:bg-surface-alt/60">
              <td class="max-w-0 px-4 py-3">
                <button type="button" class="block max-w-full text-left" @click="goEdit(post)">
                  <p class="truncate font-serif text-[15px] font-medium text-ink group-hover:text-primary">
                    {{ post.title }}
                  </p>
                  <p class="mt-0.5 truncate font-mono text-xs text-ink-faint">/{{ post.slug }}</p>
                </button>
              </td>
              <td class="px-4 py-3">
                <StatusBadge :status="post.status" size="sm" />
              </td>
              <td class="px-4 py-3 font-mono text-sm tabular text-ink-soft">
                <span class="inline-flex items-center gap-1">
                  <Eye :size="13" />
                  {{ post.viewCount }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-ink-soft" :title="formatDateTime(post.updatedAt)">
                {{ formatRelative(post.updatedAt) }}
              </td>
              <td class="px-4 py-3">
                <div class="flex justify-end gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100">
                  <button
                    type="button"
                    class="rounded-md p-1.5 text-ink-soft hover:bg-surface hover:text-primary"
                    aria-label="编辑"
                    title="编辑"
                    @click="goEdit(post)"
                  >
                    <Pencil :size="15" />
                  </button>
                  <button
                    type="button"
                    class="rounded-md p-1.5 text-ink-soft hover:bg-danger-soft hover:text-danger"
                    aria-label="删除"
                    title="删除"
                    @click="requestDelete(post)"
                  >
                    <Trash2 :size="15" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 移动端卡片 -->
      <div class="space-y-2.5 md:hidden">
        <div
          v-for="post in filteredPosts"
          :key="post.id"
          class="rounded-xl border border-border bg-surface p-4"
        >
          <div class="flex items-start justify-between gap-3">
            <button type="button" class="min-w-0 flex-1 text-left" @click="goEdit(post)">
              <p class="truncate font-serif text-[15px] font-medium text-ink">{{ post.title }}</p>
              <p class="mt-0.5 truncate font-mono text-xs text-ink-faint">/{{ post.slug }}</p>
            </button>
            <StatusBadge :status="post.status" size="sm" />
          </div>
          <div class="mt-3 flex items-center justify-between text-xs text-ink-soft">
            <span class="inline-flex items-center gap-1 font-mono tabular">
              <Eye :size="12" />
              {{ post.viewCount }} · {{ formatRelative(post.updatedAt) }}
            </span>
            <div class="flex gap-1">
              <button
                type="button"
                class="rounded-md p-1.5 text-ink-soft hover:bg-surface-alt hover:text-primary"
                aria-label="编辑"
                @click="goEdit(post)"
              >
                <Pencil :size="15" />
              </button>
              <button
                type="button"
                class="rounded-md p-1.5 text-ink-soft hover:bg-danger-soft hover:text-danger"
                aria-label="删除"
                @click="requestDelete(post)"
              >
                <Trash2 :size="15" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="store.hasMore" class="mt-5 flex justify-center">
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg border border-border px-4 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-surface-alt disabled:opacity-60"
          :disabled="store.loadingMore"
          @click="store.loadMore()"
        >
          <Loader2 v-if="store.loadingMore" :size="15" class="animate-spin" />
          {{ store.loadingMore ? '加载中…' : '加载更多' }}
        </button>
      </div>
    </template>
  </main>

  <ConfirmDialog
    :open="Boolean(deleteTarget)"
    danger
    title="删除这篇文章？"
    :description="deleteTarget ? `《${deleteTarget.title}》删除后无法恢复。` : ''"
    confirm-label="删除"
    :loading="deleting"
    @cancel="deleteTarget = null"
    @confirm="confirmDelete"
  />
</template>
