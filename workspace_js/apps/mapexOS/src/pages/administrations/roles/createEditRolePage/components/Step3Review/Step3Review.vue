<template>
  <FormReview
    :sections="reviewSections"
    :description="t.reviewStep.subtitle.value"
    :show-success-banner="true"
    :success-message="isEditMode ? t.reviewStep.successMessageEdit.value : t.reviewStep.successMessage.value"
    @edit-section="emit('editSection', $event)"
  />
</template>

<script setup lang="ts">
defineOptions({
  name: 'Step3Review',
});

/** TYPE IMPORTS */
import type { Step3ReviewProps, Step3ReviewEmits } from './interfaces/Step3Review.interface';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useRolesTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<Step3ReviewProps>(), {
  isEditMode: false,
});
const emit = defineEmits<Step3ReviewEmits>();

/** COMPOSABLES & STORES */
const t = useRolesTranslations();

/** COMPUTED */

/**
 * Get total selected permissions count
 */
const selectedPermissionsCount = computed(() => {
  let count = 0;
  props.resourcePermissions.forEach(resource => {
    resource.actions.forEach(action => {
      if (action.granted) count++;
    });
  });
  return count;
});

/**
 * Get enabled resources summary
 */
const enabledResourcesSummary = computed(() => {
  const enabledResources = props.resourcePermissions
    .filter(r => r.enabled)
    .map(r => r.label);

  if (enabledResources.length === 0) {
    return t.reviewStep.values.noPermissions.value;
  }

  return enabledResources.join(', ');
});

/**
 * Get scope display value
 */
const scopeDisplay = computed(() => {
  const scopeMap: Record<string, string> = {
    global: t.scopeOptions.global.value,
    local: t.scopeOptions.local.value,
  };
  return props.modelValue.scope ? scopeMap[props.modelValue.scope] || props.modelValue.scope : t.reviewStep.values.notSelected.value;
});

/**
 * Build review sections from role data
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;

  return [
    // Basic Information Section
    {
      stepNumber: STEP.BASIC_INFO,
      label: t.reviewStep.sections.basicInfo.value,
      icon: { name: 'badge', color: 'primary' },
      fields: [
        {
          label: t.reviewStep.fields.name.value,
          value: data.name || t.reviewStep.values.notProvided.value,
          type: 'text',
          colSize: 6,
        },
        {
          label: t.reviewStep.fields.scope.value,
          value: scopeDisplay.value,
          type: 'chip',
          badgeColors: data.scope === 'global' ? 'purple' : 'orange',
          icon: data.scope === 'global' ? 'public' : 'place',
          colSize: 6,
        },
        {
          label: t.reviewStep.fields.description.value,
          value: data.description || t.reviewStep.values.notProvided.value,
          type: 'text',
          colSize: 12,
        },
        {
          label: t.reviewStep.fields.isTemplate.value,
          value: data.isTemplate,
          type: 'boolean',
          colSize: 6,
        },
      ],
    },
    // Permissions Section
    {
      stepNumber: STEP.PERMISSIONS,
      label: t.reviewStep.sections.permissions.value,
      icon: { name: 'vpn_key', color: 'primary' },
      fields: [
        {
          label: t.reviewStep.fields.totalPermissions.value,
          value: `${selectedPermissionsCount.value} ${selectedPermissionsCount.value === 1 ? t.labels.permission.value : t.labels.permissions.value}`,
          type: 'chip',
          badgeColors: selectedPermissionsCount.value > 0 ? 'primary' : 'grey',
          icon: 'check_circle',
          colSize: 6,
        },
        {
          label: t.reviewStep.fields.enabledResources.value,
          value: enabledResourcesSummary.value,
          type: 'text',
          colSize: 12,
        },
      ],
    },
  ];
});
</script>
