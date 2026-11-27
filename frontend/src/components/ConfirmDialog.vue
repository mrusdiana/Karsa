<script setup>
import ModalDialog from '@/components/ModalDialog.vue'

defineProps({ open: Boolean, title: String, text: String, confirmLabel: { type: String, default: 'Hapus' }, danger: { type: Boolean, default: true }, busy: Boolean })
defineEmits(['confirm', 'cancel'])
</script>

<template>
  <ModalDialog :open="open" :title="title" @close="$emit('cancel')">
    <p class="text-sm text-muted">{{ text }}</p>
    <slot />
    <div class="mt-6 flex justify-end gap-2">
      <button type="button" class="btn-ghost" @click="$emit('cancel')">Batal</button>
      <button type="button" :class="danger ? 'btn-danger' : 'btn-primary'" :disabled="busy" @click="$emit('confirm')">{{ busy ? 'Memproses…' : confirmLabel }}</button>
    </div>
  </ModalDialog>
</template>
