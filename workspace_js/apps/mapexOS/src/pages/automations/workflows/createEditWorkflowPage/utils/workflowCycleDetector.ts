import type { WorkflowNode, WorkflowEdge } from '../interfaces';

/** Node types that represent async pause points — cycles through these are allowed */
const ASYNC_NODE_TYPES = new Set([
  'core/wait_signal',
  'core/wait_for',
  'core/delay',
]);

/** Annotation types excluded from graph analysis */
const ANNOTATION_TYPES = new Set(['core/text_note', 'core/group_frame']);

/**
 * Check whether an edge is an annotation edge (note connector).
 *
 * @param {WorkflowEdge} edge - Edge to test
 * @returns {boolean} True if the edge connects an annotation handle
 */
function isAnnotationEdge(edge: WorkflowEdge): boolean {
  return edge.sourceHandle === '__note_out' || edge.targetHandle === '__note';
}

/**
 * Check whether an edge is a loop-body back-edge that should be excluded.
 * Loop nodes (core/loop) have a 'body' sourceHandle that feeds back
 * into the loop body — this is an intentional structural cycle.
 *
 * @param {WorkflowEdge} edge - Edge to test
 * @param {Map<string, string>} nodeTypeMap - nodeId → nodeType lookup
 * @returns {boolean} True if the edge is a loop body back-edge
 */
function isLoopBodyEdge(
  edge: WorkflowEdge,
  nodeTypeMap: Map<string, string>,
): boolean {
  return (
    nodeTypeMap.get(edge.source) === 'core/loop' &&
    edge.sourceHandle === 'body'
  );
}

/**
 * Find all elementary cycles in a directed graph using iterative DFS.
 * Each cycle is returned as an array of node IDs forming the loop path.
 *
 * @param {Map<string, string[]>} adjacency - Adjacency list (nodeId → target nodeIds)
 * @returns {string[][]} Array of cycles, each cycle is an array of node IDs
 */
function findAllCycles(adjacency: Map<string, string[]>): string[][] {
  const cycles: string[][] = [];

  /** 0 = white (unvisited), 1 = gray (in stack), 2 = black (done) */
  const color = new Map<string, number>();
  const parent = new Map<string, string | null>();

  for (const nodeId of adjacency.keys()) {
    color.set(nodeId, 0);
  }

  /**
   * Reconstruct cycle path from the back-edge target up to the current node.
   *
   * @param {string} cycleStart - Node where the cycle begins (back-edge target)
   * @param {string} cycleEnd - Node where the back-edge was found
   * @returns {string[]} Ordered list of node IDs in the cycle
   */
  function reconstructCycle(cycleStart: string, cycleEnd: string): string[] {
    const path: string[] = [cycleEnd];
    let current = cycleEnd;

    while (current !== cycleStart) {
      const p = parent.get(current);
      if (p === null || p === undefined) break;
      path.push(p);
      current = p;
    }

    path.reverse();
    return path;
  }

  for (const startNode of adjacency.keys()) {
    if (color.get(startNode) !== 0) continue;

    /** Iterative DFS using explicit stack: [nodeId, neighborIndex] */
    const stack: [string, number][] = [[startNode, 0]];
    color.set(startNode, 1);
    parent.set(startNode, null);

    while (stack.length > 0) {
      const top = stack[stack.length - 1] as [string, number];
      const nodeId = top[0];
      const neighbors = adjacency.get(nodeId) ?? [];

      if (top[1] < neighbors.length) {
        const neighbor = neighbors[top[1]] as string;
        top[1]++;

        const neighborColor = color.get(neighbor) ?? 0;

        if (neighborColor === 0) {
          color.set(neighbor, 1);
          parent.set(neighbor, nodeId);
          stack.push([neighbor, 0]);
        } else if (neighborColor === 1) {
          cycles.push(reconstructCycle(neighbor, nodeId));
        }
      } else {
        color.set(nodeId, 2);
        stack.pop();
      }
    }
  }

  return cycles;
}

/**
 * Detect tight cycles (no async pause point) in the workflow graph.
 * Cycles that pass through wait_signal, wait_for, or delay are allowed.
 * Loop body back-edges (core/loop → body handle) are excluded.
 * Annotation nodes and edges are excluded.
 *
 * @param {WorkflowNode[]} nodes - All workflow nodes
 * @param {WorkflowEdge[]} edges - All workflow edges
 * @returns {Set<string>} Node IDs participating in tight (dangerous) cycles
 */
export function detectTightCycles(
  nodes: WorkflowNode[],
  edges: WorkflowEdge[],
): Set<string> {
  /** 1. Build nodeTypeMap: nodeId → nodeType */
  const nodeTypeMap = new Map<string, string>();
  for (const node of nodes) {
    nodeTypeMap.set(node.id, node.type);
  }

  /** 2. Filter edges — exclude annotations and loop body back-edges */
  const filteredEdges = edges.filter((edge) => {
    if (isAnnotationEdge(edge)) return false;

    const sourceType = nodeTypeMap.get(edge.source) ?? '';
    const targetType = nodeTypeMap.get(edge.target) ?? '';
    if (ANNOTATION_TYPES.has(sourceType) || ANNOTATION_TYPES.has(targetType)) {
      return false;
    }

    if (isLoopBodyEdge(edge, nodeTypeMap)) return false;

    return true;
  });

  /** 3. Build adjacency list */
  const adjacency = new Map<string, string[]>();
  for (const node of nodes) {
    if (!ANNOTATION_TYPES.has(node.type)) {
      adjacency.set(node.id, []);
    }
  }
  for (const edge of filteredEdges) {
    const targets = adjacency.get(edge.source);
    if (targets) targets.push(edge.target);
  }

  /** 4. Find all cycles */
  const cycles = findAllCycles(adjacency);

  /** 5. Filter tight cycles (no async node in cycle) */
  const tightCycleNodes = new Set<string>();

  for (const cycle of cycles) {
    const hasAsync = cycle.some((nodeId) => {
      const nodeType = nodeTypeMap.get(nodeId) ?? '';
      return ASYNC_NODE_TYPES.has(nodeType);
    });

    if (!hasAsync) {
      for (const nodeId of cycle) {
        tightCycleNodes.add(nodeId);
      }
    }
  }

  return tightCycleNodes;
}
