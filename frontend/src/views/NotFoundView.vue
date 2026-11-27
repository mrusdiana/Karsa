<script setup>
import StatusState from '@/components/StatusState.vue'
</script>

<template>
  <div class="container-app py-16">
    <StatusState title="Halaman tidak ditemukan" text="Tautannya mungkin salah atau halamannya sudah dipindah.">
      <RouterLink :to="{ name: 'home' }" class="btn-primary mt-2">Ke beranda</RouterLink>
    </StatusState>
  </div>
</template>
