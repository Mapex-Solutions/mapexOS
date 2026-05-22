<template>
  <FormReview
    :sections="previewData"
    :description="t.steps.step6.subtitle.value"
    :show-success-banner="true"
    :success-message="t.steps.step6.successMessage.value"
    @edit-section="emit('editSection', $event)"
  />
</template>

<script setup lang="ts">
/** TYPE IMPORTS */
import type { Step5ReviewProps, Step5ReviewEmits } from './interfaces/Step5Review.interface';

defineOptions({
  name: 'Step5Review'
});

import type { ReviewSectionDef } from '@components/forms/review/interfaces';

import { computed } from 'vue';

import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { FormReview } from '@components/forms';

const props = defineProps<Step5ReviewProps>();
const emit = defineEmits<Step5ReviewEmits>();

const t = useAddAssetTranslations();

const previewData = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;
  const { selectedTemplate, selectedRouteGroups } = props.formState;

  return [
    // Identification Section
    {
      stepNumber: 1,
      label: t.steps.step6.sections.identification.value,
      icon: { name: 'fingerprint', color: 'primary' },
      testId: 'review-identification-section',
      fields: [
        {
          label: t.steps.step6.fields.name.value,
          value: data.name,
          type: 'text',
          colSize: 6
        },
        {
          label: t.steps.step6.fields.assetId.value,
          value: data.assetId,
          type: 'text',
          colSize: 6
        },
        {
          label: t.steps.step6.fields.status.value,
          value: data.enabled
            ? t.statusOptions.active.label.value
            : t.statusOptions.inactive.label.value,
          type: 'badge',
          badgeColors: {
            [t.statusOptions.active.label.value]: 'positive',
            [t.statusOptions.inactive.label.value]: 'negative'
          },
          colSize: 6,
        },
        {
          label: t.steps.step6.fields.description.value,
          value: data.description || '-',
          type: 'text',
          colSize: 12
        },
      ],
    },
    // Asset Template Section
    {
      stepNumber: 2,
      label: t.steps.step6.sections.assetTemplate.value,
      icon: { name: 'description', color: 'secondary' },
      testId: 'review-template-section',
      fields: [
        {
          label: t.steps.step6.fields.assetTemplate.value,
          value: selectedTemplate?.name || '-',
          type: 'text',
          colSize: 12,
        },
        {
          label: t.steps.step6.fields.manufacturer.value,
          value: selectedTemplate?.manufacturerName || '-',
          type: 'text',
          colSize: 4,
        },
        {
          label: t.steps.step6.fields.model.value,
          value: selectedTemplate?.modelName || '-',
          type: 'text',
          colSize: 4,
        },
        {
          label: t.steps.step6.fields.version.value,
          value: selectedTemplate?.version || '-',
          type: 'text',
          colSize: 4,
        },
      ],
    },
    // Route Groups Section
    {
      stepNumber: 3,
      label: t.steps.step6.sections.routeGroups.value,
      icon: { name: 'route', color: 'primary' },
      testId: 'review-routegroups-section',
      fields: [
        {
          label: t.steps.step6.fields.routeGroups.value,
          value: selectedRouteGroups.length > 0
            ? selectedRouteGroups.map(rg => rg.name).join(', ')
            : '-',
          type: 'text',
          colSize: 12,
        },
      ],
    },
    // Connectivity Section
    {
      stepNumber: 4,
      label: t.steps.step6.sections.connectivity.value,
      icon: { name: 'wifi', color: 'primary' },
      testId: 'review-connectivity-section',
      fields: [
        {
          label: t.steps.step6.fields.protocol.value,
          value: data.protocol,
          type: 'badge',
          badgeColors: {
            'HTTP': 'blue',
            'MQTT': 'green',
            'LoRaWAN': 'purple'
          },
          colSize: 12,
        },
        // MQTT Configuration fields (only shown when protocol is MQTT)
        ...(data.protocol === 'MQTT' ? [
          {
            label: t.steps.step6.fields.mqttUsername.value,
            value: data.mqttConfig?.username || '-',
            type: 'text' as const,
            colSize: 6,
          },
          {
            label: t.steps.step6.fields.mqttClientId.value,
            value: data.mqttConfig?.clientId || '-',
            type: 'text' as const,
            colSize: 6,
          },
        ] : []),
        {
          label: t.steps.step6.fields.latitude.value,
          value: data.latitude?.toString() || '-',
          type: 'text',
          colSize: 6,
        },
        {
          label: t.steps.step6.fields.longitude.value,
          value: data.longitude?.toString() || '-',
          type: 'text',
          colSize: 6,
        },
      ],
    },
    // Health Monitoring Section
    {
      stepNumber: 5,
      label: t.steps.step6.sections.healthMonitoring.value,
      icon: { name: 'monitor_heart', color: 'primary' },
      testId: 'review-health-monitoring-section',
      fields: [
        {
          label: t.steps.step6.fields.healthMonitoringEnabled.value,
          value: data.healthMonitor?.enabled
            ? t.statusOptions.active.label.value
            : t.statusOptions.inactive.label.value,
          type: 'badge',
          badgeColors: {
            [t.statusOptions.active.label.value]: 'positive',
            [t.statusOptions.inactive.label.value]: 'grey',
          },
          colSize: 12,
        },
        ...(data.healthMonitor?.enabled ? [
          {
            label: t.steps.step6.fields.threshold.value,
            value: `${data.healthMonitor.thresholdMinutes} min`,
            type: 'text' as const,
            colSize: 6,
          },
          {
            label: t.steps.step6.fields.requiredMisses.value,
            value: String(data.healthMonitor.requiredMisses),
            type: 'text' as const,
            colSize: 6,
          },
          {
            label: t.steps.step6.fields.offlineRouteGroups.value,
            value: data.healthMonitor.selectedOfflineRouteGroups?.length
              ? data.healthMonitor.selectedOfflineRouteGroups.map(rg => rg.name).join(', ')
              : '-',
            type: 'text' as const,
            colSize: 6,
          },
          {
            label: t.steps.step6.fields.onlineRouteGroups.value,
            value: data.healthMonitor.selectedOnlineRouteGroups?.length
              ? data.healthMonitor.selectedOnlineRouteGroups.map(rg => rg.name).join(', ')
              : '-',
            type: 'text' as const,
            colSize: 6,
          },
        ] : []),
      ],
    },
  ];
});
</script>
