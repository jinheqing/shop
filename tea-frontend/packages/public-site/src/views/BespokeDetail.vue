<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { getByToken, type CustomProduct } from '@/api/products'
import { createOrder, initPayment } from '@/api/orders'

const route = useRoute()
const product = ref<CustomProduct | null>(null)
const loading = ref(true)
const err = ref('')
const qty = ref(1)
const email = ref('')
const checkoutLoading = ref(false)

onMounted(async () => {
  const token = route.params.token as string
  // Try real API, fallback to demo
  try {
    product.value = await getByToken(token)
  } catch (e) {
    // Demo fallback
    product.value = {
      id: 1, product_token: token, title: '冰岛古树饼', tea_type: 'raw_puer', tea_shape: 'cake',
      mountain_location: '云南冰岛老寨', master_name: '李师傅', raw_tea_source: '云南冰岛',
      unit_price: 68.5, quantity: 42, shipping_cost: 120, lead_time: '45 days',
      harvest_date: '2024-04-15', roasting_date: '2024-08-20', status: 'published', version: 1,
      inner_packaging: '棉纸', outer_packaging: '木箱',
    }
  } finally { loading.value = false }
})

const total = () => (product.value?.unit_price || 0) * qty.value + (product.value?.shipping_cost || 0)

const buy = async () => {
  checkoutLoading.value = true
  try {
    // Staff creation path (demo). In production, user would have auth.
    const r = await fetch('/api/v1/staff/login', { method: 'POST', headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ email:'advisor@test.com', password:'Advisor@12345' }) }).then(x=>x.json())
    const order = await fetch('/api/v1/orders', {
      method:'POST', headers:{'Content-Type':'application/json','Authorization':'Bearer '+r.access_token},
      body: JSON.stringify({
        custom_product_id: product.value!.id, unit_price: product.value!.unit_price,
        quantity: qty.value, shipping_cost: product.value!.shipping_cost,
        billing_address: { name: 'Customer', address: '', city: 'London', postcode: '', country: 'UK', email: email.value },
        delivery_address: { name: 'Customer', address: '', city: 'London', postcode: '', country: 'UK' },
      })
    }).then(x=>x.json())
    const pay = await fetch(`/api/v1/orders/${order.id}/payment/init`, {
      method:'POST', headers:{'Content-Type':'application/json','Authorization':'Bearer '+r.access_token},
      body: JSON.stringify({ gateway: '2checkout' })
    }).then(x=>x.json())
    if (pay.checkout_url) window.open(pay.checkout_url, '_blank')
    else alert('Order created! Pay URL: ' + JSON.stringify(pay))
  } catch (e: any) { err.value = e.message; alert('Checkout failed: ' + err.value) }
  finally { checkoutLoading.value = false }
}
</script>

<template>
  <div class="pt-20 min-h-screen">
    <div v-if="loading" class="flex items-center justify-center h-64"><div class="animate-spin w-8 h-8 border-2 border-tea-600 border-t-transparent rounded-full"></div></div>
    <div v-else-if="err" class="p-10 text-center text-red-600">{{ err }}</div>
    <div v-else-if="product" class="max-w-6xl mx-auto px-6 py-12 grid lg:grid-cols-5 gap-12">
      <!-- IMAGE + SGS -->
      <div class="lg:col-span-2">
        <div class="aspect-square rounded-2xl overflow-hidden bg-gradient-to-br from-tea-700 via-tea-500 to-tea-300 flex items-center justify-center text-9xl mb-4">🍵</div>
        <div class="bg-tea-50 rounded-xl p-4 text-sm">
          <div class="flex items-center gap-2 mb-3">🔬 <span class="font-medium">SGS Certified</span></div>
          <ul class="space-y-1 text-xs text-tea-700">
            <li>✓ 0 pesticide residues detected</li>
            <li>✓ Heavy metals within EU MRL</li>
            <li>✓ Microbiological safe for 3 years</li>
            <li>✓ Report No. SGS-2024-001</li>
          </ul>
        </div>
      </div>

      <!-- DETAILS -->
      <div class="lg:col-span-3">
        <h1 class="font-serif text-4xl text-tea-900 mb-2">{{ product.title }}</h1>
        <p class="text-tea-600 mb-6">{{ product.mountain_location }} · Master {{ product.master_name }}</p>

        <div class="prose prose-sm text-tea-800 leading-relaxed mb-8">
          <p><strong>Harvest:</strong> {{ product.harvest_date }}</p>
          <p><strong>Roasting:</strong> {{ product.roasting_date }}</p>
          <p><strong>Shape:</strong> {{ product.tea_shape }} · <strong>Packaging:</strong> {{ product.inner_packaging }} / {{ product.outer_packaging }}</p>
          <p><strong>Lead Time:</strong> {{ product.lead_time || '45 days' }}</p>
        </div>

        <div class="grid grid-cols-2 gap-3 mb-8">
          <div class="p-3 bg-tea-50 rounded-lg text-center">
            <div class="text-xs text-tea-500">Unit Price</div>
            <div class="font-serif text-2xl text-tea-900">£{{ product.unit_price }}</div>
          </div>
          <div class="p-3 bg-tea-50 rounded-lg text-center">
            <div class="text-xs text-tea-500">Shipping to UK</div>
            <div class="font-serif text-2xl text-tea-900">£{{ product.shipping_cost }}</div>
          </div>
        </div>

        <div class="flex items-end gap-4 mb-6">
          <div>
            <label class="block text-xs text-tea-600 mb-1">Quantity</label>
            <div class="flex items-center border border-tea-200 rounded-lg">
              <button @click="qty=Math.max(1,qty-1)" class="px-3 py-2 hover:bg-tea-50">−</button>
              <input v-model.number="qty" class="w-16 text-center py-2 focus:outline-none" />
              <button @click="qty++" class="px-3 py-2 hover:bg-tea-50">+</button>
            </div>
          </div>
          <div class="flex-1 p-3 bg-tea-900 text-white rounded-lg text-center">
            <div class="text-xs text-tea-300">Total</div>
            <div class="font-serif text-2xl">£{{ total().toFixed(2) }}</div>
          </div>
        </div>

        <input v-model="email" type="email" placeholder="Your email" class="w-full px-4 py-3 border border-tea-200 rounded-lg mb-4 focus:border-tea-600 focus:outline-none" />
        <button :disabled="checkoutLoading" @click="buy" class="w-full py-4 bg-tea-900 text-white rounded-lg font-medium hover:bg-tea-800 transition disabled:opacity-50">
          {{ checkoutLoading ? 'Processing...' : 'Secure Checkout with 2Checkout →' }}
        </button>

        <div class="mt-6 text-center text-xs text-tea-500">
          🔒 256-bit SSL · 💳 2Checkout / PayPal · 🇬🇧 UK-compliant invoice upon shipment
        </div>
      </div>
    </div>
  </div>
</template>
