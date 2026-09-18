<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()
const room = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const visibility = ref<string>('registered')

const roomId = computed(() => route.params.id as string | undefined)

onMounted(async () => {
  if (!roomId.value) { error.value = 'Missing room ID'; loading.value = false; return }
  try {
    // 带 auth header 尝试访问（后端 LiveAccessMiddleware 会处理可见性）
    const data: any = await api.get(`/live-rooms/${roomId.value}`)
    room.value = data
    visibility.value = data?.visibility || 'registered'
  } catch (err: any) {
    const status = err?.code || 0
    if (status === 403 || (typeof err === 'string' && err.includes('not accessible'))) {
      // 403 → restricted
      visibility.value = 'restricted'
      error.value = 'restricted'
    } else if (status === 401) {
      // 401 → registered + 未登录
      visibility.value = 'registered'
      error.value = 'login_required'
    } else {
      error.value = err?.message || 'Failed to load broadcast'
    }
  } finally {
    loading.value = false
  }
})

const isLoggedIn = () => !!localStorage.getItem('user_token')

function goLogin() {
  const redirect = encodeURIComponent(route.fullPath)
  router.push({ path: '/magic-link', query: { redirect } })
}
</script>

<template>
  <div class="detail-root">
    <!-- ========== 加载中 ========== -->
    <div v-if="loading" class="lux-loading">
      <span class="lux-loading-mark uppercase-caps">Curating · Session</span>
    </div>

    <!-- ========== Restricted + 未登录 ========== -->
    <div v-else-if="error === 'restricted' && !isLoggedIn()" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">This broadcast is not publicly accessible</h1>
        <p class="lux-denied-body">
          Your advisor has reserved this session for specific members.
          Sign in with your registered email to request entry.
        </p>
        <button class="lux-btn lux-btn-primary" @click="goLogin">
          <span class="uppercase-caps">Request · Access</span>
        </button>
        <p class="lux-denied-sub">A member of our advisory team will respond within 24 hours.</p>
      </div>
    </div>

    <!-- ========== Restricted + 已登录但无权限 ========== -->
    <div v-else-if="error === 'restricted' && isLoggedIn()" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">Private · Broadcast</h1>
        <p class="lux-denied-body">
          Your advisor has not yet granted access to this session.
          Reach out directly via the Conversation panel — they will add you to the guest list.
        </p>
        <button class="lux-btn lux-btn-ghost" @click="router.push('/chat')">
          <span class="uppercase-caps">Speak with · Your Advisor</span>
        </button>
      </div>
    </div>

    <!-- ========== Registered + 未登录 ========== -->
    <div v-else-if="error === 'login_required'" class="lux-denied">
      <div class="lux-denied-card">
        <h1 class="lux-denied-title uppercase-caps">Sign · In · Required</h1>
        <p class="lux-denied-body">
          This session is available to registered members. Sign in with your magic link to view.
        </p>
        <button class="lux-btn lux-btn-primary" @click="goLogin">
          <span class="uppercase-caps">Sign · In</span>
        </button>
      </div>
    </div>

    <!-- ========== 正常显示 ========== -->
    <div v-else-if="room" class="lux-live">
      <header class="lux-live-bar">
        <span class="lux-live-dot"></span>
        <span class="lux-live-title uppercase-caps">{{ room.room_name || 'Private Session' }}</span>
        <span class="lux-live-visibility lux-live-visibility--{{ visibility }} uppercase-caps">{{ visibility }}</span>
      </header>
      <div class="lux-live-stage">
        <div class="lux-video-placeholder" style="
          background: #0B0A09;
          display: flex; align-items: center; justify-content: center;
          font-family: 'Cormorant Garamond', serif;
          color: #C5A572; font-size: 18px; letter-spacing: 0.15em;
          border-radius: 2px;
          min-height: 400px;
        ">
          Live stream · Connecting
        </div>
      </div>
      <p class="lux-live-desc">{{ room.description || 'A curated private session.' }}</p>
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;500;600&family=Inter:wght@300;400&display=swap');

.detail-root {
  min-height: 100vh;
  background: #F8F5EF;
  font-family: 'Inter', sans-serif;
  color: #1a1714;
  padding: 48px 32px;
  max-width: 960px;
  margin: 0 auto;
}
.uppercase-caps { text-transform: uppercase; letter-spacing: 0.28em; }

.lux-loading {
  height: 60vh; display: flex; align-items: center; justify-content: center;
}
.lux-loading-mark {
  font-family: 'Cormorant Garamond', serif;
  font-size: 14px; color: #C5A572;
}

.lux-denied {
  height: 60vh; display: flex; align-items: center; justify-content: center;
}
.lux-denied-card {
  text-align: center;
  padding: 64px 48px;
  border: 1px solid #D8D0C4;
  border-radius: 2px;
  background: #F8F5EF;
  max-width: 520px;
}
.lux-denied-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 16px; font-weight: 500;
  color: #C5A572;
  margin: 0 0 24px;
}
.lux-denied-body {
  font-size: 15px; line-height: 1.7;
  color: #44403a;
  font-family: 'Cormorant Garamond', serif;
  margin: 0 0 32px;
}
.lux-denied-sub {
  font-size: 12px; color: #8a8578;
  margin-top: 20px;
}

.lux-btn {
  padding: 14px 32px;
  border-radius: 2px;
  font-size: 11px; font-weight: 400;
  cursor: pointer;
  transition: all 600ms ease-out;
  letter-spacing: 0.15em;
  font-family: 'Inter', sans-serif;
}
.lux-btn-primary {
  background: #0B0A09; color: #C5A572;
  border: 1px solid #0B0A09;
}
.lux-btn-primary:hover { background: #2A2520; }
.lux-btn-ghost {
  background: transparent; color: #0B0A09;
  border: 1px solid #0B0A09;
}
.lux-btn-ghost:hover { background: rgba(11,10,9,0.05); }
.lux-btn-label { font-family: 'Cormorant Garamond', serif; font-size: 13px; }

.lux-live-bar {
  display: flex; align-items: center; gap: 14px;
  padding-bottom: 20px;
  border-bottom: 1px solid #D8D0C4;
  margin-bottom: 24px;
}
.lux-live-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: #C5A572;
}
.lux-live-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 14px; color: #0B0A09;
}
.lux-live-visibility {
  margin-left: auto;
  font-size: 10px; color: #8a8578;
}
.lux-live-visibility--restricted { color: #B4645A; }
.lux-live-visibility--public { color: #4a7c59; }
.lux-live-stage { margin-bottom: 20px; }
.lux-live-desc {
  font-family: 'Cormorant Garamond', serif;
  font-size: 16px; line-height: 1.6;
  color: #44403a;
}
</style>
