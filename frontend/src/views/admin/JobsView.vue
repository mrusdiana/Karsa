<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Eye, Lock, Search, Trash2 } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
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
import { JOB_STATUS, formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const company = ref(route.query.company ?? '')
const closing = ref(null)
const deleting = ref(null)
const reason = ref('')
const reasonError = ref('')
const busy = ref(false)
const { data, loading, error, load } = useLoad(async () => (await api.get('/admin/jobs', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })

function update(changes) {
  const query = { ...route.query, ...changes }
  delete query.page
  for (const key of Object.keys(query)) if (!query[key]) delete query[key]
  router.push({ query })
}

async function close() {
  busy.value = true
  reasonError.value = ''
  try {
    await api.post(`/admin/jobs/${closing.value.id}/close`, { reason: reason.value })
    toast.show('Lowongan ditutup.')
    closing.value = null
    load()
  } catch (failure) {
    reasonError.value = fieldErrors(failure).reason ?? ''
    if (!reasonError.value) toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
  }
}

async function remove() {
  busy.value = true
  try {
    await api.delete(`/admin/jobs/${deleting.value.id}`)
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
  <PageHeader title="Lowongan" text="Moderasi lowongan dari semua perusahaan." />
  <div class="mb-4 flex flex-wrap gap-2.5">
    <form class="flex-1" role="search" @submit.prevent="update({ company: company.trim() })">
      <label class="input flex max-w-md items-center gap-2.5"><Search class="size-[18px] text-muted" aria-hidden="true" /><span class="sr-only">Cari perusahaan</span><input v-model="company" class="w-full bg-transparent outline-none focus-visible:shadow-none" placeholder="Cari nama perusahaan" /></label>
    </form>
    <label class="sr-only" for="job-status">Status</label>
    <select id="job-status" class="input w-auto" :value="route.query.status ?? ''" @change="update({ status: $event.target.value })"><option value="">Semua status</option><option v-for="(status, key) in JOB_STATUS" :key="key" :value="key">{{ status.label }}</option></select>
  </div>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="4" height="h-14" />
  <StatusState v-else-if="!data.data.length" title="Tidak ada lowongan yang cocok" />
  <template v-else>
    <div class="card overflow-x-auto">
      <table class="w-full min-w-[760px] text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Lowongan</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 text-right font-semibold">Pelamar</th><th class="px-4 py-3 font-semibold">Dibuat</th><th class="px-4 py-3"><span class="sr-only">Aksi</span></th></tr></thead>
        <tbody>
          <tr v-for="job in data.data" :key="job.id" class="border-t border-line">
            <td class="px-4 py-3"><div class="font-semibold">{{ job.title }}</div><div class="text-[13.5px] text-muted">{{ job.company.name }} · {{ job.category.name }}</div><div v-if="job.moderation_note" class="text-[13.5px] text-danger">Alasan: {{ job.moderation_note }}</div></td>
            <td class="px-4 py-3"><StatusBadge :status="JOB_STATUS[job.status]" /></td>
            <td class="num px-4 py-3 text-right">{{ job.applicant_count }}</td>
            <td class="px-4 py-3 text-[14px] text-muted">{{ formatDate(job.created_at) }}</td>
            <td class="px-4 py-3">
              <div class="flex justify-end gap-0.5">
                <RouterLink v-if="job.status === 'open'" :to="{ name: 'job', params: { slug: job.slug } }" class="btn-icon size-9" :aria-label="`Lihat ${job.title}`" title="Lihat"><Eye class="size-[18px]" aria-hidden="true" /></RouterLink>
                <button v-if="job.status === 'open'" type="button" class="btn-icon size-9" :aria-label="`Tutup ${job.title}`" title="Tutup" @click="closing = job; reason = ''; reasonError = ''"><Lock class="size-[18px]" aria-hidden="true" /></button>
                <button type="button" class="btn-icon size-9 hover:text-danger" :aria-label="`Hapus ${job.title}`" title="Hapus" @click="deleting = job"><Trash2 class="size-[18px]" aria-hidden="true" /></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>

  <ModalDialog :open="!!closing" title="Tutup lowongan" @close="closing = null">
    <form class="flex flex-col gap-4" @submit.prevent="close">
      <p class="text-sm text-muted">{{ closing?.title }} dari {{ closing?.company.name }} akan ditutup dan tidak bisa dibuka lagi oleh perusahaan.</p>
      <FormField id="close-reason" label="Alasan" :error="reasonError" required v-slot="{ describedBy }"><textarea id="close-reason" v-model="reason" class="input" rows="3" maxlength="500" :aria-describedby="describedBy" /></FormField>
      <div class="flex justify-end gap-2"><button type="button" class="btn-ghost" @click="closing = null">Batal</button><button type="submit" class="btn-danger" :disabled="busy">{{ busy ? 'Memproses…' : 'Tutup lowongan' }}</button></div>
    </form>
  </ModalDialog>
  <ConfirmDialog :open="!!deleting" title="Hapus lowongan?" :text="`${deleting?.title} akan dihapus permanen. Lowongan yang sudah punya lamaran hanya bisa ditutup.`" :busy="busy" @confirm="remove" @cancel="deleting = null" />
</template>
