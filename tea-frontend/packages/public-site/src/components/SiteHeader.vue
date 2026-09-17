<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { ref, watchEffect } from 'vue'
const route = useRoute()
const scrolled = ref(false)
watchEffect(() => { window.onscroll = () => scrolled.value = window.scrollY > 10 })
</script>

<template>
  <header :class="[
    'fixed top-0 inset-x-0 z-50 transition-all duration-300',
    scrolled ? 'bg-white/95 backdrop-blur shadow-sm' : 'bg-transparent'
  ]">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
      <RouterLink to="/" class="flex items-center gap-2">
        <span class="text-2xl">🍃</span>
        <span :class="['font-serif text-lg tracking-wide', scrolled ? 'text-tea-900' : 'text-white']">UK Tea House</span>
      </RouterLink>
      <nav :class="['hidden md:flex items-center gap-8 text-sm font-medium', scrolled ? 'text-tea-800' : 'text-white/90']">
        <RouterLink to="/mountains" class="hover:text-tea-500 transition">Tea Mountains</RouterLink>
        <RouterLink to="/bespoke" class="hover:text-tea-500 transition">Bespoke</RouterLink>
        <RouterLink to="/live" class="hover:text-tea-500 transition flex items-center gap-1">
          <span class="inline-block w-2 h-2 rounded-full bg-red-500 live-dot"></span> Live
        </RouterLink>
        <RouterLink to="/quality" class="hover:text-tea-500 transition">Quality</RouterLink>
      </nav>
      <a href="mailto:hello@ukteahouse.co.uk" :class="[
        'px-4 py-2 rounded-full text-sm font-medium transition',
        scrolled ? 'bg-tea-700 text-white hover:bg-tea-800' : 'bg-white/20 text-white hover:bg-white/30 backdrop-blur'
      ]">Contact</a>
    </div>
  </header>
</template>
