<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
const auth = useAuth()
const router = useRouter()
const menus: any[] = [
  { label: 'Dashboard', path: '/dashboard', icon: 'DataAnalysis' },
  { label: 'Orders', path: '/orders', icon: 'ShoppingCart' },
  { label: 'Bespoke Products', path: '/custom-products', icon: 'Document' },
  { label: 'Live Rooms', path: '/live-rooms', icon: 'VideoCamera' },
  { label: 'LiveKit Rooms', path: '/livekit', icon: 'Monitor' },
  { label: 'Slow Presets', path: '/slow-presets', icon: 'Clock' },
  { label: 'Media Nodes', path: '/nodes', icon: 'MonitorHeart' },
  { label: 'SGS Reports', path: '/sgs', icon: 'DocumentChecked' },
  { label: 'Invoices', path: '/invoices', icon: 'Tickets' },
  { label: 'Declarations', path: '/declarations', icon: 'Tickets' },
  { label: 'FX Ledgers', path: '/ledgers', icon: 'Money' },
  { label: 'IM Conversations', path: '/conversations', icon: 'ChatDotRound' },
  { label: 'Audit Logs', path: '/audit-logs', icon: 'Historic' },
]
</script>
<template>
  <el-container class="h-screen">
    <el-aside width="220px" class="bg-slate-900 text-slate-100 overflow-auto">
      <div class="p-5 border-b border-slate-800">
        <span class="text-2xl mr-2">🍃</span>
        <span class="font-serif text-lg">UK Tea House</span>
      </div>
      <el-menu :default-active="router.currentRoute.value.path" background-color="#0f172a" text-color="#cbd5e1" active-text-color="#fbbf24" router>
        <el-menu-item v-for="m in menus" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <span>{{ m.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="bg-white border-b flex items-center justify-between">
        <div class="text-lg font-medium text-slate-900">{{ router.currentRoute.value.meta.title }}</div>
        <div class="flex items-center gap-4">
          <span class="text-sm text-slate-600">{{ auth.user?.email }}</span>
          <el-tag :type="auth.user?.role==='admin'?'danger':auth.user?.role==='farmer'?'success':'primary'" size="small">{{ auth.user?.role }}</el-tag>
          <el-button size="small" @click="auth.logout(); router.push('/login')">Logout</el-button>
        </div>
      </el-header>
      <el-main class="bg-slate-50 overflow-auto"><router-view /></el-main>
    </el-container>
  </el-container>
</template>
