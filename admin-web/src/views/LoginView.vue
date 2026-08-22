<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Eye, EyeOff, Loader2, SquarePen } from '@lucide/vue'
import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const toast = useToastStore()

const form = reactive({ username: '', password: '' })
const showPassword = ref(false)
const submitting = ref(false)
const errorMessage = ref('')

async function handleSubmit() {
  if (!form.username.trim() || !form.password) {
    errorMessage.value = '请输入用户名和密码'
    return
  }
  errorMessage.value = ''
  submitting.value = true
  try {
    const { token } = await login(form.username.trim(), form.password)
    auth.setSession(token, form.username.trim())
    toast.success('登录成功')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/posts'
    router.push(redirect)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : '登录失败'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-paper px-4 py-12">
    <div class="w-full max-w-sm">
      <div class="mb-8 flex flex-col items-center gap-3 text-center">
        <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-ink text-paper shadow-card">
          <SquarePen :size="22" :stroke-width="2.2" />
        </div>
        <div>
          <h1 class="font-serif text-2xl font-semibold tracking-tight text-ink">笔耕</h1>
          <p class="mt-1 text-sm text-ink-soft">博客管理后台</p>
        </div>
      </div>

      <form
        class="rounded-2xl border border-border bg-surface p-6 shadow-card"
        novalidate
        @submit.prevent="handleSubmit"
      >
        <div class="space-y-4">
          <div>
            <label for="username" class="mb-1.5 block text-sm font-medium text-ink">用户名</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              autocomplete="username"
              class="w-full rounded-lg border border-border bg-paper px-3.5 py-2.5 text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
              placeholder="输入管理员用户名"
              :disabled="submitting"
            />
          </div>

          <div>
            <label for="password" class="mb-1.5 block text-sm font-medium text-ink">密码</label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                class="w-full rounded-lg border border-border bg-paper px-3.5 py-2.5 pr-10 text-sm text-ink placeholder:text-ink-faint focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20"
                placeholder="输入登录密码"
                :disabled="submitting"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 flex w-10 items-center justify-center text-ink-faint hover:text-ink-soft"
                :aria-label="showPassword ? '隐藏密码' : '显示密码'"
                @click="showPassword = !showPassword"
              >
                <EyeOff v-if="showPassword" :size="17" />
                <Eye v-else :size="17" />
              </button>
            </div>
          </div>

          <p v-if="errorMessage" class="text-sm text-danger" role="alert">{{ errorMessage }}</p>

          <button
            type="submit"
            class="flex w-full items-center justify-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-primary-hover disabled:cursor-not-allowed disabled:opacity-70"
            :disabled="submitting"
          >
            <Loader2 v-if="submitting" :size="16" class="animate-spin" />
            {{ submitting ? '登录中…' : '登录' }}
          </button>
        </div>
      </form>

      <p class="mt-5 text-center text-xs text-ink-faint">登录凭证有效期 10 分钟，过期后需重新登录</p>
    </div>
  </div>
</template>
