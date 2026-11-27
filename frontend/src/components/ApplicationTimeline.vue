<script setup>
import { APPLICATION_STATUS, formatDateTime } from '@/utils/format'

defineProps({ events: { type: Array, required: true }, byLabels: { type: Object, default: () => ({ company: 'Perusahaan', you: 'kamu', candidate: 'kandidat' }) } })
const dots = { primary: 'bg-primary', info: 'bg-info', warning: 'bg-warning', success: 'bg-success', danger: 'bg-danger', neutral: 'bg-neutral' }
</script>

<template>
  <ol class="ml-1.5 flex flex-col gap-3.5 border-l-2 border-line pl-5 text-[14.5px]">
    <li v-for="(event, index) in events" :key="index" class="relative" :class="index === events.length - 1 && 'font-semibold'">
      <span class="absolute top-1.5 -left-[27px] size-2.5 rounded-full ring-4 ring-surface" :class="dots[APPLICATION_STATUS[event.to_status].tone]" />
      {{ APPLICATION_STATUS[event.to_status].label }}
      <span class="font-normal text-muted"> · oleh {{ byLabels[event.by] }} · {{ formatDateTime(event.created_at) }}</span>
    </li>
  </ol>
</template>
