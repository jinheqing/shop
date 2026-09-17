<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const tickets = ref<any[]>([
  { id: 1, user_id: 1, user_email: 'james.l@london.com', request_type: 'access', status: 'pending', created_at: '2026-09-15T11:00:00Z', due_at: '2026-09-19T11:00:00Z' },
  { id: 2, user_id: 3, user_email: 'chen.xs@manchester.uk', request_type: 'erasure', status: 'processing', created_at: '2026-09-16T09:00:00Z', due_at: '2026-09-20T09:00:00Z' },
  { id: 3, user_id: 2, user_email: 'sarah.w@edinburgh.scot', request_type: 'portability', status: 'completed', created_at: '2026-09-10T14:00:00Z', completed_at: '2026-09-12T10:00:00Z', due_at: '2026-09-14T14:00:00Z' },
  { id: 4, user_id: 5, user_email: 'anonymized@example.com', request_type: 'rectification', status: 'pending', created_at: '2026-09-17T08:00:00Z', due_at: '2026-09-21T08:00:00Z' },
])

async function exportData(id: number) { await api.post(`/dsar/requests/${id}/export`); ElMessage.success('Data exported and queued for email') }
async function erase(id: number) { await api.post(`/dsar/requests/${id}/delete`); ElMessage.success('User data erased'); tickets.value.find(t=>t.id===id).status='completed' }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🔒 GDPR / DSAR Tickets</span>
      <div class="flex gap-2">
        <el-tag type="warning">{{ tickets.filter(t=>t.status==='pending'||t.status==='processing').length }} Open</el-tag>
        <el-button @click="">📁 Export All Open</el-button>
      </div>
    </div></template>
    <el-alert type="info" :closable="false" class="mb-4">
      GDPR requires response within 30 days. "Due At" column shows deadline. Erasure requests block future marketing.
    </el-alert>
    <el-table :data="tickets" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column prop="user_email" label="User Email" width="220" />
      <el-table-column prop="request_type" label="Type" width="130">
        <template #default="{ row }">
          <el-tag :type="{'access':'primary','erasure':'danger','rectification':'warning','portability':'info','restriction':'warning'}[row.request_type]" effect="dark">{{ row.request_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="Status" width="130">
        <template #default="{ row }">
          <el-tag :type="{'pending':'warning','processing':'primary','completed':'success','overdue':'danger'}[row.status]" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="Requested" width="170" />
      <el-table-column prop="due_at" label="Due By" width="170">
        <template #default="{ row }">
          <span :class="new Date(row.due_at) < new Date() ? 'text-red-600 font-medium' : ''">{{ row.due_at }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="220">
        <template #default="{ row }">
          <el-button size="small" @click="exportData(row.id)">📤 Export Data</el-button>
          <el-button v-if="row.request_type==='erasure'" size="small" type="danger" @click="erase(row.id)">🗑 Erase</el-button>
          <el-button size="small">Reply</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
