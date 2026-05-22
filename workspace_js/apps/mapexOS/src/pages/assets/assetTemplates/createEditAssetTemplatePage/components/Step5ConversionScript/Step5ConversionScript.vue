<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step5ConversionScriptProps, Step5ConversionScriptEmits } from './interfaces/Step5ConversionScript.interface';

defineOptions({
  name: 'Step5ConversionScript'
});

import { computed } from 'vue';
import { useAddAssetTemplateTranslations } from '@src/composables/i18n/pages/assets/addAssetTemplate/useAddAssetTemplateTranslations';
import { MonacoScriptEditor } from '../MonacoScriptEditor';

// Props
// Emits
const props = withDefaults(defineProps<Step5ConversionScriptProps>(), {
  errorMessage: ''
});

const emit = defineEmits<Step5ConversionScriptEmits>();

// i18n
const t = useAddAssetTemplateTranslations();

// Computed v-model for scriptConversion
const scriptConversion = computed({
  get: () => props.modelValue.scriptConversion || '',
  set: (value: string) => {
    const updated = { ...props.modelValue, scriptConversion: value };
    emit('update:modelValue', updated);
  }
});

// Guidelines configuration with additional info about StandardizedPayload
const guidelines = computed(() => ({
  title: t.steps.step5.guidelines.title.value,
  availableVariables: t.steps.step5.guidelines.availableVariables.value,
  items: [
    {
      label: t.steps.step5.guidelines.payloadVar.value,
      description: ''
    }
  ],
  expectedReturn: t.steps.step5.guidelines.expectedReturn.value,
  expectedReturnValue: t.steps.step5.guidelines.expectedReturnValue.value,
  additionalInfo: `
    <div class="q-mt-md q-mb-md">
      <p class="text-weight-medium">${t.steps.step5.guidelines.standardizedPayloadFormat.value}</p>
      <p class="text-caption text-grey-7">${t.steps.step5.guidelines.standardizedPayloadDescription.value}</p>
      <pre class="standardized-payload-format">{
  eventType: string,    // Required - e.g., "sensor.reading", "device.alert"
  eventId: string,      // Required - Unique identifier for this event
  data: object,         // Required - The actual data payload
  metadata?: object,    // Optional - Additional metadata
  created: string       // Required - ISO 8601 timestamp
}</pre>
    </div>

    <p class="q-mt-sm"><strong>${t.steps.step5.guidelines.requiredFields.value}</strong></p>
    <ul>
      <li><code>${t.steps.step5.guidelines.eventType.value}</code></li>
      <li><code>${t.steps.step5.guidelines.eventId.value}</code></li>
      <li><code>${t.steps.step5.guidelines.data.value}</code></li>
      <li><code>${t.steps.step5.guidelines.created.value}</code></li>
    </ul>
    <p class="q-mt-sm"><strong>${t.steps.step5.guidelines.optionalFields.value}</strong></p>
    <ul>
      <li><code>${t.steps.step5.guidelines.metadata.value}</code></li>
    </ul>
  `
}));
</script>

<template>
  <div>
    <!-- Important Banner -->
    <q-banner dense rounded class="bg-blue-1 text-primary q-mb-md">
      <template v-slot:avatar>
        <q-icon name="info" color="primary" size="sm" />
      </template>
      <div class="text-body2">
        {{ t.steps.step5.banner.important.value }}
      </div>
    </q-banner>

    <MonacoScriptEditor
      v-model="scriptConversion"
      icon="transform"
      :title="t.steps.step5.title.value"
      :subtitle="t.steps.step5.subtitle.value"
      :guidelines="guidelines"
      :error-message="errorMessage"
    />
  </div>
</template>

<style>
.standardized-payload-format {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  background-color: var(--mapex-submenu-bg);
  padding: 12px;
  border-radius: var(--mapex-radius-xs);
  margin: 0;
  overflow-x: auto;
}
</style>
