<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <footer class="bg-ink-900 text-ivory-100 mt-16 md:mt-20">
    <div class="max-w-7xl mx-auto px-6 py-16">
      <!-- Compliance note — 极细的 serif 小字，不再用大色块 banner -->
      <div class="mb-10 text-center">
        <p class="text-[10px] uppercase tracking-lux text-sand font-sans leading-loose">
          18 &middot; Only &middot; You must be 18 or older to purchase &middot; Not a medicinal product<br>
          SGS tested &middot; VAT GB123456789 &middot; ICO registered (pending)
        </p>
      </div>

      <!-- Hairline gold divider -->
      <div class="hairline-gold mb-12"></div>

      <!-- Main columns -->
      <div class="grid md:grid-cols-5 gap-10">
        <div class="md:col-span-2">
          <div class="flex items-center gap-3 mb-5">
            <span class="w-1.5 h-1.5 rounded-full bg-gold"></span>
            <span class="font-serif text-2xl tracking-brand">UK · Tea · House</span>
          </div>
          <p class="text-sm text-sand leading-relaxed max-w-sm font-serif">
            Single-origin Pu'er tea, traceable from cloud tea garden to cup.
            Live-streamed farms, SGS-certified, bespoke blended for you inside 45 days.
          </p>
          <div class="mt-6 text-[10px] uppercase tracking-lux text-sand leading-loose font-sans">
            12 Savile Row, Mayfair, London W1J 5PA<br>
            hello@ukteahouse.co.uk &middot; +44 20 7946 0958
          </div>
        </div>

        <div>
          <h4 class="text-[10px] uppercase tracking-lux text-gold mb-5 font-sans">Shop</h4>
          <ul class="space-y-3 text-sm font-serif">
            <li><RouterLink to="/bespoke" class="text-sand hover:text-gold transition-duration-lux">Bespoke Blending</RouterLink></li>
            <li><RouterLink to="/tea-gardens" class="text-sand hover:text-gold transition-duration-lux">Tea Gardens</RouterLink></li>
            <li><RouterLink to="/live" class="text-sand hover:text-gold transition-duration-lux">Live Cameras</RouterLink></li>
            <li><RouterLink to="/quality" class="text-sand hover:text-gold transition-duration-lux">SGS Quality</RouterLink></li>
          </ul>
        </div>

        <div>
          <h4 class="text-[10px] uppercase tracking-lux text-gold mb-5 font-sans">Account</h4>
          <ul class="space-y-3 text-sm font-serif">
            <li><RouterLink to="/magic-link" class="text-sand hover:text-gold transition-duration-lux">Sign In</RouterLink></li>
            <li><RouterLink to="/account" class="text-sand hover:text-gold transition-duration-lux">My Orders</RouterLink></li>
            <li><RouterLink to="/chat" class="text-sand hover:text-gold transition-duration-lux">Chat with Advisor</RouterLink></li>
            <li><RouterLink to="/contact" class="text-sand hover:text-gold transition-duration-lux">Contact Us</RouterLink></li>
          </ul>
        </div>

        <div>
          <h4 class="text-[10px] uppercase tracking-lux text-gold mb-5 font-sans">Legal &amp; Privacy</h4>
          <ul class="space-y-3 text-sm font-serif">
            <li><RouterLink to="/privacy" class="text-sand hover:text-gold transition-duration-lux">Privacy · GDPR</RouterLink></li>
            <li><RouterLink to="/terms" class="text-sand hover:text-gold transition-duration-lux">Terms of Service</RouterLink></li>
            <li><RouterLink to="/cookie-policy" class="text-sand hover:text-gold transition-duration-lux">Cookie Policy</RouterLink></li>
            <li><RouterLink to="/gdpr-dsar" class="text-sand hover:text-gold transition-duration-lux">Data Subject Request</RouterLink></li>
            <li><RouterLink to="/faq" class="text-sand hover:text-gold transition-duration-lux">FAQ</RouterLink></li>
          </ul>
        </div>
      </div>

      <!-- Hairline gold divider -->
      <div class="hairline-gold my-10"></div>

      <!-- Bottom bar — 极细小字 -->
      <div class="flex flex-col md:flex-row justify-between gap-3 text-[10px] uppercase tracking-lux text-sand font-sans">
        <div>&copy; 2026 UK Tea House Ltd. &middot; Company No. 12345678 &middot; VAT GB123456789</div>
        <div class="hidden md:block">Registered in England &amp; Wales &middot; ICO registered (pending) &middot; Self-hosted in London</div>
      </div>
    </div>
  </footer>
</template>
