// 用户端文章：后端返回统一 snake_case 字段，且不包含内部 id
export interface PublicPost {
  title: string
  slug: string
  summary: string | null
  content: string
  coverUrl: string | null
  categoryId: number | null
  status: number
  viewCount: number
  createdAt: string
  updatedAt: string
}
