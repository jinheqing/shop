<script setup lang="ts">
// ============================================================
// MagicLink.vue — 老钱审美 magic link 登录页
// - 删除 emoji 🔬📧
// - 不用 ElMessage，用 in-page status message
// - squared border + ink/ivory/gold palette
// ============================================================
import { ref } from 'vue'

const step = ref<'request' | 'sent' | 'verify'>('request')
const email = ref('')
const token = ref('')
const verifiedEmail = ref('')
const msg = ref<{ kind: 'error' | 'ok' | 'info'; text: string }>({ kind: 'info', text: '' })

async function request() {
  msg.value = { kind: 'info', text: '' }
  try {
    await fetch('/api/v1/user/magic-link/request', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value }),
    })
    step.value = 'sent'
    verifiedEmail.value = email.value
  } catch {
    msg.value = { kind: 'error', text: 'A network error occurred. Please try again.' }
  }
}

async function verify() {
  msg.value = { kind: 'info', text: '' }
  let raw = token.value.trim()
  if (raw.includes('token=')) {
    try { raw = raw.split('token=')[1].split(/[&?#]/)[0] } catch {}
  }
  try {
    const r = await fetch('/api/v1/user/magic-link/verify', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: verifiedEmail.value, token: raw }),
    }).then(x => x.json())
    if (r.access_token) {
      localStorage.setItem('user_token', r.access_token)
      msg.value = { kind: 'ok', text: 'Welcome. Taking you to your account…' }
      setTimeout(() => { window.location.href = '/account' }, 700)
    } else {
      msg.value = { kind: 'error', text: r.message || 'Verification failed. Please request a fresh link.' }
    }
  } catch {
    msg.value = { kind: 'error', text: 'A network error occurred. Please try again.' }
  }
}
</script>

<template>
  <div class="pt-20 min-h-screen bg-ivory-100 flex items-center justify-center">
    <div class="w-full max-w-md mx-auto px-6 py-16">

      <!-- Section title -->
      <div class="mb-10 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Sign · In · By · Link</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-3xl md:text-4xl text-ink-900 mb-3 leading-tight">
          Magic <span class="italic">Link</span>
        </h1>
        <p class="text-sm text-sand font-serif">No password required. We'll send a one-time link to your email.</p>
      </div>

      <!-- Card -->
      <div class="bg-white border border-gold/20 p-6 md:p-8" style="border-radius: 2px;">

        <!-- In-page status -->
        <div v-if="msg.text"
          :class="[
            'mb-5 px-4 py-3 text-sm font-serif leading-relaxed',
            msg.kind === 'error' ? 'bg-ink-900 text-ivory-100' : 'border border-gold/30 bg-ivory-100 text-ink-900'
          ]"
          style="border-radius: 2px;">
          <div class="flex items-center gap-2 mb-1">
            <span class="w-1 h-1 rounded-full"
              :class="msg.kind === 'error' ? 'bg-red-400' : 'bg-gold'"></span>
            <span class="text-[10px] uppercase tracking-lux font-sans"
              :class="msg.kind === 'error' ? 'text-red-300' : 'text-gold'">
              {{ msg.kind === 'error' ? 'Notice' : msg.kind === 'ok' ? 'Acknowledged' : 'Info' }}
            </span>
          </div>
          <p :class="msg.kind === 'error' ? 'text-ivory-100/80' : 'text-sand'">{{ msg.text }}</p>
        </div>

        <!-- Request step -->
        <div v-if="step==='request'" class="space-y-4">
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Your · Email</label>
            <input v-model="email" type="email" placeholder="you@domain.co.uk" required
              class="w-full px-4 py-3 border border-gold/20 bg-ivory-50 text-ink-900 font-serif focus:border-gold focus:outline-none"
              style="border-radius: 2px;" />
          </div>
          <button @click="request"
            class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Send · Magic · Link
          </button>
        </div>

        <!-- Sent step -->
        <div v-if="step==='sent'" class="space-y-4">
          <!-- Letter icon — hairline framed square, no emoji -->
          <div class="flex items-center justify-center mb-3">
            <div class="w-12 h-12 border border-gold/40 flex items-center justify-center" style="border-radius: 2px;">
              <div class="w-5 h-3 border border-gold"></div>
            </div>
          </div>
          <p class="text-center text-sand font-serif text-sm leading-relaxed">
            We've sent a magic link to <strong class="text-ink-900">{{ verifiedEmail }}</strong>.
            Click the link in your email to sign in — or paste the full link or token below.
          </p>
          <div>
            <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Paste · Link · Or · Token</label>
            <input v-model="token" placeholder="https://ukteahouse.co.uk/magic?token=…"
              class="w-full px-4 py-3 border border-gold/20 bg-ivory-50 text-ink-900 font-serif focus:border-gold focus:outline-none"
              style="border-radius: 2px;" />
          </div>
          <button @click="verify"
            class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans"
            style="border-radius: 2px;">
            Verify · &amp; · Sign · In
          </button>
        </div>
      </div>

      <p class="text-center text-xs text-sand/70 mt-6 font-serif">
        Alternatively, <RouterLink to="/login" class="text-gold underline hover:text-ink-900">sign in with a password</RouterLink>.
      </p>
    </div>
  </div>
</template>
