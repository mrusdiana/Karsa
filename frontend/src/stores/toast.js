import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])
  let nextId = 0

  function show(message, tone = 'success') {
    const toast = { id: ++nextId, message, tone }
    toasts.value.push(toast)
    setTimeout(() => (toasts.value = toasts.value.filter((item) => item.id !== toast.id)), 3000)
  }

  return { toasts, show }
})
