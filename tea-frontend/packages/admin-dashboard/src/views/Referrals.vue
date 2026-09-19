<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const total = ref(0)
const loading = ref(false)
const open = ref(false)
const editing = ref<any>(null)

const filters = reactive({
  referrer_name: '',
  referred_name: '',
  reward_only: false,
  reward_pending: false,
})

const form = reactive({
  referred_user_id: 0,
  referrer_user_id: null as number | null,
  referrer_name: '',
  source: 'manual',
  notes: '',
})

async function load() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.referrer_name) params.referrer_name = filters.referrer_name
    if (filters.referred_name) params.referred_name = filters.referred_name
    if (filters.reward_only) params.reward_only = true
    if (filters.reward_pending) params.reward_pending = true
    const data = await api.get('/referrals', { params })
    list.value = data.items || data
    total.value = data.total || list.value.length
  } catch (e) {
    ElMessage.error('Failed to load referrals')
  } finally {
    loading.value = false
  }
}

function openNew() {
  editing.value = null
  Object.assign(form, { referred_user_id: 0, referrer_user_id: null, referrer_name: '', source: 'manual', notes: '' })
  open.value = true
}

async function submit() {
  if (!form.referred_user_id || !form.referrer_name) {
    ElMessage.warning('Referred user ID and referrer name are required')
    return
  }
  const payload = { ...form }
  if (editing.value) {
    // No update endpoint for now — delete and recreate
    ElMessage.info('Delete existing and create new for edits')
    return
  }
  try {
    await api.post('/referrals', payload)
    ElMessage.success('Referral created')
    open.value = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || 'Failed to create')
  }
}

async function triggerReward(row: any) {
  await ElMessageBox.confirm(
    `Mark referral reward as triggered? (Referrer: ${row.referrer_name}, Friend order amount: £${row.friend_order_amount?.toFixed(2) || '0.00'})`,
    'Confirm Reward',
    { type: 'info' }
  )
  await api.post(`/referrals/${row.id}/trigger-reward`)
  ElMessage.success('Reward triggered')
  load()
}

async function del(row: any) {
  await ElMessageBox.confirm('Delete this referral relationship?', 'Confirm', { type: 'warning' })
  await api.delete(`/referrals/${row.id}`)
  ElMessage.success('Deleted')
  load()
}

function sourceBadge(s: string) {
  const map: Record<string, { type: string; label: string }> = {
    name_share: { type: 'primary', label: 'Name Share' },
    short_code: { type: 'success', label: 'Short Code' },
    manual: { type: 'warning', label: 'Manual' },
  }
  return map[s] || { type: 'info', label: s }
}

onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-xl font-medium mb-4">Referrals Management</h1>

    <!-- Stats summary -->
    <div class="grid grid-cols-4 gap-4 mb-4">
      <div class="bg-white rounded p-4 border">
        <div class="text-xs text-slate-500">Total Referrals</div>
        <div class="text-2xl font-bold">{{ total }}</div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="text-xs text-slate-500">Rewards Triggered</div>
        <div class="text-2xl font-bold text-green-600">
          {{ list.filter(r => r.reward_triggered_at).length }}
        </div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="text-xs text-slate-500">Pending Rewards</div>
        <div class="text-2xl font-bold text-orange-500">
          {{ list.filter(r => r.friend_order_amount > 0 && !r.reward_triggered_at).length }}
        </div>
      </div>
      <div class="bg-white rounded p-4 border">
        <div class="text-xs text-slate-500">No Order Yet</div>
        <div class="text-2xl font-bold text-slate-400">
          {{ list.filter(r => !r.friend_order_amount || r.friend_order_amount === 0).length }}
        </div>
      </div>
    </div>

    <!-- Filters -->
    <el-card class="mb-4">
      <div class="flex gap-3 items-center flex-wrap">
        <el-input v-model="filters.referrer_name" placeholder="Referrer name..." clearable style="width: 200px" />
        <el-input v-model="filters.referred_name" placeholder="Referred name..." clearable style="width: 200px" />
        <el-checkbox v-model="filters.reward_only">Only rewarded</el-checkbox>
        <el-checkbox v-model="filters.reward_pending">Pending reward</el-checkbox>
        <el-button type="primary" @click="load">Search</el-button>
        <el-button @click="Object.assign(filters, { referrer_name: '', referred_name: '', reward_only: false, reward_pending: false }); load()">Reset</el-button>
        <el-button type="success" @click="openNew" class="ml-auto">+ Add Referral</el-button>
      </div>
    </el-card>

    <!-- Table -->
    <el-card>
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="Referrer">
          <template #default="{ row }">
            <div>
              <div class="font-medium">{{ row.referrer_name }}</div>
              <div class="text-xs text-slate-400" v-if="row.referrer_user_id">User #{{ row.referrer_user_id }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Referred Friend">
          <template #default="{ row }">
            <div>
              <div class="font-medium">{{ row.referred_name }}</div>
              <div class="text-xs text-slate-400">User #{{ row.referred_user_id }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Source" width="120">
          <template #default="{ row }">
            <el-tag :type="sourceBadge(row.source).type">{{ sourceBadge(row.source).label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Friend Order" width="150">
          <template #default="{ row }">
            <template v-if="row.friend_order_amount > 0">
              <div class="font-medium">£{{ row.friend_order_amount.toFixed(2) }}</div>
              <div class="text-xs text-slate-400">Order #{{ row.friend_order_id }}</div>
            </template>
            <span v-else class="text-slate-400 text-sm">No order yet</span>
          </template>
        </el-table-column>
        <el-table-column label="Reward" width="150">
          <template #default="{ row }">
            <template v-if="row.reward_triggered_at">
              <el-tag type="success" size="small">✓ Triggered</el-tag>
              <div class="text-xs text-slate-400 mt-1">{{ new Date(row.reward_triggered_at).toLocaleDateString() }}</div>
            </template>
            <template v-else-if="row.friend_order_amount > 0">
              <el-tag type="warning" size="small">Pending</el-tag>
              <el-button size="small" type="success" link @click="triggerReward(row)">Trigger Now</el-button>
            </template>
            <span v-else class="text-slate-400 text-sm">Waiting...</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created" width="170">
          <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
        <el-table-column label="Actions" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" link @click="del(row)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create Dialog -->
    <el-dialog v-model="open" title="Create Referral (Manual)" width="500px">
      <el-form label-width="120px">
        <el-form-item label="Referred User ID">
          <el-input-number v-model="form.referred_user_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Referrer User ID (optional)">
          <el-input-number v-model="form.referrer_user_id" :min="1" :allow-null="true" style="width: 100%" placeholder="If referrer exists in our system" />
        </el-form-item>
        <el-form-item label="Referrer Name">
          <el-input v-model="form.referrer_name" placeholder="Name of the person who recommended" />
        </el-form-item>
        <el-form-item label="Source">
          <el-select v-model="form.source" style="width: 100%">
            <el-option label="Manual (staff)" value="manual" />
            <el-option label="Name Share (auto)" value="name_share" />
            <el-option label="Short Code (auto)" value="short_code" />
          </el-select>
        </el-form-item>
        <el-form-item label="Notes">
          <el-input v-model="form.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="open = false">Cancel</el-button>
        <el-button type="primary" @click="submit">Create</el-button>
      </template>
    </el-dialog>
  </div>
</template>
