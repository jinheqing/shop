import { createRouter, createWebHistory } from 'vue-router'

// admin 专属路由（advisor 无权访问）
const ADMIN_ONLY = { requiredRole: 'admin' }

const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  { path: '/', component: () => import('@/views/Layout.vue'), redirect: '/dashboard', children: [
    // Dashboard — 所有 staff 可看
    { path: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: 'Dashboard' } },

    // Customers — 所有 staff
    { path: 'users', component: () => import('@/views/Users.vue'), meta: { title: 'Customers' } },
    { path: 'conversations', component: () => import('@/views/Conversations.vue'), meta: { title: 'IM Conversations' } },

    // Products — 所有 staff 可浏览，但只有 admin 能创建/发布（组件内再校验）
    { path: 'custom-products', component: () => import('@/views/CustomProducts.vue'), meta: { title: 'Bespoke Products' } },
    { path: 'custom-products/:id', component: () => import('@/views/CustomProductDetail.vue'), meta: { title: 'Product' } },
    { path: 'qrcodes', component: () => import('@/views/QRCodes.vue'), meta: { title: 'QR Codes' } },

    // Orders + Payments — 所有 staff 可看
    { path: 'orders', component: () => import('@/views/Orders.vue'), meta: { title: 'Orders' } },
    { path: 'orders/:id', component: () => import('@/views/OrderDetail.vue'), meta: { title: 'Order Detail' } },
    { path: 'transactions', component: () => import('@/views/Transactions.vue'), meta: { title: 'Payment Transactions' } },
    { path: 'payment-config', component: () => import('@/views/PaymentConfig.vue'), meta: { title: 'Payment Config', ...ADMIN_ONLY } },
    { path: 'webhooks', component: () => import('@/views/WebhookLogs.vue'), meta: { title: 'Webhook Logs', ...ADMIN_ONLY } },

    // Invoices + Declarations + FX + SGS — 所有 staff 可看
    { path: 'invoices', component: () => import('@/views/Invoices.vue'), meta: { title: 'Invoices' } },
    { path: 'declarations', component: () => import('@/views/Declarations.vue'), meta: { title: 'Declarations' } },
    { path: 'ledgers', component: () => import('@/views/Ledgers.vue'), meta: { title: 'FX Ledgers' } },
    { path: 'sgs', component: () => import('@/views/SGS.vue'), meta: { title: 'SGS Reports' } },

    // Live — 所有 staff
    { path: 'live-rooms', component: () => import('@/views/LiveRooms.vue'), meta: { title: 'Live Rooms' } },
    { path: 'live-calendar', component: () => import('@/views/LiveCalendar.vue'), meta: { title: 'Live Calendar' } },
    { path: 'customer-requests', component: () => import('@/views/CustomerRequests.vue'), meta: { title: 'Customer Live Requests' } },
    { path: 'delivery-inspection', component: () => import('@/views/DeliveryInspection.vue'), meta: { title: 'Delivery Inspection' } },
    { path: 'livekit', component: () => import('@/views/LiveKitRooms.vue'), meta: { title: 'LiveKit Rooms' } },
    { path: 'recordings', component: () => import('@/views/Recordings.vue'), meta: { title: 'Recordings' } },
    { path: 'videos', component: () => import('@/views/Videos.vue'), meta: { title: 'Video Publishing' } },
    { path: 'user-groups', component: () => import('@/views/UserGroups.vue'), meta: { title: 'User Groups' } },
    { path: 'short-links', component: () => import('@/views/ShortLinks.vue'), meta: { title: 'Short Links' } },
    { path: 'slow-presets', component: () => import('@/views/SlowPresets.vue'), meta: { title: 'Slow Presets' } },

    // Infrastructure — admin only
    { path: 'nodes', component: () => import('@/views/Nodes.vue'), meta: { title: 'Media Nodes', ...ADMIN_ONLY } },
    { path: 'translate-status', component: () => import('@/views/TranslateStatus.vue'), meta: { title: 'Translate Engine', ...ADMIN_ONLY } },

    // System — admin only
    { path: 'staff', component: () => import('@/views/Staff.vue'), meta: { title: 'Staff', ...ADMIN_ONLY } },
    { path: 'system-settings', component: () => import('@/views/SystemSettings.vue'), meta: { title: 'System Settings', ...ADMIN_ONLY } },
    { path: 'site-content', component: () => import('@/views/SiteContent.vue'), meta: { title: 'CMS', ...ADMIN_ONLY } },
    { path: 'dsar', component: () => import('@/views/DSAR.vue'), meta: { title: 'GDPR/DSAR', ...ADMIN_ONLY } },
    { path: 'audit-logs', component: () => import('@/views/AuditLogs.vue'), meta: { title: 'Audit Logs', ...ADMIN_ONLY } },
  ] },
]
const r = createRouter({ history: createWebHistory('/admin/'), routes })
r.beforeEach((to, _, next) => {
  if (to.path === '/login') return next()
  if (!localStorage.getItem('staff_token')) return next('/login')
  const role = localStorage.getItem('staff_role') || ''
  const meta: any = to.meta || {}
  if (meta.requiredRole && role !== meta.requiredRole) {
    // advisor 访问 admin 专属路由 → 跳回 dashboard
    return next('/dashboard')
  }
  next()
})
export default r
