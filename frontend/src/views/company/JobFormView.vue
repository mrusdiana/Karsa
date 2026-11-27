<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown, ArrowUp, ChevronRight, Plus, Trash2 } from 'lucide-vue-next'
import FormField from '@/components/FormField.vue'
import LoadingCards from '@/components/LoadingCards.vue'
import SkillInput from '@/components/SkillInput.vue'
import StatusState from '@/components/StatusState.vue'
import api, { errorMessage, fieldErrors } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { JOB_LEVEL, JOB_TYPE, WORK_MODE } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToastStore()
const editingId = route.params.id
const categories = ref([])
const loading = ref(true)
const loadError = ref(false)
const saving = ref('')
const errors = ref({})
const currentStatus = ref('draft')
const maxDeadline = new Date(Date.now() + 90 * 86400000).toLocaleDateString('sv')
const today = new Date().toLocaleDateString('sv')
const form = reactive({
  category_id: '', title: '', type: 'full_time', level: 'entry', work_mode: 'onsite', location: '',
  salary_min: '', salary_max: '', show_salary: true, skills: [], description: '', requirements: '', deadline: '', screening_questions: [],
})

const toMillions = (value) => (value == null ? '' : value / 1_000_000)
const fromMillions = (value) => (value === '' || value == null ? null : Math.round(Number(value) * 1_000_000))

function toQuestion(question) {
  const [, operator = '>=', number = ''] = /^(>=|<=)(\d+)$/.exec(question.expected_answer ?? '') ?? []
  return { question: question.question, type: question.type, required: question.required, yesNo: question.type === 'yes_no' ? (question.expected_answer ?? '') : 'yes', operator, number }
}

function expectedAnswer(question) {
  if (question.type === 'yes_no') return question.yesNo
  if (question.type === 'number' && question.number !== '') return `${question.operator}${question.number}`
  return ''
}

onMounted(async () => {
  try {
    categories.value = (await api.get('/categories')).data
    if (editingId) {
      const job = (await api.get(`/company/jobs/${editingId}`)).data
      Object.assign(form, {
        category_id: job.category_id, title: job.title, type: job.type, level: job.level, work_mode: job.work_mode, location: job.location ?? '',
        salary_min: toMillions(job.salary_min), salary_max: toMillions(job.salary_max), show_salary: job.show_salary, skills: [...job.skills],
        description: job.description, requirements: job.requirements, deadline: job.deadline, screening_questions: job.screening_questions.map(toQuestion),
      })
      currentStatus.value = job.status
    }
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
  }
})

const addQuestion = () => form.screening_questions.push({ question: '', type: 'yes_no', required: true, yesNo: 'yes', operator: '>=', number: '' })
function moveQuestion(index, offset) {
  const questions = form.screening_questions
  ;[questions[index], questions[index + offset]] = [questions[index + offset], questions[index]]
}

const actions = computed(() => {
  if (currentStatus.value === 'open') return [{ status: 'open', label: 'Simpan perubahan', primary: true }]
  if (currentStatus.value === 'closed') return [{ status: 'closed', label: 'Simpan perubahan', primary: true }]
  return [{ status: 'draft', label: 'Simpan draft' }, { status: 'open', label: 'Buka lowongan', primary: true, needsVerification: true }]
})

async function submit(status) {
  saving.value = status
  errors.value = {}
  const body = {
    ...form,
    category_id: Number(form.category_id) || 0,
    salary_min: fromMillions(form.salary_min),
    salary_max: fromMillions(form.salary_max),
    status,
    screening_questions: form.screening_questions.map((question) => ({ question: question.question, type: question.type, required: question.required, expected_answer: expectedAnswer(question) })),
  }
  try {
    const job = editingId ? (await api.put(`/company/jobs/${editingId}`, body)).data : (await api.post('/company/jobs', body)).data
    toast.show(status === 'open' && currentStatus.value !== 'open' ? 'Lowongan dibuka.' : 'Lowongan disimpan.')
    router.push({ name: 'company-jobs', query: { status: job.status } })
  } catch (failure) {
    errors.value = fieldErrors(failure)
    toast.show(Object.keys(errors.value).length ? 'Periksa lagi isian yang ditandai.' : errorMessage(failure), 'danger')
  } finally {
    saving.value = ''
  }
}
</script>

<template>
  <nav class="mb-4 flex items-center gap-1.5 text-[14px] text-muted" aria-label="Breadcrumb"><RouterLink :to="{ name: 'company-jobs' }" class="hover:text-text">Lowongan</RouterLink><ChevronRight class="size-3.5" aria-hidden="true" /><span class="text-text">{{ editingId ? 'Edit' : 'Baru' }}</span></nav>
  <h1 class="mb-6 text-2xl sm:text-[28px]">{{ editingId ? 'Edit lowongan' : 'Buat lowongan' }}</h1>
  <StatusState v-if="loadError" type="error" title="Lowongan tidak ditemukan" @retry="router.go(0)" />
  <LoadingCards v-else-if="loading" :count="3" height="h-48" />
  <form v-else class="flex max-w-3xl flex-col gap-5" novalidate @submit.prevent="submit(actions.at(-1).status)">
    <section class="card flex flex-col gap-5 p-5 sm:p-6">
      <h2 class="text-lg">1. Info dasar</h2>
      <FormField id="title" label="Judul posisi" :error="errors.title" required v-slot="{ describedBy }"><input id="title" v-model="form.title" class="input" maxlength="120" placeholder="Contoh: Frontend Developer Intern" :aria-describedby="describedBy" /></FormField>
      <div class="grid gap-5 sm:grid-cols-2">
        <FormField id="category" label="Kategori" :error="errors.category_id" required v-slot="{ describedBy }">
          <select id="category" v-model="form.category_id" class="input" :aria-describedby="describedBy"><option value="">Pilih kategori</option><option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option></select>
        </FormField>
        <FormField id="type" label="Tipe" :error="errors.type" required v-slot="{ describedBy }"><select id="type" v-model="form.type" class="input" :aria-describedby="describedBy"><option v-for="(label, value) in JOB_TYPE" :key="value" :value="value">{{ label }}</option></select></FormField>
        <FormField id="level" label="Level" :error="errors.level" required v-slot="{ describedBy }"><select id="level" v-model="form.level" class="input" :aria-describedby="describedBy"><option v-for="(label, value) in JOB_LEVEL" :key="value" :value="value">{{ label }}</option></select></FormField>
        <FormField id="work-mode" label="Mode kerja" :error="errors.work_mode" required v-slot="{ describedBy }"><select id="work-mode" v-model="form.work_mode" class="input" :aria-describedby="describedBy"><option v-for="(label, value) in WORK_MODE" :key="value" :value="value">{{ label }}</option></select></FormField>
      </div>
    </section>

    <section class="card flex flex-col gap-5 p-5 sm:p-6">
      <h2 class="text-lg">2. Lokasi, gaji & deadline</h2>
      <div class="grid gap-5 sm:grid-cols-2">
        <FormField id="location" label="Lokasi (kota)" :error="errors.location" :hint="form.work_mode === 'remote' ? 'Opsional untuk remote' : ''" :required="form.work_mode !== 'remote'" v-slot="{ describedBy }"><input id="location" v-model="form.location" class="input" maxlength="100" :aria-describedby="describedBy" /></FormField>
        <FormField id="deadline" label="Deadline" :error="errors.deadline" hint="Maksimal 90 hari dari hari ini" required v-slot="{ describedBy }"><input id="deadline" v-model="form.deadline" type="date" class="input" :min="today" :max="maxDeadline" :aria-describedby="describedBy" /></FormField>
        <FormField id="salary-min" label="Gaji minimal (juta/bulan)" :error="errors.salary_min" v-slot="{ describedBy }"><input id="salary-min" v-model="form.salary_min" type="number" min="0" step="0.5" class="input" placeholder="5" :aria-describedby="describedBy" /></FormField>
        <FormField id="salary-max" label="Gaji maksimal (juta/bulan)" :error="errors.salary_max" v-slot="{ describedBy }"><input id="salary-max" v-model="form.salary_max" type="number" min="0" step="0.5" class="input" placeholder="8" :aria-describedby="describedBy" /></FormField>
      </div>
      <label class="flex items-center gap-2.5 text-sm"><input v-model="form.show_salary" type="checkbox" class="size-[18px] accent-primary" />Tampilkan gaji di halaman lowongan</label>
    </section>

    <section class="card flex flex-col gap-5 p-5 sm:p-6">
      <h2 class="text-lg">3. Skill & deskripsi</h2>
      <FormField id="skills" label="Skill yang dibutuhkan" :error="errors.skills" hint="1–10 skill, dipakai untuk menghitung skor kecocokan" required v-slot="{ describedBy, invalid }"><SkillInput id="skills" v-model="form.skills" :max="10" :described-by="describedBy" :invalid="invalid" /></FormField>
      <FormField id="description" label="Deskripsi pekerjaan" :error="errors.description" required v-slot="{ describedBy }"><textarea id="description" v-model="form.description" class="input" rows="5" maxlength="10000" :aria-describedby="describedBy" /></FormField>
      <FormField id="requirements" label="Persyaratan" :error="errors.requirements" hint="Satu persyaratan per baris" required v-slot="{ describedBy }"><textarea id="requirements" v-model="form.requirements" class="input" rows="5" maxlength="10000" :aria-describedby="describedBy" /></FormField>
    </section>

    <section class="card flex flex-col gap-4 p-5 sm:p-6">
      <div><h2 class="text-lg">4. Pertanyaan screening</h2><p class="mt-1 text-[14px] text-muted">Opsional, maksimal 5. Jawaban yang diharapkan tidak ditampilkan ke pelamar.</p></div>
      <div v-for="(question, index) in form.screening_questions" :key="index" class="flex flex-col gap-3 rounded-btn border border-line bg-bg p-4">
        <div class="flex items-start gap-2">
          <FormField :id="`question-${index}`" :label="`Pertanyaan ${index + 1}`" :error="errors[`screening_questions[${index}].question`]" class="flex-1" v-slot="{ describedBy }"><input :id="`question-${index}`" v-model="question.question" class="input" maxlength="255" placeholder="Contoh: Bersedia bekerja di kantor?" :aria-describedby="describedBy" /></FormField>
          <div class="mt-6 flex">
            <button type="button" class="btn-icon size-10" :disabled="index === 0" aria-label="Naikkan" @click="moveQuestion(index, -1)"><ArrowUp class="size-4" aria-hidden="true" /></button>
            <button type="button" class="btn-icon size-10" :disabled="index === form.screening_questions.length - 1" aria-label="Turunkan" @click="moveQuestion(index, 1)"><ArrowDown class="size-4" aria-hidden="true" /></button>
            <button type="button" class="btn-icon size-10 hover:text-danger" aria-label="Hapus pertanyaan" @click="form.screening_questions.splice(index, 1)"><Trash2 class="size-4" aria-hidden="true" /></button>
          </div>
        </div>
        <div class="flex flex-wrap items-end gap-3">
          <div><label :for="`question-type-${index}`" class="label">Tipe jawaban</label><select :id="`question-type-${index}`" v-model="question.type" class="input w-auto"><option value="yes_no">Ya / Tidak</option><option value="number">Angka</option><option value="text">Teks singkat</option></select></div>
          <div v-if="question.type === 'yes_no'"><span class="label">Jawaban yang diharapkan</span>
            <div class="flex gap-2"><label v-for="option in [{ value: 'yes', label: 'Ya' }, { value: 'no', label: 'Tidak' }, { value: '', label: 'Bebas' }]" :key="option.value" class="flex min-h-11 items-center gap-2 rounded-btn border border-line bg-surface px-3 text-sm has-[:checked]:border-primary"><input v-model="question.yesNo" type="radio" :name="`expected-${index}`" :value="option.value" class="accent-primary" />{{ option.label }}</label></div>
          </div>
          <div v-else-if="question.type === 'number'" class="flex items-end gap-2">
            <div><label :for="`operator-${index}`" class="label">Syarat</label><select :id="`operator-${index}`" v-model="question.operator" class="input w-auto"><option value=">=">≥</option><option value="<=">≤</option></select></div>
            <div><label :for="`number-${index}`" class="label">Angka</label><input :id="`number-${index}`" v-model="question.number" type="number" min="0" class="input w-24" /></div>
          </div>
          <label class="flex min-h-11 items-center gap-2 text-sm"><input v-model="question.required" type="checkbox" class="size-[18px] accent-primary" />Wajib</label>
        </div>
        <p v-if="errors[`screening_questions[${index}].expected_answer`]" class="text-[14px] text-danger">{{ errors[`screening_questions[${index}].expected_answer`] }}</p>
      </div>
      <button v-if="form.screening_questions.length < 5" type="button" class="btn-secondary self-start" @click="addQuestion"><Plus class="size-4" aria-hidden="true" />Tambah pertanyaan</button>
    </section>

    <p v-if="errors.status" class="text-sm text-danger" role="alert">{{ errors.status }}</p>
    <div class="sticky bottom-0 -mx-4 flex flex-wrap justify-end gap-2 border-t border-line bg-bg/95 px-4 py-3 backdrop-blur sm:mx-0 sm:px-0">
      <RouterLink :to="{ name: 'company-jobs' }" class="btn-ghost">Batal</RouterLink>
      <button
        v-for="action in actions"
        :key="action.status"
        :type="action.primary ? 'submit' : 'button'"
        :class="action.primary ? 'btn-primary' : 'btn-secondary'"
        :disabled="!!saving || (action.needsVerification && !auth.isVerified)"
        :title="action.needsVerification && !auth.isVerified ? 'Lowongan bisa dibuka setelah perusahaan terverifikasi' : undefined"
        @click.prevent="submit(action.status)"
      >{{ saving === action.status ? 'Menyimpan…' : action.label }}</button>
    </div>
  </form>
</template>
