<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">CMS — Site Content</span>
          <el-button type="primary" :icon="Refresh" @click="load">Reload</el-button>
        </div>
      </template>
      <el-alert v-if="!loaded" title="Click Reload to fetch from backend" type="info" show-icon :closable="false" class="mb-4" />
      <el-table :data="contents" v-else stripe>
        <el-table-column prop="page_key" label="Page" width="160" />
        <el-table-column prop="section_key" label="Section" width="180" />
        <el-table-column label="Content Preview" min-width="260">
          <template #default="{ row }">
            <span class="text-xs text-gray-500">{{ JSON.stringify(row.content).slice(0, 120) }}{{ JSON.stringify(row.content).length > 120 ? '…' : '' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="Updated" width="180" />
        <el-table-column label="Actions" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">Edit</el-button>
            <el-button size="small" type="success" @click="reset(row)">Reset to defaults</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog v-model="dialog" :title="`Edit ${editing?.page_key}/${editing?.section_key}`" width="640px">
      <div v-if="editing">
        <div class="mb-2 text-xs text-gray-400">Edit JSON content directly</div>
        <el-input v-model="draft" type="textarea" :rows="12" />
      </div>
      <template #footer>
        <el-button @click="dialog = false">Cancel</el-button>
        <el-button type="primary" @click="save">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { api } from '../api/client'

const contents = ref<any[]>([])
const loaded = ref(false)
const dialog = ref(false)
const editing = ref<any>(null)
const draft = ref('')

async function load() {
  try {
    const res: any = await api.get('/site-contents')
    // res can be array or { items: [...] }
    contents.value = Array.isArray(res) ? res : (res?.items || [])
    loaded.value = true
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

function openEdit(row: any) {
  editing.value = row
  draft.value = JSON.stringify(row.content, null, 2)
  dialog.value = true
}

async function save() {
  try {
    const parsed = JSON.parse(draft.value)
    const url = editing.value.id
      ? `/site-contents/${editing.value.id}`
      : '/site-contents'
    await api.put(url, {
      page_key: editing.value.page_key,
      section_key: editing.value.section_key,
      content: parsed
    })
    ElMessage.success('Saved')
    dialog.value = false
    await load()
  } catch (e: any) {
    ElMessage.error('Save failed: ' + (e?.message || e))
  }
}

async function reset(row: any) {
  try {
    const defaults: Record<string, any> = {
      'home/hero': { title: 'Pu\'er Tea Direct from Yunnan', subtitle: 'Personalized tea experience' },
      'home/featured_presets': { items: [] },
      'home/quality_sgs': { title: 'Third-Party Lab Certified', items: [] },
      'about/story': { paragraph: 'UK-based Pu\'er tea specialist since 2024.' },
      'checkout/contact_us': { email: 'hello@ukteahouse.co.uk', phone: '+44 ...' },
      'footer/links': { privacy: '/privacy', terms: '/terms', gdpr_dsar: '/gdpr-dsar' }
    }
    const def = defaults[`${row.page_key}/${row.section_key}`] || {}
    await api.put(`/site-contents/${row.id}`, { ...row, content: def })
    ElMessage.success('Reset to defaults')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.message || 'Reset failed')
  }
}

onMounted(load)
</script>
