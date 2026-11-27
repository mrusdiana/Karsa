<script setup>
import { ref } from 'vue'
import { Check, Pencil, Trash2, X } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, fieldErrors } from '@/services/api'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
const name = ref('')
const createError = ref('')
const editing = ref(null)
const editName = ref('')
const editError = ref('')
const deleting = ref(null)
const busy = ref(false)
const { data, loading, error, load } = useLoad(async () => (await api.get('/admin/categories')).data, [])
const slugify = (text) => text.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')

async function create() {
  createError.value = ''
  try {
    await api.post('/admin/categories', { name: name.value })
    toast.show('Kategori ditambahkan.')
    name.value = ''
    load()
  } catch (failure) {
    createError.value = fieldErrors(failure).name ?? errorMessage(failure)
  }
}

async function saveEdit() {
  editError.value = ''
  try {
    await api.put(`/admin/categories/${editing.value}`, { name: editName.value })
    toast.show('Kategori diperbarui.')
    editing.value = null
    load()
  } catch (failure) {
    editError.value = fieldErrors(failure).name ?? errorMessage(failure)
  }
}

async function remove() {
  busy.value = true
  try {
    await api.delete(`/admin/categories/${deleting.value.id}`)
    toast.show('Kategori dihapus.')
    load()
  } catch (failure) {
    toast.show(errorMessage(failure), 'danger')
  } finally {
    busy.value = false
    deleting.value = null
  }
}
</script>

<template>
  <PageHeader title="Kategori" text="Kategori lowongan yang tampil di filter dan beranda." />
  <form class="card mb-4 flex flex-col gap-2 p-4 sm:flex-row sm:items-start" @submit.prevent="create">
    <div class="flex-1">
      <label for="new-category" class="sr-only">Nama kategori baru</label>
      <input id="new-category" v-model="name" class="input" maxlength="60" placeholder="Nama kategori baru" :aria-describedby="createError ? 'new-category-error' : 'new-category-slug'" />
      <p v-if="createError" id="new-category-error" class="mt-1 text-[14px] text-danger">{{ createError }}</p>
      <p v-else-if="name" id="new-category-slug" class="mt-1 text-[14px] text-muted">Slug: {{ slugify(name) }}</p>
    </div>
    <button type="submit" class="btn-primary" :disabled="!name.trim()">Tambah</button>
  </form>
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="4" height="h-14" />
  <div v-else class="card overflow-x-auto">
    <table class="w-full min-w-[560px] text-sm">
      <thead><tr class="bg-surface-2 text-left text-[13.5px] text-muted"><th class="px-4 py-3 font-semibold">Nama</th><th class="px-4 py-3 font-semibold">Slug</th><th class="px-4 py-3 text-right font-semibold">Lowongan dibuka</th><th class="px-4 py-3"><span class="sr-only">Aksi</span></th></tr></thead>
      <tbody>
        <tr v-for="category in data" :key="category.id" class="border-t border-line">
          <td class="px-4 py-2.5">
            <form v-if="editing === category.id" class="flex flex-col gap-1" @submit.prevent="saveEdit">
              <div class="flex gap-1"><label :for="`edit-${category.id}`" class="sr-only">Nama kategori</label><input :id="`edit-${category.id}`" v-model="editName" class="input min-h-9" maxlength="60" /><button type="submit" class="btn-icon size-9" aria-label="Simpan"><Check class="size-4" aria-hidden="true" /></button><button type="button" class="btn-icon size-9" aria-label="Batal" @click="editing = null"><X class="size-4" aria-hidden="true" /></button></div>
              <p v-if="editError" class="text-[14px] text-danger">{{ editError }}</p>
            </form>
            <span v-else class="font-medium">{{ category.name }}</span>
          </td>
          <td class="px-4 py-2.5 text-muted">{{ category.slug }}</td>
          <td class="num px-4 py-2.5 text-right">{{ category.job_count }}</td>
          <td class="px-4 py-2.5">
            <div class="flex justify-end gap-0.5">
              <button type="button" class="btn-icon size-9" :aria-label="`Edit ${category.name}`" @click="editing = category.id; editName = category.name; editError = ''"><Pencil class="size-[18px]" aria-hidden="true" /></button>
              <button type="button" class="btn-icon size-9 hover:text-danger" :aria-label="`Hapus ${category.name}`" @click="deleting = category"><Trash2 class="size-[18px]" aria-hidden="true" /></button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <ConfirmDialog :open="!!deleting" title="Hapus kategori?" :text="`${deleting?.name} akan dihapus. Kategori yang masih dipakai lowongan tidak bisa dihapus.`" :busy="busy" @confirm="remove" @cancel="deleting = null" />
</template>
