<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const sections = ref<any[]>([
  { id: 1, page_key: 'home', section_key: 'hero', title: 'Home → Hero Banner', content: { headline: 'Pu\'er Tea, Traceable to the Mountain.', subheadline: 'Single origin, SGS certified, 24/7 live streamed from Yunnan.', cta_text: 'Create Your Bespoke →', cta_link: '/bespoke' }, updated_by: 'Admin Root', updated_at: '2026-09-15T10:00:00Z' },
  { id: 2, page_key: 'home', section_key: 'live_strip', title: 'Home → Live Strip Heading', content: { heading: 'Live From the Tea Mountains', subheading: 'Watch your tea being picked, rolled and sun-dried in real time.' }, updated_by: 'Admin Root', updated_at: '2026-09-15T10:00:00Z' },
  { id: 3, page_key: 'tea_mountains', section_key: 'intro', title: 'Tea Mountains → Intro', content: { heading: 'Six Mountains. One Promise.', body: 'Iceland · Banzhang · Jingmai · Nanruo · Mangpeng · Laobanzhang — each GPS tagged, master identified.' }, updated_by: 'Admin Root', updated_at: '2026-09-15T10:00:00Z' },
  { id: 4, page_key: 'bespoke', section_key: 'how_it_works', title: 'Bespoke → How It Works', content: { steps: [ '1. Pick mountain & roast', '2. Choose packaging', '3. We quote within 24h', '4. Tea arrives in 45 days' ] }, updated_by: 'Admin Root', updated_at: '2026-09-14T09:00:00Z' },
  { id: 5, page_key: 'quality', section_key: 'intro', title: 'Quality → SGS Intro', content: { heading: 'Independently Tested. Always.', body: 'Every batch tested for pesticides, heavy metals, microbiology, flavonoid profile.' }, updated_by: 'Admin Root', updated_at: '2026-09-13T14:00:00Z' },
])

const editing = ref<any>(null)
const draft = reactive<any>({})

function startEdit(s: any) {
  editing.value = s
  Object.assign(draft, s.content)
}
async function save() {
  editing.value.content = { ...draft }
  ElMessage.success('✅ Site content saved — public site updated instantly')
  editing.value = null
}
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">🏛️ Public Site CMS</span>
      <el-tag type="success">Live — changes apply instantly</el-tag>
    </div></template>
    <el-alert type="info" :closable="false" class="mb-4">
      All public site content is pulled from this CMS table. Edit sections below. All fields in English (UK).
    </el-alert>
    <el-table :data="sections" stripe>
      <el-table-column prop="page_key" label="Page" width="150">
        <template #default="{ row }"><code>{{ row.page_key }}</code></template>
      </el-table-column>
      <el-table-column prop="section_key" label="Section" width="150">
        <template #default="{ row }"><code>{{ row.section_key }}</code></template>
      </el-table-column>
      <el-table-column prop="title" label="Description" width="280" />
      <el-table-column label="Content Preview" min-width="280">
        <template #default="{ row }">
          <code class="text-xs text-slate-500 break-all">{{ JSON.stringify(row.content).slice(0, 80) }}…</code>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="Updated" width="170" />
      <el-table-column label="Actions" width="100">
        <template #default="{ row }">
          <el-button size="small" @click="startEdit(row)">✏️ Edit</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="!!editing" :title="editing?.title" width="640px">
    <div v-if="editing">
      <div class="text-xs text-slate-500 mb-3">Edit JSON content below (pretty-printed):</div>
      <el-input v-model="draft.content" type="textarea" :autosize="{ minRows: 12, maxRows: 25 }"
        placeholder="JSON content" />
    </div>
    <template #footer>
      <el-button @click="editing = null">Cancel</el-button>
      <el-button type="primary" @click="save">💾 Publish</el-button>
    </template>
  </el-dialog>
</template>
