import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

export const useNotificationStore = defineStore('notifications', () => {
  const applications = ref(0)
  const invitations = ref(0)

  async function load() {
    if (useAuthStore().role !== 'seeker') {
      applications.value = 0
      invitations.value = 0
      return
    }
    try {
      const counts = (await api.get('/me/notifications')).data
      applications.value = counts.applications
      invitations.value = counts.invitations
    } catch {
      applications.value = 0
      invitations.value = 0
    }
  }

  return { applications, invitations, load }
})
