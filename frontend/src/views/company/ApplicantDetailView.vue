<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { CheckCircle2, ChevronRight, Download, ExternalLink, Mail, MapPin, Phone, XCircle } from 'lucide-vue-next'
import ApplicationTimeline from '@/components/ApplicationTimeline.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import MatchScore from '@/components/MatchScore.vue'
import ScoreBreakdown from '@/components/ScoreBreakdown.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, openFile } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import { APPLICATION_STATUS, COMPANY_MOVES, SCORE_LABEL, formatDate, scoreLevel } from '@/utils/format'

const route = useRoute()
const toast = useToastStore()
const note = ref('')
const busy = ref('')
const rejecting = ref(false)
const announcement = ref('')
const { data: applicant, loading, error, load } = useLoad(async () => {
  const result = (await api.get(`/company/applications/${route.params.id}`)).data
  note.value = result.company_note ?? ''
  return result
})
const moves = computed(() => COMPANY_MOVES[applicant.value?.status] ?? [])
const moveLabels = { reviewed: 'Tandai ditinjau', interview: 'Lanjut interview', accepted: 'Terima', rejected: 'Tolak' }

async function update(body, action) {
  busy.value = action
  try {
    applicant.value = (await api.patch(`/company/applications/${route.params.id}`, body)).data
    if (body.status) announcement.value = `Status diubah menjadi ${APPLICATION_STATUS[body.status].label}`
    toast.show(body.status ? announcement.value : 'Catatan disimpan.')
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
    if (failure.response?.status === 409) load()
  } finally {
    busy.value = ''
    rejecting.value = false
  }
}
</script>

<template>
  <p class="sr-only" aria-live="polite">{{ announcement }}</p>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="2" height="h-56" />
  <template v-else>
    <nav class="mb-4 flex flex-wrap items-center gap-1.5 text-[14px] text-muted" aria-label="Breadcrumb">
      <RouterLink :to="{ name: 'company-jobs' }" class="hover:text-text">Lowongan</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" />
      <RouterLink :to="{ name: 'company-applicants', params: { id: applicant.job.id } }" class="hover:text-text">{{ applicant.job.title }}</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" />
      <span class="text-text">{{ applicant.candidate.name }}</span>
    </nav>
    <div class="grid items-start gap-5 xl:grid-cols-[1fr_400px]">
      <div class="flex flex-col gap-5">
        <section class="card p-5 sm:p-6">
          <div class="flex flex-wrap items-start gap-4">
            <UserAvatar :name="applicant.candidate.name" :url="applicant.candidate.avatar_url" :size="56" />
            <div class="min-w-0 flex-1">
              <h1 class="text-2xl">{{ applicant.candidate.name }}</h1>
              <p class="text-sm text-muted">{{ applicant.candidate.headline }}</p>
              <div class="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-[14px] text-muted">
                <a :href="`mailto:${applicant.candidate.email}`" class="flex items-center gap-1.5 hover:text-primary"><Mail class="size-3.5" aria-hidden="true" />{{ applicant.candidate.email }}</a>
                <span v-if="applicant.candidate.phone" class="flex items-center gap-1.5"><Phone class="size-3.5" aria-hidden="true" />{{ applicant.candidate.phone }}</span>
                <span v-if="applicant.candidate.location" class="flex items-center gap-1.5"><MapPin class="size-3.5" aria-hidden="true" />{{ applicant.candidate.location }}</span>
              </div>
            </div>
            <StatusBadge :status="APPLICATION_STATUS[applicant.status]" />
          </div>
          <dl class="mt-5 grid gap-4 border-t border-line pt-5 text-sm sm:grid-cols-3">
            <div><dt class="text-[14px] text-muted">Pendidikan</dt><dd class="font-medium">{{ [applicant.candidate.education, applicant.candidate.major].filter(Boolean).join(' ') || '–' }}</dd></div>
            <div><dt class="text-[14px] text-muted">Tahun lulus</dt><dd class="font-medium">{{ applicant.candidate.graduation_year ?? '–' }}</dd></div>
            <div><dt class="text-[14px] text-muted">Pengalaman</dt><dd class="font-medium">{{ applicant.candidate.experience_years }} tahun</dd></div>
          </dl>
          <div class="mt-4 flex flex-wrap gap-1.5"><span v-for="skill in applicant.candidate.skills" :key="skill" class="chip-skill">{{ skill }}</span></div>
          <p v-if="applicant.candidate.bio" class="mt-4 text-sm whitespace-pre-line">{{ applicant.candidate.bio }}</p>
          <div class="mt-5 flex flex-wrap gap-2">
            <button type="button" class="btn-secondary" @click="openFile(applicant.cv_url)"><Download class="size-4" aria-hidden="true" />Unduh CV</button>
            <a v-if="applicant.candidate.portfolio_url" :href="applicant.candidate.portfolio_url" target="_blank" rel="noopener noreferrer" class="btn-ghost"><ExternalLink class="size-4" aria-hidden="true" />Portofolio</a>
          </div>
        </section>
        <section v-if="applicant.cover_letter" class="card p-5 sm:p-6"><h2 class="mb-2 text-lg">Surat pengantar</h2><p class="text-sm whitespace-pre-line">{{ applicant.cover_letter }}</p></section>
        <section v-if="applicant.answers.length" class="card p-5 sm:p-6">
          <h2 class="mb-2 text-lg">Jawaban screening</h2>
          <dl class="divide-y divide-line text-sm">
            <div v-for="answer in applicant.answers" :key="answer.question" class="flex items-start justify-between gap-4 py-3">
              <dt>{{ answer.question }}<span v-if="answer.required" class="text-danger"> *</span></dt>
              <dd class="flex shrink-0 items-center gap-1.5 font-semibold" :class="answer.passed ? 'text-success' : 'text-danger'">
                <CheckCircle2 v-if="answer.passed" class="size-4" aria-label="Lolos" /><XCircle v-else class="size-4" aria-label="Tidak lolos" />
                {{ answer.answer === 'yes' ? 'Ya' : answer.answer === 'no' ? 'Tidak' : answer.answer || 'Tidak dijawab' }}
              </dd>
            </div>
          </dl>
        </section>
      </div>

      <aside class="flex flex-col gap-5 xl:sticky xl:top-6">
        <section class="card p-5">
          <div class="flex items-center gap-3"><MatchScore :score="applicant.match_score" :size="48" /><div><div class="font-semibold">{{ SCORE_LABEL[scoreLevel(applicant.match_score)] }}</div><div class="text-[14px] text-muted">{{ applicant.meets_requirements ? 'Memenuhi semua syarat wajib' : 'Tidak memenuhi syarat wajib' }}</div></div></div>
          <ScoreBreakdown :breakdown="applicant.score_breakdown" class="mt-4" />
        </section>
        <section class="card p-5">
          <h2 class="mb-3 text-base">Ubah status</h2>
          <div v-if="moves.length" class="flex flex-wrap gap-2">
            <button v-for="status in moves.filter((move) => move !== 'rejected')" :key="status" type="button" class="btn-primary btn-sm" :disabled="!!busy" @click="update({ status }, status)"><CheckCircle2 class="size-4" aria-hidden="true" />{{ busy === status ? 'Menyimpan…' : moveLabels[status] }}</button>
            <button type="button" class="btn-secondary btn-sm text-danger" :disabled="!!busy" @click="rejecting = true"><XCircle class="size-4" aria-hidden="true" />Tolak</button>
          </div>
          <p v-else class="text-sm text-muted">Status sudah final.</p>
        </section>
        <section class="card p-5">
          <label for="note" class="mb-2 block text-base font-semibold">Catatan internal <span class="text-[14px] font-normal text-muted">(tidak terlihat pelamar)</span></label>
          <textarea id="note" v-model="note" class="input" rows="3" maxlength="1000" />
          <button type="button" class="btn-secondary btn-sm mt-2" :disabled="!!busy || note === (applicant.company_note ?? '')" @click="update({ company_note: note }, 'note')">{{ busy === 'note' ? 'Menyimpan…' : 'Simpan catatan' }}</button>
        </section>
        <section class="card p-5"><h2 class="mb-4 text-base">Timeline</h2><ApplicationTimeline :events="applicant.events" /><p class="mt-3 text-xs text-muted">Dilamar {{ formatDate(applicant.created_at) }}</p></section>
      </aside>
    </div>
    <ConfirmDialog :open="rejecting" title="Tolak pelamar ini?" :text="`${applicant.candidate.name} akan menerima email bahwa lamarannya tidak dilanjutkan.`" confirm-label="Tolak" :busy="busy === 'rejected'" @confirm="update({ status: 'rejected' }, 'rejected')" @cancel="rejecting = false" />
  </template>
</template>
