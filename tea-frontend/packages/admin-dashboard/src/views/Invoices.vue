<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Invoices</span>
          <el-button :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_id" label="Order" width="100" />
        <el-table-column prop="invoice_number" label="Invoice #" width="180" />
        <el-table-column label="Amount" width="130">
          <template #default="{ row }">GBP {{ Number(row.total_amount || 0).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'paid' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created" width="180" />
        <el-table-column label="Actions" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="download(row)">Download PDF</el-button>
            <el-button size="small" type="warning" @click="regenerate(row)">Regenerate</el-button>
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
  try {
    const res: any = await api.get('/orders?size=100')
    const orders = res?.items || res || []
    items.value = orders.map((o: any) => ({
      id: o.id, order_id: o.id, invoice_number: `INV-${o.id}`,
      total_amount: o.total_amount || 0, status: o.payment_status || o.state,
      created_at: o.created_at
    }))
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

async function download(row: any) {
  try {
    const res = await api.get(`/orders/${row.order_id}/invoice/pdf`, { responseType: 'blob' })
    const blob = res instanceof Blob ? res : new Blob([res as any], { type: 'application/pdf' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url; a.download = `invoice-${row.order_id}.pdf`; a.click()
    URL.revokeObjectURL(url)
  } catch (e: any) {
    ElMessage.error('PDF download failed')
  }
}

async function regenerate(row: any) {
  try {
    await api.post(`/orders/${row.order_id}/invoice/regenerate`)
    ElMessage.success('Regenerated')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Regenerate failed')
  }
}

onMounted(load)
</script>
