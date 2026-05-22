/** TYPE IMPORTS */
import type { WorkflowNode, WorkflowEdge } from '../../../interfaces/CreateEditWorkflow.interface';

/** EXTERNAL IMPORTS */
import dagre from '@dagrejs/dagre';

/**
 * Options for auto-layout computation
 */
export interface AutoLayoutOptions {
  /** Base node width for dagre (px) — grows with output count */
  nodeWidth?: number;

  /** Base node height for dagre (px) */
  nodeHeight?: number;

  /** Snap positions to grid [x, y] */
  snapGrid?: [number, number];
}

/**
 * Result of auto-layout computation
 */
export interface AutoLayoutResult {
  /** Nodes with updated positions */
  nodes: WorkflowNode[];

  /** Whether any positions actually changed */
  changed: boolean;
}

/**
 * Internal dagre graph configuration for a layout strategy
 */
interface LayoutStrategy {
  /** Layout direction */
  rankdir: 'TB' | 'LR';

  /** Horizontal separation between siblings (px) */
  nodesep: number;

  /** Vertical separation between ranks (px) */
  ranksep: number;

  /** Dagre ranker algorithm */
  ranker: 'network-simplex' | 'tight-tree' | 'longest-path';

  /** Node alignment within rank */
  align?: 'UL' | 'UR' | 'DL' | 'DR';
}

/** DEFAULT VALUES */
const DEFAULT_NODE_WIDTH = 120;
const DEFAULT_NODE_HEIGHT = 70;
const WIDTH_PER_OUTPUT = 60;

/**
 * Predefined layout strategies that cycle on each click.
 * Each strategy produces a visually different arrangement.
 */
const LAYOUT_STRATEGIES: LayoutStrategy[] = [
  // 1) Top-to-bottom, generous spacing, network-simplex (best crossing reduction)
  { rankdir: 'TB', nodesep: 120, ranksep: 140, ranker: 'network-simplex' },
  // 2) Left-to-right, wide horizontal flow
  { rankdir: 'LR', nodesep: 100, ranksep: 160, ranker: 'network-simplex' },
  // 3) Top-to-bottom, tight-tree (different node ordering)
  { rankdir: 'TB', nodesep: 140, ranksep: 120, ranker: 'tight-tree', align: 'UL' },
  // 4) Left-to-right, tight-tree with right alignment
  { rankdir: 'LR', nodesep: 120, ranksep: 140, ranker: 'tight-tree', align: 'DR' },
  // 5) Top-to-bottom, longest-path (maximizes vertical spread)
  { rankdir: 'TB', nodesep: 100, ranksep: 160, ranker: 'longest-path' },
  // 6) Left-to-right, longest-path
  { rankdir: 'LR', nodesep: 140, ranksep: 120, ranker: 'longest-path', align: 'UL' },
];

/** Current strategy index — cycles through LAYOUT_STRATEGIES */
let currentStrategyIndex = 0;

/**
 * Compute auto-layout positions for workflow nodes using dagre.
 * Each call cycles through a different layout strategy (direction, ranker,
 * alignment, spacing) so successive clicks produce varied arrangements.
 *
 * Excludes child nodes (annotations/notes) from the layout graph —
 * they keep their relative position to the parent.
 * Annotation edges (sourceHandle === '__note_out') are excluded.
 *
 * Node dimensions in dagre use generous padding so edges route
 * around nodes rather than through them.
 *
 * @param {WorkflowNode[]} nodes - All workflow nodes
 * @param {WorkflowEdge[]} edges - All workflow edges
 * @param {AutoLayoutOptions} options - Layout options
 * @returns {AutoLayoutResult} Nodes with updated positions and changed flag
 */
export function autoLayoutNodes(
  nodes: WorkflowNode[],
  edges: WorkflowEdge[],
  options: AutoLayoutOptions = {},
): AutoLayoutResult {
  const {
    nodeWidth = DEFAULT_NODE_WIDTH,
    nodeHeight = DEFAULT_NODE_HEIGHT,
    snapGrid,
  } = options;

  // Separate layout-eligible nodes from child nodes (e.g., text notes)
  const layoutNodes = nodes.filter(n => !n.parentNodeId);
  const childNodes = nodes.filter(n => n.parentNodeId);

  if (layoutNodes.length <= 1) {
    return { nodes: [...nodes], changed: false };
  }

  // Pick current strategy and advance index for next call
  const strategy = LAYOUT_STRATEGIES[currentStrategyIndex % LAYOUT_STRATEGIES.length]!;
  currentStrategyIndex = (currentStrategyIndex + 1) % LAYOUT_STRATEGIES.length;

  // Count outgoing edges per node for dynamic width
  const outEdgeCount = new Map<string, number>();
  for (const edge of edges) {
    if (edge.sourceHandle === '__note_out') continue;
    outEdgeCount.set(edge.source, (outEdgeCount.get(edge.source) || 0) + 1);
  }

  // Build dagre graph
  const g = new dagre.graphlib.Graph();
  g.setDefaultEdgeLabel(() => ({}));
  g.setGraph({
    rankdir: strategy.rankdir,
    nodesep: strategy.nodesep,
    ranksep: strategy.ranksep,
    marginx: 60,
    marginy: 60,
    ranker: strategy.ranker,
    ...(strategy.align ? { align: strategy.align } : {}),
  });

  // Add nodes with dynamic width/height based on output count.
  // Use padded dimensions so dagre reserves more space around each node,
  // preventing edges from visually passing through node boxes.
  const PADDING = 40;
  const nodeWidthMap = new Map<string, number>();
  for (const node of layoutNodes) {
    const outputs = outEdgeCount.get(node.id) || 0;
    const baseW = outputs > 1
      ? Math.max(nodeWidth, outputs * WIDTH_PER_OUTPUT)
      : nodeWidth;
    nodeWidthMap.set(node.id, baseW);
    g.setNode(node.id, {
      width: baseW + PADDING,
      height: nodeHeight + PADDING,
    });
  }

  // Add flow edges (exclude annotation edges)
  // Higher weight = dagre tries harder to keep this edge short and straight
  const layoutNodeIds = new Set(layoutNodes.map(n => n.id));
  for (const edge of edges) {
    if (edge.sourceHandle === '__note_out') continue;
    if (!layoutNodeIds.has(edge.source) || !layoutNodeIds.has(edge.target)) continue;
    g.setEdge(edge.source, edge.target, { weight: 2, minlen: 1 });
  }

  // Run layout
  dagre.layout(g);

  // Map results back to WorkflowNode positions
  let changed = false;
  const updatedNodes: WorkflowNode[] = [];

  for (const node of layoutNodes) {
    const dagreNode = g.node(node.id);
    if (!dagreNode) {
      updatedNodes.push(node);
      continue;
    }

    // Convert dagre center-based coords to top-left using per-node width
    const w = nodeWidthMap.get(node.id) || nodeWidth;
    let x = dagreNode.x - w / 2;
    let y = dagreNode.y - nodeHeight / 2;

    // Snap to grid
    if (snapGrid) {
      x = Math.round(x / snapGrid[0]) * snapGrid[0];
      y = Math.round(y / snapGrid[1]) * snapGrid[1];
    }

    if (x !== node.position.x || y !== node.position.y) {
      changed = true;
    }

    updatedNodes.push({ ...node, position: { x, y } });
  }

  // Build position delta map for parents that moved
  const deltaMap = new Map<string, { dx: number; dy: number }>();
  for (const updated of updatedNodes) {
    const original = layoutNodes.find(n => n.id === updated.id);
    if (original && (original.position.x !== updated.position.x || original.position.y !== updated.position.y)) {
      deltaMap.set(updated.id, {
        dx: updated.position.x - original.position.x,
        dy: updated.position.y - original.position.y,
      });
    }
  }

  // Append child nodes — text notes use absolute positions so shift by parent delta
  for (const child of childNodes) {
    const delta = child.parentNodeId ? deltaMap.get(child.parentNodeId) : undefined;
    if (delta && child.type === 'core/text_note') {
      updatedNodes.push({ ...child, position: { x: child.position.x + delta.dx, y: child.position.y + delta.dy } });
    } else {
      updatedNodes.push(child);
    }
  }

  return { nodes: updatedNodes, changed };
}
