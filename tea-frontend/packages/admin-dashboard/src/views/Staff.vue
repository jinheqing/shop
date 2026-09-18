<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">Staff Management</span>
          <el-button type="primary" @click="openCreate">+ Invite New Staff</el-button>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="Name" width="160" />
        <el-table-column prop="email" label="Email" width="220" />
        <el-table-column prop="role" label="Role" width="130">
          <template #default="{ row }">
            <el-tag size="small" :type="roleColor(row.role)">{{ row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'">{{ row.is_active ? 'Active' : 'Disabled' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="Last Login" width="180" />
        <el-table-column label="Actions" width="200">
          <template #default="{ row }">
            <el-button size="small" v-if="row.is_active" type="warning" @click="toggle(row, false)">Disable</el-button>
            <el-button size="small" v-else type="success" @click="toggle(row, true)">Enable</el-button>
            <el-button size="small" type="danger" @click="remove(row)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialog" title="Invite Staff" width="480px">
      <el-form :model="form">
        <el-form-item label="Name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="Email"><el-input v-model="form.email" /></el-form-item>
        <el-form-item label="Role">
          <el-select v-model="form.role" style="width:100%">
            <el-option label="Admin" value="admin" />
            <el-option label="Supervisor" value="supervisor" />
            <el-option label="Advisor" value="advisor" />
            <el-option label="Tea Farmer" value="tea_farmer" />
            <el-option label="Operations" value="operations" />
          </el-select>
        </el-form-item>
        <el-form-item label="Timezone">
          <el-input v-model="form.work_timezone" placeholder="Europe/London" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">Cancel</el-button>
        <el-button type="primary" @click="submitCreate">Create & Send Invite</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

const items = ref<any[]>([])
const dialog = ref(false)
const form = ref({ name: '', email: '', role: 'advisor', work_timezone: 'Europe/London' })

async function load() {
  const res: any = await api.get('/staff')
  items.value = res?.items || res || []
}

function roleColor(r: string) {
  const m: Record<string, string> = { admin: 'danger', supervisor: 'warning', advisor: 'primary', tea_farmer: 'success', operations: 'info' }
  return m[r] || ''
}

function openCreate() {
  form.value = { name: '', email: '', role: 'advisor', work_timezone: 'Europe/London' }
  dialog.value = true
}

async function submitCreate() {
  try {
    await api.post('/staff', form.value)
    ElMessage.success('Staff invited')
    dialog.value = false
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Create failed')
  }
}

async function toggle(row: any, active: boolean) {
  try {
    await api.post(`/staff/${row.id}/toggle`, { active })
    ElMessage.success('Updated')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Toggle failed')
  }
}

async function remove(row: any) {
  if (!confirm(`Delete ${row.email}? (soft-delete)`)) return
  try {
    await api.delete(`/staff/${row.id}`)
    ElMessage.success('Soft-deleted')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Delete failed')
  }
}

onMounted(load)
</script>
