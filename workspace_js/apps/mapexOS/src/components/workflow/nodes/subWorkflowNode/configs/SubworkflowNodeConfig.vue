<script setup lang="ts">
/** TYPE IMPORTS */
import type { WorkflowSelectorItem } from '@components/dialogs/common/workflowSelectorDialog/interfaces/workflowSelectorDialog.interface';
import type { NodeConfigComponentProps, NodeConfigComponentEmits, FieldSourceValue } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { WorkflowSelectorDialog } from '@components/dialogs/common/workflowSelectorDialog';
import { FieldSourceSelector } from '@components/forms/fieldSourceSelector';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { states } = useWorkflowContext();
const { t } = usePluginI18n('core-flow-control');

/** STATE */

/**
 * Whether the workflow selector dialog is open
 */
const workflowDialogOpen = ref(false);

/** COMPUTED */

/**
 * Current execution mode from config
 */
const executionMode = computed<string>(() =>
  (props.config.executionMode as string) ?? 'sync',
);

/**
 * Current timeout config
 */
const timeout = computed<TimeoutConfig>(() =>
  (props.config.timeout as TimeoutConfig) ?? { duration: 30, unit: 'seconds' },
);

/**
 * Current input mappings from config
 */
const inputMappings = computed<InputMapping[]>(() =>
  (props.config.inputMappings as InputMapping[]) ?? [],
);

/**
 * Current output mappings from config (child output → parent state)
 */
const outputMappings = computed<OutputMapping[]>(() =>
  (props.config.outputMappings as OutputMapping[]) ?? [],
);

/**
 * Available workflow variables for source/target dropdowns
 */
const variableOptions = computed(() =>
  states.value.map(v => ({ label: v.field, value: v.field })),
);

/**
 * Whether output mappings section is relevant (only sync mode)
 */
const showOutputMappings = computed(() =>
  !!props.config.workflowId && executionMode.value === 'sync',
);

/**
 * State fields for FieldSourceSelector
 */
const stateFields = computed(() =>
  states.value.map(v => ({ name: v.field, type: v.type })),
);

/** FUNCTIONS */

/**
 * Handle workflow selection from WorkflowSelectorDialog
 *
 * @param {WorkflowSelectorItem} workflow - Selected workflow
 */
function handleWorkflowSelect(workflow: WorkflowSelectorItem): void {
  emit('update:config', {
    ...props.config,
    workflowId: workflow.id,
    workflowName: workflow.name,
  });
}

/**
 * Clear selected workflow and reset related config
 */
function clearWorkflow(): void {
  emit('update:config', {
    ...props.config,
    workflowId: undefined,
    workflowName: undefined,
    inputMappings: [],
    outputMappings: [],
  });
}

/**
 * Update execution mode (sync/async)
 *
 * @param {string} mode - New execution mode
 */
function updateExecutionMode(mode: string): void {
  const patch: Record<string, unknown> = { ...props.config, executionMode: mode };
  if (mode === 'async') {
    patch.outputMappings = [];
  }
  emit('update:config', patch);
}

/**
 * Update timeout duration value
 *
 * @param {string | number | null} val - New duration value
 */
function updateTimeoutDuration(val: string | number | null): void {
  const duration = Math.max(1, Number(val) || 30);
  emit('update:config', {
    ...props.config,
    timeout: { ...timeout.value, duration },
  });
}

/**
 * Update timeout unit
 *
 * @param {string} unit - New time unit
 */
function updateTimeoutUnit(unit: string): void {
  emit('update:config', {
    ...props.config,
    timeout: { ...timeout.value, unit },
  });
}

// ── Input Mapping helpers ───────────────────────────────────────────────

/**
 * Add a new empty input mapping row
 */
function addInputMapping(): void {
  const current = [...inputMappings.value];
  current.push({ childVariable: '', source: { type: 'literal', value: '' } });
  emit('update:config', { ...props.config, inputMappings: current });
}

/**
 * Remove an input mapping by index
 *
 * @param {number} index - Mapping index to remove
 */
function removeInputMapping(index: number): void {
  const current = [...inputMappings.value];
  current.splice(index, 1);
  emit('update:config', { ...props.config, inputMappings: current });
}

/**
 * Update an input mapping's child variable name
 *
 * @param {number} index - Mapping index
 * @param {string} value - New child variable name
 */
function updateInputChildVariable(index: number, value: string): void {
  const current = [...inputMappings.value];
  if (!current[index]) return;
  current[index] = { ...current[index], childVariable: value };
  emit('update:config', { ...props.config, inputMappings: current });
}

/**
 * Update an input mapping's source from FieldSourceSelector
 *
 * @param {number} index - Mapping index
 * @param {FieldSourceValue} source - New source value
 */
function updateInputSource(index: number, source: FieldSourceValue): void {
  const current = [...inputMappings.value];
  if (!current[index]) return;
  current[index] = { ...current[index], source };
  emit('update:config', { ...props.config, inputMappings: current });
}

// ── Output Mapping helpers ──────────────────────────────────────────────

/**
 * Add a new empty output mapping row
 */
function addOutputMapping(): void {
  const current = [...outputMappings.value];
  current.push({ outputKey: '', targetVariable: '' });
  emit('update:config', { ...props.config, outputMappings: current });
}

/**
 * Remove an output mapping by index
 *
 * @param {number} index - Mapping index to remove
 */
function removeOutputMapping(index: number): void {
  const current = [...outputMappings.value];
  current.splice(index, 1);
  emit('update:config', { ...props.config, outputMappings: current });
}

/**
 * Update an output mapping field by index
 *
 * @param {number} index - Mapping index
 * @param {string} field - Field name ('outputKey' or 'targetVariable')
 * @param {string} value - New value
 */
function updateOutputMapping(index: number, field: string, value: string): void {
  const current = [...outputMappings.value];
  if (!current[index]) return;

  current[index] = { ...current[index], [field]: value };
  emit('update:config', { ...props.config, outputMappings: current });
}
</script>

<script lang="ts">
import type { InputMapping, OutputMapping, TimeoutConfig } from '../interfaces/subWorkflowNode.interface';
export type { InputMapping, OutputMapping, TimeoutConfig };
</script>

<template>
  <div class="subworkflow-config">
    <!-- WORKFLOW SELECTOR -->
    <div class="subworkflow-config__section">
      <div class="subworkflow-config__label">{{ t('nodes.subworkflow.config.workflowSection') }}</div>

      <!-- Selected workflow display -->
      <div
        v-if="props.config.workflowName"
        class="subworkflow-config__selected"
      >
        <q-item dense class="rounded-borders">
          <q-item-section avatar>
            <q-avatar
              icon="hub"
              color="deep-purple-6"
              text-color="white"
              size="sm"
            />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-medium ellipsis">
              {{ props.config.workflowName }}
            </q-item-label>
            <q-item-label caption>
              <q-badge color="deep-purple-6" :label="t('nodes.subworkflow.config.subworkflowBadge')" dense />
            </q-item-label>
          </q-item-section>
          <q-item-section side>
            <div class="row q-gutter-xs">
              <q-btn
                flat
                dense
                round
                icon="swap_horiz"
                size="xs"
                color="primary"
                @click="workflowDialogOpen = true"
              >
                <AppTooltip :content="t('nodes.subworkflow.config.changeWorkflow')" />
              </q-btn>
              <q-btn
                flat
                dense
                round
                icon="close"
                size="xs"
                color="grey-7"
                @click="clearWorkflow"
              >
                <AppTooltip :content="t('nodes.subworkflow.config.removeWorkflow')" />
              </q-btn>
            </div>
          </q-item-section>
        </q-item>
      </div>

      <!-- Select workflow button -->
      <q-btn
        v-else
        outline
        no-caps
        dense
        color="deep-purple-6"
        icon="hub"
        :label="t('nodes.subworkflow.config.selectWorkflow')"
        class="full-width"
        @click="workflowDialogOpen = true"
      />
    </div>

    <!-- EXECUTION MODE -->
    <div v-if="props.config.workflowId" class="subworkflow-config__section">
      <div class="subworkflow-config__label">{{ t('nodes.subworkflow.config.executionModeSection') }}</div>
      <q-btn-toggle
        :model-value="executionMode"
        toggle-color="deep-purple-6"
        :options="[
          { label: t('nodes.subworkflow.config.sync'), value: 'sync', icon: 'sync' },
          { label: t('nodes.subworkflow.config.async'), value: 'async', icon: 'sync_disabled' },
        ]"
        no-caps
        dense
        unelevated
        spread
        class="subworkflow-config__toggle"
        @update:model-value="updateExecutionMode"
      />
      <div class="subworkflow-config__hint">
        {{ executionMode === 'sync'
          ? t('nodes.subworkflow.config.syncDescription')
          : t('nodes.subworkflow.config.asyncDescription')
        }}
      </div>
    </div>

    <!-- TIMEOUT -->
    <div v-if="props.config.workflowId" class="subworkflow-config__section">
      <div class="subworkflow-config__label">{{ t('nodes.subworkflow.config.executionTimeoutSection') }}</div>
      <div class="row q-col-gutter-sm">
        <div class="col">
          <q-input
            :model-value="timeout.duration"
            outlined
            dense
            type="number"
            :min="1"
            :label="t('nodes.subworkflow.config.duration')"
            @update:model-value="updateTimeoutDuration"
          >
            <template #prepend>
              <q-icon name="timer" size="xs" />
            </template>
          </q-input>
        </div>
        <div class="col-5">
          <q-select
            :model-value="timeout.unit"
            outlined
            dense
            :options="[
              { label: t('nodes.subworkflow.config.seconds'), value: 'seconds' },
              { label: t('nodes.subworkflow.config.minutes'), value: 'minutes' },
              { label: t('nodes.subworkflow.config.hours'), value: 'hours' },
            ]"
            option-label="label"
            option-value="value"
            emit-value
            map-options
            @update:model-value="updateTimeoutUnit"
          />
        </div>
      </div>
      <div class="subworkflow-config__hint">
        {{ executionMode === 'sync'
          ? t('nodes.subworkflow.config.timeoutSyncHint')
          : t('nodes.subworkflow.config.timeoutAsyncHint')
        }}
      </div>
    </div>

    <!-- INPUT MAPPINGS (Parent → Child) -->
    <div v-if="props.config.workflowId" class="subworkflow-config__section">
      <div class="subworkflow-config__label">
        {{ t('nodes.subworkflow.config.inputMappingsSection') }} ({{ inputMappings.length }})
      </div>
      <div class="subworkflow-config__hint q-mb-sm">
        {{ t('nodes.subworkflow.config.inputMappingsHint') }}
      </div>

      <div
        v-for="(mapping, index) in inputMappings"
        :key="'in_' + index"
        class="subworkflow-config__mapping-card"
      >
        <!-- Child variable name -->
        <q-input
          :model-value="mapping.childVariable"
          outlined
          dense
          :label="t('nodes.subworkflow.config.childVariable')"
          :placeholder="t('nodes.subworkflow.config.childVariablePlaceholder')"
          class="q-mb-sm"
          @update:model-value="(val: string | number | null) => updateInputChildVariable(index, String(val ?? ''))"
        >
          <template #prepend>
            <q-icon name="input" size="xs" />
          </template>
        </q-input>

        <!-- Source type + value via FieldSourceSelector -->
        <FieldSourceSelector
          :model-value="(mapping.source as FieldSourceValue)"
          :allowed-types="['literal', 'input', 'state']"
          :state-fields="stateFields"
          @update:model-value="(val: FieldSourceValue) => updateInputSource(index, val)"
        />

        <!-- Remove mapping -->
        <div class="text-right q-mt-xs">
          <q-btn
            flat
            dense
            round
            icon="delete_outline"
            size="xs"
            color="negative"
            @click="removeInputMapping(index)"
          >
            <AppTooltip :content="t('nodes.subworkflow.config.removeMapping')" />
          </q-btn>
        </div>
      </div>

      <!-- Add input mapping button -->
      <q-btn
        outline
        no-caps
        dense
        color="deep-purple-6"
        icon="add"
        :label="t('nodes.subworkflow.config.addInput')"
        class="full-width"
        @click="addInputMapping"
      />
    </div>

    <!-- OUTPUT MAPPINGS (Child output → Parent state) — Sync only -->
    <div v-if="showOutputMappings" class="subworkflow-config__section">
      <div class="subworkflow-config__label">
        {{ t('nodes.subworkflow.config.outputMappingsSection') }} ({{ outputMappings.length }})
      </div>
      <div class="subworkflow-config__hint q-mb-xs">
        {{ t('nodes.subworkflow.config.outputMappingsHint') }}
      </div>
      <div class="subworkflow-config__info-banner q-mb-sm">
        <q-icon name="info" color="deep-purple-6" size="xs" class="q-mr-sm" />
        <span>
          Output is also available as <code>nodes.&lt;nodeId&gt;.output.*</code> in downstream expressions.
        </span>
      </div>

      <div
        v-for="(mapping, index) in outputMappings"
        :key="'out_' + index"
        class="subworkflow-config__mapping-card"
      >
        <div class="row q-col-gutter-sm">
          <!-- Child output key -->
          <div class="col">
            <q-input
              :model-value="mapping.outputKey"
              outlined
              dense
              :label="t('nodes.subworkflow.config.childOutputKey')"
              :placeholder="t('nodes.subworkflow.config.valuePlaceholder')"
              @update:model-value="(val: string | number | null) => updateOutputMapping(index, 'outputKey', String(val ?? ''))"
            >
              <template #prepend>
                <q-icon name="output" size="xs" color="deep-purple-6" />
              </template>
            </q-input>
          </div>

          <!-- Arrow indicator -->
          <div class="col-auto self-center">
            <q-icon name="arrow_forward" size="xs" color="grey-6" />
          </div>

          <!-- Parent state variable -->
          <div class="col">
            <q-select
              v-if="variableOptions.length > 0"
              :model-value="mapping.targetVariable"
              outlined
              dense
              :label="t('nodes.subworkflow.config.parentVariable')"
              :options="variableOptions"
              option-label="label"
              option-value="value"
              emit-value
              map-options
              @update:model-value="(val: string) => updateOutputMapping(index, 'targetVariable', val)"
            >
              <template #prepend>
                <q-icon name="storage" size="xs" color="purple-6" />
              </template>
            </q-select>
            <q-input
              v-else
              :model-value="mapping.targetVariable"
              outlined
              dense
              :label="t('nodes.subworkflow.config.parentVariable')"
              :placeholder="t('nodes.subworkflow.config.parentVariablePlaceholder')"
              @update:model-value="(val: string | number | null) => updateOutputMapping(index, 'targetVariable', String(val ?? ''))"
            >
              <template #prepend>
                <q-icon name="storage" size="xs" color="purple-6" />
              </template>
            </q-input>
          </div>
        </div>

        <!-- Remove mapping -->
        <div class="text-right q-mt-xs">
          <q-btn
            flat
            dense
            round
            icon="delete_outline"
            size="xs"
            color="negative"
            @click="removeOutputMapping(index)"
          >
            <AppTooltip :content="t('nodes.subworkflow.config.removeMapping')" />
          </q-btn>
        </div>
      </div>

      <!-- Add output mapping button -->
      <q-btn
        outline
        no-caps
        dense
        color="deep-purple-6"
        icon="add"
        :label="t('nodes.subworkflow.config.addOutput')"
        class="full-width"
        @click="addOutputMapping"
      />
    </div>

    <!-- EMPTY STATE -->
    <div v-if="!props.config.workflowId" class="subworkflow-config__empty">
      <q-icon name="hub" size="md" color="grey-6" />
      <div class="text-caption text-grey-6 q-mt-sm">
        {{ t('nodes.subworkflow.config.selectPrompt') }}
      </div>
    </div>

    <!-- Workflow Selector Dialog -->
    <WorkflowSelectorDialog
      v-model="workflowDialogOpen"
      :selected-workflow-id="(props.config.workflowId as string) ?? null"
      :exclude-workflow-id="(props.config._currentWorkflowId as string) ?? null"
      @select="handleWorkflowSelect"
    />
  </div>
</template>

<style lang="scss" scoped>
.subworkflow-config {
  &__section {
    margin-bottom: 16px;
  }

  &__label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    letter-spacing: 0.8px;
    margin-bottom: 8px;
  }

  &__selected {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-bg);
  }

  &__toggle {
    width: 100%;

    :deep(.q-btn) {
      font-size: 0.8rem;
    }
  }

  &__hint {
    font-size: 0.72rem;
    color: var(--mapex-text-muted);
    margin-top: 6px;
    line-height: 1.4;
  }

  &__info-banner {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    border: 1px solid var(--mapex-wf-tint-border);
    font-size: 0.72rem;
    color: var(--mapex-text-secondary);
    line-height: 1.4;

    code {
      font-family: 'Roboto Mono', monospace;
      font-size: 0.68rem;
      background: var(--mapex-wf-tint-3);
      padding: 1px 4px;
      border-radius: var(--mapex-radius-xs);
    }
  }

  &__mapping-card {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-bg);
    padding: 12px;
    margin-bottom: 10px;
  }

  &__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px 16px;
    text-align: center;
  }
}
</style>
