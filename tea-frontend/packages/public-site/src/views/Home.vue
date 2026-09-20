<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getSlowPresets } from '@/api/live'
import type { CustomProduct } from '@/api/products'
import { getPublished } from '@/api/products'

const presets = ref<any[]>([])
const featured = ref<CustomProduct[]>([])

onMounted(async () => {
  try { presets.value = (await getSlowPresets()) || [] } catch {}
  try { featured.value = (await getPublished()) || [] } catch {}
})
</script>

<template>
  <!-- HERO -->
  <section class="hero-bg min-h-screen flex items-center text-white relative overflow-hidden">
    <!-- subtle gold grain overlay -->
    <div class="absolute inset-0 pointer-events-none opacity-30 mix-blend-overlay"
         style="background-image: radial-gradient(circle at 20% 30%, rgba(212,160,79,0.15), transparent 40%), radial-gradient(circle at 80% 70%, rgba(212,160,79,0.1), transparent 40%);"></div>

    <div class="max-w-7xl mx-auto px-6 py-32 text-center md:text-left relative">
      <span class="inline-block px-5 py-1.5 bg-white/10 backdrop-blur rounded-full text-xs tracking-[0.25em] uppercase mb-8 border border-white/15">
        Single Origin · Direct From Yunnan
      </span>
      <h1 class="font-serif text-5xl md:text-7xl lg:text-8xl leading-[1.05] max-w-4xl mb-8 tracking-tight">
        Pu'er Tea,<br/> <span class="italic">Traceable</span> to the Garden.
      </h1>
      <p class="text-lg md:text-xl text-white/80 max-w-2xl mb-12 leading-relaxed">
        Scan the QR on every bespoke box. See the tea garden, the master who rolled your leaves, and a live camera showing where it grew — 24 hours a day.
      </p>
      <div class="flex flex-col sm:flex-row gap-4 justify-center md:justify-start items-center md:items-start">
        <a href="#collection"
          class="px-9 py-4 bg-white text-ink-900 rounded-full font-medium hover:bg-ivory-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-0.5 text-center tracking-wide">
          Shop The Collection →
        </a>
        <RouterLink to="/bespoke"
          class="px-9 py-4 border border-white/50 rounded-full font-medium hover:bg-white/15 transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 text-center tracking-wide">
          Design Your Own Blend
        </RouterLink>
      </div>
    </div>
  </section>

  <!-- FEATURES — 老钱审美：不用 emoji，用 hairline 金线 + serif CAPS 标签 -->
  <section class="py-24 md:py-28 bg-ivory-100">
    <div class="max-w-7xl mx-auto px-6">
      <div class="grid md:grid-cols-3 gap-12 md:gap-16">
        <div class="text-center">
          <div class="flex items-center gap-3 justify-center mb-6">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Cloud · Garden</span>
            <span class="h-px w-8 bg-gold/60"></span>
          </div>
          <h3 class="font-serif text-xl mb-3 text-ink-900">Sourced · From · The · Mountain</h3>
          <p class="text-sand text-sm leading-relaxed font-serif">Every batch linked to its tea garden village, master name, and harvest date.</p>
        </div>
        <div class="text-center">
          <div class="flex items-center gap-3 justify-center mb-6">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">SGS · Tested</span>
            <span class="h-px w-8 bg-gold/60"></span>
          </div>
          <h3 class="font-serif text-xl mb-3 text-ink-900">Independently · Certified</h3>
          <p class="text-sand text-sm leading-relaxed font-serif">All products independently tested by SGS China for pesticides, heavy metals, microbiology.</p>
        </div>
        <div class="text-center">
          <div class="flex items-center gap-3 justify-center mb-6">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Bespoke</span>
            <span class="h-px w-8 bg-gold/60"></span>
          </div>
          <h3 class="font-serif text-xl mb-3 text-ink-900">Crafted · For · You</h3>
          <p class="text-sand text-sm leading-relaxed font-serif">Choose tea garden, roast level, packaging. Your own tea — from leaf to cup — inside 45 days.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- LIVE STRIP — 第二区块主视觉 -->
  <!-- 设计参考：Dior 2026 秀场、Château Lafite Rothschild、Rolls-Royce 配置器 -->
  <!-- 原则：off-black 舞台 + ivory 面板 + 香槟金细点；零 emoji；零红色脉冲；零暴露内部 URL -->
  <section class="bg-ink-900 py-24 md:py-28 relative overflow-hidden">
    <!-- 极淡的香槟金噪点 overlay — 像老丝绒背景 -->
    <div class="absolute inset-0 pointer-events-none opacity-[0.07]"
         style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 70%, #C5A572 0%, transparent 50%);"></div>

    <div class="max-w-7xl mx-auto px-6 relative">
      <!-- 区块头：serif 大标题 + 香槟金 hairline 分隔 -->
      <div class="mb-16 md:mb-20 max-w-4xl">
        <div class="flex items-center gap-3 mb-6">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Slow · Live · Yunnan</span>
        </div>
        <h2 class="font-serif text-4xl md:text-5xl text-ivory-100 leading-[1.1] mb-5">
          A Moment,<br/> <em class="not-italic text-gold">Always</em> in the Garden.
        </h2>
        <p class="text-sand max-w-xl leading-relaxed">
          Three slow cameras. Four seasons. One cloud mountain in Yunnan.
          Watch the tea grow, 24 hours a day — quietly, without commentary.
        </p>
      </div>

      <!-- 三张卡：off-black 内联元素，零 emoji，零粗圆角 -->
      <div class="grid md:grid-cols-3 gap-px bg-gold/15">
        <RouterLink v-for="p in (presets.length ? presets : [
          { name: '云雾茶区 · 临沧', location: 'Yunnan · Lincang', status: 'live' },
          { name: '布朗山茶区 · 西双版纳', location: 'Yunnan · Xishuangbanna', status: 'live' },
          { name: '澜沧茶区 · 普洱', location: 'Yunnan · Pu\'er', status: 'live' },
        ])" :key="p.name"
            to="/live-room"
            class="group relative bg-ink-900 hover:bg-ink-800 transition-duration-lux overflow-hidden block focus:outline-none">

          <!-- 舞台视频帧：off-black + 1px 金边线框 placeholder -->
          <div class="aspect-[16/10] relative">
            <div class="absolute inset-0 bg-gradient-to-br from-ink-800 to-ink-900"></div>
            <!-- 极淡的金线框装饰 — 像老照片取景器 -->
            <div class="absolute inset-4 border border-gold/20 pointer-events-none"></div>

            <!-- LIVE 指示：一个静态香槟金小点 + SERIF CAPS 极小字 — 绝不用红色脉冲 -->
            <div class="absolute top-5 left-5 flex items-center gap-2">
              <span class="w-1 h-1 bg-gold rounded-full"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold/90 font-sans">Live</span>
            </div>

            <!-- 悬停：金线从 20% → 60%，600ms slow transition -->
            <div class="absolute inset-4 border border-gold/0 group-hover:border-gold/60 transition-duration-lux pointer-events-none"></div>
          </div>

          <!-- 象牙白丝绒感底部条 -->
          <div class="px-6 py-7 bg-ivory-100">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h3 class="font-serif text-ink-900 text-lg md:text-xl leading-tight">{{ p.name }}</h3>
                <p class="text-[11px] uppercase tracking-lux text-sand mt-2 font-sans">{{ p.location }}</p>
              </div>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans shrink-0 mt-2">View →</span>
            </div>
          </div>
        </RouterLink>
      </div>

      <!-- 底部一行极细的 serif 小字 — 像画廊导览的 closing line -->
      <p class="text-center text-[11px] uppercase tracking-lux text-sand mt-14 font-sans">
        Cameras operate on local time · The advisory reserves the right to adjust schedules
      </p>
    </div>
  </section>

  <!-- COLLECTION -->
  <section id="collection" class="py-24 bg-ivory-100">
    <div class="max-w-7xl mx-auto px-6">
      <!-- 区块头：serif 大标题 + 香槟金 hairline -->
      <div class="flex items-end justify-between mb-14 flex-wrap gap-4">
        <div>
          <div class="flex items-center gap-3 mb-4">
            <span class="h-px w-10 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Signature · From The House</span>
          </div>
          <h2 class="font-serif text-4xl text-ink-900">Our Bespoke Collection</h2>
        </div>
        <RouterLink to="/bespoke" class="text-[11px] uppercase tracking-lux text-sand hover:text-ink-900 transition font-sans">View Collection →</RouterLink>
      </div>

      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-0 bg-gold/10">
        <RouterLink v-for="cp in (featured.length ? featured : [
          { id: 1, product_token: 'demo-1', title: '邦东古树饼', tea_garden_location: '云南省 · 临沧市', master_name: '李师傅', unit_price: 68.5, tea_type: 'raw_puer' },
          { id: 2, product_token: 'demo-2', title: '班章熟砖', tea_garden_location: '云南省 · 西双版纳州', master_name: '张师傅', unit_price: 120, tea_type: 'ripe_puer' },
          { id: 3, product_token: 'demo-3', title: '景迈金瓜', tea_garden_location: '云南省 · 普洱市', master_name: '王师傅', unit_price: 88, tea_type: 'raw_puer' },
        ])" :key="cp.id" :to="`/bespoke/${cp.product_token || 'demo-' + cp.id}`"
          class="group bg-ivory-100 hover:bg-ivory-50 transition-duration-lux border border-transparent hover:border-gold/30 block">
          <!-- 卡头：off-black placeholder 帧 + 金线框装饰 → 绝不用 emoji -->
          <div class="aspect-[4/3] relative bg-ink-900 overflow-hidden">
            <div class="absolute inset-5 border border-gold/15 group-hover:border-gold/50 transition-duration-lux pointer-events-none"></div>
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="text-center">
                <div class="font-serif text-gold/50 text-2xl tracking-wider">N° {{ cp.id || '—' }}</div>
                <div class="text-[10px] uppercase tracking-lux text-sand mt-2 font-sans">{{ cp.tea_type === 'raw_puer' ? 'Raw Pu\'er' : 'Ripe Pu\'er' }}</div>
              </div>
            </div>
          </div>
          <div class="p-7">
            <h3 class="font-serif text-xl text-ink-900 mb-2">{{ cp.title }}</h3>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-6 font-sans">
              {{ cp.tea_garden_location }} · Master {{ cp.master_name }}
            </p>
            <div class="flex items-center justify-between pt-5 border-t border-gold/15">
              <span class="font-serif text-ink-900 text-xl">£{{ cp.unit_price }}</span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">View →</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </section>

  <!-- YOUR OWN BATCH — Bespoke 独立叙述，不推销，陈述式 -->
  <section class="py-24 md:py-28 bg-ivory-50 border-t border-gold/15">
    <div class="max-w-4xl mx-auto px-6 text-center">
      <!-- 一条极细的金线做呼吸感分隔 -->
      <div class="flex items-center justify-center gap-4 mb-10">
        <span class="h-px w-16 bg-gold/40"></span>
        <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Your Own Batch</span>
        <span class="h-px w-16 bg-gold/40"></span>
      </div>

      <h2 class="font-serif text-4xl md:text-5xl text-ink-900 leading-[1.15] mb-6">
        A tea made,<br/> <em class="not-italic">for</em> you.
      </h2>

      <p class="text-ink-700 leading-relaxed max-w-xl mx-auto mb-4">
        Choose the mountain village where your leaves are picked. The roast.
        The paper it's wrapped in — cotton from Jiangsu, or bamboo from the same forest.
      </p>
      <p class="text-ink-700 leading-relaxed max-w-xl mx-auto mb-10">
        One master in Yunnan will oversee your batch from harvest to box.
        You can watch the tea grow, live, before it's even picked.
      </p>

      <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
        <RouterLink to="/bespoke"
          class="px-9 py-4 bg-ink-900 text-ivory-100 rounded-full font-medium hover:bg-ink-800 transition-all duration-300 hover:shadow-xl hover:-translate-y-0.5 text-center tracking-wide">
          Design Your Blend →
        </RouterLink>
      </div>

      <p class="text-[10px] uppercase tracking-lux text-sand mt-12 font-sans">
        Bespoke batches take approximately 45 days · Minimum order £88
      </p>
    </div>
  </section>
</template>
