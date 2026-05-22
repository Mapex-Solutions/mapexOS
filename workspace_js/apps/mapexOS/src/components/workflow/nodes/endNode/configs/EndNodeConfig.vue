<script setup lang="ts">
defineOptions({
  name: 'EndNodeConfig',
});

/** TYPE IMPORTS */
import type {
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
  FieldSourceValue,
} from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { FieldSourceSelector } from '@components/forms/fieldSourceSelector';

/** COMPOSABLES */
import { usePluginI18n } from '@src/composables/workflow';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { t } = usePluginI18n('core-flow-control');

/** COMPUTED */

/**
 * Whether the node terminates with error
 */
const terminateWithError = computed<boolean>(
  () => (props.config.terminateWithError as boolean) ?? false,
);

/**
 * Error code string
 */
const errorCode = computed<string>(
  () => (props.config.errorCode as string) ?? '',
);

/**
 * Error message as FieldSourceValue
 */
const errorMessage = computed<FieldSourceValue>(() => {
  const raw = props.config.errorMessage as FieldSourceValue | undefined;
  return {
    type: raw?.type ?? 'literal',
    value: raw?.value ?? '',
    ...(raw?.nodeId != null && { nodeId: raw.nodeId }),
    ...(raw?.mode != null && { mode: raw.mode }),
  };
});

/** FUNCTIONS */

/**
 * Emit config update with partial merge
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}

/**
 * Toggle error termination mode
 *
 * @param {boolean} value - Whether to terminate with error
 */
function updateTerminateWithError(value: boolean): void {
  emitUpdate({ terminateWithError: value });
}

/**
 * Update error code
 *
 * @param {string} value - Error code string
 */
function updateErrorCode(value: string): void {
  emitUpdate({ errorCode: value });
}

/**
 * Handle error message source update from FieldSourceSelector
 *
 * @param {FieldSourceValue} value - Updated error message source
 */
function handleErrorMessageUpdate(value: FieldSourceValue): void {
  emitUpdate({ errorMessage: value });
}
</script>

<template>
  <div class="end-config">
    <!-- Info banner -->
    <div class="end-config__banner q-mb-md">
      <q-icon
        :name="terminateWithError ? 'error' : 'check_circle'"
        :color="terminateWithError ? 'negative' : 'positive'"
        size="xs"
        class="q-mr-sm"
      />
      <span class="text-caption" style="color: var(--mapex-text-secondary);">
        {{ terminateWithError
          ? t('nodes.end.config.terminateWithErrorBanner')
          : t('nodes.end.config.terminateSuccessBanner')
        }}
      </span>
    </div>

    <!-- Termination Mode -->
    <div class="end-config__section">
      <div class="end-config__section-label">{{ t('nodes.end.config.terminationModeSection') }}</div>

      <div class="row items-center justify-between">
        <div>
          <div class="text-body2 text-weight-medium" style="color: var(--mapex-text-primary);">
            {{ t('nodes.end.config.terminateWithError') }}
          </div>
          <div class="text-caption" style="color: var(--mapex-text-secondary);">
            {{ t('nodes.end.config.terminateWithErrorHint') }}
          </div>
        </div>
        <q-toggle
          :model-value="terminateWithError"
          color="negative"
          @update:model-value="updateTerminateWithError"
        />
      </div>
    </div>

    <!-- Error Configuration (only when error mode is enabled) -->
    <template v-if="terminateWithError">
      <div class="end-config__section">
        <div class="end-config__section-label">{{ t('nodes.end.config.errorCodeSection') }}</div>

        <q-input
          :model-value="errorCode"
          outlined
          dense
          :placeholder="t('nodes.end.config.errorCodePlaceholder')"
          :hint="t('nodes.end.config.errorCodeHint')"
          @update:model-value="(val: string | number | null) => updateErrorCode(String(val ?? ''))"
        >
          <template #prepend>
            <q-icon name="label" color="negative" size="xs" />
          </template>
        </q-input>
      </div>

      <div class="end-config__section">
        <div class="end-config__section-label">{{ t('nodes.end.config.errorMessageSection') }}</div>

        <FieldSourceSelector
          :model-value="errorMessage"
          :allowed-types="['literal', 'input', 'state', 'event', 'assetStatus']"
          @update:model-value="handleErrorMessageUpdate"
        />
      </div>

      <!-- Compensation hint -->
      <div class="end-config__hint">
        <q-icon name="info" color="blue-6" size="xs" class="q-mr-sm" />
        <span class="text-caption" style="color: var(--mapex-text-secondary);">
          {{ t('nodes.end.config.compensationHint') }}
        </span>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.end-config {
  &__banner {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-elevated);
    border: 1px solid var(--mapex-card-border);
  }

  &__section {
    margin-bottom: 16px;
  }

  &__section-label {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--mapex-text-secondary);
    margin-bottom: 6px;
    text-transform: uppercase;
  }

  &__hint {
    display: flex;
    align-items: flex-start;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    border: 1px solid var(--mapex-wf-tint-border);
  }
}
</style>
