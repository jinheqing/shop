<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const id = parseInt(route.params.id as string)
const order = ref<any>(null)
const inv = ref<any>(null)
const timeline = ref<any>(null)
const loading = ref(true)

// 状态流转弹窗
const transitionDialog = reactive({
  visible: false,
  targetState: '',
  reason: '',
  courier: '',
  tracking_no: '',
  shipped_at: '',
  eta_at: '',
})

// 物流独立编辑弹窗
const shippingDialog = reactive({
  visible: false,
  courier: '',
  tracking_no: '',
  shipped_at: '',
  eta_at: '',
})

// 后端正确的状态序列（和 state machine 对齐）
const STATE_FLOW = [
  'ordering',
  'paid',
  'pending_declaration',
  'producing',
  'ready_for_production',
  'pending_customs',
  'customs_clear',
  'shipped',
  'completed',
]

// 每个状态的中文标签（el-steps 显示用）
const STATE_LABELS: Record<string, string> = {
  ordering: 'Order Placed',
  paid: 'Payment Confirmed',
  pending_declaration: 'Declaration Pending',
  producing: 'In Production',
  ready_for_production: 'Ready for Shipment',
  pending_customs: 'Awaiting Customs',
  customs_clear: 'Customs Cleared',
  shipped: 'Dispatched',
  completed: 'Delivered',
  cancelled: 'Cancelled',
  disputed: 'Disputed',
}

async function load() {
  loading.value = true
  try {
    order.value = await api.get(`/orders/${id}`)
    try { inv.value = await api.get(`/orders/${id}/invoice`) } catch {}
    try { timeline.value = await api.get(`/orders/${id}/timeline`) } catch {}
  } finally { loading.value = false }
}
onMounted(load)

const stepActive = computed(() => {
  if (!order.value) return 0
  const idx = STATE_FLOW.indexOf(order.value.state)
  return idx >= 0 ? idx + 1 : 0
})

// 打开状态流转弹窗
function openTransition(targetState: string) {
  transitionDialog.targetState = targetState
  transitionDialog.reason = ''
  transitionDialog.courier = order.value?.courier || ''
  transitionDialog.tracking_no = order.value?.tracking_no || ''
  transitionDialog.shipped_at = order.value?.shipped_at || ''
  transitionDialog.eta_at = order.value?.eta_at || ''
  transitionDialog.visible = true
}

// 打开物流独立编辑弹窗
function openShippingEdit() {
  shippingDialog.courier = order.value?.courier || ''
  shippingDialog.tracking_no = order.value?.tracking_no || ''
  shippingDialog.shipped_at = order.value?.shipped_at || ''
  shippingDialog.eta_at = order.value?.eta_at || ''
  shippingDialog.visible = true
}

// 执行状态流转
async function confirmTransition() {
  const payload: any = {
    target_state: transitionDialog.targetState,
    reason: transitionDialog.reason || undefined,
  }
  // shipped 状态时附带物流信息
  if (transitionDialog.targetState === 'shipped') {
    if (transitionDialog.courier) payload.courier = transitionDialog.courier
    if (transitionDialog.tracking_no) payload.tracking_no = transitionDialog.tracking_no
    if (transitionDialog.shipped_at) payload.shipped_at = transitionDialog.shipped_at
    if (transitionDialog.eta_at) payload.eta_at = transitionDialog.eta_at
  }
  try {
    await api.post(`/orders/${id}/state`, payload)
    ElMessage.success(`→ ${STATE_LABELS[transitionDialog.targetState] || transitionDialog.targetState}`)
    transitionDialog.visible = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Transition failed')
  }
}

// 执行独立物流更新
async function saveShipping() {
  const payload: any = {}
  if (shippingDialog.courier !== order.value?.courier) payload.courier = shippingDialog.courier
  if (shippingDialog.tracking_no !== order.value?.tracking_no) payload.tracking_no = shippingDialog.tracking_no
  if (shippingDialog.shipped_at !== order.value?.shipped_at) payload.shipped_at = shippingDialog.shipped_at
  if (shippingDialog.eta_at !== order.value?.eta_at) payload.eta_at = shippingDialog.eta_at
  if (Object.keys(payload).length === 0) {
    shippingDialog.visible = false
    return
  }
  try {
    await api.put(`/orders/${id}/shipping`, payload)
    ElMessage.success('Shipping info updated')
    shippingDialog.visible = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Update failed')
  }
}

// 取消订单
async function cancel() {
  try {
    await ElMessageBox.confirm('Cancel this order? This action cannot be undone.', 'Confirm Cancel', { type: 'warning' })
  } catch { return }
  try { await api.post(`/orders/${id}/cancel`); ElMessage.success('Cancelled'); load() } catch (e: any) { ElMessage.error(e?.message || 'Cancel failed') }
}

async function regenInv() { await api.post(`/orders/${id}/invoice/regenerate`); ElMessage.success('Invoice regenerated'); load() }

// 状态流转按钮列表（按当前状态决定显示什么）
function availableTransitions(currentState: string): string[] {
  const map: Record<string, string[]> = {
    ordering: ['paid', 'cancelled'],
    paid: ['pending_declaration', 'producing', 'cancelled'],
    pending_declaration: ['producing', 'pending_customs', 'cancelled'],
    producing: ['ready_for_production', 'cancelled'],
    ready_for_production: ['pending_declaration', 'pending_customs', 'cancelled'],
    pending_customs: ['customs_clear', 'cancelled'],
    customs_clear: ['shipped', 'cancelled'],
    shipped: ['completed', 'disputed', 'cancelled'],
    disputed: ['cancelled', 'completed'],
  }
  return map[currentState] || []
}

const transitionButtons = computed(() => {
  if (!order.value) return []
  return availableTransitions(order.value.state)
})

function stateBtnClass(target: string) {
  if (target === 'paid' || target === 'completed') return 'success'
  if (target === 'cancelled' || target === 'disputed') return 'danger'
  if (target === 'shipped' || target === 'customs_clear') return 'primary'
  return 'default'
}

// Timeline 事件类型对应的图标/颜色
const timelineIcon: Record<string, string> = {
  order_created: '📦',
  state_change: '🔄',
  shipping: '🚚',
  audit: '📋',
  order_updated: '✏️',
}

function formatTime(iso: string) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleString('en-GB', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div v-if="loading" class="text-center py-20"><el-icon class="is-loading text-4xl"><Loading /></el-icon></div>
  <div v-else-if="order" class="space-y-4">
    <!-- 顶部标题 + steps -->
    <el-card>
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="font-serif text-2xl text-slate-900">{{ order.order_no }}</h2>
          <p class="text-sm text-slate-500">Order ID #{{ order.id }} · Placed {{ order.created_at }}</p>
        </div>
        <div class="flex gap-2">
          <el-button @click="router.back()">← Back</el-button>
          <el-button type="primary" @click="router.push('/orders')">All Orders</el-button>
        </div>
      </div>

      <el-steps :active="stepActive" finish-status="success" align-center>
        <el-step v-for="s in STATE_FLOW" :key="s" :title="STATE_LABELS[s]" />
      </el-steps>
    </el-card>

    <el-row :gutter="16">
      <!-- 左侧：状态操作 + 物流 + Timeline -->
      <el-col :span="14">
        <!-- 状态操作区 -->
        <el-card>
          <h3 class="font-medium mb-3">🔄 Status Actions</h3>
          <el-space wrap>
            <el-button
              v-for="t in transitionButtons"
              :key="t"
              :type="stateBtnClass(t as string) as any"
              @click="openTransition(t)"
            >
              → {{ STATE_LABELS[t] || t }}
            </el-button>
            <el-button
              v-if="order.state !== 'cancelled' && order.state !== 'completed'"
              type="danger"
              @click="cancel"
            >
              Cancel Order
            </el-button>
          </el-space>
        </el-card>

        <!-- 物流信息区 -->
        <el-card class="mt-4">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium">🚚 Shipping & Tracking</h3>
            <el-button size="small" @click="openShippingEdit">✎ Edit Shipping</el-button>
          </div>
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <div class="text-xs text-slate-500 mb-1">Courier</div>
              <div class="font-medium">{{ order.courier || '—' }}</div>
            </div>
            <div>
              <div class="text-xs text-slate-500 mb-1">Tracking No.</div>
              <div class="font-medium font-mono">{{ order.tracking_no || '—' }}</div>
            </div>
            <div>
              <div class="text-xs text-slate-500 mb-1">Shipped At</div>
              <div class="font-medium">{{ order.shipped_at ? formatTime(order.shipped_at) : '—' }}</div>
            </div>
            <div>
              <div class="text-xs text-slate-500 mb-1">ETA</div>
              <div class="font-medium">{{ order.eta_at ? formatTime(order.eta_at) : '—' }}</div>
            </div>
          </div>
        </el-card>

        <!-- Timeline（真实数据） -->
        <el-card class="mt-4">
          <h3 class="font-medium mb-3">📝 Timeline</h3>
          <div v-if="timeline && timeline.events && timeline.events.length > 0" class="space-y-3">
            <div
              v-for="(ev, i) in timeline.events"
              :key="i"
              class="flex gap-3 text-sm border-l-2 border-slate-200 pl-4 py-1"
            >
              <span class="text-lg">{{ timelineIcon[ev.type] || '•' }}</span>
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span class="font-medium capitalize">{{ ev.type.replace(/_/g, ' ') }}</span>
                  <span class="text-xs text-slate-400">{{ formatTime(ev.at) }}</span>
                </div>
                <div class="text-slate-600 mt-1">
                  <!-- state_change 详情 -->
                  <template v-if="ev.type === 'state_change'">
                    <span class="text-slate-400">{{ STATE_LABELS[ev.detail.from_state] || ev.detail.from_state }}</span>
                    <span class="mx-1">→</span>
                    <span class="font-medium">{{ STATE_LABELS[ev.detail.to_state] || ev.detail.to_state }}</span>
                    <div v-if="ev.detail.reason" class="text-slate-500 mt-1 italic">"{{ ev.detail.reason }}"</div>
                  </template>
                  <!-- shipping 详情 -->
                  <template v-else-if="ev.type === 'shipping'">
                    <span v-if="ev.detail.courier">📮 {{ ev.detail.courier }}</span>
                    <span v-if="ev.detail.tracking_no" class="font-mono ml-2">{{ ev.detail.tracking_no }}</span>
                    <span v-if="ev.detail.eta_at" class="ml-2 text-slate-500">ETA {{ formatTime(ev.detail.eta_at) }}</span>
                  </template>
                  <!-- order_created -->
                  <template v-else-if="ev.type === 'order_created'">
                    Order created for <strong>#{{ ev.detail.order_no }}</strong>
                  </template>
                  <!-- audit / order_updated / 其他 -->
                  <template v-else>
                    <span v-if="ev.detail.action">{{ ev.detail.action }}</span>
                    <span v-if="ev.detail.reason" class="italic text-slate-500">"{{ ev.detail.reason }}"</span>
                  </template>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-sm text-slate-400 text-center py-4">No timeline events yet.</div>
        </el-card>

        <!-- 定制产品快照 -->
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

      <!-- 右侧：财务 + 地址 + 发票 -->
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

    <!-- ============ 状态流转弹窗 ============ -->
    <el-dialog v-model="transitionDialog.visible" :title="`Transition to ${STATE_LABELS[transitionDialog.targetState] || transitionDialog.targetState}`" width="520px">
      <el-form label-width="100px" label-position="top">
        <el-form-item label="Reason / Note">
          <el-input
            v-model="transitionDialog.reason"
            type="textarea"
            :rows="3"
            placeholder="e.g. 'FedEx tracking #FX123456789', 'container sealed 2026-09-17', 'SGS report attached'"
          />
        </el-form-item>

        <template v-if="transitionDialog.targetState === 'shipped'">
          <el-divider content-position="left">📦 Shipping Details</el-divider>
          <el-form-item label="Courier">
            <el-select v-model="transitionDialog.courier" placeholder="Select courier" clearable style="width: 100%">
              <el-option label="FedEx" value="FedEx" />
              <el-option label="DHL" value="DHL" />
              <el-option label="UPS" value="UPS" />
              <el-option label="EMS" value="EMS" />
              <el-option label="China Post" value="China Post" />
              <el-option label="Private Courier" value="Private" />
            </el-select>
          </el-form-item>
          <el-form-item label="Tracking No.">
            <el-input v-model="transitionDialog.tracking_no" placeholder="AWB / container number" />
          </el-form-item>
          <el-form-item label="Shipped At (ISO 8601)">
            <el-input v-model="transitionDialog.shipped_at" placeholder="2026-09-17T10:30:00Z" />
          </el-form-item>
          <el-form-item label="ETA (ISO 8601)">
            <el-input v-model="transitionDialog.eta_at" placeholder="2026-09-24T18:00:00Z" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="transitionDialog.visible = false">Cancel</el-button>
        <el-button type="primary" @click="confirmTransition">Confirm Transition</el-button>
      </template>
    </el-dialog>

    <!-- ============ 独立物流编辑弹窗 ============ -->
    <el-dialog v-model="shippingDialog.visible" title="Edit Shipping Information" width="480px">
      <el-form label-position="top">
        <el-form-item label="Courier">
          <el-select v-model="shippingDialog.courier" placeholder="Select courier" clearable style="width: 100%">
            <el-option label="FedEx" value="FedEx" />
            <el-option label="DHL" value="DHL" />
            <el-option label="UPS" value="UPS" />
            <el-option label="EMS" value="EMS" />
            <el-option label="China Post" value="China Post" />
            <el-option label="Private Courier" value="Private" />
          </el-select>
        </el-form-item>
        <el-form-item label="Tracking / Container No.">
          <el-input v-model="shippingDialog.tracking_no" placeholder="e.g. FX123456789" />
        </el-form-item>
        <el-form-item label="Shipped At (ISO 8601)">
          <el-input v-model="shippingDialog.shipped_at" placeholder="2026-09-17T10:30:00Z" />
        </el-form-item>
        <el-form-item label="ETA (ISO 8601)">
          <el-input v-model="shippingDialog.eta_at" placeholder="2026-09-24T18:00:00Z" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shippingDialog.visible = false">Cancel</el-button>
        <el-button type="primary" @click="saveShipping">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>
