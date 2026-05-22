/** TYPE IMPORTS */
import type {
  WorkflowGeneralSettings,
  WorkflowVariable,
  CaptureField,
  ExternalSignal,
  ExternalVariable,
  WorkflowNode,
  WorkflowEdge,
  WorkflowDefinition,
  CanvasViewport,
} from '../interfaces';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** CONSTANTS */
import { DEFAULT_GENERAL_SETTINGS, DEFAULT_VALUE_BY_TYPE } from '../constants';
import { START_NODE_TYPE, START_NODE_ID } from '@src/components/workflow/constants';

/** UTILS */
import { resolveNodeHandles, findOrphanedEdges, detectTightCycles } from '../utils';

/** STORES */
import { usePluginRegistryStore } from '@stores/pluginRegistry';

/**
 * Default canvas viewport
 */
const DEFAULT_VIEWPORT: CanvasViewport = { x: 0, y: 0, zoom: 1 };

/**
 * Default start node — always present in every workflow
 */
const DEFAULT_START_NODE: WorkflowNode = {
  id: START_NODE_ID,
  type: START_NODE_TYPE,
  position: { x: 250, y: 50 },
  config: {},
  label: 'Start',
};

/** STATE */

/**
 * General settings for the workflow
 */
const generalSettings = ref<WorkflowGeneralSettings>({ ...DEFAULT_GENERAL_SETTINGS });

/**
 * State variables (persist during execution)
 */
const states = ref<WorkflowVariable[]>([]);

/**
 * Capture fields (stored in ClickHouse)
 */
const captureFields = ref<CaptureField[]>([]);

/**
 * External inputs (input contract for callers)
 */
const externalInputs = ref<ExternalVariable[]>([]);

/**
 * External signals (signal contract for wait_signal nodes)
 */
const externalSignals = ref<ExternalSignal[]>([]);

/**
 * DAG nodes on canvas
 */
const nodes = ref<WorkflowNode[]>([]);

/**
 * DAG edges (connections)
 */
const edges = ref<WorkflowEdge[]>([]);

/**
 * Canvas viewport metadata
 */
const viewport = ref<CanvasViewport>({ ...DEFAULT_VIEWPORT });

/**
 * Incremented every time a node config changes.
 * Used by WorkflowCanvas to sync config-driven changes (like dynamic handles) to Vue Flow.
 */
const nodeConfigVersion = ref(0);

/**
 * Installed marketplace plugin IDs (e.g., ['telegram', 'slack'])
 */
const installedPlugins = ref<string[]>([]);

/**
 * Definition status computed by the backend
 */
const definitionStatus = ref<'valid' | 'plugin_missing' | 'invalid'>('valid');

/**
 * Missing plugin IDs computed by the backend
 */
const missingPlugins = ref<string[]>([]);

/**
 * Per-node validation errors keyed by node ID
 */
const nodeValidationErrors = ref<Record<string, string[]>>({});

/**
 * Central state composable for the workflow editor page.
 * Manages all data across the 3 tabs: General, State, Workflow.
 *
 * @returns {object} Reactive state, computed values, and methods
 */
export function useWorkflowEditorState() {
  /** COMPUTED */

  /**
   * Total variables count (variables + capture fields)
   * @returns {number} Combined count
   */
  const variablesCount = computed(() => {
    return states.value.length + captureFields.value.length + externalInputs.value.length + externalSignals.value.length;
  });

  /**
   * Total nodes count on canvas
   * @returns {number} Node count
   */
  const nodesCount = computed(() => nodes.value.length);

  /**
   * Serialize all state into a WorkflowDefinition JSON
   * @returns {WorkflowDefinition} Complete workflow definition
   */
  const getCurrentWorkflow = computed<WorkflowDefinition>(() => {
    // Clean edges — omit default pathOffset values and strip legacy annotation edges
    const cleanEdges: WorkflowEdge[] = edges.value.filter(e => e.sourceHandle !== '__note_out').map(e => {
      const clean: WorkflowEdge = {
        id: e.id,
        source: e.source,
        target: e.target,
      };
      if (e.sourceHandle !== undefined) clean.sourceHandle = e.sourceHandle;
      if (e.targetHandle !== undefined) clean.targetHandle = e.targetHandle;
      if (e.label !== undefined) clean.label = e.label;
      if (e.pathOffsetX) clean.pathOffsetX = e.pathOffsetX;
      if (e.pathOffsetY) clean.pathOffsetY = e.pathOffsetY;
      return clean;
    });

    return {
      name: generalSettings.value.name,
      description: generalSettings.value.description,
      enabled: generalSettings.value.enabled,
      isTemplate: generalSettings.value.sharedWithChildren,
      definitionVersion: 1,
      timezone: { ...generalSettings.value.timezone },
      retryPolicy: { ...generalSettings.value.retryPolicy },
      states: states.value,
      captureFields: captureFields.value,
      externalInputs: externalInputs.value,
      externalSignals: externalSignals.value,
      nodes: nodes.value,
      edges: cleanEdges,
      installedPlugins: [...installedPlugins.value],
      metadata: {
        canvasViewport: { ...viewport.value },
      },
    };
  });

  /** METHODS */

  /**
   * Set all states from a complete workflow definition
   * Used when loading an existing workflow in edit mode
   *
   * @param {WorkflowDefinition} workflow - Complete workflow definition
   * @returns {void}
   */
  function setAllStates(workflow: WorkflowDefinition): void {
    generalSettings.value = {
      name: workflow.name || '',
      description: workflow.description || '',
      enabled: workflow.enabled ?? true,
      isTemplate: false,
      sharedWithChildren: workflow.isTemplate ?? false,
      timezone: workflow.timezone || { ...DEFAULT_GENERAL_SETTINGS.timezone },
      retryPolicy: workflow.retryPolicy || { ...DEFAULT_GENERAL_SETTINGS.retryPolicy },
    };

    states.value = workflow.states || [];
    captureFields.value = workflow.captureFields || [];
    externalInputs.value = workflow.externalInputs || [];
    externalSignals.value = workflow.externalSignals || [];
    installedPlugins.value = workflow.installedPlugins || [];
    definitionStatus.value = workflow.status || 'valid';
    missingPlugins.value = workflow.missingPlugins || [];

    const loadedNodes = workflow.nodes || [];
    const hasStart = loadedNodes.some(n => n.type === START_NODE_TYPE);
    const resolvedNodes = hasStart
      ? loadedNodes
      : [{ ...DEFAULT_START_NODE, position: { ...DEFAULT_START_NODE.position } }, ...loadedNodes];

    // Migrate text note positions: old format used relative coordinates
    // (e.g., {x: -90, y: -20}), new format uses absolute canvas coordinates.
    // Detect old format by negative x position (canvas coords are positive).
    const parentMap = new Map(resolvedNodes.map(n => [n.id, n]));
    for (const node of resolvedNodes) {
      if (node.type === 'core/text_note' && node.parentNodeId && node.position.x < 0) {
        const parent = parentMap.get(node.parentNodeId);
        if (parent) {
          node.position = {
            x: parent.position.x + node.position.x,
            y: parent.position.y + node.position.y,
          };
        }
      }
    }

    nodes.value = resolvedNodes;

    // Filter out legacy annotation edges (dashed lines removed — notes are standalone nodes)
    edges.value = (workflow.edges || []).filter(e => e.sourceHandle !== '__note_out');
    viewport.value = workflow.metadata?.canvasViewport || { ...DEFAULT_VIEWPORT };
  }

  /**
   * Reset all states to default values
   * Used when creating a new workflow
   *
   * @returns {void}
   */
  function resetAllStates(): void {
    generalSettings.value = { ...DEFAULT_GENERAL_SETTINGS };
    states.value = [];
    captureFields.value = [];
    externalInputs.value = [];
    externalSignals.value = [];
    installedPlugins.value = [];
    definitionStatus.value = 'valid';
    missingPlugins.value = [];
    nodes.value = [{ ...DEFAULT_START_NODE, position: { ...DEFAULT_START_NODE.position } }];
    edges.value = [];
    viewport.value = { ...DEFAULT_VIEWPORT };
  }

  /**
   * Update general settings
   *
   * @param {Partial<WorkflowGeneralSettings>} settings - Partial general settings
   * @returns {void}
   */
  function updateGeneral(settings: Partial<WorkflowGeneralSettings>): void {
    generalSettings.value = { ...generalSettings.value, ...settings };
  }

  /**
   * Add a workflow state variable
   *
   * @param {WorkflowVariable} variable - State variable to add
   * @returns {void}
   */
  function addState(variable: WorkflowVariable): void {
    states.value.push(variable);
  }

  /**
   * Update a workflow state variable by index
   *
   * @param {number} index - State variable index
   * @param {WorkflowVariable} variable - Updated state variable
   * @returns {void}
   */
  function updateState(index: number, variable: WorkflowVariable): void {
    if (index >= 0 && index < states.value.length) {
      states.value[index] = variable;
    }
  }

  /**
   * Remove a workflow state variable by index
   *
   * @param {number} index - State variable index to remove
   * @returns {void}
   */
  function removeState(index: number): void {
    if (index >= 0 && index < states.value.length) {
      states.value.splice(index, 1);
    }
  }

  /**
   * Move a state variable up or down in the list
   *
   * @param {number} index - Current index
   * @param {'up' | 'down'} direction - Move direction
   * @returns {void}
   */
  function moveState(index: number, direction: 'up' | 'down'): void {
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= states.value.length) return;
    const temp = states.value[index]!;
    states.value[index] = states.value[newIndex]!;
    states.value[newIndex] = temp;
  }

  /**
   * Add a capture field
   *
   * @param {CaptureField} field - Capture field to add
   * @returns {void}
   */
  function addCaptureField(field: CaptureField): void {
    captureFields.value.push(field);
  }

  /**
   * Update a capture field by index
   *
   * @param {number} index - Field index
   * @param {CaptureField} field - Updated field
   * @returns {void}
   */
  function updateCaptureField(index: number, field: CaptureField): void {
    if (index >= 0 && index < captureFields.value.length) {
      captureFields.value[index] = field;
    }
  }

  /**
   * Remove a capture field by index
   *
   * @param {number} index - Field index to remove
   * @returns {void}
   */
  function removeCaptureField(index: number): void {
    if (index >= 0 && index < captureFields.value.length) {
      captureFields.value.splice(index, 1);
    }
  }

  /**
   * Move a capture field up or down in the list
   *
   * @param {number} index - Current index
   * @param {'up' | 'down'} direction - Move direction
   * @returns {void}
   */
  function moveCaptureField(index: number, direction: 'up' | 'down'): void {
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= captureFields.value.length) return;
    const temp = captureFields.value[index]!;
    captureFields.value[index] = captureFields.value[newIndex]!;
    captureFields.value[newIndex] = temp;
  }

  /**
   * Add an external variable
   *
   * @param {ExternalVariable} variable - External variable to add
   * @returns {void}
   */
  function addExternalInput(variable: ExternalVariable): void {
    externalInputs.value.push(variable);
  }

  /**
   * Update an external variable by index
   *
   * @param {number} index - Variable index
   * @param {ExternalVariable} variable - Updated variable
   * @returns {void}
   */
  function updateExternalInput(index: number, variable: ExternalVariable): void {
    if (index >= 0 && index < externalInputs.value.length) {
      externalInputs.value[index] = variable;
    }
  }

  /**
   * Remove an external variable by index
   *
   * @param {number} index - Variable index to remove
   * @returns {void}
   */
  function removeExternalInput(index: number): void {
    if (index >= 0 && index < externalInputs.value.length) {
      externalInputs.value.splice(index, 1);
    }
  }

  /**
   * Move an external variable up or down in the list
   *
   * @param {number} index - Current index
   * @param {'up' | 'down'} direction - Move direction
   * @returns {void}
   */
  function moveExternalInput(index: number, direction: 'up' | 'down'): void {
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= externalInputs.value.length) return;
    const temp = externalInputs.value[index]!;
    externalInputs.value[index] = externalInputs.value[newIndex]!;
    externalInputs.value[newIndex] = temp;
  }

  /**
   * Add an external signal
   *
   * @param {ExternalSignal} signal - External signal to add
   * @returns {void}
   */
  function addExternalSignal(signal: ExternalSignal): void {
    externalSignals.value.push(signal);
  }

  /**
   * Update an external signal by index
   *
   * @param {number} index - Signal index
   * @param {ExternalSignal} signal - Updated signal
   * @returns {void}
   */
  function updateExternalSignal(index: number, signal: ExternalSignal): void {
    if (index >= 0 && index < externalSignals.value.length) {
      externalSignals.value[index] = signal;
    }
  }

  /**
   * Remove an external signal by index
   *
   * @param {number} index - Signal index to remove
   * @returns {void}
   */
  function removeExternalSignal(index: number): void {
    if (index >= 0 && index < externalSignals.value.length) {
      externalSignals.value.splice(index, 1);
    }
  }

  /**
   * Move an external signal up or down in the list
   *
   * @param {number} index - Current index
   * @param {'up' | 'down'} direction - Move direction
   * @returns {void}
   */
  function moveExternalSignal(index: number, direction: 'up' | 'down'): void {
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= externalSignals.value.length) return;
    const temp = externalSignals.value[index]!;
    externalSignals.value[index] = externalSignals.value[newIndex]!;
    externalSignals.value[newIndex] = temp;
  }

  /**
   * Update canvas nodes (called by VueFlow onChange)
   *
   * @param {WorkflowNode[]} newNodes - Updated nodes array
   * @returns {void}
   */
  function updateNodes(newNodes: WorkflowNode[]): void {
    nodes.value = newNodes;
  }

  /**
   * Update canvas edges (called by VueFlow onChange)
   *
   * @param {WorkflowEdge[]} newEdges - Updated edges array
   * @returns {void}
   */
  function updateEdges(newEdges: WorkflowEdge[]): void {
    edges.value = newEdges;
  }

  /**
   * Add a single node to the canvas
   *
   * @param {WorkflowNode} node - Node to add
   * @returns {void}
   */
  function addNode(node: WorkflowNode): void {
    nodes.value = [...nodes.value, node];
  }

  /**
   * Remove a node, its child nodes (e.g., notes) and connected edges.
   * Returns the list of all removed node IDs so callers can sync Vue Flow.
   *
   * @param {string} nodeId - Node ID to remove
   * @returns {string[]} Array of removed node IDs (empty if blocked)
   */
  function removeNode(nodeId: string): string[] {
    const node = nodes.value.find(n => n.id === nodeId);
    if (!node) return [];

    // Prevent deletion of undeletable nodes (e.g., Start)
    const pluginRegistry = usePluginRegistryStore();
    const nodeType = pluginRegistry.getNodeType(node.type);
    if (nodeType?.deletable === false) return [];

    // Collect child nodes attached to this node (e.g., text notes)
    const childIds = nodes.value
      .filter(n => n.parentNodeId === nodeId)
      .map(n => n.id);

    const idsToRemove = [nodeId, ...childIds];
    const idSet = new Set(idsToRemove);

    nodes.value = nodes.value.filter(n => !idSet.has(n.id));
    edges.value = edges.value.filter(e => !idSet.has(e.source) && !idSet.has(e.target));
    return idsToRemove;
  }

  /**
   * Update a single node's config.
   * Increments nodeConfigVersion for Vue Flow sync.
   * If the node type has dynamic resolvers, cleans orphaned edges.
   *
   * @param {string} nodeId - Node ID
   * @param {Record<string, unknown>} config - Updated config
   * @returns {void}
   */
  function updateNodeConfig(nodeId: string, config: Record<string, unknown>): void {
    const node = nodes.value.find(n => n.id === nodeId);
    if (!node) return;

    node.config = { ...node.config, ...config };
    nodeConfigVersion.value++;

    const pluginRegistry = usePluginRegistryStore();
    const nodeType = pluginRegistry.getNodeType(node.type);
    if (!nodeType?.resolveOutputs && !nodeType?.resolveInputs) return;

    const resolved = resolveNodeHandles(nodeType, node.config);
    const orphaned = findOrphanedEdges(nodeId, resolved.inputs, resolved.outputs, edges.value);

    if (orphaned.length) {
      edges.value = edges.value.filter(e => !orphaned.includes(e.id));
    }
  }

  /**
   * Replace a node's config entirely (instead of merging).
   * Used when operation changes to discard fields from the previous operation.
   *
   * @param {string} nodeId - Node ID
   * @param {Record<string, unknown>} config - New config (replaces existing)
   * @returns {void}
   */
  function replaceNodeConfig(nodeId: string, config: Record<string, unknown>): void {
    const node = nodes.value.find(n => n.id === nodeId);
    if (!node) return;

    node.config = { ...config };
    nodeConfigVersion.value++;

    const pluginRegistry = usePluginRegistryStore();
    const nodeType = pluginRegistry.getNodeType(node.type);
    if (!nodeType?.resolveOutputs && !nodeType?.resolveInputs) return;

    const resolved = resolveNodeHandles(nodeType, node.config);
    const orphaned = findOrphanedEdges(nodeId, resolved.inputs, resolved.outputs, edges.value);

    if (orphaned.length) {
      edges.value = edges.value.filter(e => !orphaned.includes(e.id));
    }
  }

  /**
   * Add an edge between nodes
   *
   * @param {WorkflowEdge} edge - Edge to add
   * @returns {void}
   */
  function addEdge(edge: WorkflowEdge): void {
    edges.value = [...edges.value, edge];
  }

  /**
   * Duplicate a node on the canvas.
   * Skips nodes marked as non-deletable (e.g., Start).
   *
   * @param {string} nodeId - ID of the node to duplicate
   * @returns {string | null} The new node ID, or null if duplication was skipped
   */
  function duplicateNode(nodeId: string): string | null {
    const sourceNode = nodes.value.find(n => n.id === nodeId);
    if (!sourceNode) return null;

    const pluginRegistry = usePluginRegistryStore();
    const nodeType = pluginRegistry.getNodeType(sourceNode.type);
    if (nodeType?.deletable === false) return null;

    const newId = `n_${sourceNode.type.replace('/', '_')}_${Date.now()}`;

    const newNode: WorkflowNode = {
      id: newId,
      type: sourceNode.type,
      position: { x: sourceNode.position.x + 40, y: sourceNode.position.y + 40 },
      config: { ...sourceNode.config },
      label: nodeType?.label || sourceNode.label || '',
    };

    nodes.value = [...nodes.value, newNode];
    return newId;
  }

  /**
   * Remove an edge by ID
   *
   * @param {string} edgeId - Edge ID to remove
   * @returns {void}
   */
  function removeEdge(edgeId: string): void {
    edges.value = edges.value.filter(e => e.id !== edgeId);
  }

  /**
   * Add a text note balloon attached to a target node.
   * Creates a TextNoteNode positioned to the left of the target
   * and an annotation edge connecting them.
   *
   * @param {string} targetNodeId - ID of the node to attach the note to
   * @returns {void}
   */
  function addNoteToNode(targetNodeId: string): void {
    const targetNode = nodes.value.find(n => n.id === targetNodeId);
    if (!targetNode) return;

    const noteId = `note_${Date.now()}`;
    const noteNode: WorkflowNode = {
      id: noteId,
      type: 'core/text_note',
      parentNodeId: targetNodeId,
      position: { x: targetNode.position.x - 90, y: targetNode.position.y + 10 },
      config: { text: '', color: 'grey' },
    };

    nodes.value = [...nodes.value, noteNode];
  }

  /**
   * Update canvas viewport
   *
   * @param {CanvasViewport} newViewport - Updated viewport
   * @returns {void}
   */
  function updateViewport(newViewport: CanvasViewport): void {
    viewport.value = newViewport;
  }

  /**
   * Add a marketplace plugin to the installed list
   *
   * @param {string} pluginId - Plugin ID to add
   * @returns {void}
   */
  function addInstalledPlugin(pluginId: string): void {
    if (!installedPlugins.value.includes(pluginId)) {
      installedPlugins.value.push(pluginId);
    }
  }

  /**
   * Remove a marketplace plugin from the installed list
   *
   * @param {string} pluginId - Plugin ID to remove
   * @returns {void}
   */
  function removeInstalledPlugin(pluginId: string): void {
    installedPlugins.value = installedPlugins.value.filter(id => id !== pluginId);
  }

  /**
   * Get default value for a variable type
   *
   * @param {string} type - Variable type
   * @returns {string | number | boolean | Record<string, unknown>} Default value
   */
  function getDefaultValue(type: string): string | number | boolean | Record<string, unknown> {
    return DEFAULT_VALUE_BY_TYPE[type] ?? '';
  }

  /**
   * Validate all nodes by calling each node type's validate function.
   * Also checks required properties for declarative nodes.
   * Stores results in nodeValidationErrors.
   *
   * @returns {number} Number of nodes with errors
   */
  function validateAllNodes(): number {
    const pluginRegistry = usePluginRegistryStore();
    const errors: Record<string, string[]> = {};

    /** Annotations — exempt from all connectivity checks */
    const ANNOTATION_TYPES = new Set(['core/text_note', 'core/group_frame']);

    /** Nodes exempt from input check (they are entry points) */
    const NO_INPUT_CHECK = new Set(['core/start']);

    /** Nodes exempt from output check (terminal or logical routing via backend) */
    const NO_OUTPUT_CHECK = new Set(['core/end', 'core/goto']);

    /** Build incoming/outgoing sets from edges */
    const hasIncoming = new Set<string>();
    const hasOutgoing = new Set<string>();
    for (const edge of edges.value) {
      hasOutgoing.add(edge.source);
      hasIncoming.add(edge.target);
    }

    /** Detect tight cycles (no async pause point) */
    const tightCycleNodeIds = detectTightCycles(nodes.value, edges.value);

    for (const node of nodes.value) {
      const nodeType = pluginRegistry.getNodeType(node.type);
      if (!nodeType) continue;

      const nodeErrors: string[] = [];

      // 1. Call plugin validate function if defined
      if (nodeType.validate) {
        const result = nodeType.validate(node.config);
        if (!result.valid) nodeErrors.push(...result.errors);
      }

      // 2. Check required properties (for declarative nodes)
      //    Only validate properties that are currently visible (respects displayOptions.show)
      if (nodeType.properties) {
        for (const prop of nodeType.properties) {
          if (!prop.required) continue;

          // Skip hidden properties based on displayOptions.show
          if (prop.displayOptions?.show) {
            const visible = Object.entries(prop.displayOptions.show).every(
              ([field, values]) => (values as string[]).includes(node.config[field] as string),
            );
            if (!visible) continue;
          }

          const val = node.config[prop.name];
          if (val === undefined || val === null || val === '') {
            nodeErrors.push(`propRequired::${prop.displayName}`);
          }
        }
      }

      // 3. Check connectivity — inputs and outputs must be connected
      if (!ANNOTATION_TYPES.has(node.type)) {
        const needsInput = !NO_INPUT_CHECK.has(node.type) && (nodeType.inputs?.length ?? 0) > 0;
        const needsOutput = !NO_OUTPUT_CHECK.has(node.type) && (nodeType.outputs?.length ?? 0) > 0;

        const isGotoReceiver = node.type === 'core/goto' && node.config.role === 'receiver';
        if (needsInput && !hasIncoming.has(node.id) && !isGotoReceiver) {
          nodeErrors.push('noIncomingConnection');
        }
        if (needsOutput && !hasOutgoing.has(node.id)) {
          nodeErrors.push('noOutgoingConnection');
        }
      }

      // 4. Check tight cycles — no async pause point in cycle
      if (tightCycleNodeIds.has(node.id)) {
        nodeErrors.push('tightCycleDetected');
      }

      if (nodeErrors.length > 0) errors[node.id] = nodeErrors;
    }

    // 5. Cross-node: goto sender/receiver pair validation
    const gotoNodes = nodes.value.filter(n => n.type === 'core/goto');
    const senderLabels = new Set<string>();
    const receiverLabels = new Set<string>();
    for (const g of gotoNodes) {
      const label = (g.config.pairLabel as string)?.trim();
      if (!label) continue;
      if (g.config.role === 'sender') senderLabels.add(label);
      else if (g.config.role === 'receiver') receiverLabels.add(label);
    }
    for (const g of gotoNodes) {
      const label = (g.config.pairLabel as string)?.trim();
      if (!label) continue;
      if (g.config.role === 'sender' && !receiverLabels.has(label)) {
        const nodeErr = errors[g.id] ?? [];
        nodeErr.push(`gotoSenderNeedsReceiver::${label}`);
        errors[g.id] = nodeErr;
      }
      if (g.config.role === 'receiver' && !senderLabels.has(label)) {
        const nodeErr = errors[g.id] ?? [];
        nodeErr.push(`gotoReceiverNeedsSender::${label}`);
        errors[g.id] = nodeErr;
      }
    }

    nodeValidationErrors.value = errors;
    if (Object.keys(errors).length > 0) {
      console.warn('[Workflow] Validation errors:', JSON.stringify(errors, null, 2));
    }
    return Object.keys(errors).length;
  }

  /**
   * Clear all per-node validation errors
   *
   * @returns {void}
   */
  function clearNodeValidationErrors(): void {
    nodeValidationErrors.value = {};
  }

  return {
    // Reactive states
    generalSettings,
    states,
    captureFields,
    externalInputs,
    externalSignals,
    installedPlugins,
    definitionStatus,
    missingPlugins,
    nodes,
    edges,
    viewport,
    nodeConfigVersion,
    nodeValidationErrors,

    // Computed
    variablesCount,
    nodesCount,
    getCurrentWorkflow,

    // General
    setAllStates,
    resetAllStates,
    updateGeneral,

    // States
    addState,
    updateState,
    removeState,
    moveState,

    // Capture Fields
    addCaptureField,
    updateCaptureField,
    removeCaptureField,
    moveCaptureField,

    // External Inputs
    addExternalInput,
    updateExternalInput,
    removeExternalInput,
    moveExternalInput,

    // External Signals
    addExternalSignal,
    updateExternalSignal,
    removeExternalSignal,
    moveExternalSignal,

    // Installed plugins
    addInstalledPlugin,
    removeInstalledPlugin,

    // Nodes
    updateNodes,
    addNode,
    removeNode,
    duplicateNode,
    updateNodeConfig,
    replaceNodeConfig,

    // Edges
    updateEdges,
    addEdge,
    removeEdge,

    // Notes
    addNoteToNode,

    // Viewport
    updateViewport,

    // Validation
    validateAllNodes,
    clearNodeValidationErrors,

    // Utils
    getDefaultValue,
  };
}
