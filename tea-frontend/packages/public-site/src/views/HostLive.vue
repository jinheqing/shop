<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// ========== 从 URL 读取开播参数 ==========
const config = {
  livekitUrl: route.query.livekit_url as string || import.meta.env.VITE_LIVEKIT_URL || '',
  hostToken:  route.query.host_token as string  || localStorage.getItem('host_token') || '',
  roomId:     route.query.room_id as string     || localStorage.getItem('room_id') || '',
  apiBase:    (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '') + '/api/v1',
}

// ========== LiveKit ==========
const lkClient = ref<any>(null)
const lkConnected = ref(false)
const lkPublishing = ref(false)
const localVideoEl = ref<HTMLVideoElement | null>(null)
const guestVideos = ref<Array<{ participant: any; el: HTMLVideoElement }>>([])

// ========== IM WebSocket ==========
const imWS = ref<WebSocket | null>(null)
const imConnected = ref(false)

// ========== 连麦状态 ==========
const incomingInvite = ref<{ user_id: number; nickname: string; link_session_id: string } | null>(null)
const activeLinkSession = ref<string>('')  // 当前活跃的连麦 session_id

// ========== ASR → 字幕广播 ==========
const asrWS = ref<WebSocket | null>(null)
const asrConnected = ref(false)
const asrEnabled = ref(true)  // 主播端默认开启字幕
const audioStream = ref<MediaStream | null>(null)
const audioProcessor = ref<ScriptProcessorNode | null>(null)
const audioSource = ref<MediaStreamAudioSourceNode | null>(null)
const latestSubtitle = ref<{ text: string; translation: string; type: 'partial' | 'final' } | null>(null)

// ========== 观众弹幕（聊天） ==========
const chatMessages = ref<Array<{ subtype: string; content: string; nickname?: string; ts: number }>>([])

const canStart = computed(() => !!config.livekitUrl && !!config.hostToken && !!config.roomId)

// ============ 生命周期 ============
onMounted(async () => {
  // 保存一份到 localStorage 便于后续恢复
  if (config.hostToken) localStorage.setItem('host_token', config.hostToken)
  if (config.roomId) localStorage.setItem('room_id', config.roomId)

  // 1) 建立 IM WS（staff_token 或 user_token 皆可 — hub 会自动按 userType 分 clientKey）
  const token = localStorage.getItem('staff_token') || localStorage.getItem('user_token') || ''
  if (token) openIM(token)

  // 2) 加入 LiveKit 并自动推流（主播端不需要 joinLiveKit 后再手动推 — 直接 publish）
  if (canStart.value) await startBroadcast()
})

onUnmounted(() => {
  stopBroadcast()
  closeIM()
  closeASR()
})

// ============ IM WebSocket ============
function openIM(jwt: string) {
  const host = window.location.host
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
  imWS.value = new WebSocket(`${proto}://${host}/ws/im?token=${jwt}`)
  imWS.value.onopen = () => {
    imConnected.value = true
    imWS.value?.send(JSON.stringify({ type: 'join_room', payload: { room_id: config.roomId } }))
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
  const env = { type: 'barrage', payload: { ...payload, room_id: config.roomId, ts: Date.now() } }
  imWS.value.send(JSON.stringify(env))
}

function handleBarrage(p: any) {
  const st = p.subtype || 'chat'
  if (st === 'chat') {
    chatMessages.value.push({ subtype: st, content: p.content, nickname: p.nickname, ts: p.ts || Date.now() })
    if (chatMessages.value.length > 100) chatMessages.value = chatMessages.value.slice(-100)
  } else if (st === 'link_invite') {
    // 观众请求连麦 → 主播弹卡片
    incomingInvite.value = {
      user_id: p.user_id,
      nickname: p.nickname || `Viewer #${p.user_id}`,
      link_session_id: p.link_session_id,
    }
  } else if (st === 'link_accept') {
    // 观众接受了主播的邀请 → 主播等对方 publish 进来即可（LiveKit TrackSubscribed 会触发）
  } else if (st === 'link_reject') {
    // 观众拒绝了邀请
    incomingInvite.value = null
  } else if (st === 'link_end') {
    activeLinkSession.value = ''
    guestVideos.value = []
  } else if (st === 'subtitle') {
    // 观众端字幕（主播端自己的字幕走本地 ASR，不用这里）
  }
}

// ============ LiveKit 开播 ============
async function startBroadcast() {
  try {
    const { Room, RoomEvent, VideoPresets } = await import('livekit-client')
    const room = new Room({ autoSubscribe: false, dynacast: true, videoCaptureDefaults: { resolution: VideoPresets.h720 } })
    lkClient.value = room

    // 远端 guest 视频（连麦观众 publish 进来后）
    room.on(RoomEvent.TrackSubscribed, (track: any, _pub: any, participant: any) => {
      if (track.kind === 'video') {
        const el = track.attach()
        guestVideos.value.push({ participant, el })
      } else if (track.kind === 'audio') {
        track.attach()  // 默认播放连麦观众的音频（主播能听到对方说话）
      }
    })
    room.on(RoomEvent.TrackUnsubscribed, (track: any) => {
      guestVideos.value = guestVideos.value.filter(g => g.participant !== track.source?.participant)
      try { track.detach() } catch {}
    })

    await room.connect(config.livekitUrl, config.hostToken)
    lkConnected.value = true

    // 立即发布自己的摄像头 + 麦克风
    await room.localParticipant.setCameraEnabled(true)
    await room.localParticipant.setMicrophoneEnabled(true)
    lkPublishing.value = true

    // 把本地视频挂到页面上
    setTimeout(() => {
      const localTrack = room.localParticipant.videoTracks.values().next().value
      if (localTrack && localVideoEl.value) {
        localVideoEl.value.srcObject = new MediaStream([localTrack.mediaStreamTrack])
      }
    }, 500)

    // 启动主播端 ASR 字幕广播
    if (asrEnabled.value) await openASRForHost()
  } catch (e: any) {
    alert('LiveKit connect failed: ' + (e?.message || String(e)))
  }
}

async function stopBroadcast() {
  try { await lkClient.value?.disconnect() } catch {}
  lkClient.value = null
  lkConnected.value = false
  lkPublishing.value = false
  closeASR()
  guestVideos.value = []
}

async function toggleMic() {
  if (!lkClient.value) return
  await lkClient.value.localParticipant.setMicrophoneEnabled(!lkClient.value.localParticipant.isMicrophoneEnabled)
}
async function toggleCam() {
  if (!lkClient.value) return
  await lkClient.value.localParticipant.setCameraEnabled(!lkClient.value.localParticipant.isCameraEnabled)
}

// ============ 连麦：邀请 / 接受 / 拒绝 / 结束 ============
function acceptLinkMic() {
  if (!incomingInvite.value) return
  // 主播接受 → 广播 link_accept + 主播自己 publish 保持不变
  sendBarrage({ subtype: 'link_accept', content: 'accept', target_user_id: incomingInvite.value.user_id, link_session_id: incomingInvite.value.link_session_id })
  activeLinkSession.value = incomingInvite.value.link_session_id
  incomingInvite.value = null
}
function rejectLinkMic() {
  if (!incomingInvite.value) return
  sendBarrage({ subtype: 'link_reject', content: 'reject', target_user_id: incomingInvite.value.user_id, link_session_id: incomingInvite.value.link_session_id })
  incomingInvite.value = null
}
function endLinkMic() {
  if (!activeLinkSession.value) return
  sendBarrage({ subtype: 'link_end', content: 'end', link_session_id: activeLinkSession.value })
  activeLinkSession.value = ''
  guestVideos.value = []
}

// ============ 主播端 ASR → 字幕广播 ============
function asrWSURL() {
  const host = window.location.hostname || 'localhost'
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
  return `${proto}://${host}:8090/translate/asr-stream`
}

async function openASRForHost() {
  if (asrWS.value?.readyState === WebSocket.OPEN) return
  try {
    // 优先拿 LiveKit 本地音频轨，不行就重新 getUserMedia
    let stream: MediaStream
    const lp = lkClient.value?.localParticipant
    if (lp?.audioTracks?.size) {
      stream = new MediaStream()
      for (const t of lp.audioTracks.values()) stream.addTrack(t.mediaStreamTrack)
    } else {
      stream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false })
    }
    audioStream.value = stream

    asrWS.value = new WebSocket(asrWSURL())
    asrWS.value.binaryType = 'arraybuffer'

    asrWS.value.onopen = () => {
      asrConnected.value = true
      asrWS.value?.send(JSON.stringify({ target_lang: 'auto' }))
      startAudioCapture(stream)
    }
    asrWS.value.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data)
        latestSubtitle.value = { text: msg.text, translation: msg.translation || '', type: msg.type }
        // 以 barrage subtype=subtitle 广播给所有观众
        sendBarrage({
          subtype: 'subtitle',
          content: msg.text,
          speaker: 'host',
          text: msg.text,
          translation: msg.translation || '',
          lang: 'auto',
          transcript_type: msg.type,
          nickname: 'Host',
        })
      } catch {}
    }
    asrWS.value.onclose = () => { asrConnected.value = false; stopAudioCapture() }
  } catch (e: any) {
    console.warn('ASR init failed:', e)
  }
}

function closeASR() { stopAudioCapture(); try { asrWS.value?.close() } catch {}; asrWS.value = null; asrConnected.value = false }

function startAudioCapture(stream: MediaStream) {
  if (!asrWS.value) return
  const ctx = new (window.AudioContext || (window as any).webkitAudioContext)({ sampleRate: 16000 })
  const source = ctx.createMediaStreamSource(stream)
  const processor = ctx.createScriptProcessor(4096, 1, 1)
  processor.onaudioprocess = (ev) => {
    if (asrWS.value?.readyState !== WebSocket.OPEN) return
    const input = ev.inputBuffer.getChannelData(0)
    const int16 = new Int16Array(input.length)
    for (let i = 0; i < input.length; i++) {
      const s = Math.max(-1, Math.min(1, input[i]))
      int16[i] = s < 0 ? s * 0x8000 : s * 0x7FFF
    }
    asrWS.value.send(int16.buffer)
  }
  source.connect(processor)
  processor.connect(ctx.destination)
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
}
</script>

<template>
  <!-- ========= 手机端 H5 布局：竖屏 16:9 主视频 + 下方控制条 ========= -->
  <div class="host-root">
    <!-- 开播前配置缺失提示 -->
    <div v-if="!canStart" class="host-setup">
      <h1 class="host-setup-title">🎬 Host Setup</h1>
      <p class="host-setup-desc">Open this page from the admin dashboard "Go Live" dialog (tab: 📱 Mobile).</p>
      <div class="host-config-check">
        <div :class="['host-check', config.livekitUrl ? 'ok' : 'bad']">LiveKit URL: {{ config.livekitUrl || 'MISSING' }}</div>
        <div :class="['host-check', config.hostToken ? 'ok' : 'bad']">Host Token: {{ config.hostToken ? '✓ present' : 'MISSING' }}</div>
        <div :class="['host-check', config.roomId ? 'ok' : 'bad']">Room ID: {{ config.roomId || 'MISSING' }}</div>
      </div>
    </div>

    <!-- 开播后主界面 -->
    <div v-else class="host-live">
      <!-- 主视频 (本地摄像头镜像) -->
      <div class="host-video-stage">
        <video ref="localVideoEl" autoplay muted playsinline class="host-video-local" />

        <!-- 连麦观众小窗 (PIP) -->
        <template v-for="(g, i) in guestVideos" :key="i">
          <div class="host-video-guest" :style="{ top: (16 + i * 96) + 'px' }">
            <video ref="g.el" autoplay playsinline muted />
            <span class="host-video-guest-label">Guest</span>
          </div>
        </template>

        <!-- 状态角标 -->
        <div class="host-status-bar">
          <span class="host-live-dot"></span>
          <span class="host-live-text">LIVE</span>
          <span class="host-room">{{ config.roomId }}</span>
        </div>

        <!-- 主播端字幕（反白，像直播监视器） -->
        <div v-if="latestSubtitle" class="host-subtitle">
          <span :class="['host-sub-text', latestSubtitle.type === 'partial' ? 'partial' : '']">{{ latestSubtitle.text }}</span>
          <span v-if="latestSubtitle.translation" class="host-sub-trans">— {{ latestSubtitle.translation }}</span>
        </div>
      </div>

      <!-- 控制面板 -->
      <div class="host-controls">
        <div class="host-btn-row">
          <button class="host-btn" :class="lkClient?.localParticipant?.isMicrophoneEnabled ? 'on' : 'off'" @click="toggleMic">Mic</button>
          <button class="host-btn" :class="lkClient?.localParticipant?.isCameraEnabled ? 'on' : 'off'" @click="toggleCam">Cam</button>
          <button class="host-btn" :class="asrConnected ? 'on' : 'off'" @click="asrConnected ? closeASR() : openASRForHost()">📝 CC</button>
          <button class="host-btn danger" @click="endLinkMic" :disabled="!activeLinkSession">🔗 End</button>
          <button class="host-btn danger-big" @click="stopBroadcast">⏹ End Stream</button>
        </div>
        <div class="host-connection">
          <span :class="['host-con-dot', lkConnected ? 'good' : 'bad']"></span>LiveKit
          <span :class="['host-con-dot', imConnected ? 'good' : 'bad']"></span>IM
          <span :class="['host-con-dot', asrConnected ? 'good' : 'bad']"></span>ASR
        </div>
      </div>

      <!-- 连麦邀请弹窗 -->
      <div v-if="incomingInvite" class="host-invite">
        <div class="host-invite-card">
          <p class="host-invite-text">{{ incomingInvite.nickname }} requests to join the broadcast. Accept?</p>
          <div class="host-invite-actions">
            <button class="host-btn success" @click="acceptLinkMic">Accept</button>
            <button class="host-btn" @click="rejectLinkMic">Reject</button>
          </div>
        </div>
      </div>

      <!-- 聊天弹幕 -->
      <div class="host-chat">
        <div v-for="(m, i) in chatMessages.slice(-6)" :key="i" class="host-chat-item">
          <span v-if="m.nickname" class="host-chat-name">{{ m.nickname }}:</span>
          <span>{{ m.content }}</span>
        </div>
        <div v-if="!chatMessages.length" class="host-chat-empty">No chat yet</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
* { box-sizing: border-box; }
.host-root { min-height: 100vh; background: #0B0A09; color: #F8F5EF; font-family: -apple-system, sans-serif; }

/* ========== 配置缺失页 ========== */
.host-setup { padding: 48px 24px; text-align: center; }
.host-setup-title { font-size: 24px; margin: 0 0 12px; }
.host-setup-desc { color: #8A8578; font-size: 14px; margin-bottom: 32px; }
.host-config-check { text-align: left; max-width: 400px; margin: 0 auto; }
.host-check { padding: 10px 14px; border-radius: 4px; margin-bottom: 8px; font-size: 13px; font-family: monospace; }
.host-check.ok { background: #2d5a3d; }
.host-check.bad { background: #5a2d2d; }

/* ========== 直播主界面 ========== */
.host-live { display: flex; flex-direction: column; height: 100vh; }

.host-video-stage {
  position: relative;
  flex: 0 0 auto;
  width: 100%;
  aspect-ratio: 9/16;
  background: #000;
  overflow: hidden;
}

.host-video-local {
  width: 100%; height: 100%;
  object-fit: cover;
  transform: scaleX(-1); /* 镜像 */
}

.host-video-guest {
  position: absolute;
  right: 12px;
  width: 120px; height: 160px;
  background: #2A2520;
  border: 2px solid #C5A572;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0,0,0,0.6);
}
.host-video-guest video { width: 100%; height: 100%; object-fit: cover; }
.host-video-guest-label {
  position: absolute; bottom: 4px; left: 6px;
  font-size: 10px; background: rgba(0,0,0,0.7); padding: 2px 6px; border-radius: 2px;
}

/* 状态 */
.host-status-bar {
  position: absolute; top: 12px; left: 12px;
  display: flex; align-items: center; gap: 8px;
  background: rgba(0,0,0,0.6); padding: 6px 12px; border-radius: 20px;
}
.host-live-dot { width: 8px; height: 8px; border-radius: 50%; background: #e11d48; animation: pulse 1.5s infinite; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }
.host-live-text { font-size: 11px; font-weight: 700; letter-spacing: 0.2em; color: #e11d48; }
.host-room { font-size: 11px; color: #8a8578; }

/* 字幕 */
.host-subtitle {
  position: absolute; bottom: 12px; left: 12px; right: 12px;
  background: rgba(0,0,0,0.75);
  padding: 10px 16px; border-radius: 6px;
  text-align: center;
}
.host-sub-text { font-size: 16px; color: #F8F5EF; }
.host-sub-text.partial { font-style: italic; opacity: 0.7; }
.host-sub-trans { display: block; font-size: 13px; color: #C5A572; margin-top: 2px; }

/* ========== 控制面板 ========== */
.host-controls {
  flex: 1; background: #161513; padding: 20px 16px;
  display: flex; flex-direction: column; gap: 14px;
}

.host-btn-row { display: flex; gap: 10px; flex-wrap: wrap; }
.host-btn {
  flex: 1 1 auto; min-width: 70px;
  padding: 14px 10px;
  border-radius: 8px;
  border: 1px solid #3a3530;
  background: #1e1b18;
  color: #D8D0C4;
  font-size: 13px; font-weight: 500;
  cursor: pointer;
  transition: all 150ms;
}
.host-btn.on { background: #4a7c59; border-color: #4a7c59; }
.host-btn.off { background: #2A2520; }
.host-btn.danger { background: #5a2d2d; border-color: #7a3d3d; }
.host-btn.danger-big { background: #7a2d2d; border-color: #e11d48; color: #fff; font-weight: 700; flex: 2; }
.host-btn.success { background: #2d5a3d; border-color: #4a7c59; }
.host-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.host-connection {
  display: flex; gap: 16px; justify-content: center;
  font-size: 11px; color: #8a8578;
  padding-top: 8px;
}
.host-con-dot {
  display: inline-block; width: 8px; height: 8px; border-radius: 50%;
  margin-right: 6px;
}
.host-con-dot.good { background: #4a7c59; }
.host-con-dot.bad { background: #e11d48; }

/* ========== 邀请弹窗 ========== */
.host-invite {
  position: fixed; top: 0; left: 0; right: 0;
  padding: 24px;
  background: rgba(11,10,9,0.95);
  z-index: 100;
}
.host-invite-card {
  background: #2A2520; border-radius: 12px;
  padding: 24px;
}
.host-invite-text { font-size: 16px; margin-bottom: 16px; line-height: 1.5; }
.host-invite-actions { display: flex; gap: 12px; }
.host-invite-actions .host-btn { flex: 1; }

/* ========== 聊天弹幕 ========== */
.host-chat {
  max-height: 120px; overflow: hidden;
  background: rgba(0,0,0,0.3);
  padding: 8px 12px;
  border-radius: 4px;
}
.host-chat-item { font-size: 12px; padding: 3px 0; color: #D8D0C4; }
.host-chat-name { color: #C5A572; margin-right: 6px; }
.host-chat-empty { font-size: 12px; color: #6b6459; text-align: center; padding: 12px; }
</style>
