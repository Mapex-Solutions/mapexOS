<script setup lang="ts">
defineOptions({
  name: 'WaitForNodeConfig',
});

/** TYPE IMPORTS */
import type { NodeConfigComponentProps, NodeConfigComponentEmits } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed } from 'vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { WAIT_FOR_OPERATORS } from '../constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { states } = useWorkflowContext();
const { t } = usePluginI18n('core-timers');

/** COMPUTED */

/**
 * State variable options for dropdown (from workflow state tab)
 */
const stateFields = computed(() =>
  states.value.map(v => ({
    label: `state.${v.field}`,
    value: v.field,
    type: v.type,
  })),
);

/**
 * Currently selected state variable name
 */
const fieldValue = computed({
  get: () => (props.config.field as string) || '',
  set: (val: string) => emitUpdate({ field: val }),
});

/**
 * Currently selected operator
 */
const operatorValue = computed({
  get: () => (props.config.operator as string) || 'equals',
  set: (val: string) => emitUpdate({ operator: val }),
});

/**
 * Compare-to source type (literal or variable)
 */
const compareSource = computed({
  get: () => {
    const ct = props.config.compareTo as { source?: string; value?: string } | undefined;
    return ct?.source || 'literal';
  },
  set: (val: string) => {
    const ct = props.config.compareTo as { source?: string; value?: string } | undefined;
    emitUpdate({
      compareTo: { source: val, value: ct?.value || '' },
    });
  },
});

/**
 * Compare-to value
 */
const compareValue = computed({
  get: () => {
    const ct = props.config.compareTo as { source?: string; value?: string } | undefined;
    return ct?.value || '';
  },
  set: (val: string) => {
    emitUpdate({
      compareTo: { source: compareSource.value, value: val },
    });
  },
});

/**
 * Polling interval value
 */
const intervalValue = computed({
  get: () => (props.config.interval as string) || '30s',
  set: (val: string) => emitUpdate({ interval: val }),
});

/**
 * Whether to hide the COMPARE TO section (unary operators)
 */
const hideCompare = computed(() => {
  const unary = ['isEmpty', 'isNotEmpty', 'isTrue', 'isFalse'];
  return unary.includes(operatorValue.value);
});

/**
 * Operator options with symbol display
 */
const operatorOptions = computed(() =>
  WAIT_FOR_OPERATORS.map(op => ({
    label: op.label,
    value: op.value,
    symbol: op.symbol,
  })),
);

/** FUNCTIONS */

/**
 * Emit config update with merged values
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}
</script>

<template>
  <div class="wait-for-config">
    <!-- CONDITION section -->
    <div class="wait-for-config__section">
      <div class="wait-for-config__label">
        <q-icon name="rule" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_for.config.conditionSection') }}
      </div>

      <!-- Variable selector -->
      <div class="wait-for-config__card">
        <div class="wait-for-config__field-label">{{ t('nodes.wait_for.config.variable') }}</div>
        <q-select
          v-model="fieldValue"
          outlined
          dense
          emit-value
          map-options
          :options="stateFields"
          option-label="label"
          option-value="value"
          :placeholder="t('nodes.wait_for.config.selectStateVariable')"
        >
          <template #prepend>
            <q-icon name="data_object" size="18px" />
          </template>
          <template #no-option>
            <q-item>
              <q-item-section class="text-grey-6 text-caption">
                {{ t('nodes.wait_for.config.noStateVariables') }}
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>
    </div>

    <!-- OPERATOR section -->
    <div class="wait-for-config__section">
      <div class="wait-for-config__label">
        <q-icon name="compare_arrows" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_for.config.operatorSection') }}
      </div>

      <q-select
        v-model="operatorValue"
        outlined
        dense
        emit-value
        map-options
        :options="operatorOptions"
        option-label="label"
        option-value="value"
      >
        <template #option="{ itemProps, opt }">
          <q-item v-bind="itemProps">
            <q-item-section side class="wait-for-config__operator-symbol">
              {{ opt.symbol }}
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ opt.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <!-- COMPARE TO section (hidden for unary operators) -->
    <div v-if="!hideCompare" class="wait-for-config__section">
      <div class="wait-for-config__label">
        <q-icon name="swap_horiz" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_for.config.compareToSection') }}
      </div>

      <div class="wait-for-config__card">
        <div class="wait-for-config__field-label">{{ t('nodes.wait_for.config.source') }}</div>
        <q-btn-toggle
          v-model="compareSource"
          no-caps
          dense
          unelevated
          toggle-color="primary"
          class="q-mb-sm"
          :options="[
            { label: t('nodes.wait_for.config.sourceLiteral'), value: 'literal' },
            { label: t('nodes.wait_for.config.sourceVariable'), value: 'variable' },
          ]"
        />

        <div class="wait-for-config__field-label">{{ t('nodes.wait_for.config.value') }}</div>
        <q-select
          v-if="compareSource === 'variable'"
          v-model="compareValue"
          outlined
          dense
          emit-value
          map-options
          :options="stateFields"
          option-label="label"
          option-value="value"
          :placeholder="t('nodes.wait_for.config.selectStateVariable')"
        >
          <template #prepend>
            <q-icon name="data_object" size="18px" />
          </template>
        </q-select>
        <q-input
          v-else
          v-model="compareValue"
          outlined
          dense
          :placeholder="t('nodes.wait_for.config.valuePlaceholder')"
        />
      </div>
    </div>

    <!-- POLLING section -->
    <div class="wait-for-config__section">
      <div class="wait-for-config__label">
        <q-icon name="schedule" size="14px" class="q-mr-xs" />
        {{ t('nodes.wait_for.config.timingSection') }}
      </div>

      <div class="wait-for-config__card">
        <div class="wait-for-config__field-label">{{ t('nodes.wait_for.config.pollingInterval') }}</div>
        <q-input
          v-model="intervalValue"
          outlined
          dense
          :hint="t('nodes.wait_for.config.pollingIntervalPlaceholder')"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.wait-for-config {
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

  &__operator-symbol {
    font-size: 16px;
    font-weight: 700;
    min-width: 28px;
    text-align: center;
    color: var(--mapex-text-secondary);
  }
}
</style>
