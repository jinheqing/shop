<script setup lang="ts">
// ============================================================
// QuotePage.vue — 顾问确认模式（Old Money Edition）
// ============================================================
// 老钱定制惯例（Hermès Bespoke / Rolls-Royce Private Office）：
//   网页只展示"配置概览" + "Request the Advisor to Confirm"
//   不在网页直接显示总价 — 总价由顾问私信交付
//   单一主行动，二级动作收进文字链接
// ============================================================

import { onMounted, ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const token = route.params.token as string
const product = ref<any>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    product.value = await api.get(`/custom-products/by-token/${token}`)
  } catch {
    error.value = 'This commission record was not found, or is no longer available.'
  } finally {
    loading.value = false
  }
})

function goChat() {
  const t = localStorage.getItem('user_token')
  if (!t) router.push(`/magic-link?redirect=/chat`)
  else router.push('/chat')
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100">

    <!-- Loading -->
    <div v-if="loading" class="text-center py-32">
      <div class="inline-block w-10 h-10 border border-gold/30 border-t-gold rounded-full animate-spin"></div>
      <div class="mt-4 text-[11px] uppercase tracking-lux text-sand font-sans">Loading · Your · Commission · Record</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="max-w-md mx-auto mt-20 p-10 bg-white border border-gold/20 text-center" style="border-radius: 2px;">
      <div class="flex items-center gap-3 justify-center mb-6">
        <span class="h-px w-10 bg-gold/50"></span>
        <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Record · Unavailable</span>
        <span class="h-px w-10 bg-gold/50"></span>
      </div>
      <h2 class="font-serif text-2xl text-ink-900 mb-3">This record cannot be shown.</h2>
      <p class="text-sm text-sand mb-8 font-serif">{{ error }}</p>
      <RouterLink to="/bespoke" class="inline-block px-8 py-3 border border-gold/40 text-gold hover:bg-gold/5 transition text-[11px] uppercase tracking-lux font-sans">Begin · A · New · Commission</RouterLink>
    </div>

    <!-- Commission record -->
    <div v-else-if="product" class="max-w-5xl mx-auto px-6 py-12 md:py-16">

      <!-- Section title -->
      <div class="mb-10 md:mb-12 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Commission · {{ token }}</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-3 leading-tight">{{ product.title }}</h1>
        <p class="text-sand text-sm font-serif">A commission record, prepared by your advisor. The quotation will follow by private message.</p>
      </div>

      <div class="grid md:grid-cols-2 gap-8">

        <!-- LEFT — Provenance, off-black 主舞台 -->
        <div class="bg-ink-900 text-ivory-100 p-8 md:p-10 relative overflow-hidden" style="border-radius: 2px;">
          <div class="absolute inset-0 pointer-events-none opacity-[0.07]"
               style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 65%, #C5A572 0%, transparent 50%);"></div>

          <div class="relative">
            <div class="flex items-center gap-3 mb-5">
              <span class="h-px w-8 bg-gold/60"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">From · The · Mountain</span>
            </div>
            <h2 class="font-serif text-3xl md:text-4xl text-ivory-100 leading-tight mb-7">{{ product.tea_garden_location }}</h2>

            <div class="space-y-4 text-sm">
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Master</div>
                <div class="font-serif text-ivory-100">{{ product.master_name }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Harvested</div>
                <div class="font-serif text-ivory-100">{{ product.harvest_date }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Roasted</div>
                <div class="font-serif text-ivory-100">{{ product.roasting_date }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Stored · At</div>
                <div class="font-serif text-ivory-100">{{ product.storage_location }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">SGS · Independently · Tested</div>
                <div class="font-serif text-gold">Passed</div>
              </div>
            </div>

            <div class="mt-8 pt-6 border-t border-gold/20">
              <RouterLink to="/live"
                class="inline-flex items-center gap-2 text-[11px] uppercase tracking-lux text-gold hover:text-gold-soft transition-duration-lux font-sans">
                <span class="w-1 h-1 rounded-full bg-gold"></span>
                Watch · Live · Camera · From · The · Mountain &rarr;
              </RouterLink>
            </div>
          </div>
        </div>

        <!-- RIGHT — Configuration summary, 象牙白面板 -->
        <div class="bg-white border border-gold/15 p-8 md:p-10" style="border-radius: 2px;">
          <div class="flex items-center gap-3 mb-5">
            <span class="h-px w-8 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Configuration · Summary</span>
          </div>

          <h3 class="font-serif text-2xl md:text-3xl text-ink-900 mb-1">{{ product.title }}</h3>
          <div class="text-[11px] uppercase tracking-lux text-sand font-sans mb-7">{{ product.raw_tea_source }}</div>

          <!-- 叙述式描述列表，不用 el-descriptions 的表格风格 -->
          <dl class="space-y-4 mb-7">
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Tea · Type</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.tea_type }}</dd>
            </div>
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Shape</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.tea_shape }}{{ product.tea_shape_weight ? ' · ' + product.tea_shape_weight + 'g' : '' }}</dd>
            </div>
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Flower · Smoked</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.smoked_with_flower ? 'Yes — ' + product.flower_type : 'No' }}</dd>
            </div>
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Inner · Packaging</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.inner_packaging }}</dd>
            </div>
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Outer · Packaging</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.outer_packaging }}</dd>
            </div>
            <div class="flex justify-between gap-4 pb-3 border-b border-gold/10">
              <dt class="text-[11px] uppercase tracking-lux text-sand font-sans">Lead · Time</dt>
              <dd class="font-serif text-ink-900 text-right">{{ product.lead_time }}</dd>
            </div>
          </dl>

          <!-- Notes from customer -->
          <div v-if="product.custom_requirement" class="bg-ivory-100 p-5 mb-7" style="border-radius: 2px;">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">Your · Notes</div>
            <div class="text-sm text-ink-900 font-serif leading-relaxed">{{ product.custom_requirement }}</div>
          </div>

          <!-- 报价提示：不在网页直接显示总价 -->
          <div class="bg-ink-900 text-ivory-100 p-5 mb-7" style="border-radius: 2px;">
            <div class="flex items-center gap-2 mb-2">
              <span class="w-1 h-1 rounded-full bg-gold"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Quotation · By · Private · Message</span>
            </div>
            <p class="text-sm text-sand leading-relaxed font-serif">
              Final pricing — including any bespoke additions and shipping —
              will be confirmed by your advisor before roasting begins.
            </p>
          </div>

          <!-- 单一主行动 -->
          <button @click="goChat"
            class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Request · The · Advisor · To · Confirm
          </button>

          <!-- 二级动作 — 文字链接 -->
          <div class="mt-6 text-center">
            <RouterLink :to="`/traceability/${token}`"
              class="text-[11px] uppercase tracking-lux text-gold hover:text-gold-soft transition-duration-lux font-sans">
              View · Full · Provenance · Record &rarr;
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- 极细 disclaimer -->
      <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-10 leading-loose font-sans">
        All prices in GBP &middot; <RouterLink to="/privacy" class="hover:text-gold transition">Privacy</RouterLink>
        &middot; <RouterLink to="/faq" class="hover:text-gold transition">FAQ</RouterLink>
      </p>
    </div>
  </div>
</template>
