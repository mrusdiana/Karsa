<script setup>
defineProps({ items: { type: Array, required: true } })
</script>

<template>
  <nav class="fixed inset-x-0 bottom-0 z-40 border-t border-line bg-surface/95 pb-[env(safe-area-inset-bottom)] backdrop-blur" aria-label="Navigasi utama">
    <ul class="mx-auto flex max-w-lg">
      <li v-for="item in items" :key="item.label" class="flex-1">
        <component
          :is="item.to ? 'RouterLink' : 'button'"
          :to="item.to"
          :type="item.to ? undefined : 'button'"
          class="relative flex h-16 w-full flex-col items-center justify-center gap-1 text-xs font-medium transition"
          :class="item.active ? 'text-primary' : 'text-muted hover:text-text'"
          :aria-current="item.active ? 'page' : undefined"
          @click="item.action?.()"
        >
          <span class="relative">
            <component :is="item.icon" class="size-[22px]" :stroke-width="item.active ? 2.4 : 2" aria-hidden="true" />
            <span v-if="item.badge" class="absolute -top-1.5 -right-2.5 grid h-[18px] min-w-[18px] place-items-center rounded-full border-2 border-surface bg-danger px-1 text-[11.5px] font-semibold text-white">{{ item.badge }}</span>
          </span>
          {{ item.label }}
        </component>
      </li>
    </ul>
  </nav>
</template>
