<script setup lang="ts">
import { ref } from 'vue'

const logs = ref<any[]>([
  { id: 1, gateway: '2checkout', event_type: 'ORDER_CREATED', gateway_ref: '20389123456', order_id: 9, raw_status: 'CREATED', signature_valid: true, processed: true, created_at: '2026-09-17T13:34:39Z', response_code: 200 },
  { id: 2, gateway: '2checkout', event_type: 'ORDER_PAYMENT_STATUS_CHANGED', gateway_ref: '20389123456', order_id: 9, raw_status: 'APPROVED', signature_valid: true, processed: true, created_at: '2026-09-17T13:34:41Z', response_code: 200 },
  { id: 3, gateway: 'paypal', event_type: 'PAYMENT.CAPTURE.COMPLETED', gateway_ref: 'PAYPAL-7821-ABC', order_id: 8, raw_status: 'COMPLETED', signature_valid: true, processed: true, created_at: '2026-09-16T09:12:30Z', response_code: 200 },
  { id: 4, gateway: '2checkout', event_type: 'ORDER_PAYMENT_STATUS_CHANGED', gateway_ref: '20389000001', order_id: 7, raw_status: 'REFUND_ISSUED', signature_valid: true, processed: true, created_at: '2026-09-15T19:00:00Z', response_code: 200 },
  { id: 5, gateway: 'paypal', event_type: 'PAYMENT.CAPTURE.DECLINED', gateway_ref: 'PAYPAL-6666-DEF', order_id: 10, raw_status: 'DECLINED', signature_valid: true, processed: true, created_at: '2026-09-17T14:01:01Z', response_code: 200 },
  { id: 6, gateway: 'paypal', event_type: 'PAYMENT.CAPTURE.COMPLETED', gateway_ref: 'PAYPAL-XXXXX', order_id: null, raw_status: 'COMPLETED', signature_valid: false, processed: false, created_at: '2026-09-17T14:05:00Z', response_code: 400 },
])
const selected = ref<any>(null)
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🔔 Webhook Delivery Logs</span>
      <div class="flex items-center gap-2">
        <el-tag type="success">{{ logs.filter(l=>l.signature_valid).length }} ✅ Signed</el-tag>
        <el-tag type="danger">{{ logs.filter(l=>!l.signature_valid).length }} ❌ Invalid Sig</el-tag>
        <el-tag type="warning">{{ logs.filter(l=>!l.processed).length }} ⚠️ Unprocessed</el-tag>
        <el-button @click="">🔄 Retry All Failed</el-button>
      </div>
    </div></template>
    <el-table :data="logs" stripe @row-click="selected = $event">
      <el-table-column prop="created_at" label="Received" width="170" />
      <el-table-column prop="gateway" label="Gateway" width="100">
        <template #default="{ row }"><el-tag :type="row.gateway==='2checkout'?'primary':'warning'" size="small">{{ row.gateway }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="event_type" label="Event Type" width="240">
        <template #default="{ row }"><code class="text-xs">{{ row.event_type }}</code></template>
      </el-table-column>
      <el-table-column prop="gateway_ref" label="Gateway Ref" show-overflow-tooltip />
      <el-table-column prop="order_id" label="Order ID" width="90" />
      <el-table-column prop="signature_valid" label="Sig" width="80">
        <template #default="{ row }">
          <el-icon :class="row.signature_valid ? 'text-green-600' : 'text-red-600'"><component :is="row.signature_valid ? 'Check' : 'Close'" /></el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="processed" label="Processed" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.processed" type="success" size="small">✅ Yes</el-tag>
          <el-tag v-else type="danger" size="small">❌ No</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="response_code" label="HTTP" width="80" />
      <el-table-column label="Action" width="80">
        <template #default="{ row }"><el-button size="small" @click="selected = row">Raw</el-button></template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-drawer v-model="!!selected" title="Webhook Payload" size="600px">
    <pre v-if="selected" class="bg-slate-900 text-green-300 p-4 rounded text-xs overflow-auto">{{ JSON.stringify(selected, null, 2) }}</pre>
  </el-drawer>
</template>
