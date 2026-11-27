import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import api from '@/services/api'

export const HOME_BY_ROLE = { seeker: 'home', company: 'company-dashboard', admin: 'admin-dashboard' }

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const ready = ref(false)

  const role = computed(() => user.value?.role ?? null)
  const company = computed(() => user.value?.company ?? null)
  const isVerified = computed(() => company.value?.status === 'verified')

  async function init() {
    if (ready.value) return
    if (localStorage.getItem('token')) {
      try {
        user.value = (await api.get('/auth/me')).data.user
      } catch {
        clearSession()
      }
    }
    ready.value = true
  }

  async function refresh() {
    user.value = (await api.get('/auth/me')).data.user
  }

  function setSession(data) {
    localStorage.setItem('token', data.token)
    user.value = data.user
  }

  function clearSession() {
    localStorage.removeItem('token')
    user.value = null
  }

  const login = async (form) => setSession((await api.post('/auth/login', form)).data)
  const register = async (form) => setSession((await api.post('/auth/register', form)).data)

  return { user, ready, role, company, isVerified, init, refresh, login, register, clearSession }
})
