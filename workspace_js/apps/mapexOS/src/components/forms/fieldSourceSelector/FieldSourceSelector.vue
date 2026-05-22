<script setup lang="ts">
defineOptions({
  name: 'FieldSourceSelector',
});

/**
 * FieldSourceSelector - Reusable source type selector component
 *
 * Supports event/state/variable/literal/nodeOutput source types.
 * Follows the EventFieldInput emit-based pattern: emits events for the parent
 */

/** TYPE IMPORTS */
import type {
  FieldSourceSelectorProps,
  FieldSourceSelectorEmits,
  FieldSourceValue,
  SourceType,
} from './interfaces/fieldSourceSelector.interface';
import type { AssetStatusFieldOption } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed, ref, watch } from 'vue';

/** COMPOSABLES */
import { useWorkflowEditorState } from '@src/pages/automations/workflows/createEditWorkflowPage/composables';
import { useTS } from '@utils/translation';

/** SERVICES */
import { apis } from '@services/mapex';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { SOURCE_TYPE_OPTIONS, DEFAULT_FIELD_SOURCE_VALUE } from './constants';
import { ASSET_STATUS_FIELD_OPTIONS } from '@src/components/workflow/constants';

/** PROPS & EMITS */
const props = withDefaults(defineProps<FieldSourceSelectorProps>(), {
  label: '',
  placeholder: '',
  disabled: false,
  hasTemplates: false,
  templateCount: 0,
  stateFields: () => [],
  nodeOutputOptions: () => [],
  credentialId: '',
  fetchOptionsKey: '',
  fetchOptionsLabel: '',
});

const emit = defineEmits<FieldSourceSelectorEmits>();

/** COMPOSABLES & STORES */
const { states, externalInputs } = useWorkflowEditorState();
const pluginRegistry = usePluginRegistryStore();
const ts = useTS({ capitalize: true });
const assetStatusBasePath = 'pages.automations.createEditWorkflow.assetStatusFields';

/** COMPUTED */

/**
 * Current value with fallback
 */
const currentValue = computed<FieldSourceValue>(
  () => props.modelValue || DEFAULT_FIELD_SOURCE_VALUE,
);

/**
 * Filtered source type options based on allowedTypes prop
 */
const filteredTypeOptions = computed(() => {
  const allowed = new Set(props.allowedTypes);
  // Auto-include assetStatus when event is allowed
  if (allowed.has('event')) {
    allowed.add('assetStatus');
  }
  return SOURCE_TYPE_OPTIONS
    .filter(o => allowed.has(o.value))
    .map(o => {
      if (o.value === 'fetchOptions' && props.fetchOptionsLabel) {
        return { ...o, label: props.fetchOptionsLabel };
      }
      return o;
    });
});

/**
 * Whether to show the type selector (hidden when only 1 type allowed)
 */
const showTypeSelector = computed(() => filteredTypeOptions.value.length > 1);

/**
 * Current source type option config (icon, color)
 */
const currentTypeConfig = computed(() => {
  if (currentValue.value.type === 'event' && currentValue.value.mode === 'assetStatus') {
    return SOURCE_TYPE_OPTIONS.find(o => o.value === 'assetStatus') || SOURCE_TYPE_OPTIONS[3];
  }
  return SOURCE_TYPE_OPTIONS.find(o => o.value === currentValue.value.type) || SOURCE_TYPE_OPTIONS[3];
});

/**
 * External input variables mapped to select options (input source only)
 */
const inputOptions = computed(() =>
  externalInputs.value.map(v => ({
    label: `${v.label} (input.${v.field})`,
    value: v.field,
    icon: v.icon || 'input',
    color: 'cyan-6',
  })),
);

/**
 * State field names for autocomplete.
 * Uses props.stateFields when provided by the caller, otherwise falls back
 * to the global states from useWorkflowEditorState.
 */
const stateFieldNames = computed(() => {
  if (props.stateFields.length > 0) {
    return props.stateFields.map(f => f.name);
  }
  return states.value.map(v => v.field);
});

/**
 * Selected node output option (for nodeOutput type)
 */
const selectedNodeOption = computed(() =>
  props.nodeOutputOptions.find(n => n.id === currentValue.value.nodeId),
);

/**
 * Output hints for the selected node (from plugin registry)
 */
const selectedNodeOutputHints = computed(() => {
  if (!selectedNodeOption.value) return [];
  const nodeType = pluginRegistry.nodeTypeMap.get(selectedNodeOption.value.type);
  return nodeType?.availableOutputs ?? [];
});

/**
 * Formatted output hints string for display
 */
const outputHintText = computed(() => {
  if (selectedNodeOutputHints.value.length === 0) return '';
  return selectedNodeOutputHints.value
    .map(h => h.path === '*' ? `* (${h.description})` : h.path)
    .join(', ');
});

/**
 * Whether event input is readonly (dynamic mode)
 */
const isEventReadonly = computed(() =>
  currentValue.value.type === 'event'
  && (currentValue.value.mode || 'dynamic') === 'dynamic'
  && currentValue.value.mode !== 'assetStatus',
);

/**
 * Current event mode (dynamic or manual)
 */
const eventMode = computed(() => {
  if (currentValue.value.type !== 'event') return null;
  if (currentValue.value.mode === 'assetStatus') return 'assetStatus';
  return currentValue.value.mode || 'dynamic';
});

/** FUNCTIONS */

/**
 * Emit updated value with partial merge
 *
 * @param {Partial<FieldSourceValue>} partial - Partial value to merge
 */
function emitValue(partial: Partial<FieldSourceValue>): void {
  emit('update:modelValue', { ...currentValue.value, ...partial });
}

/**
 * Update source type and reset value
 *
 * @param {SourceType} type - New source type
 */
function updateType(type: SourceType): void {
  if (type === 'assetStatus') {
    emit('update:modelValue', { type: 'event', value: '', mode: 'assetStatus' });
    return;
  }
  const newValue: FieldSourceValue = { type, value: '' };
  if (type === 'event') {
    newValue.mode = 'dynamic';
  }
  emit('update:modelValue', newValue);
}

/**
 * Update the value string
 *
 * @param {string} val - New value
 */
function updateValue(val: string): void {
  emitValue({ value: val });
}

/**
 * Handle event input click — emit appropriate event for parent
 */
function handleEventClick(): void {
  if (eventMode.value !== 'dynamic') return;

  if (!props.hasTemplates) {
    emit('openTemplateSelector');
  } else {
    emit('openEventSelector');
  }
}

/**
 * Switch event mode to manual
 */
function switchToManualMode(): void {
  emitValue({ mode: 'manual' });
}

/**
 * Switch event mode to dynamic
 */
function switchToDynamicMode(): void {
  emitValue({ mode: 'dynamic' });
}

/**
 * Asset status field options with translated labels
 */
const translatedAssetStatusOptions = computed(() =>
  ASSET_STATUS_FIELD_OPTIONS.map((opt: AssetStatusFieldOption) => ({
    ...opt,
    translatedLabel: ts(`${assetStatusBasePath}.${opt.value}`),
    translatedHint: ts(`${assetStatusBasePath}.${opt.value}Hint`),
    availabilityHint: opt.availability === 'offline'
      ? ts(`${assetStatusBasePath}.offlineOnly`)
      : ts(`${assetStatusBasePath}.allEvents`),
  })),
);

/**
 * Update selected node for nodeOutput type
 *
 * @param {string} nodeId - Selected node ID
 */
function updateNodeId(nodeId: string): void {
  emitValue({ nodeId, value: '' });
}

/**
 * Update the output path for nodeOutput type
 *
 * @param {string} path - Output path
 */
function updateOutputPath(path: string): void {
  emitValue({ value: path });
}

/** ── fetchOptions state ── */

/**
 * Options fetched from the fetchOptions proxy endpoint
 */
const fetchOptionsItems = ref<Array<{ label: string; value: unknown }>>([]);

/**
 * Loading state for fetchOptions
 */
const fetchOptionsLoading = ref(false);

/**
 * Whether fetchOptions has the required props to fetch
 */
const canFetchOptions = computed(() =>
  !!props.credentialId && !!props.fetchOptionsKey,
);

/**
 * Fetch options from the fetchOptions proxy endpoint
 *
 * @returns {Promise<void>}
 */
async function fetchOptions(): Promise<void> {
  if (!canFetchOptions.value) {
    fetchOptionsItems.value = [];
    return;
  }
  fetchOptionsLoading.value = true;
  try {
    const response = await apis.workflows.http.post('/api/v1/load_options', {
      credentialId: props.credentialId,
      resourceKey: props.fetchOptionsKey,
      dependsOn: {},
    });
    const items = response?.data?.data ?? response?.data ?? [];
    fetchOptionsItems.value = (Array.isArray(items) ? items : []).map((item: { label: string; value: unknown }) => ({
      label: item.label,
      value: item.value,
    }));
  } catch (error) {
    console.error('[FieldSourceSelector] fetchOptions failed:', error);
    fetchOptionsItems.value = [];
  } finally {
    fetchOptionsLoading.value = false;
  }
}

/** WATCHERS */

watch(
  [() => props.credentialId, () => props.fetchOptionsKey, () => currentValue.value.type],
  () => {
    const isFetchType = currentValue.value.type === 'fetchOptions';
    if (isFetchType && canFetchOptions.value) {
      void fetchOptions();
    }
  },
  { immediate: true },
);

/** ── literal-template autocomplete ── */

/**
 * Static example syntax for each namespace. Lives outside i18n because
 * vue-i18n's message parser treats {{ }} as its own interpolation syntax,
 * which would conflict with showing literal {{namespace.path}} examples.
 */
const LITERAL_NAMESPACE_EXAMPLES = {
  event: '{{event.<path>}}',
  state: '{{state.<path>}}',
  input: '{{input.<path>}}',
  output: '{{output.<nodeId>.<path>}}',
} as const;

/**
 * Static autocomplete-insert text for each namespace prefix.
 */
const LITERAL_AUTOCOMPLETE_PREFIXES = {
  event: '{{event.',
  state: '{{state.',
  input: '{{input.',
  output: '{{output.',
} as const;

/**
 * Catalog of node types that populate {{output.<nodeId>.<field>}}.
 * Sourced from runtime/domain/executors/* — see workflow service docs/context.md.
 * Type names are technical (not translated); shape descriptions are pulled
 * from i18n at render time via outputCatalog.shapes.<type>.
 */
const LITERAL_OUTPUT_CATALOG_PRODUCING = [
  'loop', 'code', 'subworkflow', 'plugin', 'wait_signal', 'wait_for', 'trigger_event',
] as const;

const LITERAL_OUTPUT_CATALOG_EMPTY: readonly string[] = [
  'start', 'end', 'condition', 'switch', 'set_state',
  'goto', 'log', 'fanout', 'merge', 'sequence',
] as const;

/**
 * Whether the namespace autocomplete menu is currently open
 */
const literalAutocompleteOpen = ref(false);

/**
 * Watches the literal textarea content; opens the autocomplete menu when the
 * user just typed `{{` and there is no closing `}}` on the trailing segment.
 */
watch(
  () => currentValue.value.value,
  (newVal) => {
    if (currentValue.value.type !== 'literal') {
      literalAutocompleteOpen.value = false;
      return;
    }
    const lastOpen = newVal.lastIndexOf('{{');
    if (lastOpen === -1) {
      literalAutocompleteOpen.value = false;
      return;
    }
    const tail = newVal.slice(lastOpen);
    if (tail.includes('}}')) {
      literalAutocompleteOpen.value = false;
      return;
    }
    if (/\{\{[a-zA-Z]/.test(tail)) {
      literalAutocompleteOpen.value = false;
      return;
    }
    literalAutocompleteOpen.value = true;
  },
);

/**
 * Inserts a namespace prefix (event./state./input./output.) at the trailing
 * `{{` of the current value, replacing any whitespace between `{{` and the
 * caret. Closes the autocomplete menu after insertion.
 *
 * @param {'event' | 'state' | 'input' | 'output'} prefix - Namespace prefix
 */
function insertNamespacePrefix(prefix: 'event' | 'state' | 'input' | 'output'): void {
  const current = currentValue.value.value;
  const lastOpen = current.lastIndexOf('{{');
  if (lastOpen === -1) {
    literalAutocompleteOpen.value = false;
    return;
  }
  const head = current.slice(0, lastOpen);
  const newValue = `${head}{{${prefix}.`;
  updateValue(newValue);
  literalAutocompleteOpen.value = false;
}
</script>

<template>
  <div class="field-source-selector">
    <!-- Source type selector (hidden when only 1 type allowed) -->
    <q-select
      v-if="showTypeSelector"
      :model-value="currentValue.type === 'event' && currentValue.mode === 'assetStatus' ? 'assetStatus' : currentValue.type"
      :options="[...filteredTypeOptions]"
      outlined
      dense
      emit-value
      map-options
      options-dense
      option-value="value"
      option-label="label"
      :label="label || undefined"
      :disable="disabled"
      class="field-source-selector__type-select q-mb-sm"
      @update:model-value="(val: string) => updateType(val as SourceType)"
    >
      <template #prepend>
        <q-icon
          :name="currentTypeConfig?.icon || 'event'"
          :color="currentTypeConfig?.color || 'blue-6'"
          size="xs"
        />
      </template>
      <template #option="scope">
        <q-item v-bind="scope.itemProps">
          <q-item-section avatar>
            <q-icon :name="scope.opt.icon" :color="scope.opt.color" size="xs" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- EVENT type: three modes (assetStatus / dynamic / manual) -->
    <template v-if="currentValue.type === 'event'">
      <!-- ASSET STATUS MODE: q-select dropdown with predefined fields -->
      <template v-if="eventMode === 'assetStatus'">
        <q-select
          :model-value="currentValue.value || null"
          :options="[...translatedAssetStatusOptions]"
          outlined
          dense
          emit-value
          map-options
          option-value="value"
          option-label="translatedLabel"
          placeholder="Select health monitoring field..."
          :disable="disabled"
          @update:model-value="(val: string) => emitValue({ value: val })"
        >
          <template #prepend>
            <q-icon name="monitor_heart" color="orange-7" size="xs" />
          </template>
          <template #option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section avatar>
                <q-icon :name="scope.opt.icon" color="orange-7" size="xs" />
              </q-item-section>
              <q-item-section>
                <q-item-label>
                  {{ scope.opt.translatedLabel }}
                  <q-badge
                    :color="scope.opt.availability === 'offline' ? 'negative' : 'positive'"
                    :label="scope.opt.availabilityHint"
                    class="q-ml-sm"
                    dense
                  />
                </q-item-label>
                <q-item-label caption class="text-grey-6">
                  {{ scope.opt.translatedHint }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </template>
          <template #append>
            <q-btn flat round dense size="sm" icon="more_vert" color="grey-7" class="q-ml-xs" @click.stop>
              <q-menu anchor="bottom end" self="top end" :offset="[0, 4]">
                <q-list dense style="min-width: 240px; padding: 8px 0;">
                  <q-item clickable v-close-popup class="q-py-sm q-px-md" @click="switchToDynamicMode">
                    <q-item-section avatar><q-icon name="folder" color="primary" size="sm" /></q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">Browse Mode</q-item-label>
                      <q-item-label caption class="text-grey-7">Select from template fields</q-item-label>
                    </q-item-section>
                  </q-item>
                  <q-separator class="q-my-xs" />
                  <q-item clickable v-close-popup class="q-py-sm q-px-md" @click="switchToManualMode">
                    <q-item-section avatar><q-icon name="edit" color="orange-6" size="sm" /></q-item-section>
                    <q-item-section>
                      <q-item-label class="text-weight-medium">Manual Mode</q-item-label>
                      <q-item-label caption class="text-grey-7">Type field path manually</q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </template>
        </q-select>
      </template>

      <!-- DYNAMIC / MANUAL MODE (existing behavior) -->
      <template v-else>
        <q-input
          :model-value="currentValue.value"
          outlined
          dense
          :readonly="isEventReadonly"
          :placeholder="isEventReadonly ? 'Click to select...' : 'e.g. event.payload.status'"
          :class="{ 'cursor-pointer': isEventReadonly }"
          :disable="disabled"
          @click="handleEventClick"
          @update:model-value="(val: string | number | null) => updateValue(String(val ?? ''))"
        >
          <template #prepend>
            <q-icon name="event" color="blue-6" size="xs" />
          </template>
          <template #append>
            <q-chip v-if="hasTemplates" dense size="sm" color="blue-1" text-color="blue-9">
              <q-icon name="folder" size="xs" class="q-mr-xs" />{{ templateCount }}
            </q-chip>
            <q-btn flat round dense size="sm" icon="more_vert" color="grey-7" class="q-ml-xs" @click.stop>
              <q-menu anchor="bottom end" self="top end" :offset="[0, 4]">
                <q-list dense style="min-width: 240px; padding: 8px 0;">
                  <!-- DYNAMIC MODE OPTIONS -->
                  <template v-if="eventMode === 'dynamic'">
                    <q-item v-if="hasTemplates" clickable v-close-popup class="q-py-sm q-px-md" @click="emit('openEventSelector')">
                      <q-item-section avatar><q-icon name="search" color="primary" size="sm" /></q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Select Field</q-item-label>
                        <q-item-label caption class="text-grey-7">Browse {{ templateCount }} template{{ templateCount !== 1 ? 's' : '' }}</q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-separator v-if="hasTemplates" class="q-my-xs" />
                    <q-item clickable v-close-popup class="q-py-sm q-px-md" @click="emit('openTemplateSelector')">
                      <q-item-section avatar><q-icon name="add" color="primary" size="sm" /></q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Search Templates</q-item-label>
                        <q-item-label caption class="text-grey-7">Select asset templates</q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-separator class="q-my-xs" />
                    <q-item clickable v-close-popup class="q-py-sm q-px-md" @click="switchToManualMode">
                      <q-item-section avatar><q-icon name="edit" color="orange-6" size="sm" /></q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Manual Mode</q-item-label>
                        <q-item-label caption class="text-grey-7">Type field path manually</q-item-label>
                      </q-item-section>
                    </q-item>
                  </template>
                  <!-- MANUAL MODE OPTIONS -->
                  <template v-else-if="eventMode === 'manual'">
                    <q-item clickable v-close-popup class="q-py-sm q-px-md" @click="switchToDynamicMode">
                      <q-item-section avatar><q-icon name="folder" color="primary" size="sm" /></q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Browse Mode</q-item-label>
                        <q-item-label caption class="text-grey-7">Select from template fields</q-item-label>
                      </q-item-section>
                    </q-item>
                  </template>
                </q-list>
              </q-menu>
            </q-btn>
          </template>
        </q-input>
      </template>
    </template>

    <!-- STATE type: q-select autocomplete -->
    <q-select
      v-else-if="currentValue.type === 'state'"
      :model-value="currentValue.value"
      outlined
      dense
      use-input
      hide-selected
      fill-input
      input-debounce="0"
      :placeholder="placeholder || 'Select variable...'"
      :options="stateFieldNames"
      :disable="disabled"
      @update:model-value="(val: string) => updateValue(val)"
    >
      <template #prepend>
        <q-icon name="storage" color="purple-6" size="xs" />
      </template>
      <template #no-option>
        <q-item>
          <q-item-section class="text-grey-6 text-caption">No state variables defined</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- INPUT type: q-select with external input variables only -->
    <q-select
      v-else-if="currentValue.type === 'input' || currentValue.type === 'variable'"
      :model-value="currentValue.value"
      outlined
      dense
      use-input
      hide-selected
      fill-input
      input-debounce="0"
      :placeholder="placeholder || 'Select input...'"
      :options="inputOptions"
      option-label="label"
      option-value="value"
      emit-value
      map-options
      :disable="disabled"
      @update:model-value="(val: string) => updateValue(val)"
    >
      <template #prepend>
        <q-icon name="input" color="cyan-6" size="xs" />
      </template>
      <template #option="scope">
        <q-item v-bind="scope.itemProps">
          <q-item-section avatar>
            <q-icon :name="scope.opt.icon" :color="scope.opt.color" size="xs" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template #no-option>
        <q-item>
          <q-item-section class="text-grey-6 text-caption">No external inputs defined</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- NODE_OUTPUT type: node selector + output path -->
    <template v-else-if="currentValue.type === 'nodeOutput'">
      <!-- Node selector -->
      <q-select
        :model-value="currentValue.nodeId || null"
        :options="nodeOutputOptions"
        outlined
        dense
        emit-value
        map-options
        option-value="id"
        option-label="label"
        :placeholder="placeholder || 'Select node...'"
        :disable="disabled"
        class="q-mb-sm"
        @update:model-value="(val: string) => updateNodeId(val)"
      >
        <template #prepend>
          <q-icon name="hub" color="teal-6" size="xs" />
        </template>
        <template #append>
          <q-icon name="help_outline" size="xs" class="cursor-pointer">
            <q-menu anchor="bottom end" self="top end" :offset="[0, 8]">
              <q-list dense style="min-width: 320px; padding: 12px 16px;">
                <q-item-label class="text-weight-medium q-mb-sm">
                  {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.title') }}
                </q-item-label>
                <q-item-label caption class="text-weight-medium q-mt-sm q-mb-xs">
                  {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.producingHeader') }}
                </q-item-label>
                <q-item
                  v-for="nodeType in LITERAL_OUTPUT_CATALOG_PRODUCING"
                  :key="nodeType"
                  class="q-pa-none"
                >
                  <q-item-section>
                    <q-item-label class="text-weight-medium">
                      {{ nodeType }}
                    </q-item-label>
                    <q-item-label caption>
                      {{ ts(`components.forms.fieldSourceSelector.literal.outputCatalog.shapes.${nodeType}`) }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
                <q-separator class="q-my-sm" />
                <q-item-label caption class="text-weight-medium q-mb-xs">
                  {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.emptyHeader') }}
                </q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ LITERAL_OUTPUT_CATALOG_EMPTY.join(', ') }}
                </q-item-label>
                <q-separator class="q-my-sm" />
                <q-item-label caption class="text-grey-7">
                  {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.footerNote') }}
                </q-item-label>
              </q-list>
            </q-menu>
          </q-icon>
        </template>
        <template #option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section avatar>
              <q-icon name="hub" color="teal-6" size="xs" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ scope.opt.label }}</q-item-label>
              <q-item-label caption class="text-grey-6">
                {{ scope.opt.id }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </template>
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey-6 text-caption">No nodes available</q-item-section>
          </q-item>
        </template>
      </q-select>

      <!-- Output path input -->
      <q-input
        v-if="currentValue.nodeId"
        :model-value="currentValue.value"
        outlined
        dense
        placeholder="e.g. item, response.data"
        :hint="outputHintText ? `Available: ${outputHintText}` : 'Type the output path'"
        :disable="disabled"
        @update:model-value="(val: string | number | null) => updateOutputPath(String(val ?? ''))"
      >
        <template #prepend>
          <q-icon name="data_object" color="teal-6" size="xs" />
        </template>
      </q-input>
    </template>

    <!-- FETCH_OPTIONS / LOAD_OPTIONS type: dynamic dropdown from API -->
    <template v-else-if="currentValue.type === 'fetchOptions'">
      <q-select
        v-if="canFetchOptions"
        :model-value="currentValue.value || null"
        :options="fetchOptionsItems"
        :loading="fetchOptionsLoading"
        outlined
        dense
        emit-value
        map-options
        use-input
        input-debounce="300"
        option-value="value"
        option-label="label"
        :placeholder="placeholder || 'Select option...'"
        :disable="disabled"
        @update:model-value="(val: unknown) => updateValue(String(val ?? ''))"
        @filter="(val: string, update: (fn: () => void) => void) => update(() => {})"
      >
        <template #prepend>
          <q-icon name="cloud_download" color="purple-6" size="xs" />
        </template>
        <template #append>
          <q-btn flat dense round icon="refresh" size="xs" @click.stop="fetchOptions" />
        </template>
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey-6 text-caption">
              {{ fetchOptionsLoading ? 'Loading...' : 'No options available' }}
            </q-item-section>
          </q-item>
        </template>
      </q-select>
      <div v-else class="text-caption text-grey-6 q-pa-sm" style="border: 1px dashed var(--mapex-card-border); border-radius: var(--mapex-radius-sm);">
        <q-icon name="vpn_key" size="xs" class="q-mr-xs" />
        Select a credential first
      </div>
    </template>

    <!-- LITERAL type: multi-line q-input (accepts double-brace templates) -->
    <q-input
      v-else
      :model-value="currentValue.value"
      outlined
      dense
      type="textarea"
      autogrow
      :placeholder="placeholder || ts('components.forms.fieldSourceSelector.literal.placeholder')"
      :disable="disabled"
      input-style="max-height: 240px; overflow-y: auto;"
      @update:model-value="(val: string | number | null) => updateValue(String(val ?? ''))"
    >
      <q-menu
        v-model="literalAutocompleteOpen"
        no-focus
        no-parent-event
        anchor="bottom start"
        self="top start"
        :offset="[0, 4]"
      >
        <q-list dense style="min-width: 240px;">
          <q-item
            v-for="ns in (['event', 'state', 'input', 'output'] as const)"
            :key="ns"
            v-close-popup
            clickable
            @click="insertNamespacePrefix(ns)"
          >
            <q-item-section>
              <q-item-label class="text-weight-medium">
                {{ ts('components.forms.fieldSourceSelector.literal.autocompleteHintPrefix') }} {{ LITERAL_AUTOCOMPLETE_PREFIXES[ns] }}
              </q-item-label>
              <q-item-label caption>
                {{ ts(`components.forms.fieldSourceSelector.literal.namespaces.${ns}.hint`) }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-menu>
      <template #prepend>
        <q-icon
          :name="currentTypeConfig?.icon || 'format_quote'"
          :color="currentTypeConfig?.color || 'green-6'"
          size="xs"
        />
      </template>
      <template #append>
        <q-icon
          name="info"
          size="xs"
          class="cursor-pointer"
        >
          <q-menu anchor="bottom end" self="top end" :offset="[0, 8]">
            <q-list dense style="min-width: 320px; padding: 12px 16px;">
              <q-item-label class="text-weight-medium q-mb-sm">
                {{ ts('components.forms.fieldSourceSelector.literal.infoTitle') }}
              </q-item-label>
              <q-item-label caption class="q-mb-md">
                {{ ts('components.forms.fieldSourceSelector.literal.infoDescription') }}
              </q-item-label>
              <q-item
                v-for="ns in (['event', 'state', 'input', 'output'] as const)"
                :key="ns"
                class="q-pa-none q-mb-xs"
              >
                <q-item-section>
                  <q-item-label class="text-weight-medium">
                    {{ LITERAL_NAMESPACE_EXAMPLES[ns] }}
                  </q-item-label>
                  <q-item-label caption>
                    {{ ts(`components.forms.fieldSourceSelector.literal.namespaces.${ns}.hint`) }}
                  </q-item-label>
                </q-item-section>
                <q-item-section v-if="ns === 'output'" side>
                  <q-icon name="help_outline" size="xs" class="cursor-pointer">
                    <q-menu anchor="top end" self="top start" :offset="[8, 0]">
                      <q-list dense style="min-width: 320px; padding: 12px 16px;">
                        <q-item-label class="text-weight-medium q-mb-sm">
                          {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.title') }}
                        </q-item-label>
                        <q-item-label caption class="text-weight-medium q-mt-sm q-mb-xs">
                          {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.producingHeader') }}
                        </q-item-label>
                        <q-item
                          v-for="nodeType in LITERAL_OUTPUT_CATALOG_PRODUCING"
                          :key="nodeType"
                          class="q-pa-none"
                        >
                          <q-item-section>
                            <q-item-label class="text-weight-medium">
                              {{ nodeType }}
                            </q-item-label>
                            <q-item-label caption>
                              {{ ts(`components.forms.fieldSourceSelector.literal.outputCatalog.shapes.${nodeType}`) }}
                            </q-item-label>
                          </q-item-section>
                        </q-item>
                        <q-separator class="q-my-sm" />
                        <q-item-label caption class="text-weight-medium q-mb-xs">
                          {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.emptyHeader') }}
                        </q-item-label>
                        <q-item-label caption class="text-grey-7">
                          {{ LITERAL_OUTPUT_CATALOG_EMPTY.join(', ') }}
                        </q-item-label>
                        <q-separator class="q-my-sm" />
                        <q-item-label caption class="text-grey-7">
                          {{ ts('components.forms.fieldSourceSelector.literal.outputCatalog.footerNote') }}
                        </q-item-label>
                      </q-list>
                    </q-menu>
                  </q-icon>
                </q-item-section>
              </q-item>
              <q-separator class="q-my-sm" />
              <q-item-label caption class="text-grey-7">
                {{ ts('components.forms.fieldSourceSelector.literal.missingPathNote') }}
              </q-item-label>
              <q-item-label caption class="text-grey-7">
                {{ ts('components.forms.fieldSourceSelector.literal.escapeNote') }}
              </q-item-label>
            </q-list>
          </q-menu>
        </q-icon>
      </template>
    </q-input>
  </div>
</template>

<style lang="scss" scoped>
.field-source-selector {
  &__type-select {
    max-width: 100%;
  }
}

.cursor-pointer {
  cursor: pointer;
}
</style>
