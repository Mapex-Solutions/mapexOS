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
import type { AssetAttributeForm, AssetAttributeGeoValue } from '../../interfaces';

import { computed } from 'vue';

import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { useTS } from '@utils/translation';
import { FormReview } from '@components/forms';
import { STEP } from '../../constants';
import { resolveProtocolHealthPolicy } from '../../helpers';

const props = defineProps<Step5ReviewProps>();
const emit = defineEmits<Step5ReviewEmits>();

const tsRaw = useTS({ capitalize: false });

/**
 * Renders an attribute value for the read-only review, per its kind.
 * @param {AssetAttributeForm} attr - The attribute to format.
 * @returns {string} A display string.
 */
function attributeValueText(attr: AssetAttributeForm): string {
  switch (attr.kind) {
    case 'boolean':
      return attr.value ? 'true' : 'false';
    case 'geo': {
      const geo = attr.value as AssetAttributeGeoValue | null;
      return geo && geo.lat !== null && geo.lon !== null ? `${geo.lat}, ${geo.lon}` : '-';
    }
    case 'integer':
      return typeof attr.value === 'number' ? String(attr.value) : '-';
    case 'string':
    case 'date':
    default:
      return typeof attr.value === 'string' && attr.value !== '' ? attr.value : '-';
  }
}

const t = useAddAssetTranslations();

const previewData = computed((): ReviewSectionDef[] => {
  const data = props.modelValue;
  const { selectedTemplate, selectedRouteGroups } = props.formState;

  // Each section navigates by its step id resolved to the live 1-based position
  // in the visible list, so the edit links stay correct under any step order.
  const pos = (id: string): number => props.visibleStepIds.indexOf(id) + 1;

  // A LoRaWAN gateway has no data model and no data routing, so the template and
  // route-group steps are skipped — drop their review sections too.
  const isGateway = data.protocol === 'LORAWAN' && data.lorawanConfig?.kind === 'gateway';

  // Connection-driven protocols (MQTT broker, LoRaWAN gateway via LNS) report
  // online/offline directly, so heartbeat-absence thresholds do not apply.
  const isAlwaysMonitored = resolveProtocolHealthPolicy(data.protocol, data.lorawanConfig?.kind).alwaysEnabled;

  const sections: ReviewSectionDef[] = [
    // Identification Section
    {
      stepNumber: pos(STEP.IDENTIFICATION),
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
    // Attributes Section (only when the operator added custom fields)
    ...(data.attributes?.length
      ? [{
          stepNumber: pos(STEP.IDENTIFICATION),
          label: tsRaw('pages.assets.addAsset.steps.step1.attributes.title'),
          icon: { name: 'label', color: 'primary' },
          testId: 'review-attributes-section',
          fields: data.attributes.map((a) => ({
            label: a.label || '-',
            value: attributeValueText(a),
            type: 'text' as const,
            colSize: 6,
          })),
        }]
      : []),
    // Asset Template Section
    {
      stepNumber: pos(STEP.ASSET_TEMPLATE),
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
      stepNumber: pos(STEP.ROUTE_GROUPS),
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
      stepNumber: pos(STEP.CONNECTIVITY),
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
            'LORAWAN': 'purple'
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
        // LoRaWAN Configuration fields (only shown when protocol is LoRaWAN).
        // The kind selects the disjoint set: gateway (radio plan + auth) or
        // sensor/device (identity + region/class/activation).
        ...(data.protocol === 'LORAWAN' ? [
          {
            label: t.steps.step6.fields.lorawanKind.value,
            value: data.lorawanConfig?.kind === 'gateway'
              ? t.steps.step6.fields.lorawanKindGateway.value
              : t.steps.step6.fields.lorawanKindDevice.value,
            type: 'text' as const,
            colSize: 12,
          },
          ...(data.lorawanConfig?.kind === 'gateway'
            ? [
                {
                  label: t.steps.step6.fields.lorawanFrequencyPlan.value,
                  value: data.lorawanConfig?.gateway?.frequencyPlanId || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
                {
                  label: t.steps.step6.fields.lorawanAuthMode.value,
                  value: data.lorawanConfig?.gateway?.authMode || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
              ]
            : [
                {
                  label: t.steps.step6.fields.lorawanDevEui.value,
                  value: data.assetId || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
                {
                  label: t.steps.step6.fields.lorawanRegion.value,
                  value: data.lorawanConfig?.region || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
                {
                  label: t.steps.step6.fields.lorawanClass.value,
                  value: data.lorawanConfig?.class || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
                {
                  label: t.steps.step6.fields.lorawanActivation.value,
                  value: data.lorawanConfig?.activation || '-',
                  type: 'text' as const,
                  colSize: 6,
                },
              ]),
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
      stepNumber: pos(STEP.HEALTH_MONITORING),
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
          // Heartbeat-absence thresholds only apply to heartbeat-based liveness.
          ...(!isAlwaysMonitored ? [
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
          ] : []),
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

  if (isGateway) {
    return sections.filter(
      (s) => s.testId !== 'review-template-section' && s.testId !== 'review-routegroups-section'
    );
  }
  return sections;
});
</script>
