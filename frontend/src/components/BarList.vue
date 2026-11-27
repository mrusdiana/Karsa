<script setup>
import { computed } from 'vue'

const props = defineProps({ items: { type: Array, required: true }, empty: { type: String, default: 'Belum ada data.' } })
const maximum = computed(() => Math.max(1, ...props.items.map((item) => item.value)))
</script>

<template>
  <ul v-if="items.length" class="flex flex-col gap-3.5">
    <li v-for="item in items" :key="item.label">
      <div class="flex justify-between gap-3 text-sm"><span class="truncate">{{ item.label }}</span><b class="num">{{ item.value }}</b></div>
      <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-surface-2"><div class="h-full rounded-full bg-primary" :style="{ width: `${(item.value / maximum) * 100}%` }" /></div>
    </li>
  </ul>
  <p v-else class="text-sm text-muted">{{ empty }}</p>
</template>
