<script setup>
import { computed } from 'vue'
import { SCORE_LABEL, scoreLevel } from '@/utils/format'

const props = defineProps({ score: { type: Number, required: true }, size: { type: Number, default: 40 } })
const colors = { high: '#15803D', medium: '#B45309', low: '#6B7280' }
const level = computed(() => scoreLevel(props.score))
const circumference = 2 * Math.PI * 16
</script>

<template>
  <span class="relative inline-grid shrink-0 place-items-center" :style="{ width: `${size}px`, height: `${size}px` }" role="img" :aria-label="`${SCORE_LABEL[level]}, skor ${score}`">
    <svg viewBox="0 0 40 40" class="absolute inset-0 -rotate-90" aria-hidden="true">
      <circle cx="20" cy="20" r="16" fill="none" stroke="#F3F4F1" stroke-width="4" />
      <circle cx="20" cy="20" r="16" fill="none" :stroke="colors[level]" stroke-width="4" stroke-linecap="round" :stroke-dasharray="`${(circumference * score) / 100} ${circumference}`" />
    </svg>
    <span class="num relative text-[13.5px] font-bold" :style="{ fontSize: `${size * 0.31}px` }">{{ score }}</span>
  </span>
</template>
