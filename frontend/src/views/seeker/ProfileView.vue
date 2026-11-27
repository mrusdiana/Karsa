<script setup>
import { computed, reactive, ref } from 'vue'
import { Camera, Check, Eye, FileText, Lock, LogOut, UploadCloud } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useSignOut } from '@/composables/useSignOut'
import { useAuthStore } from '@/stores/auth'
import FormField from '@/components/FormField.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import PageHeader from '@/components/PageHeader.vue'
import SkillInput from '@/components/SkillInput.vue'
import StatusState from '@/components/StatusState.vue'
import { useLoad } from '@/composables/useLoad'
import api, { errorMessage, fieldErrors, openFile } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import { EDUCATION, MISSING_LABEL, formatDate } from '@/utils/format'

const toast = useToastStore()
const auth = useAuthStore()
const { confirming, signOut } = useSignOut()
const avatarInput = ref()
const avatarUploading = ref(false)
const avatarError = ref('')
const form = reactive({ headline: '', phone: '', location: '', education: '', major: '', graduation_year: null, experience_years: 0, skills: [], bio: '', portfolio_url: '', is_discoverable: false })
const errors = ref({})
const saving = ref(false)
const uploading = ref(false)
const cvError = ref('')
const fileInput = ref()

function fill(profile) {
  for (const key of Object.keys(form)) form[key] = profile[key] ?? (key === 'skills' ? [] : key === 'experience_years' ? 0 : key === 'graduation_year' ? null : '')
  return profile
}

const { data: profile, loading, error, load } = useLoad(async () => fill((await api.get('/me/profile')).data))
const completeness = computed(() => profile.value?.completeness ?? { percent: 0, missing: [] })
const canDiscover = computed(() => completeness.value.missing.length === 0)

async function save(changes = {}) {
  saving.value = true
  errors.value = {}
  try {
    const body = { ...form, ...changes, graduation_year: form.graduation_year ? Number(form.graduation_year) : null, experience_years: Number(form.experience_years) || 0 }
    profile.value = fill((await api.put('/me/profile', body)).data)
    toast.show('Profil disimpan.')
  } catch (failure) {
    errors.value = fieldErrors(failure)
    if (errors.value.is_discoverable) toast.show(errors.value.is_discoverable, 'danger')
    else toast.show(Object.keys(errors.value).length ? 'Periksa lagi isian yang ditandai.' : errorMessage(failure), 'danger')
  } finally {
    saving.value = false
  }
}

async function uploadAvatar(event) {
  const file = event.target.files[0]
  event.target.value = ''
  avatarError.value = ''
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) return (avatarError.value = 'Foto harus JPG, PNG, atau WebP.')
  if (file.size > 2 * 1024 * 1024) return (avatarError.value = 'Ukuran foto maksimal 2 MB.')
  avatarUploading.value = true
  try {
    const body = new FormData()
    body.append('avatar', file)
    auth.user = (await api.post('/me/avatar', body)).data.user
    toast.show('Foto profil diperbarui.')
  } catch (failure) {
    avatarError.value = fieldErrors(failure).avatar ?? errorMessage(failure)
  } finally {
    avatarUploading.value = false
  }
}

async function uploadCV(event) {
  const file = event.target.files[0]
  event.target.value = ''
  cvError.value = ''
  if (!file) return
  if (file.type !== 'application/pdf') return (cvError.value = 'CV harus berupa file PDF.')
  if (file.size > 2 * 1024 * 1024) return (cvError.value = 'Ukuran CV maksimal 2 MB.')
  uploading.value = true
  try {
    const body = new FormData()
    body.append('cv', file)
    profile.value = fill((await api.post('/me/cv', body)).data)
    toast.show('CV berhasil diunggah.')
  } catch (failure) {
    cvError.value = fieldErrors(failure).cv ?? errorMessage(failure)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="container-app py-8 pb-16">
    <PageHeader title="Profil & CV" text="Profil lengkap dipakai untuk melamar dan menghitung skor kecocokanmu." />
    <StatusState v-if="error" type="error" @retry="load" />
    <LoadingCards v-else-if="loading" :count="2" height="h-64" />
    <div v-else class="grid items-start gap-6 lg:grid-cols-[1fr_340px]">
      <form class="card flex flex-col gap-5 p-5 sm:p-7" novalidate @submit.prevent="save()">
        <h2 class="text-lg">Data diri</h2>
        <FormField id="headline" label="Headline" :error="errors.headline" hint="Contoh: Mahasiswa Informatika, fokus frontend" required v-slot="{ describedBy }">
          <input id="headline" v-model="form.headline" class="input" maxlength="150" :aria-describedby="describedBy" />
        </FormField>
        <div class="grid gap-5 sm:grid-cols-2">
          <FormField id="location" label="Lokasi (kota)" :error="errors.location" required v-slot="{ describedBy }"><input id="location" v-model="form.location" class="input" maxlength="100" :aria-describedby="describedBy" /></FormField>
          <FormField id="phone" label="Nomor HP" :error="errors.phone" hint="Hanya terlihat oleh perusahaan yang kamu lamar" v-slot="{ describedBy }"><input id="phone" v-model="form.phone" class="input" inputmode="tel" maxlength="20" :aria-describedby="describedBy" /></FormField>
          <FormField id="education" label="Pendidikan terakhir" :error="errors.education" required v-slot="{ describedBy }">
            <select id="education" v-model="form.education" class="input" :aria-describedby="describedBy"><option value="">Pilih</option><option v-for="level in EDUCATION" :key="level" :value="level">{{ level }}</option></select>
          </FormField>
          <FormField id="major" label="Jurusan" :error="errors.major" v-slot="{ describedBy }"><input id="major" v-model="form.major" class="input" maxlength="100" :aria-describedby="describedBy" /></FormField>
          <FormField id="graduation-year" label="Tahun lulus" :error="errors.graduation_year" hint="Isi perkiraan kalau masih kuliah" required v-slot="{ describedBy }"><input id="graduation-year" v-model="form.graduation_year" type="number" min="1970" max="2100" class="input" :aria-describedby="describedBy" /></FormField>
          <FormField id="experience" label="Pengalaman kerja (tahun)" :error="errors.experience_years" v-slot="{ describedBy }"><input id="experience" v-model="form.experience_years" type="number" min="0" max="30" class="input" :aria-describedby="describedBy" /></FormField>
        </div>
        <FormField id="skills" label="Skill" :error="errors.skills" hint="1–15 skill, tekan Enter setelah mengetik" required v-slot="{ describedBy, invalid }">
          <SkillInput id="skills" v-model="form.skills" :max="15" :described-by="describedBy" :invalid="invalid" />
        </FormField>
        <FormField id="bio" label="Bio singkat" :error="errors.bio" v-slot="{ describedBy }"><textarea id="bio" v-model="form.bio" class="input" rows="4" maxlength="2000" :aria-describedby="describedBy" /></FormField>
        <FormField id="portfolio" label="Tautan portofolio" :error="errors.portfolio_url" v-slot="{ describedBy }"><input id="portfolio" v-model="form.portfolio_url" type="url" class="input" placeholder="https://" :aria-describedby="describedBy" /></FormField>
        <div class="flex justify-end"><button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Menyimpan…' : 'Simpan profil' }}</button></div>
      </form>

      <div class="order-first flex flex-col gap-5 lg:sticky lg:top-24 lg:order-none">
        <section class="card flex items-center gap-4 p-5">
          <div class="relative">
            <UserAvatar :name="auth.user.name" :url="auth.user.avatar_url" :size="72" />
            <button type="button" class="absolute -right-1 -bottom-1 grid size-8 place-items-center rounded-full border-2 border-surface bg-primary text-white hover:bg-primary-hover" :disabled="avatarUploading" aria-label="Ganti foto profil" @click="avatarInput.click()"><Camera class="size-4" aria-hidden="true" /></button>
          </div>
          <div class="min-w-0">
            <div class="truncate font-semibold">{{ auth.user.name }}</div>
            <div class="truncate text-sm text-muted">{{ auth.user.email }}</div>
            <button type="button" class="link mt-1 text-sm" :disabled="avatarUploading" @click="avatarInput.click()">{{ avatarUploading ? 'Mengunggah…' : 'Ganti foto' }}</button>
          </div>
          <input ref="avatarInput" type="file" accept="image/jpeg,image/png,image/webp" class="sr-only" aria-label="Pilih foto profil" @change="uploadAvatar" />
        </section>
        <p v-if="avatarError" class="-mt-3 text-sm text-danger" role="alert">{{ avatarError }}</p>
        <section class="card p-5">
          <div class="flex items-baseline justify-between"><h2 class="text-base">Kelengkapan profil</h2><b class="num text-primary">{{ completeness.percent }}%</b></div>
          <div class="mt-2.5 h-2 overflow-hidden rounded-full bg-surface-2"><div class="h-full rounded-full bg-primary transition-all" :style="{ width: `${completeness.percent}%` }" /></div>
          <p v-if="canDiscover" class="mt-3 flex items-center gap-1.5 text-[14px] text-success"><Check class="size-4" aria-hidden="true" />Profil lengkap, kamu bisa melamar.</p>
          <ul v-else class="mt-3 flex flex-col gap-1 text-[14px] text-muted">
            <li v-for="item in completeness.missing" :key="item" class="flex items-center gap-2"><span class="size-1.5 rounded-full bg-warning" />{{ MISSING_LABEL[item] }}</li>
          </ul>
        </section>

        <section class="card p-5">
          <h2 class="text-base">CV</h2>
          <div v-if="profile.cv" class="mt-3 flex items-center gap-3 rounded-btn border border-line p-3">
            <FileText class="size-8 shrink-0 text-danger" aria-hidden="true" />
            <div class="min-w-0 flex-1"><div class="truncate text-sm font-medium">{{ profile.cv.name }}</div><div class="text-xs text-muted">{{ Math.ceil(profile.cv.size / 1024) }} KB · {{ formatDate(profile.cv.uploaded_at) }}</div></div>
            <button type="button" class="btn-icon size-9" aria-label="Lihat CV" @click="openFile(profile.cv.url)"><Eye class="size-[18px]" aria-hidden="true" /></button>
          </div>
          <button type="button" class="mt-3 flex w-full flex-col items-center gap-1.5 rounded-btn border-2 border-dashed border-line px-4 py-5 text-sm transition hover:border-primary hover:bg-primary-soft/40" :disabled="uploading" @click="fileInput.click()">
            <UploadCloud class="size-6 text-primary" aria-hidden="true" />
            <span class="font-semibold">{{ uploading ? 'Mengunggah…' : profile.cv ? 'Ganti CV' : 'Unggah CV' }}</span>
            <span class="text-xs text-muted">Hanya PDF, maksimal 2 MB</span>
          </button>
          <input ref="fileInput" type="file" accept="application/pdf" class="sr-only" aria-label="Pilih file CV" @change="uploadCV" />
          <p v-if="cvError" class="mt-2 text-[14px] text-danger" role="alert">{{ cvError }}</p>
        </section>

        <section class="card p-5">
          <label class="flex cursor-pointer items-start justify-between gap-4" :class="!canDiscover && 'cursor-not-allowed opacity-60'">
            <span><span class="block text-base font-semibold">Bisa ditemukan perusahaan</span><span class="mt-1 block text-[14px] text-muted">Perusahaan terverifikasi bisa melihat nama, headline, skill, lokasi, dan pendidikanmu.</span></span>
            <input type="checkbox" role="switch" class="peer sr-only" :checked="form.is_discoverable" :disabled="!canDiscover || saving" @change="save({ is_discoverable: $event.target.checked })" />
            <span class="relative mt-1 h-6 w-11 shrink-0 rounded-full bg-line-strong transition peer-checked:bg-primary peer-focus-visible:shadow-[0_0_0_3px_rgb(15_118_110/0.35)] after:absolute after:top-0.5 after:left-0.5 after:size-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
          </label>
          <p class="mt-3 flex items-center gap-1.5 text-[14px] text-muted"><Lock class="size-3.5" aria-hidden="true" />Email, nomor HP, dan CV tetap tersembunyi.</p>
          <p v-if="!canDiscover" class="mt-1.5 text-[14px] text-warning">Lengkapi profil dan CV untuk mengaktifkan ini.</p>
        </section>
        <button type="button" class="btn-secondary text-danger md:hidden" @click="confirming = true"><LogOut class="size-4" aria-hidden="true" />Keluar</button>
      </div>
    </div>
    <ConfirmDialog :open="confirming" title="Keluar dari Karsa?" text="Kamu perlu masuk lagi untuk melamar dan melihat lamaranmu." confirm-label="Keluar" @confirm="signOut" @cancel="confirming = false" />
  </div>
</template>
