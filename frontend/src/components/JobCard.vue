<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { BadgeCheck, Bookmark, Laptop, MapPin } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import api, { errorMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { APPLICATION_STATUS, JOB_LEVEL, WORK_MODE, formatRelative, formatSalary } from '@/utils/format'

const props = defineProps({ job: { type: Object, required: true }, closed: Boolean })
const emit = defineEmits(['unsaved'])
const auth = useAuthStore()
const router = useRouter()
const toast = useToastStore()
const saved = ref(props.job.is_saved ?? false)
const salary = computed(() => formatSalary(props.job))
const extraSkills = computed(() => props.job.skills.length - 4)

async function toggleSave() {
  if (!auth.role) return router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
  const willSave = !saved.value
  saved.value = willSave
  try {
    await api[willSave ? 'put' : 'delete'](`/me/saved-jobs/${props.job.id}`)
    toast.show(willSave ? 'Lowongan disimpan.' : 'Dihapus dari tersimpan.')
    if (!willSave) emit('unsaved', props.job.id)
  } catch (error) {
    saved.value = !willSave
    toast.show(errorMessage(error), 'danger')
  }
}
</script>

<template>
  <article class="group relative flex gap-4 card p-4 transition duration-150 hover:border-line-strong hover:shadow-card sm:p-5">
    <CompanyLogo :name="job.company.name" :url="job.company.logo_url" />
    <div class="min-w-0 flex-1">
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <h3 class="text-base leading-snug font-semibold">
            <RouterLink :to="{ name: 'job', params: { slug: job.slug } }" class="after:absolute after:inset-0 after:rounded-card focus-visible:shadow-none">{{ job.title }}</RouterLink>
          </h3>
          <p class="mt-0.5 flex items-center gap-1 text-[14px] font-medium text-muted">
            {{ job.company.name }}<BadgeCheck v-if="job.company.verified" class="size-[15px] text-primary" aria-label="Terverifikasi" />
          </p>
        </div>
        <StatusBadge v-if="closed" :status="{ label: 'Ditutup', icon: 'Lock', tone: 'neutral' }" />
        <StatusBadge v-else-if="job.application_status" :status="{ ...APPLICATION_STATUS[job.application_status], label: 'Sudah dilamar', icon: 'Send', tone: 'info' }" />
      </div>
      <div class="mt-2.5 mb-3 flex flex-wrap gap-1.5">
        <span class="chip bg-primary-soft text-primary">{{ JOB_LEVEL[job.level] }}</span>
        <span class="chip"><Laptop class="size-3.5 text-muted" aria-hidden="true" />{{ WORK_MODE[job.work_mode] }}</span>
        <span v-if="job.location" class="chip"><MapPin class="size-3.5 text-muted" aria-hidden="true" />{{ job.location }}</span>
      </div>
      <div class="flex items-center justify-between gap-3 pr-10">
        <span class="num text-[15.5px]" :class="salary ? 'font-semibold' : 'font-medium text-muted'">{{ salary ?? 'Gaji dirahasiakan' }}</span>
        <span class="text-[14px] text-muted">{{ formatRelative(job.created_at) }}</span>
      </div>
      <div class="mt-3 flex flex-wrap gap-1.5 pr-10">
        <span v-for="skill in job.skills.slice(0, 4)" :key="skill" class="chip-skill">{{ skill }}</span>
        <span v-if="extraSkills > 0" class="chip-skill text-muted">+{{ extraSkills }}</span>
      </div>
    </div>
    <button
      v-if="auth.role !== 'company' && auth.role !== 'admin'"
      type="button"
      class="btn-icon absolute right-3 bottom-3 z-10"
      :class="saved && 'text-primary hover:text-primary'"
      :aria-label="saved ? 'Hapus dari tersimpan' : 'Simpan lowongan'"
      :aria-pressed="saved"
      @click="toggleSave"
    >
      <Bookmark class="size-5" :fill="saved ? 'currentColor' : 'none'" aria-hidden="true" />
    </button>
  </article>
</template>
