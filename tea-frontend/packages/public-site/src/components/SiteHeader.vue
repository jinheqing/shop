<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const token = localStorage.getItem('user_token')
const route = useRoute()
const mobileOpen = ref(false)
const scrolled = ref(false)

function onScroll() {
  scrolled.value = window.scrollY > 20
}
onMounted(() => { window.addEventListener('scroll', onScroll) })
onUnmounted(() => { window.removeEventListener('scroll', onScroll) })

function closeMobile() { mobileOpen.value = false }

const mobileNav = [
  { to: '/bespoke', label: t('header.bespoke_blending') },
  { to: '/tea-gardens', label: t('header.tea_gardens') },
  { to: '/live', label: t('header.live_cameras') },
  { to: '/quality', label: t('header.quality') },
  { to: '/about', label: t('header.about') },
  { to: '/faq', label: t('header.faq') },
  { to: '/contact', label: t('header.contact') },
]
</script>

<template>
  <header
    class="fixed top-0 left-0 right-0 z-50 transition-all duration-300"
    :class="scrolled
      ? 'bg-white/90 backdrop-blur shadow-sm border-b border-tea-100'
      : 'bg-transparent'"
  >
    <div class="max-w-7xl mx-auto px-4 md:px-6 py-3 md:py-4 flex items-center justify-between">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 shrink-0">
        <span class="text-2xl">🍃</span>
        <span
          class="font-display text-xl tracking-brand"
          :class="scrolled ? 'text-tea-900' : 'text-white'"
        >UK Tea House</span>
      </RouterLink>

      <!-- Desktop nav -->
      <nav class="hidden md:flex items-center gap-7 text-sm tracking-wide"
           :class="scrolled ? 'text-tea-700' : 'text-white/90'">
        <RouterLink to="/bespoke" class="hover:text-tea-900 transition">{{ t('header.bespoke') }}</RouterLink>
        <RouterLink to="/tea-gardens" class="hover:text-tea-900 transition">{{ t('header.tea_gardens') }}</RouterLink>
        <RouterLink to="/live" class="hover:text-tea-900 transition">{{ t('header.live') }}</RouterLink>
        <RouterLink to="/quality" class="hover:text-tea-900 transition">{{ t('header.quality') }}</RouterLink>
        <RouterLink to="/about" class="hover:text-tea-900 transition">{{ t('header.about') }}</RouterLink>
        <RouterLink to="/faq" class="hover:text-tea-900 transition">{{ t('header.faq') }}</RouterLink>
        <RouterLink to="/contact" class="hover:text-tea-900 transition">{{ t('header.contact') }}</RouterLink>
      </nav>

      <!-- Right CTA -->
      <div class="hidden md:flex items-center gap-2 text-sm">
        <RouterLink v-if="!token" to="/login"
          class="px-4 py-2 rounded-full transition"
          :class="scrolled
            ? 'text-tea-800 hover:bg-tea-100'
            : 'text-white hover:bg-white/10'">{{ t('header.sign_in') }}</RouterLink>
        <RouterLink v-if="!token" to="/magic-link"
          class="px-4 py-2 rounded-full border transition"
          :class="scrolled
            ? 'border-tea-300 text-tea-800 hover:bg-tea-100'
            : 'border-white/40 text-white hover:bg-white/10'">{{ t('header.magic_link') }}</RouterLink>
        <!-- Concierge — 老钱核心：永远可见，永远能叫顾问，不依赖登录 -->
        <RouterLink to="/chat"
          :class="scrolled
            ? 'text-tea-700 hover:text-tea-900'
            : 'text-white hover:text-white/80'"
          :title="t('header.speak_advisor')">
          <span class="hidden lg:inline-flex items-center gap-2">
            <span class="w-1 h-1 rounded-full" :class="scrolled ? 'bg-gold' : 'bg-gold'"></span>
            <span class="text-[11px] uppercase tracking-lux">{{ t('header.speak_advisor') }}</span>
          </span>
          <span class="lg:hidden">Chat</span>
        </RouterLink>
        <template v-if="token">
          <RouterLink to="/account"
            class="px-4 py-2 rounded-full transition"
            :class="scrolled
              ? 'bg-tea-800 text-white hover:bg-tea-900'
              : 'bg-white text-tea-900 hover:bg-tea-100'">{{ t('header.my_account') }}</RouterLink>
        </template>
      </div>

      <!-- Mobile hamburger -->
      <button
        @click="mobileOpen = !mobileOpen"
        class="md:hidden w-10 h-10 rounded-lg flex items-center justify-center transition"
        :class="scrolled ? 'text-tea-900 hover:bg-tea-100' : 'text-white hover:bg-white/10'"
        aria-label="Toggle menu"
      >
        <svg v-if="!mobileOpen" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
        <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>

    <!-- Mobile menu -->
    <transition name="slide">
      <div v-if="mobileOpen"
        class="md:hidden bg-white border-t border-tea-100 shadow-lg">
        <nav class="max-w-7xl mx-auto px-6 py-4 flex flex-col gap-1">
          <RouterLink v-for="item in mobileNav" :key="item.to" :to="item.to" @click="closeMobile"
            class="py-3 border-b border-tea-50 text-tea-800 hover:text-tea-900 transition"
            :class="route.path === item.to ? 'font-medium text-tea-900' : ''">{{ item.label }}</RouterLink>
          <div class="pt-4 flex flex-col gap-2">
            <!-- Concierge 永远在最前 — 老钱核心 -->
            <RouterLink to="/chat" @click="closeMobile"
              class="w-full py-3 text-center rounded-full bg-tea-900 text-ivory-100 hover:bg-tea-800 transition">
              {{ t('header.speak_advisor') }}
            </RouterLink>
            <RouterLink v-if="!token" to="/login" @click="closeMobile"
              class="w-full py-3 text-center rounded-full border border-tea-300 text-tea-800 hover:bg-tea-100 transition">{{ t('header.sign_in') }}</RouterLink>
            <RouterLink v-if="!token" to="/magic-link" @click="closeMobile"
              class="w-full py-3 text-center rounded-full border border-tea-300 text-tea-800 hover:bg-tea-100 transition">{{ t('header.magic_link') }}</RouterLink>
            <RouterLink v-if="token" to="/account" @click="closeMobile"
              class="w-full py-3 text-center rounded-full bg-tea-800 text-white hover:bg-tea-900 transition">{{ t('header.my_account') }}</RouterLink>
          </div>
        </nav>
      </div>
    </transition>
  </header>
</template>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
