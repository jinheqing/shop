<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const orderId = parseInt(route.query.order_id as string || '0')
const status = ref<'success' | 'pending' | 'failed' | 'cancelled'>('pending')
const order = ref<any>(null)

onMounted(async () => {
  if (!orderId) { status.value = 'failed'; return }
  try {
    order.value = await api.get(`/orders/${orderId}`)
    if (order.value.state === 'paid' || order.value.state !== 'ordering') status.value = 'success'
    else status.value = 'pending'
  } catch { status.value = 'failed' }
})
</script>
<template>
  <div class="pt-20 min-h-screen bg-tea-50">
    <div class="max-w-lg mx-auto py-20 px-6">
      <el-card v-if="status==='success'" class="text-center p-8">
        <div class="text-7xl mb-6">✅</div>
        <h1 class="font-serif text-3xl text-green-700 mb-2">Payment Successful</h1>
        <p class="text-tea-700 mb-6">Thank you! Your bespoke tea order is now confirmed and being prepared.</p>
        <div v-if="order" class="bg-tea-50 rounded-xl p-5 mb-6 text-left">
          <div class="text-xs text-tea-500 mb-1">Order Reference</div>
          <div class="font-mono font-semibold text-lg">{{ order.order_no }}</div>
          <div class="mt-3 text-sm"><strong>{{ order.custom_product_snapshot?.title }}</strong></div>
          <div class="text-sm text-tea-600">{{ order.custom_product_snapshot?.tea_garden_location }}</div>
          <div class="mt-3 font-serif text-xl">£{{ order.total_amount }}</div>
        </div>
        <div class="space-y-3 text-sm text-tea-600">
          <div>📧 Order confirmation sent to your email</div>
          <div>🚚 Estimated delivery: 45 days</div>
          <div>🔬 SGS certificate will be available in your account</div>
        </div>
        <div class="mt-8 space-y-2">
          <RouterLink to="/account" class="block w-full text-center py-3 bg-tea-800 text-white rounded-xl">View My Orders →</RouterLink>
          <RouterLink to="/chat" class="block w-full text-center py-3 border border-tea-300 rounded-xl">💬 Chat with Your Advisor</RouterLink>
        </div>
      </el-card>

      <el-card v-else-if="status==='pending'" class="text-center p-8">
        <div class="text-7xl mb-6">⏳</div>
        <h1 class="font-serif text-3xl text-amber-700 mb-2">Payment Pending</h1>
        <p class="text-tea-700 mb-6">Your payment hasn't been confirmed yet. This may take a few minutes.</p>
        <button @click="window.location.reload()" class="px-6 py-3 bg-tea-800 text-white rounded-xl">🔄 Check Again</button>
      </el-card>

      <el-card v-else class="text-center p-8">
        <div class="text-7xl mb-6">{{ status==='cancelled' ? '🚫' : '❌' }}</div>
        <h1 class="font-serif text-3xl text-red-700 mb-2">{{ status==='cancelled' ? 'Payment Cancelled' : 'Payment Failed' }}</h1>
        <p class="text-tea-700 mb-6">Your payment was not completed. Your quote is still available.</p>
        <div class="space-y-2">
          <button class="w-full py-3 bg-tea-800 text-white rounded-xl">🔁 Try Again</button>
          <RouterLink to="/contact" class="block w-full text-center py-3 border border-tea-300 rounded-xl">Need Help?</RouterLink>
        </div>
      </el-card>
    </div>
  </div>
</template>
