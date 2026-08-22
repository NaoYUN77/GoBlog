import { defineStore } from 'pinia'

// 与后端 modles.TokenExpireDuration 保持一致（10 分钟），仅用于前端倒计时提醒，
// 真正的过期判定始终以后端接口返回的 401/token 错误为准。
export const TOKEN_TTL_MS = 10 * 60 * 1000

const STORAGE_KEY = 'blog-admin-auth'

interface AuthState {
  token: string | null
  username: string | null
  loginAt: number | null
}

function loadInitial(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return { token: null, username: null, loginAt: null }
    const parsed = JSON.parse(raw) as AuthState
    return {
      token: parsed.token ?? null,
      username: parsed.username ?? null,
      loginAt: parsed.loginAt ?? null,
    }
  } catch {
    return { token: null, username: null, loginAt: null }
  }
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => loadInitial(),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    expiresAt: (state) => (state.loginAt ? state.loginAt + TOKEN_TTL_MS : null),
  },
  actions: {
    setSession(token: string, username: string) {
      this.token = token
      this.username = username
      this.loginAt = Date.now()
      this.persist()
    },
    logout() {
      this.token = null
      this.username = null
      this.loginAt = null
      localStorage.removeItem(STORAGE_KEY)
    },
    persist() {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({ token: this.token, username: this.username, loginAt: this.loginAt }),
      )
    },
  },
})
