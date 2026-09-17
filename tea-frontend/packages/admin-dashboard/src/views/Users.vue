<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const search = ref('')

async function load() {
  try { const d: any = await api.get('/users'); list.value = d?.items || d || [] }
  catch {
    list.value = [
      { id: 1, name: 'James Lovelace', email: 'james.l@london.com', phone: '+447700123456', preferred_language: 'en', preferred_timezone: 'Europe/London', orders_count: 3, total_spent: 652.00, preferred_advisor_id: 5, last_login_at: '2026-09-17T10:20:00Z', created_at: '2026-05-01T09:00:00Z' },
      { id: 2, name: 'Sarah Whitfield', email: 'sarah.w@edinburgh.scot', phone: '+441315550123', preferred_language: 'en', preferred_timezone: 'Europe/London', orders_count: 7, total_spent: 1890.00, preferred_advisor_id: 5, last_login_at: '2026-09-15T16:45:00Z', created_at: '2026-03-20T11:30:00Z' },
      { id: 3, name: '陳曉珊', email: 'chen.xs@manchester.uk', phone: '+441615559876', preferred_language: 'zh', preferred_timezone: 'Europe/London', orders_count: 1, total_spent: 188.00, preferred_advisor_id: 6, last_login_at: '2026-09-10T08:15:00Z', created_at: '2026-07-12T14:00:00Z' },
    ]
  }
}
load()

async function dsar(id: number) {
  await ElMessageBox.confirm('Send DSAR data export request for this user?', 'GDPR', { type: 'warning' })
  await api.post(`/dsar/requests`, { user_id: id, request_type: 'access' })
  ElMessage.success('DSAR request queued — will export to user email within 48h')
}
async function del(id: number) {
  await ElMessageBox.confirm('Soft-delete this user? (GDPR erasure)', 'Warning', { type: 'warning' })
  await api.post(`/users/${id}/dsar-delete`)
  ElMessage.success('Erasure request sent')
  load()
}

const filtered = () => list.value.filter(u => !search.value || u.name.includes(search.value) || u.email.includes(search.value))
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">👥 Customer Management ({{ list.length }})</span>
      <div class="flex gap-2">
        <el-input v-model="search" placeholder="Search by name/email" style="width:220px" clearable />
        <el-button type="primary">+ New User</el-button>
      </div>
    </div></template>
    <el-table :data="filtered()" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column prop="name" label="Name" width="180" />
      <el-table-column prop="email" label="Email" width="220" />
      <el-table-column prop="phone" label="Phone" width="160" />
      <el-table-column label="Orders" width="110">
        <template #default="{ row }">
          <div class="font-medium">{{ row.orders_count }}</div>
          <div class="text-xs text-slate-500">£{{ row.total_spent?.toFixed(2) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="preferred_language" label="Lang" width="80" />
      <el-table-column prop="last_login_at" label="Last Login" width="170" />
      <el-table-column label="Actions" width="260">
        <template #default="{ row }">
          <el-button size="small">Orders</el-button>
          <el-button size="small" type="success">💬 IM</el-button>
          <el-button size="small" type="warning" @click="dsar(row.id)">DSAR</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Erase</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
