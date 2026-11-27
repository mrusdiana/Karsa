<script setup>
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { INVITATION_STATUS, formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const { data, loading, error, load } = useLoad(async () => (await api.get('/company/invitations', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })
</script>

<template>
  <PageHeader title="Undangan Terkirim" :text="data ? `${data.remaining_today} dari 20 undangan tersisa hari ini` : ''">
    <RouterLink :to="{ name: 'company-talents' }" class="btn-primary">Cari Talenta</RouterLink>
  </PageHeader>
  <label class="sr-only" for="invitation-status">Status</label>
  <select id="invitation-status" class="input mb-4 w-auto" :value="route.query.status ?? ''" @change="router.push({ query: $event.target.value ? { status: $event.target.value } : {} })">
    <option value="">Semua status</option><option v-for="(status, key) in INVITATION_STATUS" :key="key" :value="key">{{ status.label }}</option>
  </select>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="3" height="h-16" />
  <StatusState v-else-if="!data.data.length" title="Belum ada undangan" text="Undang kandidat dari halaman Cari Talenta." />
  <template v-else>
    <div class="card overflow-x-auto">
      <table class="w-full min-w-[640px] text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Kandidat</th><th class="px-4 py-3 font-semibold">Lowongan</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 font-semibold">Dikirim</th></tr></thead>
        <tbody>
          <tr v-for="invitation in data.data" :key="invitation.id" class="border-t border-line">
            <td class="px-4 py-3.5"><div class="font-semibold">{{ invitation.candidate.name }}</div><div class="text-[13.5px] text-muted">{{ invitation.candidate.headline }}</div></td>
            <td class="px-4 py-3.5"><RouterLink :to="{ name: 'company-applicants', params: { id: invitation.job.id } }" class="hover:text-primary">{{ invitation.job.title }}</RouterLink></td>
            <td class="px-4 py-3.5"><StatusBadge :status="INVITATION_STATUS[invitation.status]" /></td>
            <td class="px-4 py-3.5 text-[14px] text-muted">{{ formatDate(invitation.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>
</template>
