<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/api/client'
const route = useRoute()
const router = useRouter()
const id = parseInt(route.params.id as string)
const inv = ref<any>(null)
const loading = ref(true)

onMounted(async () => {
  try { inv.value = await api.get(`/orders/${id}/invoice`) }
  catch { alert('Invoice not found'); router.push('/account') }
  finally { loading.value = false }
})
</script>
<template>
  <div class="pt-20 py-12 bg-tea-50">
    <div class="max-w-3xl mx-auto px-6">
      <div v-if="loading" class="text-center py-20"><div class="animate-spin w-8 h-8 border-2 border-tea-600 border-t-transparent rounded-full mx-auto"></div></div>
      <div v-else-if="inv" class="bg-white rounded-2xl shadow-sm border border-tea-100 p-10">
        <div class="flex items-start justify-between mb-10 pb-6 border-b border-tea-200">
          <div>
            <div class="text-xs text-tea-500 tracking-widest uppercase mb-2">Invoice</div>
            <div class="font-serif text-2xl text-tea-900">UK Tea House Ltd.</div>
            <div class="text-xs text-tea-600 mt-1">Mayfair, London W1J · VAT: GB123456789</div>
          </div>
          <div class="text-right">
            <div class="text-xs text-tea-500 mb-1">No.</div>
            <div class="font-mono text-lg">{{ inv.invoice_no }}</div>
            <div class="text-xs text-tea-500 mt-2">Issued {{ new Date(inv.created_at).toLocaleDateString() }}</div>
          </div>
        </div>
        <div class="grid md:grid-cols-2 gap-6 mb-8">
          <div>
            <div class="text-xs text-tea-500 mb-1">Billed To</div>
            <div class="text-sm text-tea-800">{{ inv.custom_product_snapshot?.title || 'Customer' }}</div>
          </div>
          <div>
            <div class="text-xs text-tea-500 mb-1">Product</div>
            <div class="text-sm text-tea-800">{{ inv.custom_product_snapshot?.title }}</div>
            <div class="text-xs text-tea-500 mt-1">Mountain: {{ inv.custom_product_snapshot?.tea_garden_location }} · Master: {{ inv.custom_product_snapshot?.master_name }}</div>
          </div>
        </div>
        <table class="w-full text-sm">
          <thead><tr class="border-b border-tea-200 text-left text-tea-500"><th class="py-2">Description</th><th>HS Code</th><th class="text-right">Amount</th></tr></thead>
          <tbody>
            <tr><td class="py-3">Bespoke Pu'er Tea</td><td>{{ inv.hs_code }}</td><td class="text-right">£{{ inv.total_amount }}</td></tr>
            <tr class="border-b border-tea-100"><td class="py-3">Shipping to UK</td><td>—</td><td class="text-right">Included</td></tr>
          </tbody>
          <tfoot><tr><td colspan="2" class="text-right pt-4 font-medium">Total</td><td class="text-right pt-4 font-serif text-xl">£{{ inv.total_amount }}</td></tr></tfoot>
        </table>
        <div class="mt-10 pt-6 border-t border-tea-200 flex justify-between items-center">
          <div class="text-xs text-tea-500">Origin: {{ inv.country_of_origin }} · This document is legally valid for UK customs declarations.</div>
          <div class="flex gap-2">
            <button class="px-4 py-2 border border-tea-300 rounded-lg text-sm hover:bg-tea-50">🖨️ Print</button>
            <button class="px-4 py-2 bg-tea-800 text-white rounded-lg text-sm hover:bg-tea-900">⬇ Download PDF</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
