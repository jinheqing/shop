<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const step = ref<'request' | 'sent' | 'verify'>('request')
const email = ref('')
const token = ref('')
const verifiedEmail = ref('')

// 推荐人（全部可选，低调放在底部）
const referralSource = ref('')
const referrerName = ref('')
// 短链 code — 从 URL query ?r=XXX 里取
const urlParams = new URLSearchParams(window.location.search)
const referralCodeFromURL = ref(urlParams.get('r') || '')

async function request() {
  const payload: any = { email: email.value }
  // 只有填了才发
  if (referralSource.value) payload.referral_source = referralSource.value
  if (referrerName.value.trim()) payload.referrer_name = referrerName.value.trim()
  if (referralCodeFromURL.value) payload.referral_code = referralCodeFromURL.value

  await fetch('/api/v1/user/magic-link/request', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
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
  const r = await fetch('/api/v1/user/magic-link/verify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: verifiedEmail.value, token: raw }),
  }).then(x => x.json())
  if (r.access_token) {
    localStorage.setItem('user_token', r.access_token)
    ElMessage.success('Welcome!')
    window.location.href = '/account'
  } else {
    ElMessage.error(r.message || 'Verification failed')
  }
}

// 用户选了 "friend" 才显示 referrer name 输入框
function showReferrerNameInput(): boolean {
  return referralSource.value === 'friend'
}
</script>
<template>
  <div class="min-h-screen pt-20 bg-tea-50">
    <div class="max-w-md mx-auto p-8 bg-white rounded-2xl shadow-sm border border-tea-100">
      <h1 class="font-serif text-3xl text-tea-900 mb-6 text-center">🔬 Magic Link Login</h1>
      <p v-if="step==='request'" class="text-sm text-tea-600 mb-6">No password needed. We'll email you a one-time link.</p>
      <div v-if="step==='request'">
        <input v-model="email" type="email" placeholder="Your email" required class="w-full px-4 py-3 rounded-xl border border-tea-200 mb-4" />

        <!-- ===== 推荐人字段（低调放在底部，折叠式）===== -->
        <details class="mb-4 border-t border-tea-100 pt-4">
          <summary class="text-xs text-tea-500 cursor-pointer hover:text-tea-700">How did you hear about us? (optional)</summary>
          <div class="mt-3 space-y-3">
            <select v-model="referralSource" class="w-full px-4 py-2.5 rounded-xl border border-tea-200 text-sm text-tea-700 bg-white">
              <option value="">Select...</option>
              <option value="friend">A friend recommended</option>
              <option value="youtube">YouTube</option>
              <option value="search">Search engine</option>
              <option value="other">Other</option>
            </select>
            <input
              v-if="showReferrerNameInput()"
              v-model="referrerName"
              type="text"
              placeholder="Your friend's name"
              class="w-full px-4 py-2.5 rounded-xl border border-tea-200 text-sm"
            />
            <p v-if="referralCodeFromURL" class="text-xs text-tea-500">We detected a referral link — thank you to your friend!</p>
          </div>
        </details>

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
