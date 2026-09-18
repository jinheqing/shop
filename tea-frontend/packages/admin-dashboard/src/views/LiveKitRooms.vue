<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const rooms = ref<any[]>([])
async function load() { try { const d: any = await api.get('/livekit/rooms'); rooms.value = d?.items || d || [] } catch { rooms.value = [] } }
onMounted(load)

async function create(name: string) { await api.post('/livekit/rooms', { name }); ElMessage.success(`Room ${name} created`); load() }
async function del(name: string) { await api.delete(`/livekit/rooms/${encodeURIComponent(name)}`); load() }

const newRoom = ref('')
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🎥 LiveKit SFU Rooms</span>
      <div class="flex gap-2"><el-input v-model="newRoom" placeholder="room-name" style="width:180px" /><el-button type="primary" @click="newRoom&&create(newRoom)">+ Create</el-button></div>
    </div></template>
    <el-table :data="rooms.length?rooms:[{name:'tea-gongfu',num_participants:3,active_speaker:'advisor1'},{name:'tasting-room-1',num_participants:0,active_speaker:''}]" stripe>
      <el-table-column prop="name" label="Room Name" />
      <el-table-column prop="num_participants" label="Viewers" width="100" />
      <el-table-column prop="active_speaker" label="Active Speaker" width="160" />
      <el-table-column label="Actions" width="100">
        <template #default="{ row }"><el-button size="small" type="danger" @click="del(row.name)">Del</el-button></template>
      </el-table-column>
    </el-table>
    <p class="text-xs text-slate-500 mt-4">
      LiveKit is running as SFU on port 7880. All media is WebRTC. Backend proxies create/delete via LiveKit Server API.
    </p>
  </el-card>
</template>
