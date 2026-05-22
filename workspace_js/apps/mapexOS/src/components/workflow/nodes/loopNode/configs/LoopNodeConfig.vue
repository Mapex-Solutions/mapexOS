<script setup lang="ts">
defineOptions({
  name: 'LoopNodeConfig',
});

/** TYPE IMPORTS */
import type {
  NodeConfigComponentProps,
  NodeConfigComponentEmits,
  FieldSourceValue,
  NodeOutputOption,
} from '@src/components/workflow/interfaces';
import type { AssetTemplateResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed, onMounted } from 'vue';

/** COMPONENTS */
import { FieldSourceSelector } from '@components/forms/fieldSourceSelector';
import { AssetTemplateSelectorDialog } from '@components/dialogs/common/assetTemplateSelectorDialog';
import { GenericSelectorDialog } from '@components/dialogs/common/genericSelectorDialog';

/** COMPOSABLES */
import { useWorkflowContext, usePluginI18n } from '@src/composables/workflow';

/** SERVICES */
import { apis } from '@services/mapex';

/** PROPS & EMITS */
const props = defineProps<NodeConfigComponentProps>();
const emit = defineEmits<NodeConfigComponentEmits>();

/** COMPOSABLES & STORES */
const { nodes, getNodeType } = useWorkflowContext();
const { t } = usePluginI18n('core-flow-control');

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
 * Selected asset template IDs — restored from config
 */
const selectedTemplateIds = ref<string[]>(
  (props.config.selectedTemplateIds as string[]) ?? [],
);

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
 * Source value as FieldSourceValue (preserves nodeId/mode for nodeOutput/event)
 */
const source = computed<FieldSourceValue>(() => {
  const raw = props.config.source as FieldSourceValue | undefined;
  return {
    type: raw?.type ?? 'state',
    value: raw?.value ?? '',
    ...(raw?.nodeId != null && { nodeId: raw.nodeId }),
    ...(raw?.mode != null && { mode: raw.mode }),
  };
});

/**
 * Whether templates have been selected (for event mode)
 */
const hasTemplates = computed(() => selectedTemplateIds.value.length > 0);

/**
 * Node output options — only nodes with outputHints (produce data), excluding self
 */
const nodeOutputOptions = computed<NodeOutputOption[]>(() =>
  nodes.value
    .filter(n => {
      if (n.id === props.config._nodeId) return false;
      const def = getNodeType(n.type);
      return def?.availableOutputs && def.availableOutputs.length > 0;
    })
    .map(n => ({
      id: n.id,
      label: `${n.label || n.id} (${n.type.split('/').pop() || 'node'})`,
      type: n.type,
    })),
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
  if (!source.value.value || source.value.type !== 'event') return [];
  return fieldItems.value
    .filter(item => item.path === source.value.value)
    .map(item => item.id);
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
 * Handle source value update from FieldSourceSelector
 *
 * @param {FieldSourceValue} value - Updated source value
 */
function handleSourceUpdate(value: FieldSourceValue): void {
  emitUpdate({ source: value });
}

/**
 * Handle open event selector — open field selector if templates exist
 */
function handleOpenEventSelector(): void {
  fieldSearchQuery.value = '';
  fieldSelectorOpen.value = true;
}

/**
 * Handle open template selector
 */
function handleOpenTemplateSelector(): void {
  templateDialogOpen.value = true;
}

/**
 * Fetch available fields for selected templates from API and cache locally
 *
 * @param {string[]} templateIds - Template IDs to fetch fields for
 * @param {AssetTemplateResponse[]} templates - Optional full template objects for name caching
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
 * Handle field selection from GenericSelectorDialog
 *
 * @param {any[]} selectedItems - Selected items from dialog
 */
function handleFieldSelect(selectedItems: any[]): void {
  const item = selectedItems[0];
  if (!item) return;

  emitUpdate({ source: { type: 'event', value: item.path as string } });
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
  if (selectedTemplateIds.value.length > 0) {
    await fetchTemplateFields(selectedTemplateIds.value);
  }
});
</script>

<template>
  <div class="loop-config">
    <!-- Source Section -->
    <div class="loop-config__section">
      <div class="loop-config__section-label">{{ t('nodes.loop.config.sourceSection') }}</div>

      <FieldSourceSelector
        :model-value="source"
        :allowed-types="['event', 'assetStatus', 'state', 'input', 'literal', 'nodeOutput']"
        :has-templates="hasTemplates"
        :template-count="selectedTemplateIds.length"
        :node-output-options="nodeOutputOptions"
        @update:model-value="handleSourceUpdate"
        @open-event-selector="handleOpenEventSelector"
        @open-template-selector="handleOpenTemplateSelector"
      />
    </div>

    <!-- Usage hint -->
    <q-banner dense rounded class="q-mt-sm" style="background: var(--mapex-wf-tint-2); border: 1px solid var(--mapex-wf-tint-border);">
      <template #avatar>
        <q-icon name="lightbulb" color="amber-8" size="xs" />
      </template>
      <span class="text-caption" style="color: var(--mapex-text-secondary);">
        In <strong>Body</strong> nodes, select <strong>Node Output</strong> &rarr; this Loop:
        <br />
        <code>item</code> &mdash; current element (string, number, or object)
        <br />
        <code>item.field</code> &mdash; property of an object element
        <br />
        <code>index</code> &mdash; 0-based position
      </span>
    </q-banner>

    <!-- Asset Template Selector Dialog -->
    <AssetTemplateSelectorDialog
      v-model="templateDialogOpen"
      :selected-template-ids="selectedTemplateIds"
      @select="handleTemplateSelect"
    />

    <!-- Event Field Selector Dialog -->
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
          <div class="loop-config__field-dialog-banner">
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
.loop-config {
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

  &__banner {
    display: flex;
    align-items: flex-start;
    padding: 8px 12px;
    border-radius: var(--mapex-radius-md);
    background: var(--mapex-wf-tint-2);
    border: 1px solid var(--mapex-wf-tint-border);

    code {
      font-size: 0.7rem;
      padding: 1px 4px;
      border-radius: var(--mapex-radius-xs);
      background: var(--mapex-surface-elevated);
    }
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
