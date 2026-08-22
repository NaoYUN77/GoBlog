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
    port: 5173,
    proxy: {
      // 后端 Gin 服务默认监听 8080，未配置 CORS，开发环境走代理避免跨域
      '/admin': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
})
