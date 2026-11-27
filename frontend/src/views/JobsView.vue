<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MapPin, Search, SlidersHorizontal, X } from 'lucide-vue-next'
import JobCard from '@/components/JobCard.vue'
import JobFilters from '@/components/JobFilters.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { JOB_LEVEL, JOB_TYPE, WORK_MODE } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const listKeys = ['level', 'mode', 'type']
const sortOptions = { newest: 'Terbaru', salary: 'Gaji tertinggi', deadline: 'Deadline terdekat' }
const drawerOpen = ref(false)
const keyword = ref('')
const location = ref('')

const filters = computed(() => ({
  q: route.query.q ?? '',
  category: route.query.category ?? '',
  location: route.query.location ?? '',
  salary_min: route.query.salary_min ?? '',
  sort: route.query.sort ?? 'newest',
  ...Object.fromEntries(listKeys.map((key) => [key, (route.query[key] ?? '').split(',').filter(Boolean)])),
}))

const categories = useLoad(async () => (await api.get('/categories')).data, [])
const jobs = useLoad(async () => (await api.get('/jobs', { params: route.query })).data, null, { immediate: false })

const activeChips = computed(() => {
  const chips = []
  const category = categories.data.value.find((item) => item.slug === filters.value.category)
  if (category) chips.push({ key: 'category', label: category.name })
  const labels = { level: JOB_LEVEL, mode: WORK_MODE, type: JOB_TYPE }
  for (const key of listKeys) for (const value of filters.value[key]) chips.push({ key, value, label: labels[key][value] })
  if (filters.value.location) chips.push({ key: 'location', label: filters.value.location })
  if (filters.value.salary_min) chips.push({ key: 'salary_min', label: `Gaji ≥ Rp ${Number(filters.value.salary_min) / 1_000_000} jt` })
  return chips
})

function update(changes) {
  const query = { ...route.query, ...changes }
  delete query.page
  for (const key of Object.keys(query)) if (!query[key] || query[key] === 'newest') delete query[key]
  router.push({ name: 'jobs', query })
}

const setFilter = (key, value) => update({ [key]: value })
function toggleFilter(key, value) {
  const values = filters.value[key].includes(value) ? filters.value[key].filter((item) => item !== value) : [...filters.value[key], value]
  update({ [key]: values.join(',') })
}
function removeChip(chip) {
  if (chip.value) toggleFilter(chip.key, chip.value)
  else setFilter(chip.key, '')
}
const clearFilters = () => router.push({ name: 'jobs', query: filters.value.q ? { q: filters.value.q } : {} })
const changePage = (page) => router.push({ name: 'jobs', query: { ...route.query, page } })

watch(
  () => route.query,
  () => {
    keyword.value = filters.value.q
    location.value = filters.value.location
    jobs.load()
  },
  { immediate: true },
)
</script>

<template>
  <div class="border-b border-line bg-surface py-5 sm:py-6">
    <form class="container-app flex flex-col gap-2.5 sm:flex-row" role="search" @submit.prevent="update({ q: keyword.trim(), location: location.trim() })">
      <label class="input flex flex-1 items-center gap-2.5"><Search class="size-[18px] text-muted" aria-hidden="true" /><span class="sr-only">Kata kunci</span><input v-model="keyword" class="w-full bg-transparent outline-none focus-visible:shadow-none" placeholder="Cari posisi, skill, atau perusahaan" /></label>
      <label class="input flex items-center gap-2.5 sm:w-60"><MapPin class="size-[18px] text-muted" aria-hidden="true" /><span class="sr-only">Lokasi</span><input v-model="location" class="w-full bg-transparent outline-none focus-visible:shadow-none" placeholder="Semua lokasi" /></label>
      <button type="submit" class="btn-primary px-7">Cari</button>
    </form>
  </div>

  <div class="container-app flex items-start gap-8 py-8 pb-16">
    <aside class="sticky top-24 hidden w-[280px] shrink-0 card p-5 lg:block" aria-label="Filter lowongan">
      <JobFilters :filters="filters" :categories="categories.data.value" @toggle="toggleFilter" @set="setFilter" @clear="clearFilters" />
    </aside>

    <section class="min-w-0 flex-1" aria-live="polite">
      <div class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 class="text-[22px] sm:text-2xl">{{ jobs.data.value ? `${jobs.data.value.total} lowongan ditemukan` : 'Cari lowongan' }}</h1>
          <p v-if="filters.q" class="mt-0.5 text-[14px] text-muted">untuk "{{ filters.q }}"</p>
        </div>
        <div class="flex gap-2">
          <button type="button" class="btn-secondary lg:hidden" @click="drawerOpen = true"><SlidersHorizontal class="size-4" aria-hidden="true" />Filter ({{ activeChips.length }})</button>
          <label class="sr-only" for="sort">Urutkan</label>
          <select id="sort" class="input w-auto" :value="filters.sort" @change="setFilter('sort', $event.target.value)">
            <option v-for="(label, value) in sortOptions" :key="value" :value="value">{{ label }}</option>
          </select>
        </div>
      </div>
      <div v-if="activeChips.length" class="mt-4 flex flex-wrap gap-2">
        <button v-for="chip in activeChips" :key="chip.key + (chip.value ?? '')" type="button" class="chip bg-primary-soft text-primary hover:bg-primary/15" :aria-label="`Hapus filter ${chip.label}`" @click="removeChip(chip)">{{ chip.label }}<X class="size-3.5" aria-hidden="true" /></button>
      </div>

      <div class="mt-5">
        <StatusState v-if="jobs.error.value" type="error" @retry="jobs.load" />
        <LoadingCards v-else-if="jobs.loading.value" :count="4" />
        <StatusState v-else-if="!jobs.data.value.data.length" title="Tidak ada lowongan yang cocok" text="Coba kata kunci lain atau kurangi filter.">
          <button type="button" class="btn-secondary mt-2" @click="router.push({ name: 'jobs' })">Hapus filter</button>
        </StatusState>
        <template v-else>
          <div class="flex flex-col gap-3.5"><JobCard v-for="job in jobs.data.value.data" :key="job.id" :job="job" /></div>
          <PaginationBar :page="jobs.data.value" @change="changePage" />
        </template>
      </div>
    </section>
  </div>

  <div v-if="drawerOpen" class="fixed inset-0 z-50 flex items-end bg-text/40 lg:hidden" @click.self="drawerOpen = false" @keydown.esc="drawerOpen = false">
    <div class="max-h-[85dvh] w-full overflow-y-auto rounded-t-[20px] bg-surface p-5" role="dialog" aria-label="Filter lowongan">
      <JobFilters :filters="filters" :categories="categories.data.value" @toggle="toggleFilter" @set="setFilter" @clear="clearFilters" />
      <button type="button" class="btn-primary mt-6 w-full" @click="drawerOpen = false">Tampilkan {{ jobs.data.value?.total ?? 0 }} lowongan</button>
    </div>
  </div>
</template>
