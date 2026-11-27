<script setup>
import { computed } from 'vue'
import { ArrowRight, BadgeCheck, Briefcase, Building2, Send, Users } from 'lucide-vue-next'
import BarList from '@/components/BarList.vue'
import ColumnChart from '@/components/ColumnChart.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'

const { data, loading, error, load } = useLoad(async () => (await api.get('/admin/dashboard')).data)
const weekFormat = new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short' })
const weekly = computed(() => data.value.weekly.map((week) => ({ label: weekFormat.format(new Date(`${week.week_start}T00:00:00`)), value: week.count })))
const totalUsers = computed(() => Object.values(data.value.users).reduce((sum, count) => sum + count, 0))
</script>

<template>
  <PageHeader title="Dashboard" text="Ringkasan aktivitas Karsa" />
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="3" />
  <template v-else>
    <div v-if="data.companies.pending" class="mb-5 flex flex-col gap-3 rounded-card border border-warning/25 bg-warning/[0.06] p-4 sm:flex-row sm:items-center">
      <span class="grid size-9 place-items-center rounded-full bg-warning/12 text-warning"><Building2 class="size-[18px]" aria-hidden="true" /></span>
      <p class="flex-1 font-semibold">{{ data.companies.pending }} perusahaan menunggu verifikasi</p>
      <RouterLink :to="{ name: 'admin-companies' }" class="btn-primary btn-sm">Tinjau sekarang <ArrowRight class="size-4" aria-hidden="true" /></RouterLink>
    </div>
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <StatCard label="Pengguna" :value="totalUsers" :icon="Users" :note="`${data.users.seeker} pencari kerja · ${data.users.company} perusahaan`" />
      <StatCard label="Perusahaan terverifikasi" :value="data.companies.verified" :icon="BadgeCheck" :note="`${data.companies.pending} menunggu · ${data.companies.rejected} ditolak`" />
      <StatCard label="Lowongan dibuka" :value="data.open_jobs" :icon="Briefcase" />
      <StatCard label="Total lamaran" :value="data.applications" :icon="Send" />
    </div>
    <div class="mt-4 grid gap-4 xl:grid-cols-[1.5fr_1fr]">
      <section class="card p-5 sm:p-6"><div class="flex items-center justify-between"><h2 class="text-xl">Lamaran per minggu</h2><span class="text-[14px] text-muted">8 minggu terakhir</span></div><ColumnChart :data="weekly" title="Lamaran per minggu" /></section>
      <section class="card p-5 sm:p-6"><h2 class="mb-4 text-xl">Lowongan per kategori</h2><BarList :items="data.jobs_by_category.map((item) => ({ label: item.name, value: item.count }))" /></section>
    </div>
    <section class="card mt-4 p-5 sm:p-6 xl:max-w-xl">
      <h2 class="text-xl">Perusahaan teratas</h2><p class="mb-4 text-[14px] text-muted">Berdasarkan jumlah lamaran</p>
      <BarList :items="data.top_companies.map((item) => ({ label: item.name, value: item.applications }))" />
    </section>
  </template>
</template>
