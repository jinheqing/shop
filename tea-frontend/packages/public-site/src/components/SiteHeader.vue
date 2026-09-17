<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
const token = localStorage.getItem('user_token')
</script>
<template>
  <header class="fixed top-0 left-0 right-0 z-50 transition-all" :class="{'bg-white shadow-sm': false}">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
      <RouterLink to="/" class="flex items-center gap-2">
        <span class="text-2xl">🍃</span>
        <span class="font-serif text-xl text-tea-900">UK Tea House</span>
      </RouterLink>
      <nav class="hidden md:flex items-center gap-7 text-sm text-tea-700">
        <RouterLink to="/bespoke" class="hover:text-tea-900">Bespoke</RouterLink>
        <RouterLink to="/tea-gardens" class="hover:text-tea-900">Tea Gardens</RouterLink>
        <RouterLink to="/live" class="hover:text-tea-900">Live</RouterLink>
        <RouterLink to="/quality" class="hover:text-tea-900">Quality</RouterLink>
        <RouterLink to="/about" class="hover:text-tea-900">About</RouterLink>
        <RouterLink to="/faq" class="hover:text-tea-900">FAQ</RouterLink>
        <RouterLink to="/contact" class="hover:text-tea-900">Contact</RouterLink>
      </nav>
      <div class="flex items-center gap-3 text-sm">
        <RouterLink v-if="!token" to="/login" class="px-4 py-2 text-tea-800 hover:bg-tea-100 rounded-full transition">Sign In</RouterLink>
        <RouterLink v-if="!token" to="/magic-link" class="px-4 py-2 border border-tea-300 rounded-full hover:bg-tea-100 transition">Magic Link</RouterLink>
        <template v-if="token">
          <RouterLink to="/chat" class="px-3 py-2 text-tea-700 hover:text-tea-900" title="Chat with Advisor">💬 Chat</RouterLink>
          <RouterLink to="/account" class="px-4 py-2 bg-tea-800 text-white rounded-full hover:bg-tea-900 transition">My Account</RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>
