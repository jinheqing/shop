<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const codes = ref<any[]>([
  { id: 1, custom_product_id: 16, product_token: 'XK92AB38M7f2...f3A', mountain_location: '冰岛老寨', master_name: '王师傅', qr_code_position: 'outer_back', generated_at: '2026-09-17T13:40:00Z', public_url: '/trace/XK92AB38M7f2' },
  { id: 2, custom_product_id: 18, product_token: 'PHX91KL5N8b1...a7D', mountain_location: '凤凰山大乌岽', master_name: '李师傅', qr_code_position: 'outer_front', generated_at: '2026-09-17T12:00:00Z', public_url: '/trace/PHX91KL5N8b1' },
  { id: 3, custom_product_id: 20, product_token: 'JM72QR3T6V...e9F', mountain_location: '景迈山', master_name: '李师傅', qr_code_position: 'hidden', generated_at: '2026-09-16T09:00:00Z', public_url: '/trace/JM72QR3T6V' },
])

async function generate(productId: number) { await api.post(`/qrcodes/generate`, { custom_product_id: productId }); ElMessage.success('QR code regenerated') }
async function print(id: number) { ElMessage.success('🖨️ Print batch QR sheet ready') }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🔲 Traceability QR Codes</span>
      <el-button type="primary">+ Generate QR for Product</el-button>
    </div></template>
    <el-alert type="info" :closable="false" class="mb-4">
      Each bespoke product gets a unique QR. Scanning opens the public traceability page — GPS, master, harvest date, live camera embed.
    </el-alert>
    <el-table :data="codes" stripe>
      <el-table-column label="QR Preview" width="90">
        <template #default="{ row }">
          <div class="w-14 h-14 bg-slate-900 rounded flex items-center justify-center text-green-400 font-mono text-[8px] leading-none text-center">QR<br/>CODE</div>
        </template>
      </el-table-column>
      <el-table-column label="Product" width="250">
        <template #default="{ row }"><div class="font-medium">Product #{{ row.custom_product_id }}</div><div class="text-xs text-slate-500">{{ row.mountain_location }} · {{ row.master_name }}</div></template>
      </el-table-column>
      <el-table-column label="Public URL / Token" min-width="280">
        <template #default="{ row }">
          <code class="text-xs break-all">{{ row.public_url }}</code>
          <div class="text-xs text-slate-500 mt-1">{{ row.product_token }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="qr_code_position" label="Position" width="130">
        <template #default="{ row }">
          <el-tag effect="plain">{{ row.qr_code_position }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="generated_at" label="Generated" width="170" />
      <el-table-column label="Actions" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="print(row.id)">🖨️ Print</el-button>
          <el-button size="small" @click="generate(row.custom_product_id)">🔄 Regenerate</el-button>
          <el-button size="small" type="info">🔗 Copy URL</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
