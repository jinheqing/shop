<script setup lang="ts">
// ============================================================
// ConciergeWidget — 全站浮动顾问入口（Old Money Edition）
// ============================================================
// Ritz-Carlton "Please ask the concierge" aesthetic:
// off-black base + champagne gold accent, serif tiny-caps labels,
// 600ms slow transition.
//
// 信任原则：永远不收集却假装提交。用户点击 → 直接跳 /chat
// 已登录则进入顾问对话；未登录则进入 magic-link 后再回到 /chat
// ============================================================

import { ref } from 'vue'
import { useRouter } from 'vue-router'

const open = ref(false)
const scrolled = ref(false)
const router = useRouter()

window.addEventListener('scroll', () => { scrolled.value = window.scrollY > 400 })

function toggle() { open.value = !open.value }

function goChat() {
  const token = localStorage.getItem('user_token')
  open.value = false
  if (!token) {
    router.push('/magic-link?redirect=/chat')
  } else {
    router.push('/chat')
  }
}

function goContact() {
  open.value = false
  router.push('/contact')
}
</script>

<template>
  <div
    class="fixed right-7 z-[60] font-sans"
    :class="scrolled ? 'bottom-6' : 'bottom-8'">

    <!-- ===== slide-in 顾问联系面板 ===== -->
    <transition name="concierge-panel">
      <div
        v-if="open"
        class="absolute bottom-16 right-0 w-[340px] md:w-[380px] bg-ink-900 text-ivory-100 shadow-2xl overflow-hidden"
        style="border-radius: 2px;">

        <div class="h-px bg-gold/40"></div>

        <div class="px-8 py-9">

          <div class="flex items-start justify-between mb-7">
            <div>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">At Your Service</span>
              <h3 class="font-serif text-2xl text-ivory-100 mt-2 leading-tight">
                Speak with<br> a Tea Advisor
              </h3>
            </div>
            <button
              @click="open = false"
              class="text-sand hover:text-gold transition text-xl leading-none p-1"
              aria-label="Close">&times;</button>
          </div>

          <p class="font-serif text-sm text-sand leading-relaxed mb-7">
            One advisor oversees your enquiry from first conversation to
            final brew. No call centres. No bots. A named person, by email
            or private message — Monday through Saturday.
          </p>

          <div class="h-px bg-gold/20 my-6"></div>

          <div class="flex flex-col gap-2">
            <button
              @click="goChat"
              class="w-full py-3 text-[11px] uppercase tracking-lux text-ink-900 bg-ivory-100 hover:bg-gold transition font-sans">
              Open · Private · Conversation
            </button>
            <button
              @click="goContact"
              class="w-full py-3 text-[11px] uppercase tracking-lux text-gold border border-gold/40 hover:border-gold hover:bg-gold/5 transition font-sans">
              Other · Enquiries
            </button>
          </div>

          <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-7 font-sans leading-loose">
            Mon&ndash;Sat &middot; 09:00 &ndash; 18:00 GMT<br>
            hello@ukteahouse.co.uk &middot; +44 20 7946 0958
          </p>
        </div>
      </div>
    </transition>

    <!-- ===== 触发按钮 ===== -->
    <button
      @click="toggle"
      class="group relative w-14 h-14 rounded-full bg-ink-900 hover:bg-ink-800 border border-gold/30 hover:border-gold transition-all duration-lux shadow-xl flex items-center justify-center"
      :aria-expanded="open"
      aria-label="Talk to a tea advisor">

      <svg width="18" height="18" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"
           class="text-gold group-hover:text-gold-soft transition">
        <path d="M6 8a6 6 0 1 1 12 0c0 7 3 9 3 9H3s3-2 3-9"></path>
        <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"></path>
      </svg>

      <span
        class="absolute right-full mr-3 whitespace-nowrap text-[10px] uppercase tracking-lux text-gold font-sans opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-duration-lux pointer-events-none">
        Speak with an Advisor
      </span>

      <span class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 bg-gold rounded-full border-2 border-ink-900"></span>
    </button>
  </div>
</template>

<style scoped>
.concierge-panel-enter-active,
.concierge-panel-leave-active {
  transition: opacity 300ms ease, transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
}
.concierge-panel-enter-from,
.concierge-panel-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.98);
}
</style>
