<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const messages = ref<any[]>([])
const activeConv = ref<any>(null)
const input = ref('')
const ws = ref<WebSocket | null>(null)

async function loadConvs() { try { const d: any = await api.get('/conversations'); list.value = d?.items || d || [] } catch {} }
onMounted(() => { loadConvs() })

async function selectConv(id: number) {
  activeConv.value = list.value.find(c => c.id === id)
  try { const d: any = await api.get(`/conversations/${id}/messages`); messages.value = d?.items || d || [] } catch {}
  // Try WebSocket
  try {
    const token = localStorage.getItem('staff_token')
    ws.value = new WebSocket(`ws://localhost:8080/api/v1/ws/im?token=${token}&conv_id=${id}`)
    ws.value.onmessage = (ev) => {
      try { const m = JSON.parse(ev.data); messages.value.push(m) } catch {}
    }
  } catch {}
}

async function send() {
  if (!activeConv.value || !input.value) return
  await api.post('/conversations', { conversation_type: 'single' }) // placeholder — actual message endpoint separate
  ElMessage.success('Sent (demo)')
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
            <div class="font-medium text-sm">{{ c.id }}</div>
            <div class="text-xs text-slate-500">single · {{ c.participants?.length || 0 }} participants</div>
          </div>
        </el-scrollbar>
      </el-col>
      <el-col :span="18" class="flex flex-col h-full">
        <div class="p-3 border-b font-medium flex items-center justify-between">
          <span>IM Chat {{ activeConv ? `#${activeConv.id}` : '—' }}</span>
          <el-tag v-if="ws?.readyState===1" size="small" type="success">WebSocket Connected</el-tag>
        </div>
        <el-scrollbar class="flex-1 p-4 bg-slate-50">
          <div v-if="!activeConv" class="text-center text-slate-400 mt-20">Select a conversation on the left to start chatting</div>
          <div v-for="(m,i) in messages" :key="i" class="mb-4" :class="m.sender_id==4 ? 'text-right' : ''">
            <div class="inline-block px-4 py-2 rounded-2xl max-w-[70%]" :class="m.sender_id==4 ? 'bg-tea-700 text-white' : 'bg-white border'">
              {{ m.content }}
            </div>
            <div v-if="m.translated_content" class="text-xs text-slate-500 mt-1">🇬🇧 {{ m.translated_content }}</div>
          </div>
        </el-scrollbar>
        <div class="p-3 border-t flex gap-2">
          <el-input v-model="input" placeholder="Type message... (auto-translated via FastAPI)" @keyup.enter="send" />
          <el-button type="primary" @click="send">Send</el-button>
        </div>
      </el-col>
    </el-row>
  </el-card>
</template>
