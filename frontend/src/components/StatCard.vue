<script setup>
defineProps({ label: String, value: [Number, String], note: String, icon: Object, noteTone: String })
</script>

<template>
  <div class="card p-5">
    <div class="flex items-center justify-between text-[14px] font-medium text-muted">{{ label }}<component :is="icon" v-if="icon" class="size-[18px]" aria-hidden="true" /></div>
    <div class="num mt-1.5 font-display text-[30px] leading-tight font-extrabold">{{ value }}</div>
    <div v-if="note" class="mt-0.5 text-[14px] font-medium" :class="noteTone === 'success' ? 'text-success' : 'text-muted'">{{ note }}</div>
  </div>
</template>
