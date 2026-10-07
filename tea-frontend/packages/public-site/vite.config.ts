import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      // IM WebSocket（后端注册在根路由 /ws/im）
      '/ws': { target: 'http://localhost:8080', ws: true, changeOrigin: true },
      // 上传文件静态服务
      '/uploads': { target: 'http://localhost:8080', changeOrigin: true }
    }
  }
})
