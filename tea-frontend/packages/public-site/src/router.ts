import { createRouter, createWebHistory } from 'vue-router'

export const routes = [
  { path: '/', name: 'home', component: () => import('@/views/Home.vue'), meta: { title: 'UK Tea House — Pu\'er From Cloud Mountains' } },
  { path: '/mountains', name: 'mountains', component: () => import('@/views/Mountains.vue'), meta: { title: 'Tea Mountains' } },
  { path: '/bespoke', name: 'bespoke', component: () => import('@/views/Bespoke.vue'), meta: { title: 'Bespoke Blending' } },
  { path: '/bespoke/:token', name: 'bespoke-detail', component: () => import('@/views/BespokeDetail.vue'), props: true, meta: { title: 'Bespoke Product' } },
  { path: '/live', name: 'live', component: () => import('@/views/Live.vue'), meta: { title: 'Live Tea Cameras' } },
  { path: '/quality', name: 'quality', component: () => import('@/views/Quality.vue'), meta: { title: 'SGS Quality Assurance' } },
  { path: '/checkout/:token', name: 'checkout', component: () => import('@/views/Checkout.vue'), props: true, meta: { title: 'Secure Checkout' } },
]

const router = createRouter({ history: createWebHistory(), routes, scrollBehavior: () => ({ top: 0 }) })
router.afterEach(to => { if (to.meta.title) document.title = to.meta.title as string })
export default router
