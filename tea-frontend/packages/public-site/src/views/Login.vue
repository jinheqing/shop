<script setup lang="ts">
// ============================================================
// Login.vue — 老钱审美登录页
// ============================================================
// - off-black 主舞台，象牙白卡片
// - 不用 emoji（旧版 "Welcome Back 🍃" 已删）
// - 不用 alert()，用 in-page 错误条
// ============================================================

import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const err = ref('')

async function submit() {
  loading.value = true
  err.value = ''
  try {
    const r = await fetch('/api/v1/user/login', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, password: password.value }),
    }).then(x => x.json())
    if (r.access_token) { localStorage.setItem('user_token', r.access_token); router.push('/account') }
    else err.value = r.message || 'Login failed. Please check your credentials, or request a magic link.'
  } catch {
    err.value = 'A network error occurred. Please try again, or write to hello@ukteahouse.co.uk.'
  }
  finally { loading.value = false }
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100 flex items-center justify-center">
    <div class="w-full max-w-md mx-auto px-6 py-16">

      <!-- Section title -->
      <div class="mb-10 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Welcome · Back</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-3xl md:text-4xl text-ink-900">Sign · In.</h1>
      </div>

      <!-- Card -->
      <div class="bg-white border border-gold/15 p-8 md:p-10" style="border-radius: 2px;">

        <!-- Error banner -->
        <div v-if="err" class="mb-6 p-4 bg-ink-900 text-ivory-100 text-sm font-serif" style="border-radius: 2px;">
          <span class="text-gold text-[10px] uppercase tracking-lux font-sans">Attention</span><br>
          {{ err }}
        </div>

        <form @submit.prevent="submit" class="space-y-6">
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Your · Email</label>
            <input v-model="email" type="email" placeholder="Email" required
              class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
              style="border-radius: 2px;" />
          </div>
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Your · Password</label>
            <input v-model="password" type="password" placeholder="Password" required
              class="w-full px-4 py-3 bg-white border border-gold/20 focus:border-gold focus:outline-none font-serif text-ink-900"
              style="border-radius: 2px;" />
          </div>

          <button :disabled="loading"
            class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans disabled:opacity-50 disabled:cursor-not-allowed"
            style="border-radius: 2px;">
            {{ loading ? 'Signing · In · · ·' : 'Sign · In' }}
          </button>
        </form>

        <!-- Hairline divider -->
        <div class="hairline-gold my-7"></div>

        <!-- Magic link — 二级动作，文字链接 -->
        <div class="text-center">
          <RouterLink to="/magic-link"
            class="text-[11px] uppercase tracking-lux text-gold hover:text-gold-soft transition-duration-lux font-sans">
            Prefer · A · Magic · Link · No · Password · Needed &rarr;
          </RouterLink>
        </div>
      </div>

      <p class="text-[10px] uppercase tracking-lux text-sand text-center mt-6 leading-loose font-sans">
        Forgot your password &middot; <RouterLink to="/magic-link" class="hover:text-gold transition">Request a magic link</RouterLink>
      </p>
    </div>
  </div>
</template>
