<script setup lang="ts">
// ============================================================
// FAQ.vue — 叙述式长文（Old Money Edition）
// ============================================================
// 老钱审美：折叠 +号 是 SaaS 风格，奢侈品 FAQ 是
// 分类用 serif 小标题 + 叙述式长文段落
// 参考：Hermès FAQ / Rolls-Royce Private Office
// ============================================================

import { RouterLink } from 'vue-router'

type QA = { q: string; a: string }
type Section = { title: string; items: QA[] }

const sections: Section[] = [
  {
    title: 'Notes · On · Bespoke',
    items: [
      { q: 'How long does a bespoke commission take?',
        a: 'Forty-five days from confirmation. The batch moves through harvest, roasting, packing, and shipping. UK delivery takes a further three to five days after customs clearance.' },
      { q: 'May I visit the tea farm?',
        a: 'Yes — we run two escorted journeys to Yunnan each year, in spring and autumn. Write to hello@ukteahouse.co.uk for the next available dates and itinerary.' },
      { q: 'Can I return or cancel a bespoke blend?',
        a: 'Bespoke blends, made to your specification, cannot be cancelled or returned. Non-bespoke items may be returned within fourteen calendar days of delivery, unopened, in accordance with the Consumer Rights Act 2015. See our <RouterLink to="/terms" class="text-gold hover:text-gold-soft">Terms of Service</RouterLink>.' },
      { q: 'Do you accept corporate commissions?',
        a: 'Yes. Custom-branded wooden boxes, gift certificates, and volume pricing are available. Write to gifting@ukteahouse.co.uk.' },
    ],
  },
  {
    title: 'Notes · On · Provenance · & · Quality',
    items: [
      { q: 'Is my tea traceable?',
        a: 'Always. Scan the QR on every box for the tea garden village, master profile, harvest date, and a twenty-four-hour live camera of the tea garden. The traceability record reflects the production batch, not every individual cup.' },
      { q: 'Are your farms organic?',
        a: 'All six farms practise natural farming — zero pesticides since 2018. They are independently audited each year, and SGS certificates are available on request.' },
      { q: 'Is tea from Yunnan safe to drink?',
        a: 'Every batch is tested by SGS China for pesticides, heavy metals, and microbiology. Live streams show where the tea grew. Tea is a food item, not a medicinal product — test results are not medical advice.' },
    ],
  },
  {
    title: 'Notes · On · Delivery · & · Payment',
    items: [
      { q: 'Do you ship outside the United Kingdom?',
        a: 'Currently the United Kingdom, the European Union, the United States, and Australia. International shipping begins at £150. Customs duties are the responsibility of the recipient.' },
      { q: 'How is my payment protected?',
        a: 'We use 2Checkout and PayPal — both PCI-DSS Level 1 compliant. We never see or store your card details. 256-bit SSL throughout.' },
    ],
  },
  {
    title: 'Notes · On · Privacy · & · The · Law',
    items: [
      { q: 'Are you registered with the ICO?',
        a: 'Registration with the UK Information Commissioner\'s Office is pending. Our data processing follows GDPR and UK GDPR principles. Our Data Protection Officer can be reached at privacy@ukteahouse.co.uk.' },
      { q: 'What are my GDPR rights?',
        a: 'You may access, correct, delete, port, restrict, or object to the processing of your data. Exercise these rights via our <RouterLink to="/gdpr-dsar" class="text-gold hover:text-gold-soft">Data Subject Access Request page</RouterLink>. We respond within thirty calendar days of receipt.' },
      { q: 'Do you store my payment card?',
        a: 'No. Payments are routed exclusively through 2Checkout and PayPal, both PCI-DSS Level 1 certified. We receive only an order reference.' },
      { q: 'How long do you keep my data?',
        a: 'Orders: seven years (UK VAT). Staff audit logs: ten years (GDPR accountability). Your account: until you delete it, plus three years post-deletion for tax records. Cookie consent: five years.' },
      { q: 'Must I be eighteen to order?',
        a: 'Yes. By placing an order you confirm you are eighteen or older. Tea products are food items, not medicinal products — not intended for children.' },
      { q: 'Where can I read the legal matters in full?',
        a: '<RouterLink to="/privacy" class="text-gold hover:text-gold-soft">Privacy &amp; GDPR</RouterLink> · <RouterLink to="/terms" class="text-gold hover:text-gold-soft">Terms of Service</RouterLink> · <RouterLink to="/cookie-policy" class="text-gold hover:text-gold-soft">Cookie Policy</RouterLink> · <RouterLink to="/gdpr-dsar" class="text-gold hover:text-gold-soft">Data Subject Requests</RouterLink>' },
    ],
  },
]
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100">
    <div class="max-w-3xl mx-auto px-6 py-16 md:py-20">

      <!-- Section title -->
      <div class="mb-12 md:mb-16 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Notes · On · The · House</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-3">Frequently · Asked.</h1>
        <p class="text-sand text-sm font-serif">
          Should your enquiry not be answered below, please
          <RouterLink to="/contact" class="text-gold hover:text-gold-soft">write to the house</RouterLink>.
        </p>
      </div>

      <!-- Sections — 叙述式长文，不用折叠 +号 -->
      <div class="space-y-14">

        <section v-for="(sec, si) in sections" :key="si">
          <!-- Section title -->
          <div class="mb-7">
            <div class="flex items-center gap-3 mb-4">
              <span class="h-px w-8 bg-gold/60"></span>
              <span class="text-[10px] uppercase tracking-lux text-gold font-sans">{{ sec.title }}</span>
            </div>
          </div>

          <!-- Questions -->
          <div class="space-y-7">
            <div v-for="(qa, qi) in sec.items" :key="qi">
              <h3 class="font-serif text-xl text-ink-900 mb-2 leading-snug">{{ qa.q }}</h3>
              <p class="text-sand text-sm leading-relaxed font-serif" v-html="qa.a"></p>
            </div>
          </div>
        </section>
      </div>

      <!-- Closing — 极细 disclaimer -->
      <div class="hairline-gold my-12"></div>
      <p class="text-[10px] uppercase tracking-lux text-sand text-center leading-loose font-sans">
        All answers reflect our standard practice &middot; Specific commissions may vary<br>
        Tea is a food item, not a medicinal product &middot; SGS results are not medical advice
      </p>
    </div>
  </div>
</template>
