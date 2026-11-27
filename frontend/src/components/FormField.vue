<script setup>
defineProps({ label: String, id: { type: String, required: true }, error: String, hint: String, required: Boolean })
</script>

<template>
  <div>
    <label v-if="label" :for="id" class="label">{{ label }}<span v-if="required" class="text-danger" aria-hidden="true"> *</span></label>
    <slot :described-by="error || hint ? `${id}-note` : undefined" :invalid="!!error" />
    <p v-if="error" :id="`${id}-note`" class="mt-1 text-[14px] text-danger">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-note`" class="mt-1 text-[14px] text-muted">{{ hint }}</p>
  </div>
</template>
