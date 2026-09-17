import { createRouter, createWebHistory } from 'vue-router'
const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  { path: '/', component: () => import('@/views/Layout.vue'), redirect: '/dashboard', children: [
    { path: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: 'Dashboard' } },
    { path: 'orders', component: () => import('@/views/Orders.vue'), meta: { title: 'Orders' } },
    { path: 'orders/:id', component: () => import('@/views/OrderDetail.vue'), meta: { title: 'Order Detail' } },
    { path: 'custom-products', component: () => import('@/views/CustomProducts.vue'), meta: { title: 'Bespoke Products' } },
    { path: 'custom-products/:id', component: () => import('@/views/CustomProductDetail.vue'), meta: { title: 'Product' } },
    { path: 'live-rooms', component: () => import('@/views/LiveRooms.vue'), meta: { title: 'Live Rooms' } },
    { path: 'livekit', component: () => import('@/views/LiveKitRooms.vue'), meta: { title: 'LiveKit Rooms' } },
    { path: 'slow-presets', component: () => import('@/views/SlowPresets.vue'), meta: { title: 'Slow Presets' } },
    { path: 'nodes', component: () => import('@/views/Nodes.vue'), meta: { title: 'Media Nodes' } },
    { path: 'sgs', component: () => import('@/views/SGS.vue'), meta: { title: 'SGS Reports' } },
    { path: 'invoices', component: () => import('@/views/Invoices.vue'), meta: { title: 'Invoices' } },
    { path: 'declarations', component: () => import('@/views/Declarations.vue'), meta: { title: 'Declarations' } },
    { path: 'ledgers', component: () => import('@/views/Ledgers.vue'), meta: { title: 'FX Ledgers' } },
    { path: 'conversations', component: () => import('@/views/Conversations.vue'), meta: { title: 'IM Conversations' } },
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
