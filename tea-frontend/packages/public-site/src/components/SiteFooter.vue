<template>
  <footer class="bg-tea-950 text-tea-200 py-16">
    <div class="max-w-7xl mx-auto px-6 grid md:grid-cols-4 gap-10">
      <div>
        <div class="flex items-center gap-2 mb-4">
          <span class="text-2xl">🍃</span>
          <span class="font-serif text-lg text-white">UK Tea House</span>
        </div>
        <p class="text-sm text-tea-300 leading-relaxed">Single-origin Pu'er tea, traceable from Yunnan cloud mountains to your London door.</p>
      </div>
      <div>
        <h4 class="text-white font-medium mb-4">Explore</h4>
        <ul class="space-y-2 text-sm">
          <li><RouterLink to="/mountains" class="hover:text-tea-300">Tea Mountains</RouterLink></li>
          <li><RouterLink to="/bespoke" class="hover:text-tea-300">Bespoke Blending</RouterLink></li>
          <li><RouterLink to="/live" class="hover:text-tea-300">Live Cameras</RouterLink></li>
          <li><RouterLink to="/quality" class="hover:text-tea-300">SGS Quality</RouterLink></li>
        </ul>
      </div>
      <div>
        <h4 class="text-white font-medium mb-4">Company</h4>
        <ul class="space-y-2 text-sm text-tea-300">
          <li>hello@ukteahouse.co.uk</li>
          <li>+44 20 7946 0958</li>
          <li>Mayfair, London W1J</li>
        </ul>
      </div>
      <div>
        <h4 class="text-white font-medium mb-4">Certifications</h4>
        <div class="flex gap-3 text-xs">
          <span class="px-2 py-1 border border-tea-600 rounded">🇬🇧 UK VAT</span>
          <span class="px-2 py-1 border border-tea-600 rounded">GDPR</span>
          <span class="px-2 py-1 border border-tea-600 rounded">SGS</span>
        </div>
      </div>
    </div>
    <div class="border-t border-tea-800 mt-10 pt-6 text-center text-xs text-tea-500">
      © {{ new Date().getFullYear() }} UK Tea House Ltd. All rights reserved.
    </div>
  </footer>
</template>
<script setup>
import { RouterLink } from 'vue-router'
</script>
