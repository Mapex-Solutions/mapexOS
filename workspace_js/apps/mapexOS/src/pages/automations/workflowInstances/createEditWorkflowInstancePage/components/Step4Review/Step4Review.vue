<script setup lang="ts">
defineOptions({ name: 'Step4Review' });

/** TYPE IMPORTS */
import type { WorkflowInstanceFormData, WorkflowInstanceFormState, ExternalInputDefinition } from '../../interfaces';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useCreateEditWorkflowInstanceTranslations } from '@src/composables/i18n/pages/automations/workflowInstances/createEditWorkflowInstancePage/useCreateEditWorkflowInstanceTranslations';

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: WorkflowInstanceFormData;
  formState: WorkflowInstanceFormState;
}>();
const emit = defineEmits<{
  (e: 'edit-section', step: number): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowInstanceTranslations();

/** COMPUTED */

const inputDefinitions = computed((): ExternalInputDefinition[] => {
  const def = props.modelValue.selectedDefinition as any;
  if (!def?.externalInputs || !Array.isArray(def.externalInputs)) return [];
  return def.externalInputs as ExternalInputDefinition[];
});

/**
 * Review sections following the FormReview pattern
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;
  const selectedDef = data.selectedDefinition;

  const sections: ReviewSectionDef[] = [
    // Section 1: Identification
    {
      stepNumber: 1,
      label: t.review.identification.value,
      icon: { name: 'badge', color: 'primary' },
      fields: [
        {
          label: t.review.name.value,
          value: data.name || '—',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.review.enabled.value,
          value: data.enabled ? t.review.yes.value : t.review.no.value,
          type: 'badge',
          badgeColors: {
            [t.review.yes.value]: 'positive',
            [t.review.no.value]: 'negative',
          },
          colSize: 6,
        },
        {
          label: t.review.description.value,
          value: data.description || '—',
          type: 'text',
          colSize: 12,
        },
        {
          label: t.review.isTemplate.value,
          value: data.isTemplate ? t.review.yes.value : t.review.no.value,
          type: 'badge',
          badgeColors: {
            [t.review.yes.value]: 'blue',
            [t.review.no.value]: 'grey',
          },
          colSize: 6,
        },
      ],
    },

    // Section 2: Definition
    {
      stepNumber: 2,
      label: t.review.definition.value,
      icon: { name: 'account_tree', color: 'primary' },
      fields: [
        {
          label: t.review.name.value,
          value: selectedDef?.name || '—',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.review.version.value,
          value: selectedDef ? `v${selectedDef.definitionVersion}` : '—',
          type: 'badge',
          badgeColors: 'blue',
          colSize: 6,
        },
        {
          label: t.review.status.value,
          value: selectedDef ? getHealthLabel(String(selectedDef.status || 'valid')) : '—',
          type: 'badge',
          badgeColors: {
            [t.review.healthValid.value]: 'positive',
            [t.review.healthPluginMissing.value]: 'warning',
            [t.review.healthInvalid.value]: 'negative',
          },
          colSize: 6,
        },
      ],
    },

    // Section 3: External Inputs
    {
      stepNumber: 3,
      label: t.review.externalInputs.value,
      icon: { name: 'input', color: 'primary' },
      fields: inputDefinitions.value.length > 0
        ? inputDefinitions.value.map((inputDef) => ({
            label: inputDef.label + (inputDef.required ? ' *' : ''),
            value: getInputDisplayValue(inputDef),
            type: 'text' as const,
            colSize: 6,
          }))
        : [
            {
              label: t.review.externalInputs.value,
              value: t.review.noInputs.value,
              type: 'text' as const,
              colSize: 12,
            },
          ],
    },
  ];

  return sections;
});

/** FUNCTIONS */

/**
 * Get translated health label
 * @param {string} status - Raw status value
 * @returns {string} Translated label
 */
function getHealthLabel(status: string): string {
  if (status === 'plugin_missing') return t.review.healthPluginMissing.value;
  if (status === 'invalid') return t.review.healthInvalid.value;
  return t.review.healthValid.value;
}

/**
 * Get display value for an external input
 * @param {ExternalInputDefinition} inputDef - Input definition
 * @returns {string} Display value
 */
function getInputDisplayValue(inputDef: ExternalInputDefinition): string {
  const val = props.modelValue.externalInputs[inputDef.field];
  if (val === undefined || val === null || val === '') return '—';
  if (typeof val === 'boolean') return val ? t.review.yes.value : t.review.no.value;
  return String(val);
}
</script>

<template>
  <FormReview
    :sections="reviewSections"
    :show-success-banner="true"
    :success-message="t.notifications.created.value"
    @edit-section="emit('edit-section', $event)"
  />
</template>
