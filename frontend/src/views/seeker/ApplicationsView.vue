<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChevronRight } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { APPLICATION_STATUS, formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const tabs = [
  { key: '', label: 'Semua' },
  { key: 'applied,reviewed,interview', label: 'Aktif' },
  { key: 'accepted', label: 'Diterima' },
  { key: 'rejected,withdrawn', label: 'Ditolak/Ditarik' },
]
const activeTab = computed(() => route.query.status ?? '')
const { data, loading, error, load } = useLoad(async () => (await api.get('/me/applications', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })
</script>

<template>
  <div class="container-app py-8 pb-16">
    <PageHeader title="Lamaran Saya" text="Pantau status setiap lamaranmu." />
    <div class="mb-5 flex gap-1 overflow-x-auto border-b border-line" role="tablist">
      <button v-for="tab in tabs" :key="tab.key" type="button" role="tab" :aria-selected="activeTab === tab.key" class="-mb-px min-h-11 border-b-2 px-3.5 text-sm font-medium whitespace-nowrap" :class="activeTab === tab.key ? 'border-primary text-primary' : 'border-transparent text-muted hover:text-text'" @click="router.push({ query: tab.key ? { status: tab.key } : {} })">{{ tab.label }}</button>
    </div>
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="3" height="h-20" />
    <StatusState v-else-if="!data.data.length" title="Belum ada lamaran di sini" text="Lamaran yang kamu kirim akan muncul di halaman ini.">
      <RouterLink :to="{ name: 'jobs' }" class="btn-primary mt-2">Cari lowongan</RouterLink>
    </StatusState>
    <template v-else>
      <ul class="flex flex-col gap-3">
        <li v-for="application in data.data" :key="application.id">
          <RouterLink :to="{ name: 'application', params: { id: application.id } }" class="card flex items-center gap-4 p-4 transition hover:border-line-strong hover:shadow-card">
            <CompanyLogo :name="application.job.company.name" :url="application.job.company.logo_url" />
            <div class="min-w-0 flex-1">
              <div class="truncate font-semibold">{{ application.job.title }}</div>
              <div class="text-[14px] text-muted">{{ application.job.company.name }} · Dilamar {{ formatDate(application.created_at) }}</div>
            </div>
            <div class="flex flex-col items-end gap-1.5 sm:flex-row sm:items-center sm:gap-3">
              <span v-if="application.has_update" class="rounded-full bg-danger px-2 py-0.5 text-xs font-semibold text-white">Status baru</span>
              <span v-if="application.source === 'invited'" class="hidden text-xs font-semibold text-info sm:inline">Diundang</span>
              <StatusBadge :status="APPLICATION_STATUS[application.status]" />
            </div>
            <ChevronRight class="hidden size-5 text-muted sm:block" aria-hidden="true" />
          </RouterLink>
        </li>
      </ul>
      <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
    </template>
  </div>
</template>
