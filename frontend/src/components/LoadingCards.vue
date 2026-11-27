<script setup>
defineProps({ count: { type: Number, default: 3 }, height: { type: String, default: 'h-36' } })
</script>

<template>
  <div class="flex flex-col gap-3.5" aria-busy="true" aria-label="Memuat">
    <div v-for="index in count" :key="index" class="card animate-pulse bg-surface-2/60" :class="height" />
  </div>
</template>
