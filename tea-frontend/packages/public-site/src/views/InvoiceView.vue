<script setup lang="ts">
// ============================================================
// InvoiceView.vue — 老钱审美发票页
// - 删除 emoji 🖨️⬇
// - 不用 alert()，用 in-page error
// - ink/ivory/gold palette + squared corners
// ============================================================
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const id = parseInt(route.params.id as string)
const inv = ref<any>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    inv.value = await api.get(`/orders/${id}/invoice`)
  } catch {
    error.value = 'Invoice not found. Please check your account or contact your advisor.'
    setTimeout(() => router.push('/account'), 2000)
  } finally {
    loading.value = false
  }
})

function print() {
  window.print()
}
</script>

<template>
  <div class="pt-20 py-12 bg-ivory-100">
    <div class="max-w-3xl mx-auto px-6">

      <!-- Loading -->
      <div v-if="loading" class="text-center py-20">
        <div class="w-8 h-8 border border-gold/40 border-t-gold rounded-full mx-auto animate-spin"></div>
        <p class="mt-3 text-[11px] uppercase tracking-lux text-sand font-sans">Loading · · ·</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="bg-ink-900 text-ivory-100 p-5" style="border-radius: 2px;">
        <div class="flex items-center gap-2 mb-1">
          <span class="w-1 h-1 rounded-full bg-gold"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Notice</span>
        </div>
        <p class="text-sm text-ivory-100/80 font-serif">{{ error }}</p>
      </div>

      <!-- Invoice -->
      <div v-else-if="inv" class="bg-white border border-gold/20 p-8 md:p-12" style="border-radius: 2px;">
        <!-- Header -->
        <div class="flex items-start justify-between mb-10 pb-6 border-b border-gold/20">
          <div>
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-2">Invoice</div>
            <div class="font-serif text-2xl text-ink-900">UK · Tea · House · Ltd.</div>
            <div class="text-xs text-sand mt-1 font-serif">12 Savile Row, Mayfair, London W1J 5PA · VAT GB123456789</div>
          </div>
          <div class="text-right">
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-1">No.</div>
            <div class="font-mono text-lg text-ink-900">{{ inv.invoice_no }}</div>
            <div class="text-[10px] uppercase tracking-lux text-sand font-sans mt-2">
              Issued {{ new Date(inv.created_at).toLocaleDateString('en-GB') }}
            </div>
          </div>
        </div>

        <!-- Parties -->
        <div class="grid md:grid-cols-2 gap-6 mb-8">
          <div>
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-1">Billed · To</div>
            <div class="text-sm text-ink-900 font-serif">{{ inv.custom_product_snapshot?.title || 'Customer' }}</div>
          </div>
          <div>
            <div class="text-[10px] uppercase tracking-lux text-gold font-sans mb-1">Product</div>
            <div class="text-sm text-ink-900 font-serif">{{ inv.custom_product_snapshot?.title }}</div>
            <div class="text-xs text-sand mt-1 font-serif">
              Tea Garden: {{ inv.custom_product_snapshot?.tea_garden_location }} · Master: {{ inv.custom_product_snapshot?.master_name }}
            </div>
          </div>
        </div>

        <!-- Table -->
        <table class="w-full text-sm font-serif">
          <thead>
            <tr class="border-b border-gold/20">
              <th class="py-3 text-left text-[10px] uppercase tracking-lux text-gold font-sans">Description</th>
              <th class="py-3 text-[10px] uppercase tracking-lux text-gold font-sans">HS · Code</th>
              <th class="py-3 text-right text-[10px] uppercase tracking-lux text-gold font-sans">Amount</th>
            </tr>
          </thead>
          <tbody class="text-sand">
            <tr><td class="py-3">Bespoke Pu'er Tea</td><td>{{ inv.hs_code }}</td><td class="text-right text-ink-900">£{{ inv.total_amount }}</td></tr>
            <tr class="border-t border-gold/10"><td class="py-3">Shipping to UK</td><td>—</td><td class="text-right">Included</td></tr>
          </tbody>
          <tfoot>
            <tr class="border-t border-gold/30">
              <td colspan="2" class="text-right pt-4 text-[10px] uppercase tracking-lux text-gold font-sans">Total</td>
              <td class="text-right pt-4 font-serif text-xl text-ink-900">£{{ inv.total_amount }}</td>
            </tr>
          </tfoot>
        </table>

        <!-- Footer -->
        <div class="mt-10 pt-6 border-t border-gold/20 flex flex-col md:flex-row md:justify-between md:items-center gap-4">
          <div class="text-xs text-sand font-serif leading-relaxed">
            Origin: {{ inv.country_of_origin }} · This document is legally valid for UK customs declarations.
          </div>
          <div class="flex gap-2">
            <button @click="print"
              class="px-5 py-2 border border-gold/30 text-[11px] uppercase tracking-lux text-ink-900 font-sans hover:bg-ivory-50"
              style="border-radius: 2px;">Print</button>
            <button @click="print"
              class="px-5 py-2 bg-ink-900 text-ivory-100 text-[11px] uppercase tracking-lux font-sans hover:bg-ink-800"
              style="border-radius: 2px;">Download · PDF</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
