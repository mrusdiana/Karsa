<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, BadgeCheck, Bookmark, Briefcase, ChevronDown, ChevronRight, Clock, Laptop, Layers, MapPin, Users } from 'lucide-vue-next'
import ApplyDialog from '@/components/ApplyDialog.vue'
import CompanyLogo from '@/components/CompanyLogo.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import MatchScore from '@/components/MatchScore.vue'
import ScoreBreakdown from '@/components/ScoreBreakdown.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'
import { useToastStore } from '@/stores/toast'
import { APPLICATION_STATUS, JOB_LEVEL, JOB_TYPE, SCORE_LABEL, WORK_MODE, deadlineInfo, formatRelative, formatSalary, scoreLevel } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const notifications = useNotificationStore()
const breakdownOpen = ref(false)
const applyOpen = ref(false)
const saved = ref(false)
const { data: job, loading, error, load } = useLoad(async () => (await api.get(`/jobs/${route.params.slug}`)).data, null, { immediate: false })

const salary = computed(() => formatSalary(job.value))
const deadline = computed(() => deadlineInfo(job.value.deadline))
const invitationId = computed(() => (route.query.undangan ? Number(route.query.undangan) : undefined))
const matchedSkills = computed(() => new Set(job.value.match_preview?.breakdown.skill.matched ?? []))
const paragraphs = (text) => text.split('\n').map((line) => line.trim()).filter(Boolean)

async function refreshAfterApply() {
  job.value = (await api.get(`/jobs/${route.params.slug}`)).data
  notifications.load()
  if (invitationId.value) router.replace({ query: {} })
}

function requireSeeker() {
  if (auth.role) return true
  router.push({ name: 'login', query: { redirect: route.fullPath } })
  return false
}

function openApply() {
  if (requireSeeker()) applyOpen.value = true
}

async function toggleSave() {
  if (!requireSeeker()) return
  const willSave = !saved.value
  saved.value = willSave
  try {
    await api[willSave ? 'put' : 'delete'](`/me/saved-jobs/${job.value.id}`)
    toast.show(willSave ? 'Lowongan disimpan.' : 'Dihapus dari tersimpan.')
  } catch (failure) {
    saved.value = !willSave
    toast.show(errorMessage(failure), 'danger')
  }
}

watch(
  () => route.params.slug,
  async () => {
    await load()
    saved.value = job.value?.is_saved ?? false
    if (invitationId.value && job.value && !job.value.application_status) applyOpen.value = true
  },
  { immediate: true },
)
</script>

<template>
  <div class="container-app py-6 pb-28 lg:pb-16">
    <StatusState v-if="error" type="error" title="Lowongan tidak ditemukan atau sudah ditutup" @retry="load">
      <RouterLink :to="{ name: 'jobs' }" class="btn-primary mt-2">Cari lowongan lain</RouterLink>
    </StatusState>
    <LoadingCards v-else-if="loading || !job" :count="2" height="h-64" />
    <template v-else>
      <nav class="mb-5 flex flex-wrap items-center gap-1.5 text-[14px] text-muted" aria-label="Breadcrumb">
        <RouterLink :to="{ name: 'jobs' }" class="hover:text-text">Cari Lowongan</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" />
        <RouterLink :to="{ name: 'jobs', query: { category: job.category.slug } }" class="hover:text-text">{{ job.category.name }}</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" />
        <span class="text-text" aria-current="page">{{ job.title }}</span>
      </nav>

      <div class="grid items-start gap-6 lg:grid-cols-[1fr_380px] lg:gap-8">
        <div class="flex flex-col gap-5">
          <section class="card p-5 sm:p-7">
            <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:gap-5">
              <CompanyLogo :name="job.company.name" :url="job.company.logo_url" :size="64" />
              <div>
                <h1 class="text-2xl leading-tight sm:text-[28px]">{{ job.title }}</h1>
                <p class="mt-1 flex flex-wrap items-center gap-1.5 font-medium text-muted">
                  <RouterLink :to="{ name: 'company', params: { slug: job.company.slug } }" class="hover:text-primary">{{ job.company.name }}</RouterLink>
                  <BadgeCheck v-if="job.company.verified" class="size-4 text-primary" aria-label="Terverifikasi" />
                  <span>· Dipasang {{ formatRelative(job.created_at).toLowerCase() }} · {{ job.views }} kali dilihat</span>
                </p>
                <div class="mt-3.5 flex flex-wrap gap-1.5">
                  <span class="chip bg-primary-soft text-primary">{{ JOB_LEVEL[job.level] }}</span>
                  <span class="chip"><Briefcase class="size-3.5 text-muted" aria-hidden="true" />{{ JOB_TYPE[job.type] }}</span>
                  <span class="chip"><Laptop class="size-3.5 text-muted" aria-hidden="true" />{{ WORK_MODE[job.work_mode] }}</span>
                  <span v-if="job.location" class="chip"><MapPin class="size-3.5 text-muted" aria-hidden="true" />{{ job.location }}</span>
                </div>
              </div>
            </div>
          </section>

          <section class="card flex flex-col gap-7 p-5 sm:p-7">
            <div><h2 class="mb-2.5 text-xl">Deskripsi</h2><p v-for="(line, index) in paragraphs(job.description)" :key="index" class="mb-2 max-w-[68ch] text-[#2A2F3A]">{{ line }}</p></div>
            <div><h2 class="mb-2.5 text-xl">Persyaratan</h2><ul class="flex max-w-[68ch] list-disc flex-col gap-1.5 pl-5 text-[#2A2F3A]"><li v-for="(line, index) in paragraphs(job.requirements)" :key="index">{{ line }}</li></ul></div>
            <div>
              <h2 class="mb-3 text-xl">Skill</h2>
              <div class="flex flex-wrap gap-2">
                <span v-for="skill in job.skills" :key="skill" class="chip-skill" :class="matchedSkills.has(skill) && 'border-success text-success'">{{ skill }}<span v-if="matchedSkills.has(skill)" class="sr-only"> (kamu punya)</span></span>
              </div>
            </div>
          </section>

          <section class="card flex gap-4 p-5 sm:p-6">
            <CompanyLogo :name="job.company.name" :url="job.company.logo_url" :size="48" />
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center justify-between gap-2"><h2 class="text-[18px]">Tentang {{ job.company.name }}</h2><RouterLink :to="{ name: 'company', params: { slug: job.company.slug } }" class="link">Lihat profil <ArrowRight class="size-4" aria-hidden="true" /></RouterLink></div>
              <p class="mt-0.5 text-[14px] text-muted">{{ [job.company.industry, job.company.size && `${job.company.size} karyawan`, job.company.location].filter(Boolean).join(' · ') }}</p>
              <p v-if="job.company.description" class="mt-2.5 text-sm">{{ job.company.description }}</p>
            </div>
          </section>
        </div>

        <aside class="card order-first p-5 sm:p-6 lg:sticky lg:top-24 lg:order-none">
          <div class="text-[14px] text-muted">Gaji</div>
          <div class="num mt-0.5 font-display text-2xl font-extrabold" :class="!salary && 'text-lg text-muted'">{{ salary ? salary.replace('/bulan', '') : 'Gaji dirahasiakan' }}<span v-if="salary" class="font-sans text-sm font-medium text-muted">/bulan</span></div>
          <dl class="my-5 grid grid-cols-[auto_1fr] gap-x-4 gap-y-3 border-y border-line py-4 text-sm">
            <dt class="flex items-center gap-2 text-muted"><Briefcase class="size-4" aria-hidden="true" />Tipe</dt><dd class="text-right font-medium">{{ JOB_TYPE[job.type] }}</dd>
            <dt class="flex items-center gap-2 text-muted"><Layers class="size-4" aria-hidden="true" />Level</dt><dd class="text-right font-medium">{{ JOB_LEVEL[job.level] }}</dd>
            <dt class="flex items-center gap-2 text-muted"><Laptop class="size-4" aria-hidden="true" />Mode</dt><dd class="text-right font-medium">{{ WORK_MODE[job.work_mode] }}</dd>
            <dt class="flex items-center gap-2 text-muted"><MapPin class="size-4" aria-hidden="true" />Lokasi</dt><dd class="text-right font-medium">{{ job.location ?? 'Remote' }}</dd>
            <dt class="flex items-center gap-2 text-muted"><Clock class="size-4" aria-hidden="true" />Deadline</dt><dd class="text-right font-medium" :class="deadline.urgent && 'text-warning'">{{ deadline.text }}</dd>
            <dt class="flex items-center gap-2 text-muted"><Users class="size-4" aria-hidden="true" />Pelamar</dt><dd class="text-right font-medium">{{ job.applicant_count }} orang</dd>
          </dl>

          <template v-if="job.match_preview">
            <button type="button" class="flex w-full items-center gap-3.5 text-left" :aria-expanded="breakdownOpen" @click="breakdownOpen = !breakdownOpen">
              <MatchScore :score="job.match_preview.score" />
              <span class="flex-1"><span class="block text-sm font-semibold">{{ SCORE_LABEL[scoreLevel(job.match_preview.score)] }}</span><span class="text-[14px] text-muted">Perkiraan skor kecocokanmu</span></span>
              <ChevronDown class="size-5 text-muted transition" :class="breakdownOpen && 'rotate-180'" aria-hidden="true" />
            </button>
            <ScoreBreakdown v-if="breakdownOpen" :breakdown="job.match_preview.breakdown" preview class="mt-3.5" />
          </template>

          <template v-if="auth.role !== 'company' && auth.role !== 'admin'">
            <div v-if="job.application_status" class="mt-5 flex items-center justify-between rounded-btn bg-info/10 px-4 py-3 text-sm font-medium text-info">
              Sudah dilamar <StatusBadge :status="APPLICATION_STATUS[job.application_status]" />
            </div>
            <div class="fixed inset-x-0 bottom-16 z-30 flex gap-2.5 border-t border-line bg-surface p-3 md:bottom-0 lg:static lg:mt-5 lg:flex-col lg:border-0 lg:p-0">
              <button type="button" class="btn-secondary size-12 shrink-0 px-0 lg:hidden" :aria-pressed="saved" :aria-label="saved ? 'Hapus dari tersimpan' : 'Simpan lowongan'" @click="toggleSave"><Bookmark class="size-5" :fill="saved ? 'currentColor' : 'none'" aria-hidden="true" /></button>
              <RouterLink v-if="job.application_status" :to="{ name: 'applications' }" class="btn-secondary h-12 flex-1">Lihat lamaran</RouterLink>
              <button v-else type="button" class="btn-primary h-12 flex-1 text-[16px]" @click="openApply">{{ invitationId ? 'Terima & lamar' : 'Lamar sekarang' }}</button>
              <button type="button" class="btn-secondary hidden lg:flex" :aria-pressed="saved" @click="toggleSave"><Bookmark class="size-4" :fill="saved ? 'currentColor' : 'none'" aria-hidden="true" />{{ saved ? 'Tersimpan' : 'Simpan' }}</button>
            </div>
          </template>
        </aside>
      </div>
      <ApplyDialog v-if="auth.role === 'seeker'" :open="applyOpen" :job="job" :invitation-id="invitationId" @close="applyOpen = false" @applied="refreshAfterApply" />
    </template>
  </div>
</template>
