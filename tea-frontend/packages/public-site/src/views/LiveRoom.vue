<script setup lang="ts">
import { onMounted, onUnmounted, ref, nextTick, computed } from 'vue'
import { api } from '@/api/client'

// ============ Room State ============
const rooms = ref<any[]>([])
const activeRoom = ref<any>(null)
const presets = ref<any[]>([])
const roomID = computed(() => activeRoom.value?.room_id || `demo-${activeRoom.value?.id || 0}`)

// ============ Barrage State ============
const barrages = ref<{id: number; text: string; nickname: string; ts: number; color: string; top: number}[]>([])
const barrageInput = ref('')
const barrageWS = ref<WebSocket | null>(null)
const barrageConnected = ref(false)
const barrageSubscribed = ref(false)

const EMOJI_SET = ['😀','😂','🥰','😍','👏','🔥','💪','❤️','🎉','🍵','🌿','✨','🫖','🥳']

// ============ Chat Panel State ============
const chatOpen = ref(false)
const chatConvs = ref<any[]>([])
const chatMessages = ref<any[]>([])
const chatInput = ref('')
const chatWS = ref<WebSocket | null>(null)
const wsConnected = ref(false)
const emojiPickerOpen = ref(false)

// ============ LiveKit State ============
const lkClient = ref<any>(null)
const lkConnected = ref(false)
const lkLocalVideo = ref<HTMLVideoElement | null>(null)
const lkLocalAudio = ref(false)
const lkRemoteParticipants = ref<any[]>([])
const lkJoinLoading = ref(false)
const hasLiveKitSDK = ref(false)

// ============ Demo Room ============
const DEMO_ROOM = {
  id: 999, room_id: 'demo-gongfu', room_name: 'Gongfu Tea Demo', room_type: 'live',
  status: 'live', location: 'Yunnan · Lincang', camera_rtmp_url: 'rtmp://localhost:1935/live/demo'
}

// ============ Init ============
onMounted(async () => {
  try {
    rooms.value = ((await api.get('/live-rooms') as any)?.items || [])
  } catch { rooms.value = [] }
  try {
    presets.value = ((await api.get('/slow-presets') as any)?.items || [])
  } catch { presets.value = [] }
  if (rooms.value.length > 0) {
    await selectRoom(rooms.value[0])
  } else {
    await selectRoom(DEMO_ROOM)
  }

  // Check livekit SDK availability
  try {
    const mod = await import('livekit-client').catch(() => null)
    if (mod?.Room) { hasLiveKitSDK.value = true }
  } catch { hasLiveKitSDK.value = false }
})

onUnmounted(() => {
  barrageWS.value?.close()
  chatWS.value?.close()
  leaveLiveKit()
})

// ============ Room Switch ============
async function selectRoom(room: any) {
  activeRoom.value = room
  // Auto-crash a system conversation for this room's chat
  if (chatOpen.value) await ensureRoomChat()
  // Connect barrage ws
  connectBarrageWS()
}

// ============ Barrage WebSocket ============
function connectBarrageWS() {
  const token = localStorage.getItem('user_token') || localStorage.getItem('staff_token') || 'guest'
  const host = window.location.hostname
  barrageWS.value?.close()
  barrageWS.value = new WebSocket(`ws://${host}:8080/api/v1/ws/im?token=${token}`)
  barrageWS.value.onopen = () => {
    barrageConnected.value = true
    if (roomID.value) subscribeRoomBarrage()
  }
  barrageWS.value.onmessage = (ev) => {
    try {
      const env = JSON.parse(ev.data)
      if (env.type === 'barrage') {
        addBarrage(env.payload.content, env.payload.nickname || 'Viewer')
      }
    } catch {}
  }
  barrageWS.value.onclose = () => { barrageConnected.value = false; barrageSubscribed.value = false }
}

function subscribeRoomBarrage() {
  if (!barrageWS.value || !roomID.value || barrageSubscribed.value) return
  barrageWS.value.send(JSON.stringify({ type: 'join_room', payload: { room_id: roomID.value } }))
  barrageSubscribed.value = true
}

let barrageID = 0
function addBarrage(text: string, nickname: string) {
  const colors = ['#fff', '#ffd700', '#ff6b6b', '#4ecdc4', '#a8e6cf', '#ff8a5c']
  barrages.value.push({
    id: ++barrageID,
    text, nickname,
    ts: Date.now(),
    color: colors[Math.floor(Math.random() * colors.length)],
    top: Math.floor(Math.random() * 70) + 5 // 5-75% height
  })
  // auto-clean after 12s
  setTimeout(() => {
    barrages.value = barrages.value.filter(b => b.id !== barrageID)
  }, 12000)
}

function sendBarrage() {
  if (!barrageInput.value.trim() || !barrageWS.value) return
  if (!barrageConnected.value) connectBarrageWS()
  const nickname = localStorage.getItem('user_name') || 'Tea Lover'
  barrageWS.value.send(JSON.stringify({
    type: 'barrage',
    payload: { room_id: roomID.value, content: barrageInput.value, nickname }
  }))
  addBarrage(barrageInput.value, 'You')
  barrageInput.value = ''
}

// ============ Chat Panel (reuses IM WS) ============
async function ensureRoomChat() {
  // Find or create system conversation for this room
  let conv: any = null
  try {
    const list = (await api.get('/conversations') as any)?.items || []
    conv = list.find((c: any) => c.title === `room-${roomID.value}`)
  } catch {}
  if (!conv) {
    try {
      conv = await api.post('/conversations', { title: `room-${roomID.value}`, conversation_type: 'group' }) as any
    } catch {}
  }
  if (!conv) return

  chatOpen.value = true
  chatConvs.value = [conv]
  try {
    const resp = await api.get(`/conversations/${conv.id}/messages`) as any
    chatMessages.value = resp.items || resp || []
  } catch { chatMessages.value = [] }

  // Connect chat WS
  const token = localStorage.getItem('user_token') || localStorage.getItem('staff_token')
  if (!token) return
  const host = window.location.hostname
  chatWS.value?.close()
  chatWS.value = new WebSocket(`ws://${host}:8080/api/v1/ws/im?token=${token}`)
  chatWS.value.onopen = () => {
    wsConnected.value = true
    chatWS.value?.send(JSON.stringify({ type: 'join_conversation', payload: { conversation_id: conv.id } }))
  }
  chatWS.value.onclose = () => { wsConnected.value = false }
  chatWS.value.onmessage = (ev) => {
    try {
      const env = JSON.parse(ev.data)
      if (env.type === 'chat_message') {
        chatMessages.value.push(env.payload)
        nextTick(() => { const el = document.getElementById('chat-scroll'); if (el) el.scrollTop = el.scrollHeight })
      }
    } catch {}
  }
}

function closeChat() { chatOpen.value = false; chatWS.value?.close() }

function sendChat() {
  if (!chatInput.value.trim() || !chatConvs.value.length || !chatWS.value) return
  const payload = {
    conversation_id: chatConvs.value[0].id,
    content: chatInput.value.trim(),
    message_type: 'text'
  }
  chatWS.value.send(JSON.stringify({ type: 'send_message', payload }))
  chatMessages.value.push({
    content: chatInput.value.trim(), message_type: 'text', sender_type: 'me', created_at: new Date().toISOString()
  })
  chatInput.value = ''
}

// ============ LiveKit 1:1 Connect ============
async function joinLiveKit() {
  if (!hasLiveKitSDK.value) {
    alert('LiveKit SDK not loaded. Install livekit-client package and configure a LiveKit Cloud or self-hosted server.')
    return
  }
  lkJoinLoading.value = true
  try {
    // Get token from backend
    const resp: any = await api.get('/livekit/token')
    const token = resp?.token || resp?.access_token
    if (!token) throw new Error('no livekit token')

    const { Room, RoomEvent, ConnectOptions, VideoPresets } = await import('livekit-client')
    const room = new Room({ videoCaptureDefaults: { resolution: VideoPresets.h720 }, audioCaptureDefaults: {} })
    lkClient.value = room

    room.on(RoomEvent.TrackSubscribed, (track: any, publication: any, participant: any) => {
      console.log('Track subscribed from', participant.identity)
    })
    room.on(RoomEvent.ParticipantConnected, (p: any) => {
      lkRemoteParticipants.value.push(p)
    })
    room.on(RoomEvent.ParticipantDisconnected, (p: any) => {
      lkRemoteParticipants.value = lkRemoteParticipants.value.filter((x: any) => x.identity !== p.identity)
    })

    await room.connect(import.meta.env.VITE_LIVEKIT_URL || 'wss://tea.livekit.cloud', token, {
      autoSubscribe: true,
      dynacast: true
    })
    lkConnected.value = true

    // Publish local camera+mic
    await room.localParticipant.setCameraEnabled(true)
    await room.localParticipant.setMicrophoneEnabled(true)

    console.log('LiveKit connected:', room.name)
  } catch (e: any) {
    console.error('LiveKit join failed:', e)
    alert('LiveKit join failed: ' + (e?.message || String(e)))
  } finally {
    lkJoinLoading.value = false
  }
}

async function leaveLiveKit() {
  try { await lkClient.value?.disconnect() } catch {}
  lkConnected.value = false
  lkLocalVideo.value = null
  lkRemoteParticipants.value = []
}

async function toggleMic() {
  if (!lkClient.value) return
  lkLocalAudio.value = !lkLocalAudio.value
  await lkClient.value.localParticipant.setMicrophoneEnabled(!lkLocalAudio.value)
}

async function toggleCam() {
  if (!lkClient.value) return
  const enabled = !lkClient.value.localParticipant.isCameraEnabled
  await lkClient.value.localParticipant.setCameraEnabled(enabled)
}

// ============ Quick Room Select ============
function selectPresetAsRoom(p: any) {
  selectRoom({ id: p.id, room_id: `preset-${p.id}`, room_name: p.name, room_type: 'live', status: 'live', location: p.location, camera_rtmp_url: p.camera_rtmp_url })
}
</script>

<template>
  <div class="min-h-screen bg-tea-950 text-tea-50">
    {/* Hero */}
    <section class="bg-gradient-to-br from-tea-900 via-tea-800 to-slate-900 py-12">
      <div class="max-w-7xl mx-auto px-6">
        <h1 class="font-serif text-4xl md:text-5xl mb-2">🔴 Live From the Gardens</h1>
        <p class="text-tea-300 max-w-2xl">Slow-live cameras show you where your tea grows. Join a 1:1 tasting or chat with other tea lovers.</p>
      </div>
    </section>

    <!-- Main layout: video + sidebar -->
    <div class="max-w-7xl mx-auto px-4 md:px-6 py-6">
      <div class="flex flex-col lg:flex-row gap-6">
        <!-- Left: Video + barrage -->
        <div class="flex-1 min-w-0">
          <!-- Video stage -->
          <div class="relative bg-black rounded-2xl overflow-hidden aspect-video shadow-2xl border border-tea-700">
            <!-- Demo video area -->
            <div class="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-tea-900 via-emerald-950 to-slate-900">
              <div class="text-center">
                <div class="text-8xl mb-3">🍃</div>
                <div class="text-tea-400 font-mono text-sm mb-2">{{ activeRoom?.camera_rtmp_url || 'rtmp://localhost:1935/live/demo' }}</div>
                <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-red-600/80 text-sm font-medium">
                  <span class="w-2 h-2 rounded-full bg-white animate-pulse"></span>
                  LIVE · {{ activeRoom?.location || 'Yunnan, China' }}
                </div>
              </div>
            </div>

            <!-- LiveKit local video (small picture-in-picture) -->
            <div v-if="lkConnected" class="absolute bottom-4 right-4 w-32 h-24 md:w-40 md:h-28 rounded-lg overflow-hidden border-2 border-white/30 bg-black z-20">
              <video ref="lkLocalVideo" autoplay muted playsinline class="w-full h-full object-cover" />
            </div>

            <!-- Barrage overlay -->
            <div class="absolute inset-0 pointer-events-none overflow-hidden z-10">
              <div v-for="b in barrages" :key="b.id"
                   class="absolute whitespace-nowrap font-medium text-lg md:text-xl barrage-item drop-shadow-lg"
                   :style="{ top: b.top + '%', color: b.color }">
                {{ b.nickname }}: {{ b.text }}
              </div>
            </div>

            <!-- Top control bar -->
            <div class="absolute top-0 left-0 right-0 p-4 flex items-center justify-between bg-gradient-to-b from-black/80 to-transparent z-30">
              <div class="flex items-center gap-2 text-sm">
                <span class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></span>
                <span class="font-medium">{{ activeRoom?.room_name || 'Live Tea Room' }}</span>
                <span class="opacity-60">· {{ activeRoom?.location || '' }}</span>
              </div>
              <div class="flex items-center gap-2 text-xs opacity-80">
                <span>🎥 {{ barrages.length + 42 }} watching</span>
              </div>
            </div>

            <!-- Bottom barrage input -->
            <div class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/90 to-transparent z-30">
              <div class="flex gap-2 items-center">
                <div class="flex gap-1">
                  <button v-for="e in EMOJI_SET.slice(0, 6)" :key="e" @click="barrageInput += e" class="text-lg opacity-80 hover:opacity-100 transition pointer-events-auto">{{ e }}</button>
                </div>
                <input v-model="barrageInput" @keyup.enter="sendBarrage" :placeholder="barrageConnected ? 'Send bullet comment...' : 'Connecting...'"
                       class="flex-1 bg-white/10 backdrop-blur border border-white/20 rounded-full px-4 py-2 text-sm text-white placeholder:text-white/50 focus:outline-none focus:border-tea-400 pointer-events-auto" />
                <button @click="sendBarrage" class="px-4 py-2 rounded-full bg-tea-600 hover:bg-tea-500 text-sm font-medium transition pointer-events-auto">发送</button>
              </div>
            </div>
          </div>

          <!-- LiveKit controls row -->
          <div class="mt-4 flex items-center gap-3 flex-wrap">
            <button v-if="!lkConnected" @click="joinLiveKit" :disabled="lkJoinLoading || !hasLiveKitSDK"
                    class="px-5 py-2.5 bg-tea-600 hover:bg-tea-500 rounded-xl font-medium transition disabled:opacity-40 text-sm flex items-center gap-2">
              {{ lkJoinLoading ? 'Connecting...' : '🎥 Join 1:1 Tasting' }}
            </button>
            <template v-else>
              <button @click="toggleCam" class="px-3 py-2 bg-tea-700 hover:bg-tea-600 rounded-lg text-sm transition">📷 Camera</button>
              <button @click="toggleMic" class="px-3 py-2 bg-tea-700 hover:bg-tea-600 rounded-lg text-sm transition">🎤 Mic</button>
              <button @click="leaveLiveKit" class="px-3 py-2 bg-red-600 hover:bg-red-500 rounded-lg text-sm transition">✋ Leave</button>
            </template>
            <button @click="ensureRoomChat" class="px-5 py-2.5 bg-white/10 hover:bg-white/20 border border-white/20 rounded-xl font-medium transition text-sm">💬 Open Chat</button>
            <span v-if="!hasLiveKitSDK" class="text-xs text-amber-400">⚠️ LiveKit SDK 未检测</span>
            <span v-else-if="lkConnected" class="text-xs text-green-400">🟢 已连接 LiveKit</span>
            <span v-else class="text-xs text-tea-400">🟡 可用 LiveKit Cloud / 自建 LiveKit</span>
          </div>

          <!-- Active Rooms -->
          <div class="mt-8">
            <h2 class="font-serif text-2xl mb-4">🎥 Your Active Rooms</h2>
            <div class="grid md:grid-cols-3 gap-4">
              <div v-for="r in [...rooms, ...presets.map(p => ({...p, id: p.id || `p-${p.id}`, room_name: p.name, status: 'live'}))]" :key="r.id"
                   @click="selectRoom(r)"
                   :class="['p-4 rounded-xl cursor-pointer transition border', activeRoom?.room_id === r.room_id ? 'bg-tea-700 border-tea-500' : 'bg-tea-900 border-tea-700 hover:bg-tea-800']">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium">{{ r.room_name || r.name }}</span>
                  <span :class="['text-[10px] px-2 py-0.5 rounded-full', r.status==='live' ? 'bg-green-600 text-white' : 'bg-tea-700 text-tea-300']">{{ r.status }}</span>
                </div>
                <div class="text-xs opacity-70">{{ r.location }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right sidebar: Slow-presets / upcoming -->
        <div class="lg:w-80 flex-shrink-0 space-y-6">
          <div class="bg-tea-900 rounded-2xl border border-tea-700 p-5">
            <h3 class="font-serif text-lg mb-4">🌱 24/7 Garden Cameras</h3>
            <div class="space-y-2">
              <div v-for="p in presets.length ? presets : [
                { name: 'Lincang · Manzhang Village', location: '临沧', status: 'live' },
                { name: 'Bannezhang · Bulang Mountain', location: '勐海', status: 'live' },
                { name: 'Jingmai · Lahu Heritage', location: '澜沧', status: 'live' }
              ]" :key="p.id || p.name" @click="selectPresetAsRoom(p)"
                   class="p-3 rounded-lg bg-tea-800/50 hover:bg-tea-800 cursor-pointer transition border border-transparent hover:border-tea-600">
                <div class="text-sm font-medium">{{ p.name }}</div>
                <div class="text-xs opacity-60 mt-0.5">{{ p.location }} · {{ p.status || '24/7 live' }}</div>
              </div>
            </div>
          </div>

          <div class="bg-tea-900 rounded-2xl border border-tea-700 p-5">
            <h3 class="font-serif text-lg mb-4">🎤 Book a 1:1 Tasting</h3>
            <p class="text-xs text-tea-400 mb-3">Private LiveKit video call with our tea advisor — 45 min from £50.</p>
            <button class="w-full py-2.5 bg-tea-600 hover:bg-tea-500 rounded-xl font-medium transition text-sm">Request Session</button>
          </div>

          <div class="bg-tea-900 rounded-2xl border border-tea-700 p-5">
            <h3 class="font-serif text-lg mb-4">📢 Live Etiquette</h3>
            <ul class="text-xs text-tea-400 space-y-1.5">
              <li>• 30-sec broadcast delay for safety</li>
              <li>• Be kind in bullet comments 💚</li>
              <li>• Tea weights may vary by season</li>
              <li>• All rooms are recorded for quality</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Chat panel (slide-in) -->
    <div v-if="chatOpen" class="fixed inset-0 bg-black/60 z-40 flex justify-end" @click.self="closeChat">
      <div class="w-full sm:w-96 bg-white text-tea-900 h-full flex flex-col shadow-2xl">
        <div class="p-4 border-b bg-tea-50 flex items-center justify-between">
          <h3 class="font-serif text-lg">💬 Room Chat</h3>
          <button @click="closeChat" class="text-tea-600 hover:text-tea-900">✕</button>
        </div>
        <div id="chat-scroll" class="flex-1 overflow-auto p-4 bg-slate-50 space-y-3">
          <div v-for="(m, i) in chatMessages" :key="m.id || i"
               :class="['flex', m.sender_type === 'me' ? 'justify-end' : 'justify-start']">
            <div :class="['max-w-[80%] px-3 py-2 rounded-2xl text-sm',
              m.sender_type === 'me' ? 'bg-tea-700 text-white rounded-tr-sm' : 'bg-white border border-tea-200 rounded-tl-sm']">
              {{ m.content }}
              <div class="text-[10px] mt-0.5 opacity-60">{{ new Date(m.created_at).toLocaleTimeString() }}</div>
            </div>
          </div>
        </div>
        <div class="border-t p-3 bg-white">
          <div class="flex gap-2">
            <input v-model="chatInput" @keyup.enter="sendChat" :placeholder="wsConnected ? 'Type a message...' : 'Connecting...'"
                   class="flex-1 px-4 py-2 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none text-sm" />
            <button @click="sendChat" class="px-4 py-2 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition text-sm">Send</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.barrage-item {
  animation: barrage-fly 12s linear forwards;
  will-change: transform;
}
@keyframes barrage-fly {
  from { transform: translateX(100%); }
  to { transform: translateX(calc(-100vw - 200px)); }
}
</style>
