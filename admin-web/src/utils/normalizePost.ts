import type { Post, PostStatus } from '@/types/post'

/**
 * 后端两个读取接口返回的字段风格不一致：
 * - GET /admin/list  直接序列化 Postdb（无 json tag），字段是大驼峰：ID/Title/CoverURL/CreatedAt...
 * - GET /admin/get/:id 手动拼装 gin.H，字段是 snake_case，且时间字段是 create_at/update_at（少一个 d）
 * 这里统一做兼容归一化，前端只使用 Post 类型，不关心后端具体返回哪一种。
 */
export function normalizePost(raw: Record<string, unknown>): Post {
  const pick = (...keys: string[]): unknown => {
    for (const key of keys) {
      if (raw[key] !== undefined && raw[key] !== null) return raw[key]
    }
    return undefined
  }

  const id = pick('id', 'ID')
  const status = pick('status', 'Status')
  const viewCount = pick('view_count', 'ViewCount')
  const categoryId = pick('category_id', 'CategoryID')

  return {
    id: id === undefined ? '' : String(id),
    title: String(pick('title', 'Title') ?? ''),
    slug: String(pick('slug', 'Slug') ?? ''),
    summary: (pick('summary', 'Summary') as string | null | undefined) ?? null,
    content: String(pick('content', 'Content') ?? ''),
    coverUrl: (pick('cover_url', 'CoverURL') as string | null | undefined) ?? null,
    categoryId: categoryId === undefined ? null : Number(categoryId),
    status: (Number(status ?? 0) as PostStatus) ?? 0,
    viewCount: Number(viewCount ?? 0),
    createdAt: (pick('created_at', 'create_at', 'CreatedAt') as string | undefined) ?? null,
    updatedAt: (pick('updated_at', 'update_at', 'UpdatedAt') as string | undefined) ?? null,
  }
}
