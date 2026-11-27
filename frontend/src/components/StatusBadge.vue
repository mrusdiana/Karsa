<script setup>
import { computed } from 'vue'
import { BadgeCheck, Ban, CalendarClock, CheckCircle2, Clock, Eye, FilePen, Hourglass, Lock, Send, Undo2, XCircle } from 'lucide-vue-next'

const props = defineProps({ status: { type: Object, required: true }, small: Boolean })
const icons = { BadgeCheck, Ban, CalendarClock, CheckCircle2, Clock, Eye, FilePen, Hourglass, Lock, Send, Undo2, XCircle }
const tones = {
  primary: 'bg-primary/10 text-primary',
  info: 'bg-info/10 text-info',
  success: 'bg-success/10 text-success',
  warning: 'bg-warning/10 text-warning',
  danger: 'bg-danger/10 text-danger',
  neutral: 'bg-neutral/12 text-neutral',
}
const icon = computed(() => icons[props.status.icon])
</script>

<template>
  <span class="inline-flex h-[26px] items-center gap-1.5 rounded-full px-2.5 text-[13.5px] font-semibold whitespace-nowrap" :class="tones[status.tone]">
    <component :is="icon" class="size-3.5" aria-hidden="true" />{{ status.label }}
  </span>
</template>
