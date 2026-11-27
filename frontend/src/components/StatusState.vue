<script setup>
import { AlertTriangle, SearchX } from 'lucide-vue-next'

defineProps({ type: { type: String, default: 'empty' }, title: String, text: String })
defineEmits(['retry'])
</script>

<template>
  <div class="card flex flex-col items-center gap-3 px-6 py-12 text-center">
    <AlertTriangle v-if="type === 'error'" class="size-10 text-danger" aria-hidden="true" />
    <SearchX v-else class="size-10 text-muted" aria-hidden="true" />
    <h2 class="text-lg">{{ title ?? (type === 'error' ? 'Terjadi kesalahan. Coba lagi.' : 'Belum ada data') }}</h2>
    <p v-if="text" class="max-w-sm text-sm text-muted">{{ text }}</p>
    <button v-if="type === 'error'" type="button" class="btn-secondary mt-2" @click="$emit('retry')">Coba lagi</button>
    <slot />
  </div>
</template>
