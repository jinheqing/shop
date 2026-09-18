<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ ledger_no: '', order_id: 0, amount_gbp: 0, currency: 'GBP', fx_rate: 9.1, payment_gateway: '2checkout', gateway_transaction_id: '' })
async function load() { try { const d: any = await api.get('/ledgers'); list.value = d?.items || d || [] } catch {} }
onMounted(load)
async function save() { await api.post('/ledgers', form); ElMessage.success('Saved'); open.value = false; load() }
async function del(id: number) { await api.delete(`/ledgers/${id}`); load() }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">💱 收汇台账 Foreign Exchange Ledgers</span><el-button type="primary" @click="open=true">+ New</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="ledger_no" label="Ledger No" width="180" />
      <el-table-column prop="order_id" label="Order ID" width="100" />
      <el-table-column prop="amount_gbp" label="GBP" width="100"><template #default="{ row }">£{{ row.amount_gbp }}</template></el-table-column>
      <el-table-column prop="currency" label="Currency" width="90" />
      <el-table-column prop="fx_rate" label="FX Rate" width="90" />
      <el-table-column prop="payment_gateway" label="Gateway" width="120" />
      <el-table-column prop="gateway_transaction_id" label="Tx ID" show-overflow-tooltip />
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="80"><template #default="{ row }"><el-button size="small" type="danger" @click="del(row.id)">Del</el-button></template></el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="open" title="New FX Ledger Entry" width="560px">
    <el-form :model="form" label-width="120px">
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="Ledger No"><el-input v-model="form.ledger_no" placeholder="FX-YYYYNNNN" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Order ID"><el-input-number v-model="form.order_id" class="w-full" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8"><el-form-item label="Amount GBP"><el-input-number v-model="form.amount_gbp" :precision="2" class="w-full" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Currency"><el-input v-model="form.currency" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="FX Rate"><el-input-number v-model="form.fx_rate" :precision="4" class="w-full" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="Gateway"><el-select v-model="form.payment_gateway" class="w-full"><el-option value="2checkout">2Checkout</el-option><el-option value="paypal">PayPal</el-option><el-option value="bank_transfer">Bank Transfer</el-option></el-select></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Tx ID"><el-input v-model="form.gateway_transaction_id" /></el-form-item></el-col>
      </el-row>
    </el-form>
    <template #footer><el-button @click="open=false">Cancel</el-button><el-button type="primary" @click="save">Save</el-button></template>
  </el-dialog>
</template>
