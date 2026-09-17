<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'

const router = useRouter()
const user = ref<any>(null)
const orders = ref<any[]>([])
const tab = ref<'orders' | 'addresses' | 'sgs' | 'privacy'>('orders')

async function load() {
  try {
    orders.value = (await api.get('/orders')).items || []
  } catch {}
}
onMounted(load)

function logout() { localStorage.removeItem('user_token'); router.push('/') }
</script>
<template>
  <div class="pt-20 min-h-screen bg-tea-50">
    <div class="max-w-5xl mx-auto px-6 py-12">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="font-serif text-3xl text-tea-900">My Account</h1>
          <p class="text-sm text-tea-600">Welcome back, tea lover 🍵</p>
        </div>
        <button @click="logout" class="text-sm text-tea-700 hover:underline">Sign out</button>
      </div>

      <el-tabs v-model="tab">
        <el-tab-pane label="📦 My Orders" name="orders">
          <div v-if="orders.length===0" class="text-center py-12 text-tea-500">
            <div class="text-4xl mb-3">📭</div>
            No orders yet. <RouterLink to="/bespoke" class="text-tea-700 font-medium">Create your first bespoke →</RouterLink>
          </div>
          <div v-else class="space-y-4">
            <div v-for="o in orders" :key="o.id" class="bg-white border border-tea-100 rounded-2xl p-6 card-hover">
              <div class="flex items-start justify-between mb-3">
                <div>
                  <div class="font-mono text-xs text-tea-500 mb-1">{{ o.order_no }}</div>
                  <div class="font-serif text-lg">{{ o.custom_product_snapshot?.title }}</div>
                  <div class="text-sm text-tea-600">{{ o.custom_product_snapshot?.tea_garden_location }} · {{ o.custom_product_snapshot?.master_name }}</div>
                </div>
                <div class="text-right">
                  <div class="font-serif text-xl">£{{ o.total_amount }}</div>
                  <el-tag size="small" :type="{'ordering':'info','paid':'success','producing':'warning','shipped':'primary','completed':'success','cancelled':'danger'}[o.state]||'info'" effect="dark" class="mt-2">{{ o.state }}</el-tag>
                </div>
              </div>
              <div class="flex gap-3 text-sm">
                <RouterLink :to="`/orders/${o.id}/invoice`" class="text-tea-700 hover:underline">🧾 Invoice</RouterLink>
                <span v-if="o.state==='paid'" class="text-tea-600">· Producing now, shipping within 45 days</span>
                <span v-if="o.state==='producing'" class="text-tea-600">· Your tea is being roasted in Yunnan 🫖</span>
                <span v-if="o.state==='shipped'" class="text-tea-600">· Track your package →</span>
              </div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="🏠 Addresses" name="addresses">
          <div class="bg-white border border-tea-100 rounded-2xl p-8 text-center text-tea-500">Addresses added automatically at checkout. Edit here soon.</div>
        </el-tab-pane>
        <el-tab-pane label="🔬 SGS Certificates" name="sgs">
          <div class="bg-white border border-tea-100 rounded-2xl p-8 text-center text-tea-500">Your batch certificates appear here. Scan the QR on your box for instant access.</div>
        </el-tab-pane>
        <el-tab-pane label="🔒 GDPR / Data" name="privacy">
          <el-card>
            <h3 class="font-medium mb-2">Your Data Rights</h3>
            <ul class="text-sm text-tea-700 space-y-2">
              <li>✓ Access your personal data — <a href="mailto:privacy@ukteahouse.co.uk" class="text-tea-700 underline">request a copy</a></li>
              <li>✓ Correct or update your data at any time in Account → Addresses</li>
              <li>✓ Delete your data — <a href="mailto:privacy@ukteahouse.co.uk" class="text-tea-700 underline">request erasure</a></li>
              <li>✓ Object to or restrict processing — email <a href="mailto:privacy@ukteahouse.co.uk" class="text-tea-700 underline">privacy@ukteahouse.co.uk</a></li>
              <li>✓ Portability — we can export your data in machine-readable format</li>
            </ul>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
