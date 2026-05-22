<script setup lang="ts">
/** TYPE IMPORTS */
import type {
  DynamicNodeFormProps,
  DynamicNodeFormEmits,
  NodePropertyDefinition,
} from '../../interfaces/CreateEditWorkflow.interface';
import type { FieldSourceValue, SourceType, NodeOutputOption } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, nextTick } from 'vue';

/** COMPONENTS */
import { FieldSourceSelector } from '@components/forms';

/** UTILS */
import { buildDefaultConfig } from '@src/utils/workflow/buildDefaultConfig';

/** COMPOSABLES */
import { useWorkflowContext } from '@src/composables/workflow';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** PROPS & EMITS */
const props = defineProps<DynamicNodeFormProps>();
const emit = defineEmits<DynamicNodeFormEmits>();

/** COMPOSABLES & STORES */
const pluginRegistry = usePluginRegistryStore();
const { nodes, getNodeType } = useWorkflowContext();

/** Node output options — nodes with availableOutputs, excluding self */
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

/** STATE */
const form = ref<Record<string, unknown>>(buildInitialForm());

/**
 * Guard flag to prevent circular updates between props.config ↔ form watchers.
 * When true, the form watcher skips emitting back to the parent.
 */
let syncing = true;

/** Skip first emission — mount should not dirty the parent config */
void nextTick(() => { syncing = false; });

/** WATCHERS */
watch(() => props.config, (val) => {
  syncing = true;
  form.value = buildInitialForm(val);
  void nextTick(() => { syncing = false; });
}, { deep: true });

watch(form, (val, oldVal) => {
  if (syncing) return;

  // Operation changed — rebuild config from scratch for the new operation
  if (oldVal && val.operation !== oldVal.operation && val.operation && props.nodeType) {
    const nodeType = pluginRegistry.getNodeType(props.nodeType);
    if (nodeType) {
      const freshConfig = buildDefaultConfig(nodeType, val.operation as string);
      freshConfig.operation = val.operation;
      // Preserve credentialId
      if (val.credentialId !== undefined) freshConfig.credentialId = val.credentialId;
      syncing = true;
      form.value = buildInitialForm(freshConfig);
      void nextTick(() => { syncing = false; });
      emit('update:config', { ...freshConfig });
      return;
    }
  }

  emit('update:config', buildVisibleConfig());
}, { deep: true });

/** FUNCTIONS */

/**
 * Build form values from properties defaults merged with a config source.
 * The form keeps ALL properties internally (needed for displayOptions evaluation).
 *
 * @param {Record<string, unknown>} [configSource] - Config to merge from (defaults to props.config)
 * @returns {Record<string, unknown>} Form state
 */
function buildInitialForm(configSource?: Record<string, unknown>): Record<string, unknown> {
  const src = configSource ?? props.config;
  const result: Record<string, unknown> = {};
  for (const prop of props.properties) {
    result[prop.name] = src[prop.name] ?? prop.default;
  }
  // Carry over non-property keys (e.g., credentialId, __description)
  for (const [key, val] of Object.entries(src)) {
    if (!(key in result)) {
      result[key] = val;
    }
  }
  return result;
}

/**
 * Build config with only visible properties (based on displayOptions).
 * This is what gets emitted to the parent — no invisible field pollution.
 *
 * @returns {Record<string, unknown>} Config with only visible fields
 */
function buildVisibleConfig(): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  for (const prop of props.properties) {
    if (prop.type === 'notice') continue;
    if (isVisible(prop) || prop.type === 'hidden') {
      result[prop.name] = form.value[prop.name];
    }
  }
  // Preserve keys not from properties (e.g., credentialId, __description)
  for (const [key, val] of Object.entries(form.value)) {
    if (!(key in result) && !props.properties.some(p => p.name === key)) {
      result[key] = val;
    }
  }
  return result;
}

/**
 * Check if a property should be visible based on displayOptions.show
 *
 * @param {NodePropertyDefinition} prop - Property definition to check
 * @returns {boolean} Whether the property is visible
 */
function isVisible(prop: NodePropertyDefinition): boolean {
  if (prop.type === 'hidden') return false;
  if (!prop.displayOptions?.show) return true;
  return Object.entries(prop.displayOptions.show).every(
    ([field, values]) => (values as string[]).includes(form.value[field] as string),
  );
}

/**
 * Determine if a string field should render as textarea
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {boolean} Whether to use textarea
 */
function isMultiline(prop: NodePropertyDefinition): boolean {
  return prop.rendering?.multiline === true;
}

/**
 * Determine if a string field should render as password
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {boolean} Whether to use password input
 */
function isPassword(prop: NodePropertyDefinition): boolean {
  return prop.rendering?.password === true || prop.isSecret === true;
}

/**
 * Get the placeholder for a property
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {string | undefined} Placeholder text
 */
function getPlaceholder(prop: NodePropertyDefinition): string | undefined {
  return prop.rendering?.placeholder;
}

/**
 * Get textarea rows for a property
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {number} Number of rows
 */
function getRows(prop: NodePropertyDefinition): number {
  return prop.rendering?.rows ?? 3;
}

/**
 * Get min/max values for number fields
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {{ min: number | undefined; max: number | undefined }} Min/max values
 */
function getMinMax(prop: NodePropertyDefinition): { min: number | undefined; max: number | undefined } {
  return {
    min: prop.rendering?.min,
    max: prop.rendering?.max,
  };
}

/**
 * Get allowed source types for a fieldSource property
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {SourceType[]} Allowed source types
 */
function getAllowedSources(prop: NodePropertyDefinition): SourceType[] {
  const sources = prop.allowedSources ?? ['literal', 'state', 'event', 'assetStatus', 'nodeOutput'];

  if (prop.fetchOptions?.rules?.length && !sources.includes('fetchOptions')) {
    return [...sources, 'fetchOptions'];
  }
  return sources;
}

/**
 * Get the fetchOptions key and label for the current form state
 *
 * @param {NodePropertyDefinition} prop - Property definition
 * @returns {{ key: string; label: string }} Resolved fetchOptions config
 */
function resolveFetchOptions(prop: NodePropertyDefinition): { key: string; label: string } {
  // New format: fetchOptions.rules
  if (prop.fetchOptions?.rules?.length) {
    for (const rule of prop.fetchOptions.rules) {
      const matches = Object.entries(rule.when).every(([field, values]) =>
        values.includes(form.value[field]),
      );
      if (matches || Object.keys(rule.when).length === 0) {
        return { key: rule.key, label: rule.label };
      }
    }
  }
  return { key: '', label: '' };
}

/**
 * Get the credential ID from the current node config
 *
 * @returns {string} Credential ID or empty string
 */
function getCredentialId(): string {
  return (props.config.credentialId as string) ?? '';
}

/**
 * Handle FieldSourceSelector v-model update
 *
 * @param {string} name - Property name
 * @param {FieldSourceValue} value - New field source value
 */
function updateFieldSource(name: string, value: FieldSourceValue): void {
  form.value[name] = value;
}

/**
 * Get the banner color for a notice type
 *
 * @param {string | undefined} noticeType - Notice type
 * @returns {string} Quasar color string
 */
function getNoticeColor(noticeType?: string): string {
  switch (noticeType) {
    case 'warning': return 'warning';
    case 'success': return 'positive';
    default: return 'info';
  }
}

/**
 * Get the icon for a notice type
 *
 * @param {string | undefined} noticeType - Notice type
 * @returns {string} Material icon name
 */
function getNoticeIcon(noticeType?: string): string {
  switch (noticeType) {
    case 'warning': return 'warning';
    case 'success': return 'check_circle';
    default: return 'info';
  }
}
</script>

<template>
  <div class="dynamic-node-form">
    <template v-for="prop in properties" :key="prop.name">
      <div v-if="isVisible(prop)" class="q-mb-md">
        <!-- ── String (text, textarea, password) ── -->
        <q-input
          v-if="prop.type === 'string' && isMultiline(prop)"
          v-model="form[prop.name] as string"
          :label="prop.displayName"
          :hint="prop.hint"
          :placeholder="getPlaceholder(prop)"
          :rows="getRows(prop)"
          type="textarea"
          outlined
          dense
        />
        <q-input
          v-else-if="prop.type === 'string' && isPassword(prop)"
          v-model="form[prop.name] as string"
          :label="prop.displayName"
          :hint="prop.hint"
          :placeholder="getPlaceholder(prop)"
          type="password"
          outlined
          dense
        />
        <q-input
          v-else-if="prop.type === 'string'"
          v-model="form[prop.name] as string"
          :label="prop.displayName"
          :hint="prop.hint"
          :placeholder="getPlaceholder(prop)"
          outlined
          dense
        />

        <!-- ── Number ── -->
        <q-input
          v-else-if="prop.type === 'number'"
          v-model.number="form[prop.name] as number"
          :label="prop.displayName"
          :hint="prop.hint"
          :placeholder="getPlaceholder(prop)"
          :min="getMinMax(prop).min"
          :max="getMinMax(prop).max"
          type="number"
          outlined
          dense
        />

        <!-- ── Boolean ── -->
        <q-toggle
          v-else-if="prop.type === 'boolean'"
          v-model="form[prop.name] as boolean"
          :label="prop.displayName"
        />

        <!-- ── Options (single select) ── -->
        <q-select
          v-else-if="prop.type === 'options'"
          v-model="form[prop.name]"
          :label="prop.displayName"
          :hint="prop.hint"
          :options="prop.options"
          emit-value
          map-options
          outlined
          dense
        />

        <!-- ── MultiOptions (multi select) ── -->
        <q-select
          v-else-if="prop.type === 'multiOptions'"
          v-model="form[prop.name]"
          :label="prop.displayName"
          :hint="prop.hint"
          :options="prop.options"
          emit-value
          map-options
          multiple
          use-chips
          outlined
          dense
        />

        <!-- ── JSON (textarea) ── -->
        <q-input
          v-else-if="prop.type === 'json'"
          v-model="form[prop.name] as string"
          :label="prop.displayName"
          :hint="prop.hint"
          type="textarea"
          outlined
          dense
          autogrow
        />

        <!-- ── FieldSource (source type selector) ── -->
        <div v-else-if="prop.type === 'fieldSource'">
          <div class="text-caption text-weight-medium q-mb-xs" style="color: var(--mapex-text-secondary)">
            {{ prop.displayName }}
            <span v-if="prop.required" class="text-negative">*</span>
          </div>
          <FieldSourceSelector
            :model-value="(form[prop.name] as FieldSourceValue) ?? { type: 'literal', value: '' }"
            :allowed-types="getAllowedSources(prop)"
            :placeholder="getPlaceholder(prop) ?? ''"
            :credential-id="getCredentialId()"
            :fetch-options-key="resolveFetchOptions(prop).key"
            :fetch-options-label="resolveFetchOptions(prop).label"
            :node-output-options="nodeOutputOptions"
            @update:model-value="(val: FieldSourceValue) => updateFieldSource(prop.name, val)"
          />
          <div v-if="prop.hint" class="text-caption q-mt-xs" style="color: var(--mapex-text-muted)">
            {{ prop.hint }}
          </div>
        </div>

        <!-- ── DateTime ── -->
        <q-input
          v-else-if="prop.type === 'dateTime'"
          v-model="form[prop.name] as string"
          :label="prop.displayName"
          :hint="prop.hint"
          outlined
          dense
        >
          <template #prepend>
            <q-icon name="event" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-date
                  :model-value="form[prop.name] as string"
                  mask="YYYY-MM-DD"
                  @update:model-value="(val: string | null) => { form[prop.name] = val ?? '' }"
                />
              </q-popup-proxy>
            </q-icon>
          </template>
          <template v-if="!prop.rendering?.dateOnly" #append>
            <q-icon name="access_time" class="cursor-pointer">
              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                <q-time
                  :model-value="form[prop.name] as string"
                  mask="YYYY-MM-DD HH:mm"
                  format24h
                  @update:model-value="(val: string | null) => { form[prop.name] = val ?? '' }"
                />
              </q-popup-proxy>
            </q-icon>
          </template>
        </q-input>

        <!-- ── Notice (informational banner) ── -->
        <q-banner
          v-else-if="prop.type === 'notice'"
          :class="`bg-${getNoticeColor(prop.noticeType)}`"
          rounded
          dense
          class="text-white"
        >
          <template #avatar>
            <q-icon :name="getNoticeIcon(prop.noticeType)" color="white" />
          </template>
          {{ prop.default }}
        </q-banner>
      </div>
    </template>
  </div>
</template>
