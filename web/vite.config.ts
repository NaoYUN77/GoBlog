import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5174,
    proxy: {
      // 用户端公开接口：/posts 列表、/post/:slug 详情（Gin 监听 8080，开发环境走代理）
      '/post': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
})
