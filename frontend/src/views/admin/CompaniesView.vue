<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CheckCircle2, ExternalLink, XCircle } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import FormField from '@/components/FormField.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import ModalDialog from '@/components/ModalDialog.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, fieldErrors } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import { COMPANY_STATUS, formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const selected = ref(null)
const rejecting = ref(false)
const reason = ref('')
const reasonError = ref('')
const busy = ref('')
const status = () => route.query.status ?? 'pending'
const { data, loading, error, load } = useLoad(async () => (await api.get('/admin/companies', { params: { ...route.query, status: status() } })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })

function open(company) {
  selected.value = company
  rejecting.value = false
  reason.value = ''
  reasonError.value = ''
}

async function decide(action) {
  busy.value = action
  reasonError.value = ''
  try {
    await api.post(`/admin/companies/${selected.value.id}/${action}`, action === 'reject' ? { reason: reason.value } : undefined)
    toast.show(action === 'verify' ? `${selected.value.name} terverifikasi.` : `${selected.value.name} ditolak.`)
    selected.value = null
    load()
  } catch (failure) {
    reasonError.value = fieldErrors(failure).reason ?? ''
    if (!reasonError.value) toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = ''
  }
}
</script>

<template>
  <PageHeader title="Verifikasi Perusahaan" text="Pastikan hanya perusahaan asli yang bisa membuka lowongan." />
  <div class="mb-4 flex gap-1 border-b border-line" role="tablist">
    <button v-for="(item, key) in COMPANY_STATUS" :key="key" type="button" role="tab" :aria-selected="status() === key" class="-mb-px min-h-11 border-b-2 px-3.5 text-sm font-medium whitespace-nowrap" :class="status() === key ? 'border-primary text-primary' : 'border-transparent text-muted hover:text-text'" @click="router.push({ query: { status: key } })">{{ key === 'pending' ? 'Menunggu' : item.label }}</button>
  </div>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="3" height="h-16" />
  <StatusState v-else-if="!data.data.length" title="Tidak ada perusahaan di sini" />
  <template v-else>
    <ul class="flex flex-col gap-3">
      <li v-for="company in data.data" :key="company.id">
        <button type="button" class="card flex w-full items-center gap-4 p-4 text-left transition hover:border-line-strong hover:shadow-card" @click="open(company)">
          <CompanyLogo :name="company.name" :url="company.logo_url" />
          <div class="min-w-0 flex-1"><div class="font-semibold">{{ company.name }}</div><div class="truncate text-[14px] text-muted">{{ company.owner.name }} · {{ company.owner.email }} · daftar {{ formatDate(company.created_at) }}</div></div>
          <StatusBadge :status="COMPANY_STATUS[company.status]" />
        </button>
      </li>
    </ul>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>

  <ModalDialog :open="!!selected" :title="selected?.name ?? ''" wide @close="selected = null">
    <template v-if="selected">
      <div class="flex items-center gap-4"><CompanyLogo :name="selected.name" :url="selected.logo_url" :size="64" /><StatusBadge :status="COMPANY_STATUS[selected.status]" /></div>
      <dl class="mt-5 grid gap-3 text-sm sm:grid-cols-2">
        <div><dt class="text-[14px] text-muted">Pemilik akun</dt><dd class="font-medium">{{ selected.owner.name }}</dd><dd class="text-[14px] text-muted">{{ selected.owner.email }}</dd></div>
        <div><dt class="text-[14px] text-muted">Industri · Ukuran</dt><dd class="font-medium">{{ [selected.industry, selected.size && `${selected.size} karyawan`].filter(Boolean).join(' · ') || '–' }}</dd></div>
        <div><dt class="text-[14px] text-muted">Lokasi</dt><dd class="font-medium">{{ selected.location || '–' }}</dd></div>
        <div><dt class="text-[14px] text-muted">Website</dt><dd><a v-if="selected.website" :href="selected.website" target="_blank" rel="noopener noreferrer" class="link font-medium">{{ selected.website }}<ExternalLink class="size-3.5" aria-hidden="true" /></a><span v-else>–</span></dd></div>
      </dl>
      <p class="mt-4 text-sm whitespace-pre-line">{{ selected.description || 'Belum ada deskripsi.' }}</p>
      <p v-if="selected.rejection_reason" class="mt-4 rounded-btn bg-danger/10 px-3.5 py-2.5 text-sm text-danger">Alasan ditolak: {{ selected.rejection_reason }}</p>
      <FormField v-if="rejecting" id="reject-reason" label="Alasan penolakan" :error="reasonError" required class="mt-5" v-slot="{ describedBy }">
        <textarea id="reject-reason" v-model="reason" class="input" rows="3" maxlength="500" placeholder="Contoh: Website dan deskripsi perusahaan belum diisi." :aria-describedby="describedBy" />
      </FormField>
      <div class="mt-6 flex flex-wrap justify-end gap-2">
        <template v-if="rejecting">
          <button type="button" class="btn-ghost" @click="rejecting = false">Batal</button>
          <button type="button" class="btn-danger" :disabled="!!busy" @click="decide('reject')">{{ busy === 'reject' ? 'Memproses…' : 'Tolak perusahaan' }}</button>
        </template>
        <template v-else>
          <button v-if="selected.status !== 'rejected'" type="button" class="btn-secondary text-danger" @click="rejecting = true"><XCircle class="size-4" aria-hidden="true" />Tolak</button>
          <button v-if="selected.status !== 'verified'" type="button" class="btn-primary" :disabled="!!busy" @click="decide('verify')"><CheckCircle2 class="size-4" aria-hidden="true" />{{ busy === 'verify' ? 'Memproses…' : 'Setujui' }}</button>
        </template>
      </div>
    </template>
  </ModalDialog>
</template>
