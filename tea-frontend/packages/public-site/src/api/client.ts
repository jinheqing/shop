import axios from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(cfg => {
  const t = localStorage.getItem('user_token')
  if (t) cfg.headers.Authorization = `Bearer ${t}`
  return cfg
})

api.interceptors.response.use(
  r => r.data,
  err => {
    if (err.response?.status === 401) localStorage.removeItem('user_token')
    return Promise.reject(err?.response?.data || err.message)
  }
)
