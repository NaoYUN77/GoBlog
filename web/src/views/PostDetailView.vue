<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import MarkdownIt from 'markdown-it'
import { fetchPostBySlug } from '@/api/posts'
import type { PublicPost } from '@/types/post'
import { formatDate, formatViewCount } from '@/utils/format'

const route = useRoute()

const post = ref<PublicPost | null>(null)
const loading = ref(true)
const error = ref('')

const md = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: false,
})

const html = computed(() => (post.value ? md.render(post.value.content) : ''))

async function load() {
  const slug = String(route.params.slug ?? '')
  loading.value = true
  error.value = ''
  post.value = null
  try {
    post.value = await fetchPostBySlug(slug)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '文章加载失败'
  } finally {
    loading.value = false
  }
}

watch(() => route.params.slug, load, { immediate: true })
</script>

<template>
  <article>
    <RouterLink
      to="/"
      class="inline-flex min-h-11 items-center gap-1 font-mono text-xs uppercase tracking-wider text-ink-faint transition-colors hover:text-primary"
    >
      ← 返回列表
    </RouterLink>

    <!-- 加载中 -->
    <div v-if="loading" class="mt-10 space-y-4" aria-hidden="true">
      <div class="h-9 w-4/5 animate-pulse rounded bg-surface-alt" />
      <div class="h-4 w-44 animate-pulse rounded bg-surface-alt" />
      <div class="mt-10 space-y-3">
        <div class="h-4 w-full animate-pulse rounded bg-surface-alt" />
        <div class="h-4 w-11/12 animate-pulse rounded bg-surface-alt" />
        <div class="h-4 w-3/4 animate-pulse rounded bg-surface-alt" />
      </div>
    </div>

    <!-- 错误：草稿/不存在的 slug 后端返回 404 -->
    <div v-else-if="error" class="mt-20 text-center">
      <p class="text-ink-soft">{{ error }}</p>
      <RouterLink
        to="/"
        class="mt-6 inline-flex items-center justify-center rounded-full border border-border-strong bg-surface px-5 py-2 text-sm font-medium text-ink transition-colors hover:border-primary hover:text-primary"
      >
        返回列表
      </RouterLink>
    </div>

    <!-- 正文 -->
    <template v-else-if="post">
      <header class="mt-8">
        <h1 class="font-serif text-4xl font-bold leading-snug tracking-wide text-ink">
          {{ post.title }}
        </h1>
        <p class="mt-5 font-mono text-xs tracking-wider text-ink-faint">
          {{ formatDate(post.createdAt) }} · {{ formatViewCount(post.viewCount) }} 次阅读
        </p>
      </header>

      <div
        class="prose prose-stone mt-10 max-w-none border-t border-border pt-10"
        v-html="html"
      />
    </template>
  </article>
</template>
