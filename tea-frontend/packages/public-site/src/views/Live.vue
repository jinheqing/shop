<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getSlowPresets, listLiveRooms, type LiveRoom } from '@/api/live'

const presets = ref<any[]>([])
const rooms = ref<LiveRoom[]>([])
onMounted(async () => {
  try { presets.value = (await getSlowPresets()) || [] } catch {}
  try { rooms.value = ((await listLiveRooms()) as any)?.items || [] } catch {}
})
</script>

<template>
  <div class="pt-20 min-h-screen">
    <section class="bg-tea-900 text-tea-50 py-20">
      <div class="max-w-6xl mx-auto px-6">
        <h1 class="font-serif text-5xl md:text-6xl mb-4">Live From the Tea Gardens.</h1>
        <p class="text-lg text-tea-300 max-w-2xl">24/7 slow-live cameras show the tea growing. Taste sessions streamed by appointment. Watch your own bespoke tea as it's being roasted.</p>
      </div>
    </section>

    <!-- Live disclaimer -->
    <div class="bg-amber-50 border-y border-amber-200 text-sm text-amber-900">
      <div class="max-w-6xl mx-auto px-6 py-3 text-center">
        <span class="font-medium">⚠️</span> Live feeds may be delayed by up to 30 seconds. Camera availability depends on local connectivity. Tea preparation times shown are estimates — actual production schedules vary by garden and harvest season.
      </div>
    </div>

    <section class="py-16 bg-tea-50">
      <div class="max-w-7xl mx-auto px-6">
        <h2 class="font-serif text-3xl text-tea-900 mb-8">🌱 24/7 Garden Cameras</h2>
        <div class="grid md:grid-cols-3 gap-6">
          <div v-for="p in (presets.length ? presets : [
            { name: '云南省 · 临沧市 · 云雾茶区', location: '临沧', camera_rtmp_url: 'rtmp://localhost:1935/slow/manzhang', status: 'live' },
            { name: '云南省 · 西双版纳州 · 布朗山茶区', location: '西双版纳', camera_rtmp_url: 'rtmp://localhost:1935/slow/banzhang', status: 'live' },
            { name: '云南省 · 普洱市 · 澜沧茶区', location: '普洱', camera_rtmp_url: 'rtmp://localhost:1935/slow/jingmai', status: 'live' },
          ])" :key="p.name" class="rounded-2xl overflow-hidden card-hover bg-white border border-tea-100">
            <div class="relative aspect-video bg-tea-900 flex items-center justify-center text-tea-200">
              <div class="text-center">
                <div class="text-6xl mb-2">🏞️</div>
                <div class="text-xs font-mono text-tea-400">{{ p.camera_rtmp_url }}</div>
              </div>
              <div class="absolute top-3 left-3 flex items-center gap-2 px-3 py-1 bg-red-600 text-white text-xs font-medium rounded-full">
                <span class="w-1.5 h-1.5 rounded-full bg-white live-dot"></span> LIVE · 24/7
              </div>
            </div>
            <div class="p-5">
              <h3 class="font-serif text-lg text-tea-900">{{ p.name }}</h3>
              <p class="text-sm text-tea-500">{{ p.location }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-16 bg-white">
      <div class="max-w-7xl mx-auto px-6">
        <h2 class="font-serif text-3xl text-tea-900 mb-4">🎥 Upcoming Tasting Sessions</h2>
        <p class="text-tea-600 mb-8">Book a 1:1 tasting with our tea advisors. They'll brew Pu'er from three gardens and walk you through the tasting notes.</p>
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="p-6 rounded-xl bg-tea-50 border border-tea-100">
            <div class="text-xs text-tea-500 mb-2">SAT · Sep 20</div>
            <h4 class="font-serif text-lg mb-2">Morning Gongfu</h4>
            <p class="text-xs text-tea-600 mb-4">10:00 GMT · 45 min</p>
            <button class="w-full py-2 text-sm bg-tea-700 text-white rounded-lg hover:bg-tea-800">Book Now</button>
          </div>
          <div class="p-6 rounded-xl bg-tea-50 border border-tea-100">
            <div class="text-xs text-tea-500 mb-2">SUN · Sep 21</div>
            <h4 class="font-serif text-lg mb-2">Rare Tasting</h4>
            <p class="text-xs text-tea-600 mb-4">15:00 GMT · 60 min</p>
            <button class="w-full py-2 text-sm bg-tea-700 text-white rounded-lg hover:bg-tea-800">Book Now</button>
          </div>
          <div class="p-6 rounded-xl bg-tea-50 border border-tea-100">
            <div class="text-xs text-tea-500 mb-2">SAT · Sep 27</div>
            <h4 class="font-serif text-lg mb-2">Masterclass</h4>
            <p class="text-xs text-tea-600 mb-4">14:00 GMT · 90 min</p>
            <button class="w-full py-2 text-sm bg-tea-700 text-white rounded-lg hover:bg-tea-800">Book Now</button>
          </div>
          <div class="p-6 rounded-xl bg-tea-50 border border-tea-100">
            <div class="text-xs text-tea-500 mb-2">PRIVATE</div>
            <h4 class="font-serif text-lg mb-2">Your Bespoke Stream</h4>
            <p class="text-xs text-tea-600 mb-4">From £50 · 30 min</p>
            <button class="w-full py-2 text-sm bg-tea-700 text-white rounded-lg hover:bg-tea-800">Request</button>
          </div>
        </div>

        <h3 class="font-serif text-xl text-tea-900 mt-16 mb-4">Your Active Rooms</h3>
        <div v-if="rooms.length" class="grid md:grid-cols-3 gap-4">
          <div v-for="r in rooms" :key="r.id" class="p-4 bg-white rounded-xl border border-tea-100 flex items-center justify-between">
            <div>
              <div class="font-medium text-tea-900">{{ r.room_name }}</div>
              <div class="text-xs text-tea-500">{{ r.room_type }}</div>
            </div>
            <span :class="['text-xs px-2 py-1 rounded-full', r.status==='live'?'bg-green-100 text-green-700':'bg-tea-100 text-tea-500']">{{ r.status }}</span>
          </div>
        </div>
        <div v-else class="text-sm text-tea-500">No active rooms. Book a session above.</div>
      </div>
    </section>
  </div>
</template>
