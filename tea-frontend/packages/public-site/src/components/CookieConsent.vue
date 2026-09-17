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
  <div v-if="show" class="fixed bottom-4 right-4 max-w-md bg-white border shadow-xl rounded-lg p-5 z-50">
    <div class="font-semibold mb-2">🍪 We use cookies</div>
    <p class="text-sm text-gray-600 mb-4">
      Essential cookies are always on. We also use analytics and marketing cookies to improve your experience.
    </p>
    <div class="flex gap-2 flex-wrap">
      <button class="bg-gray-800 text-white px-4 py-2 rounded text-sm" @click="submit(false, false)">Essential Only</button>
      <button class="bg-blue-600 text-white px-4 py-2 rounded text-sm" @click="submit(true, false)">+ Analytics</button>
      <button class="bg-emerald-600 text-white px-4 py-2 rounded text-sm" @click="submit(true, true)">Accept All</button>
    </div>
    <div class="mt-3 text-xs text-gray-400">See our <a href="/privacy" class="underline">Privacy Policy</a></div>
  </div>
</template>
