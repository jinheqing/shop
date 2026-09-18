<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])

async function load() {
  try { list.value = await api.get('/short-links') } catch { list.value = [] }
}

async function copyCode(row: any) {
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  const full = `${base}/s/${row.code}`
  await navigator.clipboard.writeText(full)
  ElMessage.success(`Copied: ${full}`)
}

async function del(row: any) {
  if (!confirm(`Delete short link ${row.code}?`)) return
  await api.delete(`/short-links/${row.id}`)
  ElMessage.success('Deleted')
  load()
}

async function create() {
  const url = prompt('Target URL (relative or absolute):')
  if (!url) return
  try {
    
  const r: any = await api.post('/short-links', { target_url: url })
    ElMessage.success(`Created: /s/${r.code}`)
    load()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || 'Failed')
  }
}

onMounted(load)
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">Short Links</span>
        <el-button size="small" @click="create">+ Create Short Link</el-button>
      </div>
    </template>

    <el-table :data="list" stripe>
      <el-table-column prop="code" label="Code" width="100">
        <template #default="{ row }">
          <code class="bg-slate-100 px-2 py-0.5 rounded text-xs">{{ row.code }}</code>
        </template>
      </el-table-column>
      <el-table-column prop="target_url" label="Target URL" show-overflow-tooltip />
      <el-table-column prop="click_count" label="Clicks" width="100" sortable />
      <el-table-column label="Expires" width="160">
        <template #default="{ row }">
          <span v-if="row.expires_at">{{ row.expires_at.slice(0,10) }}</span>
          <span v-else class="text-slate-400">never</span>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="140">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="copyCode(row)">Copy</el-button>
          <el-button size="small" type="danger" link @click="del(row)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
