<script setup>
import { computed, ref } from 'vue'
import { Clock } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage } from '@/services/api'
import { useNotificationStore } from '@/stores/notifications'
import { useToastStore } from '@/stores/toast'
import { INVITATION_STATUS, JOB_LEVEL, deadlineInfo, formatDate } from '@/utils/format'

const toast = useToastStore()
const tab = ref('pending')
const declining = ref(null)
const busy = ref(false)
const { data, loading, error, load } = useLoad(async () => (await api.get('/me/invitations')).data, [])
const visible = computed(() => data.value.filter((invitation) => (tab.value === 'pending') === (invitation.status === 'pending')))
const pendingCount = computed(() => data.value.filter((invitation) => invitation.status === 'pending').length)

async function decline() {
  busy.value = true
  try {
    const updated = (await api.post(`/me/invitations/${declining.value.id}/decline`)).data
    data.value = data.value.map((invitation) => (invitation.id === updated.id ? updated : invitation))
    toast.show('Undangan ditolak.')
    useNotificationStore().load()
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
    declining.value = null
  }
}
</script>

<template>
  <div class="container-app py-8 pb-16">
    <PageHeader title="Undangan" text="Perusahaan yang tertarik dengan profilmu dan mengundangmu melamar." />
    <div class="mb-5 flex gap-1 border-b border-line" role="tablist">
      <button v-for="option in [{ key: 'pending', label: `Menunggu (${pendingCount})` }, { key: 'history', label: 'Riwayat' }]" :key="option.key" type="button" role="tab" :aria-selected="tab === option.key" class="-mb-px min-h-11 border-b-2 px-3.5 text-sm font-medium" :class="tab === option.key ? 'border-primary text-primary' : 'border-transparent text-muted hover:text-text'" @click="tab = option.key">{{ option.label }}</button>
    </div>
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="2" />
    <StatusState v-else-if="!visible.length" :title="tab === 'pending' ? 'Belum ada undangan baru' : 'Belum ada riwayat undangan'" text="Aktifkan &quot;Bisa ditemukan perusahaan&quot; di profil supaya perusahaan bisa mengundangmu.">
      <RouterLink :to="{ name: 'profile' }" class="btn-secondary mt-2">Buka profil</RouterLink>
    </StatusState>
    <ul v-else class="grid gap-4 md:grid-cols-2">
      <li v-for="invitation in visible" :key="invitation.id" class="card flex flex-col gap-4 p-5">
        <div class="flex items-start gap-3.5">
          <CompanyLogo :name="invitation.company.name" :url="invitation.company.logo_url" />
          <div class="min-w-0 flex-1">
            <div class="text-[14px] font-medium text-muted">{{ invitation.company.name }}</div>
            <RouterLink :to="{ name: 'job', params: { slug: invitation.job.slug } }" class="font-semibold hover:text-primary">{{ invitation.job.title }}</RouterLink>
            <div class="text-[14px] text-muted">{{ JOB_LEVEL[invitation.job.level] }} · {{ formatDate(invitation.created_at) }}</div>
          </div>
          <StatusBadge :status="INVITATION_STATUS[invitation.status]" />
        </div>
        <p v-if="invitation.message" class="rounded-btn bg-surface-2 px-3.5 py-2.5 text-sm">"{{ invitation.message }}"</p>
        <template v-if="invitation.status === 'pending'">
          <p class="flex items-center gap-1.5 text-[14px] text-muted"><Clock class="size-3.5" aria-hidden="true" />Berlaku sampai {{ formatDate(invitation.expires_at) }} ({{ deadlineInfo(invitation.expires_at.slice(0, 10)).days }} hari lagi)</p>
          <div class="flex gap-2">
            <RouterLink :to="{ name: 'job', params: { slug: invitation.job.slug }, query: { undangan: invitation.id } }" class="btn-primary flex-1">Terima & lamar</RouterLink>
            <button type="button" class="btn-ghost" @click="declining = invitation">Tolak</button>
          </div>
        </template>
      </li>
    </ul>
    <ConfirmDialog :open="!!declining" title="Tolak undangan ini?" :text="`${declining?.company.name} akan melihat bahwa kamu menolak undangan.`" confirm-label="Tolak undangan" :busy="busy" @confirm="decline" @cancel="declining = null" />
  </div>
</template>
