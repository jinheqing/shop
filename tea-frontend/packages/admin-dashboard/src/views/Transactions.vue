<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const filter = ref({ gateway: '', status: '', order_id: '' })

async function load() {
  try { const d: any = await api.get('/payment/transactions', { params: filter.value }); list.value = d?.items || d || [] }
  catch {
    list.value = [
      { id: 1, order_id: 9, order_no: 'ORD-20260917-398535', payment_gateway: '2checkout', gateway_transaction_id: '20389123456', amount: 188.00, currency: 'GBP', status: 'success', paid_at: '2026-09-17T13:34:41Z', created_at: '2026-09-17T13:34:40Z' },
      { id: 2, order_id: 8, order_no: 'ORD-20260916-112233', payment_gateway: 'paypal', gateway_transaction_id: 'PAYPAL-7821-ABC', amount: 376.00, currency: 'GBP', status: 'success', paid_at: '2026-09-16T09:12:30Z', created_at: '2026-09-16T09:12:28Z' },
      { id: 3, order_id: 7, order_no: 'ORD-20260915-445566', payment_gateway: '2checkout', gateway_transaction_id: '20389000001', amount: 94.00, currency: 'GBP', status: 'refunded', paid_at: '2026-09-15T18:02:10Z', created_at: '2026-09-15T18:02:09Z' },
      { id: 4, order_id: 10, order_no: 'ORD-20260917-778899', payment_gateway: 'paypal', gateway_transaction_id: 'PAYPAL-6666-DEF', amount: 250.00, currency: 'GBP', status: 'failed', paid_at: null, created_at: '2026-09-17T14:01:00Z' },
    ]
  }
}
onMounted(load)

async function refund(id: number) {
  await api.post(`/payment/transactions/${id}/refund`, { reason: 'customer_request' })
  ElMessage.success('Refund initiated')
  load()
}

const totals = {
  total: list.value.reduce((s, t) => s + (t.status === 'success' ? t.amount : 0), 0),
  count: list.value.filter(t => t.status === 'success').length,
  refunds: list.value.filter(t => t.status === 'refunded').reduce((s, t) => s + t.amount, 0),
}
</script>
<template>
  <div class="space-y-4">
    <el-row :gutter="16">
      <el-col :span="6"><el-card shadow="hover"><div class="text-xs text-slate-500 mb-1">Successful Revenue</div><div class="font-serif text-2xl text-green-700">£{{ totals.total.toFixed(2) }}</div><div class="text-xs text-slate-500">{{ totals.count }} transactions</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-xs text-slate-500 mb-1">Refunded</div><div class="font-serif text-2xl text-red-600">£{{ totals.refunds.toFixed(2) }}</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-xs text-slate-500 mb-1">2Checkout Volume</div><div class="font-serif text-2xl">£{{ list.filter(t=>t.payment_gateway==='2checkout'&&t.status==='success').reduce((s,t)=>s+t.amount,0).toFixed(2) }}</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><div class="text-xs text-slate-500 mb-1">PayPal Volume</div><div class="font-serif text-2xl">£{{ list.filter(t=>t.payment_gateway==='paypal'&&t.status==='success').reduce((s,t)=>s+t.amount,0).toFixed(2) }}</div></el-card></el-col>
    </el-row>

    <el-card>
      <template #header><div class="flex justify-between items-center"><span class="font-medium">💳 Payment Transactions</span>
        <div class="flex gap-2">
          <el-select v-model="filter.gateway" placeholder="All gateways" clearable style="width:150px">
            <el-option value="2checkout">2Checkout</el-option>
            <el-option value="paypal">PayPal</el-option>
          </el-select>
          <el-select v-model="filter.status" placeholder="All status" clearable style="width:150px">
            <el-option value="pending">Pending</el-option>
            <el-option value="processing">Processing</el-option>
            <el-option value="success">Success</el-option>
            <el-option value="failed">Failed</el-option>
            <el-option value="refunded">Refunded</el-option>
          </el-select>
          <el-button type="primary" @click="load">🔍 Filter</el-button>
          <el-button @click="">📥 Export CSV</el-button>
        </div>
      </div></template>
      <el-table :data="list" stripe>
        <el-table-column prop="created_at" label="Created" width="170" />
        <el-table-column prop="order_no" label="Order" width="170" />
        <el-table-column prop="payment_gateway" label="Gateway" width="120">
          <template #default="{ row }">
            <el-tag :type="row.payment_gateway==='2checkout'?'primary':'warning'" size="small">{{ row.payment_gateway }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="gateway_transaction_id" label="Tx ID" show-overflow-tooltip />
        <el-table-column label="Amount" width="120">
          <template #default="{ row }">£{{ row.amount.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="currency" label="Cur" width="60" />
        <el-table-column prop="status" label="Status" width="110">
          <template #default="{ row }">
            <el-tag :type="{'success':'success','pending':'info','processing':'warning','failed':'danger','refunded':'info','disputed':'danger'}[row.status] || 'info'" effect="dark">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paid_at" label="Paid At" width="170" />
        <el-table-column label="Actions" width="120">
          <template #default="{ row }">
            <el-button size="small" @click="">View</el-button>
            <el-button v-if="row.status==='success'" size="small" type="danger" @click="refund(row.id)">Refund</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
