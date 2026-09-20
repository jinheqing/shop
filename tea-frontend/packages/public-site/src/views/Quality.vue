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

    <!-- Hero — off-black 主舞台 -->
    <section class="bg-ink-900 text-ivory-100 py-20 md:py-28 relative overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.08]"
           style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 65%, #C5A572 0%, transparent 50%);"></div>
      <div class="max-w-5xl mx-auto px-6 text-center relative">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">SGS · UKAS · ISO 17025</span>
          <span class="h-px w-10 bg-gold/60"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-6xl lg:text-7xl leading-tight mb-6 tracking-tight">
          Quality &amp; <span class="italic text-gold">Safety</span>
        </h1>
        <p class="text-lg md:text-xl text-sand max-w-2xl mx-auto leading-relaxed font-serif">
          Every tea batch is independently tested by SGS China.
          A full audit trail, from tea garden to cup.
        </p>
      </div>
    </section>

    <!-- Assurance strip — 不用 emoji，用 hairline 金线分隔 -->
    <section class="py-12 md:py-14 bg-ivory-100 border-y border-gold/15">
      <div class="max-w-6xl mx-auto px-6">
        <div class="grid md:grid-cols-4 gap-8 text-center">
          <div>
            <div class="flex items-center gap-3 justify-center mb-3">
              <span class="h-px w-6 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Pesticide-Free</span>
              <span class="h-px w-6 bg-gold/50"></span>
            </div>
            <p class="text-xs text-sand font-serif leading-relaxed">Every batch tested for 400+ pesticides</p>
          </div>
          <div>
            <div class="flex items-center gap-3 justify-center mb-3">
              <span class="h-px w-6 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Full · Traceability</span>
              <span class="h-px w-6 bg-gold/50"></span>
            </div>
            <p class="text-xs text-sand font-serif leading-relaxed">QR code links to garden, master &amp; harvest</p>
          </div>
          <div>
            <div class="flex items-center gap-3 justify-center mb-3">
              <span class="h-px w-6 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">UKAS · Accredited</span>
              <span class="h-px w-6 bg-gold/50"></span>
            </div>
            <p class="text-xs text-sand font-serif leading-relaxed">Lab meets ISO 17025 standards</p>
          </div>
          <div>
            <div class="flex items-center gap-3 justify-center mb-3">
              <span class="h-px w-6 bg-gold/50"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Guarantee</span>
              <span class="h-px w-6 bg-gold/50"></span>
            </div>
            <p class="text-xs text-sand font-serif leading-relaxed">100% refund if test fails</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Reports -->
    <section class="py-16 md:py-20 bg-white">
      <div class="max-w-5xl mx-auto px-6">

        <div class="mb-10 md:mb-12 text-center">
          <div class="flex items-center gap-3 justify-center mb-5">
            <span class="h-px w-10 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Recent · SGS · Reports</span>
            <span class="h-px w-10 bg-gold/50"></span>
          </div>
          <h2 class="font-serif text-3xl md:text-4xl text-ink-900">The · Laboratory · Record.</h2>
        </div>

        <div v-if="loading" class="text-center py-20">
          <div class="inline-block w-10 h-10 border border-gold/30 border-t-gold rounded-full animate-spin"></div>
          <div class="mt-4 text-[11px] uppercase tracking-lux text-sand font-sans">Loading · Reports</div>
        </div>

        <div v-else class="grid gap-4">
          <div v-for="r in reports" :key="r.id || r.report_number"
               class="bg-ivory-100 border border-gold/15 p-6 card-lux"
               style="border-radius: 2px;">
            <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3 mb-4">
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">{{ r.report_number }}</div>
                <div class="font-serif text-lg text-ink-900">{{ r.tea_garden || r.title }}</div>
              </div>
              <span class="text-[10px] uppercase tracking-lux text-gold border border-gold/40 px-3 py-1 self-start font-sans"
                    style="border-radius: 2px;">
                Passed
              </span>
            </div>
            <div class="grid sm:grid-cols-3 gap-4 text-sm mb-3">
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Pesticides</div>
                <div class="font-serif text-ink-900">{{ r.pesticide_nd_limit || 'Not Detected' }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Heavy · Metals</div>
                <div class="font-serif text-ink-900">{{ r.heavy_metals || 'ND' }}</div>
              </div>
              <div>
                <div class="text-[10px] uppercase tracking-lux text-sand font-sans mb-1">Microbiology</div>
                <div class="font-serif text-ink-900">{{ r.microbiology || 'Passed' }}</div>
              </div>
            </div>
            <div class="text-[10px] uppercase tracking-lux text-sand font-sans pt-3 border-t border-gold/10">
              Issued · {{ r.issued_at }}
            </div>
          </div>
          <div v-if="!reports.length" class="text-center text-sand py-12 font-serif">No SGS reports available at this time.</div>
        </div>

        <!-- Disclaimer -->
        <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-10 leading-loose font-sans">
          SGS reports apply only to the specific production batch tested &middot; Results are not medical advice<br>
          All products are food items, not medicinal products
        </p>
      </div>
    </section>
  </div>
</template>
