<script setup lang="ts">
/** TYPE IMPORTS */
import type { NodeConfigPanelProps, NodeConfigPanelEmits } from './interfaces/NodeConfigPanel.interface';
import type { PluginNodeType, WorkflowNode } from '../../interfaces/CreateEditWorkflow.interface';
import type { IconSectionNavItem } from '@components/navigation';
import type { NodeTimeoutConfig, NodeErrorHandlerConfig } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { computed, ref, onMounted } from 'vue';

/** COMPONENTS */
import DynamicNodeForm from '../DynamicNodeForm/DynamicNodeForm.vue';
import { IconSectionNav } from '@components/navigation';
import { CredentialSelector } from '../CredentialSelector';
import { TimeoutConfig } from '@components/workflow/TimeoutConfig';
import { ErrorHandlerConfig } from '@components/workflow/ErrorHandlerConfig';

/** COMPOSABLES */
import { useWorkflowEditorState, useWorkflowHistory } from '../../composables';
import { useCreateEditWorkflowTranslations } from '@composables/i18n/pages/automations/workflows';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { CONFIG_NAV_WIDTH } from './constants';

/** PROPS & EMITS */
const props = defineProps<NodeConfigPanelProps>();
const emit = defineEmits<NodeConfigPanelEmits>();

/** COMPOSABLES & STORES */
const t = useCreateEditWorkflowTranslations();
const pluginRegistry = usePluginRegistryStore();
const { nodes, edges, updateNodeConfig, replaceNodeConfig, removeNode, nodeValidationErrors, nodeConfigVersion } = useWorkflowEditorState();
const { pushSnapshot } = useWorkflowHistory();

/** STATE */
const activeSection = ref('config');

/** COMPUTED */

/**
 * Currently selected node data
 */
const selectedNode = computed<WorkflowNode | undefined>(() =>
  nodes.value.find(n => n.id === props.nodeId),
);

/**
 * Node type definition from plugin registry
 */
const nodeType = computed<PluginNodeType | undefined>(() => {
  if (!selectedNode.value) return undefined;
  return pluginRegistry.getNodeType(selectedNode.value.type);
});

/**
 * Config component from plugin (for dynamic form rendering)
 */
const configComponent = computed(() => nodeType.value?.configComponent);

/**
 * Credential definition from the parent plugin (if the plugin requires credentials)
 */
const credentialDefs = computed(() => {
  if (!nodeType.value?._pluginId) return [];
  const plugin = pluginRegistry.plugins.get(nodeType.value._pluginId);
  return plugin?.credentials ?? [];
});

/**
 * Plugin name for the credential selector dialog
 */
const pluginName = computed(() => {
  if (!nodeType.value?._pluginId) return '';
  const plugin = pluginRegistry.plugins.get(nodeType.value._pluginId);
  return plugin?.name ?? '';
});

/**
 * Config object enriched with internal _nodeId for self-filtering in nodeOutput selectors
 */
const enrichedConfig = computed(() => {
  if (!selectedNode.value) return {};
  return { ...selectedNode.value.config, _nodeId: props.nodeId };
});

/**
 * Count of incoming connections
 */
const incomingConnections = computed(() =>
  edges.value.filter(e => e.target === props.nodeId).length,
);

/**
 * Count of outgoing connections
 */
const outgoingConnections = computed(() =>
  edges.value.filter(e => e.source === props.nodeId).length,
);

/**
 * Raw validation error keys for the currently selected node
 */
const nodeErrors = computed<string[]>(() =>
  nodeValidationErrors.value[props.nodeId] ?? [],
);

/**
 * Translated validation errors resolved through i18n
 */
const translatedErrors = computed<string[]>(() =>
  nodeErrors.value.map(key => t.resolveValidationError(key)),
);

/**
 * Node description text from config
 */
const descriptionValue = computed({
  get: () => (selectedNode.value?.config.__description as string) || '',
  set: (val: string) => updateNodeConfig(props.nodeId, { __description: val }),
});

/**
 * Section navigation items for the icon rail.
 * Built-in sections: config + notes. Future: plugin-provided sections.
 */
const sectionItems = computed<IconSectionNavItem[]>(() => {
  const items: IconSectionNavItem[] = [
    {
      name: 'config',
      icon: 'tune',
      tooltip: t.nodeConfig.configTab.value,
    },
  ];

  // Timeout section (only for async nodes)
  if (hasTimeout.value) {
    items.push({
      name: 'timeout',
      icon: 'timer',
      tooltip: t.timeout.sectionTitle.value,
    });
  }

  // Error handler section (only for nodes with error output handle)
  if (hasErrorHandler.value) {
    items.push({
      name: 'errorHandler',
      icon: 'replay',
      tooltip: t.errorHandler.sectionTitle.value,
    });
  }

  items.push({
    name: 'notes',
    icon: 'description',
    tooltip: t.nodeConfig.notesTab.value,
    badge: !!descriptionValue.value,
  });

  return items;
});

/**
 * Whether the node type supports timeout configuration (is async)
 */
const hasTimeout = computed(() => !!nodeType.value?.timeout);

/**
 * Current node timeout config (from node instance, fallback to nodeType default)
 */
const nodeTimeout = computed<NodeTimeoutConfig>(() => {
  const nodeVal = selectedNode.value?.timeout;
  if (nodeVal) return nodeVal;
  const typeDefault = nodeType.value?.timeout;
  if (typeDefault) return { ...typeDefault };
  return { duration: 30, unit: 'seconds', enableOutput: false };
});

/**
 * Handle timeout config changes
 * @param {NodeTimeoutConfig} timeout - Updated timeout
 */
function handleTimeoutUpdate(timeout: NodeTimeoutConfig): void {
  if (!selectedNode.value) return;
  selectedNode.value.timeout = timeout;
  // Increment configVersion so WorkflowCanvas recalculates handles (timeout output)
  nodeConfigVersion.value++;
}

/**
 * Whether the node type supports error handler (has an "error" output handle)
 */
const hasErrorHandler = computed(() => {
  if (!nodeType.value?.outputs) return false;
  return nodeType.value.outputs.some(o => o.id === 'error');
});

/**
 * Default error handler config
 */
const DEFAULT_ERROR_HANDLER: NodeErrorHandlerConfig = {
  enabled: false,
  maxAttempts: 3,
  initialInterval: 5,
  intervalUnit: 'seconds',
  backoffMultiplier: 2.0,
};

/**
 * Current node error handler config (from node instance, fallback to nodeType default or global default)
 */
const nodeErrorHandler = computed<NodeErrorHandlerConfig>(() => {
  const nodeVal = selectedNode.value?.errorHandler;
  if (nodeVal) return nodeVal;
  const typeDefault = nodeType.value?.errorHandler;
  if (typeDefault) return { ...typeDefault };
  return { ...DEFAULT_ERROR_HANDLER };
});

/**
 * Handle error handler config changes
 * @param {NodeErrorHandlerConfig} errorHandler - Updated error handler
 */
function handleErrorHandlerUpdate(errorHandler: NodeErrorHandlerConfig): void {
  if (!selectedNode.value) return;
  selectedNode.value.errorHandler = errorHandler;
}

/** FUNCTIONS */

/**
 * Track previous operation to detect operation changes
 */
let lastOperation: unknown = selectedNode.value?.config?.operation;

/**
 * Handle config changes from the dynamic form.
 * Strips the runtime-only _nodeId before persisting.
 * When operation changes, replaces config entirely instead of merging.
 *
 * @param {Record<string, unknown>} config - Updated config
 * @returns {void}
 */
function handleConfigUpdate(config: Record<string, unknown>): void {
  const cleanConfig = Object.fromEntries(
    Object.entries(config).filter(([key]) => key !== '_nodeId'),
  );

  const operationChanged = cleanConfig.operation !== lastOperation;
  lastOperation = cleanConfig.operation;

  if (operationChanged) {
    replaceNodeConfig(props.nodeId, cleanConfig);
  } else {
    updateNodeConfig(props.nodeId, cleanConfig);
  }
}

/**
 * Handle node deletion
 *
 * @returns {void}
 */
function handleDelete(): void {
  pushSnapshot('Delete node');
  removeNode(props.nodeId);
  emit('close');
}

/** LIFECYCLE HOOKS */
onMounted(() => {
  pushSnapshot('Edit config');
});
</script>

<template>
  <div class="node-config-panel">
    <template v-if="selectedNode && nodeType">
      <div class="node-config-panel__body">
        <!-- Left: Icon section nav -->
        <IconSectionNav
          v-model="activeSection"
          :items="sectionItems"
          :width="CONFIG_NAV_WIDTH"
        />

        <!-- Right: Content column -->
        <div class="node-config-panel__content">
          <!-- Header: icon + title left, info + close right -->
          <div class="node-config-panel__header">
            <q-icon :name="nodeType.icon" :color="nodeType.color" size="24px" />
            <div class="q-ml-sm" style="min-width: 0;">
              <div class="text-subtitle1 text-weight-medium ellipsis">{{ nodeType.label }}</div>
              <div class="text-caption text-grey-6 ellipsis">{{ selectedNode.type }}</div>
            </div>
            <q-space />
            <!-- Node Info Popover -->
            <q-btn flat dense round icon="info_outline" size="sm" color="grey-6">
              <q-menu anchor="bottom right" self="top right">
                <q-card flat bordered style="min-width: 240px;">
                  <q-card-section class="q-pa-sm">
                    <div class="text-overline text-grey-7 q-mb-xs">{{ t.nodeConfig.nodeInfo.value }}</div>
                    <div class="text-caption text-grey-6">
                      <div class="q-mb-xs"><span class="text-weight-medium">{{ t.nodeConfig.nodeId.value }}</span> {{ selectedNode.id }}</div>
                      <div class="q-mb-xs"><span class="text-weight-medium">{{ t.nodeConfig.connections.value }}</span> {{ incomingConnections }} in, {{ outgoingConnections }} out</div>
                      <div><span class="text-weight-medium">{{ t.nodeConfig.position.value }}</span> {{ Math.round(selectedNode.position.x) }}, {{ Math.round(selectedNode.position.y) }}</div>
                    </div>
                  </q-card-section>
                </q-card>
              </q-menu>
            </q-btn>
            <q-btn flat dense round icon="close" size="sm" @click="emit('close')" />
          </div>

          <!-- Config section -->
          <div v-show="activeSection === 'config'" class="node-config-panel__form q-pa-md">
            <!-- Validation errors banner -->
            <q-banner v-if="nodeErrors.length > 0" rounded dense class="node-config-panel__error-banner q-mb-md">
              <template #avatar>
                <q-icon name="error" color="negative" />
              </template>
              <div class="node-config-panel__error-title">
                {{ t.nodeConfig.validationErrorCount(nodeErrors.length) }}
              </div>
              <ul class="node-config-panel__error-list">
                <li v-for="(err, i) in translatedErrors" :key="i">{{ err }}</li>
              </ul>
            </q-banner>
            <!-- Credential selector (when plugin requires credentials) -->
            <CredentialSelector
              v-if="credentialDefs.length > 0 && nodeType?._pluginId"
              :plugin-id="nodeType._pluginId"
              :plugin-name="pluginName"
              :credential-defs="credentialDefs"
              :model-value="(selectedNode?.config?.credentialId as string) ?? null"
              @update:model-value="handleConfigUpdate({ credentialId: $event })"
            />

            <!-- Declarative form from properties[] (preferred) -->
            <DynamicNodeForm
              v-if="nodeType?.properties?.length"
              :properties="nodeType.properties"
              :config="enrichedConfig"
              :node-type="selectedNode?.type"
              @update:config="handleConfigUpdate"
            />

            <!-- Legacy: plugin config component -->
            <component
              v-else-if="configComponent"
              :is="configComponent"
              :config="enrichedConfig"
              @update:config="handleConfigUpdate"
            />

            <!-- Fallback: JSON editor for config -->
            <div v-else>
              <q-input
                :model-value="JSON.stringify(selectedNode.config, null, 2)"
                type="textarea"
                outlined
                dense
                autogrow
                readonly
                :label="t.nodeConfig.configJson.value"
              />
              <div class="text-caption text-grey-6 q-mt-sm">
                {{ t.nodeConfig.noConfigForm.value }}
              </div>
            </div>
          </div>

          <!-- Timeout section (async nodes only) -->
          <div v-show="activeSection === 'timeout'" class="node-config-panel__form q-pa-md">
            <TimeoutConfig
              v-if="hasTimeout"
              :model-value="nodeTimeout"
              @update:model-value="handleTimeoutUpdate"
            />
          </div>

          <!-- Error Handler section (nodes with error output only) -->
          <div v-show="activeSection === 'errorHandler'" class="node-config-panel__form q-pa-md">
            <ErrorHandlerConfig
              v-if="hasErrorHandler"
              :model-value="nodeErrorHandler"
              @update:model-value="handleErrorHandlerUpdate"
            />
          </div>

          <!-- Notes section -->
          <div v-show="activeSection === 'notes'" class="node-config-panel__form q-pa-md">
            <div class="node-config-panel__notes-header q-mb-sm">
              <q-icon name="info" size="16px" color="primary" class="q-mr-xs" />
              <span class="text-caption text-grey-6">
                {{ t.nodeConfig.notesInfo.value }}
              </span>
            </div>
            <q-input
              v-model="descriptionValue"
              type="textarea"
              outlined
              dense
              autogrow
              :placeholder="t.nodeConfig.notesPlaceholder.value"
              :input-style="{ minHeight: '120px' }"
            />
          </div>

          <q-separator />

          <!-- Actions -->
          <div class="node-config-panel__actions q-pa-md">
            <div class="row q-gutter-sm">
              <q-btn
                flat
                no-caps
                color="negative"
                icon="delete"
                :label="t.nodeConfig.delete.value"
                class="col"
                @click="handleDelete"
              />
              <q-btn
                unelevated
                no-caps
                color="primary"
                icon="check"
                :label="t.nodeConfig.apply.value"
                class="col"
                @click="emit('close')"
              />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Node not found -->
    <template v-else>
      <div class="node-config-panel__header">
        <q-space />
        <q-btn flat dense round icon="close" size="sm" @click="emit('close')" />
      </div>
      <div class="text-center q-pa-xl text-grey-6">
        <q-icon name="warning" size="32px" />
        <p class="text-caption q-mt-sm">{{ t.nodeConfig.nodeNotFound.value }}</p>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.node-config-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--mapex-surface-bg);

  &__body {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  &__content {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }

  &__header {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid var(--mapex-card-border);
    flex-shrink: 0;
  }

  &__form {
    flex: 1;
    overflow-y: auto;
  }

  &__actions {
    flex-shrink: 0;
  }

  &__notes-header {
    display: flex;
    align-items: flex-start;
  }

  &__error-banner {
    background: color-mix(in srgb, var(--q-negative) 10%, var(--mapex-surface-bg));
    border: 1px solid color-mix(in srgb, var(--q-negative) 30%, transparent);
  }

  &__error-title {
    font-size: var(--mapex-font-xs);
    font-weight: var(--mapex-font-weight-medium);
    color: var(--q-negative);
  }

  &__error-list {
    padding-left: var(--mapex-spacing-lg);
    margin: var(--mapex-spacing-xs) 0;
    font-size: var(--mapex-font-2xs);
    color: var(--mapex-text-secondary);
  }
}
</style>
