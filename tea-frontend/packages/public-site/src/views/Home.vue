<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getSlowPresets } from '@/api/live'
import type { CustomProduct } from '@/api/products'
import { getPublished } from '@/api/products'

const presets = ref<any[]>([])
const featured = ref<CustomProduct[]>([])

onMounted(async () => {
  try { presets.value = (await getSlowPresets()) || [] } catch {}
  try { featured.value = (await getPublished()) || [] } catch {}
})
</script>

<template>
  <!-- HERO -->
  <section class="hero-bg min-h-screen flex items-center text-white relative overflow-hidden">
    <!-- subtle gold grain overlay -->
    <div class="absolute inset-0 pointer-events-none opacity-30 mix-blend-overlay"
         style="background-image: radial-gradient(circle at 20% 30%, rgba(212,160,79,0.15), transparent 40%), radial-gradient(circle at 80% 70%, rgba(212,160,79,0.1), transparent 40%);"></div>

    <div class="max-w-7xl mx-auto px-6 py-32 text-center md:text-left relative">
      <span class="inline-block px-5 py-1.5 bg-white/10 backdrop-blur rounded-full text-xs tracking-[0.25em] uppercase mb-8 border border-white/15">
        Single Origin · Direct From Yunnan
      </span>
      <h1 class="font-display font-semibold text-5xl md:text-7xl lg:text-8xl leading-[1.05] max-w-4xl mb-8 tracking-tight">
        Pu'er Tea,<br/> <span class="italic">Traceable</span> to the Garden.
      </h1>
      <p class="text-lg md:text-xl text-white/80 max-w-2xl mb-12 leading-relaxed">
        Scan the QR on every bespoke box. See the tea garden, the master who rolled your leaves, and a live camera showing where it grew — 24 hours a day.
      </p>
      <div class="flex flex-col sm:flex-row gap-4 justify-center md:justify-start items-center md:items-start">
        <RouterLink to="/bespoke"
          class="px-9 py-4 bg-white text-tea-900 rounded-full font-medium hover:bg-tea-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-0.5 text-center tracking-wide">
          Create Your Blend →
        </RouterLink>
        <RouterLink to="/live"
          class="px-9 py-4 border border-white/50 rounded-full font-medium hover:bg-white/15 transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 text-center tracking-wide">
          📷 Watch Live Gardens
        </RouterLink>
      </div>
    </div>
  </section>

  <!-- FEATURES -->
  <section class="py-24 md:py-28 bg-white">
    <div class="max-w-7xl mx-auto px-6">
      <div class="grid md:grid-cols-3 gap-12 md:gap-16">
        <div class="text-center">
          <div class="w-16 h-16 bg-tea-100 rounded-full flex items-center justify-center text-2xl mx-auto mb-6">🏔️</div>
          <h3 class="font-display text-xl mb-3 text-tea-900">Cloud Garden Sourced</h3>
          <p class="text-tea-700 text-sm leading-relaxed">Every batch linked to its tea garden village, master name, and harvest date.</p>
        </div>
        <div class="text-center">
          <div class="w-16 h-16 bg-tea-100 rounded-full flex items-center justify-center text-2xl mx-auto mb-6">🔬</div>
          <h3 class="font-display text-xl mb-3 text-tea-900">SGS Certified</h3>
          <p class="text-tea-700 text-sm leading-relaxed">All products independently tested by SGS China for pesticides, heavy metals, microbiology.</p>
        </div>
        <div class="text-center">
          <div class="w-16 h-16 bg-tea-100 rounded-full flex items-center justify-center text-2xl mx-auto mb-6">📦</div>
          <h3 class="font-display text-xl mb-3 text-tea-900">Bespoke Blending</h3>
          <p class="text-tea-700 text-sm leading-relaxed">Choose tea garden, roast level, packaging. Your own tea — from leaf to cup — inside 45 days.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- LIVE STRIP -->
  <section class="py-20 bg-tea-50 border-y border-tea-100">
    <div class="max-w-7xl mx-auto px-6">
      <div class="flex items-center justify-between mb-10 flex-wrap gap-4">
        <div>
          <h2 class="font-display text-3xl md:text-4xl text-tea-900 mb-2">Live From the Tea Gardens</h2>
          <p class="text-tea-600 text-sm md:text-base">Watch your tea being picked, rolled and sun-dried in real time.</p>
        </div>
        <RouterLink to="/live" class="hidden md:inline-block text-tea-700 hover:text-tea-900 font-medium">View All →</RouterLink>
      </div>
      <div class="grid md:grid-cols-3 gap-6">
        <div v-for="p in (presets.length ? presets : [
          { name: '云南省临沧市临翔区邦东乡曼岗村茶园', location: '临沧', camera_rtmp_url: 'rtmp://localhost:1935/slow/manzhang' },
          { name: '云南省西双版纳州勐海县布朗山乡班章村茶园', location: '勐海', camera_rtmp_url: 'rtmp://localhost:1935/slow/banzhang' },
          { name: '云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园', location: '澜沧', camera_rtmp_url: 'rtmp://localhost:1935/slow/jingmai' },
        ])" :key="p.name" class="relative rounded-2xl overflow-hidden card-hover group">
          <div class="aspect-video bg-tea-900 flex items-center justify-center text-tea-200">
            <div class="text-center">
              <div class="text-5xl mb-2">🏞️</div>
              <div class="text-xs text-tea-400 font-mono">{{ p.camera_rtmp_url || 'rtmp://...' }}</div>
            </div>
          </div>
          <div class="absolute top-3 left-3 flex items-center gap-2 px-2.5 py-1 bg-red-600 text-white text-xs rounded-full font-medium shadow-sm">
            <span class="w-1.5 h-1.5 rounded-full bg-white live-dot"></span> LIVE
          </div>
          <div class="absolute bottom-0 inset-x-0 p-4 bg-gradient-to-t from-black/80 to-transparent text-white">
            <h3 class="font-display text-base md:text-lg">{{ p.name }}</h3>
            <p class="text-xs text-white/70">{{ p.location }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- FEATURED -->
  <section class="py-24 bg-white">
    <div class="max-w-7xl mx-auto px-6">
      <h2 class="font-display text-3xl md:text-4xl text-tea-900 mb-12">Our Bespoke Collection</h2>
      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-8">
        <RouterLink v-for="cp in (featured.length ? featured : [
          { id: 1, product_token: 'demo-1', title: '邦东古树饼', tea_garden_location: '云南省临沧市临翔区邦东乡曼岗村茶园', master_name: '李师傅', unit_price: 68.5, tea_type: 'raw_puer' },
          { id: 2, product_token: 'demo-2', title: '班章熟砖', tea_garden_location: '云南省西双版纳州勐海县布朗山乡班章村茶园', master_name: '张师傅', unit_price: 120, tea_type: 'ripe_puer' },
          { id: 3, product_token: 'demo-3', title: '景迈金瓜', tea_garden_location: '云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园', master_name: '王师傅', unit_price: 88, tea_type: 'raw_puer' },
        ])" :key="cp.id" :to="`/bespoke/${cp.product_token || 'demo-' + cp.id}`" class="group block card-hover rounded-2xl overflow-hidden border border-tea-100 bg-white">
          <div class="aspect-[4/3] bg-gradient-to-br from-tea-800 to-tea-500 flex items-center justify-center text-5xl group-hover:scale-105 transition-transform duration-500">🍵</div>
          <div class="p-6">
            <h3 class="font-display text-lg text-tea-900 mb-1">{{ cp.title }}</h3>
            <p class="text-xs text-tea-500 mb-3">{{ cp.tea_garden_location }} · Master {{ cp.master_name }}</p>
            <div class="flex items-center justify-between">
              <span class="text-tea-800 font-medium">£{{ cp.unit_price }}</span>
              <span class="text-xs text-tea-600 group-hover:text-tea-900">View →</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </section>
</template>
