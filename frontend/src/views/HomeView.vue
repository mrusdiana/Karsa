<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowRight, Briefcase, CalendarClock, CheckCircle2, Cloud, Code2, LineChart, Mail, MapPin, Megaphone, PenTool, Search, TrendingUp } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import JobCard from '@/components/JobCard.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import MatchScore from '@/components/MatchScore.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const keyword = ref('')
const location = ref('')
const { data, loading, error, load } = useLoad(async () => (await api.get('/home')).data)

const categoryIcons = { 'software-engineering': Code2, data: LineChart, 'ui-ux-design': PenTool, product: Briefcase, qa: CheckCircle2, 'devops-cloud': Cloud, 'digital-marketing': Megaphone }
const popular = [
  { label: 'Frontend', query: { q: 'frontend' } },
  { label: 'Data Analyst', query: { q: 'data analyst' } },
  { label: 'UI/UX', query: { category: 'ui-ux-design' } },
  { label: 'Magang Remote', query: { level: 'internship', mode: 'remote' } },
]
const totalJobs = computed(() => data.value?.categories.reduce((sum, category) => sum + category.job_count, 0) ?? 0)

function search() {
  const query = {}
  if (keyword.value.trim()) query.q = keyword.value.trim()
  if (location.value.trim()) query.location = location.value.trim()
  router.push({ name: 'jobs', query })
}
</script>

<template>
  <section class="hero-grid border-b border-line">
    <div class="container-app relative flex items-center gap-12 py-12 sm:py-[72px]">
      <div class="min-w-0 flex-1">
        <span class="inline-flex h-7 items-center gap-1.5 rounded-full bg-primary/10 px-3 text-[13.5px] font-semibold text-primary"><TrendingUp class="size-3.5" aria-hidden="true" />Lowongan magang, entry, dan junior</span>
        <h1 class="mt-5 max-w-2xl text-[28px] leading-tight font-extrabold tracking-[-0.02em] sm:text-[40px]">Temukan magang dan kerja pertamamu di dunia tech.</h1>
        <p class="mt-4 max-w-xl text-base text-muted sm:text-[18px]">Lowongan dari perusahaan terverifikasi. Status lamaranmu selalu jelas, tanpa di-ghosting.</p>
        <form class="mt-8 flex max-w-2xl flex-col gap-2 rounded-card border border-line bg-surface p-2 shadow-card sm:flex-row" role="search" @submit.prevent="search">
          <label class="flex min-h-11 flex-1 items-center gap-2.5 px-3 text-muted">
            <Search class="size-[18px] shrink-0" aria-hidden="true" /><span class="sr-only">Kata kunci</span>
            <input v-model="keyword" class="w-full bg-transparent text-[16px] text-text outline-none placeholder:text-muted focus-visible:shadow-none" placeholder="Cari posisi, skill, atau perusahaan" />
          </label>
          <span class="hidden w-px self-stretch bg-line sm:my-2 sm:block" />
          <label class="flex min-h-11 items-center gap-2.5 px-3 text-muted sm:w-52">
            <MapPin class="size-[18px] shrink-0" aria-hidden="true" /><span class="sr-only">Lokasi</span>
            <input v-model="location" class="w-full bg-transparent text-[16px] text-text outline-none placeholder:text-muted focus-visible:shadow-none" placeholder="Semua lokasi" />
          </label>
          <button type="submit" class="btn-primary px-7">Cari</button>
        </form>
        <div class="mt-4 flex flex-wrap items-center gap-2">
          <span class="text-[14px] font-medium text-muted">Populer:</span>
          <RouterLink v-for="item in popular" :key="item.label" :to="{ name: 'jobs', query: item.query }" class="chip hover:bg-primary-soft hover:text-primary">{{ item.label }}</RouterLink>
        </div>
        <dl v-if="data" class="mt-7 flex gap-8">
          <div><dt class="sr-only">Lowongan dibuka</dt><dd class="num font-display text-[22px] font-extrabold">{{ totalJobs }}</dd><dd class="text-[14px] text-muted">Lowongan dibuka</dd></div>
          <div><dt class="sr-only">Perusahaan</dt><dd class="num font-display text-[22px] font-extrabold">{{ data.companies.length }}</dd><dd class="text-[14px] text-muted">Perusahaan merekrut</dd></div>
          <div><dt class="sr-only">Kategori</dt><dd class="num font-display text-[22px] font-extrabold">{{ data.categories.length }}</dd><dd class="text-[14px] text-muted">Kategori</dd></div>
        </dl>
      </div>
      <div class="hidden w-[380px] flex-col gap-3 lg:flex" aria-hidden="true">
        <div class="card p-[18px] shadow-card">
          <div class="flex items-center justify-between text-[14px] text-muted">Lamaran: Frontend Developer Intern<span class="inline-flex h-[26px] items-center gap-1.5 rounded-full bg-warning/10 px-2.5 text-[13.5px] font-semibold text-warning"><CalendarClock class="size-3.5" />Interview</span></div>
          <ul class="mt-3.5 flex flex-col gap-2.5 text-[14px]">
            <li class="flex items-center gap-2"><span class="size-2.5 rounded-full bg-primary" />Dilamar<span class="ml-auto text-muted">Senin</span></li>
            <li class="flex items-center gap-2"><span class="size-2.5 rounded-full bg-info" />Ditinjau perusahaan<span class="ml-auto text-muted">Selasa</span></li>
            <li class="flex items-center gap-2 font-semibold"><span class="size-2.5 rounded-full bg-warning" />Lanjut ke interview<span class="ml-auto font-medium text-muted">Rabu</span></li>
          </ul>
        </div>
        <div class="card flex items-center gap-3.5 p-[18px] shadow-card"><MatchScore :score="86" /><div><div class="text-sm font-semibold">Cocok tinggi untukmu</div><div class="text-[14px] text-muted">3 dari 3 skill cocok · level sesuai</div></div></div>
        <div class="card flex items-center gap-3.5 p-[18px] shadow-card"><span class="grid size-9 place-items-center rounded-full bg-info/10 text-info"><Mail class="size-[18px]" /></span><div><div class="text-sm font-semibold">Perusahaan mengundangmu melamar</div><div class="text-[14px] text-muted">Aktifkan "Bisa ditemukan" di profil</div></div></div>
      </div>
    </div>
  </section>

  <div class="container-app py-14">
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="4" />
    <template v-else>
      <section aria-labelledby="categories-title">
        <h2 id="categories-title" class="mb-5 text-xl">Jelajahi kategori</h2>
        <div class="flex flex-wrap gap-2.5">
          <RouterLink v-for="category in data.categories" :key="category.slug" :to="{ name: 'jobs', query: { category: category.slug } }" class="inline-flex min-h-11 items-center gap-2 rounded-full border border-line bg-surface px-[18px] text-sm font-medium transition hover:border-primary hover:text-primary">
            <component :is="categoryIcons[category.slug] ?? Briefcase" class="size-4" aria-hidden="true" />{{ category.name }}<span class="text-muted">· {{ category.job_count }}</span>
          </RouterLink>
        </div>
      </section>

      <section class="mt-14" aria-labelledby="latest-title">
        <div class="mb-5 flex items-end justify-between gap-4">
          <div><h2 id="latest-title" class="text-xl">Lowongan terbaru</h2><p class="mt-1 text-[14px] text-muted">Dari perusahaan yang sudah diverifikasi admin Karsa</p></div>
          <RouterLink :to="{ name: 'jobs' }" class="link">Lihat semua <ArrowRight class="size-4" aria-hidden="true" /></RouterLink>
        </div>
        <div v-if="data.latest.length" class="grid gap-4 md:grid-cols-2"><JobCard v-for="job in data.latest" :key="job.id" :job="job" /></div>
        <StatusState v-else title="Belum ada lowongan dibuka" />
      </section>
    </template>
  </div>

  <section v-if="data?.entry_level.length" class="border-y border-line bg-surface py-14" aria-labelledby="entry-title">
    <div class="container-app">
      <div class="mb-5 flex items-end justify-between gap-4">
        <div><h2 id="entry-title" class="text-xl">Magang & Fresh Graduate</h2><p class="mt-1 text-[14px] text-muted">Level magang dan entry, cocok untuk langkah pertama</p></div>
        <RouterLink :to="{ name: 'jobs', query: { level: 'internship,entry' } }" class="link">Lihat semua <ArrowRight class="size-4" aria-hidden="true" /></RouterLink>
      </div>
      <div class="grid gap-4 md:grid-cols-2"><JobCard v-for="job in data.entry_level" :key="job.id" :job="job" /></div>
    </div>
  </section>

  <div v-if="data" class="container-app py-14">
    <section v-if="data.companies.length" aria-labelledby="companies-title">
      <h2 id="companies-title" class="mb-5 text-xl">Perusahaan yang sedang merekrut</h2>
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        <RouterLink v-for="company in data.companies" :key="company.slug" :to="{ name: 'company', params: { slug: company.slug } }" class="card flex flex-col items-center gap-2.5 p-5 text-center transition hover:border-line-strong hover:shadow-card">
          <CompanyLogo :name="company.name" :url="company.logo_url" :size="56" />
          <span class="text-sm font-semibold">{{ company.name }}</span>
          <span class="text-[14px] text-muted">{{ company.open_jobs }} lowongan dibuka</span>
        </RouterLink>
      </div>
    </section>
    <div v-if="auth.role !== 'seeker'" class="mt-14 flex flex-col gap-5 rounded-[20px] bg-primary px-6 py-9 text-white sm:px-12 md:flex-row md:items-center md:justify-between">
      <div><h2 class="text-[22px] sm:text-[26px]">Sedang mencari talenta tech muda?</h2><p class="mt-1.5 opacity-85">Pasang lowongan, saring pelamar dengan skor kecocokan, dan undang kandidat duluan.</p></div>
      <RouterLink :to="auth.role === 'company' ? { name: 'company-job-new' } : { name: 'register', query: { role: 'company' } }" class="btn bg-white text-primary hover:bg-primary-soft">Pasang lowongan gratis</RouterLink>
    </div>
  </div>
</template>
