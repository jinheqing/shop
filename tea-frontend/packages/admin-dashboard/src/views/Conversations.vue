<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import { api } from '@/api/client'

// IM 会话管理页（staff 端）
// 协议与后端 /ws/im 对齐：
//   上行 {type:'send_message', payload:{conversation_id, content, message_type, client_msg_id}}
//   下行 {type:'chat_message'|'ack'|'error', payload:{...}}

const list = ref<any[]>([])
const messages = ref<any[]>([])
const activeConv = ref<any>(null)
const input = ref('')
const wsConnected = ref(false)
const ws = ref<WebSocket | null>(null)
const scrollRef = ref<HTMLElement | null>(null)

function wsURL() {
  const token = localStorage.getItem('staff_token')
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  // 后端 WS 注册在根路由 /ws/im（不在 /api/v1 下），dev 由 vite 代理转发
  return `${proto}://${location.host}/ws/im?token=${token}`
}

function connectWS(convID?: number) {
  ws.value?.close()
  const token = localStorage.getItem('staff_token')
  if (!token) return
  ws.value = new WebSocket(wsURL())
  ws.value.onopen = () => {
    wsConnected.value = true
    if (convID) {
      ws.value?.send(JSON.stringify({ type: 'join_conversation', payload: { conversation_id: convID } }))
    }
  }
  ws.value.onclose = () => { wsConnected.value = false }
  ws.value.onerror = () => { wsConnected.value = false }
  ws.value.onmessage = (ev) => {
    try {
      const env = JSON.parse(ev.data)
      if (env.type === 'chat_message') handleIncoming(env.payload)
      else if (env.type === 'ack') {
        const idx = messages.value.findIndex(m => m.client_msg_id === env.payload?.client_msg_id && m.local)
        if (idx >= 0) messages.value.splice(idx, 1)
      }
    } catch { /* 非 JSON 帧忽略 */ }
  }
}

function handleIncoming(msg: any) {
  if (msg?.id) {
    const exists = messages.value.find(m => m.id === msg.id)
    if (exists) {
      // 翻译结果回填
      if (msg.translation_status === 'translated' && (msg.translation_en || msg.translation_zh)) Object.assign(exists, msg)
      return
    }
  }
  messages.value.push(msg)
  scrollToBottom()
}

function scrollToBottom() { nextTick(() => { if (scrollRef.value) scrollRef.value.scrollTop = scrollRef.value.scrollHeight }) }

async function loadConvs() {
  // 后端 /conversations 返回 {code:0, data:[...]}（无 items 包装）
  try { const d: any = await api.get('/conversations'); list.value = d?.data || d?.items || [] } catch {}
}
onMounted(loadConvs)
onUnmounted(() => ws.value?.close())

async function selectConv(id: number) {
  activeConv.value = list.value.find(c => c.id === id)
  // 历史消息：后端返回 {code:0, data:{count, messages}}
  try {
    const d: any = await api.get(`/conversations/${id}/messages`)
    messages.value = d?.data?.messages || d?.items || d || []
  } catch { messages.value = [] }
  connectWS(id)
  scrollToBottom()
}

function makeClientID() { return 's_' + Date.now() + '_' + Math.random().toString(36).slice(2, 7) }

function send() {
  if (!activeConv.value || !input.value.trim()) return
  const cid = makeClientID()
  const content = input.value.trim()
  // 乐观插入本地 pending 消息，ack 后移除
  messages.value.push({ content, message_type: 'text', sender_type: 'staff', sender_id: -1, created_at: new Date().toISOString(), local: true, client_msg_id: cid })
  scrollToBottom()
  ws.value?.send(JSON.stringify({
    type: 'send_message',
    payload: { conversation_id: activeConv.value.id, content, message_type: 'text', client_msg_id: cid }
  }))
  input.value = ''
}
</script>

<template>
  <el-card class="h-[calc(100vh-180px)]">
    <el-row class="h-full">
      <el-col :span="6" class="border-r">
        <div class="p-3 border-b font-medium text-slate-700">Conversations ({{ list.length }})</div>
        <el-scrollbar class="h-[calc(100%-48px)]">
          <div v-for="c in list" :key="c.id" @click="selectConv(c.id)"
            :class="['p-3 border-b cursor-pointer hover:bg-slate-50', activeConv?.id===c.id?'bg-amber-50':'']">
            <div class="font-medium text-sm">#{{ c.id }} · {{ c.conversation_type || 'single' }}</div>
            <div class="text-xs text-slate-500">{{ c.participants?.length || 0 }} participants</div>
          </div>
        </el-scrollbar>
      </el-col>
      <el-col :span="18" class="flex flex-col h-full">
        <div class="p-3 border-b font-medium flex items-center justify-between">
          <span>IM Chat {{ activeConv ? `#${activeConv.id}` : '—' }}</span>
          <el-tag v-if="wsConnected" size="small" type="success">WebSocket Connected</el-tag>
          <el-tag v-else size="small" type="info">Offline</el-tag>
        </div>
        <div ref="scrollRef" class="flex-1 p-4 bg-slate-50 overflow-y-auto">
          <div v-if="!activeConv" class="text-center text-slate-400 mt-20">Select a conversation on the left to start chatting</div>
          <div v-for="(m,i) in messages" :key="m.id || m.client_msg_id || i"
               :class="['mb-4 flex', m.sender_type === 'staff' ? 'justify-end' : 'justify-start']">
            <div class="max-w-[70%]">
              <div class="inline-block px-4 py-2 rounded-2xl"
                   :class="m.sender_type === 'staff' ? 'bg-tea-700 text-white' : 'bg-white border'">
                <div class="text-sm">{{ m.content }}</div>
              </div>
              <!-- 实时翻译：后端异步翻译完成后经 WS 回填 translation_en / translation_zh -->
              <div v-if="m.translation_en || m.translation_zh" class="text-xs text-slate-500 mt-1">
                🌐 {{ m.translation_en || m.translation_zh }}
              </div>
              <div class="text-[10px] text-slate-400 mt-0.5">{{ m.sender_type }} · {{ new Date(m.created_at).toLocaleTimeString() }}</div>
            </div>
          </div>
        </div>
        <div class="p-3 border-t flex gap-2">
          <el-input v-model="input" placeholder="Type message... (auto-translated via FastAPI)" :disabled="!activeConv" @keyup.enter="send" />
          <el-button type="primary" :disabled="!activeConv || !input.trim()" @click="send">Send</el-button>
        </div>
      </el-col>
    </el-row>
  </el-card>
</template>
