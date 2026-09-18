<script setup lang="ts">
import { onMounted, onUnmounted, ref, nextTick, computed, watch } from 'vue'
// MessageBubble is defined as a local <script> block component below — accessible in <template>
import { api } from '@/api/client'

// ============ Types ============
type Attachment = { url: string; type: string; name?: string; size?: number }
type CardPayload = { kind: string; title?: string; sku?: string; price?: number; lead_time?: string; product_token?: string; order_no?: string; status?: string; total?: number; redirect_url?: string } | null
type ChatMsg = {
  id?: number
  conversation_id?: number
  message_type: string
  content: string
  attachments?: Attachment[]
  card_payload?: CardPayload
  sender_type: string
  sender_id: number
  sender_name?: string
  translation_zh?: string
  translation_en?: string
  translation_status?: string
  created_at: string
  local?: boolean
  client_msg_id?: string
}

// ============ State ============
const convs = ref<any[]>([])
const active = ref<any>(null)
const messages = ref<ChatMsg[]>([])
const input = ref('')
const ws = ref<WebSocket | null>(null)
const wsConnected = ref(false)
const isLoggedIn = computed(() => !!localStorage.getItem('user_token'))
const sidebarOpen = ref(true)
const emojiPickerOpen = ref(false)
const showUploadMenu = ref(false)
const showCardMenu = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const pendingAttachments = ref<Attachment[]>([])
const showTranslation = ref<Record<number, boolean>>({})
const scrollRef = ref<HTMLElement | null>(null)

// ============ Emoji Set (lightweight, no external dep) ============
const EMOJIS = [
  '😀','😂','🤣','😊','😍','🥰','😘','😎','🤔','😢','😭','😡','🥳','😴','🤗','😇',
  '👍','👎','👏','🙏','💪','✌️','🤝','👀','❤️','💔','💕','💖','💯','🔥','✨','🌟',
  '☕','🍵','🍃','🌿','🌸','🌹','🍜','🥟','🧧','🎋','🏮','🫖','🍵','🫙','🍶','🥢',
  '🎁','🎉','🥳','🎊','💐','🏆','📦','📮','💰','💷','💴','💵','🪙','📜','📄','📎'
]

// ============ 我是谁（从 localStorage token 类型推断） ============
// JWT payload 里有 subject_type: "user" | "staff"，用 token key 名直接推断即可
const meUserType: 'user' | 'staff' | '' = (() => {
  if (localStorage.getItem('user_token')) return 'user'
  if (localStorage.getItem('staff_token')) return 'staff'
  return ''
})()

// ============ WebSocket ============
function wsURL(token: string) {
  const host = window.location.host
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
  // 后端 /ws/im 注册在根路由，不在 /api/v1 分组里
  return `${proto}://${host}/ws/im?token=${token}`
}

function connectWS() {
  const token = localStorage.getItem('user_token') || localStorage.getItem('staff_token')
  if (!token) { console.warn('no token for ws'); return }
  ws.value?.close()
  ws.value = new WebSocket(wsURL(token))
  ws.value.onopen = () => {
    wsConnected.value = true
    console.log('WS connected')
    if (active.value) {
      ws.value?.send(JSON.stringify({ type: 'join_conversation', payload: { conversation_id: active.value.id } }))
    }
  }
  ws.value.onclose = () => { wsConnected.value = false; console.log('WS closed') }
  ws.value.onerror = () => { wsConnected.value = false }
  ws.value.onmessage = (ev) => {
    try {
      const env = JSON.parse(ev.data)
      if (env.type === 'chat_message') {
        handleIncomingMessage(env.payload)
      } else if (env.type === 'barrage') {
        // barrages handled in LiveRoom
      } else if (env.type === 'ack') {
        // remove local pending msg
        const idx = messages.value.findIndex(m => m.client_msg_id === env.payload.client_msg_id && m.local)
        if (idx >= 0) messages.value.splice(idx, 1)
      } else if (env.type === 'error') {
        console.warn('WS error:', env.payload)
      }
    } catch (e) { /* raw text? */ }
  }
}

function handleIncomingMessage(msg: ChatMsg) {
  // avoid dup
  if (msg.id) {
    const exists = messages.value.find(m => m.id === msg.id)
    if (exists) {
      // update translation
      if (msg.translation_status === 'translated' && (msg.translation_en || msg.translation_zh)) {
        Object.assign(exists, msg)
        showTranslation.value[msg.id!] = true
      }
      return
    }
  }
  messages.value.push(msg)
  scrollToBottom()
}

function scrollToBottom() { nextTick(() => { if (scrollRef.value) scrollRef.value.scrollTop = scrollRef.value.scrollHeight }) }

// ============ Conversations ============
async function loadConvs() {
  try { convs.value = (await api.get('/conversations') as any).items || [] } catch {}
}

async function selectConv(c: any) {
  active.value = c
  showCardMenu.value = false
  emojiPickerOpen.value = false
  if (!c) return
  try {
    const resp = await api.get(`/conversations/${c.id}/messages`) as any
    messages.value = resp.items || resp || []
  } catch { messages.value = [] }
  // join ws room
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({ type: 'join_conversation', payload: { conversation_id: c.id } }))
  } else {
    connectWS()
  }
  await nextTick()
  scrollToBottom()
}

// ============ Upload ============
async function uploadFile(file: File, type: 'image' | 'video' | 'file'): Promise<Attachment | null> {
  const fd = new FormData()
  fd.append('file', file)
  uploading.value = true
  uploadProgress.value = 0
  try {
    const resp: any = await api.post(`/upload?type=${type}`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e: any) => { uploadProgress.value = Math.round((e.loaded / e.total) * 100) }
    })
    const f = resp.file
    return { url: f.url, type: f.type, name: f.name, size: f.size }
  } catch (e: any) {
    console.warn('upload failed:', e)
    return null
  } finally { uploading.value = false; uploadProgress.value = 0 }
}

async function handleFiles(files: FileList | null, type: 'image' | 'video' | 'file') {
  if (!files?.length) return
  const atts: Attachment[] = []
  for (const f of Array.from(files)) {
    const a = await uploadFile(f, type)
    if (a) atts.push(a)
  }
  pendingAttachments.value.push(...atts)
  showUploadMenu.value = false
}

// Paste image from clipboard
async function onPaste(e: ClipboardEvent) {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of Array.from(items)) {
    if (item.kind === 'file' && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) { e.preventDefault(); const a = await uploadFile(file, 'image'); if (a) pendingAttachments.value.push(a) }
    }
  }
}

function removePendingAttachment(i: number) { pendingAttachments.value.splice(i, 1) }

// ============ Send ============
function makeClientID() { return 'c_' + Date.now() + '_' + Math.random().toString(36).slice(2, 7) }

function send() {
  if (!active.value) return
  const content = input.value.trim()
  const atts = pendingAttachments.value
  if (!content && atts.length === 0) return

  let msgType = 'text'
  let finalContent = content
  if (atts.length > 0 && atts.every(a => a.type === 'image')) msgType = 'image'
  else if (atts.length > 0 && atts.every(a => a.type === 'video')) msgType = 'video'
  else if (atts.length > 0) msgType = 'file'
  else if (content && [...content].every(c => isEmoji(c) || /\s/.test(c))) msgType = 'emoji'

  const cid = makeClientID()
  const payload = {
    conversation_id: active.value.id,
    content: finalContent,
    message_type: msgType,
    attachments: atts.length ? atts : undefined,
    client_msg_id: cid
  }

  // optimistic local msg
  messages.value.push({
    content: finalContent, message_type: msgType, attachments: atts.length ? atts : undefined,
    sender_type: 'me', sender_id: 0, created_at: new Date().toISOString(), local: true, client_msg_id: cid
  })
  scrollToBottom()

  ws.value?.send(JSON.stringify({ type: 'send_message', payload }))

  input.value = ''
  pendingAttachments.value = []
  emojiPickerOpen.value = false
}

// ============ Emoji ============
function insertEmoji(e: string) { input.value += e }
function isEmoji(c: string) { return /[\u{1F300}-\u{1FAFF}\u{2600}-\u{27BF}\u{FE00}-\u{FE0F}\u{200D}]/u.test(c) }

// ============ Translation Toggle ============
function toggleTranslation(id: number) { showTranslation.value[id] = !showTranslation.value[id] }
function displayTranslation(msg: ChatMsg): string | undefined {
  if (!msg) return undefined
  if ((msg.sender_type || '').toLowerCase() === 'user' && msg.sender_id !== 1) return msg.translation_zh
  return msg.translation_en
}
function showTranslate(msg: ChatMsg): boolean {
  return !!(msg.sender_type && msg.sender_id !== 1) && (msg.translation_zh || msg.translation_en) && msg.message_type === 'text'
}

// ============ Card Message (advisor side) ============
async function sendQuoteCard() {
  showCardMenu.value = false
  try {
    const products = (await api.get('/custom-products') as any).items || []
    const me = (await api.get('/custom-products') as any).items?.[0] || products?.[0]
    if (!me) return
    const payload = {
      conversation_id: active.value.id,
      message_type: 'quote_card',
      content: '',
      card_payload: {
        kind: 'quote_card',
        title: me.title || me.product_name || "Bespoke Pu\'er",
        sku: me.sku, price: me.unit_price, lead_time: me.lead_time,
        product_token: me.product_token, redirect_url: `/bespoke/${me.product_token || me.id}`
      }
    }
    ws.value?.send(JSON.stringify({ type: 'send_message', payload }))
  } catch {}
}

async function sendOrderCard() {
  showCardMenu.value = false
  try {
    const orders = (await api.get('/orders') as any).items || []
    const o = orders?.[0]
    if (!o) return
    const payload = {
      conversation_id: active.value.id,
      message_type: 'order_card',
      content: '',
      card_payload: {
        kind: 'order_card',
        order_no: o.order_no, status: o.state, total: o.total_amount,
        redirect_url: `/account/orders/${o.id}`
      }
    }
    ws.value?.send(JSON.stringify({ type: 'send_message', payload }))
  } catch {}
}

// ============ Lifecycle ============
onMounted(async () => {
  if (!isLoggedIn.value) return
  await loadConvs()
  if (!convs.value.length) {
    // auto-create conversation with default advisor
    // 自动创建与默认 advisor（staff_id=1）的会话
    try { convs.value = (await api.post('/conversations', { other_staff_id: 1 }) as any) || [] } catch {}
    await loadConvs()
  }
  if (convs.value.length) await selectConv(convs.value[0])
})

onUnmounted(() => { ws.value?.close() })
</script>


<!-- ============ Inline MessageBubble Sub-component ============ -->
<script lang="ts">
import { defineComponent, computed, h } from 'vue'

export const MessageBubble = defineComponent({
  name: 'MessageBubble',
  props: {
    msg: { type: Object as () => any, required: true },
    showTranslation: { type: Boolean, default: false }
  },
  emits: ['toggle-translation'],
  setup(props, { emit }) {
    // isMine 判断: 后端 sender_type 是 "user" / "staff" / "system"
    // 用 module-level meUserType 对比（从 localStorage token 类型推断）
    const isMine = computed(() => {
      const t = (props.msg.sender_type || '').toLowerCase()
      if (meUserType && t === meUserType) return true
      // 本地 pending 消息 sender_type 硬编码为 "me"
      if (t === 'me') return true
      return false
    })
    const isCard = computed(() => ['quote_card', 'order_card'].includes(props.msg.message_type))
    const atts = computed(() => props.msg.attachments || [])

    return () => {
      const m = props.msg
      const bubbleClass = isMine.value
        ? 'bg-tea-700 text-white rounded-2xl rounded-tr-sm'
        : 'bg-white border border-tea-200 rounded-2xl rounded-tl-sm'
      const wrapperAlign = isMine.value ? 'items-end' : 'items-start'

      // Build children
      const children: any[] = []

      // Card payload (quote/order)
      if (isCard.value && m.card_payload) {
        const cp = m.card_payload
        if (m.message_type === 'quote_card') {
          children.push(
            h('div', { class: `max-w-xs shadow-lg overflow-hidden rounded-xl ${isMine.value ? 'bg-tea-600/95 text-white' : 'bg-gradient-to-br from-tea-50 to-white text-tea-900 border border-tea-200'}` }, [
              h('div', { class: 'p-4' }, [
                h('div', { class: 'text-xs opacity-70 uppercase tracking-wider mb-1' }, '📋 Quote'),
                h('div', { class: 'font-serif text-lg font-semibold' }, cp.title || "Bespoke Pu-er Tea"),
                cp.sku ? h('div', { class: 'text-xs opacity-70 mt-1' }, `SKU: ${cp.sku}`) : null,
                h('div', { class: 'flex items-baseline gap-2 mt-2' }, [
                  h('span', { class: 'text-2xl font-bold' }, `£${cp.price ?? '—'}`),
                  cp.lead_time ? h('span', { class: 'text-xs opacity-60' }, `${cp.lead_time || '10-15 days'} lead time`) : null
                ]),
                h('button', {
                  class: `mt-3 w-full py-2 rounded-lg text-sm font-medium transition ${isMine.value ? 'bg-white text-tea-700 hover:bg-tea-50' : 'bg-tea-700 text-white hover:bg-tea-800'}`,
                  onClick: () => window.open(cp.redirect_url || '#', '_blank')
                }, 'View Details →')
              ])
            ])
          )
        } else {
          // order_card
          const statusColor: any = { ordering:'bg-blue-100 text-blue-700', paid:'bg-green-100 text-green-700', producing:'bg-amber-100 text-amber-700', shipped:'bg-indigo-100 text-indigo-700', completed:'bg-green-100 text-green-700', cancelled:'bg-red-100 text-red-700' }
          children.push(
            h('div', { class: `max-w-xs shadow-lg overflow-hidden rounded-xl ${isMine.value ? 'bg-tea-600/95 text-white' : 'bg-gradient-to-br from-white to-slate-50 text-tea-900 border border-tea-200'}` }, [
              h('div', { class: 'p-4' }, [
                h('div', { class: 'flex items-center justify-between mb-2' }, [
                  h('div', { class: 'text-xs opacity-70 uppercase tracking-wider' }, '🧾 Order'),
                  h('span', { class: `text-xs px-2 py-0.5 rounded-full ${statusColor[m.status] || 'bg-gray-100 text-gray-700'}` }, m.status || '—')
                ]),
                h('div', { class: 'font-serif text-lg font-semibold' }, cp.order_no || 'Order'),
                h('div', { class: 'mt-3 text-2xl font-bold' }, `£${cp.total ?? '—'}`),
                h('button', {
                  class: `mt-3 w-full py-2 rounded-lg text-sm font-medium transition ${isMine.value ? 'bg-white text-tea-700 hover:bg-tea-50' : 'bg-tea-700 text-white hover:bg-tea-800'}`,
                  onClick: () => window.open(cp.redirect_url || '#', '_blank')
                }, 'View Order →')
              ])
            ])
          )
        }
      } else {
        // Regular text/emoji/image/video/file bubble
        // Attachments
        if (atts.value.length) {
          atts.value.forEach((a: any) => {
            if (a.type === 'image') {
              children.push(h('img', { src: a.url, class: 'max-w-xs rounded-lg block shadow-sm', alt: a.name || '' }))
            } else if (a.type === 'video') {
              children.push(h('video', { src: a.url, controls: true, class: 'max-w-xs rounded-lg block shadow-sm' }))
            } else {
              children.push(h('a', { href: a.url, download: a.name, class: `max-w-xs p-3 rounded-lg flex items-center gap-2 transition ${isMine.value ? 'bg-tea-800/60 text-white hover:bg-tea-800' : 'bg-tea-50 text-tea-800 hover:bg-tea-100'}` }, [
                h('span', { class: 'text-xl' }, '📎'),
                h('div', null, [
                  h('div', { class: 'text-sm font-medium truncate max-w-[180px]' }, a.name || 'file'),
                  a.size ? h('div', { class: 'text-xs opacity-60' }, `${Math.round(a.size/1024)} KB`) : null
                ])
              ]))
            }
          })
        }

        // Content text
        if (m.content && m.content.trim()) {
          // emoji-only: big
          const bigEmoji = m.message_type === 'emoji'
          children.push(h('div', {
            class: `${bigEmoji ? 'text-4xl leading-none' : 'text-sm leading-relaxed'} whitespace-pre-wrap break-words`,
            innerHTML: escapeHtml(m.content)
          }))
        }

        // Translation badge
        const hasTrans = m.translation_status === 'translated' && (m.translation_en || m.translation_zh)
        if (hasTrans) {
          // 后端 asyncTranslate 只往"相反语言"字段写：原文中文→translation_en，原文英文→translation_zh
          // 不管是谁发的，显示有内容的那个（让用户看到与原文不同的译文）
          const transText = m.translation_en || m.translation_zh
          if (transText) {
            children.push(h('div', {
              class: 'mt-1 text-xs opacity-70 cursor-pointer underline',
              onClick: () => emit('toggle-translation'),
              title: 'Click to toggle original'
            }, props.showTranslation ? transText : '[点击显示翻译 / Click for translation]'))
          }
        }

        // Timestamp
        children.push(h('div', { class: `text-[10px] mt-1 opacity-50 ${isMine.value ? 'text-right' : ''}` }, new Date(m.created_at).toLocaleTimeString()))
      }

      return h('div', { class: `flex flex-col gap-1 ${wrapperAlign}` }, [
        // sender label
        !isMine.value ? h('div', { class: 'text-[10px] text-tea-500 px-1' }, m.sender_type === 'staff' ? 'Tea Advisor' : 'You') : null,
        h('div', { class: isCard.value ? '' : bubbleClass, style: isCard.value ? '' : 'max-width: 75%' }, children)
      ])
    }
  }
})

function escapeHtml(s: string) {
  const div = document.createElement('div')
  div.textContent = s
  return div.innerHTML
}
</script>

<template>
  <div class="pt-16 min-h-screen bg-gradient-to-br from-tea-50 via-white to-tea-50">
    <div class="max-w-6xl mx-auto px-4 md:px-6 py-6">
      <div class="mb-4 flex items-center gap-3">
        <button @click="sidebarOpen = !sidebarOpen" class="md:hidden p-2 rounded-lg hover:bg-tea-100">☰</button>
        <div>
          <h1 class="font-serif text-2xl md:text-3xl text-tea-900">💬 Chat with Your Advisor</h1>
          <p class="text-sm text-tea-600">Real-time messaging — text, images, files, video · auto-translated between English & Chinese</p>
        </div>
        <span class="ml-auto flex items-center gap-1 text-xs px-3 py-1 rounded-full" :class="wsConnected ? 'bg-green-100 text-green-700' : 'bg-tea-100 text-tea-500'">
          <span class="w-1.5 h-1.5 rounded-full" :class="wsConnected ? 'bg-green-500 animate-pulse' : 'bg-tea-400'"></span>
          {{ wsConnected ? 'Online' : 'Offline' }}
        </span>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-tea-100 overflow-hidden">
        <div class="flex h-[calc(100vh-140px)] min-h-[560px]">
          <!-- Sidebar: Conversations -->
          <aside v-if="sidebarOpen" class="w-60 md:w-64 border-r bg-tea-50/60 overflow-auto flex-shrink-0">
            <div class="p-3 border-b bg-white/80 text-xs font-semibold text-tea-700 uppercase tracking-wider">Conversations</div>
            <div v-for="c in convs" :key="c.id" @click="selectConv(c)"
                 :class="['p-3 border-b cursor-pointer transition text-sm', active?.id === c.id ? 'bg-tea-700 text-white' : 'hover:bg-tea-100 text-tea-800']">
              <div class="font-medium">Conversation #{{ c.id }}</div>
              <div class="text-xs opacity-70 mt-0.5">2 participants</div>
            </div>
          </aside>

          <!-- Chat Area -->
          <section class="flex-1 flex flex-col min-w-0">
            <div class="flex-1 overflow-auto p-4 md:p-6 bg-slate-50" ref="scrollRef">
              <div v-if="!isLoggedIn" class="text-center mt-16">
                <div class="text-5xl mb-4">🫖</div>
                <h2 class="font-serif text-2xl text-tea-900 mb-2">Speak with your Tea Advisor</h2>
                <p class="text-tea-600 text-sm mb-6 max-w-md mx-auto">
                  Sign in to chat with a Pu'er tea expert. Bespoke blending, garden stories,
                  and instant answers — in English or Chinese.
                </p>
                <div class="flex items-center justify-center gap-3">
                  <router-link to="/magic-link?redirect=/chat"
                               class="px-6 py-2.5 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition">
                    Sign in to Chat →
                  </router-link>
                  <router-link to="/about"
                               class="px-5 py-2.5 text-tea-700 hover:bg-tea-50 rounded-xl transition text-sm">
                    Learn more
                  </router-link>
                </div>
              </div>

              <div v-else-if="!active" class="text-center text-tea-500 mt-20">
                <div class="text-5xl mb-3">💬</div>
                Select a conversation to start chatting →
              </div>

              <div v-else class="space-y-3 max-w-3xl mx-auto">
                <div v-for="(m, i) in messages" :key="m.id || m.client_msg_id || i"
                     :class="['flex', m.sender_type === 'me' ? 'justify-end' : 'justify-start']">
                  <MessageBubble :msg="m" :show-translation="showTranslation[m.id!]" @toggle-translation="toggleTranslation(m.id!)" />
                </div>
              </div>
            </div>

            <!-- Pending attachments preview -->
            <div v-if="pendingAttachments.length" class="border-t bg-tea-50/50 px-4 py-3 flex gap-3 flex-wrap">
              <div v-for="(a, i) in pendingAttachments" :key="i" class="relative group">
                <img v-if="a.type === 'image'" :src="a.url" class="h-16 w-16 object-cover rounded-lg border" />
                <div v-else class="h-16 w-16 rounded-lg bg-white border flex flex-col items-center justify-center text-xs text-tea-600">
                  <span>{{ a.type === 'video' ? '🎬' : '📎' }}</span>
                  <span class="truncate px-1 w-full text-center">{{ a.name }}</span>
                </div>
                <button @click="removePendingAttachment(i)" class="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-red-500 text-white text-xs opacity-0 group-hover:opacity-100 transition">×</button>
              </div>
              <span class="text-xs text-tea-600 self-center">{{ pendingAttachments.length }} file(s) attached</span>
            </div>

            <!-- Input Area -->
            <div class="border-t p-3 md:p-4 bg-white">
              <div class="flex items-end gap-2">
                <!-- Toolbar -->
                <div class="flex items-center gap-1 text-tea-500">
                  <button @click="emojiPickerOpen = !emojiPickerOpen; showUploadMenu = false; showCardMenu = false" title="Emoji" class="p-2 hover:bg-tea-100 rounded-lg transition relative">😊</button>
                  <button @click="showUploadMenu = !showUploadMenu; emojiPickerOpen = false; showCardMenu = false" title="Upload" class="p-2 hover:bg-tea-100 rounded-lg transition relative">📎</button>
                  <button @click="showCardMenu = !showCardMenu; showUploadMenu = false; emojiPickerOpen = false" title="Send Quote/Order" class="p-2 hover:bg-tea-100 rounded-lg transition relative">💳</button>
                </div>

                <!-- Text input -->
                <div class="flex-1 relative">
                  <textarea v-model="input" @keyup.enter.exact.prevent="send" @paste="onPaste" rows="1"
                            placeholder="Type a message... (paste image to attach)"
                            class="w-full px-4 py-2.5 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none resize-none text-sm bg-tea-50/30 leading-relaxed"></textarea>
                </div>

                <button @click="send" :disabled="!input.trim() && !pendingAttachments.length"
                        class="px-6 py-2.5 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition disabled:opacity-40 disabled:cursor-not-allowed whitespace-nowrap text-sm">
                  Send
                </button>
              </div>

              <!-- Emoji picker popover -->
              <div v-if="emojiPickerOpen" class="absolute bg-white border rounded-xl shadow-lg p-3 z-10 mt-1 max-w-[320px]">
                <div class="grid grid-cols-8 gap-1 text-xl">
                  <button v-for="e in EMOJIS" :key="e" @click="insertEmoji(e); emojiPickerOpen = false" class="hover:bg-tea-50 rounded p-1 transition text-center">{{ e }}</button>
                </div>
              </div>

              <!-- Upload menu -->
              <div v-if="showUploadMenu" class="absolute bg-white border rounded-xl shadow-lg p-2 z-10 mt-1 flex flex-col gap-1 text-sm w-36">
                <label class="cursor-pointer hover:bg-tea-50 px-3 py-2 rounded-lg">🖼️ Upload Image<input type="file" accept="image/*" class="hidden" @change="(e:any) => handleFiles(e.target.files, 'image')" /></label>
                <label class="cursor-pointer hover:bg-tea-50 px-3 py-2 rounded-lg">🎬 Upload Video<input type="file" accept="video/*" class="hidden" @change="(e:any) => handleFiles(e.target.files, 'video')" /></label>
                <label class="cursor-pointer hover:bg-tea-50 px-3 py-2 rounded-lg">📎 Upload File<input type="file" class="hidden" @change="(e:any) => handleFiles(e.target.files, 'file')" /></label>
              </div>

              <!-- Card menu -->
              <div v-if="showCardMenu" class="absolute bg-white border rounded-xl shadow-lg p-2 z-10 mt-1 flex flex-col gap-1 text-sm w-40">
                <button @click="sendQuoteCard" class="hover:bg-tea-50 px-3 py-2 rounded-lg text-left">📋 Send Quote Card</button>
                <button @click="sendOrderCard" class="hover:bg-tea-50 px-3 py-2 rounded-lg text-left">🧾 Send Order Card</button>
              </div>

              <!-- Upload progress -->
              <div v-if="uploading" class="mt-2 h-1.5 bg-tea-100 rounded overflow-hidden">
                <div class="h-full bg-tea-600 transition-all" :style="{ width: uploadProgress + '%' }"></div>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>
