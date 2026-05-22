<script setup lang="ts">
defineOptions({
  name: 'ConditionItemCard',
});

/** TYPE IMPORTS */
import type { FieldSourceValue, NodeOutputOption } from '@src/components/workflow/interfaces';
import type { WorkflowConditionItem } from '../interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, nextTick } from 'vue';

/** COMPONENTS */
import { FieldSourceSelector } from '@components/forms/fieldSourceSelector';
import { AppTooltip } from '@components/tooltips/appTooltip';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ComparisonOperator } from '../constants/conditionNode.constant';
import { CONDITION_OPERATOR_OPTIONS } from '../constants';
import { SOURCE_TYPE_OPTIONS } from '@components/forms/fieldSourceSelector';

/** PROPS & EMITS */
const props = defineProps<{
  /** Condition item data */
  condition: WorkflowConditionItem;
  /** Workflow state fields for dropdowns */
  stateFields: Array<{ name: string; type: string }>;
  /** Whether this condition can be removed (false when only 1 in group) */
  canRemove: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:condition', condition: WorkflowConditionItem): void;
  (e: 'remove'): void;
  (e: 'select-event-field', payload: { side: 'field' | 'value' }): void;
}>();

/** COMPOSABLES & STORES */
const { nodes, getNodeType } = useWorkflowContext();
const { t } = usePluginI18n('core-logic');

/** STATE */

/**
 * Whether the condition body is expanded
 */
const isExpanded = ref(false);

/**
 * Whether inline name editing is active
 */
const isEditingName = ref(false);

/**
 * Editable name buffer
 */
const editableName = ref(props.condition.name);

/**
 * Original name before editing (for cancel)
 */
const originalName = ref('');

/**
 * Reference to name input for autofocus
 */
const nameInputRef = ref<{ focus: () => void } | null>(null);

/** COMPUTED */

/**
 * Whether the current operator is unary (no value side needed)
 */
const isUnaryOperator = computed(() => {
  const op = props.condition.operator;
  return op === ComparisonOperator.IsNull || op === ComparisonOperator.IsNotNull;
});

/**
 * Collapsed summary text showing field symbol value
 *
 * @returns {string} Summary like "status = COMPLETE"
 */
const collapsedSummary = computed((): string => {
  const fieldText = props.condition.field.value || 'field';
  const opEntry = CONDITION_OPERATOR_OPTIONS.find(o => o.value === props.condition.operator);
  const opSymbol = opEntry?.symbol || '=';

  if (isUnaryOperator.value) return `${fieldText} ${opSymbol}`;

  const valueText = props.condition.value.value || 'value';
  return `${fieldText} ${opSymbol} ${valueText}`;
});

/**
 * Source type config for field side icon/color
 */
const fieldSourceConfig = computed(() =>
  SOURCE_TYPE_OPTIONS.find(o => o.value === props.condition.field.type) || SOURCE_TYPE_OPTIONS[0],
);

/**
 * Field value as FieldSourceValue.
 * Pass-through preserves all fields (type, value, mode, nodeId) so the
 * FieldSourceSelector can detect modes like `assetStatus`.
 */
const fieldSource = computed<FieldSourceValue>(() => ({ ...props.condition.field }));

/**
 * Value as FieldSourceValue.
 * Pass-through preserves all fields (see fieldSource above).
 */
const valueSource = computed<FieldSourceValue>(() => ({ ...props.condition.value }));

/**
 * Node output options — only nodes with outputHints (produce data)
 */
const nodeOutputOptions = computed<NodeOutputOption[]>(() =>
  nodes.value
    .filter(n => {
      const def = getNodeType(n.type);
      return def?.availableOutputs && def.availableOutputs.length > 0;
    })
    .map(n => ({
      id: n.id,
      label: `${n.label || n.id} (${n.type.split('/').pop() || 'node'})`,
      type: n.type,
    })),
);

/** WATCHERS */

watch(() => props.condition.name, (newName) => {
  if (!isEditingName.value) {
    editableName.value = newName;
  }
});

/** FUNCTIONS */

/**
 * Toggle expand/collapse state
 */
function toggleExpanded(): void {
  isExpanded.value = !isExpanded.value;
}

/**
 * Start inline name editing
 */
function startEditingName(): void {
  originalName.value = editableName.value;
  isEditingName.value = true;
  void nextTick(() => nameInputRef.value?.focus());
}

/**
 * Save edited name
 */
function saveName(): void {
  const trimmed = editableName.value.trim();
  editableName.value = trimmed || originalName.value;
  isEditingName.value = false;
  if (trimmed && trimmed !== props.condition.name) {
    emitUpdate({ name: trimmed });
  }
}

/**
 * Cancel name editing
 */
function cancelEditName(): void {
  editableName.value = originalName.value;
  isEditingName.value = false;
}

/**
 * Emit condition update with merged values
 *
 * @param {Partial<WorkflowConditionItem>} partial - Partial condition to merge
 */
function emitUpdate(partial: Partial<WorkflowConditionItem>): void {
  emit('update:condition', { ...props.condition, ...partial });
}

/**
 * Handle field source update from FieldSourceSelector.
 * Pass the entire FieldSourceValue (preserves `mode`, `nodeId`, etc.) so
 * modes like `assetStatus` and `manual`/`dynamic` survive the round-trip.
 *
 * @param {FieldSourceValue} value - Updated field source
 */
function handleFieldUpdate(value: FieldSourceValue): void {
  emitUpdate({ field: { ...value } });
}

/**
 * Handle value source update from FieldSourceSelector.
 * Same pass-through semantics as handleFieldUpdate.
 *
 * @param {FieldSourceValue} value - Updated value source
 */
function handleValueUpdate(value: FieldSourceValue): void {
  emitUpdate({ value: { ...value } });
}

/**
 * Update the comparison operator
 *
 * @param {string} operator - New operator value
 */
function updateOperator(operator: string): void {
  emitUpdate({ operator: operator as ComparisonOperator });
}

/**
 * Handle event field selector request — bubble up to parent (SwitchNodeConfig)
 *
 * @param {'field' | 'value'} side - Which side triggered the request
 */
function handleEventFieldRequest(side: 'field' | 'value'): void {
  emit('select-event-field', { side });
}
</script>

<template>
  <div class="condition-item-card" :class="{ 'condition-item-card--expanded': isExpanded }">
    <!-- Header -->
    <div class="condition-item-card__header" @click="toggleExpanded">
      <div class="condition-item-card__header-row">
        <!-- Expand/collapse chevron -->
        <q-icon
          :name="isExpanded ? 'expand_more' : 'chevron_right'"
          size="18px"
          color="grey-6"
          class="condition-item-card__chevron"
        />

        <!-- Condition type icon -->
        <q-icon
          :name="fieldSourceConfig?.icon || 'event'"
          :color="fieldSourceConfig?.color || 'blue-6'"
          size="16px"
        />

        <!-- Name display / edit -->
        <div v-if="!isEditingName" class="condition-item-card__name" @click.stop="startEditingName">
          {{ editableName }}
          <q-icon name="edit" size="12px" class="condition-item-card__edit-icon" />
          <AppTooltip :content="editableName" />
        </div>
        <div v-else class="condition-item-card__name-editor" @click.stop>
          <q-input
            ref="nameInputRef"
            v-model="editableName"
            dense
            borderless
            input-class="condition-item-card__name-input"
            @blur="saveName"
            @keyup.enter="saveName"
            @keyup.esc="cancelEditName"
          >
            <template #append>
              <q-icon name="check" size="14px" color="positive" class="cursor-pointer" @click="saveName" />
              <q-icon name="close" size="14px" color="grey-6" class="cursor-pointer q-ml-xs" @click="cancelEditName" />
            </template>
          </q-input>
        </div>

        <!-- Context menu -->
        <q-btn
          flat
          dense
          round
          icon="more_vert"
          size="xs"
          color="grey-6"
          class="condition-item-card__menu-btn"
          @click.stop
        >
          <q-menu>
            <q-list dense style="min-width: 140px;">
              <q-item clickable v-close-popup @click="startEditingName">
                <q-item-section side><q-icon name="edit" size="xs" /></q-item-section>
                <q-item-section>{{ t('nodes.condition.config.rename') }}</q-item-section>
              </q-item>
              <q-separator />
              <q-item
                clickable
                v-close-popup
                :disable="!canRemove"
                @click="emit('remove')"
              >
                <q-item-section side><q-icon name="delete" size="xs" color="negative" /></q-item-section>
                <q-item-section class="text-negative">{{ t('nodes.condition.config.deleteItem') }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>

      <!-- Collapsed summary (second line) -->
      <div v-if="!isExpanded" class="condition-item-card__summary-row">
        <span class="condition-item-card__summary">
          {{ collapsedSummary }}
        </span>
      </div>
    </div>

    <!-- Expanded body -->
    <div v-if="isExpanded" class="condition-item-card__body">
      <!-- FIELD section -->
      <div class="condition-item-card__section">
        <div class="condition-item-card__section-label">{{ t('nodes.condition.config.ifLabel') }}</div>

        <FieldSourceSelector
          :model-value="fieldSource"
          :allowed-types="['event', 'assetStatus', 'state', 'input', 'literal', 'nodeOutput']"
          :state-fields="stateFields"
          :node-output-options="nodeOutputOptions"
          @update:model-value="handleFieldUpdate"
          @open-event-selector="handleEventFieldRequest('field')"
          @open-template-selector="handleEventFieldRequest('field')"
        />
      </div>

      <!-- OPERATOR section -->
      <div class="condition-item-card__section">
        <div class="condition-item-card__section-label">{{ t('nodes.condition.config.operatorLabel') }}</div>
        <q-select
          :model-value="condition.operator"
          :options="[...CONDITION_OPERATOR_OPTIONS]"
          outlined
          dense
          emit-value
          map-options
          options-dense
          option-value="value"
          option-label="label"
          @update:model-value="updateOperator"
        >
          <template #prepend>
            <q-icon name="compare_arrows" color="amber-8" size="xs" />
          </template>
          <template #option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section side class="condition-item-card__operator-symbol">
                {{ scope.opt.symbol }}
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ scope.opt.label }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>
      </div>

      <!-- VALUE section (hidden for unary operators) -->
      <div v-if="!isUnaryOperator" class="condition-item-card__section">
        <div class="condition-item-card__section-label">{{ t('nodes.condition.config.compareToLabel') }}</div>

        <FieldSourceSelector
          :model-value="valueSource"
          :allowed-types="['event', 'assetStatus', 'state', 'input', 'literal', 'nodeOutput']"
          :state-fields="stateFields"
          :node-output-options="nodeOutputOptions"
          @update:model-value="handleValueUpdate"
          @open-event-selector="handleEventFieldRequest('value')"
          @open-template-selector="handleEventFieldRequest('value')"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.condition-item-card {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  background: var(--mapex-surface-bg);
  margin-bottom: 8px;
  transition: border-color var(--mapex-transition-base);

  &--expanded {
    border-color: var(--mapex-text-muted);
  }

  &__header {
    display: flex;
    flex-direction: column;
    padding: 8px 10px;
    cursor: pointer;
    min-height: 36px;

    &:hover {
      background: var(--mapex-wf-tint-1);
    }
  }

  &__header-row {
    display: flex;
    align-items: center;
    gap: 6px;
    width: 100%;
  }

  &__chevron {
    flex-shrink: 0;
  }

  &__name {
    font-size: 0.8rem;
    font-weight: 500;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
    flex: 1;
    cursor: pointer;

    &:hover .condition-item-card__edit-icon {
      opacity: 1;
    }
  }

  &__edit-icon {
    opacity: 0;
    color: var(--mapex-text-secondary);
    transition: opacity var(--mapex-transition-base);
    margin-left: 2px;
  }

  &__name-editor {
    flex: 1;
    min-width: 0;
  }

  &__name-input {
    font-size: 0.8rem !important;
    padding: 2px 4px !important;
  }

  &__summary-row {
    padding: 2px 0 0 28px;
  }

  &__summary {
    font-family: 'Roboto Mono', monospace;
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
  }

  &__menu-btn {
    flex-shrink: 0;
    opacity: 0.5;
    transition: opacity var(--mapex-transition-base);

    &:hover {
      opacity: 1;
    }
  }

  &__body {
    padding: 12px;
    border-top: 1px solid var(--mapex-card-border);
  }

  &__section {
    margin-bottom: 14px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  &__section-label {
    display: flex;
    align-items: center;
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--mapex-text-secondary);
    letter-spacing: 0.5px;
    margin-bottom: 6px;
  }

  &__operator-symbol {
    font-size: 14px;
    font-weight: 700;
    min-width: 24px;
    text-align: center;
    color: var(--mapex-text-secondary);
  }
}
</style>
