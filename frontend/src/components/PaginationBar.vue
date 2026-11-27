<script setup>
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({ page: { type: Object, required: true } })
const emit = defineEmits(['change'])
const range = computed(() => {
  const first = (props.page.page - 1) * props.page.per_page + 1
  return { first, last: Math.min(props.page.total, first + props.page.data.length - 1) }
})
const pages = computed(() => {
  const numbers = new Set([1, props.page.last_page, props.page.page - 1, props.page.page, props.page.page + 1])
  return [...numbers].filter((number) => number >= 1 && number <= props.page.last_page).sort((left, right) => left - right)
})
</script>

<template>
  <nav v-if="page.total > 0" class="mt-6 flex flex-col items-center gap-2" aria-label="Halaman">
    <div v-if="page.last_page > 1" class="flex items-center gap-1.5">
      <button type="button" class="btn-secondary size-10 min-h-10 px-0" :disabled="page.page === 1" aria-label="Sebelumnya" @click="emit('change', page.page - 1)"><ChevronLeft class="size-4" aria-hidden="true" /></button>
      <button
        v-for="number in pages"
        :key="number"
        type="button"
        class="size-10 min-h-10 px-0"
        :class="number === page.page ? 'btn-primary' : 'btn-secondary'"
        :aria-current="number === page.page ? 'page' : undefined"
        @click="emit('change', number)"
      >{{ number }}</button>
      <button type="button" class="btn-secondary size-10 min-h-10 px-0" :disabled="page.page === page.last_page" aria-label="Berikutnya" @click="emit('change', page.page + 1)"><ChevronRight class="size-4" aria-hidden="true" /></button>
    </div>
    <p class="text-[14px] text-muted">Menampilkan {{ range.first }}–{{ range.last }} dari {{ page.total }}</p>
  </nav>
</template>
