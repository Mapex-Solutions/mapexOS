<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step3PreprocessorScriptProps, Step3PreprocessorScriptEmits } from './interfaces/Step3PreprocessorScript.interface';

defineOptions({
  name: 'Step3PreprocessorScript'
});

import { computed } from 'vue';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { MonacoScriptEditor } from '../MonacoScriptEditor';

// Props
// Emits
const props = defineProps<Step3PreprocessorScriptProps>();
const emit = defineEmits<Step3PreprocessorScriptEmits>();

// i18n
const t = useAddAssetTemplateTranslations();

// Computed v-model for scriptProcessor
const scriptProcessor = computed({
  get: () => props.modelValue.scriptProcessor || '',
  set: (value: string) => {
    const updated = { ...props.modelValue, scriptProcessor: value };
    emit('update:modelValue', updated);
  }
});

// Guidelines configuration
const guidelines = computed(() => ({
  title: t.steps.step3.guidelines.title.value,
  availableVariables: t.steps.step3.guidelines.availableVariables.value,
  items: [
    {
      label: t.steps.step3.guidelines.payloadVar.value,
      description: ''
    },
    {
      label: t.steps.step3.guidelines.consoleLog.value,
      description: ''
    }
  ],
  expectedReturn: t.steps.step3.guidelines.expectedReturn.value,
  expectedReturnValue: t.steps.step3.guidelines.expectedReturnValue.value,
}));
</script>

<template>
  <MonacoScriptEditor
    v-model="scriptProcessor"
    icon="code"
    :title="t.steps.step3.title.value"
    :subtitle="t.steps.step3.subtitle.value"
    :guidelines="guidelines"
  />
</template>
