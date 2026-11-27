<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Eye, Lock, Pencil, Plus, Trash2, Unlock, Users } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { JOB_STATUS, deadlineInfo } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const deleting = ref(null)
const busy = ref(false)
const { data, loading, error, load } = useLoad(async () => (await api.get('/company/jobs', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })

async function setOpen(job, open) {
  try {
    await api.post(`/company/jobs/${job.id}/${open ? 'open' : 'close'}`)
    toast.show(open ? 'Lowongan dibuka.' : 'Lowongan ditutup.')
    load()
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  }
}

async function remove() {
  busy.value = true
  try {
    await api.delete(`/company/jobs/${deleting.value.id}`)
    toast.show('Lowongan dihapus.')
    load()
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
    deleting.value = null
  }
}
</script>

<template>
  <PageHeader title="Lowongan" text="Kelola lowongan dan lihat pelamarnya.">
    <RouterLink :to="{ name: 'company-job-new' }" class="btn-primary"><Plus class="size-4" aria-hidden="true" />Buat lowongan</RouterLink>
  </PageHeader>
  <div class="mb-4 flex gap-1 overflow-x-auto border-b border-line" role="tablist">
    <button v-for="option in [{ key: '', label: 'Semua' }, ...Object.entries(JOB_STATUS).map(([key, status]) => ({ key, label: status.label }))]" :key="option.key" type="button" role="tab" :aria-selected="(route.query.status ?? '') === option.key" class="-mb-px min-h-11 border-b-2 px-3.5 text-sm font-medium" :class="(route.query.status ?? '') === option.key ? 'border-primary text-primary' : 'border-transparent text-muted hover:text-text'" @click="router.push({ query: option.key ? { status: option.key } : {} })">{{ option.label }}</button>
  </div>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="3" height="h-16" />
  <StatusState v-else-if="!data.data.length" title="Belum ada lowongan" text="Buat lowongan pertamamu. Draft bisa disiapkan sebelum perusahaan terverifikasi.">
    <RouterLink :to="{ name: 'company-job-new' }" class="btn-primary mt-2">Buat lowongan</RouterLink>
  </StatusState>
  <template v-else>
    <div class="card overflow-x-auto">
      <table class="w-full min-w-[760px] text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Lowongan</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 text-right font-semibold">Pelamar</th><th class="px-4 py-3 text-right font-semibold">Dilihat</th><th class="px-4 py-3 font-semibold">Deadline</th><th class="px-4 py-3 text-right font-semibold"><span class="sr-only">Aksi</span></th></tr></thead>
        <tbody>
          <tr v-for="job in data.data" :key="job.id" class="border-t border-line align-middle">
            <td class="px-4 py-3.5"><div class="font-semibold">{{ job.title }}</div><div class="text-[13.5px] text-muted">{{ job.category.name }}</div><div v-if="job.moderation_note" class="mt-1 text-[13.5px] text-danger">Ditutup admin: {{ job.moderation_note }}</div></td>
            <td class="px-4 py-3.5"><StatusBadge :status="JOB_STATUS[job.status]" /></td>
            <td class="num px-4 py-3.5 text-right">{{ job.applicant_count }}</td>
            <td class="num px-4 py-3.5 text-right">{{ job.views }}</td>
            <td class="px-4 py-3.5 text-[14px]" :class="job.status === 'open' && deadlineInfo(job.deadline).urgent ? 'text-warning' : 'text-muted'">{{ job.status === 'open' ? deadlineInfo(job.deadline).text : '–' }}</td>
            <td class="px-4 py-3.5">
              <div class="flex justify-end gap-0.5">
                <RouterLink :to="{ name: 'company-applicants', params: { id: job.id } }" class="btn-icon size-9" :aria-label="`Pelamar ${job.title}`" title="Lihat pelamar"><Users class="size-[18px]" aria-hidden="true" /></RouterLink>
                <RouterLink v-if="job.status === 'open'" :to="{ name: 'job', params: { slug: job.slug } }" class="btn-icon size-9" :aria-label="`Lihat ${job.title} di halaman publik`" title="Lihat di halaman publik"><Eye class="size-[18px]" aria-hidden="true" /></RouterLink>
                <RouterLink :to="{ name: 'company-job-edit', params: { id: job.id } }" class="btn-icon size-9" :aria-label="`Edit ${job.title}`" title="Edit"><Pencil class="size-[18px]" aria-hidden="true" /></RouterLink>
                <button v-if="job.status === 'open'" type="button" class="btn-icon size-9" :aria-label="`Tutup ${job.title}`" title="Tutup" @click="setOpen(job, false)"><Lock class="size-[18px]" aria-hidden="true" /></button>
                <button v-else-if="!job.moderation_note" type="button" class="btn-icon size-9" :disabled="!auth.isVerified" :aria-label="`Buka ${job.title}`" :title="auth.isVerified ? 'Buka' : 'Tunggu verifikasi perusahaan'" @click="setOpen(job, true)"><Unlock class="size-[18px]" aria-hidden="true" /></button>
                <button type="button" class="btn-icon size-9 hover:text-danger" :aria-label="`Hapus ${job.title}`" title="Hapus" @click="deleting = job"><Trash2 class="size-[18px]" aria-hidden="true" /></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>
  <ConfirmDialog :open="!!deleting" title="Hapus lowongan?" :text="`${deleting?.title} akan dihapus permanen. Lowongan yang sudah punya pelamar tidak bisa dihapus; tutup saja.`" :busy="busy" @confirm="remove" @cancel="deleting = null" />
</template>
