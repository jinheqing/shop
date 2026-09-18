<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Customer Requests — Pending</span>
          <el-button :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-alert v-if="!loaded" title="Click Reload to fetch from backend" type="info" show-icon :closable="false" class="mb-4" />
      <el-empty v-else-if="!items.length" description="No pending customer requests" />
      <el-table v-else :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="Request / Room" min-width="220">
          <template #default="{ row }">
            <div>{{ row.title || '(no title)' }}</div>
            <div class="text-xs text-gray-400">{{ row.description?.slice(0, 100) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="Location" width="160" />
        <el-table-column prop="status" label="Status" width="140">
          <template #default="{ row }">
            <el-tag :type="row.status === 'customer_request' ? 'warning' : 'success'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="240">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="approve(row)">Approve</el-button>
            <el-button size="small" @click="reject(row)">Reject</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Approve Dialog -->
    <el-dialog v-model="dialog" title="Approve Customer Request" width="480px">
      <div class="mb-2">Room ID: {{ editing?.id }}</div>
      <el-form :model="form">
        <el-form-item label="Scheduled Start (RFC3339)">
          <el-input v-model="form.scheduled_start" placeholder="2026-10-01T14:00:00Z" />
        </el-form-item>
        <el-form-item label="Note">
          <el-input v-model="form.note" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">Cancel</el-button>
        <el-button type="primary" @click="submitApprove">Confirm Approve</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { api } from '../api/client'

const items = ref<any[]>([])
const loaded = ref(false)
const dialog = ref(false)
const editing = ref<any>(null)
const form = ref({ scheduled_start: '', note: '' })

async function load() {
  try {
    const res: any = await api.get('/live-rooms/customer-requests')
    items.value = res?.items || res || []
    loaded.value = true
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function approve(row: any) {
  editing.value = row
  form.value = { scheduled_start: '', note: '' }
  dialog.value = true
}

function reject(row: any) {
  ElMessage.info(`Reject logic not yet wired — would change room ${row.id} status to rejected`)
}

async function submitApprove() {
  try {
    await api.post(`/live-rooms/customer-requests/${editing.value.id}/approve`, form.value)
    ElMessage.success('Approved')
    dialog.value = false
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Approve failed')
  }
}

onMounted(load)
</script>
