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

    <!-- Filter Bar -->
    <div class="flex gap-2 mb-3 flex-wrap">
      <el-select v-model="filter.visibility" placeholder="Visibility" clearable size="small" style="width:160px" @change="load">
        <el-option label="Public (all visitors)" value="public" />
        <el-option label="Registered users" value="registered" />
        <el-option label="Restricted (selected)" value="restricted" />
      </el-select>
      <el-select v-model="filter.type" placeholder="Access Type" clearable size="small" style="width:160px" @change="load">
        <el-option label="Slow Live" value="slow_live" />
        <el-option label="Scheduled" value="scheduled" />
        <el-option label="Advisor Room" value="advisor" />
        <el-option label="Admin Room" value="admin" />
      </el-select>
      <el-select v-model="filter.status" placeholder="Status" clearable size="small" style="width:140px" @change="load">
        <el-option label="Configuring" value="configuring" />
        <el-option label="Scheduled" value="scheduled" />
        <el-option label="Live" value="live" />
        <el-option label="Offline" value="offline" />
        <el-option label="Ended" value="ended" />
      </el-select>
      <el-input v-model="filter.keyword" placeholder="Search name/room_id" clearable size="small" style="width:220px" @keyup.enter="load" />
      <el-button size="small" @click="load">🔍 Search</el-button>
      <el-button size="small" text @click="resetFilter">Reset</el-button>
    </div>

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
      <!-- Visibility Column (NEW) -->
      <el-table-column label="Visibility" width="150">
        <template #default="{ row }">
          <el-tag
            :type="({public:'success', registered:'info', restricted:'danger'} as any)[row.visibility || 'registered']"
            effect="dark" size="small">
            {{ row.visibility || 'registered' }}
          </el-tag>
          <div v-if="row.visibility === 'restricted'" class="text-[10px] text-slate-400 mt-1">
            {{ (row.visible_user_ids || []).length }} users · {{ (row.visible_group_ids || []).length }} groups
          </div>
        </template>
      </el-table-column>
      <!-- Access Type Column (NEW) -->
      <el-table-column prop="type" label="Access" width="110">
        <template #default="{ row }">
          <span class="text-xs text-slate-500">{{ row.type || 'scheduled' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="Status" width="110">
        <template #default="{ row }">
          <el-tag :type="row.status === 'live' ? 'danger' : row.status === 'ended' ? 'success' : 'info'"
                   effect="dark" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="scheduled_start" label="Scheduled" width="170">
        <template #default="{ row }">{{ row.scheduled_start?.slice(0,16).replace('T',' ') || '—' }}</template>
      </el-table-column>
      <el-table-column prop="location" label="Location" show-overflow-tooltip />
      <el-table-column label="Rec." width="80">
        <template #default="{ row }">
          <el-tag v-if="row.enable_recording !== false" type="success" size="small" effect="plain">ON</el-tag>
          <el-tag v-else size="small" effect="plain">OFF</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Recording" width="100">
        <template #default="{ row }">
          <el-link v-if="row.recording_url" type="primary" :href="resolveUrl(row.recording_url)" target="_blank" size="small">View</el-link>
          <span v-else class="text-slate-300 text-xs">—</span>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="340">
        <template #default="{ row }">
          <el-button size="small" v-if="row.status !== 'live'" type="success" :loading="goLiveLoading" @click="goLive(row)">🎬 Go Live</el-button>
          <el-button size="small" v-else type="warning" @click="end(row)">End</el-button>
          <el-button size="small" v-if="row.status !== 'live'" @click="openEdit(row)">Edit</el-button>
          <el-button size="small" v-if="row.status !== 'live'" @click="openVisEdit(row)">👁 Vis</el-button>
          <el-button size="small" @click="copyKey(row)">Key</el-button>
          <el-button size="small" type="danger" @click="del(row.id)">Del</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- Create / Edit Dialog -->
  <el-dialog v-model="open" :title="isEditing ? 'Edit Live Room' : 'Create Live Room'" width="680px">
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

      <!-- NEW: Visibility -->
      <el-form-item label="Visibility">
        <el-radio-group v-model="form.visibility">
          <el-radio-button value="public">🌐 Public</el-radio-button>
          <el-radio-button value="registered">🔒 Registered</el-radio-button>
          <el-radio-button value="restricted">👥 Restricted</el-radio-button>
        </el-radio-group>
        <div class="text-[11px] text-slate-400 mt-1">
          public=任何人可看 · registered=登录即可 · restricted=必须在白名单用户或用户组里
        </div>
      </el-form-item>

      <!-- NEW: Restricted access controls -->
      <template v-if="form.visibility === 'restricted'">
        <el-form-item label="Visible User IDs">
          <el-select v-model="form.visible_user_ids" multiple filterable allow-create default-first-option style="width:100%"
                     placeholder="Type user ID and press Enter">
          </el-select>
          <div class="text-[11px] text-slate-400">直接输入数字 ID 回车添加</div>
        </el-form-item>
        <el-form-item label="Visible Groups">
          <el-select v-model="form.visible_group_ids" multiple style="width:100%" placeholder="Select user groups">
            <el-option v-for="g in groups" :key="g.id" :value="String(g.id)" :label="g.name" />
          </el-select>
        </el-form-item>
      </template>

      <!-- NEW: Access Type + Enable Recording -->
      <el-form-item label="Access Type">
        <el-select v-model="form.type" style="width:100%">
          <el-option label="Slow Live (24/7 camera)" value="slow_live" />
          <el-option label="Scheduled (appointment)" value="scheduled" />
          <el-option label="Advisor Private" value="advisor" />
          <el-option label="Admin Only" value="admin" />
        </el-select>
      </el-form-item>
      <el-form-item label="Enable Recording">
        <el-switch v-model="form.enable_recording" />
        <span class="ml-2 text-xs text-slate-400">结束后自动生成 Recording 回放（可单独配置可见性）</span>
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

  <!-- Quick Visibility Dialog (shortcut) -->
  <el-dialog v-model="openVis" :title="`Quick Visibility for #${editingVis?.id}`" width="520px">
    <el-form :model="visForm" label-width="140px">
      <el-form-item label="Visibility">
        <el-radio-group v-model="visForm.visibility">
          <el-radio-button value="public">🌐 Public</el-radio-button>
          <el-radio-button value="registered">🔒 Registered</el-radio-button>
          <el-radio-button value="restricted">👥 Restricted</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <template v-if="visForm.visibility === 'restricted'">
        <el-form-item label="Visible User IDs">
          <el-select v-model="visForm.visible_user_ids" multiple filterable allow-create default-first-option style="width:100%" />
        </el-form-item>
        <el-form-item label="Visible Groups">
          <el-select v-model="visForm.visible_group_ids" multiple style="width:100%">
            <el-option v-for="g in groups" :key="g.id" :value="String(g.id)" :label="g.name" />
          </el-select>
        </el-form-item>
      </template>
    </el-form>
    <template #footer>
      <el-button @click="openVis = false">Cancel</el-button>
      <el-button type="primary" :loading="saving" @click="saveVis">Save Visibility</el-button>
    </template>
  </el-dialog>

  <!-- ========== 开播配置弹窗 ========== -->
  <el-dialog v-model="openGoLive" :title="`🎬 Go Live · ${goLiveRoom?.room_name || ''}`" width="680px">
    <div v-if="goLiveConfig" class="space-y-4">
      <!-- 直播状态 -->
      <el-alert type="success" show-icon :closable="false">
        Room status: <strong>LIVE</strong> · room_id: <code>{{ goLiveRoom?.room_id }}</code>
      </el-alert>

      <!-- 方案选择 -->
      <el-tabs :model-value="goLiveRoom?.push_source || 'app_webrtc'" @update:model-value="(v) => { if(goLiveRoom) goLiveRoom.push_source = v }">
        <!-- === 手机 App WebRTC 推流 (推荐, 支持连麦+翻译) === -->
        <el-tab-pane label="📱 Mobile App WebRTC" name="app_webrtc">
          <div class="text-xs text-slate-500 mb-3">
            手机端打开 App / H5 主播页，填入 LiveKit URL + Host Token 即可开播。
            支持互动连麦、实时字幕翻译。
          </div>
          <div class="space-y-3">
            <!-- 📱 一键打开主播 H5 (手机扫码 / 桌面新窗口) -->
            <div class="flex items-center gap-2 p-3 bg-green-50 border border-green-200 rounded">
              <div class="flex-1">
                <div class="text-sm font-medium text-green-700">📱 Quick Start</div>
                <div class="text-xs text-green-600 mt-1">手机扫码或新窗口打开主播 H5（已自动填入所有参数）</div>
              </div>
              <el-button type="success" size="small" @click="openHostH5">📱 Open Host H5</el-button>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-32 text-xs text-slate-500">LiveKit URL</span>
              <el-input :model-value="goLiveConfig.livekit_url" readonly size="small" />
              <el-button size="small" @click="copy(goLiveConfig.livekit_url, 'LiveKit URL')">Copy</el-button>
            </div>
            <div class="flex items-start gap-2">
              <span class="w-32 text-xs text-slate-500 mt-1">Host Token</span>
              <el-input :model-value="goLiveConfig.host_token" readonly size="small" type="textarea" :autosize="{ minRows: 3, maxRows: 5 }" />
              <el-button size="small" @click="copy(goLiveConfig.host_token, 'Host Token')">Copy</el-button>
            </div>
            <div class="text-[11px] text-amber-600">
              ⚠ Token 有效期 {{ goLiveConfig.expires || 3600 }}s（约 1h）。过期后需要重新点击 "🎬 Go Live" 获取新 token。
            </div>
          </div>
        </el-tab-pane>

        <!-- === OBS RTMP 推流 === -->
        <el-tab-pane label="🖥 OBS RTMP" name="obs_rtmp">
          <div class="text-xs text-slate-500 mb-3">
            OBS Studio → 设置 → 直播 → 自定义：粘贴 URL + 密钥。推流到 MediaMTX，后端自动转 WebRTC。
          </div>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-32 text-xs text-slate-500">服务器</span>
              <el-input :model-value="goLiveConfig.obs_rtmp_url" readonly size="small" />
              <el-button size="small" @click="copy(goLiveConfig.obs_rtmp_url, 'RTMP URL')">Copy</el-button>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-32 text-xs text-slate-500">串流密钥</span>
              <el-input :model-value="goLiveConfig.obs_rtmp_key" readonly size="small" show-password />
              <el-button size="small" @click="copy(goLiveConfig.obs_rtmp_key, 'Stream Key')">Copy</el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- === Camera RTMP (慢直播) === -->
        <el-tab-pane label="📹 Camera RTMP" name="camera_rtmp">
          <div class="text-xs text-slate-500 mb-3">
            IP Camera → RTMP 推流到 MediaMTX slow 入口。24/7 慢直播专用。
          </div>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-32 text-xs text-slate-500">RTMP URL</span>
              <el-input :model-value="goLiveConfig.obs_rtmp_url" readonly size="small" />
              <el-button size="small" @click="copy(goLiveConfig.obs_rtmp_url, 'Camera RTMP')">Copy</el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    <template #footer>
      <el-button @click="openGoLive = false">Close</el-button>
      <el-button type="danger" @click="end(goLiveRoom)" :loading="goLiveLoading">End Stream</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api, upload } from '@/api/client'

const list = ref<any[]>([])
const groups = ref<any[]>([])
const open = ref(false)
const openVis = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const editingVis = ref<any>(null)
const uploadingCover = ref(false)
const uploadingRec = ref(false)
const saving = ref(false)
const coverRef = ref<any>(null)
const recRef = ref<any>(null)

const filter = reactive({ visibility: '', type: '', status: '', keyword: '' })

function triggerFile(ref: any) { ref?.click() }
const form = reactive({
  room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp',
  order_id: null as number | null, location: '', description: '',
  cover_image: '', recording_url: '',
  // 2026-09 新增字段
  visibility: 'registered',
  type: 'scheduled',
  visible_user_ids: [] as string[],
  visible_group_ids: [] as string[],
  enable_recording: true,
})
const visForm = reactive({
  visibility: 'registered',
  visible_user_ids: [] as string[],
  visible_group_ids: [] as string[],
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

function resetFilter() {
  Object.assign(filter, { visibility: '', type: '', status: '', keyword: '' })
  load()
}

async function load() {
  try {
    const params: any = {}
    if (filter.visibility) params.visibility = filter.visibility
    if (filter.type) params.type = filter.type
    if (filter.status) params.status = filter.status
    if (filter.keyword) params.keyword = filter.keyword
    const d: any = await api.get('/live-rooms', { params })
    list.value = d?.items || d || []
  } catch {}
  // 加载用户组列表（给 restricted 用）
  if (!groups.value.length) {
    try { groups.value = await api.get('/user-groups') } catch { groups.value = [] }
  }
}
onMounted(load)

function openCreate() {
  isEditing.value = false
  editingId.value = null
  Object.assign(form, {
    room_name: '', room_type: 'obs_tasting', push_source: 'obs_rtmp',
    order_id: null, location: '', description: '', cover_image: '', recording_url: '',
    visibility: 'registered', type: 'scheduled',
    visible_user_ids: [], visible_group_ids: [],
    enable_recording: true,
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
    recording_url: row.recording_url ?? '',
    visibility: row.visibility ?? 'registered',
    type: row.type ?? 'scheduled',
    visible_user_ids: (row.visible_user_ids || []).map(String),
    visible_group_ids: (row.visible_group_ids || []).map(String),
    enable_recording: row.enable_recording !== false,
  })
  open.value = true
}

function openVisEdit(row: any) {
  editingVis.value = row
  Object.assign(visForm, {
    visibility: row.visibility ?? 'registered',
    visible_user_ids: (row.visible_user_ids || []).map(String),
    visible_group_ids: (row.visible_group_ids || []).map(String),
  })
  openVis.value = true
}

async function saveVis() {
  saving.value = true
  try {
    const payload: any = {
      visibility: visForm.visibility,
      visible_user_ids: visForm.visible_user_ids,
      visible_group_ids: visForm.visible_group_ids,
    }
    await api.put(`/live-rooms/${editingVis.value.id}`, payload)
    ElMessage.success('Visibility updated')
    openVis.value = false
    load()
  } catch (err: any) { ElMessage.error(err?.message || 'Failed') }
  finally { saving.value = false }
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

// ========== 开播配置弹窗 ==========
const openGoLive = ref(false)
const goLiveRoom = ref<any>(null)
const goLiveConfig = ref<any>(null)  // { host_token, obs_rtmp_url, obs_rtmp_key, web_url }
const goLiveLoading = ref(false)

async function goLive(row: any) {
  goLiveRoom.value = row
  goLiveConfig.value = null
  goLiveLoading.value = true
  try {
    // 1) 先 Start（后端状态机 + 触发 LiveKit CreateRoom + 确保 host token）
    const started = await api.post(`/live-rooms/${row.id}/start`)
    // 2) 拿 OBS/Host token + RTMP/Web 配置
    const tokenResp = await api.post('/livekit/token-for-obs', {
      room_name: row.room_id,
      identity: 'host-' + row.room_id,
    })
    goLiveConfig.value = {
      host_token: tokenResp?.token || (started?.livekit_token_for_host) || '',
      expires: tokenResp?.expires,
      obs_rtmp_url: tokenResp?.obs_rtmp_url || row.obs_rtmp_url,
      obs_rtmp_key: tokenResp?.obs_rtmp_key || row.obs_rtmp_key,
      web_url: tokenResp?.web_url,
      livekit_url: import.meta.env.VITE_LIVEKIT_URL || '',
    }
    openGoLive.value = true
    await load()  // 刷新列表显示 live 状态
    ElMessage.success('Room started. Stream config ready.')
  } catch (err: any) {
    // 如果已经是 live（重复 start），直接拿配置
    if (err?.message?.includes('cannot start') || err?.code === 400) {
      try {
        const tokenResp = await api.post('/livekit/token-for-obs', {
          room_name: row.room_id,
          identity: 'host-' + row.room_id,
        })
        goLiveConfig.value = {
          host_token: tokenResp?.token || '',
          expires: tokenResp?.expires,
          obs_rtmp_url: tokenResp?.obs_rtmp_url || row.obs_rtmp_url,
          obs_rtmp_key: tokenResp?.obs_rtmp_key || row.obs_rtmp_key,
          web_url: tokenResp?.web_url,
          livekit_url: import.meta.env.VITE_LIVEKIT_URL || '',
        }
        openGoLive.value = true
      } catch (e2: any) {
        ElMessage.error(err?.message || e2?.message || 'Failed to start')
      }
    } else {
      ElMessage.error(err?.message || 'Failed to start')
    }
  } finally { goLiveLoading.value = false }
}

async function end(row: any) {
  await api.post(`/live-rooms/${row.id}/end`)
  ElMessage.success('Ended')
  openGoLive.value = false
  await load()
}

function copy(val: string, label = 'Value') {
  if (!val) { ElMessage.info(`${label} is empty`); return }
  navigator.clipboard.writeText(val).then(() => ElMessage.success(`${label} copied`))
}

function copyKey(row: any) {
  const k = row.obs_rtmp_key || row.camera_rtmp_key || ''
  if (k) copy(k, 'RTMP Key')
  else ElMessage.info('No RTMP key yet — start stream first')
}

function openHostH5() {
  if (!goLiveConfig.value || !goLiveRoom.value) return
  const adminBase = window.location.origin
  // admin-dashboard 和 public-site 可能不同 host；尝试相对路径 fallback
  const base = adminBase.includes('admin') ? adminBase.replace('admin', '') : adminBase
  const params = new URLSearchParams({
    host_token:  goLiveConfig.value.host_token,
    livekit_url: goLiveConfig.value.livekit_url,
    room_id:     goLiveRoom.value.room_id,
  })
  const url = `${base}/#/host?${params.toString()}`
  window.open(url, '_blank', 'noopener,noreferrer')
  copy(url, 'Host H5 URL (copied to clipboard)')
}
</script>
