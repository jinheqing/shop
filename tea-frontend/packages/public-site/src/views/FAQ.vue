<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'

type FAQ = { q: string; a: string }
const faqs: FAQ[] = [
  { q: 'How long does a bespoke order take?', a: '45 days from confirmation. Harvest → Roast → Pack → Ship. UK delivery takes 3–5 days after customs clearance.' },
  { q: 'Can I visit the tea farm?', a: 'Yes — we run 2-week escorted trips to Yunnan twice a year. Email hello@ukteahouse.co.uk for the next available dates.' },
  { q: 'Is my tea traceable?', a: 'Absolutely. Scan the QR on every box for tea garden village, master profile, harvest date, and a 24/7 live camera of the tea garden. Traceability reflects the production batch, not every individual cup.' },
  { q: 'Do you ship outside the UK?', a: 'Currently UK, EU, USA, and Australia. International shipping from £150. Customs duties are your responsibility.' },
  { q: 'How is payment protected?', a: 'We use 2Checkout and PayPal — both PCI-DSS Level 1 compliant. We never see or store your card details. 256-bit SSL throughout.' },
  { q: 'Can I return or cancel?', a: 'You have 14 calendar days from delivery for unopened, non-bespoke products (Consumer Rights Act 2015). Bespoke blends made to your specification cannot be cancelled or returned. See our <a href="/terms" class="underline hover:text-tea-950">Terms of Service</a>.' },
  { q: 'Are your farms organic?', a: 'All six farms practice natural farming — zero pesticides since 2018. Independently audited annually. SGS certificates available on request.' },
  { q: 'Do you do corporate orders?', a: 'Yes — custom-branded wooden boxes, gift certificates, and volume pricing. Contact gifting@ukteahouse.co.uk.' },

  // --- Compliance / Legal section ---
  { q: 'Are you registered with the ICO?', a: 'Registration reference with the UK Information Commissioner\'s Office (ICO) is pending. Our data processing follows GDPR and UK GDPR principles. Our Data Protection Officer can be reached at privacy@ukteahouse.co.uk.' },
  { q: 'What are my GDPR rights?', a: 'You can access, correct, delete, port, restrict, or object to the processing of your data. Exercise these rights via our <a href="/gdpr-dsar" class="underline hover:text-tea-950">Data Subject Access Request page</a>. SLA: 30 calendar days from receipt.' },
  { q: 'Do you store my payment card?', a: 'No. We never store or process your card details directly. Payments are routed exclusively through 2Checkout and PayPal, both PCI-DSS Level 1 certified. We receive only an order reference.' },
  { q: 'How long do you keep my data?', a: 'Orders: 7 years (UK VAT). Staff audit logs: 10 years (GDPR accountability). Your account: until you delete it, plus 3 years post-deletion for tax records. Cookie consent: 5 years.' },
  { q: 'Do I need to be 18 to order?', a: 'Yes. By placing an order you confirm you are 18 or older. Tea products are food items, not medicinal products — not intended for children.' },
  { q: 'Is tea from Yunnan safe to drink?', a: 'All batches are tested by SGS China for pesticides, heavy metals, and microbiology. Our live streams show where the tea grew. However, tea is not a medicinal product and results are not medical advice.' },
  { q: 'Where can I read the legal stuff?', a: '<a href="/privacy" class="underline hover:text-tea-950">Privacy & GDPR</a> · <a href="/terms" class="underline hover:text-tea-950">Terms of Service</a> · <a href="/cookie-policy" class="underline hover:text-tea-950">Cookie Policy</a> · <a href="/gdpr-dsar" class="underline hover:text-tea-950">Data Subject Requests</a>' },
]

const open = ref<number | null>(0)
function toggle(i: number) { open.value = open.value === i ? null : i }
</script>

<template>
  <div class="pt-20 py-20 bg-tea-50">
    <div class="max-w-3xl mx-auto px-6">
      <div class="text-center mb-12">
        <h1 class="font-display font-semibold text-4xl md:text-5xl text-tea-900 mb-4">Frequently <span class="italic">Asked</span></h1>
        <p class="text-tea-600">Can't find your answer? <RouterLink to="/contact" class="text-tea-700 underline hover:text-tea-950">Contact us</RouterLink>.</p>
      </div>

      <div class="rounded-2xl border border-tea-100 bg-white overflow-hidden divide-y divide-tea-100">
        <button v-for="(f, i) in faqs" :key="i"
          @click="toggle(i)"
          class="w-full text-left px-6 py-5 transition hover:bg-tea-50/60"
          :class="open === i ? 'bg-tea-50/80' : ''">
          <div class="flex items-start justify-between gap-4">
            <span class="font-medium text-tea-900">{{ f.q }}</span>
            <span class="text-tea-400 shrink-0 text-xl leading-none transition-transform duration-200"
                  :class="open === i ? 'rotate-45' : ''">+</span>
          </div>
          <div v-if="open === i" class="mt-3 text-sm text-tea-700 leading-relaxed prose prose-sm max-w-none"
               v-html="f.a"></div>
        </button>
      </div>
    </div>
  </div>
</template>
