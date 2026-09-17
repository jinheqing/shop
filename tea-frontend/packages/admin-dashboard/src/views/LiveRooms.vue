<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-bold text-lg">Live Rooms</span>
        <el-button type="primary" @click="open = true">+ New Room</el-button>
      </div>
    </template>
    <el-table :data="list" stripe>
      <el-table-column prop="room_name" label="Name" width="180" />
      <el-table-column prop="room_type" label="Type" width="160">
        <template #default="{ row }">
          <el-tag size="small">{{ row.room_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Push URL" min-width="240">
        <template #default="{ row }">
          <code class="text-xs text-slate-600">{{ row.obs_rtmp_url || row.camera_rtmp_url || '—' }}</code>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status === 'live' ? 'danger' : row.status === 'ended' ? 'success' : 'info'" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="scheduled_start" label="Scheduled" width="170" />
      <el-table-column label="Actions" width="280">
        <template #default="{ row }">
          <el-button size="small" v-if="row.status !== 'live'" type="success" @click="start(row)">Start</el-button>
          <el-button size="small" v-else type="warning" @click="end(row)">End</el-button>
          <el-button size="small" @click="copyKey(row)">Key</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="Create Live Room" width="560px">
    <div class="mb-2 text-xs text-gray-500">Select a preset (matches backend RoomType constants)</div>
    <el-radio-group v-model="form.room_type" class="mb-4">
      <el-radio-button v-for="p in presets" :key="p.room_type" :value="p.room_type">{{ p.label }}</el-radio-button>
    </el-radio-group>
    <el-form :model="form" label-width="120px">
      <el-form-item label="Room Name"><el-input v-model="form.room_name" /></el-form-item>
      <el-form-item label="Push Source">
        <el-select v-model="form.push_source" style="width:100%">
          <el-option label="OBS RTMP" value="obs_rtmp" />
          <el-option label="Camera RTMP (slow)" value="camera_rtmp" />
          <el-option label="App WebRTC" value="app_webrtc" />
        </el-select>
      </el-form-item>
      <el-form-item label="Order ID (optional)"><el-input-number v-model="form.order_id" :min="0" style="width:100%" /></el-form-item>
      <el-form-item label="Location"><el-input v-model="form.location" placeholder="e.g. 云南省临沧市临翔区邦东乡曼岗村茶园" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">Cancel</el-button>
      <el-button type="primary" @click="save">Create</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp', order_id: null as number | null, location: '' })

// ✅ 与后端 RoomType 常量对齐
const presets = [
  { room_type: 'obs_tasting', label: 'OBS 品鉴' },
  { room_type: 'open_calendar', label: '开放预约' },
  { room_type: 'delivery_inspection', label: '验货间' },
  { room_type: 'customer_request', label: '客户请求' },
  { room_type: 'slow_preset', label: '24/7 慢直播' },
  { room_type: 'custom_private', label: '定制私密' },
]

async function load() {
  const d: any = await api.get('/live-rooms')
  list.value = d?.items || d || []
}
onMounted(load)

async function save() {
  await api.post('/live-rooms', form)
  ElMessage.success('Created!')
  open.value = false
  load()
}
async function del(id: number) {
  await api.delete(`/live-rooms/${id}`); load()
}
async function start(row: any) {
  await api.post(`/live-rooms/${row.id}/start`)
  ElMessage.success('Started')
  load()
}
async function end(row: any) {
  await api.post(`/live-rooms/${row.id}/end`)
  ElMessage.success('Ended')
  load()
}
function copyKey(row: any) {
  const k = row.obs_rtmp_key || row.camera_rtmp_key || ''
  if (k) { navigator.clipboard.writeText(k); ElMessage.success('Key copied') }
  else ElMessage.info('No RTMP key yet — start stream first')
}
</script>
