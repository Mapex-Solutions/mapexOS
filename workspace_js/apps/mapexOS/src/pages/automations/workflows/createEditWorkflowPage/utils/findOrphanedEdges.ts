/** TYPE IMPORTS */
import type { HandleDefinition, WorkflowEdge } from '../interfaces/CreateEditWorkflow.interface';

/**
 * Find edges that reference handles which no longer exist on a node.
 * Returns the IDs of orphaned edges that should be removed.
 *
 * @param {string} nodeId - The node ID to check
 * @param {HandleDefinition[]} resolvedInputs - Current resolved input handles
 * @param {HandleDefinition[]} resolvedOutputs - Current resolved output handles
 * @param {WorkflowEdge[]} edges - All edges in the workflow
 * @returns {string[]} IDs of orphaned edges
 */
export function findOrphanedEdges(
  nodeId: string,
  resolvedInputs: HandleDefinition[],
  resolvedOutputs: HandleDefinition[],
  edges: WorkflowEdge[],
): string[] {
  const inputIds = new Set(resolvedInputs.map(h => h.id));
  const outputIds = new Set(resolvedOutputs.map(h => h.id));

  const orphaned: string[] = [];

  for (const edge of edges) {
    if (edge.source === nodeId && edge.sourceHandle && !outputIds.has(edge.sourceHandle)) {
      orphaned.push(edge.id);
    }
    if (edge.target === nodeId && edge.targetHandle && !inputIds.has(edge.targetHandle)) {
      orphaned.push(edge.id);
    }
  }

  return orphaned;
}
