<script setup lang="ts">
/** TYPE IMPORTS (ALL types first, grouped) */
import type { FieldInputSelectorProps, FieldInputSelectorEmits, FieldType } from './interfaces';

defineOptions({
  name: 'FieldInputSelector'
});

/** VUE IMPORTS */
import { computed } from 'vue';

/** LOCAL IMPORTS (constants and handlers ONLY - NO types here!) */
import { FIELD_TYPE_SELECTOR_OPTIONS, DEFAULT_FIELD_TYPE } from './constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<FieldInputSelectorProps>(), {
  modelValue: DEFAULT_FIELD_TYPE,
  label: 'Type',
  disabled: false,
  showIcons: true,
});

const emit = defineEmits<FieldInputSelectorEmits>();

/** COMPUTED */

/**
 * Current selected field type
 */
const selectedType = computed({
  get: () => props.modelValue || DEFAULT_FIELD_TYPE,
  set: (value: FieldType) => emit('update:modelValue', value),
});

/**
 * Current option details (icon, color, etc.)
 */
const currentOption = computed(() => {
  return FIELD_TYPE_SELECTOR_OPTIONS.find(opt => opt.value === selectedType.value) || FIELD_TYPE_SELECTOR_OPTIONS[3];
});
</script>

<template>
  <q-select
    v-model="selectedType"
    :options="FIELD_TYPE_SELECTOR_OPTIONS"
    :label="label"
    outlined
    dense
    option-label="label"
    option-value="value"
    emit-value
    map-options
    options-dense
    :disable="disabled"
  >
    <template v-if="showIcons" #prepend>
      <q-icon v-if="currentOption" :name="currentOption.icon" :color="currentOption.color" />
    </template>

    <template v-if="showIcons" #option="scope">
      <q-item v-bind="scope.itemProps">
        <q-item-section avatar>
          <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ scope.opt.label }}</q-item-label>
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>
