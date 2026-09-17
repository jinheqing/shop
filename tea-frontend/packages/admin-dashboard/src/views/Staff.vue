<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([
  { id: 1, name: 'Admin Root', email: 'admin@ukteahouse.co.uk', role: 'admin', mfa_enabled: true, is_active: true, last_login_at: '2026-09-17T13:20:00Z' },
  { id: 5, name: '顾问小王', email: 'advisor.w@ukteahouse.co.uk', role: 'advisor', mfa_enabled: false, is_active: true, last_login_at: '2026-09-17T10:02:00Z' },
  { id: 6, name: '顾问小李', email: 'advisor.li@ukteahouse.co.uk', role: 'advisor', mfa_enabled: false, is_active: true, last_login_at: '2026-09-16T15:40:00Z' },
  { id: 10, name: '王师傅', email: 'farmer.wang@ukteahouse.co.uk', role: 'tea_farmer', mfa_enabled: false, is_active: true, assigned_farm_id: 1, last_login_at: '2026-09-17T06:10:00Z' },
  { id: 20, name: 'Operations Jane', email: 'ops.jane@ukteahouse.co.uk', role: 'operations', mfa_enabled: true, is_active: true, last_login_at: '2026-09-17T09:00:00Z' },
])
const open = ref(false)
const form = reactive({ name: '', email: '', role: 'advisor', mfa_enabled: false })

async function save() { await api.post('/staff', form); ElMessage.success('Staff created'); open.value = false }
async function toggle(id: number) { await api.post(`/staff/${id}/toggle`); ElMessage.success('Toggled') }
async function resetMfa(id: number) { ElMessage.success('MFA reset token sent to staff email') }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🧑‍💼 Staff Management</span>
      <el-button type="primary" @click="open=true">+ Add Staff</el-button>
    </div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="name" label="Name" width="160" />
      <el-table-column prop="email" label="Email" width="200" />
      <el-table-column prop="role" label="Role" width="130">
        <template #default="{ row }">
          <el-tag :type="{'admin':'danger','supervisor':'warning','advisor':'primary','tea_farmer':'success','operations':'info'}[row.role]" effect="dark">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="mfa_enabled" label="MFA" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.mfa_enabled" type="success" size="small">✅ ON</el-tag>
          <el-tag v-else type="info" size="small">OFF</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_active" label="Active" width="100">
        <template #default="{ row }">
          <el-switch :model-value="row.is_active" @change="toggle(row.id)" />
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="Last Login" width="170" />
      <el-table-column label="Actions" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="">Edit</el-button>
          <el-button size="small" type="warning" @click="resetMfa(row.id)">Reset MFA</el-button>
          <el-button size="small" type="danger" @click="">Password Reset</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="Add New Staff" width="480px">
    <el-form :model="form" label-width="100px">
      <el-form-item label="Name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="Email"><el-input v-model="form.email" /></el-form-item>
      <el-form-item label="Role">
        <el-select v-model="form.role" class="w-full">
          <el-option value="admin">Admin</el-option>
          <el-option value="supervisor">Supervisor</el-option>
          <el-option value="advisor">Advisor</el-option>
          <el-option value="tea_farmer">Tea Farmer</el-option>
          <el-option value="operations">Operations</el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="Force MFA">
        <el-switch v-model="form.mfa_enabled" />
      </el-form-item>
    </el-form>
    <template #footer><el-button @click="open=false">Cancel</el-button><el-button type="primary" @click="save">Create & Send Invite</el-button></template>
  </el-dialog>
</template>
