<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ name: '', location: '', description: '', camera_rtmp_url: '' })

async function load() { try { const d: any = await api.get('/slow-presets'); list.value = d?.items || d || [] } catch {} }
onMounted(load)

async function save() { await api.post('/slow-presets', form); ElMessage.success('Created'); open.value = false; load() }
async function enable(id: number) { await api.post(`/slow-presets/${id}/enable`); ElMessage.success('Streaming'); load() }
async function disable(id: number) { await api.post(`/slow-presets/${id}/disable`); ElMessage.success('Stopped'); load() }
async function del(id: number) { await api.delete(`/slow-presets/${id}`); load() }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🌱 24/7 Slow-Live Camera Presets</span><el-button type="primary" @click="open=true">+ New Preset</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="name" label="Name" width="160" />
      <el-table-column prop="location" label="Location" width="160" />
      <el-table-column prop="description" label="Description" show-overflow-tooltip />
      <el-table-column prop="status" label="Status" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status==='live'?'success':'info'" effect="dark">
            <span v-if="row.status==='live'" class="flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full bg-white animate-pulse"></span>LIVE</span>
            <span v-else>{{ row.status }}</span>
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="RTMP URL" min-width="260">
        <template #default="{ row }"><code class="text-xs text-slate-600">{{ row.camera_rtmp_url || '—' }}</code></template>
      </el-table-column>
      <el-table-column label="Actions" width="200">
        <template #default="{ row }">
          <el-button v-if="row.status!=='live'" size="small" type="success" @click="enable(row.id)">▶ Enable</el-button>
          <el-button v-else size="small" @click="disable(row.id)">⏸ Stop</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="New Slow-Live Preset" width="480px">
    <el-form :model="form" label-width="100px">
      <el-form-item label="Name"><el-input v-model="form.name" placeholder="云南省临沧市临翔区邦东乡曼岗村茶园" /></el-form-item>
      <el-form-item label="Location"><el-input v-model="form.location" placeholder="临沧" /></el-form-item>
      <el-form-item label="Description"><el-input v-model="form.description" type="textarea" /></el-form-item>
      <el-form-item label="RTMP URL"><el-input v-model="form.camera_rtmp_url" placeholder="rtmp://camera-ip/live/key" /></el-form-item>
    </el-form>
    <template #footer><el-button @click="open=false">Cancel</el-button><el-button type="primary" @click="save">Save & Enable Later</el-button></template>
  </el-dialog>
</template>
