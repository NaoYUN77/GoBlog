import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore, TOKEN_TTL_MS } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      redirect: { name: 'posts' },
      children: [
        {
          path: 'posts',
          name: 'posts',
          component: () => import('@/views/PostListView.vue'),
        },
        {
          path: 'posts/new',
          name: 'post-new',
          component: () => import('@/views/PostEditorView.vue'),
        },
        {
          path: 'posts/:id/edit',
          name: 'post-edit',
          component: () => import('@/views/PostEditorView.vue'),
          props: true,
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { public: true },
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  // 登录已超过后端 token 有效期（10 分钟），本地先行判定为过期，避免带着失效 token 进入页面
  if (auth.isAuthenticated && auth.loginAt && Date.now() - auth.loginAt > TOKEN_TTL_MS) {
    auth.logout()
  }

  if (!to.meta.public && !auth.isAuthenticated) {
    return { name: 'login', query: to.fullPath !== '/' ? { redirect: to.fullPath } : undefined }
  }

  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'posts' }
  }

  return true
})

export default router
