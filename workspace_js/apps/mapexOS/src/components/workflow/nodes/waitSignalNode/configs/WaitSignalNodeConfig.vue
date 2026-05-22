<script setup lang="ts">
defineOptions({
  name: 'WaitSignalNodeConfig',
});

/** TYPE IMPORTS */
import type { NodeConfigComponentProps, NodeConfigComponentEmits } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { states, externalSignals } = useWorkflowContext();
const { t } = usePluginI18n('core-timers');

/** COMPUTED */

/**
 * State variable options for mapping dropdown
 */
const stateFields = computed(() =>
  states.value.map(v => ({
    label: `state.${v.field}`,
    value: v.field,
  })),
);

/**
 * External signal options for signal name dropdown
 */
const signalOptions = computed(() =>
  externalSignals.value.map(s => ({
    label: s.name,
    value: s.name,
  })),
);

/**
 * Signal name value
 */
const signalName = computed({
  get: () => (props.config.signalName as string) || '',
  set: (val: string) => emitUpdate({ signalName: val }),
});

/**
 * Variable mappings array
 */
const mappings = computed({
  get: () => (props.config.mappings as Array<{ from: string; to: string }>) || [],
  set: (val: Array<{ from: string; to: string }>) => emitUpdate({ mappings: val }),
});

/** FUNCTIONS */

/**
 * Emit config update with merged values
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}

/**
 * Add a new empty mapping row
 */
function addMapping(): void {
  const current = [...mappings.value];
  current.push({ from: '', to: '' });
  emitUpdate({ mappings: current });
}

/**
 * Remove a mapping row by index
 *
 * @param {number} index - Mapping index to remove
 */
function removeMapping(index: number): void {
  const current = [...mappings.value];
  current.splice(index, 1);
  emitUpdate({ mappings: current });
}

/**
 * Update a mapping's "from" value
 *
 * @param {number} index - Mapping index
 * @param {string} value - New from value
 */
function updateMappingFrom(index: number, value: string): void {
  const current = [...mappings.value];
  const existing = current[index];
  if (!existing) return;
  current[index] = { from: value, to: existing.to };
  emitUpdate({ mappings: current });
}

/**
 * Update a mapping's "to" value
 *
 * @param {number} index - Mapping index
 * @param {string} value - New to value
 */
function updateMappingTo(index: number, value: string): void {
  const current = [...mappings.value];
  const existing = current[index];
  if (!existing) return;
  current[index] = { from: existing.from, to: value };
  emitUpdate({ mappings: current });
}
</script>

<template>
  <div class="wait-signal-config">
    <!-- SIGNAL section -->
    <div class="wait-signal-config__section">
      <div class="wait-signal-config__label">
        <q-icon name="sensors" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_signal.config.signalSection') }}
      </div>

      <q-select
        :model-value="signalName"
        outlined
        dense
        emit-value
        map-options
        :options="signalOptions"
        option-label="label"
        option-value="value"
        :placeholder="t('nodes.wait_signal.config.signalNamePlaceholder')"
        :hint="t('nodes.wait_signal.config.signalNameHint')"
        @update:model-value="signalName = String($event)"
      >
        <template #prepend>
          <q-icon name="notifications_active" size="18px" />
        </template>
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey-6 text-caption">
              {{ t('nodes.wait_signal.config.noSignalsDefined') }}
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <!-- VARIABLE MAPPINGS section -->
    <div class="wait-signal-config__section">
      <div class="wait-signal-config__label">
        <q-icon name="sync_alt" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_signal.config.variableMappingsSection') }}
        <q-icon
          name="info"
          size="14px"
          color="grey-6"
          class="q-ml-xs cursor-pointer"
        >
          <AppTooltip :content="t('nodes.wait_signal.config.variableMappingsHint')" max-width="280px" />
        </q-icon>
      </div>

      <!-- Mapping rows -->
      <div v-if="mappings.length > 0" class="wait-signal-config__card">
        <div
          v-for="(mapping, index) in mappings"
          :key="`mapping-${index}`"
          class="wait-signal-config__mapping-row"
        >
          <!-- FROM -->
          <div class="wait-signal-config__mapping-field">
            <div class="wait-signal-config__field-label">{{ t('nodes.wait_signal.config.fromLabel') }}</div>
            <q-input
              :model-value="mapping.from"
              outlined
              dense
              :placeholder="t('nodes.wait_signal.config.fromPlaceholder')"
              @update:model-value="updateMappingFrom(index, String($event))"
            />
          </div>

          <!-- Arrow -->
          <q-icon name="arrow_forward" size="18px" color="grey-6" class="wait-signal-config__mapping-arrow" />

          <!-- TO -->
          <div class="wait-signal-config__mapping-field">
            <div class="wait-signal-config__field-label">{{ t('nodes.wait_signal.config.toLabel') }}</div>
            <q-select
              :model-value="mapping.to"
              outlined
              dense
              emit-value
              map-options
              :options="stateFields"
              option-label="label"
              option-value="value"
              :placeholder="t('nodes.wait_signal.config.toPlaceholder')"
              @update:model-value="updateMappingTo(index, String($event))"
            >
              <template #no-option>
                <q-item>
                  <q-item-section class="text-grey-6 text-caption">
                    {{ t('nodes.wait_signal.config.noStateVariables') }}
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <!-- Delete -->
          <q-btn
            flat
            dense
            round
            icon="close"
            size="sm"
            color="negative"
            class="wait-signal-config__mapping-delete"
            @click="removeMapping(index)"
          >
            <AppTooltip :content="t('nodes.wait_signal.config.removeMapping')" />
          </q-btn>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="wait-signal-config__empty">
        <q-icon name="info" size="16px" color="grey-6" class="q-mr-xs" />
        <span class="text-caption text-grey-6">
          {{ t('nodes.wait_signal.config.noMappingsHint') }}
        </span>
      </div>

      <!-- Add mapping button -->
      <q-btn
        flat
        dense
        no-caps
        color="primary"
        icon="add"
        :label="t('nodes.wait_signal.config.addMapping')"
        class="q-mt-sm"
        @click="addMapping"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.wait-signal-config {
  &__section {
    margin-bottom: 16px;
  }

  &__label {
    display: flex;
    align-items: center;
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    letter-spacing: 0.8px;
    margin-bottom: 8px;
  }

  &__card {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-bg);
    padding: 12px;
  }

  &__field-label {
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    letter-spacing: 0.5px;
    margin-bottom: 4px;

    &:not(:first-child) {
      margin-top: 10px;
    }
  }

  &__mapping-row {
    display: flex;
    align-items: flex-end;
    gap: 8px;
    padding-bottom: 10px;
    margin-bottom: 10px;
    border-bottom: 1px solid var(--mapex-divider);

    &:last-child {
      border-bottom: none;
      margin-bottom: 0;
      padding-bottom: 0;
    }
  }

  &__mapping-field {
    flex: 1;
    min-width: 0;
  }

  &__mapping-arrow {
    flex-shrink: 0;
    margin-bottom: 8px;
  }

  &__mapping-delete {
    flex-shrink: 0;
    margin-bottom: 4px;
  }

  &__empty {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-bg);
    border: 1px dashed var(--mapex-card-border);
  }
}
</style>
