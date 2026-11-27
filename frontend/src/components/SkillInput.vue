<script setup>
import { ref } from 'vue'
import { X } from 'lucide-vue-next'

const skills = defineModel({ type: Array, default: () => [] })
const props = defineProps({ id: String, max: { type: Number, default: 15 }, describedBy: String, invalid: Boolean })
const draft = ref('')

function add() {
  const skill = draft.value.replace(/,/g, '').trim()
  draft.value = ''
  if (!skill || skills.value.length >= props.max) return
  if (skills.value.some((existing) => existing.toLowerCase() === skill.toLowerCase())) return
  skills.value = [...skills.value, skill]
}

function onKeydown(event) {
  if (event.key === 'Enter' || event.key === ',') {
    event.preventDefault()
    add()
  } else if (event.key === 'Backspace' && !draft.value && skills.value.length) {
    skills.value = skills.value.slice(0, -1)
  }
}
</script>

<template>
  <div class="input flex min-h-11 flex-wrap items-center gap-1.5 py-1.5" :class="invalid && 'border-danger'">
    <span v-for="skill in skills" :key="skill" class="chip h-7 bg-primary-soft pr-1 text-primary">
      {{ skill }}
      <button type="button" class="grid size-5 place-items-center rounded-full hover:bg-primary/15" :aria-label="`Hapus ${skill}`" @click="skills = skills.filter((item) => item !== skill)">
        <X class="size-3.5" aria-hidden="true" />
      </button>
    </span>
    <input
      :id="id"
      v-model="draft"
      class="min-w-32 flex-1 bg-transparent outline-none placeholder:text-muted focus-visible:shadow-none"
      :placeholder="skills.length >= max ? `Maksimal ${max} skill` : 'Ketik lalu tekan Enter'"
      :disabled="skills.length >= max"
      :aria-describedby="describedBy"
      @keydown="onKeydown"
      @blur="add"
    />
  </div>
</template>
