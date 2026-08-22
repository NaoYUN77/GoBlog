import http from './http'
import type { ApiEnvelope } from './http'
import type { PublicPost } from '@/types/post'

export interface PostListResult {
  posts: PublicPost[]
  page: number
  pageSize: number
}

// 后端字段为 snake_case（cover_url/category_id/view_count/created_at/updated_at），归一化为 camelCase
function normalize(raw: Record<string, unknown>): PublicPost {
  return {
    title: String(raw.title ?? ''),
    slug: String(raw.slug ?? ''),
    summary: (raw.summary as string | null | undefined) ?? null,
    content: String(raw.content ?? ''),
    coverUrl: (raw.cover_url as string | null | undefined) ?? null,
    categoryId:
      raw.category_id === null || raw.category_id === undefined ? null : Number(raw.category_id),
    status: Number(raw.status ?? 0),
    viewCount: Number(raw.view_count ?? 0),
    createdAt: String(raw.created_at ?? ''),
    updatedAt: String(raw.updated_at ?? ''),
  }
}

export async function fetchPostList(page: number, pageSize: number): Promise<PostListResult> {
  const { data } = await http.get<ApiEnvelope<Record<string, unknown>[]>>('/posts', {
    params: { page, page_size: pageSize },
  })
  const rawList = Array.isArray(data?.data) ? data.data : []
  return {
    posts: rawList.map((item) => normalize(item)),
    page: Number(data?.page ?? page),
    pageSize: Number(data?.page_size ?? pageSize),
  }
}

export async function fetchPostBySlug(slug: string): Promise<PublicPost> {
  const { data } = await http.get<ApiEnvelope<Record<string, unknown>>>(
    `/post/${encodeURIComponent(slug)}`,
  )
  return normalize(data?.data ?? {})
}
