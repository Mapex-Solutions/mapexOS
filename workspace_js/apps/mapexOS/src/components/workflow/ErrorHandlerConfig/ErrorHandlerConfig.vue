<script setup lang="ts">
defineOptions({ name: 'ErrorHandlerConfig' });

/** TYPE IMPORTS */
import type { NodeErrorHandlerConfig } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** PROPS & EMITS */
const props = defineProps<{
  /** Current error handler config */
  modelValue: NodeErrorHandlerConfig;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: NodeErrorHandlerConfig): void;
}>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();

/** COMPUTED */

const unitOptions = computed(() => [
  { label: t.errorHandler.units.seconds.value, value: 'seconds' },
  { label: t.errorHandler.units.minutes.value, value: 'minutes' },
  { label: t.errorHandler.units.hours.value, value: 'hours' },
]);

/** FUNCTIONS */

function updateEnabled(val: boolean): void {
  emit('update:modelValue', { ...props.modelValue, enabled: val });
}

function updateMaxAttempts(val: string | number | null): void {
  const maxAttempts = Math.max(1, Math.min(10, Number(val) || 3));
  emit('update:modelValue', { ...props.modelValue, maxAttempts });
}

function updateInitialInterval(val: string | number | null): void {
  const initialInterval = Math.max(1, Number(val) || 5);
  emit('update:modelValue', { ...props.modelValue, initialInterval });
}

function updateIntervalUnit(val: string): void {
  emit('update:modelValue', { ...props.modelValue, intervalUnit: val });
}

function updateBackoffMultiplier(val: string | number | null): void {
  const backoffMultiplier = Math.max(1, Number(val) || 2);
  emit('update:modelValue', { ...props.modelValue, backoffMultiplier });
}
</script>

<template>
  <div class="error-handler-config">
    <!-- Banner -->
    <q-banner dense rounded class="q-mb-md text-caption" style="background: var(--mapex-surface-hover)">
      <template #avatar>
        <q-icon name="info" color="blue-6" size="20px" />
      </template>
      {{ t.errorHandler.banner.value }}
    </q-banner>

    <!-- Enable toggle -->
    <q-checkbox
      :model-value="modelValue.enabled"
      dense
      size="sm"
      class="q-mb-sm"
      @update:model-value="updateEnabled"
    >
      <span class="text-caption text-weight-medium">{{ t.errorHandler.enabled.value }}</span>
    </q-checkbox>

    <!-- Retry config (shown when enabled) -->
    <template v-if="modelValue.enabled">
      <q-separator class="q-my-sm" />

      <!-- Max Attempts -->
      <q-input
        :model-value="modelValue.maxAttempts"
        type="number"
        outlined
        dense
        hide-bottom-space
        class="q-mb-sm"
        :label="t.errorHandler.maxAttempts.value"
        :hint="t.errorHandler.maxAttemptsHint.value"
        :min="1"
        :max="10"
        @update:model-value="updateMaxAttempts"
      >
        <template #prepend>
          <q-icon name="replay" size="18px" />
        </template>
      </q-input>

      <!-- Initial Interval + Unit -->
      <div class="row q-gutter-sm q-mb-sm">
        <q-input
          :model-value="modelValue.initialInterval"
          type="number"
          outlined
          dense
          hide-bottom-space
          class="col"
          :label="t.errorHandler.initialInterval.value"
          :min="1"
          @update:model-value="updateInitialInterval"
        >
          <template #prepend>
            <q-icon name="schedule" size="18px" />
          </template>
        </q-input>

        <q-select
          :model-value="modelValue.intervalUnit"
          :options="unitOptions"
          outlined
          dense
          hide-bottom-space
          emit-value
          map-options
          class="col"
          :label="t.errorHandler.intervalUnit.value"
          @update:model-value="updateIntervalUnit"
        />
      </div>

      <!-- Backoff Multiplier -->
      <q-input
        :model-value="modelValue.backoffMultiplier"
        type="number"
        outlined
        dense
        hide-bottom-space
        :label="t.errorHandler.backoffMultiplier.value"
        :hint="t.errorHandler.backoffMultiplierHint.value"
        :min="1"
        :step="0.5"
        @update:model-value="updateBackoffMultiplier"
      >
        <template #prepend>
          <q-icon name="trending_up" size="18px" />
        </template>
      </q-input>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.error-handler-config {
  padding: 8px 0;
}
</style>
