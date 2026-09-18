import axios from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' }
})

// upload — 文件上传 helper，调用 POST /api/v1/upload?type=image|video|file
// 返回后端响应里的 file.url（形如 /uploads/image/xxx.jpg）
export async function upload(file: File, type: 'image' | 'video' | 'file' = 'image'): Promise<string> {
  const fd = new FormData()
  fd.append('file', file)
  const res = await axios.post(`${api.defaults.baseURL}/upload?type=${type}`, fd, {
    headers: {
      'Content-Type': 'multipart/form-data',
      ...(localStorage.getItem('staff_token') ? { Authorization: `Bearer ${localStorage.getItem('staff_token')}` } : {})
    },
    timeout: type === 'video' ? 120000 : 30000
  })
  return res.data?.file?.url || res.data?.url || ''
}

api.interceptors.request.use(cfg => {
  const t = localStorage.getItem('staff_token')
  if (t) cfg.headers.Authorization = `Bearer ${t}`
  return cfg
})
api.interceptors.response.use(
  r => r.data,
  err => { if (err.response?.status === 401) { localStorage.removeItem('staff_token'); router.push('/login') } return Promise.reject(err?.response?.data || err.message) }
)

import router from '@/router'
