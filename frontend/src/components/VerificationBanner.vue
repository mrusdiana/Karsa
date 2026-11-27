<script setup>
import { AlertTriangle, Clock } from 'lucide-vue-next'

defineProps({ company: { type: Object, required: true } })
</script>

<template>
  <div v-if="company.status === 'pending'" class="mb-6 flex items-start gap-3 rounded-card border border-warning/25 bg-warning/[0.06] p-4 text-sm">
    <Clock class="mt-0.5 size-5 shrink-0 text-warning" aria-hidden="true" />
    <p><b>Profil perusahaan sedang ditinjau admin.</b> Kamu bisa menyiapkan draft lowongan. Lowongan bisa dibuka setelah perusahaan terverifikasi.</p>
  </div>
  <div v-else-if="company.status === 'rejected'" class="mb-6 flex flex-col gap-3 rounded-card border border-danger/25 bg-danger/[0.06] p-4 text-sm sm:flex-row sm:items-center">
    <AlertTriangle class="size-5 shrink-0 text-danger" aria-hidden="true" />
    <p class="flex-1"><b>Verifikasi ditolak.</b> Alasan: {{ company.rejection_reason }}</p>
    <RouterLink :to="{ name: 'company-profile' }" class="btn-primary btn-sm">Perbaiki & ajukan ulang</RouterLink>
  </div>
</template>
