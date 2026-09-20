<script setup lang="ts">
// ============================================================
// PaymentResult.vue — 老钱审美支付结果页
// - 删除全部 emoji ✅⏳🚫❌🔄🔁📧🚚🔬
// - 不用 el-card，用纯 div + squared borders
// - ink/ivory/gold palette
// ============================================================
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
  <div class="pt-20 min-h-screen bg-ivory-100">
    <div class="max-w-lg mx-auto py-20 px-6">

      <!-- SUCCESS -->
      <div v-if="status==='success'" class="bg-white border border-gold/30 p-8 md:p-10 text-center" style="border-radius: 2px;">
        <!-- Mark — hairline gold square with seal dot, no emoji -->
        <div class="flex items-center justify-center mb-6">
          <div class="w-14 h-14 border border-gold flex items-center justify-center" style="border-radius: 2px;">
            <span class="w-2 h-2 rounded-full bg-gold"></span>
          </div>
        </div>
        <div class="flex items-center gap-3 justify-center mb-3">
          <span class="h-px w-8 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Acknowledged</span>
          <span class="h-px w-8 bg-gold/60"></span>
        </div>
        <h1 class="font-serif text-3xl text-ink-900 mb-3 leading-tight">Payment · Received</h1>
        <p class="text-sand font-serif mb-6 leading-relaxed">
          Thank you. Your bespoke tea order is now confirmed and being prepared for roasting.
        </p>

        <!-- Order summary -->
        <div v-if="order" class="bg-ink-900 text-ivory-100 p-5 mb-6 text-left" style="border-radius: 2px;">
          <div class="flex items-center gap-2 mb-3">
            <span class="w-1 h-1 rounded-full bg-gold"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Order · Reference</span>
          </div>
          <div class="font-mono text-lg text-ivory-100 mb-3">{{ order.order_no }}</div>
          <div class="text-sm text-ivory-100/80 font-serif">{{ order.custom_product_snapshot?.title }}</div>
          <div class="text-xs text-ivory-100/60 font-serif mt-1">{{ order.custom_product_snapshot?.tea_garden_location }}</div>
          <div class="mt-3 font-serif text-xl text-gold">£{{ order.total_amount }}</div>
        </div>

        <!-- Narrative next-steps, no emoji -->
        <div class="text-left space-y-2 mb-8">
          <div class="flex items-start gap-3">
            <span class="w-1 h-1 rounded-full bg-gold mt-2 shrink-0"></span>
            <p class="text-sm text-sand font-serif leading-relaxed">A confirmation has been dispatched to your email.</p>
          </div>
          <div class="flex items-start gap-3">
            <span class="w-1 h-1 rounded-full bg-gold mt-2 shrink-0"></span>
            <p class="text-sm text-sand font-serif leading-relaxed">Estimated delivery: 45 days from roasting.</p>
          </div>
          <div class="flex items-start gap-3">
            <span class="w-1 h-1 rounded-full bg-gold mt-2 shrink-0"></span>
            <p class="text-sm text-sand font-serif leading-relaxed">Your SGS certificate will be available in your account once the batch is tested.</p>
          </div>
        </div>

        <div class="space-y-2">
          <RouterLink to="/account"
            class="block w-full text-center py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            View · My · Orders
          </RouterLink>
          <RouterLink to="/chat"
            class="block w-full text-center py-4 border border-gold/30 text-ink-900 hover:bg-ivory-50 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Speak · With · Your · Advisor
          </RouterLink>
        </div>
      </div>

      <!-- PENDING -->
      <div v-else-if="status==='pending'" class="bg-white border border-gold/20 p-8 md:p-10 text-center" style="border-radius: 2px;">
        <div class="flex items-center justify-center mb-6">
          <div class="w-14 h-14 border border-gold/40 flex items-center justify-center" style="border-radius: 2px;">
            <span class="w-2 h-2 rounded-full bg-gold/60"></span>
          </div>
        </div>
        <div class="flex items-center gap-3 justify-center mb-3">
          <span class="h-px w-8 bg-gold/40"></span>
          <span class="text-[10px] uppercase tracking-lux text-sand font-sans">Awaiting</span>
          <span class="h-px w-8 bg-gold/40"></span>
        </div>
        <h1 class="font-serif text-3xl text-ink-900 mb-3 leading-tight">Payment · Pending</h1>
        <p class="text-sand font-serif mb-6 leading-relaxed">
          Your payment hasn't been confirmed yet. This may take a few minutes.
        </p>
        <button @click="window.location.reload()"
          class="px-6 py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
          style="border-radius: 2px;">
          Check · Again
        </button>
      </div>

      <!-- FAILED / CANCELLED -->
      <div v-else class="bg-white border border-gold/20 p-8 md:p-10 text-center" style="border-radius: 2px;">
        <div class="flex items-center justify-center mb-6">
          <div class="w-14 h-14 border border-gold/30 flex items-center justify-center" style="border-radius: 2px;">
            <span class="h-3 w-px bg-gold/60 rotate-45"></span>
            <span class="h-3 w-px bg-gold/60 -rotate-45 -ml-px"></span>
          </div>
        </div>
        <div class="flex items-center gap-3 justify-center mb-3">
          <span class="h-px w-8 bg-gold/40"></span>
          <span class="text-[10px] uppercase tracking-lux text-sand font-sans">{{ status === 'cancelled' ? 'Cancelled' : 'Declined' }}</span>
          <span class="h-px w-8 bg-gold/40"></span>
        </div>
        <h1 class="font-serif text-3xl text-ink-900 mb-3 leading-tight">
          {{ status==='cancelled' ? 'Payment · Cancelled' : 'Payment · Declined' }}
        </h1>
        <p class="text-sand font-serif mb-6 leading-relaxed">
          Your payment was not completed. Your quote is still available with your advisor.
        </p>
        <div class="space-y-2">
          <RouterLink to="/chat"
            class="block w-full text-center py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Speak · With · Your · Advisor
          </RouterLink>
          <RouterLink to="/contact"
            class="block w-full text-center py-4 border border-gold/30 text-ink-900 hover:bg-ivory-50 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Need · Help
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
