<script setup lang="ts">
// ============================================================
// GdprDsar.vue — 老钱审美 GDPR DSAR 页
// - 全部去掉 emoji（旧版 👁️ ✏️ 🗑️ 📦 ⏸️ 🛑 ✅ 已删）
// - 状态标签改用 in-page ink/ivory/gold 风格
// - 表单/列表改用 squared border + hairline dividers
// ============================================================
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'

const requests = ref<any[]>([])
const loading = ref(false)
const submitting = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const form = ref({
  request_type: 'access',
  reason: '',
})

function statusLabel(s: string) {
  return s ? s.charAt(0).toUpperCase() + s.slice(1) : '—'
}
function daysUntil(due?: string) {
  if (!due) return ''
  const d = new Date(due).getTime() - Date.now()
  if (d < 0) return 'overdue'
  return Math.ceil(d / 86400000) + 'd left'
}

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const r: any = await api.get('/dsar/requests')
    requests.value = (r?.items || r || []) as any[]
  } catch (e: any) {
    errorMsg.value = e?.response?.status === 401
      ? 'Please sign in to view your data requests.'
      : e?.response?.data?.error || e.message
  } finally {
    loading.value = false
  }
}

async function submit() {
  submitting.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await api.post('/dsar/requests', form.value)
    successMsg.value = 'Your GDPR request has been submitted. Reference will appear once processed.'
    form.value = { request_type: 'access', reason: '' }
    await load()
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.error || e.message || 'Submission failed'
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="pt-20 py-16 bg-ivory-100">
    <div class="max-w-4xl mx-auto px-6">

      <!-- Hero -->
      <div class="mb-12 text-center">
        <div class="flex items-center gap-3 justify-center mb-5">
          <span class="h-px w-10 bg-gold/50"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">GDPR · Data Subject Access Rights</span>
          <span class="h-px w-10 bg-gold/50"></span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-5 leading-tight">
          Your Data, <span class="italic">Your Rights</span>
        </h1>
        <p class="text-sand font-serif max-w-2xl mx-auto leading-relaxed">
          Under the EU General Data Protection Regulation (GDPR) and UK GDPR, you have clear rights over your personal data. Submit a request below and our Data Protection Officer will respond within 30 calendar days.
        </p>
      </div>

      <!-- Rights — narrative list, no emoji -->
      <section class="mb-12">
        <div class="grid sm:grid-cols-2 gap-x-12 gap-y-8">
          <div v-for="r in [
            { title: 'Access', desc: 'Get a copy of all personal data we hold about you.' },
            { title: 'Rectification', desc: 'Correct inaccurate or incomplete data.' },
            { title: 'Erasure (Right to be Forgotten)', desc: 'Request deletion of your personal data, subject to legal exceptions.' },
            { title: 'Portability', desc: 'Receive your data in a structured, machine-readable format, or request transfer to another controller.' },
            { title: 'Restriction of Processing', desc: 'Limit how we process your data while a dispute is resolved.' },
            { title: 'Object', desc: 'Object to certain types of processing, including direct marketing.' },
          ]" :key="r.title">
            <div class="flex items-center gap-3 mb-2">
              <span class="w-1 h-1 rounded-full bg-gold"></span>
              <h3 class="font-serif text-base text-ink-900">{{ r.title }}</h3>
            </div>
            <p class="text-sm text-sand font-serif leading-relaxed pl-4">{{ r.desc }}</p>
          </div>
        </div>
      </section>

      <!-- Form -->
      <section class="mb-12 bg-white border border-gold/20 p-6 md:p-8" style="border-radius: 2px;">
        <div class="flex items-center gap-3 mb-6">
          <span class="h-px w-8 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Submit · A · Data · Request</span>
        </div>

        <!-- In-page status messages -->
        <div v-if="errorMsg" class="mb-5 px-4 py-3 bg-ink-900 text-ivory-100 text-sm font-serif leading-relaxed" style="border-radius: 2px;">
          <div class="flex items-center gap-2">
            <span class="w-1 h-1 rounded-full bg-gold"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Notice</span>
          </div>
          <p class="mt-1 text-ivory-100/80">{{ errorMsg }}</p>
        </div>
        <div v-if="successMsg" class="mb-5 px-4 py-3 border border-gold/30 bg-ivory-100 text-sm text-ink-900 font-serif leading-relaxed" style="border-radius: 2px;">
          <div class="flex items-center gap-2 mb-1">
            <span class="w-1 h-1 rounded-full bg-gold"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Acknowledged</span>
          </div>
          <p>{{ successMsg }}</p>
        </div>

        <div class="mb-6">
          <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Type · Of · Request</label>
          <select v-model="form.request_type"
            class="w-full px-4 py-3 border border-gold/20 bg-ivory-50 text-ink-900 font-serif focus:border-gold focus:outline-none"
            style="border-radius: 2px;">
            <option value="access">Access — get a copy of your data</option>
            <option value="rectification">Rectification — correct your data</option>
            <option value="erasure">Erasure — delete your data</option>
            <option value="portability">Portability — get a machine-readable copy</option>
          </select>
        </div>

        <div class="mb-7">
          <label class="block text-[11px] uppercase tracking-lux text-gold font-sans mb-3">Reason · Or · Additional · Details · (Optional)</label>
          <textarea v-model="form.reason" rows="3" placeholder="Tell us more about what you'd like us to do…"
            class="w-full px-4 py-3 border border-gold/20 bg-ivory-50 text-ink-900 font-serif focus:border-gold focus:outline-none resize-none"
            style="border-radius: 2px;"></textarea>
        </div>

        <button @click="submit" :disabled="submitting"
          class="w-full py-4 bg-ink-900 text-ivory-100 hover:bg-ink-800 transition-duration-lux text-[11px] uppercase tracking-lux font-sans disabled:opacity-50 disabled:cursor-not-allowed"
          style="border-radius: 2px;">
          {{ submitting ? 'Submitting · · ·' : 'Submit · Request' }}
        </button>

        <p class="text-xs text-sand/70 mt-4 text-center font-serif">
          Submitting this form creates a ticket tracked in our audit log. SLA: 30 calendar days from receipt.
        </p>
      </section>

      <!-- My requests -->
      <section class="bg-white border border-gold/20 p-6 md:p-8" style="border-radius: 2px;">
        <div class="flex items-center justify-between mb-5">
          <div class="flex items-center gap-3">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">My · Requests</span>
          </div>
          <button @click="load" class="text-[11px] uppercase tracking-lux text-sand hover:text-gold font-sans">Refresh</button>
        </div>

        <div v-if="loading" class="text-center py-10 text-sand text-sm font-serif">Loading…</div>
        <div v-else-if="!requests.length" class="text-center py-10 text-sand text-sm font-serif">No requests yet.</div>

        <table v-else class="w-full text-sm">
          <thead>
            <tr class="border-b border-gold/20">
              <th class="text-left py-3 text-[10px] uppercase tracking-lux text-gold font-sans">Type</th>
              <th class="text-left py-3 text-[10px] uppercase tracking-lux text-gold font-sans">Created</th>
              <th class="text-left py-3 text-[10px] uppercase tracking-lux text-gold font-sans">Due</th>
              <th class="text-left py-3 text-[10px] uppercase tracking-lux text-gold font-sans">Status</th>
            </tr>
          </thead>
          <tbody class="text-ink-900 font-serif">
            <tr v-for="r in requests" :key="r.id" class="border-b border-gold/10">
              <td class="py-3 capitalize text-sand">{{ r.request_type }}</td>
              <td class="py-3 text-sand">{{ new Date(r.created_at).toLocaleDateString('en-GB') }}</td>
              <td class="py-3 text-sand">{{ new Date(r.due_at).toLocaleDateString('en-GB') }} <span class="text-xs text-sand/70">({{ daysUntil(r.due_at) }})</span></td>
              <td class="py-3">
                <span class="inline-flex items-center gap-2 text-[11px] uppercase tracking-lux font-sans"
                  :class="{
                    'text-gold': r.status === 'processing',
                    'text-ink-900': r.status === 'completed',
                    'text-sand': !['processing', 'completed'].includes(r.status),
                  }">
                  <span class="w-1 h-1 rounded-full"
                    :class="{
                      'bg-gold': r.status === 'processing',
                      'bg-ink-900': r.status === 'completed',
                      'bg-sand': !['processing', 'completed'].includes(r.status),
                    }"></span>
                  {{ statusLabel(r.status) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <!-- Footer note -->
      <div class="mt-12 text-center text-xs text-sand font-serif space-y-2">
        <p>Alternatively, you can email your request to <a href="mailto:privacy@ukteahouse.co.uk" class="text-gold underline hover:text-ink-900">privacy@ukteahouse.co.uk</a>.</p>
        <p>If you are not satisfied with our response, you have the right to lodge a complaint with the <a href="https://ico.org.uk" target="_blank" rel="noopener" class="text-gold underline hover:text-ink-900">UK Information Commissioner's Office</a>.</p>
      </div>
    </div>
  </div>
</template>
