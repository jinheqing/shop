<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const requests = ref<any[]>([
  { id: 1, user_id: 1, user_email: 'james.l@london.com', user_name: 'James', request_type: 'custom_live', preferred_date: '2026-09-25', custom_product_id: 16, description: '希望顾问能带我看一下冰岛茶园今秋的采摘情况', status: 'pending', created_at: '2026-09-17T10:30:00Z' },
  { id: 2, user_id: 2, user_email: 'sarah.w@edinburgh.scot', user_name: 'Sarah', request_type: 'tasting', preferred_date: '2026-09-28', custom_product_id: null, description: '想要顾问帮我对比 2024 vs 2025 春茶', status: 'approved', assigned_staff_id: 5, created_at: '2026-09-16T15:00:00Z' },
  { id: 3, user_id: 3, user_email: 'chen.xs@manchester.uk', user_name: '陳曉珊', request_type: 'delivery_inspection', preferred_date: '2026-09-30', custom_product_id: 18, description: '收货前想看验货直播', status: 'scheduled', live_room_id: 12, created_at: '2026-09-17T09:00:00Z' },
])

async function approve(id: number) { await api.post(`/live-rooms/customer-requests/${id}/approve`); ElMessage.success('Approved & scheduled'); requests.value.find(r=>r.id===id).status='approved' }
async function reject(id: number) { ElMessage.info('Rejected'); requests.value.find(r=>r.id===id).status='rejected' }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🎥 Customer Live Request Queue</span>
      <el-tag type="warning">{{ requests.filter(r=>r.status==='pending').length }} pending</el-tag>
    </div></template>
    <el-alert type="info" :closable="false" class="mb-4 text-sm">
      Customers submit these from the Bespoke/Quote page. Each approved auto-creates a live room & sends calendar invite.
    </el-alert>
    <el-table :data="requests" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column label="Customer" width="220">
        <template #default="{ row }"><div class="font-medium">{{ row.user_name }}</div><div class="text-xs text-slate-500">{{ row.user_email }}</div></template>
      </el-table-column>
      <el-table-column prop="request_type" label="Type" width="180">
        <template #default="{ row }">
          <el-tag :type="{'custom_live':'primary','tasting':'success','delivery_inspection':'danger'}[row.request_type]" effect="dark">{{ row.request_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="preferred_date" label="Preferred Date" width="140" />
      <el-table-column prop="description" label="Description" show-overflow-tooltip min-width="250" />
      <el-table-column prop="status" label="Status" width="110">
        <template #default="{ row }">
          <el-tag :type="{'pending':'warning','approved':'primary','scheduled':'success','rejected':'danger'}[row.status]" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="200">
        <template #default="{ row }">
          <el-button v-if="row.status==='pending'" size="small" type="success" @click="approve(row.id)">✅ Approve & Schedule</el-button>
          <el-button v-if="row.status==='pending'" size="small" type="danger" @click="reject(row.id)">❌ Reject</el-button>
          <el-button v-if="row.status==='scheduled'" size="small" type="warning">Go Live</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
