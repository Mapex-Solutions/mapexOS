<script setup lang="ts">
defineOptions({
  name: 'SwitchNodeConfig',
});

/** TYPE IMPORTS */
import type { FieldSourceValue, NodeConfigComponentProps, NodeConfigComponentEmits } from '@src/components/workflow/interfaces';
import type { AssetTemplateResponse } from '@mapexos/schemas';
import type {
  ConditionGroupItem,
  WorkflowConditionItem,
  GroupLogicOperator,
} from '../../conditionNode/interfaces';
import type { SwitchCase } from '../interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, onMounted } from 'vue';

/** COMPONENTS */
import { AssetTemplateSelectorDialog } from '@components/dialogs/common/assetTemplateSelectorDialog';
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';
import ConditionGroupCard from '../../conditionNode/configs/ConditionGroupCard.vue';
import ConditionItemCard from '../../conditionNode/configs/ConditionItemCard.vue';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** SERVICES */
import { apis } from '@services/mapex';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { ComparisonOperator } from '../../conditionNode/constants/conditionNode.constant';
import { GROUP_LOGIC_OPTIONS } from '../../conditionNode/constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { states } = useWorkflowContext();
const { t } = usePluginI18n('core-flow-control');

/** STATE */

/**
 * Index of the currently selected case in the q-select
 */
const activeCaseIndex = ref<number | null>(null);

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
 * Match mode options for case evaluation strategy
 */
const matchModeOptions = computed(() => [
  { label: t('nodes.switch.config.firstMatch'), value: 'first', icon: 'first_page' },
  { label: t('nodes.switch.config.allMatches'), value: 'all', icon: 'select_all' },
]);

/**
 * Workflow variables mapped to state fields format
 */
const stateFields = computed(() =>
  states.value.map(v => ({ name: v.field, type: v.type })),
);

/**
 * Current switch cases from config
 */
const cases = computed<SwitchCase[]>(
  () => (props.config.cases as SwitchCase[]) || [],
);

/**
 * Options for the case selector q-select
 */
const caseOptions = computed(() =>
  cases.value.map((c, i) => ({
    label: `Case ${i + 1}`,
    value: i,
    itemCount: c.condition.items.length,
  })),
);

/**
 * Currently active case (derived from activeCaseIndex)
 */
const activeCase = computed<SwitchCase | null>(() => {
  if (activeCaseIndex.value === null) return null;
  return cases.value[activeCaseIndex.value] ?? null;
});

/**
 * Root logic operator of the active case
 */
const activeCaseLogic = computed<GroupLogicOperator>(
  () => activeCase.value?.condition.logic || 'AND',
);

/**
 * Case evaluation mode: 'first' = stop at first match, 'all' = activate all matches
 */
const matchMode = computed<'first' | 'all'>(
  () => (props.config.matchMode as 'first' | 'all') || 'first',
);

/**
 * Current logic operator config (label, icon, color) for the active case
 */
const currentOperator = computed(() =>
  GROUP_LOGIC_OPTIONS.find(o => o.value === activeCaseLogic.value) || GROUP_LOGIC_OPTIONS[0],
);

/**
 * Items of the active case's root group
 */
const activeCaseItems = computed<ConditionGroupItem[]>(
  () => activeCase.value?.condition.items || [],
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

/** WATCHERS */

watch(cases, (newCases) => {
  if (newCases.length === 0) {
    activeCaseIndex.value = null;
  } else if (activeCaseIndex.value !== null && activeCaseIndex.value >= newCases.length) {
    activeCaseIndex.value = newCases.length - 1;
  }
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
 * Update match mode (first match vs all matches)
 *
 * @param {'first' | 'all'} mode - New match mode
 */
function updateMatchMode(mode: 'first' | 'all'): void {
  emitUpdate({ matchMode: mode });
}

/**
 * Update the active case's data and emit
 *
 * @param {Partial<SwitchCase>} partial - Partial case data to merge
 */
function updateActiveCase(partial: Partial<SwitchCase>): void {
  if (activeCaseIndex.value === null || !activeCase.value) return;
  const newCases = [...cases.value];
  newCases[activeCaseIndex.value] = { ...activeCase.value, ...partial };
  emitUpdate({ cases: newCases });
}

/**
 * Find a condition by ID across all items of the active case (root + sub-groups)
 *
 * @param {string} conditionId - Condition ID to find
 * @returns {WorkflowConditionItem | undefined} Found condition
 */
function findConditionById(conditionId: string): WorkflowConditionItem | undefined {
  for (const item of activeCaseItems.value) {
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
 * Add a new empty case
 */
function addCase(): void {
  const caseNum = cases.value.length + 1;
  const newCase: SwitchCase = {
    id: `case_${Date.now()}`,
    name: `Case ${caseNum}`,
    condition: {
      logic: 'AND',
      items: [],
    },
  };
  const newCases = [...cases.value, newCase];
  emitUpdate({ cases: newCases });
  activeCaseIndex.value = newCases.length - 1;
}

/**
 * Remove the currently active case
 */
function removeCase(): void {
  if (activeCaseIndex.value === null) return;
  const newCases = cases.value.filter((_, i) => i !== activeCaseIndex.value);
  emitUpdate({ cases: newCases });
  if (newCases.length === 0) {
    activeCaseIndex.value = null;
  } else {
    activeCaseIndex.value = Math.min(activeCaseIndex.value, newCases.length - 1);
  }
}

/**
 * Update root logic operator of the active case
 *
 * @param {GroupLogicOperator} logic - New logic operator
 */
function updateRootLogic(logic: GroupLogicOperator): void {
  if (!activeCase.value) return;
  updateActiveCase({
    condition: { ...activeCase.value.condition, logic },
  });
}

/**
 * Update an item at a given index in the active case's root items
 *
 * @param {number} index - Item index in root items
 * @param {ConditionGroupItem} updated - Updated item
 */
function updateItem(index: number, updated: ConditionGroupItem): void {
  if (!activeCase.value) return;
  const newItems = [...activeCaseItems.value];
  newItems[index] = updated;
  updateActiveCase({
    condition: { ...activeCase.value.condition, items: newItems },
  });
}

/**
 * Remove an item at a given index in the active case's root items
 *
 * @param {number} index - Item index to remove
 */
function removeItem(index: number): void {
  if (!activeCase.value) return;
  const newItems = activeCaseItems.value.filter((_, i) => i !== index);
  updateActiveCase({
    condition: { ...activeCase.value.condition, items: newItems },
  });
}

/**
 * Add a new empty condition to the active case's root items
 */
function addCondition(): void {
  if (!activeCase.value) return;
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
  updateActiveCase({
    condition: { ...activeCase.value.condition, items: [...activeCaseItems.value, newItem] },
  });
}

/**
 * Add a new sub-group to the active case's root items
 */
function addGroup(): void {
  if (!activeCase.value) return;
  const groupNum = activeCaseItems.value.filter(i => i.type === 'group').length + 1;
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
  updateActiveCase({
    condition: { ...activeCase.value.condition, items: [...activeCaseItems.value, newItem] },
  });
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
  if (!item || !activeConditionId.value || !activeCase.value) return;

  const fieldPath = item.path as string;
  const updated: FieldSourceValue = { type: 'event', value: fieldPath, mode: 'dynamic' };

  // Search root items and sub-group items for the active condition
  const newItems = activeCaseItems.value.map(rootItem => {
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

  updateActiveCase({
    condition: { ...activeCase.value.condition, items: newItems },
  });
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

/** LIFECYCLE HOOKS */

onMounted(async () => {
  if (cases.value.length > 0) {
    activeCaseIndex.value = 0;
  }

  if (selectedTemplateIds.value.length > 0) {
    await fetchTemplateFields(selectedTemplateIds.value);
  }
});
</script>

<template>
  <div class="switch-config">
    <!-- Empty state: centered + button -->
    <div v-if="cases.length === 0" class="switch-config__empty">
      <q-icon name="alt_route" size="32px" color="grey-5" class="q-mb-sm" />
      <span class="switch-config__empty-text">{{ t('nodes.switch.config.noCasesYet') }}</span>

      <q-btn
        round
        color="primary"
        icon="add"
        size="md"
        class="q-mt-md"
        @click="addCase"
      />
    </div>

    <!-- Populated state -->
    <template v-else>
      <!-- Evaluation Mode -->
      <div class="switch-config__section">
        <div class="switch-config__section-label">{{ t('nodes.switch.config.evaluationModeSection') }}</div>
        <q-btn-toggle
          :model-value="matchMode"
          :options="matchModeOptions"
          spread
          no-caps
          dense
          unelevated
          toggle-color="purple-7"
          @update:model-value="(val: string) => updateMatchMode(val as 'first' | 'all')"
        />
        <div class="text-caption q-mt-xs" style="color: var(--mapex-text-secondary);">
          {{ matchMode === 'first' ? t('nodes.switch.config.firstMatchDesc') : t('nodes.switch.config.allMatchesDesc') }}
        </div>
      </div>

      <!-- Case selector section -->
      <div class="switch-config__section">
        <div class="switch-config__section-label">{{ t('nodes.switch.config.casesSection') }}</div>

        <div class="switch-config__selector-row">
          <!-- Case selector q-select -->
          <q-select
            :model-value="activeCaseIndex"
            :options="caseOptions"
            outlined
            dense
            emit-value
            map-options
            options-dense
            option-value="value"
            option-label="label"
            class="switch-config__case-select"
            @update:model-value="(val: number) => activeCaseIndex = val"
          >
            <template #prepend>
              <q-icon name="alt_route" color="purple-6" size="xs" />
            </template>
            <template #option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section side>
                  <q-icon name="alt_route" color="purple-6" size="xs" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ scope.opt.label }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge
                    v-if="scope.opt.itemCount > 0"
                    color="purple-6"
                    text-color="white"
                    :label="scope.opt.itemCount"
                  />
                </q-item-section>
              </q-item>
            </template>
          </q-select>

          <!-- Add case button -->
          <q-btn
            round
            dense
            flat
            icon="add"
            color="primary"
            size="sm"
            @click="addCase"
          />

          <!-- Case context menu -->
          <q-btn
            flat
            dense
            round
            icon="more_vert"
            size="sm"
            color="grey-6"
            :disable="!activeCase"
          >
            <q-menu>
              <q-list dense style="min-width: 140px;">
                <q-item clickable v-close-popup @click="removeCase">
                  <q-item-section side><q-icon name="delete" size="xs" color="negative" /></q-item-section>
                  <q-item-section class="text-negative">{{ t('nodes.switch.config.deleteCase') }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </div>

      <!-- Active case conditions editor -->
      <div v-if="activeCase" class="switch-config__editor">
        <!-- Empty conditions state -->
        <div v-if="activeCaseItems.length === 0" class="switch-config__case-empty">
          <q-icon name="fact_check" size="24px" color="grey-5" class="q-mb-xs" />
          <span class="switch-config__case-empty-text">{{ t('nodes.switch.config.noConditionsInCase') }}</span>

          <q-btn
            round
            color="primary"
            icon="add"
            size="sm"
            class="q-mt-sm"
          >
            <q-menu>
              <q-list dense style="min-width: 160px;">
                <q-item clickable v-close-popup @click="addCondition">
                  <q-item-section side>
                    <q-icon name="rule" color="blue-6" size="xs" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ t('nodes.switch.config.conditionLabel') }}</q-item-label>
                    <q-item-label caption>{{ t('nodes.switch.config.conditionDescription') }}</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="addGroup">
                  <q-item-section side>
                    <q-icon name="folder" color="purple-6" size="xs" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ t('nodes.switch.config.groupLabel') }}</q-item-label>
                    <q-item-label caption>{{ t('nodes.switch.config.groupDescription') }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>

        <!-- Populated conditions -->
        <template v-else>
          <!-- Root operator chip -->
          <div class="switch-config__root-operator">
            <q-chip
              :color="currentOperator.color"
              text-color="white"
              dense
              size="sm"
              class="switch-config__operator-chip"
            >
              {{ activeCaseLogic }}

              <q-menu>
                <q-list dense style="min-width: 200px;">
                  <template v-for="(opt, index) in GROUP_LOGIC_OPTIONS" :key="opt.value">
                    <q-separator v-if="index === 2" />
                    <q-item
                      clickable
                      v-close-popup
                      :active="activeCaseLogic === opt.value"
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

            <span class="switch-config__root-description">
              {{ currentOperator.description }}
            </span>
          </div>

          <!-- Items list (conditions and sub-groups) -->
          <div class="switch-config__items">
            <template v-for="(item, index) in activeCaseItems" :key="item.data.id">
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
            :label="t('nodes.switch.config.addButton')"
            size="sm"
            class="q-mt-xs"
          >
            <q-list dense style="min-width: 160px;">
              <q-item clickable v-close-popup @click="addCondition">
                <q-item-section side>
                  <q-icon name="rule" color="blue-6" size="xs" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ t('nodes.switch.config.conditionLabel') }}</q-item-label>
                  <q-item-label caption>{{ t('nodes.switch.config.conditionDescription') }}</q-item-label>
                </q-item-section>
              </q-item>
              <q-item clickable v-close-popup @click="addGroup">
                <q-item-section side>
                  <q-icon name="folder" color="purple-6" size="xs" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ t('nodes.switch.config.groupLabel') }}</q-item-label>
                  <q-item-label caption>{{ t('nodes.switch.config.groupDescription') }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </template>
      </div>

      <!-- Default indicator -->
      <div class="switch-config__default">
        <q-icon name="last_page" size="16px" color="grey-6" />
        <span class="switch-config__default-text">{{ t('nodes.switch.config.defaultCase') }}</span>
      </div>
    </template>

    <!-- Asset Template Selector Dialog (shared across all cases) -->
    <AssetTemplateSelectorDialog
      v-model="templateDialogOpen"
      :selected-template-ids="selectedTemplateIds"
      @select="handleTemplateSelect"
    />

    <!-- Event Field Selector Dialog (shared across all cases) -->
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
          <div class="switch-config__field-dialog-banner">
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
.switch-config {
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

  &__selector-row {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  &__case-select {
    flex: 1;
    min-width: 0;
  }

  &__editor {
    padding: 12px;
    border: 2px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-elevated);
    margin-bottom: 12px;
  }

  &__case-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px 12px;
    border: 2px dashed var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
  }

  &__case-empty-text {
    font-size: 0.8rem;
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

  &__default {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    border: 1px dashed var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    margin-bottom: 8px;
  }

  &__default-text {
    font-size: 0.8rem;
    color: var(--mapex-text-secondary);
    font-style: italic;
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
