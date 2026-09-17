<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Delivery Inspection Rooms</span>
          <el-button type="primary" @click="dialog = true">+ Auto-Create Room for Order</el-button>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="Title" min-width="180" />
        <el-table-column prop="order_id" label="Order" width="100" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="statusColor(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scheduled_start" label="Scheduled" width="180" />
        <el-table-column label="Actions" width="200">
          <template #default="{ row }">
            <el-button size="small" v-if="row.status !== 'live'" type="success" @click="start(row)">Start</el-button>
            <el-button size="small" v-if="row.status === 'live'" type="danger" @click="end(row)">End</el-button>
            <el-button size="small" @click="openInLiveKit(row)">Open in LiveKit</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialog" title="Auto-Create Delivery Inspection Room" width="480px">
      <el-form :model="form">
        <el-form-item label="Order ID">
          <el-input-number v-model="form.order_id" :min="1" style="width:100%" />
        </el-form-item>
        <el-form-item label="Location">
          <el-input v-model="form.location" placeholder="Tea garden / Farm warehouse" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">Cancel</el-button>
        <el-button type="primary" @click="submitCreate">Create</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

const items = ref<any[]>([])
const dialog = ref(false)
const form = ref({ order_id: 1, location: '', description: '' })

async function load() {
  try {
    const res: any = await api.get('/live-rooms')
    const all = res?.items || res || []
    items.value = all.filter((r: any) => r.room_type === 'delivery_inspection')
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function statusColor(s: string) {
  const m: Record<string, string> = { scheduled: 'info', configuring: 'warning', live: 'danger', ended: 'success' }
  return m[s] || 'info'
}

async function submitCreate() {
  try {
    await api.post('/live-rooms/system-create', {
      order_id: form.value.order_id,
      room_type: 'delivery_inspection',
      location: form.value.location,
      description: form.value.description
    })
    ElMessage.success('Room created')
    dialog.value = false
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Create failed')
  }
}

async function start(row: any) {
  await api.post(`/live-rooms/${row.id}/start`)
  ElMessage.success('Stream started')
  await load()
}

async function end(row: any) {
  await api.post(`/live-rooms/${row.id}/end`)
  ElMessage.success('Stream ended')
  await load()
}

function openInLiveKit(row: any) {
  window.open(`https://cloud.livekit.io/rooms/${encodeURIComponent(row.name || row.room_name || '')}`, '_blank')
}

onMounted(load)
</script>
