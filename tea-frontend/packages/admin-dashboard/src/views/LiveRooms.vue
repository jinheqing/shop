<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp', order_id: null as number | null })
const presets = [
  { room_type: 'obs_tasting', push_source: 'obs_rtmp', label: 'OBS 品鉴间' },
  { room_type: 'app_scheduled', push_source: 'app_webrtc', label: '顾问排期 App' },
  { room_type: 'delivery_inspection', push_source: 'camera_rtmp', label: '验货间' },
  { room_type: 'customer_requested', push_source: 'app_webrtc', label: '用户请求' },
  { room_type: 'slow_247', push_source: 'camera_rtmp', label: '24/7 慢直播' },
  { room_type: 'custom_private', push_source: 'obs_rtmp', label: '定制私密' },
]

async function load() { const d: any = await api.get('/live-rooms'); list.value = d?.items || d || [] }
onMounted(load)

async function save() {
  await api.post('/live-rooms', form)
  ElMessage.success('Created!')
  open.value = false
  load()
}
async function del(id: number) { await api.delete(`/live-rooms/${id}`); load() }

function selectPreset(p: any) { form.room_type = p.room_type; form.push_source = p.push_source }
</script>

<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">Live Rooms</span><el-button type="primary" @click="open=true">+ New Room</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="room_name" label="Name" width="180" />
      <el-table-column prop="room_type" label="Type" width="160" />
      <el-table-column label="Push URL" min-width="240">
        <template #default="{ row }">
          <code class="text-xs text-slate-600">{{ row.obs_rtmp_url || row.camera_rtmp_url || '—' }}</code>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status==='live'?'success':'info'" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Key" width="140">
        <template #default="{ row }">
          <el-tooltip :content="row.obs_rtmp_key || ''" placement="top">
            <span class="text-xs font-mono">{{ (row.obs_rtmp_key || '—').slice(0, 12) }}...</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="100">
        <template #default="{ row }"><el-button size="small" type="danger" @click="del(row.id)">Del</el-button></template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="Create Live Room" width="560px">
    <div class="grid grid-cols-3 gap-2 mb-4">
      <el-button v-for="p in presets" :key="p.room_type" size="small" @click="selectPreset(p)" :type="form.room_type===p.room_type?'primary':''">{{ p.label }}</el-button>
    </div>
    <el-form :model="form" label-width="100px">
      <el-form-item label="Name"><el-input v-model="form.room_name" /></el-form-item>
      <el-form-item label="Type"><el-input v-model="form.room_type" /></el-form-item>
      <el-form-item label="Source"><el-input v-model="form.push_source" /></el-form-item>
      <el-form-item label="Order ID (optional)"><el-input-number v-model="form.order_id" :min="0" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open=false">Cancel</el-button>
      <el-button type="primary" @click="save">Create →</el-button>
    </template>
  </el-dialog>
</template>
