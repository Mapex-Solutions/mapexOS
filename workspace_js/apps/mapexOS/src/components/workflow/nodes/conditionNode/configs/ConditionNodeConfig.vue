<script setup lang="ts">
defineOptions({
  name: 'ConditionNodeConfig',
});

/** TYPE IMPORTS */
import type { FieldSourceValue, NodeConfigComponentProps, NodeConfigComponentEmits } from '@src/components/workflow/interfaces';
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type {
  ConditionGroupItem,
  WorkflowConditionItem,
  GroupLogicOperator,
} from '../interfaces';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { AssetTemplateSelectorDialog } from '@components/dialogs/common/assetTemplateSelectorDialog';
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';
import ConditionGroupCard from './ConditionGroupCard.vue';
import ConditionItemCard from './ConditionItemCard.vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ComparisonOperator } from '../constants/conditionNode.constant';
import { GROUP_LOGIC_OPTIONS } from '../constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { states } = useWorkflowContext();
const { t } = usePluginI18n('core-logic');

/** STATE */

/**
 * Whether asset template selector dialog is open
 */
const templateDialogOpen = ref(false);

/**
 * Whether event field selector dialog is open
 */
const fieldSelectorOpen = ref(false);

/**
 * Which side of the condition the dialog was opened from
 */
const activeFieldSide = ref<'field' | 'value'>('field');

/**
 * Which condition ID triggered the event field dialog
 */
const activeConditionId = ref('');

/**
 * Selected asset template IDs — restored from config, persisted on change
 */
const selectedTemplateIds = ref<string[]>(
  (props.config.selectedTemplateIds as string[]) ?? [],
);

/**
 * Template fields cache: templateId -> field paths
 */
const templateFieldsCache = ref<Map<string, string[]>>(new Map());

/**
 * Template names cache: templateId -> name
 */
const templateNamesCache = ref<Map<string, string>>(new Map());

/**
 * Loading state for field fetching
 */
const fetchingFields = ref(false);

/**
 * Search query for field selector dialog
 */
const fieldSearchQuery = ref('');

/** COMPUTED */

/**
 * Workflow variables mapped to state fields format
 */
const stateFields = computed(() =>
  states.value.map(v => ({ name: v.field, type: v.type })),
);

/**
 * Current root group logic operator
 */
const rootLogic = computed<GroupLogicOperator>(
  () => (props.config.logic as GroupLogicOperator) || 'AND',
);

/**
 * Current root group operator config (label, icon, color)
 */
const currentOperator = computed(() =>
  GROUP_LOGIC_OPTIONS.find(o => o.value === rootLogic.value) || GROUP_LOGIC_OPTIONS[0],
);

/**
 * Current root items (conditions and/or sub-groups)
 */
const items = computed<ConditionGroupItem[]>(
  () => (props.config.items as ConditionGroupItem[]) || [],
);

/**
 * All available fields from cached templates
 */
const fieldItems = computed(() => {
  const result: Array<{ id: string; path: string; templateName: string }> = [];
  for (const [templateId, fields] of templateFieldsCache.value.entries()) {
    const name = templateNamesCache.value.get(templateId) ?? 'Unknown';
    for (const field of fields) {
      result.push({ id: `${templateId}:${field}`, path: field, templateName: name });
    }
  }
  return result;
});

/**
 * Filtered field items by search query
 */
const filteredFieldItems = computed(() => {
  if (!fieldSearchQuery.value.trim()) return fieldItems.value;
  const q = fieldSearchQuery.value.toLowerCase();
  return fieldItems.value.filter(item => item.path.toLowerCase().includes(q));
});

/**
 * Currently selected field ID for GenericSelectorDialog highlighting
 */
const selectedFieldIds = computed(() => {
  const condition = findConditionById(activeConditionId.value);
  if (!condition) return [];
  const currentVal = activeFieldSide.value === 'field'
    ? condition.field.value
    : condition.value.value;
  if (!currentVal) return [];
  return fieldItems.value
    .filter(item => item.path === currentVal)
    .map(item => item.id);
});

/** FUNCTIONS */

/**
 * Emit config update
 *
 * @param {Record<string, unknown>} partial - Partial config to merge
 */
function emitUpdate(partial: Record<string, unknown>): void {
  emit('update:config', { ...props.config, ...partial });
}

/**
 * Find a condition by ID across all items (root + sub-groups)
 *
 * @param {string} conditionId - Condition ID to find
 * @returns {WorkflowConditionItem | undefined} Found condition
 */
function findConditionById(conditionId: string): WorkflowConditionItem | undefined {
  for (const item of items.value) {
    if (item.type === 'condition' && item.data.id === conditionId) {
      return item.data;
    }
    if (item.type === 'group') {
      for (const subItem of item.data.items) {
        if (subItem.type === 'condition' && subItem.data.id === conditionId) {
          return subItem.data;
        }
      }
    }
  }
  return undefined;
}

/**
 * Update root logic operator
 *
 * @param {GroupLogicOperator} logic - New logic operator
 */
function updateRootLogic(logic: GroupLogicOperator): void {
  emitUpdate({ logic });
}

/**
 * Update an item at a given index
 *
 * @param {number} index - Item index in root items
 * @param {ConditionGroupItem} updated - Updated item
 */
function updateItem(index: number, updated: ConditionGroupItem): void {
  const newItems = [...items.value];
  newItems[index] = updated;
  emitUpdate({ items: newItems });
}

/**
 * Remove an item at a given index
 *
 * @param {number} index - Item index to remove
 */
function removeItem(index: number): void {
  const newItems = items.value.filter((_, i) => i !== index);
  emitUpdate({ items: newItems });
}

/**
 * Add a new empty condition at root level
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
  emitUpdate({ items: [...items.value, newItem] });
}

/**
 * Add a new sub-group at root level
 */
function addGroup(): void {
  const groupNum = items.value.filter(i => i.type === 'group').length + 1;
  const newItem: ConditionGroupItem = {
    type: 'group',
    data: {
      id: `g_${Date.now()}`,
      name: `Group ${groupNum}`,
      logic: 'AND',
      items: [
        {
          type: 'condition',
          data: {
            id: `c_${Date.now()}_1`,
            name: 'condition',
            field: { type: 'event', value: '' },
            operator: ComparisonOperator.Equals,
            value: { type: 'input', value: '' },
          },
        },
      ],
    },
  };
  emitUpdate({ items: [...items.value, newItem] });
}

/**
 * Handle event field selection request from a condition
 *
 * @param {{ side: 'field' | 'value'; conditionId: string }} payload - Selection context
 */
function handleEventFieldRequest(payload: { side: 'field' | 'value'; conditionId: string }): void {
  activeFieldSide.value = payload.side;
  activeConditionId.value = payload.conditionId;
  if (selectedTemplateIds.value.length === 0) {
    templateDialogOpen.value = true;
  } else {
    fieldSearchQuery.value = '';
    fieldSelectorOpen.value = true;
  }
}

/**
 * Fetch available fields for selected templates from API and cache locally
 *
 * @param {string[]} templateIds - Template IDs to fetch fields for
 * @param {AssetTemplateResponse[]} templates - Optional full template objects for name caching
 * @returns {Promise<void>}
 */
async function fetchTemplateFields(templateIds: string[], templates?: AssetTemplateResponse[]): Promise<void> {
  fetchingFields.value = true;
  try {
    for (const templateId of templateIds) {
      const tpl = templates?.find(t => t.id === templateId);
      if (tpl?.name) {
        templateNamesCache.value.set(templateId, tpl.name);
      }

      try {
        const response = await apis.assets?.assetTemplate.getAvailableFields({ assetTemplateId: templateId });
        if (response?.availableFields) {
          templateFieldsCache.value.set(templateId, response.availableFields);
        } else {
          templateFieldsCache.value.set(templateId, []);
        }
      } catch {
        templateFieldsCache.value.set(templateId, []);
      }

      if (!templateNamesCache.value.has(templateId)) {
        try {
          const template = await apis.assets?.assetTemplate.getById({ assetTemplateId: templateId });
          if (template?.name) {
            templateNamesCache.value.set(templateId, template.name);
          }
        } catch { /* ignore */ }
      }
    }
  } finally {
    fetchingFields.value = false;
  }
}

/**
 * Handle template selection from AssetTemplateSelectorDialog
 *
 * @param {AssetTemplateResponse[]} templates - Selected templates
 * @returns {Promise<void>}
 */
async function handleTemplateSelect(templates: AssetTemplateResponse[]): Promise<void> {
  const ids = templates.map(t => t.id!).filter(Boolean);
  selectedTemplateIds.value = ids;
  emitUpdate({ selectedTemplateIds: ids });
  templateDialogOpen.value = false;

  await fetchTemplateFields(ids, templates);

  fieldSearchQuery.value = '';
  fieldSelectorOpen.value = true;
}

/**
 * Handle field selection from GenericSelectorDialog — writes to the tracked condition
 *
 * @param {any[]} selectedItems - Selected items from dialog (single-select)
 */
function handleFieldSelect(selectedItems: any[]): void {
  const item = selectedItems[0];
  if (!item || !activeConditionId.value) return;

  const fieldPath = item.path as string;
  const updated: FieldSourceValue = { type: 'event', value: fieldPath, mode: 'dynamic' };

  // Search root items and sub-group items for the active condition
  const newItems = items.value.map(rootItem => {
    if (rootItem.type === 'condition' && rootItem.data.id === activeConditionId.value) {
      const cond = rootItem.data;
      const updatedCond = activeFieldSide.value === 'field'
        ? { ...cond, field: updated }
        : { ...cond, value: updated };
      return { type: 'condition' as const, data: updatedCond };
    }
    if (rootItem.type === 'group') {
      const newSubItems = rootItem.data.items.map(subItem => {
        if (subItem.type === 'condition' && subItem.data.id === activeConditionId.value) {
          const cond = subItem.data;
          const updatedCond = activeFieldSide.value === 'field'
            ? { ...cond, field: updated }
            : { ...cond, value: updated };
          return { type: 'condition' as const, data: updatedCond };
        }
        return subItem;
      });
      return { type: 'group' as const, data: { ...rootItem.data, items: newSubItems } };
    }
    return rootItem;
  });

  emitUpdate({ items: newItems });
  fieldSelectorOpen.value = false;
}

/**
 * Handle search query from field selector dialog
 *
 * @param {string} query - Search query
 */
function handleFieldSearch(query: string): void {
  fieldSearchQuery.value = query;
}

/**
 * Migrate old config formats to the new root-group format
 */
function migrateOldFormat(): void {
  // v1: flat format (field + operator + value at root)
  if (props.config.field && !props.config.groups && !props.config.items) {
    const migrated = {
      logic: 'AND',
      items: [
        {
          type: 'condition',
          data: {
            id: `c_${Date.now()}`,
            name: 'condition',
            field: props.config.field || { type: 'event', value: '' },
            operator: props.config.operator || ComparisonOperator.Equals,
            value: props.config.value || { type: 'input', value: '' },
          },
        },
      ],
      selectedTemplateIds: props.config.selectedTemplateIds || [],
    };
    emit('update:config', migrated);
    return;
  }

  // v2: grouped format (rootLogic + groups[])
  if (props.config.rootLogic && props.config.groups && !props.config.items) {
    const oldGroups = props.config.groups as Array<{
      id: string;
      name: string;
      logic: string;
      conditions: WorkflowConditionItem[];
    }>;

    let migratedItems: ConditionGroupItem[];

    if (oldGroups.length === 1) {
      // Single group: unwrap conditions as root items
      migratedItems = oldGroups[0]!.conditions.map(c => ({
        type: 'condition' as const,
        data: c,
      }));
    } else {
      // Multiple groups: each becomes a sub-group
      migratedItems = oldGroups.map(g => ({
        type: 'group' as const,
        data: {
          id: g.id,
          name: g.name,
          logic: g.logic as GroupLogicOperator,
          items: g.conditions.map(c => ({
            type: 'condition' as const,
            data: c,
          })),
        },
      }));
    }

    const migrated = {
      logic: oldGroups.length === 1
        ? oldGroups[0]!.logic
        : props.config.rootLogic,
      items: migratedItems,
      selectedTemplateIds: props.config.selectedTemplateIds || [],
    };
    emit('update:config', migrated);
  }
}

/** LIFECYCLE HOOKS */

onMounted(async () => {
  migrateOldFormat();

  if (selectedTemplateIds.value.length > 0) {
    await fetchTemplateFields(selectedTemplateIds.value);
  }
});
</script>

<template>
  <div class="condition-config">
    <!-- Empty state: centered + button -->
    <div v-if="items.length === 0" class="condition-config__empty">
      <q-icon name="fact_check" size="32px" color="grey-5" class="q-mb-sm" />
      <span class="condition-config__empty-text">{{ t('nodes.condition.config.noConditionsYet') }}</span>

      <q-btn
        round
        color="primary"
        icon="add"
        size="md"
        class="q-mt-md"
      >
        <q-menu>
          <q-list dense style="min-width: 160px;">
            <q-item clickable v-close-popup @click="addCondition">
              <q-item-section side>
                <q-icon name="rule" color="blue-6" size="xs" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ t('nodes.condition.config.conditionLabel') }}</q-item-label>
                <q-item-label caption>{{ t('nodes.condition.config.conditionDescription') }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="addGroup">
              <q-item-section side>
                <q-icon name="folder" color="purple-6" size="xs" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ t('nodes.condition.config.groupLabel') }}</q-item-label>
                <q-item-label caption>{{ t('nodes.condition.config.groupDescription') }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>

    <!-- Populated state -->
    <template v-else>
      <!-- Root operator chip -->
      <div class="condition-config__root-operator">
        <q-chip
          :color="currentOperator.color"
          text-color="white"
          dense
          size="sm"
          class="condition-config__operator-chip"
        >
          {{ rootLogic }}

          <q-menu>
            <q-list dense style="min-width: 200px;">
              <template v-for="(opt, index) in GROUP_LOGIC_OPTIONS" :key="opt.value">
                <q-separator v-if="index === 2" />
                <q-item
                  clickable
                  v-close-popup
                  :active="rootLogic === opt.value"
                  @click="updateRootLogic(opt.value as GroupLogicOperator)"
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

        <span class="condition-config__root-description">
          {{ currentOperator.description }}
        </span>
      </div>

      <!-- Items list (conditions and sub-groups) -->
      <div class="condition-config__items">
        <template v-for="(item, index) in items" :key="item.data.id">
          <ConditionItemCard
            v-if="item.type === 'condition'"
            :condition="item.data"
            :state-fields="stateFields"
            :can-remove="true"
            @update:condition="(updated) => updateItem(index, { type: 'condition', data: updated })"
            @remove="removeItem(index)"
            @select-event-field="(payload) => handleEventFieldRequest({ ...payload, conditionId: item.data.id })"
          />

          <ConditionGroupCard
            v-else
            :group="item.data"
            :can-remove="true"
            :state-fields="stateFields"
            :is-sub-group="true"
            @update:group="(updated) => updateItem(index, { type: 'group', data: updated })"
            @remove="removeItem(index)"
            @select-event-field="handleEventFieldRequest"
          />
        </template>
      </div>

      <!-- Add button with Condition / Group choice -->
      <q-btn-dropdown
        flat
        dense
        no-caps
        color="primary"
        icon="add"
        :label="t('nodes.condition.config.addButton')"
        size="sm"
        class="q-mt-xs"
      >
        <q-list dense style="min-width: 160px;">
          <q-item clickable v-close-popup @click="addCondition">
            <q-item-section side>
              <q-icon name="rule" color="blue-6" size="xs" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ t('nodes.condition.config.conditionLabel') }}</q-item-label>
              <q-item-label caption>{{ t('nodes.condition.config.conditionDescription') }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-item clickable v-close-popup @click="addGroup">
            <q-item-section side>
              <q-icon name="folder" color="purple-6" size="xs" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ t('nodes.condition.config.groupLabel') }}</q-item-label>
              <q-item-label caption>{{ t('nodes.condition.config.groupDescription') }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-btn-dropdown>
    </template>

    <!-- Asset Template Selector Dialog (shared across all conditions) -->
    <AssetTemplateSelectorDialog
      v-model="templateDialogOpen"
      :selected-template-ids="selectedTemplateIds"
      @select="handleTemplateSelect"
    />

    <!-- Event Field Selector Dialog (shared across all conditions) -->
    <GenericSelectorDialog
      :model-value="fieldSelectorOpen"
      title="Select Event Field"
      icon="list_alt"
      icon-color="blue-6"
      :items="filteredFieldItems"
      :multi-select="false"
      :selected-ids="selectedFieldIds"
      :loading="fetchingFields"
      search-placeholder="Search fields..."
      empty-text="No fields available from selected templates."
      empty-icon="inbox"
      results-icon="list_alt"
      footer-icon="list_alt"
      item-noun-singular="field"
      item-noun-plural="fields"
      @update:model-value="fieldSelectorOpen = $event"
      @select="handleFieldSelect"
      @search="handleFieldSearch"
    >
      <template #filters>
        <div class="col-12">
          <div class="condition-config__field-dialog-banner">
            <q-icon name="memory" color="blue-6" size="xs" class="q-mr-sm" />
            <span class="text-caption" style="color: var(--mapex-text-secondary);">
              {{ selectedTemplateIds.length }} template{{ selectedTemplateIds.length !== 1 ? 's' : '' }} selected
            </span>
            <q-space />
            <q-btn
              flat
              dense
              no-caps
              size="sm"
              color="primary"
              label="Change"
              icon="swap_horiz"
              @click="fieldSelectorOpen = false; templateDialogOpen = true"
            />
          </div>
        </div>
      </template>

      <template #item="{ item }">
        <q-item-section avatar>
          <q-icon name="data_object" color="blue-6" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium" style="font-family: 'Roboto Mono', monospace; font-size: 0.85rem;">
            {{ item.path }}
          </q-item-label>
          <q-item-label caption>
            From: {{ item.templateName }}
          </q-item-label>
        </q-item-section>
      </template>
    </GenericSelectorDialog>
  </div>
</template>

<style lang="scss" scoped>
.condition-config {
  &__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 32px 16px;
    border: 2px dashed var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-elevated);
  }

  &__empty-text {
    font-size: 0.85rem;
    color: var(--mapex-text-secondary);
  }

  &__root-operator {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }

  &__operator-chip {
    flex-shrink: 0;
    cursor: pointer;
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.5px;
  }

  &__root-description {
    font-size: 0.75rem;
    color: var(--mapex-text-secondary);
  }

  &__items {
    margin-bottom: 4px;
  }

  &__field-dialog-banner {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    border: 1px solid var(--mapex-wf-tint-border);
  }
}
</style>
