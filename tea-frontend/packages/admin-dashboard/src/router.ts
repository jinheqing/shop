import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  { path: '/', component: () => import('@/views/Layout.vue'), redirect: '/dashboard', children: [
    { path: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: 'Dashboard' } },
    { path: 'custom-products', component: () => import('@/views/CustomProducts.vue'), meta: { title: 'Bespoke Products' } },
    { path: 'orders', component: () => import('@/views/Orders.vue'), meta: { title: 'Orders' } },
    { path: 'live-rooms', component: () => import('@/views/LiveRooms.vue'), meta: { title: 'Live Rooms' } },
    { path: 'nodes', component: () => import('@/views/Nodes.vue'), meta: { title: 'Media Nodes' } },
    { path: 'invoices', component: () => import('@/views/Invoices.vue'), meta: { title: 'Invoices' } },
    { path: 'sgs', component: () => import('@/views/SGS.vue'), meta: { title: 'SGS Reports' } },
  ] },
]

const r = createRouter({ history: createWebHistory(), routes })
r.beforeEach((to, _, next) => {
  if (to.path === '/login') return next()
  if (!localStorage.getItem('staff_token')) return next('/login')
  next()
})
export default r
