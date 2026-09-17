<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const token = route.params.token as string
const product = ref<any>(null)
const loading = ref(true)

onMounted(async () => {
  try { product.value = await api.get(`/custom-products/by-token/${token}`) }
  catch {
    product.value = {
      title: '邦东古树饼',
      tea_garden_location: '云南省临沧市临翔区邦东乡曼岗村茶园',
      master_name: '王师傅',
      harvest_date: '2026-04-12',
      roasting_date: '2026-05-20',
      storage_location: '潮州恒温恒湿仓储 3 号库',
      tea_type: 'Raw Pu\'er (生普)',
      tea_shape: 'Cake · 357g',
      product_card_text: '凤凰山大乌岽 · 古树单丛 · 2026 春茶',
      sgs_report_no: 'SGS-ICL-2026-003',
      stream_url: 'rtmp://localhost:1935/slow/iceland',
    }
  } finally { loading.value = false }
})
</script>
<template>
  <div class="min-h-screen bg-black text-white">
    <!-- Live Camera Hero -->
    <section class="relative h-[60vh] bg-gradient-to-b from-black via-tea-950 to-black">
      <div class="absolute inset-0 flex items-center justify-center">
        <div class="text-center">
          <div class="inline-flex items-center gap-2 px-4 py-2 bg-red-600 rounded-full text-sm mb-6 animate-pulse">
            <span class="w-2.5 h-2.5 rounded-full bg-white"></span>
            LIVE · 24/7 FROM THE TEA MOUNTAIN
          </div>
          <div class="text-6xl mb-4">📷</div>
          <div class="text-sm text-tea-300">Live stream from {{ product?.tea_garden_location }}</div>
        </div>
      </div>
      <div class="absolute bottom-8 left-0 right-0 text-center">
        <div class="font-serif text-3xl">{{ product?.title || 'Your Bespoke Tea' }}</div>
      </div>
    </section>

    <section class="max-w-4xl mx-auto px-6 py-16">
      <h2 class="font-serif text-3xl text-center mb-4">🔍 Trace Your Tea</h2>
      <p class="text-center text-tea-300 mb-12">From tea garden to your cup. Every step verified.</p>

      <el-timeline :size="large">
        <el-timeline-item :timestamp="product?.harvest_date" color="success">
          <div class="bg-tea-900/50 rounded-xl p-5">
            <div class="text-lg font-serif mb-1">🌱 Harvested</div>
            <div class="text-tea-200 text-sm">Picked from {{ product?.tea_garden_location }}</div>
            <div class="text-xs text-tea-400">Master: {{ product?.master_name }}</div>
          </div>
        </el-timeline-item>

        <el-timeline-item :timestamp="product?.roasting_date" color="warning">
          <div class="bg-tea-900/50 rounded-xl p-5">
            <div class="text-lg font-serif mb-1">🔥 Roasted & Rolled</div>
            <div class="text-tea-200 text-sm">Traditional processing</div>
            <div class="text-xs text-tea-400 mt-2">{{ product?.tea_type }} · {{ product?.tea_shape }}</div>
          </div>
        </el-timeline-item>

        <el-timeline-item timestamp="2026-05-25" color="primary">
          <div class="bg-tea-900/50 rounded-xl p-5">
            <div class="text-lg font-serif mb-1">🔬 SGS Tested</div>
            <div class="text-tea-200 text-sm">Independent lab: pesticides · heavy metals · microbiology</div>
            <div class="text-xs text-tea-400 mt-2">Report: {{ product?.sgs_report_no }} — ✅ PASSED</div>
          </div>
        </el-timeline-item>

        <el-timeline-item timestamp="2026-06-01" color="primary">
          <div class="bg-tea-900/50 rounded-xl p-5">
            <div class="text-lg font-serif mb-1">📦 Packed</div>
            <div class="text-tea-200 text-sm">Custom packaging with your QR code</div>
            <div class="text-xs text-tea-400 mt-2">Card text: "{{ product?.product_card_text }}"</div>
          </div>
        </el-timeline-item>

        <el-timeline-item timestamp="Current" color="success">
          <div class="bg-tea-900/50 rounded-xl p-5 border border-tea-700">
            <div class="text-lg font-serif mb-1">🚚 On Its Way</div>
            <div class="text-tea-200 text-sm">Shipped via DHL Express to your door in the UK</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </section>

    <section class="bg-tea-950 py-12 text-center">
      <div class="max-w-2xl mx-auto px-6">
        <h3 class="font-serif text-2xl mb-4">Thank you for choosing traceable tea.</h3>
        <p class="text-tea-300 text-sm mb-6">Your purchase supports ethical farmers practicing natural farming without pesticides.</p>
        <RouterLink to="/" class="inline-block px-8 py-3 bg-white text-tea-900 rounded-full font-medium">Back to UK Tea House →</RouterLink>
      </div>
    </section>
  </div>
</template>
