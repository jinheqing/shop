<script setup lang="ts">
import { ref, computed } from 'vue'
import { api } from '@/api/client'

const rooms = ref<any[]>([
  { id: 1, room_name: '24/7 冰岛老寨慢直播', room_type: 'slow_preset', status: 'live', scheduled_start: null, host_staff_id: null },
  { id: 5, room_name: '凤凰山大乌岽 · 品鉴会', room_type: 'obs_tasting', status: 'scheduled', scheduled_start: '2026-09-20T20:00:00Z', host_staff_id: 5 },
  { id: 7, room_name: 'James 先生 · 定制茶交付验货', room_type: 'delivery_inspection', status: 'scheduled', scheduled_start: '2026-09-22T14:00:00Z', host_staff_id: 5, order_id: 9 },
  { id: 8, room_name: 'Sarah 女士 · 私享咨询', room_type: 'customer_request', status: 'scheduled', scheduled_start: '2026-09-24T18:00:00Z', host_staff_id: 6 },
  { id: 9, room_name: '秋季新茶发布会', room_type: 'open_calendar', status: 'scheduled', scheduled_start: '2026-09-28T19:00:00Z', host_staff_id: 5 },
])

const typeColors: any = {
  slow_preset: 'bg-green-100 border-green-400 text-green-800',
  obs_tasting: 'bg-blue-100 border-blue-400 text-blue-800',
  open_calendar: 'bg-purple-100 border-purple-400 text-purple-800',
  customer_request: 'bg-yellow-100 border-yellow-400 text-yellow-800',
  delivery_inspection: 'bg-red-100 border-red-400 text-red-800',
  custom_private: 'bg-pink-100 border-pink-400 text-pink-800',
}

const roomsByDay = computed(() => {
  const map: any = {}
  for (let d = 1; d <= 30; d++) map[d] = []
  for (const r of rooms.value) {
    if (!r.scheduled_start) continue
    const day = new Date(r.scheduled_start).getUTCDate()
    if (map[day]) map[day].push(r)
  }
  return map
})
</script>
<template>
  <div class="space-y-4">
    <el-card>
      <template #header><div class="flex justify-between items-center">
        <span class="font-medium">📅 Live Room Schedule — September 2026</span>
        <el-button type="primary" @click="">+ Schedule Room</el-button>
      </div></template>

      <div class="grid grid-cols-7 gap-1 text-center text-xs text-slate-500 mb-2">
        <div>Mon</div><div>Tue</div><div>Wed</div><div>Thu</div><div>Fri</div><div>Sat</div><div>Sun</div>
      </div>
      <div class="grid grid-cols-7 gap-1">
        <template v-for="(empty, i) in 0" :key="i"></template>
        <div v-for="day in 30" :key="day" class="min-h-[100px] border border-slate-200 p-1 rounded relative hover:bg-slate-50">
          <div class="text-xs text-slate-400 mb-1">{{ day }}</div>
          <div v-for="r in (roomsByDay[day] || [])" :key="r.id"
            :class="['text-xs p-1 mb-0.5 rounded border truncate cursor-pointer', typeColors[r.room_type]]"
            :title="r.room_name">
            <span class="font-medium">{{ r.room_type==='delivery_inspection' ? '📦' : r.room_type==='slow_preset' ? '🌱' : r.room_type==='obs_tasting' ? '🍵' : r.room_type==='customer_request' ? '👤' : r.room_type==='open_calendar' ? '📺' : '🔒' }}</span>
            {{ r.room_name.slice(0, 15) }}
          </div>
        </div>
      </div>
    </el-card>

    <el-card>
      <template #header><span class="font-medium">📋 Upcoming Live Rooms</span></template>
      <el-table :data="rooms" stripe>
        <el-table-column prop="room_name" label="Name" width="250" />
        <el-table-column prop="room_type" label="Type" width="170">
          <template #default="{ row }">
            <el-tag :type="{'slow_preset':'success','obs_tasting':'primary','open_calendar':'','customer_request':'warning','delivery_inspection':'danger','custom_private':''}[row.room_type]" effect="dark" size="small">{{ row.room_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="Status" width="110">
          <template #default="{ row }">
            <el-tag :type="{'live':'success','scheduled':'info','offline':'warning','ended':''}[row.status]" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scheduled_start" label="Scheduled" width="180" />
        <el-table-column label="Actions" width="160">
          <template #default="{ row }">
            <el-button v-if="row.status==='scheduled'" size="small" type="success">▶ Go Live</el-button>
            <el-button size="small">Edit</el-button>
            <el-button v-if="row.order_id" size="small" type="warning">Order {{ row.order_id }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
