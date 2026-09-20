<script setup lang="ts">
import { useRouter } from 'vue-router'
import { computed } from 'vue'
import { useAuth } from '@/stores/auth'
const auth = useAuth()
const router = useRouter()

type MenuItem = { label: string; path: string; icon: string; roles?: string[] }
type MenuGroup = { label: string; items: MenuItem[] }

const menuGroups: MenuGroup[] = [
  { label: 'Overview', items: [
    { label: 'Dashboard', path: '/dashboard', icon: 'DataAnalysis' },
  ]},
  { label: 'Customers & IM', items: [
    { label: 'Customers', path: '/users', icon: 'User' },
    { label: 'IM Conversations', path: '/conversations', icon: 'ChatDotRound' },
  ]},
  { label: 'Bespoke Products', items: [
    { label: 'Product List', path: '/custom-products', icon: 'Document' },
    { label: 'Traceability QR Codes', path: '/qrcodes', icon: 'Connection' },
  ]},
  { label: 'Orders & Payments', items: [
    { label: 'Orders', path: '/orders', icon: 'ShoppingCart' },
    { label: 'Payment Transactions', path: '/transactions', icon: 'Wallet' },
    { label: 'Webhook Logs', path: '/webhooks', icon: 'Bell', roles: ['admin', 'supervisor'] },
    { label: 'Payment Gateway Config', path: '/payment-config', icon: 'Setting', roles: ['admin', 'supervisor'] },
  ]},
  { label: 'Finance & Compliance', items: [
    { label: 'Invoices', path: '/invoices', icon: 'Tickets' },
    { label: 'Declarations', path: '/declarations', icon: 'Files' },
    { label: 'FX Ledgers', path: '/ledgers', icon: 'Money' },
    { label: 'SGS Reports', path: '/sgs', icon: 'DocumentChecked' },
  ]},
  { label: 'Live Streaming', items: [
    { label: 'Live Rooms', path: '/live-rooms', icon: 'VideoCamera' },
    { label: 'Schedule Calendar', path: '/live-calendar', icon: 'Calendar' },
    { label: 'Customer Live Requests', path: '/customer-requests', icon: 'EditPen' },
    { label: 'Delivery Inspection', path: '/delivery-inspection', icon: 'Van' },
    { label: 'Slow Presets (24/7)', path: '/slow-presets', icon: 'Clock' },
    { label: 'LiveKit SFU', path: '/livekit', icon: 'Monitor' },
    { label: 'Recordings', path: '/recordings', icon: 'VideoPlay' },
    { label: 'Video Publishing', path: '/videos', icon: 'Film' },
    { label: 'Short Links', path: '/short-links', icon: 'Link' },
    { label: 'User Groups', path: '/user-groups', icon: 'User', roles: ['admin', 'supervisor'] },
  ]},
  { label: 'Infrastructure', items: [
    { label: 'Media Nodes', path: '/nodes', icon: 'MonitorHeart', roles: ['admin', 'supervisor'] },
    { label: 'Translate Engine', path: '/translate-status', icon: 'Monitor', roles: ['admin', 'supervisor'] },
  ]},
  { label: 'System', items: [
    { label: 'System Settings', path: '/system-settings', icon: 'Setting', roles: ['admin', 'supervisor'] },
    { label: 'Staff', path: '/staff', icon: 'UserFilled', roles: ['admin'] },
    { label: 'CMS (Public Site)', path: '/site-content', icon: 'Reading', roles: ['admin', 'supervisor'] },
    { label: 'GDPR / DSAR', path: '/dsar', icon: 'Lock', roles: ['admin', 'supervisor'] },
    { label: 'Audit Logs', path: '/audit-logs', icon: 'Historic', roles: ['admin', 'supervisor'] },
  ]},
]

// 根据当前用户角色过滤菜单
const visibleGroups = computed(() => {
  const role = auth.user?.role || ''
  return menuGroups
    .map(g => ({
      ...g,
      items: g.items.filter(it => !it.roles || it.roles.includes(role))
    }))
    .filter(g => g.items.length > 0)
})
</script>
<template>
  <el-container class="h-screen">
    <el-aside width="240px" class="bg-slate-900 text-slate-100 overflow-auto">
      <div class="p-4 border-b border-slate-800">
        <span class="text-2xl mr-2">🍃</span>
        <span class="font-serif text-base">UK Tea House Admin</span>
      </div>
      <div v-for="g in visibleGroups" :key="g.label" class="mb-2">
        <div class="px-4 pt-3 pb-1 text-[11px] uppercase tracking-widest text-slate-400">{{ g.label }}</div>
        <el-menu :default-active="router.currentRoute.value.path" background-color="#0f172a" text-color="#cbd5e1" active-text-color="#fbbf24" router>
          <el-menu-item v-for="it in g.items" :key="it.path" :index="it.path">
            <el-icon><component :is="it.icon" /></el-icon>
            <span>{{ it.label }}</span>
          </el-menu-item>
        </el-menu>
      </div>
    </el-aside>
    <el-container>
      <el-header class="bg-white border-b flex items-center justify-between py-3 h-auto">
        <div class="text-lg font-medium text-slate-900">{{ router.currentRoute.value.meta.title }}</div>
        <div class="flex items-center gap-4">
          <span class="text-sm text-slate-600">{{ auth.user?.email }}</span>
          <el-tag :type="auth.user?.role==='admin'?'danger':auth.user?.role==='supervisor'?'warning':auth.user?.role==='farmer'?'success':'primary'" size="small">{{ auth.user?.role }}</el-tag>
          <el-button size="small" @click="auth.logout(); router.push('/login')">Logout</el-button>
        </div>
      </el-header>
      <el-main class="bg-slate-50 overflow-auto"><router-view /></el-main>
    </el-container>
  </el-container>
</template>
