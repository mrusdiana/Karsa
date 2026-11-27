<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const keyword = ref(route.query.q ?? '')
const target = ref(null)
const busy = ref(false)
const roles = { seeker: 'Pencari kerja', company: 'Perusahaan', admin: 'Admin' }
const { data, loading, error, load } = useLoad(async () => (await api.get('/admin/users', { params: route.query })).data, null, { immediate: false })
watch(() => route.query, load, { immediate: true })

function update(changes) {
  const query = { ...route.query, ...changes }
  delete query.page
  for (const key of Object.keys(query)) if (!query[key]) delete query[key]
  router.push({ query })
}

async function toggle() {
  busy.value = true
  try {
    const updated = (await api.patch(`/admin/users/${target.value.id}`, { is_active: !target.value.is_active })).data
    data.value.data = data.value.data.map((user) => (user.id === updated.id ? updated : user))
    toast.show(updated.is_active ? 'Akun diaktifkan.' : 'Akun dinonaktifkan.')
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
    target.value = null
  }
}
</script>

<template>
  <PageHeader title="Pengguna" text="Cari pengguna dan nonaktifkan akun yang bermasalah." />
  <div class="mb-4 flex flex-wrap gap-2.5">
    <form class="flex-1" role="search" @submit.prevent="update({ q: keyword.trim() })">
      <label class="input flex max-w-md items-center gap-2.5"><Search class="size-[18px] text-muted" aria-hidden="true" /><span class="sr-only">Cari nama atau email</span><input v-model="keyword" class="w-full bg-transparent outline-none focus-visible:shadow-none" placeholder="Cari nama atau email" /></label>
    </form>
    <label class="sr-only" for="role-filter">Peran</label>
    <select id="role-filter" class="input w-auto" :value="route.query.role ?? ''" @change="update({ role: $event.target.value })"><option value="">Semua peran</option><option v-for="(label, key) in roles" :key="key" :value="key">{{ label }}</option></select>
  </div>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="4" height="h-14" />
  <StatusState v-else-if="!data.data.length" title="Tidak ada pengguna yang cocok" />
  <template v-else>
    <div class="card overflow-x-auto">
      <table class="w-full min-w-[720px] text-sm">
        <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Nama</th><th class="px-4 py-3 font-semibold">Peran</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 font-semibold">Daftar</th><th class="px-4 py-3"><span class="sr-only">Aksi</span></th></tr></thead>
        <tbody>
          <tr v-for="user in data.data" :key="user.id" class="border-t border-line">
            <td class="px-4 py-3"><div class="font-semibold">{{ user.name }}</div><div class="text-[13.5px] text-muted">{{ user.email }}<span v-if="user.company"> · {{ user.company.name }}</span></div></td>
            <td class="px-4 py-3">{{ roles[user.role] }}</td>
            <td class="px-4 py-3"><StatusBadge :status="user.is_active ? { label: 'Aktif', icon: 'CheckCircle2', tone: 'success' } : { label: 'Nonaktif', icon: 'Ban', tone: 'danger' }" /></td>
            <td class="px-4 py-3 text-[14px] text-muted">{{ formatDate(user.created_at) }}</td>
            <td class="px-4 py-3 text-right"><button v-if="user.id !== auth.user.id" type="button" class="btn-secondary btn-sm" :class="user.is_active && 'text-danger'" @click="target = user">{{ user.is_active ? 'Nonaktifkan' : 'Aktifkan' }}</button></td>
          </tr>
        </tbody>
      </table>
    </div>
    <PaginationBar :page="data" @change="(page) => router.push({ query: { ...route.query, page } })" />
  </template>
  <ConfirmDialog :open="!!target" :title="target?.is_active ? 'Nonaktifkan akun?' : 'Aktifkan akun?'" :text="target?.is_active ? `${target?.name} tidak bisa masuk lagi, dan lowongannya hilang dari pencarian.` : `${target?.name} bisa masuk kembali.`" :confirm-label="target?.is_active ? 'Nonaktifkan' : 'Aktifkan'" :danger="target?.is_active" :busy="busy" @confirm="toggle" @cancel="target = null" />
</template>
