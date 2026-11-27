<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { CheckCircle2, FileText } from 'lucide-vue-next'
import FormField from '@/components/FormField.vue'
import ModalDialog from '@/components/ModalDialog.vue'
import api, { errorMessage, fieldErrors } from '@/services/api'
import { MISSING_LABEL } from '@/utils/format'

const props = defineProps({ open: Boolean, job: { type: Object, required: true }, invitationId: Number })
const emit = defineEmits(['close', 'applied'])
const profile = ref(null)
const loading = ref(false)
const submitting = ref(false)
const done = ref(false)
const errors = ref({})
const message = ref('')
const form = reactive({ cover_letter: '', answers: {} })
const missing = computed(() => profile.value?.completeness.missing ?? [])

watch(
  () => props.open,
  async (isOpen) => {
    if (!isOpen) return
    done.value = false
    errors.value = {}
    message.value = ''
    loading.value = true
    try {
      profile.value = (await api.get('/me/profile')).data
    } catch (error) {
      message.value = errorMessage(error)
    } finally {
      loading.value = false
    }
  },
)

async function submit() {
  submitting.value = true
  errors.value = {}
  message.value = ''
  try {
    await api.post(`/jobs/${props.job.id}/apply`, {
      cover_letter: form.cover_letter,
      answers: props.job.screening_questions.map((question) => ({ question_id: question.id, answer: String(form.answers[question.id] ?? '') })),
      invitation_id: props.invitationId ?? null,
    })
    done.value = true
    emit('applied')
  } catch (error) {
    errors.value = fieldErrors(error)
    message.value = errors.value.profile ? 'Lengkapi profil dan CV sebelum melamar.' : errorMessage(error)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <ModalDialog :open="open" :title="done ? 'Lamaran terkirim' : `Lamar: ${job.title}`" wide @close="emit('close')">
    <div v-if="done" class="flex flex-col items-center gap-3 py-4 text-center">
      <CheckCircle2 class="size-12 text-success" aria-hidden="true" />
      <p class="font-semibold">Lamaranmu untuk {{ job.title }} sudah terkirim.</p>
      <p class="text-sm text-muted">Kamu akan mendapat email setiap kali statusnya berubah.</p>
      <RouterLink :to="{ name: 'applications' }" class="btn-primary mt-2">Lihat Lamaran Saya</RouterLink>
    </div>
    <p v-else-if="loading" class="text-sm text-muted">Memuat profil…</p>
    <div v-else-if="missing.length">
      <p class="text-sm">Profilmu belum lengkap. Lengkapi dulu bagian berikut:</p>
      <ul class="mt-3 flex flex-col gap-1.5 text-sm">
        <li v-for="item in missing" :key="item" class="flex items-center gap-2"><span class="size-1.5 rounded-full bg-danger" />{{ MISSING_LABEL[item] }}</li>
      </ul>
      <RouterLink :to="{ name: 'profile' }" class="btn-primary mt-5 w-full">Lengkapi profil</RouterLink>
    </div>
    <form v-else class="flex flex-col gap-5" @submit.prevent="submit">
      <div class="flex items-center gap-3 rounded-btn border border-line bg-bg p-3 text-sm">
        <FileText class="size-5 text-danger" aria-hidden="true" />
        <span class="flex-1 truncate font-medium">{{ profile?.cv?.name }}</span>
        <RouterLink :to="{ name: 'profile' }" class="link text-[14px]">Ganti di profil</RouterLink>
      </div>
      <FormField id="cover-letter" label="Surat pengantar (opsional)" :error="errors.cover_letter" :hint="`${form.cover_letter.length}/2000`" v-slot="{ describedBy }">
        <textarea id="cover-letter" v-model="form.cover_letter" class="input" rows="4" maxlength="2000" :aria-describedby="describedBy" placeholder="Ceritakan singkat kenapa kamu cocok untuk posisi ini." />
      </FormField>
      <fieldset v-for="question in job.screening_questions" :key="question.id">
        <legend class="label">{{ question.question }}<span v-if="question.required" class="text-danger" aria-hidden="true"> *</span></legend>
        <div v-if="question.type === 'yes_no'" class="flex gap-2">
          <label v-for="option in [{ value: 'yes', label: 'Ya' }, { value: 'no', label: 'Tidak' }]" :key="option.value" class="flex min-h-11 flex-1 cursor-pointer items-center gap-2 rounded-btn border border-line px-3.5 text-sm has-[:checked]:border-primary has-[:checked]:bg-primary-soft">
            <input v-model="form.answers[question.id]" type="radio" :name="`question-${question.id}`" :value="option.value" class="accent-primary" />{{ option.label }}
          </label>
        </div>
        <input v-else-if="question.type === 'number'" v-model="form.answers[question.id]" type="number" min="0" step="any" class="input" :aria-label="question.question" />
        <textarea v-else v-model="form.answers[question.id]" class="input" rows="2" maxlength="500" :aria-label="question.question" />
        <p v-if="errors[`answers.${question.id}`]" class="mt-1 text-[14px] text-danger">{{ errors[`answers.${question.id}`] }}</p>
      </fieldset>
      <p v-if="message" class="rounded-btn bg-danger/10 px-3.5 py-2.5 text-sm text-danger" role="alert">{{ message }}</p>
      <div class="flex justify-end gap-2">
        <button type="button" class="btn-ghost" @click="emit('close')">Batal</button>
        <button type="submit" class="btn-primary" :disabled="submitting">{{ submitting ? 'Mengirim…' : 'Kirim lamaran' }}</button>
      </div>
    </form>
  </ModalDialog>
</template>
