<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
const form = ref({ name: '', email: '', subject: 'General', message: '' })
function submit() { ElMessage.success('Thanks! We\'ll reply within 24 hours.'); form.value = { name:'', email:'', subject:'General', message:'' } }
</script>
<template>
  <div class="pt-20 py-20 bg-tea-50">
    <div class="max-w-3xl mx-auto px-6">
      <h1 class="font-serif text-5xl text-tea-900 mb-4 text-center">Get In Touch</h1>
      <p class="text-tea-600 text-center mb-12">Advisors respond within 24 hours. Or chat with us on WhatsApp +44 7700 900 123.</p>

      <div class="grid md:grid-cols-3 gap-6 mb-10">
        <div class="bg-white rounded-xl p-6 text-center">
          <div class="text-3xl mb-3">📧</div><div class="font-medium mb-1">Email</div><div class="text-xs text-tea-600">hello@ukteahouse.co.uk</div>
        </div>
        <div class="bg-white rounded-xl p-6 text-center">
          <div class="text-3xl mb-3">📞</div><div class="font-medium mb-1">Phone</div><div class="text-xs text-tea-600">+44 20 7946 0958</div>
        </div>
        <div class="bg-white rounded-xl p-6 text-center">
          <div class="text-3xl mb-3">📍</div><div class="font-medium mb-1">Office</div><div class="text-xs text-tea-600">Mayfair, London W1J</div>
        </div>
      </div>

      <form @submit.prevent="submit" class="bg-white rounded-2xl p-8 border border-tea-100">
        <div class="grid md:grid-cols-2 gap-4 mb-4">
          <input v-model="form.name" placeholder="Your name" required class="px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
          <input v-model="form.email" type="email" placeholder="Email" required class="px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
        </div>
        <select v-model="form.subject" class="w-full mb-4 px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none">
          <option value="General">General Enquiry</option><option value="Bespoke">Bespoke Blending</option><option value="Shipping">Shipping</option><option value="Complaint">Complaint</option><option value="GDPR">GDPR / Privacy</option>
        </select>
        <textarea v-model="form.message" rows="5" placeholder="How can we help?" required class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none mb-4" />
        <button type="submit" class="w-full py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition">Send Message →</button>
      </form>
    </div>
  </div>
</template>
