<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'

const status = ref<any>({
  fastapi_url: '',
  fastapi_health: 'unknown',
  gpu_available: false,
  queue_depth: 0,
  avg_latency_ms: 0,
  languages_supported: 0,
  ws_sessions_active: 0,
  total_translations_today: 0,
  last_translation_at: '',
})

const recent = ref<any[]>([])
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    // 1) 读 translate service_url 配置（从 System Settings 写入 DB 的值）
    const cfgRes: any = await api.get('/system/config/translate')
    if (cfgRes?.value?.service_url) {
      status.value.fastapi_url = cfgRes.value.service_url
    }

    // 2) 查真实 translate engine 状态
    try {
      const s: any = await api.get('/translate/status')
      status.value.fastapi_health = s?.status || 'ok'
      status.value.gpu_available = !!s?.gpu_available
      status.value.queue_depth = s?.queue_depth ?? 0
      status.value.avg_latency_ms = s?.avg_latency_ms ?? 0
      status.value.total_translations_today = s?.total_translations_today ?? 0
      if (Array.isArray(s?.recent)) {
        recent.value = s.recent
      }
    } catch {
      // translate engine 不可用（dev 环境常见），保持 fastapi_health=unknown
      status.value.fastapi_health = status.value.fastapi_url ? 'unreachable' : 'not_configured'
    }
  } catch (e: any) {
    errorMsg.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
<template>
  <div class="space-y-4">
    <el-alert v-if="errorMsg" :title="errorMsg" type="error" show-icon :closable="false" class="mb-2" />
    <el-alert v-if="loading" title="Loading translate engine status…" type="info" show-icon :closable="false" class="mb-2" />

    <el-row :gutter="16">
      <el-col :span="6"><el-card shadow="hover">
        <div class="text-xs text-slate-500 mb-1">Engine Health</div>
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 rounded-full" :class="{
            'bg-green-500': status.fastapi_health === 'ok' || status.fastapi_health === 'healthy',
            'bg-red-500': status.fastapi_health === 'unreachable' || status.fastapi_health === 'down',
            'bg-slate-400': true,
          }"></span>
          <span class="font-serif text-xl">{{ status.fastapi_health }}</span>
        </div>
      </el-card></el-col>
      <el-col :span="6"><el-card shadow="hover">
        <div class="text-xs text-slate-500 mb-1">Avg Latency</div>
        <div class="font-serif text-xl">{{ status.avg_latency_ms }} ms</div>
      </el-card></el-col>
      <el-col :span="6"><el-card shadow="hover">
        <div class="text-xs text-slate-500 mb-1">Today's Translations</div>
        <div class="font-serif text-xl">{{ status.total_translations_today.toLocaleString() }}</div>
      </el-card></el-col>
      <el-col :span="6"><el-card shadow="hover">
        <div class="text-xs text-slate-500 mb-1">Active Sessions</div>
        <div class="font-serif text-xl">{{ status.ws_sessions_active }} · Q {{ status.queue_depth }}</div>
      </el-card></el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">🌐 Translation Engine Status</span>
          <div class="flex gap-2">
            <el-button size="small" @click="load" :disabled="loading">Reload</el-button>
            <el-link href="#/system-settings" :underline="false">→ Configure in System Settings</el-link>
          </div>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="FastAPI Endpoint"><code>{{ status.fastapi_url || '(not set)' }}</code></el-descriptions-item>
        <el-descriptions-item label="GPU Accelerated">{{ status.gpu_available ? '✅ Yes' : '❌ No (CPU)' }}</el-descriptions-item>
        <el-descriptions-item label="Languages Supported">{{ status.languages_supported || '—' }}</el-descriptions-item>
        <el-descriptions-item label="Last Translation">{{ status.last_translation_at || '—' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card>
      <template #header><span class="font-medium">📜 Recent Translation Jobs</span></template>
      <el-alert v-if="!recent.length" title="No recent jobs (translate engine not running)" type="info" show-icon :closable="false" class="mb-2" />
      <el-table :data="recent" stripe>
        <el-table-column prop="direction" label="Direction" width="140" />
        <el-table-column prop="source_type" label="Source" width="150" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="({'completed':'success','translating':'warning','failed':'danger','pending':'info'} as Record<string,string>)[row.status]" effect="dark">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Latency" width="140">
          <template #default="{ row }">{{ row.latency_ms ? row.latency_ms + ' ms' : '—' }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
