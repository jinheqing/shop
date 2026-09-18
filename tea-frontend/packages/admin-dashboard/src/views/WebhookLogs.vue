<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Webhook & Integration Logs</span>
          <el-tag type="info">via Audit Logs (target_type=webhook)</el-tag>
        </div>
      </template>
      <el-alert title="Webhook logs are stored as audit logs with target_type='webhook' or 'payment_webhook'. Filter below." type="info" show-icon :closable="false" class="mb-4" />
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="action" label="Action" width="180" />
        <el-table-column prop="target_type" label="Target" width="160" />
        <el-table-column prop="target_id" label="Target ID" width="100" />
        <el-table-column label="Detail" min-width="240">
          <template #default="{ row }">
            <span class="text-xs text-gray-500">{{ JSON.stringify(row.detail).slice(0, 180) }}{{ JSON.stringify(row.detail).length > 180 ? '…' : '' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP" width="140" />
        <el-table-column prop="created_at" label="Time" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

const items = ref<any[]>([])

async function load() {
  try {
    const res: any = await api.get('/staff/audit-logs', { params: { target_type: 'payment_webhook', size: 100 } })
    items.value = res?.items || []
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

onMounted(load)
</script>
