<script setup lang="ts">
defineOptions({
  name: 'Step5Review'
});

/** TYPE IMPORTS */
import type { Trigger, TriggerFormState } from '../../interfaces';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useCreateEditTriggerTranslations } from '@composables/i18n/pages/automations/triggers/createEditTrigger';

/** LOCAL IMPORTS */
import { CATEGORY_OPTIONS, TRIGGER_TYPE_OPTIONS } from '../../constants';

/** COMPOSABLES & STORES */
const t = useCreateEditTriggerTranslations();

/** PROPS & EMITS */
const props = defineProps<{
  modelValue: Trigger;
  formState: TriggerFormState;
}>();

const emit = defineEmits<{
  'edit-section': [step: number];
}>();

/** COMPUTED */

/**
 * Get category display info
 */
const categoryInfo = computed(() => {
  return CATEGORY_OPTIONS.find((c) => c.value === props.modelValue.category);
});

/**
 * Get trigger type display info
 */
const typeInfo = computed(() => {
  return TRIGGER_TYPE_OPTIONS.find((t) => t.value === props.modelValue.triggerType);
});

/**
 * Build review sections from trigger data
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;

  return [
    // Category & Type Section
    {
      stepNumber: 1,
      label: t.steps.step5.sections.categoryType.value,
      icon: { name: 'category', color: 'primary' },
      fields: [
        {
          label: t.steps.step5.fields.category.value,
          value: categoryInfo.value?.label || data.category,
          type: 'chip',
          icon: categoryInfo.value?.emoji === '📡' ? 'router' : 'chat',
          badgeColors: categoryInfo.value?.emoji === '📡' ? 'blue-6' : 'purple-6',
          colSize: 6,
        },
        {
          label: t.steps.step5.fields.triggerType.value,
          value: typeInfo.value?.label || data.triggerType,
          type: 'chip',
          ...(typeInfo.value?.icon ? { icon: typeInfo.value.icon } : {}),
          badgeColors: 'primary',
          colSize: 6,
        },
      ],
    },
    // Basic Information Section
    {
      stepNumber: 3,
      label: t.steps.step5.sections.basicInfo.value,
      icon: { name: 'info', color: 'primary' },
      fields: [
        {
          label: t.steps.step5.fields.name.value,
          value: data.name,
          type: 'text',
          colSize: 12,
        },
        {
          label: t.steps.step5.fields.description.value,
          value: data.description,
          type: 'text',
          colSize: 12,
        },
        {
          label: t.steps.step5.fields.status.value,
          value: data.enabled,
          type: 'boolean',
          colSize: 6,
        },
      ],
    },
    // Configuration Section
    {
      stepNumber: 4,
      label: `${data.triggerType.toUpperCase()} ${t.steps.step5.sections.configuration.value}`,
      icon: { name: 'settings', color: 'primary' },
      fields: [
        {
          label: t.steps.step5.fields.configuration.value,
          value: data.config,
          type: 'json',
          colSize: 12,
        },
      ],
    },
  ];
});
</script>

<template>
  <FormReview
    :sections="reviewSections"
    :description="t.steps.step5.subtitle.value"
    :show-success-banner="true"
    :success-message="t.steps.step5.successMessage.value"
    @edit-section="emit('edit-section', $event)"
  />
</template>
