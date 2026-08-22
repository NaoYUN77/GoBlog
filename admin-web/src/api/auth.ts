import http from './http'

export interface LoginResult {
  token: string
}

export async function login(username: string, password: string): Promise<LoginResult> {
  const { data } = await http.post('/login', { username, password })
  if (!data?.token) {
    throw new Error(data?.msg || '登录失败')
  }
  return { token: data.token as string }
}
