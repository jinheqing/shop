<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { ref } from 'vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import enLocale from 'element-plus/es/locale/lang/en'
import { setLocale } from '@/i18n'

const auth = useAuth()
const router = useRouter()
const { t, locale } = useI18n()
const epLocale = ref(locale.value === 'zh' ? zhCn : enLocale)

type MenuItem = { key: string; path: string; icon: string }
type MenuGroup = { key: string; items: MenuItem[] }

const menuGroups: MenuGroup[] = [
  { key: 'menu.overview', items: [
    { key: 'menu.dashboard', path: '/dashboard', icon: 'DataAnalysis' },
  ]},
  { key: 'menu.customers_im', items: [
    { key: 'menu.customers', path: '/users', icon: 'User' },
    { key: 'menu.conversations', path: '/conversations', icon: 'ChatDotRound' },
  ]},
  { key: 'menu.bespoke_products', items: [
    { key: 'menu.product_list', path: '/custom-products', icon: 'Document' },
    { key: 'menu.qr_codes', path: '/qrcodes', icon: 'Connection' },
  ]},
  { key: 'menu.orders_payments', items: [
    { key: 'menu.orders', path: '/orders', icon: 'ShoppingCart' },
    { key: 'menu.transactions', path: '/transactions', icon: 'Wallet' },
    { key: 'menu.webhooks', path: '/webhooks', icon: 'Bell' },
    { key: 'menu.payment_config', path: '/payment-config', icon: 'Setting' },
  ]},
  { key: 'menu.finance_compliance', items: [
    { key: 'menu.invoices', path: '/invoices', icon: 'Tickets' },
    { key: 'menu.declarations', path: '/declarations', icon: 'Files' },
    { key: 'menu.ledgers', path: '/ledgers', icon: 'Money' },
    { key: 'menu.sgs_reports', path: '/sgs', icon: 'DocumentChecked' },
  ]},
  { key: 'menu.live_streaming', items: [
    { key: 'menu.live_rooms', path: '/live-rooms', icon: 'VideoCamera' },
    { key: 'menu.schedule_calendar', path: '/live-calendar', icon: 'Calendar' },
    { key: 'menu.customer_requests', path: '/customer-requests', icon: 'EditPen' },
    { key: 'menu.delivery_inspection', path: '/delivery-inspection', icon: 'Van' },
    { key: 'menu.slow_presets', path: '/slow-presets', icon: 'Clock' },
    { key: 'menu.livekit', path: '/livekit', icon: 'Monitor' },
    { key: 'menu.recordings', path: '/recordings', icon: 'VideoPlay' },
    { key: 'menu.short_links', path: '/short-links', icon: 'Link' },
    { key: 'menu.user_groups', path: '/user-groups', icon: 'User' },
    { key: 'menu.referrals', path: '/referrals', icon: 'Share' },
  ]},
  { key: 'menu.infrastructure', items: [
    { key: 'menu.media_nodes', path: '/nodes', icon: 'MonitorHeart' },
    { key: 'menu.translate_engine', path: '/translate-status', icon: 'Monitor' },
  ]},
  { key: 'menu.system', items: [
    { key: 'menu.system_settings', path: '/system-settings', icon: 'Setting' },
    { key: 'menu.staff', path: '/staff', icon: 'UserFilled' },
    { key: 'menu.cms', path: '/site-content', icon: 'Reading' },
    { key: 'menu.gdpr_dsar', path: '/dsar', icon: 'Lock' },
    { key: 'menu.audit_logs', path: '/audit-logs', icon: 'Historic' },
  ]},
]

// 路径 → 标题翻译 key 映射，用于在顶栏替代 router.meta.title
const pathTitleMap: Record<string, string> = {
  '/dashboard': 'menu.dashboard',
  '/users': 'menu.customers',
  '/conversations': 'menu.conversations',
  '/custom-products': 'menu.product_list',
  '/qrcodes': 'menu.qr_codes',
  '/orders': 'menu.orders',
  '/transactions': 'menu.transactions',
  '/payment-config': 'menu.payment_config',
  '/webhooks': 'menu.webhooks',
  '/invoices': 'menu.invoices',
  '/declarations': 'menu.declarations',
  '/ledgers': 'menu.ledgers',
  '/sgs': 'menu.sgs_reports',
  '/live-rooms': 'menu.live_rooms',
  '/live-calendar': 'menu.schedule_calendar',
  '/customer-requests': 'menu.customer_requests',
  '/delivery-inspection': 'menu.delivery_inspection',
  '/slow-presets': 'menu.slow_presets',
  '/livekit': 'menu.livekit',
  '/recordings': 'menu.recordings',
  '/short-links': 'menu.short_links',
  '/user-groups': 'menu.user_groups',
  '/referrals': 'menu.referrals',
  '/nodes': 'menu.media_nodes',
  '/translate-status': 'menu.translate_engine',
  '/staff': 'menu.staff',
  '/system-settings': 'menu.system_settings',
  '/site-content': 'menu.cms',
  '/dsar': 'menu.gdpr_dsar',
  '/audit-logs': 'menu.audit_logs',
}

function pageTitle(): string {
  const path = router.currentRoute.value.path
  const direct = pathTitleMap[path]
  if (direct) return t(direct)
  // 子路由（如 /orders/:id）回退到父级路径对应的 key
  for (const prefix of Object.keys(pathTitleMap)) {
    if (path.startsWith(prefix + '/')) {
      return t(pathTitleMap[prefix])
    }
  }
  const meta = router.currentRoute.value.meta.title
  return meta ? String(meta) : ''
}

function switchLocale(lang: string) {
  setLocale(lang)
  epLocale.value = lang === 'zh' ? zhCn : enLocale
}
</script>
<template>
  <el-config-provider :locale="epLocale">
    <el-container class="h-screen">
      <el-aside width="240px" class="bg-slate-900 text-slate-100 overflow-auto">
        <div class="p-4 border-b border-slate-800">
          <span class="text-2xl mr-2">🍃</span>
          <span class="font-serif text-base">{{ t('layout.brand') }}</span>
        </div>
        <div v-for="g in menuGroups" :key="g.key" class="mb-2">
          <div class="px-4 pt-3 pb-1 text-[11px] uppercase tracking-widest text-slate-400">{{ t(g.key) }}</div>
          <el-menu :default-active="router.currentRoute.value.path" background-color="#0f172a" text-color="#cbd5e1" active-text-color="#fbbf24" router>
            <el-menu-item v-for="it in g.items" :key="it.path" :index="it.path">
              <el-icon><component :is="it.icon" /></el-icon>
              <span>{{ t(it.key) }}</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-aside>
      <el-container>
        <el-header class="bg-white border-b flex items-center justify-between py-3 h-auto">
          <div class="text-lg font-medium text-slate-900">{{ pageTitle() }}</div>
          <div class="flex items-center gap-4">
            <el-dropdown trigger="click" @command="switchLocale">
              <el-button size="small" text>
                <span class="text-sm">{{ locale === 'zh' ? '中文' : 'English' }}</span>
                <el-icon class="ml-1"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="zh" :disabled="locale === 'zh'">中文</el-dropdown-item>
                  <el-dropdown-item command="en" :disabled="locale === 'en'">English</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <span class="text-sm text-slate-600">{{ auth.user?.email }}</span>
            <el-tag :type="auth.user?.role==='admin'?'danger':auth.user?.role==='supervisor'?'warning':auth.user?.role==='farmer'?'success':'primary'" size="small">{{ auth.user?.role }}</el-tag>
            <el-button size="small" @click="auth.logout(); router.push('/login')">{{ t('layout.logout') }}</el-button>
          </div>
        </el-header>
        <el-main class="bg-slate-50 overflow-auto"><router-view /></el-main>
      </el-container>
    </el-container>
  </el-config-provider>
</template>
