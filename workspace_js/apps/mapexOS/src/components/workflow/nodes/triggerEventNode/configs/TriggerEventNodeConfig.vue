<script setup lang="ts">
/** TYPE IMPORTS */
import type { TriggerResponse, AssetTemplateResponse } from '@mapexos/schemas';
import type { FieldSourceValue, NodeConfigComponentProps, NodeConfigComponentEmits } from '@src/components/workflow/interfaces';
/** Inline type — was in deleted rules page */
interface ParsedVariable {
  name: string;
  type: string;
  path: string;
  [key: string]: any;
}

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { AppTooltip } from '@components/tooltips';
import { TriggerSelectorDialog } from '@components/dialogs/common/triggerSelectorDialog';
import { AssetTemplateSelectorDialog } from '@components/dialogs/common/assetTemplateSelectorDialog';
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';
import { EventFieldInput } from '@components/forms/eventFieldInput';

/** COMPOSABLES */
import { usePluginI18n } from '@src/composables/workflow';

/** UTILS */
import { handleApiError } from '@utils/error';

/** Stub utils — were in deleted rules page */
function extractFieldsFromConfig(config: any): string[] { void config; return []; }
function parseAllFieldVariables(fields: string[]): Record<string, ParsedVariable[]> { void fields; return {}; }

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { t } = usePluginI18n('core-triggers');

/** STATE */

/**
 * Whether the trigger selector dialog is open
 */
const triggerDrawerOpen = ref(false);

/**
 * Whether the asset template selector dialog is open
 */
const templateDrawerOpen = ref(false);

/**
 * Whether the event field selector dialog is open
 */
const fieldSelectorOpen = ref(false);

/**
 * Variable path currently being edited (for event field selection)
 */
const activeVariablePath = ref<string | null>(null);

/**
 * Selected asset template IDs — restored from config, persisted on change
 */
const selectedTemplateIds = ref<string[]>(
  (props.config.selectedTemplateIds as string[]) ?? [],
);

/**
 * Full trigger response for the selected trigger (for variable extraction)
 */
const selectedTriggerFull = ref<TriggerResponse | null>(null);

/**
 * Template fields cache: templateId → field paths
 */
const templateFieldsCache = ref<Map<string, string[]>>(new Map());

/**
 * Template names cache: templateId → name
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
 * Parse variables from selected trigger config
 */
const parsedVariables = computed<Record<string, ParsedVariable[]>>(() => {
  if (!selectedTriggerFull.value?.config) return {};
  const fields = extractFieldsFromConfig(selectedTriggerFull.value.config);
  return parseAllFieldVariables(fields);
});

/**
 * Unique variable list across all fields (deduplicated by path)
 */
const uniqueVariables = computed(() => {
  const seen = new Set<string>();
  const result: Array<{ path: string; placeholder: string; fieldKey: string }> = [];
  for (const [fieldKey, vars] of Object.entries(parsedVariables.value)) {
    for (const v of vars) {
      if (!seen.has(v.path)) {
        seen.add(v.path);
        result.push({ path: v.path, placeholder: v.placeholder, fieldKey });
      }
    }
  }
  return result;
});

/**
 * Current trigger variables from config
 */
const triggerVariables = computed<Record<string, { path: string; placeholder: string; fieldKey: string; value: FieldSourceValue }>>(() => {
  return (props.config.variables as Record<string, { path: string; placeholder: string; fieldKey: string; value: FieldSourceValue }>) ?? {};
});

/**
 * Count of configured variables (non-empty value)
 */
const configuredCount = computed(() => {
  return Object.values(triggerVariables.value).filter(v => v.value?.value).length;
});

/**
 * Whether templates are available for field browsing
 */
const hasTemplates = computed(() => selectedTemplateIds.value.length > 0);

/**
 * All available fields from cached templates as items for GenericSelectorDialog
 */
const fieldItems = computed(() => {
  const items: Array<{ id: string; path: string; templateName: string }> = [];
  for (const [templateId, fields] of templateFieldsCache.value.entries()) {
    const name = templateNamesCache.value.get(templateId) ?? 'Unknown';
    for (const field of fields) {
      items.push({ id: `${templateId}:${field}`, path: field, templateName: name });
    }
  }
  return items;
});

/**
 * Filtered field items by search query (client-side)
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
  if (!activeVariablePath.value) return [];
  const currentValue = getVariableValue(activeVariablePath.value).value;
  if (!currentValue) return [];
  return fieldItems.value
    .filter(item => item.path === currentValue)
    .map(item => item.id);
});

/** FUNCTIONS */

/**
 * Get icon for trigger type
 *
 * @param {string | undefined} triggerType - Trigger type string
 * @returns {string} Material icon name
 */
function getTriggerTypeIcon(triggerType?: string): string {
  const icons: Record<string, string> = {
    http: 'http',
    mqtt: 'router',
    rabbitmq: 'cloud_queue',
    nats: 'cloud',
    websocket: 'cable',
    email: 'email',
    teams: 'groups',
    slack: 'chat',
  };
  return icons[triggerType ?? ''] ?? 'notifications_active';
}

/**
 * Get color for trigger category
 *
 * @param {string | undefined} category - Trigger category string
 * @returns {string} Quasar color name
 */
function getCategoryColor(category?: string): string {
  return category === 'communication' ? 'purple-6' : 'blue-6';
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
      // Cache template name from full objects if available
      const tpl = templates?.find(t => t.id === templateId);
      if (tpl?.name) {
        templateNamesCache.value.set(templateId, tpl.name);
      }

      // Fetch available fields from API
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

      // Fetch template name if not already cached
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
 * Get the FieldSourceValue for a variable by path
 *
 * @param {string} path - Variable path
 * @returns {FieldSourceValue} Current field value or default
 */
function getVariableValue(path: string): FieldSourceValue {
  return triggerVariables.value[path]?.value ?? { type: 'event', value: '' };
}

/**
 * Handle trigger selection from TriggerSelectorDrawer
 *
 * @param {TriggerResponse} trigger - Selected trigger
 * @returns {Promise<void>}
 */
async function handleTriggerSelect(trigger: TriggerResponse): Promise<void> {
  emit('update:config', {
    ...props.config,
    triggerId: trigger.id,
    triggerName: trigger.name,
    triggerType: trigger.triggerType,
    variables: {},
  });

  selectedTriggerFull.value = trigger;

  // Fetch full config for variable extraction if needed
  if (!trigger.config) {
    try {
      const full = await apis.triggers?.trigger.getById({ triggerId: trigger.id! });
      if (full) selectedTriggerFull.value = full;
    } catch (error) {
      handleApiError(error, {
        defaultMessage: 'Failed to fetch trigger details',
      });
    }
  }
}

/**
 * Clear trigger selection and variables
 *
 * @returns {void}
 */
function clearTrigger(): void {
  emit('update:config', {
    ...props.config,
    triggerId: undefined,
    triggerName: undefined,
    triggerType: undefined,
    variables: {},
    selectedTemplateIds: [],
  });
  selectedTriggerFull.value = null;
  selectedTemplateIds.value = [];
  templateFieldsCache.value.clear();
  templateNamesCache.value.clear();
}

/**
 * Handle variable value update from EventFieldInput
 *
 * @param {string} path - Variable path
 * @param {string} placeholder - Original placeholder
 * @param {string} fieldKey - Field key
 * @param {FieldSourceValue} fieldValue - Updated field value
 * @returns {void}
 */
function handleVariableUpdate(path: string, placeholder: string, fieldKey: string, fieldValue: FieldSourceValue): void {
  const currentVars = { ...triggerVariables.value };
  currentVars[path] = {
    path,
    placeholder,
    fieldKey,
    value: fieldValue,
  };
  emit('update:config', { ...props.config, variables: currentVars });
}

/**
 * Handle openEventSelector from EventFieldInput — opens field selector dialog
 *
 * @param {string} path - Variable path being edited
 * @returns {void}
 */
function handleOpenEventSelector(path: string): void {
  activeVariablePath.value = path;
  fieldSearchQuery.value = '';
  fieldSelectorOpen.value = true;
}

/**
 * Handle openTemplateSelector from EventFieldInput — opens template selector dialog
 *
 * @param {string} path - Variable path being edited
 * @returns {void}
 */
function handleOpenTemplateSelector(path: string): void {
  activeVariablePath.value = path;
  templateDrawerOpen.value = true;
}

/**
 * Handle template selection from AssetTemplateSelectorDrawer
 * Fetches available fields for each template and opens the field selector
 *
 * @param {AssetTemplateResponse[]} templates - Selected templates
 * @returns {Promise<void>}
 */
async function handleTemplateSelect(templates: AssetTemplateResponse[]): Promise<void> {
  const ids = templates.map(t => t.id!).filter(Boolean);
  selectedTemplateIds.value = ids;
  emit('update:config', { ...props.config, selectedTemplateIds: ids });
  templateDrawerOpen.value = false;

  // Fetch fields for selected templates
  await fetchTemplateFields(ids, templates);

  // Open field selector if we were in the middle of selecting a field
  if (activeVariablePath.value) {
    fieldSearchQuery.value = '';
    fieldSelectorOpen.value = true;
  }
}

/**
 * Handle field selection from GenericSelectorDialog
 *
 * @param {any[]} items - Selected items from dialog (single-select)
 * @returns {void}
 */
function handleFieldSelect(items: any[]): void {
  const item = items[0];
  if (!item || !activeVariablePath.value) return;

  const fieldPath = item.path as string;
  const varInfo = uniqueVariables.value.find(v => v.path === activeVariablePath.value);
  if (!varInfo) return;

  const currentVars = { ...triggerVariables.value };
  currentVars[activeVariablePath.value] = {
    path: varInfo.path,
    placeholder: varInfo.placeholder,
    fieldKey: varInfo.fieldKey,
    value: { type: 'event', value: fieldPath, mode: 'dynamic' },
  };
  emit('update:config', { ...props.config, variables: currentVars });
  fieldSelectorOpen.value = false;
  activeVariablePath.value = null;
}

/**
 * Handle search query from field selector dialog
 *
 * @param {string} query - Search query
 * @returns {void}
 */
function handleFieldSearch(query: string): void {
  fieldSearchQuery.value = query;
}

/** LIFECYCLE HOOKS */

/**
 * Restore selected trigger and template fields from API when editing existing node
 */
onMounted(async () => {
  const triggerId = props.config.triggerId as string | undefined;
  if (triggerId) {
    try {
      const trigger = await apis.triggers?.trigger.getById({ triggerId });
      if (trigger) {
        selectedTriggerFull.value = trigger;
      }
    } catch (error) {
      handleApiError(error, {
        defaultMessage: 'Failed to restore trigger details',
      });
    }
  }

  // Restore template fields cache if templates are already selected
  if (selectedTemplateIds.value.length > 0) {
    await fetchTemplateFields(selectedTemplateIds.value);
  }
});
</script>

<template>
  <div class="trigger-event-config">
    <!-- TRIGGER SELECTOR -->
    <div class="trigger-event-config__section">
      <div class="trigger-event-config__label">{{ t('nodes.trigger_event.config.triggerSection') }}</div>

      <!-- Selected trigger display -->
      <div
        v-if="props.config.triggerName"
        class="trigger-event-config__selected"
      >
        <q-item dense class="rounded-borders">
          <q-item-section avatar>
            <q-avatar
              :icon="getTriggerTypeIcon(props.config.triggerType as string)"
              :color="getCategoryColor(props.config.triggerType as string)"
              text-color="white"
              size="sm"
            />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-medium ellipsis">
              {{ props.config.triggerName }}
            </q-item-label>
            <q-item-label caption>
              <q-badge
                :color="getCategoryColor(props.config.triggerType as string)"
                :label="props.config.triggerType as string"
                dense
              />
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
                @click="triggerDrawerOpen = true"
              >
                <AppTooltip :content="t('nodes.trigger_event.config.changeTrigger')" />
              </q-btn>
              <q-btn
                flat
                dense
                round
                icon="close"
                size="xs"
                color="grey-7"
                @click="clearTrigger"
              >
                <AppTooltip :content="t('nodes.trigger_event.config.removeTrigger')" />
              </q-btn>
            </div>
          </q-item-section>
        </q-item>
      </div>

      <!-- Select trigger button -->
      <q-btn
        v-else
        outline
        no-caps
        dense
        color="primary"
        icon="notifications_active"
        :label="t('nodes.trigger_event.config.selectTrigger')"
        class="full-width"
        @click="triggerDrawerOpen = true"
      />
    </div>

    <!-- VARIABLES -->
    <div v-if="uniqueVariables.length > 0" class="trigger-event-config__section">
      <div class="trigger-event-config__label">
        {{ t('nodes.trigger_event.config.variablesSection') }} ({{ configuredCount }}/{{ uniqueVariables.length }})
      </div>

      <div
        v-for="v in uniqueVariables"
        :key="v.path"
        class="trigger-event-config__var-card"
      >
        <!-- Variable header -->
        <div class="trigger-event-config__var-header">
          <span class="trigger-event-config__var-placeholder">{{ v.placeholder }}</span>
          <q-badge outline color="grey-7" class="trigger-event-config__var-path">
            {{ v.path }}
          </q-badge>
        </div>

        <!-- EventFieldInput — shared component with dropdown menu -->
        <EventFieldInput
          :model-value="getVariableValue(v.path)"
          :label="v.path"
          :placeholder="'e.g. event.' + v.path.split('.').pop()"
          :has-templates="hasTemplates"
          :template-count="selectedTemplateIds.length"
          :has-state-fields="false"
          :state-fields="[]"
          :state-field-count="0"
          @update:model-value="(val: FieldSourceValue) => handleVariableUpdate(v.path, v.placeholder, v.fieldKey, val)"
          @open-event-selector="handleOpenEventSelector(v.path)"
          @open-template-selector="handleOpenTemplateSelector(v.path)"
        />
      </div>
    </div>

    <!-- EMPTY STATE -->
    <div v-if="!props.config.triggerId" class="trigger-event-config__empty">
      <q-icon name="touch_app" size="md" color="grey-6" />
      <div class="text-caption text-grey-6 q-mt-sm">
        {{ t('nodes.trigger_event.config.selectPrompt') }}
      </div>
    </div>

    <!-- Trigger Selector Dialog (centered modal — workflow uses dialog, not drawer) -->
    <TriggerSelectorDialog
      v-model="triggerDrawerOpen"
      :selected-trigger-id="(props.config.triggerId as string) ?? null"
      @select="handleTriggerSelect"
    />

    <!-- Asset Template Selector Dialog (centered modal — workflow uses dialog, not drawer) -->
    <AssetTemplateSelectorDialog
      v-model="templateDrawerOpen"
      :selected-template-ids="selectedTemplateIds"
      @select="handleTemplateSelect"
    />

    <!-- Event Field Selector Dialog (centered modal — workflow uses dialog, not drawer) -->
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
      <!-- Manage templates button in filters area -->
      <template #filters>
        <div class="col-12">
          <div class="trigger-event-config__field-dialog-banner">
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
              @click="fieldSelectorOpen = false; templateDrawerOpen = true"
            />
          </div>
        </div>
      </template>

      <!-- Item rendering -->
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
.trigger-event-config {
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

  &__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px 16px;
    text-align: center;
  }

  &__var-card {
    border: 1px solid var(--mapex-card-border);
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-surface-bg);
    padding: 12px;
    margin-bottom: 10px;

    // Stack EventFieldInput vertically in narrow panel
    :deep(.row.q-col-gutter-sm) {
      flex-direction: column;

      > .col-auto,
      > .col {
        width: 100%;
        flex: none;
        padding-left: 0;
      }

      > .col-auto {
        margin-bottom: 8px;
      }
    }
  }

  &__var-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
    gap: 8px;
  }

  &__var-placeholder {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--mapex-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }

  &__var-path {
    font-family: monospace;
    font-size: 0.65rem;
    flex-shrink: 0;
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
