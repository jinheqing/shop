<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const room = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const visibility = ref<string>('registered')

const roomId = computed(() => route.params.id as string | undefined)

// ============ LiveKit ============
const lkClient = ref<any>(null)
const lkConnected = ref(false)
const lkJoinLoading = ref(false)
const hasLiveKitSDK = ref(false)
const remoteVideoEl = ref<HTMLVideoElement | null>(null)

// ============ 字幕 / 翻译 ============
const asrWS = ref<WebSocket | null>(null)
const asrConnected = ref(false)
const microphoneEnabled = ref(false)
const audioStream = ref<MediaStream | null>(null)
const audioProcessor = ref<ScriptProcessorNode | null>(null)
const audioSource = ref<MediaStreamAudioSourceNode | null>(null)

// 字幕队列（时间戳有序，latest 在最后）
const captions = ref<Array<{
  id: number
  type: 'partial' | 'final'
  speaker: 'host' | 'me' | 'guest'
  text: string
  translation?: string
  lang: string
  ts: number
}>>([])
let captionSeq = 0

const showTranslation = ref(true)
const sourceLang = ref<'zh' | 'en' | 'auto'>('auto')

const targetLang = computed(() => {
  if (sourceLang.value === 'zh') return 'en'
  if (sourceLang.value === 'en') return 'zh'
  return 'auto'
})

// ============ IM WebSocket（订阅主播字幕广播 + 连麦信令） ============
const imWS = ref<WebSocket | null>(null)
const imConnected = ref(false)

// ============ 连麦状态 ============
const linked = ref(false)           // 观众是否正在连麦
const linkSessionId = ref<string>('')
const linkRequested = ref(false)

async function openIM() {
  const token = localStorage.getItem('user_token') || localStorage.getItem('staff_token')
  if (!token) return
  const host = window.location.host
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
  imWS.value = new WebSocket(`${proto}://${host}/ws/im?token=${token}`)
  imWS.value.onopen = () => {
    imConnected.value = true
    imWS.value?.send(JSON.stringify({ type: 'join_room', payload: { room_id: room.value?.room_id } }))
  }
  imWS.value.onclose = () => { imConnected.value = false }
  imWS.value.onmessage = (ev) => {
    try {
      const env = JSON.parse(ev.data)
      if (env.type === 'barrage') handleBarrage(env.payload)
    } catch {}
  }
}

function closeIM() { try { imWS.value?.close() } catch {}; imWS.value = null; imConnected.value = false }

function sendBarrage(payload: any) {
  if (!imWS.value || imWS.value.readyState !== WebSocket.OPEN) return
  const env = { type: 'barrage', payload: { ...payload, room_id: room.value?.room_id, ts: Date.now() } }
  imWS.value.send(JSON.stringify(env))
}

function handleBarrage(p: any) {
  const st = p.subtype || 'chat'
  if (st === 'subtitle') {
    // 主播/连麦者发来的字幕广播 → 追加到字幕队列
    const speaker = p.speaker === 'host' ? 'host' : 'guest'
    pushCaption({
      speaker, type: p.transcript_type === 'partial' ? 'partial' : 'final',
      text: p.text, translation: p.translation, lang: p.lang || 'auto', ts: p.ts || Date.now(),
    })
  } else if (st === 'link_invite') {
    // 主播邀请观众连麦 — 这里用"观众请求主播接受"的模型，主播拒绝/接受用 link_reject/link_accept 广播
    // 观众作为"被邀请者"，主播是主动方 — 如果是反向场景（观众请求主播邀请），主播会发 invite
  } else if (st === 'link_accept' && p.target_user_id && p.target_user_id === 0 /* 观众不知道自己 ID，简单处理 */) {
    // 简化：收到 link_accept（无 target 或 target 为 0 时）→ 视为主播同意了观众的连麦请求
    handleLinkAccepted(p.link_session_id)
  } else if (st === 'link_reject') {
    linkRequested.value = false
  } else if (st === 'link_end') {
    linked.value = false
    linkSessionId.value = ''
  }
}

// 观众请求连麦 → 主播会收到 link_invite（观众作为"请求方"，主播作为"被请求方"）
// 这里复用同一 barrage subtype，主播收到 invite 后决定 accept/reject
async function requestLinkMic() {
  linkRequested.value = true
  linkSessionId.value = 'l-' + Math.random().toString(36).slice(2, 10)
  sendBarrage({ subtype: 'link_invite', content: 'request', link_session_id: linkSessionId.value, nickname: 'Viewer' })
  // 同时自己也 publish audio 到 LiveKit（主播 accept 后就能听到）
  try {
    if (lkClient.value) await lkClient.value.localParticipant.setMicrophoneEnabled(true)
  } catch {}
}

async function handleLinkAccepted(sessionId: string) {
  linkRequested.value = false
  linked.value = true
  linkSessionId.value = sessionId
  try {
    if (lkClient.value) {
      await lkClient.value.localParticipant.setMicrophoneEnabled(true)
      await lkClient.value.localParticipant.setCameraEnabled(true)
    }
  } catch {}
}

async function endLinkMic() {
  sendBarrage({ subtype: 'link_end', content: 'end', link_session_id: linkSessionId.value })
  linked.value = false
  linkRequested.value = false
  linkSessionId.value = ''
  try {
    if (lkClient.value) await lkClient.value.localParticipant.setMicrophoneEnabled(false)
  } catch {}
}

const isLoggedIn = () => !!localStorage.getItem('user_token')

onMounted(async () => {
  if (!roomId.value) { error.value = 'Missing room ID'; loading.value = false; return }
  try {
    const data: any = await api.get(`/live-rooms/${roomId.value}`)
    room.value = data
    visibility.value = data?.visibility || 'registered'
  } catch (err: any) {
    const status = err?.code || 0
    if (status === 403 || (typeof err === 'string' && err.includes('not accessible'))) {
      visibility.value = 'restricted'
      error.value = 'restricted'
    } else if (status === 401) {
      visibility.value = 'registered'
      error.value = 'login_required'
    } else {
      error.value = err?.message || 'Failed to load broadcast'
    }
  } finally {
    loading.value = false
  }

  try {
    const mod = await import('livekit-client').catch(() => null)
    if (mod?.Room) { hasLiveKitSDK.value = true }
  } catch {}
})

onUnmounted(() => {
  leaveLiveKit()
  closeASR()
  closeIM()
})

function goLogin() {
  const redirect = encodeURIComponent(route.fullPath)
  router.push({ path: '/magic-link', query: { redirect } })
}

// ============ LiveKit join (subscribe-only viewer) ============
async function joinLiveKit() {
  if (!hasLiveKitSDK.value || !room.value) { alert('LiveKit SDK unavailable'); return }
  lkJoinLoading.value = true
  try {
    const identity = localStorage.getItem('user_token')
      ? 'viewer-' + room.value.room_id
      : 'guest-' + room.value.room_id
    const resp: any = await api.post('/livekit/token', { room_name: room.value.room_id, identity })
    const token = resp?.token
    if (!token) throw new Error('no token')
    const { Room, RoomEvent, VideoPresets } = await import('livekit-client')
    const r = new Room({ autoSubscribe: true, dynacast: true, videoCaptureDefaults: { resolution: VideoPresets.h720 } })
    lkClient.value = r

    // 订阅远端视频轨（host 发布的）
    r.on(RoomEvent.TrackSubscribed, (track: any, _pub: any, participant: any) => {
      if (track.kind === 'video' && remoteVideoEl.value) {
        const el = track.attach()
        remoteVideoEl.value.appendChild(el)
      }
    })
    r.on(RoomEvent.TrackUnsubscribed, (track: any) => {
      try { track.detach() } catch {}
    })

    await r.connect(import.meta.env.VITE_LIVEKIT_URL || 'wss://tea.livekit.cloud', token)
    lkConnected.value = true
    // LiveKit join 成功 → 打开 IM WS 订阅 barrage（主播字幕广播 + 连麦信令）
    if (isLoggedIn()) openIM()
  } catch (e: any) {
    alert('LiveKit connect failed: ' + (e?.message || String(e)))
  } finally { lkJoinLoading.value = false }
}

async function leaveLiveKit() {
  try {
    await lkClient.value?.disconnect()
  } catch {}
  lkConnected.value = false
  // detach 所有挂在 video 元素上的旧 track
  if (remoteVideoEl.value) {
    remoteVideoEl.value.innerHTML = ''
  }
}

// ============ ASR + 翻译插件 ============
// 前端采集麦克风 PCM → tea-translate /translate/asr-stream WS
// 得到 partial（实时识别）和 final（识别+翻译）字幕

function asrWSURL() {
  // tea-translate 直接暴露 /translate/asr-stream
  // 或走 Go 代理：这里直接直连 translate 服务（默认 8090）
  const host = window.location.hostname || 'localhost'
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
  return `${proto}://${host}:8090/translate/asr-stream`
}

async function openASR() {
  if (asrWS.value?.readyState === WebSocket.OPEN) return
  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: { channelCount: 1, sampleRate: 16000, echoCancellation: true, noiseSuppression: true },
      video: false,
    })
    audioStream.value = stream
    microphoneEnabled.value = true

    // 建立到 tea-translate 的 WebSocket
    asrWS.value = new WebSocket(asrWSURL())
    asrWS.value.binaryType = 'arraybuffer'

    asrWS.value.onopen = () => {
      asrConnected.value = true
      // 首帧发配置
      const cfg = { target_lang: targetLang.value === 'auto' ? 'auto' : targetLang.value }
      asrWS.value?.send(JSON.stringify(cfg))
      startAudioCapture(stream)
    }
    asrWS.value.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data)
        if (msg.type === 'partial') {
          pushCaption({ type: 'partial', speaker: 'me', text: msg.text, lang: 'auto', ts: Date.now() })
        } else if (msg.type === 'final') {
          pushCaption({
            type: 'final', speaker: 'me',
            text: msg.text,
            translation: msg.translation,
            lang: 'auto', ts: Date.now(),
          })
        }
      } catch {}
    }
    asrWS.value.onclose = () => {
      asrConnected.value = false
      stopAudioCapture()
    }
    asrWS.value.onerror = () => {
      asrConnected.value = false
    }
  } catch (e: any) {
    alert('Microphone access failed: ' + (e?.message || String(e)))
  }
}

function closeASR() {
  stopAudioCapture()
  try { asrWS.value?.close() } catch {}
  asrWS.value = null
  asrConnected.value = false
  microphoneEnabled.value = false
}

async function toggleASR() {
  if (asrConnected.value) closeASR()
  else await openASR()
}

function startAudioCapture(stream: MediaStream) {
  if (!asrWS.value) return
  const ctx = new (window.AudioContext || (window as any).webkitAudioContext)({ sampleRate: 16000 })
  const source = ctx.createMediaStreamSource(stream)
  const processor = ctx.createScriptProcessor(4096, 1, 1)

  processor.onaudioprocess = (ev) => {
    if (asrWS.value?.readyState !== WebSocket.OPEN) return
    const input = ev.inputBuffer.getChannelData(0)
    // float32 → int16 PCM
    const int16 = new Int16Array(input.length)
    for (let i = 0; i < input.length; i++) {
      const s = Math.max(-1, Math.min(1, input[i]))
      int16[i] = s < 0 ? s * 0x8000 : s * 0x7FFF
    }
    asrWS.value.send(int16.buffer)
  }

  source.connect(processor)
  processor.connect(ctx.destination) // keep flowing
  audioSource.value = source
  audioProcessor.value = processor
}

function stopAudioCapture() {
  try { audioProcessor.value?.disconnect() } catch {}
  try { audioSource.value?.disconnect() } catch {}
  audioProcessor.value = null
  audioSource.value = null
  audioStream.value?.getTracks().forEach(t => t.stop())
  audioStream.value = null
  microphoneEnabled.value = false
}

function pushCaption(cap: Omit<typeof captions.value[number], 'id'>) {
  captions.value.push({ id: ++captionSeq, ...cap })
  // 只保留最近 50 条 final + 最近 3 条 partial，避免刷屏
  const finals = captions.value.filter(c => c.type === 'final').slice(-50)
  const partials = captions.value.filter(c => c.type === 'partial').slice(-3)
  captions.value = [...finals, ...partials]
}

function clearCaptions() { captions.value = [] }
</script>

<template>
  <div class="detail-root">
    <!-- ========== 加载中 ========== -->
    <div v-if="loading" class="lux-loading">
      <span class="lux-loading-mark uppercase-caps">Curating · Session</span>
    </div>

    <!-- ========== Restricted + 未登录 ========== -->
    <div v-else-if="error === 'restricted' && !isLoggedIn()" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">This broadcast is not publicly accessible</h1>
        <p class="lux-denied-body">
          Your advisor has reserved this session for specific members.
          Sign in with your registered email to request entry.
        </p>
        <button class="lux-btn lux-btn-primary" @click="goLogin">
          <span class="uppercase-caps">Request · Access</span>
        </button>
        <p class="lux-denied-sub">A member of our advisory team will respond within 24 hours.</p>
      </div>
    </div>

    <!-- ========== Restricted + 已登录但无权限 ========== -->
    <div v-else-if="error === 'restricted' && isLoggedIn()" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">Private · Broadcast</h1>
        <p class="lux-denied-body">
          Your advisor has not yet granted access to this session.
          Reach out directly via the Conversation panel — they will add you to the guest list.
        </p>
        <button class="lux-btn lux-btn-ghost" @click="router.push('/chat')">
          <span class="uppercase-caps">Speak with · Your Advisor</span>
        </button>
      </div>
    </div>

    <!-- ========== Registered + 未登录 ========== -->
    <div v-else-if="error === 'login_required'" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">Sign · In · Required</h1>
        <p class="lux-denied-body">
          This session is available to registered members. Sign in with your magic link to view.
        </p>
        <button class="lux-btn lux-btn-primary" @click="goLogin">
          <span class="uppercase-caps">Sign · In</span>
        </button>
      </div>
    </div>

    <!-- ========== 正常显示 ========== -->
    <div v-else-if="room" class="lux-live">
      <header class="lux-live-bar">
        <span class="lux-live-dot"></span>
        <span class="lux-live-title uppercase-caps">{{ room.room_name || 'Private Session' }}</span>
        <span class="lux-live-visibility lux-live-visibility--{{ visibility }} uppercase-caps">{{ visibility }}</span>
      </header>

      <!-- 主视频舞台 -->
      <div class="lux-live-stage">
        <!-- LiveKit 远端视频 -->
        <div v-if="lkConnected" ref="remoteVideoEl" class="lux-video-host"></div>

        <!-- PIP: 本地摄像头（只在 host 权限时有用） -->
        <video v-show="lkConnected" autoplay muted playsinline class="lux-video-pip" />

        <!-- 未连接时的占位 -->
        <div v-if="!lkConnected" class="lux-video-placeholder">
          <span class="lux-placeholder-text">Live stream · Awaiting connection</span>
        </div>

        <!-- 字幕叠层（舞台下缘） -->
        <div v-if="captions.length" class="lux-caption-overlay">
          <template v-for="(cap, i) in captions.slice(-3)" :key="cap.id">
            <div :class="['lux-caption-line', cap.type === 'partial' ? 'lux-caption-line--partial' : '']">
              <span class="lux-caption-src">{{ cap.text }}</span>
              <span v-if="showTranslation && cap.translation" class="lux-caption-dst">— {{ cap.translation }}</span>
            </div>
          </template>
        </div>
      </div>

      <!-- 控制条 -->
      <div class="lux-controls">
        <button v-if="!lkConnected" @click="joinLiveKit" :disabled="lkJoinLoading" class="lux-btn lux-btn-primary">
          <span class="uppercase-caps">{{ lkJoinLoading ? 'Connecting...' : 'Join · Broadcast' }}</span>
        </button>
        <template v-else>
          <button @click="leaveLiveKit" class="lux-btn lux-btn-ghost">
            <span class="uppercase-caps">Leave</span>
          </button>

          <!-- 连麦按钮 -->
          <template v-if="!linked">
            <button v-if="!linkRequested" @click="requestLinkMic" class="lux-btn lux-btn-ghost">
              <span class="uppercase-caps">Request · Link Mic</span>
            </button>
            <button v-else disabled class="lux-btn">
              <span class="uppercase-caps">Awaiting · Host</span>
            </button>
          </template>
          <button v-else @click="endLinkMic" class="lux-btn lux-btn-active">
            <span class="uppercase-caps">🔗 In · Call · End</span>
          </button>
        </template>

        <!-- 翻译插件按钮 -->
        <button :class="['lux-btn', asrConnected ? 'lux-btn-active' : 'lux-btn-ghost']" @click="toggleASR">
          <span class="uppercase-caps">{{ asrConnected ? 'Listening · Live' : 'Translate · Speech' }}</span>
        </button>

        <div v-if="asrConnected" class="lux-lang-switch">
          <button :class="['lux-lang-btn', sourceLang === 'auto' && 'lux-lang-btn--active']" @click="sourceLang = 'auto'">Auto</button>
          <button :class="['lux-lang-btn', sourceLang === 'zh' && 'lux-lang-btn--active']" @click="sourceLang = 'zh'">中文</button>
          <button :class="['lux-lang-btn', sourceLang === 'en' && 'lux-lang-btn--active']" @click="sourceLang = 'en'">EN</button>
          <button class="lux-lang-btn" @click="showTranslation = !showTranslation">
            {{ showTranslation ? 'Hide Trans' : 'Show Trans' }}
          </button>
          <button class="lux-lang-btn" @click="clearCaptions">Clear</button>
        </div>

        <span class="lux-status-pill" :class="asrConnected ? 'lux-status-pill--on' : 'lux-status-pill--off'">
          <span class="lux-status-dot"></span>
          {{ asrConnected ? 'ASR Streaming' : 'ASR Offline' }}
        </span>
      </div>

      <p class="lux-live-desc">{{ room.description || 'A curated private session.' }}</p>
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;500;600&family=Inter:wght@300;400&display=swap');

.detail-root {
  min-height: 100vh;
  background: #F8F5EF;
  font-family: 'Inter', sans-serif;
  color: #1a1714;
  padding: 48px 32px;
  max-width: 1200px;
  margin: 0 auto;
}
.uppercase-caps { text-transform: uppercase; letter-spacing: 0.28em; }

/* ====== Loading / Denied ====== */
.lux-loading { height: 60vh; display: flex; align-items: center; justify-content: center; }
.lux-loading-mark { font-family: 'Cormorant Garamond', serif; font-size: 14px; color: #C5A572; }

.lux-denied { height: 60vh; display: flex; align-items: center; justify-content: center; }
.lux-denied-card {
  text-align: center; padding: 64px 48px;
  border: 1px solid #D8D0C4; border-radius: 2px;
  background: #F8F5EF; max-width: 520px;
}
.lux-denied-title { font-family: 'Cormorant Garamond', serif; font-size: 16px; font-weight: 500; color: #C5A572; margin: 0 0 24px; }
.lux-denied-body { font-size: 15px; line-height: 1.7; color: #44403a; font-family: 'Cormorant Garamond', serif; margin: 0 0 32px; }
.lux-denied-sub { font-size: 12px; color: #8a8578; margin-top: 20px; }

/* ====== Buttons ====== */
.lux-btn {
  padding: 14px 32px; border-radius: 2px;
  font-size: 11px; font-weight: 400; cursor: pointer;
  transition: all 600ms ease-out; letter-spacing: 0.15em;
  font-family: 'Inter', sans-serif;
}
.lux-btn-primary { background: #0B0A09; color: #C5A572; border: 1px solid #0B0A09; }
.lux-btn-primary:hover { background: #2A2520; }
.lux-btn-ghost { background: transparent; color: #0B0A09; border: 1px solid #0B0A09; }
.lux-btn-ghost:hover { background: rgba(11,10,9,0.05); }
.lux-btn-active { background: #4a7c59; color: #F8F5EF; border: 1px solid #4a7c59; }
.lux-btn-label { font-family: 'Cormorant Garamond', serif; font-size: 13px; }

/* ====== Live Header ====== */
.lux-live-bar {
  display: flex; align-items: center; gap: 14px;
  padding-bottom: 20px; border-bottom: 1px solid #D8D0C4; margin-bottom: 24px;
}
.lux-live-dot { width: 8px; height: 8px; border-radius: 50%; background: #C5A572; }
.lux-live-title { font-family: 'Cormorant Garamond', serif; font-size: 14px; color: #0B0A09; }
.lux-live-visibility { margin-left: auto; font-size: 10px; color: #8a8578; }
.lux-live-visibility--restricted { color: #B4645A; }
.lux-live-visibility--public { color: #4a7c59; }

/* ====== Video Stage ====== */
.lux-live-stage {
  position: relative;
  aspect-ratio: 16 / 9;
  background: #0B0A09;
  border: 1px solid rgba(197,165,114,0.18);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 20px;
}
.lux-video-placeholder {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
}
.lux-placeholder-text {
  font-family: 'Cormorant Garamond', serif;
  color: #C5A572; font-size: 18px; letter-spacing: 0.15em; opacity: 0.5;
}
.lux-video-host {
  width: 100%; height: 100%; object-fit: cover;
}
.lux-video-pip {
  position: absolute; right: 16px; bottom: 80px;
  width: 120px; height: 80px; border: 1px solid rgba(197,165,114,0.4);
  border-radius: 2px; object-fit: cover; background: #2A2520;
}

/* ====== 字幕叠层 ====== */
.lux-caption-overlay {
  position: absolute; left: 0; right: 0; bottom: 0;
  padding: 20px 32px 60px;
  background: linear-gradient(transparent 0%, rgba(11,10,9,0.85) 100%);
  pointer-events: none;
  display: flex; flex-direction: column; gap: 4px;
  text-align: center;
}
.lux-caption-line {
  font-family: 'Cormorant Garamond', serif;
  color: #F8F5EF; font-size: 22px; line-height: 1.3;
  text-shadow: 0 2px 16px rgba(0,0,0,0.6);
  transition: opacity 300ms;
}
.lux-caption-line--partial { opacity: 0.6; font-style: italic; }
.lux-caption-dst {
  color: #C5A572; font-size: 18px; letter-spacing: 0.02em;
  margin-left: 8px;
}

/* ====== Controls ====== */
.lux-controls {
  display: flex; flex-wrap: wrap; gap: 10px; align-items: center;
  padding-bottom: 20px; border-bottom: 1px solid #D8D0C4; margin-bottom: 20px;
}
.lux-lang-switch {
  display: flex; gap: 4px; margin-left: auto;
}
.lux-lang-btn {
  padding: 6px 14px; font-size: 10px; letter-spacing: 0.12em;
  border: 1px solid #D8D0C4; background: transparent; color: #44403a;
  border-radius: 2px; cursor: pointer; text-transform: uppercase;
  font-family: 'Inter', sans-serif; font-weight: 400;
}
.lux-lang-btn--active { border-color: #C5A572; color: #C5A572; }
.lux-status-pill {
  font-size: 10px; letter-spacing: 0.15em;
  display: flex; align-items: center; gap: 6px;
  padding: 4px 12px; border-radius: 2px;
  border: 1px solid #D8D0C4;
  font-family: 'Inter', sans-serif;
  text-transform: uppercase;
}
.lux-status-pill--on { color: #4a7c59; border-color: #4a7c59; }
.lux-status-pill--off { color: #8a8578; }
.lux-status-dot { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }

.lux-live-desc {
  font-family: 'Cormorant Garamond', serif;
  font-size: 16px; line-height: 1.6; color: #44403a;
}
</style>
