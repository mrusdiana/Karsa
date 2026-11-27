<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronRight, FileText } from 'lucide-vue-next'
import ApplicationTimeline from '@/components/ApplicationTimeline.vue'
import CompanyLogo from '@/components/CompanyLogo.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, openFile } from '@/services/api'
import { useNotificationStore } from '@/stores/notifications'
import { useToastStore } from '@/stores/toast'
import { APPLICATION_STATUS, JOB_LEVEL, WORK_MODE, formatDate } from '@/utils/format'

const route = useRoute()
const toast = useToastStore()
const confirming = ref(false)
const busy = ref(false)
const notifications = useNotificationStore()
const { data: application, loading, error, load } = useLoad(async () => {
  const result = (await api.get(`/me/applications/${route.params.id}`)).data
  notifications.load()
  return result
})

async function withdraw() {
  busy.value = true
  try {
    application.value = (await api.post(`/me/applications/${route.params.id}/withdraw`)).data
    toast.show('Lamaran ditarik.')
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
    confirming.value = false
  }
}
</script>

<template>
  <div class="container-app py-8 pb-16">
    <nav class="mb-5 flex items-center gap-1.5 text-[14px] text-muted" aria-label="Breadcrumb"><RouterLink :to="{ name: 'applications' }" class="hover:text-text">Lamaran Saya</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" /><span class="text-text">Detail</span></nav>
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="2" height="h-48" />
    <div v-else class="grid items-start gap-6 lg:grid-cols-[1fr_360px]">
      <div class="flex flex-col gap-5">
        <section class="card flex flex-col gap-4 p-5 sm:flex-row sm:items-center sm:p-6">
          <CompanyLogo :name="application.job.company.name" :url="application.job.company.logo_url" :size="56" />
          <div class="min-w-0 flex-1">
            <h1 class="text-xl sm:text-2xl"><RouterLink :to="{ name: 'job', params: { slug: application.job.slug } }" class="hover:text-primary">{{ application.job.title }}</RouterLink></h1>
            <p class="text-sm text-muted">{{ application.job.company.name }} · {{ JOB_LEVEL[application.job.level] }} · {{ WORK_MODE[application.job.work_mode] }}</p>
            <p class="mt-1 text-[14px] text-muted">Dilamar {{ formatDate(application.created_at) }}<span v-if="application.source === 'invited'"> · lewat undangan</span></p>
          </div>
          <StatusBadge :status="APPLICATION_STATUS[application.status]" />
        </section>
        <section class="card p-5 sm:p-6">
          <h2 class="text-lg">Yang kamu kirim</h2>
          <button type="button" class="mt-4 flex w-full items-center gap-3 rounded-btn border border-line p-3 text-left text-sm hover:border-line-strong" @click="openFile(application.cv_url)">
            <FileText class="size-6 text-danger" aria-hidden="true" /><span class="flex-1 font-medium">CV saat melamar</span><span class="link text-[14px]">Lihat</span>
          </button>
          <div v-if="application.cover_letter" class="mt-5"><h3 class="label">Surat pengantar</h3><p class="text-sm whitespace-pre-line">{{ application.cover_letter }}</p></div>
          <div v-if="application.answers.length" class="mt-5">
            <h3 class="label">Jawaban screening</h3>
            <dl class="divide-y divide-line text-sm">
              <div v-for="answer in application.answers" :key="answer.question" class="flex justify-between gap-4 py-2.5"><dt>{{ answer.question }}</dt><dd class="font-medium">{{ answer.answer === 'yes' ? 'Ya' : answer.answer === 'no' ? 'Tidak' : answer.answer || '–' }}</dd></div>
            </dl>
          </div>
        </section>
      </div>
      <aside class="card p-5 sm:p-6 lg:sticky lg:top-24">
        <h2 class="mb-4 text-lg">Timeline</h2>
        <ApplicationTimeline :events="application.events" />
        <button v-if="application.status === 'applied'" type="button" class="btn-secondary mt-6 w-full text-danger" @click="confirming = true">Tarik lamaran</button>
      </aside>
    </div>
    <ConfirmDialog :open="confirming" title="Tarik lamaran ini?" text="Lamaran yang sudah ditarik tidak bisa diajukan ulang untuk lowongan yang sama." confirm-label="Tarik lamaran" :busy="busy" @confirm="withdraw" @cancel="confirming = false" />
  </div>
</template>
