import type { OrganizationCoverageItem, OrganizationTreeNode } from '@stores/organization/types';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('treeBuilder');

/**
 * Calculate depth from pathKey
 * Example: "000001" = 1, "000001/000002" = 2, "000001/000002/000003" = 3
 *
 * @param pathKey - The hierarchical path key
 * @returns Depth level (1-based)
 */
export function getDepthFromPathKey(pathKey: string): number {
  return pathKey.split('/').length;
}

/**
 * Get parent pathKey
 * Example: "000001/000002/000003" -> "000001/000002"
 *
 * @param pathKey - The hierarchical path key
 * @returns Parent pathKey or null if root level
 */
export function getParentPathKey(pathKey: string): string | null {
  const segments = pathKey.split('/');
  if (segments.length <= 1) return null;
  return segments.slice(0, -1).join('/');
}

/**
 * Check if one pathKey is ancestor of another
 * Example: isAncestor("000001", "000001/000002/000003") -> true
 *
 * @param ancestorPath - Potential ancestor pathKey
 * @param descendantPath - Potential descendant pathKey
 * @returns true if ancestorPath is ancestor of descendantPath
 */
export function isAncestor(ancestorPath: string, descendantPath: string): boolean {
  if (ancestorPath === descendantPath) return false;
  return descendantPath.startsWith(ancestorPath + '/');
}

/**
 * Build hierarchical tree from flat list of organizations
 * Uses pathKey to determine parent-child relationships
 *
 * @param organizations - Flat list of organization coverage items
 * @returns Array of root-level tree nodes with nested children
 */
export function buildOrganizationTree(
  organizations: OrganizationCoverageItem[]
): OrganizationTreeNode[] {
  // Create a map for quick lookup by pathKey
  const orgMap = new Map<string, OrganizationTreeNode>();

  // Convert all organizations to tree nodes
  organizations.forEach(org => {
    orgMap.set(org.pathKey, {
      ...org,
      depth: getDepthFromPathKey(org.pathKey),
      enabled: true, // Assuming all are enabled (can be enhanced with actual status)
      children: [],
    });
  });

  // Build parent-child relationships
  const rootNodes: OrganizationTreeNode[] = [];

  orgMap.forEach((node) => {
    const parentPathKey = getParentPathKey(node.pathKey);

    if (parentPathKey) {
      // Find parent and add this node as child
      const parent = orgMap.get(parentPathKey);
      if (parent) {
        if (!parent.children) parent.children = [];
        parent.children.push(node);
      } else {
        // Parent not in coverage but node exists - treat as orphan root
        logger.warn(`Parent ${parentPathKey} not found for ${node.pathKey}`);
        rootNodes.push(node);
      }
    } else {
      // No parent means this is a root node (vendor level)
      rootNodes.push(node);
    }
  });

  // Sort children by pathKey at each level (ensures consistent ordering)
  function sortChildren(nodes: OrganizationTreeNode[]) {
    nodes.sort((a, b) => a.pathKey.localeCompare(b.pathKey));
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        sortChildren(node.children);
      }
    });
  }

  sortChildren(rootNodes);

  return rootNodes;
}

/**
 * Find organization in tree by ID
 * Recursively searches the tree structure
 *
 * @param nodes - Tree nodes to search
 * @param id - Organization ID to find
 * @returns Found organization node or null
 */
export function findOrganizationInTree(
  nodes: OrganizationTreeNode[],
  id: string
): OrganizationTreeNode | null {
  for (const node of nodes) {
    if (node.id === id) return node;
    if (node.children && node.children.length > 0) {
      const found = findOrganizationInTree(node.children, id);
      if (found) return found;
    }
  }
  return null;
}

/**
 * Get all ancestors of an organization
 * Returns path from root to the organization (excluding the org itself)
 *
 * @param nodes - Tree nodes to search
 * @param id - Organization ID
 * @returns Array of ancestor nodes from root to parent
 */
export function getAncestors(
  nodes: OrganizationTreeNode[],
  id: string
): OrganizationTreeNode[] {
  const ancestors: OrganizationTreeNode[] = [];

  function traverse(currentNodes: OrganizationTreeNode[], path: OrganizationTreeNode[]): boolean {
    for (const node of currentNodes) {
      if (node.id === id) {
        ancestors.push(...path);
        return true;
      }
      if (node.children && node.children.length > 0) {
        if (traverse(node.children, [...path, node])) {
          return true;
        }
      }
    }
    return false;
  }

  traverse(nodes, []);
  return ancestors;
}

/**
 * Flatten tree back to list
 * Useful for searching or filtering
 *
 * @param nodes - Tree nodes
 * @returns Flat array of all nodes
 */
export function flattenTree(nodes: OrganizationTreeNode[]): OrganizationTreeNode[] {
  const result: OrganizationTreeNode[] = [];

  function traverse(currentNodes: OrganizationTreeNode[]) {
    currentNodes.forEach(node => {
      result.push(node);
      if (node.children && node.children.length > 0) {
        traverse(node.children);
      }
    });
  }

  traverse(nodes);
  return result;
}
