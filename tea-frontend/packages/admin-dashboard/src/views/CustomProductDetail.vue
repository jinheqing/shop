<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const id = parseInt(route.params.id as string)
const product = ref<any>(null)
const loading = ref(true)

async function load() {
  loading.value = true
  try { product.value = await api.get(`/custom-products/${id}`) }
  finally { loading.value = false }
}
onMounted(load)

async function publish() { await api.post(`/custom-products/${id}/publish`); ElMessage.success('Published'); load() }
async function approve() { await api.post(`/custom-products/${id}/review`, { approved: true }); ElMessage.success('Approved'); load() }
async function reject() { await api.post(`/custom-products/${id}/review`, { approved: false }); ElMessage.success('Rejected'); load() }
</script>
<template>
  <div v-if="loading" class="text-center py-20"><el-icon class="is-loading text-4xl"><Loading /></el-icon></div>
  <div v-else-if="product" class="space-y-4">
    <el-card>
      <div class="flex items-center justify-between mb-2">
        <div>
          <h2 class="font-serif text-2xl">{{ product.title }}</h2>
          <p class="text-sm text-slate-500">SKU: {{ product.sku }} · Version: {{ product.version }} · Status: <el-tag :type="product.status==='published'?'success':'info'">{{ product.status }}</el-tag></p>
        </div>
        <div class="flex gap-2">
          <el-button @click="router.back()">← Back</el-button>
          <el-button v-if="product.status!=='published'" type="primary" @click="publish">🚀 Publish →</el-button>
        </div>
      </div>
    </el-card>

    <el-row :gutter="16">
      <el-col :span="14">
        <el-card>
          <h3 class="font-medium mb-3">📋 Full Bespoke Fields (26+)</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="Tea Type">{{ product.tea_type }}</el-descriptions-item>
            <el-descriptions-item label="Shape">{{ product.tea_shape }}</el-descriptions-item>
            <el-descriptions-item label="Tea Garden">{{ product.tea_garden_location }}</el-descriptions-item>
            <el-descriptions-item label="Master">{{ product.master_name }}</el-descriptions-item>
            <el-descriptions-item label="Source">{{ product.raw_tea_source }}</el-descriptions-item>
            <el-descriptions-item label="Inner Pack">{{ product.inner_packaging }}</el-descriptions-item>
            <el-descriptions-item label="Outer Pack">{{ product.outer_packaging }}</el-descriptions-item>
            <el-descriptions-item label="QR Position">{{ product.qr_code_position }}</el-descriptions-item>
            <el-descriptions-item label="Harvest">{{ product.harvest_date }}</el-descriptions-item>
            <el-descriptions-item label="Roast">{{ product.roasting_date }}</el-descriptions-item>
            <el-descriptions-item label="Lead Time">{{ product.lead_time }}</el-descriptions-item>
            <el-descriptions-item label="Storage">{{ product.storage_location }}</el-descriptions-item>
            <el-descriptions-item label="Custom Req" :span="2">{{ product.custom_requirement }}</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card class="mt-4">
          <h3 class="font-medium mb-3">🎛️ Staff Review</h3>
          <el-alert v-if="product.status==='pending_review'" title="This bespoke needs staff review before publishing." type="warning" show-icon class="mb-3" />
          <div class="flex gap-2">
            <el-button type="success" @click="approve">✅ Approve</el-button>
            <el-button type="danger" @click="reject">❌ Reject</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card>
          <h3 class="font-medium mb-2">💰 Pricing</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between"><span>Unit Price</span><span>£{{ product.unit_price }}</span></div>
            <div class="flex justify-between"><span>Quantity</span><span>{{ product.quantity }}</span></div>
            <div class="flex justify-between"><span>Shipping</span><span>£{{ product.shipping_cost }}</span></div>
            <el-divider class="my-2" />
            <div class="flex justify-between font-serif text-lg"><span>Total</span><span>£{{ product.total_amount }}</span></div>
          </div>
        </el-card>

        <el-card class="mt-4">
          <h3 class="font-medium mb-2">🔗 Public Share Link</h3>
          <div v-if="product.product_token" class="break-all">
            <code class="text-xs">/bespoke/{{ product.product_token }}</code>
            <el-button size="small" class="ml-2">📋 Copy</el-button>
          </div>
          <div v-else class="text-xs text-slate-400">Publish to generate token</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
