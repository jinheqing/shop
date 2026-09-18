<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const filter = reactive({ visibility: '', live_room_id: '' })
const openVis = ref(false)
const editing = ref<any>(null)
const groups = ref<any[]>([])

const VIS_OPTIONS = [
  { value: 'public', label: 'Public (all visitors)' },
  { value: 'registered', label: 'Registered users' },
  { value: 'restricted', label: 'Restricted (selected users/groups)' },
]

async function load() {
  try { list.value = await api.get('/recordings') } catch { list.value = [] }
  try { groups.value = await api.get('/user-groups') } catch { groups.value = [] }
}

function openVisEdit(row: any) {
  editing.value = reactive({
    id: row.id,
    visibility: row.visibility || 'registered',
    visible_user_ids: row.visible_user_ids || [],
    visible_group_ids: row.visible_group_ids || [],
  })
  openVis.value = true
}
async function saveVis() {
  await api.put(`/recordings/${editing.value.id}/visibility`, {
    visibility: editing.value.visibility,
    visible_user_ids: editing.value.visible_user_ids,
    visible_group_ids: editing.value.visible_group_ids,
  })
  ElMessage.success('Visibility updated')
  openVis.value = false
  load()
}
async function del(row: any) {
  if (!confirm(`Delete recording "${row.title || row.id}"?`)) return
  await api.delete(`/recordings/${row.id}`)
  ElMessage.success('Deleted')
  load()
}

function resolveUrl(url: string) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  return base + url
}

function fmtSize(bytes: number) {
  if (!bytes) return '—'
  if (bytes > 1024*1024) return (bytes/1024/1024).toFixed(1) + ' MB'
  return bytes + ' B'
}

onMounted(load)
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">Recordings (独立管理)</span>
        <div class="flex gap-2">
          <el-select v-model="filter.visibility" placeholder="Visibility" clearable size="small" style="width:160px">
            <el-option v-for="o in VIS_OPTIONS" :key="o.value" :value="o.value" :label="o.label" />
          </el-select>
          <el-button size="small" @click="load">Refresh</el-button>
        </div>
      </div>
    </template>

    <el-alert type="info" show-icon :closable="false" class="mb-3">
      一对一：一个直播停止后自动生成一条 Recording。Restricted 级别的回放只有指定用户 / 用户组能看。
    </el-alert>

    <el-table :data="list" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column label="Title" min-width="200">
        <template #default="{ row }">
          <div>{{ row.title || '(from LiveRoom #' + row.live_room_id + ')' }}</div>
          <a :href="resolveUrl(row.file_url)" target="_blank" class="text-xs text-slate-400 hover:text-blue-500">▶ Play</a>
        </template>
      </el-table-column>
      <el-table-column prop="duration_sec" label="Duration" width="100">
        <template #default="{ row }">{{ row.duration_sec ? Math.floor(row.duration_sec/60)+'m' : '—' }}</template>
      </el-table-column>
      <el-table-column label="Size" width="120">
        <template #default="{ row }">{{ fmtSize(row.file_size) }}</template>
      </el-table-column>
      <el-table-column label="Visibility" width="160">
        <template #default="{ row }">
          <el-tag :type="({public:'success', registered:'', restricted:'danger'} as any)[row.visibility]" effect="dark" size="small">{{ row.visibility }}</el-tag>
          <div v-if="row.visibility === 'restricted'" class="text-[10px] text-slate-400 mt-1">
            {{ (row.visible_user_ids || []).length }} users · {{ (row.visible_group_ids || []).length }} groups
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Created" width="170">
        <template #default="{ row }">{{ row.created_at?.slice(0,16).replace('T',' ') }}</template>
      </el-table-column>
      <el-table-column label="Actions" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openVisEdit(row)">👁 Visibility</el-button>
          <el-button size="small" type="danger" @click="del(row)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- Visibility 编辑 -->
  <el-dialog v-model="openVis" :title="`Visibility for Recording #${editing?.id}`" width="520px">
    <el-form :model="editing" label-width="160px">
      <el-form-item label="Visibility">
        <el-select v-model="editing.visibility" class="w-full">
          <el-option v-for="o in VIS_OPTIONS" :key="o.value" :value="o.value" :label="o.label" />
        </el-select>
      </el-form-item>
      <template v-if="editing?.visibility === 'restricted'">
        <el-form-item label="Visible User IDs">
          <el-input v-model="editing.visible_user_ids" placeholder="[1, 42, 108]" />
          <div class="text-[11px] text-slate-400">数组格式，也可以 JSON array 粘贴</div>
        </el-form-item>
        <el-form-item label="Visible Groups">
          <el-select v-model="editing.visible_group_ids" multiple class="w-full">
            <el-option v-for="g in groups" :key="g.id" :value="g.id" :label="g.name" />
          </el-select>
        </el-form-item>
      </template>
    </el-form>
    <template #footer>
      <el-button @click="openVis=false">Cancel</el-button>
      <el-button type="primary" @click="saveVis">Save</el-button>
    </template>
  </el-dialog>
</template>
