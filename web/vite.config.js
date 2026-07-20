import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 开发环境下将接口请求转发到后端，避免跨域。
      // 接口统一以 /api 前缀，避免与前端 SPA 路由 /article、/article/:title 冲突。
      '/api': 'http://localhost:6000',
    },
  },
})
