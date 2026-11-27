<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { GraduationCap, Lock, MapPin, Send } from 'lucide-vue-next'
import FormField from '@/components/FormField.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import ModalDialog from '@/components/ModalDialog.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import SkillInput from '@/components/SkillInput.vue'
import StatusState from '@/components/StatusState.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { EDUCATION } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const filters = reactive({
  skills: (route.query.skills ?? '').split(',').filter(Boolean),
  location: route.query.location ?? '',
  education: route.query.education ?? '',
  grad_from: route.query.grad_from ?? '',
  grad_to: route.query.grad_to ?? '',
  max_experience: route.query.max_experience ?? '',
})
const inviting = ref(null)
const invite = reactive({ job_id: '', message: '' })
const sending = ref(false)
const inviteError = ref('')
const openJobs = ref([])
const remaining = ref(null)
const highlighted = computed(() => new Set(filters.skills.map((skill) => skill.toLowerCase())))
const { data, loading, error, load } = useLoad(async () => (await api.get('/company/talents', { params: route.query })).data, null, { immediate: false })

function search() {
  const query = { ...filters, skills: filters.skills.join(',') }
  for (const key of Object.keys(query)) if (!query[key]) delete query[key]
  router.push({ query })
}

async function openInvite(talent) {
  inviting.value = talent
  inviteError.value = ''
  Object.assign(invite, { job_id: '', message: '' })
  const [jobs, sent] = await Promise.all([api.get('/company/jobs', { params: { status: 'open' } }), api.get('/company/invitations')])
  openJobs.value = jobs.data.data
  remaining.value = sent.data.remaining_today
}

async function sendInvite() {
  sending.value = true
  inviteError.value = ''
  try {
    await api.post('/company/invitations', { user_id: inviting.value.id, job_id: Number(invite.job_id), message: invite.message })
    inviting.value.invited_job_ids.push(Number(invite.job_id))
    toast.show(`Undangan terkirim ke ${inviting.value.name}.`)
    inviting.value = null
  } catch (failure) {
    inviteError.value = errorMessage(failure)
  } finally {
    sending.value = false
  }
}

watch(() => route.query, () => auth.isVerified && load(), { immediate: true })
</script>

<template>
  <PageHeader title="Cari Talenta" text="Kandidat yang membuka profilnya untuk perusahaan terverifikasi." />
  <StatusState v-if="!auth.isVerified" title="Fitur ini tersedia setelah perusahaanmu terverifikasi." text="Lengkapi profil perusahaan supaya admin bisa memverifikasinya.">
    <RouterLink :to="{ name: 'company-profile' }" class="btn-primary mt-2">Buka profil perusahaan</RouterLink>
  </StatusState>
  <div v-else class="grid items-start gap-6 xl:grid-cols-[280px_1fr]">
    <form class="card flex flex-col gap-4 p-5 xl:sticky xl:top-6" @submit.prevent="search">
      <FormField id="talent-skills" label="Skill (harus dimiliki semua)" v-slot="{ describedBy }"><SkillInput id="talent-skills" v-model="filters.skills" :max="5" :described-by="describedBy" /></FormField>
      <FormField id="talent-location" label="Lokasi"><input id="talent-location" v-model="filters.location" class="input" placeholder="Contoh: Bandung" /></FormField>
      <FormField id="talent-education" label="Pendidikan"><select id="talent-education" v-model="filters.education" class="input"><option value="">Semua</option><option v-for="level in EDUCATION" :key="level" :value="level">{{ level }}</option></select></FormField>
      <div class="grid grid-cols-2 gap-3">
        <FormField id="grad-from" label="Lulus dari"><input id="grad-from" v-model="filters.grad_from" type="number" class="input" placeholder="2024" /></FormField>
        <FormField id="grad-to" label="Sampai"><input id="grad-to" v-model="filters.grad_to" type="number" class="input" placeholder="2027" /></FormField>
      </div>
      <FormField id="max-experience" label="Pengalaman maks. (tahun)"><input id="max-experience" v-model="filters.max_experience" type="number" min="0" class="input" /></FormField>
      <button type="submit" class="btn-primary">Cari kandidat</button>
    </form>

    <section aria-live="polite">
      <StatusState v-if="error" type="error" @retry="load" />
      <LoadingCards v-else-if="loading || !data" :count="3" />
      <StatusState v-else-if="!data.data.length" title="Belum ada kandidat yang cocok" text="Coba kurangi filter skill atau lokasi." />
      <template v-else>
        <p class="mb-3 text-sm text-muted">{{ data.total }} kandidat ditemukan</p>
        <div class="grid gap-4 md:grid-cols-2">
          <article v-for="talent in data.data" :key="talent.id" class="card flex flex-col gap-3.5 p-5">
            <div class="flex items-start gap-3">
              <UserAvatar :name="talent.name" :url="talent.avatar_url" :size="44" />
              <div class="min-w-0"><h2 class="font-sans text-base font-semibold">{{ talent.name }}</h2><p class="text-[14px] text-muted">{{ talent.headline }}</p></div>
            </div>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-[14px] text-muted">
              <span v-if="talent.location" class="flex items-center gap-1.5"><MapPin class="size-3.5" aria-hidden="true" />{{ talent.location }}</span>
              <span class="flex items-center gap-1.5"><GraduationCap class="size-3.5" aria-hidden="true" />{{ [talent.education, talent.major].filter(Boolean).join(' ') }} · lulus {{ talent.graduation_year }}</span>
              <span>{{ talent.experience_years }} tahun pengalaman</span>
            </div>
            <div class="flex flex-wrap gap-1.5"><span v-for="skill in talent.skills" :key="skill" class="chip-skill" :class="highlighted.has(skill.toLowerCase()) && 'border-primary bg-primary-soft text-primary'">{{ skill }}</span></div>
            <div class="mt-auto flex items-center justify-between gap-3 border-t border-line pt-3.5">
              <span class="flex items-center gap-1.5 text-xs text-muted"><Lock class="size-3.5" aria-hidden="true" />Kontak disembunyikan</span>
              <button type="button" class="btn-primary btn-sm" @click="openInvite(talent)"><Send class="size-4" aria-hidden="true" />{{ talent.invited_job_ids.length ? `Diundang (${talent.invited_job_ids.length})` : 'Undang' }}</button>
            </div>
          </article>
        </div>
        <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
      </template>
    </section>
  </div>

  <ModalDialog :open="!!inviting" :title="`Undang ${inviting?.name ?? ''}`" @close="inviting = null">
    <form class="flex flex-col gap-4" @submit.prevent="sendInvite">
      <FormField id="invite-job" label="Lowongan" required>
        <select id="invite-job" v-model="invite.job_id" class="input" required>
          <option value="" disabled>Pilih lowongan yang masih dibuka</option>
          <option v-for="job in openJobs" :key="job.id" :value="job.id" :disabled="inviting?.invited_job_ids.includes(job.id)">{{ job.title }}{{ inviting?.invited_job_ids.includes(job.id) ? ' (sudah diundang)' : '' }}</option>
        </select>
      </FormField>
      <p v-if="!openJobs.length" class="text-[14px] text-warning">Belum ada lowongan yang dibuka.</p>
      <FormField id="invite-message" label="Pesan (opsional)" :hint="`${invite.message.length}/500`"><textarea id="invite-message" v-model="invite.message" class="input" rows="3" maxlength="500" placeholder="Profilmu cocok dengan tim kami. Yuk, melamar!" /></FormField>
      <p v-if="remaining != null" class="text-[14px] text-muted">{{ remaining }} undangan tersisa hari ini</p>
      <p v-if="inviteError" class="rounded-btn bg-danger/10 px-3.5 py-2.5 text-sm text-danger" role="alert">{{ inviteError }}</p>
      <div class="flex justify-end gap-2"><button type="button" class="btn-ghost" @click="inviting = null">Batal</button><button type="submit" class="btn-primary" :disabled="sending || !invite.job_id || remaining === 0">{{ sending ? 'Mengirim…' : 'Kirim undangan' }}</button></div>
    </form>
  </ModalDialog>
</template>
