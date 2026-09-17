<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const mountainSlug = route.params.slug as string

const mountains: any = {
  'iceland-old-village': {
    name: '冰岛老寨', region: 'Yunnan · Lincang', altitude: 2200, master: '王师傅',
    description: 'Iceland Old Village — arguably the most sought-after Pu\'er terroir in the world. Ancient tea trees (300+ years) grow wild on steep mountain slopes, producing tea with a unique floral sweetness.',
    gps: '23.93°N, 100.15°E',
    stream_url: 'rtmp://localhost:1935/slow/iceland',
    harvest_season: 'Late March → Early April',
    sgs_report: 'SGS-ICL-2026-003',
    certifications: ['Organic', 'Natural Farming', 'Single Origin Certified'],
    teas_available: ['Raw Pu\'er (生普)', 'Ancient Tree (古树)'],
  },
  'banzhang': {
    name: '班章村', region: 'Yunnan · Menghai', altitude: 1700, master: '李师傅',
    description: 'Banzhang Village — the "King of Pu\'er". Rich, mineral-forward character with legendary aging potential.',
    gps: '22.21°N, 100.65°E',
    stream_url: 'rtmp://localhost:1935/slow/banzhang',
    harvest_season: 'April',
    sgs_report: 'SGS-BZ-2026-007',
    certifications: ['Organic', 'Natural Farming'],
    teas_available: ['Raw Pu\'er', 'Ripe Pu\'er (熟普)'],
  },
  'jingmai': {
    name: '景迈山', region: 'Yunnan · Lancang', altitude: 1400, master: '陶师傅',
    description: 'Jingmai Mountain — 1300-year-old tea heritage site. Wild forest tea with pronounced orchid fragrance.',
    gps: '22.28°N, 100.02°E',
    stream_url: 'rtmp://localhost:1935/slow/jingmai',
    harvest_season: 'March',
    sgs_report: 'SGS-JM-2026-012',
    certifications: ['UNESCO Heritage Candidate', 'Organic', 'Natural Farming'],
    teas_available: ['Raw Pu\'er', 'Tea Cake', 'Loose Leaf'],
  },
}
const mountain = mountains[mountainSlug] || mountains['iceland-old-village']
</script>
<template>
  <div class="pt-20 min-h-screen">
    <section class="relative h-96 bg-gradient-to-br from-tea-800 to-tea-950 flex items-center">
      <div class="absolute inset-0 bg-black/40"></div>
      <div class="relative max-w-6xl mx-auto px-6 text-white">
        <RouterLink to="/mountains" class="text-sm text-tea-200 hover:text-white">← All Mountains</RouterLink>
        <h1 class="font-serif text-5xl mt-4">{{ mountain.name }}</h1>
        <div class="mt-2 text-tea-200">{{ mountain.region }} · {{ mountain.altitude }}m · 📍 {{ mountain.gps }}</div>
      </div>
    </section>

    <section class="max-w-6xl mx-auto px-6 py-12 grid md:grid-cols-3 gap-8">
      <div class="md:col-span-2">
        <h2 class="font-serif text-2xl text-tea-900 mb-4">About This Mountain</h2>
        <p class="text-tea-700 leading-relaxed mb-8">{{ mountain.description }}</p>

        <div class="aspect-video bg-black rounded-2xl overflow-hidden relative">
          <div class="absolute inset-0 flex items-center justify-center text-tea-100">
            <div class="text-center">
              <div class="inline-flex items-center gap-2 px-3 py-1 bg-red-600 rounded-full text-xs mb-4">
                <span class="w-2 h-2 rounded-full bg-white animate-pulse"></span> LIVE FROM MOUNTAIN
              </div>
              <div class="text-4xl mb-2">📷</div>
              <div class="text-xs text-tea-300">Streaming 24/7 from {{ mountain.name }}</div>
              <div class="text-[10px] text-tea-400 mt-1">{{ mountain.stream_url }}</div>
            </div>
          </div>
        </div>
      </div>

      <div>
        <div class="bg-white rounded-2xl border border-tea-100 p-6 mb-4">
          <h3 class="font-serif text-lg mb-4">Mountain Info</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between"><span class="text-tea-500">Master Tea Maker</span><span class="font-medium">{{ mountain.master }}</span></div>
            <div class="flex justify-between"><span class="text-tea-500">Altitude</span><span>{{ mountain.altitude }}m</span></div>
            <div class="flex justify-between"><span class="text-tea-500">Harvest</span><span>{{ mountain.harvest_season }}</span></div>
            <div class="flex justify-between"><span class="text-tea-500">SGS Report</span><span class="font-mono text-xs">{{ mountain.sgs_report }}</span></div>
          </div>
        </div>

        <div class="bg-white rounded-2xl border border-tea-100 p-6 mb-4">
          <h3 class="font-serif text-lg mb-4">Certifications</h3>
          <div class="flex flex-wrap gap-2">
            <span v-for="c in mountain.certifications" :key="c" class="px-2 py-1 bg-tea-100 text-tea-800 rounded text-xs">{{ c }}</span>
          </div>
        </div>

        <div class="bg-tea-800 rounded-2xl p-6 text-white text-center">
          <div class="text-sm mb-2">Want tea from this mountain?</div>
          <RouterLink to="/bespoke" class="inline-block px-6 py-3 bg-white text-tea-900 rounded-full font-medium mt-2">Create Your Bespoke →</RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>
