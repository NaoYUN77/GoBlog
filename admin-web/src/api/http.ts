import axios, { AxiosError } from 'axios'
import JSONbig from 'json-bigint'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import router from '@/router'

// 文章 id 由雪花算法生成，超过 JS Number 安全整数范围，
// 用 json-bigint 解析响应，避免大整数精度丢失导致 id 错误。
const parseBig = JSONbig({ storeAsString: true })

export interface ApiEnvelope<T = unknown> {
  code?: number
  msg?: string
  data?: T
  page?: number
  page_size?: number
}

const http = axios.create({
  baseURL: '/admin',
  timeout: 15000,
  transformResponse: [
    (data: string) => {
      if (typeof data !== 'string' || data.length === 0) return data
      try {
        return parseBig.parse(data)
      } catch {
        return data
      }
    },
  ],
})

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.set('Authorization', `Bearer ${auth.token}`)
  }
  return config
})

const TOKEN_ERROR_HINTS = ['token', '请求头', '无法解析']

function looksLikeAuthError(msg: string | undefined): boolean {
  if (!msg) return false
  return TOKEN_ERROR_HINTS.some((hint) => msg.includes(hint))
}

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiEnvelope>) => {
    const status = error.response?.status
    const msg = error.response?.data?.msg

    const isAuthFailure = status === 401 || (status === 400 && looksLikeAuthError(msg))

    if (isAuthFailure) {
      const auth = useAuthStore()
      const wasAuthenticated = auth.isAuthenticated
      auth.logout()
      if (wasAuthenticated) {
        useToastStore().error('登录状态已失效，请重新登录')
        router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
      }
      return Promise.reject(new Error(msg || '登录状态已失效'))
    }

    if (!error.response) {
      return Promise.reject(new Error('无法连接服务器，请检查网络或后端服务'))
    }

    return Promise.reject(new Error(msg || '请求失败，请稍后重试'))
  },
)

export default http
