<script setup>
import { useRoute, useRouter } from 'vue-router'
import { watch } from 'vue'
import JobCard from '@/components/JobCard.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()
const { data, loading, error, load } = useLoad(async () => (await api.get('/me/saved-jobs', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })
const remove = (id) => {
  data.value.data = data.value.data.filter((job) => job.id !== id)
  data.value.total -= 1
}
</script>

<template>
  <div class="container-app py-8 pb-16">
    <PageHeader title="Lowongan Tersimpan" text="Lowongan yang kamu tandai untuk dilamar nanti." />
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="3" />
    <StatusState v-else-if="!data.data.length" title="Belum ada lowongan tersimpan" text="Tekan ikon bookmark di kartu lowongan untuk menyimpannya.">
      <RouterLink :to="{ name: 'jobs' }" class="btn-primary mt-2">Cari lowongan</RouterLink>
    </StatusState>
    <template v-else>
      <div class="grid gap-4 md:grid-cols-2"><JobCard v-for="job in data.data" :key="job.id" :job="job" :closed="!job.is_open" @unsaved="remove" /></div>
      <PaginationBar :page="data" @change="(page) => router.push({ query: { page } })" />
    </template>
  </div>
</template>
