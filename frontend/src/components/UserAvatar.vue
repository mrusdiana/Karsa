<script setup>
import { initials } from '@/utils/format'

defineProps({ name: { type: String, default: '' }, url: String, size: { type: Number, default: 36 } })
</script>

<template>
  <img v-if="url" :src="url" :alt="`Foto ${name}`" class="shrink-0 rounded-full border border-line object-cover" :style="{ width: `${size}px`, height: `${size}px` }" />
  <span v-else class="grid shrink-0 place-items-center rounded-full bg-primary-soft font-semibold text-primary" :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${Math.max(12, size * 0.36)}px` }" role="img" :aria-label="`Foto ${name}`">{{ initials(name) }}</span>
</template>
