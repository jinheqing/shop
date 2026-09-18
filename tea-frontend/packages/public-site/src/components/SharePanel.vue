<script setup lang="ts">
import { ref } from 'vue'
import { api } from '@/api/client'

const props = defineProps<{
  url: string
  title?: string
  description?: string
}>()

const shortUrl = ref<string>('')
const loading = ref(false)
const copied = ref(false)

// 获取或创建短链
async function resolveShortLink() {
  if (shortUrl.value) return shortUrl.value
  loading.value = true
  try {
    // 先尝试登录后创建短链（需要 JWT）
    const token = localStorage.getItem('user_token')
    if (token) {
      const r: any = await api.post('/short-links', { target_url: props.url })
      if (r?.code) {
        // 后端可能包装 {code, data}
        shortUrl.value = r.data ? buildShortUrl(r.data.code) : buildShortUrl(r.code)
      } else if (r?.code) {
        shortUrl.value = buildShortUrl(r.code)
      } else {
        shortUrl.value = buildShortUrl(r.code || r.Code)
      }
    }
    // 登录失败或无 token，直接用完整 URL
    if (!shortUrl.value) shortUrl.value = props.url
  } catch {
    shortUrl.value = props.url
  } finally {
    loading.value = false
  }
  return shortUrl.value
}

function buildShortUrl(code: string) {
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  return `${base}/s/${code}`
}

async function copyLink() {
  const url = await resolveShortLink()
  try {
    await navigator.clipboard.writeText(url)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  } catch {
    // fallback
    const ta = document.createElement('textarea')
    ta.value = url; document.body.appendChild(ta); ta.select()
    document.execCommand('copy'); document.body.removeChild(ta)
  }
}

function shareOn(platform: 'whatsapp' | 'twitter' | 'facebook' | 'linkedin') {
  const encoded = encodeURIComponent(props.url)
  const title = encodeURIComponent(props.title || 'UK Tea House')
  const text = encodeURIComponent(`${props.title || 'UK Tea House'} — ${props.description || 'Discover bespoke Pu\'er from Yunnan'}`)

  let href = ''
  switch (platform) {
    case 'whatsapp':
      href = `https://wa.me/?text=${text}%20${encoded}`
      break
    case 'twitter':
      href = `https://twitter.com/intent/tweet?url=${encoded}&text=${text}`
      break
    case 'facebook':
      href = `https://www.facebook.com/sharer/sharer.php?u=${encoded}`
      break
    case 'linkedin':
      href = `https://www.linkedin.com/sharing/share-offsite/?url=${encoded}`
      break
  }
  window.open(href, '_blank', 'width=600,height=500')
}
</script>

<template>
  <div class="share-panel">
    <span class="share-label uppercase-caps">Share · This · Moment</span>
    <div class="share-buttons">
      <button @click="copyLink" :disabled="loading" class="share-btn share-btn-copy">
        <template v-if="copied">✓ Copied</template>
        <template v-else-if="loading">…</template>
        <template v-else>Copy · Link</template>
      </button>
      <button @click="shareOn('whatsapp')" class="share-btn share-btn-whatsapp">WhatsApp</button>
      <button @click="shareOn('twitter')" class="share-btn share-btn-twitter">X</button>
      <button @click="shareOn('facebook')" class="share-btn share-btn-facebook">Facebook</button>
    </div>
    <div v-if="shortUrl" class="share-short">
      {{ shortUrl }}
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;500;600&family=Inter:wght@300;400&display=swap');

.share-panel {
  padding: 20px 0;
  border-top: 1px solid #D8D0C4;
  display: flex; flex-direction: column; gap: 12px;
}
.share-label {
  font-family: 'Inter', sans-serif;
  font-size: 10px; font-weight: 400;
  color: #C5A572;
  letter-spacing: 0.3em;
}
.uppercase-caps { text-transform: uppercase; }

.share-buttons { display: flex; gap: 8px; flex-wrap: wrap; }

.share-btn {
  font-family: 'Inter', sans-serif;
  font-size: 11px; font-weight: 400;
  padding: 8px 16px;
  border-radius: 2px;
  border: 1px solid #2A2520;
  background: transparent;
  color: #0B0A09;
  cursor: pointer;
  letter-spacing: 0.15em;
  transition: all 600ms ease-out;
}
.share-btn:hover { background: #0B0A09; color: #F8F5EF; border-color: #0B0A09; }
.share-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.share-btn-copy { border-color: #C5A572; color: #C5A572; }
.share-btn-copy:hover { background: #C5A572; color: #0B0A09; }

.share-short {
  font-family: 'Cormorant Garamond', serif;
  font-size: 13px; color: #8a8578;
  letter-spacing: 0.05em;
}
</style>
