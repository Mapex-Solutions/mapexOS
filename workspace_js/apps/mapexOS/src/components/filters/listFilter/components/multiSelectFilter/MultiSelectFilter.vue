<script setup lang="ts">
defineOptions({
  name: 'MultiSelectFilter'
});

import type { MultiSelectFilterProps, MultiSelectFilterEmits } from './interfaces';

import { computed } from 'vue';

const props = withDefaults(defineProps<MultiSelectFilterProps>(), {
  clearable: true,
  disabled: false,
});

const emit = defineEmits<MultiSelectFilterEmits>();

const value = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
});

const hasIcons = computed(() => props.options.some(opt => opt.icon));

function handleEnter() {
  emit('enter');
}
</script>

<template>
  <q-select
    v-model="value"
    outlined
    dense
    multiple
    use-chips
    options-dense
    emit-value
    map-options
    class="rounded-borders"
    :label="label"
    :options="options"
    :clearable="clearable"
    :disable="disabled"
    @keyup.enter="handleEnter"
  >
    <template v-slot:prepend>
      <q-icon color="primary" :name="icon" />
    </template>

    <template v-if="hasIcons" v-slot:option="scope">
      <q-item v-bind="scope.itemProps">
        <q-item-section v-if="scope.opt.icon" avatar>
          <q-icon :name="scope.opt.icon" :color="scope.opt.color || 'primary'" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ scope.opt.label }}</q-item-label>
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>

<style lang="scss" scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
