<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { ref, computed, onMounted } from 'vue'
import { api } from '@/api/client'

// ===== Data =====
const gardens = ref<any[]>([])
const loading = ref(true)

// ===== Videos =====
const videos = ref<any[]>([])
const categories = ref<any[]>([])
const activeCategory = ref<number | null>(null) // null = all
const activeVideo = ref<any>(null) // selected video for player
const videoLoading = ref(false)

// ===== Filtered videos =====
const filteredVideos = computed(() => {
  if (activeCategory.value === null) return videos.value
  return videos.value.filter(v => v.category_id === activeCategory.value)
})

async function load() {
  loading.value = true
  // Gardens
  try {
    const d: any = await api.get('/public/slow-presets')
    const got = d?.items || d || []
    if (got.length) { gardens.value = got }
    else { gardens.value = fallbackGardens }
  } catch (e) {
    gardens.value = fallbackGardens
  } finally { loading.value = false }

  // Videos
  videoLoading.value = true
  try {
    const vresp: any = await api.get('/public/videos')
    videos.value = vresp?.items || []
  } catch { videos.value = [] }
  try {
    categories.value = await api.get('/public/video-categories') || []
  } catch { categories.value = [] }
  videoLoading.value = false

  // Auto-play first video
  if (videos.value.length && !activeVideo.value) {
    activeVideo.value = videos.value[0]
  }
}

function selectVideo(v: any) {
  activeVideo.value = v
  // Scroll to player on mobile
  if (window.innerWidth < 768) {
    const el = document.getElementById('video-player')
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function resolveUrl(url: string) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '') + url
}

function fmtDuration(sec: number) {
  if (!sec) return ''
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

const fallbackGardens = [
  { name: '云南省 · 临沧市 · 云雾茶区', location: '临沧 · 云南', description: '古树普洱核心产区，高海拔终年云雾缭绕' },
  { name: '云南省 · 西双版纳州 · 布朗山茶区', location: '西双版纳 · 云南', description: '王者之地，乔木古树，浓烈霸道' },
  { name: '云南省 · 普洱市 · 澜沧茶区', location: '普洱 · 云南', description: '千年万亩古茶园，布朗族与傣族世代守护' },
]

onMounted(load)
</script>

<template>
  <div class="pt-20">
    <!-- ===== HERO — 与 Live.vue 一致的 off-black 主舞台 ===== -->
    <section class="bg-ink-900 text-ivory-100 py-24 md:py-32 relative overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.05]"
           style="background-image: radial-gradient(circle at 20% 50%, #C5A572 0%, transparent 50%), radial-gradient(circle at 80% 30%, #C5A572 0%, transparent 40%);"></div>
      <div class="max-w-5xl mx-auto px-6 relative">
        <div class="flex items-center gap-3 mb-8">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Single Origin · Yunnan</span>
        </div>
        <h1 class="font-serif text-5xl md:text-7xl leading-[1.05] mb-8">
          Our Tea <em class="not-italic text-gold">Gardens</em>
        </h1>
        <p class="text-sand max-w-xl text-lg leading-relaxed">
          Direct from village-level tea gardens across Yunnan. Each garden has a 24/7 live camera
          so you can watch where your tea grew — and now, film stories from the mountain.
        </p>
      </div>
    </section>

    <!-- ===== VIDEO SECTION — 与 Live 页面设计一致 ===== -->
    <section class="py-20 md:py-24 bg-ink-900">
      <div class="max-w-7xl mx-auto px-4 md:px-6">
        <!-- Section header -->
        <div class="flex items-end justify-between mb-12 flex-wrap gap-4">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <span class="h-px w-10 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">From The Mountain</span>
            </div>
            <h2 class="font-serif text-3xl md:text-4xl text-ivory-100">Garden Films</h2>
          </div>
          <!-- Category filter -->
          <div v-if="categories.length" class="flex flex-wrap gap-2">
            <button
              @click="activeCategory = null"
              :class="[
                'px-4 py-2 text-[11px] uppercase tracking-lux font-sans rounded-full border transition-all duration-300',
                activeCategory === null
                  ? 'bg-gold/20 border-gold text-gold'
                  : 'border-ivory-100/20 text-sand hover:border-gold/50 hover:text-gold'
              ]"
            >All</button>
            <button
              v-for="cat in categories" :key="cat.id"
              @click="activeCategory = cat.id"
              :class="[
                'px-4 py-2 text-[11px] uppercase tracking-lux font-sans rounded-full border transition-all duration-300',
                activeCategory === cat.id
                  ? 'bg-gold/20 border-gold text-gold'
                  : 'border-ivory-100/20 text-sand hover:border-gold/50 hover:text-gold'
              ]"
            >{{ cat.name }}</button>
          </div>
        </div>

        <!-- No videos -->
        <div v-if="!videoLoading && !filteredVideos.length" class="text-center py-16">
          <p class="text-sand text-sm font-sans tracking-wide">Films are being curated. Please return soon.</p>
        </div>

        <!-- Player + List layout: PC side-by-side, mobile stacked -->
        <div v-if="filteredVideos.length" class="grid lg:grid-cols-[1fr_360px] gap-6">
          <!-- ===== Player ===== -->
          <div id="video-player" class="bg-ink-900">
            <!-- Video frame — 16:9, off-black velvet stage like LiveRoom -->
            <div class="relative bg-black border border-gold/15 overflow-hidden" style="aspect-ratio: 16/9;">
              <video
                v-if="activeVideo"
                :key="activeVideo.id"
                :src="resolveUrl(activeVideo.video_url)"
                :poster="activeVideo.cover_url ? resolveUrl(activeVideo.cover_url) : ''"
                controls
                playsinline
                class="w-full h-full object-contain"
              ></video>
              <div v-else class="flex items-center justify-center h-full">
                <span class="text-gold/40 text-sm font-sans tracking-lux uppercase">No film selected</span>
              </div>
            </div>
            <!-- Caption — like LiveRoom stage caption -->
            <div class="mt-5">
              <div v-if="activeVideo" class="flex items-center gap-3 mb-3">
                <span class="h-px w-8 bg-gold/40"></span>
                <span v-if="activeVideo.category" class="text-[10px] uppercase tracking-lux text-gold font-sans">{{ activeVideo.category.name }}</span>
                <span v-if="activeVideo.duration_sec" class="text-[10px] uppercase tracking-lux text-sand font-sans">{{ fmtDuration(activeVideo.duration_sec) }}</span>
              </div>
              <h3 class="font-serif text-2xl md:text-3xl text-ivory-100 leading-tight mb-2">
                {{ activeVideo?.title || '—' }}
              </h3>
              <p v-if="activeVideo?.description" class="text-sand text-sm leading-relaxed max-w-2xl">
                {{ activeVideo.description }}
              </p>
            </div>
          </div>

          <!-- ===== Video list — sidebar on PC, horizontal scroll on mobile ===== -->
          <aside class="lg:border-l lg:border-ivory-100/10 lg:pl-6">
            <h4 class="text-[10px] uppercase tracking-lux text-gold font-sans mb-4 pb-3 border-b border-ivory-100/10">
              All Films · {{ filteredVideos.length }}
            </h4>
            <!-- Mobile: horizontal scroll; Desktop: vertical list -->
            <div class="flex lg:flex-col gap-3 overflow-x-auto lg:overflow-x-visible pb-2 lg:pb-0">
              <button
                v-for="v in filteredVideos" :key="v.id"
                @click="selectVideo(v)"
                :class="[
                  'shrink-0 lg:shrink flex gap-3 p-3 rounded-sm border transition-all duration-300 text-left',
                  activeVideo?.id === v.id
                    ? 'border-gold/50 bg-gold/5'
                    : 'border-ivory-100/10 hover:border-gold/30'
                ]"
              >
                <!-- Thumbnail -->
                <div class="relative w-24 lg:w-20 shrink-0" style="aspect-ratio: 16/10;">
                  <img
                    v-if="v.cover_url"
                    :src="resolveUrl(v.cover_url)"
                    class="w-full h-full object-cover rounded-sm"
                  />
                  <div v-else class="w-full h-full bg-ink-800 rounded-sm flex items-center justify-center">
                    <span class="text-gold/30 text-xs">▶</span>
                  </div>
                  <span v-if="v.duration_sec" class="absolute bottom-1 right-1 px-1.5 py-0.5 bg-black/70 text-ivory-100 text-[9px] font-sans rounded-sm">
                    {{ fmtDuration(v.duration_sec) }}
                  </span>
                </div>
                <!-- Text -->
                <div class="flex-1 min-w-0">
                  <p class="font-serif text-sm text-ivory-100 leading-snug truncate">{{ v.title }}</p>
                  <p v-if="v.category" class="text-[10px] uppercase tracking-lux text-sand mt-1 font-sans">{{ v.category.name }}</p>
                </div>
              </button>
            </div>
          </aside>
        </div>
      </div>
    </section>

    <!-- ===== TEA GARDENS GRID — 保持原有的茶园卡片 ===== -->
    <section class="py-20 md:py-24 bg-ivory-100">
      <div class="max-w-6xl mx-auto px-6">
        <div class="flex items-end justify-between mb-14 flex-wrap gap-4">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <span class="h-px w-10 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Cloud · Tea · Mountains</span>
            </div>
            <h2 class="font-serif text-4xl text-ink-900">The Gardens</h2>
          </div>
          <RouterLink to="/live" class="text-[11px] uppercase tracking-lux text-sand hover:text-ink-900 transition font-sans">
            Watch Live Cameras →
          </RouterLink>
        </div>

        <div v-if="loading" class="text-center py-20 text-tea-500 font-sans text-sm tracking-wide">Loading tea gardens…</div>

        <div v-else class="grid md:grid-cols-2 gap-px bg-gold/15">
          <div v-for="(g, i) in gardens" :key="g.id || g.name"
               :class="[
                 'group bg-white hover:bg-ivory-50 transition-all duration-500',
                 'p-8 flex flex-col'
               ]">
            <!-- Visual -->
            <div class="aspect-[16/9] mb-6 relative overflow-hidden"
                 :class="[
                   'flex items-center justify-center',
                   i % 3 === 0 ? 'bg-gradient-to-br from-ink-900 to-ink-700' :
                   i % 3 === 1 ? 'bg-gradient-to-br from-ink-800 to-ink-600' :
                                 'bg-gradient-to-br from-ink-700 to-ink-500'
                 ]">
              <span class="text-gold/30 font-serif text-4xl">⛰</span>
              <div v-if="g.camera_rtmp_url"
                   class="absolute top-3 left-3 flex items-center gap-2 px-2.5 py-1 bg-ink-900/80 text-gold text-[10px] uppercase tracking-lux font-sans rounded-full">
                <span class="w-1 h-1 rounded-full bg-gold"></span> Live
              </div>
            </div>
            <!-- Content -->
            <h3 class="font-serif text-xl text-ink-900 mb-2 leading-snug">{{ g.name || 'Unnamed Garden' }}</h3>
            <p class="text-[11px] uppercase tracking-lux text-sand mb-3 font-sans">{{ g.location || g.region || '—' }}</p>
            <p class="text-sm text-ink-700 leading-relaxed mb-6 flex-1">{{ g.description || 'Visit our live streams to see this tea garden in action.' }}</p>
            <RouterLink to="/live" class="text-[10px] uppercase tracking-lux text-gold hover:text-ink-900 transition font-sans">
              Watch Live →
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== Footer note ===== -->
    <section class="bg-ivory-100 py-14 border-t border-gold/10">
      <div class="max-w-4xl mx-auto px-6 text-center">
        <p class="text-[11px] uppercase tracking-lux text-sand leading-loose font-sans">
          Films may be subject to seasonal availability. The advisory reserves the right to adjust schedules.
        </p>
      </div>
    </section>
  </div>
</template>
