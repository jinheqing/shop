<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
// 后端 DeclarationCreateRequest 字段:
// order_id(required), customs_declaration_no, hs_code(required), commodity_desc(required),
// gross_weight, net_weight, declared_value, customs_status, declaration_date, cleared_date, remarks
const form = reactive({
  customs_declaration_no: '',
  order_id: 0,
  hs_code: '0902.10',
  commodity_desc: '普洱茶饼',
  gross_weight: null as number | null,
  net_weight: null as number | null,
  declared_value: null as number | null,
  customs_status: 'pending',
  declaration_date: '',
  cleared_date: '',
  remarks: '',
})

async function load() {
  try {
    const d: any = await api.get('/declarations')
    list.value = d?.items || d || []
  } catch {}
}
onMounted(load)

async function save() {
  await api.post('/declarations', form)
  ElMessage.success('Created')
  open.value = false
  load()
}
async function del(id: number) {
  await api.delete(`/declarations/${id}`)
  load()
}
</script>
<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">📑 报关台账 Declarations</span>
        <el-button type="primary" @click="open = true">+ New</el-button>
      </div>
    </template>
    <el-table :data="list" stripe>
      <el-table-column prop="customs_declaration_no" label="Decl No" width="180" />
      <el-table-column prop="order_id" label="Order ID" width="100" />
      <el-table-column prop="hs_code" label="HS Code" width="110" />
      <el-table-column prop="commodity_desc" label="Commodity" />
      <el-table-column prop="declared_value" label="Declared" width="100">
        <template #default="{ row }">£{{ row.declared_value }}</template>
      </el-table-column>
      <el-table-column prop="customs_status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="({pending:'warning',cleared:'success',held:'danger'} as Record<string,string>)[row.customs_status] || 'info'" effect="dark">
            {{ row.customs_status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="80">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="open" title="New Declaration" width="620px">
    <el-form :model="form" label-width="140px">
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Customs Decl No">
            <el-input v-model="form.customs_declaration_no" placeholder="DEC-YYYYNNNN" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Order ID" required>
            <el-input-number v-model="form.order_id" class="w-full" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="HS Code" required>
            <el-input v-model="form.hs_code" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Commodity" required>
            <el-input v-model="form.commodity_desc" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="Gross Weight (g)">
            <el-input-number v-model="form.gross_weight" :precision="1" class="w-full" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Net Weight (g)">
            <el-input-number v-model="form.net_weight" :precision="1" class="w-full" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Declared (£)">
            <el-input-number v-model="form.declared_value" :precision="2" class="w-full" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="Customs Status">
            <el-select v-model="form.customs_status" class="w-full">
              <el-option value="pending" label="Pending" />
              <el-option value="cleared" label="Cleared" />
              <el-option value="held" label="Held" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Declaration Date">
            <el-date-picker v-model="form.declaration_date" type="date" class="w-full" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Cleared Date">
            <el-date-picker v-model="form.cleared_date" type="date" class="w-full" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="Remarks">
        <el-input v-model="form.remarks" type="textarea" :rows="2" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">Cancel</el-button>
      <el-button type="primary" @click="save">Save</el-button>
    </template>
  </el-dialog>
</template>
