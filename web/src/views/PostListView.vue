<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchPostList } from '@/api/posts'
import type { PublicPost } from '@/types/post'
import { formatDate, formatViewCount } from '@/utils/format'

const PAGE_SIZE = 100

const posts = ref<PublicPost[]>([])
const page = ref(1)
const loading = ref(true)
const loadingMore = ref(false)
const error = ref('')
const hasMore = ref(true)

async function loadFirstPage() {
  loading.value = true
  error.value = ''
  try {
    const { posts: list } = await fetchPostList(1, PAGE_SIZE)
    posts.value = list
    page.value = 1
    hasMore.value = list.length === PAGE_SIZE
  } catch (e) {
    error.value = e instanceof Error ? e.message : '文章加载失败'
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value) return
  loadingMore.value = true
  try {
    const next = page.value + 1
    const { posts: list } = await fetchPostList(next, PAGE_SIZE)
    posts.value = [...posts.value, ...list]
    page.value = next
    hasMore.value = list.length === PAGE_SIZE
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loadingMore.value = false
  }
}

onMounted(loadFirstPage)
</script>

<template>
  <section aria-label="文章列表">
    <h1 class="font-serif text-3xl font-bold tracking-wide text-ink">文章</h1>

    <!-- 加载中：骨架屏 -->
    <div v-if="loading" class="mt-10 space-y-2" aria-hidden="true">
      <div v-for="n in 4" :key="n" class="animate-pulse border-b border-border py-7">
        <div class="h-6 w-2/3 rounded bg-surface-alt" />
        <div class="mt-3 h-4 w-5/6 rounded bg-surface-alt" />
        <div class="mt-4 h-3 w-24 rounded bg-surface-alt" />
      </div>
    </div>

    <!-- 错误 -->
    <div v-else-if="error" class="mt-20 text-center">
      <p class="text-ink-soft">{{ error }}</p>
      <button
        type="button"
        class="mt-6 inline-flex items-center justify-center rounded-full border border-border-strong bg-surface px-5 py-2 text-sm font-medium text-ink transition-colors hover:border-primary hover:text-primary"
        @click="loadFirstPage"
      >
        重试
      </button>
    </div>

    <!-- 空状态 -->
    <p v-else-if="posts.length === 0" class="mt-20 text-center text-ink-soft">
      还没有文章，请稍后再来。
    </p>

    <!-- 文章列表 -->
    <template v-else>
      <ul class="mt-6">
        <li v-for="post in posts" :key="post.slug" class="border-b border-border">
          <RouterLink
            :to="`/post/${post.slug}`"
            class="group block py-7 outline-offset-4"
          >
            <h2
              class="font-serif text-2xl font-semibold leading-snug text-ink transition-colors duration-150 group-hover:text-primary"
            >
              {{ post.title }}
            </h2>
            <p v-if="post.summary" class="mt-2 line-clamp-2 text-sm leading-relaxed text-ink-soft">
              {{ post.summary }}
            </p>
            <p
              class="mt-4 font-mono text-xs tracking-wider text-ink-faint transition-colors group-hover:text-ink-soft"
            >
              {{ formatDate(post.createdAt) }} · {{ formatViewCount(post.viewCount) }} 次阅读
            </p>
          </RouterLink>
        </li>
      </ul>

      <div class="mt-10 text-center">
        <button
          v-if="hasMore"
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-full border border-border-strong bg-surface px-6 py-2 text-sm font-medium text-ink transition-colors hover:border-primary hover:text-primary disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="loadingMore"
          @click="loadMore"
        >
          {{ loadingMore ? '加载中…' : '加载更多' }}
        </button>
        <p v-else class="font-mono text-xs tracking-wider text-ink-faint">—— 已显示全部文章 ——</p>
      </div>
    </template>
  </section>
</template>
