<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
async function load() { const d: any = await api.get('/orders'); list.value = d?.items || d || [] }
onMounted(load)

async function invoice(id: number) {
  const d: any = await api.get(`/orders/${id}/invoice`)
  ElMessage.success(`Invoice generated: ${d.invoice_no}`)
}
async function transition(id: number, state: string) {
  await api.post(`/orders/${id}/state`, { state })
  ElMessage.success(`→ ${state}`)
  load()
}

async function showTimeline(id: number) {
  const d: any = await api.get(`/orders/${id}/timeline`)
  ElMessageBox.alert(
    '<pre style="max-height:400px;overflow:auto;font-size:11px;white-space:pre-wrap">'+JSON.stringify(d, null, 2)+'</pre>',
    `Timeline — Order #${id}`,
    { dangerouslyUseHTMLString: true }
  )
}
</script>

<template>
  <el-card>
    <el-table :data="list" stripe>
      <el-table-column prop="order_no" label="Order No" width="180" />
      <el-table-column prop="hs_code" label="HS Code" width="100" />
      <el-table-column prop="total_amount" label="Amount" width="100" />
      <el-table-column prop="state" label="State" width="140">
        <template #default="{ row }">
          <el-tag :type="{ordering:'info',paid:'success',pending_declaration:'warning',producing:'warning',ready_for_production:'primary',shipped:'primary',completed:'success',cancelled:'danger'}[row.state] || 'info'">{{ row.state }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="300">
        <template #default="{ row }">
          <el-button size="small" @click="showTimeline(row.id)">⏱ Timeline</el-button>
          <el-button size="small" @click="invoice(row.id)">📄 Invoice</el-button>
          <el-button size="small" @click="showDecl(row.id)">📋 Decl</el-button>
          <el-dropdown size="small" @command="(s:string)=>transition(row.id,s)">
            <el-button size="small">Transition ▾</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="paid">→ paid</el-dropdown-item>
                <el-dropdown-item command="producing">→ producing</el-dropdown-item>
                <el-dropdown-item command="ready_for_production">→ ready_for_production</el-dropdown-item>
                <el-dropdown-item command="pending_customs">→ pending_customs</el-dropdown-item>
                <el-dropdown-item command="shipped">→ shipped</el-dropdown-item>
                <el-dropdown-item command="completed">→ completed</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
