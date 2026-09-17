<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'

const gardens = ref<any[]>([])
const loading = ref(true)

async function load() {
  try {
    const d: any = await api.get('/public/slow-presets')
    gardens.value = d?.items || d || []
  } catch (e) {
    // Fallback demo — will be replaced by real data once backend runs
    gardens.value = [
      { name: '云南省临沧市临翔区邦东乡曼岗村茶园', location: '临沧 · 云南', description: '古树普洱核心产区' },
      { name: '广东省潮州市潮安区凤凰镇大乌岽村茶园', location: '潮州 · 广东', description: '凤凰单丛' },
    ]
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>

<template>
  <div class="max-w-6xl mx-auto py-10 px-4">
    <h1 class="text-4xl font-bold mb-4">Our Tea Gardens</h1>
    <p class="text-gray-600 mb-8">Direct from village-level tea gardens in Yunnan & Guangdong. Each garden has a 24/7 live camera.</p>

    <div v-if="loading" class="text-center py-20 text-gray-400">Loading…</div>

    <div v-else class="grid md:grid-cols-2 gap-6">
      <div v-for="g in gardens" :key="g.id || g.name" class="border rounded-lg p-5 shadow-sm hover:shadow-md transition">
        <h3 class="text-lg font-semibold mb-2">{{ g.name || 'Unnamed Garden' }}</h3>
        <p class="text-sm text-gray-500 mb-2">{{ g.location || g.region || '—' }}</p>
        <p class="text-gray-700 text-sm mb-4">{{ g.description || 'Visit our live streams below to see this tea garden in action.' }}</p>
        <div class="flex gap-2 text-xs">
          <span v-if="g.status === 'live'" class="bg-red-600 text-white px-2 py-0.5 rounded flex items-center gap-1">
            <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse"></span> LIVE
          </span>
          <a href="/live" class="text-emerald-700 hover:underline">Watch Live →</a>
        </div>
      </div>
    </div>
  </div>
</template>
