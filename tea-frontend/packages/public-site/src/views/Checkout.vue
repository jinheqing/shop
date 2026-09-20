<script setup lang="ts">
// ============================================================
// Checkout.vue — 老钱审美结算页
// ============================================================
// - off-black Hero + 象牙白表单
// - 不用 emoji 当 icon，用 hairline 金线分隔
// - 不用 alert() 弹原生对话框，用 in-page 状态条
// - 单一主行动 "Confirm · Order · With · 2Checkout"
// ============================================================

import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getByToken, type CustomProduct } from '@/api/products'

const route = useRoute()
const router = useRouter()
const product = ref<CustomProduct | null>(null)
const qty = ref(1)
const form = ref({
  name: '', email: '', phone: '',
  billing: { address: '', city: 'London', postcode: '', country: 'UK' },
  delivery: { address: '', city: 'London', postcode: '', country: 'UK' },
})
const loading = ref(false)
const errorMsg = ref('')

onMounted(async () => {
  try { product.value = await getByToken(route.params.token as string) }
  catch {
    product.value = { id: 1, unit_price: 68.5, shipping_cost: 120, title: 'Demo Product', tea_type: 'raw_puer', tea_shape: 'cake', status: 'published', version: 1 } as any
  }
})

const total = () => (product.value?.unit_price || 0) * qty.value + (product.value?.shipping_cost || 0)

const place = async () => {
  errorMsg.value = ''
  loading.value = true
  try {
    const userToken = localStorage.getItem('user_token')
    if (!userToken) {
      router.push('/magic-link?redirect=' + encodeURIComponent('/checkout/' + route.params.token))
      return
    }
    const order = await fetch('/api/v1/orders', {
      method:'POST', headers:{'Content-Type':'application/json','Authorization':'Bearer '+userToken},
      body: JSON.stringify({
        custom_product_id: product.value!.id, unit_price: product.value!.unit_price,
        quantity: qty.value, shipping_cost: product.value!.shipping_cost,
        billing_address: { ...form.value.billing, name: form.value.name, email: form.value.email },
        delivery_address: { ...form.value.delivery, name: form.value.name },
      })
    }).then(x=>x.json())
    if (order.code && order.code !== 0 && order.code !== 201) {
      errorMsg.value = 'Order could not be placed: ' + (order.message || 'Please try again, or speak with your advisor.')
      return
    }
    const pay = await fetch(`/api/v1/orders/${order.id}/payment/init`, {
      method:'POST', headers:{'Content-Type':'application/json','Authorization':'Bearer '+userToken},
      body: JSON.stringify({ gateway: '2checkout' })
    }).then(x=>x.json())
    if (pay.checkout_url) window.location.href = pay.checkout_url
    else router.push('/')
  } catch (e: any) {
    errorMsg.value = 'An error occurred: ' + e.message + '. Please try again, or speak with your advisor.'
  } finally { loading.value = false }
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100">
    <div class="max-w-5xl mx-auto px-6 py-12 md:py-16">

      <!-- Section title -->
      <div class="mb-10 md:mb-12 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Secure · Checkout</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-3">Confirm · Your · Commission.</h1>
        <p class="text-sand text-sm font-serif">A secure payment, processed by 2Checkout. Your card details are never stored on our servers.</p>
      </div>

      <!-- Error banner — in-page, not native alert -->
      <div v-if="errorMsg" class="mb-8 p-5 bg-ink-900 text-ivory-100" style="border-radius: 2px;">
        <div class="flex items-center gap-2 mb-2">
          <span class="w-1 h-1 rounded-full bg-gold"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Attention</span>
        </div>
        <p class="text-sm text-sand font-serif leading-relaxed">{{ errorMsg }}</p>
      </div>

      <div class="grid md:grid-cols-3 gap-8">

        <!-- LEFT — Form -->
        <div class="md:col-span-2 space-y-6">

          <!-- Your details -->
          <div class="bg-white border border-gold/15 p-6 md:p-8" style="border-radius: 2px;">
            <div class="flex items-center gap-3 mb-5">
              <span class="h-px w-8 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Your · Details</span>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <input v-model="form.name" placeholder="Full name"
                class="px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;" />
              <input v-model="form.email" type="email" placeholder="Email"
                class="px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;" />
            </div>
          </div>

          <!-- Billing address -->
          <div class="bg-white border border-gold/15 p-6 md:p-8" style="border-radius: 2px;">
            <div class="flex items-center gap-3 mb-5">
              <span class="h-px w-8 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Billing · Address</span>
            </div>
            <input v-model="form.billing.address" placeholder="Street address"
              class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900 mb-3"
              style="border-radius: 2px;" />
            <div class="grid grid-cols-3 gap-3">
              <input v-model="form.billing.city" placeholder="City"
                class="px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;" />
              <input v-model="form.billing.postcode" placeholder="Postcode"
                class="px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;" />
              <input :value="'UK'" placeholder="Country"
                class="px-4 py-3 bg-white border border-gold/20 font-serif text-ink-900"
                style="border-radius: 2px;" />
            </div>
          </div>
        </div>

        <!-- RIGHT — Order Summary (sticky) -->
        <div class="bg-white border border-gold/15 p-6 md:p-8 h-fit md:sticky md:top-24" style="border-radius: 2px;">
          <div class="flex items-center gap-3 mb-5">
            <span class="h-px w-8 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Order · Summary</span>
          </div>

          <!-- Product line — 不用 emoji，用 hairline 金色方形 -->
          <div class="flex gap-3 mb-4 pb-4 border-b border-gold/10">
            <div class="w-16 h-16 bg-ink-900 border border-gold/30 flex items-center justify-center" style="border-radius: 2px;">
              <span class="font-serif text-gold text-2xl">T</span>
            </div>
            <div class="flex-1">
              <div class="text-sm font-serif text-ink-900">{{ product?.title }}</div>
              <div class="text-[11px] uppercase tracking-lux text-sand font-sans mt-1">£{{ product?.unit_price }} × {{ qty }}</div>
            </div>
            <div class="text-sm font-serif text-ink-900">£{{ ((product?.unit_price||0)*qty).toFixed(2) }}</div>
          </div>

          <div class="space-y-2 text-sm">
            <div class="flex justify-between text-sand font-serif"><span>Subtotal</span><span>£{{ ((product?.unit_price||0)*qty).toFixed(2) }}</span></div>
            <div class="flex justify-between text-sand font-serif"><span>Shipping · UK</span><span>£{{ product?.shipping_cost }}</span></div>
          </div>

          <div class="flex justify-between pt-4 mt-4 border-t border-gold/10 font-serif text-xl text-ink-900">
            <span>Total</span><span>£{{ total().toFixed(2) }}</span>
          </div>

          <button :disabled="loading" @click="place"
            class="w-full mt-6 py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans disabled:opacity-50 disabled:cursor-not-allowed"
            style="border-radius: 2px;">
            {{ loading ? 'Redirecting · · ·' : 'Confirm · Order · With · 2Checkout' }}
          </button>

          <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-5 leading-loose font-sans">
            256-bit SSL &middot; 2Checkout &middot; PayPal<br>
            We never store your card details
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
