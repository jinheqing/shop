<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getByToken, type CustomProduct } from '@/api/products'
import { createOrder, initPayment } from '@/api/orders'

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

onMounted(async () => {
  try { product.value = await getByToken(route.params.token as string) }
  catch { product.value = { id: 1, unit_price: 68.5, shipping_cost: 120, title: 'Demo Product', tea_type: 'raw_puer', tea_shape: 'cake', status: 'published', version: 1 } as any }
})

const total = () => (product.value?.unit_price || 0) * qty.value + (product.value?.shipping_cost || 0)

const place = async () => {
  loading.value = true
  try {
    // 使用用户自己的 user_token 下单（不伪造 staff 身份）
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
      alert('Order failed: ' + (order.message || JSON.stringify(order)))
      return
    }
    const pay = await fetch(`/api/v1/orders/${order.id}/payment/init`, {
      method:'POST', headers:{'Content-Type':'application/json','Authorization':'Bearer '+userToken},
      body: JSON.stringify({ gateway: '2checkout' })
    }).then(x=>x.json())
    if (pay.checkout_url) window.location.href = pay.checkout_url
    else router.push('/')
  } catch (e: any) { alert('Error: ' + e.message) }
  finally { loading.value = false }
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-tea-50">
    <div class="max-w-4xl mx-auto px-6 py-12">
      <h1 class="font-serif text-4xl text-tea-900 mb-8">Secure Checkout</h1>
      <div class="grid md:grid-cols-3 gap-8">
        <div class="md:col-span-2 space-y-6">
          <div class="bg-white rounded-xl p-6 border border-tea-100">
            <h3 class="font-medium text-tea-900 mb-4">Your Details</h3>
            <div class="grid grid-cols-2 gap-4">
              <input v-model="form.name" placeholder="Full name" class="px-4 py-3 border border-tea-200 rounded-lg focus:border-tea-600 focus:outline-none" />
              <input v-model="form.email" type="email" placeholder="Email" class="px-4 py-3 border border-tea-200 rounded-lg focus:border-tea-600 focus:outline-none" />
            </div>
          </div>

          <div class="bg-white rounded-xl p-6 border border-tea-100">
            <h3 class="font-medium text-tea-900 mb-4">Billing Address</h3>
            <input v-model="form.billing.address" placeholder="Street address" class="w-full px-4 py-3 border border-tea-200 rounded-lg mb-3 focus:border-tea-600 focus:outline-none" />
            <div class="grid grid-cols-3 gap-3">
              <input v-model="form.billing.city" placeholder="City" class="px-4 py-3 border border-tea-200 rounded-lg focus:border-tea-600 focus:outline-none" />
              <input v-model="form.billing.postcode" placeholder="Postcode" class="px-4 py-3 border border-tea-200 rounded-lg focus:border-tea-600 focus:outline-none" />
              <input v-model="form.billing.country" placeholder="Country" :value="'UK'" class="px-4 py-3 border border-tea-200 rounded-lg focus:border-tea-600 focus:outline-none" />
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl p-6 border border-tea-100 h-fit sticky top-24">
          <h3 class="font-medium text-tea-900 mb-4">Order Summary</h3>
          <div class="flex gap-3 mb-4 pb-4 border-b border-tea-100">
            <div class="w-16 h-16 bg-gradient-to-br from-tea-700 to-tea-400 rounded-lg flex items-center justify-center text-2xl">🍵</div>
            <div class="flex-1">
              <div class="text-sm font-medium text-tea-900">{{ product?.title }}</div>
              <div class="text-xs text-tea-500">£{{ product?.unit_price }} × {{ qty }}</div>
            </div>
            <div class="text-sm font-medium">£{{ ((product?.unit_price||0)*qty).toFixed(2) }}</div>
          </div>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between text-tea-600"><span>Subtotal</span><span>£{{ ((product?.unit_price||0)*qty).toFixed(2) }}</span></div>
            <div class="flex justify-between text-tea-600"><span>Shipping (UK)</span><span>£{{ product?.shipping_cost }}</span></div>
          </div>
          <div class="flex justify-between pt-4 mt-4 border-t border-tea-100 font-serif text-xl">
            <span>Total</span><span>£{{ total().toFixed(2) }}</span>
          </div>
          <button :disabled="loading" @click="place" class="w-full mt-6 py-4 bg-tea-900 text-white rounded-lg font-medium hover:bg-tea-800 transition disabled:opacity-50">
            {{ loading ? 'Redirecting to 2Checkout...' : 'Pay Securely →' }}
          </button>
          <p class="text-xs text-center text-tea-500 mt-4">🔒 256-bit SSL · 2Checkout / PayPal · We never store your card details</p>
        </div>
      </div>
    </div>
  </div>
</template>
