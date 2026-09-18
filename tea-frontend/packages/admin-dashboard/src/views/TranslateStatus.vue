<script setup lang="ts">
import { ref } from 'vue'

const status = ref({
  fastapi_url: 'http://localhost:8090',
  fastapi_health: 'healthy',
  model_nllb: 'nllb-200-1.3B-ct2-int8',
  model_whisper: 'whisper-base-ct2',
  gpu_available: false,
  queue_depth: 3,
  avg_latency_ms: 420,
  languages_supported: 202,
  ws_sessions_active: 12,
  total_translations_today: 1847,
  last_translation_at: '2026-09-17T14:05:22Z',
})

const recent = ref([
  { id: 1, direction: 'zh_en', source_type: 'im_text', status: 'completed', latency_ms: 380 },
  { id: 2, direction: 'en_zh', source_type: 'im_text', status: 'completed', latency_ms: 412 },
  { id: 3, direction: 'zh_en', source_type: 'live_subtitle', status: 'completed', latency_ms: 560 },
  { id: 4, direction: 'en_zh', source_type: 'im_text', status: 'failed', latency_ms: null },
])
</script>
<template>
  <div class="space-y-4">
    <el-row :gutter="16">
      <el-col :span="6"><el-card shadow="hover">
        <div class="text-xs text-slate-500 mb-1">Engine Health</div>
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 rounded-full" :class="status.fastapi_health==='healthy'?'bg-green-500':'bg-red-500'"></span>
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
      <template #header><span class="font-medium">🌐 Translation Engine Status</span></template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="FastAPI Endpoint"><code>{{ status.fastapi_url }}</code></el-descriptions-item>
        <el-descriptions-item label="GPU Accelerated">{{ status.gpu_available ? '✅ Yes' : '❌ No (CPU)' }}</el-descriptions-item>
        <el-descriptions-item label="NLLB Model"><code>{{ status.model_nllb }}</code></el-descriptions-item>
        <el-descriptions-item label="Whisper Model"><code>{{ status.model_whisper }}</code></el-descriptions-item>
        <el-descriptions-item label="Languages Supported">{{ status.languages_supported }}</el-descriptions-item>
        <el-descriptions-item label="Last Translation">{{ status.last_translation_at }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card>
      <template #header><span class="font-medium">📜 Recent Translation Jobs</span></template>
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
