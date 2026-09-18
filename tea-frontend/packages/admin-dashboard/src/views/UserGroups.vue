<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const editing = ref<any>(null)
const showMembers = ref<any>(null)
const members = ref<any[]>([])

const form = reactive({
  name: '', description: '',
  privileges: [] as string[],
  reciprocal_str: '',
})

const PRIVILEGE_OPTIONS = [
  { value: 'preorder_priority', label: '优先锁货权 (稀缺茶)' },
  { value: 'private_masterclass', label: '闭门品鉴 & Masterclass' },
  { value: 'vip_concierge', label: '专属顾问直连' },
  { value: 'deep_traceability', label: '深度溯源 (完整茶园/茶农/SGS)' },
  { value: 'offline_invite', label: '年度私宴邀请' },
]

async function load() {
  try { list.value = await api.get('/user-groups') } catch { list.value = [] }
}

function openNew() {
  editing.value = null
  Object.assign(form, { name: '', description: '', privileges: [], reciprocal_str: '' })
  open.value = true
}
function openEdit(row: any) {
  editing.value = row
  Object.assign(form, {
    name: row.name, description: row.description || '',
    privileges: row.privileges || [], reciprocal_str: (row.reciprocal || []).join(', '),
  })
  open.value = true
}

async function submit() {
  const payload = { ...form, reciprocal: form.reciprocal_str.split(',').map(s=>s.trim()).filter(Boolean) }
  if (editing.value) {
    await api.put(`/user-groups/${editing.value.id}`, payload)
    ElMessage.success('Updated')
  } else {
    await api.post('/user-groups', payload)
    ElMessage.success('Created')
  }
  open.value = false
  load()
}

async function del(row: any) {
  await ElMessageBox.confirm(`Delete group "${row.name}"? All members will be removed.`, 'Confirm')
  await api.delete(`/user-groups/${row.id}`)
  ElMessage.success('Deleted')
  load()
}

async function viewMembers(row: any) {
  showMembers.value = row
  try { members.value = await api.get(`/user-groups/${row.id}/members`) } catch { members.value = [] }
}
async function removeMember(userId: number) {
  await api.delete(`/user-groups/${showMembers.value.id}/members/${userId}`)
  viewMembers(showMembers.value)
}
async function addMember() {
  const input = prompt('Enter user ID to add:')
  if (!input) return
  const uid = parseInt(input)
  if (isNaN(uid)) { ElMessage.error('Invalid ID'); return }
  try {
    await api.post(`/user-groups/${showMembers.value.id}/members`, { user_id: uid })
    ElMessage.success('Added')
    viewMembers(showMembers.value)
  } catch (e: any) { ElMessage.error(e?.response?.data?.message || 'Failed') }
}
async function bulkAdd() {
  const input = prompt('Enter comma-separated user IDs (e.g. 1,2,3):')
  if (!input) return
  const ids = input.split(',').map(s => parseInt(s.trim())).filter(n => !isNaN(n))
  if (ids.length === 0) { ElMessage.error('No valid IDs'); return }
  await api.post(`/user-groups/${showMembers.value.id}/members/bulk`, { user_ids: ids })
  ElMessage.success(`Bulk added ${ids.length}`)
  viewMembers(showMembers.value)
}
async function autoSync() {
  await api.post('/user-groups/auto-sync')
  ElMessage.success('Auto-sync triggered')
  load()
}

onMounted(load)
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">User Groups</span>
        <div class="flex gap-2">
          <el-button size="small" @click="autoSync">🔄 Auto-Sync (spend_threshold)</el-button>
          <el-button type="primary" size="small" @click="openNew">+ New Group</el-button>
        </div>
      </div>
    </template>

    <el-alert type="info" show-icon :closable="false" class="mb-3">
      Non-price privileges — 老钱不公开等级。Privileges: preorder_priority (稀缺茶锁货) · private_masterclass (闭门品鉴) · vip_concierge (专属顾问) · deep_traceability (深度溯源) · offline_invite (年度私宴)
    </el-alert>

    <el-table :data="list" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column prop="name" label="Group" width="240" />
      <el-table-column prop="description" label="Description" show-overflow-tooltip />
      <el-table-column label="Privileges" min-width="280">
        <template #default="{ row }">
          <el-tag v-for="p in (row.privileges || [])" :key="p" size="small" type="warning" class="mr-1 mb-1">{{ p }}</el-tag>
          <span v-if="!row.privileges?.length" class="text-slate-400 text-xs">—</span>
        </template>
      </el-table-column>
      <el-table-column label="Members" width="100">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="viewMembers(row)">{{ (row.members || []).length || 'View' }}</el-button>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">Edit</el-button>
          <el-button size="small" type="danger" @click="del(row)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 新建/编辑 对话框 -->
  <el-dialog v-model="open" :title="editing ? `Edit ${editing.name}` : 'New User Group'" width="640px">
    <el-form :model="form" label-width="120px">
      <el-form-item label="Name" required>
        <el-input v-model="form.name" placeholder="e.g. VIP 2026 / 勐海古树 2026 批次买家" />
      </el-form-item>
      <el-form-item label="Description">
        <el-input v-model="form.description" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item label="Privileges">
        <el-checkbox-group v-model="form.privileges">
          <el-checkbox v-for="opt in PRIVILEGE_OPTIONS" :key="opt.value" :value="opt.value">{{ opt.label }}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="Reciprocal Clubs">
        <el-input v-model="form.reciprocal_str" placeholder="e.g. annabels, soho_house" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open=false">Cancel</el-button>
      <el-button type="primary" @click="submit">Save</el-button>
    </template>
  </el-dialog>

  <!-- 成员管理 对话框 -->
  <el-dialog v-model="showMembers" :title="`Members of ${showMembers?.name || ''}`" width="720px">
    <div class="flex gap-2 mb-3">
      <el-button size="small" @click="addMember">+ Add Member</el-button>
      <el-button size="small" @click="bulkAdd">Bulk Add (CSV)</el-button>
    </div>
    <el-table :data="members" stripe>
      <el-table-column prop="user_id" label="User ID" width="100" />
      <el-table-column label="Name / Email">
        <template #default="{ row }">
          <span>{{ row.user?.name || '(unknown)' }}</span>
          <div class="text-xs text-slate-400">{{ row.user?.email }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="added_at" label="Added At" width="180" />
      <el-table-column label="Action" width="100">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="removeMember(row.user_id)">Remove</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>
