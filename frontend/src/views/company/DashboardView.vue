<script setup>
import { computed, onMounted, ref } from 'vue'
import { Briefcase, Mail, Plus, UserPlus, Users } from 'lucide-vue-next'
import BarList from '@/components/BarList.vue'
import ColumnChart from '@/components/ColumnChart.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
import StatusState from '@/components/StatusState.vue'
import VerificationBanner from '@/components/VerificationBanner.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { APPLICATION_STATUS } from '@/utils/format'

const auth = useAuthStore()
const company = ref(null)
const { data, loading, error, load } = useLoad(async () => (await api.get('/company/dashboard')).data)
const weekFormat = new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short' })
const weekly = computed(() => data.value.weekly.map((week) => ({ label: weekFormat.format(new Date(`${week.week_start}T00:00:00`)), value: week.count })))
const byStatus = computed(() => Object.entries(data.value.by_status).map(([status, count]) => ({ label: APPLICATION_STATUS[status].label, value: count })))

onMounted(async () => {
  company.value = (await api.get('/company/profile')).data
  auth.refresh()
})
</script>

<template>
  <PageHeader title="Dashboard" :text="auth.company?.name">
    <RouterLink :to="{ name: 'company-job-new' }" class="btn-primary"><Plus class="size-4" aria-hidden="true" />Buat lowongan</RouterLink>
  </PageHeader>
  <VerificationBanner v-if="company" :company="company" />
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="3" />
  <template v-else>
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <StatCard label="Lowongan dibuka" :value="data.open_jobs" :icon="Briefcase" />
      <StatCard label="Total pelamar" :value="data.total_applicants" :icon="Users" />
      <StatCard label="Pelamar baru 7 hari" :value="data.new_applicants_7d" :icon="UserPlus" :note="data.new_applicants_7d ? 'Perlu ditinjau' : ''" note-tone="success" />
      <StatCard label="Undangan menunggu" :value="data.pending_invitations" :icon="Mail" />
    </div>
    <div class="mt-4 grid gap-4 xl:grid-cols-[1.5fr_1fr]">
      <section class="card p-5 sm:p-6"><div class="flex items-center justify-between"><h2 class="text-xl">Lamaran per minggu</h2><span class="text-[14px] text-muted">8 minggu terakhir</span></div><ColumnChart :data="weekly" title="Lamaran per minggu" /></section>
      <section class="card p-5 sm:p-6"><h2 class="mb-4 text-xl">Pelamar per status</h2><BarList :items="byStatus" /></section>
    </div>
    <section class="card mt-4 overflow-hidden">
      <h2 class="px-5 pt-5 pb-3 text-xl sm:px-6">Lowongan dengan pelamar terbanyak</h2>
      <table v-if="data.top_jobs.length" class="w-full text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-5 py-3 font-semibold sm:px-6">Lowongan</th><th class="px-5 py-3 text-right font-semibold sm:px-6">Pelamar</th></tr></thead>
        <tbody>
          <tr v-for="job in data.top_jobs" :key="job.id" class="border-t border-line">
            <td class="px-5 py-3.5 sm:px-6"><RouterLink :to="{ name: 'company-applicants', params: { id: job.id } }" class="font-medium hover:text-primary">{{ job.title }}</RouterLink></td>
            <td class="num px-5 py-3.5 text-right sm:px-6">{{ job.applicants }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="px-5 pb-5 text-sm text-muted sm:px-6">Belum ada pelamar.</p>
    </section>
  </template>
</template>
