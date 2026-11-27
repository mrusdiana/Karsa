<script setup>
import { BadgeCheck, ExternalLink } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import CompanyLogo from '@/components/CompanyLogo.vue'
import JobCard from '@/components/JobCard.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'

const route = useRoute()
const { data: company, loading, error, load } = useLoad(async () => (await api.get(`/companies/${route.params.slug}`)).data)
</script>

<template>
  <StatusState v-if="error" type="error" title="Perusahaan tidak ditemukan" class="container-app my-10" @retry="load" />
  <LoadingCards v-else-if="loading" class="container-app my-10" :count="2" />
  <template v-else>
    <div class="h-28 border-b border-line bg-primary-soft sm:h-36" />
    <div class="container-app pb-16">
      <div class="-mt-12 flex flex-col gap-4 sm:flex-row sm:items-end sm:gap-6">
        <CompanyLogo :name="company.name" :url="company.logo_url" :size="96" class="border-4 border-bg" />
        <div class="pb-1">
          <h1 class="flex items-center gap-2 text-2xl sm:text-[28px]">{{ company.name }}<BadgeCheck class="size-6 text-primary" aria-label="Terverifikasi" /></h1>
          <p class="text-sm text-muted">{{ [company.industry, company.size && `${company.size} karyawan`, company.location].filter(Boolean).join(' · ') }}</p>
        </div>
        <a v-if="company.website" :href="company.website" target="_blank" rel="noopener noreferrer" class="btn-secondary sm:ml-auto"><ExternalLink class="size-4" aria-hidden="true" />Website</a>
      </div>
      <p v-if="company.description" class="mt-6 max-w-[68ch]">{{ company.description }}</p>
      <h2 class="mt-10 mb-4 text-xl">{{ company.jobs.length }} lowongan dibuka</h2>
      <div v-if="company.jobs.length" class="grid gap-4 md:grid-cols-2"><JobCard v-for="job in company.jobs" :key="job.id" :job="job" /></div>
      <StatusState v-else title="Belum ada lowongan dibuka" />
    </div>
  </template>
</template>
