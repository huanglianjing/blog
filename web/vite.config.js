import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 开发环境下将接口请求转发到后端，避免跨域。
      // 只代理具体接口路径，避免与前端 SPA 路由 /article、/article/:title 冲突。
      '/article/list': 'http://localhost:6000',
      '/article/detail': 'http://localhost:6000',
      '/category/overview': 'http://localhost:6000',
      '/category/list': 'http://localhost:6000',
      '/tag/overview': 'http://localhost:6000',
      '/tag/list': 'http://localhost:6000',
    },
  },
})
