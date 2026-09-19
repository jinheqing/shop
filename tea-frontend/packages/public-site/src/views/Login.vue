<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
    else err.value = r.message || 'Login failed'
  } catch { err.value = t('login.network_error') }
  finally { loading.value = false }
}
</script>
<template>
  <div class="min-h-screen pt-20 bg-tea-50">
    <div class="max-w-md mx-auto p-8 bg-white rounded-2xl shadow-sm border border-tea-100">
      <h1 class="font-serif text-3xl text-tea-900 mb-6 text-center">{{ t('login.welcome') }}</h1>
      <form @submit.prevent="submit" class="space-y-4">
        <input v-model="email" type="email" :placeholder="t('login.email')" required class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
        <input v-model="password" type="password" :placeholder="t('login.password')" required class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none" />
        <p v-if="err" class="text-red-600 text-sm">{{ err }}</p>
        <button :disabled="loading" class="w-full py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition disabled:opacity-50">
          {{ loading ? t('login.signing_in') : t('login.sign_in') }}
        </button>
      </form>
      <div class="mt-6 text-center text-sm">
        <RouterLink to="/magic-link" class="text-tea-700 hover:underline">{{ t('login.magic_link') }}</RouterLink>
      </div>
    </div>
  </div>
</template>
