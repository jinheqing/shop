<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
// 后端 LedgerCreateRequest:
// order_id(required), payment_gateway(required), gateway_transaction_id(required),
// amount_gbp(required), exchange_rate, amount_cny, bank_account, remarks
const form = reactive({
  order_id: 0,
  amount_gbp: 0,
  exchange_rate: null as number | null,
  amount_cny: null as number | null,
  payment_gateway: '2checkout',
  gateway_transaction_id: '',
  bank_account: '',
  remarks: '',
})

async function load() {
  try {
    const d: any = await api.get('/ledgers')
    list.value = d?.items || d || []
  } catch {}
}
onMounted(load)

async function save() {
  await api.post('/ledgers', form)
  ElMessage.success('Saved')
  open.value = false
  load()
}
async function del(id: number) {
  await api.delete(`/ledgers/${id}`)
  load()
}
</script>
<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">💱 收汇台账 Foreign Exchange Ledgers</span>
        <el-button type="primary" @click="open = true">+ New</el-button>
      </div>
    </template>
    <el-table :data="list" stripe>
      <el-table-column prop="order_id" label="Order ID" width="100" />
      <el-table-column prop="amount_gbp" label="GBP" width="100">
        <template #default="{ row }">£{{ row.amount_gbp }}</template>
      </el-table-column>
      <el-table-column prop="exchange_rate" label="FX Rate" width="90" />
      <el-table-column prop="amount_cny" label="CNY" width="100" />
      <el-table-column prop="payment_gateway" label="Gateway" width="120" />
      <el-table-column prop="gateway_transaction_id" label="Tx ID" show-overflow-tooltip />
      <el-table-column prop="bank_account" label="Bank" width="130" />
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="80">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="open" title="New FX Ledger Entry" width="580px">
    <el-form :model="form" label-width="150px">
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Order ID" required>
            <el-input-number v-model="form.order_id" class="w-full" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Amount GBP" required>
            <el-input-number v-model="form.amount_gbp" :precision="2" class="w-full" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Exchange Rate">
            <el-input-number v-model="form.exchange_rate" :precision="4" class="w-full" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Amount CNY">
            <el-input-number v-model="form.amount_cny" :precision="2" class="w-full" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Payment Gateway" required>
            <el-select v-model="form.payment_gateway" class="w-full">
              <el-option value="2checkout">2Checkout</el-option>
              <el-option value="paypal">PayPal</el-option>
              <el-option value="bank_transfer">Bank Transfer</el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Gateway Tx ID" required>
            <el-input v-model="form.gateway_transaction_id" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Bank Account">
            <el-input v-model="form.bank_account" placeholder="e.g. HSBC XXX" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Remarks">
            <el-input v-model="form.remarks" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
      <el-button @click="open = false">Cancel</el-button>
      <el-button type="primary" @click="save">Save</el-button>
    </template>
  </el-dialog>
</template>
