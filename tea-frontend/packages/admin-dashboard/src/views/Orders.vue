<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const loading = ref(false)

// ===== Create Order =====
const createOpen = ref(false)
const cpOptions = ref<any[]>([])
const creating = ref(false)
const createForm = ref({
  custom_product_id: 0 as number,
  user_id: '' as string,
  unit_price: '' as string,
  quantity: '' as string,
  shipping_cost: '' as string,
  hs_code: '' as string,
  country_of_origin: '' as string,
})

async function load() {
  loading.value = true
  try {
    const d: any = await api.get('/orders')
    list.value = d?.items || d || []
  } finally {
    loading.value = false
  }
}

async function loadCpOptions() {
  const d: any = await api.get('/custom-products')
  cpOptions.value = (d?.items || d || []).filter((p: any) => p.status !== 'cancelled')
}

async function openCreate() {
  await loadCpOptions()
  createForm.value = {
    custom_product_id: cpOptions.value[0]?.id || 0,
    user_id: '',
    unit_price: '',
    quantity: '',
    shipping_cost: '',
    hs_code: '',
    country_of_origin: '',
  }
  createOpen.value = true
}

async function submitCreate() {
  if (!createForm.value.custom_product_id) {
    ElMessage.warning('Please select a custom product')
    return
  }
  creating.value = true
  try {
    const payload: any = { custom_product_id: createForm.value.custom_product_id }
    if (createForm.value.user_id) payload.user_id = Number(createForm.value.user_id)
    if (createForm.value.unit_price) payload.unit_price = Number(createForm.value.unit_price)
    if (createForm.value.quantity) payload.quantity = Number(createForm.value.quantity)
    if (createForm.value.shipping_cost) payload.shipping_cost = Number(createForm.value.shipping_cost)
    if (createForm.value.hs_code) payload.hs_code = createForm.value.hs_code
    if (createForm.value.country_of_origin) payload.country_of_origin = createForm.value.country_of_origin

    await api.post('/orders', payload)
    ElMessage.success('Order created')
    createOpen.value = false
    load()
  } catch (e: any) {
    ElMessage.error('Create failed: ' + (e?.message || e))
  } finally {
    creating.value = false
  }
}

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
    <template #header>
      <div style="display:flex;align-items:center;justify-content:space-between">
        <span style="font-weight:600">Orders</span>
        <el-button type="primary" @click="openCreate">+ Create Order</el-button>
      </div>
    </template>
    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="order_no" label="Order No" width="180" />
      <el-table-column prop="hs_code" label="HS Code" width="100" />
      <el-table-column prop="total_amount" label="Amount" width="100" />
      <el-table-column prop="state" label="State" width="140">
        <template #default="{ row }">
          <el-tag :type="({ordering:'info',paid:'success',pending_declaration:'warning',producing:'warning',ready_for_production:'primary',shipped:'primary',completed:'success',cancelled:'danger'} as Record<string,string>)[row.state] || 'info'">{{ row.state }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="Created" width="180" />
      <el-table-column label="Actions" width="300">
        <template #default="{ row }">
          <el-button size="small" @click="showTimeline(row.id)">⏱ Timeline</el-button>
          <el-button size="small" @click="invoice(row.id)">📄 Invoice</el-button>
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

  <!-- Create Order Dialog -->
  <el-dialog v-model="createOpen" title="+ Create Order" width="560px">
    <el-form :model="createForm" label-width="130px">
      <el-form-item label="Custom Product" required>
        <el-select v-model="createForm.custom_product_id" placeholder="Select..." style="width:100%">
          <el-option v-for="cp in cpOptions" :key="cp.id" :label="`#${cp.id} ${cp.title} (${cp.sku || 'no-sku'})`" :value="cp.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="User ID (optional)">
        <el-input v-model="createForm.user_id" placeholder="If placing for a specific user, enter their user_id" />
      </el-form-item>
      <el-form-item label="Unit Price (override)">
        <el-input v-model="createForm.unit_price" placeholder="Leave empty = use CP default" />
      </el-form-item>
      <el-form-item label="Quantity (override)">
        <el-input v-model="createForm.quantity" placeholder="Leave empty = use CP default" />
      </el-form-item>
      <el-form-item label="Shipping Cost (override)">
        <el-input v-model="createForm.shipping_cost" placeholder="Leave empty = use CP default" />
      </el-form-item>
      <el-form-item label="HS Code">
        <el-input v-model="createForm.hs_code" placeholder="e.g. 0902.10" />
      </el-form-item>
      <el-form-item label="Country of Origin">
        <el-input v-model="createForm.country_of_origin" placeholder="e.g. China" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="createOpen = false">Cancel</el-button>
      <el-button type="primary" :loading="creating" @click="submitCreate">Create</el-button>
    </template>
  </el-dialog>
</template>
