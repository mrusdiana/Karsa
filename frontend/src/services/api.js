import axios from 'axios'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const api = axios.create({ baseURL: import.meta.env.VITE_API_URL, headers: { Accept: 'application/json' } })

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    const isDisabled = status === 403 && error.response?.data?.message === 'Akun dinonaktifkan.'
    if ((status === 401 || isDisabled) && localStorage.getItem('token')) {
      useAuthStore().clearSession()
      if (isDisabled) useToastStore().show('Akunmu dinonaktifkan admin.', 'danger')
      const current = router.currentRoute.value
      if (current.meta.role) router.push({ name: 'login', query: { redirect: current.fullPath } })
    }
    return Promise.reject(error)
  },
)

export function errorMessage(error) {
  if (!error.response) return 'Terjadi kesalahan. Coba lagi.'
  if (error.response.status === 404) return 'Data tidak ditemukan.'
  return error.response.data?.message ?? 'Terjadi kesalahan. Coba lagi.'
}

export function fieldErrors(error) {
  return error.response?.status === 422 ? (error.response.data?.errors ?? {}) : {}
}

export async function openFile(url) {
  const fileWindow = window.open('', '_blank')
  try {
    const response = await api.get(url, { responseType: 'blob' })
    fileWindow.location = URL.createObjectURL(response.data)
  } catch (error) {
    fileWindow.close()
    useToastStore().show(errorMessage(error), 'danger')
  }
}

export default api
