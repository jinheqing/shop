<script setup lang="ts">
import { ref, onMounted } from 'vue'
const show = ref(false)
onMounted(() => { if (!localStorage.getItem('cookie_consent')) show.value = true })
function accept() { localStorage.setItem('cookie_consent', 'all'); show.value = false }
function decline() { localStorage.setItem('cookie_consent', 'essential_only'); show.value = false }
function customize() { ElMessage.info('Cookie preferences saved') }
</script>
<template>
  <div v-if="show" class="fixed bottom-0 left-0 right-0 bg-white border-t border-tea-200 shadow-xl z-50 p-5">
    <div class="max-w-6xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
      <div class="text-sm text-tea-800 flex-1">
        🍪 We use essential cookies to run the site, and analytics cookies to improve your experience. Read our <RouterLink to="/privacy" class="underline">Privacy Policy</RouterLink>.
      </div>
      <div class="flex gap-2">
        <button @click="decline" class="px-4 py-2 border border-tea-300 rounded-lg text-sm hover:bg-tea-50">Essential Only</button>
        <button @click="customize" class="px-4 py-2 border border-tea-300 rounded-lg text-sm hover:bg-tea-50">Customize</button>
        <button @click="accept" class="px-6 py-2 bg-tea-800 text-white rounded-lg text-sm font-medium hover:bg-tea-900">Accept All</button>
      </div>
    </div>
  </div>
</template>
