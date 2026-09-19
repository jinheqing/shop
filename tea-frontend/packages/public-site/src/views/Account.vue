<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'

const router = useRouter()
const tab = ref<'orders' | 'chat' | 'broadcasts' | 'garden' | 'referrals'>('orders')

const orders = ref<any[]>([])
const quotes = ref<any[]>([])
const conversations = ref<any[]>([])
const recordings = ref<any[]>([])
const userGroups = ref<any[]>([])
const myReferrals = ref<any>({ items: [], stats: {}, referral_code: '', referral_url: '' })

// Garden privileges derived from userGroups
const gardenPrivileges = computed(() => {
  const list: Array<{ type: string; label: string; sourceGroup: string }> = []
  for (const g of userGroups.value) {
    const privs = g.privileges || []
    for (const p of privs) {
      if (typeof p === 'string') {
        if (p.startsWith('offline_garden_tour:')) {
          list.push({ type: 'garden_tour', label: `Garden Tour · ${p.split(':')[1]}`, sourceGroup: g.name })
        } else if (p === 'private_masterclass') {
          list.push({ type: 'masterclass', label: 'Private Masterclass · Invitation', sourceGroup: g.name })
        } else if (p === 'offline_invite') {
          list.push({ type: 'dinner', label: 'Annual Private Dinner · Invitation', sourceGroup: g.name })
        }
      }
    }
  }
  return list
})

async function load() {
  // Orders + Quotes
  try { orders.value = (await api.get('/orders') as any)?.items || [] } catch {}
  try { quotes.value = (await api.get('/custom-products/published') as any)?.items || [] } catch {}

  // Conversations
  try { conversations.value = (await api.get('/conversations') as any)?.items || [] } catch {}

  // Recordings (会员可见的)
  try { recordings.value = (await api.get('/my/recordings') as any) || [] } catch {}

  // User groups
  try { userGroups.value = (await api.get('/users/me/groups') as any) || [] } catch {}

  // My referrals (老客户看自己推荐了谁)
  try { myReferrals.value = (await api.get('/user/referrals') as any) || { items: [], stats: {}, referral_code: '', referral_url: '' } } catch {}
}
onMounted(load)

function logout() { localStorage.removeItem('user_token'); router.push('/') }

function resolveUrl(url: string) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  return base + url
}
</script>

<template>
  <div class="ac-root pt-20">
    <!-- Header -->
    <header class="ac-head">
      <div class="ac-head-inner">
        <div>
          <h1 class="ac-title">My · Account</h1>
          <p class="ac-sub uppercase-caps">A curated record of your journey with UK Tea House</p>
        </div>
        <div class="ac-head-right">
          <div v-if="userGroups.length" class="ac-badge-group">
            <span v-for="g in userGroups" :key="g.id" class="ac-badge">{{ g.name }}</span>
          </div>
          <button @click="logout" class="ac-logout uppercase-caps">Sign · Out</button>
        </div>
      </div>
    </header>

    <!-- Tabs -->
    <nav class="ac-tabs">
      <button v-for="t in [
        { key: 'orders', label: 'Orders & Bespoke' },
        { key: 'chat', label: 'Conversations' },
        { key: 'broadcasts', label: 'Private Broadcasts' },
        { key: 'garden', label: 'Garden & Invitations' },,
        { key: 'referrals', label: 'Share · The · Tea' },
      ]" :key="t.key"
        :class="['ac-tab', tab === t.key ? 'ac-tab--active' : '']"
        @click="tab = t.key">
        <span class="uppercase-caps">{{ t.label }}</span>
      </button>
    </nav>

    <main class="ac-body">
      <!-- ====== Tab 1: Orders & Bespoke ====== -->
      <section v-if="tab === 'orders'">
        <!-- Orders -->
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Orders</h2>
          <div v-if="orders.length === 0" class="ac-empty">
            <p>Your first order begins with a conversation.</p>
            <router-link to="/bespoke" class="ac-btn ac-btn-primary">Begin · Bespoke</router-link>
          </div>
          <div v-else class="ac-card-grid">
            <article v-for="o in orders" :key="o.id" class="ac-card">
              <header class="ac-card-head">
                <span class="ac-mono">{{ o.order_no }}</span>
                <span class="ac-state">{{ o.state }}</span>
              </header>
              <h3 class="ac-card-title">{{ o.custom_product_snapshot?.title || 'Bespoke Pu\'er' }}</h3>
              <p class="ac-card-meta">{{ o.custom_product_snapshot?.raw_tea_source?.slice(0, 60) || '' }}</p>
              <footer class="ac-card-foot">
                <span class="ac-price">£{{ o.total_amount }}</span>
                <a :href="`/orders/${o.id}/invoice`" class="ac-link">Invoice →</a>
              </footer>
            </article>
          </div>
        </div>

        <!-- Published Bespoke Quotes -->
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Your · Bespoke · Quotes</h2>
          <div v-if="quotes.length === 0" class="ac-empty">
            <p>No active bespoke quotes yet. Ask your advisor to prepare one.</p>
          </div>
          <div v-else class="ac-card-grid">
            <article v-for="q in quotes" :key="q.id" class="ac-card ac-card--quote">
              <h3 class="ac-card-title">{{ q.title }}</h3>
              <p class="ac-card-meta">{{ q.tea_type }} · {{ q.tea_shape }} · {{ q.unit_price ? '£'+q.unit_price+' / '+q.tea_shape_weight+'g' : '' }}</p>
              <footer class="ac-card-foot">
                <a :href="`/#/quote/${q.product_token}`" class="ac-link">View Quote →</a>
              </footer>
            </article>
          </div>
        </div>
      </section>

      <!-- ====== Tab 2: Chat History ====== -->
      <section v-if="tab === 'chat'">
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Conversations</h2>
          <div v-if="conversations.length === 0" class="ac-empty">
            <p>Your advisory conversations appear here.</p>
            <router-link to="/chat" class="ac-btn ac-btn-primary">Speak · With · Your · Advisor</router-link>
          </div>
          <div v-else class="ac-list">
            <a v-for="c in conversations" :key="c.id" href="#/chat" class="ac-list-item">
              <span class="ac-list-dot"></span>
              <span class="ac-list-label">{{ c.title || 'Private Conversation' }}</span>
            </a>
          </div>
        </div>
      </section>

      <!-- ====== Tab 3: Private Broadcasts & Recordings ====== -->
      <section v-if="tab === 'broadcasts'">
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Private · Broadcasts · & · Recordings</h2>
          <div v-if="recordings.length === 0" class="ac-empty">
            <p>No recordings yet. Your advisor will share private sessions here.</p>
          </div>
          <div v-else class="ac-card-grid">
            <article v-for="r in recordings" :key="r.id" class="ac-card ac-card--broadcast">
              <header class="ac-card-head">
                <span class="ac-rec-dot"></span>
                <span class="ac-mono">{{ r.live_room?.room_name || r.title || 'Private Session' }}</span>
              </header>
              <div class="ac-rec-meta">
                <span v-if="r.duration_sec">{{ Math.floor(r.duration_sec / 60) }}m</span>
                <span>{{ r.visibility }}</span>
              </div>
              <footer class="ac-card-foot">
                <a :href="resolveUrl(r.file_url)" target="_blank" class="ac-link">▶ Watch Recording</a>
              </footer>
            </article>
          </div>
        </div>
      </section>

      <!-- ====== Tab 4: Garden Visits & Masterclass ====== -->
      <section v-if="tab === 'garden'">
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Garden · Visits · & · Invitations</h2>
          <div v-if="gardenPrivileges.length === 0" class="ac-empty">
            <p>Exclusive invitations are extended to members of our private circles.</p>
            <p class="ac-empty-sub">Speak with your advisor about joining a bespoke group.</p>
          </div>
          <div v-else class="ac-list">
            <div v-for="(item, i) in gardenPrivileges" :key="i" class="ac-list-item ac-list-item--garden">
              <span class="ac-list-dot ac-list-dot--gold"></span>
              <span class="ac-list-label">{{ item.label }}</span>
              <span class="ac-list-source uppercase-caps">{{ item.sourceGroup }}</span>
            </div>
          </div>
        </div>

        <!-- Quick Reciprocal Clubs note -->
        <div v-if="userGroups.some(g => (g.reciprocal || []).length)" class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Reciprocal · Access</h2>
          <div v-for="g in userGroups" :key="'r-'+g.id" class="ac-list">
            <div v-for="rc in (g.reciprocal || [])" :key="rc" class="ac-list-item">
              <span class="ac-list-dot ac-list-dot--gold"></span>
              <span class="ac-list-label">{{ rc.replace(/_/g, ' ') }}</span>
              <span class="ac-list-source uppercase-caps">{{ g.name }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- ====== Tab 5: Share The Tea (Referrals) ====== -->
      <section v-if="tab === 'referrals'">
        <div class="ac-section">
          <h2 class="ac-section-title uppercase-caps">Share · The · Tea</h2>
          <p class="ac-sub">Share your love for our mountain tea with friends. No discount codes, no promotions — simply a personal introduction.</p>

          <!-- 专属短链 -->
          <div v-if="myReferrals.referral_code" class="ac-referral-card">
            <div class="ac-referral-code">
              <span class="ac-mono">{{ myReferrals.referral_url }}</span>
              <button @click="navigator.clipboard?.writeText(myReferrals.referral_url)" class="ac-btn ac-btn-primary">Copy · Link</button>
            </div>
            <p class="ac-referral-hint">Share this link privately with friends via WhatsApp or email. When they sign up, your name will be attached.</p>
          </div>

          <!-- 统计 -->
          <div v-if="myReferrals.stats && myReferrals.stats.total_referred > 0" class="ac-referral-stats">
            <div>
              <div class="ac-stat-num">{{ myReferrals.stats.total_referred }}</div>
              <div class="ac-stat-label uppercase-caps">Friends Invited</div>
            </div>
            <div>
              <div class="ac-stat-num">{{ myReferrals.stats.valid_orders }}</div>
              <div class="ac-stat-label uppercase-caps">With Orders</div>
            </div>
            <div>
              <div class="ac-stat-num">{{ myReferrals.stats.rewarded }}</div>
              <div class="ac-stat-label uppercase-caps">Rewards Sent</div>
            </div>
          </div>

          <!-- 推荐列表 -->
          <div v-if="myReferrals.items && myReferrals.items.length > 0" class="ac-list">
            <div v-for="r in myReferrals.items" :key="r.id" class="ac-list-item">
              <span class="ac-list-dot ac-list-dot--gold"></span>
              <span class="ac-list-label">{{ r.referred_name }}</span>
              <span v-if="r.friend_order_amount > 0" class="ac-list-source uppercase-caps">£{{ r.friend_order_amount.toFixed(2) }}</span>
              <span v-if="r.reward_triggered_at" class="ac-list-source uppercase-caps ac-rewarded">✓ Thank you sent</span>
            </div>
          </div>

          <div v-if="!myReferrals.items || myReferrals.items.length === 0" class="ac-empty">
            <p>No referrals yet. Share your favourite tea with a friend — your name travels farther than any advertisement.</p>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer links: address + GDPR -->
    <footer class="ac-foot">
      <a href="#/privacy">Privacy · & · GDPR</a>
      <span>·</span>
      <a href="#/gdpr-dsar">Data · Access · Request</a>
      <span>·</span>
      <a href="mailto:advisor@ukteahouse.co.uk">advisor@ukteahouse.co.uk</a>
    </footer>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;500;600;700&family=Inter:wght@300;400;500&display=swap');

.ac-root { min-height: 100vh; background: #F8F5EF; color: #0B0A09; }
.uppercase-caps { text-transform: uppercase; letter-spacing: 0.28em; }

/* Header */
.ac-head {
  border-bottom: 1px solid rgba(11,10,9,0.08);
  padding: 48px 48px 0;
}
.ac-head-inner {
  max-width: 960px; margin: 0 auto;
  display: flex; align-items: flex-end; justify-content: space-between; gap: 24px;
}
.ac-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 44px; font-weight: 500;
  letter-spacing: -0.01em;
  margin: 0 0 8px;
}
.ac-sub {
  font-family: 'Inter', sans-serif;
  font-size: 11px; color: #8a8578;
  margin: 0; letter-spacing: 0.3em;
}
.ac-head-right { display: flex; flex-direction: column; align-items: flex-end; gap: 12px; }
.ac-badge-group { display: flex; gap: 8px; flex-wrap: wrap; justify-content: flex-end; }
.ac-badge {
  font-family: 'Cormorant Garamond', serif;
  font-size: 13px; font-weight: 500;
  padding: 4px 12px;
  border: 1px solid #C5A572;
  border-radius: 2px;
  color: #C5A572;
  letter-spacing: 0.08em;
}
.ac-logout {
  background: none; border: none;
  font-family: 'Inter', sans-serif;
  font-size: 11px; color: #8a8578;
  cursor: pointer; letter-spacing: 0.15em;
  transition: color 300ms ease-out;
}
.ac-logout:hover { color: #0B0A09; }

/* Tabs */
.ac-tabs {
  max-width: 960px; margin: 0 auto;
  display: flex; gap: 32px;
  border-bottom: 1px solid rgba(11,10,9,0.08);
  padding: 0 48px;
}
.ac-tab {
  background: none; border: none;
  font-family: 'Inter', sans-serif;
  font-size: 10px; color: #8a8578;
  padding: 20px 0; cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 300ms ease-out;
}
.ac-tab:hover { color: #0B0A09; }
.ac-tab--active { color: #0B0A09; border-bottom-color: #C5A572; }

/* Body */
.ac-body { max-width: 960px; margin: 0 auto; padding: 40px 48px 80px; }
.ac-section { margin-bottom: 64px; }
.ac-section-title {
  font-family: 'Inter', sans-serif;
  font-size: 10px; color: #C5A572;
  margin: 0 0 24px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(197,165,114,0.2);
}

/* Cards */
.ac-card-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 24px; }
.ac-card {
  padding: 28px;
  border: 1px solid rgba(11,10,9,0.08);
  border-radius: 2px;
  background: #FBF9F4;
  transition: border-color 600ms ease-out;
}
.ac-card:hover { border-color: rgba(197,165,114,0.5); }
.ac-card-head {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 14px;
}
.ac-mono { font-family: 'Inter', monospace; font-size: 11px; color: #8a8578; letter-spacing: 0.05em; }
.ac-state {
  font-family: 'Inter', sans-serif; font-size: 10px; color: #C5A572;
  padding: 2px 8px; border: 1px solid #C5A572; border-radius: 2px;
  text-transform: uppercase; letter-spacing: 0.2em;
}
.ac-card-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 22px; font-weight: 500; color: #0B0A09;
  margin: 0 0 6px;
}
.ac-card-meta {
  font-family: 'Cormorant Garamond', serif;
  font-size: 14px; color: #6b6459;
  margin: 0 0 18px; line-height: 1.5;
}
.ac-card-foot {
  display: flex; align-items: center; justify-content: space-between;
  padding-top: 14px;
  border-top: 1px solid rgba(11,10,9,0.06);
}
.ac-price {
  font-family: 'Cormorant Garamond', serif;
  font-size: 20px; color: #0B0A09;
}
.ac-link {
  font-size: 11px; color: #C5A572;
  text-decoration: none; letter-spacing: 0.1em;
  transition: color 300ms ease-out;
}
.ac-link:hover { color: #0B0A09; }

/* Broadcast card */
.ac-card--broadcast .ac-rec-dot {
  width: 8px; height: 8px; border-radius: 50%; background: #B4645A;
}
.ac-rec-meta {
  font-size: 11px; color: #8a8578; letter-spacing: 0.1em;
  display: flex; gap: 12px; margin: 6px 0 14px;
}

/* List items */
.ac-list { display: flex; flex-direction: column; gap: 0; }
.ac-list-item {
  display: flex; align-items: center; gap: 14px;
  padding: 18px 0;
  border-bottom: 0.5px solid rgba(11,10,9,0.06);
  text-decoration: none; color: #0B0A09;
  transition: border-color 600ms ease-out;
}
.ac-list-item:hover { border-bottom-color: #C5A572; }
.ac-list-dot { width: 6px; height: 6px; border-radius: 50%; background: #0B0A09; opacity: 0.4; }
.ac-list-dot--gold { background: #C5A572; opacity: 1; }
.ac-list-label {
  font-family: 'Cormorant Garamond', serif;
  font-size: 18px; color: #0B0A09;
  flex: 1;
}
.ac-list-source {
  font-size: 10px; color: #8a8578; letter-spacing: 0.15em;
}

/* Empty */
.ac-empty {
  text-align: center;
  padding: 48px 24px;
  border: 1px dashed rgba(11,10,9,0.1);
  border-radius: 2px;
  color: #6b6459;
  font-family: 'Cormorant Garamond', serif;
  font-size: 16px; line-height: 1.6;
}
.ac-empty-sub { font-size: 13px; color: #8a8578; margin-top: 6px; }

/* Buttons */
.ac-btn {
  display: inline-block; margin-top: 16px;
  font-family: 'Inter', sans-serif;
  font-size: 10px; font-weight: 400;
  padding: 12px 24px;
  border-radius: 2px; text-decoration: none;
  letter-spacing: 0.2em; cursor: pointer;
  transition: all 600ms ease-out;
  border: 1px solid transparent;
}
.ac-btn-primary { background: #0B0A09; color: #C5A572; }
.ac-btn-primary:hover { background: #2A2520; }

/* Referrals */
.ac-sub { font-family: 'Cormorant Garamond', serif; font-size: 14px; color: #6b6459; line-height: 1.6; margin: 0 0 28px; }
.ac-referral-card {
  padding: 28px; border: 1px solid #C5A572; border-radius: 2px;
  background: #FBF9F4; margin-bottom: 32px;
}
.ac-referral-code {
  display: flex; align-items: center; justify-content: space-between; gap: 16px; flex-wrap: wrap;
}
.ac-referral-code .ac-mono { font-size: 13px; }
.ac-referral-hint { font-family: 'Cormorant Garamond', serif; font-size: 13px; color: #8a8578; margin: 12px 0 0; }
.ac-referral-stats { display: flex; gap: 48px; margin: 0 0 32px; }
.ac-stat-num { font-family: 'Cormorant Garamond', serif; font-size: 36px; color: #0B0A09; }
.ac-stat-label { font-size: 10px; color: #8a8578; }
.ac-rewarded { color: #C5A572; }

/* Footer */
.ac-foot {
  border-top: 1px solid rgba(11,10,9,0.08);
  padding: 24px 48px;
  text-align: center;
  font-size: 11px; color: #8a8578; letter-spacing: 0.1em;
  display: flex; justify-content: center; gap: 14px; flex-wrap: wrap;
}
.ac-foot a { color: #8a8578; text-decoration: none; transition: color 300ms; }
.ac-foot a:hover { color: #0B0A09; }
</style>
