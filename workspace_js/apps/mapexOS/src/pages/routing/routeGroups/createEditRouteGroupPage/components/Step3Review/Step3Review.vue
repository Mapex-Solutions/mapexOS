<script setup lang="ts">
defineOptions({
  name: 'Step3Review'
});

/** TYPE IMPORTS */
import type { Step3ReviewProps } from './interfaces/Step3Review.interface';
import type { RouterFormState } from '../../interfaces';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useRouteGroupsTranslations } from '@composables/i18n/pages/routing/routeGroups/useRouteGroupsTranslations';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<Step3ReviewProps>();

const emit = defineEmits<{
  'edit-section': [step: number];
}>();

/** COMPOSABLES & STORES */
const t = useRouteGroupsTranslations();

/** FUNCTIONS */

/**
 * Get router kind label from translations
 * @param {string} kind - Router kind value
 * @returns {string} Router kind label
 */
function getRouterKindLabel(kind: string): string {
  const kindMap: Record<string, string> = {
    'lake_house': t.routerKinds.lake_house.label.value,
    'notification': t.routerKinds.notification.label.value,
    'save_event': t.routerKinds.save_event.label.value,
    'workflow': t.routerKinds.workflow.label.value,
  };
  return kindMap[kind] || kind;
}

/**
 * Format a single router for display
 * @param {RouterFormState} router - Router form state
 * @returns {string} Formatted router string
 */
function formatRouter(router: RouterFormState): string {
  const kindLabel = getRouterKindLabel(router.kind);
  let destination = '';

  if (router.kind === 'lake_house' && router.lakeHouse) {
    destination = router.lakeHouseName || router.lakeHouse.lakeHouseId || '';
  } else if (router.kind === 'notification' && router.notification) {
    destination = router.notificationName || router.notification.notificationId || '';
  } else if (router.kind === 'save_event') {
    destination = 'Save Event';
  } else if (router.kind === 'workflow' && router.workflow) {
    const instanceId = router.workflow.data?.instanceId as string | undefined;
    destination = router.workflow.mode + (instanceId ? ` (${instanceId.substring(0, 8)}...)` : '');
  }

  const conditional = router.hasConditionalRouting ? ' (Conditional)' : '';
  return `${kindLabel}${destination ? ' - ' + destination : ''}${conditional}`;
}

/** COMPUTED */

/**
 * Build router fields for the review section
 */
const routerFields = computed(() => {
  if (props.routerForms.length === 0) {
    return [{
      label: 'Routers',
      value: t.createEdit.reviewStep.values.notConfigured.value,
      type: 'text' as const,
      colSize: 12,
    }];
  }

  return props.routerForms.map((router, index) => ({
    label: `Router #${index + 1}`,
    value: formatRouter(router),
    type: 'text' as const,
    colSize: 12,
  }));
});

/**
 * Build review sections from route group data
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.formData;

  return [
    // Basic Information Section
    {
      stepNumber: STEP.BASIC_INFO,
      label: t.createEdit.reviewStep.sections.basicInfo.value,
      icon: { name: 'info', color: 'primary' },
      fields: [
        {
          label: t.createEdit.reviewStep.fields.name.value,
          value: data.name,
          type: 'text',
          colSize: 6,
        },
        {
          label: t.createEdit.reviewStep.fields.status.value,
          value: data.enabled,
          type: 'boolean',
          colSize: 6,
        },
        {
          label: t.createEdit.reviewStep.fields.templateSource.value,
          value: data.isTemplate
            ? t.createEdit.reviewStep.values.shared.value
            : t.createEdit.reviewStep.values.local.value,
          type: 'chip',
          badgeColors: data.isTemplate ? 'orange-6' : 'green-6',
          colSize: 6,
        },
        {
          label: t.createEdit.reviewStep.fields.description.value,
          value: data.description || t.createEdit.reviewStep.values.none.value,
          type: 'text',
          colSize: 12,
        },
      ],
    },
    // Routers Section
    {
      stepNumber: STEP.ROUTERS_CONFIG,
      label: t.createEdit.reviewStep.sections.routers.value.replace('{count}', props.routerForms.length.toString()),
      icon: { name: 'route', color: 'primary' },
      fields: routerFields.value,
    },
  ];
});
</script>

<template>
  <FormReview
    :sections="reviewSections"
    :description="t.createEdit.reviewStep.subtitle.value"
    :show-success-banner="true"
    :success-message="t.createEdit.reviewStep.successMessage.value"
    @edit-section="emit('edit-section', $event)"
  />
</template>
