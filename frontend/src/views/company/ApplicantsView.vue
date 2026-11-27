<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { AlertTriangle, ChevronRight, Mail } from 'lucide-vue-next'
import LoadingCards from '@/components/LoadingCards.vue'
import MatchScore from '@/components/MatchScore.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'
import { APPLICATION_STATUS, JOB_STATUS, formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const job = ref(null)
const skill = ref(route.query.skill ?? '')
const minimumScore = ref(Number(route.query.min_score ?? 0))
const { data, loading, error, load } = useLoad(async () => (await api.get(`/company/jobs/${route.params.id}/applications`, { params: route.query })).data, null, { immediate: false })
const total = computed(() => (data.value ? Object.values(data.value.counts).reduce((sum, count) => sum + count, 0) : 0))

function update(changes) {
  const query = { ...route.query, ...changes }
  delete query.page
  for (const key of Object.keys(query)) if (query[key] === '' || query[key] == null || query[key] === 0) delete query[key]
  router.push({ query })
}

watch(() => route.query, load, { immediate: true })
onMounted(async () => {
  job.value = (await api.get(`/company/jobs/${route.params.id}`)).data
})
</script>

<template>
  <nav class="mb-2 flex items-center gap-1.5 text-[14px] text-muted" aria-label="Breadcrumb"><RouterLink :to="{ name: 'company-jobs' }" class="hover:text-text">Lowongan</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" /><span class="text-text">{{ job?.title ?? 'Pelamar' }}</span></nav>
  <div class="flex flex-wrap items-center justify-between gap-3">
    <h1 class="text-2xl sm:text-[26px]">Pelamar · {{ job?.title }}</h1>
    <StatusBadge v-if="job" :status="JOB_STATUS[job.status]" />
  </div>

  <div class="mt-5 mb-4 flex gap-1 overflow-x-auto border-b border-line" role="tablist">
    <button v-for="option in [{ key: '', label: 'Semua', count: total }, ...Object.entries(APPLICATION_STATUS).map(([key, status]) => ({ key, label: status.label, count: data?.counts[key] ?? 0 }))]" :key="option.key" type="button" role="tab" :aria-selected="(route.query.status ?? '') === option.key" class="-mb-px flex min-h-11 items-center gap-2 border-b-2 px-3 text-sm font-medium whitespace-nowrap" :class="(route.query.status ?? '') === option.key ? 'border-primary text-primary' : 'border-transparent text-muted hover:text-text'" @click="update({ status: option.key })">
      {{ option.label }}<span class="num rounded-full px-2 text-xs font-semibold" :class="(route.query.status ?? '') === option.key ? 'bg-primary-soft' : 'bg-surface-2 text-text'">{{ option.count }}</span>
    </button>
  </div>

  <div class="mb-4 flex flex-wrap gap-2.5">
    <label class="input flex w-auto items-center gap-3"><span class="text-[14px] whitespace-nowrap">Skor min. <b class="num">{{ minimumScore }}</b></span><input v-model.number="minimumScore" type="range" min="0" max="100" step="5" class="w-24 accent-primary" @change="update({ min_score: minimumScore })" /></label>
    <form @submit.prevent="update({ skill: skill.trim() })"><label class="sr-only" for="skill-filter">Skill</label><input id="skill-filter" v-model="skill" class="input w-40" placeholder="Skill, mis. Vue" @blur="update({ skill: skill.trim() })" /></form>
    <label class="input flex w-auto cursor-pointer items-center gap-2 text-[14px]"><input type="checkbox" class="size-[18px] accent-primary" :checked="route.query.meets === 'true'" @change="update({ meets: $event.target.checked ? 'true' : '' })" />Hanya yang memenuhi syarat</label>
    <label class="sr-only" for="source-filter">Sumber</label>
    <select id="source-filter" class="input w-auto" :value="route.query.source ?? ''" @change="update({ source: $event.target.value })"><option value="">Semua sumber</option><option value="applied">Melamar</option><option value="invited">Diundang</option></select>
    <label class="sr-only" for="sort-filter">Urutkan</label>
    <select id="sort-filter" class="input w-auto sm:ml-auto" :value="route.query.sort ?? ''" @change="update({ sort: $event.target.value })"><option value="">Skor tertinggi</option><option value="newest">Terbaru</option></select>
  </div>

  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="4" height="h-16" />
  <StatusState v-else-if="!data.data.length" title="Belum ada pelamar yang cocok" :text="total ? 'Coba ubah filter.' : 'Bagikan lowongan atau undang kandidat lewat Cari Talenta.'">
    <RouterLink v-if="!total" :to="{ name: 'company-talents' }" class="btn-secondary mt-2">Cari Talenta</RouterLink>
  </StatusState>
  <template v-else>
    <div class="card overflow-x-auto">
      <table class="w-full min-w-[780px] text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Pelamar</th><th class="px-4 py-3 font-semibold">Skor</th><th class="px-4 py-3 font-semibold">Syarat</th><th class="px-4 py-3 font-semibold">Sumber</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 font-semibold">Dilamar</th></tr></thead>
        <tbody>
          <tr v-for="applicant in data.data" :key="applicant.id" class="relative border-t border-line transition hover:bg-surface-2/60">
            <td class="px-4 py-3">
              <div class="flex items-center gap-3">
                <UserAvatar :name="applicant.candidate.name" :url="applicant.candidate.avatar_url" />
                <div class="min-w-0">
                  <RouterLink :to="{ name: 'company-applicant', params: { id: applicant.id } }" class="font-semibold after:absolute after:inset-0">{{ applicant.candidate.name }}</RouterLink>
                  <div class="truncate text-[13.5px] text-muted">{{ [applicant.candidate.headline, applicant.candidate.location].filter(Boolean).join(' · ') }}</div>
                </div>
              </div>
            </td>
            <td class="px-4 py-3"><MatchScore :score="applicant.match_score" /></td>
            <td class="px-4 py-3"><span v-if="applicant.meets_requirements" class="text-[14px] text-muted">Memenuhi</span><span v-else class="inline-flex h-[26px] items-center gap-1.5 rounded-full bg-danger/10 px-2.5 text-[13.5px] font-semibold text-danger"><AlertTriangle class="size-3.5" aria-hidden="true" />Tidak memenuhi syarat</span></td>
            <td class="px-4 py-3"><span v-if="applicant.source === 'invited'" class="inline-flex h-[26px] items-center gap-1.5 rounded-full bg-info/10 px-2.5 text-[13.5px] font-semibold text-info"><Mail class="size-3.5" aria-hidden="true" />Diundang</span><span v-else class="text-[14px] text-muted">Melamar</span></td>
            <td class="px-4 py-3"><StatusBadge :status="APPLICATION_STATUS[applicant.status]" /></td>
            <td class="px-4 py-3 text-[14px] text-muted">{{ formatDate(applicant.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>
</template>
