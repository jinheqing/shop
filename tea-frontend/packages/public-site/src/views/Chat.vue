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
const isLoggedIn = ref(!!localStorage.getItem('user_token'))
// 监听登录状态变化（localStorage 非响应式，需手动监听）
window.addEventListener('storage', (e) => {
  if (e.key === 'user_token') isLoggedIn.value = !!e.newValue
})
const sidebarOpen = ref(true)
const mobileSidebar = ref(false)
const showUploadMenu = ref(false)
const showCardMenu = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const pendingAttachments = ref<Attachment[]>([])
const showTranslation = ref<Record<number, boolean>>({})
const scrollRef = ref<HTMLElement | null>(null)

// ============ Emoji Set (lightweight, no external dep) ============

// ============ 我是谁（从 localStorage token 类型推断） ============
// ⚠️ meUserType 定义在下方普通 <script> 块（module level），供 MessageBubble 使用
// <script setup> 与普通 <script> 作用域不互通，勿在此重复声明

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
  // 后端 /conversations 返回 {code:0, data:[...]}（无 items 包装）
  try { convs.value = (await api.get('/conversations') as any).data || [] } catch {}
}

async function selectConv(c: any) {
  active.value = c
  showCardMenu.value = false
  if (!c) return
  try {
    // 后端消息列表返回 {code:0, data:{messages:[...], count:n}}
    const resp = await api.get(`/conversations/${c.id}/messages`) as any
    messages.value = resp?.data?.messages || []
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
}

// ============ Emoji ============

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
  isLoggedIn.value = !!localStorage.getItem('user_token')
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

<script lang="ts">
import { defineComponent, computed, h } from 'vue'

// 我是谁（module level）—— 从 localStorage token 类型推断，供 MessageBubble 判断消息归属
// JWT payload 里有 subject_type: "user" | "staff"，用 token key 名直接推断
const meUserType: 'user' | 'staff' | '' = (() => {
  if (typeof localStorage !== 'undefined') {
    if (localStorage.getItem('user_token')) return 'user'
    if (localStorage.getItem('staff_token')) return 'staff'
  }
  return ''
})()

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
          children.push(h('div', {
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
  <div class="chat-page min-h-screen bg-gradient-to-b from-tea-50 to-white">

    <!-- 未登录引导页 -->
    <div v-if="!isLoggedIn" class="pt-24 px-4 text-center">
      <div class="max-w-md mx-auto">
        <div class="text-6xl mb-5">🫖</div>
        <h1 class="font-serif text-3xl text-tea-900 mb-3">Speak with your Tea Advisor</h1>
        <p class="text-tea-600 text-sm leading-relaxed mb-8">
          Sign in to chat with a Pu'er tea expert. Bespoke blending, garden stories,
          and instant answers — in English or Chinese.
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
          <router-link to="/magic-link?redirect=/chat"
                       class="w-full sm:w-auto px-6 py-3 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 transition">
            Sign in to Chat →
          </router-link>
          <router-link to="/about"
                       class="w-full sm:w-auto px-5 py-3 text-tea-700 hover:bg-tea-50 rounded-xl transition text-sm border border-tea-200">
            Learn about us
          </router-link>
        </div>
      </div>
    </div>

    <!-- 已登录：聊天主界面 -->
    <div v-else class="flex flex-col h-screen pt-14 md:pt-0">

      <!-- 顶部栏 -->
      <header class="flex items-center gap-3 px-3 md:px-6 py-3 border-b border-tea-100 bg-white sticky top-0 z-20">
        <button @click="mobileSidebar = true" class="md:hidden p-2 -ml-1 rounded-lg hover:bg-tea-100 text-tea-700" aria-label="Conversations">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
        </button>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <div class="w-9 h-9 md:w-10 md:h-10 rounded-full bg-gradient-to-br from-tea-600 to-tea-800 text-white flex items-center justify-center text-sm font-medium flex-shrink-0">🍃</div>
            <div class="min-w-0">
              <h1 class="font-serif text-base md:text-lg text-tea-900 truncate">Tea Advisor</h1>
              <div class="flex items-center gap-1.5 text-[11px] text-tea-500">
                <span class="w-1.5 h-1.5 rounded-full" :class="wsConnected ? 'bg-green-500 animate-pulse' : 'bg-tea-300'"></span>
                {{ wsConnected ? 'Online · replies in minutes' : 'Reconnecting...' }}
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- 主体 -->
      <div class="flex flex-1 min-h-0">

        <!-- PC 左侧会话列表 -->
        <aside class="hidden md:flex w-64 lg:w-72 flex-col border-r border-tea-100 bg-tea-50/40">
          <div class="px-4 py-3 border-b border-tea-100 text-xs font-semibold text-tea-600 uppercase tracking-wider">Conversations</div>
          <div class="flex-1 overflow-y-auto">
            <div v-if="convs.length === 0" class="p-6 text-center text-sm text-tea-400">Loading conversations...</div>
            <button v-for="c in convs" :key="c.id" @click="selectConv(c)"
                    :class="['w-full text-left px-4 py-3 border-b border-tea-100/60 transition',
                             active?.id === c.id ? 'bg-tea-700 text-white' : 'hover:bg-tea-100/60 text-tea-800']">
              <div class="font-medium text-sm">Conversation #{{ c.id }}</div>
              <div :class="['text-xs mt-0.5', active?.id === c.id ? 'text-tea-200' : 'text-tea-500']">With your tea advisor</div>
            </button>
          </div>
        </aside>

        <!-- 手机底部弹出会话列表 -->
        <div v-if="mobileSidebar" class="md:hidden fixed inset-0 z-40" @click.self="mobileSidebar = false">
          <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="mobileSidebar = false"></div>
          <div class="absolute bottom-0 inset-x-0 max-h-[70vh] bg-white rounded-t-2xl shadow-xl flex flex-col">
            <div class="flex items-center justify-between px-5 py-4 border-b border-tea-100">
              <span class="font-semibold text-tea-800">Conversations</span>
              <button @click="mobileSidebar = false" class="w-8 h-8 rounded-full bg-tea-100 hover:bg-tea-200 flex items-center justify-center text-tea-600">×</button>
            </div>
            <div class="flex-1 overflow-y-auto">
              <button v-for="c in convs" :key="c.id" @click="selectConv(c); mobileSidebar = false"
                      :class="['w-full text-left px-5 py-4 border-b border-tea-50 transition flex items-center gap-3',
                               active?.id === c.id ? 'bg-tea-50' : '']">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-tea-500 to-tea-700 text-white flex items-center justify-center text-sm flex-shrink-0">🍃</div>
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-sm text-tea-800">Conversation #{{ c.id }}</div>
                  <div class="text-xs text-tea-500">Your tea advisor</div>
                </div>
                <div v-if="active?.id === c.id" class="w-2 h-2 rounded-full bg-tea-600"></div>
              </button>
            </div>
          </div>
        </div>

        <!-- 聊天区 -->
        <section class="flex-1 flex flex-col min-w-0 bg-slate-50/50">
          <div class="flex-1 overflow-y-auto px-3 sm:px-6 py-4 md:py-6" ref="scrollRef">
            <div v-if="!active" class="h-full flex items-center justify-center">
              <div class="text-center text-tea-400">
                <div class="text-5xl mb-3">💬</div>
                <p class="text-sm">Loading your conversation...</p>
              </div>
            </div>
            <div v-else class="max-w-3xl mx-auto space-y-2.5 md:space-y-3">
              <div v-if="messages.length === 0" class="flex justify-start mb-4">
                <div class="max-w-[85%] md:max-w-[70%] rounded-2xl rounded-tl-sm px-4 py-3 bg-white border border-tea-100 shadow-sm">
                  <p class="text-sm text-tea-800 leading-relaxed">
                    👋 Hi! I'm your personal tea advisor. Ask me about Pu'er blends, traceability, SGS reports,
                    or anything about Yunnan tea gardens · 你好！我是你的普洱茶顾问 😊
                  </p>
                </div>
              </div>
              <div v-for="(m, i) in messages" :key="m.id || m.client_msg_id || i"
                   :class="['flex', m.sender_type === 'me' ? 'justify-end' : 'justify-start']">
                <MessageBubble :msg="m" :show-translation="showTranslation[m.id!]" @toggle-translation="toggleTranslation(m.id!)" />
              </div>
            </div>
          </div>

          <!-- 附件预览 -->
          <div v-if="pendingAttachments.length" class="border-t border-tea-100 bg-white px-3 md:px-4 py-2.5 flex gap-2 flex-wrap">
            <div v-for="(a, i) in pendingAttachments" :key="i" class="relative group">
              <img v-if="a.type === 'image'" :src="a.url" class="h-14 w-14 object-cover rounded-lg border border-tea-200" />
              <div v-else class="h-14 w-14 rounded-lg bg-tea-50 border border-tea-200 flex flex-col items-center justify-center text-[10px] text-tea-600">
                <span class="text-base">{{ a.type === 'video' ? '🎬' : '📎' }}</span>
              </div>
              <button @click="removePendingAttachment(i)" class="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-red-500 text-white text-xs leading-none opacity-0 group-hover:opacity-100 transition">×</button>
            </div>
          </div>

          <!-- 输入区 -->
          <div class="border-t border-tea-100 bg-white px-3 md:px-4 py-2.5 md:py-3 pb-[env(safe-area-inset-bottom)]">
            <div class="flex items-end gap-2">
              <div class="relative flex-shrink-0">
                <button @click="showUploadMenu = !showUploadMenu; showCardMenu = false" title="Attach" class="p-2 md:p-2.5 hover:bg-tea-100 rounded-lg transition text-tea-500">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
                </button>
                <div v-if="showUploadMenu" class="absolute bottom-full mb-2 left-0 bg-white border border-tea-100 rounded-xl shadow-xl p-1.5 z-30 w-36">
                  <label class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-tea-50 cursor-pointer text-sm text-tea-700">
                    🖼️ Image<input type="file" accept="image/*" class="hidden" @change="(e) => { handleFiles(e.target.files, 'image'); showUploadMenu = false }" />
                  </label>
                  <label class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-tea-50 cursor-pointer text-sm text-tea-700">
                    🎬 Video<input type="file" accept="video/*" class="hidden" @change="(e) => { handleFiles(e.target.files, 'video'); showUploadMenu = false }" />
                  </label>
                  <label class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-tea-50 cursor-pointer text-sm text-tea-700">
                    📎 File<input type="file" class="hidden" @change="(e) => { handleFiles(e.target.files, 'file'); showUploadMenu = false }" />
                  </label>
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <textarea v-model="input" @keyup.enter.exact.prevent="send" @paste="onPaste" rows="1"
                          placeholder="Type a message..."
                          class="w-full px-3.5 md:px-4 py-2 md:py-2.5 rounded-xl border border-tea-200 focus:border-tea-500 focus:outline-none focus:ring-2 focus:ring-tea-100 resize-none text-sm bg-tea-50/50 leading-relaxed max-h-32"></textarea>
              </div>
              <button @click="send" :disabled="!input.trim() && !pendingAttachments.length"
                      class="flex-shrink-0 px-4 md:px-5 py-2 md:py-2.5 bg-tea-700 text-white rounded-xl font-medium hover:bg-tea-800 active:bg-tea-900 transition disabled:opacity-40 disabled:cursor-not-allowed text-sm flex items-center gap-1.5">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="m22 2-7 20-4-9-9-4Z"/><path d="M22 2 11 13"/></svg>
                <span class="hidden sm:inline">Send</span>
              </button>
            </div>
            <div v-if="uploading" class="mt-2 h-1 bg-tea-100 rounded-full overflow-hidden">
              <div class="h-full bg-tea-600 transition-all duration-200" :style="{ width: uploadProgress + '%' }"></div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
