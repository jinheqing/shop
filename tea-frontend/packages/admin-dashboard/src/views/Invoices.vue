<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/api/client'

const orders = ref<any[]>([])
const inv = ref<any>(null)
const loading = ref(false)

onMounted(async () => {
  const d: any = await api.get('/orders')
  orders.value = d?.items || d || []
})

async function fetchInv(id: number) {
  loading.value = true
  try { inv.value = await api.get(`/orders/${id}/invoice`) }
  finally { loading.value = false }
}
</script>

<template>
  <el-row :gutter="16">
    <el-col :span="10">
      <el-card>
        <template #header><span class="font-medium">Select Order</span></template>
        <el-table :data="orders" stripe @row-click="(r)=>fetchInv(r.id)" highlight-current-row>
          <el-table-column prop="order_no" label="Order" />
          <el-table-column prop="state" label="State" />
          <el-table-column prop="total_amount" label="Amount" />
        </el-table>
      </el-card>
    </el-col>
    <el-col :span="14">
      <el-card>
        <template #header><span class="font-medium">Invoice Preview</span></template>
        <div v-if="loading" class="text-center text-slate-400 py-20">Loading...</div>
        <div v-else-if="inv" class="p-6 border rounded-lg">
          <div class="flex justify-between items-start mb-6">
            <div>
              <div class="font-serif text-xl">UK Tea House Ltd.</div>
              <div class="text-xs text-slate-500">Mayfair, London W1J · VAT: GB123456789</div>
            </div>
            <div class="text-right">
              <div class="text-xs text-slate-500">Invoice No.</div>
              <div class="font-mono font-medium">{{ inv.invoice_no }}</div>
            </div>
          </div>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="HS Code">{{ inv.hs_code }}</el-descriptions-item>
            <el-descriptions-item label="Origin">{{ inv.country_of_origin }}</el-descriptions-item>
            <el-descriptions-item label="Total Amount">£{{ inv.total_amount }}</el-descriptions-item>
            <el-descriptions-item label="PDF Path">{{ inv.pdf_url }}</el-descriptions-item>
          </el-descriptions>
          <div class="mt-6 text-center">
            <el-button type="primary" :disabled="!inv.pdf_url">⬇ Download PDF</el-button>
          </div>
        </div>
        <div v-else class="text-center text-slate-400 py-20">Click an order on the left to generate/view invoice</div>
      </el-card>
    </el-col>
  </el-row>
</template>
