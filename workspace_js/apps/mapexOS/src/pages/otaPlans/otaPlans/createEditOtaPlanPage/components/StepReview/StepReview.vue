<template>
  <div class="q-gutter-md">
    <div class="text-subtitle1 text-weight-medium">{{ t.review.title.value }}</div>

    <!-- Details + window -->
    <div class="review-section">
      <div class="review-header">
        <q-icon name="schedule" color="primary" size="sm" class="q-mr-sm" />
        <span class="text-subtitle2">{{ t.review.window.value }}</span>
      </div>
      <q-separator class="q-my-sm" />
      <div class="text-body2 q-mb-sm">
        <span class="text-weight-medium">{{ modelValue.name }}</span>
        <span v-if="modelValue.description" class="text-grey-7"> — {{ modelValue.description }}</span>
      </div>
      <div class="text-body2 q-mb-sm">
        {{ modelValue.startAt ? formatDate(modelValue.startAt) : t.review.runNow.value }}
        → {{ modelValue.maxTime ? formatDate(modelValue.maxTime) : '' }}
      </div>
      <div class="row q-gutter-sm">
        <DetailChip icon="speed" color="orange" :label="`${modelValue.rolloutConfig.ratePerMinute}/min`" />
        <DetailChip icon="report" color="negative" :label="`${modelValue.rolloutConfig.abortThresholdPct}%`" />
      </div>
    </div>

    <!-- Firmware -->
    <div class="review-section">
      <div class="review-header">
        <q-icon name="memory" color="primary" size="sm" class="q-mr-sm" />
        <span class="text-subtitle2">{{ t.review.firmware.value }}</span>
      </div>
      <q-separator class="q-my-sm" />
      <div class="row q-gutter-sm">
        <DetailChip icon="new_releases" color="blue" :label="modelValue.version" />
        <DetailChip v-if="modelValue.file" icon="attach_file" color="grey" :label="modelValue.file.name" />
        <DetailChip icon="check_circle" color="positive" :label="t.firmware.ready.value" />
      </div>
    </div>

    <!-- Devices -->
    <div class="review-section">
      <div class="review-header">
        <q-icon name="devices" color="primary" size="sm" class="q-mr-sm" />
        <span class="text-subtitle2">{{ t.review.devices.value }}</span>
      </div>
      <q-separator class="q-my-sm" />
      <DetailChip
        icon="group"
        color="indigo"
        :label="`${modelValue.selectedAssetIds.length} ${t.review.deviceCount.value}`"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'StepReview'
});

/** TYPE IMPORTS */
import type { OtaPlanFormData } from '../../interfaces/createEditOtaPlan.interface';

/** COMPONENTS */
import { DetailChip } from '@components/chips/DetailChip';

/** COMPOSABLES */
import { useCreateEditOtaPlanTranslations } from '@composables/i18n/pages/otaPlans/createEditOtaPlan/useCreateEditOtaPlanTranslations';

/** PROPS & EMITS */
defineProps<{
  modelValue: OtaPlanFormData;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditOtaPlanTranslations();

/** FUNCTIONS */

/** Locale-formatted date-time for the window summary */
function formatDate(value: string): string {
  if (!value) return '';
  return new Date(value).toLocaleString();
}
</script>

<style lang="scss" scoped>
.review-section {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  padding: var(--mapex-spacing-md);
}

.review-header {
  display: flex;
  align-items: center;
}
</style>
