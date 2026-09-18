<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const step = ref<'request' | 'sent' | 'verify'>('request')
const email = ref('')
const token = ref('')
const verifiedEmail = ref('')

async function request() {
  await fetch('/api/v1/user/magic-link/request', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email: email.value }) })
  step.value = 'sent'
  verifiedEmail.value = email.value
}
async function verify() {
  // 支持粘贴完整 URL 或裸 token
  let raw = token.value.trim()
  if (raw.includes('token=')) {
    try {
      raw = raw.split('token=')[1].split(/[&?#]/)[0]
    } catch {}
  }
  const r = await fetch('/api/v1/user/magic-link/verify', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email: verifiedEmail.value, token: raw }) }).then(x => x.json())
  if (r.access_token) { localStorage.setItem('user_token', r.access_token); ElMessage.success('Welcome!'); window.location.href = '/account' }
  else { ElMessage.error(r.message || 'Verification failed') }
}
</script>
<template>
  <div class="min-h-screen pt-20 bg-tea-50">
    <div class="max-w-md mx-auto p-8 bg-white rounded-2xl shadow-sm border border-tea-100">
      <h1 class="font-serif text-3xl text-tea-900 mb-6 text-center">🔬 Magic Link Login</h1>
      <p v-if="step==='request'" class="text-sm text-tea-600 mb-6">No password needed. We'll email you a one-time link.</p>
      <div v-if="step==='request'">
        <input v-model="email" type="email" placeholder="Your email" required class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-4" />
        <button @click="request" class="w-full py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition">Send Magic Link →</button>
      </div>
      <div v-if="step==='sent'">
        <div class="text-center text-5xl mb-4">📧</div>
        <p class="text-center text-tea-700 mb-6">We sent a magic link to <strong>{{ verifiedEmail }}</strong>. Click the link in your email to sign in — or paste the full link or token below.</p>
        <input v-model="token" placeholder="Paste magic link or token from email" class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-4" />
        <button @click="verify" class="w-full py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition">Verify & Sign In →</button>
      </div>
    </div>
  </div>
</template>
