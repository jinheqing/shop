<script setup lang="ts">import { RouterLink } from 'vue-router'</script>
<template>
  <footer class="bg-tea-900 text-tea-100 mt-20">
    <div class="max-w-7xl mx-auto px-6 py-16">
      <div class="grid md:grid-cols-5 gap-10">
        <div class="md:col-span-2">
          <div class="flex items-center gap-2 mb-4"><span class="text-3xl">🍃</span><span class="font-serif text-2xl">UK Tea House</span></div>
          <p class="text-sm text-tea-300 leading-relaxed max-w-sm">Single-origin Pu'er tea, traceable from cloud mountain to cup. Live-streamed farms, SGS-certified, bespoke blended for you inside 45 days.</p>
          <div class="mt-6 text-xs text-tea-400">Mayfair, London W1J · hello@ukteahouse.co.uk · +44 20 7946 0958</div>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Shop</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/bespoke" class="hover:text-white">Bespoke Blending</RouterLink></li>
            <li><RouterLink to="/mountains" class="hover:text-white">Tea Mountains</RouterLink></li>
            <li><RouterLink to="/live" class="hover:text-white">Live Cameras</RouterLink></li>
            <li><RouterLink to="/quality" class="hover:text-white">SGS Quality</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Account</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/login" class="hover:text-white">Sign In</RouterLink></li>
            <li><RouterLink to="/account" class="hover:text-white">My Orders</RouterLink></li>
            <li><RouterLink to="/chat" class="hover:text-white">Chat with Advisor</RouterLink></li>
            <li><RouterLink to="/contact" class="hover:text-white">Contact Us</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">Legal</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/privacy" class="hover:text-white">Privacy · GDPR</RouterLink></li>
            <li><a href="#" class="hover:text-white">Terms of Service</a></li>
            <li><a href="#" class="hover:text-white">Cookie Policy</a></li>
            <li><RouterLink to="/faq" class="hover:text-white">FAQ</RouterLink></li>
          </ul>
        </div>
      </div>
      <div class="mt-12 pt-8 border-t border-tea-800 flex flex-col md:flex-row justify-between gap-4 text-xs text-tea-400">
        <div>© 2026 UK Tea House Ltd. · Company No. 12345678 · VAT GB123456789</div>
        <div>Registered in England & Wales · ICO registered (pending) · Self-hosted infrastructure in London</div>
      </div>
    </div>
  </footer>
</template>
