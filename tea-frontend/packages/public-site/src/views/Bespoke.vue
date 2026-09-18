<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { api } from '@/api/client'

const step = ref(1)
const submitting = ref(false)
const success = ref(false)
const errorMsg = ref('')

const data = ref({
  tea_type: 'raw_puer',
  teaGarden: '',
  roast: 'light',
  shape: 'cake',
  weight: 357,
  quantity: 10,
  inner_packaging: 'cotton_paper',
  outer_packaging: 'wooden_box',
  message: '',
  contact_email: '',
  contact_name: '',
})

const teaGardens = ['云南省临沧市临翔区邦东乡曼岗村茶园', '云南省西双版纳州勐海县布朗山乡班章村茶园', '云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园', '云南省西双版纳州勐海县勐海镇贺开村茶园', '云南省临沧市临翔区南美乡多依村茶园', '云南省普洱市澜沧拉祜族自治县惠民镇翁基村茶园']
const teaTypes = [{v:'raw_puer', l:'生普洱 Raw Pu\'er'},{v:'ripe_puer', l:'熟普洱 Ripe Pu\'er'}]
const roasts = [{v:'light', l:'轻火 Light Roast'},{v:'medium', l:'中火 Medium Roast'},{v:'heavy', l:'重火 Heavy Roast'}]
const shapes = [{v:'cake', l:'饼 Tea Cake'},{v:'brick', l:'砖 Tea Brick'},{v:'tuo', l:'沱 Tuo Cha'},{v:'golden_brick', l:'金砖 Golden Brick'}]
const packagings = [{v:'cotton_paper', l:'棉纸 Cotton Paper'},{v:'bamboo', l:'竹筒 Bamboo Tube'},{v:'wooden_box', l:'木箱 Wooden Box + £20'},{v:'silk', l:'丝绸 Silk Wrap + £35'}]

async function submit() {
  errorMsg.value = ''
  submitting.value = true
  try {
    const token = localStorage.getItem('user_token')
    const payload: any = {
      title: `Bespoke ${data.value.tea_type === 'raw_puer' ? 'Raw' : 'Ripe'} Pu'er`,
      tea_type: data.value.tea_type,
      tea_shape: data.value.shape,
      unit_price: 60 + (data.value.weight / 100),
      quantity: data.value.quantity,
      shipping_cost: 120,
      inner_packaging: data.value.inner_packaging,
      outer_packaging: data.value.outer_packaging,
      tea_garden_location: data.value.teaGarden || 'Master Selection',
      master_name: 'Master Selection',
      raw_tea_source: data.value.teaGarden || 'Master Selection',
      custom_requirement: data.value.message,
    }
    if (token) {
      await api.post('/custom-products', payload)
    } else {
      // 未登录：存到 localStorage，提示稍后登录后确认提交
      const queue = JSON.parse(localStorage.getItem('bespoke_queue') || '[]')
      queue.push({ ...payload, queued_at: new Date().toISOString() })
      localStorage.setItem('bespoke_queue', JSON.stringify(queue))
    }
    success.value = true
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.error || e.message || 'Submission failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="pt-20">
    <section class="bg-tea-900 text-tea-50 py-20 md:py-28 relative overflow-hidden">
      <div class="absolute inset-0 opacity-30" style="background-image: radial-gradient(circle at 30% 40%, rgba(212,160,79,0.15), transparent 50%);"></div>
      <div class="max-w-4xl mx-auto px-6 text-center relative">
        <h1 class="font-display text-5xl md:text-6xl lg:text-7xl mb-6 font-semibold tracking-tight">Your Tea, <span class="italic">Your Way</span>.</h1>
        <p class="text-lg md:text-xl text-tea-200 max-w-2xl mx-auto">Tell us what you want. Our master tea roasters in Yunnan will craft it within 45 days — and you'll watch every step live.</p>
      </div>
    </section>

    <section class="py-16 md:py-20 bg-tea-50">
      <div class="max-w-2xl mx-auto px-6">
        <!-- Success state -->
        <div v-if="success" class="bg-white rounded-2xl border border-tea-100 p-10 text-center shadow-sm">
          <div class="text-5xl mb-4">🍵</div>
          <h3 class="font-display text-2xl text-tea-900 mb-3">Bespoke request received</h3>
          <p class="text-tea-700 mb-2">Our tea advisor will contact you within 24 hours with pricing and next steps.</p>
          <p class="text-sm text-tea-500">Token will appear in your account once confirmed.</p>
          <RouterLink to="/" class="inline-block mt-6 px-8 py-3 bg-tea-800 text-white rounded-full hover:bg-tea-900 transition">Back to Home</RouterLink>
        </div>

        <template v-else>
          <!-- Steps -->
          <div class="flex items-center justify-center gap-2 mb-10 flex-wrap">
            <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium shrink-0', step>=1?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">1</span>
            <span class="text-xs text-tea-500 shrink-0">Tea Type</span>
            <span class="w-8 md:w-12 h-px bg-tea-200 shrink-0"></span>
            <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium shrink-0', step>=2?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">2</span>
            <span class="text-xs text-tea-500 shrink-0">Blend</span>
            <span class="w-8 md:w-12 h-px bg-tea-200 shrink-0"></span>
            <span :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium shrink-0', step>=3?'bg-tea-700 text-white':'bg-tea-200 text-tea-500']">3</span>
            <span class="text-xs text-tea-500 shrink-0">Packaging</span>
          </div>

          <form class="bg-white rounded-2xl shadow-sm border border-tea-100 p-6 md:p-8" @submit.prevent="step<3?step++:submit()">
            <div v-if="errorMsg" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-xl text-sm">{{ errorMsg }}</div>

            <template v-if="step===1">
              <label class="block text-sm font-medium text-tea-900 mb-2">Tea Type</label>
              <div class="grid grid-cols-2 gap-3 mb-6">
                <button v-for="t in teaTypes" :key="t.v" type="button" @click="data.tea_type = t.v"
                  :class="['p-4 rounded-xl border text-left transition', data.tea_type===t.v ? 'border-tea-600 bg-tea-50' : 'border-tea-200 hover:border-tea-400']">
                  {{ t.l }}
                </button>
              </div>
              <label class="block text-sm font-medium text-tea-900 mb-2">Preferred Tea Garden (Village)</label>
              <select v-model="data.teaGarden" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none bg-white">
                <option value="">Not sure — let the master choose</option>
                <option v-for="m in teaGardens" :key="m" :value="m">{{ m }}</option>
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
              <select v-model="data.inner_packaging" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-4 focus:border-tea-600 focus:outline-none bg-white">
                <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
              </select>
              <label class="block text-sm font-medium text-tea-900 mb-2">Outer Packaging</label>
              <select v-model="data.outer_packaging" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none bg-white">
                <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
              </select>
              <label class="block text-sm font-medium text-tea-900 mb-2">Any Special Notes?</label>
              <textarea v-model="data.message" rows="3" placeholder="e.g. For a 60th birthday, family crest, gift message..." class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-6 focus:border-tea-600 focus:outline-none resize-none"></textarea>
              <div class="flex gap-3">
                <button type="button" @click="step--" class="flex-1 py-3 border border-tea-300 text-tea-700 rounded-xl">← Back</button>
                <button type="submit" :disabled="submitting"
                  class="flex-1 py-3 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition disabled:opacity-50 disabled:cursor-not-allowed">
                  {{ submitting ? 'Submitting...' : 'Submit Bespoke Request →' }}
                </button>
              </div>
            </template>
          </form>

          <p class="text-xs text-tea-500 text-center mt-4 leading-relaxed">
            <strong>Note:</strong> Bespoke production takes up to 45 days. Your advisor will confirm final pricing and lead time before roasting begins. Actual tea may differ slightly from the sample description shown here.
          </p>

          <div class="mt-8 text-center text-sm text-tea-500">
            Already have a bespoke token? <RouterLink to="/bespoke/demo-1" class="text-tea-700 font-medium hover:underline">View existing product →</RouterLink>
          </div>
        </template>
      </div>
    </section>
  </div>
</template>
