<script setup lang="ts">
import { provide, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppSidebar from '@/components/AppSidebar.vue'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { MOBILE_MENU_KEY } from '@/utils/injectionKeys'

const mobileOpen = ref(false)
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()

provide(MOBILE_MENU_KEY, () => {
  mobileOpen.value = true
})

function handleLogout() {
  auth.logout()
  toast.info('已退出登录')
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-paper">
    <AppSidebar :mobile-open="mobileOpen" @close="mobileOpen = false" @logout="handleLogout" />

    <div class="flex min-w-0 flex-1 flex-col overflow-y-auto">
      <router-view />
    </div>
  </div>
</template>
