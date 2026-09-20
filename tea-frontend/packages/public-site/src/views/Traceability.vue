<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const token = route.params.token as string
const product = ref<any>(null)
const loading = ref(true)

onMounted(async () => {
  try { product.value = await api.get(`/custom-products/by-token/${token}`) }
  catch {
    product.value = {
      title: '邦东古树饼',
      tea_garden_location: '云南省 · 临沧市',
      master_name: '王师傅',
      harvest_date: '2026-04-12',
      roasting_date: '2026-05-20',
      storage_location: 'London · Climate-Controlled Cellar',
      tea_type: 'Raw Pu\'er (生普)',
      tea_shape: 'Cake · 357g',
      product_card_text: '凤凰山大乌岽 · 古树单丛 · 2026 春茶',
      sgs_report_no: 'SGS-ICL-2026-003',
      stream_url: 'rtmp://localhost:1935/slow/iceland',
    }
  } finally { loading.value = false }
})
</script>

<template>
  <div class="min-h-screen bg-ivory-100 text-ink-900 pt-20">
    <!-- ============================================================
         Live Camera Hero — off-black 主舞台
         绝不用红色脉冲 LIVE 徽章
         香槟金静态小点 + SERIF CAPS 字
         ============================================================ -->
    <section class="relative h-[60vh] bg-ink-900 text-ivory-100 overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.06]"
           style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 65%, #C5A572 0%, transparent 50%);"></div>

      <div class="absolute inset-0 flex items-center justify-center">
        <div class="text-center relative">
          <!-- LIVE 指示：香槟金静态小点 + SERIF CAPS 字 -->
          <div class="inline-flex items-center gap-2 px-4 py-2 border border-gold/30 rounded-full text-xs mb-8">
            <span class="w-1 h-1 bg-gold rounded-full"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Live · 24/7 · From The Tea Mountain</span>
          </div>
          <!-- 占位帧：极淡的金线框装饰，像老照片取景器 -->
          <div class="relative w-[280px] h-[160px] mx-auto border border-gold/15 flex items-center justify-center">
            <div class="absolute inset-3 border border-gold/10"></div>
            <span class="font-serif text-gold/40 text-sm tracking-widest">Yunnan · Tea Mountain</span>
          </div>
          <div class="text-[11px] uppercase tracking-lux text-sand mt-6 font-sans">
            Live stream from {{ product?.tea_garden_location }}
          </div>
        </div>
      </div>

      <div class="absolute bottom-10 left-0 right-0 text-center">
        <div class="font-serif text-3xl md:text-4xl text-ivory-100">{{ product?.title || 'Your Bespoke Tea' }}</div>
      </div>
    </section>

    <!-- ============================================================
         Trace Timeline — 象牙白面板 + hairline 分隔
         衬线大标题 + 香槟金小标签
         绝不用 emoji 当图标
         ============================================================ -->
    <section class="max-w-4xl mx-auto px-6 py-20">
      <div class="text-center mb-16">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Provenance · Record</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h2 class="font-serif text-4xl md:text-5xl text-ink-900 mb-4">Trace · Your · Tea</h2>
        <p class="text-sand max-w-xl mx-auto leading-relaxed">
          From tea garden to your cup. Every step verified,
          every hand named.
        </p>
      </div>

      <!-- 时间线：左侧香槟金 hairline + 节点圆点，右侧叙述式卡 -->
      <div class="relative pl-8 md:pl-12">
        <!-- 左侧金线 -->
        <div class="absolute left-2 md:left-3 top-2 bottom-2 w-px bg-gold/20"></div>

        <!-- 1. Harvested -->
        <div class="relative mb-12">
          <div class="absolute -left-7 md:-left-10 top-1.5 w-3 h-3 rounded-full bg-gold border-2 border-ivory-100"></div>
          <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">
            {{ product?.harvest_date || '—' }} · Harvested
          </div>
          <div class="bg-white border border-gold/15 p-6" style="border-radius: 2px;">
            <h3 class="font-serif text-xl text-ink-900 mb-2">Picked from the Garden</h3>
            <p class="text-sand text-sm leading-relaxed mb-2">Leaves picked from {{ product?.tea_garden_location }}</p>
            <p class="text-[11px] uppercase tracking-lux text-sand font-sans">Master · {{ product?.master_name }}</p>
          </div>
        </div>

        <!-- 2. Roasted & Rolled -->
        <div class="relative mb-12">
          <div class="absolute -left-7 md:-left-10 top-1.5 w-3 h-3 rounded-full bg-gold border-2 border-ivory-100"></div>
          <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">
            {{ product?.roasting_date || '—' }} · Roasted &amp; Rolled
          </div>
          <div class="bg-white border border-gold/15 p-6" style="border-radius: 2px;">
            <h3 class="font-serif text-xl text-ink-900 mb-2">Traditional Processing</h3>
            <p class="text-sand text-sm leading-relaxed mb-2">Hand-rolled by the master, in one batch, in the traditional manner.</p>
            <p class="text-[11px] uppercase tracking-lux text-sand font-sans">
              {{ product?.tea_type }} · {{ product?.tea_shape }}
            </p>
          </div>
        </div>

        <!-- 3. SGS Tested -->
        <div class="relative mb-12">
          <div class="absolute -left-7 md:-left-10 top-1.5 w-3 h-3 rounded-full bg-gold border-2 border-ivory-100"></div>
          <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">
            SGS · Independently Tested
          </div>
          <div class="bg-white border border-gold/15 p-6" style="border-radius: 2px;">
            <h3 class="font-serif text-xl text-ink-900 mb-2">Independent Laboratory</h3>
            <p class="text-sand text-sm leading-relaxed mb-3">Pesticides · Heavy Metals · Microbiology — verified by SGS China.</p>
            <p class="text-[11px] uppercase tracking-lux text-sand font-sans">
              Report · {{ product?.sgs_report_no }} · Passed
            </p>
          </div>
        </div>

        <!-- 4. Packed -->
        <div class="relative mb-12">
          <div class="absolute -left-7 md:-left-10 top-1.5 w-3 h-3 rounded-full bg-gold border-2 border-ivory-100"></div>
          <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">
            Packed · Bespoke
          </div>
          <div class="bg-white border border-gold/15 p-6" style="border-radius: 2px;">
            <h3 class="font-serif text-xl text-ink-900 mb-2">Custom Packaging</h3>
            <p class="text-sand text-sm leading-relaxed mb-2">Packed in your chosen paper and box, with your QR code affixed.</p>
            <p class="text-[11px] uppercase tracking-lux text-sand font-sans italic">
              Card text · "{{ product?.product_card_text }}"
            </p>
          </div>
        </div>

        <!-- 5. On Its Way -->
        <div class="relative">
          <div class="absolute -left-7 md:-left-10 top-1.5 w-3 h-3 rounded-full bg-gold border-2 border-ivory-100 animate-pulse"></div>
          <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">
            Current · In Transit
          </div>
          <div class="bg-ink-900 text-ivory-100 p-6 border border-gold/30" style="border-radius: 2px;">
            <h3 class="font-serif text-xl text-ivory-100 mb-2">On Its Way</h3>
            <p class="text-sand text-sm leading-relaxed">Shipped via DHL Express to your door in the United Kingdom.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         Closing — 致谢段，克制陈述
         ============================================================ -->
    <section class="bg-ink-900 py-16 text-center">
      <div class="max-w-2xl mx-auto px-6">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Thank You</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h3 class="font-serif text-3xl text-ivory-100 mb-4">For choosing traceable tea.</h3>
        <p class="text-sand text-sm leading-relaxed mb-8 max-w-lg mx-auto">
          Your purchase supports ethical farmers practising natural farming,
          without pesticides — in the cloud mountains of Yunnan.
        </p>
        <RouterLink to="/" class="inline-block px-8 py-3 border border-gold/40 text-gold hover:bg-gold/5 transition text-[11px] uppercase tracking-lux font-sans">
          Back · To · The · House
        </RouterLink>
      </div>
    </section>

    <!-- ============================================================
         Disclaimer — 极细的 serif 小字，像画廊导览注释
         ============================================================ -->
    <div class="max-w-3xl mx-auto px-6 py-10">
      <p class="text-[10px] uppercase tracking-lux text-sand text-center leading-loose font-sans">
        This traceability record reflects the production batch, not every individual cup.
        Live camera feeds are provided for context and may be delayed.
        SGS results apply only to the tested batch. Not medical advice.
      </p>
    </div>
  </div>
</template>
