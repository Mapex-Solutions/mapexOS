<script setup lang="ts">
/** TYPE IMPORTS */
import type { WorkflowCanvasProps, WorkflowCanvasEmits } from './interfaces/WorkflowCanvas.interface';
import type { WorkflowNode, WorkflowEdge } from '../../interfaces/CreateEditWorkflow.interface';

/** VUE IMPORTS */
import { computed, nextTick, onMounted, onUnmounted, ref, watch, markRaw } from 'vue';
import { VueFlow, useVueFlow, MarkerType } from '@vue-flow/core';
import { MiniMap } from '@vue-flow/minimap';
import { Controls } from '@vue-flow/controls';
import { Background } from '@vue-flow/background';
import { useQuasar } from 'quasar';

/** COMPONENTS */
import AdjustableEdge from '../AdjustableEdge/AdjustableEdge.vue';
import { GenericWorkflowNode } from '@components/workflow/nodes/GenericWorkflowNode';

/** COMPOSABLES */
import { useWorkflowEditorState, useWorkflowHistory } from '../../composables';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** UTILS */
import { createConnectionValidator, resolveNodeHandles } from '../../utils';
import { buildDefaultConfig } from '@utils/workflow';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { SNAP_GRID_SIZE } from '../../constants';
import { MIN_ZOOM, MAX_ZOOM } from './constants';
import { createWorkflowHotkeyHandler } from './utils';

/** PROPS & EMITS */
defineProps<WorkflowCanvasProps>();
const emit = defineEmits<WorkflowCanvasEmits>();

/** COMPOSABLES & STORES */
const $q = useQuasar();
const pluginRegistry = usePluginRegistryStore();
const { nodes, edges, addNode, addEdge, removeEdge, removeNode, duplicateNode, nodeConfigVersion, nodeValidationErrors, updateViewport } = useWorkflowEditorState();
const { pushSnapshot, undo, redo, finishRestore, isBatchRestoring } = useWorkflowHistory();

const {
  onConnect,
  addEdges,
  removeEdges: vfRemoveEdges,
  onNodeClick,
  onPaneClick,
  onEdgesChange,
  onNodeDragStart,
  onNodeDragStop,
  screenToFlowCoordinate,
  getNodes,
  getSelectedNodes,
  getSelectedEdges,
  setNodes,
  setEdges,
  addNodes: vfAddNodes,
  removeNodes: vfRemoveNodes,
  updateNodeInternals,
  onViewportChangeEnd,
  fitView,
} = useVueFlow();

const isValidConnection = createConnectionValidator(() => ({
  getNodeType: (typeId: string) => pluginRegistry.getNodeType(typeId),
  nodes: nodes.value,
  edges: edges.value,
  resolveHandles: resolveNodeHandles,
}));

/** STATE */

/**
 * Counter for generating unique node IDs
 */
const nodeIdCounter = ref(0);

/**
 * Clipboard buffer for copy/paste — stores copied nodes and internal edges
 */
const clipboard = ref<{ nodes: WorkflowNode[]; edges: WorkflowEdge[] } | null>(null);

/** COMPUTED */

/**
 * Grid dot color — dark dots on light bg, subtle white dots on dark bg
 */
const gridPatternColor = computed(() =>
  $q.dark.isActive ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.08)',
);

/**
 * Edge accent color — uses MapexOS primary palette
 */
const edgeColor = computed(() =>
  $q.dark.isActive ? '#4CAF7D' : '#3b6d5e',
);

/**
 * Build Vue Flow nodeTypes map from plugin registry.
 * Accesses nodeTypeMap directly (not via action) for proper reactivity tracking.
 */
const nodeTypes = computed(() => {
  const types: Record<string, any> = {};
  for (const [type, nodeType] of pluginRegistry.nodeTypeMap.entries()) {
    types[type] = nodeType.canvasComponent ?? GenericWorkflowNode;
  }
  return markRaw(types);
});

/**
 * Custom edge type map for Vue Flow.
 * Registers 'adjustable' edge type using our AdjustableEdge component.
 */
const edgeTypes = computed(() => {
  const types: Record<string, any> = {
    adjustable: AdjustableEdge,
  };
  return markRaw(types);
});

/**
 * Map workflow nodes to Vue Flow nodes.
 * Map workflow nodes to Vue Flow nodes.
 * Group frame children use Vue Flow parentNode (relative positioning).
 * Text notes use absolute positions — parentNode breaks edge routing.
 */
const flowNodes = computed(() =>
  nodes.value.map(node => {
    // Text notes use absolute positions — Vue Flow parentNode breaks edge routing
    // between child and parent. Only group_frame children use Vue Flow parentNode.
    const useVfParent = node.parentNodeId && node.type !== 'core/text_note';

    return {
      id: node.id,
      type: node.type,
      position: node.position,
      data: { config: node.config, label: node.label ?? '', __nodeType: node.type, hasErrors: !!(nodeValidationErrors.value[node.id]?.length) },
      ...(useVfParent ? { parentNode: node.parentNodeId } : {}),
      ...(node.type === 'core/group_frame' ? { zIndex: -1 } : {}),
    };
  }),
);

/**
 * Map workflow edges to Vue Flow edges with animated bezier style.
 */
const flowEdges = computed(() =>
  edges.value.map(edge => ({
    id: edge.id,
    source: edge.source,
    ...(edge.sourceHandle !== undefined && { sourceHandle: edge.sourceHandle }),
    target: edge.target,
    ...(edge.targetHandle !== undefined && { targetHandle: edge.targetHandle }),
    ...(edge.label !== undefined && { label: edge.label }),
    type: 'adjustable',
    animated: true,
    style: { stroke: edgeColor.value, strokeWidth: 2.5 },
    markerEnd: { type: MarkerType.ArrowClosed, color: edgeColor.value },
  })),
);

/** WATCHERS */

/**
 * Sync composable nodes → Vue Flow internal state.
 * Watches for node additions/removals (ID list changes) and pushes
 * updated list to Vue Flow, preserving any user-dragged positions.
 */
const nodeIdList = computed(() => nodes.value.map(n => n.id).join(','));

watch(nodeIdList, () => {
  if (isBatchRestoring.value) return;

  const positionMap = new Map(getNodes.value.map(n => [n.id, n.position]));

  setNodes(nodes.value.map(node => {
    const useVfParent = node.parentNodeId && node.type !== 'core/text_note';
    return {
      id: node.id,
      type: node.type,
      position: useVfParent ? node.position : (positionMap.get(node.id) || node.position),
      data: { config: node.config, label: node.label ?? '', __nodeType: node.type, hasErrors: !!(nodeValidationErrors.value[node.id]?.length) },
      ...(useVfParent ? { parentNode: node.parentNodeId } : {}),
      ...(node.type === 'core/group_frame' ? { zIndex: -1 } : {}),
    };
  }));

  void nextTick(() => {
    updateNodeInternals(nodes.value.map(n => n.id));
  });
});

/**
 * Sync composable edges → Vue Flow internal state.
 */
const edgeIdList = computed(() => edges.value.map(e => e.id).join(','));

watch(edgeIdList, () => {
  if (isBatchRestoring.value) return;

  void nextTick(() => {
    setEdges(edges.value.map(edge => ({
      id: edge.id,
      source: edge.source,
      ...(edge.sourceHandle !== undefined && { sourceHandle: edge.sourceHandle }),
      target: edge.target,
      ...(edge.targetHandle !== undefined && { targetHandle: edge.targetHandle }),
      ...(edge.label !== undefined && { label: edge.label }),
      type: 'adjustable',
      animated: true,
      style: { stroke: edgeColor.value, strokeWidth: 2.5 },
      markerEnd: { type: MarkerType.ArrowClosed, color: edgeColor.value },
    })));
  });
});

/**
 * Sync node config changes (e.g., dynamic handle updates) to Vue Flow.
 * Triggered when nodeConfigVersion increments after updateNodeConfig.
 * Preserves Vue Flow's internal dragged positions via positionMap.
 * Calls updateNodeInternals after setNodes so Vue Flow recalculates
 * handle positions and re-routes edges accordingly.
 */
watch(nodeConfigVersion, () => {
  if (isBatchRestoring.value) return;

  const positionMap = new Map(getNodes.value.map(n => [n.id, n.position]));

  setNodes(nodes.value.map(node => {
    const useVfParent = node.parentNodeId && node.type !== 'core/text_note';
    return {
      id: node.id,
      type: node.type,
      position: useVfParent ? node.position : (positionMap.get(node.id) || node.position),
      data: { config: node.config, label: node.label ?? '', __nodeType: node.type, hasErrors: !!(nodeValidationErrors.value[node.id]?.length) },
      ...(useVfParent ? { parentNode: node.parentNodeId } : {}),
      ...(node.type === 'core/group_frame' ? { zIndex: -1 } : {}),
    };
  }));

  void nextTick(() => {
    updateNodeInternals(nodes.value.map(n => n.id));
  });
});

/** FUNCTIONS */

/**
 * Handle new connection between nodes.
 * Syncs to both Vue Flow and our composable state.
 */
onConnect((connection) => {
  pushSnapshot('Connect');

  const edgeId = `e_${crypto.randomUUID()}`;

  addEdges([{ ...connection, id: edgeId }]);

  addEdge({
    id: edgeId,
    source: connection.source,
    ...(connection.sourceHandle ? { sourceHandle: connection.sourceHandle } : {}),
    target: connection.target,
    ...(connection.targetHandle ? { targetHandle: connection.targetHandle } : {}),
  });
});

/**
 * Sync Vue Flow edge removals (select + Delete key) to our composable state
 */
onEdgesChange((changes) => {
  if (isBatchRestoring.value) return;

  const removals = changes.filter(c => c.type === 'remove');
  if (removals.length) pushSnapshot('Remove edge');

  for (const change of changes) {
    if (change.type === 'remove') {
      removeEdge(change.id);
    }
  }
});

/**
 * Capture node positions before drag begins (for undo)
 */
onNodeDragStart(() => {
  pushSnapshot('Move node');
});

/**
 * Sync final dragged positions back to composable state after drag ends.
 * For text notes: they use absolute positions (not Vue Flow parentNode),
 * so when a parent node moves we shift its child notes by the same delta.
 */
onNodeDragStop(({ nodes: draggedNodes }) => {
  for (const vfNode of draggedNodes) {
    const composableNode = nodes.value.find(n => n.id === vfNode.id);
    if (!composableNode) continue;

    const dx = vfNode.position.x - composableNode.position.x;
    const dy = vfNode.position.y - composableNode.position.y;
    composableNode.position = { ...vfNode.position };

    // Shift child text notes by the same delta (they use absolute positions)
    if (dx !== 0 || dy !== 0) {
      for (const child of nodes.value) {
        if (child.parentNodeId === composableNode.id && child.type === 'core/text_note') {
          child.position = { x: child.position.x + dx, y: child.position.y + dy };
        }
      }
    }
  }
});

/**
 * Handle node click — emit selection only if the node type is configurable
 * and CTRL/Meta is NOT pressed (multi-select mode should not open config panel)
 */
onNodeClick(({ node, event }) => {
  if (event.ctrlKey || event.metaKey) return;
  const nodeType = pluginRegistry.getNodeType(node.data?.__nodeType as string);
  if (nodeType?.configurable === false) return;
  emit('node-select', node.id);
});

/**
 * Handle pane click — deselect
 */
onPaneClick(() => {
  emit('canvas-click');
});

/**
 * Sync Vue Flow viewport to composable state so other components
 * (e.g., PluginCatalogGroup) can position nodes in the visible area.
 */
onViewportChangeEnd(({ x, y, zoom }) => {
  updateViewport({ x, y, zoom });
});

/**
 * Handle drop from catalog — create new node at drop position
 *
 * @param {DragEvent} event - HTML5 drop event
 * @returns {void}
 */
function handleDrop(event: DragEvent): void {
  if (!event.dataTransfer) return;

  const nodeTypeStr = event.dataTransfer.getData('application/workflow-node-type');
  if (!nodeTypeStr) return;

  const nodeType = pluginRegistry.getNodeType(nodeTypeStr);
  if (!nodeType) return;

  const position = screenToFlowCoordinate({
    x: event.clientX,
    y: event.clientY,
  });

  const newNodeId = `n_${nodeTypeStr.replace('/', '_')}_${++nodeIdCounter.value}`;

  // Detect default operation from properties to build a filtered config
  const operationProp = nodeType.properties?.find(p => p.name === 'operation');
  const defaultOperation = operationProp?.default as string | undefined;
  const config = buildDefaultConfig(nodeType, defaultOperation);

  pushSnapshot('Add node');
  addNode({
    id: newNodeId,
    type: nodeTypeStr,
    position,
    config,
    label: nodeType.label,
  });
}

/**
 * Prevent default on dragover to allow drop
 *
 * @param {DragEvent} event - Drag event
 * @returns {void}
 */
function handleDragOver(event: DragEvent): void {
  event.preventDefault();
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move';
  }
}

/** HOTKEYS */

const hotkeys = createWorkflowHotkeyHandler({
  onDelete: () => {
    const selectedNodes = getSelectedNodes.value;
    const selectedEdges = getSelectedEdges.value;

    if (!selectedNodes.length && !selectedEdges.length) return;

    pushSnapshot('Delete');

    for (const node of selectedNodes) {
      const removedIds = removeNode(node.id);
      if (removedIds.length) vfRemoveNodes(removedIds);
    }

    for (const edge of selectedEdges) {
      removeEdge(edge.id);
      vfRemoveEdges([edge.id]);
    }

    emit('canvas-click');
  },
  onDuplicate: () => {
    const selected = getSelectedNodes.value;
    if (!selected.length) return;

    pushSnapshot('Duplicate');

    for (const node of selected) {
      const newId = duplicateNode(node.id);
      if (!newId) continue;

      const newNode = nodes.value.find(n => n.id === newId);
      if (!newNode) continue;

      vfAddNodes([{
        id: newNode.id,
        type: newNode.type,
        position: newNode.position,
        data: { config: newNode.config, label: newNode.label ?? '', __nodeType: newNode.type },
      }]);
    }
  },
  onCopy: () => {
    const selected = getSelectedNodes.value;
    if (!selected.length) return;

    const selectedIds = new Set(selected.map(n => n.id));

    // Deep-clone selected nodes
    const copiedNodes: WorkflowNode[] = selected
      .filter(n => {
        const nt = pluginRegistry.getNodeType(n.data?.__nodeType as string);
        return nt?.deletable !== false;
      })
      .map(n => ({
        id: n.id,
        type: n.type,
        position: { ...n.position },
        config: { ...(n.data?.config as Record<string, unknown> || {}) },
        label: (n.data?.label as string) || '',
      }));

    // Capture edges where both source and target are in the selection
    const copiedEdges: WorkflowEdge[] = edges.value
      .filter(e => selectedIds.has(e.source) && selectedIds.has(e.target))
      .map(e => ({ ...e }));

    clipboard.value = { nodes: copiedNodes, edges: copiedEdges };
  },
  onPaste: () => {
    if (!clipboard.value || !clipboard.value.nodes.length) return;

    pushSnapshot('Paste');

    const idMap = new Map<string, string>();
    const offset = 60;

    // Create new nodes with remapped IDs
    for (const srcNode of clipboard.value.nodes) {
      const newId = `n_${srcNode.type.replace('/', '_')}_${Date.now()}_${++nodeIdCounter.value}`;
      idMap.set(srcNode.id, newId);

      const nodeType = pluginRegistry.getNodeType(srcNode.type);

      addNode({
        id: newId,
        type: srcNode.type,
        position: { x: srcNode.position.x + offset, y: srcNode.position.y + offset },
        config: { ...srcNode.config },
        label: srcNode.label || nodeType?.label || '',
      });
    }

    // Re-create internal edges with remapped IDs
    for (const srcEdge of clipboard.value.edges) {
      const newSource = idMap.get(srcEdge.source);
      const newTarget = idMap.get(srcEdge.target);
      if (!newSource || !newTarget) continue;

      const newEdgeId = `e_${crypto.randomUUID()}`;
      addEdge({
        id: newEdgeId,
        source: newSource,
        ...(srcEdge.sourceHandle ? { sourceHandle: srcEdge.sourceHandle } : {}),
        target: newTarget,
        ...(srcEdge.targetHandle ? { targetHandle: srcEdge.targetHandle } : {}),
      });
    }

    // Update clipboard positions so next paste offsets further
    clipboard.value = {
      nodes: clipboard.value.nodes.map(n => ({
        ...n,
        position: { x: n.position.x + offset, y: n.position.y + offset },
      })),
      edges: clipboard.value.edges,
    };
  },
  onUndo: () => {
    const snapshot = undo();
    if (!snapshot) return;
    syncSnapshotToVueFlow(snapshot.nodes, snapshot.edges);
    void nextTick(() => finishRestore());
    emit('canvas-click');
  },
  onRedo: () => {
    const snapshot = redo();
    if (!snapshot) return;
    syncSnapshotToVueFlow(snapshot.nodes, snapshot.edges);
    void nextTick(() => finishRestore());
    emit('canvas-click');
  },
  onEscape: () => {
    emit('canvas-click');
  },
});

/**
 * Sync a snapshot's nodes and edges directly to Vue Flow internal state.
 * Bypasses the normal watcher-based sync to avoid positionMap overrides.
 *
 * @param {WorkflowNode[]} snapshotNodes - Nodes to restore
 * @param {WorkflowEdge[]} snapshotEdges - Edges to restore
 * @returns {void}
 */
function syncSnapshotToVueFlow(snapshotNodes: WorkflowNode[], snapshotEdges: WorkflowEdge[]): void {
  setNodes(snapshotNodes.map(node => {
    const useVfParent = node.parentNodeId && node.type !== 'core/text_note';
    return {
      id: node.id,
      type: node.type,
      position: { ...node.position },
      data: { config: node.config, label: node.label ?? '', __nodeType: node.type, hasErrors: !!(nodeValidationErrors.value[node.id]?.length) },
      ...(useVfParent ? { parentNode: node.parentNodeId } : {}),
      ...(node.type === 'core/group_frame' ? { zIndex: -1 } : {}),
    };
  }));

  setEdges(snapshotEdges.map(edge => ({
    id: edge.id,
    source: edge.source,
    ...(edge.sourceHandle !== undefined && { sourceHandle: edge.sourceHandle }),
    target: edge.target,
    ...(edge.targetHandle !== undefined && { targetHandle: edge.targetHandle }),
    ...(edge.label !== undefined && { label: edge.label }),
    type: 'adjustable',
    animated: true,
    style: { stroke: edgeColor.value, strokeWidth: 2.5 },
    markerEnd: { type: MarkerType.ArrowClosed, color: edgeColor.value },
  })));

  void nextTick(() => {
    updateNodeInternals(snapshotNodes.map(n => n.id));
  });
}

/** EXPOSE */
defineExpose({
  /**
   * Fit all nodes into the viewport with smooth animation
   */
  fitView: () => fitView({ padding: 0.15, duration: 300 }),

  /**
   * Apply layout positions directly to Vue Flow nodes.
   * Bypasses watchers — directly sets each node's position in Vue Flow
   * and syncs composable state, then recalculates internals.
   *
   * @param {WorkflowNode[]} layoutNodes - Nodes with updated positions
   */
  applyLayout: (layoutNodes: WorkflowNode[]) => {
    // Update composable state
    nodes.value = layoutNodes;

    // Directly set positions on Vue Flow internal nodes
    const posMap = new Map(layoutNodes.map(n => [n.id, n.position]));
    for (const vfNode of getNodes.value) {
      const pos = posMap.get(vfNode.id);
      if (pos) {
        vfNode.position = { ...pos };
      }
    }

    void nextTick(() => {
      updateNodeInternals(layoutNodes.map(n => n.id));
    });
  },

  /**
   * Restore a history snapshot to Vue Flow internal state.
   * Called by WorkflowTab after undo/redo updates composable state.
   *
   * @param {WorkflowNode[]} snapshotNodes - Nodes from the restored snapshot
   * @param {WorkflowEdge[]} snapshotEdges - Edges from the restored snapshot
   */
  restoreSnapshot: (snapshotNodes: WorkflowNode[], snapshotEdges: WorkflowEdge[]) => {
    syncSnapshotToVueFlow(snapshotNodes, snapshotEdges);
  },
});

/** LIFECYCLE HOOKS */

onMounted(() => hotkeys.attach());
onUnmounted(() => hotkeys.detach());
</script>

<template>
  <div
    class="workflow-canvas"
    @drop="handleDrop"
    @dragover="handleDragOver"
  >
    <VueFlow
      :nodes="flowNodes"
      :edges="flowEdges"
      :node-types="nodeTypes"
      :edge-types="edgeTypes"
      :min-zoom="MIN_ZOOM"
      :max-zoom="MAX_ZOOM"
      :snap-to-grid="true"
      :snap-grid="SNAP_GRID_SIZE"
      :nodes-draggable="!toolbarState.locked"
      :nodes-connectable="!toolbarState.locked"
      :edges-updatable="!toolbarState.locked"
      :delete-key-code="null"
      :is-valid-connection="isValidConnection"
      fit-view-on-init
      class="workflow-vue-flow"
    >
      <!-- Background grid (N8N-style dots) -->
      <Background
        v-if="toolbarState.showGrid"
        variant="dots"
        :gap="20"
        :size="2"
        :pattern-color="gridPatternColor"
      />

      <!-- Minimap -->
      <MiniMap
        v-if="toolbarState.showMinimap"
        position="bottom-right"
      />

      <!-- Controls -->
      <Controls position="bottom-left" />
    </VueFlow>
  </div>
</template>

<style lang="scss" scoped>
.workflow-canvas {
  width: 100%;
  height: 100%;
}

.workflow-vue-flow {
  width: 100%;
  height: 100%;
  background: var(--mapex-page-bg);
}
</style>

<style lang="scss">
/* Vue Flow global overrides for theme */
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/minimap/dist/style.css';
@import '@vue-flow/controls/dist/style.css';
@import '@vue-flow/node-resizer/dist/style.css';

/* Workflow SDK design tokens + utility classes */
@import 'src/css/workflow-tokens.scss';
@import 'src/css/workflow-utilities.scss';

/* ── Workflow-scoped design tokens ── */

body.body--light .workflow-vue-flow {
  --wf-edge-color: var(--mapex-primary);
  --wf-edge-hover: #2d5548;
  --wf-edge-glow: rgba(var(--mapex-primary-rgb), 0.3);
  --wf-handle-border: var(--mapex-surface-bg);
  --wf-node-border: var(--mapex-card-border);
  --wf-node-border-hover: rgba(0, 0, 0, 0.2);
  --wf-node-border-selected: rgba(0, 0, 0, 0.3);
}

body.body--dark .workflow-vue-flow {
  --wf-edge-color: var(--mapex-primary);
  --wf-edge-hover: #5BC08D;
  --wf-edge-glow: rgba(var(--mapex-primary-rgb), 0.4);
  --wf-handle-border: var(--mapex-text-on-primary, #fff);
  --wf-node-border: rgba(255, 255, 255, 0.15);
  --wf-node-border-hover: rgba(255, 255, 255, 0.3);
  --wf-node-border-selected: rgba(255, 255, 255, 0.4);
}

/* ── Minimap ── */

.vue-flow__minimap {
  background: var(--mapex-surface-bg);
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-sm);
}

/* ── Controls ── */

.vue-flow__controls {
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-sm);
  overflow: hidden;
  box-shadow: var(--mapex-shadow-sm);
}

.vue-flow__controls-button {
  background: var(--mapex-surface-elevated) !important;
  border-color: var(--mapex-divider) !important;
  color: var(--mapex-text-secondary) !important;
  fill: var(--mapex-text-secondary) !important;
  transition: var(--mapex-transition-fast);

  &:hover {
    background: var(--mapex-surface-bg) !important;
    color: var(--mapex-text-primary) !important;
    fill: var(--mapex-text-primary) !important;
  }

  svg {
    fill: inherit !important;
  }
}

/* ── Edges ── */

.vue-flow__edge-path {
  stroke-linecap: round;
}

.vue-flow__edge.animated .vue-flow__edge-path {
  stroke-dasharray: 5;
  animation: edge-flow 0.5s linear infinite;
}

@keyframes edge-flow {
  to {
    stroke-dashoffset: -10;
  }
}

.vue-flow__edge:hover .vue-flow__edge-path {
  stroke-width: 4 !important;
  filter: drop-shadow(0 0 6px var(--wf-edge-glow));
}

.vue-flow__edge.selected .vue-flow__edge-path {
  stroke: var(--wf-edge-hover) !important;
  stroke-width: 3.5 !important;
  filter: drop-shadow(0 0 4px var(--wf-edge-glow));
}

.vue-flow__edge-interaction {
  stroke-width: 20 !important;
}

/* ── Connection line (while dragging from handle) ── */

.vue-flow__connection-path {
  stroke: var(--wf-edge-color);
  stroke-width: 2.5;
  stroke-dasharray: 5;
  animation: edge-flow 0.5s linear infinite;
}

/* ── Handles ── */

.vue-flow__handle {
  transition: transform 0.15s, box-shadow 0.15s;
}

.vue-flow__handle:hover {
  transform: scale(1.4);
  box-shadow: 0 0 6px var(--wf-edge-glow);
}
</style>
