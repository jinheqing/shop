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
        ? 'bg-ink-900 text-ivory-100 rounded-2xl rounded-tr-sm'
        : 'bg-white border border-gold/20 rounded-2xl rounded-tl-sm'
      const wrapperAlign = isMine.value ? 'items-end' : 'items-start'

      // Build children
      const children: any[] = []

      // Card payload (quote/order)
      if (isCard.value && m.card_payload) {
        const cp = m.card_payload
        if (m.message_type === 'quote_card') {
          children.push(
            h('div', { class: `max-w-xs shadow-lg overflow-hidden rounded-xl ${isMine.value ? 'bg-ink-900/95 text-ivory-100' : 'bg-white text-ink-900 border border-gold/20'}` }, [
              h('div', { class: 'p-4' }, [
                h('div', { class: 'text-[10px] uppercase tracking-lux text-gold font-sans mb-1' }, 'Quotation'),
                h('div', { class: 'font-serif text-lg font-semibold' }, cp.title || "Bespoke Pu-er Tea"),
                cp.sku ? h('div', { class: 'text-xs opacity-70 mt-1' }, `SKU: ${cp.sku}`) : null,
                h('div', { class: 'flex items-baseline gap-2 mt-2' }, [
                  h('span', { class: 'text-2xl font-bold' }, `£${cp.price ?? '—'}`),
                  cp.lead_time ? h('span', { class: 'text-xs opacity-60' }, `${cp.lead_time || '10-15 days'} lead time`) : null
                ]),
                h('button', {
                  class: `mt-3 w-full py-2 rounded-lg text-sm font-medium transition ${isMine.value ? 'bg-ivory-100 text-ink-900 hover:bg-ivory-50' : 'bg-ink-900 text-ivory-100 hover:bg-ink-800'}`,
                  onClick: () => window.open(cp.redirect_url || '#', '_blank')
                }, 'View Details →')
              ])
            ])
          )
        } else {
          // order_card
          const statusColor: any = { ordering:'bg-ivory-100 text-sand', paid:'bg-ink-900 text-ivory-100', producing:'bg-gold/20 text-gold', shipped:'bg-ink-800 text-ivory-100', completed:'bg-ink-900 text-ivory-100', cancelled:'bg-ink-900/40 text-ivory-100' }
          children.push(
            h('div', { class: `max-w-xs shadow-lg overflow-hidden rounded-xl ${isMine.value ? 'bg-ink-900/95 text-ivory-100' : 'bg-white text-ink-900 border border-gold/20'}` }, [
              h('div', { class: 'p-4' }, [
                h('div', { class: 'flex items-center justify-between mb-2' }, [
                  h('div', { class: 'text-[10px] uppercase tracking-lux text-gold font-sans' }, 'Order'),
                  h('span', { class: `text-[10px] uppercase tracking-lux px-2 py-0.5 rounded-full ${statusColor[m.status] || 'bg-ivory-100 text-sand'}` }, m.status || '—')
                ]),
                h('div', { class: 'font-serif text-lg font-semibold' }, cp.order_no || 'Order'),
                h('div', { class: 'mt-3 text-2xl font-bold' }, `£${cp.total ?? '—'}`),
                h('button', {
                  class: `mt-3 w-full py-2 rounded-lg text-sm font-medium transition ${isMine.value ? 'bg-ivory-100 text-ink-900 hover:bg-ivory-50' : 'bg-ink-900 text-ivory-100 hover:bg-ink-800'}`,
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
              children.push(h('a', { href: a.url, download: a.name, class: `max-w-xs p-3 rounded-lg flex items-center gap-2 transition ${isMine.value ? 'bg-ink-800/60 text-ivory-100 hover:bg-ink-800' : 'bg-ivory-50 text-sand hover:bg-ivory-100'}` }, [
                h('span', { class: 'text-xl' }, '·'),
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
        !isMine.value ? h('div', { class: 'text-[10px] text-sand px-1 font-sans' }, m.sender_type === 'staff' ? 'Tea Advisor' : 'You') : null,
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
  <div class="chat-page min-h-screen bg-ivory-100">

    <!-- 未登录引导页 -->
    <div v-if="!isLoggedIn" class="pt-24 px-4 text-center">
      <div class="max-w-md mx-auto">
        <div class="flex items-center justify-center mb-5">
          <div class="w-14 h-14 border border-gold/40 flex items-center justify-center" style="border-radius: 2px;">
            <span class="w-2 h-2 rounded-full bg-gold"></span>
          </div>
        </div>
        <h1 class="font-serif text-3xl text-ink-900 mb-3">Speak · With · Your · Tea · Advisor</h1>
        <p class="text-sand text-sm leading-relaxed mb-8 font-serif">
          Sign in to chat with a Pu'er tea expert. Bespoke blending, garden stories,
          and instant answers — in English or Chinese.
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
          <router-link to="/magic-link?redirect=/chat"
                       class="w-full sm:w-auto px-6 py-4 bg-ink-900 text-ivory-100 text-[11px] uppercase tracking-lux font-sans hover:bg-ink-800 transition"
                       style="border-radius: 2px;">
            Sign · In · To · Chat
          </router-link>
          <router-link to="/about"
                       class="w-full sm:w-auto px-5 py-4 border border-gold/30 text-ink-900 hover:bg-ivory-50 transition text-[11px] uppercase tracking-lux font-sans"
                       style="border-radius: 2px;">
            Learn · About · Us
          </router-link>
        </div>
      </div>
    </div>

    <!-- 已登录：聊天主界面 -->
    <div v-else class="flex flex-col h-screen pt-14 md:pt-0">

      <!-- 顶部栏 -->
      <header class="flex items-center gap-3 px-3 md:px-6 py-3 border-b border-gold/15 bg-white sticky top-0 z-20">
        <button @click="mobileSidebar = true" class="md:hidden p-2 -ml-1 hover:bg-ivory-100 text-ink-900" aria-label="Conversations" style="border-radius: 2px;">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
        </button>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <div class="w-9 h-9 md:w-10 md:h-10 bg-ink-900 text-ivory-100 flex items-center justify-center text-sm font-medium flex-shrink-0 font-serif" style="border-radius: 2px;">TH</div>
            <div class="min-w-0">
              <h1 class="font-serif text-base md:text-lg text-ink-900 truncate">Tea · Advisor</h1>
              <div class="flex items-center gap-1.5 text-[11px] text-sand font-sans">
                <span class="w-1.5 h-1.5 rounded-full" :class="wsConnected ? 'bg-gold' : 'bg-sand/50'"></span>
                {{ wsConnected ? 'Online · replies in minutes' : 'Reconnecting...' }}
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- 主体 -->
      <div class="flex flex-1 min-h-0">

        <!-- PC 左侧会话列表 -->
        <aside class="hidden md:flex w-64 lg:w-72 flex-col border-r border-gold/15 bg-ivory-50">
          <div class="px-4 py-3 border-b border-gold/15 text-[10px] uppercase tracking-lux text-gold font-sans">Conversations</div>
          <div class="flex-1 overflow-y-auto">
            <div v-if="convs.length === 0" class="p-6 text-center text-sm text-sand/70 font-serif">Loading conversations...</div>
            <button v-for="c in convs" :key="c.id" @click="selectConv(c)"
                    :class="['w-full text-left px-4 py-3 border-b border-gold/10 transition',
                             active?.id === c.id ? 'bg-ink-900 text-ivory-100' : 'hover:bg-ivory-100 text-ink-900']">
              <div class="font-medium text-sm font-serif">Conversation #{{ c.id }}</div>
              <div :class="['text-xs mt-0.5 font-sans', active?.id === c.id ? 'text-ivory-100/60' : 'text-sand']">With your tea advisor</div>
            </button>
          </div>
        </aside>

        <!-- 手机底部弹出会话列表 -->
        <div v-if="mobileSidebar" class="md:hidden fixed inset-0 z-40" @click.self="mobileSidebar = false">
          <div class="absolute inset-0 bg-ink-900/40 backdrop-blur-sm" @click="mobileSidebar = false"></div>
          <div class="absolute bottom-0 inset-x-0 max-h-[70vh] bg-white shadow-xl flex flex-col" style="border-top-left-radius: 2px; border-top-right-radius: 2px;">
            <div class="flex items-center justify-between px-5 py-4 border-b border-gold/15">
              <span class="font-serif text-ink-900">Conversations</span>
              <button @click="mobileSidebar = false" class="w-8 h-8 bg-ivory-100 hover:bg-ivory-50 flex items-center justify-center text-sand" style="border-radius: 2px;">×</button>
            </div>
            <div class="flex-1 overflow-y-auto">
              <button v-for="c in convs" :key="c.id" @click="selectConv(c); mobileSidebar = false"
                      :class="['w-full text-left px-5 py-4 border-b border-gold/10 transition flex items-center gap-3',
                               active?.id === c.id ? 'bg-ivory-100' : '']">
                <div class="w-10 h-10 bg-ink-800 text-ivory-100 flex items-center justify-center text-sm font-serif flex-shrink-0" style="border-radius: 2px;">TH</div>
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-sm text-ink-900 font-serif">Conversation #{{ c.id }}</div>
                  <div class="text-xs text-sand font-sans">Your tea advisor</div>
                </div>
                <div v-if="active?.id === c.id" class="w-2 h-2 rounded-full bg-gold"></div>
              </button>
            </div>
          </div>
        </div>

        <!-- 聊天区 -->
        <section class="flex-1 flex flex-col min-w-0 bg-ivory-50">
          <div class="flex-1 overflow-y-auto px-3 sm:px-6 py-4 md:py-6" ref="scrollRef">
            <div v-if="!active" class="h-full flex items-center justify-center">
              <div class="text-center text-sand/70">
                <div class="flex items-center justify-center mb-3">
                  <div class="w-10 h-10 border border-gold/30 flex items-center justify-center" style="border-radius: 2px;">
                    <span class="w-1.5 h-1.5 rounded-full bg-gold/60"></span>
                  </div>
                </div>
                <p class="text-sm font-serif">Loading your conversation...</p>
              </div>
            </div>
            <div v-else class="max-w-3xl mx-auto space-y-2.5 md:space-y-3">
              <div v-if="messages.length === 0" class="flex justify-start mb-4">
                <div class="max-w-[85%] md:max-w-[70%] rounded-2xl rounded-tl-sm px-4 py-3 bg-white border border-gold/15 shadow-sm">
                  <p class="text-sm text-ink-900 leading-relaxed font-serif">
                    Welcome. I'm your personal tea advisor. Ask me about Pu'er blends, traceability, SGS reports, or anything about Yunnan tea gardens. · 你好！我是你的普洱茶顾问。
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
          <div v-if="pendingAttachments.length" class="border-t border-gold/15 bg-white px-3 md:px-4 py-2.5 flex gap-2 flex-wrap">
            <div v-for="(a, i) in pendingAttachments" :key="i" class="relative group">
              <img v-if="a.type === 'image'" :src="a.url" class="h-14 w-14 object-cover border border-gold/20" style="border-radius: 2px;" />
              <div v-else class="h-14 w-14 bg-ivory-50 border border-gold/20 flex flex-col items-center justify-center text-[10px] text-sand font-sans" style="border-radius: 2px;">
                <span class="text-[10px] uppercase tracking-lux">{{ a.type === 'video' ? 'VID' : 'DOC' }}</span>
              </div>
              <button @click="removePendingAttachment(i)" class="absolute -top-1.5 -right-1.5 w-5 h-5 bg-ink-900 text-ivory-100 text-xs leading-none opacity-0 group-hover:opacity-100 transition" style="border-radius: 2px;">×</button>
            </div>
          </div>

          <!-- 输入区 -->
          <div class="border-t border-gold/15 bg-white px-3 md:px-4 py-2.5 md:py-3 pb-[env(safe-area-inset-bottom)]">
            <div class="flex items-end gap-2">
              <div class="relative flex-shrink-0">
                <button @click="showUploadMenu = !showUploadMenu; showCardMenu = false" title="Attach" class="p-2 md:p-2.5 hover:bg-ivory-100 text-sand transition" style="border-radius: 2px;">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
                </button>
                <div v-if="showUploadMenu" class="absolute bottom-full mb-2 left-0 bg-white border border-gold/20 shadow-xl p-1.5 z-30 w-36" style="border-radius: 2px;">
                  <label class="flex items-center gap-2 px-3 py-2 hover:bg-ivory-50 cursor-pointer text-sm text-sand font-serif">
                    Image<input type="file" accept="image/*" class="hidden" @change="(e) => { handleFiles(e.target.files, 'image'); showUploadMenu = false }" />
                  </label>
                  <label class="flex items-center gap-2 px-3 py-2 hover:bg-ivory-50 cursor-pointer text-sm text-sand font-serif">
                    Video<input type="file" accept="video/*" class="hidden" @change="(e) => { handleFiles(e.target.files, 'video'); showUploadMenu = false }" />
                  </label>
                  <label class="flex items-center gap-2 px-3 py-2 hover:bg-ivory-50 cursor-pointer text-sm text-sand font-serif">
                    File<input type="file" class="hidden" @change="(e) => { handleFiles(e.target.files, 'file'); showUploadMenu = false }" />
                  </label>
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <textarea v-model="input" @keyup.enter.exact.prevent="send" @paste="onPaste" rows="1"
                          placeholder="Type a message..."
                          class="w-full px-3.5 md:px-4 py-2 md:py-2.5 border border-gold/20 focus:border-gold focus:outline-none resize-none text-sm bg-ivory-50 leading-relaxed max-h-32 font-serif" style="border-radius: 2px;"></textarea>
              </div>
              <button @click="send" :disabled="!input.trim() && !pendingAttachments.length"
                      class="flex-shrink-0 px-4 md:px-5 py-2 md:py-2.5 bg-ink-900 text-ivory-100 hover:bg-ink-800 active:bg-ink-950 transition disabled:opacity-40 disabled:cursor-not-allowed text-sm flex items-center gap-1.5 font-sans" style="border-radius: 2px;">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="m22 2-7 20-4-9-9-4Z"/><path d="M22 2 11 13"/></svg>
                <span class="hidden sm:inline text-[11px] uppercase tracking-lux">Send</span>
              </button>
            </div>
            <div v-if="uploading" class="mt-2 h-1 bg-ivory-100 overflow-hidden" style="border-radius: 2px;">
              <div class="h-full bg-gold transition-all duration-200" :style="{ width: uploadProgress + '%' }"></div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
