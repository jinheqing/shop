<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'

const show = ref(false)

onMounted(() => {
  if (!localStorage.getItem('cookie_consent_seen')) {
    show.value = true
  }
})

async function submit(analytics: boolean, marketing: boolean) {
  const visitor = localStorage.getItem('visitor_id') || crypto.randomUUID()
  localStorage.setItem('visitor_id', visitor)
  localStorage.setItem('cookie_consent_seen', '1')
  localStorage.setItem('cookie_consent', JSON.stringify({ analytics, marketing, at: new Date().toISOString() }))
  try {
    await api.post('/cookie-consent', {
      consent_essential: true,
      consent_analytics: analytics,
      consent_marketing: marketing,
      visitor_id: visitor,
    })
  } catch {}
  show.value = false
}
</script>

<template>
  <Transition name="cookie-slide">
    <div v-if="show"
      class="fixed bottom-4 left-4 right-4 md:left-auto md:right-6 md:bottom-6 md:max-w-md bg-ink-900 text-ivory-100 border border-gold/30 p-5 z-50"
      style="border-radius: 2px;">
      <!-- Header — no emoji, hairline + serif title -->
      <div class="flex items-center gap-3 mb-3">
        <span class="w-1 h-1 rounded-full bg-gold"></span>
        <span class="font-serif text-base text-ivory-100 leading-snug">Cookies · On · The · House</span>
      </div>
      <p class="text-sm text-ivory-100/70 mb-5 leading-relaxed font-serif">
        Essential cookies keep the site working. We also use privacy-first analytics to understand how visitors experience our tea.
      </p>
      <!-- Buttons — squared, ink/ivory/gold, no rounded-full -->
      <div class="flex flex-wrap gap-2">
        <button class="px-4 py-2 text-[11px] uppercase tracking-lux font-sans transition border border-gold/30 text-ivory-100/80 hover:bg-gold/10"
          style="border-radius: 2px;"
          @click="submit(false, false)">Essential · Only</button>
        <button class="px-4 py-2 text-[11px] uppercase tracking-lux font-sans transition border border-gold/40 text-gold hover:bg-gold/10"
          style="border-radius: 2px;"
          @click="submit(true, false)">+ · Analytics</button>
        <button class="px-4 py-2 text-[11px] uppercase tracking-lux font-sans transition bg-gold text-ink-900 hover:bg-ivory-100"
          style="border-radius: 2px;"
          @click="submit(true, true)">Accept · All</button>
      </div>
      <div class="mt-4 text-[10px] uppercase tracking-lux text-ivory-100/40 font-sans">
        See our <a href="/privacy" class="text-gold underline hover:text-ivory-100">Privacy Policy</a>
        · <a href="/cookie-policy" class="text-gold underline hover:text-ivory-100">Cookie Details</a>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.cookie-slide-enter-active, .cookie-slide-leave-active { transition: opacity 0.3s ease, transform 0.3s ease; }
.cookie-slide-enter-from, .cookie-slide-leave-to { opacity: 0; transform: translateY(12px); }
</style>
