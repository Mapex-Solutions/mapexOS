<script setup lang="ts">
defineOptions({ name: 'TimeoutConfig' });

/** TYPE IMPORTS */
import type { NodeTimeoutConfig } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** PROPS & EMITS */
const props = defineProps<{
  /** Current timeout config */
  modelValue: NodeTimeoutConfig;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: NodeTimeoutConfig): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/** COMPUTED */

const unitOptions = computed(() => [
  { label: t.timeout.units.seconds.value, value: 'seconds' },
  { label: t.timeout.units.minutes.value, value: 'minutes' },
  { label: t.timeout.units.hours.value, value: 'hours' },
  { label: t.timeout.units.days.value, value: 'days' },
  { label: t.timeout.units.months.value, value: 'months' },
  { label: t.timeout.units.years.value, value: 'years' },
]);

/** FUNCTIONS */

/**
 * Update duration value
 * @param {string | number | null} val - New duration
 */
function updateDuration(val: string | number | null): void {
  const duration = Math.max(1, Number(val) || 30);
  emit('update:modelValue', { ...props.modelValue, duration });
}

/**
 * Update unit value
 * @param {string} val - New unit
 */
function updateUnit(val: string): void {
  emit('update:modelValue', { ...props.modelValue, unit: val });
}

/**
 * Toggle enableOutput
 * @param {boolean} val - New value
 */
function updateEnableOutput(val: boolean): void {
  emit('update:modelValue', { ...props.modelValue, enableOutput: val });
}
</script>

<template>
  <div class="timeout-config">
    <!-- Section header -->
    <div class="timeout-config__header row items-center q-mb-sm">
      <q-icon name="timer" size="18px" color="orange-6" class="q-mr-xs" />
      <span class="text-caption text-weight-medium" style="color: var(--mapex-text-secondary)">
        {{ t.timeout.sectionTitle.value }}
      </span>
    </div>

    <!-- Duration + Unit -->
    <div class="row q-gutter-sm q-mb-sm">
      <q-input
        :model-value="modelValue.duration"
        type="number"
        outlined
        dense
        hide-bottom-space
        class="col"
        :label="t.timeout.duration.value"
        :min="1"
        @update:model-value="updateDuration"
      >
        <template #prepend>
          <q-icon name="schedule" size="18px" />
        </template>
      </q-input>

      <q-select
        :model-value="modelValue.unit"
        :options="unitOptions"
        outlined
        dense
        hide-bottom-space
        emit-value
        map-options
        class="col"
        :label="t.timeout.unit.value"
        @update:model-value="updateUnit"
      />
    </div>

    <!-- Enable Output -->
    <q-checkbox
      :model-value="modelValue.enableOutput"
      dense
      size="sm"
      class="timeout-config__checkbox"
      @update:model-value="updateEnableOutput"
    >
      <span class="text-caption">{{ t.timeout.enableOutput.value }}</span>
      <q-tooltip max-width="280px">{{ t.timeout.enableOutputHint.value }}</q-tooltip>
    </q-checkbox>
  </div>
</template>

<style lang="scss" scoped>
.timeout-config {
  padding: 8px 0;
  border-top: 1px solid var(--mapex-card-border);

  &__header {
    opacity: 0.85;
  }

  &__checkbox {
    margin-left: -2px;
  }
}
</style>
