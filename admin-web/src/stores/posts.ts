import { defineStore } from 'pinia'
import { fetchPostList } from '@/api/posts'
import type { Post } from '@/types/post'

// 后端 /admin/list 单页最多 100 条，且不返回总数，
// 因此这里以「分批拉取全部 + 前端筛选/排序」的方式提供更好的检索体验。
const PAGE_SIZE = 100

export const usePostsStore = defineStore('posts', {
  state: () => ({
    posts: [] as Post[],
    loading: false,
    loadingMore: false,
    error: '' as string,
    nextPage: 1,
    hasMore: true,
    loadedOnce: false,
  }),
  actions: {
    async loadFirstPage(force = false) {
      if (this.loadedOnce && !force) return
      this.loading = true
      this.error = ''
      try {
        const result = await fetchPostList(1, PAGE_SIZE)
        this.posts = result.posts
        this.nextPage = 2
        this.hasMore = result.posts.length >= PAGE_SIZE
        this.loadedOnce = true
      } catch (err) {
        this.error = err instanceof Error ? err.message : '加载文章列表失败'
      } finally {
        this.loading = false
      }
    },
    async loadMore() {
      if (!this.hasMore || this.loadingMore) return
      this.loadingMore = true
      try {
        const result = await fetchPostList(this.nextPage, PAGE_SIZE)
        this.posts = [...this.posts, ...result.posts]
        this.nextPage += 1
        this.hasMore = result.posts.length >= PAGE_SIZE
      } catch (err) {
        this.error = err instanceof Error ? err.message : '加载更多文章失败'
      } finally {
        this.loadingMore = false
      }
    },
    upsertLocal(post: Post) {
      const idx = this.posts.findIndex((p) => p.id === post.id)
      if (idx >= 0) {
        this.posts.splice(idx, 1, post)
      } else {
        this.posts.unshift(post)
      }
    },
    removeLocal(id: string) {
      this.posts = this.posts.filter((p) => p.id !== id)
    },
  },
})
