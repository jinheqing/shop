<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

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

// 非首页没有暗色 Hero 背景，始终用深色文字 + 象牙白底
const isHome = computed(() => route.path === '/')
const dark = computed(() => scrolled.value || !isHome.value)
</script>

<template>
  <header
    class="fixed top-0 left-0 right-0 z-50 transition-all duration-lux"
    :class="dark
      ? 'bg-ivory-100/95 backdrop-blur border-b border-gold/15'
      : 'bg-transparent'">
    <div class="max-w-7xl mx-auto px-4 md:px-6 py-4 md:py-5 flex items-center justify-between">
      <!-- Logo — 不用 emoji，serif 字标 + 香槟金小点 -->
      <RouterLink to="/" class="flex items-center gap-3 shrink-0">
        <span class="w-1.5 h-1.5 rounded-full" :class="dark ? 'bg-gold' : 'bg-gold'"></span>
        <span
          class="font-serif text-xl tracking-brand"
          :class="dark ? 'text-ink-900' : 'text-ivory-100'">
          UK · Tea · House
        </span>
      </RouterLink>

      <!-- Desktop nav — serif + tracking-wide -->
      <nav class="hidden md:flex items-center gap-8 text-sm tracking-wide"
           :class="dark ? 'text-sand-deep' : 'text-ivory-100/90'">
        <RouterLink to="/bespoke" class="hover:text-gold transition-duration-lux">Bespoke</RouterLink>
        <RouterLink to="/tea-gardens" class="hover:text-gold transition-duration-lux">Tea Gardens</RouterLink>
        <RouterLink to="/live" class="hover:text-gold transition-duration-lux">Live</RouterLink>
        <RouterLink to="/quality" class="hover:text-gold transition-duration-lux">Quality</RouterLink>
        <RouterLink to="/about" class="hover:text-gold transition-duration-lux">About</RouterLink>
        <RouterLink to="/faq" class="hover:text-gold transition-duration-lux">FAQ</RouterLink>
        <RouterLink to="/contact" class="hover:text-gold transition-duration-lux">Contact</RouterLink>
      </nav>

      <!-- Right CTA — 唯一突出的是 "Speak with an Advisor" -->
      <div class="hidden md:flex items-center gap-3 text-sm">
        <RouterLink v-if="!token" to="/magic-link"
          class="text-[11px] uppercase tracking-lux transition-duration-lux"
          :class="dark
            ? 'text-sand-deep hover:text-gold'
            : 'text-ivory-100/80 hover:text-gold'">Sign In</RouterLink>

        <!-- Concierge — 老钱核心：永远可见 -->
        <RouterLink to="/chat"
          class="flex items-center gap-2 px-4 py-2 border transition-duration-lux"
          :class="dark
            ? 'border-gold/40 text-gold hover:bg-gold/5'
            : 'border-gold/40 text-gold hover:bg-gold/10'"
          title="Speak with an Advisor">
          <span class="w-1 h-1 rounded-full bg-gold"></span>
          <span class="text-[11px] uppercase tracking-lux font-sans">Speak · With · An · Advisor</span>
        </RouterLink>

        <template v-if="token">
          <RouterLink to="/account"
            class="px-4 py-2 transition-duration-lux text-[11px] uppercase tracking-lux"
            :class="dark
              ? 'bg-ink-900 text-ivory-100 hover:bg-ink-800'
              : 'bg-ivory-100 text-ink-900 hover:bg-ivory-200'">My Account</RouterLink>
        </template>
      </div>

      <!-- Mobile hamburger — 44px 触摸区 -->
      <button
        @click="mobileOpen = !mobileOpen"
        class="md:hidden w-12 h-12 flex items-center justify-center transition-duration-lux"
        :class="dark ? 'text-ink-900 hover:text-gold' : 'text-ivory-100 hover:text-gold'"
        aria-label="Toggle menu"
      >
        <svg v-if="!mobileOpen" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
        <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>

    <!-- Mobile menu — 象牙白面板 + hairline 分隔 -->
    <transition name="slide">
      <div v-if="mobileOpen"
        class="md:hidden bg-ivory-100 border-t border-gold/15 shadow-lg">
        <nav class="max-w-7xl mx-auto px-6 py-4 flex flex-col">
          <RouterLink v-for="item in ([
            {to:'/bespoke', label:'Bespoke Blending'},
            {to:'/tea-gardens', label:'Tea Gardens'},
            {to:'/live', label:'Live Cameras'},
            {to:'/quality', label:'SGS Quality'},
            {to:'/about', label:'About'},
            {to:'/faq', label:'FAQ'},
            {to:'/contact', label:'Contact'},
          ])" :key="item.to" :to="item.to" @click="closeMobile"
            class="py-4 border-b border-gold/10 text-ink-900 hover:text-gold transition-duration-lux font-serif text-lg"
            :class="route.path === item.to ? 'text-gold' : ''">{{ item.label }}</RouterLink>

          <div class="pt-6 flex flex-col gap-3">
            <!-- Concierge 永远在最前 — 老钱核心 -->
            <RouterLink to="/chat" @click="closeMobile"
              class="w-full py-4 text-center text-[11px] uppercase tracking-lux bg-ink-900 text-gold hover:bg-ink-800 transition-duration-lux"
              style="border-radius: 2px;">
              Speak · With · An · Advisor
            </RouterLink>
            <RouterLink v-if="!token" to="/magic-link" @click="closeMobile"
              class="w-full py-4 text-center text-[11px] uppercase tracking-lux border border-gold/40 text-gold hover:bg-gold/5 transition-duration-lux"
              style="border-radius: 2px;">Sign In</RouterLink>
            <RouterLink v-if="token" to="/account" @click="closeMobile"
              class="w-full py-4 text-center text-[11px] uppercase tracking-lux bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux"
              style="border-radius: 2px;">My Account</RouterLink>
          </div>
        </nav>
      </div>
    </transition>
  </header>
</template>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: opacity 300ms ease, transform 300ms ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
