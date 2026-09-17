<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const inspections = ref<any[]>([
  { id: 1, live_room_id: 7, order_id: 9, order_no: 'ORD-20260917-398535', user_name: 'James Lovelace', product_title: '邦东古树饼', scheduled_start: '2026-09-22T14:00:00Z', status: 'scheduled', livekit_token_for_host: 'eyJhbGciOi...', livekit_token_for_viewers: 'eyJhbGciOi...' },
  { id: 2, live_room_id: 12, order_id: 11, order_no: 'ORD-20260918-112233', user_name: '陳曉珊', product_title: '凤凰山大乌岽单丛', scheduled_start: '2026-09-30T10:00:00Z', status: 'scheduled' },
  { id: 3, live_room_id: 5, order_id: 8, order_no: 'ORD-20260916-112233', user_name: 'Sarah Whitfield', product_title: '布朗山古树普洱', scheduled_start: '2026-09-18T09:00:00Z', status: 'completed', recording_url: 'https://cdn.ourdomain.com/recordings/room5.mp4' },
])

async function autoCreate(orderId: number) {
  ElMessage.success('🔄 Delivery inspection room auto-created for order ' + orderId + ' — LiveKit room + LiveKit tokens generated')
}
</script>
<template>
  <div class="space-y-4">
    <el-card>
      <template #header><div class="flex justify-between items-center">
        <span class="font-medium">📦 Delivery Inspection Live Rooms</span>
        <el-button type="primary" @click="autoCreate(9)">🔄 System: Auto-Create For Order</el-button>
      </div></template>
      <el-alert type="info" :closable="false" class="mb-4">
        Automatically created by the order state machine when order transitions to <code>ready_for_delivery</code>. Bind to order → customer gets live link in email.
      </el-alert>
      <el-table :data="inspections" stripe>
        <el-table-column prop="order_no" label="Order" width="170" />
        <el-table-column prop="user_name" label="Customer" width="180" />
        <el-table-column prop="product_title" label="Product" width="220" />
        <el-table-column prop="scheduled_start" label="Scheduled Start" width="180" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="{ row }">
            <el-tag :type="{'scheduled':'warning','live':'success','completed':'primary','expired':'info'}[row.status]" effect="dark">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="220">
          <template #default="{ row }">
            <el-button v-if="row.status==='scheduled'" size="small" type="success">▶ Go Live (Host)</el-button>
            <el-button size="small">View Tokens</el-button>
            <el-button v-if="row.recording_url" size="small" type="info">🎬 Replay</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
