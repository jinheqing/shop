<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import axios from 'axios'

const forms = ref<HTMLFormElement | null>(null)
const step = ref(1)
const data = ref({
  tea_type: 'raw_puer',
  mountain: '',
  roast: 'light',
  shape: 'cake',
  weight: 357,
  quantity: 10,
  inner_packaging: 'cotton_paper',
  outer_packaging: 'wooden_box',
  message: '',
})

const mountains = ['冰岛老寨', '班章村', '景迈山', '老班章', '南糯山', '忙肺村']
const teaTypes = [{v:'raw_puer', l:'生普洱 Raw Pu\'er'},{v:'ripe_puer', l:'熟普洱 Ripe Pu\'er'}]
const roasts = [{v:'light', l:'轻火 Light Roast'},{v:'medium', l:'中火 Medium Roast'},{v:'heavy', l:'重火 Heavy Roast'}]
const shapes = [{v:'cake', l:'饼 Tea Cake'},{v:'brick', l:'砖 Tea Brick'},{v:'tuo', l:'沱 Tuo Cha'},{v:'golden_brick', l:'金砖 Golden Brick'}]
const packagings = [{v:'cotton_paper', l:'棉纸 Cotton Paper'},{v:'bamboo', l:'竹筒 Bamboo Tube'},{v:'wooden_box', l:'木箱 Wooden Box + £20'},{v:'silk', l:'丝绸 Silk Wrap + £35'}]

const submit = async () => {
  try {
    const r = await axios.post('/api/v1/staff/login', { email: 'advisor@test.com', password: 'Advisor@12345' })
    await axios.post('/api/v1/custom-products', {
      title: `定制茶 ${data.value.tea_type === 'raw_puer' ? '生普' : '熟普'}`,
      tea_type: data.value.tea_type,
      tea_shape: data.value.shape,
      unit_price: 60 + (data.value.weight/100),
      quantity: data.value.quantity,
      shipping_cost: 120,
      inner_packaging: data.value.inner_packaging,
      outer_packaging: data.value.outer_packaging,
      mountain_location: data.value.mountain,
      master_name: '李师傅',
      raw_tea_source: data.value.mountain,
      custom_requirement: data.value.message,
    }, { headers: { Authorization: `Bearer ${r.data.access_token}` } })
    alert('✅ Bespoke request submitted! Your advisor will contact you within 24h.')
  } catch (e: any) {
    alert('Submission failed: ' + (e?.response?.data?.message || e.message))
  }
}
</script>

<template>
  <div class="pt-20">
    <section class="bg-tea-900 text-tea-50 py-24">
      <div class="max-w-4xl mx-auto px-6 text-center">
        <h1 class="font-serif text-5xl md:text-6xl mb-6">Your Tea, Your Way.</h1>
        <p class="text-lg text-tea-200">Tell us what you want. Our master tea roasters in Yunnan will craft it within 45 days — and you'll watch every step live.</p>
      </div>
    </section>

    <section class="py-16 bg-tea-50">
      <div class="max-w-2xl mx-auto px-6">
        <!-- Steps -->
        <div class="flex items-center justify-center gap-2 mb-10">
          <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium', step>=1?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">1</span>
          <span class="text-xs text-tea-500">Tea Type</span>
          <span class="w-12 h-px bg-tea-200"></span>
          <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium', step>=2?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">2</span>
          <span class="text-xs text-tea-500">Blend</span>
          <span class="w-12 h-px bg-tea-200"></span>
          <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium', step>=3?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">3</span>
          <span class="text-xs text-tea-500">Packaging</span>
        </div>

        <form class="bg-white rounded-2xl shadow-sm border border-tea-100 p-8" @submit.prevent="step<3?step++:submit()">
          <template v-if="step===1">
            <label class="block text-sm font-medium text-tea-900 mb-2">Tea Type</label>
            <div class="grid grid-cols-2 gap-3 mb-6">
              <button v-for="t in teaTypes" :key="t.v" type="button" @click="data.tea_type = t.v"
                :class="['p-4 rounded-xl border text-left transition', data.tea_type===t.v ? 'border-tea-600 bg-tea-50' : 'border-tea-200 hover:border-tea-400']">
                {{ t.l }}
              </button>
            </div>
            <label class="block text-sm font-medium text-tea-900 mb-2">Preferred Mountain</label>
            <select v-model="data.mountain" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none">
              <option value="">Not sure — let the master choose</option>
              <option v-for="m in mountains" :key="m" :value="m">{{ m }}</option>
            </select>
            <button type="submit" class="w-full py-3 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition">Continue →</button>
          </template>
          <template v-if="step===2">
            <label class="block text-sm font-medium text-tea-900 mb-2">Roast Level</label>
            <div class="grid grid-cols-3 gap-3 mb-6">
              <button v-for="r in roasts" :key="r.v" type="button" @click="data.roast = r.v"
                :class="['p-3 rounded-xl border text-sm text-left', data.roast===r.v ? 'border-tea-600 bg-tea-50' : 'border-tea-200']">{{ r.l }}</button>
            </div>
            <label class="block text-sm font-medium text-tea-900 mb-2">Shape</label>
            <div class="grid grid-cols-2 gap-3 mb-6">
              <button v-for="s in shapes" :key="s.v" type="button" @click="data.shape = s.v"
                :class="['p-3 rounded-xl border text-left', data.shape===s.v ? 'border-tea-600 bg-tea-50' : 'border-tea-200']">{{ s.l }}</button>
            </div>
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div>
                <label class="block text-sm text-tea-900 mb-2">Weight (g)</label>
                <input v-model.number="data.weight" type="number" class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
              </div>
              <div>
                <label class="block text-sm text-tea-900 mb-2">Quantity</label>
                <input v-model.number="data.quantity" type="number" min="1" class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
              </div>
            </div>
            <div class="flex gap-3">
              <button type="button" @click="step--" class="flex-1 py-3 border border-tea-300 text-tea-700 rounded-xl">← Back</button>
              <button type="submit" class="flex-1 py-3 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800">Continue →</button>
            </div>
          </template>
          <template v-if="step===3">
            <label class="block text-sm font-medium text-tea-900 mb-2">Inner Packaging</label>
            <select v-model="data.inner_packaging" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-4 focus:border-tea-600 focus:outline-none">
              <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
            </select>
            <label class="block text-sm font-medium text-tea-900 mb-2">Outer Packaging</label>
            <select v-model="data.outer_packaging" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none">
              <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
            </select>
            <label class="block text-sm font-medium text-tea-900 mb-2">Any Special Notes?</label>
            <textarea v-model="data.message" rows="3" placeholder="e.g. For a 60th birthday, family crest, gift message..." class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none"></textarea>
            <div class="flex gap-3">
              <button type="button" @click="step--" class="flex-1 py-3 border border-tea-300 text-tea-700 rounded-xl">← Back</button>
              <button type="submit" class="flex-1 py-3 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800">Submit Bespoke Request →</button>
            </div>
          </template>
        </form>

        <div class="mt-8 text-center text-sm text-tea-500">
          Already have a bespoke token? <RouterLink to="/bespoke/demo-1" class="text-tea-700 font-medium hover:underline">View existing product →</RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>
