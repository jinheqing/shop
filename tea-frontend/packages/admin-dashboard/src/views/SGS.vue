<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
// 后端 SgsReportCreateRequest: report_no, batch_no, tea_type, test_date, issue_date, pdf_url, test_items
// 注: 后端没有 title / issuer / supplier_name 字段
const form = reactive({
  report_no: '',
  batch_no: '',
  tea_type: 'raw_puer',
  test_date: '',
  issue_date: '',
  pdf_url: '',
})

async function load() {
  try {
    const d: any = await api.get('/sgs-reports')
    list.value = d?.items || d || []
  } catch {}
}
onMounted(load)

async function save() {
  await api.post('/sgs-reports', form)
  ElMessage.success('Uploaded')
  open.value = false
  Object.assign(form, Object.fromEntries(Object.keys(form).map(k => [k, k === 'tea_type' ? 'raw_puer' : ''])))
  load()
}
</script>

<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-medium">SGS Certificates</span>
        <el-button type="primary" @click="open = true">+ Upload Report</el-button>
      </div>
    </template>
    <el-table :data="list" stripe>
      <el-table-column prop="report_no" label="Report No" width="160" />
      <el-table-column prop="batch_no" label="Batch" width="120" />
      <el-table-column prop="tea_type" label="Tea Type" width="120" />
      <el-table-column prop="test_date" label="Test Date" width="140" />
      <el-table-column prop="issue_date" label="Issued" width="140" />
      <el-table-column prop="pdf_url" label="PDF" min-width="200">
        <template #default="{ row }">
          <el-link v-if="row.pdf_url" type="primary" :href="row.pdf_url" target="_blank">View PDF →</el-link>
          <span v-else class="text-slate-300">—</span>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="open" title="Upload SGS Report" width="560px">
    <el-form :model="form" label-width="120px">
      <el-form-item label="Report No" required>
        <el-input v-model="form.report_no" placeholder="SGS-2024-001" />
      </el-form-item>
      <el-form-item label="Batch No" required>
        <el-input v-model="form.batch_no" />
      </el-form-item>
      <el-form-item label="Tea Type" required>
        <el-select v-model="form.tea_type" class="w-full">
          <el-option value="raw_puer" label="生普洱" />
          <el-option value="ripe_puer" label="熟普洱" />
          <el-option value="green" label="绿茶" />
          <el-option value="black" label="红茶" />
        </el-select>
      </el-form-item>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Test Date" required>
            <el-date-picker v-model="form.test_date" type="date" class="w-full" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Issue Date" required>
            <el-date-picker v-model="form.issue_date" type="date" class="w-full" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="PDF URL" required>
        <el-input v-model="form.pdf_url" placeholder="https://..." />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">Cancel</el-button>
      <el-button type="primary" @click="save">Upload</el-button>
    </template>
  </el-dialog>
</template>
