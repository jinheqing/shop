<script setup lang="ts">
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

function statusClass(s: string) {
  switch (s) {
    case 'completed': return 'bg-emerald-100 text-emerald-800'
    case 'processing': return 'bg-amber-100 text-amber-800'
    case 'pending': return 'bg-tea-100 text-tea-700'
    default: return 'bg-slate-100 text-slate-700'
  }
}
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
    // 未登录时后端会返回 401 — 静默处理
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
    successMsg.value = '✅ Your GDPR request has been submitted. Reference will appear once processed.'
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
  <div class="pt-20 py-16 bg-tea-50">
    <div class="max-w-4xl mx-auto px-6">
      <!-- Hero -->
      <div class="mb-12 text-center">
        <div class="text-xs tracking-[0.25em] uppercase text-tea-500 mb-3">GDPR · Data Subject Access Rights</div>
        <h1 class="font-display font-semibold text-4xl md:text-5xl text-tea-900 mb-4">Your Data, <span class="italic">Your Rights</span></h1>
        <p class="text-tea-700 max-w-2xl mx-auto">
          Under the EU General Data Protection Regulation (GDPR) and UK GDPR, you have clear rights over your personal data. Submit a request below and our Data Protection Officer will respond within 30 calendar days.
        </p>
      </div>

      <!-- Rights grid -->
      <section class="mb-12 grid sm:grid-cols-2 gap-4">
        <div v-for="r in [
          { title: 'Access', icon: '👁️', desc: 'Get a copy of all personal data we hold about you.' },
          { title: 'Rectification', icon: '✏️', desc: 'Correct inaccurate or incomplete data.' },
          { title: 'Erasure (Right to be Forgotten)', icon: '🗑️', desc: 'Request deletion of your personal data, subject to legal exceptions.' },
          { title: 'Portability', icon: '📦', desc: 'Receive your data in a structured, machine-readable format, or request transfer to another controller.' },
          { title: 'Restriction of Processing', icon: '⏸️', desc: 'Limit how we process your data while a dispute is resolved.' },
          { title: 'Object', icon: '🛑', desc: 'Object to certain types of processing, including direct marketing.' },
        ]" :key="r.title"
             class="rounded-2xl border border-tea-100 bg-white p-5">
          <div class="flex items-start gap-3 mb-2">
            <span class="text-2xl leading-none">{{ r.icon }}</span>
            <h3 class="font-display text-lg text-tea-900">{{ r.title }}</h3>
          </div>
          <p class="text-sm text-tea-700 leading-relaxed">{{ r.desc }}</p>
        </div>
      </section>

      <!-- Form -->
      <section class="mb-12 bg-white rounded-2xl border border-tea-100 p-6 md:p-8">
        <h2 class="font-display text-2xl text-tea-900 mb-6">Submit a data request</h2>

        <div v-if="errorMsg" class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-xl text-sm">{{ errorMsg }}</div>
        <div v-if="successMsg" class="mb-4 p-3 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-xl text-sm">{{ successMsg }}</div>

        <div class="mb-5">
          <label class="block text-sm font-medium text-tea-900 mb-2">Type of request</label>
          <select v-model="form.request_type" class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none bg-white">
            <option value="access">Access — get a copy of your data</option>
            <option value="rectification">Rectification — correct your data</option>
            <option value="erasure">Erasure — delete your data</option>
            <option value="portability">Portability — get a machine-readable copy</option>
          </select>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-tea-900 mb-2">Reason or additional details (optional)</label>
          <textarea v-model="form.reason" rows="3" placeholder="Tell us more about what you'd like us to do..." class="w-full px-4 py-3 rounded-xl border border-tea-200 focus:border-tea-600 focus:outline-none resize-none"></textarea>
        </div>

        <button @click="submit" :disabled="submitting"
          class="w-full py-3 bg-tea-800 text-white rounded-xl font-medium hover:bg-tea-900 transition disabled:opacity-50 disabled:cursor-not-allowed">
          {{ submitting ? 'Submitting…' : 'Submit request →' }}
        </button>

        <p class="text-xs text-tea-500 mt-4 text-center">
          Submitting this form creates a ticket tracked in our audit log. SLA: 30 calendar days from receipt.
        </p>
      </section>

      <!-- My requests -->
      <section class="bg-white rounded-2xl border border-tea-100 p-6 md:p-8">
        <div class="flex items-center justify-between mb-4">
          <h2 class="font-display text-2xl text-tea-900">My requests</h2>
          <button @click="load" class="text-sm text-tea-700 hover:text-tea-950">Refresh</button>
        </div>

        <div v-if="loading" class="text-center py-10 text-tea-500 text-sm">Loading…</div>
        <div v-else-if="!requests.length" class="text-center py-10 text-tea-500 text-sm">No requests yet.</div>

        <table v-else class="w-full text-sm">
          <thead class="text-xs text-tea-500 border-b border-tea-100">
            <tr>
              <th class="text-left py-2 font-medium">Type</th>
              <th class="text-left py-2 font-medium">Created</th>
              <th class="text-left py-2 font-medium">Due</th>
              <th class="text-left py-2 font-medium">Status</th>
            </tr>
          </thead>
          <tbody class="text-tea-800">
            <tr v-for="r in requests" :key="r.id" class="border-b border-tea-50">
              <td class="py-3 capitalize">{{ r.request_type }}</td>
              <td class="py-3">{{ new Date(r.created_at).toLocaleDateString('en-GB') }}</td>
              <td class="py-3">{{ new Date(r.due_at).toLocaleDateString('en-GB') }} <span class="text-xs text-tea-500">({{ daysUntil(r.due_at) }})</span></td>
              <td class="py-3">
                <span :class="['text-xs px-2 py-1 rounded-full font-medium', statusClass(r.status)]">{{ statusLabel(r.status) }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <!-- Footer note -->
      <div class="mt-12 text-center text-xs text-tea-500 space-y-1">
        <p>Alternatively, you can email your request to <a href="mailto:privacy@ukteahouse.co.uk" class="underline hover:text-tea-900">privacy@ukteahouse.co.uk</a>.</p>
        <p>If you are not satisfied with our response, you have the right to lodge a complaint with the <a href="https://ico.org.uk" target="_blank" rel="noopener" class="underline hover:text-tea-900">UK Information Commissioner's Office</a>.</p>
      </div>
    </div>
  </div>
</template>
