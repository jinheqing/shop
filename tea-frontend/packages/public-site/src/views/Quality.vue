<template>
  <div class="pt-20">
    <section class="bg-tea-900 text-tea-50 py-24">
      <div class="max-w-5xl mx-auto px-6">
        <div class="flex items-center gap-4 mb-6">
          <span class="text-5xl">🔬</span>
          <span class="font-serif text-5xl">SGS Quality Assurance.</span>
        </div>
        <p class="text-lg text-tea-200 max-w-3xl">Every single batch we ship to the UK passes independent third-party testing by SGS China. Your peace of mind comes with every tin.</p>
      </div>
    </section>

    <section class="py-20 bg-white">
      <div class="max-w-6xl mx-auto px-6 grid md:grid-cols-2 gap-12">
        <div>
          <h3 class="font-serif text-2xl text-tea-900 mb-4">What We Test For</h3>
          <ul class="space-y-3 text-tea-700">
            <li class="flex items-start gap-3">
              <span class="text-green-600 mt-0.5">✓</span>
              <div><strong>Pesticide Residues</strong> — 280+ pesticides tested, 0 detected on every batch</div>
            </li>
            <li class="flex items-start gap-3">
              <span class="text-green-600 mt-0.5">✓</span>
              <div><strong>Heavy Metals</strong> — Lead, arsenic, cadmium within EU MRL limits</div>
            </li>
            <li class="flex items-start gap-3">
              <span class="text-green-600 mt-0.5">✓</span>
              <div><strong>Microbiology</strong> — E. coli, salmonella, yeast, mold — all within food grade</div>
            </li>
            <li class="flex items-start gap-3">
              <span class="text-green-600 mt-0.5">✓</span>
              <div><strong>Chemical Authenticity</strong> — HPLC 验证真普洱，无掺假</div>
            </li>
            <li class="flex items-start gap-3">
              <span class="text-green-600 mt-0.5">✓</span>
              <div><strong>Water Activity</strong> — a<sub>w</sub> ≤ 0.75，确保存储安全</div>
            </li>
          </ul>
        </div>
        <div class="bg-tea-50 rounded-2xl p-8">
          <h3 class="font-serif text-2xl text-tea-900 mb-4">Scanning Your Tea Box</h3>
          <div class="space-y-4 text-sm text-tea-700">
            <div class="flex gap-3">
              <span class="w-8 h-8 bg-tea-700 text-white rounded-full flex items-center justify-center text-xs flex-shrink-0">1</span>
              <div>Scan the QR on your bespoke box with any phone camera</div>
            </div>
            <div class="flex gap-3">
              <span class="w-8 h-8 bg-tea-700 text-white rounded-full flex items-center justify-center text-xs flex-shrink-0">2</span>
              <div>See SGS report, tea garden village, master profile</div>
            </div>
            <div class="flex gap-3">
              <span class="w-8 h-8 bg-tea-700 text-white rounded-full flex items-center justify-center text-xs flex-shrink-0">3</span>
              <div>Stream the slow-live camera of where it grew — 24/7</div>
            </div>
            <div class="flex gap-3">
              <span class="w-8 h-8 bg-tea-700 text-white rounded-full flex items-center justify-center text-xs flex-shrink-0">4</span>
              <div>See its journey: harvest → roast → pack → ship to your door</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20 bg-tea-50">
      <div class="max-w-6xl mx-auto px-6">
        <h3 class="font-serif text-3xl text-tea-900 mb-10 text-center">Recent SGS Reports</h3>
        <div class="grid md:grid-cols-3 gap-6">
          <div v-for="(r,i) in [
            {no:'SGS-2024-001', batch:'曼岗村-2024春', date:'2024-05-20', result:'PASS'},
            {no:'SGS-2024-002', batch:'班章-2024春', date:'2024-05-25', result:'PASS'},
            {no:'SGS-2024-003', batch:'景迈-2024春', date:'2024-06-01', result:'PASS'},
          ]" :key="i" class="bg-white rounded-xl p-6 border border-tea-100 card-hover">
            <div class="flex items-center justify-between mb-4">
              <span class="font-mono text-xs text-tea-500">{{ r.no }}</span>
              <span class="text-xs px-2 py-1 rounded-full bg-green-100 text-green-700">{{ r.result }}</span>
            </div>
            <div class="font-serif text-lg text-tea-900 mb-1">{{ r.batch }}</div>
            <div class="text-xs text-tea-500">Issued {{ r.date }}</div>
            <button class="mt-4 text-sm text-tea-700 hover:text-tea-900">View PDF →</button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
