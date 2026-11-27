<script setup>
import { computed } from 'vue'
import { initials } from '@/utils/format'

const props = defineProps({ name: { type: String, default: '' }, url: String, size: { type: Number, default: 40 } })
const palette = [
  ['#E6F2F0', '#0F766E'],
  ['#E8EEFD', '#2563EB'],
  ['#FDF0E3', '#B45309'],
  ['#F1EAFB', '#7C3AED'],
  ['#FCE8EC', '#BE123C'],
  ['#E7F5EA', '#15803D'],
]
const colors = computed(() => palette[[...props.name].reduce((sum, letter) => sum + letter.charCodeAt(0), 0) % palette.length])
</script>

<template>
  <img v-if="url" :src="url" :alt="`Logo ${name}`" class="shrink-0 rounded-logo border border-line bg-surface object-cover" :style="{ width: `${size}px`, height: `${size}px` }" />
  <span
    v-else
    class="grid shrink-0 place-items-center rounded-logo font-display font-extrabold"
    :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${size * 0.36}px`, background: colors[0], color: colors[1] }"
    role="img"
    :aria-label="`Logo ${name}`"
  >{{ initials(name) }}</span>
</template>
