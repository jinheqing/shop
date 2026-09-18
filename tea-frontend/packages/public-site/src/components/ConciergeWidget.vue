<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

// ============================================================
// ConciergeWidget — 全站浮动顾问入口
// ============================================================
// 任何页面右下角可见。Ritz-Carlton "Please ask bell" /
// Rolls-Royce Private Office aesthetic: off-black base +
// champagne gold accent, serif tiny-caps labels, 600ms slow transition.
// ============================================================

const open = ref(false)
const scrolled = ref(false)
const form = ref({ name: '', email: '', message: '' })
const router = useRouter()

window.addEventListener('scroll', () => { scrolled.value = window.scrollY > 400 })

function goChat() { open.value = false; router.push('/chat') }

async function submit() {
  const key = 'concierge_inquiries'
  const arr = JSON.parse(localStorage.getItem(key) || '[]')
  arr.unshift({
    ...form.value,
    contact_preference: 'email',
    created_at: new Date().toISOString(),
    source: 'public_site_concierge'
  })
  localStorage.setItem(key, JSON.stringify(arr.slice(0, 50)))
  form.value = { name: '', email: '', message: '' }
  open.value = false
  alert('Thank you. An advisor will contact you within 24 hours.')
}
</script>

<template>
  <div
    class="fixed right-7 z-[60] font-sans"
    :class="scrolled ? 'bottom-6' : 'bottom-8'">

    <!-- ===== slide-in 预约表单 ===== -->
    <transition name="concierge-panel">
      <div
        v-if="open"
        class="absolute bottom-16 right-0 w-[360px] md:w-[400px] bg-ink-900 text-ivory-100 shadow-2xl overflow-hidden"
        style="border-radius: 2px;">

        <div class="h-px bg-gold/40"></div>

        <div class="px-8 py-8">

          <div class="flex items-start justify-between mb-6">
            <div>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">At Your Service</span>
              <h3 class="font-serif text-2xl text-ivory-100 mt-2 leading-tight">
                Request a<br> Private Consultation
              </h3>
            </div>
            <button
              @click="open = false"
              class="text-sand hover:text-gold transition text-xl leading-none p-1"
              aria-label="Close">&times;</button>
          </div>

          <div class="space-y-5">
            <div>
              <label class="block text-[10px] uppercase tracking-lux text-gold font-sans mb-2">Name</label>
              <input
                v-model="form.name"
                type="text"
                class="w-full bg-transparent border-b border-gold/20 focus:border-gold text-ivory-100 text-sm py-2 outline-none transition"
                style="border-top:none;border-left:none;border-right:none;border-radius:0;" />
            </div>
            <div>
              <label class="block text-[10px] uppercase tracking-lux text-gold font-sans mb-2">Email</label>
              <input
                v-model="form.email"
                type="email"
                class="w-full bg-transparent border-b border-gold/20 focus:border-gold text-ivory-100 text-sm py-2 outline-none transition"
                style="border-top:none;border-left:none;border-right:none;border-radius:0;" />
            </div>
            <div>
              <label class="block text-[10px] uppercase tracking-lux text-gold font-sans mb-2">How may we help?</label>
              <textarea
                v-model="form.message"
                rows="2"
                class="w-full bg-transparent border-b border-gold/20 focus:border-gold text-ivory-100 text-sm py-2 outline-none transition resize-none"
                style="border-top:none;border-left:none;border-right:none;border-radius:0;"></textarea>
            </div>
          </div>

          <div class="h-px bg-gold/20 my-7"></div>

          <div class="flex flex-col gap-2">
            <button
              @click="goChat"
              class="w-full py-3 text-[11px] uppercase tracking-lux text-ink-900 bg-ivory-100 hover:bg-gold transition font-sans">
              Chat with Advisor Now
            </button>
            <button
              @click="submit"
              :disabled="!form.email"
              class="w-full py-3 text-[11px] uppercase tracking-lux text-gold border border-gold/40 hover:border-gold hover:bg-gold/5 transition font-sans disabled:opacity-40 disabled:cursor-not-allowed">
              Request Callback
            </button>
          </div>

          <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-7 font-sans">
            Mon&ndash;Sat &middot; 09:00 &ndash; 18:00 GMT &middot; hello@ukteahouse.co.uk
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
