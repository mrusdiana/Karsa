<script setup>
defineProps({ title: String, text: String })
</script>

<template>
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <slot name="before" />
      <h1 class="text-2xl leading-tight sm:text-[28px]">{{ title }}</h1>
      <p v-if="text" class="mt-1 text-sm text-muted">{{ text }}</p>
    </div>
    <div class="flex flex-wrap gap-2"><slot /></div>
  </div>
</template>
