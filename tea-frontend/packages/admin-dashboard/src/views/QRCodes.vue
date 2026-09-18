<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">QR Code Generator</span>
        </div>
      </template>
      <div class="mb-4 text-sm text-gray-500">
        Select a custom product, then generate a QR code that links to the quote page.
      </div>

      <el-form inline class="mb-4">
        <el-form-item label="Product">
          <el-select v-model="selectedProduct" placeholder="Select product" style="width: 320px" filterable>
            <el-option v-for="p in products" :key="p.id" :label="`#${p.id} · ${p.name_en}`" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :disabled="!selectedProduct" @click="generate">Generate QR</el-button>
        </el-form-item>
      </el-form>

      <div v-if="qrUrl" class="flex gap-6 items-start">
        <div class="border rounded p-4 bg-white">
          <img :src="qrUrl" class="w-48 h-48" />
          <div class="text-xs text-gray-400 mt-2 text-center">Product #{{ selectedProduct?.id }}</div>
        </div>
        <div>
          <div class="mb-2"><strong>URL:</strong></div>
          <div class="text-blue-600 break-all">{{ qrTargetUrl }}</div>
          <el-button class="mt-3" @click="downloadQR">Download PNG</el-button>
        </div>
      </div>
      <el-empty v-else description="Select a product and click Generate" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

const products = ref<any[]>([])
const selectedProduct = ref<any>(null)
const qrUrl = ref('')
const qrTargetUrl = ref('')

async function loadProducts() {
  try {
    const res: any = await api.get('/custom-products')
    products.value = res?.items || res || []
  } catch (e: any) {
    ElMessage.error('Products load failed: ' + (e?.message || e))
  }
}

async function generate() {
  try {
    const token = selectedProduct.value.public_token || selectedProduct.value.id
    const url = `${window.location.origin}/quote/${token}`
    qrTargetUrl.value = url
    const res: any = await api.post('/qrcodes/generate', { url })
    qrUrl.value = res?.data_url || res?.qr_code || res?.image_url || res
  } catch (e: any) {
    ElMessage.error('QR generate failed: ' + (e?.message || e))
  }
}

function downloadQR() {
  if (!qrUrl.value.startsWith('data:')) {
    ElMessage.warning('Only data URLs supported for download')
    return
  }
  const a = document.createElement('a')
  a.href = qrUrl.value
  a.download = `qr-product-${selectedProduct.value.id}.png`
  a.click()
}

onMounted(loadProducts)
</script>
