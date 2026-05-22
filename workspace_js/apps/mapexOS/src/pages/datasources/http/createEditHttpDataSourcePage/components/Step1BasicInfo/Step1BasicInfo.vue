<template>
  <div class="row q-col-gutter-md">
    <!-- Name -->
    <div class="col-12 col-sm-6">
      <q-input
        v-model="localData.name"
        outlined
        dense
        hide-bottom-space
        class="rounded-borders"
        :label="`${t.basicInfo.name.value} *`"
        :rules="[(val) => !!val || t.basicInfo.nameRequired.value]"
        @update:model-value="updateValue"
      >
        <template #prepend>
          <q-icon name="label" color="primary"/>
        </template>
      </q-input>
    </div>

    <!-- Status -->
    <div class="col-12 col-sm-6">
      <q-select
        v-model="localData.enabled"
        outlined
        dense
        hide-bottom-space
        class="rounded-borders"
        :label="`${t.basicInfo.status.value} *`"
        :options="enabledOptions"
        emit-value
        map-options
        @update:model-value="updateValue"
      >
        <template #prepend>
          <q-icon name="toggle_on" color="primary"/>
        </template>
      </q-select>
    </div>

    <!-- Description - Full width -->
    <div class="col-12">
      <q-input
        v-model="localData.description"
        outlined
        dense
        hide-bottom-space
        type="textarea"
        class="rounded-borders"
        :label="t.basicInfo.description.value"
        @update:model-value="updateValue"
      >
        <template #prepend>
          <q-icon name="notes" color="primary"/>
        </template>
      </q-input>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step1BasicInfo'
});

/** TYPE IMPORTS */
import type { StepEmits, StepProps } from '../../interfaces/httpDataSource.interface';

/** VUE IMPORTS */
import { reactive, watch, computed } from 'vue';

/** COMPOSABLES */
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

/** PROPS & EMITS */
const props = defineProps<StepProps>();
const emit = defineEmits<StepEmits>();

/** COMPOSABLES & STORES */
const t = useHttpDataSourceCreateEditTranslations();

/** COMPUTED */
const enabledOptions = computed(() => t.basicInfo.statusOptions.value);

/** STATE */
const localData = reactive({
  name: props.modelValue.name || '',
  description: props.modelValue.description || '',
  enabled: props.modelValue.enabled ?? true,
});

watch(() => props.modelValue, (newVal) => {
  localData.name = newVal.name || '';
  localData.description = newVal.description || '';
  localData.enabled = newVal.enabled ?? true;
}, { deep: true, immediate: true });

/**
 * Emit updated values to parent component
 * Merges local form data with existing model value
 * @returns {void}
 */
function updateValue(): void {
  emit('update:modelValue', {
    ...props.modelValue,
    ...localData
  });
}
</script>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
