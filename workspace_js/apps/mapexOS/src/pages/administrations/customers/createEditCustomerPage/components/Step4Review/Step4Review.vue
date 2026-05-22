<template>
  <FormReview
    :sections="previewData"
    :description="t.formDescriptions.review.value"
    :show-success-banner="true"
    :success-message="successMessage"
    @edit-section="emit('editSection', $event)"
  />
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step4Review',
});

/** TYPE IMPORTS */
import type { Step4ReviewProps, Step4ReviewEmits } from './interfaces/Step4Review.interface';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useAddCustomerTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import {
  ROLE_POLICY_OPTIONS,
  DEFAULT_SCOPE_OPTIONS,
  ORG_TYPE_CONFIG,
} from '../../constants';

const props = withDefaults(defineProps<Step4ReviewProps>(), {
  isEditMode: false,
});
const emit = defineEmits<Step4ReviewEmits>();

/** COMPOSABLES & STORES */
const t = useAddCustomerTranslations();

/** COMPUTED */

/**
 * Success message for the review banner
 */
const successMessage = computed(() =>
  props.isEditMode
    ? t.messages.reviewEditSummary.value
    : t.messages.reviewCreateSummary.value,
);

/**
 * Get role policy label from value
 *
 * @param {string} value - Role policy value
 * @returns {string} Policy label
 */
function getRolePolicyLabel(value: string): string {
  const option = ROLE_POLICY_OPTIONS.find(opt => opt.value === value);
  return option?.label || value;
}

/**
 * Get default scope label from value
 *
 * @param {string} value - Default scope value
 * @returns {string} Scope label
 */
function getDefaultScopeLabel(value: string): string {
  const option = DEFAULT_SCOPE_OPTIONS.find(opt => opt.value === value);
  return option?.label || value;
}

/**
 * The step number for access policy (varies by type)
 */
const accessPolicyStepNumber = computed(() =>
  props.typeConfig.hasAddress ? 3 : 2,
);

/**
 * Preview data sections for FormReview component
 * Maps organization data to review sections with edit functionality
 * Sections are dynamic based on org type configuration
 */
const previewData = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;
  const sections: ReviewSectionDef[] = [];

  // Basic Information Section - always present
  const basicFields = [
    {
      label: t.fields.name.value,
      value: data.name || '-',
      type: 'text' as const,
      colSize: 6,
    },
    {
      label: 'Type',
      value: ORG_TYPE_CONFIG[props.orgType].label,
      type: 'badge' as const,
      badgeColors: { [ORG_TYPE_CONFIG[props.orgType].label]: ORG_TYPE_CONFIG[props.orgType].iconColor },
      colSize: 6,
    },
  ];

  // Add phone only for types that support it
  if (props.typeConfig.hasPhone) {
    basicFields.push({
      label: t.fields.phone.value,
      value: data.phone || '-',
      type: 'text' as const,
      colSize: 6,
    });
  }

  basicFields.push({
    label: t.fields.enabled.value,
    value: data.enabled ? t.status.enabled.value : t.status.disabled.value,
    type: 'badge' as const,
    badgeColors: { [t.status.enabled.value]: 'positive', [t.status.disabled.value]: 'negative' },
    colSize: 6,
  });

  sections.push({
    stepNumber: 1,
    label: t.sections.basicInfo.value,
    icon: { name: props.typeConfig.icon, color: props.typeConfig.iconColor },
    fields: basicFields,
    testId: 'review-basic-section',
  });

  // Address Section - only for types with address
  if (props.typeConfig.hasAddress) {
    sections.push({
      stepNumber: 2,
      label: t.sections.address.value,
      icon: { name: 'location_on', color: 'secondary' },
      testId: 'review-address-section',
      fields: [
        {
          label: t.fields.country.value,
          value: data.address?.country || '-',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.fields.state.value,
          value: data.address?.state || '-',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.fields.city.value,
          value: data.address?.city || '-',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.fields.zipCode.value,
          value: data.address?.zipCode || '-',
          type: 'text',
          colSize: 6,
        },
      ],
    });
  }

  // Access Policy Section - always present
  sections.push({
    stepNumber: accessPolicyStepNumber.value,
    label: t.sections.accessPolicy.value,
    icon: { name: 'policy', color: 'primary' },
    testId: 'review-access-section',
    fields: [
      {
        label: t.fields.authProvider.value,
        value: 'Internal',
        type: 'badge',
        badgeColors: { 'Internal': 'grey-7' },
        colSize: 6,
      },
      {
        label: t.fields.rolePolicy.value,
        value: getRolePolicyLabel(data.accessPolicy?.rolePolicy || 'strict'),
        type: 'badge',
        badgeColors: { 'Strict': 'orange', 'Merge': 'blue' },
        colSize: 6,
      },
      {
        label: t.fields.defaultScope.value,
        value: getDefaultScopeLabel(data.accessPolicy?.defaultScope || 'local'),
        type: 'badge',
        badgeColors: { 'Local': 'teal', 'Recursive': 'purple' },
        colSize: 6,
      },
    ],
  });

  return sections;
});
</script>
