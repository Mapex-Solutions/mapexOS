<script setup lang="ts">
defineOptions({
  name: 'Step6Review'
});

/** TYPE IMPORTS */
import type { Step6ReviewProps } from './interfaces/Step6Review.interface';
import type { ReviewSectionDef } from '@components/forms/review/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FormReview } from '@components/forms';

/** COMPOSABLES */
import { useHttpDataSourceCreateEditTranslations } from '@composables/i18n/pages/datasources/http';

/** LOCAL IMPORTS */
import { STEP } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<Step6ReviewProps>();

const emit = defineEmits<{
  'edit-section': [step: number];
}>();

/** COMPOSABLES & STORES */
const t = useHttpDataSourceCreateEditTranslations();

/** COMPUTED */

/**
 * Get formatted days of week string
 */
const formattedDays = computed(() => {
  if (!props.dataSource.daysOfWeek || props.dataSource.daysOfWeek.length === 0) {
    return t.review.values.notConfigured.value;
  }

  const dayNames: Record<string, string> = {
    '0': t.workingHours.daysOfWeek.sunday.value,
    '1': t.workingHours.daysOfWeek.monday.value,
    '2': t.workingHours.daysOfWeek.tuesday.value,
    '3': t.workingHours.daysOfWeek.wednesday.value,
    '4': t.workingHours.daysOfWeek.thursday.value,
    '5': t.workingHours.daysOfWeek.friday.value,
    '6': t.workingHours.daysOfWeek.saturday.value,
  };

  return props.dataSource.daysOfWeek
    .map(day => dayNames[day] || day)
    .join(', ');
});

/**
 * Get formatted time interval string
 */
const formattedTimeInterval = computed(() => {
  if (!props.dataSource.timeIntervals || props.dataSource.timeIntervals.length === 0) {
    return t.review.values.notConfigured.value;
  }

  const interval = props.dataSource.timeIntervals[0];
  if (!interval) {
    return t.review.values.notConfigured.value;
  }
  return `${interval.startTime} - ${interval.endTime}`;
});

/**
 * Get formatted rate limit type with value
 */
const formattedRateLimit = computed(() => {
  if (!props.dataSource.rateLimitType) {
    return t.review.values.notConfigured.value;
  }

  const typeLabels: Record<string, string> = {
    'second': t.review.values.perSecond.value,
    'minute': t.review.values.perMinute.value,
    'hour': t.review.values.perHour.value,
  };

  return `${props.dataSource.rateLimitValue} ${typeLabels[props.dataSource.rateLimitType] || props.dataSource.rateLimitType}`;
});

/**
 * Get formatted action on exceed
 */
const formattedActionOnExceed = computed(() => {
  if (!props.dataSource.actionOnExceed) {
    return t.review.values.notConfigured.value;
  }

  const actionLabels: Record<string, string> = {
    'drop': t.review.values.drop.value,
    'queue': t.review.values.queue.value,
  };

  return actionLabels[props.dataSource.actionOnExceed] || props.dataSource.actionOnExceed;
});

/**
 * Get formatted auth type display
 */
const formattedAuthType = computed(() => {
  if (!props.dataSource.authType || props.dataSource.authType === 'none') {
    return t.review.values.none.value;
  }

  const authLabels: Record<string, string> = {
    'apiKey': 'API Key',
    'jwt': 'JWT',
    'ip_whitelist': 'IP Whitelist',
    'oauth2': 'OAuth2',
  };

  return authLabels[props.dataSource.authType] || props.dataSource.authType;
});

/**
 * Get formatted binding mode display
 */
const formattedBindingMode = computed(() => {
  if (!props.dataSource.bindingMode) {
    return t.review.values.notConfigured.value;
  }

  const bindingLabels: Record<string, string> = {
    'fixedAssetId': t.review.values.fixedAsset.value,
    'uuidField': t.review.values.uuidField.value,
  };

  return bindingLabels[props.dataSource.bindingMode] || props.dataSource.bindingMode;
});

/**
 * Get formatted UUID paths
 */
const formattedUuidPaths = computed(() => {
  if (!props.dataSource.finalUuidPaths || props.dataSource.finalUuidPaths.length === 0) {
    return t.review.values.notConfigured.value;
  }

  return props.dataSource.finalUuidPaths.join(', ');
});

/**
 * Build review sections from data source
 */
const reviewSections = computed((): ReviewSectionDef[] => {
  const data = props.dataSource;
  const sections: ReviewSectionDef[] = [];

  // Section 1: Basic Information
  sections.push({
    stepNumber: STEP.BASIC_INFO,
    label: t.review.sections.basicInfo.value,
    icon: { name: 'info', color: 'primary' },
    fields: [
      {
        label: t.review.fields.name.value,
        value: data.name || t.review.values.notConfigured.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.status.value,
        value: data.enabled,
        type: 'boolean',
        colSize: 6,
      },
      {
        label: t.review.fields.description.value,
        value: data.description || t.review.values.noDescription.value,
        type: 'text',
        colSize: 12,
      },
    ],
  });

  // Section 2: Working Hours & Rate Limit
  const workingHoursFields: ReviewSectionDef['fields'] = [
    {
      label: t.review.fields.workingHoursEnabled.value,
      value: data.enableWorkingHours,
      type: 'boolean',
      colSize: 6,
    },
  ];

  if (data.enableWorkingHours) {
    workingHoursFields.push(
      {
        label: t.review.fields.daysOfWeek.value,
        value: formattedDays.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.timeInterval.value,
        value: formattedTimeInterval.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.timezone.value,
        value: data.timezone || t.review.values.notConfigured.value,
        type: 'text',
        colSize: 6,
      },
    );
  }

  workingHoursFields.push({
    label: t.review.fields.rateLimitEnabled.value,
    value: data.enableRateLimit,
    type: 'boolean',
    colSize: 6,
  });

  if (data.enableRateLimit) {
    workingHoursFields.push(
      {
        label: t.review.fields.rateLimitValue.value,
        value: formattedRateLimit.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.burstCapacity.value,
        value: data.burstCapacity?.toString() || t.review.values.notConfigured.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.actionOnExceed.value,
        value: formattedActionOnExceed.value,
        type: 'text',
        colSize: 6,
      },
    );
  }

  sections.push({
    stepNumber: STEP.WORKING_HOURS,
    label: t.review.sections.workingHours.value,
    icon: { name: 'schedule', color: 'primary' },
    fields: workingHoursFields,
  });

  // Section 3: Authentication
  const authFields: ReviewSectionDef['fields'] = [
    {
      label: t.review.fields.authType.value,
      value: formattedAuthType.value,
      type: 'chip',
      badgeColors: data.authType === 'none' ? 'red-6' : 'primary',
      colSize: 6,
    },
  ];

  if (data.authType === 'apiKey') {
    authFields.push(
      {
        label: t.review.fields.apiKeyHeader.value,
        value: data.apiKey?.headerApiKey || t.review.values.notConfigured.value,
        type: 'text',
        colSize: 6,
      },
      {
        label: t.review.fields.apiKeyValue.value,
        value: data.apiKey?.valueApiKey ? '••••••••' : t.review.values.notConfigured.value,
        type: 'text',
        colSize: 6,
      },
    );
  } else if (data.authType === 'jwt') {
    authFields.push({
      label: t.review.fields.jwtSecret.value,
      value: data.jwt?.secretKey ? '••••••••' : t.review.values.notConfigured.value,
      type: 'text',
      colSize: 6,
    });
  } else if (data.authType === 'ip_whitelist') {
    authFields.push({
      label: t.review.fields.ipAddresses.value,
      value: data.ipWhitelist?.addresses?.join(', ') || t.review.values.notConfigured.value,
      type: 'text',
      colSize: 12,
    });
  } else if (data.authType === 'oauth2') {
    authFields.push({
      label: t.review.fields.jwksUrl.value,
      value: data.oauth2?.jwksUrl || t.review.values.notConfigured.value,
      type: 'text',
      colSize: 12,
    });
  }

  sections.push({
    stepNumber: STEP.AUTHENTICATION,
    label: t.review.sections.authentication.value,
    icon: { name: 'lock', color: 'primary' },
    fields: authFields,
  });

  // Section 4: Asset Binding
  const bindingFields: ReviewSectionDef['fields'] = [
    {
      label: t.review.fields.bindingMode.value,
      value: formattedBindingMode.value,
      type: 'chip',
      badgeColors: 'primary',
      colSize: 6,
    },
  ];

  if (data.bindingMode === 'fixedAssetId') {
    bindingFields.push({
      label: t.review.fields.selectedAsset.value,
      value: data.directAssetId || t.review.values.notConfigured.value,
      type: 'text',
      colSize: 6,
    });
  } else if (data.bindingMode === 'uuidField') {
    bindingFields.push({
      label: t.review.fields.uuidPaths.value,
      value: formattedUuidPaths.value,
      type: 'text',
      colSize: 12,
    });
  }

  sections.push({
    stepNumber: STEP.ASSET_BINDING,
    label: t.review.sections.assetBinding.value,
    icon: { name: 'device_unknown', color: 'primary' },
    fields: bindingFields,
  });

  return sections;
});
</script>

<template>
  <FormReview
    :sections="reviewSections"
    :description="t.review.subtitle.value"
    :show-success-banner="true"
    :success-message="t.review.successMessage.value"
    @edit-section="emit('edit-section', $event)"
  />
</template>
