<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Live Calendar</span>
          <el-button :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-alert v-if="!loaded" title="Click Reload to fetch from backend" type="info" show-icon :closable="false" class="mb-4" />
      <el-table :data="sessions" v-else stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="Title" min-width="180" />
        <el-table-column prop="room_type" label="Type" width="160">
          <template #default="{ row }">
            <el-tag size="small">{{ row.room_type || 'obs_tasting' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scheduled_start" label="Scheduled" width="180" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="statusColor(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="Location" min-width="200" />
        <el-table-column prop="order_id" label="Order" width="100" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { api } from '../api/client'

const sessions = ref<any[]>([])
const loaded = ref(false)

async function load() {
  try {
    const res: any = await api.get('/live-rooms/schedule/calendar')
    sessions.value = res?.items || res?.sessions || res || []
    loaded.value = true
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function statusColor(s: string) {
  const m: Record<string, string> = { scheduled: 'info', configuring: 'warning', live: 'danger', ended: 'success' }
  return m[s] || 'info'
}

onMounted(load)
</script>
