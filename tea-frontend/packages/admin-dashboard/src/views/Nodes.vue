<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ name: '', ip: '', username: 'root', password: '', port: 22, node_type: 'livekit_edge' })

async function load() { const d: any = await api.get('/nodes'); list.value = d?.items || d || [] }
onMounted(load)
async function save() { await api.post('/nodes/deploy', form); ElMessage.success('Deploying...'); open.value = false; load() }
async function del(id: number) { await api.delete(`/nodes/${id}`); load() }
</script>

<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">Media Nodes · WireGuard Mesh</span><el-button type="primary" @click="open=true">+ Add Node</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="node_name" label="Name" width="140" />
      <el-table-column prop="node_type" label="Type" width="140">
        <template #default="{ row }"><el-tag size="small">{{ row.node_type }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="public_ip" label="Public IP" width="140" />
      <el-table-column prop="wireguard_ip" label="WG IP" width="140" />
      <el-table-column prop="status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status==='online'?'success':row.status==='failed'?'danger':'warning'" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Resource" min-width="240">
        <template #default="{ row }">
          <div class="flex items-center gap-4 text-xs">
            <span>CPU: {{ row.last_cpu_percent ?? '—' }}%</span>
            <span>MEM: {{ row.last_memory_percent ?? '—' }}%</span>
            <span>👁️ {{ row.last_viewers_count ?? 0 }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="100">
        <template #default="{ row }"><el-button size="small" type="danger" @click="del(row.id)">Del</el-button></template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="One-Click Node Deploy" width="520px">
    <p class="text-xs text-slate-500 mb-4">SSH password is ephemeral — used once, then destroyed. WireGuard private key never leaves this server.</p>
    <el-form :model="form" label-width="100px">
      <el-form-item label="Name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="Public IP"><el-input v-model="form.ip" /></el-form-item>
      <el-form-item label="Type">
        <el-select v-model="form.node_type" class="w-full">
          <el-option value="livekit_edge" label="LiveKit Edge" />
          <el-option value="mediamtx" label="MediaMTX RTMP" />
          <el-option value="cdn" label="CDN Edge" />
        </el-select>
      </el-form-item>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="SSH User"><el-input v-model="form.username" /></el-form-item></el-col>
        <el-col :span="6"><el-form-item label="Port"><el-input-number v-model="form.port" class="w-full" /></el-form-item></el-col>
        <el-col :span="6"><el-form-item label="Pwd"><el-input v-model="form.password" type="password" show-password /></el-form-item></el-col>
      </el-row>
    </el-form>
    <template #footer>
      <el-button @click="open=false">Cancel</el-button>
      <el-button type="primary" @click="save">🚀 Deploy</el-button>
    </template>
  </el-dialog>
</template>
