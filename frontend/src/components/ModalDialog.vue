<script setup>
import { nextTick, ref, watch } from 'vue'
import { X } from 'lucide-vue-next'

const props = defineProps({ open: Boolean, title: String, wide: Boolean })
const emit = defineEmits(['close'])
const dialog = ref()

watch(
  () => props.open,
  async (isOpen) => {
    await nextTick()
    if (isOpen && !dialog.value?.open) dialog.value?.showModal()
    if (!isOpen) dialog.value?.close()
  },
  { immediate: true },
)
</script>

<template>
  <dialog
    ref="dialog"
    class="m-auto max-h-[90dvh] w-[calc(100%-2rem)] overflow-y-auto rounded-card border border-line bg-surface p-0 text-text shadow-card backdrop:bg-text/40"
    :class="wide ? 'max-w-xl' : 'max-w-md'"
    :aria-label="title"
    @close="open && emit('close')"
  >
    <div class="flex items-center justify-between gap-4 border-b border-line px-6 py-4">
      <h2 class="text-lg">{{ title }}</h2>
      <button type="button" class="btn-icon -mr-2" aria-label="Tutup" @click="emit('close')"><X class="size-5" aria-hidden="true" /></button>
    </div>
    <div class="p-6"><slot /></div>
  </dialog>
</template>
