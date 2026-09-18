<script setup lang="ts">
import { onMounted, onUnmounted, ref, nextTick, computed } from 'vue'
import { api } from '@/api/client'

// ================ 奢侈配色（参考 Rolls-Royce / Château Lafite） ================
// 不碰任何抖音快手式的红/粉/霓虹
//   主视频舞台:  off-black #0B0A09（哑黑，像珠宝丝绒盒内衬）
//   面板背景:    ivory #F8F5EF（象牙白，Hermès 丝巾色）
//   强调色:      champagne gold #C5A572（香槟金，统一所有点缀）
//   辅助:        oak #2A2520 + sand #D8D0C4
// ================

const rooms = ref<any[]>([])
const presets = ref<any[]>([])
const activeRoom = ref<any>(null)
const roomID = computed(() => activeRoom.value?.room_id || `demo-${activeRoom.value?.id || 0}`)
const sessionElapsed = ref(0)
let sessionTimer: number | null = null

// Chat (restrained — 欧洲奢侈品牌不做弹幕墙)
const chatOpen = ref(false)
const chatConvs = ref<any[]>([])
const chatMessages = ref<any[]>([])
const chatInput = ref('')
const chatWS = ref<WebSocket | null>(null)
const wsConnected = ref(false)

// LiveKit
const lkClient = ref<any>(null)
const lkConnected = ref(false)
const lkJoinLoading = ref(false)
const hasLiveKitSDK = ref(false)
const attendeeCount = ref(12) // 克制显示，不像抖音"2.3w 人正在看"

onMounted(async () => {
  try {
    rooms.value = ((await api.get('/live-rooms') as any)?.items || [])
  } catch { rooms.value = [] }
  try {
    presets.value = ((await api.get('/slow-presets') as any)?.items || [])
  } catch { presets.value = [] }

  const DEMO = {
    id: 999, room_id: 'demo-gongfu', room_name: 'Gongfu · Morning Ceremony', room_type: 'live',
    status: 'live', location: 'Yunnan · Lincang · Tea Mountain', camera_rtmp_url: 'rtmp://localhost:1935/live/demo'
  }
  await selectRoom(rooms.value[0] || DEMO)
  sessionTimer = window.setInterval(() => { sessionElapsed.value++ }, 1000)

  try {
    const mod = await import('livekit-client').catch(() => null)
    if (mod?.Room) { hasLiveKitSDK.value = true }
  } catch {}
})

onUnmounted(() => {
  clearInterval(sessionTimer || 0)
  chatWS.value?.close()
  leaveLiveKit()
})

async function selectRoom(room: any) {
  activeRoom.value = room
  if (chatOpen.value) await ensureRoomChat()
}

// —— Chat (restrained, 衬线气泡，像美术馆导览评论) ——
async function ensureRoomChat() {
  let conv: any = null
  try {
    const list = (await api.get('/conversations') as any)?.items || []
    conv = list.find((c: any) => c.title === `room-${roomID.value}`)
  } catch {}
  if (!conv) {
    try { conv = await api.post('/conversations', { title: `room-${roomID.value}`, conversation_type: 'group' }) as any } catch {}
  }
  if (!conv) return
  chatOpen.value = true
  chatConvs.value = [conv]
  try {
    const resp = await api.get(`/conversations/${conv.id}/messages`) as any
    chatMessages.value = resp.items || resp || []
  } catch { chatMessages.value = [] }
  const token = localStorage.getItem('user_token') || localStorage.getItem('staff_token')
  if (!token) return
  const host = window.location.hostname
  chatWS.value?.close()
  chatWS.value = new WebSocket(`ws://${host}:8080/api/v1/ws/im?token=${token}`)
  chatWS.value.onopen = () => { wsConnected.value = true; chatWS.value?.send(JSON.stringify({ type: 'join_conversation', payload: { conversation_id: conv.id } })) }
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

function sendChat() {
  if (!chatInput.value.trim() || !chatConvs.value.length || !chatWS.value) return
  chatWS.value.send(JSON.stringify({
    type: 'send_message',
    payload: { conversation_id: chatConvs.value[0].id, content: chatInput.value.trim(), message_type: 'text' }
  }))
  chatMessages.value.push({ content: chatInput.value.trim(), message_type: 'text', sender_type: 'me', created_at: new Date().toISOString() })
  chatInput.value = ''
}

// —— LiveKit ——
async function joinLiveKit() {
  if (!hasLiveKitSDK.value) { alert('LiveKit SDK unavailable. Configure VITE_LIVEKIT_URL.'); return }
  lkJoinLoading.value = true
  try {
    const resp: any = await api.get('/livekit/token')
    const token = resp?.token || resp?.access_token
    if (!token) throw new Error('no token')
    const { Room, RoomEvent, VideoPresets } = await import('livekit-client')
    const room = new Room({ videoCaptureDefaults: { resolution: VideoPresets.h720 } })
    lkClient.value = room
    room.on(RoomEvent.ParticipantConnected, () => {})
    await room.connect(import.meta.env.VITE_LIVEKIT_URL || 'wss://tea.livekit.cloud', token, { autoSubscribe: true, dynacast: true })
    lkConnected.value = true
    await room.localParticipant.setCameraEnabled(true)
    await room.localParticipant.setMicrophoneEnabled(true)
  } catch (e: any) {
    alert('Connection failed: ' + (e?.message || String(e)))
  } finally { lkJoinLoading.value = false }
}
async function leaveLiveKit() { try { await lkClient.value?.disconnect() } catch {} lkConnected.value = false }
async function toggleMic() { if (!lkClient.value) return; await lkClient.value.localParticipant.setMicrophoneEnabled(!lkClient.value.localParticipant.isMicrophoneEnabled) }
async function toggleCam() { if (!lkClient.value) return; await lkClient.value.localParticipant.setCameraEnabled(!lkClient.value.localParticipant.isCameraEnabled) }

const formattedElapsed = computed(() => {
  const h = Math.floor(sessionElapsed.value / 3600)
  const m = Math.floor((sessionElapsed.value % 3600) / 60)
  const s = sessionElapsed.value % 60
  return `${h.toString().padStart(2,'0')}:${m.toString().padStart(2,'0')}:${s.toString().padStart(2,'0')}`
})
</script>

<template>
  <!-- 
    ========== DESIGN INTENT ==========
    参考 Rolls-Royce 配置器 + Château Lafite + Hermès 丝巾盒:
    
    · 主舞台 off-black (#0B0A09) 像珠宝丝绒盒，只一个焦点
    · 面板象牙白 (#F8F5EF) 衬线字体，像美术馆策展标签
    · 强调只用 champagne gold (#C5A572) — 极克制
    · 无 emoji, 无弹幕浮层, 无 pulsing 红点, 无倒计时, 无"xxx 人正在看"
    · 交互 slow 600ms ease-out, 像高级丝绒移动
    · 衬线 serif 大字 + 小型 CAPS label — typography as architecture
  -->
  <div class="lux-root">
    
    <!-- 
      ===== 顶部品牌带 — 1px 香槟金线分隔, 小型 SERIF CAPS =====
      像爱马仕丝巾包装上的金色箔印
    -->
    <header class="lux-brand-bar">
      <div class="lux-brand-inner">
        <span class="lux-brand-mark">UK · Tea · House</span>
        <span class="lux-brand-sep">·</span>
        <span class="lux-brand-label uppercase-caps">A Bespoke Session</span>
      </div>
    </header>

    <!-- 
      ===== 主舞台: off-black 丝绒盒 =====
      只有一个焦点 — 视频
      所有状态标签都像美术馆策展牌 — 小号 SERIF, 低对比度
    -->
    <main class="lux-stage-wrap">
      <!-- 16:9 视频舞台 -->
      <div class="lux-stage">
        
        <!-- 舞台中央视频区 -->
        <div class="lux-stage-video">
          <video 
            ref="v => null"
            class="lux-video-placeholder"
            autoplay muted playsinline loop
            :poster="`data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 800 450'><rect fill='%230B0A09' width='800' height='450'/><text x='400' y='210' font-family='Georgia' font-size='24' fill='%23C5A572' text-anchor='middle' opacity='0.45'>Yunnan · Tea Mountain</text><text x='400' y='250' font-family='Georgia' font-size='14' fill='%238A8578' text-anchor='middle' opacity='0.6'>rtmp://live.ukteahouse.co.uk/{{ activeRoom?.room_id || 'demo' }}</text></svg>`"
          >
          </video>

          <!-- LiveKit PIP — 右下角小窗, 像望远镜目镜 -->
          <div v-if="lkConnected" class="lux-lk-pip">
            <video autoplay muted playsinline class="lux-pip-video" />
          </div>
        </div>

        <!-- 
          舞台下缘 — 像老电影胶片底部文字带
          serif 小号 + 香槟金 dot — 绝不用红色脉冲
        -->
        <div class="lux-stage-rail">
          <div class="lux-stage-meta">
            <!-- LIVE 指示 — 一个静态金圆点 + SERIF CAPS, 绝不用 pulsing red -->
            <span class="lux-live-pill">
              <span class="lux-live-dot"></span>
              <span class="lux-live-caps uppercase-caps">Live</span>
            </span>
            <span class="lux-rail-sep">|</span>
            <span class="lux-session-time">{{ formattedElapsed }}</span>
            <span class="lux-rail-sep">|</span>
            <span class="lux-garden-label">{{ activeRoom?.location || 'Yunnan' }}</span>
          </div>

          <div class="lux-stage-actions">
            <!-- 1:1 连麦 — "Request Private Tasting", serif ghost 按钮 -->
            <button 
              v-if="!lkConnected" 
              @click="joinLiveKit"
              :disabled="lkJoinLoading"
              class="lux-btn lux-btn-ghost"
            >
              <span class="lux-btn-label uppercase-caps">{{ lkJoinLoading ? 'Connecting...' : 'Request · Private Tasting' }}</span>
            </button>
            <template v-else>
              <button @click="toggleCam" class="lux-btn lux-btn-mini">Cam</button>
              <button @click="toggleMic" class="lux-btn lux-btn-mini">Mic</button>
              <button @click="leaveLiveKit" class="lux-btn lux-btn-mini lux-btn-danger">Leave</button>
            </template>

            <button @click="ensureRoomChat" class="lux-btn lux-btn-ghost">
              <span class="lux-btn-label uppercase-caps">Conversation</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 
        ===== 舞台下方 — 象牙白面板, 衬线大标题 =====
        像爱马仕丝巾盒底部的说明书
      -->
      <section class="lux-stage-caption">
        <h1 class="lux-caption-title">{{ activeRoom?.room_name || 'Gongfu · Morning Ceremony' }}</h1>
        <p class="lux-caption-sub">A curated moment from the tea mountains of Yunnan · Slow Camera · 24/7</p>
      </section>
    </main>

    <!-- 
      ===== 象牙白侧栏: 正在进行的场次 + 茶园列表 =====
      不叫 "Upcoming Tasting" — 叫 "Curated Sessions"
    -->
    <aside class="lux-sidebar">
      <div class="lux-panel">
        <h2 class="lux-panel-title uppercase-caps">Curated Sessions</h2>
        <ul class="lux-list">
          <li v-for="r in [...rooms, ...presets.map(p => ({...p, id: p.id || `p-${p.id}`, room_name: p.name, room_id: p.id || `p-${p.id}`}))]" 
              :key="r.id"
              @click="selectRoom(r)"
              :class="['lux-list-item', activeRoom?.room_id === r.room_id ? 'lux-list-item--active' : '']">
            <span class="lux-list-status">{{ r.status === 'live' ? '·' : '—' }}</span>
            <span class="lux-list-label">{{ r.room_name || r.name }}</span>
            <span class="lux-list-meta">{{ r.location }}</span>
          </li>
          <li v-if="!rooms.length && !presets.length" class="lux-list-empty uppercase-caps">
            Sessions being curated
          </li>
        </ul>
      </div>

      <div class="lux-panel">
        <h2 class="lux-panel-title uppercase-caps">Garden Region</h2>
        <address class="lux-address">
          <span>Yunnan · Lincang</span>
          <span>Tea Mountains · High Elevation</span>
          <span>Province of Yunnan · China</span>
        </address>
      </div>

      <div class="lux-panel lux-panel-minimal">
        <p class="lux-panel-note">
          <em>Tea should not be rushed.</em> 
          Cameras operate on local time. The advisory reserves the right to reschedule.
        </p>
      </div>
    </aside>

    <!-- 
      ===== Conversation 面板 — slide-in, 衬线气泡 =====
      绝不像抖音弹幕 — 这是美术馆导览器式的评论区
    -->
    <transition name="lux-panel-slide">
      <aside v-if="chatOpen" class="lux-chat">
        <div class="lux-chat-head">
          <h3 class="lux-chat-title uppercase-caps">Conversation</h3>
          <span class="lux-chat-count">{{ attendeeCount }} guests · {{ wsConnected ? 'Connected' : 'Standby' }}</span>
          <button @click="chatOpen = false" class="lux-chat-close">×</button>
        </div>
        <div id="chat-scroll" class="lux-chat-body">
          <div v-for="(m, i) in chatMessages" :key="m.id || i" :class="['lux-bubble', m.sender_type === 'me' ? 'lux-bubble--me' : 'lux-bubble--them']">
            <p class="lux-bubble-text">{{ m.content }}</p>
            <span class="lux-bubble-time">{{ new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}</span>
          </div>
          <p v-if="!chatMessages.length" class="lux-chat-empty uppercase-caps">A quiet room · Wait for the first word</p>
        </div>
        <div class="lux-chat-foot">
          <input v-model="chatInput" @keyup.enter="sendChat" placeholder="Speak to the room..." class="lux-chat-input" />
          <button @click="sendChat" class="lux-btn lux-btn-mini">Send</button>
        </div>
      </aside>
    </transition>
  </div>
</template>

<!-- 
  ========= 奢侈风格 CSS =========
  · 只用 serif (Georgia/Cormorant/Playfair) + Inter
  · 无任何 emoji 动画, 无 hover lift, 无 bounce
  · 过渡 600ms ease-out (像高级丝绒), 不是 150ms cubic-bezier
  · Hairline borders (0.5px) 像爱马仕丝巾的金线
  · 颜色只用 off-black / ivory / champagne gold / oak
-->
<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;500;600;700&family=Inter:wght@300;400;500&display=swap');

:root {
  --off-black:   #0B0A09;
  --ivory:       #F8F5EF;
  --ivory-dark:  #EDE8DD;
  --champagne:   #C5A572;
  --champagne-dim: rgba(197, 165, 114, 0.5);
  --oak:         #2A2520;
  --sand:        #D8D0C4;
  --ink-on-ivory:#1a1714;
}

.lux-root {
  min-height: 100vh;
  background: var(--ivory);
  font-family: 'Inter', -apple-system, sans-serif;
  color: var(--ink-on-ivory);
  display: grid;
  grid-template-columns: 1fr 340px;
  grid-template-rows: auto 1fr;
  gap: 0;
}

/* ========== 顶部品牌带 ========== */
.lux-brand-bar {
  grid-column: 1 / -1;
  background: var(--off-black);
  color: var(--champagne);
  padding: 18px 40px;
  border-bottom: 1px solid var(--champagne-dim);
  display: flex;
  align-items: center;
}
.lux-brand-inner { display: flex; align-items: baseline; gap: 16px; }
.lux-brand-mark {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 22px; font-weight: 500;
  letter-spacing: 0.15em;
}
.lux-brand-sep { color: var(--champagne-dim); font-size: 12px; }
.lux-brand-label { font-size: 11px; font-weight: 300; letter-spacing: 0.3em; }
.uppercase-caps { text-transform: uppercase; letter-spacing: 0.28em; }

/* ========== 主舞台 ========== */
.lux-stage-wrap { padding: 32px 40px 40px; }
.lux-stage {
  position: relative;
  background: var(--off-black);
  border: 1px solid rgba(197,165,114,0.18);
  overflow: hidden;
  aspect-ratio: 16 / 9;
  border-radius: 2px; /* 绝不用 16px — 奢侈品牌用微圆角或直角 */
}
.lux-stage-video { width: 100%; height: 100%; position: relative; }
.lux-video-placeholder {
  width: 100%; height: 100%; object-fit: cover;
  filter: contrast(0.95) brightness(0.95);
}

/* LiveKit PIP — 右下角望远镜窗 */
.lux-lk-pip {
  position: absolute; right: 24px; bottom: 84px;
  width: 140px; height: 90px;
  border: 1px solid var(--champagne-dim);
  border-radius: 2px; overflow: hidden;
  background: var(--oak);
}
.lux-pip-video { width: 100%; height: 100%; object-fit: cover; }

/* ========== 舞台底部 rail ========== */
.lux-stage-rail {
  position: absolute; left: 0; right: 0; bottom: 0;
  padding: 16px 28px;
  background: linear-gradient(180deg, transparent 0%, rgba(11,10,9,0.9) 60%, rgba(11,10,9,0.97) 100%);
  display: flex; align-items: center; justify-content: space-between;
  border-top: 1px solid rgba(197,165,114,0.15);
}
.lux-stage-meta { display: flex; align-items: center; gap: 14px; }

/* LIVE pill — 金圆点, 绝不用 pulsing red */
.lux-live-pill { display: flex; align-items: center; gap: 8px; }
.lux-live-dot {
  width: 6px; height: 6px; border-radius: 50%;
  background: var(--champagne);
}
.lux-live-caps {
  font-family: 'Cormorant Garamond', serif;
  font-size: 13px; font-weight: 500;
  letter-spacing: 0.32em;
  color: var(--champagne);
}
.lux-rail-sep { color: rgba(216,208,196,0.35); font-size: 10px; }
.lux-session-time {
  font-family: 'Cormorant Garamond', serif;
  font-size: 14px; font-weight: 400;
  color: var(--sand);
  letter-spacing: 0.08em;
}
.lux-garden-label {
  font-family: 'Cormorant Garamond', serif;
  font-size: 14px; color: var(--sand);
}

/* ========== 按钮 — ghost, hairline border, 600ms ========== */
.lux-stage-actions { display: flex; align-items: center; gap: 10px; }
.lux-btn {
  font-family: 'Inter', sans-serif;
  font-size: 11px; font-weight: 400;
  padding: 10px 20px;
  border-radius: 2px;
  border: 1px solid rgba(197,165,114,0.35);
  background: transparent;
  color: var(--champagne);
  cursor: pointer;
  transition: all 600ms ease-out;
  letter-spacing: 0.15em;
}
.lux-btn:hover {
  background: rgba(197,165,114,0.08);
  border-color: var(--champagne);
}
.lux-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.lux-btn-label { font-family: 'Cormorant Garamond', serif; font-size: 13px; letter-spacing: 0.22em; }
.lux-btn-mini { padding: 8px 14px; font-size: 11px; }
.lux-btn-danger { border-color: rgba(180,100,90,0.4); color: #B4645A; }
.lux-btn-danger:hover { background: rgba(180,100,90,0.08); border-color: #B4645A; }

/* ========== 舞台下方 caption ========== */
.lux-stage-caption { margin-top: 28px; }
.lux-caption-title {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 44px; font-weight: 500;
  letter-spacing: -0.01em;
  color: var(--off-black);
  margin: 0 0 8px;
  line-height: 1.15;
}
.lux-caption-sub {
  font-family: 'Inter', sans-serif;
  font-size: 13px; font-weight: 300;
  color: #6b6459;
  letter-spacing: 0.04em;
  margin: 0;
}

/* ========== 侧栏 (象牙白) ========== */
.lux-sidebar {
  padding: 32px 32px 40px 0;
  border-left: 1px solid rgba(11,10,9,0.08);
  display: flex; flex-direction: column; gap: 32px;
}
.lux-panel { }
.lux-panel-title {
  font-family: 'Inter', sans-serif;
  font-size: 10px; font-weight: 400;
  color: var(--champagne);
  margin: 0 0 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--sand);
}
.lux-panel-minimal { margin-top: auto; }

/* 列表项 — hover 只是极细金线下移, 绝不背景翻色 */
.lux-list { list-style: none; padding: 0; margin: 0; }
.lux-list-item {
  padding: 14px 0;
  cursor: pointer;
  border-bottom: 0.5px solid rgba(11,10,9,0.05);
  display: flex; flex-direction: column; gap: 2px;
  transition: border-color 600ms ease-out;
}
.lux-list-item:hover, .lux-list-item--active { border-bottom-color: var(--champagne); }
.lux-list-item:hover .lux-list-label { color: var(--off-black); }
.lux-list-status {
  color: var(--champagne);
  font-size: 11px; letter-spacing: 0.1em;
  margin-bottom: 4px;
}
.lux-list-label {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 20px; font-weight: 400;
  color: #44403a;
  transition: color 600ms ease-out;
}
.lux-list-meta {
  font-size: 11px; color: #8a8578; font-weight: 300; letter-spacing: 0.05em;
}
.lux-list-empty { font-size: 11px; color: #a8a195; padding: 20px 0; }

.lux-address {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 16px; font-weight: 400;
  color: #44403a;
  font-style: normal;
  line-height: 1.7;
}

.lux-panel-note {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 14px; font-style: italic;
  color: #6b6459; line-height: 1.7;
  border-top: 1px solid var(--sand);
  padding-top: 18px;
}

/* ========== Conversation slide-in ========== */
.lux-chat {
  position: fixed; top: 0; right: 0; bottom: 0;
  width: 380px; max-width: 90vw;
  background: var(--ivory);
  border-left: 1px solid var(--sand);
  display: flex; flex-direction: column;
  z-index: 50;
}
.lux-panel-slide-enter-active, .lux-panel-slide-leave-active { transition: transform 600ms ease-out; }
.lux-panel-slide-enter-from, .lux-panel-slide-leave-to { transform: translateX(100%); }

.lux-chat-head {
  padding: 24px 28px 18px;
  border-bottom: 1px solid var(--sand);
  display: flex; flex-direction: column; gap: 6px;
}
.lux-chat-title {
  font-family: 'Inter', sans-serif;
  font-size: 10px; font-weight: 400; color: var(--champagne);
  margin: 0;
}
.lux-chat-count { font-size: 11px; color: #8a8578; letter-spacing: 0.05em; }
.lux-chat-close {
  position: absolute; top: 18px; right: 22px;
  background: transparent; border: none;
  font-size: 22px; color: #a8a195; cursor: pointer;
  line-height: 1;
  font-family: 'Inter', sans-serif;
  transition: color 600ms ease-out;
}
.lux-chat-close:hover { color: var(--off-black); }

.lux-chat-body {
  flex: 1; overflow-y: auto;
  padding: 20px 24px;
  background: var(--ivory);
  scroll-behavior: smooth;
}
.lux-bubble {
  margin-bottom: 18px;
  max-width: 85%;
  padding: 12px 18px;
  border-radius: 2px;
  border: 1px solid rgba(11,10,9,0.08);
  background: var(--ivory);
  transition: border-color 600ms ease-out;
}
.lux-bubble--me { margin-left: auto; border-color: var(--champagne-dim); background: #FBF9F4; }
.lux-bubble-text {
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 16px; line-height: 1.5;
  color: var(--ink-on-ivory);
  margin: 0;
}
.lux-bubble-time { font-size: 10px; color: #a8a195; letter-spacing: 0.08em; margin-top: 6px; display: block; }
.lux-chat-empty { font-size: 11px; color: #a8a195; text-align: center; padding: 40px 0; }

.lux-chat-foot {
  padding: 18px 24px 24px;
  border-top: 1px solid var(--sand);
  background: var(--ivory);
  display: flex; gap: 12px;
}
.lux-chat-input {
  flex: 1;
  background: transparent;
  border: none; border-bottom: 1px solid var(--sand);
  padding: 8px 0;
  font-family: 'Cormorant Garamond', Georgia, serif;
  font-size: 16px; color: var(--ink-on-ivory);
  outline: none;
  transition: border-color 600ms ease-out;
}
.lux-chat-input:focus { border-bottom-color: var(--champagne); }
.lux-chat-input::placeholder { color: #a8a195; }

/* ========== 响应式 ========== */
@media (max-width: 960px) {
  .lux-root { grid-template-columns: 1fr; }
  .lux-sidebar {
    padding: 24px;
    border-left: none;
    border-top: 1px solid var(--sand);
  }
  .lux-stage-wrap { padding: 20px; }
  .lux-brand-bar { padding: 14px 20px; }
  .lux-caption-title { font-size: 32px; }
  .lux-stage-rail { flex-direction: column; gap: 12px; align-items: flex-start; }
  .lux-stage-actions { width: 100%; justify-content: flex-start; flex-wrap: wrap; }
}
</style>
