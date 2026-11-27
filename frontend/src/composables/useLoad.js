import { onMounted, ref } from 'vue'

export function useLoad(fetcher, initial = null, { immediate = true } = {}) {
  const data = ref(initial)
  const loading = ref(immediate)
  const error = ref(false)

  async function load() {
    loading.value = true
    error.value = false
    try {
      data.value = await fetcher()
    } catch {
      error.value = true
    } finally {
      loading.value = false
    }
  }

  if (immediate) onMounted(load)
  return { data, loading, error, load }
}
