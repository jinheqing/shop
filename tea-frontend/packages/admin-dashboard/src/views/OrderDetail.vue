<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElTabs } from 'element-plus'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const id = parseInt(route.params.id as string)
const order = ref<any>(null)
const inv = ref<any>(null)
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    order.value = await api.get(`/orders/${id}`)
    try { inv.value = await api.get(`/orders/${id}/invoice`) } catch {}
  } finally { loading.value = false }
}
onMounted(load)

const trackingNumber = ref('')
const shippingCarrier = ref('')

async function transition(state: string, extra: Record<string, string> = {}) {
  try {
    await api.post(`/orders/${id}/state`, { target_state: state, ...extra })
    ElMessage.success(`→ ${state}`)
    if (state === 'shipped') { trackingNumber.value = ''; shippingCarrier.value = '' }
    load()
  } catch (e: any) { ElMessage.error(e?.message || 'Transition failed') }
}
async function shipWithTracking() {
  if (!trackingNumber.value || !shippingCarrier.value) {
    ElMessage.warning('Please enter tracking number and carrier')
    return
  }
  await transition('shipped', {
    tracking_number: trackingNumber.value,
    shipping_carrier: shippingCarrier.value,
  })
}
async function cancel() { await api.post(`/orders/${id}/cancel`); ElMessage.success('Cancelled'); load() }
async function regenInv() { await api.post(`/orders/${id}/invoice/regenerate`); ElMessage.success('Invoice regenerated'); load() }

const timeline = computed(() => {
  if (!order.value) return []
  const states = ['ordering','paid','pending_declaration','producing','ready_for_production','pending_customs','customs_clear','shipped','completed']
  const idx = states.indexOf(order.value.state)
  return states.map((s, i) => ({ name: s, done: i < idx, current: i === idx }))
})
</script>
<template>
  <div v-if="loading" class="text-center py-20"><el-icon class="is-loading text-4xl"><Loading /></el-icon></div>
  <div v-else-if="order" class="space-y-4">
    <el-card>
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="font-serif text-2xl text-slate-900">{{ order.order_no }}</h2>
          <p class="text-sm text-slate-500">Order ID #{{ order.id }} · Placed {{ order.created_at }}</p>
        </div>
        <div class="flex gap-2">
          <el-button @click="router.back()">← Back</el-button>
        </div>
      </div>

      <el-steps :active="timeline.findIndex(s=>s.current)+1" finish-status="success" align-center>
        <el-step v-for="s in timeline" :key="s.name" :title="s.name" :status="s.current?'process':s.done?'success':'wait'" />
      </el-steps>
    </el-card>

    <el-row :gutter="16">
      <el-col :span="14">
        <el-card>
          <h3 class="font-medium mb-3">🔄 Status Actions</h3>
          <el-space wrap>
            <el-button v-if="order.state==='ordering'" type="success" @click="transition('paid')">Mark Paid</el-button>
            <el-button v-if="order.state==='paid'" @click="transition('pending_declaration')">→ Declaration</el-button>
            <el-button v-if="order.state==='paid' || order.state==='pending_declaration'" type="warning" @click="transition('producing')">→ Producing</el-button>
            <el-button v-if="order.state==='producing'" type="warning" @click="transition('ready_for_production')">→ Ready for Production</el-button>
            <el-button v-if="order.state==='ready_for_production'" @click="transition('pending_customs')">→ Pending Customs</el-button>
            <el-button v-if="order.state==='pending_customs'" type="primary" @click="transition('customs_clear')">→ Customs Clear</el-button>
            <el-button v-if="order.state==='customs_clear'" type="primary" @click="shipWithTracking">→ Shipped</el-button>
            <el-button v-if="order.state==='shipped'" type="success" @click="transition('completed')">→ ✅ Completed</el-button>
            <el-button v-if="order.state!=='cancelled' && order.state!=='completed'" type="danger" @click="cancel">Cancel Order</el-button>
          </el-space>
          <div v-if="order.state==='customs_clear'" class="mt-4 flex flex-wrap items-end gap-3 p-3 bg-slate-50 rounded">
            <div>
              <label class="block text-xs text-slate-500 mb-1">Tracking Number</label>
              <el-input v-model="trackingNumber" placeholder="e.g. DHL123456789" style="width:220px" />
            </div>
            <div>
              <label class="block text-xs text-slate-500 mb-1">Shipping Carrier</label>
              <el-input v-model="shippingCarrier" placeholder="e.g. DHL / FedEx / UPS" style="width:180px" />
            </div>
            <el-button type="primary" @click="shipWithTracking">→ Mark Shipped</el-button>
          </div>
          <div v-if="order.state==='shipped' || order.state==='completed'" class="mt-3 text-sm text-slate-600">
            <span class="text-xs text-slate-400 mr-2">Tracking:</span>
            <span>{{ order.shipping_carrier || '—' }} · {{ order.tracking_number || '—' }}</span>
          </div>
        </el-card>

        <el-card class="mt-4">
          <h3 class="font-medium mb-3">📄 Bespoke Product Snapshot</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="Title">{{ order.custom_product_snapshot?.title }}</el-descriptions-item>
            <el-descriptions-item label="Tea Type">{{ order.custom_product_snapshot?.tea_type }}</el-descriptions-item>
            <el-descriptions-item label="Tea Garden">{{ order.custom_product_snapshot?.tea_garden_location }}</el-descriptions-item>
            <el-descriptions-item label="Master">{{ order.custom_product_snapshot?.master_name }}</el-descriptions-item>
            <el-descriptions-item label="Inner Pack">{{ order.custom_product_snapshot?.inner_packaging }}</el-descriptions-item>
            <el-descriptions-item label="Outer Pack">{{ order.custom_product_snapshot?.outer_packaging }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card>
          <h3 class="font-medium mb-3">💷 Financials</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between"><span>Unit Price</span><span>£{{ order.unit_price }}</span></div>
            <div class="flex justify-between"><span>Quantity</span><span>{{ order.quantity }}</span></div>
            <div class="flex justify-between"><span>Shipping</span><span>£{{ order.shipping_cost }}</span></div>
            <el-divider class="my-2" />
            <div class="flex justify-between font-serif text-lg"><span>Total</span><span>£{{ order.total_amount }}</span></div>
          </div>
          <div class="mt-4 pt-4 border-t">
            <div class="flex justify-between text-xs text-slate-500 mb-2"><span>HS Code</span><span>{{ order.hs_code }}</span></div>
            <div class="flex justify-between text-xs text-slate-500 mb-3"><span>Origin</span><span>{{ order.country_of_origin }}</span></div>
            <el-button size="small" @click="regenInv">🔄 Regenerate Invoice</el-button>
          </div>
        </el-card>

        <el-card class="mt-4">
          <h3 class="font-medium mb-3">📍 Addresses</h3>
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <div class="text-xs text-slate-500 mb-1">Billing</div>
              <div>{{ order.billing_address_snapshot?.name }}</div>
              <div class="text-slate-600">{{ order.billing_address_snapshot?.address }}, {{ order.billing_address_snapshot?.city }}</div>
              <div class="text-slate-600">{{ order.billing_address_snapshot?.postcode }}, {{ order.billing_address_snapshot?.country }}</div>
            </div>
            <div>
              <div class="text-xs text-slate-500 mb-1">Delivery</div>
              <div>{{ order.delivery_address_snapshot?.name }}</div>
              <div class="text-slate-600">{{ order.delivery_address_snapshot?.address }}, {{ order.delivery_address_snapshot?.city }}</div>
              <div class="text-slate-600">{{ order.delivery_address_snapshot?.postcode }}, {{ order.delivery_address_snapshot?.country }}</div>
            </div>
          </div>
        </el-card>

        <el-card v-if="inv" class="mt-4">
          <h3 class="font-medium mb-2">🧾 Invoice</h3>
          <div class="text-sm space-y-1">
            <div>No. <code>{{ inv.invoice_no }}</code></div>
            <div class="text-xs text-slate-500">PDF: {{ inv.pdf_url }}</div>
            <el-button size="small" type="primary">⬇ Download PDF</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
