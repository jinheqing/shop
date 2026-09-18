<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/stores/auth'

const router = useRouter()
const auth = useAuth()
const form = ref(import.meta.env.DEV ? { email: 'admin@ukteahouse.co.uk', password: 'Admin!Tea2026' } : { email: '', password: '' })
const loading = ref(false)

async function submit() {
  loading.value = true
  try { await auth.login(form.value.email, form.value.password); ElMessage.success('Welcome back!'); router.push('/') }
  catch { ElMessage.error('Invalid credentials') }
  finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-amber-900 to-stone-800">
    <div class="bg-white rounded-2xl shadow-2xl p-10 w-full max-w-md">
      <div class="text-center mb-8">
        <div class="text-5xl mb-2">🍃</div>
        <h1 class="font-serif text-2xl text-slate-900">UK Tea House Admin</h1>
        <p class="text-sm text-slate-500 mt-2">Staff portal · Pu'er Private Domain</p>
      </div>
      <el-form :model="form" label-position="top" @submit.prevent="submit">
        <el-form-item label="Email"><el-input v-model="form.email" size="large" /></el-form-item>
        <el-form-item label="Password"><el-input v-model="form.password" type="password" show-password size="large" /></el-form-item>
        <el-button type="primary" size="large" :loading="loading" class="w-full bg-slate-900 hover:bg-slate-800" @click="submit">Sign In</el-button>
      </el-form>
      <p class="text-center text-xs text-slate-400 mt-6">Default admin (dev only) — change immediately on first login</p>
    </div>
  </div>
</template>
