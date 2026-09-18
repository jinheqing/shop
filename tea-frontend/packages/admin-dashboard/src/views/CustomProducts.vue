<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api, upload } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const editingId = ref<number | null>(null)
const tab = ref<'list' | 'create'>('list')
const uploading = ref(false)
const coverFileRef = ref<any>(null)

function triggerFile(r: any) { r?.click() }
function resolveUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  return base + url
}
async function onImageFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const url = await upload(file, 'image')
    ;(form as any).product_image_url = url
    ElMessage.success('Product image uploaded')
  } catch (err: any) {
    ElMessage.error(`Upload failed: ${err?.message || err}`)
  } finally {
    uploading.value = false
    input.value = ''
  }
}

// 完整 26 字段表单 — 和 Go struct / DB DDL 100% 对齐
const form = reactive({
  // 基本 (4)
  title: '', raw_tea_source: '', custom_requirement: '',
  // 普洱参数 (6)
  tea_type: 'raw_puer', tea_shape: 'cake', tea_shape_weight: 357,
  smoked_with_flower: false, flower_type: 'jasmine',
  // 包装 (2)
  inner_packaging: '竹编内层 + 棉纸', outer_packaging: '哑光纸质礼盒 + 丝绒内衬', product_image_url: '',
  // 产品卡 (3)
  product_card_text: '', product_card_format: 'vertical', qr_code_position: 'outer_back',
  // 价格 (4)
  unit_price: 0, quantity: 1, shipping_cost: 0, lead_time: '45 days from confirmation',
  // 溯源 (5)
  harvest_date: '', roasting_date: '', tea_garden_location: '', master_name: '', storage_location: '',
  // SGS (1)
  sgs_report_id: null as number | null,
  // 直播 (2)
  include_custom_live: false, live_scheduled_date: '',
})

async function load() {
  try { const d: any = await api.get('/custom-products'); list.value = d?.items || d || [] }
  catch {
    list.value = [
      { id: 16, title: '潮州凤凰单丛古树茶', tea_type: 'raw_puer', tea_shape: 'cake', tea_shape_weight: 357, tea_garden_location: '广东省潮州市潮安区凤凰镇大乌岽村茶园', master_name: '李师傅', status: 'published', version: 1, product_token: 'XK92AB38...', unit_price: 188, quantity: 2, total_amount: 376, harvest_date: '2026-04-15', roasting_date: '2026-05-20', raw_tea_source: '广东省潮州市凤凰镇大乌岽村，树龄 800 年', custom_requirement: '想要古树生普，压制 357g 饼，茉莉花香熏制' },
      { id: 18, title: '临沧邦东古树熟茶', tea_type: 'ripe_puer', tea_shape: 'brick', tea_shape_weight: 500, tea_garden_location: '云南省临沧市临翔区邦东乡曼岗村茶园', master_name: '王师傅', status: 'draft', version: 1, product_token: null, unit_price: 120, quantity: 3, total_amount: 360, harvest_date: '2026-03-28', roasting_date: '2026-06-10', raw_tea_source: '云南省临沧市临翔区邦东乡曼岗村，茶园树龄 300 年', custom_requirement: '熟普砖茶，便于存放' },
    ]
  }
}
onMounted(load)

async function submit() {
  const method = editingId.value ? 'put' : 'post'
  const url = editingId.value ? `/custom-products/${editingId.value}` : '/custom-products'
  if (method === 'post') await api.post(url, form)
  else await api.put(url, form)
  ElMessage.success('✅ Saved')
  open.value = false; editingId.value = null; load()
}
function openNew() {
  Object.assign(form, { title: '', raw_tea_source: '', custom_requirement: '', tea_type: 'raw_puer', tea_shape: 'cake', tea_shape_weight: 357, smoked_with_flower: false, flower_type: 'jasmine', inner_packaging: '', outer_packaging: '', product_image_url: '', product_card_text: '', product_card_format: 'vertical', qr_code_position: 'outer_back', unit_price: 0, quantity: 1, shipping_cost: 0, lead_time: '', harvest_date: '', roasting_date: '', tea_garden_location: '', master_name: '', storage_location: '', sgs_report_id: null, include_custom_live: false, live_scheduled_date: '' })
  editingId.value = null; open.value = true
}
function openEdit(row: any) { editingId.value = row.id; Object.assign(form, row); open.value = true }
async function publish(id: number) { await api.post(`/custom-products/${id}/publish`); ElMessage.success('🚀 Published → product_token generated'); load() }
async function del(id: number) { await api.delete(`/custom-products/${id}`); load() }

function total() {
  return (form.unit_price * form.quantity + form.shipping_cost).toFixed(2)
}
</script>

<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">📦 Bespoke Products (26 fields, DB-aligned)</span>
      <el-button type="primary" @click="openNew">+ New Bespoke Quote</el-button>
    </div></template>

    <el-alert type="info" :closable="false" class="mb-4">
      All {{ 26 }} form field names match Go API JSON (snake_case) and DB columns exactly. Auto-gen fields: version, sku, product_token, created_by_staff_id, reviewed_by_staff_id.
    </el-alert>

    <el-table :data="list" stripe>
      <el-table-column prop="id" label="#" width="60" />
      <el-table-column label="Cover" width="100">
        <template #default="{ row }">
          <el-image v-if="row.product_image_url" :src="resolveUrl(row.product_image_url)" fit="cover"
                   style="width:72px;height:44px;border-radius:2px;border:1px solid #e5e7eb" />
          <div v-else class="w-[72px] h-[44px] bg-slate-100 flex items-center justify-center text-slate-300 text-[10px]">&mdash;</div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="Title" min-width="260" />
      <el-table-column prop="tea_type" label="Type" width="130" />
      <el-table-column prop="tea_shape" label="Shape" width="100" />
      <el-table-column prop="tea_garden_location" label="Tea Garden" width="220" />
      <el-table-column prop="master_name" label="Master" width="120" />
      <el-table-column prop="status" label="Status" width="110">
        <template #default="{ row }">
          <el-tag :type="({'draft':'info','published':'success','archived':''} as Record<string,string>)[row.status]" effect="dark">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Version" width="80">
        <template #default="{ row }">v{{ row.version || 1 }}</template>
      </el-table-column>
      <el-table-column label="Price" width="100">
        <template #default="{ row }">£{{ row.total_amount || row.unit_price }}</template>
      </el-table-column>
      <el-table-column label="Actions" width="200">
        <template #default="{ row }">
          <el-button v-if="!row.product_token" size="small" type="success" @click="publish(row.id)">🚀 Publish</el-button>
          <el-button size="small" @click="openEdit(row)">Edit</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 完整 26 字段表单对话框 -->
  <el-dialog v-model="open" :title="editingId ? `Edit Bespoke #${editingId}` : 'Create New Bespoke Quote — 26 fields'" width="720px" top="5vh">
    <el-form :model="form" label-width="160px">
      <!-- Group 1: 基本信息 -->
      <el-divider content-position="left">📝 Basic Info (4)</el-divider>
      <el-form-item label="Title" required><el-input v-model="form.title" placeholder="潮州凤凰单丛古树茶" /></el-form-item>
      <el-form-item label="Raw Tea Source" required><el-input v-model="form.raw_tea_source" type="textarea" :rows="2" /></el-form-item>
      <el-form-item label="Customer Requirement" required><el-input v-model="form.custom_requirement" type="textarea" :rows="2" /></el-form-item>

      <!-- Group 2: 普洱参数 -->
      <el-divider content-position="left">🍵 Tea Parameters (6)</el-divider>
      <el-row :gutter="12">
        <el-col :span="12"><el-form-item label="Tea Type" required>
          <el-select v-model="form.tea_type" class="w-full">
            <el-option value="raw_puer">Raw Pu'er (生普)</el-option>
            <el-option value="ripe_puer">Ripe Pu'er (熟普)</el-option>
            <el-option value="ancient_tree">Ancient Tree (古树)</el-option>
            <el-option value="vintage">Vintage ( vintage )</el-option>
          </el-select>
        </el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Tea Shape" required>
          <el-select v-model="form.tea_shape" class="w-full">
            <el-option value="loose">Loose Leaf</el-option>
            <el-option value="cake">Cake (饼)</el-option>
            <el-option value="brick">Brick (砖)</el-option>
            <el-option value="tuo">Tuo (沱)</el-option>
          </el-select>
        </el-form-item></el-col>
      </el-row>
      <el-row :gutter="12">
        <el-col :span="8"><el-form-item label="Shape Weight (g)">
          <el-input-number v-model="form.tea_shape_weight" :min="50" :max="5000" class="w-full" />
        </el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Flower Smoked">
          <el-switch v-model="form.smoked_with_flower" />
        </el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Flower Type">
          <el-select v-model="form.flower_type" class="w-full" :disabled="!form.smoked_with_flower">
            <el-option value="jasmine">Jasmine 茉莉</el-option>
            <el-option value="osmanthus">Osmanthus 桂花</el-option>
            <el-option value="orchid">Orchid 兰花</el-option>
            <el-option value="custom">Custom</el-option>
          </el-select>
        </el-form-item></el-col>
      </el-row>

      <!-- Group 3: 包装 -->
      <el-divider content-position="left">📦 Packaging (2)</el-divider>
      <el-form-item label="Inner Packaging" required><el-input v-model="form.inner_packaging" /></el-form-item>
      <el-form-item label="Outer Packaging" required><el-input v-model="form.outer_packaging" /></el-form-item>

      <!-- Product Image Upload -->
      <el-form-item label="Product Image">
        <div class="flex items-start gap-4">
          <div v-if="(form as any).product_image_url" class="relative">
            <el-image :src="resolveUrl((form as any).product_image_url)" fit="cover"
                     style="width:160px;height:100px;border-radius:2px;border:1px solid #e5e7eb" />
            <button type="button" @click="(form as any).product_image_url = ''"
                    class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full text-xs flex items-center justify-center hover:bg-red-600">&times;</button>
          </div>
          <div class="flex flex-col gap-2">
            <el-button type="primary" plain size="small" :disabled="uploading" @click="triggerFile(coverFileRef)">
              <span v-if="uploading">Uploading&hellip;</span>
              <span v-else>{{ (form as any).product_image_url ? 'Replace Image' : 'Upload Product Image' }}</span>
            </el-button>
            <div class="text-[11px] text-slate-400">JPG/PNG/WEBP &middot; max 10MB &middot; 推荐 16:10 &middot; 800&times;500</div>
          </div>
          <input ref="coverFileRef" type="file" accept="image/jpeg,image/png,image/webp" class="hidden"
                 @change="onImageFileChange" />
        </div>
      </el-form-item>

      <!-- Group 4: 产品卡 -->
      <el-divider content-position="left">💳 Product Card (3)</el-divider>
      <el-form-item label="Card Text"><el-input v-model="form.product_card_text" type="textarea" /></el-form-item>
      <el-row :gutter="12">
        <el-col :span="12"><el-form-item label="Card Format">
          <el-select v-model="form.product_card_format" class="w-full">
            <el-option value="vertical">Vertical</el-option>
            <el-option value="horizontal">Horizontal</el-option>
          </el-select>
        </el-form-item></el-col>
        <el-col :span="12"><el-form-item label="QR Code Position" required>
          <el-select v-model="form.qr_code_position" class="w-full">
            <el-option value="outer_front">Outer Front</el-option>
            <el-option value="outer_back">Outer Back</el-option>
            <el-option value="inner_front">Inner Front</el-option>
            <el-option value="hidden">Hidden</el-option>
          </el-select>
        </el-form-item></el-col>
      </el-row>

      <!-- Group 5: 价格 -->
      <el-divider content-position="left">💷 Pricing (4)</el-divider>
      <el-row :gutter="12">
        <el-col :span="6"><el-form-item label="Unit Price (£)" required><el-input-number v-model="form.unit_price" :precision="2" class="w-full" /></el-form-item></el-col>
        <el-col :span="6"><el-form-item label="Quantity" required><el-input-number v-model="form.quantity" :min="1" class="w-full" /></el-form-item></el-col>
        <el-col :span="6"><el-form-item label="Shipping (£)"><el-input-number v-model="form.shipping_cost" :precision="2" class="w-full" /></el-form-item></el-col>
        <el-col :span="6"><el-form-item label="Total (£)"><span class="font-serif text-lg text-tea-700">£{{ total() }}</span></el-form-item></el-col>
      </el-row>
      <el-form-item label="Lead Time" required><el-input v-model="form.lead_time" placeholder="45 days from confirmation" /></el-form-item>

      <!-- Group 6: 溯源 -->
      <el-divider content-position="left">🏔️ Traceability (5)</el-divider>
      <el-row :gutter="12">
        <el-col :span="12"><el-form-item label="Tea Garden Location" required><el-input v-model="form.tea_garden_location" placeholder="云南省临沧市临翔区邦东乡曼岗村茶园" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="Master Name" required><el-input v-model="form.master_name" /></el-form-item></el-col>
      </el-row>
      <el-row :gutter="12">
        <el-col :span="8"><el-form-item label="Harvest Date" required>
          <el-date-picker v-model="form.harvest_date" type="date" value-format="YYYY-MM-DD" class="w-full" />
        </el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Roasting Date" required>
          <el-date-picker v-model="form.roasting_date" type="date" value-format="YYYY-MM-DD" class="w-full" />
        </el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Storage Location" required><el-input v-model="form.storage_location" /></el-form-item></el-col>
      </el-row>

      <!-- Group 7: SGS + Live -->
      <el-divider content-position="left">🔬 SGS + 🎥 Custom Live (3)</el-divider>
      <el-row :gutter="12">
        <el-col :span="8"><el-form-item label="SGS Report ID"><el-input-number v-model="form.sgs_report_id" :min="0" class="w-full" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Include Custom Live"><el-switch v-model="form.include_custom_live" /></el-form-item></el-col>
        <el-col :span="8"><el-form-item label="Live Scheduled">
          <el-date-picker v-model="form.live_scheduled_date" type="date" value-format="YYYY-MM-DD" class="w-full" :disabled="!form.include_custom_live" />
        </el-form-item></el-col>
      </el-row>
    </el-form>

    <template #footer>
      <el-button @click="open=false">Cancel</el-button>
      <el-button type="primary" @click="submit">💾 Save Bespoke Quote</el-button>
    </template>
  </el-dialog>
</template>
