<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step4ValidationScriptProps, Step4ValidationScriptEmits } from './interfaces/Step4ValidationScript.interface';

defineOptions({
  name: 'Step4ValidationScript'
});

import { computed } from 'vue';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { MonacoScriptEditor } from '../MonacoScriptEditor';

// Props
// Emits
const props = defineProps<Step4ValidationScriptProps>();
const emit = defineEmits<Step4ValidationScriptEmits>();

// i18n
const t = useAddAssetTemplateTranslations();

// Computed v-model for scriptValidator
const scriptValidator = computed({
  get: () => props.modelValue.scriptValidator || '',
  set: (value: string) => {
    const updated = { ...props.modelValue, scriptValidator: value };
    emit('update:modelValue', updated);
  }
});

// Guidelines configuration
const guidelines = computed(() => ({
  title: t.steps.step4.guidelines.title.value,
  availableVariables: t.steps.step4.guidelines.availableVariables.value,
  items: [
    {
      label: t.steps.step4.guidelines.payloadVar.value,
      description: ''
    }
  ],
  expectedReturn: t.steps.step4.guidelines.expectedReturn.value,
  expectedReturnValue: t.steps.step4.guidelines.expectedReturnValue.value,
}));
</script>

<template>
  <MonacoScriptEditor
    v-model="scriptValidator"
    icon="verified"
    :title="t.steps.step4.title.value"
    :subtitle="t.steps.step4.subtitle.value"
    :guidelines="guidelines"
  />
</template>
