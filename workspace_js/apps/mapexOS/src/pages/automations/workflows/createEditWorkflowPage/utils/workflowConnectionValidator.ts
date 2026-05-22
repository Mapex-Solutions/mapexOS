/** TYPE IMPORTS */
import type { Connection } from '@vue-flow/core';
import type {
  WorkflowNode,
  WorkflowEdge,
  PluginNodeType,
  HandleDefinition,
  ResolvedHandles,
} from '../interfaces/CreateEditWorkflow.interface';

/**
 * Context needed by the connection validator.
 * Provided as a getter so it always reads fresh reactive state.
 */
export interface ConnectionValidationContext {
  /** Look up a plugin node type definition by its type string */
  getNodeType: (typeId: string) => PluginNodeType | undefined;

  /** Current nodes on canvas */
  nodes: WorkflowNode[];

  /** Current edges on canvas */
  edges: WorkflowEdge[];

  /** Resolve dynamic handles for a node type + config + timeout */
  resolveHandles: (nodeType: PluginNodeType, config: Record<string, unknown>, timeout?: { duration: number; unit: string; enableOutput: boolean }) => ResolvedHandles;
}

/**
 * Create a connection validator for the workflow canvas.
 * Returns a function compatible with Vue Flow's `isValidConnection` prop.
 *
 * IMPORTANT: Vue Flow calls this function in TWO contexts:
 * 1. During drag — receives a Connection (no `id`)
 * 2. Inside createGraphEdges (on setEdges/addEdges) — receives a GraphEdge (has `id`)
 * In context 2, the edge already exists in our state, so duplicate/maxConnection
 * checks must exclude the edge itself to avoid self-rejection.
 *
 * Generic rules applied to ALL node types:
 * 1. No self-connections (node → itself)
 * 2. Output → Input only (source handle must be an output, target handle must be an input)
 * 3. No duplicate connections (same source+handle → same target+handle)
 * 4. Max connections per handle (if `maxConnections` is defined in HandleDefinition)
 *
 * @param {() => ConnectionValidationContext} getContext - Getter for fresh validation context
 * @returns {(connection: Connection) => boolean} Validator function for Vue Flow
 */
export function createConnectionValidator(
  getContext: () => ConnectionValidationContext,
): (connection: Connection) => boolean {
  return function isValidConnection(connection: Connection): boolean {
    try {
      const ctx = getContext();

      // Rule 1: No self-connections
      if (connection.source === connection.target) return false;

      // When called from createGraphEdges, the parameter is a GraphEdge (has `id`).
      // We need this to exclude the edge itself from duplicate/maxConnection checks.
      const edgeId = (connection as unknown as { id?: string }).id;

      // Look up nodes in our composable state
      const sourceNode = ctx.nodes.find(n => n.id === connection.source);
      const targetNode = ctx.nodes.find(n => n.id === connection.target);

      if (!sourceNode || !targetNode) {
        // Fail-open: allow if nodes not synced yet
        return true;
      }

      // Look up node type definitions from plugin registry
      const sourceType = ctx.getNodeType(sourceNode.type);
      const targetType = ctx.getNodeType(targetNode.type);

      if (!sourceType || !targetType) {
        // Fail-open: allow if plugin not loaded
        return true;
      }

      // Resolve handles (supports dynamic handles like Fanout and Timeout)
      const sourceHandles = ctx.resolveHandles(sourceType, sourceNode.config, sourceNode.timeout);
      const targetHandles = ctx.resolveHandles(targetType, targetNode.config, targetNode.timeout);

      // Rule 2: Source handle must be an OUTPUT, target handle must be an INPUT
      const sourceHandle = findHandle(sourceHandles.outputs, connection.sourceHandle);
      const targetHandle = findHandle(targetHandles.inputs, connection.targetHandle);

      if (!sourceHandle || !targetHandle) {
        // Blocks IN→IN and OUT→OUT connections
        return false;
      }

      // Rule 3: No duplicate connections (exclude self when re-validating existing edges)
      const isDuplicate = ctx.edges.some(
        e =>
          e.id !== edgeId &&
          e.source === connection.source &&
          e.sourceHandle === connection.sourceHandle &&
          e.target === connection.target &&
          e.targetHandle === connection.targetHandle,
      );
      if (isDuplicate) return false;

      // Rule 4: Max connections per handle (exclude self when re-validating)
      if (!isWithinMaxConnections(sourceHandle, connection.source, 'source', connection.sourceHandle, ctx.edges, edgeId)) {
        return false;
      }
      if (!isWithinMaxConnections(targetHandle, connection.target, 'target', connection.targetHandle, ctx.edges, edgeId)) {
        return false;
      }

      return true;
    } catch (err) {
      console.warn('[ConnectionValidator] Validation error, allowing connection.', err);
      return true;
    }
  };
}

/**
 * Find a handle definition by ID.
 * When handleId is null/undefined, returns the first handle (single-handle nodes).
 *
 * @param {HandleDefinition[]} handles - Handle definitions to search
 * @param {string | null | undefined} handleId - Handle ID to find
 * @returns {HandleDefinition | undefined} Found handle or undefined
 */
function findHandle(
  handles: HandleDefinition[],
  handleId: string | null | undefined,
): HandleDefinition | undefined {
  if (!handleId) return handles[0];
  return handles.find(h => h.id === handleId);
}

/**
 * Check if adding one more connection to a handle stays within its maxConnections limit.
 * Excludes the edge itself when re-validating existing edges (avoids self-counting).
 *
 * @param {HandleDefinition} handle - Handle definition with optional maxConnections
 * @param {string} nodeId - Node ID owning the handle
 * @param {'source' | 'target'} side - Which side of the edge this handle is on
 * @param {string | null | undefined} handleId - Handle ID
 * @param {WorkflowEdge[]} edges - Current edges
 * @param {string | undefined} excludeEdgeId - Edge ID to exclude from count (self)
 * @returns {boolean} True if within limit
 */
function isWithinMaxConnections(
  handle: HandleDefinition,
  nodeId: string,
  side: 'source' | 'target',
  handleId: string | null | undefined,
  edges: WorkflowEdge[],
  excludeEdgeId?: string,
): boolean {
  if (handle.maxConnections == null) return true;

  const current = edges.filter(e => {
    if (excludeEdgeId && e.id === excludeEdgeId) return false;
    if (side === 'source') {
      return e.source === nodeId && e.sourceHandle === handleId;
    }
    return e.target === nodeId && e.targetHandle === handleId;
  }).length;

  return current < handle.maxConnections;
}
