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
    <div v-if="show" class="fixed bottom-4 left-4 right-4 md:left-auto md:right-6 md:bottom-6 md:max-w-md bg-white border border-tea-100 shadow-2xl rounded-2xl p-5 z-50">
      <div class="flex items-start gap-3 mb-3">
        <span class="text-2xl leading-none">🍪</span>
        <div class="font-display text-base text-tea-900 leading-snug">We use cookies on tea</div>
      </div>
      <p class="text-sm text-tea-700 mb-4 leading-relaxed">
        Essential cookies keep the site working. We also use analytics to understand how visitors experience our tea.
      </p>
      <div class="flex flex-wrap gap-2">
        <button class="px-4 py-2 rounded-full text-sm transition border border-tea-300 text-tea-800 hover:bg-tea-100"
          @click="submit(false, false)">Essential Only</button>
        <button class="px-4 py-2 rounded-full text-sm transition bg-tea-700 text-white hover:bg-tea-800"
          @click="submit(true, false)">+ Analytics</button>
        <button class="px-4 py-2 rounded-full text-sm transition bg-tea-900 text-white hover:bg-tea-950"
          @click="submit(true, true)">Accept All</button>
      </div>
      <div class="mt-3 text-xs text-tea-500">See our <a href="/privacy" class="underline hover:text-tea-900">Privacy Policy</a> · <a href="/privacy#cookies" class="underline hover:text-tea-900">Cookie Details</a></div>
    </div>
  </Transition>
</template>

<style scoped>
.cookie-slide-enter-active, .cookie-slide-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.cookie-slide-enter-from, .cookie-slide-leave-to { opacity: 0; transform: translateY(12px); }
</style>
