import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'

export const useAuth = defineStore('auth', () => {
  const token = ref(localStorage.getItem('staff_token') || '')
  const user = ref<any>(null)

  async function login(email: string, password: string) {
    const data: any = await api.post('/staff/login', { email, password })
    token.value = data.access_token
    user.value = { email, role: data.role }
    localStorage.setItem('staff_token', token.value)
    return data
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('staff_token')
  }

  return { token, user, login, logout }
})
