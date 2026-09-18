<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Customer Users</span>
          <el-button :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="Name" width="160" />
        <el-table-column prop="email" label="Email" width="220" />
        <el-table-column prop="country" label="Country" width="130" />
        <el-table-column prop="is_email_verified" label="Verified" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_email_verified" type="success">Yes</el-tag>
            <el-tag v-else type="warning">No</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Signed Up" width="180" />
        <el-table-column label="Actions" width="140">
          <template #default="{ row }">
            <el-button size="small" @click="viewOrders(row)">Orders</el-button>
            <el-button size="small" type="danger" @click="exportData(row)">Export</el-button>
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

const items = ref<any[]>([])

async function load() {
  const res: any = await api.get('/users')
  items.value = res?.items || res || []
}

function viewOrders(row: any) {
  window.open(`/#/orders?user_id=${row.id}`, '_blank')
}

async function exportData(row: any) {
  ElMessage.info(`Would trigger DSAR export for user ${row.id}`)
}

onMounted(load)
</script>
