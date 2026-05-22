<script setup lang="ts">
defineOptions({
  name: 'InputFilter'
});

import type { InputFilterProps, InputFilterEmits } from './interfaces';

import { computed } from 'vue';

const props = withDefaults(defineProps<InputFilterProps>(), {
  clearable: true,
  disabled: false,
  debounce: 0,
});

const emit = defineEmits<InputFilterEmits>();

const value = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

function handleEnter() {
  emit('enter');
}
</script>

<template>
  <q-input
    v-model="value"
    outlined
    dense
    class="rounded-borders"
    :label="label"
    :clearable="clearable"
    :disable="disabled"
    :debounce="debounce"
    :mask="mask"
    :type="type"
    :placeholder="placeholder"
    @keyup.enter="handleEnter"
  >
    <template v-slot:prepend>
      <q-icon color="primary" :name="icon" />
    </template>
  </q-input>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
