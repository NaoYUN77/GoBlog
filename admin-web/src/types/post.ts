// 文章状态：0 草稿，1 已发布（与后端 tb_post.status 保持一致）
export type PostStatus = 0 | 1

export interface Post {
  id: string
  title: string
  slug: string
  summary: string | null
  content: string
  coverUrl: string | null
  categoryId: number | null
  status: PostStatus
  viewCount: number
  createdAt: string | null
  updatedAt: string | null
}

export interface CreatePostPayload {
  title: string
  slug: string
  summary?: string | null
  content: string
  cover_url?: string | null
  category_id?: number | null
  status: PostStatus
}

export interface UpdatePostPayload {
  title?: string
  slug?: string
  summary?: string | null
  content?: string
  cover_url?: string | null
  category_id?: number | null
  status?: PostStatus
}
