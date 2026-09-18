<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getSlowPresets, listLiveRooms, type LiveRoom } from '@/api/live'

const presets = ref<any[]>([])
const rooms = ref<LiveRoom[]>([])
onMounted(async () => {
  try { presets.value = (await getSlowPresets()) || [] } catch {}
  try { rooms.value = ((await listLiveRooms()) as any)?.items || [] } catch {}
})
</script>

<!--
  直播列表页 — Luxury Edition
  设计参考：Château Lafite Rothschild 酒庄导览册 / Dior 后台 / Rolls-Royce Private Office
  原则：
    - off-black + ivory + champagne 三色
    - 零 emoji，零粗圆角，零红色脉冲，零暴露内部 RTMP URL
    - serif 大标题 + 10px SERIF CAPS 小标签 + 香槟金 hairline
    - 慢过渡（600ms），悬停金线从 20% → 60%
-->
<template>
  <div class="pt-20 min-h-screen bg-ivory-100">

    <!-- ===== HERO — off-black 主舞台 ===== -->
    <section class="bg-ink-900 text-ivory-100 py-24 md:py-32 relative overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.05]"
           style="background-image: radial-gradient(circle at 20% 50%, #C5A572 0%, transparent 50%), radial-gradient(circle at 80% 30%, #C5A572 0%, transparent 40%);"></div>
      <div class="max-w-5xl mx-auto px-6 relative">
        <div class="flex items-center gap-3 mb-8">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Slow · Live · 24 Hours</span>
        </div>
        <h1 class="font-serif text-5xl md:text-7xl leading-[1.05] mb-8">
          The Garden,<br/> <em class="not-italic text-gold">Uninterrupted.</em>
        </h1>
        <p class="text-sand max-w-xl text-lg leading-relaxed">
          Three slow cameras positioned in the cloud mountains of Yunnan.
          They show the tea growing — nothing more, nothing less.
          No commentary. No editing. Just the mountain, in real time.
        </p>
      </div>
    </section>

    <!-- ===== 24/7 SLOW CAMERAS ===== -->
    <section class="py-24 bg-ivory-100">
      <div class="max-w-7xl mx-auto px-6">
        <div class="flex items-end justify-between mb-14 flex-wrap gap-4">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <span class="h-px w-10 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Three Cameras · Always On</span>
            </div>
            <h2 class="font-serif text-4xl text-ink-900">Garden Cameras</h2>
          </div>
          <p class="text-[11px] uppercase tracking-lux text-sand font-sans">
            {{ presets.length || 3 }} stations · 24/7
          </p>
        </div>

        <!-- 三联卡：gap-px 香槟金线分隔 → 像美术馆三幅挂画 -->
        <div class="grid md:grid-cols-3 gap-px bg-gold/15">
          <RouterLink v-for="p in (presets.length ? presets : [
            { name: '云雾茶区 · 临沧', location: 'Yunnan · Lincang', status: 'live' },
            { name: '布朗山茶区 · 西双版纳', location: 'Yunnan · Xishuangbanna', status: 'live' },
            { name: '澜沧茶区 · 普洱', location: 'Yunnan · Pu\'er', status: 'live' },
          ])" :key="p.name"
              to="/live-room"
              class="group bg-ink-900 hover:bg-ink-800 transition-duration-lux block">
            <!-- 视频帧 -->
            <div class="aspect-[16/10] relative">
              <div class="absolute inset-5 border border-gold/15 group-hover:border-gold/50 transition-duration-lux pointer-events-none"></div>
              <!-- LIVE 指示：香槟金小点 + SERIF CAPS 字 -->
              <div class="absolute top-6 left-6 flex items-center gap-2">
                <span class="w-1 h-1 bg-gold rounded-full"></span>
                <span class="text-[10px] uppercase tracking-lux text-gold/90 font-sans">Live · 24/7</span>
              </div>
            </div>
            <!-- 象牙白标签 -->
            <div class="px-6 py-7 bg-ivory-100">
              <h3 class="font-serif text-ink-900 text-lg">{{ p.name }}</h3>
              <p class="text-[11px] uppercase tracking-lux text-sand mt-2 font-sans">{{ p.location }}</p>
            </div>
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ===== PRIVATE TASTING — 预约品鉴 ===== -->
    <!-- 不是 "Upcoming Tasting Sessions · Book Now" 紧迫感 → 是 "By Appointment · Request" 的克制 -->
    <section class="py-24 bg-white border-y border-gold/10">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center mb-16">
          <div class="flex items-center gap-3 justify-center mb-6">
            <span class="h-px w-10 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Private · Confidential</span>
            <span class="h-px w-10 bg-gold/50"></span>
          </div>
          <h2 class="font-serif text-4xl md:text-5xl text-ink-900 mb-5">A Private Tasting</h2>
          <p class="text-sand max-w-2xl mx-auto leading-relaxed">
            Request a one-to-one session with our tea advisor.
            Three gardens, three brews, three hours — on your schedule.
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-0 bg-gold/10">
          <div class="bg-white p-8 transition-duration-lux hover:bg-ivory-100 group">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-4">Saturday · Sep 20</div>
            <h4 class="font-serif text-xl text-ink-900 mb-2">Morning Gongfu</h4>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-8 font-sans">10:00 BST · 45 min</p>
            <RouterLink to="/chat" class="text-[10px] uppercase tracking-lux text-gold font-sans group-hover:text-ink-900 transition">Request →</RouterLink>
          </div>
          <div class="bg-white p-8 transition-duration-lux hover:bg-ivory-100 group">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-4">Sunday · Sep 21</div>
            <h4 class="font-serif text-xl text-ink-900 mb-2">Rare Tasting</h4>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-8 font-sans">15:00 BST · 60 min</p>
            <RouterLink to="/chat" class="text-[10px] uppercase tracking-lux text-gold font-sans group-hover:text-ink-900 transition">Request →</RouterLink>
          </div>
          <div class="bg-white p-8 transition-duration-lux hover:bg-ivory-100 group">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-4">Saturday · Sep 27</div>
            <h4 class="font-serif text-xl text-ink-900 mb-2">Masterclass</h4>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-8 font-sans">14:00 BST · 90 min</p>
            <RouterLink to="/chat" class="text-[10px] uppercase tracking-lux text-gold font-sans group-hover:text-ink-900 transition">Request →</RouterLink>
          </div>
          <div class="bg-ink-900 p-8 group">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-4">By Request</div>
            <h4 class="font-serif text-xl text-ivory-100 mb-2">Your Bespoke Stream</h4>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-8 font-sans">From £50 · 30 min</p>
            <RouterLink to="/chat" class="text-[10px] uppercase tracking-lux text-gold font-sans group-hover:text-ivory-100 transition">Request →</RouterLink>
          </div>
        </div>

        <!-- 已激活的房间（如果有） -->
        <div v-if="rooms.length" class="mt-24">
          <h3 class="font-serif text-2xl text-ink-900 mb-8 text-center">Your Active Rooms</h3>
          <div class="grid md:grid-cols-3 gap-0 bg-gold/10 max-w-3xl mx-auto">
            <div v-for="r in rooms" :key="r.id"
                 class="bg-white p-6 flex items-center justify-between group">
              <div>
                <div class="font-serif text-ink-900">{{ r.room_name }}</div>
                <div class="text-[10px] uppercase tracking-lux text-sand mt-1 font-sans">{{ r.room_type }}</div>
              </div>
              <span :class="[
                'text-[10px] uppercase tracking-lux font-sans flex items-center gap-2',
                r.status === 'live' ? 'text-gold' : 'text-sand'
              ]">
                <span :class="['w-1 h-1 rounded-full', r.status === 'live' ? 'bg-gold' : 'bg-sand']"></span>
                {{ r.status }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 底部极细的注释 — 取代原来的 amber 警告条 ===== -->
    <section class="bg-ivory-100 py-14">
      <div class="max-w-4xl mx-auto px-6 text-center">
        <p class="text-[11px] uppercase tracking-lux text-sand leading-loose font-sans">
          Feeds may be delayed by up to 30 seconds. Camera availability depends on local connectivity.
          The advisory reserves the right to reschedule sessions. Private tastings are confirmed within 24 hours.
        </p>
      </div>
    </section>
  </div>
</template>
