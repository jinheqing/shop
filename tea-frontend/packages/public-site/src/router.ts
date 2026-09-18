import { createRouter, createWebHistory } from 'vue-router'

export const routes = [
  // Core public pages
  { path: '/', name: 'home', component: () => import('@/views/Home.vue'), meta: { title: 'UK Tea House' } },
  { path: '/tea-gardens', name: 'tea-gardens', component: () => import('@/views/TeaGardens.vue'), meta: { title: 'Tea Garden Stories' } },
  { path: '/live', name: 'live', component: () => import('@/views/LiveRoom.vue'), meta: { title: 'Live Tea Cameras' } },
  { path: '/quality', name: 'quality', component: () => import('@/views/Quality.vue'), meta: { title: 'SGS Quality Assurance' } },
  { path: '/bespoke', name: 'bespoke', component: () => import('@/views/Bespoke.vue'), meta: { title: 'Bespoke Blending' } },

  // Quote (shareable) — design doc: /quote/{token}
  { path: '/quote/:token', name: 'quote', component: () => import('@/views/QuotePage.vue'), props: true, meta: { title: 'Bespoke Quote' } },
  // Legacy bespoke-detail, redirect via alias below
  { path: '/bespoke/:token', name: 'bespoke-detail', component: () => import('@/views/QuotePage.vue'), props: true, meta: { title: 'Bespoke Quote' } },

  // Checkout + Payment Result
  { path: '/checkout/:token', name: 'checkout', component: () => import('@/views/Checkout.vue'), props: true },
  { path: '/order/result', name: 'payment-result', component: () => import('@/views/PaymentResult.vue'), meta: { title: 'Payment Result' } },

  // Traceability QR scan page
  { path: '/trace/:token', name: 'trace', component: () => import('@/views/Traceability.vue'), props: true, meta: { title: 'Trace Your Tea' } },

  // User auth + account
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue'), meta: { title: 'Sign In' } },
  { path: '/magic-link', name: 'magic-link', component: () => import('@/views/MagicLink.vue'), meta: { title: 'Magic Link Login' } },
  { path: '/account', name: 'account', component: () => import('@/views/Account.vue'), meta: { title: 'My Account' } },

  // IM chat
  { path: '/chat', name: 'chat', component: () => import('@/views/Chat.vue'), meta: { title: 'Chat with Advisor' } },

  // Company / legal
  { path: '/about', name: 'about', component: () => import('@/views/About.vue'), meta: { title: 'About Us' } },
  { path: '/contact', name: 'contact', component: () => import('@/views/Contact.vue'), meta: { title: 'Contact' } },
  { path: '/faq', name: 'faq', component: () => import('@/views/FAQ.vue'), meta: { title: 'FAQ' } },
  { path: '/privacy', name: 'privacy', component: () => import('@/views/Privacy.vue'), meta: { title: 'Privacy · GDPR' } },
  { path: '/cookie-policy', name: 'cookie-policy', component: () => import('@/views/CookiePolicy.vue'), meta: { title: 'Cookie Policy' } },
  { path: '/terms', name: 'terms', component: () => import('@/views/Terms.vue'), meta: { title: 'Terms of Service' } },
  { path: '/gdpr-dsar', name: 'gdpr-dsar', component: () => import('@/views/GdprDsar.vue'), meta: { title: 'GDPR · Data Requests' } },

  // Invoice
  { path: '/orders/:id/invoice', name: 'invoice', component: () => import('@/views/InvoiceView.vue'), props: true },
]

const router = createRouter({ history: createWebHistory(), routes, scrollBehavior: () => ({ top: 0 }) })
router.afterEach(to => { if (to.meta.title) document.title = to.meta.title as string + ' · UK Tea House' })
export default router
