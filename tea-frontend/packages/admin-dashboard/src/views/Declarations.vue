<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive({ declaration_no: '', order_id: 0, hs_code: '0902.10', commodity_desc: '普洱茶饼', origin: 'China', destination: 'UK', amount: 0, customs_status: 'pending' })
async function load() { try { const d: any = await api.get('/declarations'); list.value = d?.items || d || [] } catch {} }
onMounted(load)
async function save() { await api.post('/declarations', form); ElMessage.success('Created'); open.value = false; load() }
async function del(id: number) { await api.delete(`/declarations/${id}`); load() }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">📑 报关台账 Declarations</span><el-button type="primary" @click="open=true">+ New</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="declaration_no" label="Declaration No" width="180" />
      <el-table-column prop="order_id" label="Order ID" width="100" />
      <el-table-column prop="hs_code" label="HS Code" width="110" />
      <el-table-column prop="origin" label="Origin" width="100" />
      <el-table-column prop="destination" label="Dest" width="80" />
      <el-table-column prop="amount" label="Amount" width="100">
        <template #default="{ row }">£{{ row.amount }}</template>
      </el-table-column>
      <el-table-column prop="customs_status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="{'pending':'warning','cleared':'success','held':'danger'}[row.customs_status] || 'info'" effect="dark">{{ row.customs_status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="80">
        <template #default="{ row }"><el-button size="small" type="danger" @click="del(row.id)">Del</el-button></template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="open" title="New Declaration" width="560px">
    <el-form :model="form" label-width="120px">
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="Declaration No"><el-input v-model="form.declaration_no" placeholder="DEC-YYYYNNNN" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Order ID"><el-input-number v-model="form.order_id" class="w-full" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="HS Code"><el-input v-model="form.hs_code" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Commodity"><el-input v-model="form.commodity_desc" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8"><el-form-item label="Origin"><el-input v-model="form.origin" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Destination"><el-input v-model="form.destination" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Amount (£)"><el-input-number v-model="form.amount" :precision="2" class="w-full" /></el-form-item></el-col>
      </el-row>
    </el-form>
    <template #footer><el-button @click="open=false">Cancel</el-button><el-button type="primary" @click="save">Save</el-button></template>
  </el-dialog>
</template>
