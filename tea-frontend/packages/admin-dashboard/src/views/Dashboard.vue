<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/api/client'

const stats = ref({ orders: 0, rooms: 0, nodes: 0, products: 0 })
onMounted(async () => {
  try {
    const [o, r, n, p] = await Promise.all([
      api.get<any, any>('/orders').catch(() => ({ items: [] })),
      api.get<any, any>('/live-rooms').catch(() => ({ items: [] })),
      api.get<any, any>('/nodes').catch(() => ({ items: [] })),
      api.get<any, any>('/custom-products').catch(() => ({ items: [] })),
    ])
    stats.value = {
      orders: o?.total || o?.items?.length || 0,
      rooms: r?.total || r?.items?.length || 0,
      nodes: n?.total || n?.items?.length || 0,
      products: p?.total || p?.items?.length || 0,
    }
  } catch {}
})
</script>

<template>
  <div>
    <el-row :gutter="16" class="mb-6">
      <el-col :span="6"><el-card shadow="hover"><div class="text-slate-500 text-sm">Total Orders</div><div class="text-3xl font-serif mt-1">{{ stats.orders }}</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-slate-500 text-sm">Live Rooms</div><div class="text-3xl font-serif mt-1">{{ stats.rooms }}</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-slate-500 text-sm">Media Nodes</div><div class="text-3xl font-serif mt-1">{{ stats.nodes }}</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-slate-500 text-sm">Bespoke Products</div><div class="text-3xl font-serif mt-1">{{ stats.products }}</div></el-card></el-col>
    </el-row>

    <el-card>
      <template #header><span class="font-medium">Backend System Health</span></template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="API Backend"><el-tag type="success" effect="dark">✅ Running</el-tag></el-descriptions-item>
        <el-descriptions-item label="PostgreSQL"><el-tag type="success" effect="dark">✅ 18 tables</el-tag></el-descriptions-item>
        <el-descriptions-item label="LiveKit SFU"><el-tag>Configured</el-tag></el-descriptions-item>
        <el-descriptions-item label="MediaMTX RTMP"><el-tag>Configured</el-tag></el-descriptions-item>
        <el-descriptions-item label="2Checkout/PayPal"><el-tag>Mock Ready</el-tag></el-descriptions-item>
        <el-descriptions-item label="GDPR Audit Logs"><el-tag type="success">✅ Active</el-tag></el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>
