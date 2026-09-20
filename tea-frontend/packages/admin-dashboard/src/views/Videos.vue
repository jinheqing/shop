<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api, upload } from '@/api/client'

// ===== Data =====
const videos = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)

// ===== Category dialog =====
const catDialog = ref(false)
const catEditing = ref<any>(null)
const catForm = reactive({ name: '', slug: '', description: '', sort_order: 0 })

// ===== Video dialog =====
const videoDialog = ref(false)
const videoEditing = ref<any>(null)
const videoForm = reactive({
  id: null as any,
  title: '',
  description: '',
  video_url: '',
  cover_url: '',
  duration_sec: 0,
  category_id: null as any,
  sort_order: 0,
  status: 'draft',
})
const videoUploading = ref(false)
const coverUploading = ref(false)

const filterCategory = ref('')
const filterStatus = ref('')

const filteredVideos = computed(() => {
  return videos.value.filter(v => {
    if (filterCategory.value && v.category_id !== Number(filterCategory.value)) return false
    if (filterStatus.value && v.status !== filterStatus.value) return false
    return true
  })
})

async function load() {
  loading.value = true
  try {
    videos.value = await api.get('/videos') || []
  } catch { videos.value = [] }
  try {
    categories.value = await api.get('/video-categories') || []
  } catch { categories.value = [] }
  loading.value = false
}

// ===== Category CRUD =====
function openCatCreate() {
  catEditing.value = null
  Object.assign(catForm, { name: '', slug: '', description: '', sort_order: 0 })
  catDialog.value = true
}
function openCatEdit(row: any) {
  catEditing.value = row
  Object.assign(catForm, { name: row.name, slug: row.slug, description: row.description || '', sort_order: row.sort_order || 0 })
  catDialog.value = true
}
async function saveCat() {
  if (!catForm.name) { ElMessage.warning('Name required'); return }
  try {
    if (catEditing.value) {
      await api.put(`/video-categories/${catEditing.value.id}`, { ...catForm })
    } else {
      await api.post('/video-categories', { ...catForm })
    }
    ElMessage.success('Saved')
    catDialog.value = false
    load()
  } catch (e: any) { ElMessage.error(e.message || 'Failed') }
}
async function delCat(row: any) {
  try { await ElMessageBox.confirm(`Delete category "${row.name}"? Videos will become uncategorized.`) } catch { return }
  await api.delete(`/video-categories/${row.id}`)
  ElMessage.success('Deleted')
  load()
}

// ===== Video CRUD =====
function openVideoCreate() {
  videoEditing.value = null
  Object.assign(videoForm, { id: null, title: '', description: '', video_url: '', cover_url: '', duration_sec: 0, category_id: null, sort_order: 0, status: 'draft' })
  videoDialog.value = true
}
function openVideoEdit(row: any) {
  videoEditing.value = row
  Object.assign(videoForm, {
    id: row.id,
    title: row.title || '',
    description: row.description || '',
    video_url: row.video_url || '',
    cover_url: row.cover_url || '',
    duration_sec: row.duration_sec || 0,
    category_id: row.category_id || null,
    sort_order: row.sort_order || 0,
    status: row.status || 'draft',
  })
  videoDialog.value = true
}
async function saveVideo() {
  if (!videoForm.title) { ElMessage.warning('Title required'); return }
  if (!videoForm.video_url) { ElMessage.warning('Video file required'); return }
  try {
    if (videoEditing.value) {
      await api.put(`/videos/${videoForm.id}`, { ...videoForm })
    } else {
      await api.post('/videos', { ...videoForm })
    }
    ElMessage.success('Saved')
    videoDialog.value = false
    load()
  } catch (e: any) { ElMessage.error(e.message || 'Failed') }
}
async function delVideo(row: any) {
  try { await ElMessageBox.confirm(`Delete video "${row.title}"?`) } catch { return }
  await api.delete(`/videos/${row.id}`)
  ElMessage.success('Deleted')
  load()
}
async function togglePublish(row: any) {
  const newStatus = row.status === 'published' ? 'draft' : 'published'
  await api.put(`/videos/${row.id}`, { status: newStatus })
  ElMessage.success(newStatus === 'published' ? 'Published' : 'Unpublished')
  load()
}

// ===== Upload =====
async function handleVideoUpload(file: File) {
  videoUploading.value = true
  try {
    const url = await upload(file, 'video')
    videoForm.video_url = url
    ElMessage.success('Video uploaded')
  } catch (e: any) { ElMessage.error('Upload failed: ' + (e.message || '')) }
  videoUploading.value = false
}
async function handleCoverUpload(file: File) {
  coverUploading.value = true
  try {
    const url = await upload(file, 'image')
    videoForm.cover_url = url
    ElMessage.success('Cover uploaded')
  } catch (e: any) { ElMessage.error('Upload failed: ' + (e.message || '')) }
  coverUploading.value = false
}

function resolveUrl(url: string) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return (import.meta.env.VITE_API_BASE || '').replace(/\/api\/v1$/, '') + url
}

function fmtDuration(sec: number) {
  if (!sec) return '—'
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

function onVideoFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files[0]) handleVideoUpload(input.files[0])
}
function onCoverFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files[0]) handleCoverUpload(input.files[0])
}

onMounted(load)
</script>

<template>
  <div>
    <!-- Categories section -->
    <el-card class="mb-4">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-medium">Video Categories</span>
          <el-button size="small" type="primary" @click="openCatCreate">+ Add Category</el-button>
        </div>
      </template>
      <el-table :data="categories" stripe size="small">
        <el-table-column prop="id" label="#" width="60" />
        <el-table-column prop="name" label="Name" min-width="150" />
        <el-table-column prop="slug" label="Slug" width="150" />
        <el-table-column prop="description" label="Description" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="Order" width="80" />
        <el-table-column label="Actions" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openCatEdit(row)">Edit</el-button>
            <el-button size="small" type="danger" @click="delCat(row)">Del</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Videos section -->
    <el-card>
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-medium">Videos</span>
          <div class="flex gap-2">
            <el-select v-model="filterCategory" placeholder="All Categories" clearable size="small" style="width:180px">
              <el-option v-for="c in categories" :key="c.id" :value="c.id" :label="c.name" />
            </el-select>
            <el-select v-model="filterStatus" placeholder="All Status" clearable size="small" style="width:120px">
              <el-option value="draft" label="Draft" />
              <el-option value="published" label="Published" />
            </el-select>
            <el-button size="small" @click="load">Refresh</el-button>
            <el-button size="small" type="primary" @click="openVideoCreate">+ Add Video</el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredVideos" stripe v-loading="loading">
        <el-table-column prop="id" label="#" width="60" />
        <el-table-column label="Cover" width="100">
          <template #default="{ row }">
            <img v-if="row.cover_url" :src="resolveUrl(row.cover_url)" class="w-16 h-10 object-cover rounded" />
            <span v-else class="text-slate-300 text-xs">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="Title" min-width="200" />
        <el-table-column label="Category" width="140">
          <template #default="{ row }">{{ row.category?.name || '—' }}</template>
        </el-table-column>
        <el-table-column label="Duration" width="80">
          <template #default="{ row }">{{ fmtDuration(row.duration_sec) }}</template>
        </el-table-column>
        <el-table-column label="Status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="Order" width="70" />
        <el-table-column label="Actions" width="220">
          <template #default="{ row }">
            <el-button size="small" @click="openVideoEdit(row)">Edit</el-button>
            <el-button size="small" :type="row.status === 'published' ? 'warning' : 'success'" @click="togglePublish(row)">
              {{ row.status === 'published' ? 'Unpublish' : 'Publish' }}
            </el-button>
            <el-button size="small" type="danger" @click="delVideo(row)">Del</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Category Dialog -->
    <el-dialog v-model="catDialog" :title="catEditing ? 'Edit Category' : 'New Category'" width="480px">
      <el-form :model="catForm" label-width="100px">
        <el-form-item label="Name">
          <el-input v-model="catForm.name" placeholder="e.g. Garden Stories" />
        </el-form-item>
        <el-form-item label="Slug">
          <el-input v-model="catForm.slug" placeholder="garden-stories" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="catForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="Sort Order">
          <el-input-number v-model="catForm.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="catDialog = false">Cancel</el-button>
        <el-button type="primary" @click="saveCat">Save</el-button>
      </template>
    </el-dialog>

    <!-- Video Dialog -->
    <el-dialog v-model="videoDialog" :title="videoEditing ? 'Edit Video' : 'New Video'" width="640px">
      <el-form :model="videoForm" label-width="100px">
        <el-form-item label="Title">
          <el-input v-model="videoForm.title" placeholder="Video title" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="videoForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="Category">
          <el-select v-model="videoForm.category_id" clearable placeholder="Uncategorized" style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :value="c.id" :label="c.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="Video File">
          <div class="flex items-center gap-2 w-full">
            <el-input v-model="videoForm.video_url" placeholder="/uploads/video/xxx.mp4" class="flex-1" />
            <label class="cursor-pointer">
              <el-button :loading="videoUploading" size="small" type="primary">Upload</el-button>
              <input type="file" accept="video/*" class="hidden" @change="onVideoFileChange" />
            </label>
          </div>
        </el-form-item>
        <el-form-item label="Cover Image">
          <div class="flex items-center gap-2 w-full">
            <el-input v-model="videoForm.cover_url" placeholder="/uploads/image/xxx.jpg" class="flex-1" />
            <label class="cursor-pointer">
              <el-button :loading="coverUploading" size="small">Upload</el-button>
              <input type="file" accept="image/*" class="hidden" @change="onCoverFileChange" />
            </label>
          </div>
          <img v-if="videoForm.cover_url" :src="resolveUrl(videoForm.cover_url)" class="mt-2 w-32 h-20 object-cover rounded" />
        </el-form-item>
        <el-form-item label="Duration (s)">
          <el-input-number v-model="videoForm.duration_sec" :min="0" />
        </el-form-item>
        <el-form-item label="Sort Order">
          <el-input-number v-model="videoForm.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="Status">
          <el-radio-group v-model="videoForm.status">
            <el-radio value="draft">Draft</el-radio>
            <el-radio value="published">Published</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="videoDialog = false">Cancel</el-button>
        <el-button type="primary" @click="saveVideo">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>
