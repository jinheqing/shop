<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <footer class="bg-tea-900 text-tea-100 mt-16 md:mt-20">
    <div class="max-w-7xl mx-auto px-6 py-16">
      <!-- Age + compliance banner -->
      <div class="mb-10 rounded-2xl border border-tea-700 bg-tea-950/40 px-5 py-4 text-sm text-tea-200 flex flex-col md:flex-row md:items-center gap-3 md:gap-6">
        <div class="flex items-center gap-2">
          <span class="font-display text-base">18+ only</span>
          <span class="text-tea-400">· You must be 18 or older to purchase. Not a medicinal product.</span>
        </div>
        <div class="hidden md:block text-tea-600">|</div>
        <div class="text-xs text-tea-400">SGS tested · VAT GB123456789 · ICO registered (pending)</div>
      </div>

      <!-- Main columns -->
      <div class="grid md:grid-cols-5 gap-10">
        <div class="md:col-span-2">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-3xl">🍃</span>
            <span class="font-display text-2xl tracking-brand">UK Tea House</span>
          </div>
          <p class="text-sm text-tea-300 leading-relaxed max-w-sm">Single-origin Pu'er tea, traceable from cloud tea garden to cup. Live-streamed farms, SGS-certified, bespoke blended for you inside 45 days.</p>
          <div class="mt-6 text-xs text-tea-400 leading-relaxed">
            12 Savile Row, Mayfair, London W1J 5PA
            <br/>hello@ukteahouse.co.uk · +44 20 7946 0958
          </div>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Shop</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/bespoke" class="hover:text-white transition">Bespoke Blending</RouterLink></li>
            <li><RouterLink to="/tea-gardens" class="hover:text-white transition">Tea Gardens</RouterLink></li>
            <li><RouterLink to="/live" class="hover:text-white transition">Live Cameras</RouterLink></li>
            <li><RouterLink to="/quality" class="hover:text-white transition">SGS Quality</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Account</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/login" class="hover:text-white transition">Sign In</RouterLink></li>
            <li><RouterLink to="/account" class="hover:text-white transition">My Orders</RouterLink></li>
            <li><RouterLink to="/chat" class="hover:text-white transition">Chat with Advisor</RouterLink></li>
            <li><RouterLink to="/contact" class="hover:text-white transition">Contact Us</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Legal &amp; Privacy</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/privacy" class="hover:text-white transition">Privacy · GDPR</RouterLink></li>
            <li><RouterLink to="/terms" class="hover:text-white transition">Terms of Service</RouterLink></li>
            <li><RouterLink to="/cookie-policy" class="hover:text-white transition">Cookie Policy</RouterLink></li>
            <li><RouterLink to="/gdpr-dsar" class="hover:text-white transition">Data Subject Request</RouterLink></li>
            <li><RouterLink to="/faq" class="hover:text-white transition">FAQ</RouterLink></li>
          </ul>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="mt-12 pt-8 border-t border-tea-800 flex flex-col md:flex-row justify-between gap-4 text-xs text-tea-400">
        <div>© 2026 UK Tea House Ltd. · Company No. 12345678 · VAT GB123456789</div>
        <div class="hidden md:block">Registered in England &amp; Wales · ICO registered (pending) · Self-hosted infrastructure in London</div>
      </div>
    </div>
  </footer>
</template>
