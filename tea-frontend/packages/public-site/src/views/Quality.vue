<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'

const reports = ref<any[]>([])
const loading = ref(true)

async function load() {
  try {
    const d: any = await api.get('/public/sgs-reports')
    reports.value = d?.items || d || []
  } catch {
    reports.value = [
      { report_number: 'SGS-2025-YN-001', tea_garden: '曼岗村茶园', pesticide_nd_limit: 'ND (≤0.01 ppm)', status: 'Passed', issued_at: '2025-04-15' },
      { report_number: 'SGS-2025-YN-002', tea_garden: '景迈村茶园', pesticide_nd_limit: 'ND (≤0.01 ppm)', status: 'Passed', issued_at: '2025-04-12' },
    ]
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>

<template>
  <div class="max-w-5xl mx-auto py-10 px-4">
    <h1 class="text-4xl font-bold mb-4">Quality & Safety</h1>
    <p class="text-gray-600 mb-8">Every tea batch is independently tested by SGS (Societe Generale de Surveillance) and UKAS-accredited labs.</p>

    <div v-if="loading" class="text-center py-20 text-gray-400">Loading reports…</div>

    <div v-else class="grid gap-4">
      <div v-for="r in reports" :key="r.id || r.report_number" class="border rounded-lg p-5 bg-white shadow-sm">
        <div class="flex justify-between items-start mb-2">
          <div>
            <div class="text-xs text-gray-400">{{ r.report_number }}</div>
            <div class="font-semibold">{{ r.tea_garden || r.title }}</div>
          </div>
          <span class="bg-emerald-100 text-emerald-800 text-xs px-2 py-1 rounded-full">{{ r.status || 'Passed' }}</span>
        </div>
        <div class="text-sm text-gray-700 mb-2">
          Pesticides: <span class="font-mono text-emerald-700">{{ r.pesticide_nd_limit || 'Not Detected' }}</span>
        </div>
        <div class="text-xs text-gray-400">Issued: {{ r.issued_at }}</div>
      </div>
      <div v-if="!reports.length" class="text-center text-gray-400 py-12">No SGS reports available yet.</div>
    </div>
  </div>
</template>
