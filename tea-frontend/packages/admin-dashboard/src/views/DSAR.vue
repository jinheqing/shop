<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">GDPR — DSAR Requests</span>
          <el-button type="primary" :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-alert v-if="!loaded" title="Click Reload to fetch from backend" type="info" show-icon :closable="false" class="mb-4" />
      <el-table :data="requests" v-else stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="email" label="Email" width="200" />
        <el-table-column prop="request_type" label="Type" width="120" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitted_at" label="Submitted" width="180" />
        <el-table-column label="SLA" width="120">
          <template #default="{ row }">
            <span v-if="row.sla_deadline" :class="isOverdue(row) ? 'text-red-500 font-bold' : ''">
              {{ row.sla_deadline }}
            </span>
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="220">
          <template #default="{ row }">
            <el-button size="small" :disabled="row.status !== 'completed'" @click="exportData(row)">Export</el-button>
            <el-button size="small" type="danger" @click="erase(row)">Erase</el-button>
            <el-button size="small" type="success" @click="markDone(row)">Done</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { api } from '../api/client'

const requests = ref<any[]>([])
const loaded = ref(false)

async function load() {
  try {
    const res: any = await api.get('/dsar/requests')
    requests.value = res?.items || res || []
    loaded.value = true
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function statusTag(s: string) {
  const m: Record<string, string> = { pending: 'warning', processing: 'primary', completed: 'success', errored: 'danger' }
  return m[s] || 'info'
}

function isOverdue(row: any) {
  if (!row.sla_deadline || row.status === 'completed') return false
  return new Date(row.sla_deadline).getTime() < Date.now()
}

async function exportData(row: any) {
  try {
    const res = await api.get(`/dsar/requests/${row.id}/export`)
    const data = (res as any)?.data_export || res
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url; a.download = `dsar-${row.id}.json`; a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('Exported')
  } catch (e: any) {
    ElMessage.error('Export failed: ' + (e?.message || e))
  }
}

async function erase(row: any) {
  if (!confirm(`Erase all personal data for ${row.email}? This cannot be undone.`)) return
  try {
    await api.post(`/dsar/requests/${row.id}/delete`)
    ElMessage.success('Erased')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Erase failed')
  }
}

async function markDone(row: any) {
  try {
    await api.put(`/dsar/requests/${row.id}`, { status: 'completed' })
    ElMessage.success('Marked completed')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Update failed')
  }
}

onMounted(load)
</script>
