import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

export function useSignOut() {
  const confirming = ref(false)
  const router = useRouter()

  function signOut() {
    confirming.value = false
    useAuthStore().clearSession()
    useToastStore().show('Kamu sudah keluar.')
    router.push({ name: 'home' })
  }

  return { confirming, signOut }
}
