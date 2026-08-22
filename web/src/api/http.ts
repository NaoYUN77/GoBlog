import axios, { AxiosError } from 'axios'

export interface ApiEnvelope<T = unknown> {
  msg?: string
  data?: T
  page?: number
  page_size?: number
}

const http = axios.create({
  baseURL: '/',
  timeout: 15000,
})

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiEnvelope>) => {
    if (!error.response) {
      return Promise.reject(new Error('无法连接服务器，请稍后重试'))
    }
    return Promise.reject(new Error(error.response.data?.msg || '请求失败，请稍后重试'))
  },
)

export default http
