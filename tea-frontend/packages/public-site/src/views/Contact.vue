<script setup lang="ts">
// ============================================================
// Contact.vue — 老钱审美联系页（Old Money Edition）
// ============================================================
// - 不放 WhatsApp（熟人通讯工具，不符奢侈品定位）
// - 改 "By Appointment" 静默联系
// - 不用 emoji，用 hairline 金线分隔
// - 不用 ElMessage toast，用 in-page 状态条
// ============================================================

import { ref } from 'vue'

const form = ref({ name: '', email: '', subject: 'General', message: '' })
const sent = ref(false)
const errorMsg = ref('')

function submit() {
  errorMsg.value = ''
  if (!form.value.name || !form.value.email || !form.value.message) {
    errorMsg.value = 'Please complete all fields, or write to hello@ukteahouse.co.uk directly.'
    return
  }
  // 此处暂用本地确认（未来可对接 /public/contact-inquiries）
  sent.value = true
  form.value = { name:'', email:'', subject:'General', message:'' }
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100">
    <div class="max-w-3xl mx-auto px-6 py-16 md:py-20">

      <!-- Section title -->
      <div class="mb-12 md:mb-14 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">By · Appointment</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-3">Speak · With · The · House.</h1>
        <p class="text-sand text-sm font-serif max-w-xl mx-auto leading-relaxed">
          An advisor will respond within twenty-four hours, Monday through Saturday.
          For private commissions or large orders, we recommend a brief telephone conversation.
        </p>
      </div>

      <!-- Three contact channels — 不用 emoji，用 hairline 金线分隔 -->
      <div class="grid md:grid-cols-3 gap-6 mb-12">
        <div class="bg-white border border-gold/15 p-6 text-center" style="border-radius: 2px;">
          <div class="flex items-center gap-3 justify-center mb-4">
            <span class="h-px w-6 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">By · Letter</span>
            <span class="h-px w-6 bg-gold/50"></span>
          </div>
          <div class="font-serif text-ink-900 mb-1">Email</div>
          <div class="text-[11px] uppercase tracking-lux text-sand font-sans">hello@ukteahouse.co.uk</div>
        </div>

        <div class="bg-white border border-gold/15 p-6 text-center" style="border-radius: 2px;">
          <div class="flex items-center gap-3 justify-center mb-4">
            <span class="h-px w-6 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">By · Telephone</span>
            <span class="h-px w-6 bg-gold/50"></span>
          </div>
          <div class="font-serif text-ink-900 mb-1">London Office</div>
          <div class="text-[11px] uppercase tracking-lux text-sand font-sans">+44 20 7946 0958</div>
        </div>

        <div class="bg-white border border-gold/15 p-6 text-center" style="border-radius: 2px;">
          <div class="flex items-center gap-3 justify-center mb-4">
            <span class="h-px w-6 bg-gold/50"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">By · Appointment</span>
            <span class="h-px w-6 bg-gold/50"></span>
          </div>
          <div class="font-serif text-ink-900 mb-1">Mayfair · London</div>
          <div class="text-[11px] uppercase tracking-lux text-sand font-sans">12 Savile Row · W1J 5PA</div>
        </div>
      </div>

      <!-- Success state -->
      <div v-if="sent" class="bg-white border border-gold/20 p-10 text-center" style="border-radius: 2px;">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Received</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h3 class="font-serif text-2xl text-ink-900 mb-3">Thank you for writing.</h3>
        <p class="text-sand text-sm font-serif">An advisor will respond within twenty-four hours.</p>
      </div>

      <!-- Contact form -->
      <form v-else @submit.prevent="submit" class="bg-white border border-gold/15 p-8 md:p-10" style="border-radius: 2px;">

        <div v-if="errorMsg" class="mb-6 p-4 bg-ink-900 text-ivory-100 text-sm font-serif" style="border-radius: 2px;">
          <span class="text-gold text-[10px] uppercase tracking-lux font-sans">Attention</span><br>
          {{ errorMsg }}
        </div>

        <div class="grid md:grid-cols-2 gap-4 mb-6">
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Your · Name</label>
            <input v-model="form.name" placeholder="Full name" required
              class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
              style="border-radius: 2px;" />
          </div>
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Your · Email</label>
            <input v-model="form.email" type="email" placeholder="Email" required
              class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
              style="border-radius: 2px;" />
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Subject</label>
          <select v-model="form.subject"
            class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
            style="border-radius: 2px;">
            <option value="General">General · Enquiry</option>
            <option value="Bespoke">Bespoke · Blending</option>
            <option value="Shipping">Shipping</option>
            <option value="Complaint">Complaint</option>
            <option value="GDPR">GDPR · Privacy</option>
          </select>
        </div>

        <div class="mb-6">
          <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Message</label>
          <textarea v-model="form.message" rows="6" placeholder="How may we assist you?" required
            class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none resize-none font-serif text-ink-900"
            style="border-radius: 2px;"></textarea>
        </div>

        <button type="submit"
          class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
          style="border-radius: 2px;">
          Send · Your · Letter
        </button>
      </form>

      <!-- Office hours — 极细小字 -->
      <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-8 leading-loose font-sans">
        Office Hours &middot; Monday &ndash; Saturday &middot; 09:00 &ndash; 18:00 GMT<br>
        Closed Sundays &amp; UK Bank Holidays
      </p>
    </div>
  </div>
</template>
