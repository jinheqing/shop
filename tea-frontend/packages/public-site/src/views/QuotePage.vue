<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const token = route.params.token as string
const product = ref<any>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    product.value = await api.get(`/custom-products/by-token/${token}`)
  } catch (e: any) {
    error.value = 'Quote not found or no longer available'
  } finally {
    loading.value = false
  }
})
</script>
<template>
  <div class="pt-20 min-h-screen bg-tea-50">
    <div v-if="loading" class="text-center py-32"><div class="animate-spin w-10 h-10 border-4 border-tea-600 border-t-transparent rounded-full mx-auto"></div><div class="mt-4 text-tea-600">Loading your bespoke quote...</div></div>

    <div v-else-if="error" class="max-w-md mx-auto mt-20 p-8 bg-white rounded-2xl border border-tea-200 text-center">
      <div class="text-5xl mb-4">😔</div>
      <h2 class="font-serif text-xl mb-2">Quote Unavailable</h2>
      <p class="text-sm text-tea-600 mb-6">{{ error }}</p>
      <RouterLink to="/bespoke" class="inline-block px-6 py-3 bg-tea-800 text-white rounded-full">Create Your Own Bespoke</RouterLink>
    </div>

    <div v-else-if="product" class="max-w-5xl mx-auto px-6 py-12">
      <div class="mb-4 text-center">
        <span class="inline-block px-4 py-1.5 bg-tea-200 text-tea-800 rounded-full text-xs tracking-widest">BESPOKE QUOTE · SHAREABLE</span>
      </div>

      <div class="grid md:grid-cols-2 gap-8 mb-10">
        <!-- LEFT: Tea mountain visualization -->
        <div class="bg-gradient-to-br from-tea-700 to-tea-900 rounded-3xl p-8 text-white relative overflow-hidden">
          <div class="absolute top-0 right-0 text-[180px] opacity-10 leading-none">🏔️</div>
          <div class="relative">
            <div class="text-xs tracking-widest text-tea-200 mb-2">FROM THE MOUNTAIN</div>
            <h1 class="font-serif text-4xl leading-tight mb-6">{{ product.mountain_location }}</h1>
            <div class="space-y-3 text-tea-100 text-sm">
              <div class="flex items-center gap-2">
                <span>👨‍🍳</span><span><strong>Master:</strong> {{ product.master_name }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span>🌱</span><span><strong>Harvest:</strong> {{ product.harvest_date }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span>🔥</span><span><strong>Roast:</strong> {{ product.roasting_date }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span>📍</span><span><strong>Stored at:</strong> {{ product.storage_location }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span>🔬</span><span><strong>SGS Certified:</strong> ✅ Yes</span>
              </div>
            </div>
            <div class="mt-6">
              <el-button type="warning" @click="$router.push('/live')">📷 Watch Live Camera →</el-button>
            </div>
          </div>
        </div>

        <!-- RIGHT: Product details + order -->
        <div class="bg-white rounded-3xl p-8 border border-tea-100">
          <h2 class="font-serif text-3xl text-tea-900 mb-2">{{ product.title }}</h2>
          <div class="text-sm text-tea-600 mb-6">{{ product.raw_tea_source }}</div>

          <el-descriptions :column="1" border size="small" class="mb-6">
            <el-descriptions-item label="Tea Type">{{ product.tea_type }}</el-descriptions-item>
            <el-descriptions-item label="Shape">{{ product.tea_shape }}{{ product.tea_shape_weight ? ' · ' + product.tea_shape_weight + 'g' : '' }}</el-descriptions-item>
            <el-descriptions-item label="Flower Smoked">{{ product.smoked_with_flower ? 'Yes — ' + product.flower_type : 'No' }}</el-descriptions-item>
            <el-descriptions-item label="Inner Packaging">{{ product.inner_packaging }}</el-descriptions-item>
            <el-descriptions-item label="Outer Packaging">{{ product.outer_packaging }}</el-descriptions-item>
            <el-descriptions-item label="Lead Time">{{ product.lead_time }}</el-descriptions-item>
          </el-descriptions>

          <div class="bg-tea-50 rounded-xl p-5 mb-6">
            <div class="text-xs text-tea-500 mb-1">CUSTOMER REQUIREMENT</div>
            <div class="text-sm text-tea-800">{{ product.custom_requirement }}</div>
          </div>

          <div class="border-t pt-6">
            <div class="flex justify-between mb-1 text-sm"><span>Unit Price</span><span>£{{ product.unit_price.toFixed(2) }}</span></div>
            <div class="flex justify-between mb-1 text-sm"><span>Quantity</span><span>× {{ product.quantity }}</span></div>
            <div class="flex justify-between mb-1 text-sm"><span>Shipping (UK)</span><span>£{{ product.shipping_cost.toFixed(2) }}</span></div>
            <div class="flex justify-between font-serif text-2xl text-tea-900 mt-3 pt-3 border-t">
              <span>Total</span><span>£{{ product.total_amount.toFixed(2) }}</span>
            </div>
          </div>

          <div class="mt-6 space-y-3">
            <RouterLink :to="`/checkout/${token}`" class="block w-full text-center py-4 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition">
              🔒 Secure Checkout
            </RouterLink>
            <button class="w-full py-3 border border-tea-300 text-tea-800 rounded-xl hover:bg-tea-50">
              💬 Chat with Advisor About This Quote
            </button>
          </div>
        </div>
      </div>

      <div class="text-center text-sm text-tea-500">
        This bespoke quote was prepared by our advisor. All prices in GBP. <RouterLink to="/privacy" class="underline">Privacy</RouterLink> · <RouterLink to="/faq" class="underline">FAQ</RouterLink>
      </div>
    </div>
  </div>
</template>
