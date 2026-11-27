<script setup>
import { AlertCircle, CheckCircle2 } from 'lucide-vue-next'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
</script>

<template>
  <div class="pointer-events-none fixed inset-x-4 bottom-20 z-[70] flex flex-col items-center gap-2 sm:inset-x-auto sm:top-20 sm:right-6 sm:bottom-auto sm:items-end" aria-live="polite">
    <TransitionGroup enter-from-class="opacity-0 translate-y-2" leave-to-class="opacity-0" enter-active-class="transition duration-200" leave-active-class="transition duration-150">
      <div v-for="item in toast.toasts" :key="item.id" class="pointer-events-auto flex items-center gap-2.5 rounded-btn border border-line bg-surface px-4 py-3 text-sm font-medium shadow-card">
        <CheckCircle2 v-if="item.tone === 'success'" class="size-4 shrink-0 text-success" aria-hidden="true" />
        <AlertCircle v-else class="size-4 shrink-0 text-danger" aria-hidden="true" />
        {{ item.message }}
      </div>
    </TransitionGroup>
  </div>
</template>
