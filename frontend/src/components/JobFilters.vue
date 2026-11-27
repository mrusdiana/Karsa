<script setup>
import { JOB_LEVEL, JOB_TYPE, WORK_MODE } from '@/utils/format'

defineProps({ filters: { type: Object, required: true }, categories: { type: Array, default: () => [] } })
defineEmits(['toggle', 'set', 'clear'])
const salaryOptions = [
  { value: '', label: 'Semua gaji' },
  { value: '3000000', label: '≥ Rp 3 jt' },
  { value: '5000000', label: '≥ Rp 5 jt' },
  { value: '8000000', label: '≥ Rp 8 jt' },
  { value: '12000000', label: '≥ Rp 12 jt' },
]
const groups = [
  { key: 'level', title: 'Level', options: JOB_LEVEL },
  { key: 'mode', title: 'Mode kerja', options: WORK_MODE },
  { key: 'type', title: 'Tipe', options: JOB_TYPE },
]
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center justify-between"><h2 class="text-[18px]">Filter</h2><button type="button" class="link text-[14px]" @click="$emit('clear')">Hapus filter</button></div>
    <div>
      <label for="filter-category" class="label">Kategori</label>
      <select id="filter-category" class="input" :value="filters.category" @change="$emit('set', 'category', $event.target.value)">
        <option value="">Semua kategori</option>
        <option v-for="category in categories" :key="category.slug" :value="category.slug">{{ category.name }} ({{ category.job_count }})</option>
      </select>
    </div>
    <fieldset v-for="group in groups" :key="group.key">
      <legend class="label">{{ group.title }}</legend>
      <label v-for="(label, value) in group.options" :key="value" class="flex min-h-9 cursor-pointer items-center gap-2.5 text-sm">
        <input type="checkbox" class="size-[18px] accent-primary" :checked="filters[group.key].includes(value)" @change="$emit('toggle', group.key, value)" />{{ label }}
      </label>
    </fieldset>
    <div>
      <label for="filter-location" class="label">Lokasi</label>
      <input id="filter-location" class="input" :value="filters.location" placeholder="Contoh: Jakarta" @change="$emit('set', 'location', $event.target.value.trim())" />
    </div>
    <div>
      <label for="filter-salary" class="label">Gaji minimal</label>
      <select id="filter-salary" class="input" :value="filters.salary_min" @change="$emit('set', 'salary_min', $event.target.value)">
        <option v-for="option in salaryOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
      </select>
    </div>
  </div>
</template>
