import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  // ⚠️ 关键: admin-dashboard 部署在 /admin/ 子路径下
  // 没有这个，打包后所有资源路径都是 /assets/xxx.js → 404
  base: '/admin/',
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: { port: 5174, proxy: { '/api': { target: 'http://localhost:8080', changeOrigin: true } } }
})
