<script setup lang="ts">
defineOptions({ name: 'ExecutionDetailModal' });

/** TYPE IMPORTS */
import type { WorkflowExecutionItem } from '../../interfaces';
import type { WorkflowNode, WorkflowEdge, IWorkflowEditorContext } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed, watch, provide, markRaw, nextTick } from 'vue';
import { VueFlow, useVueFlow } from '@vue-flow/core';
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';
import '@vue-flow/controls/dist/style.css';

/** COMPONENT IMPORTS */
import { ExecutionNodeWrapper } from '../ExecutionNodeWrapper';
import { TextNoteNode } from '@src/components/workflow/nodes/textNoteNode';
import GroupFrameNode from '@src/components/workflow/nodes/groupFrameNode/GroupFrameNode.vue';

/** SERVICE IMPORTS */
import { apis } from '@services/mapex';
import { useLogger } from '@src/composables/useLogger';
import { bootMarketplacePlugins } from '@src/pages/automations/workflows/createEditWorkflowPage/utils/manifestLoader';

/** COMPOSABLE IMPORTS */
import { WORKFLOW_CONTEXT_KEY } from '@src/composables/workflow';
import { useWorkflowExecutionsPageTranslations } from '@src/composables/i18n/pages/logs/workflowExecutionsPage';

/** STORE IMPORTS */
import { usePluginRegistryStore } from '@src/stores/pluginRegistry';

/** CONSTANT IMPORTS */
import { STATUS_COLORS, STATUS_ICONS } from '../../constants';

/** PROPS & EMITS */
const props = defineProps<{
  execution: WorkflowExecutionItem;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

/** COMPOSABLES & STORES */
const logger = useLogger('ExecutionDetailModal');
const t = useWorkflowExecutionsPageTranslations();
const registry = usePluginRegistryStore();
const { fitView } = useVueFlow({ id: 'execution-viewer' });

/** CONSTANT IMPORTS */
import { HOT_STATUSES } from '../../constants';

/** STATE */
const isLoading = ref(true);
const pluginsReady = ref(false);
const activeTab = ref('dag');
const selectedNodeId = ref<string | null>(null);
const showStepPanel = ref(false);
const definitionNodes = ref<WorkflowNode[]>([]);
const definitionEdges = ref<WorkflowEdge[]>([]);
const definitionStates = ref<Array<{ field: string; type: string; defaultValue?: unknown; description?: string; durable?: boolean }>>([]);
const fullExecution = ref<Record<string, unknown> | null>(null);

/** PROVIDE — Read-only workflow context for GenericWorkflowNode */
provide<IWorkflowEditorContext>(WORKFLOW_CONTEXT_KEY, {
  nodes: computed(() => definitionNodes.value),
  edges: computed(() => definitionEdges.value),
  states: computed(() => []),
  externalSignals: computed(() => []),
  updateNodeConfig: () => { /* read-only */ },
  addNoteToNode: () => { /* read-only */ },
  pushSnapshot: () => { /* read-only */ },
  getNodeType: (type: string) => registry.getNodeType(type),
});

/** COMPUTED — Parse execution data (fullExecution first, fallback to props.execution) */

/**
 * Parse a field that may be a JSON string or already an object/array
 */
function parseField(value: unknown): unknown {
  if (!value) return null;
  if (typeof value === 'string') {
    try { return JSON.parse(value); } catch { return null; }
  }
  return value;
}

const executionPath = computed(() => {
  const full = fullExecution.value;
  const raw = full?.executionPath ?? full?.ExecutionPath ?? props.execution?.executionPath;
  const parsed = parseField(raw);
  return Array.isArray(parsed) ? parsed : [];
});

const nodeOutputs = computed(() => {
  const full = fullExecution.value;
  const raw = full?.nodeOutputs ?? full?.NodeOutputs ?? props.execution?.nodeOutputs;
  const parsed = parseField(raw);
  return (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) ? parsed as Record<string, unknown> : {};
});

const errorInfo = computed(() => {
  const full = fullExecution.value;
  const raw = full?.errorInfo ?? full?.ErrorInfo ?? props.execution?.errorInfo;
  return parseField(raw) as Record<string, unknown> | null;
});

const executionState = computed(() => {
  const full = fullExecution.value;
  const raw = full?.state ?? full?.State;
  const parsed = parseField(raw);
  return (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) ? parsed as Record<string, unknown> : {};
});

/**
 * Merge definition state metadata with execution current values.
 * Each entry has: field, type, description, durable, defaultValue, currentValue, changed.
 */
const stateVariables = computed(() => {
  const state = executionState.value;
  const defStates = definitionStates.value;

  if (defStates.length === 0 && Object.keys(state).length === 0) return [];

  // If definition states exist, use them as the source of truth for metadata
  if (defStates.length > 0) {
    return defStates.map(ds => {
      const currentValue = state[ds.field] ?? ds.defaultValue;
      const changed = state[ds.field] !== undefined && JSON.stringify(state[ds.field]) !== JSON.stringify(ds.defaultValue);
      return {
        field: ds.field,
        type: ds.type || 'string',
        description: ds.description || '',
        durable: ds.durable ?? false,
        defaultValue: ds.defaultValue,
        currentValue,
        changed,
      };
    });
  }

  // Fallback: no definition states, just show execution state as-is
  return Object.entries(state).map(([field, val]) => ({
    field,
    type: typeof val === 'number' ? 'number' : typeof val === 'boolean' ? 'boolean' : 'string',
    description: '',
    durable: false,
    defaultValue: '',
    currentValue: val,
    changed: false,
  }));
});

const eventPayload = computed(() => {
  const full = fullExecution.value;
  const raw = full?.eventPayload ?? full?.EventPayload ?? props.execution?.eventPayload;
  return parseField(raw) as Record<string, unknown> | null;
});

const nodeStatusMap = computed(() => {
  const map: Record<string, { status: string; durationMs: number; error?: string; outputHandle?: string }> = {};
  for (const entry of executionPath.value) {
    map[entry.nodeId] = { status: entry.status, durationMs: entry.durationMs || 0, error: entry.error, outputHandle: entry.outputHandle };
  }
  return map;
});

/** FUNCTIONS — Status helpers */

/**
 * Icon for each node status emitted by the workflow runtime.
 * Backend sources: runtime_walker.go, runtime_handler_resume.go, runtime_lifecycle.go
 */
function getStepIcon(status: string): string {
  const m: Record<string, string> = {
    completed: 'check_circle',
    waiting:   'hourglass_top',
    retrying:  'replay',
    timeout:   'timer_off',
    error:     'error',
    cancelled: 'cancel',
  };
  return m[status] || 'radio_button_unchecked';
}

/**
 * Quasar color token for each node status (used by side-panel badge and logs icon).
 */
function getStepColor(status: string): string {
  const m: Record<string, string> = {
    completed: 'positive',  // green
    waiting:   'orange-7',  // orange (suspended, will resume)
    retrying:  'amber-6',   // yellow-amber (about to re-run)
    timeout:   'yellow-8',  // yellow (gave up waiting)
    error:     'negative',  // red
    cancelled: 'grey-7',    // dark grey (killed by system, not a failure)
  };
  return m[status] || 'grey-5';
}

function formatDuration(ms: number): string {
  if (!ms || ms <= 0) return '-';
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  return `${(ms / 60000).toFixed(1)}m`;
}

/**
 * Format a camelCase status into a human-readable uppercase label.
 * e.g. "notReached" → "NOT REACHED", "completed" → "COMPLETED"
 */
function formatStatusLabel(status: string): string {
  return status.replace(/([A-Z])/g, ' $1').toUpperCase();
}

function formatJson(data: unknown): string {
  if (!data) return '';
  try { return JSON.stringify(data, null, 2); } catch {
    if (typeof data === 'object') return JSON.stringify(data);
    if (typeof data === 'string' || typeof data === 'number' || typeof data === 'boolean') return String(data);
    return '';
  }
}

function copyToClipboard(text: string): void {
  void window.navigator.clipboard.writeText(text);
}

/** COMPUTED — Vue Flow nodeTypes from plugin registry (same as WorkflowCanvas) */

const nodeTypes = computed(() => {
  if (!pluginsReady.value) return {};
  const types: Record<string, unknown> = {};
  for (const [type] of registry.nodeTypeMap.entries()) {
    types[type] = markRaw(ExecutionNodeWrapper);
  }
  for (const node of definitionNodes.value) {
    if (!types[node.type]) {
      types[node.type] = markRaw(ExecutionNodeWrapper);
    }
  }
  // Annotations render as their own component (no execution wrapper)
  types['core/text_note'] = markRaw(TextNoteNode);
  types['core/group_frame'] = markRaw(GroupFrameNode);
  return markRaw(types);
});

/** COMPUTED — Vue Flow nodes with execution status border styling */

const flowNodes = computed(() => {
  if (!pluginsReady.value) return [];
  return definitionNodes.value.map((node) => {
    const isAnnotation = node.type === 'core/text_note' || node.type === 'core/group_frame';
    const visited = !isAnnotation && !!nodeStatusMap.value[node.id];
    const status = nodeStatusMap.value[node.id]?.status || '';
    const selected = selectedNodeId.value === node.id;

    // Aligned with runtime emitted statuses: completed/waiting/retrying/timeout/error/cancelled.
    // "failed" and "running" are NOT emitted at node level — removed to avoid dead branches.
    const STATUS_CLASS: Record<string, string> = {
      completed: 'exec-node-completed',
      waiting:   'exec-node-waiting',
      retrying:  'exec-node-retrying',
      timeout:   'exec-node-timeout',
      error:     'exec-node-error',
      cancelled: 'exec-node-cancelled',
    };

    const classes: string[] = [];
    if (!isAnnotation) {
      classes.push(visited ? (STATUS_CLASS[status] || 'exec-node-off') : 'exec-node-off');
    }
    if (selected) classes.push('exec-node-selected');

    return {
      id: node.id,
      type: node.type,
      position: node.position || { x: 0, y: 0 },
      data: { config: node.config, label: node.label ?? '', __nodeType: node.type, __execStatus: status },
      draggable: false,
      connectable: false,
      class: classes.join(' '),
    };
  });
});

/** COMPUTED — Vue Flow edges with visited path highlighting */

/**
 * Build a set of actually traversed edges by walking the executionPath sequence.
 * For each consecutive pair (nodeA → nodeB), find the edge that connects them
 * using nodeA's outputHandle (or infer from the edge if outputHandle is missing).
 */
const traversedEdgeIds = computed(() => {
  const path = executionPath.value;
  const ids = new Set<string>();

  for (let i = 0; i < path.length - 1; i++) {
    const currentNodeId = path[i].nodeId;
    const nextNodeId = path[i + 1].nodeId;
    const usedHandle = path[i].outputHandle;

    // Find the edge that connects current → next
    for (const edge of definitionEdges.value) {
      if (edge.source !== currentNodeId || edge.target !== nextNodeId) continue;

      if (usedHandle) {
        // Backend provided the outputHandle — match exactly
        if ((edge.sourceHandle || 'out') === usedHandle) {
          ids.add(edge.id);
          break;
        }
      } else {
        // Fallback for old data without outputHandle — infer by source→target match
        ids.add(edge.id);
        break;
      }
    }
  }

  return ids;
});

const flowEdges = computed(() =>
  definitionEdges.value.filter(e => e.sourceHandle !== '__note_out').map((edge) => {
    const visited = traversedEdgeIds.value.has(edge.id);
    return {
      id: edge.id,
      source: edge.source,
      ...(edge.sourceHandle ? { sourceHandle: edge.sourceHandle } : {}),
      target: edge.target,
      ...(edge.targetHandle ? { targetHandle: edge.targetHandle } : {}),
      animated: visited,
      style: { stroke: visited ? '#4caf50' : '#555', strokeWidth: visited ? 2.5 : 1 },
    };
  }),
);

/** COMPUTED — Selected node detail for step panel */

const selectedNodeDetail = computed(() => {
  if (!selectedNodeId.value) return null;
  const p = nodeStatusMap.value[selectedNodeId.value];
  const o = nodeOutputs.value[selectedNodeId.value];
  const n = definitionNodes.value.find(n => n.id === selectedNodeId.value);
  return {
    nodeId: selectedNodeId.value,
    nodeType: n?.type || '',
    label: n?.label || selectedNodeId.value,
    status: p?.status || 'notReached',
    durationMs: p?.durationMs || 0,
    outputHandle: p?.outputHandle || '',
    error: p?.error || '',
    outputs: o || null,
    config: (n?.config as Record<string, unknown>) || null,
  };
});

/** FUNCTIONS — Actions */

/**
 * Lazily boot core + marketplace workflow plugins
 */
async function ensurePluginsBooted(): Promise<void> {
  if (registry.nodeTypeCount === 0) {
    // Boot core plugins (start, end, condition, etc.)
    const { bootWorkflowPlugins } = await import('@src/components/workflow/constants/corePlugins.constant');
    bootWorkflowPlugins((plugin) => registry.registerPlugin(plugin));
    // Boot marketplace plugins (telegram, slack, etc.)
    await bootMarketplacePlugins((plugin) => registry.registerPlugin(plugin));
  }
  pluginsReady.value = true;
}

/**
 * Load workflow definition (for DAG) and full execution data (for path highlighting)
 */
async function loadData(): Promise<void> {
  if (!props.execution?.definitionId) return;
  isLoading.value = true;
  fullExecution.value = null;

  try {
    await ensurePluginsBooted();

    // Fetch definition + execution in parallel
    // Both services use the same executionId (MongoDB _id hex) for lookup
    const executionId = props.execution.id;
    const isHot = HOT_STATUSES.includes(props.execution.status as typeof HOT_STATUSES[number]);

    const [defRes, execRes] = await Promise.allSettled([
      // 1. Definition (for DAG nodes/edges)
      apis.workflows.definition.getById({ workflowId: props.execution.definitionId }),
      // 2. Full execution (for executionPath, nodeOutputs, state)
      isHot
        ? apis.workflows.execution.getById({ executionId })
        : apis.events.events.getWorkflowByExecutionId({ executionId }),
    ]);

    // Process definition
    if (defRes.status === 'fulfilled') {
      const def = defRes.value as Record<string, unknown>;
      definitionNodes.value = (def.nodes as WorkflowNode[]) || [];
      definitionEdges.value = (def.edges as WorkflowEdge[]) || [];
      definitionStates.value = (def.states as typeof definitionStates.value) || [];
    } else {
      logger.error('Failed to load definition', defRes.reason);
    }

    // Process execution
    if (execRes.status === 'fulfilled') {
      fullExecution.value = execRes.value;
    } else {
      logger.warn('Failed to load full execution — using list data', execRes.reason);
    }
  } catch (err) {
    logger.error('Failed to load data', err);
  } finally {
    isLoading.value = false;
    void nextTick(() => {
      void nextTick(() => fitView({ padding: 0.15, duration: 300 }));
    });
  }
}

function handleNodeClick(event: { node: { id: string } }): void {
  selectedNodeId.value = event.node.id;
  showStepPanel.value = true;
}

function handlePaneClick(): void {
  selectedNodeId.value = null;
  showStepPanel.value = false;
}

function closeStepPanel(): void {
  showStepPanel.value = false;
  selectedNodeId.value = null;
}

/** WATCHERS */
watch(() => props.execution, (exec) => {
  if (exec) {
    selectedNodeId.value = null;
    showStepPanel.value = false;
    activeTab.value = 'dag';
    void loadData();
  }
}, { immediate: true });
</script>

<template>
  <div class="exec-viewer">
    <!-- Header -->
    <div class="exec-viewer__header row items-center no-wrap q-py-sm q-px-md">
      <q-btn flat round dense icon="arrow_back" class="q-mr-sm" @click="emit('close')" />
      <q-icon :name="STATUS_ICONS[execution.status] || 'timeline'" size="24px" class="q-mr-sm" />
      <span class="text-h6 ellipsis">{{ t.modal.title.value }}</span>
      <span class="q-ml-sm text-subtitle2 text-grey-6">— {{ execution.instanceName || execution.workflowName }}</span>
      <q-badge :color="STATUS_COLORS[execution.status] || 'grey-6'" :label="execution.status.toUpperCase()" class="q-ml-sm" />
      <q-space />
      <span class="text-caption text-grey-5" style="font-family: monospace">{{ (execution.workflowUUID || execution.id).substring(0, 16) }}...</span>
      <q-btn flat round dense size="xs" icon="content_copy" color="grey-5" class="q-ml-xs" @click="copyToClipboard(execution.workflowUUID || execution.id)" />
    </div>

    <q-separator />

    <!-- Tabs -->
    <q-tabs v-model="activeTab" dense class="exec-viewer__tabs text-grey-7" active-color="primary" indicator-color="primary" align="right">
      <q-tab name="dag" icon="account_tree" :label="t.modal.tabs.dag.value" />
      <q-tab name="event" icon="input" :label="t.modal.tabs.event.value" />
      <q-tab name="state" icon="data_object" :label="t.modal.tabs.state.value" />
      <q-tab name="logs" icon="format_list_numbered" :label="t.modal.tabs.logs.value" />
    </q-tabs>

    <q-separator />

    <!-- ═══ DAG Tab ═══ -->
    <div class="exec-viewer__dag" :class="{ 'exec-viewer__dag--hidden': activeTab !== 'dag' }">
      <!-- Loading -->
      <div v-if="isLoading" class="absolute-center">
        <q-spinner color="primary" size="48px" />
      </div>

      <!-- Vue Flow Canvas (same pattern as WorkflowCanvas) -->
      <VueFlow
        v-if="!isLoading && pluginsReady && flowNodes.length > 0"
        id="execution-viewer"
        :nodes="flowNodes"
        :edges="flowEdges"
        :node-types="(nodeTypes as any)"
        :nodes-draggable="false"
        :nodes-connectable="false"
        :edges-updatable="false"
        :delete-key-code="null"
        fit-view-on-init
        class="exec-viewer__flow"
        @nodeClick="handleNodeClick"
        @paneClick="handlePaneClick"
      >
        <Background variant="dots" :gap="20" :size="1" />
        <Controls position="top-left" />
      </VueFlow>

      <!-- Empty -->
      <div v-if="!isLoading && flowNodes.length === 0" class="absolute-center text-grey-5">
        {{ t.modal.noData.value }}
      </div>

      <!-- Legend (LEFT side) -->
      <div v-if="!isLoading && flowNodes.length > 0" class="exec-viewer__legend">
        <div v-for="item in [
          { color: '#4caf50', label: t.modal.legend.completed.value },
          { color: '#ff9800', label: t.modal.legend.waiting.value },
          { color: '#ffc107', label: t.modal.legend.retrying.value },
          { color: '#fdd835', label: t.modal.legend.timeout.value },
          { color: '#f44336', label: t.modal.legend.error.value },
          { color: '#616161', label: t.modal.legend.cancelled.value },
        ]" :key="item.color" class="row items-center q-gutter-xs q-mb-xs">
          <div class="exec-viewer__legend-dot" :style="{ background: item.color }" />
          <span class="text-caption">{{ item.label }}</span>
        </div>
        <div class="row items-center q-gutter-xs">
          <div class="exec-viewer__legend-dot" style="background: #9e9e9e; opacity: 0.4" />
          <span class="text-caption">{{ t.modal.legend.notReached.value }}</span>
        </div>
      </div>

      <!-- Step Detail Panel (RIGHT side overlay) -->
      <transition name="slide-right">
        <div v-if="showStepPanel && selectedNodeDetail" class="exec-viewer__step-panel">
          <!-- Panel header -->
          <div class="exec-viewer__step-header row items-center no-wrap q-py-sm q-px-md">
            <q-icon :name="getStepIcon(selectedNodeDetail.status)" :color="getStepColor(selectedNodeDetail.status)" size="24px" class="q-mr-sm" />
            <div class="col">
              <div class="text-subtitle1 text-weight-bold">{{ selectedNodeDetail.label }}</div>
              <div class="text-caption text-grey-6">{{ selectedNodeDetail.nodeType }}</div>
            </div>
            <q-btn flat round dense icon="close" size="sm" @click="closeStepPanel" />
          </div>

          <q-separator />

          <!-- Panel content -->
          <div class="exec-viewer__step-content">
            <!-- Status badge -->
            <div class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.status.value }}</div>
              <q-badge :color="getStepColor(selectedNodeDetail.status)" :label="formatStatusLabel(selectedNodeDetail.status)" />
            </div>

            <!-- Node Type -->
            <div v-if="selectedNodeDetail.nodeType" class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.nodeType.value }}</div>
              <q-badge outline color="blue-grey-6" :label="selectedNodeDetail.nodeType" />
            </div>

            <!-- Node ID -->
            <div class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.nodeId.value }}</div>
              <div style="font-family: monospace; font-size: 12px; color: var(--mapex-text-secondary); word-break: break-all">{{ selectedNodeDetail.nodeId }}</div>
            </div>

            <!-- Duration -->
            <div v-if="selectedNodeDetail.durationMs > 0" class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.duration.value }}</div>
              <div>{{ formatDuration(selectedNodeDetail.durationMs) }}</div>
            </div>

            <!-- Output handle -->
            <div v-if="selectedNodeDetail.outputHandle" class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.outputHandle.value }}</div>
              <q-badge color="blue-grey-6" :label="selectedNodeDetail.outputHandle" />
            </div>

            <!-- Configuration (what the user configured on this node) -->
            <div class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.configuration.value }}</div>
              <pre
                v-if="selectedNodeDetail.config && Object.keys(selectedNodeDetail.config).length > 0"
                class="exec-viewer__code"
              >{{ formatJson(selectedNodeDetail.config) }}</pre>
              <div v-else class="text-caption text-grey-5" style="font-style: italic">
                {{ t.modal.drawer.configurationEmpty.value }}
              </div>
            </div>

            <!-- Error -->
            <div v-if="selectedNodeDetail.error" class="q-mb-md">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.error.value }}</div>
              <q-banner dense rounded class="bg-red-1 text-negative text-caption">
                {{ selectedNodeDetail.error }}
              </q-banner>
            </div>

            <!-- Outputs -->
            <div v-if="selectedNodeDetail.outputs">
              <div class="text-caption text-grey-6 q-mb-xs">{{ t.modal.drawer.outputs.value }}</div>
              <pre class="exec-viewer__code">{{ formatJson(selectedNodeDetail.outputs) }}</pre>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <!-- ═══ Event Tab ═══ -->
    <div v-if="activeTab === 'event'" class="exec-viewer__scroll q-pa-md">
      <!-- Error banner -->
      <q-banner v-if="errorInfo" rounded class="bg-red-1 text-negative q-mb-md">
        <template #avatar><q-icon name="error" color="negative" /></template>
        <div class="text-weight-medium">{{ errorInfo.code }}: {{ errorInfo.message }}</div>
      </q-banner>

      <!-- Event Payload -->
      <q-card v-if="eventPayload" flat bordered class="q-mb-md">
        <q-card-section>
          <div class="row items-center q-mb-sm">
            <q-icon name="input" size="20px" color="blue-grey-6" class="q-mr-sm" />
            <div class="text-subtitle2">{{ t.modal.states.eventPayload.value }}</div>
          </div>
          <pre class="exec-viewer__code">{{ formatJson(eventPayload) }}</pre>
        </q-card-section>
      </q-card>

      <!-- Node Outputs -->
      <q-card v-if="Object.keys(nodeOutputs).length > 0" flat bordered class="q-mb-md">
        <q-card-section>
          <div class="row items-center q-mb-sm">
            <q-icon name="output" size="20px" color="teal-6" class="q-mr-sm" />
            <div class="text-subtitle2">{{ t.modal.states.nodeOutputs.value }}</div>
          </div>
          <div v-for="(output, nid) in nodeOutputs" :key="String(nid)" class="q-mb-sm">
            <div class="text-caption text-weight-medium text-grey-7 q-mb-xs">{{ nid }}</div>
            <pre class="exec-viewer__code">{{ formatJson(output) }}</pre>
          </div>
        </q-card-section>
      </q-card>

      <!-- No data -->
      <q-card v-if="!errorInfo && !eventPayload && Object.keys(nodeOutputs).length === 0" flat bordered>
        <q-card-section class="text-center q-pa-xl">
          <q-icon name="inbox" size="48px" color="grey-5" class="q-mb-sm" />
          <div class="text-grey-5">{{ t.modal.noData.value }}</div>
        </q-card-section>
      </q-card>
    </div>

    <!-- ═══ State Tab ═══ -->
    <div v-if="activeTab === 'state'" class="exec-viewer__scroll q-pa-md">
      <!-- State Variables Table -->
      <q-card v-if="stateVariables.length > 0" flat bordered>
        <q-card-section class="q-pb-none">
          <div class="row items-center">
            <q-icon name="data_object" size="20px" color="primary" class="q-mr-sm" />
            <div class="text-subtitle2">{{ t.modal.states.workflowState.value }}</div>
            <q-space />
            <q-badge outline color="grey-6" :label="`${stateVariables.length} var${stateVariables.length > 1 ? 's' : ''}`" />
          </div>
        </q-card-section>

        <q-card-section>
          <q-markup-table flat bordered separator="horizontal" dense class="exec-viewer__state-table">
            <thead>
              <tr>
                <th class="text-left">{{ t.modal.states.colName.value }}</th>
                <th class="text-left">{{ t.modal.states.colType.value }}</th>
                <th class="text-left">{{ t.modal.states.colDefault.value }}</th>
                <th class="text-left">{{ t.modal.states.colCurrent.value }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="v in stateVariables" :key="v.field">
                <td>
                  <div class="row items-center no-wrap">
                    <span class="text-weight-medium" style="font-family: monospace; font-size: 12px">{{ v.field }}</span>
                    <q-icon v-if="v.durable" name="save" size="14px" color="blue-grey-5" class="q-ml-xs">
                      <q-tooltip>{{ t.modal.states.durable.value }}</q-tooltip>
                    </q-icon>
                  </div>
                  <div v-if="v.description" class="text-caption text-grey-6" style="font-size: 10px; line-height: 1.2">{{ v.description }}</div>
                </td>
                <td>
                  <q-badge outline :color="v.type === 'number' ? 'blue-6' : v.type === 'boolean' ? 'orange-7' : v.type === 'json' ? 'purple-6' : 'grey-7'" :label="v.type" />
                </td>
                <td style="font-family: monospace; font-size: 11px; color: var(--mapex-text-muted)">
                  {{ v.defaultValue !== undefined && v.defaultValue !== '' ? String(v.defaultValue) : '—' }}
                </td>
                <td>
                  <span
                    style="font-family: monospace; font-size: 12px"
                    :style="{ color: v.changed ? 'var(--q-primary)' : 'var(--mapex-text-secondary)', fontWeight: v.changed ? '600' : '400' }"
                  >
                    {{ typeof v.currentValue === 'object' ? formatJson(v.currentValue) : String(v.currentValue ?? '—') }}
                  </span>
                </td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-card-section>
      </q-card>

      <!-- No data -->
      <q-card v-if="stateVariables.length === 0" flat bordered>
        <q-card-section class="text-center q-pa-xl">
          <q-icon name="inbox" size="48px" color="grey-5" class="q-mb-sm" />
          <div class="text-grey-5">{{ t.modal.noData.value }}</div>
        </q-card-section>
      </q-card>
    </div>

    <!-- ═══ Logs Tab ═══ -->
    <div v-if="activeTab === 'logs'" class="exec-viewer__scroll">
      <q-list separator>
        <q-item v-for="(entry, idx) in executionPath" :key="idx" clickable @click="selectedNodeId = entry.nodeId; showStepPanel = true; activeTab = 'dag'">
          <q-item-section avatar>
            <q-icon :name="getStepIcon(entry.status)" :color="getStepColor(entry.status)" size="22px" />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-medium">{{ definitionNodes.find(n => n.id === entry.nodeId)?.label || entry.nodeId }}</q-item-label>
            <q-item-label caption>{{ entry.nodeType }} — {{ entry.status }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-item-label caption>{{ formatDuration(entry.durationMs || 0) }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
      <div v-if="executionPath.length === 0" class="text-center text-grey-5 q-pa-xl">{{ t.modal.noData.value }}</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.exec-viewer {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 50px);
  min-height: 400px;

  &__header {
    flex-shrink: 0;
    background: var(--mapex-surface-primary);
  }

  &__tabs {
    flex-shrink: 0;
    background: var(--mapex-surface-primary);
  }

  /* ── DAG: same as WorkflowCanvas container ── */
  &__dag {
    flex: 1;
    position: relative;
    overflow: hidden;

    &--hidden {
      display: none;
    }
  }

  &__flow {
    width: 100%;
    height: 100%;
    background: var(--mapex-page-bg);
  }

  /* ── Scrollable tabs ── */
  &__scroll {
    flex: 1;
    overflow-y: auto;
  }

  /* ── Legend (LEFT side, below zoom controls) ── */
  &__legend {
    position: absolute;
    bottom: 16px;
    left: 12px;
    z-index: 10;
    padding: 8px 12px;
    background: var(--mapex-surface-primary);
    border: 1px solid var(--mapex-border-primary);
    border-radius: var(--mapex-radius-md, 8px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  }

  &__legend-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  /* ── Step Detail Panel (right overlay — matches NodeConfigPanel pattern) ── */
  &__step-panel {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: 380px;
    z-index: 20;
    display: flex;
    flex-direction: column;
    background: var(--mapex-surface-bg);
    border-left: 1px solid var(--mapex-card-border);
    box-shadow: -4px 0 16px rgba(0, 0, 0, 0.2);
  }

  &__step-header {
    flex-shrink: 0;
    background: var(--mapex-surface-bg);
  }

  &__step-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  /* ── State variables table ── */
  &__state-table {
    background: transparent !important;

    th {
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.3px;
      color: var(--mapex-text-muted) !important;
      padding: 8px 12px !important;
    }

    td {
      padding: 10px 12px !important;
      vertical-align: top;
    }
  }

  /* ── Code block ── */
  &__code {
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 11px;
    white-space: pre-wrap;
    word-break: break-all;
    background: var(--mapex-surface-secondary);
    border: 1px solid var(--mapex-border-primary);
    border-radius: 6px;
    padding: 8px 12px;
    max-height: 300px;
    overflow: auto;
    margin: 0;
  }
}

/* ── Execution status: "turn on/off" the node's internal light ──
   Palette aligned with the 6 statuses emitted by the workflow runtime:
     completed / waiting / retrying / timeout / error / cancelled
   Plus the frontend-computed "notReached" (exec-node-off). */
.exec-viewer__dag {

  /* Not reached — greyscale + dimmed (node exists in DAG but was never visited) */
  :deep(.vue-flow__node.exec-node-off .wf-node__icon) {
    filter: grayscale(0.8) brightness(0.6);
    opacity: 0.55;
    box-shadow: none;
  }
  :deep(.vue-flow__node.exec-node-off .wf-node__label) {
    opacity: 0.5;
  }

  /* Completed — green glow */
  :deep(.vue-flow__node.exec-node-completed .wf-node__icon) {
    box-shadow:
      0 0 28px rgba(76, 175, 80, 0.45),
      0 0 56px rgba(76, 175, 80, 0.22);
  }

  /* Waiting — orange pulse (suspended, callback/signal/timer expected) */
  :deep(.vue-flow__node.exec-node-waiting .wf-node__icon) {
    animation: exec-pulse-waiting 2s ease-in-out infinite;
  }

  /* Retrying — yellow pulse (retry timer fired, re-running) */
  :deep(.vue-flow__node.exec-node-retrying .wf-node__icon) {
    animation: exec-pulse-retrying 1.4s ease-in-out infinite;
  }

  /* Timeout — yellow glow (gave up waiting) */
  :deep(.vue-flow__node.exec-node-timeout .wf-node__icon) {
    box-shadow:
      0 0 28px rgba(253, 216, 53, 0.55),
      0 0 56px rgba(253, 216, 53, 0.28);
  }

  /* Error — red glow */
  :deep(.vue-flow__node.exec-node-error .wf-node__icon) {
    box-shadow:
      0 0 28px rgba(244, 67, 54, 0.45),
      0 0 56px rgba(244, 67, 54, 0.22);
  }

  /* Cancelled — attenuated dark grey (killed by system, not a failure) */
  :deep(.vue-flow__node.exec-node-cancelled .wf-node__icon) {
    filter: grayscale(0.6) brightness(0.75);
    opacity: 0.75;
    box-shadow:
      0 0 20px rgba(97, 97, 97, 0.35),
      0 0 40px rgba(97, 97, 97, 0.18);
  }
  :deep(.vue-flow__node.exec-node-cancelled .wf-node__label) {
    opacity: 0.65;
  }

  /* Selected node — brighter + scale */
  :deep(.vue-flow__node.exec-node-selected .wf-node__icon) {
    transform: scale(1.08);
  }

  /* Hide resize handles on group frames — read-only viewer */
  :deep(.vue-flow__resize-control) {
    display: none !important;
  }

}

/* Waiting pulse — orange (#ff9800) */
@keyframes exec-pulse-waiting {
  0%, 100% {
    box-shadow:
      0 0 24px rgba(255, 152, 0, 0.35),
      0 0 48px rgba(255, 152, 0, 0.18);
  }
  50% {
    box-shadow:
      0 0 32px rgba(255, 152, 0, 0.55),
      0 0 60px rgba(255, 152, 0, 0.30);
  }
}

/* Retrying pulse — yellow-amber (#ffc107), faster cadence to suggest activity */
@keyframes exec-pulse-retrying {
  0%, 100% {
    box-shadow:
      0 0 22px rgba(255, 193, 7, 0.40),
      0 0 44px rgba(255, 193, 7, 0.20);
  }
  50% {
    box-shadow:
      0 0 34px rgba(255, 193, 7, 0.65),
      0 0 64px rgba(255, 193, 7, 0.35);
  }
}

/* ── Transitions ── */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.2s ease;
}
.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(100%);
}
</style>
