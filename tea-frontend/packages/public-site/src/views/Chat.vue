<script setup lang="ts">
import { onMounted, ref, nextTick } from 'vue'
import { api } from '@/api/client'

const convs = ref<any[]>([])
const active = ref<any>(null)
const messages = ref<any[]>([])
const input = ref('')
const ws = ref<WebSocket | null>(null)

async function loadConvs() { try { convs.value = (await api.get('/conversations')).items || [] } catch {} }
onMounted(async () => {
  await loadConvs()
  if (!convs.value.length) {
    // Create default conversation with an advisor
    await api.post('/conversations', { other_user_id: 1 })
    await loadConvs()
  }
  if (convs.value.length) selectConv(convs.value[0])
})

async function selectConv(c: any) {
  active.value = c
  try { messages.value = (await api.get(`/conversations/${c.id}/messages`)).items || [] } catch {}
  const token = localStorage.getItem('user_token')
  if (token) {
    ws.value?.close()
    ws.value = new WebSocket(`ws://localhost:8080/api/v1/ws/im?token=${token}&conv_id=${c.id}`)
    ws.value.onmessage = (ev) => { try { messages.value.push(JSON.parse(ev.data)) } catch {} }
  }
  await nextTick()
}

function send() {
  if (!input.value || !active.value) return
  ws.value?.send(JSON.stringify({ type: 'message', content: input.value, conv_id: active.value.id }))
  messages.value.push({ content: input.value, sender_role: 'user', created_at: new Date().toISOString() })
  input.value = ''
}
</script>
<template>
  <div class="pt-20 min-h-screen bg-tea-50">
    <div class="max-w-5xl mx-auto px-6 py-12">
      <h1 class="font-serif text-3xl text-tea-900 mb-2">💬 Chat with Your Advisor</h1>
      <p class="text-tea-600 mb-6">Speak to our team directly. All messages auto-translate between English and Chinese.</p>

      <div class="bg-white rounded-2xl border border-tea-100 overflow-hidden">
        <div class="flex h-[500px]">
          <aside class="w-64 border-r bg-tea-50 overflow-auto">
            <div class="p-3 border-b font-medium text-sm">Your Conversations</div>
            <div v-for="c in convs" :key="c.id" @click="selectConv(c)"
              :class="['p-3 border-b cursor-pointer text-sm hover:bg-white', active?.id===c.id?'bg-white border-l-2 border-l-tea-600':'']">
              <div class="font-medium">Conversation #{{ c.id }}</div>
              <div class="text-xs text-tea-500">2 participants</div>
            </div>
          </aside>
          <section class="flex-1 flex flex-col">
            <div class="flex-1 overflow-auto p-6 bg-slate-50">
              <div v-if="!active" class="text-center text-tea-500 mt-20">Select a conversation →</div>
              <div v-for="(m,i) in messages" :key="i" class="mb-4" :class="m.sender_role==='user' ? 'text-right' : ''">
                <div class="inline-block px-4 py-2 rounded-2xl max-w-[70%]" :class="m.sender_role==='user' ? 'bg-tea-700 text-white' : 'bg-white border border-tea-200'">
                  {{ m.content }}
                </div>
                <div v-if="m.translated_content && m.sender_role!=='user'" class="text-xs text-tea-500 mt-1">🇬🇧 {{ m.translated_content }}</div>
              </div>
            </div>
            <div class="p-4 border-t flex gap-2">
              <input v-model="input" @keyup.enter="send" placeholder="Type your message..." class="flex-1 px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
              <button @click="send" class="px-6 py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition">Send</button>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>
