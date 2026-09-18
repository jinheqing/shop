import { createRouter, createWebHistory } from 'vue-router'
const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  { path: '/', component: () => import('@/views/Layout.vue'), redirect: '/dashboard', children: [
    // Dashboard
    { path: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: 'Dashboard' } },

    // Customers
    { path: 'users', component: () => import('@/views/Users.vue'), meta: { title: 'Customers' } },
    { path: 'conversations', component: () => import('@/views/Conversations.vue'), meta: { title: 'IM Conversations' } },

    // Products
    { path: 'custom-products', component: () => import('@/views/CustomProducts.vue'), meta: { title: 'Bespoke Products' } },
    { path: 'custom-products/:id', component: () => import('@/views/CustomProductDetail.vue'), meta: { title: 'Product' } },
    { path: 'qrcodes', component: () => import('@/views/QRCodes.vue'), meta: { title: 'QR Codes' } },

    // Orders + Payments
    { path: 'orders', component: () => import('@/views/Orders.vue'), meta: { title: 'Orders' } },
    { path: 'orders/:id', component: () => import('@/views/OrderDetail.vue'), meta: { title: 'Order Detail' } },
    { path: 'transactions', component: () => import('@/views/Transactions.vue'), meta: { title: 'Payment Transactions' } },
    { path: 'payment-config', component: () => import('@/views/PaymentConfig.vue'), meta: { title: 'Payment Config' } },
    { path: 'webhooks', component: () => import('@/views/WebhookLogs.vue'), meta: { title: 'Webhook Logs' } },

    // Invoices + Declarations + FX + SGS
    { path: 'invoices', component: () => import('@/views/Invoices.vue'), meta: { title: 'Invoices' } },
    { path: 'declarations', component: () => import('@/views/Declarations.vue'), meta: { title: 'Declarations' } },
    { path: 'ledgers', component: () => import('@/views/Ledgers.vue'), meta: { title: 'FX Ledgers' } },
    { path: 'sgs', component: () => import('@/views/SGS.vue'), meta: { title: 'SGS Reports' } },

    // Live
    { path: 'live-rooms', component: () => import('@/views/LiveRooms.vue'), meta: { title: 'Live Rooms' } },
    { path: 'live-calendar', component: () => import('@/views/LiveCalendar.vue'), meta: { title: 'Live Calendar' } },
    { path: 'customer-requests', component: () => import('@/views/CustomerRequests.vue'), meta: { title: 'Customer Live Requests' } },
    { path: 'delivery-inspection', component: () => import('@/views/DeliveryInspection.vue'), meta: { title: 'Delivery Inspection' } },
    { path: 'livekit', component: () => import('@/views/LiveKitRooms.vue'), meta: { title: 'LiveKit Rooms' } },
    { path: 'recordings', component: () => import('@/views/Recordings.vue'), meta: { title: 'Recordings' } },
    { path: 'user-groups', component: () => import('@/views/UserGroups.vue'), meta: { title: 'User Groups' } },
    { path: 'short-links', component: () => import('@/views/ShortLinks.vue'), meta: { title: 'Short Links' } },
    { path: 'slow-presets', component: () => import('@/views/SlowPresets.vue'), meta: { title: 'Slow Presets' } },

    // Infrastructure
    { path: 'nodes', component: () => import('@/views/Nodes.vue'), meta: { title: 'Media Nodes' } },
    { path: 'translate-status', component: () => import('@/views/TranslateStatus.vue'), meta: { title: 'Translate Engine' } },

    // System
    { path: 'staff', component: () => import('@/views/Staff.vue'), meta: { title: 'Staff' } },
    { path: 'site-content', component: () => import('@/views/SiteContent.vue'), meta: { title: 'CMS' } },
    { path: 'dsar', component: () => import('@/views/DSAR.vue'), meta: { title: 'GDPR/DSAR' } },
    { path: 'audit-logs', component: () => import('@/views/AuditLogs.vue'), meta: { title: 'Audit Logs' } },
  ] },
]
const r = createRouter({ history: createWebHistory(), routes })
r.beforeEach((to, _, next) => {
  if (to.path === '/login') return next()
  if (!localStorage.getItem('staff_token')) return next('/login')
  next()
})
export default r
