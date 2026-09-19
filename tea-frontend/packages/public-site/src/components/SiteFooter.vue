<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <footer class="bg-tea-900 text-tea-100 mt-16 md:mt-20">
    <div class="max-w-7xl mx-auto px-6 py-16">
      <!-- Age + compliance banner -->
      <div class="mb-10 rounded-2xl border border-tea-700 bg-tea-950/40 px-5 py-4 text-sm text-tea-200 flex flex-col md:flex-row md:items-center gap-3 md:gap-6">
        <div class="flex items-center gap-2">
          <span class="font-display text-base">{{ t('footer.only_18') }}</span>
          <span class="text-tea-400">{{ t('footer.age_notice') }}</span>
        </div>
        <div class="hidden md:block text-tea-600">|</div>
        <div class="text-xs text-tea-400">{{ t('footer.sgs_notice') }}</div>
      </div>

      <!-- Main columns -->
      <div class="grid md:grid-cols-5 gap-10">
        <div class="md:col-span-2">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-3xl">🍃</span>
            <span class="font-display text-2xl tracking-brand">UK Tea House</span>
          </div>
          <p class="text-sm text-tea-300 leading-relaxed max-w-sm">{{ t('footer.brand_desc') }}</p>
          <div class="mt-6 text-xs text-tea-400 leading-relaxed">
            12 Savile Row, Mayfair, London W1J 5PA
            <br/>hello@ukteahouse.co.uk · +44 20 7946 0958
          </div>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">{{ t('footer.shop') }}</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/bespoke" class="hover:text-white transition">{{ t('header.bespoke_blending') }}</RouterLink></li>
            <li><RouterLink to="/tea-gardens" class="hover:text-white transition">{{ t('header.tea_gardens') }}</RouterLink></li>
            <li><RouterLink to="/live" class="hover:text-white transition">{{ t('header.live_cameras') }}</RouterLink></li>
            <li><RouterLink to="/quality" class="hover:text-white transition">{{ t('header.quality') }}</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">{{ t('footer.account') }}</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/login" class="hover:text-white transition">{{ t('header.sign_in') }}</RouterLink></li>
            <li><RouterLink to="/account" class="hover:text-white transition">{{ t('footer.my_orders') }}</RouterLink></li>
            <li><RouterLink to="/chat" class="hover:text-white transition">{{ t('footer.chat_advisor') }}</RouterLink></li>
            <li><RouterLink to="/contact" class="hover:text-white transition">{{ t('footer.contact_us') }}</RouterLink></li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium mb-4 text-white">{{ t('footer.legal') }}</h4>
          <ul class="space-y-2 text-sm">
            <li><RouterLink to="/privacy" class="hover:text-white transition">{{ t('footer.privacy_gdpr') }}</RouterLink></li>
            <li><RouterLink to="/terms" class="hover:text-white transition">{{ t('footer.terms') }}</RouterLink></li>
            <li><RouterLink to="/cookie-policy" class="hover:text-white transition">{{ t('footer.cookie_policy') }}</RouterLink></li>
            <li><RouterLink to="/gdpr-dsar" class="hover:text-white transition">{{ t('footer.data_request') }}</RouterLink></li>
            <li><RouterLink to="/faq" class="hover:text-white transition">{{ t('footer.faq') }}</RouterLink></li>
          </ul>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="mt-12 pt-8 border-t border-tea-800 flex flex-col md:flex-row justify-between gap-4 text-xs text-tea-400">
        <div>{{ t('footer.copyright') }}</div>
        <div class="hidden md:block">{{ t('footer.registered') }}</div>
      </div>
    </div>
  </footer>
</template>
