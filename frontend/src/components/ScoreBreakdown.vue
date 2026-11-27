<script setup>
import { Check } from 'lucide-vue-next'

defineProps({ breakdown: { type: Object, required: true }, preview: Boolean })
</script>

<template>
  <div class="flex flex-col gap-4 rounded-xl border border-line bg-bg p-4 text-[14px]">
    <div>
      <div class="flex justify-between"><span>Skill</span><b class="num">{{ breakdown.skill.points }}/60</b></div>
      <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-surface-2"><div class="h-full rounded-full bg-primary" :style="{ width: `${(breakdown.skill.points / 60) * 100}%` }" /></div>
      <div class="mt-2 flex flex-wrap gap-1.5">
        <span v-for="skill in breakdown.skill.matched" :key="skill" class="chip-skill text-success"><Check class="size-3.5" aria-hidden="true" />{{ skill }}</span>
        <span v-for="skill in breakdown.skill.missing" :key="skill" class="chip-skill text-muted"><span class="sr-only">Kurang:</span>{{ skill }}</span>
      </div>
    </div>
    <div>
      <div class="flex justify-between"><span>Level</span><b class="num">{{ breakdown.level.points }}/20</b></div>
      <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-surface-2"><div class="h-full rounded-full" :class="breakdown.level.points === 20 ? 'bg-primary' : 'bg-warning'" :style="{ width: `${(breakdown.level.points / 20) * 100}%` }" /></div>
      <p class="mt-1 text-muted">{{ breakdown.level.note }}</p>
    </div>
    <div>
      <div class="flex justify-between"><span>Screening</span><b class="num">{{ breakdown.screening.points }}/20</b></div>
      <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-surface-2"><div class="h-full rounded-full bg-primary" :style="{ width: `${(breakdown.screening.points / 20) * 100}%` }" /></div>
      <p class="mt-1 text-muted">{{ preview ? 'Dihitung dari jawabanmu saat melamar' : breakdown.screening.total ? `${breakdown.screening.passed} dari ${breakdown.screening.total} jawaban lolos` : 'Lowongan tanpa pertanyaan screening' }}</p>
    </div>
  </div>
</template>
