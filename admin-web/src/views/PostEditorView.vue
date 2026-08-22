<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import {
  ArrowLeft,
  ImageOff,
  Loader2,
  Save,
  Sparkles,
  Trash2,
} from '@lucide/vue'
import AppTopbar from '@/components/AppTopbar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { createPost, deletePost, fetchPost, updatePost } from '@/api/posts'
import { useToastStore } from '@/stores/toast'
import { usePostsStore } from '@/stores/posts'
import { estimateReadingStats, formatDateTime, slugify } from '@/utils/format'
import type { Post, PostStatus, UpdatePostPayload } from '@/types/post'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const postsStore = usePostsStore()

const postId = computed(() => (typeof route.params.id === 'string' ? route.params.id : null))
const isEditMode = computed(() => Boolean(postId.value))

const loading = ref(isEditMode.value)
const saving = ref(false)
const deleting = ref(false)
const notFound = ref(false)
const showDeleteConfirm = ref(false)
const showLeaveConfirm = ref(false)

const original = ref<Post | null>(null)

const form = reactive({
  title: '',
  slug: '',
  summary: '',
  content: '',
  coverUrl: '',
  categoryId: '' as string | number,
  status: 0 as PostStatus,
})

const coverBroken = ref(false)

function fillForm(post: Post) {
  form.title = post.title
  form.slug = post.slug
  form.summary = post.summary ?? ''
  form.content = post.content
  form.coverUrl = post.coverUrl ?? ''
  form.categoryId = post.categoryId ?? ''
  form.status = post.status
  coverBroken.value = false
}

onMounted(async () => {
  if (isEditMode.value && postId.value) {
    loading.value = true
    try {
      const post = await fetchPost(postId.value)
      original.value = post
      fillForm(post)
    } catch (err) {
      notFound.value = true
      toast.error(err instanceof Error ? err.message : '文章不存在或已被删除')
    } finally {
      loading.value = false
    }
  }
})

const isDirty = computed(() => {
  if (!isEditMode.value) {
    return Boolean(form.title || form.slug || form.summary || form.content || form.coverUrl || form.categoryId)
  }
  if (!original.value) return false
  return (
    form.title !== original.value.title ||
    form.slug !== original.value.slug ||
    form.summary !== (original.value.summary ?? '') ||
    form.content !== original.value.content ||
    form.coverUrl !== (original.value.coverUrl ?? '') ||
    String(form.categoryId) !== String(original.value.categoryId ?? '') ||
    form.status !== original.value.status
  )
})

const readingStats = computed(() => estimateReadingStats(form.content))

function generateSlug() {
  if (!form.title.trim()) {
    toast.info('请先填写标题')
    return
  }
  form.slug = slugify(form.title)
}

function validate(): string | null {
  if (!form.title.trim()) return '请填写文章标题'
  if (form.title.length > 200) return '标题最多 200 个字符'
  if (!form.slug.trim()) return '请填写 URL 标识（slug）'
  if (form.slug.length > 200) return 'slug 最多 200 个字符'
  if (!form.content.trim()) return '请填写文章正文'
  if (form.summary.length > 500) return '摘要最多 500 个字符'
  return null
}

let justSaved = false

async function handleSave() {
  const error = validate()
  if (error) {
    toast.error(error)
    return
  }

  saving.value = true
  try {
    if (isEditMode.value && postId.value) {
      const payload = buildUpdatePayload()
      if (Object.keys(payload).length === 0) {
        toast.info('没有需要保存的更改')
        return
      }
      await updatePost(postId.value, payload)
      const refreshed = await fetchPost(postId.value)
      original.value = refreshed
      fillForm(refreshed)
      postsStore.upsertLocal(refreshed)
      toast.success('文章已更新')
    } else {
      await createPost({
        title: form.title.trim(),
        slug: form.slug.trim(),
        summary: form.summary.trim() || null,
        content: form.content,
        cover_url: form.coverUrl.trim() || null,
        category_id: form.categoryId === '' ? null : Number(form.categoryId),
        status: form.status,
      })
      justSaved = true
      toast.success('文章已创建')
      postsStore.loadFirstPage(true)
      router.push({ name: 'posts' })
    }
  } catch (err) {
    toast.error(err instanceof Error ? err.message : '保存失败')
  } finally {
    saving.value = false
  }
}

function buildUpdatePayload(): UpdatePostPayload {
  const payload: UpdatePostPayload = {}
  const base = original.value
  if (!base) return payload

  if (form.title !== base.title) payload.title = form.title.trim()
  if (form.slug !== base.slug) payload.slug = form.slug.trim()
  if (form.summary !== (base.summary ?? '')) payload.summary = form.summary.trim()
  if (form.content !== base.content) payload.content = form.content
  if (form.coverUrl !== (base.coverUrl ?? '')) payload.cover_url = form.coverUrl.trim()
  // 单篇详情接口不返回 category_id，若用户未主动修改该字段则不下发，
  // 避免把后端已有的分类误清空。
  if (String(form.categoryId) !== String(base.categoryId ?? '')) {
    payload.category_id = form.categoryId === '' ? null : Number(form.categoryId)
  }
  if (form.status !== base.status) payload.status = form.status

  return payload
}

async function confirmDelete() {
  if (!postId.value) return
  deleting.value = true
  try {
    await deletePost(postId.value)
    postsStore.removeLocal(postId.value)
    justSaved = true
    toast.success('文章已删除')
    router.push({ name: 'posts' })
  } catch (err) {
    toast.error(err instanceof Error ? err.message : '删除失败')
  } finally {
    deleting.value = false
    showDeleteConfirm.value = false
  }
}

let resolveLeave: ((value: boolean) => void) | null = null

onBeforeRouteLeave(() => {
  if (justSaved || !isDirty.value) return true
  showLeaveConfirm.value = true
  return new Promise<boolean>((resolve) => {
    resolveLeave = resolve
  })
})

function handleLeaveConfirm(result: boolean) {
  showLeaveConfirm.value = false
  resolveLeave?.(result)
  resolveLeave = null
}

function handleBeforeUnload(e: BeforeUnloadEvent) {
  if (isDirty.value && !justSaved) {
    e.preventDefault()
    e.returnValue = ''
  }
}
onMounted(() => window.addEventListener('beforeunload', handleBeforeUnload))
onBeforeUnmount(() => window.removeEventListener('beforeunload', handleBeforeUnload))
</script>

<template>
  <AppTopbar :title="isEditMode ? '编辑文章' : '写文章'" :subtitle="isEditMode ? `/${form.slug || '...'}` : undefined">
    <template #actions>
      <router-link
        :to="{ name: 'posts' }"
        class="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-surface-alt"
      >
        <ArrowLeft :size="15" />
        返回列表
      </router-link>
      <button
        v-if="isEditMode"
        type="button"
        class="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm font-medium text-danger transition-colors hover:bg-danger-soft"
        @click="showDeleteConfirm = true"
      >
        <Trash2 :size="15" />
        删除
      </button>
      <button
        type="button"
        class="inline-flex items-center gap-1.5 rounded-lg bg-primary px-3.5 py-2 text-sm font-semibold text-white transition-colors hover:bg-primary-hover disabled:cursor-not-allowed disabled:opacity-70"
        :disabled="saving || loading"
        @click="handleSave"
      >
        <Loader2 v-if="saving" :size="15" class="animate-spin" />
        <Save v-else :size="15" />
        {{ saving ? '保存中…' : '保存' }}
      </button>
    </template>
  </AppTopbar>

  <main class="flex-1 px-4 py-5 sm:px-6 sm:py-6">
    <div v-if="loading" class="space-y-3">
      <div class="h-10 w-2/3 animate-pulse rounded-lg bg-surface-alt"></div>
      <div class="h-96 animate-pulse rounded-xl bg-surface-alt"></div>
    </div>

    <div v-else-if="notFound" class="rounded-xl border border-dashed border-border-strong py-16 text-center">
      <p class="text-sm font-medium text-ink">文章不存在或已被删除</p>
      <router-link :to="{ name: 'posts' }" class="mt-3 inline-block text-sm font-medium text-primary hover:underline">
        返回文章列表
      </router-link>
    </div>

    <div v-else class="grid grid-cols-1 gap-5 lg:grid-cols-[1fr_320px]">
      <!-- 主栏：标题 / slug / 正文 -->
      <div class="min-w-0 space-y-4">
        <div class="rounded-xl border border-border bg-surface p-4">
          <label for="title" class="mb-1.5 block text-sm font-medium text-ink">标题</label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            maxlength="200"
            placeholder="给文章起一个标题"
            class="w-full border-0 bg-transparent p-0 font-serif text-xl font-semibold text-ink placeholder:text-ink-faint focus:outline-none focus:ring-0"
          />

          <div class="mt-3 flex items-center gap-2 border-t border-border pt-3">
            <span class="font-mono text-sm text-ink-faint">/</span>
            <input
              id="slug"
              v-model="form.slug"
              type="text"
              maxlength="200"
              placeholder="url-slug"
              class="min-w-0 flex-1 rounded-md border border-border bg-paper px-2.5 py-1.5 font-mono text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
            />
            <button
              type="button"
              class="inline-flex shrink-0 items-center gap-1 rounded-md border border-border px-2.5 py-1.5 text-xs font-medium text-ink-soft transition-colors hover:bg-surface-alt"
              title="根据标题生成 slug"
              @click="generateSlug"
            >
              <Sparkles :size="12" />
              自动生成
            </button>
          </div>
        </div>

        <div class="overflow-hidden rounded-xl border border-border">
          <MdEditor
            v-model="form.content"
            language="zh-CN"
            :toolbars-exclude="['github']"
            style="height: 560px"
            placeholder="用 Markdown 写正文…"
          />
        </div>

        <p class="flex items-center gap-1.5 px-1 font-mono text-xs tabular text-ink-faint">
          {{ readingStats.chars.toLocaleString() }} 字 · 约 {{ readingStats.minutes }} 分钟阅读
        </p>
      </div>

      <!-- 侧栏：发布状态 / 摘要 / 封面 / 分类 -->
      <aside class="space-y-4">
        <div class="rounded-xl border border-border bg-surface p-4">
          <h2 class="mb-3 text-sm font-medium text-ink">发布状态</h2>
          <div class="inline-flex w-full rounded-lg border border-border bg-paper p-0.5 text-sm">
            <button
              type="button"
              class="flex-1 rounded-md px-3 py-1.5 font-medium transition-colors"
              :class="form.status === 0 ? 'bg-warning-soft text-warning' : 'text-ink-soft hover:text-ink'"
              @click="form.status = 0"
            >
              草稿
            </button>
            <button
              type="button"
              class="flex-1 rounded-md px-3 py-1.5 font-medium transition-colors"
              :class="form.status === 1 ? 'bg-success-soft text-success' : 'text-ink-soft hover:text-ink'"
              @click="form.status = 1"
            >
              已发布
            </button>
          </div>

          <dl v-if="isEditMode && original" class="mt-3 space-y-1 border-t border-border pt-3 text-xs text-ink-soft">
            <div class="flex justify-between">
              <dt>创建于</dt>
              <dd class="font-mono tabular">{{ formatDateTime(original.createdAt) }}</dd>
            </div>
            <div class="flex justify-between">
              <dt>最近更新</dt>
              <dd class="font-mono tabular">{{ formatDateTime(original.updatedAt) }}</dd>
            </div>
            <div class="flex justify-between">
              <dt>阅读量</dt>
              <dd class="font-mono tabular">{{ original.viewCount }}</dd>
            </div>
          </dl>
        </div>

        <div class="rounded-xl border border-border bg-surface p-4">
          <label for="summary" class="mb-1.5 flex items-center justify-between text-sm font-medium text-ink">
            摘要
            <span class="font-mono text-xs font-normal text-ink-faint">{{ form.summary.length }}/500</span>
          </label>
          <textarea
            id="summary"
            v-model="form.summary"
            maxlength="500"
            rows="4"
            placeholder="用于列表页展示的简要说明（可选）"
            class="w-full resize-none rounded-lg border border-border bg-paper px-3 py-2 text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
          />
        </div>

        <div class="rounded-xl border border-border bg-surface p-4">
          <label for="cover" class="mb-1.5 block text-sm font-medium text-ink">封面图片地址</label>
          <input
            id="cover"
            v-model="form.coverUrl"
            type="url"
            placeholder="https://…"
            class="w-full rounded-lg border border-border bg-paper px-3 py-2 text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
            @input="coverBroken = false"
          />
          <div class="mt-2.5 flex aspect-video items-center justify-center overflow-hidden rounded-lg border border-dashed border-border-strong bg-paper">
            <img
              v-if="form.coverUrl && !coverBroken"
              :src="form.coverUrl"
              alt="封面预览"
              class="h-full w-full object-cover"
              @error="coverBroken = true"
            />
            <div v-else class="flex flex-col items-center gap-1.5 text-ink-faint">
              <ImageOff :size="20" :stroke-width="1.75" />
              <span class="text-xs">{{ form.coverUrl ? '图片加载失败' : '暂无封面' }}</span>
            </div>
          </div>
        </div>

        <div class="rounded-xl border border-border bg-surface p-4">
          <label for="category" class="mb-1.5 block text-sm font-medium text-ink">分类 ID</label>
          <input
            id="category"
            v-model="form.categoryId"
            type="number"
            inputmode="numeric"
            placeholder="可选，数字分类标识"
            class="w-full rounded-lg border border-border bg-paper px-3 py-2 font-mono text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
          />
        </div>
      </aside>
    </div>
  </main>

  <ConfirmDialog 
    :open="showDeleteConfirm"
    danger
    title="删除这篇文章？"
    description="删除后无法恢复，请确认。"
    confirm-label="删除"
    :loading="deleting"
    @cancel="showDeleteConfirm = false"
    @confirm="confirmDelete"
  />

  <ConfirmDialog
    :open="showLeaveConfirm"
    title="有未保存的修改"
    description="离开后本次修改的内容将会丢失。"
    confirm-label="离开"
    cancel-label="继续编辑"
    danger
    @cancel="handleLeaveConfirm(false)"
    @confirm="handleLeaveConfirm(true)"
  />
</template>
