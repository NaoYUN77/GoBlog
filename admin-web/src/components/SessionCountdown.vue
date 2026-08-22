<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Timer } from '@lucide/vue'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const auth = useAuthStore()
const router = useRouter()
const toast = useToastStore()

const now = ref(Date.now())
let timer: number | undefined

onMounted(() => {
  timer = window.setInterval(() => {
    now.value = Date.now()
  }, 1000)
})
onBeforeUnmount(() => {
  if (timer) window.clearInterval(timer)
})

const remainingMs = computed(() => {
  if (!auth.expiresAt) return 0
  return Math.max(0, auth.expiresAt - now.value)
})

const remainingLabel = computed(() => {
  const totalSeconds = Math.ceil(remainingMs.value / 1000)
  const m = Math.floor(totalSeconds / 60)
  const s = totalSeconds % 60
  return `${m}:${String(s).padStart(2, '0')}`
})

const urgency = computed<'normal' | 'warn' | 'danger'>(() => {
  if (remainingMs.value <= 30_000) return 'danger'
  if (remainingMs.value <= 120_000) return 'warn'
  return 'normal'
})

let hasExpired = false
watchExpiry()

function watchExpiry() {
  const check = () => {
    if (!hasExpired && auth.isAuthenticated && remainingMs.value <= 0) {
      hasExpired = true
      auth.logout()
      toast.error('登录凭证已过期，请重新登录')
      router.push({ name: 'login' })
    }
  }
  setInterval(check, 1000)
}
</script>

<template>
  <div
    v-if="auth.isAuthenticated"
    class="hidden items-center gap-1.5 rounded-full border px-2.5 py-1 font-mono text-xs tabular sm:inline-flex"
    :class="{
      'border-border text-ink-soft': urgency === 'normal',
      'border-warning/50 bg-warning-soft text-warning': urgency === 'warn',
      'border-danger/50 bg-danger-soft text-danger': urgency === 'danger',
    }"
    :title="'登录凭证将在 ' + remainingLabel + ' 后过期'"
  >
    <Timer :size="13" :stroke-width="2.25" />
    {{ remainingLabel }}
  </div>
</template>
