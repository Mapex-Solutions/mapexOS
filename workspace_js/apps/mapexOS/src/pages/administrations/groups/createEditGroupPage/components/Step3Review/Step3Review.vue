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
  name: 'Step4Review',
});

/** TYPE IMPORTS */
import type { Step3ReviewProps, Step3ReviewEmits } from './interfaces/Step3Review.interface';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useGroupsTranslations } from '@composables/i18n';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<Step3ReviewProps>(), {
  initialMembersCount: 0,
  isEditMode: false,
});
const emit = defineEmits<Step3ReviewEmits>();

/** COMPOSABLES & STORES */
const t = useGroupsTranslations();

/** COMPUTED */

/**
 * Effective members count
 * In EDIT mode, use initialMembersCount if user hasn't visited Step 2 (selectedMembers empty)
 * Otherwise use selectedMembers.length
 */
const effectiveMembersCount = computed(() => {
  // In EDIT mode, if no members selected yet, use initial count from API
  if (props.isEditMode && props.selectedMembers.length === 0) {
    return props.initialMembersCount;
  }
  return props.selectedMembers.length;
});

/**
 * Get enabled status display
 */
const enabledDisplay = computed(() => {
  return props.modelValue.enabled
    ? t.reviewStep.values.enabled.value
    : t.reviewStep.values.disabled.value;
});

/**
 * Build review sections from group data
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;
  const rolesCount = props.selectedRoles.length;

  return [
    // Basic Information Section
    {
      stepNumber: STEP.BASIC_INFO,
      label: t.reviewStep.sections.basicInfo.value,
      icon: { name: 'groups', color: 'primary' },
      fields: [
        {
          label: t.reviewStep.fields.name.value,
          value: data.name || t.reviewStep.values.notProvided.value,
          type: 'text',
          colSize: 6,
        },
        {
          label: t.reviewStep.fields.status.value,
          value: enabledDisplay.value,
          type: 'chip',
          badgeColors: data.enabled ? 'positive' : 'grey',
          icon: data.enabled ? 'check_circle' : 'cancel',
          colSize: 6,
        },
        {
          label: t.reviewStep.fields.description.value,
          value: data.description || t.reviewStep.values.notProvided.value,
          type: 'text',
          colSize: 12,
        },
      ],
    },
    // Roles Section
    {
      stepNumber: STEP.ROLES,
      label: t.reviewStep.sections.roles.value,
      icon: { name: 'admin_panel_settings', color: 'primary' },
      fields: [
        {
          label: t.reviewStep.fields.totalRoles.value,
          value: `${rolesCount} ${rolesCount === 1 ? t.labels.role.value : t.labels.roles.value}`,
          type: 'chip',
          badgeColors: rolesCount > 0 ? 'primary' : 'grey',
          icon: 'admin_panel_settings',
          colSize: 12,
        },
      ],
    },
    // Members Section
    {
      stepNumber: STEP.MEMBERS,
      label: t.reviewStep.sections.members.value,
      icon: { name: 'person_add', color: 'primary' },
      fields: [
        {
          label: t.reviewStep.fields.totalMembers.value,
          value: `${effectiveMembersCount.value} ${effectiveMembersCount.value === 1 ? t.labels.member.value : t.labels.members.value}`,
          type: 'chip',
          badgeColors: effectiveMembersCount.value > 0 ? 'primary' : 'grey',
          icon: 'group',
          colSize: 12,
        },
      ],
    },
  ];
});
</script>
