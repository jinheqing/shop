<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import enLocale from 'element-plus/es/locale/lang/en'
import { useAuth } from '@/stores/auth'
import { setLocale } from '@/i18n'

const router = useRouter()
const auth = useAuth()
const { t, locale } = useI18n()
const epLocale = ref(locale.value === 'zh' ? zhCn : enLocale)
const form = ref(import.meta.env.DEV ? { email: 'admin@ukteahouse.co.uk', password: 'Admin!Tea2026' } : { email: '', password: '' })
const loading = ref(false)

async function submit() {
  loading.value = true
  try { await auth.login(form.value.email, form.value.password); ElMessage.success(t('login.welcome_back')); router.push('/') }
  catch { ElMessage.error(t('login.invalid_credentials')) }
  finally { loading.value = false }
}

function switchLocale(lang: string) {
  setLocale(lang)
  epLocale.value = lang === 'zh' ? zhCn : enLocale
}
</script>

<template>
  <el-config-provider :locale="epLocale">
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-amber-900 to-stone-800">
      <div class="bg-white rounded-2xl shadow-2xl p-10 w-full max-w-md">
        <div class="text-center mb-8">
          <div class="text-5xl mb-2">🍃</div>
          <h1 class="font-serif text-2xl text-slate-900">{{ t('login.brand') }}</h1>
          <p class="text-sm text-slate-500 mt-2">{{ t('login.subtitle') }}</p>
        </div>
        <el-form :model="form" label-position="top" @submit.prevent="submit">
          <el-form-item :label="t('login.email')"><el-input v-model="form.email" size="large" /></el-form-item>
          <el-form-item :label="t('login.password')"><el-input v-model="form.password" type="password" show-password size="large" /></el-form-item>
          <el-button type="primary" size="large" :loading="loading" class="w-full bg-slate-900 hover:bg-slate-800" @click="submit">{{ t('login.sign_in') }}</el-button>
        </el-form>
        <div class="flex items-center justify-between mt-6">
          <p class="text-xs text-slate-400">{{ t('login.dev_note') }}</p>
          <el-dropdown trigger="click" @command="switchLocale">
            <el-button size="small" text>
              <span class="text-xs">{{ locale === 'zh' ? '中文' : 'English' }}</span>
              <el-icon class="ml-1"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="zh" :disabled="locale === 'zh'">中文</el-dropdown-item>
                <el-dropdown-item command="en" :disabled="locale === 'en'">English</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>
  </el-config-provider>
</template>
