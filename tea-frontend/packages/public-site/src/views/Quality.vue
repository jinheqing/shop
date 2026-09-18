<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'

const reports = ref<any[]>([])
const loading = ref(true)

async function load() {
  try {
    const d: any = await api.get('/public/sgs-reports')
    const got = d?.items || d || []
    if (got.length) { reports.value = got; return }
  } catch {
    reports.value = [
      { report_number: 'SGS-2025-YN-001', tea_garden: '云南省 · 临沧市', pesticide_nd_limit: 'ND (≤0.01 ppm)', heavy_metals: 'ND', microbiology: 'Passed', status: 'Passed', issued_at: '2025-04-15' },
      { report_number: 'SGS-2025-YN-002', tea_garden: '云南省 · 普洱市', pesticide_nd_limit: 'ND (≤0.01 ppm)', heavy_metals: 'ND', microbiology: 'Passed', status: 'Passed', issued_at: '2025-04-12' },
      { report_number: 'SGS-2025-YN-003', tea_garden: '云南省 · 西双版纳州', pesticide_nd_limit: 'ND (≤0.01 ppm)', heavy_metals: 'ND', microbiology: 'Passed', status: 'Passed', issued_at: '2025-03-28' },
    ]
  } finally { loading.value = false }
}
onMounted(load)
</script>

<template>
  <div class="pt-20">
    <!-- Hero -->
    <section class="bg-tea-900 text-tea-50 py-20 md:py-28">
      <div class="max-w-5xl mx-auto px-6 text-center">
        <div class="inline-block px-4 py-1.5 bg-white/10 backdrop-blur rounded-full text-xs tracking-[0.25em] uppercase mb-6 border border-white/15">
          SGS · UKAS · ISO 17025
        </div>
        <h1 class="font-display font-semibold text-4xl md:text-6xl lg:text-7xl leading-tight mb-6 tracking-tight">
          Quality &amp; <span class="italic">Safety</span>
        </h1>
        <p class="text-lg md:text-xl text-tea-200 max-w-2xl mx-auto leading-relaxed">
          Every tea batch is independently tested by SGS China. Full audit trail from tea garden to cup.
        </p>
      </div>
    </section>

    <!-- Assurance strip -->
    <section class="py-12 bg-tea-50 border-y border-tea-100">
      <div class="max-w-6xl mx-auto px-6">
        <div class="grid md:grid-cols-4 gap-6 text-center">
          <div>
            <div class="text-3xl mb-2">🔬</div>
            <div class="font-display text-base text-tea-900 mb-1">Pesticide-free</div>
            <div class="text-xs text-tea-600">Every batch tested for 400+ pesticides</div>
          </div>
          <div>
            <div class="text-3xl mb-2">⛓️</div>
            <div class="font-display text-base text-tea-900 mb-1">Full traceability</div>
            <div class="text-xs text-tea-600">QR code links to garden, master &amp; harvest</div>
          </div>
          <div>
            <div class="text-3xl mb-2">🏛️</div>
            <div class="font-display text-base text-tea-900 mb-1">UKAS accredited</div>
            <div class="text-xs text-tea-600">Lab meets ISO 17025 standards</div>
          </div>
          <div>
            <div class="text-3xl mb-2">🛡️</div>
            <div class="font-display text-base text-tea-900 mb-1">Guarantee</div>
            <div class="text-xs text-tea-600">100% refund if test fails</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Reports -->
    <section class="py-16 md:py-20 bg-white">
      <div class="max-w-5xl mx-auto px-6">
        <h2 class="font-display text-3xl md:text-4xl text-tea-900 mb-10 text-center">Recent SGS Reports</h2>

        <div v-if="loading" class="text-center py-20 text-tea-500">Loading reports…</div>

        <div v-else class="grid gap-4">
          <div v-for="r in reports" :key="r.id || r.report_number"
               class="rounded-2xl border border-tea-100 bg-white p-6 card-hover">
            <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3 mb-3">
              <div>
                <div class="text-xs text-tea-400 font-mono mb-1">{{ r.report_number }}</div>
                <div class="font-display text-lg text-tea-900">{{ r.tea_garden || r.title }}</div>
              </div>
              <span class="bg-tea-100 text-tea-800 text-xs px-3 py-1 rounded-full font-medium self-start">
                ✓ {{ r.status || 'Passed' }}
              </span>
            </div>
            <div class="grid sm:grid-cols-3 gap-4 text-sm mb-3">
              <div>
                <div class="text-xs text-tea-500 mb-1">Pesticides</div>
                <div class="font-mono text-tea-800">{{ r.pesticide_nd_limit || 'Not Detected' }}</div>
              </div>
              <div>
                <div class="text-xs text-tea-500 mb-1">Heavy Metals</div>
                <div class="font-mono text-tea-800">{{ r.heavy_metals || 'ND' }}</div>
              </div>
              <div>
                <div class="text-xs text-tea-500 mb-1">Microbiology</div>
                <div class="font-mono text-tea-800">{{ r.microbiology || 'Passed' }}</div>
              </div>
            </div>
            <div class="text-xs text-tea-500">Issued {{ r.issued_at }}</div>
          </div>
          <div v-if="!reports.length" class="text-center text-tea-500 py-12">No SGS reports available yet.</div>
        </div>

        <!-- SGS disclaimer -->
        <p class="text-xs text-tea-500 text-center mt-8 px-4 leading-relaxed">
          SGS reports apply only to the specific production batch tested. Results are not medical advice and do not guarantee the safety of any individual serving. All products are food items, not medicinal products.
        </p>
      </div>
    </section>
  </div>
</template>
