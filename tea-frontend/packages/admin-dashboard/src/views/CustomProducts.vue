<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage, ElDialog } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const form = reactive<any>({ title: '', tea_type: 'raw_puer', tea_shape: 'cake', unit_price: 0, quantity: 0, shipping_cost: 0, inner_packaging: '', outer_packaging: '', qr_code_position: '', lead_time: '', mountain_location: '', master_name: '', storage_location: '', raw_tea_source: '', custom_requirement: '' })

async function load() { const d: any = await api.get('/custom-products'); list.value = d?.items || d || [] }
onMounted(load)

async function save() {
  await api.post('/custom-products', form)
  ElMessage.success('Created!')
  open.value = false
  Object.keys(form).forEach(k => form[k] = '')
  load()
}
async function publish(id: number) {
  const d: any = await api.post(`/custom-products/${id}/publish`)
  ElMessage.success(`Token: ${d.product_token?.slice(0, 16)}...`)
  load()
}
async function del(id: number) {
  await api.delete(`/custom-products/${id}`)
  ElMessage.success('Deleted')
  load()
}
</script>

<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">Bespoke Products</span><el-button type="primary" @click="open=true">+ New</el-button></div></template>
    <el-table :data="list" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="Title" />
      <el-table-column prop="tea_type" label="Type" width="110" />
      <el-table-column prop="tea_shape" label="Shape" width="100" />
      <el-table-column prop="unit_price" label="£" width="80" />
      <el-table-column prop="status" label="Status" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status==='published'?'success':'info'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="220">
        <template #default="{ row }">
          <el-button size="small" v-if="row.status!=='published'" @click="publish(row.id)">Publish</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="New Bespoke Product" width="720px">
    <el-form :model="form" label-width="120px" label-position="right">
      <el-form-item label="Title"><el-input v-model="form.title" /></el-form-item>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="Tea Type"><el-select v-model="form.tea_type" class="w-full"><el-option value="raw_puer" label="生普洱" /><el-option value="ripe_puer" label="熟普洱" /></el-select></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Shape"><el-input v-model="form.tea_shape" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8"><el-form-item label="Unit Price"><el-input-number v-model="form.unit_price" :precision="2" class="w-full" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Quantity"><el-input-number v-model="form.quantity" class="w-full" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Shipping"><el-input-number v-model="form.shipping_cost" :precision="2" class="w-full" /></el-form-item></el-col>
      </el-row>
      <el-form-item label="Mountain"><el-input v-model="form.mountain_location" /></el-form-item>
      <el-form-item label="Master"><el-input v-model="form.master_name" /></el-form-item>
      <el-form-item label="Lead Time"><el-input v-model="form.lead_time" /></el-form-item>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="Inner Pack"><el-input v-model="form.inner_packaging" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Outer Pack"><el-input v-model="form.outer_packaging" /></el-form-item></el-col>
      </el-row>
      <el-form-item label="Spec"><el-input v-model="form.custom_requirement" type="textarea" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open=false">Cancel</el-button>
      <el-button type="primary" @click="save">Save & Publish Later</el-button>
    </template>
  </el-dialog>
</template>
