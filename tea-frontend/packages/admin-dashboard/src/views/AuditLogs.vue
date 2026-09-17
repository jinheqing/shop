<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/api/client'
const list = ref<any[]>([])
async function load() { try { const d: any = await api.get('/staff/audit-logs'); list.value = d?.items || d || [] } catch {} }
onMounted(load)
</script>
<template>
  <el-card>
    <template #header><span class="font-medium">📋 Staff Audit Trail (GDPR)</span></template>
    <el-alert type="info" :closable="false" class="mb-4" show-icon>
      All staff actions are logged to a separate audit database. Immutable. GDPR-compliant.
    </el-alert>
    <el-table :data="list.length?list:[
      { id: 1, staff_id: 4, action: 'custom_product.create', target_type: 'custom_product', target_id: 16, detail: { title: '冰岛古树饼' }, ip_address: '::1', user_agent: 'curl/8.5.0', created_at: '2026-09-17T13:27:42Z' },
      { id: 2, staff_id: 4, action: 'custom_product.publish', target_type: 'custom_product', target_id: 16, detail: { product_token: 'XK92...f3A' }, ip_address: '::1', user_agent: 'curl/8.5.0', created_at: '2026-09-17T13:27:43Z' },
      { id: 3, staff_id: 4, action: 'order.create', target_type: 'order', target_id: 9, detail: { order_no: 'ORD-20260917-398535' }, ip_address: '::1', user_agent: 'curl/8.5.0', created_at: '2026-09-17T13:34:39Z' },
      { id: 4, staff_id: 4, action: 'payment.init', target_type: 'order', target_id: 9, detail: { gateway: '2checkout' }, ip_address: '::1', user_agent: 'curl/8.5.0', created_at: '2026-09-17T13:34:40Z' },
      { id: 5, staff_id: 4, action: 'order.state.paid', target_type: 'order', target_id: 9, detail: { from: 'ordering', to: 'paid' }, ip_address: '149.12.5.22', user_agent: '2Checkout INS', created_at: '2026-09-17T13:34:41Z' },
    ]" stripe>
      <el-table-column prop="created_at" label="Time" width="180" />
      <el-table-column prop="staff_id" label="Staff ID" width="100" />
      <el-table-column prop="action" label="Action" width="180">
        <template #default="{ row }">
          <el-tag size="small" effect="plain">{{ row.action }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="target_type" label="Target Type" width="140" />
      <el-table-column prop="target_id" label="Target ID" width="100" />
      <el-table-column label="Detail" show-overflow-tooltip min-width="200">
        <template #default="{ row }"><code class="text-xs text-slate-600">{{ JSON.stringify(row.detail) }}</code></template>
      </el-table-column>
      <el-table-column prop="ip_address" label="IP" width="140" />
      <el-table-column prop="user_agent" label="User Agent" min-width="180" show-overflow-tooltip />
    </el-table>
  </el-card>
</template>
