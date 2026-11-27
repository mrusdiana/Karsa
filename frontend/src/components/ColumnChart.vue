<script setup>
import { computed } from 'vue'

const props = defineProps({ data: { type: Array, required: true }, title: String })
const maximum = computed(() => Math.max(1, ...props.data.map((item) => item.value)))
</script>

<template>
  <figure>
    <div class="flex h-48 items-end gap-2 border-b border-line pt-6 sm:gap-3.5">
      <div v-for="item in data" :key="item.label" class="group flex h-full flex-1 flex-col items-center justify-end gap-1.5 outline-none" tabindex="0" :aria-label="`${item.label}: ${item.value}`">
        <span class="num text-xs font-semibold">{{ item.value }}</span>
        <span class="w-full max-w-10 rounded-t-md rounded-b-sm bg-primary/85 transition group-hover:bg-primary group-focus-visible:bg-primary" :style="{ height: `${(item.value / maximum) * 85}%`, minHeight: item.value ? '4px' : '0' }" />
      </div>
    </div>
    <div class="mt-2 flex gap-2 sm:gap-3.5">
      <span v-for="item in data" :key="item.label" class="flex-1 text-center text-[12px] text-muted">{{ item.label }}</span>
    </div>
    <figcaption v-if="title" class="sr-only">{{ title }}</figcaption>
  </figure>
</template>
