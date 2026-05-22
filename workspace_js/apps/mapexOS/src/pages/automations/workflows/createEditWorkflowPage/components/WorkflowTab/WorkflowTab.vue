<script setup lang="ts">
/** TYPE IMPORTS */
import type { CanvasToolbarState } from '../../interfaces/CreateEditWorkflow.interface';
import type { IWorkflowEditorContext } from '@src/components/workflow/interfaces';

/** VUE IMPORTS */
import { ref, computed, nextTick, provide } from 'vue';

/** COMPONENTS */
import CanvasToolbar from '../CanvasToolbar/CanvasToolbar.vue';
import PluginCatalog from '../PluginCatalog/PluginCatalog.vue';
import WorkflowCanvas from '../WorkflowCanvas/WorkflowCanvas.vue';
import NodeConfigPanel from '../NodeConfigPanel/NodeConfigPanel.vue';

/** COMPOSABLES */
import { useWorkflowEditorState, useWorkflowHistory } from '../../composables';
import { WORKFLOW_CONTEXT_KEY } from '@src/composables/workflow';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/** UTILS */
import { autoLayoutNodes } from './utils';

/** LOCAL IMPORTS (constants and handlers ONLY) */
import { LAYOUT_NODE_WIDTH, LAYOUT_NODE_HEIGHT, SNAP_GRID_SIZE } from '../../constants';
import { DEFAULT_TOOLBAR_STATE } from '../CanvasToolbar/constants';
import { CATALOG_WIDTH, CATALOG_WIDTH_COLLAPSED } from '../PluginCatalog/constants';
import { CONFIG_PANEL_WIDTH } from '../NodeConfigPanel/constants';

/** COMPOSABLES & STORES */
const { nodes, edges, states, externalSignals, updateNodeConfig, addNoteToNode } = useWorkflowEditorState();
const { pushSnapshot, undo, redo, finishRestore } = useWorkflowHistory();
const pluginRegistry = usePluginRegistryStore();

/**
 * Provide workflow editor context for plugin components.
 * All plugin nodes/configs use useWorkflowContext() which injects this.
 */
provide<IWorkflowEditorContext>(WORKFLOW_CONTEXT_KEY, {
  nodes,
  edges,
  states,
  externalSignals,
  updateNodeConfig,
  addNoteToNode,
  pushSnapshot,
  getNodeType: (type: string) => pluginRegistry.getNodeType(type),
});

/** STATE */

/**
 * Canvas toolbar state
 */
const toolbarState = ref<CanvasToolbarState>({ ...DEFAULT_TOOLBAR_STATE });

/**
 * Reference to the WorkflowCanvas component for fitView
 */
const canvasRef = ref<InstanceType<typeof WorkflowCanvas> | null>(null);

/**
 * Whether the plugin catalog sidebar is collapsed
 */
const catalogCollapsed = ref(false);

/**
 * Currently selected node ID (null = none selected)
 */
const selectedNodeId = ref<string | null>(null);

/** COMPUTED */

/**
 * Whether the config panel should be visible
 */
const showConfigPanel = computed(() => selectedNodeId.value !== null);

/**
 * Current catalog width based on collapse state
 */
const catalogWidth = computed(() =>
  catalogCollapsed.value ? CATALOG_WIDTH_COLLAPSED : CATALOG_WIDTH,
);

/** FUNCTIONS */

/**
 * Handle node selection on canvas
 *
 * @param {string | null} nodeId - Selected node ID or null
 * @returns {void}
 */
function handleNodeSelect(nodeId: string | null): void {
  selectedNodeId.value = nodeId;
}

/**
 * Handle canvas click (deselect)
 *
 * @returns {void}
 */
function handleCanvasClick(): void {
  selectedNodeId.value = null;
}

/**
 * Toggle catalog collapse
 *
 * @returns {void}
 */
function toggleCatalog(): void {
  catalogCollapsed.value = !catalogCollapsed.value;
}

/**
 * Handle undo action from toolbar button.
 * Restores previous snapshot to both composable and Vue Flow state.
 *
 * @returns {void}
 */
function handleUndo(): void {
  const snapshot = undo();
  if (!snapshot) return;
  canvasRef.value?.restoreSnapshot(snapshot.nodes, snapshot.edges);
  void nextTick(() => finishRestore());
}

/**
 * Handle redo action from toolbar button.
 * Re-applies the last undone snapshot to both composable and Vue Flow state.
 *
 * @returns {void}
 */
function handleRedo(): void {
  const snapshot = redo();
  if (!snapshot) return;
  canvasRef.value?.restoreSnapshot(snapshot.nodes, snapshot.edges);
  void nextTick(() => finishRestore());
}

/**
 * Auto-organize all nodes using dagre hierarchical layout.
 * Updates composable state and triggers fitView after sync.
 *
 * @returns {void}
 */
function handleAutoOrganize(): void {
  pushSnapshot('Auto-organize');

  const result = autoLayoutNodes(nodes.value, edges.value, {
    nodeWidth: LAYOUT_NODE_WIDTH,
    nodeHeight: LAYOUT_NODE_HEIGHT,
    snapGrid: SNAP_GRID_SIZE,
  });

  // Apply directly to Vue Flow nodes (bypasses watcher timing issues)
  canvasRef.value?.applyLayout(result.nodes);

  void nextTick(() => {
    canvasRef.value?.fitView();
  });
}
</script>

<template>
  <div class="workflow-tab">
    <!-- Inline (normal) mode -->
    <template v-if="!toolbarState.maximized">
      <!-- Canvas Toolbar -->
      <CanvasToolbar v-model="toolbarState" @auto-organize="handleAutoOrganize" @undo="handleUndo" @redo="handleRedo" />

      <!-- Three-column layout -->
      <div class="workflow-layout">
        <!-- Left: Plugin Catalog -->
        <div
          class="workflow-layout__catalog"
          :style="{ width: `${catalogWidth}px` }"
        >
          <PluginCatalog
            :collapsed="catalogCollapsed"
            @toggle-collapse="toggleCatalog"
          />
        </div>

        <!-- Center: Canvas -->
        <div class="workflow-layout__canvas">
          <WorkflowCanvas
            ref="canvasRef"
            :toolbar-state="toolbarState"
            @node-select="handleNodeSelect"
            @canvas-click="handleCanvasClick"
          />
        </div>

        <!-- Right: Node Config Panel -->
        <div
          v-if="showConfigPanel"
          class="workflow-layout__config"
          :style="{ width: `${CONFIG_PANEL_WIDTH}px` }"
        >
          <NodeConfigPanel
            :node-id="selectedNodeId!"
            @close="selectedNodeId = null"
          />
        </div>
      </div>
    </template>

    <!-- Fullscreen dialog mode -->
    <q-dialog
      :model-value="toolbarState.maximized"
      maximized
      persistent
      transition-show="fade"
      transition-hide="fade"
      @update:model-value="(val: boolean) => toolbarState = { ...toolbarState, maximized: val }"
    >
      <q-card class="workflow-fullscreen">
        <!-- Toolbar inside dialog -->
        <CanvasToolbar v-model="toolbarState" @auto-organize="handleAutoOrganize" @undo="handleUndo" @redo="handleRedo" />

        <!-- Three-column layout (same structure, full height) -->
        <div class="workflow-layout workflow-layout--fullscreen">
          <!-- Left: Plugin Catalog -->
          <div
            class="workflow-layout__catalog"
            :style="{ width: `${catalogWidth}px` }"
          >
            <PluginCatalog
              :collapsed="catalogCollapsed"
              @toggle-collapse="toggleCatalog"
            />
          </div>

          <!-- Center: Canvas -->
          <div class="workflow-layout__canvas">
            <WorkflowCanvas
              ref="canvasRef"
              :toolbar-state="toolbarState"
              @node-select="handleNodeSelect"
              @canvas-click="handleCanvasClick"
            />
          </div>

          <!-- Right: Node Config Panel -->
          <div
            v-if="showConfigPanel"
            class="workflow-layout__config"
            :style="{ width: `${CONFIG_PANEL_WIDTH}px` }"
          >
            <NodeConfigPanel
              :node-id="selectedNodeId!"
              @close="selectedNodeId = null"
            />
          </div>
        </div>
      </q-card>
    </q-dialog>
  </div>
</template>

<style lang="scss" scoped>
.workflow-tab {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 200px);
  min-height: 500px;
}

.workflow-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
  border: 1px solid var(--mapex-card-border);
  border-radius: var(--mapex-radius-md);

  &--fullscreen {
    border: none;
    border-radius: 0;
    height: calc(100vh - 40px);
  }

  &__catalog {
    flex-shrink: 0;
    border-right: 1px solid var(--mapex-card-border);
    overflow-y: auto;
    transition: width 0.2s ease;
  }

  &__canvas {
    flex: 1;
    position: relative;
    overflow: hidden;
  }

  &__config {
    flex-shrink: 0;
    border-left: 1px solid var(--mapex-card-border);
    overflow-y: auto;
    transition: width 0.2s ease;
  }
}

.workflow-fullscreen {
  display: flex;
  flex-direction: column;
  background: var(--mapex-page-bg);
  max-width: 100vw !important;
  max-height: 100vh !important;
}
</style>
