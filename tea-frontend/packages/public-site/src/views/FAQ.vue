<script setup lang="ts">
const faqs = [
  { q: 'How long does a bespoke order take?', a: '45 days from confirmation. Harvest → Roast → Pack → Ship. UK delivery takes 3-5 days after customs clear.' },
  { q: 'Can I visit the tea farm?', a: 'Yes! We run 2-week escorted trips to Yunnan twice a year. Email hello@ukteahouse.co.uk for the next date.' },
  { q: 'Is my tea traceable?', a: 'Absolutely. Scan the QR on every box for tea garden village, master profile, harvest date, and a 24/7 live camera of the tea garden.' },
  { q: 'Do you ship outside the UK?', a: 'Currently UK, EU, USA, Australia. International shipping from £150. Customs duties may apply.' },
  { q: 'How is payment protected?', a: 'We use 2Checkout and PayPal — both PCI-DSS compliant. We never store your card number.' },
  { q: 'Can I return tea?', a: 'Bespoke products are non-refundable (see terms). Standard stock can be returned unopened within 14 days.' },
  { q: 'Are your farms organic?', a: 'All six farms practice natural farming — zero pesticides since 2018. Independently audited annually.' },
  { q: 'Do you do corporate orders?', a: 'Yes — custom-branded wooden boxes, gift certificates, and volume pricing. Contact gifting@ukteahouse.co.uk.' },
]
</script>
<template>
  <div class="pt-20 py-20 bg-tea-50">
    <div class="max-w-3xl mx-auto px-6">
      <h1 class="font-serif text-4xl text-tea-900 mb-4 text-center">Frequently Asked</h1>
      <p class="text-tea-600 text-center mb-12">Can't find your answer? <RouterLink to="/contact" class="text-tea-700 underline">Contact us</RouterLink>.</p>

      <el-collapse class="bg-white rounded-2xl border border-tea-100 overflow-hidden">
        <el-collapse-item v-for="(f,i) in faqs" :key="i" :title="f.q" :name="i">
          <div class="text-tea-700 leading-relaxed">{{ f.a }}</div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>
