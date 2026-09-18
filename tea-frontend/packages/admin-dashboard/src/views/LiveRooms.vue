<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <span class="font-bold text-lg">Live Rooms</span>
        <el-button type="primary" @click="openCreate">+ New Room</el-button>
      </div>
    </template>
    <el-alert type="warning" show-icon :closable="false" class="mb-3">
      ⚠ Location field will be auto-sanitized to prefecture level on public API. Do NOT enter exact village / coordinates.
    </el-alert>
    <el-table :data="list" stripe>
      <el-table-column label="Cover" width="100">
        <template #default="{ row }">
          <el-image v-if="row.cover_image" :src="resolveUrl(row.cover_image)" fit="cover"
                   style="width:72px;height:44px;border-radius:2px;border:1px solid #e5e7eb" />
          <div v-else class="w-[72px] h-[44px] bg-slate-100 flex items-center justify-center text-slate-300 text-[10px]">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="room_name" label="Name" width="180" />
      <el-table-column prop="room_type" label="Type" width="130">
        <template #default="{ row }">
          <el-tag size="small">{{ row.room_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="Status" width="110">
        <template #default="{ row }">
          <el-tag :type="row.status === 'live' ? 'danger' : row.status === 'ended' ? 'success' : 'info'"
                   effect="dark" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="scheduled_start" label="Scheduled" width="170" />
      <el-table-column prop="location" label="Location" show-overflow-tooltip />
      <el-table-column label="Recording" width="100">
        <template #default="{ row }">
          <el-link v-if="row.recording_url" type="primary" :href="resolveUrl(row.recording_url)" target="_blank" size="small">View</el-link>
          <span v-else class="text-slate-300 text-xs">—</span>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="320">
        <template #default="{ row }">
          <el-button size="small" v-if="row.status !== 'live'" type="success" @click="start(row)">Start</el-button>
          <el-button size="small" v-else type="warning" @click="end(row)">End</el-button>
          <el-button size="small" @click="copyKey(row)">Key</el-button>
          <el-button size="small" @click="openEdit(row)">Edit</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- Create / Edit Dialog -->
  <el-dialog v-model="open" :title="isEditing ? 'Edit Live Room' : 'Create Live Room'" width="600px">
    <template v-if="!isEditing">
      <div class="mb-2 text-xs text-slate-500">Select a preset (matches backend RoomType constants)</div>
      <el-radio-group v-model="form.room_type" class="mb-4">
        <el-radio-button v-for="p in presets" :key="p.room_type" :value="p.room_type">{{ p.label }}</el-radio-button>
      </el-radio-group>
    </template>

    <el-form :model="form" label-width="140px">
      <el-form-item label="Room Name"><el-input v-model="form.room_name" placeholder="Gongfu · Morning Ceremony" /></el-form-item>
      <el-form-item label="Push Source">
        <el-select v-model="form.push_source" style="width:100%">
          <el-option label="OBS RTMP" value="obs_rtmp" />
          <el-option label="Camera RTMP (slow)" value="camera_rtmp" />
          <el-option label="App WebRTC" value="app_webrtc" />
        </el-select>
      </el-form-item>
      <el-form-item label="Order ID (optional)"><el-input-number v-model="form.order_id" :min="0" style="width:100%" /></el-form-item>
      <el-form-item label="Location">
        <el-input v-model="form.location" placeholder="云南省 · 临沧市" />
        <div class="text-[11px] text-amber-600 mt-1">⚠ 请勿填写具体村/乡/坐标。将自动脱敏到市/州级。</div>
      </el-form-item>
      <el-form-item label="Description">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="直播间描述（内部 + 前端展示）" />
      </el-form-item>

      <!-- Cover Image Upload -->
      <el-form-item label="Cover Image">
        <div class="flex items-start gap-4">
          <div v-if="form.cover_image" class="relative">
            <el-image :src="resolveUrl(form.cover_image)" fit="cover"
                     style="width:180px;height:112px;border-radius:2px;border:1px solid #e5e7eb" />
            <button type="button" @click="form.cover_image = ''"
                    class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full text-xs flex items-center justify-center hover:bg-red-600">×</button>
          </div>
          <div v-else class="w-[180px] h-[112px] bg-slate-100 border border-dashed border-slate-300 flex items-center justify-center text-slate-400 text-xs">
            No cover
          </div>
          <div class="flex flex-col gap-2">
            <el-button plain size="small" :disabled="uploadingCover" @click="triggerFile(coverRef)">
              {{ uploadingCover ? 'Uploading…' : (form.cover_image ? 'Replace Image' : 'Upload Cover') }}
            </el-button>
            <div class="text-[11px] text-slate-400">JPG/PNG · max 10MB · 16:10 recommended</div>
          </div>
          <input ref="coverRef" type="file" accept="image/jpeg,image/png,image/webp" class="hidden"
                 @change="onCoverChange" />
        </div>
      </el-form-item>

      <!-- Recording Upload -->
      <el-form-item label="Recording (optional)">
        <div class="flex items-center gap-4">
          <div v-if="form.recording_url" class="text-sm">
            <el-link type="primary" :href="resolveUrl(form.recording_url)" target="_blank">
              🎞 View Recording
            </el-link>
            <button type="button" @click="form.recording_url = ''" class="ml-3 text-red-500 text-xs">Remove</button>
          </div>
          <el-button v-else plain size="small" :disabled="uploadingRec" @click="triggerFile(recRef)">
            {{ uploadingRec ? 'Uploading…' : 'Upload Video' }}
          </el-button>
          <input ref="recRef" type="file" accept="video/mp4,video/webm,video/quicktime" class="hidden"
                 @change="onRecChange" />
        </div>
        <div class="text-[11px] text-slate-400 mt-1">MP4/WEBM/MOV · max 100MB · 回放链接将自动填入</div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="open = false">Cancel</el-button>
      <el-button type="primary" :loading="saving" @click="save">{{ isEditing ? 'Save Changes' : 'Create Room' }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api, upload } from '@/api/client'

const list = ref<any[]>([])
const open = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const uploadingCover = ref(false)
const uploadingRec = ref(false)
const saving = ref(false)
const coverRef = ref<any>(null)
const recRef = ref<any>(null)

function triggerFile(ref: any) { ref?.click() }
const form = reactive({
  room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp',
  order_id: null as number | null, location: '', description: '',
  cover_image: '', recording_url: ''
})

const presets = [
  { room_type: 'obs_tasting', label: 'OBS 品鉴' },
  { room_type: 'open_calendar', label: '开放预约' },
  { room_type: 'delivery_inspection', label: '验货间' },
  { room_type: 'customer_request', label: '客户请求' },
  { room_type: 'slow_preset', label: '24/7 慢直播' },
  { room_type: 'custom_private', label: '定制私密' },
]

function resolveUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const base = (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '')
  return base + url
}

async function load() {
  try {
    const d: any = await api.get('/live-rooms')
    list.value = d?.items || d || []
  } catch {}
}
onMounted(load)

function openCreate() {
  isEditing.value = false
  editingId.value = null
  Object.assign(form, {
    room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp',
    order_id: null, location: '', description: '', cover_image: '', recording_url: ''
  })
  open.value = true
}

function openEdit(row: any) {
  isEditing.value = true
  editingId.value = row.id
  Object.assign(form, {
    room_name: row.room_name, room_type: row.room_type, push_source: row.push_source,
    order_id: row.order_id ?? null, location: row.location ?? '',
    description: row.description ?? '', cover_image: row.cover_image ?? '',
    recording_url: row.recording_url ?? ''
  })
  open.value = true
}

async function onCoverChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]; if (!file) return
  uploadingCover.value = true
  try {
    form.cover_image = await upload(file, 'image')
    ElMessage.success('Cover uploaded')
  } catch (err: any) { ElMessage.error(err?.message || 'Upload failed') }
  finally { uploadingCover.value = false; (e.target as HTMLInputElement).value = '' }
}

async function onRecChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]; if (!file) return
  uploadingRec.value = true
  try {
    form.recording_url = await upload(file, 'video')
    ElMessage.success('Recording uploaded')
  } catch (err: any) { ElMessage.error(err?.message || 'Upload failed') }
  finally { uploadingRec.value = false; (e.target as HTMLInputElement).value = '' }
}

async function save() {
  saving.value = true
  try {
    if (isEditing.value && editingId.value) {
      await api.put(`/live-rooms/${editingId.value}`, { ...form })
      ElMessage.success('Updated')
    } else {
      await api.post('/live-rooms', { ...form })
      ElMessage.success('Created')
    }
    open.value = false
    load()
  } catch (err: any) { ElMessage.error(err?.message || 'Save failed') }
  finally { saving.value = false }
}

async function del(id: number) { await api.delete(`/live-rooms/${id}`); load() }
async function start(row: any) { await api.post(`/live-rooms/${row.id}/start`); ElMessage.success('Started'); load() }
async function end(row: any) { await api.post(`/live-rooms/${row.id}/end`); ElMessage.success('Ended'); load() }
function copyKey(row: any) {
  const k = row.obs_rtmp_key || row.camera_rtmp_key || ''
  if (k) { navigator.clipboard.writeText(k); ElMessage.success('Key copied') }
  else ElMessage.info('No RTMP key yet — start stream first')
}
</script>
