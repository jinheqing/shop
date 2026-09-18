<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Payment Transactions</span>
          <div class="flex gap-2">
            <el-select v-model="filters.status" placeholder="Status" clearable style="width: 120px">
              <el-option label="Pending" value="pending" />
              <el-option label="Success" value="success" />
              <el-option label="Failed" value="failed" />
              <el-option label="Refunded" value="refunded" />
            </el-select>
            <el-select v-model="filters.gateway" placeholder="Gateway" clearable style="width: 140px">
              <el-option label="2Checkout" value="2checkout" />
              <el-option label="PayPal" value="paypal" />
            </el-select>
            <el-button :icon="Search" @click="load">Filter</el-button>
          </div>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_id" label="Order" width="100" />
        <el-table-column prop="payment_gateway" label="Gateway" width="120" />
        <el-table-column prop="gateway_transaction_id" label="Gateway ID" width="200" show-overflow-tooltip />
        <el-table-column label="Amount" width="130">
          <template #default="{ row }">
            {{ (row.currency || 'GBP') }} {{ Number(row.amount).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="statusColor(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created" width="180" />
      </el-table>
      <div class="mt-2 text-xs text-gray-500">Total: {{ total }}</div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { api } from '../api/client'

const items = ref<any[]>([])
const total = ref(0)
const filters = ref({ status: '', gateway: '' })

async function load() {
  try {
    const params: any = {}
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.gateway) params.gateway = filters.value.gateway
    const res: any = await api.get('/payment/transactions', { params })
    items.value = res?.items || res || []
    total.value = res?.total || items.value.length
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function statusColor(s: string) {
  const m: Record<string, string> = { pending: 'warning', success: 'success', failed: 'danger', refunded: 'info' }
  return m[s] || 'info'
}

onMounted(load)
</script>
