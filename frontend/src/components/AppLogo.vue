<script setup>
defineProps({ size: { type: Number, default: 34 }, withText: { type: Boolean, default: true } })
</script>

<template>
  <span class="inline-flex items-center gap-2.5 font-display text-[21px] font-extrabold tracking-tight text-text">
    <svg :width="size" :height="size" viewBox="0 0 34 34" aria-hidden="true">
      <rect width="34" height="34" rx="10" fill="#0F766E" />
      <path d="M11 8v18M11 18l9-10M14 15l5.5 6.5L26 13M21.5 13H26v4.5" stroke="#fff" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
    <span v-if="withText">Karsa</span>
    <span v-else class="sr-only">Karsa</span>
  </span>
</template>
