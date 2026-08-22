import http from './http'
import { normalizePost } from '@/utils/normalizePost'
import type { CreatePostPayload, Post, UpdatePostPayload } from '@/types/post'

export interface PostListResult {
  posts: Post[]
  page: number
  pageSize: number
}

export async function fetchPostList(page: number, pageSize: number): Promise<PostListResult> {
  const { data } = await http.get('/list', { params: { page, page_size: pageSize } })
  const rawList = Array.isArray(data?.data) ? data.data : []
  return {
    posts: rawList.map((item: Record<string, unknown>) => normalizePost(item)),
    page: Number(data?.page ?? page),
    pageSize: Number(data?.page_size ?? pageSize),
  }
}

export async function fetchPost(id: string): Promise<Post> {
  const { data } = await http.get(`/get/${id}`)
  return normalizePost(data?.data ?? {})
}

export async function createPost(payload: CreatePostPayload): Promise<void> {
  await http.post('/post', payload)
}

export async function updatePost(id: string, payload: UpdatePostPayload): Promise<void> {
  await http.patch(`/patch/${id}`, payload)
}

export async function deletePost(id: string): Promise<void> {
  await http.delete(`/delete/${id}`)
}
