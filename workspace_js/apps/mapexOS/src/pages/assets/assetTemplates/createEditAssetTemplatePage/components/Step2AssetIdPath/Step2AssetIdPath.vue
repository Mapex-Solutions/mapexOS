<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step2AssetIdPathProps, Step2AssetIdPathEmits } from './interfaces/Step2AssetIdPath.interface';

defineOptions({
  name: 'Step2AssetIdPath'
});

import type { AssetTemplateData } from '../../interfaces';
import { computed } from 'vue';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';

// Props
// Emits
const props = defineProps<Step2AssetIdPathProps>();
const emit = defineEmits<Step2AssetIdPathEmits>();

// i18n
const t = useAddAssetTemplateTranslations();

// Computed v-model for data binding
const data = computed({
  get: () => props.modelValue,
  set: (value: AssetTemplateData) => emit('update:modelValue', value)
});
</script>

<template>
  <div>
    <div class="q-mb-md">
      <div class="text-subtitle1 text-weight-medium q-mb-xs">
        <q-icon name="route" color="primary" class="q-mr-xs" />
        {{ t.steps.step2.title.value }}
      </div>
      <div class="text-body2 text-grey-7">
        {{ t.steps.step2.subtitle.value }}
      </div>
    </div>

    <q-banner rounded class="bg-blue-1 text-primary q-mb-md">
      <template v-slot:avatar>
        <q-icon name="info" color="primary" />
      </template>
      <div class="text-subtitle2 text-weight-medium q-mb-xs">{{ t.steps.step2.banner.title.value }}</div>
      <div class="text-body2">
        {{ t.steps.step2.banner.description.value }}
      </div>
      <div class="text-body2 q-mt-sm">
        <strong>{{ t.steps.step2.banner.examplesTitle.value }}</strong> {{ t.steps.step2.banner.examples.value.join(', ') }}
      </div>
    </q-banner>

    <q-input
      v-model="data.assetIdPath"
      outlined
      dense
      class="rounded-borders"
      :label="t.steps.step2.fields.assetIdPath.label.value + ' *'"
      :placeholder="t.steps.step2.fields.assetIdPath.placeholder.value"
      :hint="t.steps.step2.fields.assetIdPath.hint.value"
      :rules="[
        (val) => !!val || t.steps.step2.fields.assetIdPath.required.value,
      ]"
    >
      <template v-slot:prepend>
        <q-icon name="route" color="primary" />
      </template>
    </q-input>
  </div>
</template>

<style scoped>
.rounded-borders {
  border-radius: var(--mapex-radius-md);
}
</style>
