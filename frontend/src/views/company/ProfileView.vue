<script setup>
import { reactive, ref } from 'vue'
import { ImageUp } from 'lucide-vue-next'
import CompanyLogo from '@/components/CompanyLogo.vue'
import FormField from '@/components/FormField.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import StatusState from '@/components/StatusState.vue'
import VerificationBanner from '@/components/VerificationBanner.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, fieldErrors } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { COMPANY_SIZE, COMPANY_STATUS } from '@/utils/format'

const auth = useAuthStore()
const toast = useToastStore()
const form = reactive({ name: '', industry: '', size: '', location: '', website: '', description: '' })
const errors = ref({})
const saving = ref(false)
const uploading = ref(false)
const logoError = ref('')
const fileInput = ref()

function fill(company) {
  for (const key of Object.keys(form)) form[key] = company[key] ?? ''
  return company
}

const { data: company, loading, error, load } = useLoad(async () => fill((await api.get('/company/profile')).data))

async function save() {
  saving.value = true
  errors.value = {}
  const wasRejected = company.value.status === 'rejected'
  try {
    company.value = fill((await api.put('/company/profile', form)).data)
    await auth.refresh()
    toast.show(wasRejected ? 'Profil disimpan dan diajukan ulang untuk verifikasi.' : 'Profil perusahaan disimpan.')
  } catch (failure) {
    errors.value = fieldErrors(failure)
    if (!Object.keys(errors.value).length) toast.show(errorMessage(failure), 'danger')
  } finally {
    saving.value = false
  }
}

async function uploadLogo(event) {
  const file = event.target.files[0]
  event.target.value = ''
  logoError.value = ''
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) return (logoError.value = 'Logo harus JPG, PNG, atau WebP.')
  if (file.size > 1024 * 1024) return (logoError.value = 'Ukuran logo maksimal 1 MB.')
  uploading.value = true
  try {
    const body = new FormData()
    body.append('logo', file)
    company.value = (await api.post('/company/logo', body)).data
    toast.show('Logo diperbarui.')
  } catch (failure) {
    logoError.value = fieldErrors(failure).logo ?? errorMessage(failure)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <PageHeader title="Profil Perusahaan" text="Profil yang lengkap membantu admin memverifikasi perusahaanmu." />
  <StatusState v-if="error" type="error" @retry="load" />
  <LoadingCards v-else-if="loading" :count="2" height="h-56" />
  <template v-else>
    <VerificationBanner :company="company" />
    <div class="grid items-start gap-6 lg:grid-cols-[1fr_300px]">
      <form class="card flex flex-col gap-5 p-5 sm:p-7" novalidate @submit.prevent="save">
        <FormField id="company-name" label="Nama perusahaan" :error="errors.name" required v-slot="{ describedBy }"><input id="company-name" v-model="form.name" class="input" maxlength="150" :aria-describedby="describedBy" /></FormField>
        <div class="grid gap-5 sm:grid-cols-2">
          <FormField id="industry" label="Industri" :error="errors.industry" v-slot="{ describedBy }"><input id="industry" v-model="form.industry" class="input" maxlength="100" placeholder="Contoh: E-commerce" :aria-describedby="describedBy" /></FormField>
          <FormField id="size" label="Ukuran" :error="errors.size" v-slot="{ describedBy }">
            <select id="size" v-model="form.size" class="input" :aria-describedby="describedBy"><option value="">Pilih</option><option v-for="size in COMPANY_SIZE" :key="size" :value="size">{{ size }} karyawan</option></select>
          </FormField>
          <FormField id="company-location" label="Lokasi" :error="errors.location" v-slot="{ describedBy }"><input id="company-location" v-model="form.location" class="input" maxlength="100" :aria-describedby="describedBy" /></FormField>
          <FormField id="website" label="Website" :error="errors.website" v-slot="{ describedBy }"><input id="website" v-model="form.website" type="url" class="input" placeholder="https://" :aria-describedby="describedBy" /></FormField>
        </div>
        <FormField id="description" label="Deskripsi" :error="errors.description" v-slot="{ describedBy }"><textarea id="description" v-model="form.description" class="input" rows="5" maxlength="5000" :aria-describedby="describedBy" /></FormField>
        <div class="flex justify-end"><button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Menyimpan…' : company.status === 'rejected' ? 'Simpan & ajukan ulang' : 'Simpan profil' }}</button></div>
      </form>
      <aside class="card flex flex-col items-center gap-4 p-6 text-center">
        <CompanyLogo :name="company.name" :url="company.logo_url" :size="96" />
        <StatusBadge :status="COMPANY_STATUS[company.status]" />
        <button type="button" class="btn-secondary w-full" :disabled="uploading" @click="fileInput.click()"><ImageUp class="size-4" aria-hidden="true" />{{ uploading ? 'Mengunggah…' : 'Ganti logo' }}</button>
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" class="sr-only" aria-label="Pilih file logo" @change="uploadLogo" />
        <p class="text-xs text-muted">JPG, PNG, atau WebP, maksimal 1 MB</p>
        <p v-if="logoError" class="text-[14px] text-danger" role="alert">{{ logoError }}</p>
      </aside>
    </div>
  </template>
</template>
