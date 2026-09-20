<script setup lang="ts">
// ============================================================
// Bespoke.vue — 叙述式定制表单（Old Money Edition）
// ============================================================
// 设计参考：Hermès Bespoke / Rolls-Royce Private Office
//   - 不分步圆点（不要 SaaS-style 1-2-3 progress）
//   - 单页叙述，serif CAPS 标签，hairline 分隔
//   - off-black 主舞台 + 象牙白表单
//   - 绝不用 emoji 当 icon
// ============================================================

import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { api } from '@/api/client'

const router = useRouter()

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

const teaGardens = ['云南省 · 临沧市 · 云雾茶区', '云南省 · 西双版纳州 · 布朗山茶区', '云南省 · 普洱市 · 澜沧茶区', '云南省 · 西双版纳州 · 勐海茶区', '云南省 · 临沧市 · 雪山茶区', '云南省 · 普洱市 · 澜沧茶区']
const teaTypes = [{v:'raw_puer', l:'Raw Pu\'er'}, {v:'ripe_puer', l:'Ripe Pu\'er'}]
const roasts = [{v:'light', l:'Light Roast'}, {v:'medium', l:'Medium Roast'}, {v:'heavy', l:'Heavy Roast'}]
const shapes = [{v:'cake', l:'Tea Cake'}, {v:'brick', l:'Tea Brick'}, {v:'tuo', l:'Tuo Cha'}, {v:'golden_brick', l:'Golden Brick'}]
const packagings = [{v:'cotton_paper', l:'Cotton Paper'}, {v:'bamboo', l:'Bamboo Tube'}, {v:'wooden_box', l:'Wooden Box (+£20)'}, {v:'silk', l:'Silk Wrap (+£35)'}]

async function submit() {
  errorMsg.value = ''
  submitting.value = true
  try {
    const token = localStorage.getItem('user_token')
    const now = new Date()
    const pad = (n: number) => String(n).padStart(2, '0')
    const today = `${now.getFullYear()}-${pad(now.getMonth()+1)}-${pad(now.getDate())}`
    const basePayload: any = {
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
      const payload = {
        ...basePayload,
        qr_code_position: 'bottom',
        lead_time: '45 days from confirmation',
        harvest_date: today,
        roasting_date: today,
        storage_location: 'London (temporary)',
      }
      const resp: any = await api.post('/custom-products', payload)
      const quoteToken = resp?.token || resp?.product?.product_token || resp?.data?.token
      if (quoteToken) { router.push(`/quote/${quoteToken}`); return }
    } else {
      const resp: any = await api.post('/public/custom-products', basePayload)
      const quoteToken = resp?.token || resp?.data?.token
      if (quoteToken) { router.push(`/quote/${quoteToken}`); return }
    }
    success.value = true
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.message || e?.response?.data?.error || e.message || 'Submission failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="pt-20">
    <!-- ============================================================
         Hero — off-black 主舞台，与首页 / TeaGardens 同源
         ============================================================ -->
    <section class="bg-ink-900 text-ivory-100 py-20 md:py-28 relative overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.08]"
           style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 70%, #C5A572 0%, transparent 50%);"></div>
      <div class="max-w-4xl mx-auto px-6 text-center relative">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">By · Appointment · 45 Days</span>
          <span class="h-px w-10 bg-gold/60"></span>
        </div>
        <h1 class="font-serif text-5xl md:text-6xl lg:text-7xl mb-7 tracking-tight leading-[1.05]">
          Your · Tea,<br><span class="italic text-gold">Crafted</span> · For · You.
        </h1>
        <p class="text-lg md:text-xl text-sand max-w-2xl mx-auto leading-relaxed font-serif">
          A single advisor oversees your commission from first letter to final brew.
          One master in Yunnan crafts your batch by hand, in one season.
        </p>
      </div>
    </section>

    <!-- ============================================================
         Main Form Section — 象牙白面板，叙述式不分步
         ============================================================ -->
    <section class="py-16 md:py-20 bg-ivory-100">
      <div class="max-w-2xl mx-auto px-6">

        <!-- Success state — 极克制的感谢段，无 emoji -->
        <div v-if="success" class="bg-white border border-gold/20 p-10 md:p-12 text-center" style="border-radius: 2px;">
          <div class="flex items-center gap-3 justify-center mb-6">
            <span class="h-px w-10 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Received</span>
            <span class="h-px w-10 bg-gold/50"></span>
          </div>
          <h3 class="font-serif text-3xl md:text-4xl text-ink-900 mb-4">Thank you for your enquiry.</h3>
          <p class="text-sand mb-2 leading-relaxed font-serif">An advisor will write to you within twenty-four hours, by email or private message.</p>
          <p class="text-[11px] uppercase tracking-lux text-sand font-sans mb-8">Commission · Pending · Confirmation</p>
          <RouterLink to="/" class="inline-block px-8 py-3 border border-gold/40 text-gold hover:bg-gold/5 transition text-[11px] uppercase tracking-lux font-sans">Back · To · The · House</RouterLink>
        </div>

        <template v-else>
          <!-- Section title -->
          <div class="mb-10 md:mb-12 text-center">
            <div class="flex items-center gap-3 justify-center mb-5">
              <span class="h-px w-10 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">The · Commission</span>
              <span class="h-px w-10 bg-gold/50"></span>
            </div>
            <h2 class="font-serif text-3xl md:text-4xl text-ink-900 mb-3">Tell · Us · Of · Your · Wishes.</h2>
            <p class="text-sand text-sm leading-relaxed max-w-lg mx-auto font-serif">
              Complete this letter to your advisor. We shall respond within twenty-four hours,
              with a quotation and the name of the master assigned to your batch.
            </p>
          </div>

          <form class="bg-white border border-gold/15 p-8 md:p-10" style="border-radius: 2px;" @submit.prevent="submit()">

            <div v-if="errorMsg" class="mb-6 p-4 bg-ink-900 text-ivory-100 text-sm font-serif" style="border-radius: 2px;">
              <span class="text-gold text-[10px] uppercase tracking-lux font-sans">Error</span><br>
              {{ errorMsg }}
            </div>

            <!-- Tea Type -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Tea · Type</label>
              <div class="grid grid-cols-2 gap-3">
                <button v-for="t in teaTypes" :key="t.v" type="button" @click="data.tea_type = t.v"
                  :class="['p-4 border transition-duration-lux text-left font-serif', data.tea_type===t.v ? 'border-gold bg-ivory-100 text-ink-900' : 'border-gold/20 text-sand hover:border-gold/50']"
                  style="border-radius: 2px;">{{ t.l }}</button>
              </div>
            </div>

            <!-- Tea Garden -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Preferred · Tea · Garden · (Optional)</label>
              <select v-model="data.teaGarden"
                class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;">
                <option value="">— Master Selection —</option>
                <option v-for="m in teaGardens" :key="m" :value="m">{{ m }}</option>
              </select>
            </div>

            <!-- Hairline divider -->
            <div class="hairline-gold my-8"></div>

            <!-- Roast -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Roast · Level</label>
              <div class="grid grid-cols-3 gap-3">
                <button v-for="r in roasts" :key="r.v" type="button" @click="data.roast = r.v"
                  :class="['p-3 border transition-duration-lux text-sm font-serif text-left', data.roast===r.v ? 'border-gold bg-ivory-100 text-ink-900' : 'border-gold/20 text-sand hover:border-gold/50']"
                  style="border-radius: 2px;">{{ r.l }}</button>
              </div>
            </div>

            <!-- Shape -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Shape</label>
              <div class="grid grid-cols-2 gap-3">
                <button v-for="s in shapes" :key="s.v" type="button" @click="data.shape = s.v"
                  :class="['p-3 border transition-duration-lux font-serif text-left', data.shape===s.v ? 'border-gold bg-ivory-100 text-ink-900' : 'border-gold/20 text-sand hover:border-gold/50']"
                  style="border-radius: 2px;">{{ s.l }}</button>
              </div>
            </div>

            <!-- Weight + Quantity -->
            <div class="grid grid-cols-2 gap-4 mb-8">
              <div>
                <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Weight · (g)</label>
                <input v-model.number="data.weight" type="number"
                  class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                  style="border-radius: 2px;" />
              </div>
              <div>
                <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Quantity</label>
                <input v-model.number="data.quantity" type="number" min="1"
                  class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                  style="border-radius: 2px;" />
              </div>
            </div>

            <!-- Hairline divider -->
            <div class="hairline-gold my-8"></div>

            <!-- Packaging -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Inner · Packaging</label>
              <select v-model="data.inner_packaging"
                class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900 mb-4"
                style="border-radius: 2px;">
                <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
              </select>

              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Outer · Packaging</label>
              <select v-model="data.outer_packaging"
                class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
                style="border-radius: 2px;">
                <option v-for="p in packagings" :key="p.v" :value="p.v">{{ p.l }}</option>
              </select>
            </div>

            <!-- Hairline divider -->
            <div class="hairline-gold my-8"></div>

            <!-- Special Notes -->
            <div class="mb-8">
              <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Notes · To · Your · Advisor</label>
              <textarea v-model="data.message" rows="4"
                placeholder="e.g. A commission for a sixtieth birthday. Family crest on the box. Brief message enclosed..."
                class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none resize-none font-serif text-ink-900"
                style="border-radius: 2px;"></textarea>
            </div>

            <!-- Submit — 单一主行动，二级动作放小字链接 -->
            <button type="submit" :disabled="submitting"
              class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans disabled:opacity-50 disabled:cursor-not-allowed"
              style="border-radius: 2px;">
              {{ submitting ? 'Submitting · · ·' : 'Send · To · Your · Advisor' }}
            </button>
          </form>

          <!-- 极细 disclaimer -->
          <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-8 leading-loose font-sans">
            Bespoke production takes up to 45 days &middot; Final quotation by advisor before roasting<br>
            Tea may differ slightly from sample description &middot; Not a medicinal product
          </p>

          <!-- 二级动作 — 不显眼 -->
          <div class="mt-6 text-center">
            <RouterLink to="/contact"
              class="text-[11px] uppercase tracking-lux text-gold hover:text-gold-soft transition-duration-lux font-sans">
              Already · Commissioned · Before &middot; View · Your · Orders &rarr;
            </RouterLink>
          </div>
        </template>
      </div>
    </section>
  </div>
</template>
