<script setup lang="ts">
defineOptions({
  name: 'ConditionGroupCard',
});

/** TYPE IMPORTS */
import type {
  WorkflowConditionGroup,
  ConditionGroupItem,
  GroupLogicOperator,
} from '../interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, nextTick } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips/appTooltip';
import ConditionItemCard from './ConditionItemCard.vue';

/** COMPOSABLES */
import { usePluginI18n } from '@src/composables/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ComparisonOperator } from '../constants/conditionNode.constant';
import { GROUP_LOGIC_OPTIONS } from '../constants';

/** PROPS & EMITS */
const props = defineProps<{
  /** Condition group data */
  group: WorkflowConditionGroup;
  /** Whether this group can be removed */
  canRemove: boolean;
  /** Workflow state fields for condition dropdowns */
  stateFields: Array<{ name: string; type: string }>;
  /** Whether this is a sub-group (prevents nesting groups inside) */
  isSubGroup?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:group', group: WorkflowConditionGroup): void;
  (e: 'remove'): void;
  (e: 'select-event-field', payload: { side: 'field' | 'value'; conditionId: string }): void;
}>();

/** COMPOSABLES & STORES */
const { t } = usePluginI18n('core-logic');

/** STATE */

/**
 * Whether the group body is expanded (groups start expanded)
 */
const isExpanded = ref(true);

/**
 * Whether inline name editing is active
 */
const isEditingName = ref(false);

/**
 * Editable name buffer
 */
const editableName = ref(props.group.name ?? '');

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
 * Current operator config (label, icon, color)
 */
const currentOperator = computed(() =>
  GROUP_LOGIC_OPTIONS.find(o => o.value === props.group.logic) || GROUP_LOGIC_OPTIONS[0],
);

/**
 * Items count text for collapsed state
 */
const itemsCountText = computed(() => {
  const count = props.group.items.length;
  const noun = count === 1 ? t('nodes.condition.config.itemSingular') : t('nodes.condition.config.itemPlural');
  return `${count} ${noun}`;
});

/** WATCHERS */

watch(() => props.group.name, (newName) => {
  if (!isEditingName.value) {
    editableName.value = newName ?? '';
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
  originalName.value = editableName.value ?? '';
  isEditingName.value = true;
  void nextTick(() => nameInputRef.value?.focus());
}

/**
 * Save edited name
 */
function saveName(): void {
  const trimmed = (editableName.value ?? '').trim();
  editableName.value = trimmed || originalName.value;
  isEditingName.value = false;
  if (trimmed && trimmed !== props.group.name) {
    emitGroupUpdate({ name: trimmed });
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
 * Emit group update with merged values
 *
 * @param {Partial<WorkflowConditionGroup>} partial - Partial group to merge
 */
function emitGroupUpdate(partial: Partial<WorkflowConditionGroup>): void {
  emit('update:group', { ...props.group, ...partial });
}

/**
 * Update the group logic operator
 *
 * @param {GroupLogicOperator} logic - New logic operator
 */
function updateLogic(logic: GroupLogicOperator): void {
  emitGroupUpdate({ logic });
}

/**
 * Update an item at a given index
 *
 * @param {number} index - Item index
 * @param {ConditionGroupItem} updated - Updated item
 */
function updateItem(index: number, updated: ConditionGroupItem): void {
  const newItems = [...props.group.items];
  newItems[index] = updated;
  emitGroupUpdate({ items: newItems });
}

/**
 * Remove an item at a given index
 *
 * @param {number} index - Item index to remove
 */
function removeItem(index: number): void {
  const newItems = props.group.items.filter((_, i) => i !== index);
  emitGroupUpdate({ items: newItems });
}

/**
 * Add a new empty condition to this group
 */
function addCondition(): void {
  const newItem: ConditionGroupItem = {
    type: 'condition',
    data: {
      id: `c_${Date.now()}`,
      name: 'condition',
      field: { type: 'event', value: '' },
      operator: ComparisonOperator.Equals,
      value: { type: 'input', value: '' },
    },
  };
  emitGroupUpdate({ items: [...props.group.items, newItem] });
}

/**
 * Handle event field selection from a condition item
 *
 * @param {string} conditionId - Condition that triggered the event
 * @param {{ side: 'field' | 'value' }} payload - Selection payload
 */
function handleConditionEventField(conditionId: string, payload: { side: 'field' | 'value' }): void {
  emit('select-event-field', { ...payload, conditionId });
}
</script>

<template>
  <div class="condition-group-card" :class="{ 'condition-group-card--collapsed': !isExpanded }">
    <!-- Header -->
    <div class="condition-group-card__header" @click="toggleExpanded">
      <!-- Expand/collapse chevron -->
      <q-icon
        :name="isExpanded ? 'expand_more' : 'chevron_right'"
        size="20px"
        color="grey-6"
        class="condition-group-card__chevron"
      />

      <!-- Operator badge with dropdown -->
      <q-chip
        :color="currentOperator.color"
        text-color="white"
        dense
        size="sm"
        class="condition-group-card__operator-chip"
        @click.stop
      >
        {{ group.logic }}

        <q-menu>
          <q-list dense style="min-width: 200px;">
            <template v-for="(opt, index) in GROUP_LOGIC_OPTIONS" :key="opt.value">
              <q-separator v-if="index === 2" />
              <q-item
                clickable
                v-close-popup
                :active="group.logic === opt.value"
                @click="updateLogic(opt.value as GroupLogicOperator)"
              >
                <q-item-section side>
                  <q-icon :name="opt.icon" :color="opt.color" size="xs" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ opt.label }}</q-item-label>
                  <q-item-label caption>{{ opt.description }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-list>
        </q-menu>
      </q-chip>

      <!-- Name display / edit -->
      <div v-if="!isEditingName" class="condition-group-card__name" @click.stop="startEditingName">
        {{ editableName }}
        <q-icon name="edit" size="12px" class="condition-group-card__edit-icon" />
        <AppTooltip :content="editableName ?? ''" />
      </div>
      <div v-else class="condition-group-card__name-editor" @click.stop>
        <q-input
          ref="nameInputRef"
          v-model="editableName"
          dense
          borderless
          input-class="condition-group-card__name-input"
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

      <!-- Collapsed count -->
      <span v-if="!isExpanded" class="condition-group-card__count">
        {{ itemsCountText }}
      </span>

      <!-- Context menu -->
      <q-btn
        flat
        dense
        round
        icon="more_vert"
        size="xs"
        color="grey-6"
        class="condition-group-card__menu-btn"
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
              <q-item-section class="text-negative">{{ t('nodes.condition.config.deleteGroup') }}</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>

    <!-- Body (expanded) -->
    <div v-if="isExpanded" class="condition-group-card__body">
      <!-- Empty state -->
      <div v-if="group.items.length === 0" class="condition-group-card__empty">
        <q-btn
          round
          outline
          color="primary"
          icon="add"
          size="sm"
          @click="addCondition"
        />
        <span class="condition-group-card__empty-text">{{ t('nodes.condition.config.addConditionHint') }}</span>
      </div>

      <!-- Conditions list -->
      <template v-else>
        <template v-for="(item, index) in group.items" :key="item.data.id">
          <ConditionItemCard
            v-if="item.type === 'condition'"
            :condition="item.data"
            :state-fields="stateFields"
            :can-remove="true"
            @update:condition="(updated) => updateItem(index, { type: 'condition', data: updated })"
            @remove="removeItem(index)"
            @select-event-field="(payload) => handleConditionEventField(item.data.id, payload)"
          />
        </template>

        <!-- Add condition button (sub-groups only allow conditions) -->
        <q-btn
          flat
          dense
          no-caps
          color="primary"
          icon="add"
          :label="t('nodes.condition.config.addConditionButton')"
          size="sm"
          class="q-mt-xs"
          @click="addCondition"
        />
      </template>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.condition-group-card {
  border: 2px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);
  background: var(--mapex-surface-elevated);
  margin-bottom: 12px;
  transition: border-color var(--mapex-transition-base);

  &--collapsed {
    border-bottom-width: 2px;
  }

  &__header {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    cursor: pointer;
    gap: 6px;
    min-height: 42px;

    &:hover {
      background: var(--mapex-wf-tint-1);
      border-radius: var(--mapex-radius-md);
    }
  }

  &__chevron {
    flex-shrink: 0;
  }

  &__operator-chip {
    flex-shrink: 0;
    cursor: pointer;
    font-size: 0.7rem;
    font-weight: 600;
  }

  &__name {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
    flex: 1;
    cursor: pointer;

    &:hover .condition-group-card__edit-icon {
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
    font-size: 0.85rem !important;
    padding: 2px 4px !important;
  }

  &__count {
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
    white-space: nowrap;
    flex-shrink: 0;
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

  &__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px 8px;
    gap: 8px;
  }

  &__empty-text {
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
  }
}
</style>
