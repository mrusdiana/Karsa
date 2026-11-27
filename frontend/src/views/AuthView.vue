<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Building2, Eye, EyeOff, GraduationCap } from 'lucide-vue-next'
import AppLogo from '@/components/AppLogo.vue'
import FormField from '@/components/FormField.vue'
import { errorMessage, fieldErrors } from '@/services/api'
import { HOME_BY_ROLE, useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const props = defineProps({ mode: { type: String, required: true } })
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const isLogin = computed(() => props.mode === 'login')
const form = reactive({ name: '', email: '', password: '', role: route.query.role === 'company' ? 'company' : '', company_name: '' })
const errors = ref({})
const message = ref('')
const busy = ref(false)
const showPassword = ref(false)

watch(() => props.mode, () => {
  errors.value = {}
  message.value = ''
})

async function submit() {
  busy.value = true
  errors.value = {}
  message.value = ''
  try {
    if (isLogin.value) await auth.login({ email: form.email, password: form.password })
    else await auth.register(form)
    toast.show(isLogin.value ? `Selamat datang, ${auth.user.name.split(' ')[0]}!` : 'Akun berhasil dibuat.')
    const redirect = typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/') ? route.query.redirect : null
    router.push(redirect && auth.role === 'seeker' ? redirect : { name: auth.role === 'seeker' && !isLogin.value ? 'profile' : HOME_BY_ROLE[auth.role] })
  } catch (error) {
    errors.value = fieldErrors(error)
    if (!Object.keys(errors.value).length) message.value = errorMessage(error)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="hero-grid flex min-h-[calc(100dvh-68px)] items-center py-10">
    <div class="container-app relative flex justify-center">
      <div class="card w-full max-w-md p-6 shadow-card sm:p-8">
        <AppLogo :with-text="false" :size="40" />
        <h1 class="mt-5 text-2xl">{{ isLogin ? 'Masuk ke Karsa' : 'Buat akun Karsa' }}</h1>
        <p class="mt-1 text-sm text-muted">{{ isLogin ? 'Lanjutkan lamaran dan cek statusnya.' : 'Gratis untuk pencari kerja dan perusahaan.' }}</p>

        <div v-if="!isLogin && !form.role" class="mt-6 grid gap-3">
          <button v-for="option in [{ role: 'seeker', icon: GraduationCap, title: 'Saya mencari kerja', text: 'Magang, entry level, atau junior' }, { role: 'company', icon: Building2, title: 'Saya merekrut', text: 'Pasang lowongan dan cari talenta' }]" :key="option.role" type="button" class="flex items-center gap-4 rounded-card border border-line p-4 text-left transition hover:border-primary hover:bg-primary-soft/50" @click="form.role = option.role">
            <span class="grid size-11 place-items-center rounded-btn bg-primary-soft text-primary"><component :is="option.icon" class="size-5" aria-hidden="true" /></span>
            <span><span class="block font-semibold">{{ option.title }}</span><span class="text-[14px] text-muted">{{ option.text }}</span></span>
          </button>
        </div>

        <form v-else class="mt-6 flex flex-col gap-4" novalidate @submit.prevent="submit">
          <div v-if="!isLogin" class="flex items-center justify-between rounded-btn bg-surface-2 px-3.5 py-2.5 text-sm">
            <span>Daftar sebagai <b>{{ form.role === 'company' ? 'perusahaan' : 'pencari kerja' }}</b></span>
            <button type="button" class="link text-[14px]" @click="form.role = ''">Ganti</button>
          </div>
          <FormField v-if="!isLogin" id="name" label="Nama lengkap" :error="errors.name" required v-slot="{ describedBy, invalid }">
            <input id="name" v-model="form.name" class="input" autocomplete="name" :aria-describedby="describedBy" :aria-invalid="invalid" />
          </FormField>
          <FormField v-if="!isLogin && form.role === 'company'" id="company-name" label="Nama perusahaan" :error="errors.company_name" required v-slot="{ describedBy, invalid }">
            <input id="company-name" v-model="form.company_name" class="input" autocomplete="organization" :aria-describedby="describedBy" :aria-invalid="invalid" />
          </FormField>
          <FormField id="email" label="Email" :error="errors.email" required v-slot="{ describedBy, invalid }">
            <input id="email" v-model="form.email" type="email" class="input" autocomplete="email" :aria-describedby="describedBy" :aria-invalid="invalid" />
          </FormField>
          <FormField id="password" label="Password" :error="errors.password" :hint="isLogin ? '' : 'Minimal 8 karakter'" required v-slot="{ describedBy, invalid }">
            <div class="relative">
              <input id="password" v-model="form.password" :type="showPassword ? 'text' : 'password'" class="input pr-12" :autocomplete="isLogin ? 'current-password' : 'new-password'" :aria-describedby="describedBy" :aria-invalid="invalid" />
              <button type="button" class="btn-icon absolute top-0 right-0" :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'" @click="showPassword = !showPassword">
                <EyeOff v-if="showPassword" class="size-[18px]" aria-hidden="true" /><Eye v-else class="size-[18px]" aria-hidden="true" />
              </button>
            </div>
          </FormField>
          <p v-if="message" class="rounded-btn bg-danger/10 px-3.5 py-2.5 text-sm text-danger" role="alert">{{ message }}</p>
          <button type="submit" class="btn-primary mt-1 h-12" :disabled="busy">{{ busy ? 'Memproses…' : isLogin ? 'Masuk' : 'Daftar' }}</button>
        </form>

        <p class="mt-6 text-center text-sm text-muted">
          <template v-if="isLogin">Belum punya akun? <RouterLink :to="{ name: 'register', query: route.query }" class="link">Daftar</RouterLink></template>
          <template v-else>Sudah punya akun? <RouterLink :to="{ name: 'login', query: route.query }" class="link">Masuk</RouterLink></template>
        </p>
      </div>
    </div>
  </div>
</template>
